// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/vault"
)

func TestSysConfigState_Sanitized(t *testing.T) {
	cases := []struct {
		name                    string
		storageConfig           *server.Storage
		haStorageConfig         *server.Storage
		expectedStorageOutput   map[string]interface{}
		expectedHAStorageOutput map[string]interface{}
	}{
		{
			name: "raft storage",
			storageConfig: &server.Storage{
				Type:              "raft",
				RedirectAddr:      "http://127.0.0.1:8200",
				ClusterAddr:       "http://127.0.0.1:8201",
				DisableClustering: false,
				Config: map[string]string{
					"path":           "/storage/path/raft",
					"node_id":        "raft1",
					"max_entry_size": "2097152",
				},
			},
			haStorageConfig: nil,
			expectedStorageOutput: map[string]interface{}{
				"type":               "raft",
				"redirect_addr":      "http://127.0.0.1:8200",
				"cluster_addr":       "http://127.0.0.1:8201",
				"disable_clustering": false,
				"raft": map[string]interface{}{
					"max_entry_size": "2097152",
				},
			},
			expectedHAStorageOutput: nil,
		},
		{
			name: "inmem storage, no HA storage",
			storageConfig: &server.Storage{
				Type:              "inmem",
				RedirectAddr:      "http://127.0.0.1:8200",
				ClusterAddr:       "http://127.0.0.1:8201",
				DisableClustering: false,
			},
			haStorageConfig: nil,
			expectedStorageOutput: map[string]interface{}{
				"type":               "inmem",
				"redirect_addr":      "http://127.0.0.1:8200",
				"cluster_addr":       "http://127.0.0.1:8201",
				"disable_clustering": false,
			},
			expectedHAStorageOutput: nil,
		},
		{
			name: "inmem storage, raft HA storage",
			storageConfig: &server.Storage{
				Type:              "inmem",
				RedirectAddr:      "http://127.0.0.1:8200",
				ClusterAddr:       "http://127.0.0.1:8201",
				DisableClustering: false,
			},
			haStorageConfig: &server.Storage{
				Type:              "raft",
				RedirectAddr:      "http://127.0.0.1:8200",
				ClusterAddr:       "http://127.0.0.1:8201",
				DisableClustering: false,
				Config: map[string]string{
					"path":           "/storage/path/raft",
					"node_id":        "raft1",
					"max_entry_size": "2097152",
				},
			},
			expectedStorageOutput: map[string]interface{}{
				"type":               "inmem",
				"redirect_addr":      "http://127.0.0.1:8200",
				"cluster_addr":       "http://127.0.0.1:8201",
				"disable_clustering": false,
			},
			expectedHAStorageOutput: map[string]interface{}{
				"type":               "raft",
				"redirect_addr":      "http://127.0.0.1:8200",
				"cluster_addr":       "http://127.0.0.1:8201",
				"disable_clustering": false,
				"raft": map[string]interface{}{
					"max_entry_size": "2097152",
				},
			},
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var resp *http.Response
			confRaw := &server.Config{
				Storage:   tc.storageConfig,
				HAStorage: tc.haStorageConfig,
				SharedConfig: &configutil.SharedConfig{
					Listeners: []*configutil.Listener{
						{
							Type:    "tcp",
							Address: "127.0.0.1",
						},
					},
				},
			}

			conf := &vault.CoreConfig{
				RawConfig: confRaw,
			}

			core, _, token := vault.TestCoreUnsealedWithConfig(t, conf)
			ln, addr := TestServer(t, core)
			defer ln.Close()
			TestServerAuth(t, addr, token)

			resp = testHttpGet(t, token, addr+"/v1/sys/config/state/sanitized")
			testResponseStatus(t, resp, 200)

			var actual map[string]interface{}
			var expected map[string]interface{}

			configResp := map[string]interface{}{
				"api_addr":                            "",
				"cache_size":                          json.Number("0"),
				"cluster_addr":                        "",
				"cluster_cipher_suites":               "",
				"cluster_name":                        "",
				"default_lease_ttl":                   json.Number("0"),
				"default_max_request_duration":        json.Number("0"),
				"disable_cache":                       false,
				"disable_clustering":                  false,
				"disable_indexing":                    false,
				"disable_mlock":                       false,
				"disable_performance_standby":         false,
				"disable_printable_check":             false,
				"disable_sealwrap":                    false,
				"experiments":                         nil,
				"raw_storage_endpoint":                false,
				"detect_deadlocks":                    "",
				"introspection_endpoint":              false,
				"disable_sentinel_trace":              false,
				"enable_ui":                           false,
				"log_format":                          "",
				"log_level":                           "",
				"max_lease_ttl":                       json.Number("0"),
				"pid_file":                            "",
				"plugin_directory":                    "",
				"plugin_file_uid":                     json.Number("0"),
				"plugin_file_permissions":             json.Number("0"),
				"enable_response_header_hostname":     false,
				"enable_response_header_raft_node_id": false,
				"log_requests_level":                  "",
				"listeners": []interface{}{
					map[string]interface{}{
						"config": nil,
						"type":   "tcp",
					},
				},
				"storage":                       tc.expectedStorageOutput,
				"administrative_namespace_path": "",
				"imprecise_lease_role_tracking": false,
			}

			if tc.expectedHAStorageOutput != nil {
				configResp["ha_storage"] = tc.expectedHAStorageOutput
			}

			expected = map[string]interface{}{
				"lease_id":       "",
				"renewable":      false,
				"lease_duration": json.Number("0"),
				"wrap_info":      nil,
				"warnings":       nil,
				"auth":           nil,
				"data":           configResp,
			}

			testResponseBody(t, resp, &actual)
			expected["request_id"] = actual["request_id"]

			if diff := deep.Equal(actual, expected); len(diff) > 0 {
				t.Fatalf("bad mismatch response body: diff: %v", diff)
			}
		})
	}
}
