package vault

import (
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/uuid"
	"github.com/hashicorp/vault/logical"
)

// mockExpiration returns a mock expiration manager
func mockExpiration(t *testing.T) *ExpirationManager {
	_, ts, _ := mockTokenStore(t)
	return ts.expiration
}

func TestExpiration_Restore(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: uuid.GenerateUUID()}, view)

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
				LeaseOptions: logical.LeaseOptions{
					TTL: 20 * time.Millisecond,
				},
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
		noop.Lock()
		less := len(noop.Requests) < 3
		noop.Unlock()

		if less {
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
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
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
		LeaseOptions: logical.LeaseOptions{
			TTL: time.Hour,
		},
	}

	err = exp.RegisterAuth("auth/github/login", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestExpiration_RegisterAuth_NoLease(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.RootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	auth := &logical.Auth{
		ClientToken: root.ID,
	}

	err = exp.RegisterAuth("auth/github/login", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Should not be able to renew, no expiration
	_, err = exp.RenewToken("auth/github/login", root.ID, 0)
	if err.Error() != "lease not found or lease is not renewable" {
		t.Fatalf("err: %v", err)
	}

	// Wait and check token is not invalidated
	time.Sleep(20 * time.Millisecond)

	// Verify token does not get revoked
	out, err := exp.tokenStore.Lookup(root.ID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("missing token")
	}
}

func TestExpiration_Revoke(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: uuid.GenerateUUID()}, view)

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Hour,
			},
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
	exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: uuid.GenerateUUID()}, view)

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: 20 * time.Millisecond,
			},
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
		req = nil

		noop.Lock()
		if len(noop.Requests) > 0 {
			req = noop.Requests[0]
		}
		noop.Unlock()
		if req == nil {
			time.Sleep(5 * time.Millisecond)
			continue
		}
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
	exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: uuid.GenerateUUID()}, view)

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
				LeaseOptions: logical.LeaseOptions{
					TTL: 20 * time.Millisecond,
				},
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

func TestExpiration_RevokeByToken(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: uuid.GenerateUUID()}, view)

	paths := []string{
		"prod/aws/foo",
		"prod/aws/sub/bar",
		"prod/aws/zip",
	}
	for _, path := range paths {
		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        path,
			ClientToken: "foobarbaz",
		}
		resp := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL: 20 * time.Millisecond,
				},
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
	if err := exp.RevokeByToken("foobarbaz"); err != nil {
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
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
	}
	err = exp.RegisterAuth("auth/token/login", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Renew the token
	out, err := exp.RenewToken("auth/token/login", root.ID, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if auth.ClientToken != out.ClientToken {
		t.Fatalf("Bad: %#v", out)
	}
}

func TestExpiration_RenewToken_NotRenewable(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.RootToken()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: false,
		},
	}
	err = exp.RegisterAuth("auth/github/login", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Attempt to renew the token
	_, err = exp.RenewToken("auth/github/login", root.ID, 0)
	if err.Error() != "lease is not renewable" {
		t.Fatalf("err: %v", err)
	}
}

func TestExpiration_Renew(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: uuid.GenerateUUID()}, view)

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       20 * time.Millisecond,
				Renewable: true,
			},
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
			LeaseOptions: logical.LeaseOptions{
				TTL: 20 * time.Millisecond,
			},
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

	noop.Lock()
	defer noop.Unlock()

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

func TestExpiration_Renew_NotRenewable(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: uuid.GenerateUUID()}, view)

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       20 * time.Millisecond,
				Renewable: false,
			},
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

	_, err = exp.Renew(id, 0)
	if err.Error() != "lease is not renewable" {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

	if len(noop.Requests) != 0 {
		t.Fatalf("Bad: %#v", noop.Requests)
	}
}

func TestExpiration_Renew_RevokeOnExpire(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "prod/aws/", &MountEntry{UUID: uuid.GenerateUUID()}, view)

	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "prod/aws/foo",
	}
	resp := &logical.Response{
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL:       20 * time.Millisecond,
				Renewable: true,
			},
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
			LeaseOptions: logical.LeaseOptions{
				TTL: 20 * time.Millisecond,
			},
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
		req = nil

		noop.Lock()
		if len(noop.Requests) >= 2 {
			req = noop.Requests[1]
		}
		noop.Unlock()

		if req == nil {
			time.Sleep(5 * time.Millisecond)
			continue
		}
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
	exp.router.Mount(noop, "", &MountEntry{UUID: uuid.GenerateUUID()}, view)

	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
	}

	err := exp.revokeEntry(le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

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
		LeaseID: "foo/bar/1234",
		Auth: &logical.Auth{
			ClientToken: root.ID,
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
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
				LeaseOptions: logical.LeaseOptions{
					Renewable: true,
					TTL:       time.Hour,
				},
			},
			Data: map[string]interface{}{
				"testing": false,
			},
		},
	}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	exp.router.Mount(noop, "", &MountEntry{UUID: uuid.GenerateUUID()}, view)

	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now(),
	}

	resp, err := exp.renewEntry(le, time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

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
	if req.Secret.Increment != time.Second {
		t.Fatalf("Bad: %v", req)
	}
	if req.Secret.IssueTime.IsZero() {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_renewAuthEntry(t *testing.T) {
	exp := mockExpiration(t)

	noop := &NoopBackend{
		Response: &logical.Response{
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					Renewable: true,
					TTL:       time.Hour,
				},
			},
		},
	}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "auth/foo/")
	exp.router.Mount(noop, "auth/foo/", &MountEntry{UUID: uuid.GenerateUUID()}, view)

	le := &leaseEntry{
		LeaseID: "auth/foo/1234",
		Path:    "auth/foo/login",
		Auth: &logical.Auth{
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       time.Minute,
			},
			InternalData: map[string]interface{}{
				"MySecret": "secret",
			},
		},
		IssueTime:  time.Now(),
		ExpireTime: time.Now().Add(time.Minute),
	}

	resp, err := exp.renewAuthEntry(le, time.Second)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

	if !reflect.DeepEqual(resp, noop.Response) {
		t.Fatalf("bad: %#v", resp)
	}

	req := noop.Requests[0]
	if req.Operation != logical.RenewOperation {
		t.Fatalf("Bad: %v", req)
	}
	if req.Path != "login" {
		t.Fatalf("Bad: %v", req)
	}
	if req.Auth.Increment != time.Second {
		t.Fatalf("Bad: %v", req)
	}
	if req.Auth.IssueTime.IsZero() {
		t.Fatalf("Bad: %v", req)
	}
	if req.Auth.InternalData["MySecret"] != "secret" {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_PersistLoadDelete(t *testing.T) {
	exp := mockExpiration(t)
	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		IssueTime:  time.Now().UTC(),
		ExpireTime: time.Now().UTC(),
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
		LeaseID: "foo/bar/1234",
		Path:    "foo/bar",
		Data: map[string]interface{}{
			"testing": true,
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		IssueTime:  time.Now().UTC(),
		ExpireTime: time.Now().UTC(),
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
