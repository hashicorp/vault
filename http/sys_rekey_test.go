// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/vault"
)

// Test to check if the API errors out when wrong number of PGP keys are
// supplied for rekey
func TestSysRekey_Init_pgpKeysEntriesForRekey(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc:             Handler,
		RequestResponseCallback: schema.ResponseValidatingCallback(t),
	})
	cluster.Start()
	defer cluster.Cleanup()
	cl := cluster.Cores[0].Client

	_, err := cl.Logical().Write("sys/rekey/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
		"pgp_keys":         []string{"pgpkey1"},
	})
	if err == nil {
		t.Fatal("should have failed to write pgp key entry due to mismatched keys", err)
	}
}

func TestSysRekey_Init_Status(t *testing.T) {
	t.Run("status-barrier-default", func(t *testing.T) {
		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc:             Handler,
			RequestResponseCallback: schema.ResponseValidatingCallback(t),
		})
		cluster.Start()
		defer cluster.Cleanup()
		cl := cluster.Cores[0].Client

		resp, err := cl.Logical().Read("sys/rekey/init")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := resp.Data
		expected := map[string]interface{}{
			"started":               false,
			"t":                     json.Number("0"),
			"n":                     json.Number("0"),
			"progress":              json.Number("0"),
			"required":              json.Number("3"),
			"pgp_fingerprints":      interface{}(nil),
			"backup":                false,
			"nonce":                 "",
			"verification_required": false,
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
		}
	})
}

func TestSysRekey_Init_Setup(t *testing.T) {
	t.Run("init-barrier-barrier-key", func(t *testing.T) {
		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc:             Handler,
			RequestResponseCallback: schema.ResponseValidatingCallback(t),
		})
		cluster.Start()
		defer cluster.Cleanup()
		cl := cluster.Cores[0].Client

		// Start rekey
		resp, err := cl.Logical().Write("sys/rekey/init", map[string]interface{}{
			"secret_shares":    5,
			"secret_threshold": 3,
		})
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := resp.Data
		expected := map[string]interface{}{
			"started":               true,
			"t":                     json.Number("3"),
			"n":                     json.Number("5"),
			"progress":              json.Number("0"),
			"required":              json.Number("3"),
			"pgp_fingerprints":      interface{}(nil),
			"backup":                false,
			"verification_required": false,
		}

		if actual["nonce"].(string) == "" {
			t.Fatalf("nonce was empty")
		}
		expected["nonce"] = actual["nonce"]
		if diff := deep.Equal(actual, expected); diff != nil {
			t.Fatal(diff)
		}

		// Get rekey status
		resp, err = cl.Logical().Read("sys/rekey/init")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual = resp.Data
		expected = map[string]interface{}{
			"started":               true,
			"t":                     json.Number("3"),
			"n":                     json.Number("5"),
			"progress":              json.Number("0"),
			"required":              json.Number("3"),
			"pgp_fingerprints":      interface{}(nil),
			"backup":                false,
			"verification_required": false,
		}
		if actual["nonce"].(string) == "" {
			t.Fatalf("nonce was empty")
		}
		if actual["nonce"].(string) == "" {
			t.Fatalf("nonce was empty")
		}
		expected["nonce"] = actual["nonce"]
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
		}
	})
}

func TestSysRekey_Init_Cancel(t *testing.T) {
	t.Run("cancel-barrier-barrier-key", func(t *testing.T) {
		cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
			HandlerFunc:             Handler,
			RequestResponseCallback: schema.ResponseValidatingCallback(t),
		})
		cluster.Start()
		defer cluster.Cleanup()
		cl := cluster.Cores[0].Client

		initResp, err := cl.Logical().Write("sys/rekey/init", map[string]interface{}{
			"secret_shares":    5,
			"secret_threshold": 3,
		})
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		err = cl.Sys().RekeyCancelWithNonce(initResp.Data["nonce"].(string))
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		resp, err := cl.Logical().Read("sys/rekey/init")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := resp.Data
		expected := map[string]interface{}{
			"started":               false,
			"t":                     json.Number("0"),
			"n":                     json.Number("0"),
			"progress":              json.Number("0"),
			"required":              json.Number("3"),
			"pgp_fingerprints":      interface{}(nil),
			"backup":                false,
			"nonce":                 "",
			"verification_required": false,
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("\nexpected: %#v\nactual: %#v", expected, actual)
		}
	})
}

func TestSysRekey_badKey(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/rekey/update", map[string]interface{}{
		"key": "0123",
	})
	testResponseStatus(t, resp, 400)
}

func TestSysRekey_Update(t *testing.T) {
	t.Run("rekey-barrier-barrier-key", func(t *testing.T) {
		core, keys, token := vault.TestCoreUnsealed(t)
		ln, addr := TestServer(t, core)
		defer ln.Close()
		TestServerAuth(t, addr, token)

		resp := testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
			"secret_shares":    5,
			"secret_threshold": 3,
		})
		var rekeyStatus map[string]interface{}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &rekeyStatus)

		var actual map[string]interface{}
		var expected map[string]interface{}

		for i, key := range keys {
			resp = testHttpPut(t, token, addr+"/v1/sys/rekey/update", map[string]interface{}{
				"nonce": rekeyStatus["nonce"].(string),
				"key":   hex.EncodeToString(key),
			})

			actual = map[string]interface{}{}
			expected = map[string]interface{}{
				"started":               true,
				"nonce":                 rekeyStatus["nonce"].(string),
				"backup":                false,
				"pgp_fingerprints":      interface{}(nil),
				"required":              json.Number("3"),
				"t":                     json.Number("3"),
				"n":                     json.Number("5"),
				"progress":              json.Number(fmt.Sprintf("%d", i+1)),
				"verification_required": false,
			}
			testResponseStatus(t, resp, 200)
			testResponseBody(t, resp, &actual)

			if i+1 == len(keys) {
				delete(expected, "started")
				delete(expected, "required")
				delete(expected, "t")
				delete(expected, "n")
				delete(expected, "progress")
				expected["complete"] = true
				expected["keys"] = actual["keys"]
				expected["keys_base64"] = actual["keys_base64"]
			}

			if i+1 < len(keys) && (actual["nonce"] == nil || actual["nonce"].(string) == "") {
				t.Fatalf("expected a nonce, i is %d, actual is %#v", i, actual)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Fatalf("\nexpected: \n%#v\nactual: \n%#v", expected, actual)
			}
		}

		retKeys := actual["keys"].([]interface{})
		if len(retKeys) != 5 {
			t.Fatalf("bad: %#v", retKeys)
		}
		keysB64 := actual["keys_base64"].([]interface{})
		if len(keysB64) != 5 {
			t.Fatalf("bad: %#v", keysB64)
		}
	})
}

func TestSysRekey_ReInitUpdate(t *testing.T) {
	core, keys, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	})
	testResponseStatus(t, resp, 200)
	var initResp map[string]interface{}
	testResponseBody(t, resp, &initResp)

	resp = testHttpDeleteData(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"nonce": initResp["nonce"].(string),
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPut(t, token, addr+"/v1/sys/rekey/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	})
	testResponseStatus(t, resp, 200)

	resp = testHttpPut(t, token, addr+"/v1/sys/rekey/update", map[string]interface{}{
		"key": hex.EncodeToString(keys[0]),
	})

	testResponseStatus(t, resp, 400)
}
