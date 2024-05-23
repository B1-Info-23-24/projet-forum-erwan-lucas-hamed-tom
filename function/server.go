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
	// Initialize the database
	InitDB()
	AutoMigrate(DB)

	// Initialize the router for Mux
	r := mux.NewRouter()

	r.HandleFunc("/login", RenderLoginPageHandler).Methods("GET")

	r.HandleFunc("/signup", renderSignupPage).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Start the Mux server in a goroutine
	go func() {
		log.Println("Mux server running at http://localhost:8080/")
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	// Initialize the Gin router for the API
	ginRouter := gin.Default()

	// Configure CORS middleware
	ginRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	RegisterRoutes(ginRouter)

	// Start the Gin server
	go func() {
		log.Println("Gin server running at http://localhost:8081/api/")
		log.Fatal(ginRouter.Run(":8081"))
	}()

	// Prevent the main function from exiting
	select {}
}

func renderSignupPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/signup.html")
}
