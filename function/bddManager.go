package forum

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
}

type Image struct {
	ID        uint      `gorm:"primary_key"`
	PostID    uint      `gorm:"index"`
	URL       string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Image{})
}
