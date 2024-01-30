package gui

import "image/color"

// DesaturateColor Take a color and a float value and return a new color with less saturation
func DesaturateColor(c color.Color, f float64) color.Color {
	r, g, b, a := c.RGBA()
	r = uint32(float64(r) * f)
	g = uint32(float64(g) * f)
	b = uint32(float64(b) * f)
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
}
