package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/sfortson/fitness-tracker/server/database"
	"github.com/sfortson/fitness-tracker/server/migrations"
	"github.com/sfortson/fitness-tracker/server/models"
	"github.com/sfortson/fitness-tracker/server/pages"
)

// func loggingMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		m := httpsnoop.CaptureMetrics(next, w, r)
// 		log.Printf("%v %v %v %v", r.Method, r.URL, r.Proto, m.Code)
// 	})
// }

func authToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			// If the cookie is not set, return an unauthorized status
			log.Println("ERROR")
			// http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
			// http.Error(w, "Forbidden", http.StatusForbidden)
		}

		var session models.Session
		sessionToken := c.Value
		result := database.DB.Where("session_token = ?", sessionToken).First(&session)
		if result.Error != nil {
			// If the session token is not present in session map, return an unauthorized error
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}

		// If the session is present, but has expired, we can delete the session, and return
		// an unauthorized status
		if session.IsExpired() {
			database.DB.Delete(&session)
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}

		var user models.User
		database.DB.Where("username = ?", session.Username).First(&user)

		r = r.WithContext(context.WithValue(r.Context(), "user", c))
		next.ServeHTTP(w, r)
	})
}

func main() {
	token := base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64))
	log.Println(token)
	log.Println(base64.StdEncoding.DecodeString(token))
	log.Println("Init DB...")
	database.Open()

	log.Println("Migrating DB...")
	migrations.Migrate()

	r := mux.NewRouter()
	r.HandleFunc("/", pages.HomePage).Methods("GET", "POST")
	r.HandleFunc("/registration", pages.GetRegistration).Methods("GET")
	r.HandleFunc("/registration", pages.SubmitRegistration).Methods("POST")
	r.HandleFunc("/login", pages.Login).Methods("GET")
	r.HandleFunc("/login", pages.LoginPost).Methods("POST")
	r.Use(authToken)

	log.Println("Listening...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
