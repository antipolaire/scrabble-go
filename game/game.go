package game

import (
	"fmt"
	"go.uber.org/zap"
	"sort"
)

type Game struct {
	Board         *Board
	Bag           *Bag
	Players       []Player
	CurrentPlayer *Player
	// Temporary tiles move by the player
	TemporaryMoves map[*Player][]Move
	Dictionary     *Dictionary
}

type Move struct {
	X    int
	Y    int
	Tile *Tile
}

type GameActions interface {
	// AddTemporaryMove plays the temporary tiles of the player and returns the score
	AddTemporaryMove(player *Player, move Move)
	// RemoveTemporaryMove removes the temporary tiles of the player and returns the score
	RemoveTemporaryMove(player *Player, move Move)
	// ResetTemporaryMoves resets the temporary moves of the player
	ResetTemporaryMoves(player *Player)
	// PlayTemporaryMoves plays the temporary moves of the player and returns the score
	PlayTemporaryMoves(player *Player) int
	//CheckMove checks if the move of the player is valid. A move is a sequence of pairs of tiles along with coordinates
	CheckMove(player *Player, move []Move) bool
	// PullNewTilesFromBag pulls new tiles from the bag and adds them to the player's rack
	PullNewTilesFromBag(player *Player) []Tile
}

func (game *Game) PlayTemporaryMoves(player *Player) int {
	// Check if the player has temporary moves
	if game.TemporaryMoves[player] == nil {
		zap.L().Debug("Player has no temporary moves")
		return 0
	}
	// Check if the move is valid
	valid := game.CheckMove(player, game.TemporaryMoves[player])
	if !valid {
		zap.L().Debug("Cannot play temporary moves. Move is invalid")
		return 0
	}
	score := game.playMove(player, game.TemporaryMoves[player])
	zap.L().Debug(fmt.Sprintf("Player '%s' played temporary moves and scored %d points", player.Name, score))
	return score
}

func (game *Game) CheckMove(player *Player, move []Move) bool {
	zap.L().Debug(fmt.Sprintf("Checking move of player '%s'", player.Name))
	// First check which direction the move is in by comparing the first and the last in x as well as y direction.
	// In case x of the last and the first tile are the same, the move is in y direction and vice versa
	direction := "x"
	if move[0].X == move[len(move)-1].X {
		direction = "y"
	} else if move[0].Y == move[len(move)-1].Y {
		direction = "x"
	} else {
		zap.L().Debug("\tMove is invalid. Tiles are not in a straight line")
		return false
	}
	zap.L().Debug("\tMove is in direction", zap.String("direction", direction))
	// Iterate in the direction of the move from the first to the last ordinate in the direction of the move.
	// In case there is a gap between two consecutive moves, which is indicated if the difference between the current
	// and the previous ordinate is greater than 1, try to take the "missing" tile from the board. If the tile is not
	// on the board, the move is invalid.
	word := ""
	for i := 0; i < len(move); i++ {
		if direction == "x" {
			if i > 0 && move[i].X-move[i-1].X > 1 {
				if game.Board.IsFieldEmpty(move[i].X-1, move[i].Y) {
					zap.L().Debug("\tMove is invalid. There is a gap between two consecutive tiles")
					return false
				}
				word += game.Board.Fields[move[i].X-1][move[i].Y].Tile.Letter
				continue
			}
			word += move[i].Tile.Letter
		} else {
			if i > 0 && move[i].Y-move[i-1].Y > 1 {
				if game.Board.IsFieldEmpty(move[i].X, move[i].Y-1) {
					zap.L().Debug("\tMove is invalid. There is a gap between two consecutive tiles")
					return false
				}
				word += game.Board.Fields[move[i].X][move[i].Y-1].Tile.Letter
				continue
			}
			word += move[i].Tile.Letter
		}
	}
	zap.L().Debug("\tWord is", zap.String("word", word))
	// Check if the word is in the dictionary
	if !game.Dictionary.IsWord(word) {
		zap.L().Debug("\tMove is invalid. Word is not in the dictionary")
		return false
	}
	zap.L().Debug("\tMove is valid")
	return true
}

func (game *Game) playMove(player *Player, move []Move) int {
	score := 0
	// Iterate tiles of move and place them on the board
	for _, move := range move {
		game.Board.PlaceTile(move.Tile, move.X, move.Y)
		score += move.Tile.LetterScore
	}
	// Remove the tiles from the player's rack
	for _, move := range move {
		player.RemoveTile(*move.Tile)
	}
	game.ResetTemporaryMoves(player)
	player.Score += score
	return score
}

func (game *Game) PullNewTilesFromBag(player *Player) []Tile {
	numberOfCurrentTiles := len(player.Tiles)
	if numberOfCurrentTiles == 7 {
		zap.L().Debug(fmt.Sprintf("Player '%s' already has 7 tiles", player.Name))
		return nil
	}
	// Pull new tiles from the bag
	newTiles := game.Bag.TakeTiles(7 - numberOfCurrentTiles)

	// Add the new tiles to the player's rack
	player.Tiles = append(player.Tiles, newTiles...)
	zap.L().Debug(fmt.Sprintf("Player '%s' pulled %d new tiles", player.Name, len(newTiles)))
	return newTiles
}

func (game *Game) AddTemporaryMove(player *Player, move Move) {
	if game.TemporaryMoves[player] == nil {
		game.TemporaryMoves[player] = []Move{}
	}
	game.TemporaryMoves[player] = append(game.TemporaryMoves[player], move)
	zap.L().Debug(fmt.Sprintf(
		"Player '%s' placed tile '%s' on position (%d, %d)",
		player.Name,
		move.Tile.Letter,
		move.X,
		move.Y,
	))
	// Sort the temporary moves by x and y
	sortTemporaryMovesByPosition(player, game)
}

func sortTemporaryMovesByPosition(player *Player, game *Game) {
	sort.Slice(game.TemporaryMoves[player], func(i, j int) bool {
		if game.TemporaryMoves[player][i].X == game.TemporaryMoves[player][j].X {
			return game.TemporaryMoves[player][i].Y < game.TemporaryMoves[player][j].Y
		}
		return game.TemporaryMoves[player][i].X < game.TemporaryMoves[player][j].X
	})
}

func (game *Game) RemoveTemporaryMove(player *Player, move Move) {
	if game.TemporaryMoves[player] == nil {
		game.TemporaryMoves[player] = []Move{}
	}
	// Remove the move from the temporary moves
	for i, m := range game.TemporaryMoves[player] {
		if m == move {
			game.TemporaryMoves[player] = append(game.TemporaryMoves[player][:i], game.TemporaryMoves[player][i+1:]...)
			zap.L().Debug(fmt.Sprintf(
				"Player '%s' removed tile '%s' from from position (%d, %d)",
				player.Name,
				move.Tile.Letter,
				move.X,
				move.Y,
			))
		}
	}
	sortTemporaryMovesByPosition(player, game)
}

func (game *Game) ResetTemporaryMoves(player *Player) {
	game.TemporaryMoves[player] = []Move{}
}

func NewGame() *Game {
	return &Game{
		Board:          NewBoard(),
		Bag:            NewBag(),
		Players:        []Player{},
		TemporaryMoves: map[*Player][]Move{},
		Dictionary:     NewDictionaryFromDAWG("../assets/dicts/en.dawg"),
		CurrentPlayer:  NewPlayerWithRandomName(),
	}
}
