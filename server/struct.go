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
	PostID       string
	UserID       string
	Date         string
	Time         string
	Title        string
	Content      string
	Category     string
	ImgUrl       string
	Comments     []Comment
	NumOfComment int
	Favorite     Favorite
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
	PostID     string `json:"postID"`
	UserID     string
	React      int `json:"react"`
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

type UpdateUserData struct {
	FirstName       string `json:"firstName"`
	Nickname        string `json:"nickname"`
	LastName        string `json:"lastName"`
	Age             string `json:"age"`
	Gender          string `json:"gender"`
	Email           string `json:"email"`
	Password        string `json:"oldPassword"`
	NewPassword        string `json:"password"`
	ConfirmPasspord string `json:"confirmPassword"`
}

type UserLoginData struct {
	EmailOrNickname string `json:"emailOrUsername"`
	Password        string `json:"password"`
}

type ReturnData struct {
	User      User
	Users     []User
	ChatUsers []User
	Posts     []Post
	Port      string
	Msg       string
}


type NewMessage struct {
	Mesg         string `json:"message"`
	UserID       string `json:"userID"`
	RecieverID   string `json:"recieverID"`
	ChatID       string `json:"chatID"`
	Date         string `json:"date"`
	Notification bool   `json:"notification"`
	messageType  int
}

type SendMessage struct {
	Sender  string `json:"senderID"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

type ReturnMessages struct {
	User     string `json:"userID"`
	Reciever string `json:"recieverID"`
	ChatID   string `json:"chatID"`
	X        int    `json:"X"`
	Messages []SendMessage
}

type PostData struct {
	Title    string `json:"postTitle"`
	Category string `json:"postCategory"`
	Content  string `json:"postContent"`
}

type ResponseData struct {
	PostID  string `json:"postID"`
	Content string `json:"responseContent"`
}

type CheckTyping struct {
	Typer    string `json:"typerID"`
	Reciever string `json:"recieverID"`
	Typing   string `json:"value"`
}

type MessageRequest struct {
	ChatID string `json:"chatID"`
	Rows   string `json:"rows"`
}

type Notify struct {
	SenderID string `json:"senderID"`
	Count    int    `json:"count"`
}
