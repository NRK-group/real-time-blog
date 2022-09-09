package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
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
			http.Error(w, "403 Access Forbidden error", http.StatusForbidden)
			return

		} else {
			if strings.Contains(c.String(), "&") {
				co = strings.Split(c.Value, "&")
			}
			if !(forum.CheckSession(co[2])) {
				page = ReturnData{}
				marshallPage, err := json.Marshal(page)
				if err != nil {
					fmt.Println("Error marshalling the data: ", err)
				}
				w.Header().Set("Content-type", "application/text")
				w.WriteHeader(http.StatusOK)
				w.Write(marshallPage)
				return

			}
			chatusers, alluser, _ := forum.ArrangeUsers(co[0])
			page = ReturnData{User: forum.GetUser(co[0]), Posts: forum.AllPost("", co[0]), Msg: "Login successful", Users: alluser, ChatUsers: chatusers}
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
		DB.RegisterUser(userData)

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

			page = ReturnData{}
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

		chatusers, alluser, _ := forum.ArrangeUsers(userid)

		page = ReturnData{User: forum.GetUser(userid), Posts: forum.AllPost("", userid), Msg: "Login successful", Users: alluser, ChatUsers: chatusers}
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
			fmt.Println("removing this user from the database: ", res[1])
		}
		// Update the users status in the database
		forum.Update("User", "Status", "Offline", "userID", res[0])

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
	if r.Method == "GET" {
		page = ReturnData{Posts: forum.AllPost("", res[0]), Msg: "successful Post"}
		marshallPage, err := json.Marshal(page)
		if err != nil {
			fmt.Println("Error marshalling the data: ", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/json")
		w.Write(marshallPage)
		return
	}
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
			page = ReturnData{Posts: forum.AllPost("", res[0]), Msg: "successful Post"}
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
		if chatDetails.ChatID = forum.CheckChatID(chatDetails.User, chatDetails.Reciever); chatDetails.ChatID != "" {
			chatDetails.Messages = forum.TenMessages(chatDetails.ChatID, chatDetails.X)
			// When a message is sent check for the chat id a nd create it
			// chatDetails.ChatID = forum.CreateChatID(chatDetails.User, chatDetails.Reciever)
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

			page = ReturnData{Posts: forum.AllPost("", res[0]), Msg: "successful response--" + responseID}
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

func (forum *DB) Favorite(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/favorite" {
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
			var responseData Favorite
			err := json.NewDecoder(r.Body).Decode(&responseData)
			if err != nil {
				fmt.Print(err)
				http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			forum.CheckReactInPost(responseData.PostID, res[0], responseData.React)
			if err != nil {
				fmt.Print(err)
				http.Error(w, "500 Internal Server Error."+err.Error(), http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			page = ReturnData{Posts: forum.AllPost("", res[0]), Msg: "successful react to post--"}
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

func (forum *DB) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/updateuser" {
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
	if forum.CheckSession(res[2]) {
		if r.Method == "POST" {
			var userData UpdateUserData
			err := json.NewDecoder(r.Body).Decode(&userData)
			if err != nil {
				fmt.Print(err.Error())
				http.Error(w, "500 Internal Server Error.", http.StatusInternalServerError)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			page = ReturnData{Msg: forum.UpdateUserProfile(res[0], userData)}
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
	fmt.Println(users)
	userDetails := strings.Split(c.Value, "&")
	userIdVal := userDetails[0]
	users[userIdVal] = ws
	//Let all users know this users is online
	change := StateChange{Change: "status", Active: "Online", UserID: userDetails[0]}
	send, marshErr := json.Marshal(change)
	if marshErr != nil {
		fmt.Println("Problem marshalling change of status struct: ", marshErr)
	}
	for _, v := range users {
		v.WriteMessage(1, send)
	}
	fmt.Println(userIdVal, " is connected.")
	go forum.reader(ws)
}

func (forum *DB) Notifications(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Notify" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method == "PUT" {
		fmt.Println("Delete Info from notification table between:----------------------------------------- ")
		var chatUsers NewMessage

		err := json.NewDecoder(r.Body).Decode(&chatUsers)
		if err != nil {
			fmt.Println("Error unmarshalling the data to delete it from the database: ", err)
		}

		forum.DeleteNotification(chatUsers.UserID, chatUsers.RecieverID)

	}
	if r.Method == "GET" {
		// Get the number of notifications for the two users
		var getNotifs []Notify

		c, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "500 Internal error", http.StatusInternalServerError)
			return
		}

		username := strings.Split(c.Value, "&")[0]

		getNotifs = forum.GetNotifications(username)
		fmt.Println("getNotifs", getNotifs)
		values, marshErr := json.Marshal(getNotifs)
		if marshErr != nil {
			fmt.Println("Error marshalling notification results")
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-type", "application/json")
		w.Write(values)
	}
}
