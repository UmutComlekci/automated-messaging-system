package worker

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/umutcomlekci/automated-messaging-system/internal/activities"
	"github.com/umutcomlekci/automated-messaging-system/internal/handlers/scheduler"
	"github.com/umutcomlekci/automated-messaging-system/internal/logging"
	"github.com/umutcomlekci/automated-messaging-system/internal/repository/messages"
	"github.com/umutcomlekci/automated-messaging-system/internal/workflows"
	"github.com/umutcomlekci/automated-messaging-system/pkg/database"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func newPendingMessagesWorkerCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "pending-messages",
		RunE: pendingMessagesWorkerCommand,
	}

	return rootCmd
}

func pendingMessagesWorkerCommand(cmd *cobra.Command, args []string) error {
	workerLogger := logging.NewLogger("pending_messages_worker")

	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
		Logger:   log.NewStructuredLogger(logging.NewLogger("temporal")),
	})
	if err != nil {
		workerLogger.Error("unable to create client", slog.String("error", err.Error()))
		return err
	}
	defer c.Close()

	scheduleClient := c.ScheduleClient()
	pendingMessagesScheduleHandle := scheduleClient.GetHandle(context.Background(), scheduler.PendingMessagesScheduleName)
	_, err = pendingMessagesScheduleHandle.Describe(context.Background())
	if err != nil {
		if strings.HasPrefix(err.Error(), "workflow not found for ID") {
			duration, _ := time.ParseDuration("2m")
			_, err = scheduleClient.Create(context.Background(), client.ScheduleOptions{
				ID:      scheduler.PendingMessagesScheduleName,
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
			if err != nil {
				workerLogger.Error("temporal create schedule failed", slog.String("error", err.Error()))
				return err
			}
		}
	}

	w := worker.New(c, "sms", worker.Options{})
	w.RegisterWorkflow(workflows.PendingMessagesWorkflow)

	database, err := database.NewDatabase()
	if err != nil {
		workerLogger.Error("error connecting to database", slog.String("error", err.Error()))
		return err
	}

	messageRepository := messages.NewMessageRepository(database)
	messageRepositoryActivities := activities.NewMessageRepositoryActivities(messageRepository)
	w.RegisterActivity(messageRepositoryActivities.GetPendingMessages)
	w.RegisterActivity(messageRepositoryActivities.UpdateMessageStatus)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		workerLogger.Error("unable to start worker", slog.String("error", err.Error()))
		return err
	}

	return nil
}
