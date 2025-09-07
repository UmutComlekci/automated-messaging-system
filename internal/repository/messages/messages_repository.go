package messages

import (
	"time"

	"github.com/google/uuid"
	"github.com/umutcomlekci/automated-messaging-system/pkg/database"
)

const (
	getPendingMessagesQuery = `
		SELECT 
			id, content, phone_number, status, created_at
		FROM
			messages
		WHERE
			status = 'pending'
		ORDER BY created_at ASC
		LIMIT $1
		OFFSET $2
	`
	insertMessageQuery = `
		INSERT INTO messages (id, content, phone_number, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	updateMessageStatusQuery = `
		UPDATE messages SET status = $1 WHERE id = $2
	`
	updateMessageSentQuery = `
		UPDATE messages SET external_message_id = $1, status = $2, sent_at = $3 WHERE id = $4
	`
	getSentMessagesQuery = `
		SELECT
			id, content, phone_number, status, created_at, sent_at, external_message_id
		FROM
			messages
		WHERE
			status = 'sent'
		ORDER BY
			sent_at DESC
		LIMIT $1
		OFFSET $2
	`
)

type Repository struct {
	db database.Database
}

func NewMessageRepository(db database.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (m *Repository) Create(content, phone string) (*Message, error) {
	message := &Message{
		Id:          uuid.New(),
		Content:     content,
		PhoneNumber: phone,
		Status:      StatusPending,
		CreatedAt:   time.Now().UTC(),
	}

	_, err := m.db.Exec(insertMessageQuery, message.Id, message.Content, message.PhoneNumber, message.Status, message.CreatedAt)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (m *Repository) GetPendings(page, limit int) ([]Message, error) {
	rows, err := m.db.Query(getPendingMessagesQuery, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.Content, &message.PhoneNumber, &message.Status, &message.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (m *Repository) UpdateStatus(id uuid.UUID, status Status) error {
	_, err := m.db.Exec(updateMessageStatusQuery, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Repository) UpdateMessageAsSent(id uuid.UUID, externalMessageId string, sentAt time.Time) error {
	_, err := m.db.Exec(updateMessageSentQuery, externalMessageId, StatusSent, sentAt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Repository) GetSents(page, limit int) ([]Message, error) {
	rows, err := m.db.Query(getSentMessagesQuery, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		err := rows.Scan(&message.Id, &message.Content, &message.PhoneNumber, &message.Status, &message.CreatedAt, &message.SentAt, &message.ExternalMessageID)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}
