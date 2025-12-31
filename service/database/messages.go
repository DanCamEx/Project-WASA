package database

import "database/sql"

// CreateMessage creates a new message in the specified conversation from the sender with the given content.
// It returns the id of the newly created message.
func (db *appdbimpl) CreateMessage(senderId int, conversationId int, content string) (int, error) {
	var id int
	if err := db.c.QueryRow(`
		INSERT INTO messages (sender_id, conversation_id, content, status, is_forwarded) 
		VALUES (?, ?, ?, 'sent', false) 
		RETURNING id
	`, senderId, conversationId, content).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// GetMessage retrieves the message information for the given message id.
// It returns the sender id, content, timestamp, status, and isForwarded flag.
func (db *appdbimpl) GetMessage(id int) (*Message, error) {
	var message Message
	if err := db.c.QueryRow(`
		SELECT id, sender_id, conversation_id, content, timestamp, status, is_forwarded 
		FROM messages WHERE id = ?
	`, id).Scan(
		&message.Id,
		&message.SenderId,
		&message.ConversationId,
		&message.Content,
		&message.Timestamp,
		&message.Status,
		&message.IsForwarded); err != nil {
		return nil, err
	}
	return &message, nil
}

// DeleteMessage deletes the message with the given id from the database.
func (db *appdbimpl) DeleteMessage(id int) error {
	res, err := db.c.Exec(`
		DELETE FROM messages
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// CheckMessage checks if a message with the given id exists and was sent by the specified sender.
func (db *appdbimpl) CheckMessage(id int, senderId int) (bool, error) {
	var exists bool
	if err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM messages 
			WHERE id = ? AND sender_id = ?)
	`, id, senderId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// CheckMessageOnly checks if a message with the given id exists in the database.
func (db *appdbimpl) CheckMessageOnly(id int) (bool, error) {
	var exists bool
	if err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM messages 
			WHERE id = ?)
	`, id).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// ForwardMessage creates a new message with flag isForwarded set to true, copying the content from the original message.
// It returns the id of the newly created forwarded message.
func (db *appdbimpl) ForwardMessage(idIn int, senderId int, conversationId int) (int, error) {
	var content string
	var message *Message
	message, err := db.GetMessage(idIn)
	if err != nil {
		return 0, err
	}
	content = message.Content
	var idOut int
	if err := db.c.QueryRow(`
		INSERT INTO messages (sender_id, conversation_id, content, status, is_forwarded) 
		VALUES (?, ?, ?, 'sent', true) 
		RETURNING id
	`, senderId, conversationId, content).Scan(&idOut); err != nil {
		return 0, err
	}
	return idOut, nil
}

// RemoveStatus removes the status entry for the given message and user.
func (db *appdbimpl) RemoveStatus(messageId int, userId int) error {
	_, err := db.c.Exec(`
		DELETE FROM statuses 
		WHERE message_id = ? AND user_id = ?
	`, messageId, userId)
	if err != nil {
		return err
	}
	return nil
}

// AddStatus adds a new status entry for the given message and user with status 'sent'.
func (db *appdbimpl) AddStatus(messageId int, userId int) error {
	if _, err := db.c.Exec(`
		INSERT INTO statuses (message_id, user_id, status) 
		VALUES (?, ?, 'sent') 
	`, messageId, userId); err != nil {
		return err
	}
	return nil
}

// UpdateStatus updates the status of the message with the given id to 'read'.
func (db *appdbimpl) UpdateStatus(id int) error {
	_, err := db.c.Exec(`
		UPDATE messages 
		SET status = 'read' 
		WHERE id = ?
		`, id)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTracking updates the message status to 'read' for the given user and message id.
func (db *appdbimpl) UpdateTracking(userId int, messageId int) error {
	if _, err := db.c.Exec(`
		UPDATE statuses
		SET status = ?
		WHERE user_id = ?
		AND message_id = ?
	`, "read", userId, messageId); err != nil {
		return err
	}
	return nil
}

// IsMessageReadByAll checks if all participants in the conversation have read the specified message.
func (db *appdbimpl) IsMessageReadByAll(messageId int, conversationId int) (bool, error) {
	var totalParticipants int
	err := db.c.QueryRow(`
        SELECT COUNT(*) 
		FROM participants 
		WHERE conversation_id = ?
    `, conversationId).Scan(&totalParticipants)
	if err != nil {
		return false, err
	}

	var readCount int
	err = db.c.QueryRow(`
        SELECT COUNT(*) 
		FROM statuses 
		WHERE message_id = ? 
		AND status = 'read'
    `, messageId).Scan(&readCount)
	if err != nil {
		return false, err
	}

	return readCount == totalParticipants-1, nil
}

type Message struct {
	Id             int    `json:"id"`
	SenderId       int    `json:"senderId"`
	ConversationId int    `json:"conversationId"`
	Content        string `json:"content"`
	Timestamp      string `json:"timestamp"`
	Status         string `json:"status"`
	IsForwarded    bool   `json:"isForwarded"`
}
