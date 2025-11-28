package database

import (
	"database/sql"
	"errors"
)

func (db *DbImpl) GetIdByUserName(name string) (int, error) {
	var idFound int
	err := db.c.QueryRow("SELECT id FROM users WHERE name = ?", name).Scan(&idFound)
	if errors.Is(err, sql.ErrNoRows) {
		res, err := db.c.Exec("INSERT INTO users (name) VALUES (?)", name)
		idFound, err := res.LastInsertId()
		return int(idFound), err
	}
	return idFound, err
}
