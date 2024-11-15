package raftboltdb

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/armon/go-metrics"
	v1 "github.com/boltdb/bolt"
	"github.com/hashicorp/raft"
	"go.etcd.io/bbolt"
)

const (
	// Permissions to use on the db file. This is only used if the
	// database file does not exist and needs to be created.
	dbFileMode = 0600
)

var (
	// Bucket names we perform transactions in
	dbLogs = []byte("logs")
	dbConf = []byte("conf")

	// An error indicating a given key does not exist
	ErrKeyNotFound = errors.New("not found")
)

// BoltStore provides access to Bbolt for Raft to store and retrieve
// log entries. It also provides key/value storage, and can be used as
// a LogStore and StableStore.
type BoltStore struct {
	// conn is the underlying handle to the db.
	conn *bbolt.DB

	// The path to the Bolt database file
	path string

	msgpackUseNewTimeFormat bool
}

// Options contains all the configuration used to open the Bbolt
type Options struct {
	// Path is the file path to the Bbolt to use
	Path string

	// BoltOptions contains any specific Bbolt options you might
	// want to specify [e.g. open timeout]
	BoltOptions *bbolt.Options

	// NoSync causes the database to skip fsync calls after each
	// write to the log. This is unsafe, so it should be used
	// with caution.
	NoSync bool

	// MsgpackUseNewTimeFormat when set to true, force the underlying msgpack
	// codec to use the new format of time.Time when encoding (used in
	// go-msgpack v1.1.5 by default). Decoding is not affected, as all
	// go-msgpack v2.1.0+ decoders know how to decode both formats.
	MsgpackUseNewTimeFormat bool
}

// readOnly returns true if the contained bolt options say to open
// the DB in readOnly mode [this can be useful to tools that want
// to examine the log]
func (o *Options) readOnly() bool {
	return o != nil && o.BoltOptions != nil && o.BoltOptions.ReadOnly
}

// NewBoltStore takes a file path and returns a connected Raft backend.
func NewBoltStore(path string) (*BoltStore, error) {
	return New(Options{Path: path})
}

// New uses the supplied options to open the Bbolt and prepare it for use as a raft backend.
func New(options Options) (*BoltStore, error) {
	// Try to connect
	handle, err := bbolt.Open(options.Path, dbFileMode, options.BoltOptions)
	if err != nil {
		return nil, err
	}
	handle.NoSync = options.NoSync

	// Create the new store
	store := &BoltStore{
		conn:                    handle,
		path:                    options.Path,
		msgpackUseNewTimeFormat: options.MsgpackUseNewTimeFormat,
	}

	// If the store was opened read-only, don't try and create buckets
	if !options.readOnly() {
		// Set up our buckets
		if err := store.initialize(); err != nil {
			store.Close()
			return nil, err
		}
	}
	return store, nil
}

// initialize is used to set up all of the buckets.
func (b *BoltStore) initialize() error {
	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create all the buckets
	if _, err := tx.CreateBucketIfNotExists(dbLogs); err != nil {
		return err
	}
	if _, err := tx.CreateBucketIfNotExists(dbConf); err != nil {
		return err
	}

	return tx.Commit()
}

func (b *BoltStore) Stats() bbolt.Stats {
	return b.conn.Stats()
}

// Close is used to gracefully close the DB connection.
func (b *BoltStore) Close() error {
	return b.conn.Close()
}

// FirstIndex returns the first known index from the Raft log.
func (b *BoltStore) FirstIndex() (uint64, error) {
	tx, err := b.conn.Begin(false)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	curs := tx.Bucket(dbLogs).Cursor()
	if first, _ := curs.First(); first == nil {
		return 0, nil
	} else {
		return bytesToUint64(first), nil
	}
}

// LastIndex returns the last known index from the Raft log.
func (b *BoltStore) LastIndex() (uint64, error) {
	tx, err := b.conn.Begin(false)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	curs := tx.Bucket(dbLogs).Cursor()
	if last, _ := curs.Last(); last == nil {
		return 0, nil
	} else {
		return bytesToUint64(last), nil
	}
}

// GetLog is used to retrieve a log from Bbolt at a given index.
func (b *BoltStore) GetLog(idx uint64, log *raft.Log) error {
	defer metrics.MeasureSince([]string{"raft", "boltdb", "getLog"}, time.Now())

	tx, err := b.conn.Begin(false)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbLogs)
	val := bucket.Get(uint64ToBytes(idx))

	if val == nil {
		return raft.ErrLogNotFound
	}
	return decodeMsgPack(val, log)
}

// StoreLog is used to store a single raft log
func (b *BoltStore) StoreLog(log *raft.Log) error {
	return b.StoreLogs([]*raft.Log{log})
}

// StoreLogs is used to store a set of raft logs
func (b *BoltStore) StoreLogs(logs []*raft.Log) error {
	now := time.Now()

	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	batchSize := 0
	for _, log := range logs {
		key := uint64ToBytes(log.Index)
		val, err := encodeMsgPack(log, b.msgpackUseNewTimeFormat)
		if err != nil {
			return err
		}

		logLen := val.Len()
		bucket := tx.Bucket(dbLogs)
		if err := bucket.Put(key, val.Bytes()); err != nil {
			return err
		}
		batchSize += logLen
		metrics.AddSample([]string{"raft", "boltdb", "logSize"}, float32(logLen))
	}

	metrics.AddSample([]string{"raft", "boltdb", "logsPerBatch"}, float32(len(logs)))
	metrics.AddSample([]string{"raft", "boltdb", "logBatchSize"}, float32(batchSize))
	// Both the deferral and the inline function are important for this metrics
	// accuracy. Deferral allows us to calculate the metric after the tx.Commit
	// has finished and thus account for all the processing of the operation.
	// The inlined function ensures that we do not calculate the time.Since(now)
	// at the time of deferral but rather when the go runtime executes the
	// deferred function.
	defer func() {
		metrics.AddSample([]string{"raft", "boltdb", "writeCapacity"}, (float32(1_000_000_000)/float32(time.Since(now).Nanoseconds()))*float32(len(logs)))
		metrics.MeasureSince([]string{"raft", "boltdb", "storeLogs"}, now)
	}()

	return tx.Commit()
}

// DeleteRange is used to delete logs within a given range inclusively.
func (b *BoltStore) DeleteRange(min, max uint64) error {
	minKey := uint64ToBytes(min)

	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	curs := tx.Bucket(dbLogs).Cursor()
	for k, _ := curs.Seek(minKey); k != nil; k, _ = curs.Next() {
		// Handle out-of-range log index
		if bytesToUint64(k) > max {
			break
		}

		// Delete in-range log index
		if err := curs.Delete(); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Set is used to set a key/value set outside of the raft log
func (b *BoltStore) Set(k, v []byte) error {
	tx, err := b.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbConf)
	if err := bucket.Put(k, v); err != nil {
		return err
	}

	return tx.Commit()
}

// Get is used to retrieve a value from the k/v store by key
func (b *BoltStore) Get(k []byte) ([]byte, error) {
	tx, err := b.conn.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket(dbConf)
	val := bucket.Get(k)

	if val == nil {
		return nil, ErrKeyNotFound
	}
	return append([]byte(nil), val...), nil
}

// SetUint64 is like Set, but handles uint64 values
func (b *BoltStore) SetUint64(key []byte, val uint64) error {
	return b.Set(key, uint64ToBytes(val))
}

// GetUint64 is like Get, but handles uint64 values
func (b *BoltStore) GetUint64(key []byte) (uint64, error) {
	val, err := b.Get(key)
	if err != nil {
		return 0, err
	}
	return bytesToUint64(val), nil
}

// Sync performs an fsync on the database file handle. This is not necessary
// under normal operation unless NoSync is enabled, in which this forces the
// database file to sync against the disk.
func (b *BoltStore) Sync() error {
	return b.conn.Sync()
}

// MigrateToV2 reads in the source file path of a BoltDB file
// and outputs all the data migrated to a Bbolt destination file
func MigrateToV2(source, destination string) (*BoltStore, error) {
	_, err := os.Stat(destination)
	if err == nil {
		return nil, fmt.Errorf("file exists in destination %v", destination)
	}

	srcDb, err := v1.Open(source, dbFileMode, &v1.Options{
		ReadOnly: true,
		Timeout:  1 * time.Minute,
	})
	if err != nil {
		return nil, fmt.Errorf("failed opening source database: %v", err)
	}

	//Start a connection to the source
	srctx, err := srcDb.Begin(false)
	if err != nil {
		return nil, fmt.Errorf("failed connecting to source database: %v", err)
	}
	defer srctx.Rollback()

	//Create the destination
	destDb, err := New(Options{Path: destination})
	if err != nil {
		return nil, fmt.Errorf("failed creating destination database: %v", err)
	}
	//Start a connection to the new
	desttx, err := destDb.conn.Begin(true)
	if err != nil {
		destDb.Close()
		os.Remove(destination)
		return nil, fmt.Errorf("failed connecting to destination database: %v", err)
	}

	defer desttx.Rollback()

	//Loop over both old buckets and set them in the new
	buckets := [][]byte{dbConf, dbLogs}
	for _, b := range buckets {
		srcB := srctx.Bucket(b)
		destB := desttx.Bucket(b)
		err = srcB.ForEach(func(k, v []byte) error {
			return destB.Put(k, v)
		})
		if err != nil {
			destDb.Close()
			os.Remove(destination)
			return nil, fmt.Errorf("failed to copy %v bucket: %v", string(b), err)
		}
	}

	//If the commit fails, clean up
	if err := desttx.Commit(); err != nil {
		destDb.Close()
		os.Remove(destination)
		return nil, fmt.Errorf("failed commiting data to destination: %v", err)
	}

	return destDb, nil

}
