package vault

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/physical/inmem"
)

var (
	testImagePull sync.Once
)

// mockExpiration returns a mock expiration manager
func mockExpiration(t testing.TB) *ExpirationManager {
	c, _, _ := TestCoreUnsealed(t)
	return c.expiration
}

func mockBackendExpiration(t testing.TB, backend physical.Backend) (*Core, *ExpirationManager) {
	c, _, _ := TestCoreUnsealedBackend(t, backend)
	return c, c.expiration
}

func TestExpiration_Tidy(t *testing.T) {
	var err error

	exp := mockExpiration(t)

	// We use this later for tidy testing where we need to check the output
	logOut := new(bytes.Buffer)
	logger := log.New(&log.LoggerOptions{
		Output: logOut,
	})
	exp.logger = logger

	if err := exp.Restore(nil); err != nil {
		t.Fatal(err)
	}

	// Set up a count function to calculate number of leases
	count := 0
	countFunc := func(leaseID string) {
		count++
	}

	// Scan the storage with the count func set
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Check that there are no leases to begin with
	if count != 0 {
		t.Fatalf("bad: lease count; expected:0 actual:%d", count)
	}

	// Create a lease entry without a client token in it
	le := &leaseEntry{
		LeaseID: "lease/with/no/client/token",
		Path:    "foo/bar",
	}

	// Persist the invalid lease entry
	if err = exp.persistEntry(le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Check that the storage was successful and that the count of leases is
	// now 1
	if count != 1 {
		t.Fatalf("bad: lease count; expected:1 actual:%d", count)
	}

	// Run the tidy operation
	err = exp.Tidy()
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	if err := logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Post the tidy operation, the invalid lease entry should have been gone
	if count != 0 {
		t.Fatalf("bad: lease count; expected:0 actual:%d", count)
	}

	// Set a revoked/invalid token in the lease entry
	le.ClientToken = "invalidtoken"

	// Persist the invalid lease entry
	if err = exp.persistEntry(le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Check that the storage was successful and that the count of leases is
	// now 1
	if count != 1 {
		t.Fatalf("bad: lease count; expected:1 actual:%d", count)
	}

	// Run the tidy operation
	err = exp.Tidy()
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Post the tidy operation, the invalid lease entry should have been gone
	if count != 0 {
		t.Fatalf("bad: lease count; expected:0 actual:%d", count)
	}

	// Attach an invalid token with 2 leases
	if err = exp.persistEntry(le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	le.LeaseID = "another/invalid/lease"
	if err = exp.persistEntry(le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	// Run the tidy operation
	err = exp.Tidy()
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Post the tidy operation, the invalid lease entry should have been gone
	if count != 0 {
		t.Fatalf("bad: lease count; expected:0 actual:%d", count)
	}

	for i := 0; i < 1000; i++ {
		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        "invalid/lease/" + fmt.Sprintf("%d", i+1),
			ClientToken: "invalidtoken",
		}
		resp := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL: 100 * time.Millisecond,
				},
			},
			Data: map[string]interface{}{
				"test_key": "test_value",
			},
		}
		_, err := exp.Register(req, resp)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Check that there are 1000 leases now
	if count != 1000 {
		t.Fatalf("bad: lease count; expected:1000 actual:%d", count)
	}

	errCh1 := make(chan error)
	errCh2 := make(chan error)

	// Initiate tidy of the above 1000 invalid leases in quick succession. Only
	// one tidy operation can be in flight at any time. One of these requests
	// should error out.
	go func() {
		errCh1 <- exp.Tidy()
	}()

	go func() {
		errCh2 <- exp.Tidy()
	}()

	var err1, err2 error

	for i := 0; i < 2; i++ {
		select {
		case err1 = <-errCh1:
		case err2 = <-errCh2:
		}
	}

	if err1 != nil || err2 != nil {
		t.Fatalf("got an error: err1: %v; err2: %v", err1, err2)
	}
	if !strings.Contains(logOut.String(), "tidy operation on leases is already in progress") {
		t.Fatalf("expected to see a warning saying operation in progress, output is %s", logOut.String())
	}

	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	le.ClientToken = root.ID

	// Attach a valid token with the leases
	if err = exp.persistEntry(le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}

	// Run the tidy operation
	err = exp.Tidy()
	if err != nil {
		t.Fatal(err)
	}

	count = 0
	if err = logical.ScanView(context.Background(), exp.idView, countFunc); err != nil {
		t.Fatal(err)
	}

	// Post the tidy operation, the valid lease entry should not get affected
	if count != 1 {
		t.Fatalf("bad: lease count; expected:1 actual:%d", count)
	}
}

// To avoid pulling in deps for all users of the package, don't leave these
// uncommented in the public tree
/*
func BenchmarkExpiration_Restore_Etcd(b *testing.B) {
	addr := os.Getenv("PHYSICAL_BACKEND_BENCHMARK_ADDR")
	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())

	logger := logging.NewVaultLogger(log.Trace)
	physicalBackend, err := physEtcd.NewEtcdBackend(map[string]string{
		"address":      addr,
		"path":         randPath,
		"max_parallel": "256",
	}, logger)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	benchmarkExpirationBackend(b, physicalBackend, 10000) // 10,000 leases
}

func BenchmarkExpiration_Restore_Consul(b *testing.B) {
	addr := os.Getenv("PHYSICAL_BACKEND_BENCHMARK_ADDR")
	randPath := fmt.Sprintf("vault-%d/", time.Now().Unix())

	logger := logging.NewVaultLogger(log.Trace)
	physicalBackend, err := physConsul.NewConsulBackend(map[string]string{
		"address":      addr,
		"path":         randPath,
		"max_parallel": "256",
	}, logger)
	if err != nil {
		b.Fatalf("err: %s", err)
	}

	benchmarkExpirationBackend(b, physicalBackend, 10000) // 10,000 leases
}
*/

func BenchmarkExpiration_Restore_InMem(b *testing.B) {
	logger := logging.NewVaultLogger(log.Trace)
	inm, err := inmem.NewInmem(nil, logger)
	if err != nil {
		b.Fatal(err)
	}
	benchmarkExpirationBackend(b, inm, 100000) // 100,000 Leases
}

func benchmarkExpirationBackend(b *testing.B, physicalBackend physical.Backend, numLeases int) {
	c, _, _ := TestCoreUnsealedBackend(b, physicalBackend)
	exp := c.expiration
	noop := &NoopBackend{}
	view := NewBarrierView(c.barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		b.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		b.Fatal(err)
	}

	// Register fake leases
	for i := 0; i < numLeases; i++ {
		pathUUID, err := uuid.GenerateUUID()
		if err != nil {
			b.Fatal(err)
		}

		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        "prod/aws/" + pathUUID,
			ClientToken: "root",
		}
		resp := &logical.Response{
			Secret: &logical.Secret{
				LeaseOptions: logical.LeaseOptions{
					TTL: 400 * time.Second,
				},
			},
			Data: map[string]interface{}{
				"access_key": "xyz",
				"secret_key": "abcd",
			},
		}
		_, err = exp.Register(req, resp)
		if err != nil {
			b.Fatalf("err: %v", err)
		}
	}

	// Stop everything
	err = exp.Stop()
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err = exp.Restore(nil)
		// Restore
		if err != nil {
			b.Fatalf("err: %v", err)
		}
	}
	b.StopTimer()
}

func TestExpiration_Restore(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	paths := []string{
		"prod/aws/foo",
		"prod/aws/sub/bar",
		"prod/aws/zip",
	}
	for _, path := range paths {
		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        path,
			ClientToken: "foobar",
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
	err = exp.Stop()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Restore
	err = exp.Restore(nil)
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
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
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
	root, err := exp.tokenStore.rootToken(context.Background())
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

	err = exp.RegisterAuth("auth/github/../login", auth)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExpiration_RegisterAuth_NoLease(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.rootToken(context.Background())
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
	resp, err := exp.RenewToken(&logical.Request{}, "auth/github/login", root.ID, 0)
	if err != nil && (err != logical.ErrInvalidRequest || (resp != nil && resp.IsError() && resp.Error().Error() != "lease not found or lease is not renewable")) {
		t.Fatalf("bad: err:%v resp:%#v", err, resp)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

	// Wait and check token is not invalidated
	time.Sleep(20 * time.Millisecond)

	// Verify token does not get revoked
	out, err := exp.tokenStore.Lookup(context.Background(), root.ID)
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
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
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
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
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

	_, err = exp.Register(req, resp)
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
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	paths := []string{
		"prod/aws/foo",
		"prod/aws/sub/bar",
		"prod/aws/zip",
	}
	for _, path := range paths {
		req := &logical.Request{
			Operation:   logical.ReadOperation,
			Path:        path,
			ClientToken: "foobar",
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
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

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
	te := &logical.TokenEntry{
		ID: "foobarbaz",
	}
	if err := exp.RevokeByToken(te); err != nil {
		t.Fatalf("err: %v", err)
	}

	time.Sleep(300 * time.Millisecond)

	noop.Lock()
	defer noop.Unlock()

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

func TestExpiration_RevokeByToken_Blocking(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	// Request handle with a timeout context that simulates blocking lease revocation.
	noop.RequestHandler = func(ctx context.Context, req *logical.Request) (*logical.Response, error) {
		ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
		defer cancel()

		select {
		case <-ctx.Done():
			return noop.Response, nil
		}
	}

	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

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
					TTL: 1 * time.Minute,
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
	te := &logical.TokenEntry{
		ID: "foobarbaz",
	}
	if err := exp.RevokeByToken(te); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Lock and check that no requests has gone through yet
	noop.Lock()
	if len(noop.Requests) != 0 {
		t.Fatalf("Bad: %v", noop.Requests)
	}
	noop.Unlock()

	// Wait for a bit for timeouts to trigger and pending revocations to go
	// through and then we relock
	time.Sleep(300 * time.Millisecond)

	noop.Lock()
	defer noop.Unlock()

	// Now make sure that all requests have gone through
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
	root, err := exp.tokenStore.rootToken(context.Background())
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
	out, err := exp.RenewToken(&logical.Request{}, "auth/token/login", root.ID, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if auth.ClientToken != out.Auth.ClientToken {
		t.Fatalf("bad: %#v", out)
	}
}

func TestExpiration_RenewToken_period(t *testing.T) {
	exp := mockExpiration(t)
	root := &logical.TokenEntry{
		Policies:     []string{"root"},
		Path:         "auth/token/root",
		DisplayName:  "root",
		CreationTime: time.Now().Unix(),
		Period:       time.Minute,
	}
	if err := exp.tokenStore.create(context.Background(), root); err != nil {
		t.Fatalf("err: %v", err)
	}

	// Register a token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL:       time.Hour,
			Renewable: true,
		},
		Period: time.Minute,
	}
	err := exp.RegisterAuth("auth/token/login", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Renew the token
	out, err := exp.RenewToken(&logical.Request{}, "auth/token/login", root.ID, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if auth.ClientToken != out.Auth.ClientToken {
		t.Fatalf("bad: %#v", out)
	}

	if out.Auth.TTL > time.Minute {
		t.Fatalf("expected TTL to be less than 1 minute, got: %s", out.Auth.TTL)
	}
}

func TestExpiration_RenewToken_period_backend(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Mount a noop backend
	noop := &NoopBackend{
		Response: &logical.Response{
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL:       10 * time.Second,
					Renewable: true,
				},
				Period: 5 * time.Second,
			},
		},
		DefaultLeaseTTL: 5 * time.Second,
		MaxLeaseTTL:     5 * time.Second,
	}

	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, credentialBarrierPrefix)
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "auth/foo/", &MountEntry{Path: "auth/foo/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	// Register a token
	auth := &logical.Auth{
		ClientToken: root.ID,
		LeaseOptions: logical.LeaseOptions{
			TTL:       10 * time.Second,
			Renewable: true,
		},
		Period: 5 * time.Second,
	}

	err = exp.RegisterAuth("auth/foo/login", auth)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// Wait 3 seconds
	time.Sleep(3 * time.Second)
	resp, err := exp.RenewToken(&logical.Request{}, "auth/foo/login", root.ID, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	if resp.Auth.TTL == 0 || resp.Auth.TTL > 5*time.Second {
		t.Fatalf("expected TTL to be greater than zero and less than or equal to period, got: %s", resp.Auth.TTL)
	}

	// Wait another 3 seconds. If period works correctly, this should not fail
	time.Sleep(3 * time.Second)
	resp, err = exp.RenewToken(&logical.Request{}, "auth/foo/login", root.ID, 0)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}
	if resp.Auth.TTL < 4*time.Second || resp.Auth.TTL > 5*time.Second {
		t.Fatalf("expected TTL to be around period's value, got: %s", resp.Auth.TTL)
	}
}

func TestExpiration_RenewToken_NotRenewable(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.rootToken(context.Background())
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
	resp, err := exp.RenewToken(&logical.Request{}, "auth/github/login", root.ID, 0)
	if err != nil && (err != logical.ErrInvalidRequest || (resp != nil && resp.IsError() && resp.Error().Error() != "lease is not renewable")) {
		t.Fatalf("bad: err:%v resp:%#v", err, resp)
	}
	if resp == nil {
		t.Fatal("expected a response")
	}

}

func TestExpiration_Renew(t *testing.T) {
	exp := mockExpiration(t)
	noop := &NoopBackend{}
	_, barrier, _ := mockBarrier(t)
	view := NewBarrierView(barrier, "logical/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
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
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
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
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "prod/aws/", &MountEntry{Path: "prod/aws/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "prod/aws/foo",
		ClientToken: "foobar",
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
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "foo/bar/", &MountEntry{Path: "foo/bar/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

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

	err = exp.revokeEntry(le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	noop.Lock()
	defer noop.Unlock()

	req := noop.Requests[0]
	if req.Operation != logical.RevokeOperation {
		t.Fatalf("bad: operation; req: %#v", req)
	}
	if !reflect.DeepEqual(req.Data, le.Data) {
		t.Fatalf("bad: data; req: %#v\n le: %#v\n", req, le)
	}
}

func TestExpiration_revokeEntry_token(t *testing.T) {
	exp := mockExpiration(t)
	root, err := exp.tokenStore.rootToken(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	// N.B.: Vault doesn't allow both a secret and auth to be returned, but the
	// reason for both is that auth needs to be included in order to use the
	// token store as it's the only mounted backend, *but* RegisterAuth doesn't
	// actually create the index by token, only Register (for a Secret) does.
	// So without the Secret we don't do anything when removing the index which
	// (at the time of writing) now fails because a bug causing every token
	// expiration to do an extra delete to a non-existent key has been fixed,
	// and this test relies on this nonstandard behavior.
	le := &leaseEntry{
		LeaseID: "foo/bar/1234",
		Auth: &logical.Auth{
			ClientToken: root.ID,
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		Secret: &logical.Secret{
			LeaseOptions: logical.LeaseOptions{
				TTL: time.Minute,
			},
		},
		ClientToken: root.ID,
		Path:        "foo/bar",
		IssueTime:   time.Now(),
		ExpireTime:  time.Now(),
	}

	if err := exp.persistEntry(le); err != nil {
		t.Fatalf("error persisting entry: %v", err)
	}
	if err := exp.createIndexByToken(le.ClientToken, le.LeaseID); err != nil {
		t.Fatalf("error creating secondary index: %v", err)
	}
	exp.updatePending(le, le.Secret.LeaseTotal())

	indexEntry, err := exp.indexByToken(le.ClientToken, le.LeaseID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if indexEntry == nil {
		t.Fatalf("err: should have found a secondary index entry")
	}

	err = exp.revokeEntry(le)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	time.Sleep(300 * time.Millisecond)

	out, err := exp.tokenStore.Lookup(context.Background(), le.ClientToken)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out != nil {
		t.Fatalf("bad: %v", out)
	}

	indexEntry, err = exp.indexByToken(le.ClientToken, le.LeaseID)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if indexEntry != nil {
		t.Fatalf("err: should not have found a secondary index entry")
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
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "foo/bar/", &MountEntry{Path: "foo/bar/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

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

	resp, err := exp.renewEntry(le, 0)
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
	if !reflect.DeepEqual(req.Data, le.Data) {
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
	view := NewBarrierView(barrier, "auth/")
	meUUID, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatal(err)
	}
	err = exp.router.Mount(noop, "auth/foo/", &MountEntry{Path: "auth/foo/", Type: "noop", UUID: meUUID, Accessor: "noop-accessor"}, view)
	if err != nil {
		t.Fatal(err)
	}

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

	resp, err := exp.renewAuthEntry(&logical.Request{}, le, 0)
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
	if req.Auth.InternalData["MySecret"] != "secret" {
		t.Fatalf("Bad: %v", req)
	}
}

func TestExpiration_PersistLoadDelete(t *testing.T) {
	exp := mockExpiration(t)
	lastTime := time.Now()
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
		IssueTime:       lastTime,
		ExpireTime:      lastTime,
		LastRenewalTime: lastTime,
	}
	if err := exp.persistEntry(le); err != nil {
		t.Fatalf("err: %v", err)
	}

	out, err := exp.loadEntry("foo/bar/1234")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !le.LastRenewalTime.Equal(out.LastRenewalTime) ||
		!le.IssueTime.Equal(out.IssueTime) ||
		!le.ExpireTime.Equal(out.ExpireTime) {
		t.Fatalf("bad: expected:\n%#v\nactual:\n%#v", le, out)
	}
	le.LastRenewalTime = out.LastRenewalTime
	le.IssueTime = out.IssueTime
	le.ExpireTime = out.ExpireTime
	if !reflect.DeepEqual(out, le) {
		t.Fatalf("bad: expected:\n%#v\nactual:\n%#v", le, out)
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
				TTL:       time.Minute,
				Renewable: true,
			},
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

	// Test renewability
	le.ExpireTime = time.Time{}
	if r, _ := le.renewable(); r {
		t.Fatal("lease with zero expire time is not renewable")
	}
	le.ExpireTime = time.Now().Add(-1 * time.Hour)
	if r, _ := le.renewable(); r {
		t.Fatal("lease with expire time in the past is not renewable")
	}
	le.ExpireTime = time.Now().Add(1 * time.Hour)
	if r, err := le.renewable(); !r {
		t.Fatalf("lease with future expire time is renewable, err: %v", err)
	}
	le.Secret.LeaseOptions.Renewable = false
	if r, _ := le.renewable(); r {
		t.Fatal("secret is set to not be renewable but returns as renewable")
	}
	le.Secret = nil
	le.Auth = &logical.Auth{
		LeaseOptions: logical.LeaseOptions{
			Renewable: true,
		},
	}
	if r, err := le.renewable(); !r {
		t.Fatalf("auth is renewable but is set to not be, err: %v", err)
	}
	le.Auth.LeaseOptions.Renewable = false
	if r, _ := le.renewable(); r {
		t.Fatal("auth is set to not be renewable but returns as renewable")
	}
}

func TestExpiration_RevokeForce(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.logicalBackends["badrenew"] = badRenewFactory
	me := &MountEntry{
		Table:    mountTableType,
		Path:     "badrenew/",
		Type:     "badrenew",
		Accessor: "badrenewaccessor",
	}

	err := core.mount(context.Background(), me)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "badrenew/creds",
		ClientToken: root,
	}

	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("response was nil")
	}
	if resp.Secret == nil {
		t.Fatalf("response secret was nil, response was %#v", *resp)
	}

	req.Operation = logical.UpdateOperation
	req.Path = "sys/revoke-prefix/badrenew/creds"

	resp, err = core.HandleRequest(req)
	if err == nil {
		t.Fatal("expected error")
	}

	req.Path = "sys/revoke-force/badrenew/creds"
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}
}

func TestExpiration_RevokeForceSingle(t *testing.T) {
	core, _, root := TestCoreUnsealed(t)

	core.logicalBackends["badrenew"] = badRenewFactory
	me := &MountEntry{
		Table:    mountTableType,
		Path:     "badrenew/",
		Type:     "badrenew",
		Accessor: "badrenewaccessor",
	}

	err := core.mount(context.Background(), me)
	if err != nil {
		t.Fatal(err)
	}

	req := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "badrenew/creds",
		ClientToken: root,
	}

	resp, err := core.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("response was nil")
	}
	if resp.Secret == nil {
		t.Fatalf("response secret was nil, response was %#v", *resp)
	}
	leaseID := resp.Secret.LeaseID

	req.Operation = logical.UpdateOperation
	req.Path = "sys/leases/lookup"
	req.Data = map[string]interface{}{"lease_id": leaseID}
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.Data["id"].(string) != leaseID {
		t.Fatalf("expected id %q, got %q", leaseID, resp.Data["id"].(string))
	}

	req.Path = "sys/revoke-prefix/" + leaseID

	resp, err = core.HandleRequest(req)
	if err == nil {
		t.Fatal("expected error")
	}

	req.Path = "sys/revoke-force/" + leaseID
	resp, err = core.HandleRequest(req)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	req.Path = "sys/leases/lookup"
	req.Data = map[string]interface{}{"lease_id": leaseID}
	resp, err = core.HandleRequest(req)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid request") {
		t.Fatalf("bad error: %v", err)
	}
}

func badRenewFactory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	be := &framework.Backend{
		Paths: []*framework.Path{
			&framework.Path{
				Pattern: "creds",
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: func(context.Context, *logical.Request, *framework.FieldData) (*logical.Response, error) {
						resp := &logical.Response{
							Secret: &logical.Secret{
								InternalData: map[string]interface{}{
									"secret_type": "badRenewBackend",
								},
							},
						}
						resp.Secret.TTL = time.Second * 30
						return resp, nil
					},
				},
			},
		},

		Secrets: []*framework.Secret{
			&framework.Secret{
				Type: "badRenewBackend",
				Revoke: func(context.Context, *logical.Request, *framework.FieldData) (*logical.Response, error) {
					return nil, fmt.Errorf("always errors")
				},
			},
		},
	}

	err := be.Setup(context.Background(), conf)
	if err != nil {
		return nil, err
	}

	return be, nil
}
