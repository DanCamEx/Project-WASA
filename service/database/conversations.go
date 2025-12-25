package database

import "database/sql"

// CreateDirConversation checks if a direct conversation between two users exists, creates one if not, and returns its id.
// Both users are added as participants of the conversation in a separate table.
func (db *appdbimpl) CreateDirConversation(creatorId, participantId int) (int, error) {
	var exists bool
	var conversationId int
	if err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM conversations c
			JOIN participants p1 ON c.id = p1.conversation_id
			JOIN participants p2 ON c.id = p2.conversation_id
			WHERE c.category = 'direct' AND p1.user_id = ? AND p2.user_id = ?
		)
	`, creatorId, participantId).Scan(&exists); err != nil {
		return 0, err
	}
	if exists {
		return 0, nil
	}
	if err := db.c.QueryRow(`
		INSERT INTO conversations (category) VALUES ('direct') RETURNING id
	`).Scan(&conversationId); err != nil {
		return 0, err
	}
	if _, err := db.c.Exec(`
		INSERT INTO participants (user_id, conversation_id) VALUES (?, ?), (?, ?)
	`, creatorId, conversationId, participantId, conversationId); err != nil {
		return 0, err
	}
	return conversationId, nil
}

// CreateGrpConversation creates a new group conversation with the specified name and participants and returns its id.
// Both the creator and participants are added as participants of the conversation in a separate table.
func (db *appdbimpl) CreateGrpConversation(creatorId int, participantIds []int) (int, error) {
	var conversationId int
	if err := db.c.QueryRow(`
		INSERT INTO conversations (category) 
		VALUES ('group') 
		RETURNING id
	`).Scan(&conversationId); err != nil {
		return 0, err
	}
	if err := db.c.QueryRow(`
		INSERT INTO participants (user_id, conversation_id) 
		VALUES (?, ?) 
		RETURNING conversation_id
	`, creatorId, conversationId).Scan(&conversationId); err != nil {
		return conversationId, err
	}
	for _, participantId := range participantIds {
		_, err := db.c.Exec(`
			INSERT INTO participants (user_id, conversation_id) 
			VALUES (?, ?)
			RETURNING conversation_id
		`, participantId, conversationId)
		if err != nil {
			return conversationId, err
		}
	}
	return conversationId, nil
}

// SetGroupName updates the name of the group conversation with the given id.
// It returns an error if the conversation does not exist.
func (db *appdbimpl) SetGroupName(conversationId int, name string) error {
	res, err := db.c.Exec(`
		UPDATE conversations 
		SET name = ? 
		WHERE id = ?
	`, name, conversationId)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err == nil && rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// CheckConversation checks if a conversation with the given id exists.
func (db *appdbimpl) CheckConversation(conversationId int) (bool, error) {
	var exists bool
	if err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM conversations c
			JOIN participants p ON c.id = p.conversation_id
			WHERE c.id = ?)
	`, conversationId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// GetConversation retrieves the conversation information for the given conversation id.
// It returns the name, image, category, and latest message id of the conversation.
func (db *appdbimpl) GetConversation(conversationId int) (*Conversation, error) {
	var conversation Conversation
	if err := db.c.QueryRow(`
		SELECT id, name, image, category, latest_message_id
		FROM conversations
		WHERE id = ?
	`, conversationId).Scan(
		&conversation.Id,
		&conversation.Name,
		&conversation.Image,
		&conversation.Category,
		&conversation.LatestMessageId,
	); err != nil {
		return nil, err
	}
	return &conversation, nil
}

// CheckIfGroup checks if the conversation with the given id is a group conversation.
func (db *appdbimpl) CheckIfGroup(conversationId int) (bool, error) {
	var exists bool
	if err := db.c.QueryRow(`
		SELECT CASE WHEN category = 'group' THEN 1 ELSE 0 END
		FROM conversations
		WHERE id = ?
	`, conversationId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// GetAllUserConversations retrieves all conversations for the given user id.
func (db *appdbimpl) GetAllUserConversations(userId int) (*[]Conversation, error) {
	var conversations []Conversation
	rows, err := db.c.Query(`
		SELECT c.id, c.name, c.image, c.category, c.latest_message_id
		FROM conversations c
		JOIN participants p ON c.id = p.conversation_id
		WHERE p.user_id = ?
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id, latestMessageId int
	var name, image, category string
	for rows.Next() {
		if err := rows.Scan(&id, &name, &image, &category, &latestMessageId); err != nil {
			return nil, err
		}
		conversations = append(conversations, Conversation{
			Id:              id,
			Name:            name,
			Image:           image,
			Category:        category,
			LatestMessageId: latestMessageId,
		})
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return &conversations, nil
}

// UpdateLatestMessageId updates the latest message id of the conversation with the given id.
func (db *appdbimpl) UpdateLatestMessageId(conversationId int) error {
	if _, err := db.c.Exec(`
		UPDATE conversations
		SET latest_message_id = (
			SELECT MAX(id) FROM messages WHERE conversation_id = ?
		)
		WHERE id = ?
	`, conversationId, conversationId); err != nil {
		return err
	}
	return nil
}

// CheckParticipant checks if the user with the given id is a participant of the conversation with the given id.
func (db *appdbimpl) CheckParticipant(conversationId int, userId int) (bool, error) {
	var exists bool
	if err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM participants 
			WHERE conversation_id = ? AND user_id = ?
		)
	`, conversationId, userId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// CheckParticipantByMessageId checks if the user with the given id is a participant of the conversation
// that contains the message with the given id.
func (db *appdbimpl) CheckParticipantByMessageId(messageId int, userId int) (bool, error) {
	var exists bool
	if err := db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM participants p
			JOIN messages m ON p.conversation_id = m.conversation_id
			WHERE m.id = ? AND p.user_id = ?
		)
	`, messageId, userId).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// GetConversationMessages retrieves all the messages for the given conversation.
func (db *appdbimpl) GetConversationMessages(conversationIdIn int) (*[]Message, error) {
	rows, err := db.c.Query(`
		SELECT id, sender_id, content, timestamp, status, is_forwarded
		FROM messages
		WHERE conversation_id = ?
		ORDER BY timestamp ASC
	`, conversationIdIn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	var id, senderId, conversationIdOut int
	var content, timestamp, status string
	var isForwarded bool
	for rows.Next() {
		if err := rows.Scan(&id, &senderId, &content, &timestamp, &status, &isForwarded); err != nil {
			return nil, err
		}
		messages = append(messages, Message{
			Id:             id,
			SenderId:       senderId,
			ConversationId: conversationIdOut,
			Content:        content,
			Timestamp:      timestamp,
			Status:         status,
			IsForwarded:    isForwarded,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &messages, nil
}

func (db *appdbimpl) AddParticipant(conversationId int, userId int) error {
	if _, err := db.c.Exec(`
		INSERT INTO participants (user_id, conversation_id) 
		VALUES (?, ?)
	`, userId, conversationId); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) RemoveParticipant(conversationId int, userId int) error {
	if _, err := db.c.Exec(`
		DELETE FROM participants 
		WHERE user_id = ? AND conversation_id = ?
	`, userId, conversationId); err != nil {
		return err
	}
	return nil
}

func (db *appdbimpl) GetOtherContact(userId, conversationId int) (int, error) {
	var contactId int
	if err := db.c.QueryRow(`
		SELECT p.user_id
		FROM participants p
		WHERE p.conversation_id = ? AND p.user_id != ?
		LIMIT 1
	`, conversationId, userId).Scan(&contactId); err != nil {
		return 0, err
	}
	return contactId, nil
}

type Conversation struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Image           string `json:"image"`
	Category        string `json:"category"`
	LatestMessageId int    `json:"latestMessageId"`
}
