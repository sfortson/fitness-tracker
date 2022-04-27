package calculator

import "math"

type BodyFatCalculator struct {
	Neck float64
	Waist float64
	Weight float64
	Height float64
	Age int
}

func (bf BodyFatCalculator) Calculate() float64 {
	// 86.010×log10(abdomen-neck) - 70.041×log10(height) + 36.76
	p := (86.010 * math.Log10(bf.Waist - bf.Neck)) - (70.041 * math.Log10(bf.Height)) + 36.76
	return math.Trunc(p * 100) / 100
}

func (bf BodyFatCalculator) CalculateBMI() float64 {
	// US units: BMI = (weight (lb) ÷ height2 (in)) * 703
	bmi := (bf.Weight / math.Pow(bf.Height, 2)) * 703
	return math.Trunc(bmi * 100) /100
}