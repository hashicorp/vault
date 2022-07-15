package raft

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/cockroachdb/pebble"
	"github.com/golang/protobuf/proto"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-raftchunking"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/plugin/pb"
)

const (
	deleteOp uint32 = 1 << iota
	putOp
	restoreCallbackOp
	chunkingPrefix        = "raftchunking/"
	databaseDirectoryName = "vault-db"
)

var (
	// dataBucketPrefix is the value we use for the bucket
	dataBucketPrefix   = []byte("data/")
	configBucketPrefix = []byte("config/")
	latestIndexKey     = []byte("latest_indexes")
	latestConfigKey    = []byte("latest_config")
	localNodeConfigKey = []byte("local_node_config")
	pebbleWriteOptions = &pebble.WriteOptions{Sync: true}
)

// Verify FSM satisfies the correct interfaces
var (
	_ physical.Backend       = (*FSM)(nil)
	_ physical.Transactional = (*FSM)(nil)
	_ raft.FSM               = (*FSM)(nil)
	_ raft.BatchingFSM       = (*FSM)(nil)
)

type restoreCallback func(context.Context) error

// FSMApplyResponse is returned from an FSM apply. It indicates if the apply was
// successful or not.
type FSMApplyResponse struct {
	Success bool
}

// FSM is Vault's primary state storage. It writes updates to a pebble file
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
	noopRestore bool

	// applyCallback is used to control the pace of applies in tests
	applyCallback func()

	db *pebble.DB

	// retoreCb is called after we've restored a snapshot
	restoreCb restoreCallback

	chunker *raftchunking.ChunkingBatchingFSM

	localID         string
	desiredSuffrage string
}

// NewFSM constructs a FSM using the given directory
func NewFSM(path string, localID string, logger log.Logger) (*FSM, error) {
	// Initialize the latest term, index, and config values
	latestTerm := new(uint64)
	latestIndex := new(uint64)
	latestConfig := atomic.Value{}
	atomic.StoreUint64(latestTerm, 0)
	atomic.StoreUint64(latestIndex, 0)
	latestConfig.Store((*ConfigurationValue)(nil))

	f := &FSM{
		path:   path,
		logger: logger,

		latestTerm:   latestTerm,
		latestIndex:  latestIndex,
		latestConfig: latestConfig,
		// Assume that the default intent is to join as as voter. This will be updated
		// when this node joins a cluster with a different suffrage, or during cluster
		// setup if this is already part of a cluster with a desired suffrage.
		desiredSuffrage: "voter",
		localID:         localID,
	}

	f.chunker = raftchunking.NewChunkingBatchingFSM(f, &FSMChunkStorage{
		f:   f,
		ctx: context.Background(),
	})

	dbPath := filepath.Join(path, databaseDirectoryName)
	if err := f.openDBFile(dbPath); err != nil {
		return nil, fmt.Errorf("failed to open pebble file: %w", err)
	}

	return f, nil
}

func (f *FSM) getDB() *pebble.DB {
	f.l.RLock()
	defer f.l.RUnlock()

	return f.db
}

// SetFSMDelay adds a delay to the FSM apply. This is used in tests to simulate
// a slow apply.
func (r *RaftBackend) SetFSMDelay(delay time.Duration) {
	r.SetFSMApplyCallback(func() { time.Sleep(delay) })
}

func (r *RaftBackend) SetFSMApplyCallback(f func()) {
	r.fsm.l.Lock()
	r.fsm.applyCallback = f
	r.fsm.l.Unlock()
}

func (f *FSM) openDBFile(dbPath string) error {
	if len(dbPath) == 0 {
		return errors.New("can not open empty filename")
	}

	st, err := os.Stat(dbPath)
	switch {
	case err != nil && os.IsNotExist(err):
	case err != nil:
		return fmt.Errorf("error checking raft FSM db file %q: %v", dbPath, err)
	default:
		perms := st.Mode() & os.ModePerm
		if perms&0o077 != 0 {
			f.logger.Warn("raft FSM db directory has wider permissions than needed",
				"needed", os.FileMode(0o600), "existing", perms)
		}
	}

	start := time.Now()
	pebbleDB, err := pebble.Open(dbPath, nil)
	if err != nil {
		return err
	}
	elapsed := time.Now().Sub(start)
	f.logger.Debug("time to open database", "elapsed", elapsed, "path", dbPath)
	metrics.MeasureSince([]string{"raft_storage", "fsm", "open_db_file"}, start)

	configKey := append(configBucketPrefix, latestConfigKey...)
	indexKey := append(configBucketPrefix, latestIndexKey...)

	val, closer, err := pebbleDB.Get(indexKey)
	defer func(c io.Closer) {
		if c != nil {
			_ = c.Close()
		}
	}(closer)
	if err != nil && err != pebble.ErrNotFound {
		return err
	}

	if val != nil {
		var latest IndexValue
		newVal := make([]byte, len(val))
		copy(newVal, val)
		err := proto.Unmarshal(newVal, &latest)
		if err != nil {
			return err
		}

		f.logger.Trace("-- openDBFile", "latestIndexKey.Term", latest.Term, "latestIndexKey.Index", latest.Index)

		atomic.StoreUint64(f.latestTerm, latest.Term)
		atomic.StoreUint64(f.latestIndex, latest.Index)
	}

	val, closer, err = pebbleDB.Get(configKey)
	defer func(c io.Closer) {
		if c != nil {
			_ = c.Close()
		}
	}(closer)
	if err != nil && err != pebble.ErrNotFound {
		return err
	}
	if val != nil {
		var latest ConfigurationValue
		newVal := make([]byte, len(val))
		copy(newVal, val)
		err := proto.Unmarshal(newVal, &latest)
		if err != nil {
			return err
		}

		f.logger.Trace("-- openDBFile", "latestConfigurationValue.Index", latest.Index, "latestIndexKey.Servers", latest.Servers)
		f.latestConfig.Store(&latest)
	}

	f.db = pebbleDB
	return nil
}

func (f *FSM) Close() error {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.db.Close()
}

func writeSnapshotMetaToDB(metadata *raft.SnapshotMeta, db *pebble.DB) error {
	latestIndex := &IndexValue{
		Term:  metadata.Term,
		Index: metadata.Index,
	}
	indexBytes, err := proto.Marshal(latestIndex)
	if err != nil {
		return err
	}

	protoConfig := raftConfigurationToProtoConfiguration(metadata.ConfigurationIndex, metadata.Configuration)
	configBytes, err := proto.Marshal(protoConfig)
	if err != nil {
		return err
	}

	batch := db.NewBatch()
	defer func() { _ = batch.Close() }()

	err = batch.Set(append(configBucketPrefix, latestConfigKey...), configBytes, pebbleWriteOptions)
	if err != nil {
		return err
	}
	err = batch.Set(append(configBucketPrefix, latestIndexKey...), indexBytes, pebbleWriteOptions)
	if err != nil {
		return err
	}

	return batch.Commit(pebbleWriteOptions)
}

func (f *FSM) localNodeConfig() (*LocalNodeConfigValue, error) {
	var configBytes []byte
	key := append(configBucketPrefix, localNodeConfigKey...)

	snap := f.db.NewSnapshot()
	defer func() { _ = snap.Close() }()

	value, closer, err := snap.Get(key)
	defer func() {
		if closer != nil {
			_ = closer.Close()
		}
	}()
	if err != nil && err != pebble.ErrNotFound {
		return nil, err
	}

	if value != nil {
		configBytes = make([]byte, len(value))
		copy(configBytes, value)
	}
	if configBytes == nil {
		return nil, nil
	}

	var lnConfig LocalNodeConfigValue
	if configBytes != nil {
		err := proto.Unmarshal(configBytes, &lnConfig)
		if err != nil {
			return nil, err
		}
		f.desiredSuffrage = lnConfig.DesiredSuffrage
		return &lnConfig, nil
	}

	return nil, nil
}

func (f *FSM) DesiredSuffrage() string {
	f.l.RLock()
	defer f.l.RUnlock()
	return f.desiredSuffrage
}

func (f *FSM) upgradeLocalNodeConfig() error {
	f.l.Lock()
	defer f.l.Unlock()

	// Read the local node config
	lnConfig, err := f.localNodeConfig()
	if err != nil {
		return err
	}

	// Entry is already present. Get the suffrage value.
	if lnConfig != nil {
		f.desiredSuffrage = lnConfig.DesiredSuffrage
		return nil
	}

	//
	// This is the upgrade case where there is no entry
	//

	lnConfig = &LocalNodeConfigValue{}

	// Refer to the persisted latest raft config
	config := f.latestConfig.Load().(*ConfigurationValue)

	// If there is no config, then this is a fresh node coming up. This could end up
	// being a voter or non-voter. But by default assume that this is a voter. It
	// will be changed if this node joins the cluster as a non-voter.
	if config == nil {
		f.desiredSuffrage = "voter"
		lnConfig.DesiredSuffrage = f.desiredSuffrage
		return f.persistDesiredSuffrage(lnConfig)
	}

	// Get the last known suffrage of the node and assume that it is the desired
	// suffrage. There is no better alternative here.
	for _, srv := range config.Servers {
		if srv.Id == f.localID {
			switch srv.Suffrage {
			case int32(raft.Nonvoter):
				lnConfig.DesiredSuffrage = "non-voter"
			default:
				lnConfig.DesiredSuffrage = "voter"
			}
			// Bring the intent to the fsm instance.
			f.desiredSuffrage = lnConfig.DesiredSuffrage
			break
		}
	}

	return f.persistDesiredSuffrage(lnConfig)
}

// recordSuffrage is called when a node successfully joins the cluster. This
// intent should land in the stored configuration. If the config isn't available
// yet, we still go ahead and store the intent in the fsm. During the next
// update to the configuration, this intent will be persisted.
func (f *FSM) recordSuffrage(desiredSuffrage string) error {
	f.l.Lock()
	defer f.l.Unlock()

	if err := f.persistDesiredSuffrage(&LocalNodeConfigValue{
		DesiredSuffrage: desiredSuffrage,
	}); err != nil {
		return err
	}

	f.desiredSuffrage = desiredSuffrage
	return nil
}

func (f *FSM) persistDesiredSuffrage(lnconfig *LocalNodeConfigValue) error {
	dsBytes, err := proto.Marshal(lnconfig)
	if err != nil {
		return err
	}

	key := append(configBucketPrefix, localNodeConfigKey...)
	return f.db.Set(key, dsBytes, pebbleWriteOptions)
}

func (f *FSM) witnessSnapshot(metadata *raft.SnapshotMeta) error {
	f.l.RLock()
	defer f.l.RUnlock()

	err := writeSnapshotMetaToDB(metadata, f.db)
	if err != nil {
		return err
	}

	atomic.StoreUint64(f.latestIndex, metadata.Index)
	atomic.StoreUint64(f.latestTerm, metadata.Term)
	f.latestConfig.Store(raftConfigurationToProtoConfiguration(metadata.ConfigurationIndex, metadata.Configuration))

	return nil
}

// LatestState returns the latest index and configuration values we have seen on
// this FSM.
func (f *FSM) LatestState() (*IndexValue, *ConfigurationValue) {
	return &IndexValue{
		Term:  atomic.LoadUint64(f.latestTerm),
		Index: atomic.LoadUint64(f.latestIndex),
	}, f.latestConfig.Load().(*ConfigurationValue)
}

// Delete deletes the given key from the pebble db.
func (f *FSM) Delete(ctx context.Context, path string) error {
	defer metrics.MeasureSince([]string{"raft_storage", "fsm", "delete"}, time.Now())

	f.l.RLock()
	defer f.l.RUnlock()

	key := append(dataBucketPrefix, []byte(path)...)
	return f.db.Delete(key, nil)
}

func (f *FSM) DeletePrefix(ctx context.Context, prefix string) error {
	defer metrics.MeasureSince([]string{"raft_storage", "fsm", "delete_prefix"}, time.Now())

	f.l.RLock()
	defer f.l.RUnlock()

	batch := f.db.NewIndexedBatch()
	defer func() { _ = batch.Close() }()

	iter := batch.NewIter(nil)
	defer func() { _ = iter.Close() }()

	prefixBytes := append(dataBucketPrefix, []byte(prefix)...)

	for iter.SeekGE(prefixBytes); iter.Valid(); iter.Next() {
		k := iter.Key()

		if bytes.HasPrefix(k, prefixBytes) {
			if err := batch.Delete(k, nil); err != nil {
				return err
			}
		}
	}

	return batch.Commit(pebbleWriteOptions)
}

// Get retrieves the value at the given path from the pebble database.
func (f *FSM) Get(ctx context.Context, path string) (*physical.Entry, error) {
	// TODO: Remove this outdated metric name in an older release
	defer metrics.MeasureSince([]string{"raft", "get"}, time.Now())
	defer metrics.MeasureSince([]string{"raft_storage", "fsm", "get"}, time.Now())

	f.l.RLock()
	defer f.l.RUnlock()

	var valCopy []byte
	var found bool
	key := append(dataBucketPrefix, []byte(path)...)

	snap := f.db.NewSnapshot()
	defer func() { _ = snap.Close() }()

	val, closer, err := snap.Get(key)
	defer func() {
		if closer != nil {
			_ = closer.Close()
		}
	}()

	if err != nil && err != pebble.ErrNotFound {
		return nil, err
	}

	if val != nil {
		found = true
		valCopy = make([]byte, len(val))
		copy(valCopy, val)
	}

	if !found {
		return nil, nil
	}

	return &physical.Entry{
		Key:   path,
		Value: valCopy,
	}, nil
}

// Put writes the given entry to the pebble database.
func (f *FSM) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"raft_storage", "fsm", "put"}, time.Now())

	f.l.RLock()
	defer f.l.RUnlock()

	key := append(dataBucketPrefix, []byte(entry.Key)...)
	return f.db.Set(key, entry.Value, pebbleWriteOptions)
}

// List retrieves the set of keys with the given prefix from the pebble database.
func (f *FSM) List(ctx context.Context, prefix string) ([]string, error) {
	// TODO: Remove this outdated metric name in a future release
	defer metrics.MeasureSince([]string{"raft", "list"}, time.Now())
	defer metrics.MeasureSince([]string{"raft_storage", "fsm", "list"}, time.Now())

	f.l.RLock()
	defer f.l.RUnlock()

	var keys []string
	prefixBytes := append(dataBucketPrefix, []byte(prefix)...)

	batch := f.db.NewIndexedBatch()
	defer func() { _ = batch.Close() }()

	iter := batch.NewIter(nil)
	defer func() { _ = iter.Close() }()

	for iter.SeekGE(prefixBytes); iter.Valid(); iter.Next() {
		k := iter.Key()

		if bytes.HasPrefix(k, prefixBytes) {
			key := string(k)
			key = strings.TrimPrefix(key, string(prefixBytes))
			if i := strings.Index(key, "/"); i == -1 {
				// Add objects only from the current 'folder'
				keys = append(keys, key)
			} else {
				// Add truncated 'folder' paths
				if len(keys) == 0 || keys[len(keys)-1] != key[:i+1] {
					keys = append(keys, string(key[:i+1]))
				}
			}
		}
	}

	return keys, nil
}

// Transaction writes all the operations in the provided transaction to the pebble database
func (f *FSM) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	f.l.RLock()
	defer f.l.RUnlock()

	// Start a write transaction.
	batch := f.db.NewBatch()
	defer func() { _ = batch.Close() }()

	for _, txn := range txns {
		var err error
		key := append(dataBucketPrefix, []byte(txn.Entry.Key)...)

		switch txn.Operation {
		case physical.PutOperation:
			err = batch.Set(key, txn.Entry.Value, pebbleWriteOptions)
		case physical.DeleteOperation:
			err = batch.Delete(key, pebbleWriteOptions)
		default:
			return fmt.Errorf("%q is not a supported transaction operation", txn.Operation)
		}
		if err != nil {
			return err
		}
	}

	return batch.Commit(pebbleWriteOptions)
}

// ApplyBatch will apply a set of logs to the FSM. This is called from the raft
// library.
func (f *FSM) ApplyBatch(logs []*raft.Log) []interface{} {
	if len(logs) == 0 {
		return []interface{}{}
	}

	// Do the unmarshalling first so we don't hold locks
	var latestConfiguration *ConfigurationValue
	commands := make([]interface{}, 0, len(logs))
	for _, log := range logs {
		switch log.Type {
		case raft.LogCommand:
			command := &LogData{}
			err := proto.Unmarshal(log.Data, command)
			if err != nil {
				f.logger.Error("error proto unmarshaling log data", "error", err)
				panic("error proto unmarshaling log data")
			}
			commands = append(commands, command)
		case raft.LogConfiguration:
			configuration := raft.DecodeConfiguration(log.Data)
			config := raftConfigurationToProtoConfiguration(log.Index, configuration)

			commands = append(commands, config)

			// Update the latest configuration the fsm has received; we will
			// store this after it has been committed to storage.
			latestConfiguration = config

		default:
			panic(fmt.Sprintf("got unexpected log type: %d", log.Type))
		}
	}

	// Only advance latest pointer if this log has a higher index value than
	// what we have seen in the past.
	var logIndex []byte
	var err error
	latestIndex, _ := f.LatestState()
	lastLog := logs[len(logs)-1]
	if latestIndex.Index < lastLog.Index {
		logIndex, err = proto.Marshal(&IndexValue{
			Term:  lastLog.Term,
			Index: lastLog.Index,
		})
		if err != nil {
			f.logger.Error("unable to marshal latest index", "error", err)
			panic("unable to marshal latest index")
		}
	}

	f.l.RLock()
	defer f.l.RUnlock()

	if f.applyCallback != nil {
		f.applyCallback()
	}

	batch := f.db.NewBatch()
	defer func() { _ = batch.Close() }()

	for _, commandRaw := range commands {
		switch command := commandRaw.(type) {
		case *LogData:
			for _, op := range command.Operations {
				key := append(dataBucketPrefix, []byte(op.Key)...)
				switch op.OpType {
				case putOp:
					err = batch.Set(key, op.Value, pebbleWriteOptions)
				case deleteOp:
					err = batch.Delete(key, pebbleWriteOptions)
				case restoreCallbackOp:
					if f.restoreCb != nil {
						// Kick off the restore callback function in a go routine
						go f.restoreCb(context.Background())
					}
				default:
					panic(fmt.Errorf("%q is not a supported transaction operation", op.OpType))
				}
			}

		case *ConfigurationValue:
			key := append(configBucketPrefix, latestConfigKey...)
			configBytes, err := proto.Marshal(command)
			if err != nil {
				break
			}
			if err = batch.Set(key, configBytes, pebbleWriteOptions); err != nil {
				break
			}
		}
	}

	if len(logIndex) > 0 {
		key := append(configBucketPrefix, latestIndexKey...)
		err = batch.Set(key, logIndex, pebbleWriteOptions)
	}

	if err != nil {
		f.logger.Error("failed to store data", "error", err)
		panic("failed to store data")
	}

	err = batch.Commit(pebbleWriteOptions)
	if err != nil {
		f.logger.Error("failed to commit batch", "error", err)
		panic("failed to commit batch")
	}

	// If we advanced the latest value, update the in-memory representation too.
	if len(logIndex) > 0 {
		atomic.StoreUint64(f.latestTerm, lastLog.Term)
		atomic.StoreUint64(f.latestIndex, lastLog.Index)
	}

	// If one or more configuration changes were processed, store the latest one.
	if latestConfiguration != nil {
		f.latestConfig.Store(latestConfiguration)
	}

	// Build the responses. The logs array is used here to ensure we reply to
	// all command values; even if they are not of the types we expect. This
	// should future proof this function from more log types being provided.
	resp := make([]interface{}, len(logs))
	for i := range logs {
		resp[i] = &FSMApplyResponse{
			Success: true,
		}
	}

	return resp
}

// Apply will apply a log value to the FSM. This is called from the raft
// library.
func (f *FSM) Apply(log *raft.Log) interface{} {
	return f.ApplyBatch([]*raft.Log{log})[0]
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
	defer metrics.MeasureSince([]string{"raft_storage", "fsm", "write_snapshot"}, time.Now())

	protoWriter := NewDelimitedWriter(sink)
	metadataProtoWriter := NewDelimitedWriter(metaSink)

	f.l.RLock()
	defer f.l.RUnlock()

	// snapshot the db first, so we get a consistent view of the data during this process
	snap := f.db.NewSnapshot()
	defer func() { _ = snap.Close() }()

	var err error

	metaIter := snap.NewIter(&pebble.IterOptions{
		LowerBound: dataBucketPrefix,
	})
	copyIter, err := metaIter.Clone(pebble.CloneOptions{})
	if err != nil {
		f.logger.Error("error cloning iterator", "error", err)
	}

	// Do the first scan of the data for metadata purposes.
	for metaIter.First(); metaIter.Valid(); metaIter.Next() {
		k := metaIter.Key()
		if k == nil {
			continue
		}
		realKey := bytes.TrimPrefix(k, dataBucketPrefix)
		err = metadataProtoWriter.WriteMsg(&pb.StorageEntry{
			Key:   string(realKey),
			Value: metaIter.Value(),
		})
		if err != nil {
			metaSink.CloseWithError(err)
			break
		}
	}
	_ = metaSink.Close()
	_ = metaIter.Close()

	// Do the second scan for copy purposes.
	for copyIter.First(); copyIter.Valid(); copyIter.Next() {
		k := copyIter.Key()
		if k == nil {
			continue
		}
		realKey := bytes.TrimPrefix(k, dataBucketPrefix)
		err = protoWriter.WriteMsg(&pb.StorageEntry{
			Key:   string(realKey),
			Value: copyIter.Value(),
		})
		if err != nil {
			break
		}
	}
	_ = sink.CloseWithError(err)
	_ = copyIter.Close()
}

// Snapshot implements the FSM interface. It returns a noop snapshot object.
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &noopSnapshotter{
		fsm: f,
	}, nil
}

// SetNoopRestore is used to disable restore operations on raft startup. Because
// we are using persistent storage in our FSM we do not need to issue a restore
// on startup.
func (f *FSM) SetNoopRestore(enabled bool) {
	f.l.Lock()
	f.noopRestore = enabled
	f.l.Unlock()
}

// Restore installs a new snapshot from the provided reader. It does an atomic
// rename of the snapshot directory into the database filepath. While a restore is
// happening the FSM is locked and no writes or reads can be performed.
func (f *FSM) Restore(r io.ReadCloser) error {
	defer metrics.MeasureSince([]string{"raft_storage", "fsm", "restore_snapshot"}, time.Now())

	if f.noopRestore {
		return nil
	}

	snapshotInstaller, ok := r.(*pebbleSnapshotInstaller)
	if !ok {
		wrapper, ok := r.(raft.ReadCloserWrapper)
		if !ok {
			return fmt.Errorf("expected ReadCloserWrapper object, got: %T", r)
		}
		snapshotInstallerRaw := wrapper.WrappedReadCloser()
		snapshotInstaller, ok = snapshotInstallerRaw.(*pebbleSnapshotInstaller)
		if !ok {
			return fmt.Errorf("expected snapshot installer object, got: %T", snapshotInstallerRaw)
		}
	}

	f.l.Lock()
	defer f.l.Unlock()

	// Cache the local node config before closing the db
	lnConfig, err := f.localNodeConfig()
	if err != nil {
		return err
	}
	f.logger.Trace("local node config from existing db", "config", lnConfig)

	// Close the db
	if err := f.db.Close(); err != nil {
		f.logger.Error("failed to close database file", "error", err)
		return err
	}

	dbPath := filepath.Join(f.path, databaseDirectoryName)
	f.logger.Info("installing snapshot to FSM")

	// Install the new pebble database
	var retErr *multierror.Error
	f.logger.Trace("existing dbpath", "dbpath", dbPath)
	f.logger.Trace("new dbpath", "dbpath", snapshotInstaller.Dirname())
	if err := snapshotInstaller.Install(dbPath); err != nil {
		f.logger.Error("failed to install snapshot", "error", err)
		retErr = multierror.Append(retErr, fmt.Errorf("failed to install snapshot database: %w", err))
	} else {
		f.logger.Info("snapshot installed")
	}

	// Open the db. We want to do this regardless of if the above install
	// worked. If the install failed we should try to open the old db.
	if err := f.openDBFile(dbPath); err != nil {
		f.logger.Error("failed to open new database", "error", err)
		retErr = multierror.Append(retErr, fmt.Errorf("failed to open new pebble database: %w", err))
	}

	// Handle local node config restore. lnConfig should not be nil here, but
	// adding the nil check anyways for safety.
	if lnConfig != nil {
		// Persist the local node config on the restored fsm.
		if err := f.persistDesiredSuffrage(lnConfig); err != nil {
			f.logger.Error("failed to persist local node config from before the restore", "error", err)
			retErr = multierror.Append(retErr, fmt.Errorf("failed to persist local node config from before the restore: %w", err))
		}
	}

	return retErr.ErrorOrNil()
}

// noopSnapshotter implements the fsm.Snapshot interface. It doesn't do anything
// since our SnapshotStore reads data out of the FSM on Open().
type noopSnapshotter struct {
	fsm *FSM
}

// Persist implements the fsm.Snapshot interface. It doesn't need to persist any
// state data, but it does persist the raft metadata. This is necessary so we
// can be sure to capture indexes for operation types that are not sent to the
// FSM.
func (s *noopSnapshotter) Persist(sink raft.SnapshotSink) error {
	pebbleSnapshotSink := sink.(*PebbleSnapshotSink)

	// We are processing a snapshot, fastforward the index, term, and
	// configuration to the latest seen by the raft system.
	if err := s.fsm.witnessSnapshot(&pebbleSnapshotSink.meta); err != nil {
		return err
	}

	return nil
}

// Release doesn't do anything.
func (s *noopSnapshotter) Release() {}

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
		return false, fmt.Errorf("error encoding chunk info: %w", err)
	}

	prefix, key := f.chunkPaths(chunk)

	entry := &physical.Entry{
		Key:   key,
		Value: b,
	}

	f.f.l.RLock()
	defer f.f.l.RUnlock()

	// Start a write transaction.
	done := new(bool)
	batch := f.f.db.NewBatch()
	defer func() { _ = batch.Close() }()

	byteKey := append(dataBucketPrefix, []byte(entry.Key)...)
	if err := batch.Set(byteKey, entry.Value, pebbleWriteOptions); err != nil {
		return *done, fmt.Errorf("error storing chunk info: %w", err)
	}

	iter := batch.NewIter(&pebble.IterOptions{
		LowerBound: dataBucketPrefix,
	})

	var keys []string
	prefixBytes := []byte(prefix)
	truePrefix := append(dataBucketPrefix, prefixBytes...)

	for iter.SeekGE(prefixBytes); iter.Valid(); iter.Next() {
		k := iter.Key()
		if bytes.HasPrefix(k, truePrefix) {
			stringKey := string(k)
			stringKey = strings.TrimPrefix(stringKey, string(truePrefix))
			if i := strings.Index(key, "/"); i == -1 {
				// Add objects only from the current 'folder'
				keys = append(keys, key)
			} else {
				// Add truncated 'folder' paths
				keys = strutil.AppendIfMissing(keys, string(key[:i+1]))
			}
		}
	}

	batch.Commit(pebbleWriteOptions)
	*done = uint32(len(keys)) == chunk.NumChunks
	return *done, nil
}

func (f *FSMChunkStorage) FinalizeOp(opNum uint64) ([]*raftchunking.ChunkInfo, error) {
	ret, err := f.chunksForOpNum(opNum)
	if err != nil {
		return nil, fmt.Errorf("error getting chunks for op keys: %w", err)
	}

	prefix, _ := f.chunkPaths(&raftchunking.ChunkInfo{OpNum: opNum})
	if err := f.f.DeletePrefix(f.ctx, prefix); err != nil {
		return nil, fmt.Errorf("error deleting prefix after op finalization: %w", err)
	}

	return ret, nil
}

func (f *FSMChunkStorage) chunksForOpNum(opNum uint64) ([]*raftchunking.ChunkInfo, error) {
	prefix, _ := f.chunkPaths(&raftchunking.ChunkInfo{OpNum: opNum})

	opChunkKeys, err := f.f.List(f.ctx, prefix)
	if err != nil {
		return nil, fmt.Errorf("error fetching op chunk keys: %w", err)
	}

	if len(opChunkKeys) == 0 {
		return nil, nil
	}

	var ret []*raftchunking.ChunkInfo

	for _, v := range opChunkKeys {
		seqNum, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting seqnum to integer: %w", err)
		}

		entry, err := f.f.Get(f.ctx, prefix+v)
		if err != nil {
			return nil, fmt.Errorf("error fetching chunkinfo: %w", err)
		}

		var ci raftchunking.ChunkInfo
		if err := jsonutil.DecodeJSON(entry.Value, &ci); err != nil {
			return nil, fmt.Errorf("error decoding chunkinfo json: %w", err)
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
		return nil, fmt.Errorf("error doing recursive list for chunk saving: %w", err)
	}

	if len(opNums) == 0 {
		return nil, nil
	}

	ret := make(raftchunking.ChunkMap, len(opNums))
	for _, opNumStr := range opNums {
		opNum, err := strconv.ParseInt(opNumStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing op num during chunk saving: %w", err)
		}

		opChunks, err := f.chunksForOpNum(uint64(opNum))
		if err != nil {
			return nil, fmt.Errorf("error getting chunks for op keys during chunk saving: %w", err)
		}

		ret[uint64(opNum)] = opChunks
	}

	return ret, nil
}

func (f *FSMChunkStorage) RestoreChunks(chunks raftchunking.ChunkMap) error {
	if err := f.f.DeletePrefix(f.ctx, chunkingPrefix); err != nil {
		return fmt.Errorf("error deleting prefix for chunk restoration: %w", err)
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
				return fmt.Errorf("error storing chunk during restoration: %w", err)
			}
		}
	}

	return nil
}
