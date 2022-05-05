package pages

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/sfortson/fitness-tracker/server/calculator"
	"github.com/sfortson/fitness-tracker/server/database"
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

// func parseInt(s string) int {
// 	i, err := strconv.ParseInt(s, 10, 32)
// 	if err != nil {
// 		return 0
// 	}
// 	return int(i)
// }

type SessionToken string

var contextKeySessionToken = SessionToken("session-token")

func HomePage(w http.ResponseWriter, r *http.Request) {
	sessionToken, ok := r.Context().Value(contextKeySessionToken).(string)
	if !ok {
		log.Println("Unable to parse session token")
	}

	user, err := database.LookupUserByToken(r.Context(), sessionToken)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles("server/templates/home.html", "server/templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl := template.Must(t, err)

	// Get Age
	year2, _, _ := time.Now().Date()
	year1, _, _ := user.Birthdate.Date()
	year := math.Abs(float64(int(year2 - year1)))

	if r.Method != http.MethodPost {
		hp := homepage{
			Name: user.Username,
			FormValues: calculator.BodyFatCalculator{
				Age: int(year),
			},
		}
		tmpl.ExecuteTemplate(w, "base", hp)
		return
	}

	bf := calculator.BodyFatCalculator{
		Neck:   parseFloat(r.FormValue("neck")),
		Weight: parseFloat(r.FormValue("weight")),
		Waist:  parseFloat(r.FormValue("waist")),
		Height: parseFloat(r.FormValue("height")),
		Age:    int(year),
	}
	percentage := bf.Calculate()
	bmi := bf.CalculateBMI()
	description, healthrisk := bf.ReadIdeals(float32(percentage))

	hp := homepage{
		Name:        user.Username,
		BodyFat:     percentage,
		BMI:         bmi,
		FormValues:  bf,
		Description: description,
		HealthRisk:  healthrisk,
	}

	tmpl.ExecuteTemplate(w, "base", hp)
}
