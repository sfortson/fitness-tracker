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
			log.Println("cookie not set")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		var session database.Session
		sessionToken := c.Value
		result := database.DB.Where("session_token = ?", sessionToken).First(&session)
		if result.Error != nil {
			// If the session token is not present in session map, return an unauthorized error
			log.Println("session token not present")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		// If the session is present, but has expired, we can delete the session, and return
		// an unauthorized status
		if session.IsExpired() {
			database.DB.Delete(&session)
			log.Println("session token expired")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		log.Println("session token good to go")
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
	r.HandleFunc("/register", pages.GetRegistration).Methods("GET")
	r.HandleFunc("/register", pages.SubmitRegistration).Methods("POST")
	r.HandleFunc("/login", pages.Login).Methods("GET")
	r.HandleFunc("/login", pages.LoginPost).Methods("POST")
	r.HandleFunc("/logout", pages.Logout)
	r.Use(loggingMiddleware)

	log.Println("Listening...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
