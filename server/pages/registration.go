package pages

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/sfortson/fitness-tracker/server/database"
	"github.com/sfortson/fitness-tracker/server/models"
	"golang.org/x/crypto/bcrypt"
)

type FormValues struct {
	Username  string
	Sex       string
	Birthdate time.Time
	Email     string
	Password  string
}

func getTemplate() *template.Template {
	t, err := template.ParseFiles(
		"server/templates/registration.html",
		"server/templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
	return template.Must(t, err)
}

func GetRegistration(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()
	tmpl.ExecuteTemplate(w, "base", nil)
}

func SubmitRegistration(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate()

	hashed, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	birthdayTime, err := time.Parse(time.RFC3339, r.FormValue("birthdate")+"T00:00:00Z")
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
		Birthdate: birthdayTime,
		Password:  hashed,
		Sex:       r.FormValue("sex"),
	}

	database.DB.Create(&user)

	// database.DB.
	tmpl.ExecuteTemplate(w, "base", nil)
}
