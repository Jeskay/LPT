package data

import (
	"image"
	"image/color"

	"github.com/nfnt/resize"
)

var (
	colorPositive = [3]float64{255, 0, 0}
	colorNegative = [3]float64{0, 255, 0}
)

func gradient(value, max float64) color.Color {
	d := value / max
	r := colorPositive[0] + (colorNegative[0]-colorPositive[0])*d
	g := colorPositive[1] + (colorNegative[1]-colorPositive[1])*d
	b := colorPositive[2] + (colorNegative[2]-colorNegative[2])*d
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}

func NewGradientFromField(field *VelocityField, width, height uint) image.Image {
	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{256, 256},
	})
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			g := gradient(field.Data[x][y]+3, 6)
			img.Set(x, y, g)
		}
	}
	return resize.Resize(width, height, img, resize.Lanczos2)
}
