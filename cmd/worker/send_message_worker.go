package worker

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/umutcomlekci/automated-messaging-system/internal/activities"
	"github.com/umutcomlekci/automated-messaging-system/internal/config"
	"github.com/umutcomlekci/automated-messaging-system/internal/logging"
	"github.com/umutcomlekci/automated-messaging-system/internal/repository/messages"
	"github.com/umutcomlekci/automated-messaging-system/internal/services/sms"
	"github.com/umutcomlekci/automated-messaging-system/internal/services/sms/provider"
	"github.com/umutcomlekci/automated-messaging-system/internal/workflows"
	"github.com/umutcomlekci/automated-messaging-system/pkg/cache"
	"github.com/umutcomlekci/automated-messaging-system/pkg/database"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
)

func newSendMessageWorkerCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "send-message",
		RunE: sendMessageWorkerCommand,
	}

	return rootCmd
}

func sendMessageWorkerCommand(cmd *cobra.Command, args []string) error {
	workerLogger := logging.NewLogger("send_message_worker")

	c, err := client.Dial(client.Options{
		HostPort: config.GetTemporalHostPort(),
		Logger:   log.NewStructuredLogger(logging.NewLogger("temporal")),
	})
	if err != nil {
		workerLogger.Error("unable to create client", slog.String("error", err.Error()))
		return err
	}

	defer c.Close()
	w := worker.New(c, "sms-sender", worker.Options{})
	w.RegisterWorkflow(workflows.SendMessageWorkflow)

	database, err := database.NewDatabase()
	if err != nil {
		workerLogger.Error("error connecting to database", slog.String("error", err.Error()))
		return err
	}

	smsProvider, err := provider.NewProvider()
	if err != nil {
		workerLogger.Error("error creating sms provider", slog.String("error", err.Error()))
		return err
	}

	cacheClient, err := cache.NewCacheClient()
	if err != nil {
		workerLogger.Error("error connecting to cache", slog.String("error", err.Error()))
		return err
	}

	messageRepository := messages.NewMessageRepository(database)
	messageRepositoryActivities := activities.NewMessageRepositoryActivities(messageRepository)
	w.RegisterActivity(messageRepositoryActivities.GetPendingMessages)
	w.RegisterActivity(messageRepositoryActivities.UpdateMessageStatus)
	w.RegisterActivity(messageRepositoryActivities.MessageSent)

	smsService := sms.NewSmsService(smsProvider)
	smsServiceActivities := activities.NewSmSServiceActivities(smsService)
	w.RegisterActivity(smsServiceActivities.Send)

	cacheActivities := activities.NewCacheActivities(cacheClient)
	w.RegisterActivity(cacheActivities.SetStruct)
	w.RegisterActivity(cacheActivities.GetStruct)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		workerLogger.Error("unable to start worker", slog.String("error", err.Error()))
		return err
	}

	return nil
}
