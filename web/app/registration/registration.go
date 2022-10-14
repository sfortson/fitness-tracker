package registration

import (
	"log"
	"net/http"
	"time"

	"github.com/sfortson/fitness-tracker/internal/database"
	"github.com/sfortson/fitness-tracker/internal/helpers"
	templates "github.com/sfortson/fitness-tracker/web/app"
	"golang.org/x/crypto/bcrypt"
)

type FormValues struct {
	Username  string
	Sex       string
	Birthdate time.Time
	Email     string
	Password  string
	Height    float64
}

func GetRegistration(w http.ResponseWriter, r *http.Request) {
	tmpl := templates.WebTemplates["register"]
	tmpl.ExecuteTemplate(w, "base", nil)
}

func SubmitRegistration(w http.ResponseWriter, r *http.Request) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	birthdayTime, err := time.Parse(time.RFC3339, r.FormValue("birthdate")+"T00:00:00Z")
	if err != nil {
		log.Fatal(err)
	}

	user := database.User{
		Email:     r.FormValue("email"),
		Username:  r.FormValue("username"),
		Birthdate: birthdayTime,
		Password:  hashed,
		Sex:       r.FormValue("sex"),
		Height:    helpers.ParseFloat(r.FormValue("height")),
	}

	database.DB.Create(&user)

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
