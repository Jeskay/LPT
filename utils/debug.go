package utils

import (
	"fmt"
	"image"
	"image/draw"
	"math"
)

func FastCompare(image1, image2 image.Image) (int64, error) {
	b1, b2 := image1.Bounds(), image2.Bounds()
	img1, img2 := image.NewRGBA(image.Rect(0, 0, b1.Dx(), b1.Dy())), image.NewRGBA(image.Rect(0, 0, b2.Dx(), b2.Dy()))
	draw.Draw(img1, img1.Bounds(), image1, b1.Min, draw.Src)
	draw.Draw(img2, img2.Bounds(), image2, b2.Min, draw.Src)
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)

	for i := 0; i < len(img1.Pix); i++ {
		accumError += int64(sqDiffUInt8(img1.Pix[i], img2.Pix[i]))
	}

	return int64(math.Sqrt(float64(accumError))), nil
}

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}
