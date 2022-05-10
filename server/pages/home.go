package pages

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
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
	DateList    []string
	BFList      []float64
}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}

type SessionToken string

var contextKeySessionToken = SessionToken("session-token")

func getAge(user *database.User) int {
	// Get Age
	year2, _, _ := time.Now().Date()
	year1, _, _ := user.Birthdate.Date()
	year := int(math.Abs(float64(int(year2 - year1))))
	return year
}

func getUser(r *http.Request) (*database.User, error) {
	sessionToken, ok := r.Context().Value(contextKeySessionToken).(string)
	if !ok {
		log.Println("Unable to parse session token")
	}

	user, err := database.LookupUserByToken(r.Context(), sessionToken)

	return user, err
}

func dataLists(user *database.User) ([]string, []float64) {
	database.DB.Preload("BodyFatMeasurements").Find(&user)

	var dateList []string
	var bfList []float64
	for _, d := range user.BodyFatMeasurements {
		dateList = append(dateList, strings.Split(d.CreatedAt.String(), " ")[0])
		bfList = append(bfList, d.Percentage)
	}
	return dateList, bfList
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles("server/templates/home.html", "server/templates/base.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(t, err)

	dateList, bfList := dataLists(user)

	hp := homepage{
		Name: user.Username,
		FormValues: calculator.BodyFatCalculator{
			Age: getAge(user),
		},
		Data:     user.BodyFatMeasurements,
		DateList: dateList,
		BFList:   bfList,
	}

	tmpl.ExecuteTemplate(w, "base", hp)
}

func addBodyFatToDatabase(bodyFat database.BodyFat) {
	var foundBodyFat database.BodyFat
	result := database.DB.Where(&database.BodyFat{
		UserID: bodyFat.UserID,
		Year:   bodyFat.Year,
		Month:  bodyFat.Month,
		Day:    bodyFat.Day,
	}).First(&foundBodyFat)

	if result.Error != nil {
		database.DB.Create(&bodyFat)
	} else {
		foundBodyFat.BMI = bodyFat.BMI
		foundBodyFat.Percentage = bodyFat.Percentage
		foundBodyFat.Neck = bodyFat.Neck
		foundBodyFat.Waist = bodyFat.Waist
		foundBodyFat.Weight = bodyFat.Waist
		database.DB.Save(&foundBodyFat)
	}
}

func PostHome(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
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
		UserID:     user.ID,
		Neck:       neck,
		Weight:     weight,
		Waist:      waist,
		Height:     height,
		Year:       year2,
		Month:      month,
		Day:        day,
		Percentage: percentage,
		BMI:        bmi,
	}

	if (neck != 0 || waist != 0 || weight != 0 || height != 0) {
		addBodyFatToDatabase(bodyFat)
	}

	dateList, bfList := dataLists(user)

	hp := homepage{
		Name:        user.Username,
		BodyFat:     percentage,
		BMI:         bmi,
		FormValues:  bf,
		Description: description,
		HealthRisk:  healthrisk,
		Data:        user.BodyFatMeasurements,
		DateList:    dateList,
		BFList:      bfList,
	}

	tmpl.ExecuteTemplate(w, "base", hp)
}
