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

// Benchmark_KVV2 runs some basic benchmark on KVV2
func Benchmark_KVV2(b *testing.B) {
	cluster := vault.NewTestCluster(b, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	client := cluster.Cores[0].Client

	// mount the KVv2 backend
	err := client.Sys().Mount(v2MountPath, &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		b.Fatal(err)
	}

	data, err := client.KVv2(v2MountPath).Put(context.Background(), secretPath, secretData)
	if err != nil {
		b.Fatal(err)
	}

	data, err = client.KVv2(v2MountPath).Get(context.Background(), secretPath)
	if err != nil {
		b.Fatal(err)
	}

	require.Equal(b, "kv", data.Raw.MountType)
	require.Equal(b, secretData, data.Data)

	bench := func(b *testing.B, dataSize int) {
		data, err := uuid.GenerateRandomBytes(dataSize)
		if err != nil {
			b.Fatal(err)
		}

		testName := b.Name()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("%s/%x", secretPath, md5.Sum([]byte(fmt.Sprintf("%s-%d", testName, i))))
			_, err := client.KVv2(v2MountPath).Put(context.Background(), key, map[string]interface{}{
				"foo": string(data),
			})
			if err != nil {
				b.Fatal(err)
			}
			_, err = client.KVv2(v2MountPath).Get(context.Background(), secretPath)
			if err != nil {
				b.Fatal(err)
			}
		}
	}

	b.Run("kv-puts-and-gets", func(b *testing.B) { bench(b, 1024) })
}
