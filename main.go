package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum/server"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./server/forum.db")
	if err != nil {
		log.Fatal("Database conection error")
	}
	database := &server.DB{
		DB: server.CreateDatabase(db),
	}

	defer db.Close()

	http.HandleFunc("/", database.Home)
	http.HandleFunc("/register", database.Register)
	http.HandleFunc("/login", database.Login)
	http.HandleFunc("/vadidate", database.CheckCookie)
	frontend := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", frontend)) // handling the CSS
	fmt.Printf("Starting server at port 8800\n")
	log.Fatal(http.ListenAndServe(":8800", nil))
}
