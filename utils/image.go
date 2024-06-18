package utils

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func CombineImages(img1 image.Image, img2 image.Image) image.Image {
	rect1 := img1.Bounds()
	rect2 := img2.Bounds()
	if rect1 != rect2 {
		panic("Different image size")
	}
	img := image.NewRGBA(image.Rect(0, 0, rect1.Dx(), rect1.Dy()))
	draw.Draw(img, rect1, img1, image.Point{0, 0}, draw.Over)
	draw.Draw(img, rect1, img2, image.Point{0, 0}, draw.Over)
	return img
}

func DrawAxis(img image.Image, minX, maxX, minY, maxY float64) image.Image {
	rect := img.Bounds()
	coefX := float64(rect.Dx()) / (maxX - minX)
	coefY := float64(rect.Dy()) / (maxY - minY)
	blank := image.NewRGBA(image.Rect(0, rect.Dy(), rect.Dx(), 0))

	//draw lines
	centerX := int(math.Round((0 - minX) * coefX))
	centerY := rect.Dy() - 1 - int(math.Round((0-minY)*coefY))
	for x := 0; x < rect.Dx(); x++ {
		blank.Set(x, centerY, color.Black)
	}
	for y := 0; y < rect.Dy(); y++ {
		blank.Set(centerX, y, color.Black)
	}
	//draw marks
	step := 1.0
	if math.Abs(maxX-minX) > 10 {
		step = maxX - minX/10.0
	}
	for i := 0.0; i <= math.Abs(minX); i += step {
		x := min(centerX-int(math.Round(i*coefX)), rect.Dx()-1)
		blank.Set(x, centerY-1, color.Black)
		blank.Set(x, centerY, color.Black)
		blank.Set(x, centerY+1, color.Black)
	}
	for i := 0.0; i <= math.Abs(maxX); i += step {
		x := min(centerX+int(math.Round(i*coefX)), rect.Dx()-1)
		blank.Set(x, centerY-1, color.Black)
		blank.Set(x, centerY, color.Black)
		blank.Set(x, centerY+1, color.Black)
	}
	if math.Abs(maxY-minY) > 10 {
		step = math.Abs(maxY-minY) / 10.0
	}
	for i := 0.0; i <= math.Abs(minY); i += step {
		y := max(centerY-int(math.Round(i*coefY)), 0)
		blank.Set(centerX-1, y, color.Black)
		blank.Set(centerX, y, color.Black)
		blank.Set(centerX+1, y, color.Black)
	}
	for i := 0.0; i <= math.Abs(maxY); i += step {
		y := max(centerY+int(math.Round(i*coefY)), 0)
		blank.Set(centerX-1, y, color.Black)
		blank.Set(centerX, y, color.Black)
		blank.Set(centerX+1, y, color.Black)
	}
	return CombineImages(blank, img)
}
