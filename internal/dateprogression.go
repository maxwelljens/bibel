package bible

import (
	"fmt"
	"time"
)

// DateProgression handles calculating Bible position based on date
type DateProgression struct {
	bible       *Bible
	versesPerDay int
}

// NewDateProgression creates a new date progression calculator
func NewDateProgression(bible *Bible) *DateProgression {
	return &DateProgression{bible: bible, versesPerDay: 12}
}

// NewDateProgressionWithConfig creates a new date progression calculator with config
func NewDateProgressionWithConfig(bible *Bible, versesPerDay int) *DateProgression {
	if versesPerDay <= 0 {
		versesPerDay = 12
	}
	return &DateProgression{bible: bible, versesPerDay: versesPerDay}
}

// GetPositionForDate calculates the verse range for a given date
func (dp *DateProgression) GetPositionForDate(date time.Time) (*Bookmark, error) {
	// Calculate day of year (1-366)
	dayOfYear := date.YearDay()

	// Each day shows configured number of verses
	targetVerseOffset := (dayOfYear - 1) * dp.versesPerDay // Zero-indexed

	// Walk through Gospels to find the position
	return dp.findPositionForOffset(targetVerseOffset)
}

// findPositionForOffset finds the bookmark for a given cumulative verse offset
func (dp *DateProgression) findPositionForOffset(offset int) (*Bookmark, error) {
	// Iterate through books 40-43 (Matthew, Mark, Luke, John)
	for book := 40; book <= 43; book++ {
		// Find last chapter in this book
		maxChapter := 0
		for chapter := 1; ; chapter++ {
			if dp.bible.CountVersesInChapter(book, chapter) == 0 {
				maxChapter = chapter - 1
				break
			}
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

	// If we've gone through all Gospels and offset is still positive, wrap
	// around to beginning (start over) Calculate modulo offset within total
	// Gospel verses
	totalVerses := dp.GetTotalGospelVerses()
	if offset >= 0 {
		adjustedOffset := offset % totalVerses
		// Recursively find position for adjusted offset
		return dp.findPositionForOffset(adjustedOffset)
	}

	return nil, fmt.Errorf("could not find position for offset %d", offset)
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

// GetTotalGospelVerses returns the total number of verses in all four Gospels
func (dp *DateProgression) GetTotalGospelVerses() int {
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