package http

import (
	"reflect"
	"testing"

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
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSysMount(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPost(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"foo/": map[string]interface{}{
			"description": "foo",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSysMount_put(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp := testHttpPut(t, token, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
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
		"type":        "generic",
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
		"bar/": map[string]interface{}{
			"description": "foo",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
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
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpDelete(t, token, addr+"/v1/sys/mounts/foo")
	testResponseStatus(t, resp, 204)

	resp = testHttpGet(t, token, addr+"/v1/sys/mounts")

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
			"config": map[string]interface{}{
				"default_lease_ttl": float64(0),
				"max_lease_ttl":     float64(0),
			},
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
