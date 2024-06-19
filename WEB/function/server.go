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
	r.HandleFunc("/alpinisme", RenderAlpinismePage).Methods("GET")
	r.HandleFunc("/randonne", RenderRandonnePage).Methods("GET")
	r.HandleFunc("/treck", RenderTreckPage).Methods("GET")
	r.HandleFunc("/bivouac", RenderBivouacPage).Methods("GET")
	r.HandleFunc("/maps/{postId}", RenderMapsPage).Methods("GET")
	r.HandleFunc("/profile/{username}", RenderProfilePage).Methods("GET")
	r.HandleFunc("/edit", RenderEditPage).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	// r.HandleFunc("/login/github", githubLoginHandler)
	// r.HandleFunc("/callback/github", githubCallbackHandler)

	// r.HandleFunc("/login/facebook", facebookLoginHandler)
	// r.HandleFunc("/callback/facebook", facebookCallbackHandler)

	// r.HandleFunc("/login/google", googleLoginHandler)
	// r.HandleFunc("/callback/google", googleCallbackHandler)

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept", "UserID", "Username"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(int(12*time.Hour)),
	)

	log.Println("Web server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}

func RenderHomePage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./index.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html", "./templates/post-modif.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
	}
	template.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
func RenderAlpinismePage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/alpinisme.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html", "./templates/post-modif.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
	}
	template.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
func RenderRandonnePage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/randonné.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html", "./templates/post-modif.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
	}
	template.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
func RenderTreckPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/treck.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html", "./templates/post-modif.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
	}
	template.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
func RenderBivouacPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/bivouac.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html", "./templates/post-modif.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
	}
	template.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}

func RenderMapsPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/maps.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html", "./templates/post-modif.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
	}
	template.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
func RenderProfilePage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/profile.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html", "./templates/post-modif.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
	}
	template.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
func RenderEditPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/edit.html", "./templates/header.html", "./templates/menu.html", "./templates/login.html", "./templates/signup.html", "./templates/post.html", "./templates/post-modif.html")
	if err != nil {
		log.Println("Error parsing template files:", err)
	}
	template.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}
