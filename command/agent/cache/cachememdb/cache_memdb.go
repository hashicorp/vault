package cachememdb

import (
	"errors"
	"fmt"

	memdb "github.com/hashicorp/go-memdb"
)

const (
	tableNameIndexer = "indexer"
)

// CacheMemDB is the underlying cache database for storing indexes.
type CacheMemDB struct {
	db *memdb.MemDB
}

// New creates a new instance of CacheMemDB.
func New() (*CacheMemDB, error) {
	db, err := newDB()
	if err != nil {
		return nil, err
	}

	return &CacheMemDB{
		db: db,
	}, nil
}

func newDB() (*memdb.MemDB, error) {
	cacheSchema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableNameIndexer: &memdb.TableSchema{
				Name: tableNameIndexer,
				Indexes: map[string]*memdb.IndexSchema{
					// This index enables fetching the cached item based on the
					// identifier of the index.
					IndexNameID: &memdb.IndexSchema{
						Name:   IndexNameID,
						Unique: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "ID",
						},
					},
					// This index enables fetching all the entries in cache for
					// a given request path, in a given namespace.
					IndexNameRequestPath: &memdb.IndexSchema{
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
					IndexNameLeaseToken: &memdb.IndexSchema{
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
					IndexNameToken: &memdb.IndexSchema{
						Name:         IndexNameToken,
						Unique:       true,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "Token",
						},
					},
					// This index enables fetching all the entries in cache for
					// the given parent token.
					IndexNameTokenParent: &memdb.IndexSchema{
						Name:         IndexNameTokenParent,
						Unique:       false,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "TokenParent",
						},
					},
					// This index enables fetching all the entries in cache for
					// the given accessor.
					IndexNameTokenAccessor: &memdb.IndexSchema{
						Name:         IndexNameTokenAccessor,
						Unique:       true,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "TokenAccessor",
						},
					},
					// This index enables fetching all the entries in cache for
					// the given lease identifier.
					IndexNameLease: &memdb.IndexSchema{
						Name:         IndexNameLease,
						Unique:       true,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{
							Field: "Lease",
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
func (c *CacheMemDB) Get(indexName string, indexValues ...interface{}) (*Index, error) {
	if !validIndexName(indexName) {
		return nil, fmt.Errorf("invalid index name %q", indexName)
	}

	raw, err := c.db.Txn(false).First(tableNameIndexer, indexName, indexValues...)
	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, nil
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

	txn := c.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert(tableNameIndexer, index); err != nil {
		return fmt.Errorf("unable to insert index into cache: %v", err)
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
	iter, err := c.db.Txn(false).Get(tableNameIndexer, indexName, indexValues...)
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
	if err != nil {
		return fmt.Errorf("unable to fetch index on cache deletion: %v", err)
	}

	if index == nil {
		return nil
	}

	txn := c.db.Txn(true)
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

	c.db = newDB

	return nil
}
