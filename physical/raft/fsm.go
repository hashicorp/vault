package raft

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	metrics "github.com/armon/go-metrics"
	protoio "github.com/gogo/protobuf/io"
	proto "github.com/golang/protobuf/proto"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-raftchunking"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
	bolt "go.etcd.io/bbolt"
)

const (
	deleteOp uint32 = 1 << iota
	putOp
	restoreCallbackOp

	chunkingPrefix = "raftchunking/"
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
var _ raft.ConfigurationStore = (*FSM)(nil)

type restoreCallback func(context.Context) error

// FSMApplyResponse is returned from an FSM apply. It indicates if the apply was
// successful or not.
type FSMApplyResponse struct {
	Success bool
}

// FSM is Vault's primary state storage. It writes updates to an bolt db file
// that lives on local disk. FSM implements raft.FSM and physical.Backend
// interfaces.
type FSM struct {
	// latestIndex and latestTerm must stay at the top of this struct to be
	// properly 64-bit aligned.

	// latestIndex and latestTerm are the term and index of the last log we
	// received
	latestIndex *uint64
	latestTerm  *uint64
	// latestConfig is the latest server configuration we've seen
	latestConfig atomic.Value

	l           sync.RWMutex
	path        string
	logger      log.Logger
	permitPool  *physical.PermitPool
	noopRestore bool

	db *bolt.DB

	// retoreCb is called after we've restored a snapshot
	restoreCb restoreCallback

	// This is just used in tests to disable to storing the latest indexes and
	// configs so we can conform to the standard backend tests, which expect to
	// additional state in the backend.
	storeLatestState bool

	chunker *raftchunking.ChunkingConfigurationStore
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
			return fmt.Errorf("failed to create bucket: %v", err)
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

	f := &FSM{
		path:       conf["path"],
		logger:     logger,
		permitPool: physical.NewPermitPool(physical.DefaultParallelOperations),

		db:               boltDB,
		latestTerm:       latestTerm,
		latestIndex:      latestIndex,
		latestConfig:     latestConfig,
		storeLatestState: storeLatestState,
	}

	f.chunker = raftchunking.NewChunkingConfigurationStore(f, &FSMChunkStorage{
		f:   f,
		ctx: context.Background(),
	})

	return f, nil
}

// LatestState returns the latest index and configuration values we have seen on
// this FSM.
func (f *FSM) LatestState() (*IndexValue, *ConfigurationValue) {
	return &IndexValue{
		Term:  atomic.LoadUint64(f.latestTerm),
		Index: atomic.LoadUint64(f.latestIndex),
	}, f.latestConfig.Load().(*ConfigurationValue)
}

func (f *FSM) witnessIndex(i *IndexValue) {
	seen, _ := f.LatestState()
	if seen.Index < i.Index {
		atomic.StoreUint64(f.latestIndex, i.Index)
		atomic.StoreUint64(f.latestTerm, i.Term)
	}
}

func (f *FSM) witnessSnapshot(index, term, configurationIndex uint64, configuration raft.Configuration) error {
	var indexBytes []byte
	latestIndex, _ := f.LatestState()

	latestIndex.Index = index
	latestIndex.Term = term

	var err error
	indexBytes, err = proto.Marshal(latestIndex)
	if err != nil {
		return err
	}

	protoConfig := raftConfigurationToProtoConfiguration(configurationIndex, configuration)
	configBytes, err := proto.Marshal(protoConfig)
	if err != nil {
		return err
	}

	if f.storeLatestState {
		err = f.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(configBucketName)
			err := b.Put(latestConfigKey, configBytes)
			if err != nil {
				return err
			}

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
	atomic.StoreUint64(f.latestTerm, term)
	f.latestConfig.Store(protoConfig)

	return nil
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

// Delete deletes the given key from the bolt file.
func (f *FSM) DeletePrefix(ctx context.Context, prefix string) error {
	defer metrics.MeasureSince([]string{"raft", "delete_prefix"}, time.Now())

	f.permitPool.Acquire()
	defer f.permitPool.Release()

	f.l.RLock()
	defer f.l.RUnlock()

	err := f.db.Update(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket(dataBucketName).Cursor()

		prefixBytes := []byte(prefix)
		for k, _ := c.Seek(prefixBytes); k != nil && bytes.HasPrefix(k, prefixBytes); k, _ = c.Next() {
			if err := c.Delete(); err != nil {
				return err
			}
		}

		return nil
	})

	return err
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
			} else {
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
		f.logger.Error("error proto unmarshaling log data", "error", err)
		panic("error proto unmarshaling log data")
	}

	f.l.RLock()
	defer f.l.RUnlock()

	// Only advance latest pointer if this log has a higher index value than
	// what we have seen in the past.
	var logIndex []byte
	latestIndex, _ := f.LatestState()
	if latestIndex.Index < log.Index {
		logIndex, err = proto.Marshal(&IndexValue{
			Term:  log.Term,
			Index: log.Index,
		})
		if err != nil {
			f.logger.Error("unable to marshal latest index", "error", err)
			panic("unable to marshal latest index")
		}
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
			case restoreCallbackOp:
				if f.restoreCb != nil {
					// Kick off the restore callback function in a go routine
					go f.restoreCb(context.Background())
				}
			default:
				return fmt.Errorf("%q is not a supported transaction operation", op.OpType)
			}
			if err != nil {
				return err
			}
		}

		// TODO: benchmark so we can know how much time this adds
		if f.storeLatestState && len(logIndex) > 0 {
			b := tx.Bucket(configBucketName)
			err = b.Put(latestIndexKey, logIndex)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		f.logger.Error("failed to store data", "error", err)
		panic("failed to store data")
	}

	// If we advanced the latest value, update the in-memory representation too.
	if len(logIndex) > 0 {
		atomic.StoreUint64(f.latestTerm, log.Term)
		atomic.StoreUint64(f.latestIndex, log.Index)
	}

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

	protoReader := protoio.NewDelimitedReader(r, math.MaxInt32)
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
func (f *FSM) StoreConfiguration(index uint64, configuration raft.Configuration) {
	f.l.RLock()
	defer f.l.RUnlock()

	var indexBytes []byte
	latestIndex, _ := f.LatestState()
	// Only write the new index if we are advancing the pointer
	if index > latestIndex.Index {
		latestIndex.Index = index

		var err error
		indexBytes, err = proto.Marshal(latestIndex)
		if err != nil {
			f.logger.Error("unable to marshal latest index", "error", err)
			panic(fmt.Sprintf("unable to marshal latest index: %v", err))
		}
	}

	protoConfig := raftConfigurationToProtoConfiguration(index, configuration)
	configBytes, err := proto.Marshal(protoConfig)
	if err != nil {
		f.logger.Error("unable to marshal config", "error", err)
		panic(fmt.Sprintf("unable to marshal config: %v", err))
	}

	if f.storeLatestState {
		err = f.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(configBucketName)
			err := b.Put(latestConfigKey, configBytes)
			if err != nil {
				return err
			}

			// TODO: benchmark so we can know how much time this adds
			if len(indexBytes) > 0 {
				err = b.Put(latestIndexKey, indexBytes)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			f.logger.Error("unable to store latest configuration", "error", err)
			panic(fmt.Sprintf("unable to store latest configuration: %v", err))
		}
	}

	f.witnessIndex(latestIndex)
	f.latestConfig.Store(protoConfig)
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

type FSMChunkStorage struct {
	f   *FSM
	ctx context.Context
}

// chunkPaths returns a disk prefix and key given chunkinfo
func (f *FSMChunkStorage) chunkPaths(chunk *raftchunking.ChunkInfo) (string, string) {
	prefix := fmt.Sprintf("%s%d/", chunkingPrefix, chunk.OpNum)
	key := fmt.Sprintf("%s%d", prefix, chunk.SequenceNum)
	return prefix, key
}

func (f *FSMChunkStorage) StoreChunk(chunk *raftchunking.ChunkInfo) (bool, error) {
	b, err := jsonutil.EncodeJSON(chunk)
	if err != nil {
		return false, errwrap.Wrapf("error encoding chunk info: {{err}}", err)
	}

	prefix, key := f.chunkPaths(chunk)

	entry := &physical.Entry{
		Key:   key,
		Value: b,
	}

	f.f.permitPool.Acquire()
	defer f.f.permitPool.Release()

	f.f.l.RLock()
	defer f.f.l.RUnlock()

	// Start a write transaction.
	done := new(bool)
	if err := f.f.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(dataBucketName).Put([]byte(entry.Key), entry.Value); err != nil {
			return errwrap.Wrapf("error storing chunk info: {{err}}", err)
		}

		// Assume bucket exists and has keys
		c := tx.Bucket(dataBucketName).Cursor()

		var keys []string
		prefixBytes := []byte(prefix)
		for k, _ := c.Seek(prefixBytes); k != nil && bytes.HasPrefix(k, prefixBytes); k, _ = c.Next() {
			key := string(k)
			key = strings.TrimPrefix(key, prefix)
			if i := strings.Index(key, "/"); i == -1 {
				// Add objects only from the current 'folder'
				keys = append(keys, key)
			} else {
				// Add truncated 'folder' paths
				keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
			}
		}

		*done = uint32(len(keys)) == chunk.NumChunks

		return nil
	}); err != nil {
		return false, err
	}

	return *done, nil
}

func (f *FSMChunkStorage) FinalizeOp(opNum uint64) ([]*raftchunking.ChunkInfo, error) {
	ret, err := f.chunksForOpNum(opNum)
	if err != nil {
		return nil, errwrap.Wrapf("error getting chunks for op keys: {{err}}", err)
	}

	prefix, _ := f.chunkPaths(&raftchunking.ChunkInfo{OpNum: opNum})
	if err := f.f.DeletePrefix(f.ctx, prefix); err != nil {
		return nil, errwrap.Wrapf("error deleting prefix after op finalization: {{err}}", err)
	}

	return ret, nil
}

func (f *FSMChunkStorage) chunksForOpNum(opNum uint64) ([]*raftchunking.ChunkInfo, error) {
	prefix, _ := f.chunkPaths(&raftchunking.ChunkInfo{OpNum: opNum})

	opChunkKeys, err := f.f.List(f.ctx, prefix)
	if err != nil {
		return nil, errwrap.Wrapf("error fetching op chunk keys: {{err}}", err)
	}

	if len(opChunkKeys) == 0 {
		return nil, nil
	}

	var ret []*raftchunking.ChunkInfo

	for _, v := range opChunkKeys {
		seqNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, errwrap.Wrapf("error converting seqnum to integer: {{err}}", err)
		}

		entry, err := f.f.Get(f.ctx, prefix+v)
		if err != nil {
			return nil, errwrap.Wrapf("error fetching chunkinfo: {{err}}", err)
		}

		var ci raftchunking.ChunkInfo
		if err := jsonutil.DecodeJSON(entry.Value, &ci); err != nil {
			return nil, errwrap.Wrapf("error decoding chunkinfo json: {{err}}", err)
		}

		if ret == nil {
			ret = make([]*raftchunking.ChunkInfo, ci.NumChunks)
		}

		ret[seqNum] = &ci
	}

	return ret, nil
}

func (f *FSMChunkStorage) GetChunks() (raftchunking.ChunkMap, error) {
	opNums, err := f.f.List(f.ctx, chunkingPrefix)
	if err != nil {
		return nil, errwrap.Wrapf("error doing recursive list for chunk saving: {{err}}", err)
	}

	if len(opNums) == 0 {
		return nil, nil
	}

	ret := make(raftchunking.ChunkMap, len(opNums))
	for _, opNumStr := range opNums {
		opNum, err := strconv.ParseInt(opNumStr, 10, 64)
		if err != nil {
			return nil, errwrap.Wrapf("error parsing op num during chunk saving: {{err}}", err)
		}

		opChunks, err := f.chunksForOpNum(uint64(opNum))
		if err != nil {
			return nil, errwrap.Wrapf("error getting chunks for op keys during chunk saving: {{err}}", err)
		}

		ret[uint64(opNum)] = opChunks
	}

	return ret, nil
}

func (f *FSMChunkStorage) RestoreChunks(chunks raftchunking.ChunkMap) error {
	if err := f.f.DeletePrefix(f.ctx, chunkingPrefix); err != nil {
		return errwrap.Wrapf("error deleting prefix for chunk restoration: {{err}}", err)
	}
	if len(chunks) == 0 {
		return nil
	}

	for opNum, opChunks := range chunks {
		for _, chunk := range opChunks {
			if chunk == nil {
				continue
			}
			if chunk.OpNum != opNum {
				return errors.New("unexpected op number in chunk")
			}
			if _, err := f.StoreChunk(chunk); err != nil {
				return errwrap.Wrapf("error storing chunk during restoration: {{err}}", err)
			}
		}
	}

	return nil
}
