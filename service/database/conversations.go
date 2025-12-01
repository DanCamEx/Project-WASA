package database

// CreateDirConversation creates a new conversation between two users and returns its id.
// The conversation type is set to 'direct'.
// Both users are added as participants of the conversation in a separate table.
func (db *DbImpl) CreateDirConversation(creatorId, participantId int) (int, error) {
	var conversationId int

	err := db.c.QueryRow(`
		INSERT INTO conversations (type) VALUES ('direct') RETURNING id
	`).Scan(&conversationId)
	if err != nil {
		return conversationId, err
	}
	_, err = db.c.Exec(`
		INSERT INTO participants (user_id, conversation_id) VALUES (?, ?), (?, ?)
	`, creatorId, conversationId, participantId, conversationId)

	return conversationId, err
}

// CheckConversation checks if a conversation with the given id exists.
func (db *DbImpl) CheckConversation(conversationId int) (bool, error) {
	var exists bool

	err := db.c.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM conversations WHERE id = ?)
	`, conversationId).Scan(&exists)

	return exists, err
}
