package forum

import (
	"net/http"
	"path/filepath"
	"text/template"
)

func RenderLoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("pages", "login.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Unable to load login page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func HandleLoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	isAuthenticated, err := AuthenticateUser(db, email, password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	if isAuthenticated {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}