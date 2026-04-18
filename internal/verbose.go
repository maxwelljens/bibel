package bible

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// VerboseLogger handles verbose output with lipgloss styling
type VerboseLogger struct {
	enabled bool
	styles  *VerboseStyles
}

// VerboseStyles defines styles for verbose output
type VerboseStyles struct {
	Info    lipgloss.Style
	Warning lipgloss.Style
	Success lipgloss.Style
	Path    lipgloss.Style
	Label   lipgloss.Style
	Value   lipgloss.Style
}

// NewVerboseLogger creates a new verbose logger
func NewVerboseLogger(enabled bool) *VerboseLogger {
	return &VerboseLogger{
		enabled: enabled,
		styles:  createVerboseStyles(),
	}
}

// createVerboseStyles creates and returns the verbose output styles
func createVerboseStyles() *VerboseStyles {
	return &VerboseStyles{
		Info:    lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
		Warning: lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		Success: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		Path:    lipgloss.NewStyle().Foreground(lipgloss.Color("12")),
		Label:   lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Value:   lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
	}
}

// Info prints an informational message
func (vl *VerboseLogger) Info(format string, args ...any) {
	if !vl.enabled {
		return
	}
	message := fmt.Sprintf(format, args...)
	styled := vl.styles.Info.Render("INFO: ") + message
	fmt.Fprintln(os.Stderr, styled)
}

// Warning prints a warning message
func (vl *VerboseLogger) Warning(format string, args ...any) {
	if !vl.enabled {
		return
	}
	message := fmt.Sprintf(format, args...)
	styled := vl.styles.Warning.Render("WARN: ") + message
	fmt.Fprintln(os.Stderr, styled)
}

// Success prints a success message
func (vl *VerboseLogger) Success(format string, args ...any) {
	if !vl.enabled {
		return
	}
	message := fmt.Sprintf(format, args...)
	styled := vl.styles.Success.Render("SUCCESS: ") + message
	fmt.Fprintln(os.Stderr, styled)
}

// Path prints a file path with styling
func (vl *VerboseLogger) Path(path string) {
	if !vl.enabled {
		return
	}
	// Make path absolute for clarity
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	styled := vl.styles.Label.Render("Path: ") + vl.styles.Path.Render(absPath)
	fmt.Fprintln(os.Stderr, styled)
}

// Config prints configuration information
func (vl *VerboseLogger) Config(format string, args ...any) {
	if !vl.enabled {
		return
	}
	message := fmt.Sprintf(format, args...)
	styled := vl.styles.Label.Render("Config: ") + message
	fmt.Fprintln(os.Stderr, styled)
}

// Value prints a key-value pair
func (vl *VerboseLogger) Value(key string, value any) {
	if !vl.enabled {
		return
	}
	styled := vl.styles.Label.Render(key+": ") + vl.styles.Value.Render(fmt.Sprintf("%v", value))
	fmt.Fprintln(os.Stderr, styled)
}

// Section prints a section header
func (vl *VerboseLogger) Section(title string) {
	if !vl.enabled {
		return
	}
	styled := vl.styles.Info.Render("=== " + title + " ===")
	fmt.Fprintln(os.Stderr, styled)
}

// WithFields creates a formatted message with key-value pairs
func (vl *VerboseLogger) WithFields(fields map[string]any, message string) {
	if !vl.enabled {
		return
	}

	var parts []string
	for key, value := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}

	fieldsStr := strings.Join(parts, " ")
	styled := vl.styles.Label.Render("Fields: ") + message
	if fieldsStr != "" {
		styled += " " + vl.styles.Info.Render("["+fieldsStr+"]")
	}
	fmt.Fprintln(os.Stderr, styled)
}
