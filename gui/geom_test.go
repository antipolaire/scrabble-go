package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestGeom(t *testing.T) {
	// Test the ParseFlags function
	t.Run("Test Index Stuff", func(t *testing.T) {
		for x := 0; x < 10; x++ {
			for y := 0; y < 5; y++ {
				index := XY2I(x, y, 10)
				assert.Equal(t, x+y*10, index)
			}
		}
	})
	type PosCoords struct {
		pos fyne.Position
		X   int
		Y   int
	}

	t.Run("Pos to", func(t *testing.T) {
		expectedWithInput := []PosCoords{
			{
				pos: fyne.Position{
					X: 60, Y: 60,
				},
				X: 1,
				Y: 1,
			},
			{
				pos: fyne.Position{
					X: 61, Y: 59,
				},
				X: 1,
				Y: 0,
			},
			{
				pos: fyne.Position{
					X: 59, Y: 90,
				},
				X: 0,
				Y: 1,
			},
			{
				pos: fyne.Position{
					X: 0, Y: 0,
				},
				X: 0,
				Y: 0,
			},
			{
				pos: fyne.Position{
					X: -10, Y: 600,
				},
				X: 0,
				Y: 10,
			},
		}
		for _, expect := range expectedWithInput {
			x, y := P2C(expect.pos, 60, 60)
			assert.Equal(t, x, expect.X, fmt.Sprintf("Got %d but wanted %d from %f", x, expect.X, expect.pos.X))
			assert.Equal(t, y, expect.Y, fmt.Sprintf("Got %d but wanted %d from %f", y, expect.Y, expect.pos.Y))
		}

	})
}
