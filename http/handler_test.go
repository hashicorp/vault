// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

func TestHandler_parseMFAHandler(t *testing.T) {
	var err error
	var expectedMFACreds logical.MFACreds
	req := &logical.Request{
		Headers: make(map[string][]string),
	}

	headerName := textproto.CanonicalMIMEHeaderKey(MFAHeaderName)

	// Set TOTP passcode in the MFA header
	req.Headers[headerName] = []string{
		"my_totp:123456",
		"my_totp:111111",
		"my_second_mfa:hi=hello",
		"my_third_mfa",
	}
	err = parseMFAHeader(req)
	if err != nil {
		t.Fatal(err)
	}

	// Verify that it is being parsed properly
	expectedMFACreds = logical.MFACreds{
		"my_totp": []string{
			"123456",
			"111111",
		},
		"my_second_mfa": []string{
			"hi=hello",
		},
		"my_third_mfa": []string{},
	}
	if !reflect.DeepEqual(expectedMFACreds, req.MFACreds) {
		t.Fatalf("bad: parsed MFACreds; expected: %#v\n actual: %#v\n", expectedMFACreds, req.MFACreds)
	}

	// Split the creds of a method type in different headers and check if they
	// all get merged together
	req.Headers[headerName] = []string{
		"my_mfa:passcode=123456",
		"my_mfa:month=july",
		"my_mfa:day=tuesday",
	}
	err = parseMFAHeader(req)
	if err != nil {
		t.Fatal(err)
	}

	expectedMFACreds = logical.MFACreds{
		"my_mfa": []string{
			"passcode=123456",
			"month=july",
			"day=tuesday",
		},
	}
	if !reflect.DeepEqual(expectedMFACreds, req.MFACreds) {
		t.Fatalf("bad: parsed MFACreds; expected: %#v\n actual: %#v\n", expectedMFACreds, req.MFACreds)
	}

	// Header without method name should error out
	req.Headers[headerName] = []string{
		":passcode=123456",
	}
	err = parseMFAHeader(req)
	if err == nil {
		t.Fatalf("expected an error; actual: %#v\n", req.MFACreds)
	}

	// Header without method name and method value should error out
	req.Headers[headerName] = []string{
		":",
	}
	err = parseMFAHeader(req)
	if err == nil {
		t.Fatalf("expected an error; actual: %#v\n", req.MFACreds)
	}

	// Header without method name and method value should error out
	req.Headers[headerName] = []string{
		"my_totp:",
	}
	err = parseMFAHeader(req)
	if err == nil {
		t.Fatalf("expected an error; actual: %#v\n", req.MFACreds)
	}
}

// TestHandler_CORS_Patch verifies that http PATCH is included in the list of
// allowed request methods
func TestHandler_CORS_Patch(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	corsConfig := core.CORSConfig()
	err := corsConfig.Enable(context.Background(), []string{addr}, nil)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodOptions, addr+"/v1/sys/seal-status", nil)
	require.NoError(t, err)

	req.Header.Set("Origin", addr)
	req.Header.Set("Access-Control-Request-Method", http.MethodPatch)

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandler_cors(t *testing.T) {
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	// Enable CORS and allow from any origin for testing.
	corsConfig := core.CORSConfig()
	err := corsConfig.Enable(context.Background(), []string{addr}, nil)
	if err != nil {
		t.Fatalf("Error enabling CORS: %s", err)
	}

	req, err := http.NewRequest(http.MethodOptions, addr+"/v1/sys/seal-status", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set("Origin", "BAD ORIGIN")

	// Requests from unacceptable origins will be rejected with a 403.
	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("Bad status:\nexpected: 403 Forbidden\nactual: %s", resp.Status)
	}

	//
	// Test preflight requests
	//

	// Set a valid origin
	req.Header.Set("Origin", addr)

	// Server should NOT accept arbitrary methods.
	req.Header.Set("Access-Control-Request-Method", "FOO")

	client = cleanhttp.DefaultClient()
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Fail if an arbitrary method is accepted.
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("Bad status:\nexpected: 405 Method Not Allowed\nactual: %s", resp.Status)
	}

	// Server SHOULD accept acceptable methods.
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)

	client = cleanhttp.DefaultClient()
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	//
	// Test that the CORS headers are applied correctly.
	//
	expHeaders := map[string]string{
		"Access-Control-Allow-Origin":  addr,
		"Access-Control-Allow-Headers": strings.Join(vault.StdAllowedHeaders, ","),
		"Access-Control-Max-Age":       "300",
		"Vary":                         "Origin",
	}

	for expHeader, expected := range expHeaders {
		actual := resp.Header.Get(expHeader)
		if actual == "" {
			t.Fatalf("bad:\nHeader: %#v was not on response.", expHeader)
		}

		if actual != expected {
			t.Fatalf("bad:\nExpected: %#v\nActual: %#v\n", expected, actual)
		}
	}
}

func TestHandler_HostnameHeader(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		description   string
		config        *vault.CoreConfig
		headerPresent bool
	}{
		{
			description:   "with no header configured",
			config:        nil,
			headerPresent: false,
		},
		{
			description: "with header configured",
			config: &vault.CoreConfig{
				EnableResponseHeaderHostname: true,
			},
			headerPresent: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			var core *vault.Core

			if tc.config == nil {
				core, _, _ = vault.TestCoreUnsealed(t)
			} else {
				core, _, _ = vault.TestCoreUnsealedWithConfig(t, tc.config)
			}

			ln, addr := TestServer(t, core)
			defer ln.Close()

			req, err := http.NewRequest("GET", addr+"/v1/sys/seal-status", nil)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			client := cleanhttp.DefaultClient()
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if resp == nil {
				t.Fatal("nil response")
			}

			hnHeader := resp.Header.Get("X-Vault-Hostname")
			if tc.headerPresent && hnHeader == "" {
				t.Logf("header configured = %t", core.HostnameHeaderEnabled())
				t.Fatal("missing 'X-Vault-Hostname' header entry in response")
			}
			if !tc.headerPresent && hnHeader != "" {
				t.Fatal("didn't expect 'X-Vault-Hostname' header but it was present anyway")
			}

			rniHeader := resp.Header.Get("X-Vault-Raft-Node-ID")
			if rniHeader != "" {
				t.Fatalf("no raft node ID header was expected, since we're not running a raft cluster. instead, got %s", rniHeader)
			}
		})
	}
}

func TestHandler_CacheControlNoStore(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("GET", addr+"/v1/sys/mounts", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(consts.AuthHeaderName, token)
	req.Header.Set(WrapTTLHeaderName, "60s")

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp == nil {
		t.Fatalf("nil response")
	}

	actual := resp.Header.Get("Cache-Control")

	if actual == "" {
		t.Fatalf("missing 'Cache-Control' header entry in response writer")
	}

	if actual != "no-store" {
		t.Fatalf("bad: Cache-Control. Expected: 'no-store', Actual: %q", actual)
	}
}

func TestHandler_InFlightRequest(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	req, err := http.NewRequest("GET", addr+"/v1/sys/in-flight-req", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(consts.AuthHeaderName, token)

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if resp == nil {
		t.Fatalf("nil response")
	}

	var actual map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if actual == nil || len(actual) == 0 {
		t.Fatal("expected to get at least one in-flight request, got nil or zero length map")
	}
	for _, v := range actual {
		reqInfo, ok := v.(map[string]interface{})
		if !ok {
			t.Fatal("failed to read in-flight request")
		}
		if reqInfo["request_path"] != "/v1/sys/in-flight-req" {
			t.Fatalf("expected /v1/sys/in-flight-req in-flight request path, got %s", actual["request_path"])
		}
	}
}

// TestHandler_MissingToken tests the response / error code if a request comes
// in with a missing client token. See
// https://github.com/hashicorp/vault/issues/8377
func TestHandler_MissingToken(t *testing.T) {
	// core, _, token := vault.TestCoreUnsealed(t)
	core, _, _ := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("GET", addr+"/v1/sys/internal/ui/mounts/cubbyhole", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	req.Header.Set(WrapTTLHeaderName, "60s")

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("expected code 403, got: %d", resp.StatusCode)
	}
}

func TestHandler_Accepted(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("POST", addr+"/v1/auth/token/tidy", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(consts.AuthHeaderName, token)

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testResponseStatus(t, resp, 202)
}

// We use this test to verify header auth
func TestSysMounts_headerAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("GET", addr+"/v1/sys/mounts", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(consts.AuthHeaderName, token)

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"mount_type":     "system",
		"auth":           nil,
		"data": map[string]interface{}{
			"secret/": map[string]interface{}{
				"description":             "key/value secret storage",
				"type":                    "kv",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                map[string]interface{}{"version": "1"},
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
			},
			"sys/": map[string]interface{}{
				"description":             "system endpoints used for control, policy and debugging",
				"type":                    "system",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl":           json.Number("0"),
					"max_lease_ttl":               json.Number("0"),
					"force_no_cache":              false,
					"passthrough_request_headers": []interface{}{"Accept"},
				},
				"local":                  false,
				"seal_wrap":              true,
				"options":                interface{}(nil),
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.DefaultBuiltinVersion,
			},
			"cubbyhole/": map[string]interface{}{
				"description":             "per-token private secret storage",
				"type":                    "cubbyhole",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
				},
				"local":                  true,
				"seal_wrap":              false,
				"options":                interface{}(nil),
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "cubbyhole"),
			},
			"identity/": map[string]interface{}{
				"description":             "identity store",
				"type":                    "identity",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl":           json.Number("0"),
					"max_lease_ttl":               json.Number("0"),
					"force_no_cache":              false,
					"passthrough_request_headers": []interface{}{"Authorization"},
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                interface{}(nil),
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "identity"),
			},
			"agent-registry/": map[string]interface{}{
				"description":             "agent registry",
				"type":                    "agent_registry",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl":           json.Number("0"),
					"max_lease_ttl":               json.Number("0"),
					"force_no_cache":              false,
					"passthrough_request_headers": []interface{}{"Authorization"},
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                interface{}(nil),
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.DefaultBuiltinVersion,
			},
		},
		"secret/": map[string]interface{}{
			"description":             "key/value secret storage",
			"type":                    "kv",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                map[string]interface{}{"version": "1"},
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
		},
		"agent-registry/": map[string]interface{}{
			"description":             "agent registry",
			"type":                    "agent_registry",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl":           json.Number("0"),
				"max_lease_ttl":               json.Number("0"),
				"force_no_cache":              false,
				"passthrough_request_headers": []interface{}{"Authorization"},
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                interface{}(nil),
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.DefaultBuiltinVersion,
		},
		"sys/": map[string]interface{}{
			"description":             "system endpoints used for control, policy and debugging",
			"type":                    "system",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl":           json.Number("0"),
				"max_lease_ttl":               json.Number("0"),
				"force_no_cache":              false,
				"passthrough_request_headers": []interface{}{"Accept"},
			},
			"local":                  false,
			"seal_wrap":              true,
			"options":                interface{}(nil),
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.DefaultBuiltinVersion,
		},
		"cubbyhole/": map[string]interface{}{
			"description":             "per-token private secret storage",
			"type":                    "cubbyhole",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
			},
			"local":                  true,
			"seal_wrap":              false,
			"options":                interface{}(nil),
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "cubbyhole"),
		},
		"identity/": map[string]interface{}{
			"description":             "identity store",
			"type":                    "identity",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl":           json.Number("0"),
				"max_lease_ttl":               json.Number("0"),
				"force_no_cache":              false,
				"passthrough_request_headers": []interface{}{"Authorization"},
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                interface{}(nil),
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "identity"),
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected["request_id"] = actual["request_id"]
	for k, v := range actual["data"].(map[string]interface{}) {
		if v.(map[string]interface{})["accessor"] == "" {
			t.Fatalf("no accessor from %s", k)
		}
		if v.(map[string]interface{})["uuid"] == "" {
			t.Fatalf("no uuid from %s", k)
		}

		expected[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
		expected[k].(map[string]interface{})["uuid"] = v.(map[string]interface{})["uuid"]
		expected["data"].(map[string]interface{})[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
		expected["data"].(map[string]interface{})[k].(map[string]interface{})["uuid"] = v.(map[string]interface{})["uuid"]
	}

	if diff := deep.Equal(actual, expected); len(diff) > 0 {
		t.Fatalf("bad, diff: %#v", diff)
	}
}

// We use this test to verify header auth wrapping
func TestSysMounts_headerAuth_Wrapped(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	req, err := http.NewRequest("GET", addr+"/v1/sys/mounts", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(consts.AuthHeaderName, token)
	req.Header.Set(WrapTTLHeaderName, "60s")

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"request_id":     "",
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data":           nil,
		"wrap_info": map[string]interface{}{
			"ttl": json.Number("60"),
		},
		"warnings":   nil,
		"auth":       nil,
		"mount_type": "",
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	actualToken, ok := actual["wrap_info"].(map[string]interface{})["token"]
	if !ok || actualToken == "" {
		t.Fatal("token missing in wrap info")
	}
	expected["wrap_info"].(map[string]interface{})["token"] = actualToken

	actualCreationTime, ok := actual["wrap_info"].(map[string]interface{})["creation_time"]
	if !ok || actualCreationTime == "" {
		t.Fatal("creation_time missing in wrap info")
	}
	expected["wrap_info"].(map[string]interface{})["creation_time"] = actualCreationTime

	actualCreationPath, ok := actual["wrap_info"].(map[string]interface{})["creation_path"]
	if !ok || actualCreationPath == "" {
		t.Fatal("creation_path missing in wrap info")
	}
	expected["wrap_info"].(map[string]interface{})["creation_path"] = actualCreationPath

	actualAccessor, ok := actual["wrap_info"].(map[string]interface{})["accessor"]
	if !ok || actualAccessor == "" {
		t.Fatal("accessor missing in wrap info")
	}
	expected["wrap_info"].(map[string]interface{})["accessor"] = actualAccessor

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual: %#v\n%T %T", expected, actual, actual["warnings"], actual["data"])
	}
}

func TestHandler_sealed(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	core.Seal(token)

	resp, err := http.Get(addr + "/v1/secret/foo")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 503)
}

func TestHandler_ui_default(t *testing.T) {
	core := vault.TestCoreUI(t, false)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/ui/")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 404)
}

func TestHandler_ui_enabled(t *testing.T) {
	core := vault.TestCoreUI(t, true)
	ln, addr := TestServer(t, core)
	defer ln.Close()

	resp, err := http.Get(addr + "/ui/")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	testResponseStatus(t, resp, 200)
}

func TestHandler_error(t *testing.T) {
	w := httptest.NewRecorder()

	respondError(w, 500, errors.New("test Error"))

	if w.Code != 500 {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	// The code inside of the error should override
	// the argument to respondError
	w2 := httptest.NewRecorder()
	e := logical.CodedError(403, "error text")

	respondError(w2, 500, e)

	if w2.Code != 403 {
		t.Fatalf("expected 403, got %d", w2.Code)
	}

	// vault.ErrSealed is a special case
	w3 := httptest.NewRecorder()

	respondError(w3, 400, consts.ErrSealed)

	if w3.Code != 503 {
		t.Fatalf("expected 503, got %d", w3.Code)
	}
}

func TestHandler_requestAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)

	rootCtx := namespace.RootContext(nil)
	te, err := core.LookupToken(rootCtx, token)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	rWithAuthorization, err := http.NewRequest("GET", "v1/test/path", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	rWithAuthorization.Header.Set("Authorization", "Bearer "+token)

	rWithVault, err := http.NewRequest("GET", "v1/test/path", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	rWithVault.Header.Set(consts.AuthHeaderName, token)

	for _, r := range []*http.Request{rWithVault, rWithAuthorization} {
		req := logical.TestRequest(t, logical.ReadOperation, "test/path")
		r = r.WithContext(rootCtx)
		requestAuth(r, req)
		err = core.PopulateTokenEntry(rootCtx, req)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		if req.ClientToken != token {
			t.Fatalf("client token should be filled with %s, got %s", token, req.ClientToken)
		}
		if req.TokenEntry() == nil {
			t.Fatal("token entry should not be nil")
		}
		if !reflect.DeepEqual(req.TokenEntry(), te) {
			t.Fatalf("token entry should be the same as the core")
		}
		if req.ClientTokenAccessor == "" {
			t.Fatal("token accessor should not be empty")
		}
	}

	rNothing, err := http.NewRequest("GET", "v1/test/path", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req := logical.TestRequest(t, logical.ReadOperation, "test/path")

	requestAuth(rNothing, req)
	err = core.PopulateTokenEntry(rootCtx, req)
	if err != nil {
		t.Fatalf("expected no error, got %s", err)
	}
	if req.ClientToken != "" {
		t.Fatalf("client token should not be filled, got %s", req.ClientToken)
	}
}

func TestHandler_getTokenFromReq(t *testing.T) {
	r := http.Request{Header: http.Header{}}

	tok, _ := getTokenFromReq(&r)
	if tok != "" {
		t.Fatalf("expected '' as result, got '%s'", tok)
	}

	r.Header.Set("Authorization", "Bearer TOKEN NOT_GOOD_TOKEN")
	token, fromHeader := getTokenFromReq(&r)
	if !fromHeader {
		t.Fatal("expected from header")
	} else if token != "TOKEN NOT_GOOD_TOKEN" {
		t.Fatal("did not get expected token value")
	} else if r.Header.Get("Authorization") == "" {
		t.Fatal("expected value to be passed through")
	}

	r.Header.Set(consts.AuthHeaderName, "NEWTOKEN")
	tok, _ = getTokenFromReq(&r)
	if tok == "TOKEN" {
		t.Fatalf("%s header should be prioritized", consts.AuthHeaderName)
	} else if tok != "NEWTOKEN" {
		t.Fatalf("expected 'NEWTOKEN' as result, got '%s'", tok)
	}

	r.Header = http.Header{}
	r.Header.Set("Authorization", "Basic TOKEN")
	tok, fromHeader = getTokenFromReq(&r)
	if tok != "" {
		t.Fatalf("expected '' as result, got '%s'", tok)
	} else if fromHeader {
		t.Fatal("expected not from header")
	}
}

func TestHandler_nonPrintableChars(t *testing.T) {
	testNonPrintable(t, false)
	testNonPrintable(t, true)
}

func testNonPrintable(t *testing.T, disable bool) {
	core, _, token := vault.TestCoreUnsealedWithConfig(t, &vault.CoreConfig{
		DisableKeyEncodingChecks: disable,
	})
	ln, addr := TestListener(t)
	props := &vault.HandlerProperties{
		Core:                  core,
		DisablePrintableCheck: disable,
		ListenerConfig:        &configutil.Listener{},
	}
	TestServerWithListenerAndProperties(t, ln, addr, core, props)
	defer ln.Close()

	req, err := http.NewRequest("PUT", addr+"/v1/cubbyhole/foo\u2028bar", strings.NewReader(`{"zip": "zap"}`))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set(consts.AuthHeaderName, token)

	client := cleanhttp.DefaultClient()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if disable {
		testResponseStatus(t, resp, 204)
	} else {
		testResponseStatus(t, resp, 400)
	}
}

func TestHandler_Parse_Form(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	core := cores[0].Core
	vault.TestWaitActive(t, core)

	c := cleanhttp.DefaultClient()
	c.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: cluster.RootCAs,
		},
	}

	values := url.Values{
		"zip":   []string{"zap"},
		"abc":   []string{"xyz"},
		"multi": []string{"first", "second"},
		"empty": []string{},
	}
	req, err := http.NewRequest("POST", cores[0].Client.Address()+"/v1/secret/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Body = io.NopCloser(strings.NewReader(values.Encode()))
	req.Header.Set("x-vault-token", cluster.RootToken)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 204 {
		t.Fatalf("bad response: %#v\nrequest was: %#v\nurl was: %#v", *resp, *req, req.URL)
	}

	client := cores[0].Client
	client.SetToken(cluster.RootToken)

	apiResp, err := client.Logical().Read("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if apiResp == nil {
		t.Fatal("api resp is nil")
	}
	expected := map[string]interface{}{
		"zip":   "zap",
		"abc":   "xyz",
		"multi": "first,second",
	}
	if diff := deep.Equal(expected, apiResp.Data); diff != nil {
		t.Fatal(diff)
	}
}

// TestHandler_MaxRequestSize verifies that a request larger than the
// MaxRequestSize fails
func TestHandler_MaxRequestSize(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		DefaultHandlerProperties: vault.HandlerProperties{
			ListenerConfig: &configutil.Listener{
				MaxRequestSize: 1024,
			},
		},
		HandlerFunc: Handler,
		NumCores:    1,
	})
	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[0].Client
	_, err := client.KVv2("secret").Put(context.Background(), "foo", map[string]interface{}{
		"bar": strings.Repeat("a", 1025),
	})

	require.ErrorContains(t, err, "http: request body too large")
}

// TestHandler_MaxRequestSize_Memory sets the max request size to 1024 bytes,
// and creates a 1MB request. The test verifies that less than 1MB of memory is
// allocated when the request is sent. This test shouldn't be run in parallel,
// because it modifies GOMAXPROCS
func TestHandler_MaxRequestSize_Memory(t *testing.T) {
	ln, addr := TestListener(t)
	core, _, token := vault.TestCoreUnsealed(t)
	TestServerWithListenerAndProperties(t, ln, addr, core, &vault.HandlerProperties{
		Core: core,
		ListenerConfig: &configutil.Listener{
			Address:        addr,
			MaxRequestSize: 1024,
		},
	})
	defer ln.Close()

	data := bytes.Repeat([]byte{0x1}, 1024*1024)

	req, err := http.NewRequest("POST", addr+"/v1/sys/unseal", bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set(consts.AuthHeaderName, token)

	client := cleanhttp.DefaultClient()
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	var start, end runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&start)
	client.Do(req)
	runtime.ReadMemStats(&end)
	require.Less(t, end.TotalAlloc-start.TotalAlloc, uint64(1024*1024))
}

// Test_requiresSnapshot verifies that a request is marked as requiring a
// snapshot when it's a read, list, or create/update and has a snapshot query
// parameter
func Test_requiresSnapshot(t *testing.T) {
	testCases := []struct {
		name        string
		method      string
		queryParams map[string][]string
		expected    bool
	}{
		{
			name:        "get no snapshot",
			method:      http.MethodGet,
			queryParams: map[string][]string{"other": {"param"}},
			expected:    false,
		},
		{
			name:        "options with snapshot",
			method:      http.MethodOptions,
			queryParams: map[string][]string{VaultSnapshotRecoverParam: {"param"}},
			expected:    false,
		},
		{
			name:        "put with read snapshot",
			method:      http.MethodPut,
			queryParams: map[string][]string{VaultSnapshotReadParam: {"param"}},
			expected:    false,
		},
		{
			name:        "put with recover snapshot",
			method:      http.MethodPut,
			queryParams: map[string][]string{VaultSnapshotRecoverParam: {"param"}},
			expected:    true,
		},
		{
			name:        "list with snapshot",
			method:      "LIST",
			queryParams: map[string][]string{VaultSnapshotReadParam: {"param"}},
			expected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &http.Request{
				Method: tc.method,
				URL: &url.URL{
					RawQuery: url.Values(tc.queryParams).Encode(),
				},
			}
			require.Equal(t, tc.expected, requiresSnapshot(req))
		})
	}
}

// TestHandler_JSONLimitQuotaWrappers verifies that the handler properly orders
// the normal quota checks, JSON size limits checks, and role-based quota checks
func TestHandler_JSONLimitQuotaWrappers(t *testing.T) {
	testCases := []struct {
		name           string
		setup          func(t *testing.T, client *api.Client, roleID string)
		jsonStringSize int
		wantError      string
	}{
		{
			// set up a role-based rate limit, but don't exceed the rate
			// because the JSON is too big the request will error with
			// the JSON size error
			name: "too big json with role quota",
			setup: func(t *testing.T, client *api.Client, _ string) {
				_, err := client.Logical().Write("sys/quotas/rate-limit/my-role-quota", map[string]interface{}{
					"path": "auth/approle",
					"role": "my-role",
					"rate": 5,
				})
				require.NoError(t, err)
			},
			jsonStringSize: 5001,
			wantError:      "JSON string value exceeds allowed length",
		},
		{
			// set up a rate limit, but don't exceed the rate
			// because the JSON is too big the request will error with
			// the JSON size error
			name: "too big json with non role quota",
			setup: func(t *testing.T, client *api.Client, _ string) {
				_, err := client.Logical().Write("sys/quotas/rate-limit/my-quota", map[string]interface{}{
					"path": "auth/approle",
					"rate": 5,
				})
				require.NoError(t, err)
			},
			jsonStringSize: 5001,
			wantError:      "JSON string value exceeds allowed length",
		},
		{
			// set up a rate limit without a role and exceed it
			// even though the JSON is too big, the request will be blocked by
			// the rate limit first
			name: "too big json with non role quota blocked",
			setup: func(t *testing.T, client *api.Client, roleID string) {
				_, err := client.Logical().Write("sys/quotas/rate-limit/my-quota", map[string]interface{}{
					"path":     "auth/approle",
					"rate":     1,
					"interval": "60",
				})
				require.NoError(t, err)
			},
			jsonStringSize: 5001,
			wantError:      "rate limit quota exceeded",
		},
		{
			// set up a rate limit with a role and exceed it
			// the JSON is too big and the JSON check will trigger before
			// the role-based quota check, so we'll get a JSON size error
			name: "too big json with role quota blocked",
			setup: func(t *testing.T, client *api.Client, roleID string) {
				_, err := client.Logical().Write("sys/quotas/rate-limit/my-role-quota", map[string]interface{}{
					"path":     "auth/approle",
					"role":     "my-role",
					"rate":     1,
					"interval": "60",
				})
				require.NoError(t, err)

				// log in for the role to use up the quota
				r, err := client.Logical().Write("auth/approle/role/my-role/secret-id", nil)
				require.NoError(t, err)
				secretID := r.Data["secret_id"].(string)

				_, err = client.Logical().Write("auth/approle/login", map[string]interface{}{
					"role_id":   roleID,
					"secret_id": secretID,
				})
				require.NoError(t, err)
			},
			jsonStringSize: 5001,
			wantError:      "JSON string value exceeds allowed length",
		},
		{
			// set up a rate limit with a role and exceed it
			// the JSON is an ok size, so the role-based quota check will
			// trigger
			name: "normal json with role quota blocked",
			setup: func(t *testing.T, client *api.Client, roleID string) {
				_, err := client.Logical().Write("sys/quotas/rate-limit/my-role-quota", map[string]interface{}{
					"path":     "auth/approle",
					"role":     "my-role",
					"rate":     1,
					"interval": "60",
				})
				require.NoError(t, err)

				// log in for the role to use up the quota
				r, err := client.Logical().Write("auth/approle/role/my-role/secret-id", nil)
				require.NoError(t, err)
				secretID := r.Data["secret_id"].(string)

				_, err = client.Logical().Write("auth/approle/login", map[string]interface{}{
					"role_id":   roleID,
					"secret_id": secretID,
				})
				require.NoError(t, err)
			},
			jsonStringSize: 5,
			wantError:      "rate limit quota exceeded",
		},
		{
			// set up a rate limit with a role but don't exceed it
			// the JSON is an ok size, so the request will succeed
			name: "normal json with role quota allowed",
			setup: func(t *testing.T, client *api.Client, roleID string) {
				_, err := client.Logical().Write("sys/quotas/rate-limit/my-role-quota", map[string]interface{}{
					"path":     "auth/approle",
					"role":     "my-role",
					"rate":     1,
					"interval": "60",
				})
				require.NoError(t, err)
			},
			jsonStringSize: 5,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
				HandlerFunc: Handler,
				DefaultHandlerProperties: vault.HandlerProperties{
					ListenerConfig: &configutil.Listener{
						CustomMaxJSONStringValueLength: 5000,
					},
				},
			})
			cluster.Start()
			defer cluster.Cleanup()

			client := cluster.Cores[0].Client
			client.SetToken(cluster.RootToken)

			err := client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
				Type: "approle",
			})
			require.NoError(t, err)

			_, err = client.Logical().Write("auth/approle/role/my-role", map[string]interface{}{
				"token_policies": "default",
				"token_ttl":      "1h",
				"token_max_ttl":  "4h",
			})
			r, err := client.Logical().Read("auth/approle/role/my-role/role-id")
			require.NoError(t, err)
			roleID := r.Data["role_id"].(string)
			require.NoError(t, err)
			if tc.setup != nil {
				tc.setup(t, client, roleID)
			}
			r, err = client.Logical().Write("auth/approle/role/my-role/secret-id", nil)
			require.NoError(t, err)
			secretID := r.Data["secret_id"].(string)

			resp, err := client.Logical().Write("auth/approle/login", map[string]interface{}{
				"role_id":         roleID,
				"secret_id":       secretID,
				"additional_data": strings.Repeat("a", tc.jsonStringSize),
			})
			if tc.wantError != "" {
				require.ErrorContains(t, err, tc.wantError)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
		})
	}
}

// TestAutoSnapshotLoadForwarded tests that a request to load from a cloud
// snapshot is forwarded to the active node, rather than being redirected
func TestAutoSnapshotLoadForwarded(t *testing.T) {
	cluster := vault.NewTestCluster(t, &vault.CoreConfig{}, &vault.TestClusterOptions{
		NumCores:    2,
		HandlerFunc: Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	client := cluster.Cores[1].Client
	client.SetToken(cluster.RootToken)

	_, err := client.Logical().Write("sys/storage/raft/snapshot-auto/snapshot-load/cfg1", nil)
	// the request will fail, but all that we care about is that the error
	// doesn't indicate a redirect
	require.Error(t, err)
	require.NotContains(t, err.Error(), "redirects not allowed in these tests")
}
