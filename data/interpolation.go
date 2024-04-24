package data

import (
	"LPT/interpolation"
	"math"
)

type VelocityTimeField struct {
	Data        [][]*VelocityTime
	TimeMoments int
}

type VelocityTime struct {
	Data []float64
}

func (timeF *VelocityTimeField) InterpolateT(time, step float64) *VelocityField {
	timeData := make([]float64, timeF.TimeMoments)
	currentT := 0.0
	for i := range timeData {
		timeData[i] = currentT
		currentT += step
	}
	data := make([][]float64, len(timeF.Data))
	for i := range timeF.Data {
		data[i] = make([]float64, len(timeF.Data[i]))
	}
	for i := range timeF.Data {
		for j := range timeF.Data[i] {
			v, err := interpolation.SplineInterpolation(timeData, timeF.Data[i][j].Data, time)
			if err != nil {
				panic(err)
			}
			data[i][j] = v
		}
	}
	return NewVelocityField(data, time)
}

func InterpolateByT(fields []*VelocityField, interStep int, timeStep, t, x, y float64) float64 {
	ct := int(math.Floor(t/timeStep)) / interStep
	t1, t2, t3 := max(ct, 0), min(ct+1, len(fields)-1), min(ct+2, len(fields)-1)
	yi := []float64{
		fields[t1].GetVelocity(x, y),
		fields[t2].GetVelocity(x, y),
		fields[t3].GetVelocity(x, y),
	}
	xi := []float64{
		float64(t1) * timeStep,
		float64(t2) * timeStep,
		float64(t3) * timeStep,
	}
	v, err := interpolation.SplineInterpolation(xi, yi, t)
	if err != nil {
		panic(err)
	}
	return v
}
