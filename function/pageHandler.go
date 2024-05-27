package forum

import (
	"html/template"
	"log"
	"net/http"
)

type Userinfo struct {
	Imgpth   string
	Username string
}

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
	var userinfo Userinfo
	userinfo.Imgpth = imgpath
	userinfo.Username = GetUserFromURL(w, r)
	err = tmpl.Execute(w, userinfo)
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
