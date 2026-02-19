// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package http

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/fatih/structs"
	"github.com/go-test/deep"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/helper/versions"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestSysMounts(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
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

func TestSysMount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "kv",
		"description": "foo",
		"options": map[string]string{
			"version": "1",
		},
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"foo/": map[string]interface{}{
				"description":             "foo",
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
		"foo/": map[string]interface{}{
			"description":             "foo",
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

func TestSysMount_put(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "kv",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	// The TestSysMount test tests the thing is actually created. See that test
	// for more info.
}

// TestSysRemountSpacesFrom ensure we succeed in a remount where the 'from' mount has spaces in the name
func TestSysRemountSpacesFrom(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo%20bar", map[string]interface{}{
		"type":        "kv",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/remount", map[string]interface{}{
		"from": "foo bar",
		"to":   "baz",
	})
	testResponseStatus(t, resp, 200)
}

// TestSysRemountSpacesTo ensure we succeed in a remount where the 'to' mount has spaces in the name
func TestSysRemountSpacesTo(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo%20bar", map[string]interface{}{
		"type":        "kv",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/remount", map[string]interface{}{
		"from": "foo bar",
		"to":   "bar baz",
	})
	testResponseStatus(t, resp, 200)
}

// TestSysRemountTrailingSpaces ensures we fail on trailing spaces
func TestSysRemountTrailingSpaces(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo%20bar", map[string]interface{}{
		"type":        "kv",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/remount", map[string]interface{}{
		"from": "foo bar",
		"to":   " baz ",
	})
	testResponseStatus(t, resp, 400)

	resp = testHttpPost(t, token, addr+"/v1/sys/remount", map[string]interface{}{
		"from": " foo bar ",
		"to":   "baz",
	})
	testResponseStatus(t, resp, 400)
}

func TestSysRemount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "kv",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/remount", map[string]interface{}{
		"from": "foo",
		"to":   "bar",
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
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"bar/": map[string]interface{}{
				"description":             "foo",
				"type":                    "kv",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                map[string]interface{}{},
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
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
		"bar/": map[string]interface{}{
			"description":             "foo",
			"type":                    "kv",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                map[string]interface{}{},
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
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
		t.Fatalf("bad:\nExpected: %#v\nActual: %#v\n", expected, actual)
	}
}

func TestSysUnmount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "kv",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpDelete(t, token, addr+"/v1/sys/mounts/foo")
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
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

func TestSysTuneMount_Options(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "kv",
		"description": "foo",
	})

	testResponseStatus(t, resp, 204)
	// Mount-tune the options
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"options": map[string]string{
			"test": "true",
		},
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/foo/tune")
	testResponseStatus(t, resp, 200)

	actual := map[string]interface{}{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":       "foo",
			"default_lease_ttl": json.Number("2764800"),
			"max_lease_ttl":     json.Number("2764800"),
			"force_no_cache":    false,
			"options":           map[string]interface{}{"test": "true"},
		},
		"description":       "foo",
		"default_lease_ttl": json.Number("2764800"),
		"max_lease_ttl":     json.Number("2764800"),
		"force_no_cache":    false,
		"options":           map[string]interface{}{"test": "true"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Check that we're not allowed to unset the options map once that's set
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"options": map[string]string{},
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/foo/tune")
	testResponseStatus(t, resp, 200)

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"description":    "foo",
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":       "foo",
			"default_lease_ttl": json.Number("2764800"),
			"max_lease_ttl":     json.Number("2764800"),
			"force_no_cache":    false,
			"options":           map[string]interface{}{"test": "true"},
		},
		"default_lease_ttl": json.Number("2764800"),
		"max_lease_ttl":     json.Number("2764800"),
		"force_no_cache":    false,
		"options":           map[string]interface{}{"test": "true"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}
}

func TestSysTuneMount(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": vault.LeasedPassthroughBackendFactory,
		},
	}
	core, _, token := vault.TestCoreUnsealedWithConfig(t, coreConfig)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "kv",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"foo/": map[string]interface{}{
				"description":             "foo",
				"type":                    "kv",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                map[string]interface{}{},
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
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
		"foo/": map[string]interface{}{
			"description":             "foo",
			"type":                    "kv",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                map[string]interface{}{},
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
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

	// Shorter than system default
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"default_lease_ttl": "72h",
	})
	testResponseStatus(t, resp, 204)

	// Longer than system max
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"default_lease_ttl": "72000h",
	})
	testResponseStatus(t, resp, 204)

	// Longer than system default
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"max_lease_ttl": "72000h",
	})
	testResponseStatus(t, resp, 204)

	// Longer than backend max
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"default_lease_ttl": "72001h",
	})
	testResponseStatus(t, resp, 400)

	// Shorter than backend default
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"max_lease_ttl": "1h",
	})
	testResponseStatus(t, resp, 400)

	// Shorter than backend max, longer than system max
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"default_lease_ttl": "71999h",
	})
	testResponseStatus(t, resp, 204)

	// mark as versioned
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/foo/tune", map[string]interface{}{
		"options": map[string]string{
			"version": "1",
		},
	})
	testResponseStatus(t, resp, 200)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")
	expected = map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"foo/": map[string]interface{}{
				"description":             "foo",
				"type":                    "kv",
				"external_entropy_access": false,
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("259196400"),
					"max_lease_ttl":     json.Number("259200000"),
					"force_no_cache":    false,
				},
				"local":                  false,
				"seal_wrap":              false,
				"options":                map[string]interface{}{"version": "1"},
				"plugin_version":         "",
				"running_sha256":         "",
				"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
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
		"foo/": map[string]interface{}{
			"description":             "foo",
			"type":                    "kv",
			"external_entropy_access": false,
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("259196400"),
				"max_lease_ttl":     json.Number("259200000"),
				"force_no_cache":    false,
			},
			"local":                  false,
			"seal_wrap":              false,
			"options":                map[string]interface{}{"version": "1"},
			"plugin_version":         "",
			"running_sha256":         "",
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
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
			"running_plugin_version": versions.GetBuiltinVersion(consts.PluginTypeSecrets, "kv"),
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

	// Check simple configuration endpoint
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/foo/tune")
	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":       "foo",
			"default_lease_ttl": json.Number("259196400"),
			"max_lease_ttl":     json.Number("259200000"),
			"force_no_cache":    false,
			"options":           map[string]interface{}{"version": "1"},
		},
		"description":       "foo",
		"default_lease_ttl": json.Number("259196400"),
		"max_lease_ttl":     json.Number("259200000"),
		"force_no_cache":    false,
		"options":           map[string]interface{}{"version": "1"},
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Set a low max
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"description":       "foobar",
		"default_lease_ttl": "40s",
		"max_lease_ttl":     "80s",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":       "foobar",
			"default_lease_ttl": json.Number("40"),
			"max_lease_ttl":     json.Number("80"),
			"force_no_cache":    false,
			"options":           map[string]interface{}{"version": "1"},
		},
		"description":       "foobar",
		"default_lease_ttl": json.Number("40"),
		"max_lease_ttl":     json.Number("80"),
		"force_no_cache":    false,
		"options":           map[string]interface{}{"version": "1"},
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// First try with lease above backend max
	resp = testHttpPut(t, token, addr+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
		"ttl":  "28347h",
	})
	testResponseStatus(t, resp, 204)

	// read secret
	resp = testHttpGet(t, token, addr+"/v1/secret/foo")
	var result struct {
		LeaseID       string `json:"lease_id" structs:"lease_id"`
		LeaseDuration int    `json:"lease_duration" structs:"lease_duration"`
	}

	testResponseBody(t, resp, &result)

	expected = map[string]interface{}{
		"lease_duration": int(80),
		"lease_id":       result.LeaseID,
	}

	if !reflect.DeepEqual(structs.Map(result), expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, structs.Map(result))
	}

	// Now with lease TTL unspecified
	resp = testHttpPut(t, token, addr+"/v1/secret/foo", map[string]interface{}{
		"data": "bar",
	})
	testResponseStatus(t, resp, 204)

	// read secret
	resp = testHttpGet(t, token, addr+"/v1/secret/foo")

	testResponseBody(t, resp, &result)

	expected = map[string]interface{}{
		"lease_duration": int(40),
		"lease_id":       result.LeaseID,
	}

	if !reflect.DeepEqual(structs.Map(result), expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, structs.Map(result))
	}
}

func TestSysTuneMount_nonHMACKeys(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Mount-tune the audit_non_hmac_request_keys
	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"audit_non_hmac_request_keys": "foo",
	})
	testResponseStatus(t, resp, 204)

	// Mount-tune the audit_non_hmac_response_keys
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"audit_non_hmac_response_keys": "bar",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	testResponseStatus(t, resp, 200)

	actual := map[string]interface{}{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":                  "key/value secret storage",
			"default_lease_ttl":            json.Number("2764800"),
			"max_lease_ttl":                json.Number("2764800"),
			"force_no_cache":               false,
			"audit_non_hmac_request_keys":  []interface{}{"foo"},
			"audit_non_hmac_response_keys": []interface{}{"bar"},
			"options":                      map[string]interface{}{"version": "1"},
		},
		"description":                  "key/value secret storage",
		"default_lease_ttl":            json.Number("2764800"),
		"max_lease_ttl":                json.Number("2764800"),
		"force_no_cache":               false,
		"audit_non_hmac_request_keys":  []interface{}{"foo"},
		"audit_non_hmac_response_keys": []interface{}{"bar"},
		"options":                      map[string]interface{}{"version": "1"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Unset those mount tune values
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"audit_non_hmac_request_keys": "",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"audit_non_hmac_response_keys": "",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	testResponseStatus(t, resp, 200)

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":       "key/value secret storage",
			"default_lease_ttl": json.Number("2764800"),
			"max_lease_ttl":     json.Number("2764800"),
			"force_no_cache":    false,
			"options":           map[string]interface{}{"version": "1"},
		},
		"description":       "key/value secret storage",
		"default_lease_ttl": json.Number("2764800"),
		"max_lease_ttl":     json.Number("2764800"),
		"force_no_cache":    false,
		"options":           map[string]interface{}{"version": "1"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}
}

func TestSysTuneMount_listingVisibility(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Get original tune values, ensure that listing_visibility is not set
	resp := testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	testResponseStatus(t, resp, 200)

	actual := map[string]interface{}{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":       "key/value secret storage",
			"default_lease_ttl": json.Number("2764800"),
			"max_lease_ttl":     json.Number("2764800"),
			"force_no_cache":    false,
			"options":           map[string]interface{}{"version": "1"},
		},
		"description":       "key/value secret storage",
		"default_lease_ttl": json.Number("2764800"),
		"max_lease_ttl":     json.Number("2764800"),
		"force_no_cache":    false,
		"options":           map[string]interface{}{"version": "1"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Mount-tune the listing_visibility
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"listing_visibility": "unauth",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	testResponseStatus(t, resp, 200)

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":        "key/value secret storage",
			"default_lease_ttl":  json.Number("2764800"),
			"max_lease_ttl":      json.Number("2764800"),
			"force_no_cache":     false,
			"listing_visibility": "unauth",
			"options":            map[string]interface{}{"version": "1"},
		},
		"description":        "key/value secret storage",
		"default_lease_ttl":  json.Number("2764800"),
		"max_lease_ttl":      json.Number("2764800"),
		"force_no_cache":     false,
		"listing_visibility": "unauth",
		"options":            map[string]interface{}{"version": "1"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}
}

func TestSysTuneMount_passthroughRequestHeaders(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Mount-tune the audit_non_hmac_request_keys
	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"passthrough_request_headers": "X-Vault-Foo",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	testResponseStatus(t, resp, 200)

	actual := map[string]interface{}{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":                 "key/value secret storage",
			"default_lease_ttl":           json.Number("2764800"),
			"max_lease_ttl":               json.Number("2764800"),
			"options":                     map[string]interface{}{"version": "1"},
			"force_no_cache":              false,
			"passthrough_request_headers": []interface{}{"X-Vault-Foo"},
		},
		"description":                 "key/value secret storage",
		"default_lease_ttl":           json.Number("2764800"),
		"max_lease_ttl":               json.Number("2764800"),
		"options":                     map[string]interface{}{"version": "1"},
		"force_no_cache":              false,
		"passthrough_request_headers": []interface{}{"X-Vault-Foo"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Unset the mount tune value
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"passthrough_request_headers": "",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	testResponseStatus(t, resp, 200)

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":       "key/value secret storage",
			"default_lease_ttl": json.Number("2764800"),
			"max_lease_ttl":     json.Number("2764800"),
			"force_no_cache":    false,
			"options":           map[string]interface{}{"version": "1"},
		},
		"description":       "key/value secret storage",
		"default_lease_ttl": json.Number("2764800"),
		"max_lease_ttl":     json.Number("2764800"),
		"force_no_cache":    false,
		"options":           map[string]interface{}{"version": "1"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}
}

func TestSysTuneMount_allowedManagedKeys(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	// Mount-tune the allowed_managed_keys
	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"allowed_managed_keys": "test_key",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	testResponseStatus(t, resp, 200)

	actual := map[string]interface{}{}
	expected := map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":          "key/value secret storage",
			"default_lease_ttl":    json.Number("2764800"),
			"max_lease_ttl":        json.Number("2764800"),
			"options":              map[string]interface{}{"version": "1"},
			"force_no_cache":       false,
			"allowed_managed_keys": []interface{}{"test_key"},
		},
		"description":          "key/value secret storage",
		"default_lease_ttl":    json.Number("2764800"),
		"max_lease_ttl":        json.Number("2764800"),
		"options":              map[string]interface{}{"version": "1"},
		"force_no_cache":       false,
		"allowed_managed_keys": []interface{}{"test_key"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Unset the mount tune value
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
		"allowed_managed_keys": "",
	})
	testResponseStatus(t, resp, 204)

	// Check results
	resp = testHttpGet(t, token, addr+"/v1/sys/mounts/secret/tune")
	testResponseStatus(t, resp, 200)

	actual = map[string]interface{}{}
	expected = map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"mount_type":     "system",
		"data": map[string]interface{}{
			"description":       "key/value secret storage",
			"default_lease_ttl": json.Number("2764800"),
			"max_lease_ttl":     json.Number("2764800"),
			"force_no_cache":    false,
			"options":           map[string]interface{}{"version": "1"},
		},
		"description":       "key/value secret storage",
		"default_lease_ttl": json.Number("2764800"),
		"max_lease_ttl":     json.Number("2764800"),
		"force_no_cache":    false,
		"options":           map[string]interface{}{"version": "1"},
	}
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}
}
