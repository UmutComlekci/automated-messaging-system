package provider

import (
	"fmt"

	"github.com/umutcomlekci/automated-messaging-system/internal/config"
	"github.com/umutcomlekci/automated-messaging-system/internal/services/types"
)

type Provider interface {
	Send(*types.SmsMessage) (*types.SmsResult, error)
}

func NewProvider() (Provider, error) {
	if config.GetProvider() == "http" {
		return newHttpProvider(config.GetProviderHttpUrl()), nil
	}

	return nil, fmt.Errorf("provider %s not found", config.GetProvider())
}
