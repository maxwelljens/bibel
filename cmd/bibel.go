package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"maxwelljensen/bibel/internal"
)

// Args defines the command line arguments
type Args struct {
	Plain bool `arg:"-p,--plain" help:"output plain text without chapter name or verse numbers"`
}

// Description returns a description of the program
func (Args) Description() string {
	return "Display today's Bible reading from the four Gospels (Matthew, Mark, Luke, John).\n" +
		"Reads 12 verses per day based on the current date."
}

func main() {
	var args Args
	arg.MustParse(&args)

	// Load Bible data
	biblePath := "books/pol_nbg.json"
	bibleData, err := bible.LoadBible(biblePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading Bible data: %v\n", err)
		os.Exit(1)
	}

	// Initialise date progression
	dateProg := bible.NewDateProgression(bibleData)

	// Get current date
	currentDate := time.Now()

	// Calculate position for today
	todayBookmark, err := dateProg.GetPositionForDate(currentDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calculating date position: %v\n", err)
		os.Exit(1)
	}

	// Initialise formatter
	formatter := bible.NewFormatter()

	// Print header (unless plain mode)
	if !args.Plain {
		fmt.Println(formatter.FormatHeader(todayBookmark))
	}

	// Print snippet
	fmt.Println(formatter.ExtractAndFormat(bibleData, todayBookmark))
}

