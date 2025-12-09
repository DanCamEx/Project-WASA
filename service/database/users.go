package database

import (
	"database/sql"
	"errors"
)

// GetIdByUserName retrieves the user id for the given user name.
// If the user name does not exist, it creates a new user and returns the new id.
func (db *DbImpl) GetIdByUserName(name string) (userId int, err error) {

	err = db.c.QueryRow(`

		SELECT id 
		FROM users 
		WHERE name = ?

	`, name).Scan(&userId)
	if errors.Is(err, sql.ErrNoRows) {
		err = db.c.QueryRow(`

			INSERT INTO users (name) 
			VALUES (?) 
			RETURNING id

		`, name).Scan(&userId)
	}
	return userId, err
}

// SetUserName updates the user name for the given user id.
func (db *DbImpl) SetUserName(userId int, name string) error {
	
	_, err := db.c.Exec(`

		UPDATE users 
		SET name = ? 
		WHERE id = ?

	`, name, userId)
	return err
}

// GetUser retrieves the user information (name, image) for the given user id.
func (db *DbImpl) GetUser(userId int) (userName string, userImage string, err error) {

	err = db.c.QueryRow(`

		SELECT name, image 
		FROM users 
		WHERE id = ?
		
	`, userId).Scan(&userName, &userImage)
	return userName, userImage, err
}


// CheckUserName checks if a user name already exists in the database.
func (db *DbImpl) CheckUserName(name string) (exists bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 
			FROM users 
			WHERE name = ?)

	`, name).Scan(&exists)
	return exists, err
}

// GetAllUsers retrieves all users' names and images from the database.
func (db *DbImpl) GetAllUsers() (users [][]string, err error) {

	// Get all users table rows
	rows, err := db.c.Query(`

		SELECT name, image 
		FROM users

	`)
	if err != nil { return nil, err }

	// For each row scan the name and image and append to the users list
	var name, image string
	for rows.Next() {
		err := rows.Scan(&name, &image)
		if err != nil { return nil, err }

		users = append(users, []string{name, image})
	}

	return users, nil
}

// CheckUserId checks if a user with the given id exists in the database.
func (db *DbImpl) CheckUserId(userId int) (exists bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 
			FROM users 
			WHERE id = ?)

	`, userId).Scan(&exists)
	return exists, err
}