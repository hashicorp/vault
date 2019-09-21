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
		"APIAddr":                   "",
		"CacheSize":                 json.Number("0"),
		"ClusterAddr":               "",
		"ClusterCipherSuites":       "",
		"ClusterName":               "",
		"DefaultLeaseTTL":           json.Number("0"),
		"DefaultMaxRequestDuration": json.Number("0"),
		"DisableCache":              false,
		"DisableClustering":         false,
		"DisableIndexing":           false,
		"DisableMlock":              false,
		"DisablePerformanceStandby": false,
		"DisablePrintableCheck":     false,
		"DisableSealWrap":           false,
		"EnableRawEndpoint":         false,
		"EnableUI":                  false,
		"LogFormat":                 "",
		"LogLevel":                  "",
		"MaxLeaseTTL":               json.Number("0"),
		"PidFile":                   "",
		"PluginDirectory":           "",
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
