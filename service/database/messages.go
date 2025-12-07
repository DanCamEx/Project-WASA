package database

// CreateMessage creates a new message in the specified conversation from the sender with the given content.
// It returns the id of the newly created message.
func (db *DbImpl) CreateMessage(senderId int, conversationId int, content string) (messageId int, err error) {

	err = db.c.QueryRow(`

		INSERT INTO messages (sender_id, conversation_id, content, status, isForwarded) 
		VALUES (?, ?, ?, 'sent', false) 
		RETURNING id

	`, senderId, conversationId, content).Scan(&messageId)
	return messageId, err
}

// GetMessage retrieves the message information for the given message id.
// It returns the sender id, conversation id, content, timestamp, status, and isForwarded flag.
func (db *DbImpl) GetMessage(messageId int) (senderId int, conversationId int, content string, timestamp string, status string, isForwarded bool, err error) {

	err = db.c.QueryRow(`

		SELECT sender_id, conversation_id, content, timestamp, status, isForwarded 
		FROM messages WHERE id = ?

	`, messageId).Scan( &senderId, &conversationId, &content, &timestamp, &status, &isForwarded)
	return senderId, conversationId, content, timestamp, status, isForwarded, err
}

// DeleteMessage deletes the message with the given id from the database.
func (db *DbImpl) DeleteMessage(messageId int) error {
	_, err := db.c.Exec(`

		DELETE FROM messages
		WHERE id = ?

	`, messageId)
	return err
}

// CheckMessage checks if a message with the given id exists and was sent by the specified sender.
func (db *DbImpl) CheckMessage(messageId int, senderId int) (exists bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 
			FROM messages 
			WHERE id = ? AND sender_id = ?)

	`, messageId, senderId).Scan(&exists)
	return exists, err
}

// CheckMessageOnly checks if a message with the given id exists in the database.
func (db *DbImpl) CheckMessageOnly(messageId int) (exists bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 
			FROM messages 
			WHERE id = ?)

	`, messageId).Scan(&exists)
	return exists, err
}

// ForwardMessage creates a new message with flag isForwarded set to true, copying the content from the original message.
// It returns the id of the newly created forwarded message.
func (db *DbImpl) ForwardMessage(messageIdIn int, senderId int, conversationId int) (messageIdOut int, err error) {
	var content string

	// Get the original message using another func
	_, _, content, _, _, _, err = db.GetMessage(messageIdIn)
	if err != nil { return messageIdOut, err }

	// Creates a new message using the original content but with new sender and flag updated
	err = db.c.QueryRow(`

		INSERT INTO messages (sender_id, conversation_id, content, status, isForwarded) 
		VALUES (?, ?, ?, 'sent', true) 
		RETURNING id

	`, senderId, conversationId, content).Scan(&messageIdOut)
	return messageIdOut, err
}