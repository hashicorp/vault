// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cachememdb

import (
	"errors"
	"fmt"
	"sync/atomic"

	memdb "github.com/hashicorp/go-memdb"
)

const (
	tableNameIndexer             = "indexer"
	tableNameCapabilitiesIndexer = "capabilities-indexer"
)

// ErrCacheItemNotFound is returned on Get and GetCapabilitiesIndex calls
// when the entry is not found in the cache.
var ErrCacheItemNotFound = errors.New("cache item not found")

// CacheMemDB is the underlying cache database for storing indexes.
type CacheMemDB struct {
	db *atomic.Value
}

// New creates a new instance of CacheMemDB.
func New() (*CacheMemDB, error) {
	db, err := newDB()
	if err != nil {
		return nil, err
	}

	c := &CacheMemDB{
		db: new(atomic.Value),
	}
	c.db.Store(db)

	return c, nil
}

func newDB() (*memdb.MemDB, error) {
	cacheSchema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableNameIndexer: {
				Name: tableNameIndexer,
				Indexes: map[string]*memdb.IndexSchema{
					// This index enables fetching the cached item based on the
					// identifier of the index.
					IndexNameID: {
						Name:   IndexNameID,
						Unique: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "ID",
						},
					},
					// This index enables fetching all the entries in cache for
					// a given request path, in a given namespace.
					IndexNameRequestPath: {
						Name:   IndexNameRequestPath,
						Unique: false,
						Indexer: &memdb.CompoundIndex{
							Indexes: []memdb.Indexer{
								&memdb.StringFieldIndex{
									Field: "Namespace",
								},
								&memdb.StringFieldIndex{
									Field: "RequestPath",
								},
							},
						},
					},
					// This index enables fetching all the entries in cache
					// belonging to the leases of a given token.
					IndexNameLeaseToken: {
						Name:         IndexNameLeaseToken,
						Unique:       false,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "LeaseToken",
						},
					},
					// This index enables fetching all the entries in cache
					// that are tied to the given token, regardless of the
					// entries belonging to the token or belonging to the
					// lease.
					IndexNameToken: {
						Name:         IndexNameToken,
						Unique:       true,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "Token",
						},
					},
					// This index enables fetching all the entries in cache for
					// the given parent token.
					IndexNameTokenParent: {
						Name:         IndexNameTokenParent,
						Unique:       false,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "TokenParent",
						},
					},
					// This index enables fetching all the entries in cache for
					// the given accessor.
					IndexNameTokenAccessor: {
						Name:         IndexNameTokenAccessor,
						Unique:       true,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "TokenAccessor",
						},
					},
					// This index enables fetching all the entries in cache for
					// the given lease identifier.
					IndexNameLease: {
						Name:         IndexNameLease,
						Unique:       true,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "Lease",
						},
					},
				},
			},
			tableNameCapabilitiesIndexer: {
				Name: tableNameCapabilitiesIndexer,
				Indexes: map[string]*memdb.IndexSchema{
					// This index enables fetching the cached item based on the
					// identifier of the index.
					CapabilitiesIndexNameID: {
						Name:   CapabilitiesIndexNameID,
						Unique: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "ID",
						},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(cacheSchema)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Get returns the index based on the indexer and the index values provided.
// If the capabilities index isn't present, it will return nil, ErrCacheItemNotFound
func (c *CacheMemDB) Get(indexName string, indexValues ...interface{}) (*Index, error) {
	if !validIndexName(indexName) {
		return nil, fmt.Errorf("invalid index name %q", indexName)
	}

	txn := c.db.Load().(*memdb.MemDB).Txn(false)

	raw, err := txn.First(tableNameIndexer, indexName, indexValues...)
	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, ErrCacheItemNotFound
	}

	index, ok := raw.(*Index)
	if !ok {
		return nil, errors.New("unable to parse index value from the cache")
	}

	return index, nil
}

// Set stores the index into the cache.
func (c *CacheMemDB) Set(index *Index) error {
	if index == nil {
		return errors.New("nil index provided")
	}

	txn := c.db.Load().(*memdb.MemDB).Txn(true)
	defer txn.Abort()

	if err := txn.Insert(tableNameIndexer, index); err != nil {
		return fmt.Errorf("unable to insert index into cache: %v", err)
	}

	txn.Commit()

	return nil
}

// GetCapabilitiesIndex returns the CapabilitiesIndex from the cache.
// If the capabilities index isn't present, it will return nil, ErrCacheItemNotFound
func (c *CacheMemDB) GetCapabilitiesIndex(indexName string, indexValues ...interface{}) (*CapabilitiesIndex, error) {
	if !validCapabilitiesIndexName(indexName) {
		return nil, fmt.Errorf("invalid index name %q", indexName)
	}

	txn := c.db.Load().(*memdb.MemDB).Txn(false)

	raw, err := txn.First(tableNameCapabilitiesIndexer, indexName, indexValues...)
	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, ErrCacheItemNotFound
	}

	index, ok := raw.(*CapabilitiesIndex)
	if !ok {
		return nil, errors.New("unable to parse capabilities index value from the cache")
	}

	return index, nil
}

// SetCapabilitiesIndex stores the CapabilitiesIndex index into the cache.
func (c *CacheMemDB) SetCapabilitiesIndex(index *CapabilitiesIndex) error {
	if index == nil {
		return errors.New("nil capabilities index provided")
	}

	txn := c.db.Load().(*memdb.MemDB).Txn(true)
	defer txn.Abort()

	if err := txn.Insert(tableNameCapabilitiesIndexer, index); err != nil {
		return fmt.Errorf("unable to insert index into cache: %v", err)
	}

	txn.Commit()

	return nil
}

// EvictCapabilitiesIndex removes a capabilities index from the cache based on index name and value.
func (c *CacheMemDB) EvictCapabilitiesIndex(indexName string, indexValues ...interface{}) error {
	index, err := c.GetCapabilitiesIndex(indexName, indexValues...)
	if errors.Is(err, ErrCacheItemNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("unable to fetch index on cache deletion: %v", err)
	}

	txn := c.db.Load().(*memdb.MemDB).Txn(true)
	defer txn.Abort()

	if err := txn.Delete(tableNameCapabilitiesIndexer, index); err != nil {
		return fmt.Errorf("unable to delete index from cache: %v", err)
	}

	txn.Commit()

	return nil
}

// GetByPrefix returns all the cached indexes based on the index name and the
// value prefix.
func (c *CacheMemDB) GetByPrefix(indexName string, indexValues ...interface{}) ([]*Index, error) {
	if !validIndexName(indexName) {
		return nil, fmt.Errorf("invalid index name %q", indexName)
	}

	indexName = indexName + "_prefix"

	// Get all the objects
	txn := c.db.Load().(*memdb.MemDB).Txn(false)

	iter, err := txn.Get(tableNameIndexer, indexName, indexValues...)
	if err != nil {
		return nil, err
	}

	var indexes []*Index
	for {
		obj := iter.Next()
		if obj == nil {
			break
		}
		index, ok := obj.(*Index)
		if !ok {
			return nil, fmt.Errorf("failed to cast cached index")
		}

		indexes = append(indexes, index)
	}

	return indexes, nil
}

// Evict removes an index from the cache based on index name and value.
func (c *CacheMemDB) Evict(indexName string, indexValues ...interface{}) error {
	index, err := c.Get(indexName, indexValues...)
	if errors.Is(err, ErrCacheItemNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("unable to fetch index on cache deletion: %v", err)
	}

	txn := c.db.Load().(*memdb.MemDB).Txn(true)
	defer txn.Abort()

	if err := txn.Delete(tableNameIndexer, index); err != nil {
		return fmt.Errorf("unable to delete index from cache: %v", err)
	}

	txn.Commit()

	return nil
}

// Flush resets the underlying cache object.
func (c *CacheMemDB) Flush() error {
	newDB, err := newDB()
	if err != nil {
		return err
	}

	c.db.Store(newDB)

	return nil
}
