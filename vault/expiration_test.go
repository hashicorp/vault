package vault

import (
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
	return NewExpirationManager(router, view)
}

func TestExpiration_StartStop(t *testing.T) {
	exp := mockExpiration(t)
	err := exp.Start()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	err = exp.Restore()
	if err.Error() != "cannot restore while running" {
		t.Fatalf("err: %v", err)
	}

	err = exp.Stop()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
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
			Duration:    time.Hour,
			MaxDuration: time.Hour,
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

func TestLeaseEntry(t *testing.T) {
	le := &leaseEntry{
		VaultID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Lease: &logical.Lease{
			Renewable:   true,
			Duration:    time.Minute,
			MaxDuration: time.Hour,
		},
		IssueTime: time.Now(),
		RenewTime: time.Now(),
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
