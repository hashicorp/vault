// +build !enterprise

package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
)

func (m *ExpirationManager) leaseView(*namespace.Namespace) *BarrierView {
	return m.idView
}

func (m *ExpirationManager) tokenIndexView(*namespace.Namespace) *BarrierView {
	return m.tokenView
}

func (m *ExpirationManager) collectLeases() (map[*namespace.Namespace][]string, int, error) {
	leaseCount := 0
	existing := make(map[*namespace.Namespace][]string)
	keys, err := logical.CollectKeys(m.quitContext, m.leaseView(namespace.RootNamespace))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to scan for leases: %w", err)
	}
	existing[namespace.RootNamespace] = keys
	leaseCount += len(keys)
	return existing, leaseCount, nil
}

// add an irrevocable lease for test purposes
// returns the lease ID for the lease
func addIrrevocableLease(t *testing.T, m *ExpirationManager, pathPrefix string, ns *namespace.Namespace) string {
	t.Helper()

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		t.Fatalf("error generating uuid: %v", err)
	}

	if ns == nil {
		ns = namespace.RootNamespace
	}

	nsSuffix := ""
	if ns != namespace.RootNamespace {
		nsSuffix = fmt.Sprintf("/blah.%s", ns.ID)
	}

	le := &leaseEntry{
		LeaseID:    pathPrefix + "lease" + uuid + nsSuffix,
		Path:       pathPrefix + nsSuffix,
		namespace:  ns,
		IssueTime:  time.Now(),
		ExpireTime: time.Now().Add(time.Hour),
		RevokeErr:  "some error message",
	}

	m.pendingLock.Lock()
	defer m.pendingLock.Unlock()

	if err := m.persistEntry(context.Background(), le); err != nil {
		t.Fatalf("error persisting irrevocable lease: %v", err)
	}

	m.updatePendingInternal(le)

	return le.LeaseID
}

// InjectIrrevocableLeases injects `count` irrevocable leases (currently to a
// single mount).
// It returns a map of the mount accessor to the number of leases stored there
func (c *Core) InjectIrrevocableLeases(t *testing.T, ctx context.Context, count int) map[string]int {
	out := make(map[string]int)
	for i := 0; i < count; i++ {
		leaseID := addIrrevocableLease(t, c.expiration, "foo/", namespace.RootNamespace)

		mountAccessor := c.expiration.getLeaseMountAccessor(ctx, leaseID)
		if _, ok := out[mountAccessor]; !ok {
			out[mountAccessor] = 0
		}

		out[mountAccessor]++
	}

	return out
}
