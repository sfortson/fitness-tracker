package templates

import (
	"html/template"
	"log"
)

var WebTemplates = make(map[string]*template.Template)

func InitWebTemplates() {
	homepage, err := template.ParseFiles("web/tmpl/home.html", "web/tmpl/base.html")
	if err != nil {
		log.Fatal(err)
	}
	hptmpl := template.Must(homepage, err)

	login, err := template.ParseFiles("web/tmpl/login.html", "web/tmpl/base.html")
	if err != nil {
		log.Fatal(err)
	}
	ltmpl := template.Must(login, err)

	register, err := template.ParseFiles("web/tmpl/registration.html", "web/tmpl/base.html")
	if err != nil {
		log.Fatal(err)
	}
	rtmpl := template.Must(register, err)

	WebTemplates["homepage"] = hptmpl
	WebTemplates["login"] = ltmpl
	WebTemplates["register"] = rtmpl
}
