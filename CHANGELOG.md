# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.1] - 2026-04-20

### Fixed

- Fixed syntax error in `cmd/bibel.go` where closing brace for reading mode
validation was misplaced, preventing `-r/--reading` command line option from
taking effect

## [1.0.0] - 2026-04-18

### Added

- **Comprehensive `man` page**: `bibel.1` with 15 sections that should cover
all aspects of the program

### Changed

- **Revised documentation**: `README.md` updated to a more concise page (80
lines, was ~245)

### Fixed

- **Output mode distinction**: Plain (`-p/--plain`) and formatted
(`-f/--formatted`) modes now behave correctly:
  - Plain mode forces `useColours = false` for true plain text output
  - Formatted mode respects `config.Formatter.UseColours` setting
- **Formatter logic**: Restructured to handle all flag combinations (`-p
--numbered`, `-f --numbered`, etc.)
- `-p/--plain` and `-f/--formatted` previously operated identically (both with
ANSI colours)
- `--config` flag existed but didn't actually load custom config files
- Fixed part where `--numbered` flag alone didn't trigger formatted output mode
- Edge cases with flag combinations (`-p --numbered`, `-f --paragraphs`)

## [0.11.0] - 2026-04-18

### Added

- **New flag**: `-v/--verbose` flag to print with more information if needed.

### Fixed

- Add missing newlines at end of `easterprogression.go` and `config.go`
  (cosmetic)

## [0.10.0] - 2026-04-18

### Changed

- **ANSI colour**: Changed from hex colour codes to ANSI terminal colours
- Removed light/dark mode detection and adaptive colour logic
- Changed "Adaptive Styling" to "ANSI Colours" in `README`

### Removed

- Removed `border_colour`, `header_colour`, `text_colour`, `quit_colour` from
TOML configuration
- `HasDarkBackground()` check and `lightDark()` helper function
- All `#RRGGBB` colour codes (replaced with ANSI codes)

## [0.9.0] - 2026-04-17

### Added

- **Easter progress bar**: Fancy progress bar at bottom of TUI showing time
until Easter, which is configurable; Orthodox (default) or Roman Catholic.
Shows days/hours/minutes until next Easter with percentage progress.
- New `-l/--latin` flag to use Roman Catholic Easter instead of Eastern
Orthodox.
- Added vendored `eastertime` package with accurate Orthodox and Roman Catholic
algorithms for telling Easter time.
- **Configuration option**: `easter_type` configuration field in TOML (default:
"orthodox")

## [0.8.0] - 2026-04-17

### Added

- **Verse numbering option**: New `-n/--numbered` flag prints each verse on a
numbered line with verse number
- **Paragraph handling option**: New `-g/--paragraphs` flag renders pilcrows
(¶) as blank lines instead of ignoring them
- **TUI numbered styling**: Verse numbers use darker shade styling via
`lipgloss` for better readability
- **Configuration integration**: Both options configurable via TOML file in
`[formatter]` section (`numbered` and `paragraphs` fields)

### Changed

- **Formatter system enhanced**: New `FormatSnippetWithOptions()` and
`FormatSnippetForTUI()` methods for advanced formatting
- **TUI width handling**: Improved width calculations for numbered verses to
prevent text clipping
- **Paragraph logic refined**: Fixed newline handling to only insert blank
lines for pilcrows, not between all verses
- **Phased out short flag**: `--generate-config` no longer has the `-g` short
flag. This is now under `--paragraphs`.

## [0.7.0] - 2026-04-17

### Added

- **XDG data directory support**: Automatic Bible file discovery in
`$XDG_DATA_HOME/bibel`
  - When `bible_path` is empty in configuration, program searches
  `$XDG_DATA_HOME/bibel` for JSON files
  - Uses first `.json` file found in the directory. Provides helpful error
  messages when directory or Bible files are missing.
- **Cross-platform Bible data management**: Leverages `xdg` library for
consistent data directory paths across platforms

### Changed

- **Default configuration**: `bible_path` now defaults to empty string (uses
XDG-stadnard data directory)
- **Example configuration**: Updated `config.example.toml` with documentation
about XDG data directory usage
- **Error handling**: Improved error messages for missing Bible data files

### Fixed

- **Configuration generation**: `--generate-config` flag creates config with
empty `bible_path` for XDG compatibility
- **Edge case handling**: Prevent crashes when Bible file doesn't contain books
for selected reading mode
- **Date progression robustness**: Added check for zero total verses to prevent
divide-by-zero panics

## [0.6.0] - 2026-04-17

### Added

- **Multiple reading modes** via new `-r/--reading` command line option:
  - `evangelion` (default): Gospels only (Matthew, Mark, Luke, John)
  - `new_testament`: Entire New Testament (books 40-66)
  - `old_testament`: Old Testament only (books 1-39)
  - `bible`: The entire Bible (books 1-66)
- **Flexible Bible source support**: Can now work with any Bible translation
  - Removed hardcoded Polish book names
  - Reads `book_name` field directly from JSON data

### Changed

- **Date progression system enhanced**:
  - `NewDateProgressionWithReadingMode()` accepts reading mode parameter
  - Calculates verse positions based on selected reading scope
  - Proper wrapping at boundaries of each reading mode
- **Formatter now requires Bible reference**:
  - `NewFormatterWithConfig()` now accepts `*Bible` parameter
  - Formatter retrieves book names dynamically from loaded Bible data
  - Updated all formatter initialization calls in codebase
- **Book name handling refactored**:
  - Removed `polishBookNames` map from `verse.go`
  - `BookIndex.String()` returns numeric fallback representation
  - Actual book names sourced from Bible JSON metadata
- **Configuration system enhanced**:
  - Added `reading_mode` configuration option with validation
  - Command line `-r/--reading` overrides config setting
  - Improved default value handling for empty configuration fields

### Fixed

- **Configuration validation**: Allow empty reading mode during initial load
- **Build errors**: Fixed Go compilation issues in configuration unmarshalling
- **Backward compatibility**: Maintains existing behavior for default mode

## [0.5.0] - 2026-04-17

### Added

- **TOML configuration file** system with XDG Base Directory Specification
compliance
  - Configuration file location: `$XDG_CONFIG_HOME/bibel/config.toml` (default:
  `~/.config/bibel/config.toml`). macOS and Windows also supported.
  - Uses `xdg` for XDG compliance and `viper` for configuration management
- **Configuration generation**: `--generate-config` flag creates default
configuration file
- **Comprehensive TUI aesthetic configuration**:
  - `border_style`: "rounded", "double", "single", or "hidden"
  - `border_colour`: Custom hex colour for box borders
  - `header_colour`: Custom hex colour for header text
  - `text_colour`: Custom hex colour for verse content
  - `quit_colour`: Custom hex colour for quit message
  - `show_quit_message`: Boolean to show/hide "Press q to quit..." message
- **Formatter configuration**:
  - `header_format`: Customisable header template with variables: `{book}`,
  `{chapter}`, `{first_verse}`, `{second_verse}`
  - `use_colours`: Enable/disable ANSI colour output in formatted mode
- **Date progression configuration**:
  - `verses_per_day`: Customisable number of verses per day (default: 12)
- **Command-line override**: CLI flags (`--plain`, `--formatted`) override
configuration settings
- **Configuration validation**: Ensures valid output modes, border styles, and
positive verse counts

### Changed

- **Formatter behaviour**: `FormatHeader()` now reads `header_format` from
configuration instead of hardcoded format
- **TUI styling**: All aesthetic properties (colours, border style) now
configurable via TOML file
- **Date progression**: `NewDateProgressionWithConfig()` accepts custom
`verses_per_day` parameter
- **Application flow**: Configuration loaded before processing, with
command-line arguments taking precedence
- **Package structure**: Added `internal/config.go` for configuration
management

## [0.4.0] - 2026-04-16

### Added

- **TUI framework** using `bubbletea` for interactive terminal user interface
- **Styled output** with `lipgloss` for terminal-adaptive colours and borders
- **Interactive mode**: Verse displayed in a rounded border box with "Press q
to quit..." message
- **Smart TTY detection**: Automatically falls back to formatted output in
non-interactive environments
- **Enhanced argument handling**: `-p/--plain` now outputs plain text without
TUI or formatting
- **New package structure**: `internal/tui/` containing TUI model and styling
components

### Changed

- **Program flow**: Default mode now launches interactive TUI when run in a
terminal
- **Plain mode behaviour**: `--plain` flag now outputs only verse text (no
header or TUI)
- **Terminal detection**: Non-TTY environments show formatted output with ANSI
colours
- **Dependencies**: Added `bubbletea` and `lipgloss` for TUI functionality

## [0.3.0] - 2026-04-16

### Added

- Argument parsing using `go-arg` library
- `-p/--plain` flag to output plain text without chapter name or verse numbers
- Program description and help text

## [0.2.0] - 2026-04-16

### Added

- Complete Go migration from original OCaml implementation
- Date-based Bible verse progression system
  - Calculates position based on day of year (1st January = Matthew 1:1-12)
  - Shows 12 verses per day, cycling through all four Gospels yearly
- ANSI color-coded output (green for book/chapter, yellow for verse range)
- Smart lookahead rule: extends to end of chapter if less than 12 verses remain
- Support for Polish "Nowa Biblia Gdańska" translation (`books/pol_nbg.json``)
- Complete test suite with edge case validation
- `README`` documentation with usage examples and algorithm explanation

### Changed

- **BREAKING**: Changed from bookmark-based progression to date-based progression
  - Previous versions used `bookmark.toml` to track reading position, new
  version calculates position from current date
- Restructured codebase with modular Go packages:
  - `internal/loader.go`: JSON Bible data loading and indexing
  - `internal/dateprogression.go`: Date-to-verse position calculation
  - `internal/formatter.go`: ANSI color formatting and text extraction
  - `internal/verse.go`: Data structures and constants
  - `internal/bookmark.go`: Legacy bookmark support (unused in date-based mode)

### Fixed

- Proper handling of chapter and book boundaries
- Correct modulo wrapping for yearly cycle (total 3779 Gospel verses)
- Leap year support via Go's standard `time.YearDay()` function
- Validation of verse ranges and error handling

### Removed

- OCaml implementation
- Dependency on `bookmark.toml` file for progression
- Regex-based verse parsing (replaced with JSON Bible data)

## [0.1.0] - 2026-04-16

### Added

- Initial Go port of OCaml Bible verse utility
- Bookmark-based progression system (TOML file)
- Support for Polish "Nowa Biblia Gdańska" translation
- Basic ANSI color formatting
- Lookahead rule implementation

### Changed

- Rewritten from OCaml to Go programming language
- Changed data format from custom to JSON
- Improved error handling and validation
