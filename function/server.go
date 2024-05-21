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

	log.Println("http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", r))
}
