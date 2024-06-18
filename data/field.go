package data

import (
	"LPT/utils"
	"image"
	"image/color"
	"math"
	"math/rand"
	"sync"

	"github.com/pa-m/sklearn/interpolate"
)

type Field struct {
	particles []*Particle
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
	TimeStep       float64
	InterStepCount int
}

func NewEmptyField(size Size) *Field {
	return &Field{Size: size, particles: make([]*Particle, 0)}
}

func NewRandomField(particleCount int, size Size) *Field {
	field := &Field{Size: size}

	field.particles = make([]*Particle, particleCount)
	for i := 0; i < particleCount; i++ {
		x := field.Size.MinAxisX + rand.Float64()*(field.Size.MaxAxisX-field.Size.MinAxisX)
		y := field.Size.MinAxisY + rand.Float64()*(field.Size.MaxAxisY-field.Size.MinAxisY)
		field.particles[i] = NewParticle(x, y)
	}
	return field
}

func NewLinearField(particleCount int, size Size) *Field {
	field := &Field{Size: size}

	field.particles = make([]*Particle, particleCount)
	for i := 0; i < particleCount; i++ {
		x := field.Size.MinAxisX + rand.Float64()*(field.Size.MaxAxisX-field.Size.MinAxisX)
		y := 0.0
		field.particles[i] = NewParticle(x, y)
	}
	return field
}

func (f Field) GetNextIterationField(time float64, manager *FieldManager) *Field {
	var wg sync.WaitGroup
	for _, p := range f.particles {
		wg.Add(1)
		go func(p *Particle, t float64) {
			defer wg.Done()
			x1, y1 := p.UpdatePositionRK(t, manager)
			p.X = x1
			p.Y = y1
		}(p, time)
	}
	wg.Wait()
	return &f
	// for _, p := range f.particles {
	// 	//x2, y2 := p.UpdatePositionAnalytical(timeStep)
	// 	x1, y1 := p.UpdatePositionEuler(time, manager)
	// 	// if (x1 != x2) || y1 != y2 {
	// 	// 	fmt.Println(x1, " ", x2, " Y ", y1, " ", y2)
	// 	// }
	// 	p.X = x1
	// 	p.Y = y1
	// }
	// return &f
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

	return utils.DrawAxis(img, f.Size.MinAxisX, f.Size.MaxAxisX, f.Size.MinAxisY, f.Size.MaxAxisY)
}

func (f *Field) AddParticles(particleCount int, x, y, radius float64) {
	for i := 0; i < particleCount; i++ {
		r := radius / 2 * math.Sqrt(rand.Float64())
		theta := rand.Float64() * 2 * math.Pi
		pX := x + r*math.Cos(theta)
		pY := y + r*math.Sin(theta)
		f.particles = append(f.particles, NewParticle(pX, pY))
	}
}

func (f *Field) SetColor(x, y float64, color color.RGBA, radius float64) {
	for _, p := range f.particles {
		if math.Pow(x-p.X, 2)+math.Pow(y-p.Y, 2) < radius*radius {
			p.Color = color
		}
	}
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
