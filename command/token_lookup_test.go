package command

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestTokenLookupAccessor(t *testing.T) {
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

	// Enable the accessor flag
	args = append(args, "-accessor")

	// Expect failure if no argument is passed when accessor flag is set
	code := c.Run(args)
	if code == 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	// Add token accessor as arg
	args = append(args, resp.Auth.Accessor)
	code = c.Run(args)
	if code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestTokenLookupSelf(t *testing.T) {
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

	// Run it against itself
	code := c.Run(args)

	// Verify it worked
	if code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestTokenLookup(t *testing.T) {
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

	// Add token as arg for real test and run it
	args = append(args, resp.Auth.ClientToken)
	code := c.Run(args)

	// Verify it worked
	if code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
