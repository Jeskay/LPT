package interpolation

import (
	"errors"
	"math"

	"github.com/cnkei/gospline"
	"github.com/pa-m/sklearn/interpolate"
)

func NewBilinearInterpolation(data [][]float64) func(x, y float64) float64 {
	x, y, z := make([]float64, len(data)*len(data)), make([]float64, len(data)*len(data)), make([]float64, len(data)*len(data))
	c := 0
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data); j++ {
			x[c] = float64(i)
			y[c] = float64(j)
			z[c] = data[i][j]
			c++
		}
	}
	return interpolate.Interp2d(x, y, z)
}

func SplineInterpolation(xs, ys []float64, x float64) (float64, error) {
	s := gospline.NewCubicSpline(xs, ys)
	v := s.At(x)
	if math.IsNaN(v) {
		if xs[0] == 0 {
			return ys[0], nil
		}
		return 0, errors.New("interpolation failed")
	}
	return v, nil
}
