package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/sfortson/fitness-tracker/server/calculator"
)

type HomePage struct {
	Name    string
	BodyFat float64
	BMI float64
	FormValues calculator.BodyFatCalculator
}

func parseFloat (s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		fmt.Println(err)
		return 0.0
	}
	return f
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	t, err := template.ParseFiles("server/templates/home.html", "server/templates/base.html")
	if err != nil {
		log.Println(err)
	}
	tmpl := template.Must(t, err)

	if r.Method != http.MethodPost {
		hp := HomePage{Name: "Sam", BodyFat: 0.0}
		tmpl.ExecuteTemplate(w, "base", hp)
		return
	}

	bf := calculator.BodyFatCalculator{Neck: parseFloat(r.FormValue("neck")), Weight: 248, Waist: parseFloat(r.FormValue("waist")), Height: parseFloat(r.FormValue("height")), Age: 37}
	percentage := bf.Calculate()
	bmi := bf.CalculateBMI()

	hp := HomePage{Name: "Sam", BodyFat: percentage, BMI: bmi, FormValues: bf}
	tmpl.ExecuteTemplate(w, "base", hp)
}

func main() {
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8000", nil)
}
