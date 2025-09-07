package provider

import (
	"github.com/umutcomlekci/automated-messaging-system/internal/services/types"
	"github.com/umutcomlekci/automated-messaging-system/pkg/fasthttp"
)

type httpProvider struct {
	url string
}

func newHttpProvider(url string) Provider {
	return &httpProvider{
		url: url,
	}
}

func (p *httpProvider) Send(message *types.SmsMessage) (*types.SmsResult, error) {
	response, err := fasthttp.SendGetRequest(p.url, nil)
	if err != nil {
		return nil, err
	} else if response.StatusCode != 200 {
		return nil, err
	}

	return &types.SmsResult{
		SmsId: response.Header["X-Request-Id"],
	}, nil
}
