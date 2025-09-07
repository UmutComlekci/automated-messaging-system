package types

type (
	SmsMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Content string `json:"content"`
	}

	SmsResult struct {
		SmsId string `json:"sms_id"`
	}
)
