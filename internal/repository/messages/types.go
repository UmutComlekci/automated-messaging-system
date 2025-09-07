package messages

import (
	"time"

	"github.com/google/uuid"
)

type (
	Status string

	Message struct {
		Id                uuid.UUID  `json:"id" db:"id"`
		Content           string     `json:"content" db:"content" validate:"required,max=160"`
		PhoneNumber       string     `json:"phone_number" db:"phone_number" validate:"required,e164"`
		Status            Status     `json:"status" db:"status"`
		CreatedAt         time.Time  `json:"created_at" db:"created_at"`
		SentAt            *time.Time `json:"sent_at,omitempty" db:"sent_at"`
		ExternalMessageID *string    `json:"external_message_id,omitempty" db:"external_message_id"`
	}
)

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusSent       Status = "sent"
	StatusFailed     Status = "failed"
)
