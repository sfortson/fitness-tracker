package login

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sfortson/fitness-tracker/server/internal/database"
	templates "github.com/sfortson/fitness-tracker/server/web/app"
	"golang.org/x/crypto/bcrypt"
)

type loginform struct {
	username string
	password string
	Errors   map[string]string
}

func (lf *loginform) validateLoginForm() (database.User, bool) {
	lf.Errors = make(map[string]string)

	if lf.username == "" {
		lf.Errors["username"] = "Must enter a username"
	}

	if lf.password == "" {
		lf.Errors["password"] = "Must enter a password"
	}

	var user database.User
	result := database.DB.Where("username = ?", lf.username).First(&user)

	if result.Error != nil {
		lf.Errors["usernamePassword"] = "Username or password is incorrect"
	} else {
		err := bcrypt.CompareHashAndPassword(user.Password, []byte(lf.password))
		if err != nil {
			lf.Errors["usernamePassword"] = "Username or password is incorrect"
		}
	}

	return user, len(lf.Errors) == 0
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl := templates.WebTemplates["login"]
	tmpl.ExecuteTemplate(w, "base", nil)
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	loginForm := loginform{
		username: r.FormValue("username"),
		password: r.FormValue("password"),
	}

	user, validated := loginForm.validateLoginForm()

	if !validated {
		tmpl := templates.WebTemplates["login"]
		tmpl.ExecuteTemplate(w, "base", loginForm)
		return
	}

	var session database.Session
	oldToken := database.DB.Where("username = ?", loginForm.username).First(&session)

	if oldToken.Error == nil {
		// If a session token already exists delete it before issuing a new one
		database.DB.Delete(&session)
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(30 * time.Minute)

	newSession := database.Session{
		Username:     user.Username,
		Expiry:       expiresAt,
		SessionToken: sessionToken,
	}
	database.DB.Create(&newSession)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
