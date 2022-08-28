package server

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Delete
// is a method of the database that delete value base on table and where.
//
//	ex. Forum.Delete("User", "userID", "185c6549-caec-4eae-95b0-e16023432ef0")
func (forum *DB) Delete(table, where, value string) error {
	dlt := "DELETE FROM " + table + " WHERE " + where
	stmt, err := forum.DB.Prepare(dlt + " = (?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(value)
	if err != nil {
		return err
	}
	return nil
}

// RemoveSession
// is method of forum that removes the session based on the sessionID
func (forum *DB) RemoveSession(sessionID string) error {
	err := forum.Update("User", "sessionID", "", "sessionID", sessionID)
	if err != nil {
		return err
	}
	err = forum.Delete("Session", "sessionID", sessionID)
	if err != nil {
		return err
	}
	return nil
}

// Update
// is a method of database that update a row.
//
//	ex. Forum.Update("User", "username", "Adriell,", "userID" "7e2b4fdd-86ad-464c-a97e")
func (forum *DB) Update(table, set, to, where, id string) error {
	update := "UPDATE " + table + " SET " + set + " = '" + to + "' WHERE " + where + " = '" + id + "'"
	stmt, _ := forum.DB.Prepare(update)
	_, err := stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

// CreateSession
// is a method of database that add session in it based on the user login time.
func (forum *DB) CreateSession(userID string) (string, error) {
	date := time.Now().Format("2006 January 02 15:04:05")
	sessionID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Session (sessionID, userID, dateAndTime) values (?, ?, ?)
	`)
	_, err := stmt.Exec(sessionID, userID, date)
	if err != nil {
		return "", err
	}
	forum.Update("User", "sessionID", sessionID.String(), "userID", userID)
	return sessionID.String(), nil
}

// LoginUser
// is method of forum that checks the database if the login details match the credential
// and allow them to login if their is a match credentials
func (forum *DB) LoginUsers(emailOrNickname, pas string) string {
	var users User
	rows, err := forum.DB.Query(`SELECT * FROM User WHERE nickName = (?) OR email = (?);`, emailOrNickname, emailOrNickname)
	if err != nil {
		return err.Error()
	}
	var userID, imgUrl, firstName, lastName, nickName, gender, status, email, dateCreated, pass, sessionID string
	var age int
	for rows.Next() {
		rows.Scan(&userID, &imgUrl, &firstName, &lastName, &nickName, &gender, &age, &status, &email, &dateCreated, &pass, &sessionID)
		users = User{
			UserID:      userID,
			SessionID:   sessionID,
			Firstname:   firstName,
			Lastname:    lastName,
			Age:         age,
			Nickname:    nickName,
			Gender:      gender,
			Status:      status,
			ImgUrl:      imgUrl,
			Email:       email,
			DateCreated: dateCreated,
			Password:    pass,
		}
	}

	if users.Nickname == "" || users.Email == "" {
		return "Error - user not found"
	}
	if !(CheckPasswordHash(pas, users.Password)) {
		return "Error - password not macth"
	}
	if users.SessionID != "" {
		forum.RemoveSession(users.SessionID)
	}
	sess, err := forum.CreateSession(users.UserID)
	if err != nil {
		return err.Error()
	}
	users.SessionID = sess
	forum.Update("User", "Status", "Online", "userID", userID)
	return users.UserID + "&" + users.Nickname + "&" + users.SessionID
}

// CheckSession
// is a method of forum that checks the session in the User table and checks if it is match with the sessionID inputed
func (forum *DB) CheckSession(sessionId string) bool {
	rows, err := forum.DB.Query("SELECT sessionID FROM User WHERE sessionID = '" + sessionId + "'")
	if err != nil {
		fmt.Print(err)
		return false
	}
	var session string
	for rows.Next() {
		rows.Scan(&session)
	}
	return session != ""
}

// GetUser
// is a methond that return the user info by their user id

func (forum *DB) GetUser(uID string) User {
	var user User
	rows, err := forum.DB.Query("SELECT * FROM User WHERE userID = '" + uID + "'")
	if err != nil {
		fmt.Println(err)
		return user
	}

	var userID, imgUrl, firstName, lastName, nickName, gender, status, email, dateCreated, pass, sessionID string
	var age int
	for rows.Next() {
		rows.Scan(&userID, &imgUrl, &firstName, &lastName, &nickName, &gender, &age, &status, &email, &dateCreated, &pass, &sessionID)
		user = User{
			UserID:      userID,
			SessionID:   sessionID,
			Firstname:   firstName,
			Lastname:    lastName,
			Age:         age,
			Nickname:    nickName,
			Gender:      gender,
			Status:      status,
			ImgUrl:      imgUrl,
			Email:       email,
			DateCreated: dateCreated,
		}
	}
	forum.Update("User", "Status", "Online", "userID", userID)
	return user
}

// GetAllUser
// is a methond that return alluser nickname

func (forum *DB) GetAllUser(uID string) []User {
	var user User
	var users []User
	if uID == "" {
		return users
	}
	rows, err := forum.DB.Query("SELECT * FROM User WHERE NOT userID = '" + uID + "'")
	if err != nil {
		fmt.Println(err)
		return users
	}

	var userID, imgUrl, firstName, lastName, nickName, gender, status, email, dateCreated, pass, sessionID string
	var age int
	for rows.Next() {
		rows.Scan(&userID, &imgUrl, &firstName, &lastName, &nickName, &gender, &age, &status, &email, &dateCreated, &pass, &sessionID)
		user = User{
			UserID:      userID,
			SessionID:   sessionID,
			Firstname:   firstName,
			Lastname:    lastName,
			Age:         age,
			Nickname:    nickName,
			Gender:      gender,
			Status:      status,
			ImgUrl:      imgUrl,
			Email:       email,
			DateCreated: dateCreated,
		}
		users = append([]User{user}, users...)
	}

	return users
}

// AllPost
// is a method of forum that will return all post
func (forum *DB) AllPost(filter, uID string) []Post {
	var post Post
	var posts []Post
	var err error

	rows, err := forum.DB.Query("SELECT * FROM Post ")
	if err != nil {
		fmt.Print(err)
		return posts
	}

	var postID, userID, title, category, date, time, content, imgurl string
	for rows.Next() {
		rows.Scan(&postID, &userID, &title, &category, &date, &time, &imgurl, &content)
		post = Post{
			PostID:       postID,
			UserID:       userID,
			Date:         date,
			Time:         time,
			Content:      content,
			Category:     category,
			Title:        title,
			ImgUrl:       imgurl,
			Comments:     forum.GetComments(postID),
			NumOfComment: len(forum.GetComments(postID)),
			Favorite:     forum.GetFavoritesInPost(postID),
		}

		var username string
		rows2, err := forum.DB.Query("SELECT nickName FROM User WHERE userID = '" + userID + "'")
		if err != nil {
			fmt.Print(err)
			return posts
		}
		for rows2.Next() {
			rows2.Scan(&username)
		}
		post.UserID = username
		switch filter {
		case "go":
			if strings.Contains(category, "go") {
				posts = append([]Post{post}, posts...)
			}
		case "javascript":
			if strings.Contains(category, "javascript") {
				posts = append([]Post{post}, posts...)
			}
		case "rust":
			if strings.Contains(category, "rust") {
				posts = append([]Post{post}, posts...)
			}
		default:
			posts = append([]Post{post}, posts...)
		}
	}

	return posts
}

// Get Comments
// is a method of forum that return all the comment with that specific postID
func (forum *DB) GetComments(pID string) []Comment {
	rows, err := forum.DB.Query("SELECT * FROM Comment WHERE postID = '" + pID + "'")
	var comment Comment
	var comments []Comment
	if err != nil {
		fmt.Print(err)
		return comments
	}
	var commentID, postID, userID, date, time, content string
	for rows.Next() {
		rows.Scan(&commentID, &postID, &userID, &date, &time, &content)
		comment = Comment{
			CommentID: commentID,
			PostID:    postID,
			UserID:    userID,
			Date:      date,
			Time:      time,
			Content:   content,
		}
		var username string
		rows2, err := forum.DB.Query("SELECT nickName FROM User WHERE userID = '" + userID + "'")
		if err != nil {
			fmt.Print(err)
			return comments
		}
		for rows2.Next() {
			rows2.Scan(&username)
		}
		comment.UserID = username
		comments = append([]Comment{comment}, comments...)
	}
	return comments
}

func (forum *DB) GetFavoritesInPost(pID string) Favorite {
	rows, err := forum.DB.Query("SELECT favoriteID, postID, userID, react FROM Favorite WHERE postID = '" + pID + "'")
	var favorite Favorite

	if err != nil {
		fmt.Print(err)
		return favorite
	}
	var favoriteID, postID, userID string
	var react int
	for rows.Next() {
		rows.Scan(&favoriteID, &postID, &userID, &react)
		favorite = Favorite{
			FavoriteID: favoriteID,
			PostID:     postID,
			UserID:     userID,
			React:      react,
		}
		if react != 1 {
			react = 1
		} else {
			react = 0
		}

		favorite.React = react
	}
	return favorite
}

// CreatePost
// is a method of database that add post in it.
func (forum *DB) CreatePost(userID, title, category, imgurl, content string) (string, error) {
	date := time.Now().Format("2006 January 02")
	time := time.Now()
	postID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Post (postID, userID, title, category, date, time, imgurl, content) values (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	_, err := stmt.Exec(postID, userID, title, category, date, time, imgurl, content)
	if err != nil {
		return "", err
	}
	return postID.String(), nil
}

func (forum *DB) CheckChatID(userID, recieverID string) string {
	// Query the DB and check if there is a chatId between the two users
	searchOne := ""
	chatID, err := forum.DB.Query(`SELECT chatID from Chat WHERE user1ID = ? AND user2ID = ?`, userID, recieverID)
	if err != nil {
		fmt.Println("Error executing chatID search 1: ", err)
		return searchOne
	}

	for chatID.Next() {
		chatID.Scan(&searchOne)
	}

	searchTwo := ""
	secondChatID, err2 := forum.DB.Query(`SELECT chatID from Chat WHERE user1ID = ? AND user2ID = ?`, recieverID, userID)
	if err2 != nil {
		fmt.Println("Error executing chatID search 2: ", err2)
		return searchTwo
	}

	for secondChatID.Next() {
		secondChatID.Scan(&searchTwo)
	}
	fmt.Println("search One === ", searchOne, "search Two === ", searchTwo)
	if searchOne != "" {
		return searchOne
	}
	return searchTwo
}
