package forum

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type PostPingResponse struct {
	Post Post `json:"post"`
	Ping Ping `json:"ping"`
}

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
	r.HandleFunc("/api/profile/{username}", GetProfile).Methods("GET")
	r.HandleFunc("/api/edit", EditHandler).Methods("GET")
	r.HandleFunc("/api/editing/{username}", EditProfile).Methods("POST")
	r.HandleFunc("/api/delete/{username}", DeleteProfile).Methods("DELETE")
	r.HandleFunc("/api/post/create", CreatePost).Methods("POST")
	r.HandleFunc("/api/pings", GetPings).Methods("GET")
	r.HandleFunc("/api/post/display", DisplayPost).Methods("POST")
	r.HandleFunc("/api/post/display/{lat}/{lng}", GetCurrentPost).Methods("POST")
	r.HandleFunc("/api/post/display/{postId}", GetCurrentPostFromId).Methods("POST")
	r.HandleFunc("/api/comment/create/{postId}", CreateComment).Methods("POST")
	r.HandleFunc("/api/comment/{postId}", GetComments).Methods("GET")
	r.HandleFunc("/api/post/like/{postId}", LikePost).Methods("POST")       // New Route
	r.HandleFunc("/api/post/dislike/{postId}", DislikePost).Methods("POST") // New Route
}
func LikePost(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseUint(mux.Vars(r)["postId"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	userID := GetCoockie(w, r, "userId")

	var interaction UserPostInteraction
	if err := DB.Where("user_id = ? AND post_id = ?", userID, postId).First(&interaction).Error; err == nil {
		if interaction.Liked {
			http.Error(w, `{"error": "Post already liked"}`, http.StatusBadRequest)
			return
		}
		interaction.Liked = true
		interaction.Disliked = false
		if err := DB.Save(&interaction).Error; err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	} else {
		interaction = UserPostInteraction{
			UserID:   uint(userID),
			PostID:   uint(postId),
			Liked:    true,
			Disliked: false,
		}
		if err := DB.Create(&interaction).Error; err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	}

	var post Post
	if err := DB.First(&post, postId).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Post not found: %v"}`, err.Error()), http.StatusNotFound)
		return
	}

	post.Likes++
	if err := DB.Save(&post).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
		Post    Post   `json:"post"`
	}{
		Message: "Post liked successfully",
		Post:    post,
	}
	json.NewEncoder(w).Encode(response)
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseUint(mux.Vars(r)["postId"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	userID := GetCoockie(w, r, "userId")

	var interaction UserPostInteraction
	if err := DB.Where("user_id = ? AND post_id = ?", userID, postId).First(&interaction).Error; err == nil {
		if interaction.Disliked {
			http.Error(w, `{"error": "Post already disliked"}`, http.StatusBadRequest)
			return
		}
		interaction.Liked = false
		interaction.Disliked = true
		if err := DB.Save(&interaction).Error; err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	} else {
		interaction = UserPostInteraction{
			UserID:   uint(userID),
			PostID:   uint(postId),
			Liked:    false,
			Disliked: true,
		}
		if err := DB.Create(&interaction).Error; err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	}

	var post Post
	if err := DB.First(&post, postId).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Post not found: %v"}`, err.Error()), http.StatusNotFound)
		return
	}

	post.Dislikes++
	if err := DB.Save(&post).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
		Post    Post   `json:"post"`
	}{
		Message: "Post disliked successfully",
		Post:    post,
	}
	json.NewEncoder(w).Encode(response)
}

func register(w http.ResponseWriter, r *http.Request) {
	var M *Messages
	M = new(Messages)
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	fmt.Printf("Registering user: %+v\n", user)
	if !VerifyPassword(user.Password, M) {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, M.Messages), http.StatusBadRequest)
		return
	}
	if VerifyPassword(user.Password, M) && EmailValid(user.Email) {
		user.Password = Encrypt(user.Password)
		if err := DB.Create(&user).Error; err != nil {
			M.Messages = ""
			MessageError := err.Error()
			if MessageError == "UNIQUE constraint failed: users.email" {
				M.Messages = "Email deja utiliser"
			}
			if MessageError == "UNIQUE constraint failed: users.username" {
				M.Messages = "Pseudo deja utiliser"
			}
			log.Printf("Failed to create user: %v", err)
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, M.Messages), http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Failed to create user")
		return
	}

	fmt.Printf("New user registered: %s\n", user.Username)

	SetCookie(w, user)

	log.Printf("User details: %+v\n", user)

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
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	var user User

	if err := DB.Where("email = ?", loginInfo.Email).First(&user).Error; err != nil {
		http.Error(w, `{"error": "Email not found"}`, http.StatusUnauthorized)
		return
	}

	loginInfo.Password = Encrypt(loginInfo.Password)
	if user.Password != loginInfo.Password {
		http.Error(w, `{"error": "Invalid password"}`, http.StatusUnauthorized)
		return
	}

	DeleteCookies(w, r)
	fmt.Printf("User logged in: %s\n", user.Username)
	SetCookie(w, user)

	// Log the user details in the terminal
	log.Printf("User details: %+v\n", user)

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

	fmt.Printf("User profile deleted: %s\n", user.Username)

	w.Header().Set("Content-Type", "application/json")
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	log.Printf("Fetching profile for username: %s", username)

	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("Error retrieving user from database: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "User not found: %v"}`, err.Error()), http.StatusNotFound)
		return
	}

	log.Printf("User found: %+v", user)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Printf("Error encoding user data: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "Failed to encode user data: %v"}`, err.Error()), http.StatusInternalServerError)
	}
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	userID := GetCoockie(w, r, "userId")
	theme := r.FormValue("theme")
	content := r.FormValue("content")
	lat := r.FormValue("lat")
	lng := r.FormValue("lng")

	var user User
	if err := DB.Where("id = ?", userID).First(&user).Error; err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	newPost := Post{
		UserID:    user.ID,
		Theme:     theme,
		Content:   content,
		CreatedAt: time.Now(),
		Username:  user.Username,
	}

	if err := DB.Create(&newPost).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	newPing := Ping{
		PostID: newPost.ID,
		Lat:    lat,
		Lng:    lng,
	}

	if err := DB.Create(&newPing).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	var fileNames []string
	if err := json.Unmarshal([]byte(r.FormValue("images")), &fileNames); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to parse image file names: %v"}`, err), http.StatusBadRequest)
		return
	}

	for _, fileName := range fileNames {
		imageURL := fmt.Sprintf("../static/img/post/%s", fileName)
		newImage := Image{
			PostID: newPost.ID,
			URL:    imageURL,
		}
		if err := DB.Create(&newImage).Error; err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string `json:"message"`
		Post    Post   `json:"post"`
	}{
		Message: "Post created successfully",
		Post:    newPost,
	}
	json.NewEncoder(w).Encode(response)
}

func GetPings(w http.ResponseWriter, r *http.Request) {
	var pings []Ping
	if err := DB.Find(&pings).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pings)
}

func DisplayPost(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	// Récupérer tous les posts de la base de données
	if err := DB.Find(&posts).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Pour chaque post, récupérer les images associées
	for i := range posts {
		var images []Image
		if err := DB.Where("post_id = ?", posts[i].ID).Find(&images).Error; err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		posts[i].Images = images
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	postId, err := strconv.ParseUint(mux.Vars(r)["postId"], 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	comment.CreatedAt = time.Now()
	comment.PostID = uint(postId)
	comment.UserID = uint(GetCoockie(w, r, "userId"))
	comment.Username = GetCoockieAsString(w, r, "username")

	if err := DB.Create(&comment).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Message string  `json:"message"`
		Comment Comment `json:"comment"`
	}{
		Message: "Comment created successfully",
		Comment: comment,
	}
	json.NewEncoder(w).Encode(response)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	postID := mux.Vars(r)["postId"]

	var comments []Comment
	if err := DB.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func GetCurrentPost(w http.ResponseWriter, r *http.Request) {
	lat := mux.Vars(r)["lat"]
	lng := mux.Vars(r)["lng"]
	var ping Ping

	if err := DB.Where("lat = ? AND lng = ?", lat, lng).First(&ping).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Ping not found: %v"}`, err.Error()), http.StatusNotFound)
		return
	}

	var post Post
	if err := DB.Where("id = ?", ping.ID).First(&post).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Post not found: %v"}`, err.Error()), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func GetCurrentPostFromId(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["postId"]

	var post Post
	if err := DB.Where("id = ?", postId).First(&post).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Post not found: %v"}`, err.Error()), http.StatusNotFound)
		return
	}

	var ping Ping
	if err := DB.Where("post_id = ?", postId).First(&ping).Error; err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Ping not found: %v"}`, err.Error()), http.StatusNotFound)
		return
	}

	response := PostPingResponse{
		Post: post,
		Ping: ping,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
