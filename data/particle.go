package data

import (
	"image/color"
	"math"
	"sync"
)

type Particle struct {
	Color        color.RGBA
	X            float64
	Y            float64
	fieldManager *FieldManager
	fieldLock    sync.RWMutex
}

func NewParticle(field *FieldManager, x float64, y float64) *Particle {
	return &Particle{
		Color:        color.RGBA{100, 200, 200, 0xff},
		X:            x,
		Y:            y,
		fieldManager: field,
	}
}

func (p *Particle) UpdatePositionAnalytical(t float64) (float64, float64) {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	w := math.Tanh(r) * (1 / math.Cosh(r))
	polarAngle := math.Atan2(p.Y, p.X)
	if polarAngle < 0 {
		polarAngle = polarAngle + math.Pi*2
	}
	newX := r * math.Cos(w*t+polarAngle)
	p.fieldLock.RLock()
	if newX > p.fieldManager.GetSize().MaxAxisX {
		newX = newX - p.fieldManager.GetSize().MaxAxisX
	}
	if newX < p.fieldManager.GetSize().MinAxisX {
		newX = newX - p.fieldManager.GetSize().MinAxisX + p.fieldManager.GetSize().MaxAxisX
	}
	p.fieldLock.RUnlock()
	newY := r * math.Sin(w*t+polarAngle)
	p.fieldLock.RLock()
	if newY > p.fieldManager.GetSize().MaxAxisY {
		newY = newY - p.fieldManager.GetSize().MaxAxisY
	}
	if p.Y < p.fieldManager.GetSize().MinAxisY {
		newY = newY - p.fieldManager.GetSize().MinAxisY + p.fieldManager.GetSize().MaxAxisY
	}
	p.fieldLock.RUnlock()
	return newX, newY
}

func (p *Particle) UpdatePositionRK(t float64) (float64, float64) {
	p.fieldLock.RLock()
	size := p.fieldManager.GetSize()
	h := p.fieldManager.GetTimeStep()
	u1, v1 := p.fieldManager.GetVelocity(p.X, p.Y, t)
	x1, y1 := p.X+0.5*h*u1, p.Y+0.5*h*v1
	u2, v2 := p.fieldManager.GetVelocity(x1, y1, t+0.5*h)
	x2, y2 := p.X+0.5*h*u2, p.Y+0.5*h*v2
	u3, v3 := p.fieldManager.GetVelocity(x2, y2, t+0.5*h)
	x3, y3 := p.X+h*u3, p.Y+h*v3
	u4, v4 := p.fieldManager.GetVelocity(x3, y3, t+h)
	dX := (u1 + 2*u2 + 2*u3 + u4) * h / 6
	newX := p.X + dX
	if newX > size.MaxAxisX {
		newX = size.MinAxisX + (newX - size.MaxAxisX)
	}
	if newX < size.MinAxisX {
		newX = size.MaxAxisX - (newX - size.MinAxisX)
	}
	dY := (v1 + 2*v2 + 2*v3 + v4) * h / 6
	newY := p.Y + dY
	if newY > size.MaxAxisY {
		newY = size.MinAxisY + (newY - size.MaxAxisY)
	}
	if newY < size.MinAxisY {
		p.Y = size.MaxAxisY - (newY - size.MinAxisY)
	}
	p.fieldLock.RUnlock()
	return newX, newY
}

func (p *Particle) UpdatePositionEuler(t float64) (float64, float64) {
	u, v := p.fieldManager.GetVelocity(p.X, p.Y, t)
	newX, newY := p.X+u*p.fieldManager.GetTimeStep()*float64(p.fieldManager.GetInterStepCount()), p.Y+v*p.fieldManager.GetTimeStep()*float64(p.fieldManager.GetInterStepCount())
	return newX, newY
}
