/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database is handled here. Database-specific logic should never escape this package.
To use this package, you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high-level interface for the DB
type AppDatabase interface {

	// User-related methods
	GetIdByUserName(name string) (int, error)
	SetUserName(id int, name string) error
	GetUserInfoById(id int) (int, string, string, error)
	CheckUserName(name string) (bool, error)
	GetAllUsers() ([][]string, error)
	CheckUserId(userId int) (bool, error)

	// Conversation-related methods
	CreateDirConversation(creatorId int, participantId int) (int, error)
	CreateGrpConversation(creatorId int, participantIds []int) (int, error)
	SetGroupName(conversationId int, name string) error
	CheckConversation(conversationId int, userId int) (bool, error)
	GetConversation(conversationId int) (string, string, string, int, error)
	CheckIfGroup(conversationId int) (bool, error)
	GetAllUserConversations(userId int) ([][]string, error)

	// Message-related methods
	CreateMessage(senderId int, conversationId int, content string) (int, error)
	GetMessage(id int) (int, int, int, string, string, string, bool, error)
	DeleteMessage(id int) error
	CheckMessage(id int, senderId int) (bool, error)
}

type DbImpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Creates the database structure if it doesn't already exist.
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL DEFAULT '' UNIQUE,
			image TEXT NOT NULL DEFAULT ''
		);
		CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL DEFAULT '',
			image TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT '',
			latest_message_id INTEGER NOT NULL DEFAULT 0,
		    FOREIGN KEY(latest_message_id) REFERENCES messages(id)
		);
		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY,
			sender_id INTEGER NOT NULL DEFAULT 0,
			conversation_id INTEGER NOT NULL DEFAULT 0,
			content TEXT NOT NULL DEFAULT '',
			timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			status TEXT NOT NULL DEFAULT '',
			isForwarded BOOLEAN NOT NULL DEFAULT 0,
			FOREIGN KEY(conversation_id) REFERENCES conversations(id),
			FOREIGN KEY(sender_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS comments (
			reacter_id INTEGER NOT NULL DEFAULT 0,
			message_id INTEGER NOT NULL DEFAULT 0,
			emoji TEXT NOT NULL DEFAULT '',
			PRIMARY KEY (reacter_id, message_id),
			FOREIGN KEY(reacter_id) REFERENCES users(id),
			FOREIGN KEY(message_id) REFERENCES messages(id)
		);
		CREATE TABLE IF NOT EXISTS participants (
			user_id INTEGER NOT NULL DEFAULT 0,
			conversation_id INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (user_id, conversation_id),
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(conversation_id) REFERENCES conversations(id)
		)`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}
	return &DbImpl{
		c: db,
	}, nil
}
