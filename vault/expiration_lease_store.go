package vault

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/sdk/logical"
)

type LeaseStore struct {
	db *memdb.MemDB

	leaseCount int
}

func newLeaseStore() (*LeaseStore, error) {
	db, err := memdb.NewMemDB(schema())
	if err != nil {
		return nil, err
	}

	return &LeaseStore{
		db: db,
	}, nil
}

type memDBLeaseInfo struct {
	LeaseID string

	ExpireTimeUnix int64
	EntityID       string
	Irrevocable    bool

	revokesAttempted uint8
	issueTime        time.Time
	expireTime       time.Time
	lastRenewalTime  time.Time

	// TODO i think we can save some memory if we don't hold the full objects
	// here, even if not everything is populated?
	secret *logical.Secret
	auth   *logical.Auth
}

func newMemDBLeaseInfo(le *leaseEntry) *memDBLeaseInfo {
	ret := &memDBLeaseInfo{
		LeaseID:     le.LeaseID,
		EntityID:    le.EntityID,
		Irrevocable: le.RevokeErr,

		issueTime:       le.IssueTime,
		expireTime:      le.ExpireTime,
		lastRenewalTime: le.LastRenewalTime,
	}

	if !le.ExpireTime.IsZero() {
		// TODO should we use UnixNano?
		ret.ExpireTimeUnix = le.ExpireTime.Unix()
	}

	if le.Secret != nil {
		ret.secret = &logical.Secret{}
		ret.secret.Renewable = le.Secret.Renewable
		ret.secret.TTL = le.Secret.TTL
	}

	// TODO include policies here too
	if le.Auth != nil {
		ret.auth = &logical.Auth{}
		ret.auth.Renewable = le.Auth.Renewable
		ret.auth.TTL = le.Auth.TTL
	}

	return ret
}

func (li *memDBLeaseInfo) Clone() *memDBLeaseInfo {
	return &memDBLeaseInfo{
		LeaseID:          li.LeaseID,
		revokesAttempted: li.revokesAttempted,
		ExpireTimeUnix:   li.ExpireTimeUnix,
		EntityID:         li.EntityID,
		issueTime:        li.issueTime,
		expireTime:       li.expireTime,
		lastRenewalTime:  li.lastRenewalTime,
		secret:           li.secret,
		auth:             li.auth,
		Irrevocable:      li.Irrevocable,
	}
}

func (li *memDBLeaseInfo) ToLease() *leaseEntry {
	ret := &leaseEntry{
		IssueTime:       li.issueTime,
		ExpireTime:      li.expireTime,
		LastRenewalTime: li.lastRenewalTime,
	}
	if li.secret != nil {
		ret.Secret = &logical.Secret{}
		ret.Secret.Renewable = li.secret.Renewable
		ret.Secret.TTL = li.secret.TTL
	}
	if li.auth != nil {
		ret.Auth = &logical.Auth{}
		ret.Auth.Renewable = li.auth.Renewable
		ret.Auth.TTL = li.auth.TTL
	}

	return ret
}

func (ls *LeaseStore) Load(leaseID string) (*memDBLeaseInfo, error) {
	txn := ls.db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("leases", "id", leaseID)
	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, nil
	}

	return raw.(*memDBLeaseInfo), nil
}

func (ls *LeaseStore) Insert(le *leaseEntry) (bool, error) {
	txn := ls.db.Txn(true)
	defer txn.Abort()

	// Check for an existing lease
	raw, err := txn.First("leases", "id", le.LeaseID)
	if err != nil {
		return false, err
	}

	memLe := newMemDBLeaseInfo(le)

	// TODO handle nonexpiring leases
	if err := txn.Insert("leases", memLe); err != nil {
		return false, err
	}

	txn.Commit()
	if raw == nil {
		ls.leaseCount++
	}

	return raw == nil, nil
}

func (ls *LeaseStore) Remove(leaseID string) (bool, error) {
	txn := ls.db.Txn(true)
	defer txn.Abort()

	// Check for an existing lease
	raw, err := txn.First("leases", "id", leaseID)
	if err != nil {
		return false, err
	}

	// Nothing to delete
	if raw == nil {
		return false, nil
	}

	if err := txn.Delete("leases", raw); err != nil {
		return false, err
	}

	txn.Commit()
	ls.leaseCount--

	return true, nil
}

func (ls *LeaseStore) Update(li *memDBLeaseInfo) error {
	txn := ls.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("leases", li); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func (ls *LeaseStore) StartExpirations(ctx context.Context, expireFunc ExpireLeaseStrategy, m *ExpirationManager) {
	expireLeases := func() int {

		txn := ls.db.Txn(true)
		defer txn.Abort()

		currentTime := time.Now().Unix()
		it, err := txn.LowerBound("leases", "expiration", currentTime)
		if err != nil {
			// TODO
			panic(err)
		}

		leasesToExpire := []*memDBLeaseInfo{}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			li := obj.(*memDBLeaseInfo)
			fmt.Printf("  %s is expired: %d\n", li.LeaseID, li.ExpireTimeUnix)

			// Clone and update expire time to 0, this will prevent us from
			// re-capturing the same leases in the next timeframe.
			clone := li.Clone()
			clone.ExpireTimeUnix = 0
			if err := txn.Insert("leases", clone); err != nil {
				// TODO
				panic(err)
			}

			leasesToExpire = append(leasesToExpire, clone)
		}

		txn.Commit()

		for _, li := range leasesToExpire {
			// TODO populate namespace
			// TODO should we run this in a go routine? the current default strategy
			// should be non-blocking...
			go expireFunc(ctx, m, li.LeaseID, nil)
		}

		return len(leasesToExpire)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			expireLeases()
		}
	}
}

func schema() *memdb.DBSchema {
	// Create the DB schema
	return &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"leases": &memdb.TableSchema{
				Name: "leases",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "LeaseID"},
					},
					"irrevocable": {
						Name:    "irrevocable",
						Unique:  false,
						Indexer: &memdb.BoolFieldIndex{Field: "Irrevocable"},
					},
					"entity": {
						Name:         "entity",
						Unique:       false,
						AllowMissing: true,
						Indexer:      &memdb.UUIDFieldIndex{Field: "EntityID"},
					},
					"expiration": {
						Name:         "expiration",
						Unique:       false,
						AllowMissing: true,
						Indexer:      expirationIndexer{},
					},
				},
			},
		},
	}
}

type expirationIndexer struct{}

func (f expirationIndexer) FromArgs(args ...interface{}) ([]byte, error) {
	intArg, ok := args[0].(int64)
	if !ok {
		return nil, errors.New("arg is not int64")
	}

	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, intArg)

	return buf, nil
}

func (f expirationIndexer) FromObject(raw interface{}) (bool, []byte, error) {
	info, ok := raw.(*memDBLeaseInfo)
	if !ok {
		return false, nil, errors.New("invalid type for index")
	}

	// We want to omit non-expiring leases and leases we have already sent to
	// the expireFunc. The expire time will be reset if there are any errors and
	// the lease needs to be retried.
	if info.ExpireTimeUnix <= 0 {
		return false, nil, nil
	}

	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(buf, info.ExpireTimeUnix)

	return true, buf, nil
}
