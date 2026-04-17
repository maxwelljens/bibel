# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.3.0] - 2026-04-17

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

## [1.2.0] - 2026-04-16

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

## [1.1.0] - 2026-04-16

### Added

- Argument parsing using `go-arg` library
- `-p/--plain` flag to output plain text without chapter name or verse numbers
- Program description and help text

## [1.0.0] - 2026-04-16

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
  - Previous versions used `bookmark.toml` to track reading position
  - New version calculates position from current date
  - Backward incompatible with bookmark-based usage
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
