package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	uuid "github.com/satori/go.uuid"
)

func (DB *DB) Home(w http.ResponseWriter, r *http.Request) {
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

func (DB *DB) Register(w http.ResponseWriter, r *http.Request) {
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

		// Check if the nickname is already in use
		var allowNickname int
		nickNameErr := DB.DB.QueryRow(`SELECT 1 from User WHERE nickName = (?);`, userData.Nickname).Scan(&allowNickname)
		if nickNameErr != nil {
			fmt.Println("Error checking the nickname: ", nickNameErr.Error())
		}

		var allowEmail int
		emailErr := DB.DB.QueryRow(`SELECT 1 from User WHERE email = (?);`, userData.Email).Scan(&allowEmail)
		if emailErr != nil {
			fmt.Println("Error checking the nickname: ", emailErr.Error())
		}

		if allowNickname == 1 && allowEmail == 1 {
			// This user already exsists
			w.Write([]byte("This user is already exists"))
			return
		} else if allowEmail == 1 {
			// This email is already in use
			w.Write([]byte("This email is already in use"))
			return
		} else if allowNickname == 1 {
			// This nickname is already in use
			w.Write([]byte("This Nickname is already in use"))
			return
		}

		//Create a UserId for the new user using UUID
		userID := uuid.NewV4().String()
		//Turn age into an int
		userAge, _ := strconv.Atoi(userData.Age)
		//Gate the date of registration
		userDate := time.Now().Format("01-02-2006")
		//Hash the password
		password, hashErr := HashPassword(userData.Password)

		if hashErr != nil {
			fmt.Println("Error hashing password", hashErr)
		}
		//Valid registration so add the user to the database
		DB.DB.Exec(`INSERT INTO User VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, userID, "", userData.FirstName, userData.LastName, userData.Nickname, userData.Gender, userAge, "Offline", userData.Email, userDate, password,"" )

		w.Header().Set("Content-type", "application/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Register successful"))
		return
	}
	http.Error(w, "400 Bad Request.", http.StatusBadRequest)
}

