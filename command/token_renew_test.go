package command

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestTokenRenew(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &TokenRenewCommand{
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

	// Create a token
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

	// Renew, passing in the token
	args = append(args, resp.Auth.ClientToken)
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestTokenRenewWithIncrement(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &TokenRenewCommand{
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

	// Create a token
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

	// Renew, passing in the token
	args = append(args, resp.Auth.ClientToken)
	args = append(args, "72h")
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestTokenRenewSelf(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &TokenRenewCommand{
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

	// Create a token
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
	if resp.Auth.ClientToken == "" {
		t.Fatal("returned client token is empty")
	}

	c.Meta.ClientToken = resp.Auth.ClientToken

	// Renew using the self endpoint
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestTokenRenewSelfWithIncrement(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &TokenRenewCommand{
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

	// Create a token
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
	if resp.Auth.ClientToken == "" {
		t.Fatal("returned client token is empty")
	}

	c.Meta.ClientToken = resp.Auth.ClientToken

	args = append(args, "-increment=72h")
	// Renew using the self endpoint
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
