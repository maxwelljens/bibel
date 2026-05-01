package main

import (
	"fmt"
	"os"
	"time"

	"maxwelljensen/bibel/internal"
)

// TestCase represents a test case for date progression
type TestCase struct {
	name               string
	date               time.Time
	readingMode        bible.ReadingMode
	expectedBook       int
	expectedChapter    int
	expectedStartVerse int
	expectedEndVerse   int
}

func main() {
	// Load Bible data
	bibleData, err := bible.LoadBible("books/pol_nbg.json", false)
	if err != nil {
		fmt.Printf("FAILED: Error loading Bible data: %v\n", err)
		os.Exit(1)
	}

	// Define test cases - using known values from the fixed system
	testCases := []TestCase{
		// Day 1 (1 January) should start at Matthew 1:1-12 for evangelion mode
		{
			name:               "Day 1 - evangelion",
			date:               time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			readingMode:        bible.ReadingModeEvangelion,
			expectedBook:       40,
			expectedChapter:    1,
			expectedStartVerse: 1,
			expectedEndVerse:   12,
		},
		// Day 110 (20 April, 2026) - Mark 12:1-12
		// With cumulative lookahead, all days read 12+ verses and no short days occur.
		{
			name:               "Day 110 - evangelion",
			date:               time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC),
			readingMode:        bible.ReadingModeEvangelion,
			expectedBook:       41, // Mark
			expectedChapter:    12,
			expectedStartVerse: 1,
			expectedEndVerse:   12,
		},
		// Test New Testament mode - should also start at Matthew 1:1-12
		{
			name:               "Day 1 - new_testament",
			date:               time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			readingMode:        bible.ReadingModeNewTestament,
			expectedBook:       40,
			expectedChapter:    1,
			expectedStartVerse: 1,
			expectedEndVerse:   12,
		},
		// Test Old Testament mode - should start at Genesis 1:1-12
		{
			name:               "Day 1 - old_testament",
			date:               time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			readingMode:        bible.ReadingModeOldTestament,
			expectedBook:       1,
			expectedChapter:    1,
			expectedStartVerse: 1,
			expectedEndVerse:   12,
		},
		// Test wrap-around: Day that would exceed total Gospel verses
		// 3779 verses / 12 verses/day = ~315 days
		// With cumulative lookahead, extra verses consumed on chapter-boundary
		// days push the schedule further, so day 316 lands in Matthew 18.
		{
			name:               "Day 316 - evangelion (wrap test)",
			date:               time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC), // Day 316
			readingMode:        bible.ReadingModeEvangelion,
			expectedBook:       40,
			expectedChapter:    18,
			expectedStartVerse: 1,
			expectedEndVerse:   12,
		},
	}

	// Run tests
	passed := 0
	failed := 0

	fmt.Println("=== Running Date Progression Tests ===")
	fmt.Printf("Using Bible: %s\n", bibleData.Metadata.Name)
	fmt.Printf("Total verses in file: %d\n", len(bibleData.Verses))

	for _, tc := range testCases {
		fmt.Printf("\nTest: %s\n", tc.name)
		fmt.Printf("  Date: %s (Day of year: %d)\n", tc.date.Format("2006-01-02"), tc.date.YearDay())

		dateProg := bible.NewDateProgressionWithReadingMode(bibleData, 12, tc.readingMode)
		bookmark, err := dateProg.GetPositionForDate(tc.date)

		if err != nil {
			fmt.Printf("  FAILED: Error: %v\n", err)
			failed++
			continue
		}

		fmt.Printf("  Result: Book %d, Chapter %d, Verses %d-%d\n",
			bookmark.Book, bookmark.Chapter, bookmark.FirstVerse, bookmark.SecondVerse)

		if int(bookmark.Book) == tc.expectedBook &&
			bookmark.Chapter == tc.expectedChapter &&
			bookmark.FirstVerse == tc.expectedStartVerse &&
			bookmark.SecondVerse == tc.expectedEndVerse {
			fmt.Printf("  ✓ PASSED\n")
			passed++
		} else {
			fmt.Printf("  ✗ FAILED: Expected Book %d, Chapter %d, Verses %d-%d\n",
				tc.expectedBook, tc.expectedChapter, tc.expectedStartVerse, tc.expectedEndVerse)
			failed++
		}
	}

	fmt.Printf("\n=== Test Summary ===\n")
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)

	if failed > 0 {
		fmt.Printf("❌ Some tests FAILED\n")
		os.Exit(1)
	} else {
		fmt.Printf("✅ All tests PASSED\n")
	}
}
