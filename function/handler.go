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
