package fasthttp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"time"

	valyala "github.com/valyala/fasthttp"
)

const (
	AuthorizationHeader string = "Authorization"
	JsonContentType     string = "application/json"
	MultipartForm       string = "multipart/form-data"
)

var (
	client *valyala.Client
)

func init() {
	client = &valyala.Client{
		NoDefaultUserAgentHeader: true,
		MaxIdleConnDuration:      10 * time.Second,
		ReadTimeout:              30 * time.Second,
		MaxConnsPerHost:          2048,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

func SendRequest(url, method string, body any, headers map[string]string) (*Response, error) {
	return sendRequest(url, method, body, headers)
}

func SendGetRequest(url string, headers map[string]string) (*Response, error) {
	return sendRequest(url, valyala.MethodGet, nil, headers)
}

func SendPostRequest(url string, body any, headers map[string]string) (*Response, error) {
	return sendRequest(url, valyala.MethodPost, body, headers)
}

func SendDeleteRequest(url string, body any, headers map[string]string) (*Response, error) {
	return sendRequest(url, valyala.MethodDelete, body, headers)
}

func SendPutRequest(url string, body any, headers map[string]string) (*Response, error) {
	return sendRequest(url, valyala.MethodPut, body, headers)
}

func SendPurgeRequest(url string, body any, headers map[string]string) (*Response, error) {
	return sendRequest(url, "PURGE", body, headers)
}

func sendRequest(url, method string, body any, headers map[string]string) (*Response, error) {
	req := valyala.AcquireRequest()

	defer valyala.ReleaseRequest(req)
	req.SetRequestURI(url)

	req.Header.SetMethod(method)
	req.Header.SetContentType(JsonContentType)

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	if method != valyala.MethodGet && body != nil {
		reqBodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req.SetBody(reqBodyBytes)
	}

	resp := valyala.AcquireResponse()
	defer valyala.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(resp.Body())
	byteArray, err := io.ReadAll(buffer)
	if err != nil {
		return nil, err
	}

	header := map[string]string{}
	for key, value := range resp.Header.All() {
		header[string(key)] = string(value)
	}

	return &Response{
		url:        url,
		body:       byteArray,
		Header:     header,
		StatusCode: resp.StatusCode(),
	}, nil
}
