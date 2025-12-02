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

// CreateGrpConversation creates a new group conversation with the specified name and participants.
// The creator is also added as a participant.
// It returns the id of the newly created group conversation.
func (db *DbImpl) CreateGrpConversation(creatorId int, participantIds []int) (int, error) {
	var conversationId int

	err := db.c.QueryRow(`
		INSERT INTO conversations (type) VALUES ('group') RETURNING id
	`).Scan(&conversationId)
	if err != nil {
		return conversationId, err
	}
	_, err = db.c.Exec(`
		INSERT INTO participants (user_id, conversation_id) VALUES (?, ?)
	`, creatorId, conversationId)
	if err != nil {
		return conversationId, err
	}
	for _, participantId := range participantIds {
		_, err = db.c.Exec(`
			INSERT INTO participants (user_id, conversation_id) VALUES (?, ?)
		`, participantId, conversationId)
		if err != nil {
			return conversationId, err
		}
	}

	return conversationId, nil
}

func (db *DbImpl) SetGroupName(conversationId int, name string) error {
	_, err := db.c.Exec(`
		UPDATE conversations SET name = ? WHERE id = ?
	`, name, conversationId)
	return err
}

// CheckConversation checks if a conversation with the given id exists.
func (db *DbImpl) CheckConversation(conversationId int) (bool, error) {
	var exists bool

	err := db.c.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM conversations WHERE id = ?)
	`, conversationId).Scan(&exists)

	return exists, err
}
