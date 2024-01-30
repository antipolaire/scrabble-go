package game

import (
	"math/rand"
	"time"
)

type Bag struct {
	Tiles []Tile
}

type BagActions interface {
	// TakeTiles takes a number of random tiles from the bag and returns them. If the bag is empty, it returns an empty slice
	// of tiles. If count exceeds the number of tiles in the bag, it returns all tiles in the bag.
	TakeTiles(count int) []Tile
}

func (bag *Bag) TakeTiles(count int) []Tile {
	rand.Seed(time.Now().UnixNano())
	if count > len(bag.Tiles) {
		count = len(bag.Tiles)
	}
	tiles := make([]Tile, count)
	// take random tiles from the bag
	for i := 0; i < count; i++ {
		randomIndex := rand.Intn(len(bag.Tiles))
		tiles[i] = bag.Tiles[randomIndex]
		// remove the tile from the bag
		bag.Tiles = append(bag.Tiles[:randomIndex], bag.Tiles[randomIndex+1:]...)
	}
	return tiles
}

func NewBag() *Bag {
	bag := &Bag{
		Tiles: make([]Tile, 0),
	}
	for letter, count := range tileDistribution {
		for i := 0; i < count; i++ {
			bag.Tiles = append(bag.Tiles, *NewTile(letter, LetterScores[letter]))
		}
	}
	return bag
}
