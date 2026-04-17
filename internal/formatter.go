package bible

import (
	"fmt"
	"strings"
)

// Formatter handles formatting Bible text for display
type Formatter struct {
	useColours   bool
	headerFormat string
}

// NewFormatter creates a new formatter with default settings
func NewFormatter() *Formatter {
	return &Formatter{
		useColours:   true,
		headerFormat: "{book} {chapter}\nw. {first_verse}-{second_verse}",
	}
}

// NewFormatterWithConfig creates a new formatter with configuration
func NewFormatterWithConfig(useColours bool, headerFormat string) *Formatter {
	return &Formatter{
		useColours:   useColours,
		headerFormat: headerFormat,
	}
}

// FormatHeader formats the header for a bookmark using the configured template
func (f *Formatter) FormatHeader(bookmark *Bookmark) string {
	template := f.headerFormat
	if template == "" {
		// Fall back to default template
		template = "{book} {chapter}\nw. {first_verse}-{second_verse}"
	}
	
	// Replace template variables
	result := template
	result = strings.ReplaceAll(result, "{book}", bookmark.Book.String())
	result = strings.ReplaceAll(result, "{chapter}", fmt.Sprintf("%d", bookmark.Chapter))
	result = strings.ReplaceAll(result, "{first_verse}", fmt.Sprintf("%d", bookmark.FirstVerse))
	result = strings.ReplaceAll(result, "{second_verse}", fmt.Sprintf("%d", bookmark.SecondVerse))
	
	if !f.useColours {
		return result
	}
	
	// Add ANSI colours - simple implementation
	// For more advanced templating, we'd need a proper template engine
	const (
		reset  = "\033[0m"
		green  = "\033[32m"
		yellow = "\033[33m"
	)
	
	// Simple colouring - book and chapter in green, verse range in yellow
	lines := strings.Split(result, "\n")
	if len(lines) >= 1 {
		lines[0] = green + lines[0] + reset
	}
	if len(lines) >= 2 {
		lines[1] = yellow + lines[1] + reset
	}
	result = strings.Join(lines, "\n")
	
	return result
}

// FormatHeaderWithTemplate is deprecated - use FormatHeader instead
func (f *Formatter) FormatHeaderWithTemplate(bookmark *Bookmark) string {
	return f.FormatHeader(bookmark)
}

// FormatSnippet formats the verse range text
func (f *Formatter) FormatSnippet(verses []*Verse) string {
	if len(verses) == 0 {
		return "Error: No text matched"
	}

	var sb strings.Builder
	for i, verse := range verses {
		if i > 0 {
			sb.WriteString(" ")
		}
		// Remove any leading pilcrow (¶) and trim
		text := strings.TrimSpace(verse.Text)
		if strings.HasPrefix(text, "¶ ") {
			text = text[2:]
		} else if strings.HasPrefix(text, "¶") {
			text = text[1:]
		}
		sb.WriteString(text)
	}

	return sb.String()
}

// ExtractAndFormat extracts verses for a bookmark and formats them
func (f *Formatter) ExtractAndFormat(bible *Bible, bookmark *Bookmark) string {
	verses := bible.GetVerseRange(int(bookmark.Book), bookmark.Chapter,
		bookmark.FirstVerse, bookmark.SecondVerse)
	return f.FormatSnippet(verses)
}

// ExtractAndFormatWithHeader extracts verses and formats with header
func (f *Formatter) ExtractAndFormatWithHeader(bible *Bible, bookmark *Bookmark) string {
	header := f.FormatHeader(bookmark)
	content := f.ExtractAndFormat(bible, bookmark)
	return header + "\n" + content
}