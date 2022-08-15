package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
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

func Register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "POST" {
		var userData UserData
		err := json.NewDecoder(r.Body).Decode(&userData) // unmarshall the userdata
		if err != nil {
			fmt.Print(err)
			http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Print(userData) // this is the data that need to be inserted to the database.
		w.Header().Set("Content-type", "application/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Register successful"))
		return
	}
	http.Error(w, "400 Bad Request.", http.StatusBadRequest)
}
