package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bible "maxwelljensen/bibel/internal"
)

// Model represents the TUI application state
type Model struct {
	bibleData *bible.Bible
	bookmark  *bible.Bookmark
	width     int
	height    int
	quitting  bool
	plainMode bool
	styles    *Styles
	formatter *bible.Formatter
}

// Styles contains the lipgloss styles for the TUI
type Styles struct {
	box         lipgloss.Style
	header      lipgloss.Style
	content     lipgloss.Style
	quitMessage lipgloss.Style
}

// NewModel creates a new TUI model with the given Bible data and bookmark
func NewModel(bibleData *bible.Bible, bookmark *bible.Bookmark, plainMode bool) Model {
	m := Model{
		bibleData: bibleData,
		bookmark:  bookmark,
		plainMode: plainMode,
		styles:    createStyles(),
		formatter: &bible.Formatter{},
	}
	return m
}

// createStyles creates the lipgloss styles, adapting to terminal background color
func createStyles() *Styles {
	// Detect if terminal has dark background
	hasDarkBG := lipgloss.HasDarkBackground()

	// Helper function for light/dark colors
	lightDark := func(light, dark lipgloss.TerminalColor) lipgloss.TerminalColor {
		if hasDarkBG {
			return dark
		}
		return light
	}

	// Colors that match the existing formatter colors (green/yellow)
	// but adapted for light/dark backgrounds
	green := lightDark(lipgloss.Color("#00AA00"), lipgloss.Color("#00FF00"))
	borderColor := lightDark(lipgloss.Color("#666666"), lipgloss.Color("#999999"))
	textColor := lightDark(lipgloss.Color("#000000"), lipgloss.Color("#FFFFFF"))
	quitColor := lightDark(lipgloss.Color("#555555"), lipgloss.Color("#AAAAAA"))

	return &Styles{
		box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Margin(1, 0),

		header: lipgloss.NewStyle().
			Foreground(green).
			Bold(true).
			MarginBottom(1),

		content: lipgloss.NewStyle().
			Foreground(textColor),

		quitMessage: lipgloss.NewStyle().
			Foreground(quitColor).
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

	// In plain mode, just return the original formatted text (with ANSI colors)
	if m.plainMode {
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
	availableWidth := m.width - 10 // Leave some terminal margin
	if availableWidth < 30 {
		availableWidth = 30 // Minimum reasonable width
	}
	contentWidth := availableWidth - 6 // Account for box frame

	// Apply width constraint for text wrapping
	wrappedContent := m.styles.content.Width(contentWidth).Render(content)
	boxContent.WriteString(wrappedContent)

	// Create the box with max width constraint
	box := m.styles.box.MaxWidth(availableWidth).Render(boxContent.String())

	// Add quit message below the box
	quitMsg := m.styles.quitMessage.Render("Press q to quit...")

	// Combine box and quit message
	return lipgloss.JoinVertical(lipgloss.Center, box, quitMsg)
}

// CreateTUIProgram creates and returns a bubbletea program for the Bible verse
func CreateTUIProgram(bible *bible.Bible, bookmark *bible.Bookmark, plainMode bool) *tea.Program {
	model := NewModel(bible, bookmark, plainMode)
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

