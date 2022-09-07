package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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
	// automate the user creation
	// for i := 0; i < 100; i++ {
	// 	userData := server.UserData{
	// 		FirstName: "Tan" + strconv.Itoa(i),
	// 		LastName:  "Tan" + strconv.Itoa(i),
	// 		Nickname:  "Tan" + strconv.Itoa(i),
	// 		Email:     "tan" + strconv.Itoa(i) + "@gmail.com",
	// 		Gender:    "Male",
	// 		Age:       "22",
	// 		Password:  "hello123",
	// 	}
	// 	database.RegisterUser(userData)
	// }

	go server.SendMsgs()
	http.HandleFunc("/", database.Home)
	http.HandleFunc("/register", database.Register)
	http.HandleFunc("/login", database.Login)
	http.HandleFunc("/ws", database.WsEndpoint)
	http.HandleFunc("/vadidate", database.CheckCookie)
	http.HandleFunc("/logout", database.Logout)
	http.HandleFunc("/post", database.Post)
	http.HandleFunc("/MessageInfo", database.GetMessages)
	http.HandleFunc("/Notify", database.Notifications)
	http.HandleFunc("/response", database.Response)
	http.HandleFunc("/favorite", database.Favorite)
	http.HandleFunc("/updateuser", database.UpdateUser)
	frontend := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", frontend)) // handling the CSS
	port := ""
	port = os.Getenv("PORT")
	if port == "" {
		port = "8800"
	}
	log.Print("Listening on 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
