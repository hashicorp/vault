// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package approle

import (
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestApproleSecretId_Wrapped(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client
	client.SetToken(cluster.RootToken)

	err = client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/approle/role/test-role-1", map[string]interface{}{
		"name": "test-role-1",
	})
	require.NoError(t, err)

	client.SetWrappingLookupFunc(func(operation, path string) string {
		return "5m"
	})

	resp, err := client.Logical().Write("/auth/approle/role/test-role-1/secret-id", map[string]interface{}{})
	require.NoError(t, err)

	wrappedAccessor := resp.WrapInfo.WrappedAccessor
	wrappingToken := resp.WrapInfo.Token

	client.SetWrappingLookupFunc(func(operation, path string) string {
		return api.DefaultWrappingLookupFunc(operation, path)
	})

	unwrappedSecretid, err := client.Logical().Unwrap(wrappingToken)
	require.NoError(t, err)
	unwrappedAccessor := unwrappedSecretid.Data["secret_id_accessor"].(string)

	if wrappedAccessor != unwrappedAccessor {
		t.Fatalf("Expected wrappedAccessor (%v) to match wrapped secret_id_accessor (%v)", wrappedAccessor, unwrappedAccessor)
	}
}

func TestApproleSecretId_NotWrapped(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client
	client.SetToken(cluster.RootToken)

	err = client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/approle/role/test-role-1", map[string]interface{}{
		"name": "test-role-1",
	})
	require.NoError(t, err)

	resp, err := client.Logical().Write("/auth/approle/role/test-role-1/secret-id", map[string]interface{}{})
	require.NoError(t, err)

	if resp.WrapInfo != nil && resp.WrapInfo.WrappedAccessor != "" {
		t.Fatalf("WrappedAccessor unexpectedly set")
	}
}
