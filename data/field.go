package data

import (
	"image"
	"math"

	"github.com/pa-m/sklearn/interpolate"
)

type Field struct {
	particles []*Particle
	step      float64
	Size      Size
}

type VelocityField struct {
	Data [][]float64
	Time float64
}

type Size struct {
	MinAxisX float64
	MaxAxisX float64
	MinAxisY float64
	MaxAxisY float64
}

type FieldParams struct {
	Size           Size
	ParticleCount  int
	TimeStep       float64
	InterStepCount int
}

func (f Field) GetNextIterationField(time float64) *Field {
	// var wg sync.WaitGroup

	// for _, p := range f.particles {
	// 	wg.Add(1)
	// 	go func(p *Particle, t float64) {
	// 		defer wg.Done()
	// 		x1, y1 := p.UpdatePositionRK(t)
	// 		p.X = x1
	// 		p.Y = y1
	// 	}(p, timeStep)
	// }
	// wg.Wait()
	// return &f
	for _, p := range f.particles {
		//x2, y2 := p.UpdatePositionAnalytical(timeStep)
		x1, y1 := p.UpdatePositionEuler(time)
		// if (x1 != x2) || y1 != y2 {
		// 	fmt.Println(x1, " ", x2, " Y ", y1, " ", y2)
		// }
		p.X = x1
		p.Y = y1
	}
	return &f
}

func (f *Field) Image(imageWidth, imageHeight int) image.Image {
	width := f.Size.MaxAxisX - f.Size.MinAxisX
	height := f.Size.MinAxisY - f.Size.MaxAxisY

	widthCoef := float64(imageWidth) / width
	heightCoef := float64(imageHeight) / height

	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{int(imageWidth), int(imageHeight)},
	})

	for _, particle := range f.particles {
		x := widthCoef * (particle.X - f.Size.MinAxisX)
		y := heightCoef * (particle.Y + f.Size.MinAxisY)
		img.SetRGBA(int(x), int(y), particle.Color)
	}

	return img
}

func NewVelocityField(data [][]float64, time float64) *VelocityField {
	return &VelocityField{
		Data: data,
		Time: time,
	}
}

func (vf *VelocityField) GetVelocity(x, y, t float64) float64 {
	x1, x2 := int(math.Floor(x)), min(int(math.Ceil(x)), len(vf.Data)-1)
	y1, y2 := int(math.Floor(y)), min(int(math.Ceil(y)), len(vf.Data[0])-1)
	v1, v2, v3, v4 := vf.Data[x1][y1], vf.Data[x1][y2], vf.Data[x2][y1], vf.Data[x2][y2]
	// u1, w1 := utils.VelocityPointByFraction(float64(x2), float64(y1))
	// if v3 != u1 && v3 != w1 {
	// 	fmt.Println("Diff U ", (v3-u1)/u1, " W ", (v3-w1)/w1)
	// }
	xs := []float64{
		float64(x1),
		float64(x1),
		float64(x2),
		float64(x2),
	}
	ys := []float64{
		float64(y1),
		float64(y1),
		float64(y2),
		float64(y2),
	}
	vel := interpolate.Interp2d(xs, ys, []float64{v1, v2, v3, v4})(x, y)
	// u, w := utils.VelocityPointByFraction(t, x, y)
	// if u != vel && w != vel {
	// 	fmt.Println("Vmiss U ", (vel-u)/u, " W ", (vel-w)/w)
	// }
	return vel
}
