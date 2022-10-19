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
	MinList     []float32
	MaxList     []float32
	RegionColor string
}

func getUser(r *http.Request) (*database.User, error) {
	sessionToken, ok := r.Context().Value(session.ContextKeySessionToken).(string)
	if !ok {
		log.Println("Unable to parse session token")
	}

	user, err := database.LookupUserByToken(r.Context(), sessionToken)

	return user, err
}

func dataLists(user *database.User, ideals calculator.Ideal) ([]string, []float64, []float32, []float32) {
	database.DB.Preload("BodyFatMeasurements").Find(&user)

	var dateList []string
	var bfList []float64
	var minList []float32
	var maxList []float32

	if len(user.BodyFatMeasurements) == 0 {
		return dateList, bfList, minList, maxList
	}

	startDate := user.BodyFatMeasurements[0].CreatedAt
	endDate := time.Now()
	currentIdx := 0
	dateList = append(dateList, strings.Split(startDate.String(), " ")[0])
	bfList = append(bfList, user.BodyFatMeasurements[currentIdx].Percentage)
	minList = append(minList, ideals.CurrentMin)
	maxList = append(maxList, ideals.CurrentMax)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if currentIdx < len(user.BodyFatMeasurements)-1 {
			userTime := user.BodyFatMeasurements[currentIdx+1].CreatedAt
			if userTime.Year() <= d.Year() && userTime.YearDay() <= d.YearDay() {
				currentIdx = currentIdx + 1
			}
		}
		dateList = append(dateList, strings.Split(d.String(), " ")[0])
		bfList = append(bfList, user.BodyFatMeasurements[currentIdx].Percentage)
		minList = append(minList, ideals.CurrentMin)
		maxList = append(maxList, ideals.CurrentMax)
	}
	return dateList, bfList, minList, maxList
}

func regionColor(ideals calculator.Ideal) string {
	color := "rgba(0,40,100,0.2)"
	switch ideals.HealthRisk {
	case "Increased":
		color = "rgba(255,0,0,0.2)"
	case "Healthy":
		color = "rgba(0,255,0,0.2)"
	default:
		color = "rgba(0,40,100,0.2)"
	}
	return color
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

	bf := calculator.BodyFatCalculator{
		Age: user.GetAge(),
	}
	database.DB.Preload("BodyFatMeasurements").Find(&user)

	ideals := calculator.Ideal{
		Descripton: "",
		HealthRisk: "",
		CurrentMax: 0,
		CurrentMin: 0,
	}

	if len(user.BodyFatMeasurements) > 0 {
		ideals = bf.ReadIdeals(
			float32(user.BodyFatMeasurements[len(user.BodyFatMeasurements)-1].Percentage))
	}

	dateList, bfList, minList, maxList :=
		dataLists(user, ideals)

	hp := homepage{
		Name:        user.Username,
		Data:        user.BodyFatMeasurements,
		DateList:    dateList,
		BFList:      bfList,
		MinList:     minList,
		MaxList:     maxList,
		RegionColor: regionColor(ideals),
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
	ideals := bf.ReadIdeals(float32(percentage))

	hf := &homeForm{
		Neck:   neck,
		Waist:  waist,
		Weight: weight,
	}

	dateList, bfList, minList, maxList := dataLists(user, ideals)

	hp := homepage{
		Name:        user.Username,
		BodyFat:     percentage,
		BMI:         bmi,
		FormValues:  *hf,
		Description: ideals.Descripton,
		HealthRisk:  ideals.HealthRisk,
		Data:        user.BodyFatMeasurements,
		DateList:    dateList,
		BFList:      bfList,
		MinList:     minList,
		MaxList:     maxList,
		RegionColor: regionColor(ideals),
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
