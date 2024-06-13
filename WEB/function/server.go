package forumWeb

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartWebServer() {
	r := mux.NewRouter()

	r.HandleFunc("/", RenderHomePage).Methods("GET")
	r.HandleFunc("/login", RenderLoginPage).Methods("GET")
	r.HandleFunc("/signup", RenderSignupPage).Methods("GET")
	r.HandleFunc("/profile/{username}", RenderProfilePage).Methods("GET")
	r.HandleFunc("/edit", RenderEditPage).Methods("GET")
	r.HandleFunc("/maps/{postId}", RenderMaps).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(int(12*time.Hour)),
	)

	log.Println("Web server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}

func RenderSignupPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/signup.html")
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/login.html")
}

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func RenderProfilePage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/profile.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)

	vars := mux.Vars(r)
	username := vars["username"]
	cookie, err := r.Cookie("username")
	if err != nil || cookie.Value != username {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func RenderEditPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/edit.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func RenderMaps(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/maps.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}
