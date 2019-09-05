package http

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/command/server"
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

	expectedConfig := new(server.Config)
	configResp := structs.New(expectedConfig.Sanitized()).Map()

	var nilObject interface{}
	// Do some surgery on the expected config to line up the
	// types and string the raw fields.
	for k, v := range configResp {
		if strings.HasSuffix(k, "Raw") {
			delete(configResp, k)
			continue
		}
		switch v.(type) {
		case int:
			configResp[k] = json.Number(strconv.Itoa(v.(int)))
		case time.Duration:
			configResp[k] = json.Number(strconv.Itoa(int(v.(time.Duration))))
		}
	}
	configResp["HAStorage"] = nilObject
	configResp["Storage"] = nilObject
	configResp["Telemetry"] = nilObject

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

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad mismatch response body:\nexpected:\n%#v\nactual:\n%#v", expected, actual)
	}

}
