![bibel](/logo.png "bibel")
![GitHub](https://img.shields.io/github/license/maxwelljens/bibel?label=Licence) ![GitHub last commit](https://img.shields.io/github/last-commit/maxwelljens/bibel?label=Last%20Commit)
[![asciicast](https://asciinema.org/a/438002.svg)](https://asciinema.org/a/438002)

## What is bibel?

`bibel` (*bee•bell*) is a command line interface (CLI) utility for the Bible. Its primary use is for quick reference
from the command line and piping into other programs, as reading the whole Bible by itself in the command line is not
very comfortable.

## How do I use bibel?

**bibel** [FLAG]... [BOOK] [CHAPTER] [VERSE]...

If no arguments are specified, the entire Bible is printed. Specify BOOK to print a book from the Bible; specify BOOK
and CHAPTER to print a chapter from the Bible; specify BOOK, CHAPTER, and VERSE to print a verse or range of verses
from the Bible.

Example usage:

    $ bibel | grep "inherit the earth"
    Matthew. Chapter 5. [4] Blessed are the gentle, for they shall inherit the earth.

    $ bibel "1 mac" 1 1:6
    1 Maccabees.
    Chapter 1.
    After Alexander the Macedonian, the son of Philip, who came out of the land of Chittim [...]

    $ bibel gen
    [Prints the entirety of Genesis]

## How does bibel work?

`bibel` uses simple substring matching to find books and chapters. This works completely adequately for most of
everything, although `bibel john X Y` will match both *Book of John* and letters from Saint John (*1 John*, *2 John*,
and *3 John*) as a consequence.

Names of books used are their standard short names, (e.g. *1 Maccabees* instead of *The First Book of Maccabees*)

## How do I build bibel?

`bibel` is written in Rust. External dependencies are [clap](https://github.com/clap-rs/clap),
[colored](https://github.com/mackwic/colored), [rayon](https://github.com/rayon-rs/rayon), and
[rust-embed](https://github.com/pyros2097/rust-embed). The files in `/web_bible` are embedded during compilation of
`main.rs`, and are accessed thusly and portably at runtime.
The format of text that is expected during embedding is as follows:

- **First line**: Book title
- **Second line**: Chapter
- **Every line after**: Verse
- Each text file is a chapter

## Licence

This project is licensed under [European Union Public Licence
1.2](https://joinup.ec.europa.eu/collection/eupl/eupl-text-eupl-12). The Bible used is the World English Bible, which
is in the public domain. World English Bible is a trademark of [eBible.org](https://www.ebible.org/).
