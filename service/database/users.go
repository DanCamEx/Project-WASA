package database

import "database/sql"

// GetIdByUserName retrieves the user id for the given user name.
// If the user name does not exist, it creates a new user and returns the new id.
func (db *DbImpl) GetIdByUserName(name string) (int, error) {
	var userId int
	if err := db.c.QueryRow(`
		INSERT INTO users (name) 
		VALUES (?) 
		ON CONFLICT(name) DO UPDATE 
			SET name = excluded.name
		RETURNING id
	`, name).Scan(&userId); err != nil {
		return 0, err
	}
	return userId, nil
}

// SetUserName updates the user name for the given user id.
// It returns an error if the user does not exist.
func (db *DbImpl) SetUserName(userId int, name string) error {
	res, err := db.c.Exec(`
		UPDATE users 
		SET name = ? 
		WHERE id = ?
	`, name, userId)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetUser retrieves the user information (name, image) for the given user id.
func (db *DbImpl) GetUser(userId int) (string, string, error) {
	var userName, userImage string
	if err := db.c.QueryRow(`
		SELECT name, image 
		FROM users 
		WHERE id = ?
	`, userId).Scan(&userName, &userImage); err != nil {
		return "", "", err
	}
	return userName, userImage, nil
}

// CheckUserName checks if a user name already exists in the database.
func (db *DbImpl) CheckUserName(name string) (bool, error) {
	var exists bool
	if err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM users 
			WHERE name = ?)`, name).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// A small struct used as output for GetAllUsers method.
type User struct {
	Name  string
	Image string
}

// GetAllUsers outputs all users existing in the database.
func (db *DbImpl) GetAllUsers() ([]User, error) {
	rows, err := db.c.Query(`
		SELECT name, image 
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	var name, image string
	for rows.Next() {
		if err := rows.Scan(&name, &image); err != nil {
			return nil, err
		}
		users = append(users, User{Name: name, Image: image})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// CheckUserId checks if a user with the given id exists in the database.
func (db *DbImpl) CheckUserId(userId int) (bool, error) {
	var exists bool
	if err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM users 
			WHERE id = ?)
	`, userId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
