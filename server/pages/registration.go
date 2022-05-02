package pages

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func GetRegistration(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		"server/templates/registration.html",
		"server/templates/base.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl := template.Must(t, err)
	if r.Method != http.MethodPost {
		tmpl.ExecuteTemplate(w, "base", nil)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hashed)

	tmpl.ExecuteTemplate(w, "base", nil)
}
