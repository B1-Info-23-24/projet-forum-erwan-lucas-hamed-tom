package forum

import (
	"database/sql"
	"log"
)

// create a new user in database
func CreateUser(db *sql.DB, value [3]string) {
	if checkUser(db, value) == 0 {
		insertQuery := "INSERT INTO USER (id, pseudo, email, password) VALUES (?, ?, ?, ?)"
		_, err := db.Exec(insertQuery, nil, value[0], value[1], value[2])
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Return user information as User
func connectUser(db *sql.DB, value [2]string) User {
	var u User
	db.QueryRow("SELECT * FROM `USER` WHERE pseudo = ? AND password = ? OR email = ? AND password = ?", value[0], value[1], value[0], value[1]).Scan(&u.id, &u.pseudo, &u.email, &u.password)
	return u
}

func checkUser(db *sql.DB, value [3]string) int {
	var nbAccount int
	query := "SELECT COUNT(*) FROM USER WHERE pseudo = ? OR email = ?"
	err := db.QueryRow(query, value[0], value[1]).Scan(&nbAccount)
	if err != nil {
		log.Fatal(err)
	}
	return nbAccount
}
