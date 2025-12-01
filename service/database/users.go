package database

import (
	"database/sql"
	"errors"
)

// GetIdByUserName retrieves the user id for the given user name.
// If the user name does not exist, it creates a new user and returns the new id.
func (db *DbImpl) GetIdByUserName(name string) (int, error) {
	var id int

	err := db.c.QueryRow(`
		SELECT id FROM users WHERE name = ?
	`, name).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		err = db.c.QueryRow(`
			INSERT INTO users (name) VALUES (?) RETURNING id
		`, name).Scan(&id)
		return id, err
	}

	return id, err
}

// SetUserName updates the user name for the given user id.
func (db *DbImpl) SetUserName(id int, name string) error {

	_, err := db.c.Exec(`
		UPDATE users SET name = ? WHERE id = ?
	`, name, id)

	return err
}

// GetUserInfoById retrieves the user information (id, name, image) for the given user id.
func (db *DbImpl) GetUserInfoById(id int) (int, string, string, error) {
	var name, image string

	err := db.c.QueryRow(`
		SELECT name, COALESCE(image, '') FROM users WHERE id = ?
	`, id).Scan(&name, &image)

	return id, name, image, err
}

// CheckUserName checks if a user name already exists in the database.
func (db *DbImpl) CheckUserName(name string) (bool, error) {
	var exists bool

	err := db.c.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)
	`, name).Scan(&exists)

	return exists, err
}

// GetAllUsers retrieves all users' names and images from the database.
func (db *DbImpl) GetAllUsers() ([][2]string, error) {
	var i int
	var info [][2]string

	rows, err := db.c.Query(`
		SELECT name, image FROM users
	`)
	if err != nil {
		return info, err
	}
	for rows.Next() {
		var name, image string

		if err = rows.Scan(&name, &image); err != nil {
			rows.Close()
			return info, err
		}
		info = append(info, [2]string{name, image})
		i++
	}
	if err = rows.Err(); err != nil {
		return info, err
	}

	rows.Close()
	return info, err
}

func (db *DbImpl) CheckUserId(userId int) (bool, error) {
	var exists bool

	err := db.c.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)
	`, userId).Scan(&exists)

	return exists, err
}
