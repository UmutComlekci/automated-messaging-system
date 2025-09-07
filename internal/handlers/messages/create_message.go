package messages

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/umutcomlekci/automated-messaging-system/internal/handlers"
)

type CreateMessageRequest struct {
	Content     string `json:"content" validate:"required,max=160"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
}

// @Summary Create a new message
// @Description Add a new message to be sent automatically
// @Tags messages
// @Accept json
// @Produce json
// @Param message body CreateMessageRequest true "Message data"
// @Success 201 {object} handlers.APIResponse
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /api/v1/messages [post]
func (h *messagesHandler) CreateMessage(c *fiber.Ctx) error {
	var request CreateMessageRequest
	err := c.BodyParser(&request)
	if err != nil {
		h.logger.Error("error parsing request body", slog.String("error", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(handlers.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	err = h.validator.Struct(request)
	if err != nil {
		h.logger.Error("error validating request", slog.String("error", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(handlers.APIResponse{
			Success: false,
			Message: "Validation failed",
		})
	}

	message, err := h.messageRepository.Create(request.Content, request.PhoneNumber)
	if err != nil {
		h.logger.Error("error creating message", slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.APIResponse{
			Success: false,
			Message: "Failed to create message",
		})
	}

	h.logger.Info("message created successfully", slog.String("id", message.Id.String()))
	return c.Status(fiber.StatusCreated).JSON(message)
}
