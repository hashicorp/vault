// +build foundationdb

package foundationdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"bytes"
	"encoding/binary"

	log "github.com/hashicorp/go-hclog"
	uuid "github.com/hashicorp/go-uuid"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/physical"
)

const (
	// The minimum acceptable API version
	minAPIVersion = 520

	// The namespace under our top directory containing keys only for list operations
	metaKeysNamespace = "_meta-keys"

	// The namespace under our top directory containing the actual data
	dataNamespace = "_data"

	// The namespace under our top directory containing locks
	lockNamespace = "_lock"

	// Path hierarchy markers
	// - an entry in a directory (included in list)
	dirEntryMarker = "/\x01"
	// - a path component (excluded from list)
	dirPathMarker = "/\x02"
)

var (
	// 64bit 1 and -1 for FDB atomic Add()
	atomicArgOne      = []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	atomicArgMinusOne = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
)

// Verify FDBBackend satisfies the correct interfaces
var _ physical.Backend = (*FDBBackend)(nil)
var _ physical.Transactional = (*FDBBackend)(nil)
var _ physical.HABackend = (*FDBBackend)(nil)
var _ physical.Lock = (*FDBBackendLock)(nil)

// FDBBackend is a physical backend that stores data at a specific
// prefix within FoundationDB.
type FDBBackend struct {
	logger        log.Logger
	haEnabled     bool
	db            fdb.Database
	metaKeysSpace subspace.Subspace
	dataSpace     subspace.Subspace
	lockSpace     subspace.Subspace
	instanceUUID  string
}

func concat(a []byte, b ...byte) []byte {
	r := make([]byte, len(a)+len(b))

	copy(r, a)
	copy(r[len(a):], b)

	return r
}

func decoratePrefix(prefix string) ([]byte, error) {
	pathElements := strings.Split(prefix, "/")
	decoratedPrefix := strings.Join(pathElements[:len(pathElements)-1], dirPathMarker)

	return []byte(decoratedPrefix + dirEntryMarker), nil
}

// Turn a path string into a decorated byte array to be used as (part of) a key
// foo              /\x01foo
// foo/             /\x01foo/
// foo/bar          /\x02foo/\x01bar
// foo/bar/         /\x02foo/\x01bar/
// foo/bar/baz      /\x02foo/\x02bar/\x01baz
// foo/bar/baz/     /\x02foo/\x02bar/\x01baz/
// foo/bar/baz/quux /\x02foo/\x02bar/\x02baz/\x01quux
// This allows for range queries to retrieve the "directory" listing. The
// decoratePrefix() function builds the path leading up to the leaf.
func decoratePath(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("Invalid empty path")
	}

	path = "/" + path

	isDir := strings.HasSuffix(path, "/")
	path = strings.TrimRight(path, "/")

	lastSlash := strings.LastIndexByte(path, '/')
	decoratedPrefix, err := decoratePrefix(path[:lastSlash+1])
	if err != nil {
		return nil, err
	}

	leaf := path[lastSlash+1:]
	if isDir {
		leaf += "/"
	}

	return concat(decoratedPrefix, []byte(leaf)...), nil
}

// Turn a decorated byte array back into a path string
func undecoratePath(decoratedPath []byte) string {
	ret := strings.Replace(string(decoratedPath), dirPathMarker, "/", -1)
	ret = strings.Replace(ret, dirEntryMarker, "/", -1)

	return strings.TrimLeft(ret, "/")
}

// NewFDBBackend constructs a FoundationDB backend storing keys in the
// top-level directory designated by path
func NewFDBBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Get the top-level directory name
	path, ok := conf["path"]
	if !ok {
		path = "vault"
	}
	logger.Debug("config path set", "path", path)

	dirPath := strings.Split(strings.Trim(path, "/"), "/")

	// TLS support
	tlsCertFile, hasCertFile := conf["tls_cert_file"]
	tlsKeyFile, hasKeyFile := conf["tls_key_file"]
	tlsCAFile, hasCAFile := conf["tls_ca_file"]

	tlsEnabled := hasCertFile && hasKeyFile && hasCAFile

	if (hasCertFile || hasKeyFile || hasCAFile) && !tlsEnabled {
		return nil, fmt.Errorf("FoundationDB TLS requires all 3 of tls_cert_file, tls_key_file, and tls_ca_file")
	}

	tlsVerifyPeers, ok := conf["tls_verify_peers"]
	if !ok && tlsEnabled {
		return nil, fmt.Errorf("Required option tls_verify_peers not set in configuration")
	}

	// FoundationDB API version
	fdbApiVersionStr, ok := conf["api_version"]
	if !ok {
		return nil, fmt.Errorf("FoundationDB API version not specified")
	}

	fdbApiVersionInt, err := strconv.Atoi(fdbApiVersionStr)
	if err != nil {
		return nil, errwrap.Wrapf("failed to parse fdb_api_version parameter: {{err}}", err)
	}

	// Check requested FDB API version against minimum required API version
	if fdbApiVersionInt < minAPIVersion {
		return nil, fmt.Errorf("Configured FoundationDB API version lower than minimum required version: %d < %d", fdbApiVersionInt, minAPIVersion)
	}

	logger.Debug("FoundationDB API version set", "fdb_api_version", fdbApiVersionInt)

	// FoundationDB cluster file
	fdbClusterFile, ok := conf["cluster_file"]
	if !ok {
		return nil, fmt.Errorf("FoundationDB cluster file not specified")
	}

	haEnabled := false
	haEnabledStr, ok := conf["ha_enabled"]
	if ok {
		haEnabled, err = strconv.ParseBool(haEnabledStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed to parse ha_enabled parameter: {{err}}", err)
		}
	}

	instanceUUID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, errwrap.Wrapf("could not generate instance UUID: {{err}}", err)
	}
	logger.Debug("Instance UUID", "uuid", instanceUUID)

	if err := fdb.APIVersion(fdbApiVersionInt); err != nil {
		return nil, errwrap.Wrapf("failed to set FDB API version: {{err}}", err)
	}

	if tlsEnabled {
		opts := fdb.Options()

		tlsPassword, ok := conf["tls_password"]
		if ok {
			err := opts.SetTLSPassword(tlsPassword)
			if err != nil {
				return nil, errwrap.Wrapf("failed to set TLS password: {{err}}", err)
			}
		}

		err := opts.SetTLSCaPath(tlsCAFile)
		if err != nil {
			return nil, errwrap.Wrapf("failed to set TLS CA bundle path: {{err}}", err)
		}

		err = opts.SetTLSCertPath(tlsCertFile)
		if err != nil {
			return nil, errwrap.Wrapf("failed to set TLS certificate path: {{err}}", err)
		}

		err = opts.SetTLSKeyPath(tlsKeyFile)
		if err != nil {
			return nil, errwrap.Wrapf("failed to set TLS key path: {{err}}", err)
		}

		err = opts.SetTLSVerifyPeers([]byte(tlsVerifyPeers))
		if err != nil {
			return nil, errwrap.Wrapf("failed to set TLS peer verification criteria: {{err}}", err)
		}
	}

	db, err := fdb.Open(fdbClusterFile, []byte("DB"))
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to open database with cluster file '%s': {{err}}", fdbClusterFile), err)
	}

	topDir, err := directory.CreateOrOpen(db, dirPath, nil)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to create/open top-level directory '%s': {{err}}", path), err)
	}

	// Setup the backend
	f := &FDBBackend{
		logger:        logger,
		haEnabled:     haEnabled,
		db:            db,
		metaKeysSpace: topDir.Sub(metaKeysNamespace),
		dataSpace:     topDir.Sub(dataNamespace),
		lockSpace:     topDir.Sub(lockNamespace),
		instanceUUID:  instanceUUID,
	}
	return f, nil
}

// Increase refcount on directories in the path, from the bottom -> up
func (f *FDBBackend) incDirsRefcount(tr fdb.Transaction, path string) error {
	pathElements := strings.Split(strings.TrimRight(path, "/"), "/")

	for i := len(pathElements) - 1; i != 0; i-- {
		dPath, err := decoratePath(strings.Join(pathElements[:i], "/") + "/")
		if err != nil {
			return errwrap.Wrapf("error incrementing directories refcount: {{err}}", err)
		}

		// Atomic +1
		tr.Add(fdb.Key(concat(f.metaKeysSpace.Bytes(), dPath...)), atomicArgOne)
		tr.Add(fdb.Key(concat(f.dataSpace.Bytes(), dPath...)), atomicArgOne)
	}

	return nil
}

type DirsDecTodo struct {
	fkey   fdb.Key
	future fdb.FutureByteSlice
}

// Decrease refcount on directories in the path, from the bottom -> up, and remove empty ones
func (f *FDBBackend) decDirsRefcount(tr fdb.Transaction, path string) error {
	pathElements := strings.Split(strings.TrimRight(path, "/"), "/")

	dirsTodo := make([]DirsDecTodo, 0, len(pathElements)*2)

	for i := len(pathElements) - 1; i != 0; i-- {
		dPath, err := decoratePath(strings.Join(pathElements[:i], "/") + "/")
		if err != nil {
			return errwrap.Wrapf("error decrementing directories refcount: {{err}}", err)
		}

		metaFKey := fdb.Key(concat(f.metaKeysSpace.Bytes(), dPath...))
		dirsTodo = append(dirsTodo, DirsDecTodo{
			fkey:   metaFKey,
			future: tr.Get(metaFKey),
		})

		dataFKey := fdb.Key(concat(f.dataSpace.Bytes(), dPath...))
		dirsTodo = append(dirsTodo, DirsDecTodo{
			fkey:   dataFKey,
			future: tr.Get(dataFKey),
		})
	}

	for _, todo := range dirsTodo {
		value, err := todo.future.Get()
		if err != nil {
			return errwrap.Wrapf("error getting directory refcount while decrementing: {{err}}", err)
		}

		// The directory entry does not exist; this is not expected
		if value == nil {
			return fmt.Errorf("non-existent directory while decrementing directory refcount")
		}

		var count int64
		err = binary.Read(bytes.NewReader(value), binary.LittleEndian, &count)
		if err != nil {
			return errwrap.Wrapf("error reading directory refcount while decrementing: {{err}}", err)
		}

		if count > 1 {
			// Atomic -1
			tr.Add(todo.fkey, atomicArgMinusOne)
		} else {
			// Directory is empty, remove it
			tr.Clear(todo.fkey)
		}
	}

	return nil
}

func (f *FDBBackend) internalPut(tr fdb.Transaction, decoratedPath []byte, path string, value []byte) error {
	// Check that the meta key exists before blindly increasing the refcounts
	// in the directory hierarchy; this protects against commit_unknown_result
	// and other similar cases where a previous transaction may have gone
	// through without us knowing for sure.

	metaFKey := fdb.Key(concat(f.metaKeysSpace.Bytes(), decoratedPath...))
	metaFuture := tr.Get(metaFKey)

	dataFKey := fdb.Key(concat(f.dataSpace.Bytes(), decoratedPath...))
	tr.Set(dataFKey, value)

	value, err := metaFuture.Get()
	if err != nil {
		return errwrap.Wrapf("Put error while getting meta key: {{err}}", err)
	}

	if value == nil {
		tr.Set(metaFKey, []byte{})
		return f.incDirsRefcount(tr, path)
	}

	return nil
}

func (f *FDBBackend) internalClear(tr fdb.Transaction, decoratedPath []byte, path string) error {
	// Same as above - check existence of the meta key before taking any
	// action, to protect against a possible previous commit_unknown_result
	// error.

	metaFKey := fdb.Key(concat(f.metaKeysSpace.Bytes(), decoratedPath...))

	value, err := tr.Get(metaFKey).Get()
	if err != nil {
		return errwrap.Wrapf("Delete error while getting meta key: {{err}}", err)
	}

	if value != nil {
		dataFKey := fdb.Key(concat(f.dataSpace.Bytes(), decoratedPath...))
		tr.Clear(dataFKey)
		tr.Clear(metaFKey)
		return f.decDirsRefcount(tr, path)
	}

	return nil
}

type TxnTodo struct {
	decoratedPath []byte
	op            *physical.TxnEntry
}

// Used to run multiple entries via a transaction
func (f *FDBBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	if len(txns) == 0 {
		return nil
	}

	todo := make([]*TxnTodo, len(txns))

	for i, op := range txns {
		if op.Operation != physical.DeleteOperation && op.Operation != physical.PutOperation {
			return fmt.Errorf("%q is not a supported transaction operation", op.Operation)
		}

		decoratedPath, err := decoratePath(op.Entry.Key)
		if err != nil {
			return errwrap.Wrapf(fmt.Sprintf("could not build decorated path for transaction item %s: {{err}}", op.Entry.Key), err)
		}

		todo[i] = &TxnTodo{
			decoratedPath: decoratedPath,
			op:            op,
		}
	}

	_, err := f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		for _, txnTodo := range todo {
			var err error
			switch txnTodo.op.Operation {
			case physical.DeleteOperation:
				err = f.internalClear(tr, txnTodo.decoratedPath, txnTodo.op.Entry.Key)
			case physical.PutOperation:
				err = f.internalPut(tr, txnTodo.decoratedPath, txnTodo.op.Entry.Key, txnTodo.op.Entry.Value)
			}

			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("operation %s failed for transaction item %s: {{err}}", txnTodo.op.Operation, txnTodo.op.Entry.Key), err)
			}
		}

		return nil, nil
	})

	if err != nil {
		return errwrap.Wrapf("transaction failed: {{err}}", err)
	}

	return nil
}

// Put is used to insert or update an entry
func (f *FDBBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"foundationdb", "put"}, time.Now())

	decoratedPath, err := decoratePath(entry.Key)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("could not build decorated path to put item %s: {{err}}", entry.Key), err)
	}

	_, err = f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		err := f.internalPut(tr, decoratedPath, entry.Key, entry.Value)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("put failed for item %s: {{err}}", entry.Key), err)
	}

	return nil
}

// Get is used to fetch an entry
// Return nil for non-existent keys
func (f *FDBBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"foundationdb", "get"}, time.Now())

	decoratedPath, err := decoratePath(key)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("could not build decorated path to get item %s: {{err}}", key), err)
	}

	fkey := fdb.Key(concat(f.dataSpace.Bytes(), decoratedPath...))

	value, err := f.db.ReadTransact(func(rtr fdb.ReadTransaction) (interface{}, error) {
		value, err := rtr.Get(fkey).Get()
		if err != nil {
			return nil, err
		}

		return value, nil
	})

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("get failed for item %s: {{err}}", key), err)
	}
	if value.([]byte) == nil {
		return nil, nil
	}

	return &physical.Entry{
		Key:   key,
		Value: value.([]byte),
	}, nil
}

// Delete is used to permanently delete an entry
func (f *FDBBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"foundationdb", "delete"}, time.Now())

	decoratedPath, err := decoratePath(key)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("could not build decorated path to delete item %s: {{err}}", key), err)
	}

	_, err = f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		err := f.internalClear(tr, decoratedPath, key)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("delete failed for item %s: {{err}}", key), err)
	}

	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
// Return empty string slice for non-existent directories
func (f *FDBBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"foundationdb", "list"}, time.Now())

	prefix = strings.TrimRight("/"+prefix, "/") + "/"

	decoratedPrefix, err := decoratePrefix(prefix)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("could not build decorated path to list prefix %s: {{err}}", prefix), err)
	}

	// The beginning of the range is /\x02foo/\x02bar/\x01 (the decorated prefix) to list foo/bar/
	rangeBegin := fdb.Key(concat(f.metaKeysSpace.Bytes(), decoratedPrefix...))
	rangeEnd := fdb.Key(concat(rangeBegin, 0xff))
	pathRange := fdb.KeyRange{rangeBegin, rangeEnd}
	keyPrefixLen := len(rangeBegin)

	content, err := f.db.ReadTransact(func(rtr fdb.ReadTransaction) (interface{}, error) {
		dirList := make([]string, 0, 0)

		ri := rtr.GetRange(pathRange, fdb.RangeOptions{Mode: fdb.StreamingModeWantAll}).Iterator()

		for ri.Advance() {
			kv := ri.MustGet()

			// Strip length of the rangeBegin key off the FDB key, yielding
			// the part of the key we're interested in, which does not need
			// to be undecorated, by construction.
			dirList = append(dirList, string(kv.Key[keyPrefixLen:]))
		}

		return dirList, nil
	})

	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("could not list prefix %s: {{err}}", prefix), err)
	}

	return content.([]string), nil
}

type FDBBackendLock struct {
	f     *FDBBackend
	key   string
	value string
	fkey  fdb.Key
	lock  sync.Mutex
}

// LockWith is used for mutual exclusion based on the given key.
func (f *FDBBackend) LockWith(key, value string) (physical.Lock, error) {
	return &FDBBackendLock{
		f:     f,
		key:   key,
		value: value,
		fkey:  f.lockSpace.Pack(tuple.Tuple{key}),
	}, nil
}

func (f *FDBBackend) HAEnabled() bool {
	return f.haEnabled
}

const (
	// Position of elements in the lock content tuple
	lockContentValueIdx   = 0
	lockContentOwnerIdx   = 1
	lockContentExpiresIdx = 2

	// Number of elements in the lock content tuple
	lockTupleContentElts = 3

	lockTTL                  = 15 * time.Second
	lockRenewInterval        = 5 * time.Second
	lockAcquireRetryInterval = 5 * time.Second
)

type FDBBackendLockContent struct {
	value     string
	ownerUUID string
	expires   time.Time
}

func packLock(content *FDBBackendLockContent) []byte {
	t := tuple.Tuple{content.value, content.ownerUUID, content.expires.UnixNano()}

	return t.Pack()
}

func unpackLock(tupleContent []byte) (*FDBBackendLockContent, error) {
	t, err := tuple.Unpack(tupleContent)
	if err != nil {
		return nil, err
	}

	if len(t) != lockTupleContentElts {
		return nil, fmt.Errorf("unexpected lock content, len %d != %d", len(t), lockTupleContentElts)
	}

	return &FDBBackendLockContent{
		value:     t[lockContentValueIdx].(string),
		ownerUUID: t[lockContentOwnerIdx].(string),
		expires:   time.Unix(0, t[lockContentExpiresIdx].(int64)),
	}, nil
}

func (fl *FDBBackendLock) getLockContent(tr fdb.Transaction) (*FDBBackendLockContent, error) {
	tupleContent, err := tr.Get(fl.fkey).Get()
	if err != nil {
		return nil, err
	}

	// Lock doesn't exist
	if tupleContent == nil {
		return nil, fmt.Errorf("non-existent lock %s", fl.key)
	}

	content, err := unpackLock(tupleContent)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to unpack lock %s: {{err}}", fl.key), err)
	}

	return content, nil
}

func (fl *FDBBackendLock) setLockContent(tr fdb.Transaction, content *FDBBackendLockContent) {
	tr.Set(fl.fkey, packLock(content))
}

func (fl *FDBBackendLock) isOwned(content *FDBBackendLockContent) bool {
	return content.ownerUUID == fl.f.instanceUUID
}

func (fl *FDBBackendLock) isExpired(content *FDBBackendLockContent) bool {
	return time.Now().After(content.expires)
}

func (fl *FDBBackendLock) acquireTryLock(acquired chan struct{}, errors chan error) (bool, error) {
	wonTheRace, err := fl.f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		tupleContent, err := tr.Get(fl.fkey).Get()
		if err != nil {
			return nil, errwrap.Wrapf("could not read lock: {{err}}", err)
		}

		// Lock exists
		if tupleContent != nil {
			content, err := unpackLock(tupleContent)
			if err != nil {
				return nil, errwrap.Wrapf(fmt.Sprintf("failed to unpack lock %s: {{err}}", fl.key), err)
			}

			if fl.isOwned(content) {
				return nil, fmt.Errorf("lock %s already held", fl.key)
			}

			// The lock already exists, is not owned by us, and is not expired
			if !fl.isExpired(content) {
				return false, nil
			}
		}

		// Lock doesn't exist, or exists but is expired, we can go ahead
		content := &FDBBackendLockContent{
			value:     fl.value,
			ownerUUID: fl.f.instanceUUID,
			expires:   time.Now().Add(lockTTL),
		}

		fl.setLockContent(tr, content)

		return true, nil
	})

	if err != nil {
		errors <- err
		return false, err
	}

	if wonTheRace.(bool) {
		close(acquired)
	}

	return wonTheRace.(bool), nil
}

func (fl *FDBBackendLock) acquireLock(abandon chan struct{}, acquired chan struct{}, errors chan error) {
	ticker := time.NewTicker(lockAcquireRetryInterval)
	defer ticker.Stop()

	lockAcquired, err := fl.acquireTryLock(acquired, errors)
	if lockAcquired || err != nil {
		return
	}

	for {
		select {
		case <-abandon:
			return
		case <-ticker.C:
			lockAcquired, err := fl.acquireTryLock(acquired, errors)
			if lockAcquired || err != nil {
				return
			}
		}
	}
}

func (fl *FDBBackendLock) maintainLock(lost <-chan struct{}) {
	ticker := time.NewTicker(lockRenewInterval)
	for {
		select {
		case <-ticker.C:
			_, err := fl.f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
				content, err := fl.getLockContent(tr)
				if err != nil {
					return nil, err
				}

				// We don't own the lock
				if !fl.isOwned(content) {
					return nil, fmt.Errorf("lost lock %s", fl.key)
				}

				// The lock is expired
				if fl.isExpired(content) {
					return nil, fmt.Errorf("lock %s expired", fl.key)
				}

				content.expires = time.Now().Add(lockTTL)

				fl.setLockContent(tr, content)

				return nil, nil
			})

			if err != nil {
				fl.f.logger.Error("lock maintain", "error", err)
			}

			// Failure to renew the lock will cause another node to take over
			// and the watch to fire. DB errors will also be caught by the watch.
		case <-lost:
			ticker.Stop()
			return
		}
	}
}

func (fl *FDBBackendLock) watchLock(lost chan struct{}) {
	for {
		watch, err := fl.f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
			content, err := fl.getLockContent(tr)
			if err != nil {
				return nil, err
			}

			// We don't own the lock
			if !fl.isOwned(content) {
				return nil, fmt.Errorf("lost lock %s", fl.key)
			}

			// The lock is expired
			if fl.isExpired(content) {
				return nil, fmt.Errorf("lock %s expired", fl.key)
			}

			// Set FDB watch on the lock
			future := tr.Watch(fl.fkey)

			return future, nil
		})

		if err != nil {
			fl.f.logger.Error("lock watch", "error", err)
			break
		}

		// Wait for the watch to fire, and go again
		watch.(fdb.FutureNil).Get()
	}

	close(lost)
}

func (fl *FDBBackendLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	fl.lock.Lock()
	defer fl.lock.Unlock()

	var (
		// Inform the lock owner that we lost the lock
		lost = make(chan struct{})

		// Tell our watch and renewal routines the lock has been abandoned
		abandon = make(chan struct{})

		// Feedback from lock acquisition routine
		acquired = make(chan struct{})
		errors   = make(chan error)
	)

	// try to acquire the lock asynchronously
	go fl.acquireLock(abandon, acquired, errors)

	select {
	case <-acquired:
		// Maintain the lock after initial acquisition
		go fl.maintainLock(lost)
		// Watch the lock for changes
		go fl.watchLock(lost)
	case err := <-errors:
		// Initial acquisition failed
		close(abandon)
		return nil, err
	case <-stopCh:
		// Prospective lock owner cancelling lock acquisition
		close(abandon)
		return nil, nil
	}

	return lost, nil
}

func (fl *FDBBackendLock) Unlock() error {
	fl.lock.Lock()
	defer fl.lock.Unlock()

	_, err := fl.f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		content, err := fl.getLockContent(tr)
		if err != nil {
			return nil, errwrap.Wrapf("could not get lock content: {{err}}", err)
		}

		// We don't own the lock
		if !fl.isOwned(content) {
			return nil, nil
		}

		tr.Clear(fl.fkey)

		return nil, nil
	})

	if err != nil {
		return errwrap.Wrapf("unlock failed: {{err}}", err)
	}

	return nil
}

func (fl *FDBBackendLock) Value() (bool, string, error) {
	tupleContent, err := fl.f.db.ReadTransact(func(rtr fdb.ReadTransaction) (interface{}, error) {
		tupleContent, err := rtr.Get(fl.fkey).Get()
		if err != nil {
			return nil, errwrap.Wrapf("could not read lock: {{err}}", err)
		}

		return tupleContent, nil
	})

	if err != nil {
		return false, "", errwrap.Wrapf(fmt.Sprintf("get lock value failed for lock %s: {{err}}", fl.key), err)
	}
	if tupleContent.([]byte) == nil {
		return false, "", nil
	}

	content, err := unpackLock(tupleContent.([]byte))
	if err != nil {
		return false, "", errwrap.Wrapf(fmt.Sprintf("get lock value failed to unpack lock %s: {{err}}", fl.key), err)
	}

	return true, content.value, nil
}
