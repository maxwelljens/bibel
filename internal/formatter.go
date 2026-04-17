package bible

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Formatter handles formatting Bible text for display
type Formatter struct {
	bible        *Bible
	useColours   bool
	headerFormat string
	numbered     bool
	paragraphs   bool
}

// NewFormatter creates a new formatter with default settings
func NewFormatter(bible *Bible) *Formatter {
	return &Formatter{
		bible:        bible,
		useColours:   true,
		headerFormat: "{book} {chapter}\nv. {first_verse}-{second_verse}",
		numbered:     false,
		paragraphs:   false,
	}
}

// NewFormatterWithFullConfig creates a new formatter with full configuration
func NewFormatterWithFullConfig(bible *Bible, useColours bool, headerFormat string, numbered bool, paragraphs bool) *Formatter {
	return &Formatter{
		bible:        bible,
		useColours:   useColours,
		headerFormat: headerFormat,
		numbered:     numbered,
		paragraphs:   paragraphs,
	}
}

// NewFormatterWithConfig creates a new formatter with configuration
func NewFormatterWithConfig(bible *Bible, useColours bool, headerFormat string) *Formatter {
	return NewFormatterWithFullConfig(bible, useColours, headerFormat, false, false)
}

// FormatHeader formats the header for a bookmark using the configured template
func (f *Formatter) FormatHeader(bookmark *Bookmark) string {
	template := f.headerFormat
	if template == "" {
		// Fall back to default template
		template = "{book} {chapter}\nv. {first_verse}-{second_verse}"
	}

	// Get book name from Bible data
	bookName := f.getBookName(bookmark.Book)

	// Replace template variables
	result := template
	result = strings.ReplaceAll(result, "{book}", bookName)
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

// getBookName retrieves the book name from Bible data
func (f *Formatter) getBookName(book BookIndex) string {
	// Try to find any verse from this book to get the book name
	// Start from chapter 1, verse 1 and work forward
	for chapter := 1; chapter <= 150; chapter++ {
		for verse := 1; verse <= 200; verse++ {
			if v, ok := f.bible.GetVerse(int(book), chapter, verse); ok {
				return v.BookName
			}
		}
	}
	// If we can't find the book in the data, return a generic name
	return fmt.Sprintf("Book %d", book)
} // FormatHeaderWithTemplate is deprecated - use FormatHeader instead
func (f *Formatter) FormatHeaderWithTemplate(bookmark *Bookmark) string {
	return f.FormatHeader(bookmark)
}

// FormatSnippet formats the verse range text
func (f *Formatter) FormatSnippet(verses []*Verse) string {
	// Use the new options-aware method
	return f.FormatSnippetWithOptions(verses)
}

// FormatSnippetWithOptions formats the verse range text with configuration options
func (f *Formatter) FormatSnippetWithOptions(verses []*Verse) string {
	if len(verses) == 0 {
		return "Error: No text matched"
	}

	var sb strings.Builder
	for i, verse := range verses {
		// Check for pilcrow to handle paragraph breaks
		text := strings.TrimSpace(verse.Text)
		hasPilcrow := strings.HasPrefix(text, "¶ ") || strings.HasPrefix(text, "¶")
		
		// Remove pilcrow markers if present
		if strings.HasPrefix(text, "¶ ") {
			text = text[2:]
		} else if strings.HasPrefix(text, "¶") {
			text = text[1:]
		}
		
		// Handle paragraph breaks (blank line before verse with pilcrow)
		if f.paragraphs && hasPilcrow && i > 0 {
			// Add blank line before this verse (paragraph break)
			// If numbered mode, we'll add a newline anyway, so just add one more
			if f.numbered {
				sb.WriteString("\n\n")
			} else {
				// In non-numbered mode, need space then blank line
				if i > 0 {
					sb.WriteString(" \n\n")
				} else {
					sb.WriteString("\n\n")
				}
			}
		}
		
		// Handle numbered output
		if f.numbered {
			if i > 0 && !(f.paragraphs && hasPilcrow) {
				// Don't add newline if we already added one for paragraph
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("%d. %s", verse.VerseNum, text))
		} else {
			// Regular mode - join with spaces
			// Don't add space if we just added paragraph break
			if i > 0 && !(f.paragraphs && hasPilcrow) {
				sb.WriteString(" ")
			}
			sb.WriteString(text)
		}
	}

	return sb.String()
}

// FormatSnippetForTUI formats verse text for TUI with lipgloss styling
func (f *Formatter) FormatSnippetForTUI(verses []*Verse, style lipgloss.Style, numberStyle lipgloss.Style, contentWidth int) string {
	if len(verses) == 0 {
		return "Error: No text matched"
	}

	var resultLines []string
	
	for i, verse := range verses {
		// Process verse text
		text := strings.TrimSpace(verse.Text)
		hasPilcrow := strings.HasPrefix(text, "¶ ") || strings.HasPrefix(text, "¶")
		
		// Remove pilcrow markers if present
		if strings.HasPrefix(text, "¶ ") {
			text = text[2:]
		} else if strings.HasPrefix(text, "¶") {
			text = text[1:]
		}
		
		// Handle paragraph breaks (blank line before verse with pilcrow)
		if f.paragraphs && hasPilcrow && i > 0 {
			// Add blank line before this verse (paragraph break)
			resultLines = append(resultLines, "")
		}
		
		if f.numbered {
			// For numbered mode: "1. text"
			// Estimate number takes about 4 characters (for verse numbers up to 3 digits + ". ")
			const numberWidthEstimate = 4
			remainingWidth := max(1, contentWidth - numberWidthEstimate)
			
			numberPart := numberStyle.Render(fmt.Sprintf("%d. ", verse.VerseNum))
			textPart := style.Width(remainingWidth).Render(text)
			resultLines = append(resultLines, numberPart + textPart)
		} else {
			// For non-numbered mode, we'll handle text accumulation separately
			// Store plain text for now, will join and render later
			if len(resultLines) == 0 || resultLines[len(resultLines)-1] == "" {
				// Start new line
				resultLines = append(resultLines, text)
			} else {
				// Continue current line with space
				resultLines[len(resultLines)-1] = resultLines[len(resultLines)-1] + " " + text
			}
		}
	}
	
	// Apply styling to non-numbered lines
	if !f.numbered {
		for i, line := range resultLines {
			if line != "" {  // Keep empty lines (paragraph breaks) as is
				resultLines[i] = style.Width(contentWidth).Render(line)
			}
		}
	}
	
	return strings.Join(resultLines, "\n")
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
// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}