package main

import (
	"html/template"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("got to this route")
	t, err := template.ParseFiles("server/templates/index.html")
	if err != nil {
		log.Println(err)
		return
	}
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8000", nil)
}
