package messages

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/umutcomlekci/automated-messaging-system/internal/handlers"
)

// @Summary Get sent messages
// @Description Retrieve a list of sent messages with pagination
// @Tags messages
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} handlers.APIResponse
// @Failure 400 {object} handlers.APIResponse
// @Router /api/v1/messages/sent [get]
func (h *messagesHandler) GetSentMessages(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit > 100 {
		limit = 100
	}

	result, err := h.messageRepository.GetSents(page, limit)
	if err != nil {
		h.logger.Error("error getting sent messages", slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.APIResponse{
			Success: false,
			Message: "Failed to retrieve sent messages",
		})
	}

	h.logger.Info("sent messages retrieved successfully")
	return c.JSON(handlers.APIResponse{
		Success: true,
		Message: "Sent messages retrieved successfully",
		Data:    result,
	})
}
