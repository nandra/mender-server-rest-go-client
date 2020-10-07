package mender_rest_api_client

import (
	"io/ioutil"
	"net/http"
)

const BaseApiUrl = "https://hosted.mender.io"

// Client wraps http client, ...
type Client struct {
	jwtToken   string
	baseApiUrl string
	httpClient *http.Client
}

type ClientOpt func(*Client)

// SetHttpClient allows set custom http client
func SetHttpClient(httpClient *http.Client) ClientOpt {
	return func(c *Client) {
		if httpClient {
			panic("http client cannot be nil")
		}
		c.httpClient = httpClient
	}
}

// SetBaseApiUrl allows
func SetBaseApiUrl(baseApiUrl string) ClientOpt {
	return func(c *Client) {
		c.baseApiUrl = baseApiUrl
	}
}

// NewClient creates client to Mender Rest API. User might use custom http client, if nil is passed
// default http client will be used.
func NewClient(opts ...ClientOpt) *Client {

	// use default http client if none is set
	client := &Client{
		baseApiUrl: BaseApiUrl,
		httpClient: http.DefaultClient,
	}

	// user configurable overrides
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// Login logins to mender server using username and password
func (c *Client) Login(username, password string) error {

	headers := map[string][]string{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/jwt"},
	}

	req, err := http.NewRequest("POST", c.baseApiUrl+"/api/management/v1/useradm/auth/login", nil)
	if err != nil {
		return err
	}

	req.Header = headers
	req.SetBasicAuth(username, password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// save token for later use
	c.jwtToken = string(body)

	return nil
}
