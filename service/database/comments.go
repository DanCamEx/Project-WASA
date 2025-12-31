package database

import "database/sql"

// CreateComment creates a new comment (reaction) on a message by a user with the specified emoji.
// It returns the emoji of the created comment.
func (db *appdbimpl) CreateComment(messageId int, reacterId int, emojiIn string) (string, error) {
	var emojiOut string
	if err := db.c.QueryRow(`
		INSERT INTO comments (reacter_id, message_id, emoji) 
		VALUES (?, ?, ?) RETURNING emoji
	`, reacterId, messageId, emojiIn).Scan(&emojiOut); err != nil {
		return "", err
	}
	return emojiOut, nil
}

type Comment struct {
	ReacterId int    `json:"reacterId"`
	MessageId int    `json:"messageId"`
	Emoji     string `json:"emoji"`
}

// GetComment retrieves the comment (reaction) made by a user on a specific message.
// It returns the message id, reacter id, and emoji of the comment.
func (db *appdbimpl) GetComment(messageIdIn int, reacterIdIn int) (*Comment, error) {
	var comment Comment
	if err := db.c.QueryRow(`
		SELECT reacter_id, message_id, emoji 
		FROM comments 
		WHERE reacter_id = ? AND message_id = ?
	`, reacterIdIn, messageIdIn).Scan(&comment.ReacterId, &comment.MessageId, &comment.Emoji); err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteComment removes the comment (reaction) made by a user on a specific message.
func (db *appdbimpl) DeleteComment(messageId int, reacterId int) error {
	res, err := db.c.Exec(`
		DELETE FROM comments 
		WHERE reacter_id = ? AND message_id = ?
	`, reacterId, messageId)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// CheckComment checks if a comment (reaction) exists for a specific message by a user.
func (db *appdbimpl) CheckComment(messageId int, reacterId int) (bool, error) {
	var exists bool
	if err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM comments 
			WHERE reacter_id = ? AND message_id = ?)
	`, reacterId, messageId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// GetAllCommentsByMessageId retrieves all comments for a specific message.
func (db *appdbimpl) GetAllCommentsByMessageId(messageId int) (*[]Comment, error) {
	rows, err := db.c.Query(`
		SELECT reacter_id, message_id, emoji 
		FROM comments
		WHERE message_id = ?
	`, messageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ReacterId, &comment.MessageId, &comment.Emoji); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &comments, nil
}

func (db *appdbimpl) DeleteAllCommentsByMessageId(messageId int) error {
	_, err := db.c.Exec(`
		DELETE FROM comments 
		WHERE message_id = ?
	`, messageId)
	return err
}
