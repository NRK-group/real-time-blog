package server

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// start of the delete query

// Delete
// is a method of the database that delete value base on table and where.
//  ex. Forum.Delete("User", "userID", "185c6549-caec-4eae-95b0-e16023432ef0")
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

// Start of the Update query

// Update
//  ex. Forum.Update("User", "username", "Adriell", "userID" "7e2b4fdd-86ad-464c-a97e")
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
