package forum
package forum

import (
	"html/template"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request, imgpath string) {
	template, err := template.ParseFiles("./index.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, imgpath)
}

func Profile(w http.ResponseWriter, r *http.Request, imgpath string) {
	template, err := template.ParseFiles("./pages/profile.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, imgpath)
}

