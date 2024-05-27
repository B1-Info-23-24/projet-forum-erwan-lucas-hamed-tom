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

func Profile(w http.ResponseWriter, r *http.Request, imgpath string) error {
	tmpl, err := template.ParseFiles("./pages/profile.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
		return err
	}
	err = tmpl.Execute(w, imgpath)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}
	return nil
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	imgpath := "/path/to/user/image"
	err := Profile(w, r, imgpath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
