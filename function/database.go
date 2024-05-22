package forum

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func init() {
	db = InitDB()
	addFakeAccounts(db)
}

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS USER (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        pseudo TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func InsertUser(db *sql.DB, username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO USER (pseudo, email, password) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, email, hashedPassword)
	return err
}

func AuthenticateUser(db *sql.DB, email, password string) (bool, error) {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM USER WHERE email = ?", email).Scan(&hashedPassword)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil, err
}

func addFakeAccounts(db *sql.DB) {
	users := []struct {
		username string
		email    string
		password string
	}{
		{"testuser1", "test1@example.com", "password123"},
		{"testuser2", "test2@example.com", "password123"},
		{"testuser3", "test3@example.com", "password123"},
	}

	for _, user := range users {
		err := InsertUser(db, user.username, user.email, user.password)
		if err != nil {
			log.Printf("Error inserting user %s: %v", user.username, err)
		}
	}
}
