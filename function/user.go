package forum

import (
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"unique"`
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
