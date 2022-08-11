package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum/server"

	_ "github.com/mattn/go-sqlite3"
)

func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		http.Error(w, "500 Internal error", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, ""); err != nil {
		http.Error(w, "500 Internal error", http.StatusInternalServerError)
		return
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./server/forum.db")
	if err != nil {
		log.Fatal("Database conection error")
	}

	 server.CreateDatabase(db)

	defer db.Close()

	

	http.HandleFunc("/", Home)
	frontend := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", frontend)) // handling the CSS
	fmt.Printf("Starting server at port 8800\n")
	log.Fatal(http.ListenAndServe(":8800", nil))
}
