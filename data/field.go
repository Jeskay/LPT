package data

import (
	"errors"
	"image"
	"math"
	"math/rand"
	"sync"

	"github.com/pa-m/sklearn/interpolate"
)

type FieldManager struct {
	field           *Field
	images          []image.Image
	uFields         map[float64]*VelocityField
	wFields         map[float64]*VelocityField
	uStorage        *VelocityStorage
	wStorage        *VelocityStorage
	imageIndex      int
	velocityLen     int
	timeStep        float64
	VelocityRecords int
}
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
	Size          Size
	ParticleCount int
	TimeStep      float64
}

func NewFieldManager(params FieldParams, uFields []*VelocityField, wFields []*VelocityField) (*FieldManager, error) {
	if len(uFields) != len(wFields) {
		return nil, errors.New("amount of velocity field files must be the same")
	}
	if len(uFields) == 0 || len(uFields[0].Data) == 0 {
		return nil, errors.New("velocity fields must not be empty")
	}
	if len(uFields[0].Data) != len(wFields[0].Data) {
		return nil, errors.New("velocity fields must be the same length")
	}
	fManager := &FieldManager{
		uFields:         make(map[float64]*VelocityField),
		wFields:         make(map[float64]*VelocityField),
		velocityLen:     len(uFields[0].Data),
		timeStep:        params.TimeStep,
		VelocityRecords: len(uFields),
		images:          make([]image.Image, len(uFields)),
	}
	fManager.wStorage = NewVelocityStorage(NewVelocityTimeField(wFields), fManager.timeStep*0.5, float64(fManager.VelocityRecords*2))
	fManager.uStorage = NewVelocityStorage(NewVelocityTimeField(uFields), fManager.timeStep*0.5, float64(fManager.VelocityRecords*2))
	fManager.field = fManager.NewRandomField(params.ParticleCount, params.Size, params.TimeStep)
	fManager.images[0] = fManager.field.Image(1080, 720)
	fManager.imageIndex = 0
	var t float64 = params.TimeStep
	for i := 1; i < len(uFields); i++ {
		fManager.uFields[t] = uFields[i]
		fManager.wFields[t] = wFields[i]
		t += params.TimeStep
	}
	return fManager, nil
}

func (fm *FieldManager) GetSize() *Size       { return &fm.field.Size }
func (fm *FieldManager) GetTimeStep() float64 { return fm.timeStep }
func (fm *FieldManager) GetVelocityLen() int  { return fm.velocityLen }

func (fm *FieldManager) GetImage(index, w, h int) image.Image {
	if index > fm.imageIndex {
		for fm.imageIndex < index {
			nxtField := fm.field.GetNextIterationField(fm.timeStep)
			fm.imageIndex++
			fm.images[fm.imageIndex] = nxtField.Image(w, h)
		}
	}
	return fm.images[index]
}

func (fm *FieldManager) GetVelocity(x, y, t float64) (float64, float64) {
	if x > fm.GetSize().MaxAxisX {
		x = fm.GetSize().MinAxisX + (x - fm.GetSize().MaxAxisX)
	}
	if x < fm.GetSize().MinAxisX {
		x = fm.GetSize().MaxAxisX - (x - fm.GetSize().MinAxisX)
	}
	if y > fm.GetSize().MaxAxisY {
		y = fm.GetSize().MinAxisY + (y - fm.GetSize().MaxAxisY)
	}
	if y < fm.GetSize().MinAxisY {
		y = fm.GetSize().MaxAxisY - (y - fm.GetSize().MinAxisY)
	}
	ratioX, rationY := float64(fm.GetVelocityLen()-1)/(fm.GetSize().MaxAxisX-fm.GetSize().MinAxisX), float64(fm.GetVelocityLen()-1)/(fm.GetSize().MaxAxisY-fm.GetSize().MinAxisY)
	cX, cY := ratioX*(x-fm.GetSize().MinAxisX), rationY*(y-fm.GetSize().MinAxisY)

	u := fm.uStorage.fields[t].GetVelocity(cX, cY)
	w := fm.wStorage.fields[t].GetVelocity(cX, cY)
	return u, w
}

func (fm *FieldManager) NewRandomField(particleCount int, size Size, step float64) *Field {
	field := &Field{Size: size, step: step}

	field.particles = make([]*Particle, particleCount)
	for i := 0; i < particleCount; i++ {
		x := field.Size.MinAxisX + rand.Float64()*(field.Size.MaxAxisX-field.Size.MinAxisX)
		y := field.Size.MinAxisY + rand.Float64()*(field.Size.MaxAxisY-field.Size.MinAxisY)
		field.particles[i] = NewParticle(fm, x, y)
	}
	return field
}

func (f Field) GetNextIterationField(timeStep float64) *Field {
	var wg sync.WaitGroup

	for _, p := range f.particles {
		wg.Add(1)
		go func(p *Particle, t float64) {
			defer wg.Done()
			p.UpdatePositionRK(t)
		}(p, timeStep)
	}
	wg.Wait()
	return &f
	// for _, p := range f.particles {
	// 	p.UpdatePositionRK(timeStep)
	// }
	// return &f
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

func NewVelocityField(data [][]float64, time float64) *VelocityField {
	return &VelocityField{
		Data: data,
		Time: time,
	}
}

func (vf *VelocityField) GetVelocity(x, y float64) float64 {
	x1, x2 := int(math.Floor(x)), min(int(math.Ceil(x)), len(vf.Data)-1)
	y1, y2 := int(math.Floor(y)), min(int(math.Ceil(y)), len(vf.Data[0])-1)
	v1, v2, v3, v4 := vf.Data[x1][y1], vf.Data[x1][y2], vf.Data[x2][y1], vf.Data[x2][y2]
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
	return interpolate.Interp2d(xs, ys, []float64{v1, v2, v3, v4})(x, y)
}
