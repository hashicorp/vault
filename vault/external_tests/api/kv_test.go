// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/vault/api"
)

func TestKVV1Get(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)

	cluster.Start()
	defer cluster.Cleanup()

	client = cluster.Cores[0].Client

	// (the test cluster has already mounted the KVv1 backend at "secret")
	err := client.KVv1(v1MountPath).Put(context.Background(), secretPath, secretData)
	if err != nil {
		t.Fatal(err)
	}

	//// We use raw requests so we can check the headers for cache hit/miss.
	//req := client.NewRequest(http.MethodGet, "/v1/secret/my-secret")
	//resp1, err := client.RawRequest(req)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//resp1Map := map[string]interface{}{}
	//body, err := io.ReadAll(resp1.Body)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//err = json.Unmarshal(body, &resp1Map)
	//if err != nil {
	//	t.Fatal(err)
	//}

	data, err := client.KVv1(v1MountPath).Get(context.Background(), secretPath)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "kv", data.Raw.MountType)
	require.Equal(t, secretData, data.Data)
}

func TestKVV2Get(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, nil)

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

	data, err = client.KVv2(v2MountPath).Get(context.Background(), secretPath)
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "kv", data.Raw.MountType)
	require.Equal(t, secretData, data.Data)
}
