// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"context"
	"testing"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/vault/helper/testhelpers"

	vaulthttp "github.com/hashicorp/vault/http"

	"github.com/hashicorp/vault/vault"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/api"
)

// TestKVV1Get tests an end-to-end KVV1 get, and checks the response
func TestKVV1Get(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// (the test cluster has already mounted the KVv1 backend at "secret")
	err := client.KVv1(v1MountPath).Put(context.Background(), secretPath, secretData)
	if err != nil {
		t.Fatal(err)
	}

	data, err := client.KVv1(v1MountPath).Get(context.Background(), secretPath)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "kv", data.Raw.MountType)
	require.Equal(t, secretData, data.Data)
}

// TestKVV2Get tests an end-to-end KVV2 get, and checks the response
func TestKVV2Get(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	// mount the KVv2 backend
	err := client.Sys().Mount(v2MountPath, &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	data, err := client.KVv2(v2MountPath).Put(context.Background(), secretPath, secretData)
	if err != nil {
		t.Fatal(err)
	}

	data, err = client.KVv2(v2MountPath).Get(context.Background(), secretPath)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "kv", data.Raw.MountType)
	require.Equal(t, secretData, data.Data)
}

// TestKVV2Get_RequestForwarding tests an end-to-end KVV2 get via request forwarding, and checks the response
func TestKVV2Get_RequestForwarding(t *testing.T) {
	t.Parallel()

	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		DisablePerformanceStandby: true,
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.VersionedKVFactory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client
	testhelpers.WaitForActiveNodeAndStandbys(t, cluster)
	standbys := testhelpers.DeriveStandbyCores(t, cluster)
	standby := standbys[0].Client

	// mount the KVv2 backend
	err := client.Sys().Mount(v2MountPath, &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	data, err := client.KVv2(v2MountPath).Put(context.Background(), secretPath, secretData)
	if err != nil {
		t.Fatal(err)
	}

	data, err = client.KVv2(v2MountPath).Get(context.Background(), secretPath)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "kv", data.Raw.MountType)
	require.Equal(t, secretData, data.Data)

	// Now do the same thing to the standby
	data, err = standby.KVv2(v2MountPath).Get(context.Background(), secretPath)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "kv", data.Raw.MountType)
	require.Equal(t, secretData, data.Data)
}
