package workflows

import (
	"fmt"
	"time"

	"github.com/umutcomlekci/automated-messaging-system/internal/activities"
	"github.com/umutcomlekci/automated-messaging-system/internal/repository/messages"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func PendingMessagesWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
		TaskQueue: "sms",
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	var (
		pendingMessages             []messages.Message
		messageRepositoryActivities *activities.MessageRepositoryActivities
	)

	err := workflow.ExecuteActivity(ctx, messageRepositoryActivities.GetPendingMessages, &activities.PaginationInput{
		Page:  1,
		Limit: 2,
	}).Get(ctx, &pendingMessages)
	if err != nil {
		logger.Error("error getting pending messages", "error", err.Error())
		return err
	}

	for _, message := range pendingMessages {
		childWorkflowOptions := workflow.ChildWorkflowOptions{
			WorkflowID:        fmt.Sprintf("send-message-%s", message.Id.String()),
			TaskQueue:         "sms-sender",
			ParentClosePolicy: enums.PARENT_CLOSE_POLICY_ABANDON,
		}

		childWorkflowCtx := workflow.WithChildOptions(ctx, childWorkflowOptions)
		err = workflow.ExecuteChildWorkflow(childWorkflowCtx, SendMessageWorkflow, &SendSmsInput{
			Id:          message.Id.String(),
			Content:     message.Content,
			PhoneNumber: message.PhoneNumber,
		}).GetChildWorkflowExecution().Get(childWorkflowCtx, nil)
		if err != nil {
			if temporal.IsWorkflowExecutionAlreadyStartedError(err) {
				err = workflow.ExecuteActivity(ctx, messageRepositoryActivities.UpdateMessageStatus, &activities.UpdateMessageStatusInput{
					Id:     message.Id.String(),
					Status: messages.StatusProcessing,
				}).Get(ctx, nil)
				if err != nil {
					logger.Error("error updating message status", "error", err.Error())
					return err
				}
				continue
			}

			logger.Error("error sending message", "error", err.Error())
			return err
		}

		err = workflow.ExecuteActivity(ctx, messageRepositoryActivities.UpdateMessageStatus, &activities.UpdateMessageStatusInput{
			Id:     message.Id.String(),
			Status: messages.StatusProcessing,
		}).Get(ctx, nil)
		if err != nil {
			logger.Error("error updating message status", "error", err.Error())
			return err
		}
	}

	return nil
}
