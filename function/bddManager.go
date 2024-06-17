package forum

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"unique"`
}

type Post struct {
	ID        uint      `gorm:"primary_key"`
	UserID    uint      `gorm:"index"`
	Theme     string    `gorm:"type:varchar(50)"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Images    []Image   `gorm:"foreignkey:PostID"`
	Username  string    `gorm:"type:text"`
	Likes     int       `gorm:"default:0"`
	Dislikes  int       `gorm:"default:0"`
}
type Image struct {
	ID     uint   `gorm:"primary_key"`
	PostID uint   `gorm:"index"`
	URL    string `gorm:"type:varchar(255)"`
}

type Ping struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	PostID uint   `gorm:"index" json:"post_id"`
	Lat    string `gorm:"type:varchar(255)" json:"lat"`
	Lng    string `gorm:"type:varchar(255)" json:"lng"`
}

type UserPostInteraction struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint `gorm:"index"`
	PostID   uint `gorm:"index"`
	Liked    bool
	Disliked bool
}

type Comment struct {
	ID        uint      `gorm:"primary_key"`
	PostID    uint      `gorm:"index"`
	UserID    uint      `gorm:"index"`
	Username  string    `gorm:"type:text"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	AutoMigrate(DB)
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Ping{})
	db.AutoMigrate(&Image{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&UserPostInteraction{})
}

// Recherche par nom d'utilisateur
func RechercherUtilisateurParNom(username string) (*User, error) {
	var user User
	result := DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Recherche par thème de publication
func RechercherPublicationsParTheme(theme string) ([]Post, error) {
	var posts []Post
	result := DB.Where("theme = ?", theme).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

// Recherche par ID utilisateur
func RechercherPublicationsParUtilisateur(userID uint) ([]Post, error) {
	var posts []Post
	result := DB.Where("user_id = ?", userID).Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return posts, nil
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	// Supposons que nous recherchons par thème de publication
	posts, err := RechercherPublicationsParTheme(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
