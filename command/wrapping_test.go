package command

import (
	"os"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestWrapping_Env(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &TokenLookupCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
	}
	// Run it once for client
	c.Run(args)

	// Create a new token for us to use
	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Lease: "1h",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	prevWrapTTLEnv := os.Getenv(api.EnvVaultWrapTTL)
	os.Setenv(api.EnvVaultWrapTTL, "5s")
	defer func() {
		os.Setenv(api.EnvVaultWrapTTL, prevWrapTTLEnv)
	}()

	// Now when we do a lookup-self the response should be wrapped
	args = append(args, resp.Auth.ClientToken)

	resp, err = client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.WrapInfo == nil {
		t.Fatal("nil wrap info")
	}
	if resp.WrapInfo.Token == "" || resp.WrapInfo.TTL != 5 {
		t.Fatal("did not get token or ttl wrong")
	}
}

func TestWrapping_Flag(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &TokenLookupCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	args := []string{
		"-address", addr,
		"-wrap-ttl", "5s",
	}
	// Run it once for client
	c.Run(args)

	// Create a new token for us to use
	client, err := c.Client()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Lease: "1h",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.WrapInfo == nil {
		t.Fatal("nil wrap info")
	}
	if resp.WrapInfo.Token == "" || resp.WrapInfo.TTL != 5 {
		t.Fatal("did not get token or ttl wrong")
	}
}
