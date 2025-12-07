package database

func (db *DbImpl) CreateComment(messageId int, reacterId int, emojiIn string) (emojiOut string, err error) {

	err = db.c.QueryRow(`

		INSERT INTO comments (reacter_id, message_id, emoji) 
		VALUES (?, ?, ?) RETURNING emoji

	`, reacterId, messageId, emojiIn).Scan(&emojiOut)
	return emojiOut, err
}

func (db *DbImpl) GetComment(messageIdIn int, reacterIdIn int) (messageIdOut int, reacterIdOut int, emoji string, err error) {

	err = db.c.QueryRow(`

		SELECT reacter_id, message_id, emoji 
		FROM comments 
		WHERE reacter_id = ? AND message_id = ?

	`, reacterIdIn, messageIdIn).Scan(&reacterIdOut, &messageIdOut, &emoji)

	return reacterIdOut, messageIdOut, emoji, err
}

func (db *DbImpl) DeleteComment(messageId int, reacterId int) error {
	_, err := db.c.Exec(`

		DELETE FROM comments 
		WHERE reacter_id = ? AND message_id = ?
		
	`, reacterId, messageId)
	return err
}

func (db *DbImpl) CheckComment(messageId int, reacterId int) (exists bool, err error) {

	err = db.c.QueryRow(`

		SELECT EXISTS(
			SELECT 1 
			FROM comments 
			WHERE reacter_id = ? AND message_id = ?)

	`, reacterId, messageId).Scan(&exists)
	return exists, err
}
