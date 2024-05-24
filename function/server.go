package forum

import (
	"net/http"
)

func Server() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		imgpath := "../static/img/character.png"
		Home(w, r, imgpath)
	})
	http.HandleFunc("/profile/", func(w http.ResponseWriter, r *http.Request) {
		imgpath := "../static/img/character.png"
		username := r.URL.Path[len("/profile/"):]
		if username == "" {
			http.NotFound(w, r)
			return
		}
		Profile(w, r, imgpath)
	})

	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	http.ListenAndServe(":8080", nil)
}
