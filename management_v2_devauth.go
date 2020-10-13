package menderrc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"
)

const v2devauthScope = "api/management/v2/devauth"

// IdentityData ...
type IdentityData struct {
	Mac string `json:"mac"`
	Sku string `json:"sku"`
	Sn  string `json:"sn"`
}

// Status is admission status of the device.
type Status string

const (
	Pending       Status = "pending"
	Accepted             = "accepted"
	Rejected             = "rejected"
	Preauthorized        = "preauthorized"
	NoAuth               = "noauth"
	AllDevices           = "all devices"
)

// AuthSet ...
type AuthSet struct {
	ID           string       `json:"id"`
	Pubkey       string       `json:"pubkey"`
	IdentityData IdentityData `json:"identity_data"`
	Status       Status       `json:"status"`
	Ts           time.Time    `json:"ts"`
}

// ListDevicesFilter ...
type ListDevicesFilter struct {
	Status  Status `json:"status"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
}

// ListDevices lists devices sorted by age and optionally filter on device status.
// If filter is nil all devices are listed.
func (c *Client) ListDevicesFilterBy(ctx context.Context, filter *ListDevicesFilter) ([]Device, error) {

	var body io.Reader

	if filter != nil {
		b, err := json.Marshal(filter)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(b)
	}

	req, err := c.newAuthorizedRequest(ctx, "GET", c.p(v2devauthScope, "/devices"), body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var devices []Device
	if err := json.NewDecoder(resp.Body).Decode(&devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// ListAllDevies lists all devices.
func (c *Client) ListDevices(ctx context.Context) ([]Device, error) {
	return c.ListDevicesFilterBy(ctx, nil)
}

type PreAuthSet struct {
	IdentityData IdentityData `json:"identity_data"`
	Pubkey       string       `json:"pubkey"`
}

// Preauthorize submits a preauthorized device.
func (c *Client) Preauthorize(ctx context.Context, preAuthSet PreAuthSet) error {
	panic("TODO")
}

// Device ...
type Device struct {
	ID              string       `json:"id"`
	IdentityData    IdentityData `json:"identity_data"`
	Status          Status       `json:"status"`
	CreatedTs       time.Time    `json:"created_ts"`
	UpdatedTs       time.Time    `json:"updated_ts"`
	AuthSets        []AuthSet    `json:"auth_sets"`
	Decommissioning bool         `json:"decommissioning"`
}

// GetDevice gets a particular device by device id.
func (c *Client) GetDevice(ctx context.Context, deviceId string) (*Device, error) {

	req, err := c.newAuthorizedRequest(ctx, "GET", c.p(v2devauthScope, "/devices/%s", deviceId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	device := &Device{}
	if err = json.NewDecoder(resp.Body).Decode(device); err != nil {
		return nil, err
	}

	return device, nil
}

// DecomisionDevice removes device and associated authentication set.
func (c *Client) DecomisionDevice(ctx context.Context, deviceId string) error {

	req, err := c.newAuthorizedRequest(ctx, "DELETE", c.p(v2devauthScope, "devices/%s", deviceId), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// TODO: check response error

	return nil
}

// RejectAuthtentication remove the device authentication set.
func (c *Client) RejectAuthtentication(ctx context.Context, deviceId, authId string) error {

	req, err := c.newAuthorizedRequest(ctx, "DELETE", c.p(v2devauthScope, "devices/%s/auth/%s", deviceId, authId), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// TODO: check response error

	return nil
}

// SetAuthtenticationStatus updates the device authentication status.
func (c *Client) SetAuthtenticationStatus(ctx context.Context, deviceId, authId string) error {

	req, err := c.newAuthorizedRequest(ctx, "PUT", c.p(v2devauthScope, "/devices/%s/auth/%s/status", deviceId, authId), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// TODO: check response error

	return nil
}

// GetAuthtenticationStatus gets the device authentication set status.
func (c *Client) GetAuthtenticationStatus(ctx context.Context, deviceId, authId string) (Status, error) {

	req, err := c.newAuthorizedRequest(ctx, "GET", c.p(v2devauthScope, "/devices/%s/auth/%s/status", deviceId, authId), nil)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// TODO: check response error

	var status struct {
		Status Status `json:"status"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return "", err
	}

	return status.Status, nil
}

// CountDevicesFilterBy counts number of devices filtered by status.
// Device status filter accepts one of 'pending', 'accepted', 'rejected', 'noauth' or 'all devices'.
func (c *Client) CountDevicesFilterBy(ctx context.Context, status Status) (int, error) {

	req, err := c.newAuthorizedRequest(ctx, "GET", c.p(v2devauthScope, "/devices/count"), nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	// TODO: check response error

	var count struct {
		Count int `json:"count"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&count); err != nil {
		return 0, err
	}

	return count.Count, nil
}

// CountDevices counts number of devices.
func (c *Client) CountDevices(ctx context.Context) (int, error) {
	return c.CountDevicesFilterBy(ctx, "all devices")
}

// RewokeAPIToken revokes JWT with given id.
func (c *Client) RewokeAPIToken(id string) error {
	panic("TODO")
}

// GetDeviceLimit obtains limit of accepted devices.
func (c *Client) GetDeviceLimit(ctx context.Context) (int, error) {

	req, err := c.newAuthorizedRequest(ctx, "GET", c.p(v2devauthScope, "/limits/max_devices"), nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	// TODO: check response error

	var limit struct {
		Limit int `json:"limit"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&limit); err != nil {
		return 0, err
	}

	return limit.Limit, nil
}
