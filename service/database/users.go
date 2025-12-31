package database

// CheckUserId checks if a user with the given id exists in the database.
func (db *appdbimpl) CheckUserId(userId int) (bool, error) {
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

// GetIdByUserName retrieves the user id for the given user name.
// If the user name does not exist, it creates a new user and returns the new id.
func (db *appdbimpl) GetIdByUserName(name string) (int, error) {
	var userId int
	if err := db.c.QueryRow(`
		INSERT INTO users (name) 
		VALUES (?) 
		ON CONFLICT(name) DO UPDATE 
		SET name=excluded.name
		RETURNING id
	`, name).Scan(&userId); err != nil {
		return 0, err
	}
	return userId, nil
}

// SetUserName updates user name for the specified user id.
func (db *appdbimpl) SetUserName(userId int, name string) error {
	_, err := db.c.Exec(`
		UPDATE users
		SET name = ?
		WHERE id = ?
	`, name, userId)
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns user info by user id.
func (db *appdbimpl) GetUser(userId int) (User, error) {
	var name, photo string
	if err := db.c.QueryRow(`
		SELECT name, photo 
		FROM users 
		WHERE id = ?
	`, userId).Scan(&name, &photo); err != nil {
		return User{}, err
	}
	return User{Id: userId, Name: name, Photo: photo}, nil
}

// GetUserByName returns user info by user name. If the user does not exist, a new one is created.
func (db *appdbimpl) GetUserByName(userName string) (User, error) {
	var userId int
	var userPhoto string
	if err := db.c.QueryRow(`
		INSERT INTO users (name) 
		VALUES (?) 
		ON CONFLICT(name) DO UPDATE 
		SET name=excluded.name
		RETURNING id, photo
	`, userName).Scan(&userId, &userPhoto); err != nil {
		return User{}, err
	}
	return User{Id: userId, Name: userName, Photo: userPhoto}, nil
}

// GetUserById returns user info by user id.
func (db *appdbimpl) GetUserById(userId int) (User, error) {
	var name, photo string
	if err := db.c.QueryRow(`
		SELECT name, photo 
		FROM users 
		WHERE id = ?
	`, userId).Scan(&name, &photo); err != nil {
		return User{}, err
	}
	return User{Id: userId, Name: name, Photo: photo}, nil
}

// CheckUserName checks if a user name already exists in the db.
func (db *appdbimpl) CheckUserName(name string) (bool, error) {
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

// GetAllUsers returns all users in the db.
func (db *appdbimpl) GetAllUsers() (*[]User, error) {
	rows, err := db.c.Query(`
		SELECT id, name, photo 
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	var id int
	var name, photo string
	for rows.Next() {
		if err := rows.Scan(&id, &name, &photo); err != nil {
			return nil, err
		}
		users = append(users, User{Id: id, Name: name, Photo: photo})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}

// setUserPhoto sets or updates the user photo by id.
func (db *appdbimpl) SetUserPhoto(userId int, filename string) error {
	_, err := db.c.Exec(`
		UPDATE users 
		SET photo = ? 
		WHERE id = ?
	`, filename, userId)
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}
