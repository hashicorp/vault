package http

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault"
)

func TestSysAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, token, addr+"/v1/sys/auth")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"token/": map[string]interface{}{
				"description":             "token based credentials",
				"type":                    "token",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"token_type":        "default-service",
					"force_no_cache":    false,
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                interface{}(nil),
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeCredential, "token"),
			},
		},
		"token/": map[string]interface{}{
			"description":             "token based credentials",
			"type":                    "token",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"token_type":        "default-service",
				"force_no_cache":    false,
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                interface{}(nil),
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeCredential, "token"),
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

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected:%#v\nactual:%#v", expected, actual)
	}
}

func TestSysEnableAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/auth/foo", map[string]interface{}{
		"type":        "approle",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/auth")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"foo/": map[string]interface{}{
				"description":             "foo",
				"type":                    "approle",
				"external_entropy_access": false,
				"deprecation_status":      "supported",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"token_type":        "default-service",
					"force_no_cache":    false,
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                map[string]interface{}{},
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeCredential, "approle"),
			},
			"token/": map[string]interface{}{
				"description":             "token based credentials",
				"type":                    "token",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"token_type":        "default-service",
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                interface{}(nil),
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeCredential, "token"),
			},
		},
		"foo/": map[string]interface{}{
			"description":             "foo",
			"type":                    "approle",
			"external_entropy_access": false,
			"deprecation_status":      "supported",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"token_type":        "default-service",
				"force_no_cache":    false,
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                map[string]interface{}{},
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeCredential, "approle"),
		},
		"token/": map[string]interface{}{
			"description":             "token based credentials",
			"type":                    "token",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"token_type":        "default-service",
				"force_no_cache":    false,
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                interface{}(nil),
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeCredential, "token"),
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

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestSysDisableAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/auth/foo", map[string]interface{}{
		"type":        "noop",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpDelete(t, token, addr+"/v1/sys/auth/foo")
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/auth")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"token/": map[string]interface{}{
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"token_type":        "default-service",
					"force_no_cache":    false,
				},
				"description":             "token based credentials",
				"type":                    "token",
				"external_entropy_access": false,
				"local":                   false,
				"seal_wrap":               false,
				"options":                 interface{}(nil),
				"plugin_version":          "",
				"running_sha256":          "",
				"running_plugin_version":  versions.GetBuiltinVersion(consts.PluginTypeCredential, "token"),
			},
		},
		"token/": map[string]interface{}{
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"token_type":        "default-service",
				"force_no_cache":    false,
			},
			"description":             "token based credentials",
			"type":                    "token",
			"external_entropy_access": false,
			"local":                   false,
			"seal_wrap":               false,
			"options":                 interface{}(nil),
			"plugin_version":          "",
			"running_sha256":          "",
			"running_plugin_version":  versions.GetBuiltinVersion(consts.PluginTypeCredential, "token"),
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

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}

func TestSysTuneAuth_nonHMACKeys(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Mount-tune the audit_non_hmac_request_keys
	resp := testHttpPost(t, token, addr+"/v1/sys/auth/token/tune", map[string]interface{}{
		"audit_non_hmac_request_keys": "foo",
	})
	testResponseStatus(t, resp, 204)

	// Mount-tune the audit_non_hmac_response_keys
	resp = testHttpPost(t, token, addr+"/v1/sys/auth/token/tune", map[string]interface{}{
		"audit_non_hmac_response_keys": "bar",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/auth/token/tune")
	testResponseStatus(t, resp, 200)

	actual := map[string]interface{}{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"description":                  "token based credentials",
			"default_lease_ttl":            json.Number("2764800"),
			"max_lease_ttl":                json.Number("2764800"),
			"force_no_cache":               false,
			"audit_non_hmac_request_keys":  []interface{}{"foo"},
			"audit_non_hmac_response_keys": []interface{}{"bar"},
			"token_type":                   "default-service",
		},
		"description":                  "token based credentials",
		"default_lease_ttl":            json.Number("2764800"),
		"max_lease_ttl":                json.Number("2764800"),
		"force_no_cache":               false,
		"audit_non_hmac_request_keys":  []interface{}{"foo"},
		"audit_non_hmac_response_keys": []interface{}{"bar"},
		"token_type":                   "default-service",
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Unset those mount tune values
	resp = testHttpPost(t, token, addr+"/v1/sys/auth/token/tune", map[string]interface{}{
		"audit_non_hmac_request_keys": "",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/auth/token/tune", map[string]interface{}{
		"audit_non_hmac_response_keys": "",
	})

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/auth/token/tune")
	testResponseStatus(t, resp, 200)

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"description":       "token based credentials",
			"default_lease_ttl": json.Number("2764800"),
			"max_lease_ttl":     json.Number("2764800"),
			"force_no_cache":    false,
			"token_type":        "default-service",
		},
		"description":       "token based credentials",
		"default_lease_ttl": json.Number("2764800"),
		"max_lease_ttl":     json.Number("2764800"),
		"force_no_cache":    false,
		"token_type":        "default-service",
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}
}

func TestSysTuneAuth_showUIMount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Get original tune values, ensure that listing_visibility is not set
	resp := testHttpGet(t, token, addr+"/v1/sys/auth/token/tune")
	testResponseStatus(t, resp, 200)

	actual := map[string]interface{}{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"description":       "token based credentials",
			"default_lease_ttl": json.Number("2764800"),
			"max_lease_ttl":     json.Number("2764800"),
			"force_no_cache":    false,
			"token_type":        "default-service",
		},
		"description":       "token based credentials",
		"default_lease_ttl": json.Number("2764800"),
		"max_lease_ttl":     json.Number("2764800"),
		"force_no_cache":    false,
		"token_type":        "default-service",
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Mount-tune the listing_visibility
	resp = testHttpPost(t, token, addr+"/v1/sys/auth/token/tune", map[string]interface{}{
		"listing_visibility": "unauth",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/auth/token/tune")
	testResponseStatus(t, resp, 200)

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"description":    "token based credentials",
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"description":        "token based credentials",
			"default_lease_ttl":  json.Number("2764800"),
			"max_lease_ttl":      json.Number("2764800"),
			"force_no_cache":     false,
			"listing_visibility": "unauth",
			"token_type":         "default-service",
		},
		"default_lease_ttl":  json.Number("2764800"),
		"max_lease_ttl":      json.Number("2764800"),
		"force_no_cache":     false,
		"listing_visibility": "unauth",
		"token_type":         "default-service",
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}
}

func TestSysRemountAuth(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/auth/foo", map[string]interface{}{
		"type":        "noop",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/remount", map[string]interface{}{
		"from": "auth/foo",
		"to":   "auth/bar",
	})
	testResponseStatus(t, resp, 200)

	// Poll until the remount succeeds
	var remountResp map[string]interface{}
	testResponseBody(t, resp, &remountResp)
	corehelpers.RetryUntil(t, 5*time.Second, func() error {
		resp = testHttpGet(t, token, addr+"/v1/sys/remount/status/"+remountResp["migration_id"].(string))
		testResponseStatus(t, resp, 200)

		var remountStatusResp map[string]interface{}
		testResponseBody(t, resp, &remountStatusResp)

		status := remountStatusResp["data"].(map[string]interface{})["migration_info"].(map[string]interface{})["status"]
		if status != "success" {
			return fmt.Errorf("Expected migration status to be successful, got %q", status)
		}
		return nil
	})

	resp = testHttpGet(t, token, addr+"/v1/sys/auth")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"bar/": map[string]interface{}{
				"description":             "foo",
				"type":                    "noop",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"token_type":        "default-service",
					"force_no_cache":    false,
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                map[string]interface{}{},
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
			},
			"token/": map[string]interface{}{
				"description":             "token based credentials",
				"type":                    "token",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"token_type":        "default-service",
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                interface{}(nil),
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeCredential, "token"),
			},
		},
		"bar/": map[string]interface{}{
			"description":             "foo",
			"type":                    "noop",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"token_type":        "default-service",
				"force_no_cache":    false,
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                map[string]interface{}{},
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
		},
		"token/": map[string]interface{}{
			"description":             "token based credentials",
			"type":                    "token",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"token_type":        "default-service",
				"force_no_cache":    false,
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                interface{}(nil),
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeCredential, "token"),
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

	if diff := deep.Equal(actual, expected); diff != nil {
		t.Fatal(diff)
	}
}
