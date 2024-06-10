package forumApi

import (
	"log"
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
