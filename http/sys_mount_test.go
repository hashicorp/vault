package http

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/vault"
)

func TestSysMounts(t *testing.T) {
	core, _, token := vault.TestCoreUnsealed(t)
	ln, addr := TestServer(t, core)
	defer ln.Close()
	TestServerAuth(t, addr, token)

	resp, err := http.Get(addr + "/v1/sys/mounts")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
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

	resp := testHttpPost(t, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp, err := http.Get(addr + "/v1/sys/mounts")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"foo/": map[string]interface{}{
			"description": "foo",
			"type":        "generic",
		},
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
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

	resp := testHttpPut(t, addr+"/v1/sys/mounts/foo", map[string]interface{}{
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

	resp := testHttpPost(t, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpPost(t, addr+"/v1/sys/remount", map[string]interface{}{
		"from": "foo",
		"to":   "bar",
	})
	testResponseStatus(t, resp, 204)

	resp, err := http.Get(addr + "/v1/sys/mounts")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"bar/": map[string]interface{}{
			"description": "foo",
			"type":        "generic",
		},
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
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

	resp := testHttpPost(t, addr+"/v1/sys/mounts/foo", map[string]interface{}{
		"type":        "generic",
		"description": "foo",
	})
	testResponseStatus(t, resp, 204)

	resp = testHttpDelete(t, addr+"/v1/sys/mounts/foo")
	testResponseStatus(t, resp, 204)

	resp, err := http.Get(addr + "/v1/sys/mounts")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var actual map[string]interface{}
	expected := map[string]interface{}{
		"secret/": map[string]interface{}{
			"description": "generic secret storage",
			"type":        "generic",
		},
		"sys/": map[string]interface{}{
			"description": "system endpoints used for control, policy and debugging",
			"type":        "system",
		},
	}
	testResponseStatus(t, resp, 200)
	testResponseBody(t, resp, &actual)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
