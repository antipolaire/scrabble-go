package game

import (
	"flag"
	"fmt"
	petname "github.com/dustinkirkland/golang-petname"
	"strings"
)

type Player struct {
	Name  string
	Score int
	Tiles []Tile
}

type PlayerActions interface {
	RemoveTile(tile Tile)
}

func (player *Player) RemoveTile(tile Tile) {
	for i, t := range player.Tiles {
		if t == tile {
			player.Tiles = append(player.Tiles[:i], player.Tiles[i+1:]...)
			return
		}
	}
}

func NewPlayerWithRandomName() *Player {
	flag.Parse()
	return NewPlayer(strings.ToTitle(fmt.Sprintf("%s %s", petname.Adjective(), petname.Name())))
}

func NewPlayer(name string) *Player {
	return &Player{
		Name:  name,
		Score: 0,
		Tiles: make([]Tile, 0),
	}
}
