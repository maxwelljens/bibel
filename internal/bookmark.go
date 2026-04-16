package bible

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const (
	defaultBookmarkFile = "bookmark.toml"
)

// BookmarkManager handles reading and writing bookmarks
type BookmarkManager struct {
	FilePath string
}

// NewBookmarkManager creates a new bookmark manager
func NewBookmarkManager(filePath string) *BookmarkManager {
	if filePath == "" {
		filePath = defaultBookmarkFile
	}
	return &BookmarkManager{FilePath: filePath}
}

// ReadBookmark reads a bookmark from the file
func (bm *BookmarkManager) ReadBookmark() (*Bookmark, error) {
	data, err := os.ReadFile(bm.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default bookmark if file doesn't exist
			return &Bookmark{
				Book:        Matthew,
				Chapter:     1,
				FirstVerse:  1,
				SecondVerse: 12,
			}, nil
		}
		return nil, err
	}
	
	var bookmark Bookmark
	if err := toml.Unmarshal(data, &bookmark); err != nil {
		return nil, fmt.Errorf("failed to parse bookmark file: %w", err)
	}
	
	// Validate bookmark values
	if bookmark.Book < Matthew || bookmark.Book > John {
		bookmark.Book = Matthew
	}
	if bookmark.Chapter < 1 {
		bookmark.Chapter = 1
	}
	if bookmark.FirstVerse < 1 {
		bookmark.FirstVerse = 1
	}
	if bookmark.SecondVerse < bookmark.FirstVerse {
		bookmark.SecondVerse = bookmark.FirstVerse
	}
	
	return &bookmark, nil
}

// WriteBookmark writes a bookmark to the file
func (bm *BookmarkManager) WriteBookmark(bookmark *Bookmark) error {
	data, err := toml.Marshal(bookmark)
	if err != nil {
		return fmt.Errorf("failed to marshal bookmark: %w", err)
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(bm.FilePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	
	if err := os.WriteFile(bm.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write bookmark file: %w", err)
	}
	
	return nil
}

// AdjustForLookahead adjusts the bookmark if less than 12 verses remain after it
func (bm *BookmarkManager) AdjustForLookahead(bookmark *Bookmark, bible *Bible) *Bookmark {
	versesInChapter := bible.CountVersesInChapter(int(bookmark.Book), bookmark.Chapter)
	versesRemaining := versesInChapter - bookmark.SecondVerse
	
	// If we're not at the end of the chapter and less than 12 verses remain for NEXT snippet
	if versesRemaining > 0 && versesRemaining < 12 {
		return &Bookmark{
			Book:        bookmark.Book,
			Chapter:     bookmark.Chapter,
			FirstVerse:  bookmark.FirstVerse,
			SecondVerse: versesInChapter,
		}
	}
	
	// Otherwise return the original bookmark
	return bookmark
}

// AdvanceBookmark calculates the next bookmark based on the current one
func (bm *BookmarkManager) AdvanceBookmark(current *Bookmark, bible *Bible) *Bookmark {
	versesInChapter := bible.CountVersesInChapter(int(current.Book), current.Chapter)
	
	// Case 1: We are at the end of the chapter, move to the next chapter
	if current.SecondVerse >= versesInChapter {
		nextChapter := current.Chapter + 1
		
		// Check if next chapter exists in current book
		if bible.CountVersesInChapter(int(current.Book), nextChapter) > 0 {
			return &Bookmark{
				Book:        current.Book,
				Chapter:     nextChapter,
				FirstVerse:  1,
				SecondVerse: min(12, bible.CountVersesInChapter(int(current.Book), nextChapter)),
			}
		}
		
		// End of book, move to next book or loop to Matthew
		nextBookInt := int(current.Book) + 1
		if nextBookInt > int(John) {
			nextBookInt = int(Matthew)
		}
		
		nextBook := BookIndex(nextBookInt)
		return &Bookmark{
			Book:        nextBook,
			Chapter:     1,
			FirstVerse:  1,
			SecondVerse: min(12, bible.CountVersesInChapter(int(nextBook), 1)),
		}
	}
	
	// Case 2: We are in the middle of a chapter, advance verses
	newFirstVerse := current.SecondVerse + 1
	newSecondVerse := min(newFirstVerse+11, versesInChapter)
	
	return &Bookmark{
		Book:        current.Book,
		Chapter:     current.Chapter,
		FirstVerse:  newFirstVerse,
		SecondVerse: newSecondVerse,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}