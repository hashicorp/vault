// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func createMultipleAlgorithmTestKeys(t *testing.T, b *backend, s logical.Storage, keyData map[string][]map[string]interface{}) {
	t.Helper()

	createKey := func(name string, data map[string]interface{}) {
		t.Helper()

		_, err := b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "keys/" + name,
			Data:      data,
			Storage:   s,
		})
		require.NoError(t, err, "creating key %q", name)
	}

	changeAlgorithm := func(name string, data map[string]interface{}) {
		t.Helper()
		_, err := b.HandleRequest(context.Background(), &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "keys/" + name + "/algorithm",
			Data:      data,
			Storage:   s,
		})
		require.NoError(t, err, "rotating key %q", name)
	}

	for keyName, requests := range keyData {
		for i, data := range requests {
			if i == 0 {
				createKey(keyName, data)
			} else {
				changeAlgorithm(keyName, data)
			}
		}
	}
}

func createMultipleAlgorithmTestKeysForCluster(t *testing.T, client *api.Client, keyData map[string][]map[string]interface{}) {
	t.Helper()

	createKey := func(name string, data map[string]interface{}) {
		t.Helper()

		_, err := client.Logical().Write("transit/keys/"+name, data)
		require.NoError(t, err, "creating key %q", name)
	}

	changeAlgorithm := func(name string, data map[string]interface{}) {
		t.Helper()

		_, err := client.Logical().Write(fmt.Sprintf("transit/keys/%s/algorithm", name), data)
		require.NoError(t, err, "rotating key %q", name)
	}

	for keyName, requests := range keyData {
		for i, data := range requests {
			if i == 0 {
				createKey(keyName, data)
			} else {
				changeAlgorithm(keyName, data)
			}
		}
	}
}
