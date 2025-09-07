package scheduler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/umutcomlekci/automated-messaging-system/internal/logging"
	"go.temporal.io/sdk/client"
)

type schedulerHandler struct {
	temporalClient client.Client
	logger         *slog.Logger
}

const (
	PendingMessagesScheduleName = "pending-messages"
)

func NewHandler(temporalClient client.Client) *schedulerHandler {
	return &schedulerHandler{
		temporalClient: temporalClient,
		logger:         logging.NewLogger("scheduler_handler"),
	}
}

func (h *schedulerHandler) InitRoutes(router fiber.Router) {
	route := router.Group("/scheduler")
	route.Post("/start", h.StartScheduler)
	route.Post("/stop", h.StopScheduler)
}
