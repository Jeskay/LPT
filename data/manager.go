package data

import (
	"LPT/utils"
	"errors"
	"image"
	"math"

	"github.com/nfnt/resize"
)

type FieldManager struct {
	Field            *Field
	images           []image.Image
	VerticalFields   []*VelocityField
	HorizontalFields []*VelocityField
	imageIndex       int
	velocityLen      int
	timeStep         float64
	interStepCount   int
	VelocityRecords  int
	ImageTracing     bool
}

func NewFieldManager(params FieldParams, tracing bool, field *Field, uFields []*VelocityField, wFields []*VelocityField) (*FieldManager, error) {
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
		VerticalFields:   uFields,
		HorizontalFields: wFields,
		velocityLen:      len(uFields[0].Data),
		ImageTracing:     tracing,
		timeStep:         params.TimeStep,
		interStepCount:   params.InterStepCount,
		VelocityRecords:  len(uFields),
		images:           make([]image.Image, len(uFields)*params.InterStepCount),
		Field:            field,
	}
	fManager.imageIndex = 0
	return fManager, nil
}

func (fm *FieldManager) GetSize() *Size         { return &fm.Field.Size }
func (fm *FieldManager) GetTimeStep() float64   { return fm.timeStep }
func (fm *FieldManager) GetInterStepCount() int { return fm.interStepCount }
func (fm *FieldManager) GetVelocityLen() int    { return fm.velocityLen }

func (fm *FieldManager) GetImageById(index, w, h int) image.Image {
	if index == 0 && fm.images[index] == nil {
		fm.images[0] = resize.Resize(uint(w), uint(h), fm.Field.Image(500, 500), resize.Bilinear)
	}
	if index > fm.imageIndex {
		for fm.imageIndex < index {
			nTime := fm.timeStep * float64(fm.interStepCount) * float64(fm.imageIndex)
			nxtField := fm.Field.GetNextIterationField(nTime, fm)
			fm.imageIndex++
			img := nxtField.Image(500, 500)
			newImg := resize.Resize(uint(w), uint(h), img, resize.Bilinear)
			if fm.ImageTracing {
				fm.images[fm.imageIndex] = utils.CombineImages(fm.GetImageById(fm.imageIndex-1, w, h), newImg)
			} else {
				fm.images[fm.imageIndex] = newImg
			}

		}
	}
	return fm.images[index]
}

func (fm *FieldManager) GetCurrentFieldImage(w, h int) image.Image {
	img := fm.Field.Image(500, 500)
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
	u := InterpolateByT(fm.VerticalFields, fm.interStepCount, fm.timeStep, t, cX, cY)
	w := InterpolateByT(fm.HorizontalFields, fm.interStepCount, fm.timeStep, t, cX, cY)
	return u, w
}
