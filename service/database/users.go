package database

import (
	"database/sql"
	"errors"
)

func (db *DbImpl) GetIdByUserName(name string) (int, error) {
	var id int
	err := db.c.QueryRow("SELECT id FROM users WHERE name = ?", name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		err = db.c.QueryRow("INSERT INTO users (name) VALUES (?) RETURNING id", name).Scan(&id)
		return int(id), err
	}
	return id, err
}

func (db *DbImpl) SetMyUserName(name string) error {

	_, err := db.c.Exec("UPDATE users SET name = ? WHERE name = ?;", name)
	return err
}
