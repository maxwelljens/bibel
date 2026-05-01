package bible

import (
	"fmt"
	"time"
)

// DateProgression handles calculating Bible position based on date
type DateProgression struct {
	bible        *Bible
	versesPerDay int
	readingMode  ReadingMode
	// Pre-computed cumulative offset at the START of each day of year (1-indexed)
	dailyOffsets []int
}

// NewDateProgression creates a new date progression calculator
func NewDateProgression(bible *Bible) *DateProgression {
	dp := &DateProgression{bible: bible, versesPerDay: 12, readingMode: ReadingModeEvangelion}
	dp.precomputeDailyOffsets()
	return dp
}

// NewDateProgressionWithConfig creates a new date progression calculator with config
func NewDateProgressionWithConfig(bible *Bible, versesPerDay int) *DateProgression {
	if versesPerDay <= 0 {
		versesPerDay = 12
	}
	dp := &DateProgression{bible: bible, versesPerDay: versesPerDay, readingMode: ReadingModeEvangelion}
	dp.precomputeDailyOffsets()
	return dp
}

// NewDateProgressionWithReadingMode creates a new date progression calculator with reading mode
func NewDateProgressionWithReadingMode(bible *Bible, versesPerDay int, readingMode ReadingMode) *DateProgression {
	if versesPerDay <= 0 {
		versesPerDay = 12
	}
	dp := &DateProgression{bible: bible, versesPerDay: versesPerDay, readingMode: readingMode}
	dp.precomputeDailyOffsets()
	return dp
}

// precomputeDailyOffsets computes the cumulative verse offset for each day of the year
// using lookahead logic: when fewer than versesPerDay remain in a chapter, extend
// to read all remaining verses in that chapter. This eliminates short reading days.
func (dp *DateProgression) precomputeDailyOffsets() {
	totalVerses := dp.GetTotalVersesForMode()
	if totalVerses == 0 {
		dp.dailyOffsets = make([]int, 367)
		return
	}

	dp.dailyOffsets = make([]int, 367) // 1-indexed by day of year
	cumOffset := 0

	for day := 1; day <= 366; day++ {
		dp.dailyOffsets[day] = cumOffset % totalVerses
		cumOffset += dp.computeConsumed(dp.dailyOffsets[day])
	}
}

// computeConsumed calculates the number of verses actually read starting from a given
// cumulative offset. Uses lookahead: if fewer than versesPerDay remain in the chapter,
// extends to read all remaining verses in that chapter.
func (dp *DateProgression) computeConsumed(offset int) int {
	startBook, endBook := dp.getBookRange()

	for book := startBook; book <= endBook; book++ {
		maxChapter := dp.findMaxChapter(book)
		if maxChapter == 0 {
			continue
		}

		for chapter := 1; chapter <= maxChapter; chapter++ {
			versesInChapter := dp.bible.CountVersesInChapter(book, chapter)

			if offset < versesInChapter {
				firstVerse := offset + 1
				lastVerse := firstVerse + dp.versesPerDay - 1
				if lastVerse > versesInChapter {
					lastVerse = versesInChapter
				}
				// Lookahead: if fewer than versesPerDay remain, read to chapter end
				remaining := versesInChapter - lastVerse
				if remaining > 0 && remaining < dp.versesPerDay {
					lastVerse = versesInChapter
				}
				return lastVerse - firstVerse + 1
			}

			offset -= versesInChapter
		}
	}

	// Wrap around
	totalVerses := dp.GetTotalVersesForMode()
	if totalVerses > 0 {
		return dp.computeConsumed(offset % totalVerses)
	}
	return 0
}

// findMaxChapter finds the last chapter number in a book
func (dp *DateProgression) findMaxChapter(book int) int {
	for chapter := 1; ; chapter++ {
		if dp.bible.CountVersesInChapter(book, chapter) == 0 {
			return chapter - 1
		}
	}
}

// GetPositionForDate calculates the verse range for a given date
func (dp *DateProgression) GetPositionForDate(date time.Time) (*Bookmark, error) {
	// Calculate day of year (1-366)
	dayOfYear := date.YearDay()

	if dayOfYear < 1 || dayOfYear >= len(dp.dailyOffsets) {
		dayOfYear = 1
	}

	// Look up pre-computed cumulative offset for this day
	targetVerseOffset := dp.dailyOffsets[dayOfYear]

	return dp.findPositionForOffset(targetVerseOffset)
}

// findPositionForOffset finds the bookmark for a given cumulative verse offset.
// Applies lookahead: if fewer than versesPerDay remain in the chapter, extends
// to read to the end of the chapter.
func (dp *DateProgression) findPositionForOffset(offset int) (*Bookmark, error) {
	startBook, endBook := dp.getBookRange()

	for book := startBook; book <= endBook; book++ {
		maxChapter := dp.findMaxChapter(book)
		if maxChapter == 0 {
			continue
		}

		for chapter := 1; chapter <= maxChapter; chapter++ {
			versesInChapter := dp.bible.CountVersesInChapter(book, chapter)

			if offset < versesInChapter {
				firstVerse := offset + 1

				// Default: read versesPerDay verses, clamped to chapter end
				lastVerse := firstVerse + dp.versesPerDay - 1
				if lastVerse > versesInChapter {
					lastVerse = versesInChapter
				}

				// Lookahead: if fewer than versesPerDay remain, read to chapter end
				remaining := versesInChapter - lastVerse
				if remaining > 0 && remaining < dp.versesPerDay {
					lastVerse = versesInChapter
				}

				return &Bookmark{
					Book:        BookIndex(book),
					Chapter:     chapter,
					FirstVerse:  firstVerse,
					SecondVerse: lastVerse,
				}, nil
			}

			offset -= versesInChapter
		}
	}

	// Wrap around to beginning if offset exceeds total verses
	totalVerses := dp.GetTotalVersesForMode()
	if totalVerses == 0 {
		return nil, fmt.Errorf("no verses found for reading mode %s. The Bible data file may not contain books for this reading mode", dp.readingMode)
	}
	if offset >= 0 {
		adjustedOffset := offset % totalVerses
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

// GetTotalGospelVerses returns the total number of verses in all four Gospels.
//
// Deprecated: Use GetTotalVersesForMode instead.
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
