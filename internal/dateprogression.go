package bible

import (
	"fmt"
	"time"
)

// DateProgression handles calculating Bible position based on date
type DateProgression struct {
	bible       *Bible
	versesPerDay int
	readingMode  ReadingMode
}

// NewDateProgression creates a new date progression calculator
func NewDateProgression(bible *Bible) *DateProgression {
	return &DateProgression{bible: bible, versesPerDay: 12, readingMode: ReadingModeEvangelion}
}

// NewDateProgressionWithConfig creates a new date progression calculator with config
func NewDateProgressionWithConfig(bible *Bible, versesPerDay int) *DateProgression {
	if versesPerDay <= 0 {
		versesPerDay = 12
	}
	return &DateProgression{bible: bible, versesPerDay: versesPerDay, readingMode: ReadingModeEvangelion}
}

// NewDateProgressionWithReadingMode creates a new date progression calculator with reading mode
func NewDateProgressionWithReadingMode(bible *Bible, versesPerDay int, readingMode ReadingMode) *DateProgression {
	if versesPerDay <= 0 {
		versesPerDay = 12
	}
	return &DateProgression{bible: bible, versesPerDay: versesPerDay, readingMode: readingMode}
}

// GetPositionForDate calculates the verse range for a given date
func (dp *DateProgression) GetPositionForDate(date time.Time) (*Bookmark, error) {
	// Calculate day of year (1-366)
	dayOfYear := date.YearDay()

	// Each day shows configured number of verses
	targetVerseOffset := (dayOfYear - 1) * dp.versesPerDay // Zero-indexed

	// Walk through appropriate books based on reading mode
	return dp.findPositionForOffset(targetVerseOffset)
}

// findPositionForOffset finds the bookmark for a given cumulative verse offset
func (dp *DateProgression) findPositionForOffset(offset int) (*Bookmark, error) {
	// Determine which books to iterate through based on reading mode
	startBook, endBook := dp.getBookRange()
	
	// Iterate through books in the range
	for book := startBook; book <= endBook; book++ {
		// Find last chapter in this book
		maxChapter := 0
		for chapter := 1; ; chapter++ {
			if dp.bible.CountVersesInChapter(book, chapter) == 0 {
				maxChapter = chapter - 1
				break
			}
		}
		if maxChapter == 0 {
			continue // Book has no chapters? Shouldn't happen
		}

		// Iterate through chapters
		for chapter := 1; chapter <= maxChapter; chapter++ {
			versesInChapter := dp.bible.CountVersesInChapter(book, chapter)

			// Check if offset falls within this chapter
			if offset < versesInChapter {
				// Offset is within this chapter
				firstVerse := offset + 1 // Convert from 0-indexed to 1-indexed

				// Create initial bookmark (like AdvanceBookmark does)
				secondVerse := min(firstVerse+dp.versesPerDay-1, versesInChapter)
				bookmark := &Bookmark{
					Book:        BookIndex(book),
					Chapter:     chapter,
					FirstVerse:  firstVerse,
					SecondVerse: secondVerse,
				}

				// Apply lookahead rule (like AdjustForLookahead)
				return dp.applyLookahead(bookmark), nil
			}

			// Move to next chapter, subtracting this chapter's verses
			offset -= versesInChapter
		}
	}

	// If we've gone through all books and offset is still positive, wrap
	// around to beginning (start over) Calculate modulo offset within total verses
	totalVerses := dp.GetTotalVersesForMode()
	if totalVerses == 0 {
		return nil, fmt.Errorf("no verses found for reading mode %s. The Bible data file may not contain books for this reading mode", dp.readingMode)
	}
	if offset >= 0 {
		adjustedOffset := offset % totalVerses
		// Recursively find position for adjusted offset
		return dp.findPositionForOffset(adjustedOffset)
	}

	return nil, fmt.Errorf("could not find position for offset %d", offset)
}

// getBookRange returns the start and end book numbers for the current reading mode
func (dp *DateProgression) getBookRange() (startBook, endBook int) {
	switch dp.readingMode {
	case ReadingModeEvangelion:
		return 40, 43 // Matthew, Mark, Luke, John
	case ReadingModeNewTestament:
		return 40, 66 // Matthew through Revelation
	case ReadingModeOldTestament:
		return 1, 39 // Genesis through Malachi
	case ReadingModeBible:
		return 1, 66 // Entire Bible
	default:
		return 40, 43 // Default to Evangelion
	}
}

// applyLookahead applies the lookahead rule (same as BookmarkManager.AdjustForLookahead)
func (dp *DateProgression) applyLookahead(bookmark *Bookmark) *Bookmark {
	versesInChapter := dp.bible.CountVersesInChapter(int(bookmark.Book), bookmark.Chapter)
	versesRemaining := versesInChapter - bookmark.SecondVerse

	// If we're not at the end of the chapter and less than dp.versesPerDay verses remain for NEXT snippet
	if versesRemaining > 0 && versesRemaining < dp.versesPerDay {
		return &Bookmark{
			Book:        bookmark.Book,
			Chapter:     bookmark.Chapter,
			FirstVerse:  bookmark.FirstVerse,
			SecondVerse: versesInChapter,
		}
	}

	// Otherwise return the original bookmark
	return bookmark
}

// GetTotalVersesForMode returns the total number of verses for the current reading mode
func (dp *DateProgression) GetTotalVersesForMode() int {
	total := 0
	startBook, endBook := dp.getBookRange()
	
	for book := startBook; book <= endBook; book++ {
		for chapter := 1; ; chapter++ {
			verses := dp.bible.CountVersesInChapter(book, chapter)
			if verses == 0 {
				break
			}
			total += verses
		}
	}
	return total
}

// GetTotalGospelVerses returns the total number of verses in all four Gospels
// Deprecated: Use GetTotalVersesForMode instead
func (dp *DateProgression) GetTotalGospelVerses() int {
	// Calculate Gospels total (books 40-43)
	total := 0
	for book := 40; book <= 43; book++ {
		for chapter := 1; ; chapter++ {
			verses := dp.bible.CountVersesInChapter(book, chapter)
			if verses == 0 {
				break
			}
			total += verses
		}
	}
	return total
}