package command

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func testClient(t *testing.T, addr string, token string) *api.Client {
	config := api.DefaultConfig()
	config.Address = addr
	client, err := api.NewClient(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	client.SetToken(token)

	return client
}
