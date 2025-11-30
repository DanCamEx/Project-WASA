package database

import (
	"database/sql"
	"errors"
)

// GetIDByUserName retrieves the user ID for the given username.
// If the username does not exist, it creates a new user and returns the new ID.
func (db *DbImpl) GetIDByUserName(name string) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM users WHERE name = ?", name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		err = db.c.QueryRow("INSERT INTO users (name, image) VALUES (?, '') RETURNING id", name).Scan(&id)
		return int(id), err
	}
	return id, err
}

// SetUserName updates the username for the user in the database.
func (db *DbImpl) SetUserName(id int, name string) error {
	_, err := db.c.Exec("UPDATE users SET name = ? WHERE id = ?;", name, id)
	return err
}

// GetUserInfoById retrieves the user information (ID, name, image) for the given user ID.
func (db *DbImpl) GetUserInfoByID(id int) (int, string, string, error) {
	var name, image string
	err := db.c.QueryRow("SELECT name, image FROM users WHERE id = ?", id).Scan(&name, &image)
	return id, name, image, err
}

func (db *DbImpl) CheckUserNameExists(name string) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = ?);", name).Scan(&exists)
	return exists, err
}
