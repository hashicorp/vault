package command

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/meta"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/cli"
)

func TestRenew(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	ui := new(cli.MockUi)
	c := &RenewCommand{
		Meta: meta.Meta{
			ClientToken: token,
			Ui:          ui,
		},
	}

	// write a secret with a lease
	client := testClient(t, addr, token)
	_, err := client.Logical().Write("secret/foo", map[string]interface{}{
		"key":   "value",
		"lease": "1m",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// read the secret to get its lease ID
	secret, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	args := []string{
		"-address", addr,
		secret.LeaseID,
	}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestRenewBothWays(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := http.TestServer(t, core)
	defer ln.Close()

	// write a secret with a lease
	client := testClient(t, addr, token)
	_, err := client.Logical().Write("secret/foo", map[string]interface{}{
		"key": "value",
		"ttl": "1m",
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// read the secret to get its lease ID
	secret, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Test one renew path
	r := client.NewRequest("PUT", "/v1/sys/renew")
	body := map[string]interface{}{
		"lease_id": secret.LeaseID,
	}
	if err := r.SetJSONBody(body); err != nil {
		t.Fatal(err)
	}
	resp, err := client.RawRequest(r)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	secret, err = api.ParseSecret(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if secret.LeaseDuration != 60 {
		t.Fatal("bad lease duration")
	}

	// Test another
	r = client.NewRequest("PUT", "/v1/sys/leases/renew")
	body = map[string]interface{}{
		"lease_id": secret.LeaseID,
	}
	if err := r.SetJSONBody(body); err != nil {
		t.Fatal(err)
	}
	resp, err = client.RawRequest(r)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	secret, err = api.ParseSecret(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if secret.LeaseDuration != 60 {
		t.Fatal("bad lease duration")
	}

	// Test the other
	r = client.NewRequest("PUT", "/v1/sys/renew/"+secret.LeaseID)
	resp, err = client.RawRequest(r)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	secret, err = api.ParseSecret(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if secret.LeaseDuration != 60 {
		t.Fatalf("bad lease duration; secret is %#v\n", *secret)
	}

	// Test another
	r = client.NewRequest("PUT", "/v1/sys/leases/renew/"+secret.LeaseID)
	resp, err = client.RawRequest(r)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	secret, err = api.ParseSecret(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if secret.LeaseDuration != 60 {
		t.Fatalf("bad lease duration; secret is %#v\n", *secret)
	}
}
