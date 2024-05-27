package forum

import (
	"html/template"
	"log"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html", "./templates/header.html", "./templates/menu.html", "./templates/Signup.html", "./templates/Login.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, r)
}
