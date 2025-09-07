package fasthttp

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	body       []byte
	url        string
	StatusCode int
	Header     map[string]string
}

func (r *Response) Url() string {
	return r.url
}

func (r *Response) EndStruct(body interface{}) error {
	if !json.Valid(r.body) {
		return fmt.Errorf("json unmarshall failed, url: %s, body: %s", r.url, r.RawBodyString())
	}

	return json.Unmarshal(r.body, &body)
}

func (r *Response) RawBody() []byte {
	return r.body
}

func (r *Response) RawBodyString() string {
	return string(r.body)
}
