package forum

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "test api"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "test api 2"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
