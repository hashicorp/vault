package vault

import (
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
)

// mockExpiration returns a mock expiration manager
func mockExpiration(t *testing.T) *ExpirationManager {
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
	view := NewBarrierView(b, "expire/")

	router := NewRouter()
	logger := log.New(os.Stderr, "", log.LstdFlags)
	return NewExpirationManager(router, view, logger)
}

func TestExpiration_Register(t *testing.T) {
	exp := mockExpiration(t)
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		IsSecret: true,
		Lease: &logical.Lease{
			Duration: time.Hour,
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	id, err := exp.Register(req, resp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !strings.HasPrefix(id, req.Path) {
		t.Fatalf("bad: %s", id)
	}

	if len(id) <= len(req.Path) {
		t.Fatalf("bad: %s", id)
	}
}

func TestExpiration_revokeEntry(t *testing.T) {
	exp := mockExpiration(t)

	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "noop", "", view)

	le := &leaseEntry{
		VaultID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Lease: &logical.Lease{
			Duration: time.Minute,
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
	}

	err := exp.revokeEntry(le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	req := noop.Requests[0]
	if req.Operation != logical.RevokeOperation {
		t.Fatalf("Bad: %v", req)
	}
	if req.Path != le.Path {
		t.Fatalf("Bad: %v", req)
	}
	if !reflect.DeepEqual(req.Data, le.Data) {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_renewEntry(t *testing.T) {
	exp := mockExpiration(t)

	noop := &NoopBackend{
		Response: &logical.Response{
			IsSecret: true,
			Lease: &logical.Lease{
				Renewable: true,
				Duration:  time.Hour,
			},
			Data: map[string]interface{}{
				"testing": false,
			},
		},
	}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "noop", "", view)

	le := &leaseEntry{
		VaultID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Lease: &logical.Lease{
			Duration: time.Minute,
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
	}

	resp, err := exp.renewEntry(le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(resp, noop.Response) {
		t.Fatalf("bad: %#v", resp)
	}

	req := noop.Requests[0]
	if req.Operation != logical.RenewOperation {
		t.Fatalf("Bad: %v", req)
	}
	if req.Path != le.Path {
		t.Fatalf("Bad: %v", req)
	}
	if !reflect.DeepEqual(req.Data, le.Data) {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_PersistLoadDelete(t *testing.T) {
	exp := mockExpiration(t)
	le := &leaseEntry{
		VaultID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Lease: &logical.Lease{
			Duration: time.Minute,
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
	}
	if err := exp.persistEntry(le); err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := exp.loadEntry("foo/bar/1234")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(out, le) {
		t.Fatalf("out: %#v expect: %#v", out, le)
	}

	err = exp.deleteEntry("foo/bar/1234")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err = exp.loadEntry("foo/bar/1234")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("out: %#v", out)
	}
}

func TestLeaseEntry(t *testing.T) {
	le := &leaseEntry{
		VaultID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Lease: &logical.Lease{
			Duration: time.Minute,
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
	}

	enc, err := le.encode()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := decodeLeaseEntry(enc)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(out.Data, le.Data) {
		t.Fatalf("got: %#v, expect %#v", out, le)
	}
}
