package gui

import (
	"fyne.io/fyne/v2"
	"math"
)

// General geometric stuff

func XY2I(x int, y int, width int) int {
	return x + y*width
}

// P2C TODO Find better solution. This works but looks too complicated
func P2C(position fyne.Position, width uint16, height uint16) (int, int) {
	return int((uint16(math.Abs(float64(position.X)))) / width), int((uint16(math.Abs(float64(position.Y)))) / height)
}
