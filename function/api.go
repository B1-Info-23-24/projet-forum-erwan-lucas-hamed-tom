package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/register", register).Methods("POST")
	r.HandleFunc("/api/login", login).Methods("POST")
	r.HandleFunc("/profile/{username}", ProfileHandler).Methods("GET")
}

func register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Log user details for debugging (Remove in production)
	fmt.Printf("Registering user: %+v\n", user)

	if err := DB.Create(&user).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Printf("New user registered: %s\n", user.Username)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "User registered successfully"}`))
}

func login(w http.ResponseWriter, r *http.Request) {
	var loginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	var user User
	if err := DB.Where("email = ? AND password = ?", loginInfo.Email, loginInfo.Password).First(&user).Error; err != nil {
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	DeleteCookies(w, r)
	fmt.Printf("User logged in: %s\n", user.Username)
	SetCookie(w, user)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"message": "Login successful", "user": "%s"}`, user.Username)))
}

func UserInformation(w http.ResponseWriter, r *http.Request) {

}
