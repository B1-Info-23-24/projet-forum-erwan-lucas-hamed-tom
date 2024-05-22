package forum

import "gorm.io/gorm"

type User struct {
	gorm.Model
	id       int    `json:"id"`
	pseudo   string `json:"pseudo"`
	email    string `json:"email"`
	password string `json:"password"`
}
