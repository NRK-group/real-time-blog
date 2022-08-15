package server

import (
	"database/sql"
)

type DB struct {
	DB *sql.DB
}

type User struct {
	UserID      string
	SessionID   string
	Firstname   string
	Lastname    string
	Age         int
	Nickname    string
	Gender      string
	Status      string
	ImgUrl      string
	Email       string
	DateCreated string
	Password    string
}
type Session struct {
	SessionID   string
	UserID      string
	DateAndTime string
}
type Post struct {
	PostID   string
	UserID   string
	Date     string
	Time     string
	Title    string
	Content  string
	Category string
	ImgUrl   string
	Comments []Comment
	Favorite Favorite
}

type Comment struct {
	CommentID string
	UserID    string
	PostID    string
	Date      string
	Time      string
	ImgUrl    string
	Content   string
}

type Favorite struct {
	FavoriteID string
	PostID     string
	UserID     string
	React      int
}

type Chat struct {
	ChatID  string
	User1ID string
	User2ID string
}

type Message struct {
	MessageID string
	ChatID    string
	Content   string
	Date      string
	Time      string
	UsedID    string
}
type UserData struct {
	FirstName       string `json:"firstName"`
	Nickname        string `json:"nickname"`
	LastName        string `json:"lastName"`
	Age             string `json:"age"`
	Gender          string `string:"gender"`
	Email           string `string:"email"`
	Password        string `string:"password"`
	ConfirmPasspord string `string:"confirmPassword"`
}
