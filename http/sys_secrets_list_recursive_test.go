// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestSysSecretsListRecursive(t *testing.T) {
	core, _, rootToken := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	testMountPath1 := "test-kv-12345"
	testMountPath2 := "test-db-12346"

	// Mount two KV engines using the HTTP API
	mountKV(t, addr, rootToken, testMountPath1)
	mountKV(t, addr, rootToken, testMountPath2)

	// Write some secrets to each mount
	writeSecretHTTP(t, addr, rootToken, fmt.Sprintf("%s/alpha/beta", testMountPath1), map[string]interface{}{"key": "value1"})
	writeSecretHTTP(t, addr, rootToken, fmt.Sprintf("%s/config/db", testMountPath1), map[string]interface{}{"host": "localhost"})
	writeSecretHTTP(t, addr, rootToken, fmt.Sprintf("%s/conn", testMountPath2), map[string]interface{}{"dsn": "postgres://..."})

	// Test 1: List all secrets recursively (no filters)
	t.Run("list all recursive", func(t *testing.T) {
		resp := testHttpGet(t, rootToken, addr+"/v1/sys/secrets/list-recursive")
		testResponseStatus(t, resp, http.StatusOK)

		var secret api.Secret
		testResponseBody(t, resp, &secret)
		require.NotNil(t, secret.Data["keys"], "expected keys in response")

		keys := secret.Data["keys"].([]interface{})
		// Should have at least the 3 secrets we wrote + identity mount keys
		require.GreaterOrEqual(t, len(keys), 3, "expected at least 3 recursive keys, got %d", len(keys))

		// Verify our secrets are in the list
		found := make(map[string]bool)
		for _, k := range keys {
			key := k.(string)
			if key == fmt.Sprintf("%s/alpha/beta", testMountPath1) ||
				key == fmt.Sprintf("%s/config/db", testMountPath1) ||
				key == fmt.Sprintf("%s/conn", testMountPath2) {
				found[key] = true
			}
		}
		require.True(t, found[fmt.Sprintf("%s/alpha/beta", testMountPath1)], "expected alpha/beta in recursive list")
		require.True(t, found[fmt.Sprintf("%s/config/db", testMountPath1)], "expected config/db in recursive list")
		require.True(t, found[fmt.Sprintf("%s/conn", testMountPath2)], "expected conn in recursive list")
	})

	// Test 2: List with path filter (only one mount)
	t.Run("list filtered by path", func(t *testing.T) {
		queryParams := fmt.Sprintf("?path=%s/", testMountPath1)
		resp := testHttpGet(t, rootToken, addr+"/v1/sys/secrets/list-recursive"+queryParams)
		testResponseStatus(t, resp, http.StatusOK)

		var secret api.Secret
		testResponseBody(t, resp, &secret)
		require.NotNil(t, secret.Data["keys"])

		keys := secret.Data["keys"].([]interface{})
		for _, k := range keys {
			key := k.(string)
			require.Contains(t, key, testMountPath1, "expected all keys to be under %s, got %s", testMountPath1, key)
		}
	})

	// Test 3: List with glob pattern filter
	t.Run("list with pattern filter", func(t *testing.T) {
		queryParams := "?pattern=*conn*"
		resp := testHttpGet(t, rootToken, addr+"/v1/sys/secrets/list-recursive"+queryParams)
		testResponseStatus(t, resp, http.StatusOK)

		var secret api.Secret
		testResponseBody(t, resp, &secret)
		require.NotNil(t, secret.Data["keys"])

		keys := secret.Data["keys"].([]interface{})
		for _, k := range keys {
			key := k.(string)
			require.Contains(t, key, "conn", "expected all keys to match *conn*, got %s", key)
		}
	})

	// Test 4: List with fuzzy filter (case-insensitive substring match)
	t.Run("list with fuzzy filter", func(t *testing.T) {
		queryParams := "?fuzzy=true&pattern=ALPHA"
		resp := testHttpGet(t, rootToken, addr+"/v1/sys/secrets/list-recursive"+queryParams)
		testResponseStatus(t, resp, http.StatusOK)

		var secret api.Secret
		testResponseBody(t, resp, &secret)
		require.NotNil(t, secret.Data["keys"])

		keys := secret.Data["keys"].([]interface{})
		for _, k := range keys {
			key := k.(string)
			require.Contains(t, key, "alpha", "expected fuzzy case-insensitive match for ALPHA")
		}
	})

	// Test 5: Verify response includes key_info with mount metadata
	t.Run("response contains key info", func(t *testing.T) {
		resp := testHttpGet(t, rootToken, addr+"/v1/sys/secrets/list-recursive")
		testResponseStatus(t, resp, http.StatusOK)

		var secret api.Secret
		testResponseBody(t, resp, &secret)
		require.NotNil(t, secret.Data["key_info"], "expected key_info in response")

		keyInfo := secret.Data["key_info"].(map[string]interface{})
		require.NotEmpty(t, keyInfo, "key_info should not be empty")

		for _, info := range keyInfo {
			infoMap := info.(map[string]interface{})
			require.Contains(t, infoMap, "mount", "expected mount in key info")
			require.Contains(t, infoMap, "type", "expected type in key info")
			require.Contains(t, infoMap, "accessor", "expected accessor in key info")
		}
	})

}

// mountKV mounts a KV v2 secret engine via the HTTP API.
func mountKV(t *testing.T, addr, token, path string) {
	t.Helper()
	body := map[string]interface{}{
		"type": "kv",
	}
	resp := testHttpPutDisableRedirect(t, token, addr+"/v1/sys/mounts/"+path, body)
	testResponseStatusOr(t, resp, http.StatusNoContent)
}

// writeSecretHTTP writes a secret via the HTTP API.
func writeSecretHTTP(t *testing.T, addr, token, path string, data map[string]interface{}) {
	t.Helper()
	body := map[string]interface{}{
		"data": data,
	}
	resp := testHttpPut(t, token, addr+"/v1/"+path, body)
	testResponseStatusOr(t, resp, http.StatusNoContent)
}

// testResponseStatusOr checks that the response has one of the expected status codes.
func testResponseStatusOr(t *testing.T, resp *http.Response, expectedStatus int) {
	t.Helper()
	if resp.StatusCode != expectedStatus {
		t.Fatalf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}
}

