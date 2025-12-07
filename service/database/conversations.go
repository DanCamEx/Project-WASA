package database

// CreateDirConversation creates a new conversation between two users and returns its id.
// The conversation type is set to 'direct'.
// Both users are added as participants of the conversation in a separate table.
func (db *DbImpl) CreateDirConversation(creatorId, participantId int) (conversationCreatedId int, err error) {

	// Check if the conv already exists
	err = db.c.QueryRow(`

		SELECT c.id
		FROM conversations c
		JOIN participants p1 ON c.id = p1.conversation_id
		JOIN participants p2 ON c.id = p2.conversation_id
		WHERE c.type = 'direct' AND p1.user_id = ? AND p2.user_id = ?

	`, creatorId, participantId).Scan(&conversationCreatedId)
	if err == nil { return 0, nil }

	// Create the conv
	err = db.c.QueryRow(`

		INSERT INTO conversations (type) 
		VALUES ('direct') RETURNING id

	`).Scan(&conversationCreatedId)
	if err != nil { return conversationCreatedId, err }

	// Add both users as participants
	_, err = db.c.Exec(`

		INSERT INTO participants (user_id, conversation_id) 
		VALUES (?, ?), (?, ?)
		
	`, creatorId, conversationCreatedId, participantId, conversationCreatedId)
	return conversationCreatedId, err
}

// CreateGrpConversation creates a new group conversation with the specified name and participants and returns its id.
// The conversation type is set to 'group'.
// Both the creator and participants are added as participants of the conversation in a separate table.
func (db *DbImpl) CreateGrpConversation(creatorId int, participantIds []int) (conversationId int, err error) {

	// Create the group
	err = db.c.QueryRow(`

		INSERT INTO conversations (type) 
		VALUES ('group') 
		RETURNING id

	`).Scan(&conversationId)
	if err != nil { return conversationId, err }

	// Add the creator as a participant
	_, err = db.c.Exec(`

		INSERT INTO participants (user_id, conversation_id) 
		VALUES (?, ?)

	`, creatorId, conversationId)
	if err != nil { return conversationId, err }
	
	// Add all other members as participants
	for _, participantId := range participantIds {
		_, err = db.c.Exec(`

			INSERT INTO participants (user_id, conversation_id) 
			VALUES (?, ?)
			
		`, participantId, conversationId)
		if err != nil { return conversationId, err }
	}

	return conversationId, nil
}

// SetGroupName updates the name of the group conversation with the given id.
func (db *DbImpl) SetGroupName(conversationId int, name string) error {

	_, err := db.c.Exec(`

		UPDATE conversations 
		SET name = ? 
		WHERE id = ?

	`, name, conversationId)
	return err
}

// CheckConversation checks if a conversation with the given id exists.
func (db *DbImpl) CheckConversation(conversationId int) (conversationExists bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 FROM conversations c
			JOIN participants p ON c.id = p.conversation_id
			WHERE c.id = ?)

	`, conversationId).Scan(&conversationExists)
	return conversationExists, err
}

// GetConversation retrieves the conversation information for the given conversation id.
// It returns the name, image, type, and latest message id of the conversation.
func (db *DbImpl) GetConversation(conversationId int) (name string, image string, conversationType string, latestMessageId int, err error) {
	
	err = db.c.QueryRow(`

		SELECT name, image, type, latest_message_id
		FROM conversations
		WHERE id = ?

	`, conversationId).Scan(&name, &image, &conversationType, &latestMessageId)
	return name, image, conversationType, latestMessageId, err
}

// CheckIfGroup checks if the conversation with the given id is a group conversation.
func (db *DbImpl) CheckIfGroup(conversationId int) (isGroup bool, err error) {

	err = db.c.QueryRow(`

		SELECT CASE WHEN type = 'group' THEN 1 ELSE 0 END
		FROM conversations
		WHERE id = ?

	`, conversationId).Scan(&isGroup)
	return isGroup, err
}

// GetAllUserConversations retrieves all conversations for the given user id.
func (db *DbImpl) GetAllUserConversations(userId int) (userConversations [][]string, err error) {
	
	// Get all the conversation table rows and then takes only the ones where the user is a participant
	rows, err := db.c.Query(`

		SELECT c.id, c.name, c.image, c.type, c.latest_message_id
		FROM conversations c
		JOIN participants p ON c.id = p.conversation_id
		WHERE p.user_id = ?

	`, userId)
	if err != nil { return nil, err }

	// For each row, scan the values into variables and append to the result list
	var id, name, image, conversationType, latestMessageId string
	for rows.Next() {
		err := rows.Scan(&id, &name, &image, &conversationType, &latestMessageId)
		if err != nil { return nil, err }

		userConversations = append(userConversations, []string{id, name, image, conversationType, latestMessageId})
	}

	return userConversations, nil
}

// UpdateLatestMessageId updates the latest message id of the conversation with the given id.
func (db *DbImpl) UpdateLatestMessageId(conversationId int) error {
	_, err := db.c.Exec(`

		UPDATE conversations
		SET latest_message_id = (
			SELECT MAX(id) FROM messages WHERE conversation_id = ?
		)
		WHERE id = ?

	`, conversationId, conversationId)
	return err
}

// CheckParticipant checks if the user with the given id is a participant of the conversation with the given id.
func (db *DbImpl) CheckParticipant(conversationId int, userId int) (isParticipant bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 FROM participants 
			WHERE conversation_id = ? AND user_id = ?
		)

	`, conversationId, userId).Scan(&isParticipant)
	return isParticipant, err
}

func ( db *DbImpl) CheckParticipantByMessageId(messageId int, userId int) (isParticipant bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 
			FROM participants p
			JOIN messages m ON p.conversation_id = m.conversation_id
			WHERE m.id = ? AND p.user_id = ?
		)

	`, messageId, userId).Scan(&isParticipant)
	return isParticipant, err
}
