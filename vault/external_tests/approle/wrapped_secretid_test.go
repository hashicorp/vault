// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package approle

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

func TestApproleSecretId_Wrapped(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
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
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	err := client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
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
