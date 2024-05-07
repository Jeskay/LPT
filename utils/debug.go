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

func VelocityPoint(t, x, y float64) (float64, float64) {
	r := math.Sqrt(x*x + y*y)
	Vt := math.Tanh(r) * (1 / math.Cosh(r))
	pA := math.Atan2(y, x)
	if pA < 0 {
		pA = pA + math.Pi*2
	}
	return Vt * math.Sin(pA) * (-1), Vt * math.Cos(pA)
}

func VelocityPointByFraction(t, x, y float64) (float64, float64) {
	x1, y1 := VelocityToEuclid(x, y)
	return VelocityPoint(t, x1, y1)
}

func VelocityToNormal(x, y float64) (float64, float64) {
	vl := 256
	minL := -3
	maxL := 3
	ratioX := float64(vl-1) / (float64(maxL) - float64(minL))
	ratioY := float64(vl-1) / (float64(minL) - float64(maxL))
	cX, cY := ratioX*(x-float64(minL)), ratioY*(y+float64(minL))
	return math.Ceil(cX), math.Ceil(cY)
}

func VelocityToEuclid(x, y float64) (float64, float64) {
	maxL := 3
	minL := -3
	x1 := x * float64(maxL-minL) / 256
	y1 := y * float64(minL-maxL) / 256
	return x1 + float64(minL), y1 - float64(minL)
}

func GenerateVelocity(t, maxL, minL float64) ([][]float64, [][]float64) {
	vel1 := make([][]float64, 256)
	vel2 := make([][]float64, 256)
	for i := 0; i < 256; i++ {
		vel1[i] = make([]float64, 256)
		vel2[i] = make([]float64, 256)
		for j := 0; j < 256; j++ {
			x, y := VelocityToEuclid(float64(i), float64(j))
			u, w := VelocityPoint(t, x, y)
			xn, yn := VelocityToNormal(x, y)
			if int(xn) != i || int(yn) != j {
				fmt.Println("ALARM")
			}
			fmt.Println("Converted U ", xn, " ", i, " W ", yn, " ", j)
			vel1[i][j], vel2[i][j] = u, w
		}
	}
	return vel1, vel2
}
