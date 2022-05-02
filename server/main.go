package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"github.com/sfortson/fitness-tracker/server/pages"
	"golang.org/x/crypto/bcrypt"
)

func registration(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("server/templates/registration.html", "server/templates/base.html")
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		m := httpsnoop.CaptureMetrics(next, w, r)
		log.Printf("%v %v %v %v", r.Method, r.URL, r.Proto, m.Code)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", pages.HomePage).Methods("GET", "POST")
	r.HandleFunc("/registration", registration)
	r.Use(loggingMiddleware)

	log.Println("Listening...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
