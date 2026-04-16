package bible

// Verse represents a single Bible verse
type Verse struct {
	BookName  string `json:"book_name"`
	Book      int    `json:"book"`
	Chapter   int    `json:"chapter"`
	VerseNum  int    `json:"verse"`
	Text      string `json:"text"`
}

// VerseRange represents a range of verses (inclusive)
type VerseRange struct {
	Book         int
	Chapter      int
	FirstVerse   int
	LastVerse    int
}

// BookIndex represents the four Gospels
type BookIndex int

const (
	Matthew BookIndex = iota + 40
	Mark
	Luke
	John
)

func (b BookIndex) String() string {
	switch b {
	case Matthew:
		return "Ewangelia wg. Mateusza"
	case Mark:
		return "Ewangelia wg. Marka"
	case Luke:
		return "Ewangelia wg. Łukasza"
	case John:
		return "Ewangelia wg. Jana"
	default:
		return "Unknown"
	}
}

// Bookmark represents the current reading position
type Bookmark struct {
	Book        BookIndex `toml:"book"`
	Chapter     int       `toml:"chapter"`
	FirstVerse  int       `toml:"first_verse"`
	SecondVerse int       `toml:"second_verse"`
}