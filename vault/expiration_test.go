package vault

import (
	"log"
	"os"
	"reflect"
	"sort"
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

	_, ts, _ := mockTokenStore(t)

	router := NewRouter()
	logger := log.New(os.Stderr, "", log.LstdFlags)
	exp := NewExpirationManager(router, view, ts, logger)
	ts.SetExpirationManager(exp)
	return exp
}

func TestExpiration_Restore(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", generateUUID(), view)

	paths := []string{
		"prod/aws/foo",
		"prod/aws/sub/bar",
		"prod/aws/zip",
	}
	for _, path := range paths {
		req := &logical.Request{
			Operation: logical.ReadOperation,
			Path:      path,
		}
		resp := &logical.Response{
			Secret: &logical.Secret{
				Lease: 20 * time.Millisecond,
			},
			Data: map[string]interface{}{
				"access_key": "xyz",
				"secret_key": "abcd",
			},
		}
		_, err := exp.Register(req, resp)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Stop everything
	err := exp.Stop()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Restore
	err = exp.Restore()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Ensure all are reaped
	start := time.Now()
	for time.Now().Sub(start) < time.Second {
		if len(noop.Requests) < 3 {
			time.Sleep(5 * time.Millisecond)
			continue
		}
		break
	}
	for _, req := range noop.Requests {
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}
	}
}

func TestExpiration_Register(t *testing.T) {
	exp := mockExpiration(t)
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			Lease: time.Hour,
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

func TestExpiration_RegisterAuth(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.RootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	auth := &logical.Auth{
		ClientToken: root.ID,
		Lease:       time.Hour,
	}

	err = exp.RegisterAuth("auth/github/login", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestExpiration_Revoke(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", generateUUID(), view)

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			Lease: time.Hour,
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

	if err := exp.Revoke(id); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = noop.Requests[0]
	if req.Operation != logical.RevokeOperation {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_RevokeOnExpire(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", generateUUID(), view)

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			Lease: 20 * time.Millisecond,
		},
		Data: map[string]interface{}{
			"access_key": "xyz",
			"secret_key": "abcd",
		},
	}

	_, err := exp.Register(req, resp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	start := time.Now()
	for time.Now().Sub(start) < time.Second {
		if len(noop.Requests) == 0 {
			time.Sleep(5 * time.Millisecond)
			continue
		}
		req = noop.Requests[0]
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}
		break
	}
}

func TestExpiration_RevokePrefix(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", generateUUID(), view)

	paths := []string{
		"prod/aws/foo",
		"prod/aws/sub/bar",
		"prod/aws/zip",
	}
	for _, path := range paths {
		req := &logical.Request{
			Operation: logical.ReadOperation,
			Path:      path,
		}
		resp := &logical.Response{
			Secret: &logical.Secret{
				Lease: 20 * time.Millisecond,
			},
			Data: map[string]interface{}{
				"access_key": "xyz",
				"secret_key": "abcd",
			},
		}
		_, err := exp.Register(req, resp)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	// Should nuke all the keys
	if err := exp.RevokePrefix("prod/aws/"); err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(noop.Requests) != 3 {
		t.Fatalf("Bad: %v", noop.Requests)
	}
	for _, req := range noop.Requests {
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}
	}

	expect := []string{
		"foo",
		"sub/bar",
		"zip",
	}
	sort.Strings(noop.Paths)
	sort.Strings(expect)
	if !reflect.DeepEqual(noop.Paths, expect) {
		t.Fatalf("bad: %v", noop.Paths)
	}
}

func TestExpiration_RenewToken(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.RootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a token
	auth := &logical.Auth{
		ClientToken: root.ID,
		Lease:       time.Hour,
	}
	err = exp.RegisterAuth("auth/github/login", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Renew the token
	out, err := exp.RenewToken("auth/github/login", root.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !reflect.DeepEqual(auth, out) {
		t.Fatalf("Bad: %#v", out)
	}
}

func TestExpiration_Renew(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", generateUUID(), view)

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			Lease: 20 * time.Millisecond,
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

	noop.Response = &logical.Response{
		Secret: &logical.Secret{
			Lease: 20 * time.Millisecond,
		},
		Data: map[string]interface{}{
			"access_key": "123",
			"secret_key": "abcd",
		},
	}

	out, err := exp.Renew(id, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(out, noop.Response) {
		t.Fatalf("Bad: %#v", out)
	}

	if len(noop.Requests) != 1 {
		t.Fatalf("Bad: %#v", noop.Requests)
	}
	req = noop.Requests[0]
	if req.Operation != logical.RenewOperation {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_Renew_RevokeOnExpire(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", generateUUID(), view)

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			Lease: 20 * time.Millisecond,
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

	noop.Response = &logical.Response{
		Secret: &logical.Secret{
			Lease: 20 * time.Millisecond,
		},
		Data: map[string]interface{}{
			"access_key": "123",
			"secret_key": "abcd",
		},
	}

	_, err = exp.Renew(id, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	start := time.Now()
	for time.Now().Sub(start) < time.Second {
		if len(noop.Requests) < 2 {
			time.Sleep(5 * time.Millisecond)
			continue
		}
		req = noop.Requests[1]
		if req.Operation != logical.RevokeOperation {
			t.Fatalf("Bad: %v", req)
		}
		break
	}
}

func TestExpiration_revokeEntry(t *testing.T) {
	exp := mockExpiration(t)

	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "", generateUUID(), view)

	le := &leaseEntry{
		VaultID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			Lease: time.Minute,
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

func TestExpiration_revokeEntry_token(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.RootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	le := &leaseEntry{
		VaultID: "foo/bar/1234",
		Auth: &logical.Auth{
			ClientToken: root.ID,
			Lease:       time.Minute,
		},
		Path:       "foo/bar",
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
	}

	err = exp.revokeEntry(le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := exp.tokenStore.Lookup(root.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}
}

func TestExpiration_renewEntry(t *testing.T) {
	exp := mockExpiration(t)

	noop := &NoopBackend{
		Response: &logical.Response{
			Secret: &logical.Secret{
				Renewable: true,
				Lease:     time.Hour,
			},
			Data: map[string]interface{}{
				"testing": false,
			},
		},
	}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "", generateUUID(), view)

	le := &leaseEntry{
		VaultID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			Lease: time.Minute,
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
	}

	resp, err := exp.renewEntry(le, time.Second)
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
	if req.Secret.LeaseIncrement != time.Second {
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
		Secret: &logical.Secret{
			Lease: time.Minute,
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
		Secret: &logical.Secret{
			Lease: time.Minute,
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
