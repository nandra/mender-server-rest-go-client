package mender_rest_api_client

import (
	"encoding/json"
	"fmt"
	"path"
	"time"
)

const deviceInventoryBasePath = "/api/management/v1/inventory"

type DeviceGroup struct {
	Group string `json:"group"`
}

type DeviceInventoryList []struct {
	ID         string `json:"id"`
	Attributes []struct {
		Name        string `json:"name"`
		Scope       string `json:"scope"`
		Value       string `json:"value"`
		Description string `json:"description"`
	} `json:"attributes"`
	UpdatedTs time.Time `json:"updated_ts"`
}

type DeviceInventory struct {
	ID         string `json:"id"`
	Attributes []struct {
		Name        string `json:"name"`
		Scope       string `json:"scope"`
		Value       string `json:"value"`
		Description string `json:"description"`
	} `json:"attributes"`
	UpdatedTs time.Time `json:"updated_ts"`
}

type DeviceGroupData struct {
	Group string `json:"group"`
}

// List devices inventories
// TODO: support for queries
func (c *RestClient) ListDeviceInventories() (DeviceInventoryList, error) {
	var devInventory DeviceInventoryList = DeviceInventoryList{}
	resp, err := c.client.R().Get(path.Join(deviceInventoryBasePath, "devices"))
	if err != nil {
		return devInventory, err
	}

	if err = json.Unmarshal(resp.Body(), &devInventory); err != nil {
		return devInventory, err
	}

	return devInventory, nil
}

// Get a selected device's inventory
// TODO: test
func (c *RestClient) GetDeviceInventory(deviceId string) (DeviceInventory, error) {
	var devInventory DeviceInventory = DeviceInventory{}
	resp, err := c.client.R().Get(path.Join(deviceInventoryBasePath, "devices", deviceId))
	if err != nil {
		return devInventory, err
	}

	if err = json.Unmarshal(resp.Body(), &devInventory); err != nil {
		return devInventory, err
	}

	return devInventory, nil
}

//Remove selected device's inventory
// TODO: test
func (c *RestClient) DeleteDeviceInventory(deviceId string) error {
	_, err := c.client.R().Delete(path.Join(deviceInventoryBasePath, "devices", deviceId))
	if err != nil {
		return err
	}

	return nil
}

// Get a selected device's group
// TODO: test
func (c *RestClient) GetDeviceGroup(deviceId string) (DeviceGroupData, error) {
	var group DeviceGroupData = DeviceGroupData{}
	resp, err := c.client.R().Get(path.Join(deviceInventoryBasePath, "devices", deviceId, "group"))
	if err != nil {
		return group, err
	}

	if err = json.Unmarshal(resp.Body(), &group); err != nil {
		return group, err
	}

	return group, nil
}

// Add a device to a group
// TODO: test
func (c *RestClient) AssignGroup(deviceId, groupName string) error {
	group := DeviceGroupData{
		Group: groupName,
	}

	g, e := json.Marshal(group)
	if e != nil {
		return fmt.Errorf("Failed to marshall group %v", e)
	}

	_, err := c.client.R().SetBody(g).Put(path.Join(deviceInventoryBasePath, "devices", deviceId, "group"))
	if err != nil {
		return err
	}

	return nil
}

// Remove a device from a group
// TODO: test
func (c *RestClient) ClearGroup(deviceId, groupName string) error {
	_, err := c.client.R().Delete(path.Join(deviceInventoryBasePath, "devices", deviceId, "group", groupName))
	if err != nil {
		return err
	}

	return nil
}

// List all groups existing device groups
// TODO: test
func (c *RestClient) ListGroups() ([]string, error) {
	var listGroups []string = []string{}
	resp, err := c.client.R().Get(path.Join(deviceInventoryBasePath, "groups"))
	if err != nil {
		return listGroups, err
	}

	if err = json.Unmarshal(resp.Body(), &listGroups); err != nil {
		return listGroups, err
	}

	return listGroups, nil
}

// List the devices belonging to a given group
// TODO: test
func (c *RestClient) GetDevicesInGroup(groupName string) ([]string, error) {
	var listDevicesInGroup []string = []string{}
	resp, err := c.client.R().Get(path.Join(deviceInventoryBasePath, "groups", groupName, "devices"))
	if err != nil {
		return listDevicesInGroup, err
	}

	if err = json.Unmarshal(resp.Body(), &listDevicesInGroup); err != nil {
		return listDevicesInGroup, err
	}

	return listDevicesInGroup, nil
}

// Add devices to group
// TODO: test
func (c *RestClient) AddDevicesToGroup(groupName string, devices []string) error {

	d, e := json.Marshal(devices)
	if e != nil {
		return fmt.Errorf("Failed to marshall group %v", e)
	}

	_, err := c.client.R().SetBody(d).Patch(path.Join(deviceInventoryBasePath, "groups", groupName, "devices"))
	if err != nil {
		return err
	}

	return nil
}

// Clear devices' group
// TODO: test
func (c *RestClient) RemoveDevicesFromGroup(groupName string, devices []string) error {

	d, e := json.Marshal(devices)
	if e != nil {
		return fmt.Errorf("Failed to marshall group %v", e)
	}

	_, err := c.client.R().SetBody(d).Delete(path.Join(deviceInventoryBasePath, "groups", groupName, "devices"))
	if err != nil {
		return err
	}

	return nil
}
