package mender_rest_api_client

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/go-resty/resty/v2"
)

// base path
const deviceAuthBasePath = "/api/management/v2/devauth/"

// json responses
type ListDevices []struct {
	ID         string `json:"id"`
	Attributes []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"attributes"`
	UpdatedTs time.Time `json:"updated_ts"`
}

type Device struct {
	ID           string `json:"id"`
	IdentityData struct {
		Mac string `json:"mac"`
		Sku string `json:"sku"`
		Sn  string `json:"sn"`
	} `json:"identity_data"`
	Status    string    `json:"status"`
	CreatedTs time.Time `json:"created_ts"`
	UpdatedTs time.Time `json:"updated_ts"`
	AuthSets  []struct {
		ID           string `json:"id"`
		Pubkey       string `json:"pubkey"`
		IdentityData struct {
			Mac string `json:"mac"`
			Sku string `json:"sku"`
			Sn  string `json:"sn"`
		} `json:"identity_data"`
		Status string    `json:"status"`
		Ts     time.Time `json:"ts"`
	} `json:"auth_sets"`
	Decommissioning bool `json:"decommissioning"`
}

type DevicesCount struct {
	Count int `json:"count"`
}

func checkAndReturnError(r *resty.Response, e error) error {
	// check response error
	if r.IsError() {
		return fmt.Errorf("Error response:%v", r.StatusCode())
	}
	// check error
	if e != nil {
		return e
	}

	return nil
}

// List devices sorted by age and optionally filter on device status
// TODO: implement page, per_pare queries
func (c *RestClient) ListDevices() (ListDevices, error) {
	var devices ListDevices = ListDevices{}
	resp, err := c.client.R().Get(path.Join(deviceAuthBasePath, "devices"))
	if err = checkAndReturnError(resp, err); err != nil {
		return devices, err
	}

	if err = json.Unmarshal(resp.Body(), &devices); err != nil {
		return devices, err
	}

	return devices, nil
}

// Submit a preauthorized device.
// TODO: implement
func (c *RestClient) Preauthorize() error {
	return fmt.Errorf("Not implmplemented")
}

// Get a particular device.
func (c *RestClient) GetDevice(deviceId string) (Device, error) {
	var device Device = Device{}
	resp, err := c.client.R().Get(path.Join(deviceAuthBasePath, "devices", deviceId))
	if err = checkAndReturnError(resp, err); err != nil {
		return device, err
	}

	if err = json.Unmarshal(resp.Body(), &device); err != nil {
		return device, err
	}

	return device, nil
}

// Remove device and associated authentication set
// TODO: test
func (c *RestClient) DecomisionDevice(deviceId string) error {
	resp, err := c.client.R().Delete(path.Join(deviceAuthBasePath, "devices", deviceId))
	if err = checkAndReturnError(resp, err); err != nil {
		return err
	}

	return nil
}

// Remove the device authentication set
// TODO: test
func (c *RestClient) RejectAuthtentication(deviceId, authId string) error {
	resp, err := c.client.R().Delete(path.Join(deviceAuthBasePath, "devices", deviceId, "auth", authId))
	if err = checkAndReturnError(resp, err); err != nil {
		return err
	}

	return nil
}

// Update the device authentication set status
// TODO: test
func (c *RestClient) SetAuthtenticationStatus(deviceId, authId string) error {
	resp, err := c.client.R().Put(path.Join(deviceAuthBasePath, "devices", deviceId, authId, "status"))
	if err = checkAndReturnError(resp, err); err != nil {
		return err
	}

	return nil
}

// Get the device authentication set status
// TODO: test
func (c *RestClient) GetAuthtenticationStatus(deviceId, authId string) (string, error) {
	type AuthStatus struct {
		Status string `json:"status"`
	}
	var status AuthStatus = AuthStatus{}

	resp, err := c.client.R().Get(path.Join(deviceAuthBasePath, "devices", deviceId, authId, "status"))
	if err = checkAndReturnError(resp, err); err != nil {
		return status.Status, err
	}

	if err = json.Unmarshal(resp.Body(), &status); err != nil {
		return status.Status, err
	}

	return status.Status, nil
}

// Count number of devices, optionally filtered by status.
// TODO: added support for queries
func (c *RestClient) CountDevices() (int, error) {
	var count DevicesCount = DevicesCount{}

	resp, err := c.client.R().Get(path.Join(deviceAuthBasePath, "devices/count"))
	if err = checkAndReturnError(resp, err); err != nil {
		return 0, err
	}

	if err = json.Unmarshal(resp.Body(), &count); err != nil {
		return 0, err
	}

	return count.Count, nil
}

// Revoke JWT with given id
// TODO: implement
func (c *RestClient) RewokeAPIToken() error {
	return fmt.Errorf("Not implmplemented")
}

// Obtain limit of accepted devices.
// TODO: test
func (c *RestClient) GetDeviceLimit() (int, error) {
	type Limit struct {
		Limit int `json:"limit"`
	}

	var limit Limit

	resp, err := c.client.R().Get(path.Join(deviceAuthBasePath, "limits/max_devices"))
	if err = checkAndReturnError(resp, err); err != nil {
		return 0, err
	}

	if err = json.Unmarshal(resp.Body(), &limit); err != nil {
		return 0, err
	}

	return limit.Limit, nil
}
