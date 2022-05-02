package pages

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/sfortson/fitness-tracker/server/calculator"
)

type homepage struct {
	Name        string
	BodyFat     float64
	BMI         float64
	FormValues  calculator.BodyFatCalculator
	Description string
	HealthRisk  string
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

func parseInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int(i)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("server/templates/home.html", "server/templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl := template.Must(t, err)

	if r.Method != http.MethodPost {
		hp := homepage{Name: "Sam"}
		tmpl.ExecuteTemplate(w, "base", hp)
		return
	}

	bf := calculator.BodyFatCalculator{
		Neck:   parseFloat(r.FormValue("neck")),
		Weight: parseFloat(r.FormValue("weight")),
		Waist:  parseFloat(r.FormValue("waist")),
		Height: parseFloat(r.FormValue("height")),
		Age:    parseInt(r.FormValue("age")),
	}
	percentage := bf.Calculate()
	bmi := bf.CalculateBMI()
	description, healthrisk := bf.ReadIdeals(float32(percentage))

	hp := homepage{
		Name:        "Sam",
		BodyFat:     percentage,
		BMI:         bmi,
		FormValues:  bf,
		Description: description,
		HealthRisk:  healthrisk,
	}

	tmpl.ExecuteTemplate(w, "base", hp)
}
