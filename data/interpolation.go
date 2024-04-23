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

func InterpolateByT(fields []*VelocityField, timeStep, t, x, y float64) float64 {
	ct := math.Floor(t / timeStep)
	t1, t2, t3 := math.Max(ct-1, 0), ct, math.Min(ct+1, float64(len(fields))-1)
	yi := []float64{
		fields[int(t1)].GetVelocity(x, y),
		fields[int(t2)].GetVelocity(x, y),
		fields[int(t3)].GetVelocity(x, y),
	}
	xi := []float64{
		t1 * timeStep,
		t2 * timeStep,
		t3 * timeStep,
	}
	v, err := interpolation.SplineInterpolation(xi, yi, t)
	if err != nil {
		panic(err)
	}
	return v
}
