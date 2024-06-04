package forumApi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func StartAPIServer() {
	r := mux.NewRouter()
	RegisterRoutes(r)

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(int(12*time.Hour)),
	)

	log.Println("API server running at http://localhost:8081/")
	log.Fatal(http.ListenAndServe(":8081", corsMiddleware(r)))
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/register", register).Methods("POST")
	r.HandleFunc("/api/login", login).Methods("POST")
	r.HandleFunc("/api/profile/{username}", GetProfile).Methods("GET")
	r.HandleFunc("/api/editing/{username}", EditProfile).Methods("POST")
	r.HandleFunc("/api/delete/{username}", DeleteProfile).Methods("DELETE")
}

func register(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Log user details for debugging
	log.Printf("Registering user: %+v\n", user)

	if !VerifyPassword(user.Password) || !EmailValid(user.Email) {
		log.Printf("Invalid password or email format")
		http.Error(w, `{"error": "Invalid password or email format"}`, http.StatusBadRequest)
		return
	}
	user.Password = Encrypt(user.Password)
	if DB == nil {
		log.Fatalf("Database not initialized")
		http.Error(w, `{"error": "Database not initialized"}`, http.StatusInternalServerError)
		return
	}

	if err := DB.Create(&user).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	log.Printf("New user registered: %s\n", user.Username)

	SetCookie(w, user)

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
		User    User   `json:"user"`
	}{
		Message: "User registered and logged in successfully",
		User:    user,
	}
	json.NewEncoder(w).Encode(response)
}

func login(w http.ResponseWriter, r *http.Request) {
	var loginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginInfo); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	loginInfo.Password = Encrypt(loginInfo.Password)
	var user User

	if DB == nil {
		log.Fatalf("Database not initialized")
		http.Error(w, `{"error": "Database not initialized"}`, http.StatusInternalServerError)
		return
	}

	if err := DB.Where("email = ? AND password = ?", loginInfo.Email, loginInfo.Password).First(&user).Error; err != nil {
		log.Printf("Invalid email or password")
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	DeleteCookies(w, r)
	SetCookie(w, user)

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
		User    User   `json:"user"`
	}{
		Message: "Login successful",
		User:    user,
	}
	json.NewEncoder(w).Encode(response)
}

func EditProfile(w http.ResponseWriter, r *http.Request) {
	var editInfo struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&editInfo); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	currentUsername := mux.Vars(r)["username"]
	var user User

	if DB == nil {
		log.Fatalf("Database not initialized")
		http.Error(w, `{"error": "Database not initialized"}`, http.StatusInternalServerError)
		return
	}

	if err := DB.Where("username = ?", currentUsername).First(&user).Error; err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "User not found"}`), http.StatusNotFound)
		return
	}

	if editInfo.Username != "" {
		user.Username = editInfo.Username
	}
	if editInfo.Email != "" {
		user.Email = editInfo.Email
	}

	if err := DB.Save(&user).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	DeleteCookies(w, r)
	SetCookie(w, user)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Profile updated successfully"}`))
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	currentUsername := mux.Vars(r)["username"]
	var user User

	if DB == nil {
		log.Fatalf("Database not initialized")
		http.Error(w, `{"error": "Database not initialized"}`, http.StatusInternalServerError)
		return
	}

	if err := DB.Where("username = ?", currentUsername).First(&user).Error; err != nil {
		log.Printf("User not found: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "User not found"}`), http.StatusNotFound)
		return
	}

	if err := DB.Delete(&user).Error; err != nil {
		log.Printf("Failed to delete user: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	DeleteCookies(w, r)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Profile deleted successfully"}`))
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	log.Printf("Fetching profile for username: %s", username)

	var user User

	if DB == nil {
		log.Fatalf("Database not initialized")
		http.Error(w, `{"error": "Database not initialized"}`, http.StatusInternalServerError)
		return
	}

	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Error retrieving user from database: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "User not found"}`), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user data: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to encode user data"}`), http.StatusInternalServerError)
	}
}
