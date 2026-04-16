package main

import (
	"fmt"
	"os"
	"time"

	"maxwelljensen/bibel/internal"
)

func main() {
	// Load Bible data
	biblePath := "books/pol_nbg.json"
	bibleData, err := bible.LoadBible(biblePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading Bible data: %v\n", err)
		os.Exit(1)
	}

	// Initialize date progression
	dateProg := bible.NewDateProgression(bibleData)

	// Get current date
	currentDate := time.Now()
	
	// Calculate position for today
	todayBookmark, err := dateProg.GetPositionForDate(currentDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calculating date position: %v\n", err)
		os.Exit(1)
	}

	// Initialize formatter
	formatter := bible.NewFormatter()

	// Print header and snippet
	fmt.Println(formatter.FormatHeader(todayBookmark))
	fmt.Println(formatter.ExtractAndFormat(bibleData, todayBookmark))
}