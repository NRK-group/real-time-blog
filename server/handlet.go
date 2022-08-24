package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func (forum *DB) CheckCookie(w http.ResponseWriter, c *http.Cookie) []string {
	co := []string{}
	if strings.Contains(c.String(), "-") {
		co = strings.Split(c.Value, "-")
	}
	if len(co) != 0 {
		if !(forum.CheckSession(co[2])) {
			// Set the new token as the users `session_token` cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   "",
				Expires: time.Now(),
			})
		} else {
			return co
		}
	}
	return co
}

func (forum *DB) Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/index.html")
	if err != nil {
		http.Error(w, "500 Internal error", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, ""); err != nil {
		http.Error(w, "500 Internal error", http.StatusInternalServerError)
		return
	}

	var page ReturnData

	cookie, err := r.Cookie("session_token")

	if err != nil {
		page = ReturnData{User:User{} , Posts: forum.AllPost("", ""), }
		if err := t.Execute(w, page); err != nil {
			http.Error(w, "500 Internal error", http.StatusInternalServerError)
			return
		}
	} else {

		co := forum.CheckCookie(w, cookie)

		page = ReturnData{User: forum.GetUser(co[0]) , Posts: forum.AllPost("", "")}
	}
	if err := t.Execute(w, page); err != nil {
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
		DB.DB.QueryRow(`SELECT 1 from User WHERE nickName = (?);`, userData.Nickname).Scan(&allowNickname)

		var allowEmail int
		DB.DB.QueryRow(`SELECT 1 from User WHERE email = (?);`, userData.Email).Scan(&allowEmail)

		if allowNickname == 1 && allowEmail == 1 {
			// This user already exsists
			w.Write([]byte("This user already exists"))
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

		// Create a UserId for the new user using UUID
		userID := uuid.NewV4().String()
		// Turn age into an int
		userAge, _ := strconv.Atoi(userData.Age)
		// Gate the date of registration
		userDate := time.Now().Format("January 2, 2006")
		// Hash the password
		password, hashErr := HashPassword(userData.Password)

		if hashErr != nil {
			fmt.Println("Error hashing password", hashErr)
		}
		// Valid registration so add the user to the database
		DB.DB.Exec(`INSERT INTO User VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, userID, "", userData.FirstName, userData.LastName, userData.Nickname, userData.Gender, userAge, "Offline", userData.Email, userDate, password, "")

		w.Header().Set("Content-type", "application/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Register successful"))
		return
	}
	fmt.Println("Error in register handler")

	http.Error(w, "400 Bad Request.", http.StatusBadRequest)
}

func (forum *DB) Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "POST" {

		var userLoginData UserLoginData
		err := json.NewDecoder(r.Body).Decode(&userLoginData) // unmarshall the userdata
		if err != nil {
			fmt.Print(err)
			http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		loginResp := forum.LoginUsers(userLoginData.EmailOrNickname, userLoginData.Password)
		if loginResp[0] == 'E' {
			w.Header().Set("Content-type", "application/text")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(loginResp))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   loginResp,
			Expires: time.Now().Add(24 * time.Hour),
		})

		w.Header().Set("Content-type", "application/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Login successful"))
		return
	}
	fmt.Println("Error in login handler")
	http.Error(w, "400 Bad Request.", http.StatusBadRequest)
}
