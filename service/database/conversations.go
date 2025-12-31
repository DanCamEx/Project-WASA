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

// CreateGrpConversation creates a new group conversation adding the creator and all the members in the participants table, and returns its id.
// Both the creator and participants are added as participants of the conversation in a separate table.
func (db *appdbimpl) CreateGrpConversation(creatorId int, participantIds []int, name string) (int, error) {
	var conversationId int
	if err := db.c.QueryRow(`
		INSERT INTO conversations (category, name) 
		VALUES ('group', ?)
		RETURNING id
	`, name).Scan(&conversationId); err != nil {
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

// SetGroupName updates the name of the group conversation by id.
// If it doensn't exist or is not a group, returns an error.
func (db *appdbimpl) SetGroupName(conversationId int, name string) error {
	res, err := db.c.Exec(`
		UPDATE conversations 
		SET name = ? 
		WHERE id = ?
		AND category = 'group'
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

// CheckConversation checks if a conversation exists by id.
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

// GetConversation returns the conversation info by id.
// It returns its name, photo, category, and latest message.
func (db *appdbimpl) GetConversation(conversationId int) (*Conversation, error) {
	var conversation Conversation
	if err := db.c.QueryRow(`
		SELECT id, name, photo, category, latest_message_content, latest_message_timestamp
		FROM conversations
		WHERE id = ?
	`, conversationId).Scan(
		&conversation.Id,
		&conversation.Name,
		&conversation.Photo,
		&conversation.Category,
		&conversation.LatestMessage.Content,
		&conversation.LatestMessage.Timestamp,
	); err != nil {
		return nil, err
	}
	return &conversation, nil
}

// CheckIfGroup checks if the conversation is a group conversation.
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

// GetAllUserConversations returns all conversations of the user.
func (db *appdbimpl) GetAllUserConversations(userId int) (*[]Conversation, error) {
	var conversations []Conversation
	rows, err := db.c.Query(`
		SELECT c.id, c.name, c.photo, c.category, c.latest_message_content, c.latest_message_timestamp
		FROM conversations c
		JOIN participants p ON c.id = p.conversation_id
		WHERE p.user_id = ?
	`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var id int
	var name, photo, category, latestMessageContent, latestMessageTimestamp string
	for rows.Next() {
		if err := rows.Scan(&id, &name, &photo, &category, &latestMessageContent, &latestMessageTimestamp); err != nil {
			return nil, err
		}
		conversations = append(conversations, Conversation{
			Id:       id,
			Name:     name,
			Photo:    photo,
			Category: category,
			LatestMessage: LatestMessage{
				Content:   latestMessageContent,
				Timestamp: latestMessageTimestamp,
			},
		})
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return &conversations, nil
}

// UpdateLatestMessage updates the latest message sent in the conversation using its timestamp.
func (db *appdbimpl) UpdateLatestMessage(conversationId int) error {
	if _, err := db.c.Exec(`
		UPDATE conversations
		SET
			latest_message_content = COALESCE((
				SELECT content
				FROM messages
				WHERE conversation_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			), ''),
			latest_message_timestamp = COALESCE((
				SELECT timestamp
				FROM messages
				WHERE conversation_id = ?
				ORDER BY timestamp DESC
				LIMIT 1
			), '')
			WHERE id = ?
	`, conversationId, conversationId, conversationId); err != nil {
		return err
	}
	return nil
}

// CheckParticipant checks if the user is a participant of the conversation.
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

// CheckParticipantByMessageId checks if the user is a participant of the conversation using a message id.
// Used to know if the user can actually access the message.
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

func (db *appdbimpl) GetParticipants(conversationId int) ([]int, error) {
	rows, err := db.c.Query(`
		SELECT user_id
		FROM participants
		WHERE conversation_id = ?
	`, conversationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var participantIds []int
	var userId int
	for rows.Next() {
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}
		participantIds = append(participantIds, userId)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return participantIds, nil
}

// GetConversationMessages returns all messages of the conversation.
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

// AddParticipant adds the user id to the conversation participants.
func (db *appdbimpl) AddParticipant(conversationId int, userId int) error {
	if _, err := db.c.Exec(`
		INSERT INTO participants (user_id, conversation_id) 
		VALUES (?, ?)
	`, userId, conversationId); err != nil {
		return err
	}
	return nil
}

// RemoveParticipant removes the user id linked to the conversation from the participants.
func (db *appdbimpl) RemoveParticipant(conversationId int, userId int) error {
	if _, err := db.c.Exec(`
		DELETE FROM participants 
		WHERE user_id = ? AND conversation_id = ?
	`, userId, conversationId); err != nil {
		return err
	}
	return nil
}

// GetOtherContact returns the other user id of the conversation.
// It should be used only for direct conversations.
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

// SetGroupPhoto updates the group photo or set it if was "".
func (db *appdbimpl) SetGroupPhoto(conversationId int, filename string) error {
	_, err := db.c.Exec(`
		UPDATE conversations
		SET photo = ?
		WHERE id = ?
	`, filename, conversationId)
	if err != nil {
		return err
	}
	return nil
}

type LatestMessage struct {
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}
type Conversation struct {
	Id            int           `json:"id"`
	Name          string        `json:"name"`
	Photo         string        `json:"photo"`
	Category      string        `json:"category"`
	LatestMessage LatestMessage `json:"latestMessage"`
}
