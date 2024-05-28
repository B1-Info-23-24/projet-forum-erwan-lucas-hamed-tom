package forum

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Home(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func Profile(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("./pages/profile.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html")
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

func GetProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	log.Printf("Fetching profile for username: %s", username)

	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Error retrieving user from database: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "User not found: %v"}`, err.Error()), http.StatusNotFound)
		return
	}

	log.Printf("User found: %+v", user)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user data: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to encode user data: %v"}`, err.Error()), http.StatusInternalServerError)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) error {
	tmpl, err := template.ParseFiles("./pages/edit.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html")
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
