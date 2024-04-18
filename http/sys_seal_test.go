// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/vault/seal"
	"github.com/hashicorp/vault/version"
	"github.com/stretchr/testify/assert"
)

func TestSysSealStatus(t *testing.T) {
	core := vault.TestCore(t)
	vault.TestCoreInit(t, core)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/seal-status")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"sealed":        true,
		"t":             json.Number("3"),
		"n":             json.Number("3"),
		"progress":      json.Number("0"),
		"nonce":         "",
		"type":          "shamir",
		"recovery_seal": false,
		"initialized":   true,
		"migration":     false,
		"build_date":    version.BuildDate,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual["version"] == nil {
		t.Fatalf("expected version information")
	}
	expected["version"] = actual["version"]
	if actual["cluster_name"] == nil {
		delete(expected, "cluster_name")
	} else {
		expected["cluster_name"] = actual["cluster_name"]
	}
	if actual["cluster_id"] == nil {
		delete(expected, "cluster_id")
	} else {
		expected["cluster_id"] = actual["cluster_id"]
	}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestSysSealStatus_uninit(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/v1/sys/seal-status")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 200)
}

func TestSysSeal(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/seal", nil)
	testResponseStatus(t, resp, 204)

	if !core.Sealed() {
		t.Fatal("should be sealed")
	}
}

func TestSysSeal_unsealed(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/seal", nil)
	testResponseStatus(t, resp, 204)

	if !core.Sealed() {
		t.Fatal("should be sealed")
	}
}

func TestSysUnseal(t *testing.T) {
	core := vault.TestCore(t)
	keys, _ := vault.TestCoreInit(t, core)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	for i, key := range keys {
		resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
			"key": hex.EncodeToString(key),
		})

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"sealed":        true,
			"t":             json.Number("3"),
			"n":             json.Number("3"),
			"progress":      json.Number(fmt.Sprintf("%d", i+1)),
			"nonce":         "",
			"type":          "shamir",
			"recovery_seal": false,
			"initialized":   true,
			"migration":     false,
			"build_date":    version.BuildDate,
		}
		if i == len(keys)-1 {
			expected["sealed"] = false
			expected["progress"] = json.Number("0")
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if i < len(keys)-1 && (actual["nonce"] == nil || actual["nonce"].(string) == "") {
			t.Fatalf("got nil nonce, actual is %#v", actual)
		} else {
			expected["nonce"] = actual["nonce"]
		}
		if actual["version"] == nil {
			t.Fatalf("expected version information")
		}
		expected["version"] = actual["version"]
		if actual["cluster_name"] == nil {
			delete(expected, "cluster_name")
		} else {
			expected["cluster_name"] = actual["cluster_name"]
		}
		if actual["cluster_id"] == nil {
			delete(expected, "cluster_id")
		} else {
			expected["cluster_id"] = actual["cluster_id"]
		}
		if diff := deep.Equal(actual, expected); diff != nil {
			t.Fatal(diff)
		}
	}
}

func subtestBadSingleKey(t *testing.T, seal vault.Seal) {
	core := vault.TestCoreWithSeal(t, seal, false)
	_, err := core.Initialize(context.Background(), &vault.InitParams{
		BarrierConfig: &vault.SealConfig{
			SecretShares:    1,
			SecretThreshold: 1,
		},
		RecoveryConfig: &vault.SealConfig{
			SecretShares:    1,
			SecretThreshold: 1,
		},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	ln, addr := TestServer(t, core)
	defer ln.Close()

	testCases := []struct {
		description string
		key         string
	}{
		// hex key tests
		// hexadecimal strings have 2 symbols per byte; size(0xAA) == 1 byte
		{
			"short hex key",
			strings.Repeat("AA", 8),
		},
		{
			"long hex key",
			strings.Repeat("AA", 34),
		},
		{
			"uneven hex key byte length",
			strings.Repeat("AA", 33),
		},
		{
			"valid hex key but wrong cluster",
			"4482691dd3a710723c4f77c4920ee21b96c226bf4829fa6eb8e8262c180ae933",
		},

		// base64 key tests
		// base64 strings have min. 1 character per byte; size("m") == 1 byte
		{
			"short b64 key",
			base64.StdEncoding.EncodeToString([]byte(strings.Repeat("m", 8))),
		},
		{
			"long b64 key",
			base64.StdEncoding.EncodeToString([]byte(strings.Repeat("m", 34))),
		},
		{
			"uneven b64 key byte length",
			base64.StdEncoding.EncodeToString([]byte(strings.Repeat("m", 33))),
		},
		{
			"valid b64 key but wrong cluster",
			"RIJpHdOnEHI8T3fEkg7iG5bCJr9IKfpuuOgmLBgK6TM=",
		},

		// other key tests
		{
			"empty key",
			"",
		},
		{
			"key with bad format",
			"ThisKeyIsNeitherB64NorHex",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
				"key": tc.key,
			})

			testResponseStatus(t, resp, 400)
		})
	}
}

func subtestBadMultiKey(t *testing.T, seal vault.Seal) {
	numKeys := 3

	core := vault.TestCoreWithSeal(t, seal, false)
	_, err := core.Initialize(context.Background(), &vault.InitParams{
		BarrierConfig: &vault.SealConfig{
			SecretShares:    numKeys,
			SecretThreshold: numKeys,
		},
		RecoveryConfig: &vault.SealConfig{
			SecretShares:    numKeys,
			SecretThreshold: numKeys,
		},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	ln, addr := TestServer(t, core)
	defer ln.Close()

	testCases := []struct {
		description string
		keys        []string
	}{
		{
			"all unseal keys from another cluster",
			[]string{
				"b189d98fdec3a15bed9b1cce5088f82b92896696b788c07bdf03c73da08279a5e8",
				"0fa98232f034177d8d9c2824899a2ac1e55dc6799348533e10510b856aef99f61a",
				"5344f5caa852f9ba1967d9623ed286a45ea7c4a529522d25f05d29ff44f17930ac",
			},
		},
		{
			"mixing unseal keys from different cluster, different share config",
			[]string{
				"b189d98fdec3a15bed9b1cce5088f82b92896696b788c07bdf03c73da08279a5e8",
				"0fa98232f034177d8d9c2824899a2ac1e55dc6799348533e10510b856aef99f61a",
				"e04ea3020838c2050c4a169d7ba4d30e034eec8e83e8bed9461bf2646ee412c0",
			},
		},
		{
			"mixing unseal keys from different clusters, similar share config",
			[]string{
				"b189d98fdec3a15bed9b1cce5088f82b92896696b788c07bdf03c73da08279a5e8",
				"0fa98232f034177d8d9c2824899a2ac1e55dc6799348533e10510b856aef99f61a",
				"413f80521b393aa6c4e42e9a3a3ab7f00c2002b2c3bf1e273fc6f363f35f2a378b",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			for i, key := range tc.keys {
				resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
					"key": key,
				})

				if i == numKeys-1 {
					// last key
					testResponseStatus(t, resp, 400)
				} else {
					// unseal in progress
					testResponseStatus(t, resp, 200)
				}

			}
		})
	}
}

func TestSysUnseal_BadKeyNewShamir(t *testing.T) {
	seal := vault.NewTestSeal(t,
		&seal.TestSealOpts{StoredKeys: seal.StoredKeysSupportedShamirRoot})

	subtestBadSingleKey(t, seal)
	subtestBadMultiKey(t, seal)
}

func TestSysUnseal_BadKeyAutoUnseal(t *testing.T) {
	seal := vault.NewTestSeal(t,
		&seal.TestSealOpts{StoredKeys: seal.StoredKeysSupportedGeneric})

	subtestBadSingleKey(t, seal)
	subtestBadMultiKey(t, seal)
}

func TestSysUnseal_Reset(t *testing.T) {
	core := vault.TestCore(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	thresh := 3
	resp := testHttpPut(t, "", addr+"/v1/sys/init", map[string]interface{}{
		"secret_shares":    5,
		"secret_threshold": thresh,
	})

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	keysRaw, ok := actual["keys"]
	if !ok {
		t.Fatalf("no keys: %#v", actual)
	}
	for i, key := range keysRaw.([]interface{}) {
		if i > thresh-2 {
			break
		}

		resp := testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
			"key": key.(string),
		})

		var actual map[string]interface{}
		expected := map[string]interface{}{
			"sealed":        true,
			"t":             json.Number("3"),
			"n":             json.Number("5"),
			"progress":      json.Number(strconv.Itoa(i + 1)),
			"type":          "shamir",
			"recovery_seal": false,
			"initialized":   true,
			"migration":     false,
			"build_date":    version.BuildDate,
		}
		testResponseStatus(t, resp, 200)
		testResponseBody(t, resp, &actual)
		if actual["version"] == nil {
			t.Fatalf("expected version information")
		}
		expected["version"] = actual["version"]
		if actual["nonce"] == "" && expected["sealed"].(bool) {
			t.Fatalf("expected a nonce")
		}
		expected["nonce"] = actual["nonce"]
		if actual["cluster_name"] == nil {
			delete(expected, "cluster_name")
		} else {
			expected["cluster_name"] = actual["cluster_name"]
		}
		if actual["cluster_id"] == nil {
			delete(expected, "cluster_id")
		} else {
			expected["cluster_id"] = actual["cluster_id"]
		}
		if diff := deep.Equal(actual, expected); diff != nil {
			t.Fatal(diff)
		}
	}

	resp = testHttpPut(t, "", addr+"/v1/sys/unseal", map[string]interface{}{
		"reset": true,
	})

	actual = map[string]interface{}{}
	expected := map[string]interface{}{
		"sealed":        true,
		"t":             json.Number("3"),
		"n":             json.Number("5"),
		"progress":      json.Number("0"),
		"type":          "shamir",
		"recovery_seal": false,
		"initialized":   true,
		"build_date":    version.BuildDate,
		"migration":     false,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual["version"] == nil {
		t.Fatalf("expected version information")
	}
	expected["version"] = actual["version"]
	expected["nonce"] = actual["nonce"]
	if actual["cluster_name"] == nil {
		delete(expected, "cluster_name")
	} else {
		expected["cluster_name"] = actual["cluster_name"]
	}
	if actual["cluster_id"] == nil {
		delete(expected, "cluster_id")
	} else {
		expected["cluster_id"] = actual["cluster_id"]
	}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}

// Test Seal's permissions logic, which is slightly different than normal code
// paths in that it queries the ACL rather than having checkToken do it. This
// is because it was abusing RootPaths in logical_system, but that caused some
// haywire with code paths that expected there to be an actual corresponding
// logical.Path for it. This way is less hacky, but this test ensures that we
// have not opened up a permissions hole.
func TestSysSeal_Permissions(t *testing.T) {
	core, _, root := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, root)

	// Set the 'test' policy object to permit write access to sys/seal
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test",
		Data: map[string]interface{}{
			"rules": `path "sys/seal" { capabilities = ["read"] }`,
		},
		ClientToken: root,
	}
	resp, err := core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}

	// Create a non-root token with access to that policy
	req.Path = "auth/token/create"
	req.Data = map[string]interface{}{
		"id":       "child",
		"policies": []string{"test"},
	}

	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v %v", err, resp)
	}
	if resp.Auth.ClientToken != "child" {
		t.Fatalf("bad: %#v", resp)
	}

	// We must go through the HTTP interface since seal doesn't go through HandleRequest

	// We expect this to fail since it needs update and sudo
	httpResp := testHttpPut(t, "child", addr+"/v1/sys/seal", nil)
	testResponseStatus(t, httpResp, 403)

	// Now modify to add update capability
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test",
		Data: map[string]interface{}{
			"rules": `path "sys/seal" { capabilities = ["update"] }`,
		},
		ClientToken: root,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}

	// We expect this to fail since it needs sudo
	httpResp = testHttpPut(t, "child", addr+"/v1/sys/seal", nil)
	testResponseStatus(t, httpResp, 403)

	// Now modify to just sudo capability
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test",
		Data: map[string]interface{}{
			"rules": `path "sys/seal" { capabilities = ["sudo"] }`,
		},
		ClientToken: root,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}

	// We expect this to fail since it needs update
	httpResp = testHttpPut(t, "child", addr+"/v1/sys/seal", nil)
	testResponseStatus(t, httpResp, 403)

	// Now modify to add all needed capabilities
	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "sys/policy/test",
		Data: map[string]interface{}{
			"rules": `path "sys/seal" { capabilities = ["update", "sudo"] }`,
		},
		ClientToken: root,
	}
	resp, err = core.HandleRequest(namespace.RootContext(nil), req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil || resp.IsError() {
		t.Fatalf("bad: %#v", resp)
	}

	// We expect this to work
	httpResp = testHttpPut(t, "child", addr+"/v1/sys/seal", nil)
	testResponseStatus(t, httpResp, 204)
}

func TestSysStepDown(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/step-down", nil)
	testResponseStatus(t, resp, 204)
}

// TestSysSealStatusRedaction tests that the response from a
// a request to sys/seal-status are redacted only if no valid token
// is provided with the request
func TestSysSealStatusRedaction(t *testing.T) {
	conf := &vault.CoreConfig{
		EnableUI:        false,
		EnableRaw:       true,
		BuiltinRegistry: corehelpers.NewMockBuiltinRegistry(),
		AuditBackends: map[string]audit.Factory{
			"file": audit.NewFileBackend,
		},
	}
	core, _, token := vault.TestCoreUnsealedWithConfig(t, conf)

	// Setup new custom listener
	ln, addr := TestListener(t)
	props := &vault.HandlerProperties{
		Core: core,
		ListenerConfig: &configutil.Listener{
			RedactVersion: true,
		},
	}
	TestServerWithListenerAndProperties(t, ln, addr, core, props)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	client := cleanhttp.DefaultClient()

	// Check seal-status
	req, err := http.NewRequest("GET", addr+"/v1/sys/seal-status", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(consts.AuthHeaderName, token)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 200)

	// Verify that version exists when provided a valid token
	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	assert.NotEmpty(t, actual["version"])

	// Verify that version is redacted when no token is provided
	req, err = http.NewRequest("GET", addr+"/v1/sys/seal-status", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(consts.AuthHeaderName, "")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	assert.Empty(t, actual["version"])
}
