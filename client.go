package menderrc

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const BaseApiUrl = "https://hosted.mender.io"

// Client wraps http client, ...
type Client struct {
	jwtToken   string
	baseApiUrl string
	httpClient *http.Client
}

// NewClient creates client to Mender Rest API.
func NewClient() *Client {
	return NewCustomClient(BaseApiUrl, http.DefaultClient)
}

// NewCustomClient
func NewCustomClient(baseApiUrl string, httpClient *http.Client) *Client {
	return &Client{
		baseApiUrl: baseApiUrl,
		httpClient: httpClient,
	}
}

func (c *Client) newAuthorizedRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {

	if len(c.jwtToken) == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.jwtToken)

	return req, nil
}

func (c *Client) p(scope string, format string, args ...interface{}) string {
	return c.baseApiUrl + scope + fmt.Sprintf(format, args...)
}
