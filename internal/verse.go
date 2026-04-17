package bible

import (
	"fmt"
)

// Verse represents a single Bible verse
type Verse struct {
	BookName string `json:"book_name"`
	Book     int    `json:"book"`
	Chapter  int    `json:"chapter"`
	VerseNum int    `json:"verse"`
	Text     string `json:"text"`
}

// VerseRange represents a range of verses (inclusive)
type VerseRange struct {
	Book       int
	Chapter    int
	FirstVerse int
	LastVerse  int
}

// BookIndex represents Bible book numbers
type BookIndex int

// BookIndex.String returns the string representation of the book index.
// Note: The actual book name should be retrieved from Bible data using the Formatter.
// This is kept for backward compatibility but returns a simple numeric representation.
func (b BookIndex) String() string {
	// Return numeric representation - actual book names should come from Bible data
	return fmt.Sprintf("Book %d", b)
}

// ReadingMode defines the scope of Bible reading
type ReadingMode string

const (
	ReadingModeEvangelion   ReadingMode = "evangelion"    // Books 40-43 (Matthew, Mark, Luke, John)
	ReadingModeNewTestament ReadingMode = "new_testament" // Books 40-66
	ReadingModeOldTestament ReadingMode = "old_testament" // Books 1-39  
	ReadingModeBible        ReadingMode = "bible"         // Books 1-66
)

// Bookmark represents the current reading position
type Bookmark struct {
	Book        BookIndex `toml:"book"`
	Chapter     int       `toml:"chapter"`
	FirstVerse  int       `toml:"first_verse"`
	SecondVerse int       `toml:"second_verse"`
}