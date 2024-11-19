// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package bench

import (
	"context"
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

const (
	v2MountPath = "secret-v2"
	secretPath  = "my-secret"
)

var secretData = map[string]interface{}{
	"foo": "bar",
}

// BenchmarkSmoke_KVV2 runs basic benchmarks on writes and reads to KVV2 on an inmem test cluster.
func BenchmarkSmoke_KVV2(b *testing.B) {
	cluster := vault.NewTestCluster(b, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	// mount the KVv2 backend
	err := client.Sys().Mount(v2MountPath, &api.MountInput{
		Type: "kv-v2",
	})
	require.NoError(b, err)

	data, err := client.KVv2(v2MountPath).Put(context.Background(), secretPath, secretData)
	require.NoError(b, err)

	data, err = client.KVv2(v2MountPath).Get(context.Background(), secretPath)
	require.NoError(b, err)

	require.Equal(b, "kv", data.Raw.MountType)
	require.Equal(b, secretData, data.Data)

	bench := func(b *testing.B, dataSize int) {
		data, err := uuid.GenerateRandomBytes(dataSize)
		require.NoError(b, err)

		testName := b.Name()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("%s/%x", secretPath, md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
			_, err := client.KVv2(v2MountPath).Put(context.Background(), key, map[string]interface{}{
				"foo": string(data),
			})
			require.NoError(b, err)
			_, err = client.KVv2(v2MountPath).Get(context.Background(), secretPath)
			require.NoError(b, err)
		}
	}

	b.Run("kv-puts-and-gets", func(b *testing.B) { bench(b, 1024) })
}

// BenchmarkSmoke_ClusterCreation benchmarks the creation, start, and a cleanup of a vault.TestCluster.
// Note that the cluster created here uses inmem Physical and HAPhysical backends.
func BenchmarkSmoke_ClusterCreation(b *testing.B) {
	bench := func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cluster := vault.NewTestCluster(b, &vault.CoreConfig{}, &vault.TestClusterOptions{
				HandlerFunc: vaulthttp.Handler,
			})
			cluster.Cleanup()
		}
	}

	b.Run("cluster-creation", func(b *testing.B) { bench(b) })
}

// BenchmarkSmoke_MountUnmount runs some basic benchmarking on mounting and unmounting
func BenchmarkSmoke_MountUnmount(b *testing.B) {
	cluster := vault.NewTestCluster(b, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	bench := func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			err := client.Sys().Mount(v2MountPath, &api.MountInput{
				Type: "kv-v2",
			})
			require.NoError(b, err)
			err = client.Sys().Unmount(v2MountPath)
			require.NoError(b, err)
		}
	}

	b.Run("mount-unmount", func(b *testing.B) { bench(b) })
}

// BenchmarkSmoke_TokenCreationRevocation runs some basic benchmarking on tokens
func BenchmarkSmoke_TokenCreationRevocation(b *testing.B) {
	cluster := vault.NewTestCluster(b, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client
	rootToken := client.Token()

	bench := func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			client.SetToken(rootToken)
			secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
				Policies: []string{"default"},
				TTL:      "30m",
			})
			require.NoError(b, err)
			require.NotNil(b, secret)
			require.NotNil(b, secret.Auth)
			require.NotNil(b, secret.Auth.ClientToken)
			client.SetToken(secret.Auth.ClientToken)
			err = client.Auth().Token().RevokeSelf(secret.Auth.ClientToken)
			require.NoError(b, err)
		}
	}

	b.Run("token-creation-revocation", func(b *testing.B) { bench(b) })
}
