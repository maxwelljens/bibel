# `bibel` - Bible Verse CLI Utility

A Go reimplementation of an OCaml utility for displaying Bible verses. This
tool reads a bookmark file and displays a snippet of Bible text (typically 12
verses), then updates the bookmark for the next run.

## Features

- **Progressive Reading**: Automatically advances through the first four books
of the New Testament (Matthew, Mark, Luke, John)
- **Smart Sizing**: Default snippet size is 12 verses, extends to end of
chapter if less than 12 verses remain
- **Bookmark Persistence**: Saves reading position in a TOML file
- **Color Output**: ANSI color-coded output (green for book/chapter, yellow for
verse range)
- **Simple CLI**: Run once to display current snippet and advance bookmark

## Installation

```bash
go build ./cmd/bibel.go
```

## Usage

```bash
./bibel
```

The program will:
1. Read the current bookmark from `bookmark.toml` (creates default if not
   exists)
2. Display the Bible snippet for that bookmark with colored headers
3. Calculate and save the next bookmark

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

Books 40-43 correspond to the Evangelion:
- 40: Matthew (Mateusza)
- 41: Mark (Marek)
- 42: Luke (Łukasz)
- 43: John (Jana)

## Bookmark Format

The bookmark file (`bookmark.toml`) uses TOML format:
```toml
book = 40        # Book index (40-43)
chapter = 1      # Chapter number
first_verse = 1  # First verse in range
second_verse = 12 # Last verse in range
```

## Project Structure

```
.
├── cmd/bibel.go                 # Main CLI entry point
├── internal/bible/
│   ├── verse.go                 # Data structures
│   ├── loader.go                # JSON loading and indexing
│   ├── bookmark.go              # Bookmark management
│   └── formatter.go             # Output formatting
├── books/pol_nbg.json           # Bible data
└── old_code.ml                  # Original OCaml implementation
```

## Logic Details

- **Advancement**: After displaying verses 1-12, next bookmark is 13-24 (or to
end of chapter)
- **Chapter Boundaries**: At chapter end, moves to next chapter in same book
- **Book Boundaries**: At book end, moves to next book (loops from John back to
Matthew)
- **Lookahead Rule**: If less than 12 verses remain after current snippet,
extends current snippet to end of chapter

## Development

```bash
# Build
go build ./cmd/bibel.go

# Run with current bookmark
./bibel

# Reset bookmark to Matthew 1:1-12
echo 'book = 40
chapter = 1
first_verse = 1
second_verse = 12' > bookmark.toml
```

## License

See the original Bible data for copyright information. The code is open source.
