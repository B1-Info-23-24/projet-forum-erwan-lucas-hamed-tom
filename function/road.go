package forum

import "github.com/gorilla/mux"

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Définir les routes
	r.HandleFunc("/users", HandlerGetUsers).Methods("GET")
	r.HandleFunc("/users", HandlerCreateUser).Methods("POST")

	return r
}
