![bibel](/assets/logo.gif)

![screenshot](/assets/screenshot.svg)

![Codeberg Release](https://img.shields.io/gitea/v/release/maxwelljensen/bibel?gitea_url=https%3A%2F%2Fcodeberg.org&style=for-the-badge)
![Codeberg Licence](assets/eupl-12-badge.svg)

---

## What is bibel?

`bibel` (bee•bell) isn’t another all‑purpose Bible app. It’s a small CLI/TUI
utility with one main job: showing today’s reading from the Gospels when you
open your terminal. Think of it as a MOTD for Scripture: just 12 verses a day,
based on the calendar, without menus, search, or bloat.

That is just the default configuration, most of which can be changed to suit
one's preferences.

## How do I use bibel?

    bibel [OPTIONS]...

By default, `bibel` shows an interactive TUI with today's reading. The verses
are selected by date: day of year * 12 verses, progressing through all four
Gospels each day.

### Options

    -r, --reading        Reading mode: "evangelion" (four Gospels), "new_testament", 
                         "old_testament", "bible" (default: evangelion)
    -p, --plain          Output plain text without formatting or TUI
    -f, --formatted      Output formatted text with ANSI colours (no TUI)
    -n, --numbered       Print each verse on a numbered line with verse number
    -g, --paragraphs     Render pilcrows (¶) as blank lines instead of ignoring them
    -l, --latin          Show time until Roman Catholic Easter instead of Orthodox
    --generate-config    Generate a default configuration file and exit
    -v, --verbose        Print additional runtime information to STDOUT

### Examples

    # Default TUI mode with interactive display (shows Orthodox Easter progress)
    bibel

    # Roman Catholic Easter progress bar instead of Orthodox
    bibel --latin

    # Plain text output
    bibel --plain

    # Formatted output with ANSI colours
    bibel --formatted

    # Numbered verses (each verse on its own line with verse number)
    bibel --numbered

    # Paragraph mode (pilcrows render as blank lines)
    bibel --paragraphs

    # Generate configuration file
    bibel --generate-config

## How does `bibel` work?

`bibel` reads Bible data from a JSON file (provided example:
`books/pol_nbg.json`). It calculates today's reading position by:

1. Getting the current day of year (1-366)
2. Multiplying by 12 verses per day: `offset = (day_of_year - 1) * 12`
3. Applying modulo with total Gospel verses (3779) to cycle yearly
4. Finding the corresponding verses in the Gospels

The program extends to the end of a chapter if less than 12 verses remain for
the next day. Books 40-43 correspond to the four Gospels: Matthew, Mark, Luke,
and John.

The TUI shows a progress bar counting down to the next Easter, with Orthodox
Easter as default (calculated via Meeus Julian algorithm). Use `-l/--latin` for
Roman Catholic Easter (Delambre and Butcher's algorithm).

## Configuration

`bibel` supports configuration via a TOML file at
`$XDG_CONFIG_HOME/bibel/config.toml` (default: `~/.config/bibel/config.toml`).
macOS and Windows are also supported automatically (see
[here](https://github.com/adrg/xdg) for more details). Settings include output
mode, reading mode, border style, colours, and Easter type.

Generate a default configuration:

    bibel --generate-config

## How do I build bibel?

`bibel` is written in Go. External dependencies are `bubbletea` (TUI),
`lipgloss` (styling), `go-arg` (argument parsing), `viper` (configuration), and
`xdg` (XDG Base Directory support).

Build with:

    go build ./cmd/bibel.go

### And the Bible?

Though Bibles in the exact JSON format that `bibel` expects are not found lying
around the web everywhere, I have compiled a collection of (mainly Protestant)
Bibles in this exact JSON format over at the
[`json_bibles`](https://codeberg.org/maxwelljensen/json_bibles) repository.
This has been helpfully provided by [Bible
SuperSearch](https://www.biblesupersearch.com/bible-downloads/).

## Documentation

For more specific information, refer to the `man` page located at
`docs/bibel.1`. It is recommended to put it in `$MANPATH` for easy access via
`man bibel` on your Linux system.

## Licence

This project is licensed under [European Union Public Licence
1.2](https://joinup.ec.europa.eu/collection/eupl/eupl-text-eupl-12).
