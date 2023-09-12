// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestKVV1Get(t *testing.T) {
	t.Parallel()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	client = cluster.Cores[0].Client

	// (the test cluster has already mounted the KVv1 backend at "secret")
	err := client.KVv1(v1MountPath).Put(context.Background(), secretPath, secretData)
	if err != nil {
		t.Fatal(err)
	}

	data, err := client.KVv1(v1MountPath).Get(context.Background(), secretPath)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "kvv1", data.Raw.MountType)
	require.Equal(t, secretData, data.Data)
}

func TestKVV2Get(t *testing.T) {
	t.Parallel()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.VersionedKVFactory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	client = cluster.Cores[0].Client

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

	require.Equal(t, "kvv2", data.Raw.MountType)

	data, err = client.KVv2(v2MountPath).Get(context.Background(), secretPath)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "kvv2", data.Raw.MountType)
	require.Equal(t, secretData, data.Data)
}
