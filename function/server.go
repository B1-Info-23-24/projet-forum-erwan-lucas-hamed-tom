package forum

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartWebServer() {
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

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8080", "http://localhost:8081"}), // Autoriser les deux serveurs
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(int(12*time.Hour)),
	)

	log.Println("Web server running at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(r)))
}

func RenderSignupPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/signup.html")
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/login.html")
}
