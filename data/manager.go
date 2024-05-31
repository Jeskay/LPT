package data

import (
	"errors"
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/nfnt/resize"
)

type FieldManager struct {
	field           *Field
	images          []image.Image
	uFields         []*VelocityField
	wFields         []*VelocityField
	imageIndex      int
	velocityLen     int
	timeStep        float64
	interStepCount  int
	VelocityRecords int
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
		uFields:         uFields,
		wFields:         wFields,
		velocityLen:     len(uFields[0].Data),
		timeStep:        params.TimeStep,
		interStepCount:  params.InterStepCount,
		VelocityRecords: len(uFields),
		images:          make([]image.Image, len(uFields)*params.InterStepCount),
	}
	fManager.field = fManager.NewRandomField(params.ParticleCount, params.Size, params.TimeStep)
	fManager.images[0] = resize.Resize(uint(1080), uint(1080), fManager.field.Image(500, 500), resize.Bilinear)
	fManager.imageIndex = 0
	return fManager, nil
}

func (fm *FieldManager) GetSize() *Size         { return &fm.field.Size }
func (fm *FieldManager) GetTimeStep() float64   { return fm.timeStep }
func (fm *FieldManager) GetInterStepCount() int { return fm.interStepCount }
func (fm *FieldManager) GetVelocityLen() int    { return fm.velocityLen }
func (fm *FieldManager) SetColor(x, y float64, color color.RGBA, radius float64) {
	for _, p := range fm.field.particles {
		if math.Pow(x-p.X, 2)+math.Pow(y-p.Y, 2) < radius*radius {
			p.Color = color
		}
	}
}

func (fm *FieldManager) GetImageById(index, w, h int) image.Image {
	if index > fm.imageIndex {
		for fm.imageIndex < index {
			nxtField := fm.field.GetNextIterationField(fm.timeStep * float64(fm.interStepCount) * float64(fm.imageIndex))
			fm.imageIndex++
			img := nxtField.Image(500, 500)
			fm.images[fm.imageIndex] = resize.Resize(uint(w), uint(h), img, resize.Bilinear)
		}
	}
	return fm.images[index]
}

func (fm *FieldManager) GetCurrentFieldImage(w, h int) image.Image {
	img := fm.field.Image(500, 500)
	return resize.Resize(uint(w), uint(h), img, resize.Bilinear)
}

func (fm *FieldManager) GetVelocity(x, y, t float64) (float64, float64) {
	if x > fm.GetSize().MaxAxisX {
		x = fm.GetSize().MinAxisX + (x - fm.GetSize().MaxAxisX)
	}
	if x < fm.GetSize().MinAxisX {
		x = fm.GetSize().MaxAxisX + (x - fm.GetSize().MinAxisX)
	}
	if y > fm.GetSize().MaxAxisY {
		y = fm.GetSize().MinAxisY + (y - fm.GetSize().MaxAxisY)
	}
	if y < fm.GetSize().MinAxisY {
		y = fm.GetSize().MaxAxisY + (y - fm.GetSize().MinAxisY)
	}
	ratioX := float64(fm.GetVelocityLen()-1) / (fm.GetSize().MaxAxisX - fm.GetSize().MinAxisX)
	ratioY := float64(fm.GetVelocityLen()-1) / (fm.GetSize().MinAxisY - fm.GetSize().MaxAxisY)
	cX, cY := math.Ceil(ratioX*(x-fm.GetSize().MinAxisX)), math.Ceil(ratioY*(y+fm.GetSize().MinAxisY))
	u := InterpolateByT(fm.uFields, fm.interStepCount, fm.timeStep, t, cX, cY)
	w := InterpolateByT(fm.wFields, fm.interStepCount, fm.timeStep, t, cX, cY)
	// uft, wft := utils.VelocityPointByFraction(t, cX, cY)
	// ut, wt := utils.VelocityPoint(t, x, y)
	// umiss, wmiss := (ut-u)/ut, (wt-w)/w
	// if math.Abs(umiss) > 0.05 {
	// 	fmt.Println("ALERT U MISS ", umiss, " ", uft)
	// }
	// if math.Abs(wmiss) > 0.05 {
	// 	fmt.Println("ALERT W MISS ", wmiss, " ", wft)
	// }
	// //fmt.Println("ConvertMiss U ", (uft-ut)/ut, " W ", (wft-wt)/wt)
	// fmt.Println("Umiss ", (ut-u)/ut, " Wmiss ", (wt-w)/wt)
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

func (fm *FieldManager) NewLinearField(particleCount int, size Size, step float64) *Field {
	field := &Field{Size: size, step: step}

	field.particles = make([]*Particle, particleCount)
	for i := 0; i < particleCount; i++ {
		x := field.Size.MinAxisX + rand.Float64()*(field.Size.MaxAxisX-field.Size.MinAxisX)
		y := 0.0
		field.particles[i] = NewParticle(fm, x, y)
	}
	return field
}
