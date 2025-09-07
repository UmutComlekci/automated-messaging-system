package activities

import (
	"time"

	"github.com/umutcomlekci/automated-messaging-system/internal/repository/messages"
)

type (
	PaginationInput struct {
		Page  int `json:"page"`
		Limit int `json:"limit"`
	}

	UpdateMessageStatusInput struct {
		Id     string          `json:"id"`
		Status messages.Status `json:"status"`
	}

	MessageSentInput struct {
		Id                string    `json:"id"`
		ExternalMessageId string    `json:"external_message_id"`
		SentAt            time.Time `json:"sent_at"`
	}
)
