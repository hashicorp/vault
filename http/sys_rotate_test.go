// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"encoding/json"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/vault"
)

func TestSysRotate(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/rotate", map[string]interface{}{})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/key-status")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"term": json.Number("2"),
		},
		"term": json.Number("2"),
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	for _, field := range []string{"install_time", "encryptions"} {
		actualVal, ok := actual["data"].(map[string]interface{})[field]
		if !ok || actualVal == "" {
			t.Fatal(field, " missing in data")
		}
		expected["data"].(map[string]interface{})[field] = actualVal
		expected[field] = actualVal
	}

	expected["request_id"] = actual["request_id"]
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}
