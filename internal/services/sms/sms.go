package sms

import (
	"log/slog"

	"github.com/umutcomlekci/automated-messaging-system/internal/logging"
	"github.com/umutcomlekci/automated-messaging-system/internal/services/sms/provider"
	"github.com/umutcomlekci/automated-messaging-system/internal/services/types"
)

type Service struct {
	provider provider.Provider
	logger   *slog.Logger
}

func NewSmsService(provider provider.Provider) *Service {
	return &Service{
		provider: provider,
		logger:   logging.NewLogger("sms-service"),
	}
}

func (s *Service) Send(message *types.SmsMessage) (*types.SmsResult, error) {
	s.logger.Info("sending sms", slog.String("to", message.To), slog.String("from", message.From))
	result, err := s.provider.Send(message)
	if err != nil {
		s.logger.Error("error sending sms", "error", err.Error(), slog.String("to", message.To), slog.String("from", message.From))
		return nil, err
	}

	s.logger.Info("sms sent", slog.String("to", message.To), slog.String("from", message.From))
	return result, nil
}
