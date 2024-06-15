package utils

import (
	"image"
	"image/draw"
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
