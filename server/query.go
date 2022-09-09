package server

import (
	"fmt"
	"sort"
	"strconv"
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
			SessionID:   "sessionID",
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
			Favorite:     forum.GetFavoritesInPost(postID, uID),
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
	var commentID, postID, userID, date, time, imgUrl, content string
	for rows.Next() {
		rows.Scan(&commentID, &postID, &userID, &date, &time, &imgUrl, &content)
		comment = Comment{
			CommentID: commentID,
			PostID:    postID,
			UserID:    userID,
			Date:      date,
			Time:      time,
			ImgUrl:    imgUrl,
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

func (forum *DB) GetFavoritesInPost(pID, uID string) Favorite {
	rows, err := forum.DB.Query("SELECT favoriteID, postID, userID, react FROM Favorite WHERE postID = '" + pID + "' AND userID = '" + uID + "'")
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
	}
	return favorite
}

// ReactInPost
// is a method of database that add reaction in the post in it.
// Ex. Forum.Forum.ReactInPost("b081d711-aad2-4f90-acea-2f2842e28512", "b53124c2-39f0-4f10-8e02-b7244b406b86", 1)
func (forum *DB) ReactInPost(postID, userID string, react int) (string, error) {
	favoriteID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Favorite (favoriteID, postID, userID, react) values (?, ?, ?, ?)
	`)
	_, err := stmt.Exec(favoriteID, postID, userID, react)
	if err != nil {
		return "", err
	}
	return favoriteID.String(), nil
}

// CheckReactInPost
func (forum *DB) CheckReactInPost(pID, uID string, value int) (string, int) {
	rows, err := forum.DB.Query("SELECT favoriteID, postID, userID, react FROM Favorite WHERE postID = '" + pID + "' AND userID = '" + uID + "'")
	var reaction Favorite
	if err != nil {
		fmt.Println(err)
		return "", 0
	}
	var favoriteID, postID, userID string
	var react int
	for rows.Next() {
		rows.Scan(&favoriteID, &postID, &userID, &react)
		reaction = Favorite{
			FavoriteID: favoriteID,
			PostID:     postID,
			UserID:     userID,
			React:      react,
		}
	}
	if reaction.FavoriteID == "" {
		favoriteID, err := forum.ReactInPost(pID, uID, 1)
		fmt.Println(err)
		return favoriteID, 1
	}
	forum.Update("Favorite", "react", strconv.Itoa(value), "favoriteID", reaction.FavoriteID)
	return reaction.FavoriteID, value
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

func (forum *DB) CheckChatID(userID, recieverID string) (string, string) {
	// Query the DB and check if there is a chatId between the two users
	searchOne := ""
	searchOneDate := ""
	chatID, err := forum.DB.Query(`SELECT chatID, date from Chat WHERE user1ID = ? AND user2ID = ?`, userID, recieverID)
	if err != nil {
		fmt.Println("Error executing chatID search 1: ", err)
		return searchOne, ""
	}

	for chatID.Next() {
		chatID.Scan(&searchOne, &searchOneDate)
	}

	searchTwo := ""
	searchTwoDate := ""
	secondChatID, err2 := forum.DB.Query(`SELECT chatID, date from Chat WHERE user1ID = ? AND user2ID = ?`, recieverID, userID)
	if err2 != nil {
		fmt.Println("Error executing chatID search 2: ", err2)
		return searchTwo, ""
	}

	for secondChatID.Next() {
		secondChatID.Scan(&searchTwo, &searchTwoDate)
	}
	
	if searchOne != "" {
		return searchOne, searchOneDate
	}

	if searchTwo != "" {
		return searchTwo , searchTwoDate
	}

	return searchTwo , ""
}

func (forum *DB) CreateChatID(userID, recieverID string) string {
	// Create the chatID using uuid
	chatID := uuid.NewV4().String()

	insertChat, _ := forum.DB.Prepare(`
		INSERT INTO Chat (chatID, user1ID, user2ID, date) values (?, ?, ?, ?)
	`)
	_, err := insertChat.Exec(chatID, userID, recieverID, time.Now().Format("2006 January 02 15:04:05"))
	if err != nil {
		fmt.Println("Error inserting the chat id: ", err)
		return ""
	}
	fmt.Println("chatID added between user: ", userID, " AND user: ", recieverID)

	return chatID
}

func (forum *DB) InsertMessage(details NewMessage) {
	insertMessage, err1 := forum.DB.Prepare(`
	INSERT INTO Message (chatID, content, date, userID) VALUES (?,?,?,?)
	`)
	if err1 != nil {
		fmt.Println("Error Preparing message: ", err1)
		return
	}

	// messageID := uuid.NewV4().String()

	_, err := insertMessage.Exec(details.ChatID, details.Mesg, details.Date, details.UserID)

	forum.Update("Chat", "date", time.Now().Format("2006 January 02 15:04:05"), "chatID", details.ChatID)

	if err != nil {
		fmt.Println("Error inserting message: ", err)
	}
}

// CreateComment
// is a method of database that add comment in it.
func (forum *DB) CreateComment(userID, postID, content string) (string, error) {
	date := time.Now().Format("2006 January 02")
	time := time.Now()
	commentID := uuid.NewV4()
	stmt, _ := forum.DB.Prepare(`
		INSERT INTO Comment (commentID, postID, userID, date, time, imgUrl, content) values (?, ?, ?, ?, ?,?,?)
	`)
	_, err := stmt.Exec(commentID, postID, userID, date, time, "imgUrl", content)
	if err != nil {
		return "", err
	}
	return commentID.String(), nil
}

func (forum *DB) TenMessages(chatID string, x int) []SendMessage {
	// select the bottom 10 messages from the db
	getTen, err := forum.DB.Query(`SELECT content, date, userID FROM Message  WHERE chatID = ? ORDER BY messageID DESC LIMIT ?,10`, chatID, x)
	if err != nil {
		fmt.Println("Error selecting messages: ", err)
		return nil
	}
	result := make([]SendMessage, 0)
	// count := 1
	// start, _ := strconv.Atoi(x)
	for getTen.Next() {
		// if count > x {
		var current SendMessage
		getTen.Scan(&current.Message, &current.Date, &current.Sender)
		result = append(result, current)
		// }
		// count++
		// if count > x + 10{
		// 	break;
		// }
	}
	fmt.Println(result)
	return result
}

// ArrangeUsers
// is a method that will organized users by the last message sent, and return the all other users by alphabetic order.
func (forum *DB) ArrangeUsers(userID string) ([]User, []User, error) {
	userLastMessage := make([]User, 0)
	var userAlphabeticOrder []User
	allUsers := forum.GetAllUser(userID)

	for x := 0; x < len(allUsers); x++ {
		chatID, chatTime := forum.CheckChatID(userID, allUsers[x].UserID)
		if chatID != "" {
			allUsers[x].SessionID = chatTime
			userLastMessage = append(userLastMessage, allUsers[x])
		} else {
			userAlphabeticOrder = append([]User{allUsers[x]}, userAlphabeticOrder...)
		}
	}

	 sort.Slice(userLastMessage, func(i, j int) bool {
		one, _ := time.Parse("2006 January 02 15:04:05", userLastMessage[i].SessionID)
		two, _ := time.Parse("2006 January 02 15:04:05",userLastMessage[j].SessionID)
	 	return one.UnixNano() < two.UnixNano()
	 })

	sort.Slice(userAlphabeticOrder, func(i, j int) bool {
		return userAlphabeticOrder[i].Nickname > userAlphabeticOrder[j].Nickname
	})

	return userLastMessage, userAlphabeticOrder, nil
}

func (forum *DB) Notification(userID, recieverID string) {
	fmt.Println("Update the notifation msg", userID, recieverID)

	checkNotification, err := forum.DB.Query("SELECT 1 FROM MessageNotifications WHERE userID = ? AND recieverID = ?", userID, recieverID)
	if err != nil {
		fmt.Println("Error checking for notification row")
		return
	}
	count := 0
	for checkNotification.Next() {
		count++
	}
	if count == 0 {
		forum.CreateNotification(userID, recieverID)
	} else if count == 1 {
		forum.UpdateNotification(userID, recieverID)
	}
}

func (forum *DB) UpdateNotification(userID, recieverID string) {
	updateMsgs, err := forum.DB.Prepare(`UPDATE MessageNotifications SET number = number + 1 WHERE userID = ? AND recieverID = ?`)
	if err != nil {
		fmt.Println("Error Preparing update notification: ", err)
		return
	}

	_, errExec := updateMsgs.Exec(userID, recieverID)
	if errExec != nil {
		fmt.Println("Error executing update notifications: ", errExec)
		return
	}
}

func (forum *DB) CreateNotification(userID, recieverID string) {
	addNotification, err := forum.DB.Prepare(`INSERT INTO MessageNotifications VALUES (?,?,?)`)
	if err != nil {
		fmt.Println("Error preparing the insert statement: ", err)
		return
	}

	_, errExec := addNotification.Exec(userID, recieverID, 1)
	if errExec != nil {
		fmt.Println("Error inserting the notification")
		return
	}
}

func (forum *DB) DeleteNotification(sender, target string) {
	// Delete the row in the db
	_, deleteErr := forum.DB.Exec("DELETE FROM MessageNotifications WHERE userID = ? AND recieverID = ?", sender, target)
	if deleteErr != nil {
		fmt.Println("Error deleting from the Message Notification database")
	}
	fmt.Println("Successfully deleted notifications")
}

func (forum *DB) GetNotifications(target string) []Notify {
	getNotQry, err := forum.DB.Query("SELECT userID, number FROM MessageNotifications WHERE recieverID = ?", target)
	if err != nil {
		fmt.Println("Error querying for notification number: ", err)
	}

	result := make([]Notify, 0)
	for getNotQry.Next() {
		var temp Notify
		getNotQry.Scan(&temp.SenderID, &temp.Count)
		result = append(result, temp)
	}
	return result
}

func (forum *DB) UpdateUserProfile(userID string, userData UpdateUserData) string {
	rows, err := forum.DB.Query("SELECT password FROM User WHERE userID = '" + userID + "'")
	if err != nil {
		fmt.Print(err)
		return "User doesn't exist"
	}
	var password string
	for rows.Next() {
		rows.Scan(&password)
	}

	if !(CheckPasswordHash(userData.Password, password)) {
		return "Error - password not macth"
	}
	
	updateUser, err := forum.DB.Prepare(`UPDATE User SET firstName = ?, lastName = ?, nickName = ?, gender = ?, age = ?, password = ?, email = ?   WHERE userID = ?`)
	if err != nil {
		fmt.Println("Error Preparing update notification: ", err)
		return err.Error()
	}
	newPass, _ := HashPassword(userData.NewPassword)
	_, errExec := updateUser.Exec(userData.FirstName, userData.LastName, userData.Nickname, userData.Gender, userData.Age, newPass, userData.Email, userID)
	if errExec != nil {
		fmt.Println("Error executing update notifications: ", errExec)
		return errExec.Error()
	}

	return "Update is complete "
}

func (forum *DB) RegisterUser(userData UserData) {
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
	forum.DB.Exec(`INSERT INTO User VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`, userID, "", userData.FirstName, userData.LastName, userData.Nickname, userData.Gender, userAge, "Offline", userData.Email, userDate, password, "")
}

//ChangeStauts changes the status of a 