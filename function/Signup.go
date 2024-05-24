package forum

import (
	"html/template"
	"log"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("pages/signup.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, r)
}
