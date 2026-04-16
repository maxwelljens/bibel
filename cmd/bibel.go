package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alexflint/go-arg"
	"golang.org/x/term"
	"maxwelljensen/bibel/internal"
	"maxwelljensen/bibel/internal/tui"
)

// Args defines the command line arguments
type Args struct {
	Plain bool `arg:"-p,--plain" help:"output plain text without formatting or TUI"`
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

	// If plain mode, just print plain text (no TUI)
	if args.Plain {
		formatter := bible.NewFormatter()
		// In plain mode, we don't print the header
		fmt.Println(formatter.ExtractAndFormat(bibleData, todayBookmark))
		return
	}

	// Check if we're running in a terminal
	// If not, fall back to formatted output (similar to plain mode but with header)
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		formatter := bible.NewFormatter()
		fmt.Println(formatter.FormatHeader(todayBookmark))
		fmt.Println(formatter.ExtractAndFormat(bibleData, todayBookmark))
		return
	}

	// Create and run TUI program
	program := tui.CreateTUIProgram(bibleData, todayBookmark, false)
	if err := program.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

