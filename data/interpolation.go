package data

import "LPT/interpolation"

type VelocityTimeField struct {
	Data        [][]*VelocityTime
	TimeMoments int
}

type VelocityStorage struct {
	fields map[float64]*VelocityField
}

type VelocityTime struct {
	Data []float64
}

type TimeVelocity struct {
	time     float64
	velocity *VelocityField
}

type CoordinateVelocity struct {
	velocity float64
	x        int
	y        int
}

func NewVelocityTimeField(fields []*VelocityField) (res *VelocityTimeField) {
	res = &VelocityTimeField{Data: make([][]*VelocityTime, len(fields[0].Data))}
	for i := 0; i < len(fields[0].Data); i++ {
		res.Data[i] = make([]*VelocityTime, len(fields[0].Data[0]))
		for j := 0; j < len(fields[0].Data[0]); j++ {
			res.Data[i][j] = &VelocityTime{Data: make([]float64, len(fields))}
			for fi, f := range fields {
				res.Data[i][j].Data[fi] = f.Data[i][j]
			}
		}
	}
	res.TimeMoments = len(fields)
	return
}

func NewVelocityStorage(timeField *VelocityTimeField, step, times float64) *VelocityStorage {
	storage := &VelocityStorage{
		fields: make(map[float64]*VelocityField),
	}
	var t chan TimeVelocity = make(chan TimeVelocity)
	for i := 0.0; i < times; i++ {
		go func(time float64, step float64, c chan TimeVelocity) {
			v := timeField.InterpolateT(time, step)
			c <- TimeVelocity{time: time, velocity: v}
		}(i*step, step, t)
	}
	for i := 0; i < int(times); i++ {
		v := <-t
		storage.fields[v.time] = v.velocity
	}
	close(t)
	return storage
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
	var c chan CoordinateVelocity = make(chan CoordinateVelocity)
	for i := range timeF.Data {
		for j := range timeF.Data[i] {
			go func(i int, j int) {
				v, err := interpolation.SplineInterpolation(timeData, timeF.Data[i][j].Data, time)
				if err != nil {
					panic(err)
				}
				c <- CoordinateVelocity{velocity: v, x: i, y: j}
			}(i, j)
		}
	}
	for i := 0; i < len(timeF.Data)*len(timeF.Data[0]); i++ {
		v := <-c
		data[v.x][v.y] = v.velocity
	}
	close(c)
	return NewVelocityField(data, time)
}
