package forum

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
		handlers.AllowedOrigins([]string{"http://localhost:8080"}),
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
	r.HandleFunc("/profile/{username}", ProfileHandler).Methods("GET")
	r.HandleFunc("/api/edit", EditHandler).Methods("GET")
	r.HandleFunc("/editing/{username}", EditProfile).Methods("POST")
	r.HandleFunc("/editingPassword/{username}", EditPassword).Methods("POST")
	r.HandleFunc("/delete/{username}", DeleteProfile).Methods("DELETE")
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

	for name, values := range w.Header() {
		for _, value := range values {
			log.Printf("Header: %s=%s\n", name, value)
		}
	}
}

func EditProfile(w http.ResponseWriter, r *http.Request) {
	var editInfo struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		// Oldpassword string `json:"oldpassword"`
		// Password    string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&editInfo); err != nil {
		http.Error(w, fmt.Sprintf(`{"error newdecoder": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	currentUsername := GetUserFromURL(w, r)

	var user User
	if err := DB.Where("username = ?", currentUsername).First(&user).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "User not found"}`), http.StatusNotFound)
		return
	}

	if editInfo.Username != "" {
		user.Username = editInfo.Username
	}
	if editInfo.Email != "" {
		user.Email = editInfo.Email
	}

	// if editInfo.Oldpassword != "" {
	// 	if editInfo.Oldpassword == user.Password {
	// 		if editInfo.Password != "" {
	// 			user.Password = editInfo.Password
	// 		}
	// 	}
	// }

	if err := DB.Save(&user).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
		http.Error(w, fmt.Sprintf(`{"error in saving": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	DeleteCookies(w, r)
	SetCookie(w, user)

	fmt.Printf("User profile updated: %s\n", user.Username)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Profile updated successfully"}`))
}

func EditPassword(w http.ResponseWriter, r *http.Request) {
	var editInfo struct {
		Oldpassword string `json:"oldpassword"`
		Password    string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&editInfo); err != nil {
		http.Error(w, fmt.Sprintf(`{"error newdecoder": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	currentUsername := GetUserFromURL(w, r)

	var user User
	if err := DB.Where("username = ?", currentUsername).First(&user).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "User not found"}`), http.StatusNotFound)
		return
	}

	editInfo.Oldpassword = Encrypt(editInfo.Oldpassword)
	if VerifyPassword(editInfo.Password) {
		editInfo.Password = Encrypt(editInfo.Password)
	}

	if editInfo.Oldpassword == user.Password {
		if editInfo.Password != "" {
			user.Password = editInfo.Password
		}
	}

	if err := DB.Save(&user).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
		http.Error(w, fmt.Sprintf(`{"error in saving": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	DeleteCookies(w, r)
	SetCookie(w, user)

	fmt.Printf("User profile updated: %s\n", user.Username)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Password updated successfully"}`))
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	currentUsername := GetUserFromURL(w, r)

	var user User
	if err := DB.Where("username = ?", currentUsername).First(&user).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "User not found"}`), http.StatusNotFound)
		return
	}

	if err := DB.Delete(&user).Error; err != nil {
		log.Printf("Failed to delete user: %v", err)
		http.Error(w, fmt.Sprintf(`{"error in deletion": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	DeleteCookies(w, r)

	fmt.Printf("User profile deleted: %s\n", user.Username)

	w.Header().Set("Content-Type", "application/json")
}
