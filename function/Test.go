package forum

import (
	"html/template"
	"log"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html", "./templates/header.html", "./templates/menu.html", "./pages/testSignup.html", "./pages/testLogin.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, r)
}
