package http

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysConfigCors(t *testing.T) {
	var resp *http.Response

	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	corsConf := core.CORSConfig()

	// Try to enable CORS without providing a value for allowed_origins
	resp = testHttpPut(t, token, addr+"/v1/sys/config/cors", map[string]any{
		"allowed_headers": "X-Custom-Header",
	})

	testResponseStatus(t, resp, 500)

	// Enable CORS, but provide an origin this time.
	resp = testHttpPut(t, token, addr+"/v1/sys/config/cors", map[string]any{
		"allowed_origins": addr,
		"allowed_headers": "X-Custom-Header",
	})

	testResponseStatus(t, resp, 204)

	// Read the CORS configuration
	resp = testHttpGet(t, token, addr+"/v1/sys/config/cors")
	testResponseStatus(t, resp, 200)

	var actual map[string]any
	var expected map[string]any

	lenStdHeaders := len(corsConf.AllowedHeaders)

	expectedHeaders := make([]any, lenStdHeaders)

	for i := range corsConf.AllowedHeaders {
		expectedHeaders[i] = corsConf.AllowedHeaders[i]
	}

	expected = map[string]any{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]any{
			"enabled":         true,
			"allowed_origins": []any{addr},
			"allowed_headers": expectedHeaders,
		},
		"enabled":         true,
		"allowed_origins": []any{addr},
		"allowed_headers": expectedHeaders,
	}

	testResponseStatus(t, resp, 200)

	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected: %#v\nactual: %#v", expected, actual)
	}
}
