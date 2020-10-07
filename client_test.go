package menderrc_test

import "net/http"

func ExampleNewClient() {

	// using default http client
	c := menderrc.NewClient()
	c.Login("user", "pass")

	// using custom http client and server url
	customHttpClient := http.NewClient()

	c := menderrc.NewClient(
		menderrc.SetHttpClient(customHttpClient),
		menderrc.SetBaseApiUrl("https://mymender.example.io")
	)

}
