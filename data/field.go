package data

import (
	"fmt"
	"math/rand"
	"sync"
)

type Field struct {
	particles []*Particle
	Size      struct {
		MinAxisX float64
		MaxAxisX float64
		MinAxisY float64
		MaxAxisY float64
	}
}

func NewField(particleCount int) *Field {
	field := &Field{}

	field.particles = make([]*Particle, particleCount)
	for i := 0; i < particleCount; i++ {
		x := field.Size.MinAxisX + rand.Float64()*(field.Size.MaxAxisX-field.Size.MinAxisX)
		y := field.Size.MinAxisY + rand.Float64()*(field.Size.MaxAxisY-field.Size.MinAxisY)
		field.particles[i] = NewParticle(field, x, y)
	}
	return field
}

func (f *Field) UpdatePosition(t float64) {
	var wg sync.WaitGroup

	for _, p := range f.particles {
		wg.Add(1)
		go func(p *Particle, t float64) {
			defer wg.Done()
			p.UpdatePosition(t)
		}(p, t)
	}
	wg.Wait()
}

func (f *Field) Print() {
	fmt.Println(f.particles)
}
