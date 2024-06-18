package forum

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Messages struct {
	Messages string
}

func StartWebServer() {
	InitDB()
	AutoMigrate(DB)

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r)
	}).Methods("GET")
	r.HandleFunc("/profile/", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Path[len("/profile/"):]
		if username == "" {
			http.NotFound(w, r)
			return
		}
		Profile(w, r)
	}).Methods("GET")
	r.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		EditHandler(w, r)
	}).Methods("GET")
	r.HandleFunc("/alpinisme", func(w http.ResponseWriter, r *http.Request) {
		AlpinismeHandler(w, r)
	}).Methods("GET")
	r.HandleFunc("/randonne", func(w http.ResponseWriter, r *http.Request) {
		RandonneHandler(w, r)
	}).Methods("GET")
	r.HandleFunc("/treck", func(w http.ResponseWriter, r *http.Request) {
		TreckHandler(w, r)
	}).Methods("GET")
	r.HandleFunc("/bivouac", func(w http.ResponseWriter, r *http.Request) {
		BivouacHandler(w, r)
	}).Methods("GET")
	r.HandleFunc("/maps/{postId}", MapsHandler).Methods("GET")
	r.HandleFunc("/delete/{username}", func(w http.ResponseWriter, r *http.Request) {
		DeleteProfile(w, r)
	}).Methods("DELETE")
	r.HandleFunc("/profile/{username}", ProfileHandler).Methods("GET")
	r.HandleFunc("/editing/{username}", EditProfile).Methods("POST")
	r.HandleFunc("/post/create", CreatePost).Methods("POST")
	r.HandleFunc("/login", RenderLoginPage).Methods("GET")
	r.HandleFunc("/signup", RenderSignupPage).Methods("GET")
	r.HandleFunc("/api/pings", GetPings).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.HandleFunc("/login/github", githubLoginHandler)
	r.HandleFunc("/callback/github", githubCallbackHandler)

	r.HandleFunc("/login/facebook", facebookLoginHandler)
	r.HandleFunc("/callback/facebook", facebookCallbackHandler)

	r.HandleFunc("/login/google", googleLoginHandler)
	r.HandleFunc("/callback/google", googleCallbackHandler)

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
