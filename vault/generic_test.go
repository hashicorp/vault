package vault

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/physical"
)

// mockView returns a view attached to a barrier / backend
func mockView(t *testing.T, prefix string) *BarrierView {
	inm := physical.NewInmem()
	b, err := NewAESGCMBarrier(inm)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Initialize and unseal
	key, _ := b.GenerateKey()
	b.Initialize(key)
	b.Unseal(key)

	// Create the barrier view
	view := NewBarrierView(b, prefix)
	return view
}

// mockRequest returns a request with a real view attached
func mockRequest(t *testing.T, op Operation, path string) *Request {
	view := mockView(t, "logical/")

	// Create the request
	req := &Request{
		Operation: op,
		Path:      path,
		Data:      make(map[string]interface{}),
		View:      view,
	}
	return req
}

func TestGenericBackend_RootPaths(t *testing.T) {
	b, err := newGenericBackend(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	root := b.RootPaths()
	if len(root) != 0 {
		t.Fatalf("unexpected: %v", root)
	}
}

func TestGenericBackend_Write(t *testing.T) {
	b, err := newGenericBackend(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := mockRequest(t, WriteOperation, "foo")
	req.Data["raw"] = "test"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	out, err := req.View.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("failed to write to view")
	}
}

func TestGenericBackend_Read(t *testing.T) {
	b, err := newGenericBackend(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := mockRequest(t, WriteOperation, "foo")
	req.Data["raw"] = "test"
	req.Data["lease"] = "1h"

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req2 := mockRequest(t, ReadOperation, "foo")
	req2.View = req.View

	resp, err := b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !resp.IsSecret {
		t.Fatalf("should be secret: %#v", resp)
	}

	if resp.Lease == nil {
		t.Fatalf("should have lease: %#v", resp)
	}
	if resp.Lease.Renewable {
		t.Fatalf("bad lease: %#v", resp.Lease)
	}
	if resp.Lease.Revokable {
		t.Fatalf("bad lease: %#v", resp.Lease)
	}
	if resp.Lease.Duration != time.Hour {
		t.Fatalf("bad lease: %#v", resp.Lease)
	}
	if resp.Lease.MaxDuration != time.Hour {
		t.Fatalf("bad lease: %#v", resp.Lease)
	}
	if resp.Lease.MaxIncrement != 0 {
		t.Fatalf("bad lease: %#v", resp.Lease)
	}

	if resp.Data["raw"] != "test" {
		t.Fatalf("bad data: %#v", resp.Data)
	}
	if resp.Data["lease"] != "1h" {
		t.Fatalf("bad data: %#v", resp.Data)
	}
}

func TestGenericBackend_Delete(t *testing.T) {
	b, err := newGenericBackend(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := mockRequest(t, WriteOperation, "foo")
	req.Data["raw"] = "test"
	req.Data["lease"] = "1h"

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req2 := mockRequest(t, DeleteOperation, "foo")
	req2.View = req.View

	resp, err := b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req3 := mockRequest(t, ReadOperation, "foo")
	req3.View = req.View

	resp, err = b.HandleRequest(req3)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestGenericBackend_List(t *testing.T) {
	b, err := newGenericBackend(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := mockRequest(t, WriteOperation, "foo/bar")
	req.Data["raw"] = "test"
	req.Data["lease"] = "1h"

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req2 := mockRequest(t, ListOperation, "")
	req2.View = req.View

	resp, err := b.HandleRequest(req2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if resp.IsSecret {
		t.Fatalf("bad: %v", resp)
	}
	if resp.Lease != nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.Data["keys"] == nil {
		t.Fatalf("bad: %v", resp)
	}
	keys := resp.Data["keys"].([]string)
	if len(keys) != 1 || keys[0] != "foo/" {
		t.Fatalf("keys: %v", keys)
	}
}

func TestGenericBackend_Help(t *testing.T) {
	b, err := newGenericBackend(nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := mockRequest(t, HelpOperation, "foo")
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if resp.IsSecret {
		t.Fatalf("bad: %v", resp)
	}
	if resp.Lease != nil {
		t.Fatalf("bad: %v", resp)
	}
	if resp.Data["help"] != genericHelpText {
		t.Fatalf("bad: %v", resp)
	}
}
