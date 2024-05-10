// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

const (
	rootLeasePolicies = `
path "sys/internal/ui/*" {
capabilities = ["create", "read", "update", "delete", "list"]
}

path "auth/token/*" {
capabilities = ["create", "update", "read", "list"]
}

path "kv/foo*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
`

	dummy = `
path "/ns1/sys/leases/*" {
	capabilities = ["sudo", "create", "read", "update", "delete", "list"]
}

path "/ns1/auth/token/*" {
	capabilities = ["sudo", "create", "read", "update", "delete", "list"]
}
`
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
		t.Errorf("expected 1800 seconds, got %d", secret.Auth.LeaseDuration)
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
	if !strings.Contains(err.Error(), "invalid lease ID") {
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

// TestToken_InvalidTokenError checks that an InvalidToken error is only returned
// for tokens that have (1) exceeded the token TTL and (2) exceeded the number of uses
func TestToken_InvalidTokenError(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       logging.NewVaultLogger(hclog.Trace),
	}

	// Init new test cluster
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	// Add policy
	if err := client.Sys().PutPolicy("root-lease-policy", rootLeasePolicies); err != nil {
		t.Fatal(err)
	}
	// Add a dummy policy
	if err := client.Sys().PutPolicy("dummy", dummy); err != nil {
		t.Fatal(err)
	}

	rootToken := client.Token()

	// Enable kv secrets and mount initial secrets
	err := client.Sys().Mount("kv", &api.MountInput{Type: "kv"})
	require.NoError(t, err)

	writeSecretsToMount(t, client, "kv/foo", map[string]interface{}{
		"user":     "admin",
		"password": "password",
	})

	// Create a token that has a TTL of 5s
	tokenCreateRequest := &api.TokenCreateRequest{
		Policies: []string{"root-lease-policy"},
		TTL:      "5s",
	}
	secret, err := client.Auth().Token().CreateOrphan(tokenCreateRequest)
	token := secret.Auth.ClientToken
	client.SetToken(token)

	// Verify that token works to read from kv mount
	_, err = client.Logical().Read("kv/foo")
	require.NoError(t, err)

	time.Sleep(time.Second * 5)

	// Verify that token is expired and shows an "invalid token" error
	_, err = client.Logical().Read("kv/foo")
	require.ErrorContains(t, err, logical.ErrInvalidToken.Error())
	require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())

	// Create a second approle token with a token use limit
	client.SetToken(rootToken)
	tokenCreateRequest = &api.TokenCreateRequest{
		Policies: []string{"root-lease-policy"},
		NumUses:  5,
	}

	secret, err = client.Auth().Token().CreateOrphan(tokenCreateRequest)
	token = secret.Auth.ClientToken
	client.SetToken(token)

	for i := 0; i < 5; i++ {
		_, err = client.Logical().Read("kv/foo")
		require.NoError(t, err)
	}
	// Verify that the number of uses is exceeded so the "invalid token" error is displayed
	_, err = client.Logical().Read("kv/foo")
	require.ErrorContains(t, err, logical.ErrInvalidToken.Error())
	require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())

	// Create a third approle token that will have incorrect policy access to the subsequent request
	client.SetToken(rootToken)
	tokenCreateRequest = &api.TokenCreateRequest{
		Policies: []string{"dummy"},
	}

	secret, err = client.Auth().Token().CreateOrphan(tokenCreateRequest)
	token = secret.Auth.ClientToken
	client.SetToken(token)

	// Incorrect policy access should only return an ErrPermissionDenied error
	_, err = client.Logical().Read("kv/foo")
	require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())
	require.NotContains(t, err.Error(), logical.ErrInvalidToken)
}

func writeSecretsToMount(t *testing.T, client *api.Client, mountPath string, data map[string]interface{}) {
	_, err := client.Logical().Write(mountPath, data)
	require.NoError(t, err)
}
