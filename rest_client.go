package mender_rest_api_client

import (
	"crypto/tls"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	jwtToken      string
	username      string
	password      string
	serverUrl     string
	tlsSkipVerify bool
	client        *resty.Client
}

func NewClient(url, user, pass string, skipVerify bool) *Client {
	return &Client{serverUrl: url, username: user, password: pass, tlsSkipVerify: skipVerify, client: resty.New()}
}

//
// Login to mender server using username + password
// debug - enable rest debugging
// skipTlsVerify - disable ssl verification

func (c *Client) Login(debug, skipTlsVerify bool) error {

	client := resty.New()

	client.SetHostURL(c.serverUrl)
	// Basic Auth for all request
	client.SetBasicAuth(c.username, c.password)
	if skipTlsVerify {
		client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: c.tlsSkipVerify})
	}
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Post("/api/management/v1/useradm/auth/login")

	if err != nil {
		return err
	}

	c.jwtToken = string(resp.Body())

	// resty setup
	c.client.SetHostURL(c.serverUrl)
	// Basic Auth for all request
	if skipTlsVerify {
		c.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: c.tlsSkipVerify})
	}
	// jwt bearer token
	c.client.SetHeader("Authorization", "Bearer "+c.jwtToken)

	if debug {
		c.client.SetDebug(true)
	}

	return nil
}
