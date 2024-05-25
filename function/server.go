package forum

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Server() {
	InitDB()
	AutoMigrate(DB)

	r := mux.NewRouter()
	r.HandleFunc("/", Test)

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
