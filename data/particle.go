package data

import (
	"image/color"
	"math"
	"sync"
)

type Particle struct {
	Color     color.RGBA
	X         float64
	Y         float64
	field     *Field
	fieldLock sync.RWMutex
}

func NewParticle(field *Field, x float64, y float64) *Particle {
	return &Particle{
		Color: color.RGBA{100, 200, 200, 0xff},
		X:     x,
		Y:     y,
		field: field,
	}
}

func (p *Particle) UpdatePositionAnalytical(t float64) {
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

func (p *Particle) UpdatePositionRK(t float64) {
	p.fieldLock.RLock()
	h := p.field.step
	u1, v1 := p.field.GetVelocity(p.X, p.Y, t)
	x1, y1 := p.X+0.5*h*u1, p.Y+0.5*h*v1
	u2, v2 := p.field.GetVelocity(x1, y1, t+0.5*h)
	x2, y2 := p.X+0.5*h*u2, p.Y+0.5*h*v2
	u3, v3 := p.field.GetVelocity(x2, y2, t+0.5*h)
	x3, y3 := p.X+h*u3, p.Y+h*v3
	u4, v4 := p.field.GetVelocity(x3, y3, t+0.5*h)
	p.X += (u1 + 2*u2 + 2*u3 + u4) / 6 * h
	p.Y += (v1 + 2*v2 + 2*v3 + v4) / 6 * h
	p.fieldLock.RUnlock()
}
