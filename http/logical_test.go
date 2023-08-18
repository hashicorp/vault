// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"

	"github.com/go-test/deep"
	log "github.com/hashicorp/go-hclog"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/vault"
)

func TestLogical(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// WRITE
	resp := testHttpPut(t, token, addr+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
	})
	testResponseStatus(t, resp, 204)

	// READ
	// Bad token should return a 403
	resp = testHttpGet(t, token+"bad", addr+"/v1/secret/foo")
	testResponseStatus(t, resp, 403)

	resp = testHttpGet(t, token, addr+"/v1/secret/foo")
	var actual map[string]interface{}
	var nilWarnings interface{}
	expected := map[string]interface{}{
		"renewable":      false,
		"lease_duration": json.Number(strconv.Itoa(int((32 * 24 * time.Hour) / time.Second))),
		"data": map[string]interface{}{
			"data": "bar",
		},
		"auth":      nil,
		"wrap_info": nil,
		"warnings":  nilWarnings,
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	delete(actual, "lease_id")
	expected["request_id"] = actual["request_id"]
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}

	// DELETE
	resp = testHttpDelete(t, token, addr+"/v1/secret/foo")
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/secret/foo")
	testResponseStatus(t, resp, 404)
}

func TestLogical_noExist(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, token, addr+"/v1/secret/foo")
	testResponseStatus(t, resp, 404)
}

func TestLogical_StandbyRedirect(t *testing.T) {
	ln1, addr1 := TestListener(t)
	defer ln1.Close()
	ln2, addr2 := TestListener(t)
	defer ln2.Close()

	// Create an HA Vault
	logger := logging.NewVaultLogger(log.Debug)

	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	conf := &vault.CoreConfig{
		Physical:     inmha,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: addr1,
		DisableMlock: true,
	}
	core1, err := vault.NewCore(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer core1.Shutdown()
	keys, root := vault.TestCoreInit(t, core1)
	for _, key := range keys {
		if _, err := core1.Unseal(vault.TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	// Attempt to fix raciness in this test by giving the first core a chance
	// to grab the lock
	time.Sleep(2 * time.Second)

	// Create a second HA Vault
	conf2 := &vault.CoreConfig{
		Physical:     inmha,
		HAPhysical:   inmha.(physical.HABackend),
		RedirectAddr: addr2,
		DisableMlock: true,
	}
	core2, err := vault.NewCore(conf2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer core2.Shutdown()
	for _, key := range keys {
		if _, err := core2.Unseal(vault.TestKeyCopy(key)); err != nil {
			t.Fatalf("unseal err: %s", err)
		}
	}

	TestServerWithListener(t, ln1, addr1, core1)
	TestServerWithListener(t, ln2, addr2, core2)
	TestServerAuth(t, addr1, root)

	// WRITE to STANDBY
	resp := testHttpPutDisableRedirect(t, root, addr2+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
	})
	logger.Debug("307 test one starting")
	testResponseStatus(t, resp, 307)
	logger.Debug("307 test one stopping")

	//// READ to standby
	resp = testHttpGet(t, root, addr2+"/v1/auth/token/lookup-self")
	var actual map[string]interface{}
	var nilWarnings interface{}
	expected := map[string]interface{}{
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data": map[string]interface{}{
			"meta":             nil,
			"num_uses":         json.Number("0"),
			"path":             "auth/token/root",
			"policies":         []interface{}{"root"},
			"display_name":     "root",
			"orphan":           true,
			"id":               root,
			"ttl":              json.Number("0"),
			"creation_ttl":     json.Number("0"),
			"explicit_max_ttl": json.Number("0"),
			"expire_time":      nil,
			"entity_id":        "",
			"type":             "service",
		},
		"warnings":  nilWarnings,
		"wrap_info": nil,
		"auth":      nil,
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	actualDataMap := actual["data"].(map[string]interface{})
	delete(actualDataMap, "creation_time")
	delete(actualDataMap, "accessor")
	actual["data"] = actualDataMap
	expected["request_id"] = actual["request_id"]
	delete(actual, "lease_id")
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}

	//// DELETE to standby
	resp = testHttpDeleteDisableRedirect(t, root, addr2+"/v1/secret/foo")
	logger.Debug("307 test two starting")
	testResponseStatus(t, resp, 307)
	logger.Debug("307 test two stopping")
}

func TestLogical_CreateToken(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// WRITE
	resp := testHttpPut(t, token, addr+"/v1/auth/token/create", map[string]interface{}{
		"data": "bar",
	})

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"data":           nil,
		"wrap_info":      nil,
		"auth": map[string]interface{}{
			"policies":        []interface{}{"root"},
			"token_policies":  []interface{}{"root"},
			"metadata":        nil,
			"lease_duration":  json.Number("0"),
			"renewable":       false,
			"entity_id":       "",
			"token_type":      "service",
			"orphan":          false,
			"mfa_requirement": nil,
			"num_uses":        json.Number("0"),
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	delete(actual["auth"].(map[string]interface{}), "client_token")
	delete(actual["auth"].(map[string]interface{}), "accessor")
	delete(actual, "warnings")
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nexpected:\n%#v\nactual:\n%#v", expected, actual)
	}
}

func TestLogical_RawHTTP(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type": "http",
	})
	testResponseStatus(t, resp, 204)

	// Get the raw response
	resp = testHttpGet(t, token, addr+"/v1/foo/raw")
	testResponseStatus(t, resp, 200)

	// Test the headers
	if resp.Header.Get("Content-Type") != "plain/text" {
		t.Fatalf("Bad: %#v", resp.Header)
	}

	// Get the body
	body := new(bytes.Buffer)
	io.Copy(body, resp.Body)
	if string(body.Bytes()) != "hello world" {
		t.Fatalf("Bad: %s", body.Bytes())
	}
}

func TestLogical_RequestSizeLimit(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Write a very large object, should fail. This test works because Go will
	// convert the byte slice to base64, which makes it significantly larger
	// than the default max request size.
	resp := testHttpPut(t, token, addr+"/v1/secret/foo", map[string]interface{}{
		"data": make([]byte, DefaultMaxRequestSize),
	})
	testResponseStatus(t, resp, http.StatusRequestEntityTooLarge)
}

func TestLogical_RequestSizeDisableLimit(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestListener(t)
	props := &vault.HandlerProperties{
		Core: core,
		ListenerConfig: &configutil.Listener{
			MaxRequestSize: -1,
			Address:        "127.0.0.1",
			TLSDisable:     true,
		},
	}
	TestServerWithListenerAndProperties(t, ln, addr, core, props)

	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Write a very large object, should pass as MaxRequestSize set to -1/Negative value

	resp := testHttpPut(t, token, addr+"/v1/secret/foo", map[string]interface{}{
		"data": make([]byte, DefaultMaxRequestSize),
	})
	testResponseStatus(t, resp, http.StatusNoContent)
}

func TestLogical_ListSuffix(t *testing.T) {
	core, _, rootToken := vault.TestCoreUnsealed(t)
	req, _ := http.NewRequest("GET", "http://127.0.0.1:8200/v1/secret/foo", nil)
	req = req.WithContext(namespace.RootContext(nil))
	req.Header.Add(consts.AuthHeaderName, rootToken)

	lreq, _, status, err := buildLogicalRequest(core, nil, req)
	if err != nil {
		t.Fatal(err)
	}
	if status != 0 {
		t.Fatalf("got status %d", status)
	}
	if strings.HasSuffix(lreq.Path, "/") {
		t.Fatal("trailing slash found on path")
	}

	req, _ = http.NewRequest("GET", "http://127.0.0.1:8200/v1/secret/foo?list=true", nil)
	req = req.WithContext(namespace.RootContext(nil))
	req.Header.Add(consts.AuthHeaderName, rootToken)

	lreq, _, status, err = buildLogicalRequest(core, nil, req)
	if err != nil {
		t.Fatal(err)
	}
	if status != 0 {
		t.Fatalf("got status %d", status)
	}
	if !strings.HasSuffix(lreq.Path, "/") {
		t.Fatal("trailing slash not found on path")
	}

	req, _ = http.NewRequest("LIST", "http://127.0.0.1:8200/v1/secret/foo", nil)
	req = req.WithContext(namespace.RootContext(nil))
	req.Header.Add(consts.AuthHeaderName, rootToken)

	_, _, status, err = buildLogicalRequestNoAuth(core.PerfStandby(), nil, req)
	if err != nil || status != 0 {
		t.Fatal(err)
	}

	lreq, _, status, err = buildLogicalRequest(core, nil, req)
	if err != nil {
		t.Fatal(err)
	}
	if status != 0 {
		t.Fatalf("got status %d", status)
	}
	if !strings.HasSuffix(lreq.Path, "/") {
		t.Fatal("trailing slash not found on path")
	}
}

func TestLogical_ListWithQueryParameters(t *testing.T) {
	core, _, rootToken := vault.TestCoreUnsealed(t)

	tests := []struct {
		name          string
		requestMethod string
		url           string
		expectedData  map[string]interface{}
	}{
		{
			name:          "LIST request method parses query parameter",
			requestMethod: "LIST",
			url:           "http://127.0.0.1:8200/v1/secret/foo?key1=value1",
			expectedData: map[string]interface{}{
				"key1": "value1",
			},
		},
		{
			name:          "LIST request method parses query multiple parameters",
			requestMethod: "LIST",
			url:           "http://127.0.0.1:8200/v1/secret/foo?key1=value1&key2=value2",
			expectedData: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:          "GET request method with list=true parses query parameter",
			requestMethod: "GET",
			url:           "http://127.0.0.1:8200/v1/secret/foo?list=true&key1=value1",
			expectedData: map[string]interface{}{
				"key1": "value1",
			},
		},
		{
			name:          "GET request method with list=true parses multiple query parameters",
			requestMethod: "GET",
			url:           "http://127.0.0.1:8200/v1/secret/foo?list=true&key1=value1&key2=value2",
			expectedData: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:          "GET request method with alternate order list=true parses multiple query parameters",
			requestMethod: "GET",
			url:           "http://127.0.0.1:8200/v1/secret/foo?key1=value1&list=true&key2=value2",
			expectedData: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(tc.requestMethod, tc.url, nil)
			req = req.WithContext(namespace.RootContext(nil))
			req.Header.Add(consts.AuthHeaderName, rootToken)

			lreq, _, status, err := buildLogicalRequest(core, nil, req)
			if err != nil {
				t.Fatal(err)
			}
			if status != 0 {
				t.Fatalf("got status %d", status)
			}
			if !strings.HasSuffix(lreq.Path, "/") {
				t.Fatal("trailing slash not found on path")
			}
			if lreq.Operation != logical.ListOperation {
				t.Fatalf("expected logical.ListOperation, got %v", lreq.Operation)
			}
			if !reflect.DeepEqual(tc.expectedData, lreq.Data) {
				t.Fatalf("expected query parameter data %v, got %v", tc.expectedData, lreq.Data)
			}
		})
	}
}

func TestLogical_RespondWithStatusCode(t *testing.T) {
	resp := &logical.Response{
		Data: map[string]interface{}{
			"test-data": "foo",
		},
	}

	resp404, err := logical.RespondWithStatusCode(resp, &logical.Request{ID: "id"}, http.StatusNotFound)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	respondLogical(nil, w, nil, nil, resp404, false)

	if w.Code != 404 {
		t.Fatalf("Bad Status code: %d", w.Code)
	}

	bodyRaw, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"request_id":"id","lease_id":"","renewable":false,"lease_duration":0,"data":{"test-data":"foo"},"wrap_info":null,"warnings":null,"auth":null}`

	if string(bodyRaw[:]) != strings.Trim(expected, "\n") {
		t.Fatalf("bad response: %s", string(bodyRaw[:]))
	}
}

func TestLogical_Audit_invalidWrappingToken(t *testing.T) {
	t.Setenv("VAULT_AUDIT_DISABLE_EVENTLOGGER", "true")

	// Create a noop audit backend
	noop := corehelpers.TestNoopAudit(t, nil)
	c, _, root := vault.TestCoreUnsealedWithConfig(t, &vault.CoreConfig{
		AuditBackends: map[string]audit.Factory{
			"noop": func(ctx context.Context, config *audit.BackendConfig, _ bool, _ audit.HeaderFormatter) (audit.Backend, error) {
				return noop, nil
			},
		},
	})
	ln, addr := TestServer(t, c)
	defer ln.Close()

	// Enable the audit backend

	resp := testHttpPost(t, root, addr+"/v1/sys/audit/noop", map[string]interface{}{
		"type": "noop",
	})
	testResponseStatus(t, resp, 204)

	{
		// Make a wrapping/unwrap request with an invalid token
		resp := testHttpPost(t, root, addr+"/v1/sys/wrapping/unwrap", map[string]interface{}{
			"token": "foo",
		})
		testResponseStatus(t, resp, 400)
		body := map[string][]string{}
		testResponseBody(t, resp, &body)
		if body["errors"][0] != "wrapping token is not valid or does not exist" {
			t.Fatal(body)
		}

		// Check the audit trail on request and response
		if len(noop.ReqAuth) != 1 {
			t.Fatalf("bad: %#v", noop)
		}
		auth := noop.ReqAuth[0]
		if auth.ClientToken != root {
			t.Fatalf("bad client token: %#v", auth)
		}
		if len(noop.Req) != 1 || noop.Req[0].Path != "sys/wrapping/unwrap" {
			t.Fatalf("bad:\ngot:\n%#v", noop.Req[0])
		}

		if len(noop.ReqErrs) != 1 {
			t.Fatalf("bad: %#v", noop.RespErrs)
		}
		if noop.ReqErrs[0] != consts.ErrInvalidWrappingToken {
			t.Fatalf("bad: %#v", noop.ReqErrs)
		}
	}

	{
		resp := testHttpPostWrapped(t, root, addr+"/v1/auth/token/create", nil, 10*time.Second)
		testResponseStatus(t, resp, 200)
		body := map[string]interface{}{}
		testResponseBody(t, resp, &body)

		wrapToken := body["wrap_info"].(map[string]interface{})["token"].(string)

		// Make a wrapping/unwrap request with an invalid token
		resp = testHttpPost(t, root, addr+"/v1/sys/wrapping/unwrap", map[string]interface{}{
			"token": wrapToken,
		})
		testResponseStatus(t, resp, 200)

		// Check the audit trail on request and response
		if len(noop.ReqAuth) != 3 {
			t.Fatalf("bad: %#v", noop)
		}
		auth := noop.ReqAuth[2]
		if auth.ClientToken != root {
			t.Fatalf("bad client token: %#v", auth)
		}
		if len(noop.Req) != 3 || noop.Req[2].Path != "sys/wrapping/unwrap" {
			t.Fatalf("bad:\ngot:\n%#v", noop.Req[2])
		}

		// Make sure there is only one error in the logs
		if noop.ReqErrs[1] != nil || noop.ReqErrs[2] != nil {
			t.Fatalf("bad: %#v", noop.RespErrs)
		}
	}
}

func TestLogical_ShouldParseForm(t *testing.T) {
	const formCT = "application/x-www-form-urlencoded"

	tests := map[string]struct {
		prefix      string
		contentType string
		isForm      bool
	}{
		"JSON":                 {`{"a":42}`, formCT, false},
		"JSON 2":               {`[42]`, formCT, false},
		"JSON w/leading space": {"   \n\n\r\t  [42]  ", formCT, false},
		"Form":                 {"a=42&b=dog", formCT, true},
		"Form w/wrong CT":      {"a=42&b=dog", "application/json", false},
	}

	for name, test := range tests {
		isForm := isForm([]byte(test.prefix), test.contentType)

		if isForm != test.isForm {
			t.Fatalf("%s fail: expected isForm %t, got %t", name, test.isForm, isForm)
		}
	}
}

func TestLogical_AuditPort(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.VersionedKVFactory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": auditFile.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	core := cores[0].Core
	c := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	if err := c.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
	}

	auditLogFile, err := ioutil.TempFile("", "auditport")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": auditLogFile.Name(),
		},
	})
	if err != nil {
		t.Fatalf("failed to enable audit file, err: %#v\n", err)
	}

	writeData := map[string]interface{}{
		"data": map[string]interface{}{
			"bar": "a",
		},
	}

	// workaround kv-v2 initialization upgrade errors
	numFailures := 0
	corehelpers.RetryUntil(t, 10*time.Second, func() error {
		resp, err := c.Logical().Write("kv/data/foo", writeData)
		if err != nil {
			if strings.Contains(err.Error(), "Upgrading from non-versioned to versioned data") {
				t.Logf("Retrying fetch KV data due to upgrade error")
				time.Sleep(100 * time.Millisecond)
				numFailures += 1
				return err
			}

			t.Fatalf("write request failed, err: %#v, resp: %#v\n", err, resp)
		}

		return nil
	})

	decoder := json.NewDecoder(auditLogFile)

	var auditRecord map[string]interface{}
	count := 0
	for decoder.Decode(&auditRecord) == nil {
		count += 1

		// Skip the first line
		if count == 1 {
			continue
		}

		auditRequest := map[string]interface{}{}

		if req, ok := auditRecord["request"]; ok {
			auditRequest = req.(map[string]interface{})
		}

		if _, ok := auditRequest["remote_address"].(string); !ok {
			t.Fatalf("remote_address should be a string, not %T", auditRequest["remote_address"])
		}

		if _, ok := auditRequest["remote_port"].(float64); !ok {
			t.Fatalf("remote_port should be a number, not %T", auditRequest["remote_port"])
		}
	}

	// We expect the following items in the audit log:
	// audit log header + an entry for updating sys/audit/file
	// + request/response per failure (if any) + request/response for creating kv
	numExpectedEntries := (numFailures * 2) + 4
	if count != numExpectedEntries {
		t.Fatalf("wrong number of audit entries expected: %d got: %d", numExpectedEntries, count)
	}
}

func TestLogical_ErrRelativePath(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	core := cores[0].Core
	c := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	err := c.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	})
	if err != nil {
		t.Fatalf("failed to enable userpass, err: %v", err)
	}

	resp, err := c.Logical().Read("auth/userpass/users/user..aaa")

	if err == nil || resp != nil {
		t.Fatalf("expected read request to fail, resp: %#v, err: %v", resp, err)
	}

	respErr, ok := err.(*api.ResponseError)

	if !ok {
		t.Fatalf("unexpected error type, err: %#v", err)
	}

	if respErr.StatusCode != 400 {
		t.Errorf("expected 400 response for read, actual: %d", respErr.StatusCode)
	}

	if !strings.Contains(respErr.Error(), logical.ErrRelativePath.Error()) {
		t.Errorf("expected response for read to include %q", logical.ErrRelativePath.Error())
	}

	data := map[string]interface{}{
		"password": "abc123",
	}

	resp, err = c.Logical().Write("auth/userpass/users/user..aaa", data)

	if err == nil || resp != nil {
		t.Fatalf("expected write request to fail, resp: %#v, err: %v", resp, err)
	}

	respErr, ok = err.(*api.ResponseError)

	if !ok {
		t.Fatalf("unexpected error type, err: %#v", err)
	}

	if respErr.StatusCode != 400 {
		t.Errorf("expected 400 response for write, actual: %d", respErr.StatusCode)
	}

	if !strings.Contains(respErr.Error(), logical.ErrRelativePath.Error()) {
		t.Errorf("expected response for write to include %q", logical.ErrRelativePath.Error())
	}
}

func testBuiltinPluginMetadataAuditLog(t *testing.T, log map[string]interface{}, expectedMountClass string) {
	if mountClass, ok := log["mount_class"].(string); !ok {
		t.Fatalf("mount_class should be a string, not %T", log["mount_class"])
	} else if mountClass != expectedMountClass {
		t.Fatalf("bad: mount_class should be %s, not %s", expectedMountClass, mountClass)
	}

	if _, ok := log["mount_running_version"].(string); !ok {
		t.Fatalf("mount_running_version should be a string, not %T", log["mount_running_version"])
	}

	if _, ok := log["mount_running_sha256"].(string); ok {
		t.Fatalf("mount_running_sha256 should be nil, not %T", log["mount_running_sha256"])
	}

	if mountIsExternalPlugin, ok := log["mount_is_external_plugin"].(bool); ok && mountIsExternalPlugin {
		t.Fatalf("mount_is_external_plugin should be nil or false, not %T", log["mount_is_external_plugin"])
	}
}

// TestLogical_AuditEnabled_ShouldLogPluginMetadata_Auth tests that we have plugin metadata of a builtin auth plugin
// in audit log when it is enabled
func TestLogical_AuditEnabled_ShouldLogPluginMetadata_Auth(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		AuditBackends: map[string]audit.Factory{
			"file": auditFile.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	core := cores[0].Core
	c := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	// Enable the audit backend
	tempDir := t.TempDir()
	auditLogFile, err := os.CreateTemp(tempDir, "")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": auditLogFile.Name(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = c.Logical().Write("auth/token/create", map[string]interface{}{
		"ttl": "10s",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Check the audit trail on request and response
	decoder := json.NewDecoder(auditLogFile)
	var auditRecord map[string]interface{}
	for decoder.Decode(&auditRecord) == nil {
		auditRequest := map[string]interface{}{}
		if req, ok := auditRecord["request"]; ok {
			auditRequest = req.(map[string]interface{})
			if auditRequest["path"] != "auth/token/create" {
				continue
			}
		}
		testBuiltinPluginMetadataAuditLog(t, auditRequest, consts.PluginTypeCredential.String())

		auditResponse := map[string]interface{}{}
		if req, ok := auditRecord["response"]; ok {
			auditRequest = req.(map[string]interface{})
			if auditResponse["path"] != "auth/token/create" {
				continue
			}
		}
		testBuiltinPluginMetadataAuditLog(t, auditResponse, consts.PluginTypeCredential.String())
	}
}

// TestLogical_AuditEnabled_ShouldLogPluginMetadata_Secret tests that we have plugin metadata of a builtin secret plugin
// in audit log when it is enabled
func TestLogical_AuditEnabled_ShouldLogPluginMetadata_Secret(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.VersionedKVFactory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": auditFile.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	core := cores[0].Core
	c := cluster.Cores[0].Client
	vault.TestWaitActive(t, core)

	if err := c.Sys().Mount("kv/", &api.MountInput{
		Type: "kv-v2",
	}); err != nil {
		t.Fatalf("kv-v2 mount attempt failed - err: %#v\n", err)
	}

	// Enable the audit backend
	tempDir := t.TempDir()
	auditLogFile, err := os.CreateTemp(tempDir, "")
	if err != nil {
		t.Fatal(err)
	}

	err = c.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": auditLogFile.Name(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	{
		writeData := map[string]interface{}{
			"data": map[string]interface{}{
				"bar": "a",
			},
		}
		corehelpers.RetryUntil(t, 10*time.Second, func() error {
			resp, err := c.Logical().Write("kv/data/foo", writeData)
			if err != nil {
				t.Fatalf("write request failed, err: %#v, resp: %#v\n", err, resp)
			}
			return nil
		})
	}

	// Check the audit trail on request and response
	decoder := json.NewDecoder(auditLogFile)
	var auditRecord map[string]interface{}
	for decoder.Decode(&auditRecord) == nil {
		auditRequest := map[string]interface{}{}
		if req, ok := auditRecord["request"]; ok {
			auditRequest = req.(map[string]interface{})
			if auditRequest["path"] != "kv/data/foo" {
				continue
			}
		}
		testBuiltinPluginMetadataAuditLog(t, auditRequest, consts.PluginTypeSecrets.String())

		auditResponse := map[string]interface{}{}
		if req, ok := auditRecord["response"]; ok {
			auditRequest = req.(map[string]interface{})
			if auditResponse["path"] != "kv/data/foo" {
				continue
			}
		}
		testBuiltinPluginMetadataAuditLog(t, auditResponse, consts.PluginTypeSecrets.String())
	}
}
