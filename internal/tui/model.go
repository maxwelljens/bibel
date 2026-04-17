package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bible "maxwelljensen/bibel/internal"
)

// Model represents the TUI application state
type Model struct {
	bibleData  *bible.Bible
	bookmark   *bible.Bookmark
	width      int
	height     int
	quitting   bool
	outputMode string
	styles     *Styles
	formatter  *bible.Formatter
	config     *bible.Config
}

// Styles contains the lipgloss styles for the TUI
type Styles struct {
	box         lipgloss.Style
	header      lipgloss.Style
	content     lipgloss.Style
	quitMessage lipgloss.Style
}

// NewModel creates a new TUI model with the given Bible data and bookmark
func NewModel(bibleData *bible.Bible, bookmark *bible.Bookmark, config *bible.Config) Model {
	m := Model{
		bibleData:  bibleData,
		bookmark:   bookmark,
		outputMode: config.OutputMode,
		styles:     createStyles(config),
		formatter:  bible.NewFormatterWithConfig(config.Formatter.UseColours, config.Formatter.HeaderFormat),
		config:     config,
	}
	return m
}

// createStyles creates the lipgloss styles based on configuration
func createStyles(config *bible.Config) *Styles {
	// Detect if terminal has dark background
	hasDarkBG := lipgloss.HasDarkBackground()

	// Helper function for light/dark colors
	lightDark := func(light, dark lipgloss.TerminalColor) lipgloss.TerminalColor {
		if hasDarkBG {
			return dark
		}
		return light
	}

	// Get colours from config, or use adaptive defaults
	var borderColour, headerColour, textColour, quitColour lipgloss.TerminalColor

	// Border colour
	if config.TUI.BorderColour != "" {
		borderColour = lipgloss.Color(config.TUI.BorderColour)
	} else {
		borderColour = lightDark(lipgloss.Color("#666666"), lipgloss.Color("#999999"))
	}

	// Header colour (matching formatter's green)
	if config.TUI.HeaderColour != "" {
		headerColour = lipgloss.Color(config.TUI.HeaderColour)
	} else {
		headerColour = lightDark(lipgloss.Color("#00AA00"), lipgloss.Color("#00FF00"))
	}

	// Text colour
	if config.TUI.TextColour != "" {
		textColour = lipgloss.Color(config.TUI.TextColour)
	} else {
		textColour = lightDark(lipgloss.Color("#000000"), lipgloss.Color("#FFFFFF"))
	}

	// Quit message colour
	if config.TUI.QuitColour != "" {
		quitColour = lipgloss.Color(config.TUI.QuitColour)
	} else {
		quitColour = lightDark(lipgloss.Color("#555555"), lipgloss.Color("#AAAAAA"))
	}

	// Create border style based on config
	var border lipgloss.Border
	switch config.TUI.BorderStyle {
	case "double":
		border = lipgloss.DoubleBorder()
	case "single":
		border = lipgloss.NormalBorder()
	case "hidden":
		border = lipgloss.HiddenBorder()
	case "rounded":
		fallthrough
	default:
		border = lipgloss.RoundedBorder()
	}

	return &Styles{
		box: lipgloss.NewStyle().
			Border(border).
			BorderForeground(borderColour).
			Padding(1, 2).
			Margin(1, 0),

		header: lipgloss.NewStyle().
			Foreground(headerColour).
			Bold(true).
			MarginBottom(1),

		content: lipgloss.NewStyle().
			Foreground(textColour),

		quitMessage: lipgloss.NewStyle().
			Foreground(quitColour).
			Italic(true).
			MarginTop(1),
	}
}

// Init initializes the model (required by bubbletea)
func (m Model) Init() tea.Cmd {
	// No initial command needed
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		// Update window dimensions
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	// Get formatted content using the formatter
	header := m.formatter.FormatHeader(m.bookmark)
	content := m.formatter.ExtractAndFormat(m.bibleData, m.bookmark)

	// In plain or formatted mode, just return the original formatted text
	if m.outputMode == "plain" || m.outputMode == "formatted" {
		var output strings.Builder
		output.WriteString(header)
		output.WriteString("\n")
		output.WriteString(content)
		return output.String()
	}

	// For TUI mode, we need to parse the header to extract the book/chapter and verse range
	// The header from formatter has ANSI codes, but we want to style with lipgloss
	// Let's extract just the text without ANSI codes
	headerText := stripANSI(header)

	// Build the box content
	var boxContent strings.Builder

	// Add header styled with lipgloss
	boxContent.WriteString(m.styles.header.Render(headerText))
	boxContent.WriteString("\n\n")

	// Add verse content with wrapping based on available width
	// Calculate available width for content
	// Box has: Padding(1, 2) = 2 left + 2 right = 4
	// Border = 1 left + 1 right = 2
	// Total horizontal frame = 6
	// Use window width with some margin
	availableWidth := max(
		// Leave some terminal margin
		m.width-10,
		// Minimum reasonable width
		30)
	contentWidth := availableWidth - 6 // Account for box frame

	// Apply width constraint for text wrapping
	wrappedContent := m.styles.content.Width(contentWidth).Render(content)
	boxContent.WriteString(wrappedContent)

	// Create the box with max width constraint
	box := m.styles.box.MaxWidth(availableWidth).Render(boxContent.String())

	// Add quit message below the box if configured to show it
	var fullOutput string
	if m.config.TUI.ShowQuitMessage {
		quitMsg := m.styles.quitMessage.Render("Press q to quit...")
		fullOutput = lipgloss.JoinVertical(lipgloss.Center, box, quitMsg)
	} else {
		fullOutput = lipgloss.JoinVertical(lipgloss.Center, box)
	}

	return fullOutput
}

// CreateTUIProgram creates and returns a bubbletea program for the Bible verse
func CreateTUIProgram(bible *bible.Bible, bookmark *bible.Bookmark, config *bible.Config) *tea.Program {
	model := NewModel(bible, bookmark, config)
	return tea.NewProgram(model)
}

// stripANSI removes ANSI escape codes from a string
func stripANSI(str string) string {
	var result strings.Builder
	inEscape := false

	for i := 0; i < len(str); i++ {
		if str[i] == '\033' && i+1 < len(str) && str[i+1] == '[' {
			inEscape = true
			i++ // Skip the '['
			continue
		}

		if inEscape {
			// Look for the end of the escape sequence (letter m)
			if str[i] >= 'A' && str[i] <= 'Z' || str[i] >= 'a' && str[i] <= 'z' {
				inEscape = false
			}
			continue
		}

		result.WriteByte(str[i])
	}

	return result.String()
}