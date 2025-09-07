package scheduler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/umutcomlekci/automated-messaging-system/internal/handlers"
	"go.temporal.io/sdk/client"
)

// @Summary Stop automatic message sending
// @Description Stop the automatic message sending scheduler
// @Tags scheduler
// @Accept json
// @Produce json
// @Success 200 {object} handlers.APIResponse
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /api/v1/scheduler/stop [post]
func (h *schedulerHandler) StopScheduler(c *fiber.Ctx) error {
	scheduleClient := h.temporalClient.ScheduleClient()

	pendingMessagesScheduleHandle := scheduleClient.GetHandle(context.Background(), PendingMessagesScheduleName)
	scheduleDetail, err := pendingMessagesScheduleHandle.Describe(context.Background())
	if err != nil {
		if strings.HasPrefix(err.Error(), "workflow not found for ID") {
			return c.Status(fiber.StatusOK).JSON(handlers.APIResponse{
				Success: true,
				Message: "Scheduler stopped successfully",
			})
		}
		h.logger.Error("temporal pause schedule failed", slog.String("error", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.APIResponse{
			Success: false,
			Message: "Failed to stop scheduler",
		})
	} else if !scheduleDetail.Schedule.State.Paused {
		err = pendingMessagesScheduleHandle.Pause(context.Background(), client.SchedulePauseOptions{})
		if err != nil {
			h.logger.Error("temporal pause schedule failed", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(handlers.APIResponse{
				Success: false,
				Message: "Failed to stop scheduler",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(handlers.APIResponse{
		Success: true,
		Message: "Scheduler stopped successfully",
	})
}
