package api

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/spf13/cobra"
	"github.com/umutcomlekci/automated-messaging-system/internal/config"
	"github.com/umutcomlekci/automated-messaging-system/internal/handlers/messages"
	"github.com/umutcomlekci/automated-messaging-system/internal/handlers/scheduler"
	"github.com/umutcomlekci/automated-messaging-system/internal/logging"
	messagesR "github.com/umutcomlekci/automated-messaging-system/internal/repository/messages"
	"github.com/umutcomlekci/automated-messaging-system/pkg/database"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/log"

	_ "github.com/umutcomlekci/automated-messaging-system/docs"
)

func newApiCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:  "serve",
		RunE: apiCommand,
	}

	return rootCmd
}

// @title Message Scheduler API
// @version 1.0
// @description Automatic message sending system
// @host localhost:8080
// @BasePath /
func apiCommand(cmd *cobra.Command, args []string) error {
	apiLogger := logging.NewLogger("api")

	database, err := database.NewDatabase()
	if err != nil {
		apiLogger.Error("Error connecting to database", slog.String("error", err.Error()))
		return err
	}

	app := fiber.New(fiber.Config{})
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} [${ip}]:${port} ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes
	api := app.Group("/api/v1")

	// Scheduler control endpoints
	temporalClient, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
		Logger:   log.NewStructuredLogger(logging.NewLogger("temporal")),
	})
	if err != nil {
		apiLogger.Error("unable to create client", slog.String("error", err.Error()))
		return err
	}
	scheduler.NewHandler(temporalClient).InitRoutes(api)

	// Message endpoints
	messageRepository := messagesR.NewMessageRepository(database)
	messages.NewHandler(messageRepository).InitRoutes(api)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	port := config.GetApiPort()
	apiLogger.Info("Server starting on port", slog.String("port", port))
	if err = app.Listen(":" + port); err != nil {
		apiLogger.Error("Error starting server", slog.String("error", err.Error()))
		return err
	}

	return nil
}
