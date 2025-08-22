// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

// testVaultServer creates a test vault cluster and returns a configured API
// client and closer function.
func testVaultServer(t testing.TB) (*api.Client, func()) {
	t.Helper()

	client, _, closer := testVaultServerUnseal(t)
	return client, closer
}

// testVaultServerUnseal creates a test vault cluster and returns a configured
// API client, list of unseal keys (as strings), and a closer function.
func testVaultServerUnseal(t testing.TB) (*api.Client, []string, func()) {
	t.Helper()

	return testVaultServerCoreConfig(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": audit.NewFileBackend,
		},
		LogicalBackends: map[string]logical.Factory{
			"database":       database.Factory,
			"generic-leased": vault.LeasedPassthroughBackendFactory,
			"pki":            pki.Factory,
			"transit":        transit.Factory,
		},
		BuiltinRegistry: builtinplugins.Registry,
	})
}

// testVaultServerCoreConfig creates a new vault cluster with the given core
// configuration. This is a lower-level test helper.
func testVaultServerCoreConfig(t testing.TB, coreConfig *vault.CoreConfig) (*api.Client, []string, func()) {
	t.Helper()

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: http.Handler,
		NumCores:    1,
	})
	cluster.Start()

	// Make it easy to get access to the active
	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)

	// Get the client already setup for us!
	client := cluster.Cores[0].Client
	client.SetToken(cluster.RootToken)

	// Convert the unseal keys to base64 encoded, since these are how the user
	// will get them.
	unsealKeys := make([]string, len(cluster.BarrierKeys))
	for i := range unsealKeys {
		unsealKeys[i] = base64.StdEncoding.EncodeToString(cluster.BarrierKeys[i])
	}

	return client, unsealKeys, func() { defer cluster.Cleanup() }
}

func TestTransit_Kyber_EncryptDecrypt_RoundTrip(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		typ  string
	}{
		{"kyber512", "kyber512"},
		{"kyber768", "kyber768"},
		{"kyber1024", "kyber1024"},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			client, _, closer := testVaultServerUnseal(t)
			defer closer()

			if err := client.Sys().Mount("transit", &api.MountInput{Type: "transit"}); err != nil {
				t.Fatalf("mount transit: %v", err)
			}

			keyName := "test-" + c.name
			if _, err := client.Logical().Write("transit/keys/"+keyName, map[string]any{
				"type": c.typ,
			}); err != nil {
				t.Fatalf("create kyber key: %v", err)
			}

			msg := "hello " + c.name
			enc, err := client.Logical().Write("transit/encrypt/"+keyName, map[string]any{
				"plaintext": base64.StdEncoding.EncodeToString([]byte(msg)),
			})
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			ct, _ := enc.Data["ciphertext"].(string)
			if ct == "" {
				t.Fatalf("encrypt returned empty ciphertext")
			}

			dec, err := client.Logical().Write("transit/decrypt/"+keyName, map[string]any{
				"ciphertext": ct,
			})
			if err != nil {
				t.Fatalf("decrypt: %v", err)
			}

			gotB64, _ := dec.Data["plaintext"].(string)
			got, err := base64.StdEncoding.DecodeString(gotB64)
			if err != nil {
				t.Fatalf("decode plaintext: %v", err)
			}

			if string(got) != msg {
				t.Fatalf("round-trip mismatch: got %q, want %q", got, msg)
			}
		})
	}
}
