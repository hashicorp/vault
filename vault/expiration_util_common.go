package vault

import (
	"context"
	"fmt"
	"math/rand"
	"path"
	"testing"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

type basicLeaseTestInfo struct {
	id     string
	mount  string
	expire time.Time
}

// add an irrevocable lease for test purposes
// returns the lease ID and expire time
func addIrrevocableLease(t *testing.T, m *ExpirationManager, pathPrefix string, ns *namespace.Namespace) *basicLeaseTestInfo {
	t.Helper()

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("error generating uuid: %v", err)
	}

	if ns == nil {
		ns = namespace.RootNamespace
	}

	leaseID := path.Join(pathPrefix, "lease"+uuid)

	if ns != namespace.RootNamespace {
		leaseID = fmt.Sprintf("%s.%s", leaseID, ns.ID)
	}

	randomTimeDelta := time.Duration(rand.Int31n(24))
	le := &leaseEntry{
		LeaseID:    leaseID,
		Path:       pathPrefix,
		namespace:  ns,
		IssueTime:  time.Now(),
		ExpireTime: time.Now().Add(randomTimeDelta * time.Hour),
		RevokeErr:  "some error message",
	}

	m.pendingLock.Lock()
	defer m.pendingLock.Unlock()

	if err := m.persistEntry(context.Background(), le); err != nil {
		t.Fatalf("error persisting irrevocable lease: %v", err)
	}

	m.updatePendingInternal(le)

	return &basicLeaseTestInfo{
		id:     le.LeaseID,
		expire: le.ExpireTime,
	}
}

// InjectIrrevocableLeases injects `count` irrevocable leases (currently to a
// single mount).
// It returns a map of the mount accessor to the number of leases stored there
func (c *Core) InjectIrrevocableLeases(t *testing.T, ctx context.Context, count int) map[string]int {
	out := make(map[string]int)
	for i := 0; i < count; i++ {
		le := addIrrevocableLease(t, c.expiration, "foo/", namespace.RootNamespace)

		mountAccessor := c.expiration.getLeaseMountAccessor(ctx, le.id)
		if _, ok := out[mountAccessor]; !ok {
			out[mountAccessor] = 0
		}

		out[mountAccessor]++
	}

	return out
}

type backend struct {
	path string
	ns   *namespace.Namespace
}

// set up multiple mounts, and return a mapping of the path to the mount accessor
func mountNoopBackends(t *testing.T, c *Core, backends []*backend) map[string]string {
	t.Helper()

	// enable the noop backend
	c.logicalBackends["noop"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		return &NoopBackend{}, nil
	}

	pathToMount := make(map[string]string)
	for _, backend := range backends {
		me := &MountEntry{
			Table: mountTableType,
			Path:  backend.path,
			Type:  "noop",
		}

		nsCtx := namespace.ContextWithNamespace(context.Background(), backend.ns)
		err := c.mount(nsCtx, me)
		if err != nil {
			t.Fatalf("err mounting backend %s: %v", backend.path, err)
		}

		mount := c.router.MatchingMountEntry(nsCtx, backend.path)
		if mount == nil {
			t.Fatalf("couldn't find mount for path %s", backend.path)
		}
		pathToMount[backend.path] = mount.Accessor
	}

	return pathToMount
}
