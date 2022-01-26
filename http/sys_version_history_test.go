package http

import (
	"testing"

	"github.com/hashicorp/vault/sdk/version"
	"github.com/hashicorp/vault/vault"
)

// TestSysVersionHistory_List tests the sys/version-history endpoint. Requests
// to the endpoint must be authenticated. Without synthetically altering the
// underlying core/versions storage entries, a single version entry should
// exist.
func TestSysVersionHistory_List(t *testing.T) {
	cases := []struct {
		name           string
		auth           bool
		expectedStatus int
	}{
		{
			name:           "authenticated",
			auth:           true,
			expectedStatus: 200,
		},
		{
			name:           "unauthenticated",
			auth:           false,
			expectedStatus: 403,
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

			resp := testHttpList(t, token, addr+"/v1/sys/version-history")

			var actual map[string]interface{}

			testResponseStatus(t, resp, tc.expectedStatus)
			testResponseBody(t, resp, &actual)

			if tc.auth {
				var respData map[string]interface{}
				var ok bool
				var keys []interface{}
				var keyInfo map[string]interface{}

				if respData, ok = actual["data"].(map[string]interface{}); !ok {
					t.Fatalf("expected data key to be map, actual: %#v", actual["data"])
				}

				if keys, ok = respData["keys"].([]interface{}); !ok {
					t.Fatalf("expected keys to be array, actual: %#v", respData["keys"])
				}

				if keyInfo, ok = respData["key_info"].(map[string]interface{}); !ok {
					t.Fatalf("expected key_info to be map, actual: %#v", respData["key_info"])
				}

				if len(keys) != 1 {
					t.Fatalf("expected single version history entry for %q", version.Version)
				}

				if keyInfo[version.Version] == nil {
					t.Fatalf("expected version %s to be present in key_info, actual: %#v", version.Version, keyInfo)
				}
			}
		})
	}
}
