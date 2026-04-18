package tui

import (
	"fmt"
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
	easterProg *bible.EasterProgression
}

// Styles contains the lipgloss styles for the TUI
type Styles struct {
	box         lipgloss.Style
	header      lipgloss.Style
	content     lipgloss.Style
	numberStyle lipgloss.Style // For numbered verse output
	quitMessage lipgloss.Style
}

// NewModel creates a new TUI model with the given Bible data and bookmark
func NewModel(bibleData *bible.Bible, bookmark *bible.Bookmark, config *bible.Config) Model {
	// Initialise Easter progression
	easterProg := bible.NewEasterProgression(config.EasterType)

	m := Model{
		bibleData:  bibleData,
		bookmark:   bookmark,
		outputMode: config.OutputMode,
		styles:     createStyles(config),
		formatter:  bible.NewFormatterWithFullConfig(bibleData, config.Formatter.UseColours, config.Formatter.HeaderFormat, config.Formatter.Numbered, config.Formatter.Paragraphs),
		config:     config,
		easterProg: easterProg,
	}
	return m
}

// createStyles creates the lipgloss styles based on configuration
func createStyles(config *bible.Config) *Styles {
	// Use ANSI colours (set by terminal emulator)
	// Border colour - bright black (gray)
	borderColour := lipgloss.Color("8") // ANSI bright black (gray)
	
	// Header colour - green (matching formatter's green)
	headerColour := lipgloss.Color("2") // ANSI green
	
	// Text colour - default terminal text colour
	textColour := lipgloss.Color("") // Default terminal colour
	
	// Quit message colour - bright black (gray)
	quitColour := lipgloss.Color("8") // ANSI bright black (gray)

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

		numberStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")). // ANSI bright black (gray)
			Bold(false),

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

	// For TUI mode, use special formatting
	verses := m.bibleData.GetVerseRange(int(m.bookmark.Book), m.bookmark.Chapter,
		m.bookmark.FirstVerse, m.bookmark.SecondVerse)

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
	// Use window size with some margin
	availableWidth := max(
		// Leave some terminal margin
		m.width-10,
		// Minimum reasonable width
		30)
	contentWidth := availableWidth - 6 // Account for box frame

	// Format verses with appropriate styling
	var verseContent string
	if m.config.Formatter.Numbered || m.config.Formatter.Paragraphs {
		// Use TUI-aware formatting for numbered or paragraph mode
		verseContent = m.formatter.FormatSnippetForTUI(verses,
			m.styles.content,
			m.styles.numberStyle,
			contentWidth)
	} else {
		// Use regular formatting
		wrappedContent := m.styles.content.Width(contentWidth).Render(content)
		verseContent = wrappedContent
	}
	boxContent.WriteString(verseContent)

	// Create the box with max width constraint
	box := m.styles.box.MaxWidth(availableWidth).Render(boxContent.String())

	// Create Easter progress bar
	easterProgress := m.renderEasterProgressBar(availableWidth)

	// Add quit message below the box if configured to show it
	var fullOutput string
	if m.config.TUI.ShowQuitMessage {
		quitMsg := m.styles.quitMessage.Render("Press q to quit...")
		fullOutput = lipgloss.JoinVertical(lipgloss.Center, box, easterProgress, quitMsg)
	} else {
		fullOutput = lipgloss.JoinVertical(lipgloss.Center, box, easterProgress)
	}

	return fullOutput
}

// renderEasterProgressBar renders a fancy progress bar showing time until Easter
func (m Model) renderEasterProgressBar(width int) string {
	if m.easterProg == nil {
		return ""
	}

	// Calculate progress percentage
	progress, err := m.easterProg.GetEasterProgressPercentage()
	if err != nil {
		// If we can't calculate Easter, show error message
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")). // ANSI red
			Italic(true)
		return errorStyle.Width(width - 6).Render("Unable to calculate Easter date")
	}

	// Format time until Easter
	timeStr, err := m.easterProg.FormatEasterProgress()
	if err != nil {
		timeStr = "Calculating..."
	}

	// Create progress bar
	barWidth := max(
		// Leave some padding
		width-10,
		// Minimum bar width
		20)

	// Calculate filled width
	filledWidth := max(min(int(float64(barWidth)*progress), barWidth), 0)
	emptyWidth := barWidth - filledWidth

	// Define styles using ANSI colours
	filledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")). // Bright white
		Background(lipgloss.Color("2"))   // Green

	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).  // Bright black (gray)
		Background(lipgloss.Color("0"))   // Black

	// Create bar segments
	filledBar := filledStyle.Render(strings.Repeat("█", filledWidth))
	emptyBar := emptyStyle.Render(strings.Repeat("░", emptyWidth))
	bar := filledBar + emptyBar

	// Add percentage text
	percentage := fmt.Sprintf("%.1f%%", progress*100)
	percentageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")). // ANSI bright black (gray)
		Bold(true)

	percentageText := percentageStyle.Render(percentage)

	// Combine time text, bar, and percentage
	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")) // ANSI white (light gray)

	timeText := timeStyle.Width(barWidth - len(percentage) - 2).Render(timeStr)

	// Layout: time text on left, bar below, percentage on right of bar
	topRow := lipgloss.JoinHorizontal(lipgloss.Left, timeText, " ", percentageText)
	bottomRow := bar

	// Container for the progress bar
	containerStyle := lipgloss.NewStyle().
		Padding(0, 1).
		MarginTop(1)

	return containerStyle.Render(lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow))
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
		// Check for ESC [ sequence (ANSI escape)
		if str[i] == 27 && i+1 < len(str) && str[i+1] == '[' {
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
