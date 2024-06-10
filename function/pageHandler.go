package forum

import (
	"html/template"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Profile(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("./pages/profile.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
		return err
	}
	var user User
	user.Username = GetUserFromURL(w, r)
	err = tmpl.Execute(w, user)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}
	return nil
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	err := Profile(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Edit(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("./pages/edit.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
		return err
	}
	var user User
	user.Username = GetUserFromURL(w, r)
	err = tmpl.Execute(w, user)
	if err != nil {
		log.Println("Error executing template:", err)
		return err
	}
	return nil
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	err := Edit(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Maps(w http.ResponseWriter, r *http.Request) error {
	template, err := template.ParseFiles("./pages/maps.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html")
	if err != nil {
		log.Fatal(err)
		return err
	}
	template.Execute(w, nil)
	return nil
}

func MapsHandler(w http.ResponseWriter, r *http.Request) {
	err := Maps(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
