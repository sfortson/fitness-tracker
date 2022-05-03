package pages

import (
	"net/http"

	"github.com/sfortson/fitness-tracker/server/database"
	"github.com/sfortson/fitness-tracker/server/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate("login")
	tmpl.ExecuteTemplate(w, "base", nil)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	var user models.User
	database.DB.Where("username = ?", r.FormValue("username")).First(&user)

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.FormValue("password"))) 
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
