package http

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/fatih/structs"
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
		"data": map[string]interface{}{
			"secret/": map[string]interface{}{
				"description": "key/value secret storage",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"sys/": map[string]interface{}{
				"description": "system endpoints used for control, policy and debugging",
				"type":        "system",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"cubbyhole/": map[string]interface{}{
				"description": "per-token private secret storage",
				"type":        "cubbyhole",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     true,
				"seal_wrap": false,
			},
			"identity/": map[string]interface{}{
				"description": "identity store",
				"type":        "identity",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
		},
		"secret/": map[string]interface{}{
			"description": "key/value secret storage",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     true,
			"seal_wrap": false,
		},
		"identity/": map[string]interface{}{
			"description": "identity store",
			"type":        "identity",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	for k, v := range actual["data"].(map[string]interface{}) {
		if v.(map[string]interface{})["accessor"] == "" {
			t.Fatalf("no accessor from %s", k)
		}
		expected[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
		expected["data"].(map[string]interface{})[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected: %#v\nactual: %#v\n", expected, actual)
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
		"data": map[string]interface{}{
			"foo/": map[string]interface{}{
				"description": "foo",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"secret/": map[string]interface{}{
				"description": "key/value secret storage",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"sys/": map[string]interface{}{
				"description": "system endpoints used for control, policy and debugging",
				"type":        "system",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"cubbyhole/": map[string]interface{}{
				"description": "per-token private secret storage",
				"type":        "cubbyhole",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     true,
				"seal_wrap": false,
			},
			"identity/": map[string]interface{}{
				"description": "identity store",
				"type":        "identity",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
		},
		"foo/": map[string]interface{}{
			"description": "foo",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"secret/": map[string]interface{}{
			"description": "key/value secret storage",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     true,
			"seal_wrap": false,
		},
		"identity/": map[string]interface{}{
			"description": "identity store",
			"type":        "identity",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	for k, v := range actual["data"].(map[string]interface{}) {
		if v.(map[string]interface{})["accessor"] == "" {
			t.Fatalf("no accessor from %s", k)
		}
		expected[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
		expected["data"].(map[string]interface{})[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: expected: %#v\nactual: %#v\n", expected, actual)
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
		"data": map[string]interface{}{
			"bar/": map[string]interface{}{
				"description": "foo",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"secret/": map[string]interface{}{
				"description": "key/value secret storage",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"sys/": map[string]interface{}{
				"description": "system endpoints used for control, policy and debugging",
				"type":        "system",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"cubbyhole/": map[string]interface{}{
				"description": "per-token private secret storage",
				"type":        "cubbyhole",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     true,
				"seal_wrap": false,
			},
			"identity/": map[string]interface{}{
				"description": "identity store",
				"type":        "identity",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
		},
		"bar/": map[string]interface{}{
			"description": "foo",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"secret/": map[string]interface{}{
			"description": "key/value secret storage",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     true,
			"seal_wrap": false,
		},
		"identity/": map[string]interface{}{
			"description": "identity store",
			"type":        "identity",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	for k, v := range actual["data"].(map[string]interface{}) {
		if v.(map[string]interface{})["accessor"] == "" {
			t.Fatalf("no accessor from %s", k)
		}
		expected[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
		expected["data"].(map[string]interface{})[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
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
		"data": map[string]interface{}{
			"secret/": map[string]interface{}{
				"description": "key/value secret storage",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"sys/": map[string]interface{}{
				"description": "system endpoints used for control, policy and debugging",
				"type":        "system",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"cubbyhole/": map[string]interface{}{
				"description": "per-token private secret storage",
				"type":        "cubbyhole",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     true,
				"seal_wrap": false,
			},
			"identity/": map[string]interface{}{
				"description": "identity store",
				"type":        "identity",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
		},
		"secret/": map[string]interface{}{
			"description": "key/value secret storage",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     true,
			"seal_wrap": false,
		},
		"identity/": map[string]interface{}{
			"description": "identity store",
			"type":        "identity",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	for k, v := range actual["data"].(map[string]interface{}) {
		if v.(map[string]interface{})["accessor"] == "" {
			t.Fatalf("no accessor from %s", k)
		}
		expected[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
		expected["data"].(map[string]interface{})[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSysTuneMount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
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
		"data": map[string]interface{}{
			"foo/": map[string]interface{}{
				"description": "foo",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"secret/": map[string]interface{}{
				"description": "key/value secret storage",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"sys/": map[string]interface{}{
				"description": "system endpoints used for control, policy and debugging",
				"type":        "system",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"cubbyhole/": map[string]interface{}{
				"description": "per-token private secret storage",
				"type":        "cubbyhole",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     true,
				"seal_wrap": false,
			},
			"identity/": map[string]interface{}{
				"description": "identity store",
				"type":        "identity",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
		},
		"foo/": map[string]interface{}{
			"description": "foo",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"secret/": map[string]interface{}{
			"description": "key/value secret storage",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     true,
			"seal_wrap": false,
		},
		"identity/": map[string]interface{}{
			"description": "identity store",
			"type":        "identity",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	for k, v := range actual["data"].(map[string]interface{}) {
		if v.(map[string]interface{})["accessor"] == "" {
			t.Fatalf("no accessor from %s", k)
		}
		expected[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
		expected["data"].(map[string]interface{})[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
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

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")
	expected = map[string]interface{}{
		"lease_id":       "",
		"renewable":      false,
		"lease_duration": json.Number("0"),
		"wrap_info":      nil,
		"warnings":       nil,
		"auth":           nil,
		"data": map[string]interface{}{
			"foo/": map[string]interface{}{
				"description": "foo",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("259196400"),
					"max_lease_ttl":     json.Number("259200000"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"secret/": map[string]interface{}{
				"description": "key/value secret storage",
				"type":        "kv",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"sys/": map[string]interface{}{
				"description": "system endpoints used for control, policy and debugging",
				"type":        "system",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
			"cubbyhole/": map[string]interface{}{
				"description": "per-token private secret storage",
				"type":        "cubbyhole",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     true,
				"seal_wrap": false,
			},
			"identity/": map[string]interface{}{
				"description": "identity store",
				"type":        "identity",
				"config": map[string]interface{}{
					"default_lease_ttl": json.Number("0"),
					"max_lease_ttl":     json.Number("0"),
					"force_no_cache":    false,
					"plugin_name":       "",
				},
				"local":     false,
				"seal_wrap": false,
			},
		},
		"foo/": map[string]interface{}{
			"description": "foo",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("259196400"),
				"max_lease_ttl":     json.Number("259200000"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"secret/": map[string]interface{}{
			"description": "key/value secret storage",
			"type":        "kv",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
		"cubbyhole/": map[string]interface{}{
			"description": "per-token private secret storage",
			"type":        "cubbyhole",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     true,
			"seal_wrap": false,
		},
		"identity/": map[string]interface{}{
			"description": "identity store",
			"type":        "identity",
			"config": map[string]interface{}{
				"default_lease_ttl": json.Number("0"),
				"max_lease_ttl":     json.Number("0"),
				"force_no_cache":    false,
				"plugin_name":       "",
			},
			"local":     false,
			"seal_wrap": false,
		},
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	for k, v := range actual["data"].(map[string]interface{}) {
		if v.(map[string]interface{})["accessor"] == "" {
			t.Fatalf("no accessor from %s", k)
		}
		expected[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
		expected["data"].(map[string]interface{})[k].(map[string]interface{})["accessor"] = v.(map[string]interface{})["accessor"]
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
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
		"data": map[string]interface{}{
			"default_lease_ttl": json.Number("259196400"),
			"max_lease_ttl":     json.Number("259200000"),
			"force_no_cache":    false,
		},
		"default_lease_ttl": json.Number("259196400"),
		"max_lease_ttl":     json.Number("259200000"),
		"force_no_cache":    false,
	}

	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	expected["request_id"] = actual["request_id"]
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad:\nExpected: %#v\nActual:%#v", expected, actual)
	}

	// Set a low max
	resp = testHttpPost(t, token, addr+"/v1/sys/mounts/secret/tune", map[string]interface{}{
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
		"data": map[string]interface{}{
			"default_lease_ttl": json.Number("40"),
			"max_lease_ttl":     json.Number("80"),
			"force_no_cache":    false,
		},
		"default_lease_ttl": json.Number("40"),
		"max_lease_ttl":     json.Number("80"),
		"force_no_cache":    false,
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
