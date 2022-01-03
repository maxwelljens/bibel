// Program written by Maxwell Jensen (c) 2021
// Licensed under European Union Public Licence 1.2.
// For more information, visit <https://www.github.com/maxwelljens/bibel/>

use clap::{App, Arg};
use colored::*;
use rayon::iter::ParallelBridge;
use rayon::prelude::*;
use rust_embed::RustEmbed;

#[derive(RustEmbed)]
#[folder = "src/web_bible/"]
struct Bible;

const VERSION: &str = "1.1.0";
const LICENCE: &str = "Program written by Maxwell Jensen (c) 2021
Licensed under European Union Public Licence 1.2.
For more information, visit <https://www.github.com/maxwelljens/bibel/>

Bible used is World English Bible, which is in the public domain. It is a
trademark of <https://www.ebible.org/>.";

fn main() {
  let args = App::new("bibel")
    .version(VERSION)
    .author("Maxwell Jensen <maxwelljensen@posteo.net>")
    .about("Bible CLI utility")
    .arg(
      Arg::new("colour")
        .short('c')
        .long("colour")
        .help("Print annotations with colour"),
    )
    .arg(
      Arg::new("licence")
        .short('l')
        .long("licence")
        .help("View licence information"),
    )
    .arg(
      Arg::new("verbose")
        .short('v')
        .long("verbose")
        .help("Show output with additional information"),
    )
    .arg(Arg::new("BOOK").index(1).help("Book from the Bible"))
    .arg(Arg::new("CHAPTER").index(2).help("Chapter in <BOOK>"))
    .arg(Arg::new("VERSE").index(3).help("Verse in <CHAPTER>"))
    .get_matches();

  // If licence flag is present, ignore executing anything else
  if args.is_present("licence") {
    println!("{}", LICENCE);
  } else if args.is_present("VERSE") {
    par_print_verse(
      args.value_of("BOOK"),
      args.value_of("CHAPTER"),
      args.value_of("VERSE"),
      args.is_present("verbose"),
      args.is_present("colour"),
    );
  } else if args.is_present("CHAPTER") {
    par_print_chapter(
      args.value_of("BOOK"),
      args.value_of("CHAPTER"),
      args.is_present("colour"),
    );
  } else if args.is_present("BOOK") {
    par_print_book(args.value_of("BOOK"), args.is_present("colour"));
  } else {
    par_print_bible(args.is_present("colour"));
  }
}

/// Unwrap a file from the Bible struct
macro_rules! file_unwrap {
  ($a:expr) => {
    std::str::from_utf8(Bible::get($a.as_ref()).unwrap().data.as_ref()).unwrap()
  };
}

/// Compare input string $a and unwrapped file $b at line $nth
macro_rules! str_match {
  ($a:expr, $b:expr, $nth:expr) => {
    file_unwrap!($b)
      .lines()
      .nth($nth)
      .unwrap()
      .to_ascii_lowercase()
      .as_str()
      .contains($a.unwrap().trim().to_ascii_lowercase().as_str())
  };
}

/// Compare input number $a and unwrapped file $b at line $nth for exact number match
/// (Meaning: searching for chapter "1" will not also give you chapter "10", "11", etc.)
macro_rules! num_match {
  ($a:expr, $b:expr, $nth:expr) => {
    file_unwrap!($b)
      .lines()
      .nth($nth)
      .unwrap()
      .contains(format!(" {}.", $a.unwrap().trim()).as_str())
  };
}

/// Format lines of unwrapped file $a with full annotations
macro_rules! fmt_verses {
  ($a:expr, $title:expr, $chapter:expr, $c:expr) => {
    for (i, line) in $a {
      if i == 0 {
        $title = line;
      } else if i == 1 {
        $chapter = line;
      } else {
        if $c {
          println!(
            "{} {} [{}] {}",
            $title.magenta(),
            $chapter.blue(),
            (i - 1).to_string().blue(),
            line
          );
        } else {
          println!("{} {} [{}] {}", $title, $chapter, (i - 1), line);
        }
      }
    }
  };
}

/// Print the entire Bible, using Rayon-powered parallelism
fn par_print_bible(is_coloured: bool) {
  Bible::iter()
    .par_bridge()
    // Enter the directory of Bible
    .for_each(|dir| {
      dir
        .par_lines()
        // On each entry in directory of Bible, unwrap, iterate, and enumerate string lines:
        .for_each(|file| {
          let mut body: Vec<(usize, String)> = Vec::new();
          for (i, line) in file_unwrap!(file).lines().enumerate() {
            body.push((i, line.to_string()));
          }
          // Format the lines and print them
          let mut title: String = String::from("ERROR_TITLE");
          let mut chapter: String = String::from("ERROR_CHAPTER");
          fmt_verses!(body, title, chapter, is_coloured);
        })
    });
}

/// Print a book from the Bible, using Rayon-powered parallelism
fn par_print_book(book: Option<&str>, is_coloured: bool) {
  Bible::iter().par_bridge().for_each(|dir| {
    dir.par_lines().for_each(|file| {
      let mut body: Vec<(usize, String)> = Vec::new();
      if str_match!(book, file, 0) {
        for (i, line) in file_unwrap!(file).lines().enumerate() {
          body.push((i, line.to_string()));
        }
        // Format the lines and print them
        let mut title: String = String::from("ERROR_TITLE");
        let mut chapter: String = String::from("ERROR_CHAPTER");
        fmt_verses!(body, title, chapter, is_coloured);
      }
    })
  });
}

/// Print a chapter from the Bible, using Rayon-powered parallelism
fn par_print_chapter(book: Option<&str>, chapter: Option<&str>, is_coloured: bool) {
  Bible::iter().par_bridge().for_each(|dir| {
    dir.par_lines().for_each(|file| {
      let mut body: Vec<(usize, String)> = Vec::new();
      if str_match!(book, file, 0) && num_match!(chapter, file, 1) {
        for (i, line) in file_unwrap!(file).lines().enumerate() {
          body.push((i, line.to_string()));
        }
        // Format the lines and print them
        let mut title: String = String::from("ERROR_TITLE");
        let mut chapter: String = String::from("ERROR_CHAPTER");
        fmt_verses!(body, title, chapter, is_coloured);
      }
    })
  });
}

/// Print a verse or series of verses from the Bible, using Rayon-powered parallelism
fn par_print_verse(
  book: Option<&str>,
  chapter: Option<&str>,
  verse: Option<&str>,
  is_verbose: bool,
  is_coloured: bool,
) {
  // Do basic error checking before processing
  let parsed_verse: Vec<&str> = verse.unwrap().split_terminator(':').collect();
  if parsed_verse.len() > 2 {
    eprintln!(
      "{} More than one colon was found in your <VERSE> input.",
      "error:".red()
    );
    std::process::exit(1);
  }
  for num in parsed_verse {
    for ch in num.chars() {
      if !ch.is_numeric() {
        eprintln!(
          "{} Non-numeric characters were found in your <VERSE> input.",
          "error:".red()
        );
        std::process::exit(1);
      }
    }
  }

  Bible::iter().par_bridge().for_each(|dir| {
    dir.par_lines().for_each(|file| {
      let parsed_verse: Vec<&str> = verse.unwrap().split_terminator(':').collect();
      let verse_lower_num = parsed_verse[0].parse::<usize>().unwrap();
      let verse_upper_num = parsed_verse[1].parse::<usize>().unwrap();

      let mut body: Vec<(usize, String)> = Vec::new();
      if str_match!(book, file, 0) && num_match!(chapter, file, 1) {
        for (i, line) in file_unwrap!(file).lines().enumerate() {
          body.push((i, line.to_string()));
        }
      }

      // Check if anything was added to body at all
      if body.len() > 0 {
        // Check if verse numbers are in range
        if body.len() < verse_lower_num && is_verbose {
          eprintln!(
            "{} There's only {} verses in {}, {}, but your start <VERSE> is {}. Skipping.",
            "warning:".yellow(),
            (body.len() - 2),
            body[0].1.strip_suffix(".").unwrap(),
            body[1].1.strip_suffix(".").unwrap(),
            parsed_verse[0]
          );
        }

        if body.len() < verse_upper_num && is_verbose {
          eprintln!(
            "{} There's only {} verses in {}, {}.",
            "warning:".yellow(),
            (body.len() - 2),
            body[0].1.strip_suffix(".").unwrap(),
            body[1].1.strip_suffix(".").unwrap()
          );
        }
        // If all goes well, format the lines and print them, within bounds
        let mut title: String = String::from("ERROR_TITLE");
        let mut chapter: String = String::from("ERROR_CHAPTER");
        for (i, line) in body {
          if i == 0 {
            title = line;
          } else if i == 1 {
            chapter = line;
          } else if i > verse_lower_num && i < verse_upper_num {
            if is_coloured {
              println!(
                "{} {} [{}] {}",
                title.magenta(),
                chapter.blue(),
                (i - 1).to_string().blue(),
                line
              );
            } else {
              println!("{} {} [{}] {}", title, chapter, (i - 1), line);
            }
          }
        }
      }
    })
  });
}
