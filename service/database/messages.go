package database

// CreateMessage creates a new message in the specified conversation from the sender with the given content.
// It returns the id of the newly created message.
func (db *DbImpl) CreateMessage(senderId int, conversationId int, content string) (id int, err error) {

	err = db.c.QueryRow(`

		INSERT INTO messages (sender_id, conversation_id, content, status, is_forwarded) 
		VALUES (?, ?, ?, 'sent', false) 
		RETURNING id

	`, senderId, conversationId, content).Scan(&id)
	return id, err
}

// GetMessage retrieves the message information for the given message id.
// It returns the sender id, conversation id, content, timestamp, status, and isForwarded flag.
func (db *DbImpl) GetMessage(idIn int) (idOut int, senderId int, conversationId int, content string, timestamp string, status string, isForwarded bool, err error) {

	err = db.c.QueryRow(`

		SELECT id, sender_id, conversation_id, content, timestamp, status, is_forwarded 
		FROM messages WHERE id = ?

	`, idIn).Scan(&idOut, &senderId, &conversationId, &content, &timestamp, &status, &isForwarded)
	return idOut, senderId, conversationId, content, timestamp, status, isForwarded, err
}

// DeleteMessage deletes the message with the given id from the database.
func (db *DbImpl) DeleteMessage(id int) error {
	_, err := db.c.Exec(`

		DELETE FROM messages
		WHERE id = ?

	`, id)
	return err
}

// CheckMessage checks if a message with the given id exists and was sent by the specified sender.
func (db *DbImpl) CheckMessage(id int, senderId int) (exists bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 
			FROM messages 
			WHERE id = ? AND sender_id = ?)

	`, id, senderId).Scan(&exists)
	return exists, err
}

// CheckMessageOnly checks if a message with the given id exists in the database.
func (db *DbImpl) CheckMessageOnly(id int) (exists bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 
			FROM messages 
			WHERE id = ?)

	`, id).Scan(&exists)
	return exists, err
}

// ForwardMessage creates a new message with flag isForwarded set to true, copying the content from the original message.
// It returns the id of the newly created forwarded message.
func (db *DbImpl) ForwardMessage(idIn int, senderId int, conversationId int) (idOut int, err error) {
	var content string

	// Get the original message using another func
	_, _, _, content, _, _, _, err = db.GetMessage(idIn)
	if err != nil { return idOut, err }

	// Creates a new message using the original content but with new sender and flag updated
	err = db.c.QueryRow(`

		INSERT INTO messages (sender_id, conversation_id, content, status, is_forwarded) 
		VALUES (?, ?, ?, 'sent', true) 
		RETURNING id

	`, senderId, conversationId, content).Scan(&idOut)
	return idOut, err
}