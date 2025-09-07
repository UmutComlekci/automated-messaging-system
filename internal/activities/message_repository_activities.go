package activities

import (
	"github.com/google/uuid"
	"github.com/umutcomlekci/automated-messaging-system/internal/repository/messages"
)

type (
	MessageRepositoryActivities struct {
		messageRepository *messages.Repository
	}
)

func NewMessageRepositoryActivities(messageRepository *messages.Repository) *MessageRepositoryActivities {
	return &MessageRepositoryActivities{
		messageRepository: messageRepository,
	}
}

func (m *MessageRepositoryActivities) GetPendingMessages(input *PaginationInput) ([]messages.Message, error) {
	return m.messageRepository.GetPendings(input.Page, input.Limit)
}

func (m *MessageRepositoryActivities) UpdateMessageStatus(input *UpdateMessageStatusInput) error {
	messageId, err := uuid.Parse(input.Id)
	if err != nil {
		return err
	}

	err = m.messageRepository.UpdateStatus(messageId, input.Status)
	if err != nil {
		return err
	}

	return nil
}

func (m *MessageRepositoryActivities) MessageSent(input MessageSentInput) error {
	messageId, err := uuid.Parse(input.Id)
	if err != nil {
		return err
	}

	err = m.messageRepository.UpdateMessageAsSent(messageId, input.ExternalMessageId, input.SentAt)
	if err != nil {
		return err
	}

	return nil
}
