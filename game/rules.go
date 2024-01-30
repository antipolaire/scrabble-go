package game

import "strings"

// This represents rules and checkers for the game of scrabble

type MoveValidationResult struct {
	IsValid       bool
	MovedOnBoard  bool
	MovedOffBoard bool
}

type RuleChecker interface {
	ValidateMove(x int, y int, tile *Tile) MoveValidationResult
}

func CalculateWordScore(word string) int {
	score := 0
	for _, letter := range word {
		score += LetterScores[strings.ToUpper(string(letter))]
	}
	return score
}

// ValidateMove checks if the move is valid. A valid move is as follows:
// - Tile is not yet on the board, coordinates are within bounds and the target field is empty => Tile is placed on the board
// - Tile is already on the board, coordinates are within bounds and the target field is empty => Tile is moved on the board
// - Tile is already on the board, coordinates are out of bounds => Tile is removed from the board
// All other moves are invalid
func (game *Game) ValidateMove(x int, y int, tile *Tile) MoveValidationResult {
	newCoordinatesWithinBounds := coordinatesWithinBounds(x, y)
	isTileOnBoard := game.Board.IsTileOnBoard(tile)
	isFieldEmpty := game.Board.IsFieldEmpty(x, y)
	if !isTileOnBoard && newCoordinatesWithinBounds && isFieldEmpty {
		return MoveValidationResult{
			IsValid:       true,
			MovedOnBoard:  true,
			MovedOffBoard: false,
		}
	} else if isTileOnBoard && newCoordinatesWithinBounds && isFieldEmpty {
		return MoveValidationResult{
			IsValid:       true,
			MovedOnBoard:  true,
			MovedOffBoard: false,
		}
	} else if isTileOnBoard && !newCoordinatesWithinBounds {
		return MoveValidationResult{
			IsValid:       true,
			MovedOnBoard:  false,
			MovedOffBoard: true,
		}
	}
	return MoveValidationResult{
		IsValid:       false,
		MovedOnBoard:  false,
		MovedOffBoard: false,
	}
}

func coordinatesWithinBounds(x int, y int) bool {
	return !(x < 0 || x > 14 || y < 0 || y > 14)
}
