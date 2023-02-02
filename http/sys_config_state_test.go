package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/vault"
)

func TestSysConfigState_Sanitized(t *testing.T) {
	var resp *http.Response

	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp = testHttpGet(t, token, addr+"/v1/sys/config/state/sanitized")
	testResponseStatus(t, resp, 200)

	var actual map[string]interface{}
	var expected map[string]interface{}

	configResp := map[string]interface{}{
		"api_addr":                     "",
		"cache_size":                   json.Number("0"),
		"cluster_addr":                 "",
		"cluster_cipher_suites":        "",
		"cluster_name":                 "",
		"default_lease_ttl":            json.Number("0"),
		"default_max_request_duration": json.Number("0"),
		"disable_cache":                false,
		"disable_clustering":           false,
		"disable_indexing":             false,
		"disable_mlock":                false,
		"disable_performance_standby":  false,
		"disable_printable_check":      false,
		"disable_sealwrap":             false,
		"experiments":                  nil,
		"raw_storage_endpoint":         false,
		"detect_deadlocks":             "",
		"introspection_endpoint":       false,
		"disable_sentinel_trace":       false,
		"enable_ui":                    false,
		"ui": map[string]interface{}{
			"enabled": false,
			"dir":     "",
		},
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
}
