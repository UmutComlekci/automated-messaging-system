package activities

import (
	"github.com/umutcomlekci/automated-messaging-system/internal/services/sms"
	"github.com/umutcomlekci/automated-messaging-system/internal/services/types"
)

type SmsServiceActivities struct {
	smsService *sms.Service
}

func NewSmSServiceActivities(smsService *sms.Service) *SmsServiceActivities {
	return &SmsServiceActivities{
		smsService: smsService,
	}
}

func (a *SmsServiceActivities) Send(message *types.SmsMessage) (*types.SmsResult, error) {
	result, err := a.smsService.Send(message)
	if err != nil {
		return nil, err
	}

	return result, nil
}
