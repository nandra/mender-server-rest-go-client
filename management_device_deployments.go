package mender_rest_api_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

const deviceDeploymentsBasePath = "/api/management/v1/deployments"

type ListDeployments []struct {
	Created      time.Time `json:"created"`
	Status       string    `json:"status"`
	Name         string    `json:"name"`
	ArtifactName string    `json:"artifact_name"`
	ID           string    `json:"id"`
	Finished     time.Time `json:"finished"`
	DeviceCount  int       `json:"device_count"`
	Retries      int       `json:"retries"`
}

type DeploymentStatistics struct {
	Success          int `json:"success"`
	Pending          int `json:"pending"`
	Failure          int `json:"failure"`
	Downloading      int `json:"downloading"`
	Installing       int `json:"installing"`
	Rebooting        int `json:"rebooting"`
	Noartifact       int `json:"noartifact"`
	AlreadyInstalled int `json:"already-installed"`
	Aborted          int `json:"aborted"`
}

type DeploymentStatus struct {
	Created      time.Time `json:"created"`
	Status       string    `json:"status"`
	Name         string    `json:"name"`
	ArtifactName string    `json:"artifact_name"`
	ID           string    `json:"id"`
	Finished     time.Time `json:"finished"`
	Phases       []struct {
		BatchSize   int       `json:"batch_size,omitempty"`
		StartTs     time.Time `json:"start_ts,omitempty"`
		DeviceCount int       `json:"device_count,omitempty"`
	} `json:"phases"`
	DeviceCount int `json:"device_count"`
	Retries     int `json:"retries"`
}

type DeploymentStatusList []struct {
	ID         string    `json:"id"`
	Finished   time.Time `json:"finished"`
	Status     string    `json:"status"`
	Created    time.Time `json:"created"`
	DeviceType string    `json:"device_type"`
	Log        bool      `json:"log"`
	State      string    `json:"state"`
	Substate   string    `json:"substate"`
}

type GroupDeployment struct {
	Name         string `json:"name"`
	ArtifactName string `json:"artifact_name"`
}

type ListReleases []struct {
	Name      string `json:"name"`
	Artifacts []struct {
		Name                  string    `json:"name"`
		Description           string    `json:"description"`
		DeviceTypesCompatible []string  `json:"device_types_compatible"`
		ID                    string    `json:"id"`
		Signed                bool      `json:"signed"`
		Modified              time.Time `json:"modified"`
		Info                  struct {
			TypeInfo struct {
				Type string `json:"type"`
			} `json:"type_info"`
		} `json:"info"`
		Files []struct {
			Name     string `json:"name"`
			Checksum string `json:"checksum"`
			Size     int    `json:"size"`
			Date     string `json:"date"`
		} `json:"files"`
		Metadata struct {
		} `json:"metadata"`
	} `json:"artifacts"`
}

type ArtifactInfo struct {
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	DeviceTypesCompatible []string  `json:"device_types_compatible"`
	ID                    string    `json:"id"`
	Signed                bool      `json:"signed"`
	Modified              time.Time `json:"modified"`
	Info                  struct {
		TypeInfo struct {
			Type string `json:"type"`
		} `json:"type_info"`
	} `json:"info"`
	Files []struct {
		Name     string `json:"name"`
		Checksum string `json:"checksum"`
		Size     int    `json:"size"`
		Date     string `json:"date"`
	} `json:"files"`
	Metadata struct {
	} `json:"metadata"`
}

type StorageUsage struct {
	Limit int `json:"limit"`
	Usage int `json:"usage"`
}

// Find all deployments
// TODO: add support for queries
func (c *Client) ListDeployments() (ListDeployments, error) {
	var list ListDeployments = ListDeployments{}
	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "deployments"))
	if err != nil {
		return list, err
	}

	if err = json.Unmarshal(resp.Body(), &list); err != nil {
		return list, err
	}
	return list, nil
}

// Create a deployment
func (c *Client) CreateDeployment(deploymentName, artifactName string, devices []string, retries int) error {

	type Deployment struct {
		Name         string   `json:"name"`
		ArtifactName string   `json:"artifact_name"`
		Devices      []string `json:"devices"`
		Retries      int      `json:"retries"`
	}

	deployment := Deployment{
		Name:         deploymentName,
		ArtifactName: artifactName,
		Devices:      devices,
		Retries:      retries,
	}

	// marshal deployment to json
	d, err := json.Marshal(deployment)
	if err != nil {
		return err
	}

	_, err = c.client.R().SetBody(d).Post(path.Join(deviceDeploymentsBasePath, "deployments"))
	if err != nil {
		return err
	}

	return nil
}

// Create a deployment for a group of devices
func (c *Client) CreateDeploymentForGroup(deploymentName, artifactName, groupName string) error {

	deployment := GroupDeployment{
		Name:         deploymentName,
		ArtifactName: artifactName,
	}

	// marshal deployment to json
	d, err := json.Marshal(deployment)
	if err != nil {
		return err
	}

	_, err = c.client.R().SetBody(d).Post(path.Join(deviceDeploymentsBasePath, "deployments/group", groupName))
	if err != nil {
		return err
	}

	return nil
}

// Get the details of a selected deployment
func (c *Client) ShowDeployment(deploymentId string) (DeploymentStatus, error) {
	var stat DeploymentStatus = DeploymentStatus{}

	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "deployments", deploymentId))
	if err != nil {
		fmt.Println("failed to read list of devices", err)
	}

	if err = json.Unmarshal(resp.Body(), &stat); err != nil {
		return stat, err
	}

	return stat, nil
}

// Abort the deployment
func (c *Client) AbortDeployment(deploymentId string) error {
	type AbortDeploymentBody struct {
		Status string `json:"status"`
	}

	abort := AbortDeploymentBody{Status: "aborted"}

	// marshal deployment to json
	a, err := json.Marshal(abort)
	if err != nil {
		return err
	}

	_, err = c.client.R().SetBody(a).Put(path.Join(deviceDeploymentsBasePath, "deployments", deploymentId, "status"))
	if err != nil {
		return err
	}

	return nil
}

// Get status count for all devices in a deployment.
func (c *Client) DeploymentStatistics(deploymentId string) (DeploymentStatistics, error) {
	var stat DeploymentStatistics = DeploymentStatistics{}
	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "deployments", deploymentId, "statistics"))
	if err != nil {
		return stat, err
	}

	if err = json.Unmarshal(resp.Body(), &stat); err != nil {
		return stat, err
	}
	return stat, nil
}

// Get list of all devices and their respective status for the deployment with the given ID.
func (c *Client) ListDevicesInDeployment(deploymentId string) (DeploymentStatusList, error) {
	var list DeploymentStatusList = DeploymentStatusList{}
	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "deployments", deploymentId, "devices"))
	if err != nil {
		return list, err
	}

	if err = json.Unmarshal(resp.Body(), &list); err != nil {
		return list, err
	}

	return list, nil
}

// Get the list of device IDs being part of the deployment.
func (c *Client) ListDevicesIDsInDeployment(deploymentId string) ([]string, error) {
	var list []string = []string{}
	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "deployments", deploymentId, "device_list"))
	if err != nil {
		return list, err
	}

	if err = json.Unmarshal(resp.Body(), &list); err != nil {
		return list, err
	}
	return list, nil
}

// Get the log of a selected device's deployment
// TODO: test
func (c *Client) GetDeploymentLogForDevice(deploymentId, deviceId string) (string, error) {
	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "deployments", deploymentId, "devices", deviceId, "log"))
	if err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}

// Remove device from all deployments
func (c *Client) RemoveDeviceFromDeployment(deviceId string) error {
	_, err := c.client.R().Delete(path.Join(deviceDeploymentsBasePath, "deployments/devices", deviceId))
	if err != nil {
		return err
	}

	return nil
}

// List releases
func (c *Client) ListReleases() (ListReleases, error) {
	var releases ListReleases = ListReleases{}
	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "deployments/releases"))
	if err != nil {
		return releases, err
	}

	if err = json.Unmarshal(resp.Body(), &releases); err != nil {
		return releases, err
	}

	return releases, nil
}

// List known artifacts
func (c *Client) ListArtifacts() (ListReleases, error) {
	var releases ListReleases = ListReleases{}
	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "artifacts"))
	if err != nil {
		return releases, err
	}

	if err = json.Unmarshal(resp.Body(), &releases); err != nil {
		return releases, err
	}

	return releases, nil
}

// Upload mender artifact
func (c *Client) UploadArtifacts(artifactFilePath, artifactDescription string) error {

	artifact, err := ioutil.ReadFile(artifactFilePath)
	if err != nil {
		return fmt.Errorf("Failed to read file: %v", err)
	}

	fi, err := os.Stat(artifactFilePath)
	if err != nil {
		return err
	}
	// get the size
	size := fi.Size()

	_, err = c.client.R().
		SetFormData(map[string]string{
			"size":        strconv.FormatInt(size, 10),
			"description": artifactDescription,
		}).
		SetFileReader("artifact", fi.Name(), bytes.NewReader(artifact)).
		Post(path.Join(deviceDeploymentsBasePath, "artifacts"))
	if err != nil {
		return err
	}

	return nil
}

// TODO: implement
func (c *Client) GenerateArtifact() error {
	return fmt.Errorf("Not implmplemented")
}

// Get the details of a selected artifact
func (c *Client) ShowArtifact(artifactId string) (ArtifactInfo, error) {
	var artifact ArtifactInfo = ArtifactInfo{}
	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "artifacts", artifactId))
	if err != nil {
		return artifact, err
	}

	if err = json.Unmarshal(resp.Body(), &artifact); err != nil {
		return artifact, err
	}

	return artifact, nil
}

// Update description of a selected artifact
func (c *Client) UpdateArtifactinfo(artifactId, description string) error {
	type ArtifactDescription struct {
		Description string `json:"description"`
	}

	desc := ArtifactDescription{Description: description}

	// marshal deployment to json
	d, err := json.Marshal(desc)
	if err != nil {
		return err
	}

	_, err = c.client.R().SetBody(d).Put(path.Join(deviceDeploymentsBasePath, "artifacts", artifactId))
	if err != nil {
		return err
	}

	return nil
}

// Delete the artifact
func (c *Client) DeleteArtifact(artifactId string) error {
	_, err := c.client.R().Delete(path.Join(deviceDeploymentsBasePath, "artifacts", artifactId))
	if err != nil {
		return err
	}

	return nil
}

// Get the download link of a selected artifact
func (c *Client) DownloadArtifact(artifactId string) (string, error) {
	type ArtifactResult struct {
		URI    string    `json:"uri"`
		Expire time.Time `json:"expire"`
	}

	var result ArtifactResult = ArtifactResult{}

	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "artifacts", artifactId, "download"))
	if err != nil {
		return "", err
	}

	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	return result.URI, nil
}

// Get storage limit and current storage usage
func (c *Client) GetStorageUsage() (StorageUsage, error) {
	var usage StorageUsage = StorageUsage{}
	resp, err := c.client.R().Get(path.Join(deviceDeploymentsBasePath, "limits/storage"))
	if err != nil {
		return usage, err
	}

	if err = json.Unmarshal(resp.Body(), &usage); err != nil {
		return usage, err
	}

	return usage, nil
}
