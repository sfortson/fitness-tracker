package main

import (
	"log"
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"github.com/sfortson/fitness-tracker/server/database"
	"github.com/sfortson/fitness-tracker/server/migrations"
	"github.com/sfortson/fitness-tracker/server/pages"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(next, w, r)
		log.Printf("%v %v %v %v", r.Method, r.URL, r.Proto, m.Code)
	})
}

func main() {
	log.Println("Init DB...")
	database.Open()

	log.Println("Migrating DB...")
	migrations.Migrate()

	r := mux.NewRouter()
	r.HandleFunc("/", pages.HomePage).Methods("GET", "POST")
	r.HandleFunc("/registration", pages.GetRegistration).Methods("GET")
	r.HandleFunc("/registration", pages.SubmitRegistration).Methods("POST")
	r.Use(loggingMiddleware)

	log.Println("Listening...")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal(err)
	}
}
