package data

import (
	"LPT/interpolation"
	"math"
)

func InterpolateByT(fields []*VelocityField, interStep int, timeStep, t, x, y float64) float64 {
	ct := int(math.Floor(t/timeStep)) / interStep
	var t1, t2, t3 int
	if ct < 1 {
		t1, t2, t3 = 0, 1, 2
	} else if ct == len(fields)-1 {
		t1, t2, t3 = len(fields)-3, len(fields)-2, len(fields)-1
	} else {
		t1, t2, t3 = ct-1, ct, ct+1
	}
	tu1, tu2, tu3 := float64(t1*(1+interStep)+1)*timeStep, float64(t2*(1+interStep)+1)*timeStep, float64(t3*(1+interStep)+1)*timeStep
	// u1, _ := utils.VelocityPointByFraction(tu1, x, y)
	// u2, _ := utils.VelocityPointByFraction(tu2, x, y)
	// u3, _ := utils.VelocityPointByFraction(float64(t3*(1+interStep))*timeStep, x, y)
	yi := []float64{
		fields[t1].GetVelocity(x, y, tu1),
		fields[t2].GetVelocity(x, y, tu2),
		fields[t3].GetVelocity(x, y, tu3),
		//u1,
		//u2,
		// u3,
	}
	xi := []float64{
		tu1,
		tu2,
		tu3,
		// float64(t3*(1+interStep)) * timeStep,
	}
	//ut, wt := utils.VelocityPointByFraction(t, x, y)
	v, err := interpolation.SplineInterpolation(xi, yi, t)
	// umiss, wmiss := (ut-v)/ut, (wt-v)/wt
	// if umiss != 1 {
	// 	fmt.Println("Value ", umiss, " ", wmiss)
	// }
	// if math.Abs(umiss) > 0.05 && math.Abs(wmiss) > 0.05 {
	// 	fmt.Println("ALERT U MISS ", umiss, " ", "ALERT W MISS ", wmiss)
	// 	fmt.Println("precise ", u1, " ", u2)
	// }
	// fmt.Println("u ", ut, " w ", wt, " value ", v, " MIss ", (ut-v)/ut)
	if err != nil {
		panic(err)
	}
	return v
}
