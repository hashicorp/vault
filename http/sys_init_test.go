// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"encoding/hex"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
)

func TestSysInit_get(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	{
		// Pre-init
		resp, err := http.Get(addr + "/v1/sys/init")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"initialized": false,
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("bad: %#v", actual)
		}
	}

	vault.TestCoreInit(t, core)

	{
		// Post-init
		resp, err := http.Get(addr + "/v1/sys/init")
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"initialized": true,
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("bad: %#v", actual)
		}
	}
}

// Test to check if the API errors out when wrong number of PGP keys are
// supplied
func TestSysInit_pgpKeysEntries(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
		"pgp_keys":         []string{"pgpkey1"},
	})
	testResponseStatus(t, resp, 400)
}

// Test to check if the API errors out when wrong number of PGP keys are
// supplied for recovery config
func TestSysInit_pgpKeysEntriesForRecovery(t *testing.T) {
	core := vault.TestCoreNewSeal(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":      1,
		"secret_threshold":   1,
		"stored_shares":      1,
		"recovery_shares":    5,
		"recovery_threshold": 3,
		"recovery_pgp_keys":  []string{"pgpkey1"},
	})
	testResponseStatus(t, resp, 400)
}

func TestSysInit_put(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": 3,
	})

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	keysRaw, ok := actual["keys"]
	if !ok {
		t.Fatalf("no keys: %#v", actual)
	}

	if _, ok := actual["root_token"]; !ok {
		t.Fatal("no root token")
	}

	for _, key := range keysRaw.([]interface{}) {
		keySlice, err := hex.DecodeString(key.(string))
		if err != nil {
			t.Fatalf("bad: %s", err)
		}

		if _, err := core.Unseal(keySlice); err != nil {
			t.Fatalf("bad: %s", err)
		}
	}

	if core.Sealed() {
		t.Fatal("should not be sealed")
	}
}

func TestSysInit_Put_ValidateParams(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":      5,
		"secret_threshold":   3,
		"recovery_shares":    5,
		"recovery_threshold": 3,
	})
	testResponseStatus(t, resp, http.StatusBadRequest)
	body := map[string][]string{}
	testResponseBody(t, resp, &body)
	if body["errors"][0] != "parameters recovery_shares,recovery_threshold not applicable to seal type shamir" {
		t.Fatal(body)
	}
}

func TestSysInit_Put_ValidateParams_AutoUnseal(t *testing.T) {
	testSeal, _ := seal.NewTestSeal(&seal.TestSealOpts{Name: "transit"})
	autoSeal := vault.NewAutoSeal(testSeal)

	// Create the transit server.
	conf := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
		Seal: autoSeal,
	}
	opts := &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: Handler,
		Logger:      corehelpers.NewTestLogger(t).Named("transit-seal" + strconv.Itoa(0)),
	}
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores
	core := cores[0].Core

	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":      5,
		"secret_threshold":   3,
		"recovery_shares":    5,
		"recovery_threshold": 3,
	})
	testResponseStatus(t, resp, http.StatusBadRequest)
	body := map[string][]string{}
	testResponseBody(t, resp, &body)
	if body["errors"][0] != "parameters secret_shares,secret_threshold not applicable to seal type transit" &&
		body["errors"][0] != "parameters secret_shares,secret_threshold not applicable to seal type test-auto" {
		t.Fatal(body)
	}
}
