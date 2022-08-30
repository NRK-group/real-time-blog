package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

var page ReturnData

func (forum *DB) CheckCookie(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/vadidate" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "GET" {
		c, err := r.Cookie("session_token")
		co := []string{}

		if err != nil {
			http.Error(w, "500 Internal error", http.StatusInternalServerError)
			return

		} else {
			if strings.Contains(c.String(), "&") {
				co = strings.Split(c.Value, "&")
			}
			if !(forum.CheckSession(co[2])) {
				page = ReturnData{User: forum.GetUser(""), Posts: forum.AllPost("", ""), Msg: "", Users: forum.GetAllUser("")}
				marshallPage, err := json.Marshal(page)
				if err != nil {
					fmt.Println("Error marshalling the data: ", err)
				}
				w.Header().Set("Content-type", "application/text")
				w.WriteHeader(http.StatusOK)
				w.Write(marshallPage)
				return

			}

			page = ReturnData{User: forum.GetUser(co[0]), Posts: forum.AllPost("", ""), Msg: "Login successful", Users: forum.GetAllUser(co[0])}
			marshallPage, err := json.Marshal(page)
			if err != nil {
				fmt.Println("Error marshalling the data: ", err)
			}
			w.Header().Set("Content-type", "application/text")
			w.WriteHeader(http.StatusOK)
			w.Write(marshallPage)
			return
		}
	}
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
	/*

		var page ReturnData

		cookie, err := r.Cookie("session_token")

		if err != nil {
			page = ReturnData{User: User{}, Posts: forum.AllPost("", "")}
			if err := t.Execute(w, page); err != nil {
				http.Error(w, "500 Internal error", http.StatusInternalServerError)
				return
			}
		} else {

			co := forum.CheckCookie(w, cookie)

			page = ReturnData{User: forum.GetUser(co[0]), Posts: forum.AllPost("", "")}
		}
		if err := t.Execute(w, page); err != nil {
			http.Error(w, "500 Internal error", http.StatusInternalServerError)
			return
		}
	*/
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

		// Check if the nickname is already in use
		var allowNickname int
		DB.DB.QueryRow(`SELECT 1 from User WHERE nickName = (?);`, userData.Nickname).Scan(&allowNickname)

		var allowEmail int
		DB.DB.QueryRow(`SELECT 1 from User WHERE email = (?);`, userData.Email).Scan(&allowEmail)

		if allowNickname == 1 && allowEmail == 1 {
			// This user already exsists
			w.Header().Set("Content-type", "application/text")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This user already exists"))
			return
		} else if allowEmail == 1 {
			// This email is already in use
			w.Header().Set("Content-type", "application/text")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("This email is already in use"))
			return
		} else if allowNickname == 1 {
			// This nickname is already in use
			w.Header().Set("Content-type", "application/text")
			w.WriteHeader(http.StatusOK)
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

	http.Error(w, "400 Bad Request.", http.StatusBadRequest)
}

func (forum *DB) Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	SetupCorsResponse(w, r)

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

			page = ReturnData{User: forum.GetUser(""), Posts: forum.AllPost("", ""), Msg: loginResp, Users: forum.GetAllUser("")}
			marshallPage, err := json.Marshal(page)
			if err != nil {
				fmt.Println("Error marshalling the data: ", err)
			}
			w.Header().Set("Content-type", "application/text")
			w.WriteHeader(http.StatusOK)
			w.Write(marshallPage)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   loginResp,
			Expires: time.Now().Add(24 * time.Hour),
		})

		userid := strings.Split(loginResp, "&")[0]

		page = ReturnData{User: forum.GetUser(userid), Posts: forum.AllPost("", ""), Msg: "Login successful", Users: forum.GetAllUser(userid)}
		marshallPage, err := json.Marshal(page)
		if err != nil {
			fmt.Println("Error marshalling the data: ", err)
		}

		w.Header().Set("Content-type", "application/text")
		w.WriteHeader(http.StatusOK)
		w.Write(marshallPage)

		return
	}

	fmt.Println("Error in login handler")
	http.Error(w, "400 Bad Request.", http.StatusBadRequest)
}

func (forum *DB) Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res := strings.Split(c.Value, "&")
		err = forum.RemoveSession(res[2])
		if err != nil {
			log.Fatal(err)
		}

		// Remove the user from ws connection map
		if _, exsists := users[res[0]]; exsists {
			delete(users, res[0])
		}
		// Set the new token as the users `session_token` cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Expires: time.Now(),
		})
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/text")
		w.Write([]byte("Logout successful"))
		fmt.Println("Logout successful")
	default:
		http.Error(w, "400 Bad Request.", http.StatusBadRequest)
		return
	}
}

func (forum *DB) Post(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := strings.Split(c.Value, "&")

	if r.Method == "POST" {

		var postData PostData
		err := json.NewDecoder(r.Body).Decode(&postData)
		if err != nil {
			fmt.Print(err)
			http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if forum.CheckSession(res[2]) {

			postID, err := forum.CreatePost(res[0], postData.Title, postData.Category, "imgurl", postData.Content)
			fmt.Println(postID)
			fmt.Println(err)
			page = ReturnData{Posts: forum.AllPost("", ""), Msg: "successful Post"}
			marshallPage, err := json.Marshal(page)
			if err != nil {
				fmt.Println("Error marshalling the data: ", err)
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-type", "application/json")
			w.Write(marshallPage)
			return
		}

	}
	http.Error(w, "400 Bad Request.", http.StatusBadRequest)
}

func (forum *DB) GetMessages(w http.ResponseWriter, r *http.Request) {
	// Check the url is correct
	if r.URL.Path != "/MessageInfo" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	if r.Method == "POST" {
		// Unmarshal the data recieved
		var chatDetails ReturnMessages
		err := json.NewDecoder(r.Body).Decode(&chatDetails)
		if err != nil {
			fmt.Println("Error opening")
		}

		// Check if there is a chat between the two users
		if chatDetails.ChatID = forum.CheckChatID(chatDetails.User, chatDetails.Reciever); chatDetails.ChatID == "" {
			chatDetails.ChatID = forum.CreateChatID(chatDetails.User, chatDetails.Reciever)
		}

		marshallChat, marshErr := json.Marshal(chatDetails)
		if marshErr != nil {
			fmt.Println("Error marshalling getMessages: ", marshErr)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/json")
		w.Write(marshallChat)
	}
}

func (forum *DB) Response(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/response" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := strings.Split(c.Value, "&")

	if r.Method == "POST" {
		if forum.CheckSession(res[2]) {
			var responseData ResponseData
			err := json.NewDecoder(r.Body).Decode(&responseData)
			if err != nil {
				fmt.Print(err)
				http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			responseID, err := forum.CreateComment(res[0], responseData.PostID, responseData.Content)
			if err != nil {
				fmt.Print(err)
				http.Error(w, "500 Internal Server Error."+err.Error(), http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			page = ReturnData{Posts: forum.AllPost("", ""), Msg: "successful response--" + responseID}
			marshallPage, err := json.Marshal(page)
			if err != nil {
				fmt.Println("Error marshalling the data: ", err.Error())
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-type", "application/json")
			w.Write(marshallPage)
			return
		}
	}
	http.Error(w, "400 Bad Request.", http.StatusBadRequest)
}

func SetupCorsResponse(w http.ResponseWriter, req *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

func (forum *DB) WsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Problem upgrading", err)
		log.Println()
	}
	// Get the userId from the cookie value
	c, cookieErr := r.Cookie("session_token")
	if cookieErr != nil {
		fmt.Println("Error accessing cookie: ", cookieErr)
		return
	}

	// Store the new user in the Users map

	userIdVal := strings.Split(c.Value, "&")[0]
	users[userIdVal] = ws
	fmt.Println(userIdVal, " is connected.")
	go forum.reader(ws)
}
