package bible

import (
	"fmt"
	"strings"
)

// Formatter handles formatting Bible text for display
type Formatter struct {
}

// NewFormatter creates a new formatter
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatHeader formats the header for a bookmark with ANSI colors
func (f *Formatter) FormatHeader(bookmark *Bookmark) string {
	// ANSI color codes
	const (
		reset  = "\033[0m"
		green  = "\033[32m"
		yellow = "\033[33m"
	)

	return fmt.Sprintf("%s%s %d%s\n%s w. %d-%d%s",
		green,
		bookmark.Book.String(),
		bookmark.Chapter,
		reset,
		yellow,
		bookmark.FirstVerse,
		bookmark.SecondVerse,
		reset)
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

