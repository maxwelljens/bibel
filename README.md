# `bibel` - Bible Verse CLI Utility

A Go utility for displaying Bible verses. This tool shows Bible verses based on
the current date, displaying 12 verses per day through the four Gospels.

## Features

- **Date-Based Progression**: Automatically calculates position based on day of
year (1 January = Matthew 1:1-12)
- **Smart Sizing**: Default snippet size is 12 verses, extends to end of
chapter if less than 12 verses remain
- **Colour Output**: ANSI color-coded output (green for book/chapter, yellow
for verse range)
- **Yearly Cycle**: Progresses through all four Gospels each year, restarting
on 1 January
- **Simple CLI**: Run to display today's Bible snippet

## Installation

```bash
go build ./cmd/bibel.go
```

## Usage

```bash
./bibel
```

The program will:
1. Calculate today's date and day of year
2. Determine Bible position: day_of_year × 12 verses
3. Find corresponding verses in the Gospels
4. Display the Bible snippet with colored headers

No bookmark file is needed or created - the position is calculated from the
date alone.

## Data Format

The Bible data is in JSON format (`books/pol_nbg.json`) with the following structure:
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
2. **Verse Offset**: Multiply by 12 verses per day: `offset = (day_of_year - 1) × 12`
3. **Modulo Wrap**: Apply modulo with total Gospel verses (3779) to cycle yearly
4. **Position Mapping**: Walk through Gospels to find corresponding verses
5. **Lookahead Rule**: Extend to chapter end if less than 12 verses remain

## Project Structure

```
.
├── cmd/bibel.go                 # Main CLI entry point
├── internal/bible/
│   ├── verse.go                 # Data structures
│   ├── loader.go                # JSON loading and indexing
│   ├── dateprogression.go       # Date-based position calculation
│   ├── formatter.go             # Output formatting
│   └── bookmark.go              # Legacy bookmark management (optional)
├── books/pol_nbg.json           # Bible data
└── old_code.ml                  # Original OCaml implementation
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

# Run with today's date
./bibel

# Test with specific date (environment variable)
GOOSE_TEST_DATE=2026-01-01 ./bibel
```
