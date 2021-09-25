package main
 
import _ "embed"
import (
  "fmt"
  "encoding/gob"
  "os"
  "strconv"
  "strings"
)

// Terminal Colours
const RESET = "\033[0m"
const RED = "\033[31m"
const PURPLE = "\033[35m"
const CYAN = "\033[36m"

// Globals
const BIBLE_FILE = "bible.dat"
const HELP =
`Usage: bibel [FLAG]... [BOOK] [CHAPTER] [VERSE]...
Command line Bible utility.

  -h, --help       Show this help message and exit.
  -v, --version    Show version number and exit.

If no arguments are specified, the entire Bible is printed. Specify BOOK to
print a book from the Bible; specify BOOK and CHAPTER to print a chapter from
the Bible; specify BOOK, CHAPTER, and VERSE to print a verse or range of
verses from the Bible.

Names of the books are standard short names (e.g. "1 Maccabees" instead of
"The First Book of Maccabees")

Example usage:

  bibel | grep "inherit the earth"
  bibel "1 mac" 1 1:2
  bibel genesis 2 

For more information, visit <https://www.github.com/maxwelljens/bibel/>
Program written by Maxwell Jensen (c) 2021
Licensed under European Union Public Licence 1.2.

Bible used is World English Bible, which is in the public domain. It is a
trademark of <https://www.ebible.org/>.`
const VERSION = "VERSION: 1.0.0"
//go:embed bible.dat
var BIBEL []byte
var BIBLE []Book

type Book struct {
  Name, Chapter string
  Body []string
}

func decodeBible(encBooks []byte) ([]Book, error) {
  var decBooks []Book
  if file, err0 := os.Open(BIBLE_FILE); err0 != nil {
    fmt.Println("Opening file error:", err0)
    return nil, err0
  } else {
    defer file.Close()
    dec := gob.NewDecoder(file)
    err1 := dec.Decode(&decBooks)
    if err1 != nil {
      fmt.Println("Decoding error:", err1)
      return nil, err1
    }
  }
  return decBooks, nil
}

func insensitiveContains(a, b string) bool {
  return strings.Contains(strings.ToLower(a), strings.ToLower(b))
}

func parseVerseNum(str string) (x, y int) {
  // If no colons, the split doesn't do anything
  split := strings.Split(str, ":")
  if len(split) > 2 && len(split) != 0 {
    fmt.Println(RED + "Error: Too many colons." + RESET)
    return 0, 0
  }

  // Check if value left of colon is a parseable integer
  if i, err := strconv.Atoi(split[0]); err != nil || i <= 0 {
    fmt.Println(RED + "Error: Value left of colon is borked." + RESET)
    return 0, 0
  } else {
    x = i
  }

  // Return only x if split did nothing; no colons
  if len(split) == 1 {
    return x, 0
  }

  // Check if value right of colon is a perseable integer
  if i, err := strconv.Atoi(split[1]); err != nil || i <= 0 {
    fmt.Println(RED + "Error: Value right of colon is borked." + RESET)
    return 0, 0
  } else {
    y = i
  }

  // Final check
  if x > y {
    fmt.Println(RED + "Error: Number left of colon is greater than to the right of it." + RESET)
    return 0, 0
  }

  // If all is good
  return x, y
}

func printVerses(body []string) {
  // Pretty print each line, without brackets
  for _, v := range body {
    fmt.Println(v)
  }
}

func printVersesEnum(body []string) {
  // Pretty print each line, without brackets, with verse number
  // NOTE: Only enumerates within the context of the slice, not actual verse number
  for i, v := range body {
    i++ // Increase the index number so it counts from 1 instead of 0
    fmt.Println(CYAN + "[" + fmt.Sprint(i) + "]" + RESET, v)
  }
}

func printBible(bible []Book) {
  // Print entire Bible
  for _, v := range bible {
    for i, x := range v.Body {
      fmt.Println(PURPLE + v.Name + RESET, CYAN + v.Chapter, "[" + fmt.Sprint(i) + "]" + RESET, x)
    }
  }
}

func printBook(bible []Book, args []string) {
  // Print book from Bible
  for _, v := range bible {
    if insensitiveContains(v.Name, args[0]) {
      fmt.Println(PURPLE + v.Name + RESET)
      fmt.Println(CYAN + v.Chapter + RESET)
      printVersesEnum(v.Body)
    }
  }
}

func printChapter(bible []Book, args []string) {
  // Print chapter from Bible
  for _, v := range bible {
    if insensitiveContains(v.Name, args[0]) {
       if strings.Contains(v.Chapter, " " + args[1] + ".") {
         fmt.Println(PURPLE + v.Name + RESET)
         fmt.Println(CYAN + v.Chapter + RESET)
         printVersesEnum(v.Body)
      }
    }
  }
}

func seekVerse(bible []Book, args []string) {
  // Print verse or verse range
  for _, v := range bible {
    if insensitiveContains(v.Name, args[0]) {
      if strings.Contains(v.Chapter, " " + args[1] + ".") {
        if x, y := parseVerseNum(args[2]); x == 0 {
          fmt.Println(RED + "Error: Parsing error." + RESET)
        } else {
          x-- // User counts normally, but we index from 0
          fmt.Println(PURPLE + v.Name + RESET)
          fmt.Println(CYAN + v.Chapter + RESET)
          if y <= 0 {
            fmt.Println(v.Body[x])
          } else {
            printVerses(v.Body[x:y])
          }
        }
      }
    }
  }
}

func parseArguments(x []string) {
  switch y := len(x); y {
  case 0: printBible(BIBLE)
  case 1: printBook(BIBLE, x)
  case 2: printChapter(BIBLE, x)
  case 3: seekVerse(BIBLE, x)
  default: fmt.Println(RED + "Error: Too many arguments" + RESET)
  }
}

func parseFlags (args []string) {
  if len(args) <= 0 {
   return
  }
  switch args[0] {
  case "-h", "--help":
    fmt.Println(HELP)
    os.Exit(0)
  case "-v", "--version":
    fmt.Println(VERSION)
    os.Exit(0)
  }
}
 
func main() {
  // Decode Bible
  var err error
  if BIBLE, err = decodeBible(BIBEL); err != nil {
    fmt.Println("Bible embedding error:", err)
  // If successful, it is good
  } else {
    parseFlags(os.Args[1:])
    parseArguments(os.Args[1:])
  }
}
