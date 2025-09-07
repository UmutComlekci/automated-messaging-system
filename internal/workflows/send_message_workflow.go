package workflows

import (
	"time"

	"github.com/umutcomlekci/automated-messaging-system/internal/activities"
	"github.com/umutcomlekci/automated-messaging-system/internal/repository/messages"
	"github.com/umutcomlekci/automated-messaging-system/internal/services/types"
	"github.com/umutcomlekci/automated-messaging-system/pkg/cache"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	SendSmsInput struct {
		Id          string `json:"id"`
		Content     string `json:"content"`
		PhoneNumber string `json:"phone_number"`
	}
)

func SendMessageWorkflow(ctx workflow.Context, input SendSmsInput) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
		TaskQueue: "sms-sender",
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	var (
		messageRepositoryActivities *activities.MessageRepositoryActivities
		smsServiceActivities        *activities.SmsServiceActivities
		cacheActivities             *activities.CacheActivities
	)

	var smsResult types.SmsResult
	err := workflow.ExecuteActivity(ctx, smsServiceActivities.Send, types.SmsMessage{
		From:    "905466011474",
		To:      input.PhoneNumber,
		Content: input.Content,
	}).Get(ctx, &smsResult)
	if err != nil {
		logger.Error("error sending message", "error", err.Error())
		err = workflow.ExecuteActivity(ctx, messageRepositoryActivities.UpdateMessageStatus, &activities.UpdateMessageStatusInput{
			Id:     input.Id,
			Status: messages.StatusFailed,
		}).Get(ctx, nil)
		return err
	}

	err = workflow.ExecuteActivity(ctx, messageRepositoryActivities.MessageSent, &activities.MessageSentInput{
		Id:                input.Id,
		ExternalMessageId: smsResult.SmsId,
		SentAt:            workflow.Now(ctx),
	}).Get(ctx, nil)
	if err != nil {
		logger.Error("error updating message status", "error", err.Error())
		return err
	}

	err = workflow.ExecuteActivity(ctx, cacheActivities.SetStruct, smsResult.SmsId, &cache.SentMessageCache{
		PhoneNumber:       input.PhoneNumber,
		Message:           input.Content,
		SentAt:            workflow.Now(ctx),
		ExternalMessageId: smsResult.SmsId,
	}).Get(ctx, nil)
	if err != nil {
		logger.Error("error setting cache", "error", err.Error())
		return err
	}

	return nil
}
