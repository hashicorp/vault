package vault

import (
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

func TestSystemBackend_RootPaths(t *testing.T) {
	expected := []string{
		"mount/*",
		"remount",
	}

	b := testSystemBackend(t)
	actual := b.RootPaths()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestSystemBackend_mounts(t *testing.T) {
	b := testSystemBackend(t)
	req := logical.TestRequest(t, logical.ReadOperation, "mounts")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	exp := map[string]interface{}{
		"secret/": map[string]string{
			"type":        "generic",
			"description": "generic secret storage",
		},
		"sys/": map[string]string{
			"type":        "system",
			"description": "system endpoints used for control, policy and debugging",
		},
	}
	if !reflect.DeepEqual(resp.Data, exp) {
		t.Fatalf("got: %#v expect: %#v", resp.Data, exp)
	}
}

func TestSystemBackend_mount(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "mount/prod/secret/")
	req.Data["type"] = "generic"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_mount_invalid(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "mount/prod/secret/")
	req.Data["type"] = "nope"
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "unknown backend type: nope" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_unmount(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.DeleteOperation, "mount/secret/")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_unmount_invalid(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.DeleteOperation, "mount/foo/")
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "no matching mount" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_remount(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "remount")
	req.Data["from"] = "secret"
	req.Data["to"] = "foo"
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_remount_invalid(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "remount")
	req.Data["from"] = "unknown"
	req.Data["to"] = "foo"
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "no matching mount at 'unknown/'" {
		t.Fatalf("bad: %v", resp)
	}
}

func TestSystemBackend_remount_system(t *testing.T) {
	b := testSystemBackend(t)

	req := logical.TestRequest(t, logical.WriteOperation, "remount")
	req.Data["from"] = "sys"
	req.Data["to"] = "foo"
	resp, err := b.HandleRequest(req)
	if err != logical.ErrInvalidRequest {
		t.Fatalf("err: %v", err)
	}
	if resp.Data["error"] != "cannot remount 'sys/'" {
		t.Fatalf("bad: %v", resp)
	}
}

func testSystemBackend(t *testing.T) logical.Backend {
	c, _ := TestCoreUnsealed(t)
	return NewSystemBackend(c)
}
