package vault

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
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

	nsSuffix := ""
	if ns != namespace.RootNamespace {
		nsSuffix = fmt.Sprintf("/blah.%s", ns.ID)
	}

	randomTimeDelta := time.Duration(rand.Int31n(24))
	le := &leaseEntry{
		LeaseID:    pathPrefix + "lease" + uuid + nsSuffix,
		Path:       pathPrefix + nsSuffix,
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
