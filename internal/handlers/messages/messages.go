package messages

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/umutcomlekci/automated-messaging-system/internal/logging"
	"github.com/umutcomlekci/automated-messaging-system/internal/repository/messages"
)

type messagesHandler struct {
	messageRepository *messages.Repository
	validator         *validator.Validate
	logger            *slog.Logger
}

func NewHandler(messageRepository *messages.Repository) *messagesHandler {
	return &messagesHandler{
		messageRepository: messageRepository,
		validator:         validator.New(validator.WithRequiredStructEnabled()),
		logger:            logging.NewLogger("messages_handler"),
	}
}

func (h *messagesHandler) InitRoutes(router fiber.Router) {
	route := router.Group("/messages")
	route.Get("/sent", h.GetSentMessages)
	route.Post("/", h.CreateMessage)
}
