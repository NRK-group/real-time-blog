package server

import (
	"database/sql"
)

func initUser(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "User" (
		"userID"	TEXT NOT NULL,
		"imgUrl"	TEXT NOT NULL,
		"firstName"	CHARACTER(20) NOT NULL,
		"lastName"	CHARACTER(20) NOT NULL,
		"nickName"	CHARACTER(20) UNIQUE NOT NULL,
		"gender"	CHARACTER(20) NOT NULL,
		"age"   int,
		"status"  CHARACTER(20) NOT NULL,
		"email"	TEXT UNIQUE NOT NULL,
		"dateCreated" TEXT NOT NULL,
		"password"	TEXT NOT NULL,
		"sessionID" TEXT,
		PRIMARY KEY("userID")
		FOREIGN KEY ("sessionID")
			REFERENCES "Session" ("sessionID")
		CHECK (length("nickName") >= 3 AND length("username") <= 20 )
		CHECK (("email") LIKE '%_@__%.__%')
		CHECk (length("password") >= 8)
	);
	`)
	stmt.Exec()
}

func initSession(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Session" (
		"sessionID"	TEXT UNIQUE NOT NULL,
		"dateAndTime" TEXT NOT NULL,
		"userID" TEXT NOT NULL,
		PRIMARY KEY("sessionID")
		FOREIGN KEY ("userID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}

func initPost(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Post" (
		"postID"	TEXT UNIQUE NOT NULL,
		"userID"	TEXT NOT NULL,
		"title"     TEXT NOT NULL,
		"category"	TEXT NOT NULL,
		"date"      TEXT NOT NULL,
		"time"      TEXT NOT NULL,
		"imgUrl"	TEXT NOT NULL,
		"content"	TEXT NOT NULL,
		PRIMARY KEY("postID")
		FOREIGN KEY ("userID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}

func initComment(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Comment" (
		"commentID" TEXT UNIQUE NOT NULL,
		"postID"	TEXT NOT NULL,
		"userID"	TEXT NOT NULL,
		"date"      TEXT NOT NULL,
		"time"      TEXT NOT NULL,
		"imgUrl"	TEXT NOT NULL,
		"content"	TEXT NOT NULL,
		PRIMARY KEY("commentID")
		FOREIGN KEY ("postID")
			REFERENCES "Post" ("postID")
		FOREIGN KEY ("userID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}

func initFavorite(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Favorite" (
		"favoriteID" TEXT UNIQUE NOT NULL,
		"postID"	TEXT NOT NULL,
		"userID"	TEXT NOT NULL,
		"react"	int,
		PRIMARY KEY("favoriteID")
		FOREIGN KEY ("postID")
			REFERENCES "Post" ("postID")
		FOREIGN KEY ("userID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}

func initMessage(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Message" (
		"messageID" TEXT UNIQUE NOT NULL,
		"chatID"	TEXT NOT NULL,
		"countent"	TEXT NOT NULL,
		"date"      TEXT NOT NULL,
		"time"      TEXT NOT NULL,
		"userID"	TEXT NOT NULL,
		PRIMARY KEY("messageID")
		FOREIGN KEY ("chatID")
			REFERENCES "Chat" ("chatID")
		FOREIGN KEY ("userID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}

func initChat(db *sql.DB) {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "Chat" (
		"chatID" TEXT UNIQUE NOT NULL,
		"user1ID"	TEXT NOT NULL,
		"user2ID"	TEXT NOT NULL,
		PRIMARY KEY("chatID")
		FOREIGN KEY ("user1ID")
			REFERENCES "User" ("userID")
		FOREIGN KEY ("user2ID")
			REFERENCES "User" ("userID")
	);
	`)
	stmt.Exec()
}

func CreateDatabase(db *sql.DB) *sql.DB {
	initUser(db)
	initSession(db)
	initPost(db)
	initComment(db)
	initFavorite(db)
	initMessage(db)
	initChat(db)
	return db
}
