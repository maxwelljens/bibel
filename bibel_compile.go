package main
 
import (
  "bufio"
  "fmt"
  "encoding/gob"
  "os"
  "path/filepath"
)

const BIBLE_PATH = "/web_bible"
const BIBLE_FILE = "bible.dat"
var bible []Book

type Book struct {
  Name, Chapter string
  Body []string
}

func EncodeBible(books []Book) error {
  if file, err0 := os.Create(BIBLE_FILE); err0 != nil {
    fmt.Println("Creating file error:", err0)
    return err0
  } else {
    defer file.Close()
    enc := gob.NewEncoder(file)
    err1 := enc.Encode(bible)
    if err1 != nil {
      fmt.Println("Encoding error:", err1)
      return err1
    }
  }
  return nil
}

func CompileFile(file string) (Book, error) {
  // Open an input file, exit on error.
  inputFile, err := os.Open(file)
  if err != nil {
    fmt.Println(err)
  }

  // Closes the file when we leave the scope of the current function,
  // this makes sure we never forget to close the file if the
  // function can exit in multiple places.
  defer inputFile.Close()

  scanner := bufio.NewScanner(inputFile)
  book := Book{}

  // scanner.Scan() advances to the next token returning false if an error was encountered
  // Compile the book for returning
  for i := 0; scanner.Scan(); i++ {
    switch i {
    case 0: book.Name = scanner.Text()
    case 1: book.Chapter = scanner.Text()
    default: book.Body = append(book.Body, scanner.Text())
    }
  }

  // When finished scanning if any error other than io.EOF occured
  // it will be returned by scanner.Err().
  if err := scanner.Err(); err != nil {
    fmt.Println(scanner.Err())
  }
  return book, err
}
 
func SeekFile(fp string, fi os.FileInfo, err error) error {
  if err != nil {
    fmt.Println(err) // can't walk here,
    return nil       // but continue walking elsewhere
  }
  if fi.IsDir() {
    return nil // not a file. ignore.
  }
  matched, err := filepath.Match("*.txt", fi.Name())
  if err != nil {
    fmt.Println(err) // malformed pattern
    return err
  }
  if matched {
    if book, err := CompileFile(fp); err != nil {
      return err
    } else {
      bible = append(bible, book)
    }
  }
  return nil
}
 
func main() {
  myDir, err := os.Getwd()
  if err != nil {
    panic(err)
  }
  filepath.Walk((myDir + BIBLE_PATH), SeekFile)
  EncodeBible(bible)
}
