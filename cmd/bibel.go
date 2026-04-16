package main

import (
	"fmt"
	"os"

	"maxwelljensen/bibel/internal"
)

func main() {
	// TEST: Initialise components
	biblePath := "books/pol_nbg.json"
	bookmarkPath := "test/bookmark.toml"

	// Load Bible data
	bibleData, err := bible.LoadBible(biblePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading Bible data: %v\n", err)
		os.Exit(1)
	}

	// Initialise bookmark manager
	bookmarkMgr := bible.NewBookmarkManager(bookmarkPath)

	// Read current bookmark
	currentBookmark, err := bookmarkMgr.ReadBookmark()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading bookmark: %v\n", err)
		os.Exit(1)
	}

	// Adjust bookmark based on lookahead rule
	adjustedBookmark := bookmarkMgr.AdjustForLookahead(currentBookmark, bibleData)

	// Initialise formatter
	formatter := bible.NewFormatter()

	// Print header and snippet
	fmt.Println(formatter.FormatHeader(adjustedBookmark))
	fmt.Println(formatter.ExtractAndFormat(bibleData, adjustedBookmark))

	// Calculate next bookmark
	nextBookmark := bookmarkMgr.AdvanceBookmark(adjustedBookmark, bibleData)

	// Write next bookmark
	if err := bookmarkMgr.WriteBookmark(nextBookmark); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing next bookmark: %v\n", err)
		os.Exit(1)
	}
}

