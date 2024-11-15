// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0

package metadb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/raft-wal/types"
	"go.etcd.io/bbolt"
)

const (
	// FileName is the default file name for the bolt db file.
	FileName = "wal-meta.db"

	// *Bucket are the names used for internal bolt buckets
	MetaBucket   = "wal-meta"
	StableBucket = "stable"

	// We just need one key for now so use the byte 'm' for meta arbitrarily.
	MetaKey = "m"
)

var (
	// ErrUnintialized is returned when any call is made before Load has opened
	// the DB file.
	ErrUnintialized = errors.New("uninitialized")
)

// BoltMetaDB implements types.MetaStore using BoltDB as a reliable persistent
// store. See repo README for reasons for this design choice and performance
// implications.
type BoltMetaDB struct {
	dir string
	db  *bbolt.DB
}

func (db *BoltMetaDB) ensureOpen(dir string) error {
	if db.dir != "" && db.dir != dir {
		return fmt.Errorf("can't load dir %s, already open in dir %s", dir, db.dir)
	}
	if db.db != nil {
		return nil
	}

	fileName := filepath.Join(dir, FileName)

	open := func() error {
		bb, err := bbolt.Open(fileName, 0644, nil)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", FileName, err)
		}
		db.db = bb
		db.dir = dir
		return nil
	}

	// BoltDB can get stuck in invalid states if we crash while it's initializing.
	// We can't distinguish those as safe to just wipe it and start again because
	// we don't know for sure if it's failing due to bad init or later corruption
	// (which would loose data if we just wipe and start over). So to ensure
	// initial creation of the WAL is as crash-safe as possible we will manually
	// detect we have an atomic init procedure:
	//  1. Check if file exits already. If yes, skip init and just open it.
	//  2. Delete any existing DB file with tmp name
	//  3. Creat a new BoltDB that is empty and has the buckets with a temp name.
	//  4. Once that's committed, rename to final name and Fsync parent dir
	_, err := os.Stat(fileName)
	if err == nil {
		// File exists, just open it
		return open()
	}
	if !errors.Is(err, os.ErrNotExist) {
		// Unknown err just return that
		return fmt.Errorf("failed to stat %s: %w", FileName, err)
	}

	// File doesn't exist, initialize a new DB in a crash-safe way
	if err := safeInitBoltDB(dir); err != nil {
		return fmt.Errorf("failed initializing meta DB: %w", err)
	}

	// All good, now open it!
	return open()
}

func safeInitBoltDB(dir string) error {
	tmpFileName := filepath.Join(dir, FileName+".tmp")

	// Delete any old attempts to init that were unsuccessful
	if err := os.RemoveAll(tmpFileName); err != nil {
		return err
	}

	// Open bolt DB at tmp file name
	bb, err := bbolt.Open(tmpFileName, 0644, nil)
	if err != nil {
		return err
	}

	tx, err := bb.Begin(true)
	defer tx.Rollback()

	if err != nil {
		return err
	}
	_, err = tx.CreateBucket([]byte(MetaBucket))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucket([]byte(StableBucket))
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	// Close the file ready to rename into place and re-open. This probably isn't
	// necessary but it make it easier to reason about this code path being
	// totally separate from the common case.
	if err := bb.Close(); err != nil {
		return err
	}

	// We created the DB OK. Now rename it to the final name.
	if err := os.Rename(tmpFileName, filepath.Join(dir, FileName)); err != nil {
		return err
	}

	// And Fsync that parent dir to make sure the new new file with it's new name
	// is persisted!
	dirF, err := os.Open(dir)
	if err != nil {
		return err
	}
	err = dirF.Sync()
	closeErr := dirF.Close()
	if err != nil {
		return err
	}
	return closeErr
}

// Load loads the existing persisted state. If there is no existing state
// implementations are expected to create initialize new storage and return an
// empty state.
func (db *BoltMetaDB) Load(dir string) (types.PersistentState, error) {
	var state types.PersistentState

	if err := db.ensureOpen(dir); err != nil {
		return state, err
	}

	tx, err := db.db.Begin(false)
	if err != nil {
		return state, err
	}
	defer tx.Rollback()
	meta := tx.Bucket([]byte(MetaBucket))

	// We just need one key for now so use the byte 'm' for meta arbitrarily.
	raw := meta.Get([]byte(MetaKey))
	if raw == nil {
		// This is valid it's an "empty" log that will be initialized by the WAL.
		return state, nil
	}

	if err := json.Unmarshal(raw, &state); err != nil {
		return state, fmt.Errorf("%w: failed to parse persisted state: %s", types.ErrCorrupt, err)
	}
	return state, nil
}

// CommitState must atomically replace all persisted metadata in the current
// store with the set provided. It must not return until the data is persisted
// durably and in a crash-safe way otherwise the guarantees of the WAL will be
// compromised. The WAL will only ever call this in a single thread at one
// time and it will never be called concurrently with Load however it may be
// called concurrently with Get/SetStable operations.
func (db *BoltMetaDB) CommitState(state types.PersistentState) error {
	if db.db == nil {
		return ErrUnintialized
	}

	encoded, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to encode persisted state: %w", err)
	}

	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	meta := tx.Bucket([]byte(MetaBucket))

	if err := meta.Put([]byte(MetaKey), encoded); err != nil {
		return err
	}

	return tx.Commit()
}

// GetStable returns a value from stable store or nil if it doesn't exist. May
// be called concurrently by multiple threads.
func (db *BoltMetaDB) GetStable(key []byte) ([]byte, error) {
	if db.db == nil {
		return nil, ErrUnintialized
	}

	tx, err := db.db.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stable := tx.Bucket([]byte(StableBucket))

	val := stable.Get(key)
	if val == nil {
		return nil, nil
	}

	// Need to copy the value since bolt only guarantees the slice is valid until
	// end of txn.
	ret := make([]byte, len(val))
	copy(ret, val)
	return ret, nil
}

// SetStable stores a value from stable store. May be called concurrently with
// GetStable.
func (db *BoltMetaDB) SetStable(key []byte, value []byte) error {
	if db.db == nil {
		return ErrUnintialized
	}

	tx, err := db.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stable := tx.Bucket([]byte(StableBucket))

	if value == nil {
		err = stable.Delete(key)
	} else {
		err = stable.Put(key, value)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Close implements io.Closer
func (db *BoltMetaDB) Close() error {
	if db.db == nil {
		return nil
	}
	err := db.db.Close()
	db.db = nil
	return err
}
