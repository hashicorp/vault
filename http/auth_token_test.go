package http

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/vault"
)

func TestAuthTokenCreate(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(token)

	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Lease: "1h",
	})
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth.LeaseDuration != 3600 {
		t.Errorf("expected 1h, got %q", secret.Auth.LeaseDuration)
	}

	renewCreateRequest := &api.TokenCreateRequest{
		TTL:       "1h",
		Renewable: new(bool),
	}

	secret, err = client.Auth().Token().Create(renewCreateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth.LeaseDuration != 3600 {
		t.Errorf("expected 1h, got %q", secret.Auth.LeaseDuration)
	}
	if secret.Auth.Renewable {
		t.Errorf("expected non-renewable token")
	}

	*renewCreateRequest.Renewable = true
	secret, err = client.Auth().Token().Create(renewCreateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth.LeaseDuration != 3600 {
		t.Errorf("expected 1h, got %q", secret.Auth.LeaseDuration)
	}
	if !secret.Auth.Renewable {
		t.Errorf("expected renewable token")
	}

	explicitMaxCreateRequest := &api.TokenCreateRequest{
		TTL:            "1h",
		ExplicitMaxTTL: "1800s",
	}

	secret, err = client.Auth().Token().Create(explicitMaxCreateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth.LeaseDuration != 1800 {
		t.Errorf("expected 1800 seconds, got %q", secret.Auth.LeaseDuration)
	}

	explicitMaxCreateRequest.ExplicitMaxTTL = "2h"
	secret, err = client.Auth().Token().Create(explicitMaxCreateRequest)
	if err != nil {
		t.Fatal(err)
	}
	if secret.Auth.LeaseDuration != 3600 {
		t.Errorf("expected 3600 seconds, got %q", secret.Auth.LeaseDuration)
	}
}

func TestAuthTokenLookup(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(token)

	// Create a new token ...
	secret2, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Lease: "1h",
	})
	if err != nil {
		t.Fatal(err)
	}

	// lookup details of this token
	secret, err := client.Auth().Token().Lookup(secret2.Auth.ClientToken)
	if err != nil {
		t.Fatalf("unable to lookup details of token, err = %v", err)
	}

	if secret.Data["id"] != secret2.Auth.ClientToken {
		t.Errorf("Did not get back details about our provided token, id returned=%s", secret.Data["id"])
	}

}

func TestAuthTokenLookupSelf(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(token)

	// you should be able to lookup your own token
	secret, err := client.Auth().Token().LookupSelf()
	if err != nil {
		t.Fatalf("should be allowed to lookup self, err = %v", err)
	}

	if secret.Data["id"] != token {
		t.Errorf("Did not get back details about our own (self) token, id returned=%s", secret.Data["id"])
	}
	if secret.Data["display_name"] != "root" {
		t.Errorf("Did not get back details about our own (self) token, display_name returned=%s", secret.Data["display_name"])
	}

}

func TestAuthTokenRenew(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(token)

	// The default root token is not renewable, so this should not work
	_, err = client.Auth().Token().Renew(token, 0)
	if err == nil {
		t.Fatal("should not be allowed to renew root token")
	}
	if !strings.Contains(err.Error(), "lease is not renewable") {
		t.Fatalf("wrong error; got %v", err)
	}

	// Create a new token that should be renewable
	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Lease: "1h",
	})
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(secret.Auth.ClientToken)

	// Now attempt a renew with the new token
	secret, err = client.Auth().Token().Renew(secret.Auth.ClientToken, 3600)
	if err != nil {
		t.Fatal(err)
	}

	if secret.Auth.LeaseDuration != 3600 {
		t.Errorf("expected 1h, got %v", secret.Auth.LeaseDuration)
	}

	if secret.Auth.Renewable != true {
		t.Error("expected lease to be renewable")
	}

	// Do the same thing with the self variant
	secret, err = client.Auth().Token().RenewSelf(3600)
	if err != nil {
		t.Fatal(err)
	}

	if secret.Auth.LeaseDuration != 3600 {
		t.Errorf("expected 1h, got %v", secret.Auth.LeaseDuration)
	}

	if secret.Auth.Renewable != true {
		t.Error("expected lease to be renewable")
	}
}
