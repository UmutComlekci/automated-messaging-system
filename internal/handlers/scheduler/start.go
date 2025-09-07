package scheduler

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/umutcomlekci/automated-messaging-system/internal/handlers"
	"github.com/umutcomlekci/automated-messaging-system/internal/workflows"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

// @Summary Start automatic message sending
// @Description Start the automatic message sending scheduler
// @Tags scheduler
// @Accept json
// @Produce json
// @Success 200 {object} handlers.APIResponse
// @Failure 400 {object} handlers.APIResponse
// @Failure 500 {object} handlers.APIResponse
// @Router /api/v1/scheduler/start [post]
func (h *schedulerHandler) StartScheduler(c *fiber.Ctx) error {
	scheduleClient := h.temporalClient.ScheduleClient()

	pendingMessagesScheduleHandle := scheduleClient.GetHandle(context.Background(), PendingMessagesScheduleName)
	scheduleDetail, err := pendingMessagesScheduleHandle.Describe(context.Background())
	if err != nil {
		if strings.HasPrefix(err.Error(), "workflow not found for ID") {
			duration, _ := time.ParseDuration("2m")
			_, err = scheduleClient.Create(context.Background(), client.ScheduleOptions{
				ID:      PendingMessagesScheduleName,
				Overlap: enums.SCHEDULE_OVERLAP_POLICY_TERMINATE_OTHER,
				Action: &client.ScheduleWorkflowAction{
					Workflow:                 workflows.PendingMessagesWorkflow,
					TaskQueue:                "sms",
					WorkflowTaskTimeout:      time.Minute * 5,
					WorkflowRunTimeout:       time.Hour * 4,
					WorkflowExecutionTimeout: time.Hour * 4,
				},
				Spec: client.ScheduleSpec{
					Intervals: []client.ScheduleIntervalSpec{
						{
							Every: duration,
						},
					},
				},
			})
		}
		if err != nil {
			h.logger.Error("temporal create schedule failed", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(handlers.APIResponse{
				Success: false,
				Message: "Failed to start scheduler",
			})
		}
	} else if scheduleDetail.Schedule.State.Paused {
		err = pendingMessagesScheduleHandle.Unpause(context.Background(), client.ScheduleUnpauseOptions{})
		if err != nil {
			h.logger.Error("temporal unpause schedule failed", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(handlers.APIResponse{
				Success: false,
				Message: "Failed to start scheduler",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(handlers.APIResponse{
		Success: true,
		Message: "Scheduler started successfully",
	})
}
