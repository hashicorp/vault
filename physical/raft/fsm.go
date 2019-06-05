package raft

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	metrics "github.com/armon/go-metrics"
	protoio "github.com/gogo/protobuf/io"
	proto "github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	bolt "go.etcd.io/bbolt"
)

const (
	deleteOp uint32 = 1 << iota
	putOp
)

var (
	// dataBucketName is the value we use for the bucket
	dataBucketName   = []byte("data")
	configBucketName = []byte("config")
	latestIndexKey   = []byte("latest_indexes")
	latestConfigKey  = []byte("latest_config")
)

// Verify FSM satisfies the correct interfaces
var _ physical.Backend = (*FSM)(nil)
var _ physical.Transactional = (*FSM)(nil)
var _ raft.FSM = (*FSM)(nil)

// FSMApplyResponse is returned from an FSM apply. It indicates if the apply was
// successful or not.
type FSMApplyResponse struct {
	Success bool
}

// FSM is Vault's primary state storage. It writes updates to an bolt db file
// that lives on local disk. FSM implements raft.FSM and physical.Backend
// interfaces.
type FSM struct {
	l           sync.RWMutex
	path        string
	logger      log.Logger
	permitPool  *physical.PermitPool
	noopRestore bool

	db *bolt.DB

	// latestIndex and latestTerm are the term and index of the last log we
	// received
	latestIndex  *uint64
	latestTerm   *uint64
	latestConfig atomic.Value

	// This is just used in tests to disable to storing the latest indexes and
	// configs so we can conform to the standard backend tests, which expect to
	// additional state in the backend.
	storeLatestState bool
}

// NewFSM constructs a FSM using the given directory
func NewFSM(conf map[string]string, logger log.Logger) (*FSM, error) {
	path, ok := conf["path"]
	if !ok {
		return nil, fmt.Errorf("'path' must be set")
	}

	dbPath := filepath.Join(path, "vault.db")

	boltDB, err := bolt.Open(dbPath, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	// Initialize the latest term, index, and config values
	latestTerm := new(uint64)
	latestIndex := new(uint64)
	latestConfig := atomic.Value{}
	atomic.StoreUint64(latestTerm, 0)
	atomic.StoreUint64(latestIndex, 0)
	latestConfig.Store((*ConfigurationValue)(nil))

	err = boltDB.Update(func(tx *bolt.Tx) error {
		// make sure we have the necessary buckets created
		_, err := tx.CreateBucketIfNotExists(dataBucketName)
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		b, err := tx.CreateBucketIfNotExists(configBucketName)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		// Read in our latest index and term and populate it inmemory
		val := b.Get(latestIndexKey)
		if val != nil {
			var latest IndexValue
			err := proto.Unmarshal(val, &latest)
			if err != nil {
				return err
			}

			atomic.StoreUint64(latestTerm, latest.Term)
			atomic.StoreUint64(latestIndex, latest.Index)
		}

		// Read in our latest config and populate it inmemory
		val = b.Get(latestConfigKey)
		if val != nil {
			var latest ConfigurationValue
			err := proto.Unmarshal(val, &latest)
			if err != nil {
				return err
			}

			latestConfig.Store(&latest)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	storeLatestState := true
	if _, ok := conf["doNotStoreLatestState"]; ok {
		storeLatestState = false
	}

	return &FSM{
		path:       path,
		logger:     logger,
		permitPool: physical.NewPermitPool(physical.DefaultParallelOperations),

		db:               boltDB,
		latestTerm:       latestTerm,
		latestIndex:      latestIndex,
		latestConfig:     latestConfig,
		storeLatestState: storeLatestState,
	}, nil
}

// LatestState returns the latest index and configuration values we have seen on
// this FSM.
func (f *FSM) LatestState() (*IndexValue, *ConfigurationValue) {
	return &IndexValue{
		Term:  atomic.LoadUint64(f.latestTerm),
		Index: atomic.LoadUint64(f.latestIndex),
	}, f.latestConfig.Load().(*ConfigurationValue)
}

// Delete deletes the given key from the bolt file.
func (f *FSM) Delete(ctx context.Context, path string) error {
	defer metrics.MeasureSince([]string{"raft", "delete"}, time.Now())

	f.permitPool.Acquire()
	defer f.permitPool.Release()

	f.l.RLock()
	defer f.l.RUnlock()

	return f.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(dataBucketName).Delete([]byte(path))
	})
}

// Get retrieves the value at the given path from the bolt file.
func (f *FSM) Get(ctx context.Context, path string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"raft", "get"}, time.Now())

	f.permitPool.Acquire()
	defer f.permitPool.Release()

	f.l.RLock()
	defer f.l.RUnlock()

	var valCopy []byte
	var found bool

	err := f.db.View(func(tx *bolt.Tx) error {

		value := tx.Bucket(dataBucketName).Get([]byte(path))
		if value != nil {
			found = true
			valCopy = make([]byte, len(value))
			copy(valCopy, value)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	return &physical.Entry{
		Key:   path,
		Value: valCopy,
	}, nil
}

// Put writes the given entry to the bolt file.
func (f *FSM) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"raft", "put"}, time.Now())

	f.permitPool.Acquire()
	defer f.permitPool.Release()

	f.l.RLock()
	defer f.l.RUnlock()

	// Start a write transaction.
	return f.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(dataBucketName).Put([]byte(entry.Key), entry.Value)
	})
}

// List retrieves the set of keys with the given prefix from the bolt file.
func (f *FSM) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"raft", "list"}, time.Now())

	f.permitPool.Acquire()
	defer f.permitPool.Release()

	f.l.RLock()
	defer f.l.RUnlock()

	var keys []string

	err := f.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(dataBucketName).Cursor()

		prefixBytes := []byte(prefix)
		for k, _ := c.Seek(prefixBytes); k != nil && bytes.HasPrefix(k, prefixBytes); k, _ = c.Next() {
			key := string(k)
			key = strings.TrimPrefix(key, prefix)
			if i := strings.Index(key, "/"); i == -1 {
				// Add objects only from the current 'folder'
				keys = append(keys, key)
			} else if i != -1 {
				// Add truncated 'folder' paths
				keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
			}
		}

		return nil
	})

	return keys, err
}

// Transaction writes all the operations in the provided transaction to the bolt
// file.
func (f *FSM) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	f.permitPool.Acquire()
	defer f.permitPool.Release()

	f.l.RLock()
	defer f.l.RUnlock()

	// TODO: should this be a Batch?
	// Start a write transaction.
	err := f.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(dataBucketName)
		for _, txn := range txns {
			var err error
			switch txn.Operation {
			case physical.PutOperation:
				err = b.Put([]byte(txn.Entry.Key), txn.Entry.Value)
			case physical.DeleteOperation:
				err = b.Delete([]byte(txn.Entry.Key))
			default:
				return fmt.Errorf("%q is not a supported transaction operation", txn.Operation)
			}
			if err != nil {
				return err
			}
		}

		return nil
	})
	return err
}

// Apply will apply a log value to the FSM. This is called from the raft
// library.
func (f *FSM) Apply(log *raft.Log) interface{} {
	command := &LogData{}
	err := proto.Unmarshal(log.Data, command)
	if err != nil {
		panic("error proto unmarshaling log data")
	}

	f.l.RLock()
	defer f.l.RUnlock()

	index, err := proto.Marshal(&IndexValue{
		Term:  log.Term,
		Index: log.Index,
	})
	if err != nil {
		panic("failed to store data")
	}

	err = f.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(dataBucketName)
		for _, op := range command.Operations {
			var err error
			switch op.OpType {
			case putOp:
				err = b.Put([]byte(op.Key), op.Value)
			case deleteOp:
				err = b.Delete([]byte(op.Key))
			default:
				return fmt.Errorf("%q is not a supported transaction operation", op.OpType)
			}
			if err != nil {
				return err
			}
		}

		// TODO: benchmark so we can know how much time this adds
		if f.storeLatestState {
			err = b.Put(latestIndexKey, index)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println("failed to store data")
		panic("failed to store data")
	}

	atomic.StoreUint64(f.latestTerm, log.Term)
	atomic.StoreUint64(f.latestIndex, log.Index)

	return &FSMApplyResponse{
		Success: true,
	}
}

type writeErrorCloser interface {
	io.WriteCloser
	CloseWithError(error) error
}

// writeTo will copy the FSM's content to a remote sink. The data is written
// twice, once for use in determining various metadata attributes of the dataset
// (size, checksum, etc) and a second for the sink of the data. We also use a
// proto delimited writer so we can stream proto messages to the sink.
func (f *FSM) writeTo(ctx context.Context, metaSink writeErrorCloser, sink writeErrorCloser) {
	protoWriter := protoio.NewDelimitedWriter(sink)
	metadataProtoWriter := protoio.NewDelimitedWriter(metaSink)

	f.l.RLock()
	defer f.l.RUnlock()

	err := f.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dataBucketName)

		c := b.Cursor()

		// Do the first scan of the data for metadata purposes.
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := metadataProtoWriter.WriteMsg(&pb.StorageEntry{
				Key:   string(k),
				Value: v,
			})
			if err != nil {
				metaSink.CloseWithError(err)
				return err
			}
		}
		metaSink.Close()

		// Do the second scan for copy purposes.
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := protoWriter.WriteMsg(&pb.StorageEntry{
				Key:   string(k),
				Value: v,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	sink.CloseWithError(err)
}

// Snapshot implements the FSM interface. It returns a noop snapshot object.
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &noopSnapshotter{}, nil
}

// SetNoopRestore is used to disable restore operations on raft startup. Because
// we are using persistent storage in our FSM we do not need to issue a restore
// on startup.
func (f *FSM) SetNoopRestore(enabled bool) {
	f.l.Lock()
	f.noopRestore = enabled
	f.l.Unlock()
}

// Restore reads data from the provided reader and writes it into the FSM. It
// first deletes the existing bucket to clear all existing data, then recreates
// it so we can copy in the snapshot.
func (f *FSM) Restore(r io.ReadCloser) error {
	if f.noopRestore == true {
		return nil
	}

	protoReader := protoio.NewDelimitedReader(r, math.MaxInt64)
	defer protoReader.Close()

	f.l.Lock()
	defer f.l.Unlock()

	// Start a write transaction.
	err := f.db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(dataBucketName)
		if err != nil {
			return err
		}

		b, err := tx.CreateBucket(dataBucketName)
		if err != nil {
			return err
		}

		for {
			s := new(pb.StorageEntry)
			err := protoReader.ReadMsg(s)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
			err = b.Put([]byte(s.Key), s.Value)
			if err != nil {
				return err
			}

			// Update in-memory representation of the latest states
			switch s.Key {
			case string(latestIndexKey):
				var latest IndexValue
				err := proto.Unmarshal(s.Value, &latest)
				if err != nil {
					return err
				}

				atomic.StoreUint64(f.latestTerm, latest.Term)
				atomic.StoreUint64(f.latestIndex, latest.Index)
			case string(latestConfigKey):
				var latest ConfigurationValue
				err := proto.Unmarshal(s.Value, &latest)
				if err != nil {
					return err
				}

				f.latestConfig.Store(&latest)
			}

		}

		return nil
	})
	if err != nil {
		f.logger.Error("could not restore snapshot", "error", err)
		return err
	}

	return nil
}

// noopSnapshotter implements the fsm.Snapshot interface. It doesn't do anything
// since our SnapshotStore reads data out of the FSM on Open().
type noopSnapshotter struct{}

// Persist doesn't do anything.
func (s *noopSnapshotter) Persist(sink raft.SnapshotSink) error {
	return nil
}

// Release doesn't do anything.
func (s *noopSnapshotter) Release() {}

// StoreConfig satisfies the raft.ConfigurationStore interface and persists the
// latest raft server configuration to the bolt file.
func (f *FSM) StoreConfig(index uint64, configuration raft.Configuration) error {
	f.l.RLock()
	defer f.l.RUnlock()

	latestIndex, _ := f.LatestState()
	latestIndex.Index = index
	indexBytes, err := proto.Marshal(latestIndex)
	if err != nil {
		fmt.Println(err)
		return err
	}

	protoConfig := raftConfigurationToProtoConfiguration(index, configuration)
	configBytes, err := proto.Marshal(protoConfig)
	if err != nil {
		return err
	}

	if f.storeLatestState {
		err = f.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(dataBucketName)
			err := b.Put(latestConfigKey, configBytes)
			if err != nil {
				return err
			}

			// TODO: benchmark so we can know how much time this adds
			err = b.Put(latestIndexKey, indexBytes)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	atomic.StoreUint64(f.latestIndex, index)
	f.latestConfig.Store(protoConfig)

	return nil
}

// raftConfigurationToProtoConfiguration converts a raft configuration object to
// a proto value.
func raftConfigurationToProtoConfiguration(index uint64, configuration raft.Configuration) *ConfigurationValue {
	servers := make([]*Server, len(configuration.Servers))
	for i, s := range configuration.Servers {
		servers[i] = &Server{
			Suffrage: int32(s.Suffrage),
			Id:       string(s.ID),
			Address:  string(s.Address),
		}
	}
	return &ConfigurationValue{
		Index:   index,
		Servers: servers,
	}
}

// protoConfigurationToRaftConfiguration converts a proto configuration object
// to a raft object.
func protoConfigurationToRaftConfiguration(configuration *ConfigurationValue) (uint64, raft.Configuration) {
	servers := make([]raft.Server, len(configuration.Servers))
	for i, s := range configuration.Servers {
		servers[i] = raft.Server{
			Suffrage: raft.ServerSuffrage(s.Suffrage),
			ID:       raft.ServerID(s.Id),
			Address:  raft.ServerAddress(s.Address),
		}
	}
	return configuration.Index, raft.Configuration{
		Servers: servers,
	}
}
