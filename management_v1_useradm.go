package menderrc

import (
	"context"
	"io/ioutil"
	"net/http"
)

const v1useradmScope = "/api/management/v1/useradm"

// Login logins to mender server using username and password
func (c *Client) Login(ctx context.Context, username, password string) error {

	headers := map[string][]string{
		"Content-Type": []string{"application/json"},
		"Accept":       []string{"application/jwt"},
	}

	req, err := http.NewRequest("POST", c.p(v1useradmScope, "/auth/login"), nil)
	if err != nil {
		return err
	}

	req.Header = headers
	req.SetBasicAuth(username, password)
	req = req.WithContext(ctx)

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
