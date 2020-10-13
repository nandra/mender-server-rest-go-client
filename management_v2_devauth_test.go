package menderrc

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func restartHttpMock(action, path, jsonResponse string, response int) *Client {
	httpmock.DeactivateAndReset()
	httpmock.Activate()

	httpc := http.DefaultClient
	c := NewCustomClient("https://test_mender.com", httpc)

	httpmock.ActivateNonDefault(httpc)

	responder := httpmock.NewStringResponder(response, jsonResponse)
	httpmock.RegisterResponder(action, path, responder)

	return c
}

func newAuthorizedClient() *Client {
	c := NewClient()

	if err := c.Login(context.Background(), "username", "password"); err != nil {

	}

	return c
}

func TestListDevices(t *testing.T) {

	tests := []struct {
		jsonResponse string
		status       int
		result       []Device
	}{
		{
			jsonResponse: `[{"id":"string","identity_data":{"mac":"00:01:02:03:04:05","sku":"My Device 1","sn":"SN1234567890"},"status":"pending","created_ts":"2019-08-24T14:15:22Z","updated_ts":"2019-08-24T14:15:22Z","auth_sets":[{"id":"string","pubkey":"string","identity_data":{"mac":"00:01:02:03:04:05","sku":"My Device 1","sn":"SN1234567890"},"status":"pending","ts":"2019-08-24T14:15:22Z"}],"decommissioning":true}]`,
			status:       200,
		},
		{
			jsonResponse: `{}`,
			status:       400,
		},
		{
			jsonResponse: `{`,
			status:       200,
		},
	}

	for _, test := range tests {

		c := restartHttpMock("GET", "/api/management/v2/devauth/devices", test.jsonResponse, test.status)

		ctx := context.Background()
		_, err := c.ListDevices(ctx)
		if err != nil {
			t.Error(err)
		}

		// TODO: check result
	}
}

func TestGetDevice(t *testing.T) {
	// 	deviceId := "12345"
	// 	c := restartHttpMock("GET", path.Join(deviceAuthBasePath, "devices", deviceId), `{
	// 		"id": "12345",
	// 		"identity_data": {
	// 		  "mac": "00:01:02:03:04:05",
	// 		  "sku": "My Device 1",
	// 		  "sn": "SN1234567890"
	// 		},
	// 		"status": "pending",
	// 		"created_ts": "2019-08-24T14:15:22Z",
	// 		"updated_ts": "2019-08-24T14:15:22Z",
	// 		"auth_sets": [
	// 		  {
	// 			"id": "string",
	// 			"pubkey": "string",
	// 			"identity_data": {
	// 			  "mac": "00:01:02:03:04:05",
	// 			  "sku": "My Device 1",
	// 			  "sn": "SN1234567890"
	// 			},
	// 			"status": "pending",
	// 			"ts": "2019-08-24T14:15:22Z"
	// 		  }
	// 		],
	// 		"decommissioning": true
	// 	  }`, 200)

	// 	d, e := c.GetDevice(deviceId)
	// 	if e != nil {
	// 		t.Error(e)
	// 	}

	// 	if d.ID != deviceId {
	// 		t.Errorf("Invalid id's")
	// 	}

	// 	// error response
	// 	c = restartHttpMock("GET", path.Join(deviceAuthBasePath, "devices", deviceId), `{}`, 404)
	// 	d, e = c.GetDevice(deviceId)
	// 	if e == nil {
	// 		t.Error(e)
	// 	}

	// 	// invalid json
	// 	c = restartHttpMock("GET", path.Join(deviceAuthBasePath, "devices", deviceId), `{`, 200)
	// 	d, e = c.GetDevice(deviceId)
	// 	if e == nil {
	// 		t.Error(e)
	// 	}
}

func TestDecomisionDevice(t *testing.T) {
	// 	// ok response
	// 	deviceId := "123456"
	// 	c := restartHttpMock("DELETE", path.Join(deviceAuthBasePath, "devices", deviceId), `{}`, 204)
	// 	e := c.DecomisionDevice(deviceId)
	// 	if e != nil {
	// 		t.Error(e)
	// 	}

	// 	// device not exists
	// 	c = restartHttpMock("DELETE", path.Join(deviceAuthBasePath, "devices", deviceId), `{}`, 404)
	// 	e = c.DecomisionDevice(deviceId)
	// 	if e == nil {
	// 		t.Error(e)
	// 	}

	// }

	// func TestRejectAuthtentication(t *testing.T) {
	// 	// ok response
	// 	deviceId := "123456"
	// 	authId := "654321"
	// 	c := restartHttpMock("DELETE", path.Join(deviceAuthBasePath, "devices", deviceId, "auth", authId), `{}`, 204)
	// 	e := c.RejectAuthtentication(deviceId, authId)
	// 	if e != nil {
	// 		t.Error(e)
	// 	}

	// 	// device not exists
	// 	c = restartHttpMock("DELETE", path.Join(deviceAuthBasePath, "devices", deviceId, "auth", authId), `{}`, 404)
	// 	e = c.RejectAuthtentication(deviceId, authId)
	// 	if e == nil {
	// 		t.Error(e)
	// 	}

}

func TestSetAuthtenticationStatus(t *testing.T) {
	// 	// ok response
	// 	deviceId := "123456"
	// 	authId := "654321"
	// 	c := restartHttpMock("PUT", path.Join(deviceAuthBasePath, "devices", deviceId, authId, "status"), `{}`, 204)
	// 	e := c.SetAuthtenticationStatus(deviceId, authId)
	// 	if e != nil {
	// 		t.Error(e)
	// 	}

	// 	// device not exists
	// 	c = restartHttpMock("PUT", path.Join(deviceAuthBasePath, "devices", deviceId, authId, "status"), `{}`, 404)
	// 	e = c.SetAuthtenticationStatus(deviceId, authId)
	// 	if e == nil {
	// 		t.Error(e)
	// 	}
}

func TestGetAuthtenticationStatus(t *testing.T) {
	// 	// ok response
	// 	deviceId := "123456"
	// 	authId := "654321"
	// 	c := restartHttpMock("GET", path.Join(deviceAuthBasePath, "devices", deviceId, authId, "status"), `{"status": "accepted"}`, 200)
	// 	s, e := c.GetAuthtenticationStatus(deviceId, authId)
	// 	if e != nil || s != "accepted" {
	// 		t.Error(e)
	// 	}

	// 	// device not exists
	// 	c = restartHttpMock("GET", path.Join(deviceAuthBasePath, "devices", deviceId, authId, "status"), `{ "error": "Ivalid status"}`, 404)
	// 	_, e = c.GetAuthtenticationStatus(deviceId, authId)
	// 	if e == nil {
	// 		t.Error(e)
	// 	}

	// 	// wrong json
	// 	c = restartHttpMock("GET", path.Join(deviceAuthBasePath, "devices", deviceId, authId, "status"), `{ "error": "Ivalid status`, 200)
	// 	_, e = c.GetAuthtenticationStatus(deviceId, authId)
	// 	if e == nil {
	// 		t.Error(e)
	// 	}
}

func TestCountDevices(t *testing.T) {
	// 	// ok response
	// 	c := restartHttpMock("GET", path.Join(deviceAuthBasePath, "devices/count"), `{"count": 42}`, 200)
	// 	s, e := c.CountDevices()
	// 	if e != nil || s != 42 {
	// 		t.Error(e)
	// 	}

	// 	// device not exists
	// 	c = restartHttpMock("GET", path.Join(deviceAuthBasePath, "devices/count"), `{ "error": "Ivalid status"}`, 400)
	// 	_, e = c.CountDevices()
	// 	if e == nil {
	// 		t.Error(e)
	// 	}

	// 	// wrong json
	// 	c = restartHttpMock("GET", path.Join(deviceAuthBasePath, "devices/count"), `{"count": 4`, 200)
	// 	_, e = c.CountDevices()
	// 	if e == nil {
	// 		t.Error(e)
	// 	}
}

func TestGetDeviceLimit(t *testing.T) {
	// // ok response
	// c := restartHttpMock("GET", path.Join(deviceAuthBasePath, "limits/max_devices"), `{"limit": 123}`, 200)
	// s, e := c.GetDeviceLimit()
	// if e != nil || s != 123 {
	// 	t.Error(e)
	// }

	// // device not exists
	// c = restartHttpMock("GET", path.Join(deviceAuthBasePath, "limits/max_devices"), "", 500)
	// _, e = c.GetDeviceLimit()
	// if e == nil {
	// 	t.Error(e)
	// }

	// // wrong json
	// c = restartHttpMock("GET", path.Join(deviceAuthBasePath, "limits/max_devices"), `{"linut": 4`, 200)
	// _, e = c.GetDeviceLimit()
	// if e == nil {
	// 	t.Error(e)
	// }
}
