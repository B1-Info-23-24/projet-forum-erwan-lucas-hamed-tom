package forum

import (
	"log"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}
	var user User
	pseudo := r.FormValue("pseudo-sign")
	email := r.FormValue("email-sign")
	password := r.FormValue("password-sign")
	verifyPassword := r.FormValue("verify-password-sign")
	passwordEncrypt := Encrypt(password)
	valueCreate := [3]string{pseudo, email, passwordEncrypt}
	if verifyPassword != password {
		log.Println("ERROR : Incorrect password")
		errorMessage := "Different password"
		Home(w, r, errorMessage)
	} else if len(password) < 12 {
		log.Println("ERROR : Incorrect password")
		errorMessage := "Password is too short, 12 characters minimum"
		Home(w, r, errorMessage)
	} else if !VerifyPassword(password) {
		log.Println("ERROR : Incorrect password")
		errorMessage := "Your password must contain: lowercase, uppercase, special characters and number"
		Home(w, r, errorMessage)
	} else if !EmailValid(email) {
		log.Println("ERROR : Incorrect email")
		errorMessage := "The email entered is invalid"
		Home(w, r, errorMessage)
	} else {
		//concet user
		valueConnect := [2]string{pseudo, passwordEncrypt}
		if user.pseudo == "" {
			log.Println("ERROR : Wrong connection information")
			errorMessage := "Pseudo or email already used"
			Home(w, r, errorMessage)
		} else {
			http.Redirect(w, r, "/landingPage", http.StatusFound)
		}
	}
}
