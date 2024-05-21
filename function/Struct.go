package forum

type User struct {
	id       int    `json:"id"`
	pseudo   string `json:"pseudo"`
	email    string `json:"email"`
	password string `json:"password"`
}
