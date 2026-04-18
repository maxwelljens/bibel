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
	Reading        string `arg:"-r,--reading" help:"reading mode: evangelion (Gospels), new_testament, old_testament, bible (default: evangelion)"`
	Plain          bool   `arg:"-p,--plain" help:"output plain text without formatting or TUI"`
	Formatted      bool   `arg:"-f,--formatted" help:"output formatted text with ANSI colours (no TUI)"`
	ConfigPath     string `arg:"-c,--config" help:"path to configuration file (default: $XDG_CONFIG_HOME/bibel/config.toml)"`
	Numbered       bool   `arg:"-n,--numbered" help:"print each verse on a numbered line with verse number"`
	Paragraphs     bool   `arg:"-g,--paragraphs" help:"render pilcrows (¶) as blank lines instead of ignoring them"`
	Latin          bool   `arg:"-l,--latin" help:"show time until Roman Catholic Easter"`
	GenerateConfig bool   `arg:"--generate-config" help:"generate a default configuration file and exit"`
	Verbose        bool   `arg:"-v,--verbose" help:"print additional runtime information to STDOUT"`
}

// Description returns a description of the program
func (Args) Description() string {
	return "Display today's Bible reading based on selected reading mode.\n" +
		"Reads 12 verses per day based on the current date.\n" +
		"Reading modes: evangelion (four Gospels), new_testament, old_testament, bible\n" +
		"Configuration file: $XDG_CONFIG_HOME/bibel/config.toml"
}

func main() {
	var args Args
	arg.MustParse(&args)

	// Handle config generation
	if args.GenerateConfig {
		if err := bible.GenerateDefaultConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Default configuration generated at: %s\n", bible.GetConfigPath())
		return
	}

	// Load configuration
	config, err := bible.LoadConfig(args.ConfigPath, args.Verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Override Easter type from command line flag if specified
	if args.Latin {
		config.EasterType = "latin"
	}

	// Override reading mode from command line if specified
	if args.Reading != "" {
		config.ReadingMode = args.Reading
	}

	// Load Bible data
	bibleData, err := bible.LoadBible(config.BiblePath, args.Verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading Bible data: %v\n", err)
		os.Exit(1)
	}

	// Create date progression calculator
	// We need to convert reading mode string to ReadingMode type
	// Based on verse.go, ReadingMode is a string type with constants
	dateProg := bible.NewDateProgression(bibleData)

	// Get today's bookmark
	bookmark, err := dateProg.GetPositionForDate(time.Now())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calculating today's reading position: %v\n", err)
		os.Exit(1)
	}

			// Determine output mode: plain, formatted, or TUI
	// Plain mode forces no colours, formatted respects config
	usePlain := args.Plain
	useFormatted := args.Formatted

	// If no explicit mode but numbered/paragraphs flags are set, use formatted mode
	if !usePlain && !useFormatted && (args.Numbered || args.Paragraphs) {
		useFormatted = true
	}

	// If we're outputting text (plain or formatted)
	if usePlain || useFormatted {
		// Determine colours: plain mode = false, formatted mode = from config
		useColours := false
		if useFormatted {
			useColours = config.Formatter.UseColours
		}

		// Create formatter with appropriate settings
		formatter := bible.NewFormatterWithFullConfig(bibleData, useColours, config.Formatter.HeaderFormat, args.Numbered, args.Paragraphs)
		header := formatter.FormatHeader(bookmark)
		content := formatter.ExtractAndFormat(bibleData, bookmark)
		fmt.Println(header)
		fmt.Println(content)
		return
	}

	// Default: TUI mode (no plain or formatted flags)
	// Check if we're in a terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintln(os.Stderr, "Error: Standard input is not a terminal. Use --plain or --formatted for non-interactive output.")
		os.Exit(1)
	}

	// Create and run TUI program
	program := tui.CreateTUIProgram(bibleData, bookmark, config)
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}