package vault

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/namespace"
	"go.uber.org/atomic"
)

func Benchmark_MemDB_Insert(b *testing.B) {
	ls, err := newLeaseStore()
	if err != nil {
		b.Fatal(err)
	}
	leases := make([]*leaseEntry, b.N)
	for i := 0; i < b.N; i++ {
		u, err := uuid.GenerateUUID()
		if err != nil {
			b.Fatal(err)
		}

		expireTime := time.Now()
		leases[i] = &leaseEntry{
			LeaseID:     fmt.Sprintf("/auth/userpass/login/%d", i),
			ClientToken: "test-token-id",
			EntityID:    u,
			ExpireTime:  expireTime.Add(time.Duration(i * -1)),
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ls.Insert(leases[i])
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_MemDB_Insert_Txn(b *testing.B) {
	ls, err := newLeaseStore()
	if err != nil {
		b.Fatal(err)
	}
	leases := make([]*leaseEntry, b.N)
	for i := 0; i < b.N; i++ {
		u, err := uuid.GenerateUUID()
		if err != nil {
			b.Fatal(err)
		}

		expireTime := time.Now()
		leases[i] = &leaseEntry{
			LeaseID:     fmt.Sprintf("/auth/userpass/login/%d", i),
			ClientToken: "test-token-id",
			EntityID:    u,
			ExpireTime:  expireTime.Add(time.Duration(i * -1)),
		}
	}

	txn := ls.db.Txn(true)
	defer txn.Abort()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := txn.Insert("leases", newMemDBLeaseInfo(leases[i])); err != nil {
			b.Fatal(err)
		}
	}

	txn.Commit()
}

func Benchmark_MemDB_Load(b *testing.B) {
	ls, err := newLeaseStore()
	if err != nil {
		b.Fatal(err)
	}
	leases := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		u, err := uuid.GenerateUUID()
		if err != nil {
			b.Fatal(err)
		}

		expireTime := time.Now()
		leaseID := fmt.Sprintf("/auth/userpass/login/%d", i)
		_, err = ls.Insert(&leaseEntry{
			LeaseID:     leaseID,
			ClientToken: "test-token-id",
			EntityID:    u,
			ExpireTime:  expireTime.Add(time.Duration(i * -1)),
		})
		if err != nil {
			b.Fatal(err)
		}
		leases[i] = leaseID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val, err := ls.Load(leases[i])
		if err != nil {
			b.Fatal(err)
		}
		if val == nil {
			b.Fatal("nil")
		}
	}
}

func TestLeaseStore_Expirations(t *testing.T) {
	ls, err := newLeaseStore()
	if err != nil {
		t.Fatal(err)
	}

	leaseCount := 1000000

	txn := ls.db.Txn(true)
	defer txn.Abort()

	leases := make([]string, leaseCount)
	for i := 0; i < leaseCount; i++ {
		u, err := uuid.GenerateUUID()
		if err != nil {
			t.Fatal(err)
		}

		expireTime := time.Now()
		leaseID := fmt.Sprintf("%d", i)
		le := &leaseEntry{
			LeaseID:     leaseID,
			ClientToken: "test-token-id",
			EntityID:    u,
			ExpireTime:  expireTime.Add(time.Duration(i*2) * time.Second * -1),
		}
		if err != nil {
			t.Fatal(err)
		}
		leases[i] = leaseID

		// TODO handle nonexpiring leases
		if err := txn.Insert("leases", newMemDBLeaseInfo(le)); err != nil {
			t.Fatal(err)
		}
	}

	txn.Commit()

	eCh := make(chan string, leaseCount)
	pending := atomic.NewUint32(uint32(leaseCount))
	expireFunc := func(ctx context.Context, m *ExpirationManager, leaseID string, ns *namespace.Namespace) {
		eCh <- leaseID
		pending.Sub(1)
	}
	go ls.StartExpirations(context.Background(), expireFunc, nil)

	for {
		if pending.Load() == 0 {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	close(eCh)

	// TODO this test will become flakey if we switch to go routines
	i := leaseCount - 1
	for leaseID := range eCh {
		if leaseID != fmt.Sprintf("%d", i) {
			t.Fatal("bad", leaseID, i)
		}
		i--
	}
}
