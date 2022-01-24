package http

import (
	"testing"

	"github.com/hashicorp/vault/sdk/version"
	"github.com/hashicorp/vault/vault"
)

// TestSysVersionHistory_Get tests the sys/version-history endpoint. The
// endpoint is accessible in both an authenticated and unauthenticated fashion.
// Without synthetically altering the underlying core/versions storage entries,
// a single version entry should exist.
func TestSysVersionHistory_Get(t *testing.T) {
	cases := []struct{
		name   string
		auth   bool
	}{
		{
			name: "authenticated",
			auth: true,
		},
		{
			name: "unauthenticated",
			auth: false,
		},
		{
			name: "chicken",
			auth: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc

			t.Parallel()

			core, _, token := vault.TestCoreUnsealed(t)
			listener, addr := TestServer(t, core)
			defer listener.Close()

			if !tc.auth {
				token = ""
			}

			resp := testHttpGet(t, token, addr + "/v1/sys/version-history")

			var actual map[string]interface{}

			testResponseStatus(t, resp, 200)
			testResponseBody(t, resp, &actual)

			var respData map[string]interface{}
			var versions map[string]interface{}
			var ok bool

			if respData, ok = actual["data"].(map[string]interface{}); !ok {
				t.Fatalf("expected data key to be map, actual: %#v", actual["data"])
			}

			if versions, ok = respData["versions"].(map[string]interface{}); !ok {
				t.Fatalf("expected versions key to be map, actual: %#v", respData["versions"])
			}

			if len(versions) != 1 || versions[version.Version] == nil {
				t.Fatalf("expected single version history entry for %q", version.Version)
			}
		})
	}

}