package http

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysCapabilitiesAccessor(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Lookup the token properties
	resp := testHttpGet(t, token, addr+"/v1/auth/token/lookup/"+token)
	var lookupResp map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &lookupResp)

	// Retrieve the accessor from the token properties
	lookupData := lookupResp["data"].(map[string]interface{})
	accessor := lookupData["accessor"].(string)

	resp = testHttpPost(t, token, addr+"/v1/sys/capabilities-accessor", map[string]interface{}{
		"accessor": accessor,
		"path":     "testpath",
	})

	var actual map[string][]string
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected := map[string][]string{
		"capabilities": []string{"root"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	// Testing for non-root token's accessor
	// Create a policy first
	resp = testHttpPost(t, token, addr+"/v1/sys/policy/foo", map[string]interface{}{
		"rules": `path "testpath" {capabilities = ["read","sudo"]}`,
	})
	testResponseStatus(t, resp, 204)

	// Create a token against the test policy
	resp = testHttpPost(t, token, addr+"/v1/auth/token/create", map[string]interface{}{
		"policies": []string{"foo"},
	})

	var tokenResp map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &tokenResp)

	// Check if desired policies are present in the token
	auth := tokenResp["auth"].(map[string]interface{})
	actualPolicies := auth["policies"]
	expectedPolicies := []interface{}{"default", "foo"}
	if !reflect.DeepEqual(actualPolicies, expectedPolicies) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actualPolicies, expectedPolicies)
	}

	// Check the capabilities of non-root token using the accessor
	resp = testHttpPost(t, token, addr+"/v1/sys/capabilities-accessor", map[string]interface{}{
		"accessor": auth["accessor"],
		"path":     "testpath",
	})
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected = map[string][]string{
		"capabilities": []string{"sudo", "read"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}

func TestSysCapabilities(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Send both token and path
	resp := testHttpPost(t, token, addr+"/v1/sys/capabilities", map[string]interface{}{
		"token": token,
		"path":  "testpath",
	})

	var actual map[string][]string
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected := map[string][]string{
		"capabilities": []string{"root"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	// Send only path to capabilities-self
	resp = testHttpPost(t, token, addr+"/v1/sys/capabilities-self", map[string]interface{}{
		"path": "testpath",
	})
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}

	// Testing for non-root tokens

	// Create a policy first
	resp = testHttpPost(t, token, addr+"/v1/sys/policy/foo", map[string]interface{}{
		"rules": `path "testpath" {capabilities = ["read","sudo"]}`,
	})
	testResponseStatus(t, resp, 204)

	// Create a token against the test policy
	resp = testHttpPost(t, token, addr+"/v1/auth/token/create", map[string]interface{}{
		"policies": []string{"foo"},
	})

	var tokenResp map[string]interface{}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &tokenResp)

	// Check if desired policies are present in the token
	auth := tokenResp["auth"].(map[string]interface{})
	actualPolicies := auth["policies"]
	expectedPolicies := []interface{}{"default", "foo"}
	if !reflect.DeepEqual(actualPolicies, expectedPolicies) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actualPolicies, expectedPolicies)
	}

	// Check the capabilities with the created non-root token
	resp = testHttpPost(t, token, addr+"/v1/sys/capabilities", map[string]interface{}{
		"token": auth["client_token"],
		"path":  "testpath",
	})
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)

	expected = map[string][]string{
		"capabilities": []string{"sudo", "read"},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: got\n%#v\nexpected\n%#v\n", actual, expected)
	}
}
