package forum

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id       int    `json:"id"`
	Pseudo   string `json:"pseudo"`
	Email    string `json:"email"`
	Password string `json:"-"`
}
