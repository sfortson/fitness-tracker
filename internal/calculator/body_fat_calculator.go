package calculator

import (
	"io/ioutil"
	"log"
	"math"

	pb "github.com/sfortson/fitness-tracker/internal/calculator/proto"
	"google.golang.org/protobuf/proto"
)

type BodyFatCalculator struct {
	Neck        float64
	Waist       float64
	Weight      float64
	Height      float64
	Age         int
	Description string
	HealthRisk  string
}

type Ideal struct {
	Descripton string
	HealthRisk string
	CurrentMax float32
	CurrentMin float32
}

func (bf BodyFatCalculator) Calculate() float64 {
	// 86.010×log10(abdomen-neck) - 70.041×log10(height) + 36.76
	p := (86.010 * math.Log10(bf.Waist-bf.Neck)) - (70.041 * math.Log10(bf.Height)) + 36.76
	return math.Trunc(p*100) / 100
}

func (bf BodyFatCalculator) CalculateBMI() float64 {
	// US units: BMI = (weight (lb) ÷ height2 (in)) * 703
	bmi := (bf.Weight / math.Pow(bf.Height, 2)) * 703
	return math.Trunc(bmi*100) / 100
}

func (bf BodyFatCalculator) ReadIdeals(body_fat_percentage float32) (Ideal) {
	content, err := 
		ioutil.ReadFile("/Users/sfortson/github-projects/fitness-tracker/internal/calculator/test.proto")
	if err != nil {
		log.Fatalln("Failed to read proto:", err)
	}

	var bumps pb.Bumps
	proto.Unmarshal(content, &bumps)

	ideal := Ideal{
		Descripton: "",
		HealthRisk: "",
	}

	for _, bump := range bumps.Bump {
		if bf.Age >= int(*bump.Age.Min) && bf.Age <= int(*bump.Age.Max) {
			for _, bf := range bump.BodyFatPercentage {
				if body_fat_percentage >= float32(*bf.Min) && body_fat_percentage <= float32(*bf.Max) {
					ideal.Descripton = *bf.Description
					ideal.HealthRisk = *bf.HealthRisk
					ideal.CurrentMax = *bf.Max
					ideal.CurrentMin = *bf.Min
				}
			}
		}
	}

	return ideal
}
