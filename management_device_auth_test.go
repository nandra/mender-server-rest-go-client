package mender_rest_api_client

import (
	"testing"

	//mender_rest_api_client "login/mender"

	"github.com/jarcoal/httpmock"
)

const serverUrl = "https://test_mender.com"

func restartHttpMock(action, path, jsonResponse string, response int) *RestClient {
	httpmock.DeactivateAndReset()
	httpmock.Activate()

	c := NewRestClient(serverUrl, "", "", true)

	httpmock.ActivateNonDefault(c.client.GetClient())

	httpmock.RegisterResponder(action, path,
		httpmock.NewStringResponder(response, jsonResponse))

	return c
}

func TestListDevices(t *testing.T) {
	c := restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices"), `[
		{
		  "id": "string",
		  "identity_data": {
			"mac": "00:01:02:03:04:05",
			"sku": "My Device 1",
			"sn": "SN1234567890"
		  },
		  "status": "pending",
		  "created_ts": "2019-08-24T14:15:22Z",
		  "updated_ts": "2019-08-24T14:15:22Z",
		  "auth_sets": [
			{
			  "id": "string",
			  "pubkey": "string",
			  "identity_data": {
				"mac": "00:01:02:03:04:05",
				"sku": "My Device 1",
				"sn": "SN1234567890"
			  },
			  "status": "pending",
			  "ts": "2019-08-24T14:15:22Z"
			}
		  ],
		  "decommissioning": true
		}
	  ]`, 200)

	d, e := c.ListDevices()
	if e != nil {
		t.Error(e)
	}

	if d[0].ID != "string" {
		t.Errorf("Invalid data")
	}

	// error response
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices"), `{}`, 400)
	d, e = c.ListDevices()
	if e == nil {
		t.Error(e)
	}

	// invalid json
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices"), `{`, 200)
	d, e = c.ListDevices()
	if e == nil {
		t.Error(e)
	}
}

func TestGetDevice(t *testing.T) {
	deviceId := "12345"
	c := restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices/"+deviceId), `{
		"id": "12345",
		"identity_data": {
		  "mac": "00:01:02:03:04:05",
		  "sku": "My Device 1",
		  "sn": "SN1234567890"
		},
		"status": "pending",
		"created_ts": "2019-08-24T14:15:22Z",
		"updated_ts": "2019-08-24T14:15:22Z",
		"auth_sets": [
		  {
			"id": "string",
			"pubkey": "string",
			"identity_data": {
			  "mac": "00:01:02:03:04:05",
			  "sku": "My Device 1",
			  "sn": "SN1234567890"
			},
			"status": "pending",
			"ts": "2019-08-24T14:15:22Z"
		  }
		],
		"decommissioning": true
	  }`, 200)

	d, e := c.GetDevice(deviceId)
	if e != nil {
		t.Error(e)
	}

	if d.ID != deviceId {
		t.Errorf("Invalid id's")
	}

	// error response
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices/"+deviceId), `{}`, 404)
	d, e = c.GetDevice(deviceId)
	if e == nil {
		t.Error(e)
	}

	// invalid json
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices/"+deviceId), `{`, 200)
	d, e = c.GetDevice(deviceId)
	if e == nil {
		t.Error(e)
	}
}

func TestDecomisionDevice(t *testing.T) {
	// ok response
	deviceId := "123456"
	c := restartHttpMock("DELETE", joinURL(deviceAuthBasePath, "/devices/"+deviceId), `{}`, 204)
	e := c.DecomisionDevice(deviceId)
	if e != nil {
		t.Error(e)
	}

	// device not exists
	c = restartHttpMock("DELETE", joinURL(deviceAuthBasePath, "/devices/"+deviceId), `{}`, 404)
	e = c.DecomisionDevice(deviceId)
	if e == nil {
		t.Error(e)
	}

}

func TestRejectAuthtentication(t *testing.T) {
	// ok response
	deviceId := "123456"
	authId := "654321"
	c := restartHttpMock("DELETE", joinURL(deviceAuthBasePath, "/devices/"+deviceId+"/auth/"+authId), `{}`, 204)
	e := c.RejectAuthtentication(deviceId, authId)
	if e != nil {
		t.Error(e)
	}

	// device not exists
	c = restartHttpMock("DELETE", joinURL(deviceAuthBasePath, "/devices/"+deviceId+"/auth/"+authId), `{}`, 404)
	e = c.RejectAuthtentication(deviceId, authId)
	if e == nil {
		t.Error(e)
	}

}

func TestSetAuthtenticationStatus(t *testing.T) {
	// ok response
	deviceId := "123456"
	authId := "654321"
	c := restartHttpMock("PUT", joinURL(deviceAuthBasePath, "/devices/"+deviceId+"/"+authId+"/status"), `{}`, 204)
	e := c.SetAuthtenticationStatus(deviceId, authId)
	if e != nil {
		t.Error(e)
	}

	// device not exists
	c = restartHttpMock("PUT", joinURL(deviceAuthBasePath, "/devices/"+deviceId+"/"+authId+"/status"), `{}`, 404)
	e = c.SetAuthtenticationStatus(deviceId, authId)
	if e == nil {
		t.Error(e)
	}
}

func TestGetAuthtenticationStatus(t *testing.T) {
	// ok response
	deviceId := "123456"
	authId := "654321"
	c := restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices/"+deviceId+"/"+authId+"/status"), `{"status": "accepted"}`, 200)
	s, e := c.GetAuthtenticationStatus(deviceId, authId)
	if e != nil || s != "accepted" {
		t.Error(e)
	}

	// device not exists
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices/"+deviceId+"/"+authId+"/status"), `{ "error": "Ivalid status"}`, 404)
	_, e = c.GetAuthtenticationStatus(deviceId, authId)
	if e == nil {
		t.Error(e)
	}

	// wrong json
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices/"+deviceId+"/"+authId+"/status"), `{ "error": "Ivalid status`, 200)
	_, e = c.GetAuthtenticationStatus(deviceId, authId)
	if e == nil {
		t.Error(e)
	}
}

func TestCountDevices(t *testing.T) {
	// ok response
	c := restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices/count"), `{"count": 42}`, 200)
	s, e := c.CountDevices()
	if e != nil || s != 42 {
		t.Error(e)
	}

	// device not exists
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices/count"), `{ "error": "Ivalid status"}`, 400)
	_, e = c.CountDevices()
	if e == nil {
		t.Error(e)
	}

	// wrong json
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/devices/count"), `{"count": 4`, 200)
	_, e = c.CountDevices()
	if e == nil {
		t.Error(e)
	}
}

func TestGetDeviceLimit(t *testing.T) {
	// ok response
	c := restartHttpMock("GET", joinURL(deviceAuthBasePath, "/limits/max_devices"), `{"limit": 123}`, 200)
	s, e := c.GetDeviceLimit()
	if e != nil || s != 123 {
		t.Error(e)
	}

	// device not exists
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/limits/max_devices"), "", 500)
	_, e = c.GetDeviceLimit()
	if e == nil {
		t.Error(e)
	}

	// wrong json
	c = restartHttpMock("GET", joinURL(deviceAuthBasePath, "/limits/max_devices"), `{"linut": 4`, 200)
	_, e = c.GetDeviceLimit()
	if e == nil {
		t.Error(e)
	}
}
