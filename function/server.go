package forum

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Server() {
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/hello", HelloHandler).Methods("GET")
	r.HandleFunc("/login", RenderLoginPageHandler).Methods("GET")
	r.HandleFunc("/login", HandleLoginHandler).Methods("POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", r))
}
