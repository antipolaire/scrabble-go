package game

type Tile struct {
	Letter      string
	LetterScore int
}

func NewTile(letter string, letterScore int) *Tile {
	tile := &Tile{
		Letter:      letter,
		LetterScore: letterScore,
	}
	return tile
}
