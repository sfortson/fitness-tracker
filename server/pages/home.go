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
	Data        []database.BodyFat
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

func getAge(user *database.User) int {
	// Get Age
	year2, _, _ := time.Now().Date()
	year1, _, _ := user.Birthdate.Date()
	year := int(math.Abs(float64(int(year2 - year1))))
	return year
}

func GetHome(w http.ResponseWriter, r *http.Request) {
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

	database.DB.Preload("BodyFatMeasurements").Find(&user)

	hp := homepage{
		Name: user.Username,
		FormValues: calculator.BodyFatCalculator{
			Age: getAge(user),
		},
		Data: user.BodyFatMeasurements,
	}

	tmpl.ExecuteTemplate(w, "base", hp)
}

func PostHome(w http.ResponseWriter, r *http.Request) {
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
	year2, month, day := time.Now().Date()
	year1, _, _ := user.Birthdate.Date()
	year := math.Abs(float64(int(year2 - year1)))

	neck := parseFloat(r.FormValue("neck"))
	weight := parseFloat(r.FormValue("weight"))
	waist := parseFloat(r.FormValue("waist"))
	height := parseFloat(r.FormValue("height"))

	bf := calculator.BodyFatCalculator{
		Neck:   neck,
		Weight: weight,
		Waist:  waist,
		Height: height,
		Age:    int(year),
	}
	percentage := bf.Calculate()
	bmi := bf.CalculateBMI()
	description, healthrisk := bf.ReadIdeals(float32(percentage))

	bodyFat := database.BodyFat{
		UserID: user.ID,
		Neck:   neck,
		Weight: weight,
		Waist:  waist,
		Height: height,
		Year:   year2,
		Month:  month,
		Day:    day,
		Percentage: percentage,
		BMI: bmi,
	}

	var foundBodyFat database.BodyFat
	result := database.DB.Where(&database.BodyFat{
		UserID: user.ID,
		Year:   year2,
		Month:  month,
		Day:    day,
	}).First(&foundBodyFat)

	if result.Error != nil {
		database.DB.Create(&bodyFat)
	}
	log.Println("already added info today")

	database.DB.Preload("BodyFatMeasurements").Find(&user)

	hp := homepage{
		Name:        user.Username,
		BodyFat:     percentage,
		BMI:         bmi,
		FormValues:  bf,
		Description: description,
		HealthRisk:  healthrisk,
		Data:        user.BodyFatMeasurements,
	}

	tmpl.ExecuteTemplate(w, "base", hp)
}
