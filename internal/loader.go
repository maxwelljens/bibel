package bible

import (
	"encoding/json"
	"os"
	"strings"
)

// Bible represents the complete Bible data
type Bible struct {
	Metadata struct {
		Name string `json:"name"`
		Lang string `json:"lang_short"`
	} `json:"metadata"`
	Verses []Verse `json:"verses"`

	// Index for quick lookups
	byBookChapterVerse map[int]map[int]map[int]*Verse
}

// LoadBible loads Bible data from a JSON file
func LoadBible(filePath string, verbose bool) (*Bible, error) {
	logger := NewVerboseLogger(verbose)
	logger.Info("Loading Bible data from file")
	logger.Path(filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Warning("Failed to read Bible file: %v", err)
		return nil, err
	}

	var bible Bible
	if err := json.Unmarshal(data, &bible); err != nil {
		return nil, err
	}

	bible.buildIndex()
	logger.Success("Bible data loaded successfully")
	logger.WithFields(map[string]any{
		"metadata.name": bible.Metadata.Name,
		"metadata.lang": bible.Metadata.Lang,
		"verse_count":   len(bible.Verses),
	}, "Bible metadata")

	return &bible, nil
}

// buildIndex creates a lookup index for verses
func (b *Bible) buildIndex() {
	b.byBookChapterVerse = make(map[int]map[int]map[int]*Verse)

	for i := range b.Verses {
		verse := &b.Verses[i]

		if b.byBookChapterVerse[verse.Book] == nil {
			b.byBookChapterVerse[verse.Book] = make(map[int]map[int]*Verse)
		}

		if b.byBookChapterVerse[verse.Book][verse.Chapter] == nil {
			b.byBookChapterVerse[verse.Book][verse.Chapter] = make(map[int]*Verse)
		}

		b.byBookChapterVerse[verse.Book][verse.Chapter][verse.VerseNum] = verse
	}
}

// GetVerse retrieves a specific verse
func (b *Bible) GetVerse(book, chapter, verse int) (*Verse, bool) {
	if chapterMap, ok := b.byBookChapterVerse[book]; ok {
		if verseMap, ok := chapterMap[chapter]; ok {
			if verse, ok := verseMap[verse]; ok {
				return verse, true
			}
		}
	}
	return nil, false
}

// GetChapterText retrieves the full text of a chapter
func (b *Bible) GetChapterText(book, chapter int) string {
	if chapterMap, ok := b.byBookChapterVerse[book]; ok {
		if verseMap, ok := chapterMap[chapter]; ok {
			var text strings.Builder
			// We need to get verses in order
			// Since we don't know the max verse number, we'll iterate
			// This could be optimized but works for now
			for i := 1; ; i++ {
				if verse, ok := verseMap[i]; ok {
					text.WriteString(verse.Text + " ")
				} else {
					break
				}
			}
			return text.String()
		}
	}
	return ""
}

// CountVersesInChapter returns the number of verses in a chapter
func (b *Bible) CountVersesInChapter(book, chapter int) int {
	if chapterMap, ok := b.byBookChapterVerse[book]; ok {
		if verseMap, ok := chapterMap[chapter]; ok {
			return len(verseMap)
		}
	}
	return 0
}

// GetVerseRange retrieves verses in a range
func (b *Bible) GetVerseRange(book, chapter, firstVerse, lastVerse int) []*Verse {
	var verses []*Verse

	if chapterMap, ok := b.byBookChapterVerse[book]; ok {
		if verseMap, ok := chapterMap[chapter]; ok {
			for i := firstVerse; i <= lastVerse; i++ {
				if verse, ok := verseMap[i]; ok {
					verses = append(verses, verse)
				}
			}
		}
	}

	return verses
}

