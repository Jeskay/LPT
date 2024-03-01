package data

import (
	"math"
	"sync"
)

type Particle struct {
	Color     int
	X         float64
	Y         float64
	field     *Field
	fieldLock sync.RWMutex
}

func NewParticle(field *Field, x float64, y float64) *Particle {
	return &Particle{
		Color: 1,
		X:     x,
		Y:     y,
		field: field,
	}
}

func (p *Particle) UpdatePosition(t float64) {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	w := math.Tanh(r) * (1 / math.Cosh(r))
	polarAngle := math.Atan2(p.Y, p.X)
	if polarAngle < 0 {
		polarAngle = polarAngle + math.Pi*2
	}
	p.X = r * math.Cos(w*t+polarAngle)
	p.fieldLock.RLock()
	if p.X > p.field.Size.MaxAxisX {
		p.X = p.X - p.field.Size.MaxAxisX
	}
	if p.X < p.field.Size.MinAxisX {
		p.X = p.X - p.field.Size.MinAxisX + p.field.Size.MaxAxisX
	}
	p.fieldLock.RUnlock()
	p.Y = r * math.Sin(w*t+polarAngle)
	p.fieldLock.RLock()
	if p.Y > p.field.Size.MaxAxisY {
		p.Y = p.Y - p.field.Size.MaxAxisY
	}
	if p.Y < p.field.Size.MinAxisY {
		p.Y = p.Y - p.field.Size.MinAxisY + p.field.Size.MaxAxisY
	}
	p.fieldLock.RUnlock()
}
