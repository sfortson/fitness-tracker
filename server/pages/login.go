package pages

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sfortson/fitness-tracker/server/database"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl := getTemplate("login")
	tmpl.ExecuteTemplate(w, "base", nil)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	var user database.User
	database.DB.Where("username = ?", r.FormValue("username")).First(&user)

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(r.FormValue("password")))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(30 * time.Minute)

	session := database.Session{
		Username:     user.Username,
		Expiry:       expiresAt,
		SessionToken: sessionToken,
	}
	database.DB.Create(&session)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
