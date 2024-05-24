package forum

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

func Server() {

	InitDB()
	AutoMigrate(DB)

	r := mux.NewRouter()

	r.HandleFunc("/login", RenderLoginPageHandler).Methods("GET")
	r.HandleFunc("/signup", RenderSignupPageHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	go func() {
		log.Println("Mux server running at http://localhost:8080/")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	ginRouter := gin.Default()

	ginRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	RegisterRoutes(ginRouter)

	go func() {
		log.Println("Gin server running at http://localhost:8081/api/")
		log.Fatal(ginRouter.Run(":8081"))
	}()

	select {}
}

func RenderSignupPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/signup.html")
}

func RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/login.html")
}
