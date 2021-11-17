package vault

import (
	"context"
	"fmt"
	"math/rand"
	"path"
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
func (c *Core) AddIrrevocableLease(ctx context.Context, pathPrefix string) (*basicLeaseTestInfo, error) {
	exp := c.expiration

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating uuid: %w", err)
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting namespace from context: %w", err)
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

	exp.pendingLock.Lock()
	defer exp.pendingLock.Unlock()

	if err := exp.persistEntry(context.Background(), le); err != nil {
		return nil, fmt.Errorf("error persisting irrevocable lease: %w", err)
	}

	exp.updatePendingInternal(le)

	return &basicLeaseTestInfo{
		id:     le.LeaseID,
		expire: le.ExpireTime,
	}, nil
}

// InjectIrrevocableLeases injects `count` irrevocable leases (currently to a
// single mount).
// It returns a map of the mount accessor to the number of leases stored there
func (c *Core) InjectIrrevocableLeases(ctx context.Context, count int) (map[string]int, error) {
	out := make(map[string]int)
	for i := 0; i < count; i++ {
		le, err := c.AddIrrevocableLease(ctx, "foo/")
		if err != nil {
			return nil, err
		}

		mountAccessor := c.expiration.getLeaseMountAccessor(ctx, le.id)
		if _, ok := out[mountAccessor]; !ok {
			out[mountAccessor] = 0
		}

		out[mountAccessor]++
	}

	return out, nil
}

type backend struct {
	path string
	ns   *namespace.Namespace
}

// set up multiple mounts, and return a mapping of the path to the mount accessor
func mountNoopBackends(c *Core, backends []*backend) (map[string]string, error) {
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
		if err := c.mount(nsCtx, me); err != nil {
			return nil, fmt.Errorf("error mounting backend %s: %w", backend.path, err)
		}

		mount := c.router.MatchingMountEntry(nsCtx, backend.path)
		if mount == nil {
			return nil, fmt.Errorf("couldn't find mount for path %s", backend.path)
		}
		pathToMount[backend.path] = mount.Accessor
	}

	return pathToMount, nil
}

func (c *Core) FetchLeaseCountToRevoke() int {
	c.expiration.pendingLock.RLock()
	defer c.expiration.pendingLock.RUnlock()
	return c.expiration.leaseCount
}
