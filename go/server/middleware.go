package main

import (
	"log"
	"net/http"
)

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v %v", r.Method, r.URL, r.Proto)
		next.ServeHTTP(w, r)
	})
}

// func nosurfing(h http.Handler) http.Handler {
// 	surfing := nosurf.New(h)
// 	surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Failed to validate CSRF token:", nosurf.Reason(r))
// 		w.WriteHeader(http.StatusBadRequest)
// 	}))
// 	return surfing
// }
