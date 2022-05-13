package homepage

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sfortson/fitness-tracker/internal/calculator"
	"github.com/sfortson/fitness-tracker/internal/database"
	"github.com/sfortson/fitness-tracker/internal/helpers"
	"github.com/sfortson/fitness-tracker/internal/session"
	templates "github.com/sfortson/fitness-tracker/web/app"
)

type homeForm struct {
	Neck   float64
	Weight float64
	Waist  float64
	Errors map[string]string
}

type homepage struct {
	Name        string
	BodyFat     float64
	BMI         float64
	FormValues  homeForm
	Description string
	HealthRisk  string
	Data        []database.BodyFat
	DateList    []string
	BFList      []float64
}

func getUser(r *http.Request) (*database.User, error) {
	sessionToken, ok := r.Context().Value(session.ContextKeySessionToken).(string)
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

func (form *homeForm) validateHomeForm() bool {
	form.Errors = make(map[string]string)

	if form.Neck <= 0 {
		form.Errors["Neck"] = "Invalid neck size"
	}

	if form.Waist <= 0 {
		form.Errors["Waist"] = "Invalid waist size"
	}

	if form.Weight <= 0 {
		form.Errors["Weight"] = "Invalid weight"
	}

	return len(form.Errors) == 0
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		log.Fatal(err)
	}

	tmpl := templates.WebTemplates["homepage"]

	dateList, bfList := dataLists(user)

	hp := homepage{
		Name:     user.Username,
		Data:     user.BodyFatMeasurements,
		DateList: dateList,
		BFList:   bfList,
	}

	tmpl.ExecuteTemplate(w, "base", hp)
}

func PostHome(w http.ResponseWriter, r *http.Request) {
	refer := strings.Split(r.Header.Get("Referer"), "/")
	if refer[len(refer)-1] == "login" || refer[len(refer)-1] == "register" {
		GetHome(w, r)
		return
	}

	user, err := getUser(r)
	if err != nil {
		log.Fatal(err)
	}

	tmpl := templates.WebTemplates["homepage"]

	// Get today's date
	year2, month, day := time.Now().Date()

	neck := helpers.ParseFloat(r.FormValue("neck"))
	weight := helpers.ParseFloat(r.FormValue("weight"))
	waist := helpers.ParseFloat(r.FormValue("waist"))

	bf := calculator.BodyFatCalculator{
		Neck:   neck,
		Weight: weight,
		Waist:  waist,
		Height: user.Height,
		Age:    user.GetAge(),
	}
	percentage := bf.Calculate()
	bmi := bf.CalculateBMI()
	description, healthrisk := bf.ReadIdeals(float32(percentage))

	hf := &homeForm{
		Neck:   neck,
		Waist:  waist,
		Weight: weight,
	}

	dateList, bfList := dataLists(user)

	hp := homepage{
		Name:        user.Username,
		BodyFat:     percentage,
		BMI:         bmi,
		FormValues:  *hf,
		Description: description,
		HealthRisk:  healthrisk,
		Data:        user.BodyFatMeasurements,
		DateList:    dateList,
		BFList:      bfList,
	}

	if !hf.validateHomeForm() {
		hp.FormValues = *hf
		tmpl.ExecuteTemplate(w, "base", hp)
		return
	}

	bodyFat := database.BodyFat{
		UserID:     user.ID,
		Neck:       neck,
		Weight:     weight,
		Waist:      waist,
		Height:     user.Height,
		Year:       year2,
		Month:      month,
		Day:        day,
		Percentage: percentage,
		BMI:        bmi,
	}

	addBodyFatToDatabase(bodyFat)

	tmpl.ExecuteTemplate(w, "base", hp)
}
