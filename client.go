package freshdesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type Client interface {
	Contacts() ContactsClient
	BaseUrl() string
}

type client struct {
	apiKey   string
	baseURL  string
	contacts ContactsClient

	httpClient *http.Client
}

type Logger interface {
	retryablehttp.LeveledLogger
}

func NewClient(subdomain, apiKey string, log Logger) (Client, error) {
	rc := retryablehttp.NewClient()

	if log != nil {
		rc.Logger = log
	}

	c := &client{
		apiKey:     apiKey,
		baseURL:    fmt.Sprintf("https://%s.freshdesk.com/api/v2/", subdomain),
		httpClient: rc.StandardClient(),
	}
	c.contacts = &contactsClient{c}
	return c, nil
}

func (c *client) BaseUrl() string {
	return c.baseURL
}

func (c *client) Contacts() ContactsClient {
	return c.contacts
}

func (c *client) newRequest(method, endpoint string, body interface{}) (req *Request, err error) {
	b := make([]byte, 0)
	if body != nil {
		if b, err = json.Marshal(&body); err != nil {
			return
		}
	}

	bodyReader := bytes.NewReader(b)
	var raw *http.Request
	if raw, err = http.NewRequest(method, c.baseURL+endpoint, bodyReader); err != nil {
		return
	}

	raw.SetBasicAuth(c.apiKey, "X")
	raw.Header.Add("Content-Type", "application/json")
	return &Request{raw}, nil
}

func (c *client) do(req *Request, out interface{}, expectedStatus int) error {
	raw, err := c.httpClient.Do(req.Request)
	if err != nil {
		return err
	}

	defer raw.Body.Close()

	res := &Response{raw}

	if res.StatusCode != expectedStatus {
		return NewApiError(
			res.StatusCode,
			expectedStatus,
			req.Payload(),
			res.Payload(),
		)
	}
	if out != nil {
		if err = json.NewDecoder(res.Body).Decode(out); err != nil {
			return err
		}
	}

	return nil
}

type Request struct {
	*http.Request
}

func (r *Request) Payload() string {
	if body, err := io.ReadAll(r.Body); err == nil {
		var jsonBuffer bytes.Buffer
		if err = json.Compact(&jsonBuffer, body); err == nil {
			return jsonBuffer.String()
		}
	}
	return ""
}

type Response struct {
	*http.Response
}

func (r *Response) Payload() string {
	if body, err := io.ReadAll(r.Body); err == nil {
		var jsonBuffer bytes.Buffer
		if err = json.Compact(&jsonBuffer, body); err == nil {
			return jsonBuffer.String()
		}
	}
	return ""
}
