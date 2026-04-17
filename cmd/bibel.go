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
	Reading string `arg:"-r,--reading" help:"reading mode: evangelion (Gospels), new_testament, old_testament, bible (default: evangelion)"`
	Plain bool `arg:"-p,--plain" help:"output plain text without formatting or TUI"`
	Formatted bool `arg:"-f,--formatted" help:"output formatted text with ANSI colours (no TUI)"`
	GenerateConfig bool `arg:"-g,--generate-config" help:"generate a default configuration file and exit"`
	ConfigPath string `arg:"-c,--config" help:"path to configuration file (default: $XDG_CONFIG_HOME/bibel/config.toml)"`
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
	cfg, err := bible.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	// Override config with command line arguments
	if args.Reading != "" {
		cfg.ReadingMode = args.Reading
	}
	if args.Plain {
		cfg.OutputMode = "plain"
	} else if args.Formatted {
		cfg.OutputMode = "formatted"
	}
	if args.ConfigPath != "" {
		// Note: This would require modifying LoadConfig to accept a path
		// For now, we'll just use the default XDG location
		fmt.Fprintf(os.Stderr, "Note: Custom config path not yet implemented, using default location\n")
	}

	// Load Bible data
	biblePath := cfg.BiblePath
	bibleData, err := bible.LoadBible(biblePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading Bible data: %v\n", err)
		os.Exit(1)
	}

	// Initialise date progression with configured verses per day
	dateProg := bible.NewDateProgressionWithReadingMode(bibleData, cfg.DateProgression.VersesPerDay, bible.ReadingMode(cfg.ReadingMode))

	// Get current date
	currentDate := time.Now()

	// Calculate position for today
	todayBookmark, err := dateProg.GetPositionForDate(currentDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calculating date position: %v\n", err)
		os.Exit(1)
	}

	// Determine output mode based on configuration and terminal
	outputMode := cfg.OutputMode
	
	// If terminal detection suggests different mode, adjust
	if !term.IsTerminal(int(os.Stdout.Fd())) && outputMode == "tui" {
		// Not a terminal, fall back to formatted
		outputMode = "formatted"
	}

	// Handle different output modes
	switch outputMode {
	case "plain":
		formatter := bible.NewFormatterWithConfig(bibleData, false, cfg.Formatter.HeaderFormat)
		// In plain mode, we don't print the header
		fmt.Println(formatter.ExtractAndFormat(bibleData, todayBookmark))
		
	case "formatted":
		formatter := bible.NewFormatterWithConfig(bibleData, cfg.Formatter.UseColours, cfg.Formatter.HeaderFormat)
		fmt.Println(formatter.FormatHeader(todayBookmark))
		fmt.Println(formatter.ExtractAndFormat(bibleData, todayBookmark))
		
	case "tui":
		// Create and run TUI program with configuration
		program := tui.CreateTUIProgram(bibleData, todayBookmark, cfg)
		if err := program.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
		
	default:
		fmt.Fprintf(os.Stderr, "Unknown output mode: %s\n", outputMode)
		fmt.Fprintf(os.Stderr, "Valid modes: tui, formatted, plain\n")
		os.Exit(1)
	}
}