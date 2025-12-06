package database

// CreateMessage creates a new message in the specified conversation from the sender with the given content.
// It returns the id of the newly created message.
func (db *DbImpl) CreateMessage(senderId int, conversationId int, content string) (int, error) {
	var messageId int

	err := db.c.QueryRow(`
		INSERT INTO messages (sender_id, conversation_id, content, status, isForwarded) 
		VALUES (?, ?, ?, 'sent', false) RETURNING id
	`, senderId, conversationId, content).Scan(&messageId)
	if err != nil {
		return messageId, err
	}

	return messageId, nil
}

// GetMessage retrieves the message information for the given message id.
// It returns the message id, sender id, conversation id, content, timestamp, status, and isForwarded flag.
func (db *DbImpl) GetMessage(id int) (int, int, int, string, string, string, bool, error) {
	var messageId, senderId, conversationId int
	var content, timestamp, status string
	var isForwarded bool

	err := db.c.QueryRow(`
		SELECT id, sender_id, conversation_id, content, timestamp, status, isForwarded 
		FROM messages WHERE id = ?
	`, id).Scan(&messageId, &senderId, &conversationId, &content, &timestamp, &status, &isForwarded)

	return messageId, senderId, conversationId, content, timestamp, status, isForwarded, err
}

func (db *DbImpl) DeleteMessage(id int) error {
	_, err := db.c.Exec(`
		DELETE FROM messages WHERE id = ?
	`, id)
	return err
}

func (db *DbImpl) CheckMessage(id int, senderId int) (bool, error) {
	var exists bool

	err := db.c.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM messages WHERE id = ? AND sender_id = ?)
	`, id, senderId).Scan(&exists)

	return exists, err
}