package forum

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var templates = template.Must(template.ParseGlob("templates/header.html"))

func Server() {
	InitDB()
	AutoMigrate(DB)

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		imgpath := "../static/img/character.png"
		Home(w, r, imgpath)
	}).Methods("GET")
	r.HandleFunc("/profile/", func(w http.ResponseWriter, r *http.Request) {
		imgpath := "../static/img/character.png"
		username := r.URL.Path[len("/profile/"):]
		if username == "" {
			http.NotFound(w, r)
			return
		}
		Profile(w, r, imgpath)
	}).Methods("GET")
	r.HandleFunc("/profile/{username}", ProfileHandler).Methods("GET")
	r.HandleFunc("/login", RenderLoginPage).Methods("GET")
	r.HandleFunc("/signup", RenderSignupPage).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	RegisterRoutes(r)

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8080"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(int(12*time.Hour)),
	)

	log.Println("Server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}

func RenderSignupPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/signup.html")
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/login.html")
}
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		http.NotFound(w, r)
		return
	}
	data := map[string]interface{}{
		"Title":        "Profile Page",
		"ProfileImage": "/static/img/character.png",
		"Username":     username,
	}
	templates.ExecuteTemplate(w, "base", data)
}
