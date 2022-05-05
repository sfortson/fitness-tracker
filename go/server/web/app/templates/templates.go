package templates

import (
	"html/template"
	"log"
)

var Templates = make(map[string]*template.Template)

func InitTemplates() {
	home := "server/templates/home.html"
	t, err := template.ParseFiles(
		home,
		"server/templates/base.html")
	if err != nil {
		log.Fatal(err)
	}

	login := "server/templates/login.html"
	loginTemplate, loginErr := template.ParseFiles(
		login,
		"server/templates/base.html")
	if loginErr != nil {
		log.Fatal(err)
	}

	Templates["home"] = template.Must(t, err)
	Templates["login"] = template.Must(loginTemplate, loginErr)
}
