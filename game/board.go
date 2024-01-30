package game

const (
	// NF Special field types:
	NF = iota // NormalField
	DL        // DoubleLetterField
	TL        // TripleLetterField
	DW        // DoubleWordField
	TW        // TripleWordField
	CS        // CenterStarField
	// NFColor Special field type colors:
	NFColor = 0xffffff // White
	DLColor = 0x0000ff // Light blue
	TLColor = 0x000080 // Blue
	DWColor = 0xff00ff // Pink
	TWColor = 0xff0000 // Red
	CSColor = 0x101010
)

// Special field matrix:
var specialFields = [][]int{
	{TW, NF, NF, DL, NF, NF, NF, TW, NF, NF, NF, DL, NF, NF, TW},
	{NF, DW, NF, NF, NF, TL, NF, NF, NF, TL, NF, NF, NF, DW, NF},
	{NF, NF, DW, NF, NF, NF, DL, NF, DL, NF, NF, NF, DW, NF, NF},
	{DL, NF, NF, DW, NF, NF, NF, DL, NF, NF, NF, DW, NF, NF, DL},
	{NF, NF, NF, NF, DW, NF, NF, NF, NF, NF, DW, NF, NF, NF, NF},
	{NF, TL, NF, NF, NF, TL, NF, NF, NF, TL, NF, NF, NF, TL, NF},
	{NF, NF, DL, NF, NF, NF, DL, NF, DL, NF, NF, NF, DL, NF, NF},
	{TW, NF, NF, DL, NF, NF, NF, CS, NF, NF, NF, DL, NF, NF, TW},
	{NF, NF, DL, NF, NF, NF, DL, NF, DL, NF, NF, NF, DL, NF, NF},
	{NF, TL, NF, NF, NF, TL, NF, NF, NF, TL, NF, NF, NF, TL, NF},
	{NF, NF, NF, NF, DW, NF, NF, NF, NF, NF, DW, NF, NF, NF, NF},
	{DL, NF, NF, DW, NF, NF, NF, DL, NF, NF, NF, DW, NF, NF, DL},
	{NF, NF, DW, NF, NF, NF, DL, NF, DL, NF, NF, NF, DW, NF, NF},
	{NF, DW, NF, NF, NF, TL, NF, NF, NF, TL, NF, NF, NF, DW, NF},
	{TW, NF, NF, DL, NF, NF, NF, TW, NF, NF, NF, DL, NF, NF, TW},
}

var LetterScores = map[string]int{
	"A": 1,
	"B": 3,
	"C": 3,
	"D": 2,
	"E": 1,
	"F": 4,
	"G": 2,
	"H": 4,
	"I": 1,
	"J": 8,
	"K": 5,
	"L": 1,
	"M": 3,
	"N": 1,
	"O": 1,
	"P": 3,
	"Q": 10,
	"R": 1,
	"S": 1,
	"T": 1,
	"U": 1,
	"V": 4,
	"W": 4,
	"X": 8,
	"Y": 4,
	"Z": 10,
	"*": 0,
}

var tileDistribution = map[string]int{
	"A": 9,
	"B": 2,
	"C": 2,
	"D": 4,
	"E": 12,
	"F": 2,
	"G": 3,
	"H": 2,
	"I": 9,
	"J": 1,
	"K": 1,
	"L": 4,
	"M": 2,
	"N": 6,
	"O": 8,
	"P": 2,
	"Q": 1,
	"R": 6,
	"S": 4,
	"T": 6,
	"U": 4,
	"V": 2,
	"W": 2,
	"X": 1,
	"Y": 2,
	"Z": 1,
	"*": 2,
}

type Board struct {
	Fields [][]Field
	// Reverse map of tile positions as pairs of int
	TilePositions map[Tile][2]int
}

type BoardActions interface {
	PlaceTile(tile *Tile, x int, y int)
	RemoveTileByReference(tile *Tile)
	RemoveTileByCoordinates(x int, y int)
	GetField(x int, y int) (Field, bool)
	IsFieldEmpty(x int, y int) bool
	IsEmpty() bool
	GetTilePosition(tile *Tile) ([2]int, bool)
	IsTileOnBoard(tile *Tile) bool
	SetTilePosition(tile *Tile, x int, y int)
	UnsetTilePosition(tile *Tile)
}

func (r *Board) PlaceTile(tile *Tile, x int, y int) {
	if tile == nil {
		return
	}
	r.Fields[x][y].Tile = tile
	// Check if tile was already placed and remove it from the old position
	if oldPos, ok := r.GetTilePosition(tile); ok {
		r.RemoveTileByCoordinates(oldPos[0], oldPos[1])
	}
	r.SetTilePosition(tile, x, y)
}

func (r *Board) GetTilePosition(tile *Tile) ([2]int, bool) {
	coordinate, exists := r.TilePositions[*tile]
	return coordinate, exists
}

func (r *Board) IsTileOnBoard(tile *Tile) bool {
	_, exists := r.TilePositions[*tile]
	return exists
}

func (r *Board) SetTilePosition(tile *Tile, x int, y int) {
	r.TilePositions[*tile] = [2]int{x, y}
}

func (r *Board) UnsetTilePosition(tile *Tile) {
	delete(r.TilePositions, *tile)
}

func (r *Board) RemoveTileByReference(tile *Tile) {
	if pos, ok := r.GetTilePosition(tile); ok {
		r.RemoveTileByCoordinates(pos[0], pos[1])
	}
}

func (r *Board) RemoveTileByCoordinates(x int, y int) {
	if r.Fields[x][y].Tile == nil {
		return
	}
	r.UnsetTilePosition(r.Fields[x][y].Tile)
	r.Fields[x][y].Tile = nil
}

func (r *Board) GetField(x int, y int) (Field, bool) {
	if x < 0 || x > 14 || y < 0 || y > 14 {
		return Field{}, false
	}
	return r.Fields[x][y], true
}

func (r *Board) IsFieldEmpty(x int, y int) bool {
	// check bounds
	if x < 0 || y < 0 || x >= len(r.Fields) || y >= len(r.Fields[x]) {
		return false
	}
	return r.Fields[x][y].Tile == nil
}

func (r *Board) IsEmpty() bool {
	for i := range r.Fields {
		for j := range r.Fields[i] {
			if !r.IsFieldEmpty(i, j) {
				return false
			}
		}
	}
	return true
}

func NewBoard() *Board {
	board := &Board{
		Fields:        make([][]Field, 15),
		TilePositions: make(map[Tile][2]int),
	}
	for i := range board.Fields {
		board.Fields[i] = make([]Field, 15)
		for j := range board.Fields[i] {
			board.Fields[i][j] = *NewField(specialFields[i][j])
		}
	}
	return board
}
