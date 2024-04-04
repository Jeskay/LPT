package data

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"sync"
)

type Field struct {
	particles []*Particle
	xField    [][]float64
	yField    [][]float64
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

func NewField(particleCount int, size Size, step float64) *Field {
	field := &Field{Size: size, step: step}

	field.particles = make([]*Particle, particleCount)
	for i := 0; i < particleCount; i++ {
		x := field.Size.MinAxisX + rand.Float64()*(field.Size.MaxAxisX-field.Size.MinAxisX)
		y := field.Size.MinAxisY + rand.Float64()*(field.Size.MaxAxisY-field.Size.MinAxisY)
		field.particles[i] = NewParticle(field, x, y)
	}
	return field
}

func NewVelocityField(data [][]float64, time float64) *VelocityField {
	return &VelocityField{
		Data: data,
		Time: time,
	}
}

func (f *Field) SetVelocity(x [][]float64, y [][]float64) {
	f.xField = x
	f.yField = y
}

func (f *Field) GetVelocity(x float64, y float64, t float64) (float64, float64) {
	pLx := float64(len(f.xField))/f.Size.MaxAxisX - f.Size.MinAxisX
	pLy := float64(len(f.yField))/f.Size.MaxAxisY - f.Size.MinAxisY
	vX := int(math.Round(x * pLx))
	vY := int(math.Round(y * pLx))
	uX := int(math.Round(x * pLy))
	uY := int(math.Round(y * pLy))
	return f.xField[vX][vY], f.yField[uX][uY]
}

func (f *Field) UpdatePosition() {
	var wg sync.WaitGroup

	for _, p := range f.particles {
		wg.Add(1)
		go func(p *Particle, t float64) {
			defer wg.Done()
			p.UpdatePositionAnalytical(t)
		}(p, f.step)
	}
	wg.Wait()
}

func (f *Field) Print() {
	fmt.Println(f.particles)
}

func (f *Field) Image(imageWidth, imageHeight int) image.Image {
	width := f.Size.MaxAxisX - f.Size.MinAxisX
	height := f.Size.MaxAxisY - f.Size.MinAxisY

	widthCoef := float64(imageWidth) / width
	heightCoef := float64(imageHeight) / height

	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{int(imageWidth), int(imageHeight)},
	})

	for _, particle := range f.particles {
		x := widthCoef * (particle.X + width*0.5)
		y := heightCoef * (particle.Y + height*0.5)
		img.SetRGBA(int(x), int(y), particle.Color)
	}

	return img
}
