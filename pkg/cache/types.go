package cache

import "time"

type SentMessageCache struct {
	PhoneNumber       string    `json:"phone_number"`
	Message           string    `json:"message"`
	SentAt            time.Time `json:"sent_at"`
	ExternalMessageId string    `json:"external_message_id"`
}
