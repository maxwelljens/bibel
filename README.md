# `bibel` - Bible Verse CLI/TUI Utility

A Go utility for displaying Bible verses with interactive terminal interface.
This tool shows Bible verses based on the current date, displaying 12 verses
per day through the four Gospels.

## Features

- **Interactive TUI**: Terminal User Interface using `bubbletea` with styled
boxes
- **Date-Based Progression**: Automatically calculates position based on day of
year (1 January = Matthew 1:1-12)
- **Smart Sizing**: Default snippet size is 12 verses, extends to end of
chapter if less than 12 verses remain
- **ANSI Colours: Uses ANSI terminal colours set by terminal emulator
- **Interactive Controls**: Press `q` to quit (MOTD-like behaviour)
- **Yearly Cycle**: Progresses through all four Gospels each year, restarting
on 1 January
- **Multiple Output Modes**: Interactive TUI, formatted ANSI, or plain text
- **Verse Numbering**: Option to print each verse on a numbered line with verse
numbers
- **Paragraph Handling**: Option to render pilcrows (¶) in source JSON Bible
file as paragraph breaks with blank lines
- **Easter Progress Bar**: Fancy progress bar showing time until Easter with
visual progress indication (default: Eastern Orthodox; Catholic with `-l` flag)

## Installation

```bash
go build ./cmd/bibel.go
```

## Usage

```bash
./bibel
```

By default, the verses the program picks are done by:
1. Calculating today's date and day of year
2. Determining Bible position: day of year * 12 verses
3. Find corresponding verses in the Gospels

### Usage Examples

```bash
# Default TUI mode with interactive display (shows Orthodox Easter progress)
./bibel

# Roman Catholic Easter progress bar instead of Orthodox
./bibel --latin

# Plain text output
./bibel --plain

# Formatted output with ANSI colours
./bibel --formatted

# Numbered verses (each verse on its own line with verse number)
./bibel --numbered

# Paragraph mode (pilcrows render as blank lines)
./bibel --paragraphs

# Combine numbered and paragraph modes
./bibel --numbered --paragraphs

# Plain output with numbered verses
./bibel --plain --numbered

# Generate configuration file
./bibel --generate-config
```

This is default behaviour. Alternative configurations are available (see
below).

## Configuration

`bibel` supports configuration via a TOML file located at
`$XDG_CONFIG_HOME/bibel/config.toml` (default: `~/.config/bibel/config.toml`).
Default macOS and Windows paths are also supported.

### Configuration Options

- **reading_mode**: Reading mode: "evangelion" (four Gospels; *default*),
"new_testament", "old_testament", "bible" (default: "evangelion")
- **output_mode**: Output mode: "tui" (interactive terminal UI), "formatted"
(ANSI-coloured text), or "plain" (plain text)
- **bible_path**: Path to Bible data file (default: "books/pol_nbg.json")

#### TUI Settings (`[tui]` section)
- **show_quit_message**: Whether to show "Press q to quit..." message (default:
true)
- **border_style**: Box border style: "rounded", "double", "single", or
"hidden" (default: "rounded")

#### Formatter Settings (`[formatter]` section)
- **use_colours**: Whether to use ANSI colours in formatted output (default:
true)
- **header_format**: Header format template with variables: {book}, {chapter},
{first_verse}, {second_verse}
- **numbered**: Whether to print each verse on a numbered line corresponding to
the verse number (default: false)
- **paragraphs**: Whether to render pilcrows (¶) as blank lines instead of
ignoring them (default: false)

#### Date Progression Settings (`[date_progression]` section)
- **verses_per_day**: Number of verses to read per day (default: 12)
- **start_date**: Start date for yearly progression (format: "1 January", empty
for current year)

#### Easter Settings
- **easter_type**: Easter calculation type: "orthodox" (default) or "latin"

### Example Configuration

See `configs/config.example.toml` in the project directory for a complete example.

### Generating Configuration

Generate a default configuration file:

```bash
./bibel --generate-config
```

### Command Line Arguments

Command line arguments override configuration file settings:

- `-r, --reading`: "evangelion" (four Gospels; *default*), "new_testament",
"old_testament", "bible" (default: "evangelion")
- `-p, --plain`: Output plain text without formatting or TUI
- `-f, --formatted`: Output formatted text with ANSI colours (no TUI)
- `-g, --generate-config`: Generate a default configuration file and exit
- `-c, --config`: Path to configuration file (not yet implemented)
- `-n, --numbered`: Print each verse on a numbered line corresponding to the verse number
- `-g, --paragraphs`: Render pilcrows (¶) as blank lines instead of ignoring them
- `-l, --latin`: Use Roman Catholic Easter instead of Eastern Orthodox (shows in TUI progress bar)

## Data Format

The Bible data is in JSON format (see `books/pol_nbg.json`) with the following
structure:
```json
{
  "metadata": { ... },
  "verses": [
    {
      "book_name": "Mateusza",
      "book": 40,
      "chapter": 1,
      "verse": 1,
      "text": "¶ KSIĘGA początku Jezusa Chrystusa..."
    },
    ...
  ]
}
```

Books 40-43 correspond to the four Gospels:
- 40: Matthew (Mateusza)
- 41: Mark (Marek)
- 42: Luke (Łukasz)
- 43: John (Jana)

## Date-Based Algorithm

The program calculates reading position as follows:

1. **Day of Year**: Get current day number (1-366)
2. **Verse Offset**: Multiply by 12 verses per day: `offset = (day_of_year - 1)
   × 12`
3. **Modulo Wrap**: Apply modulo with total Gospel verses (3779) to cycle
   yearly
4. **Position Mapping**: Walk through Gospels to find corresponding verses
5. **Lookahead Rule**: Extend to chapter end if less than 12 verses remain

## Easter Progress Bar

The TUI features a progress bar showing time until the next Easter:

### Easter Type Selection
- **Default**: Eastern Orthodox Easter (calculated via Meeus Julian algorithm)
- **Catholic**: Roman Catholic Easter via `-l/--latin` flag (Delambre and Butcher's algorithm)

### Progress Bar Features
- **Visual Progress**: Filled segments (█) showing annual cycle percentage between consecutive Easters
- **Time Display**: Shows exact days, hours, and minutes until next Easter
- **Styling: Uses `lipgloss` with ANSI-coloured segments
- **Configuration**: Easter type configurable via `easter_type` field in TOML configuration

### Easter Calculations
- **Orthodox Easter**: Uses Meeus Julian algorithm (accurate for years > 325)
- **Catholic Easter**: Uses Delambre and Butcher's algorithm (Gregorian calendar)
- **Annual Cycle Progress**: Calculated as time elapsed between consecutive Easters
- **Real-time Updates**: Progress updates continuously while TUI is running

## Project Structure

```
.
├── cmd/bibel.go             # Main CLI entry point
├── configs/                 # Default configuration files
│   └── config.example.toml  # The default configuration file
├── internal/
│   ├── config.go            # Reading and writing to config
│   ├── dateprogression.go   # Date-based position calculation
│   ├── easterprogression.go # Easter date and progress calculation
│   ├── verse.go             # Data structures
│   ├── loader.go            # JSON loading and indexing
│   ├── formatter.go         # Output formatting
│   └── tui/                 # Terminal User Interface
│       └── model.go         # bubbletea TUI model and styling
├── pkg/
│   └── eastertime.go        # Easter date calculation algorithms (Orthodox/Catholic)
├── books/pol_nbg.json       # Bible data
├── go.mod                   # Go module dependencies
└── go.sum                   # Go dependency checksums
```

## Examples

- **1 January**: Matthew 1:1-12
- **2 January**: Matthew 1:13-25 (extends to end of chapter)
- **16 April**: Mark 5:41-43
- **31 December**: Matthew 18:7-18 (year wraps around)

## Development

```bash
# Build
go build ./cmd/bibel.go

# Run
./bibel
```
