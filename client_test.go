package menderrc_test

import (
	"context"

	"github.com/nandra/menderrc"
)

func ExampleNewClient() {

	// using default http client
	c := menderrc.NewClient()
	c.Login(context.Background(), "user", "pass")

}
