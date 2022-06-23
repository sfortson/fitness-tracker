package main

import (
	"context"
	"log"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"github.com/sfortson/fitness-tracker/internal/config"
	"github.com/sfortson/fitness-tracker/internal/database"
	"github.com/sfortson/fitness-tracker/internal/session"
	templates "github.com/sfortson/fitness-tracker/web/app"
	"github.com/sfortson/fitness-tracker/web/app/homepage"
	"github.com/sfortson/fitness-tracker/web/app/login"
	"github.com/sfortson/fitness-tracker/web/app/logout"
	"github.com/sfortson/fitness-tracker/web/app/registration"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)
		log.Printf("%v %v %v %v", r.Method, r.URL, r.Proto, m.Code)
	})
}

func authToken(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			// If the cookie is not set, return an unauthorized status
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		var dbSession database.Session
		sessionToken := c.Value
		result := database.DB.Where("session_token = ?", sessionToken).First(&dbSession)
		if result.Error != nil {
			// If the session token is not present in session map, return an unauthorized error
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		// If the session is present, but has expired, we can delete the session, and return
		// an unauthorized status
		if dbSession.IsExpired() {
			database.DB.Delete(&dbSession)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		var user database.User
		database.DB.Where("username = ?", dbSession.Username).First(&user)

		contextKeySessionToken := session.ContextKeySessionToken
		r = r.WithContext(context.WithValue(r.Context(), contextKeySessionToken, sessionToken))
		next.ServeHTTP(w, r)
	})
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/static/favicon.ico")
}

func main() {
	config, configerr := config.LoadConfig("config")
	if configerr != nil {
		log.Fatalln("unable to load config")
	}
	log.Println("Load Config...")

	log.Println("Init DB...")
	database.Open(config)

	log.Println("Migrating DB...")
	database.Migrate()

	log.Println("Preparing templates...")
	templates.InitWebTemplates()

	r := mux.NewRouter()
	r.HandleFunc("/favicon.ico", faviconHandler)
	r.HandleFunc("/", authToken(http.HandlerFunc(homepage.GetHome))).Methods("GET")
	r.HandleFunc("/", authToken(http.HandlerFunc(homepage.PostHome))).Methods("POST")
	r.HandleFunc("/register", registration.GetRegistration).Methods("GET")
	r.HandleFunc("/register", registration.SubmitRegistration).Methods("POST")
	r.HandleFunc("/login", login.Login).Methods("GET")
	r.HandleFunc("/login", login.LoginPost).Methods("POST")
	r.HandleFunc("/logout", logout.Logout)
	r.Use(loggingMiddleware)

	log.Println("Listening...")
	err := http.ListenAndServe(":80", r)
	if err != nil {
		log.Fatal(err)
	}
}
