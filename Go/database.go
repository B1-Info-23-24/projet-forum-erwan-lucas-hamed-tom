package forum

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func DataX() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createUsersTable := `
    CREATE TABLE IF NOT EXISTS Users (
        user_id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        last_login DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	if _, err := db.Exec(createUsersTable); err != nil {
		log.Fatal(err)
	}

	createCategoriesTable := `
    CREATE TABLE IF NOT EXISTS Categories (
        category_id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    `
	if _, err := db.Exec(createCategoriesTable); err != nil {
		log.Fatal(err)
	}

	createThreadsTable := `
    CREATE TABLE IF NOT EXISTS Threads (
        thread_id INTEGER PRIMARY KEY AUTOINCREMENT,
        category_id INTEGER,
        user_id INTEGER,
        title TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (category_id) REFERENCES Categories(category_id),
        FOREIGN KEY (user_id) REFERENCES Users(user_id)
    );
    `
	if _, err := db.Exec(createThreadsTable); err != nil {
		log.Fatal(err)
	}

	createPostsTable := `
    CREATE TABLE IF NOT EXISTS Posts (
        post_id INTEGER PRIMARY KEY AUTOINCREMENT,
        thread_id INTEGER,
        user_id INTEGER,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (thread_id) REFERENCES Threads(thread_id),
        FOREIGN KEY (user_id) REFERENCES Users(user_id)
    );
    `
	if _, err := db.Exec(createPostsTable); err != nil {
		log.Fatal(err)
	}

	createCommentsTable := `
    CREATE TABLE IF NOT EXISTS Comments (
        comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER,
        user_id INTEGER,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES Posts(post_id),
        FOREIGN KEY (user_id) REFERENCES Users(user_id)
    );
    `
	if _, err := db.Exec(createCommentsTable); err != nil {
		log.Fatal(err)
	}

	createLikesTable := `
    CREATE TABLE IF NOT EXISTS Likes (
        like_id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        post_id INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES Users(user_id),
        FOREIGN KEY (post_id) REFERENCES Posts(post_id)
    );
    `
	if _, err := db.Exec(createLikesTable); err != nil {
		log.Fatal(err)
	}

	log.Println("Database crée")
}
