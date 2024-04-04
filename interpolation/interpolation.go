package interpolation

import (
	"math"
	"slices"

	"github.com/cnkei/gospline"
)

type Kernel struct {
	data [][]float64
}

func NewKernel(data [][]float64) *Kernel {
	return &Kernel{data: data}
}

func (k *Kernel) BilinearInterpolation(x, y float64) (float64, error) {
	width, height := len(k.data[0]), len(k.data)

	i, j := int(math.Floor(x)), int(math.Floor(y))
	t, u := x-float64(i), y-float64(j)

	if i == width-1 {
		i--
	}
	if j == height-1 {
		j--
	}
	return (k.data[j][i]*(1-t)+k.data[j][i+1]*t)*(1-u) + (k.data[j+1][i]*(1-t)+k.data[j+1][i+1]*t)*u, nil
}

func SplineInterpolation(xs, ys []float64, x float64) (float64, error) {
	slices.Sort(xs)
	s := gospline.NewCubicSpline(xs, ys)
	return s.At(x), nil
}
