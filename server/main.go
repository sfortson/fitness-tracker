package main

import (
	"context"
	"log"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"github.com/sfortson/fitness-tracker/server/database"
	"github.com/sfortson/fitness-tracker/server/pages"
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

		var session database.Session
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

		var user database.User
		database.DB.Where("username = ?", session.Username).First(&user)

		contextKeySessionToken := pages.SessionToken("session-token")
		r = r.WithContext(context.WithValue(r.Context(), contextKeySessionToken, sessionToken))
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("Init DB...")
	database.Open()

	log.Println("Migrating DB...")
	database.Migrate()

	r := mux.NewRouter()
	r.HandleFunc("/", authToken(http.HandlerFunc(pages.HomePage))).Methods("GET", "POST")
	r.HandleFunc("/registration", pages.GetRegistration).Methods("GET")
	r.HandleFunc("/registration", pages.SubmitRegistration).Methods("POST")
	r.HandleFunc("/login", pages.Login).Methods("GET")
	r.HandleFunc("/login", pages.LoginPost).Methods("POST")
	r.Use(loggingMiddleware)

	log.Println("Listening...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
