package physical

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/compose/canoe"
	log "github.com/mgutz/logxi/v1"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	canoeHATTL           = (3 * time.Second).Nanoseconds()
	canoeStateUpdateTime = 500 * time.Millisecond
)

// CanoeBackend provides a raft based Backend and HABackend for Vault
type CanoeBackend struct {
	raft   *canoe.Node
	logger log.Logger

	kvStore      map[string][]byte
	kvStoreMutex sync.Mutex

	kvStoreChanMutex sync.Mutex
	kvStoreObserveID uint64
	kvStoreChans     map[uint64]chan canoeKVStoreUpdate

	haStore          map[string]canoeHAValue
	haStoreObserveID uint64
	haStoreMutex     sync.Mutex

	haStoreChanMutex sync.Mutex
	haStoreChans     map[uint64]chan canoeHAStoreUpdate

	syncTicker *time.Ticker
}

type canoeHAValue struct {
	Value string `json:"value"`
	TTL   int64  `json:"ttl"`
	Time  int64  `json:"time"`
}

func newNonInitCanoeBackend(logger log.Logger) *CanoeBackend {
	return &CanoeBackend{
		logger: logger,

		haStore:      make(map[string]canoeHAValue),
		haStoreChans: make(map[uint64]chan canoeHAStoreUpdate),

		kvStore:      make(map[string][]byte),
		kvStoreChans: make(map[uint64]chan canoeKVStoreUpdate),

		syncTicker: time.NewTicker(canoeStateUpdateTime),
	}
}

func newCanoeBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	cBackend := newNonInitCanoeBackend(logger)

	canoeConfig, err := parseCanoeConfig(conf)
	if err != nil {
		return nil, err
	}

	canoeConfig.FSM = cBackend
	canoeConfig.Logger = &canoeLogger{logger: logger}

	canoeNode, err := canoe.NewNode(canoeConfig)
	if err != nil {
		return nil, err
	}

	cBackend.raft = canoeNode

	if err := cBackend.start(); err != nil {
		return nil, err
	}

	return cBackend, nil
}

func (cb *CanoeBackend) start() error {
	if err := cb.raft.Start(); err != nil {
		return err
	}

	go func(cb *CanoeBackend) {
		if err := cb.run(); err != nil {
			panic(err)
		}
	}(cb)
	return nil
}

func (cb *CanoeBackend) run() error {
	for {
		<-cb.syncTicker.C
		if err := cb.proposeHAReleaseStaleLocks(); err != nil {
			return err
		}
	}
}

func (cb *CanoeBackend) proposeHAReleaseStaleLocks() error {
	cmd := canoeHAReleaseStaleLocksCmd{
		Time: time.Now().UnixNano(),
	}

	if err := cb.proposeCmd(canoeHAReleaseStaleLocksOp, cmd); err != nil {
		return err
	}

	return nil
}

// Put inserts or updates a key
func (cb *CanoeBackend) Put(entry *Entry) error {
	observer := cb.kvStoreObserver()

	if err := cb.proposeKVSet(entry); err != nil {
		return err
	}

	retryTicker := time.NewTicker(300 * time.Millisecond)

	for {
		select {
		case <-retryTicker.C:
			if err := cb.proposeKVSet(entry); err != nil {
				observer.destroy()
				return err
			}
		case update := <-observer.updateCh:
			if update.Type == canoeKVStoreUpdateSetType && update.Key == entry.Key && bytes.Equal(update.Value, entry.Value) {
				observer.destroy()
				return nil
			}
		}
	}
}

// Delete deletes a key if it exists
func (cb *CanoeBackend) Delete(key string) error {
	// Do it this way since we will make Get consistent later
	val, err := cb.Get(key)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	}

	observer := cb.kvStoreObserver()

	if err := cb.proposeKVDelete(key); err != nil {
		return err
	}

	retryTicker := time.NewTicker(300 * time.Millisecond)

	for {
		select {
		case <-retryTicker.C:
			if err := cb.proposeKVDelete(key); err != nil {
				observer.destroy()
				return err
			}
		case update := <-observer.updateCh:
			if update.Type == canoeKVStoreUpdateDeleteType && update.Key == key {
				observer.destroy()
				return nil
			}
		}
	}
}

// Get retrieves a key if it exists
func (cb *CanoeBackend) Get(key string) (*Entry, error) {
	cb.kvStoreMutex.Lock()
	defer cb.kvStoreMutex.Unlock()

	if val, ok := cb.kvStore[key]; ok {
		return &Entry{
			Key:   key,
			Value: val,
		}, nil
	}
	return nil, nil
}

// List will return all paths/keys from a given prefix
func (cb *CanoeBackend) List(prefix string) ([]string, error) {
	var keys []string
	seen := make(map[string]interface{})

	cb.kvStoreMutex.Lock()
	defer cb.kvStoreMutex.Unlock()

	for key, _ := range cb.kvStore {
		if strings.HasPrefix(key, prefix) {
			trimmed := strings.TrimPrefix(key, prefix)
			sep := strings.Index(trimmed, "/")
			if sep == -1 {
				keys = append(keys, trimmed)
			} else {
				trimmed = trimmed[:sep+1]
				if _, ok := seen[trimmed]; !ok {
					keys = append(keys, trimmed)
					seen[trimmed] = struct{}{}
				}
			}
		}
	}
	return keys, nil
}

// LockWith provides a lock with the given key/value pair
func (cb *CanoeBackend) LockWith(key, value string) (Lock, error) {
	cl := &CanoeLock{
		backend:     cb,
		key:         key,
		value:       value,
		releaseChan: make(chan struct{}),
		syncTicker:  time.NewTicker(canoeStateUpdateTime),
	}
	// init the chan as a closed chan so we don't have to check
	// for nil on initialization
	close(cl.releaseChan)
	return cl, nil
}

// HAEnabled reveals if HA features should be enabled. This is always true with Canoe
func (cb *CanoeBackend) HAEnabled() bool {
	return true
}

var (
	// TODO: Define read ops
	canoeKVSetOp    = "CANOE_KV_SET_OP"
	canoeKVDeleteOp = "CANOE_KV_DELETE_OP"

	canoeHATryLockOp           = "CANOE_HA_TRY_LOCK_OP"
	canoeHAReleaseLockOp       = "CANOE_HA_RELEASE_LOCK_OP"
	canoeHAReleaseStaleLocksOp = "CANOE_HA_RELEASE_STALE_LOCKS_OP"
	canoeHARefreshLockOp       = "CANOE_HA_REFRESH_LOCK_OP"
)

type canoeCmd struct {
	Op   string `json:"op"`
	Data []byte `json:"data"`
}

type canoeKVSetCmd struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

type canoeKVDeleteCmd struct {
	Key string `json:"key"`
}

type canoeHATryLockCmd struct {
	Key        string       `json:"key"`
	CanoeValue canoeHAValue `json:"val"`
}

type canoeHAReleaseLockCmd struct {
	Key   string `json:"key"`
	Value string `json:"val"`
}

type canoeHAReleaseStaleLocksCmd struct {
	Time int64 `json:"time"`
}

type canoeHARefreshLockCmd struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Time  int64  `json:"time"`
}

var (
	canoeKVStoreUpdateSetType    = "CANOE_KV_UPDATE_SET"
	canoeKVStoreUpdateDeleteType = "CANOE_KV_UPDATE_DELETE"

	canoeHAStoreUpdateTryFailType     = "CANOE_HA_UPDATE_TRY_FAILED"
	canoeHAStoreUpdateReleaseFailType = "CANOE_HA_UPDATE_RELEASE_FAILED"
	canoeHAStoreUpdateSetType         = "CANOE_HA_UPDATE_SET"
	canoeHAStoreUpdateReleaseType     = "CANOE_HA_UPDATE_RELEASED"
)

type canoeKVStoreUpdate struct {
	Type  string `json:"update_type"`
	Key   string `json:"key"`
	Value []byte `json:"val"`
}

type canoeHAStoreUpdate struct {
	Type  string `json:"update_type"`
	Key   string `json:"key"`
	Value string `json:"val"`
}

// Apply fulfills the requirements for canoe.FSM interface
func (cb *CanoeBackend) Apply(data canoe.LogData) error {
	var cmd canoeCmd
	if err := json.Unmarshal(data, &cmd); err != nil {
		return err
	}

	switch cmd.Op {
	case canoeKVSetOp:
		if err := cb.applyKVSet(cmd.Data); err != nil {
			return err
		}
	case canoeKVDeleteOp:
		if err := cb.applyKVDelete(cmd.Data); err != nil {
			return err
		}
	case canoeHATryLockOp:
		if err := cb.applyHATryLock(cmd.Data); err != nil {
			return err
		}
	case canoeHAReleaseLockOp:
		if err := cb.applyHAReleaseLock(cmd.Data); err != nil {
			return err
		}
	case canoeHAReleaseStaleLocksOp:
		if err := cb.applyHAReleaseStaleLocks(cmd.Data); err != nil {
			return err
		}
	case canoeHARefreshLockOp:
		if err := cb.applyHARefreshLock(cmd.Data); err != nil {
			return err
		}
	default:
		return errors.New("Unknown operation")
	}
	return nil
}

func (cb *CanoeBackend) applyKVSet(cmdData []byte) error {
	var cmd canoeKVSetCmd
	if err := json.Unmarshal(cmdData, &cmd); err != nil {
		return err
	}

	cb.kvStoreMutex.Lock()
	defer cb.kvStoreMutex.Unlock()

	cb.kvStore[cmd.Key] = cmd.Value

	update := &canoeKVStoreUpdate{
		Type:  canoeKVStoreUpdateSetType,
		Key:   cmd.Key,
		Value: cmd.Value,
	}

	cb.observeKVStoreUpdate(update)

	return nil
}
func (cb *CanoeBackend) applyKVDelete(cmdData []byte) error {
	var cmd canoeKVSetCmd
	if err := json.Unmarshal(cmdData, &cmd); err != nil {
		return err
	}

	cb.kvStoreMutex.Lock()
	defer cb.kvStoreMutex.Unlock()

	if val, ok := cb.kvStore[cmd.Key]; ok {
		update := &canoeKVStoreUpdate{
			Type:  canoeKVStoreUpdateDeleteType,
			Key:   cmd.Key,
			Value: val,
		}
		cb.observeKVStoreUpdate(update)
		delete(cb.kvStore, cmd.Key)
	}

	return nil
}

func (cb *CanoeBackend) applyHATryLock(cmdData []byte) error {
	var cmd canoeHATryLockCmd
	if err := json.Unmarshal(cmdData, &cmd); err != nil {
		return err
	}

	cb.haStoreMutex.Lock()
	defer cb.haStoreMutex.Unlock()
	if val, ok := cb.haStore[cmd.Key]; ok && val.Value != cmd.CanoeValue.Value {
		haUpdate := &canoeHAStoreUpdate{
			Type:  canoeHAStoreUpdateTryFailType,
			Key:   cmd.Key,
			Value: cmd.CanoeValue.Value,
		}
		cb.observeHAStoreUpdate(haUpdate)
		return nil
	}

	cb.haStore[cmd.Key] = cmd.CanoeValue
	haUpdate := &canoeHAStoreUpdate{
		Type:  canoeHAStoreUpdateSetType,
		Key:   cmd.Key,
		Value: cmd.CanoeValue.Value,
	}
	cb.observeHAStoreUpdate(haUpdate)

	return nil
}

func (cb *CanoeBackend) applyHAReleaseLock(cmdData []byte) error {
	var cmd canoeHAReleaseLockCmd
	if err := json.Unmarshal(cmdData, &cmd); err != nil {
		return err
	}

	cb.haStoreMutex.Lock()
	defer cb.haStoreMutex.Unlock()
	if val, ok := cb.haStore[cmd.Key]; ok {
		if val.Value != cmd.Value {
			haUpdate := &canoeHAStoreUpdate{
				Type:  canoeHAStoreUpdateReleaseFailType,
				Key:   cmd.Key,
				Value: val.Value,
			}
			cb.observeHAStoreUpdate(haUpdate)
		} else {
			haUpdate := &canoeHAStoreUpdate{
				Type:  canoeHAStoreUpdateReleaseType,
				Key:   cmd.Key,
				Value: val.Value,
			}
			delete(cb.haStore, cmd.Key)
			cb.observeHAStoreUpdate(haUpdate)
		}

	} else {
		haUpdate := &canoeHAStoreUpdate{
			Type:  canoeHAStoreUpdateReleaseFailType,
			Key:   cmd.Key,
			Value: cmd.Value,
		}
		cb.observeHAStoreUpdate(haUpdate)
	}

	return nil
}

func (cb *CanoeBackend) applyHAReleaseStaleLocks(cmdData []byte) error {
	var cmd canoeHAReleaseStaleLocksCmd
	if err := json.Unmarshal(cmdData, &cmd); err != nil {
		return err
	}

	cb.haStoreMutex.Lock()
	defer cb.haStoreMutex.Unlock()

	for key, val := range cb.haStore {
		if cmd.Time >= val.Time+val.TTL {
			haUpdate := &canoeHAStoreUpdate{
				Type:  canoeHAStoreUpdateReleaseType,
				Key:   key,
				Value: val.Value,
			}
			delete(cb.haStore, key)
			cb.observeHAStoreUpdate(haUpdate)
		}
	}
	return nil
}

func (cb *CanoeBackend) applyHARefreshLock(cmdData []byte) error {
	var cmd canoeHARefreshLockCmd
	if err := json.Unmarshal(cmdData, &cmd); err != nil {
		return err
	}

	cb.haStoreMutex.Lock()
	defer cb.haStoreMutex.Unlock()

	val, ok := cb.haStore[cmd.Key]
	if ok && val.Value == cmd.Value {
		val.Time = cmd.Time
		cb.haStore[cmd.Key] = val

	}

	return nil
}

func (cb *CanoeBackend) proposeCmd(op string, data interface{}) error {
	reqData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	newCmd := &canoeCmd{
		Op:   op,
		Data: reqData,
	}

	newCmdData, err := json.Marshal(newCmd)
	if err != nil {
		return err
	}

	return cb.raft.Propose(newCmdData)
}

type canoeKVUpdateObserver struct {
	updateCh <-chan canoeKVStoreUpdate
	backend  *CanoeBackend
	id       uint64
}

func (cb *CanoeBackend) kvStoreObserver() *canoeKVUpdateObserver {
	cb.kvStoreChanMutex.Lock()
	defer cb.kvStoreChanMutex.Unlock()

	ch := make(chan canoeKVStoreUpdate)
	observer := &canoeKVUpdateObserver{
		updateCh: ch,
		backend:  cb,
		id:       atomic.AddUint64(&cb.kvStoreObserveID, 1),
	}
	cb.kvStoreChans[observer.id] = ch
	return observer
}

func (uo *canoeKVUpdateObserver) destroy() {
	uo.backend.unregisterKVObserver(uo)
}

func (cb *CanoeBackend) unregisterKVObserver(uo *canoeKVUpdateObserver) {
	cb.kvStoreChanMutex.Lock()
	defer cb.kvStoreChanMutex.Unlock()

	close(cb.kvStoreChans[uo.id])
	delete(cb.kvStoreChans, uo.id)
}

func (cb *CanoeBackend) observeKVStoreUpdate(ku *canoeKVStoreUpdate) {
	cb.kvStoreChanMutex.Lock()
	defer func() {
		cb.kvStoreChanMutex.Unlock()
	}()

	deadChans := make(map[uint64]chan canoeKVStoreUpdate)

	for key, val := range cb.kvStoreChans {
		select {
		case val <- *ku:
		case <-time.Tick(150 * time.Millisecond):
			cb.logger.Warn("Found racey lock. Trying to fix")
			deadChans[key] = val
		}
	}
	for len(deadChans) > 0 {
		cb.kvStoreChanMutex.Unlock()
		// For x-platform compatibility. On linux-x86 it isn't needed. But don't know how
		// golang implements locking constructs on ARM or Windows
		time.Sleep(5 * time.Millisecond)
		cb.kvStoreChanMutex.Lock()
		for key, val := range deadChans {
			// A chan was deleted since last attempt
			if _, ok := cb.kvStoreChans[key]; !ok {
				delete(deadChans, key)
			} else {
				select {
				case val <- *ku:
					delete(deadChans, key)
				case <-time.Tick(150 * time.Millisecond):
					cb.logger.Warn("Racey lock still exists. Probably not a delete race")
				}
			}
		}
	}
}

type canoeHAUpdateObserver struct {
	updateCh <-chan canoeHAStoreUpdate
	backend  *CanoeBackend
	id       uint64
}

func (cb *CanoeBackend) haStoreObserver() *canoeHAUpdateObserver {
	cb.haStoreChanMutex.Lock()
	defer cb.haStoreChanMutex.Unlock()
	ch := make(chan canoeHAStoreUpdate)
	observer := &canoeHAUpdateObserver{
		updateCh: ch,
		backend:  cb,
		id:       atomic.AddUint64(&cb.haStoreObserveID, 1),
	}
	cb.haStoreChans[observer.id] = ch
	return observer
}

func (uo *canoeHAUpdateObserver) destroy() {
	uo.backend.unregisterHAObserver(uo)
}

func (cb *CanoeBackend) unregisterHAObserver(uo *canoeHAUpdateObserver) {
	cb.haStoreChanMutex.Lock()
	defer cb.haStoreChanMutex.Unlock()

	close(cb.haStoreChans[uo.id])
	delete(cb.haStoreChans, uo.id)
}

func (cb *CanoeBackend) observeHAStoreUpdate(hu *canoeHAStoreUpdate) {
	cb.haStoreChanMutex.Lock()
	defer func() {
		cb.haStoreChanMutex.Unlock()
	}()

	deadChans := make(map[uint64]chan canoeHAStoreUpdate)

	for key, val := range cb.haStoreChans {
		select {
		case val <- *hu:
		case <-time.Tick(150 * time.Millisecond):
			cb.logger.Warn("Found racey lock. Trying to fix")
			deadChans[key] = val
		}
	}
	for len(deadChans) > 0 {
		cb.haStoreChanMutex.Unlock()
		time.Sleep(5 * time.Millisecond)
		cb.haStoreChanMutex.Lock()
		for key, val := range deadChans {
			// A chan was deleted since last attempt
			if _, ok := cb.haStoreChans[key]; !ok {
				delete(deadChans, key)
			} else {
				select {
				case val <- *hu:
					delete(deadChans, key)
				case <-time.Tick(150 * time.Millisecond):
					cb.logger.Warn("Racey lock still exists. Probably not a delete race")
				}
			}
		}
	}
}

func (cb *CanoeBackend) proposeKVSet(entry *Entry) error {
	cmd := canoeKVSetCmd{
		Key:   entry.Key,
		Value: entry.Value,
	}

	if err := cb.proposeCmd(canoeKVSetOp, cmd); err != nil {
		return err
	}

	return nil
}

func (cb *CanoeBackend) proposeKVDelete(key string) error {
	cmd := canoeKVDeleteCmd{
		Key: key,
	}

	if err := cb.proposeCmd(canoeKVDeleteOp, cmd); err != nil {
		return err
	}

	return nil
}

// CanoeLock provides distributed locking for Vault HA
type CanoeLock struct {
	backend     *CanoeBackend
	key         string
	value       string
	releaseChan chan struct{}
	stopc       chan struct{}
	vStopc      <-chan struct{}
	mutex       sync.Mutex

	syncTicker *time.Ticker
}

func (cl *CanoeLock) proposeHATryLock() error {
	cmd := canoeHATryLockCmd{
		Key: cl.key,
		CanoeValue: canoeHAValue{
			Value: cl.value,
			TTL:   canoeHATTL,
			Time:  time.Now().UnixNano(),
		},
	}

	if err := cl.backend.proposeCmd(canoeHATryLockOp, cmd); err != nil {
		return err
	}

	return nil
}

func (cl *CanoeLock) proposeHAReleaseLock() error {
	cmd := canoeHAReleaseLockCmd{
		Key:   cl.key,
		Value: cl.value,
	}

	if err := cl.backend.proposeCmd(canoeHAReleaseLockOp, cmd); err != nil {
		return err
	}

	return nil
}

func (cl *CanoeLock) proposeHARefreshLock() error {
	cmd := canoeHARefreshLockCmd{
		Key:   cl.key,
		Value: cl.value,
		Time:  time.Now().UnixNano(),
	}

	if err := cl.backend.proposeCmd(canoeHARefreshLockOp, cmd); err != nil {
		return err
	}

	return nil
}

func (cl *CanoeLock) refreshLock() error {
	for {
		select {
		case <-cl.stopc:
			return nil
		case <-cl.vStopc:
			return nil
		case <-cl.syncTicker.C:
			if err := cl.proposeHARefreshLock(); err != nil {
				return err
			}
		}
	}
}

// Lock blocks until it is able to acquire it's lock, or the stopCh closes
func (cl *CanoeLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	cl.mutex.Lock()
	cl.stopc = make(chan struct{})
	cl.vStopc = stopCh
	cl.mutex.Unlock()

	// first pass if we already hold lock somehow
	cl.backend.haStoreMutex.Lock()

	val, isSet := cl.backend.haStore[cl.key]
	if isSet && val.Value == cl.value {
		cl.mutex.Lock()
		defer cl.mutex.Unlock()

		// Don't garbage collect an existing chan if it's already created
		select {
		case <-cl.releaseChan:
			cl.releaseChan = make(chan struct{})
		default:
		}
		cl.backend.haStoreMutex.Unlock()
		return cl.releaseChan, nil
	}

	cl.backend.haStoreMutex.Unlock()

	observer := cl.backend.haStoreObserver()
	if !isSet {
		if err := cl.proposeHATryLock(); err != nil {
			return nil, err
		}
	}

	for {
		if !isSet {
			retryTicker := time.NewTicker(300 * time.Millisecond)
			select {
			case <-cl.vStopc:
				close(cl.stopc)
				observer.destroy()
				return nil, nil
			case <-retryTicker.C:
				if err := cl.proposeHATryLock(); err != nil {
					// TODO: Investigate why defer isn't called
					observer.destroy()
					return nil, err
				}
			case lockUpdate := <-observer.updateCh:
				if lockUpdate.Type == canoeHAStoreUpdateSetType {
					isSet = true
					if lockUpdate.Key == cl.key && lockUpdate.Value == cl.value {
						cl.mutex.Lock()
						defer cl.mutex.Unlock()
						select {
						case <-cl.releaseChan:
							cl.releaseChan = make(chan struct{})
						default:
						}
						go func(cl *CanoeLock) {
							if err := cl.refreshLock(); err != nil {
								cl.backend.logger.Fatal(err.Error())
							}
						}(cl)
						go func(cl *CanoeLock) {
							if err := cl.monitorDeletion(); err != nil {
								cl.backend.logger.Fatal(err.Error())
							}
						}(cl)
						observer.destroy()
						return cl.releaseChan, nil
					}
				}
			}
		} else {
			select {
			case <-cl.vStopc:
				close(cl.stopc)
				observer.destroy()
				return nil, nil
			case lockUpdate := <-observer.updateCh:
				if lockUpdate.Type == canoeHAStoreUpdateReleaseType {
					isSet = false
					if err := cl.proposeHATryLock(); err != nil {
						observer.destroy()
						return nil, err
					}
				}
			}
		}
	}
}

// Unlock blocks until the lock's lock has been released
func (cl *CanoeLock) Unlock() error {
	cl.backend.haStoreMutex.Lock()

	val, isSet := cl.backend.haStore[cl.key]
	if !isSet || val.Value != cl.value {
		return nil
	}

	// setup observer since we now know we don't have lock
	observer := cl.backend.haStoreObserver()

	cl.backend.haStoreMutex.Unlock()

	if err := cl.proposeHAReleaseLock(); err != nil {
		return err
	}

	retryTicker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-cl.releaseChan:
			observer.destroy()
			return nil
		case <-retryTicker.C:
			if err := cl.proposeHAReleaseLock(); err != nil {
				observer.destroy()
				return err
			}
		case update := <-observer.updateCh:
			if update.Type == canoeHAStoreUpdateReleaseType {
				if update.Key == cl.key && update.Value == cl.value {
					observer.destroy()
					return nil
				}
			}
		}
	}
}

func (cl *CanoeLock) monitorDeletion() error {
	observer := cl.backend.haStoreObserver()
	for {
		select {
		case <-cl.vStopc:
			close(cl.stopc)
			close(cl.releaseChan)
			observer.destroy()
			return nil
		case lockUpdate := <-observer.updateCh:
			if lockUpdate.Type == canoeHAStoreUpdateReleaseType {
				if lockUpdate.Key == cl.key && lockUpdate.Value == cl.value {
					observer.destroy()
					close(cl.stopc)
					close(cl.releaseChan)
					return nil
				}
			}
		}
	}
}

// Value returns if the lock is held, and if it is - the value of the lock
func (cl *CanoeLock) Value() (bool, string, error) {
	//TODO: Implement ReadRequestLog since trusting the state of a FSM
	// at any read time breaks raft linearability in the few milliseconds
	// after a raft leader failover

	cl.backend.haStoreMutex.Lock()
	defer cl.backend.haStoreMutex.Unlock()
	if val, ok := cl.backend.haStore[cl.key]; ok {
		return true, val.Value, nil
	}
	return false, "", nil
}

type canoeSnapshot struct {
	HAStore map[string]canoeHAValue `json:"ha_store"`
}

// Snapshot fulfills the canoe.FSM interface, providing a snapshot from where the FSM can be restored from
func (cb *CanoeBackend) Snapshot() (canoe.SnapshotData, error) {
	cb.haStoreMutex.Lock()
	defer cb.haStoreMutex.Unlock()

	snap := &canoeSnapshot{
		HAStore: cb.haStore,
	}
	return json.Marshal(snap)
}

// Restore fulfills the canoe.FSM interface
func (cb *CanoeBackend) Restore(data canoe.SnapshotData) error {
	cb.haStoreMutex.Lock()
	defer cb.haStoreMutex.Unlock()

	return json.Unmarshal(data, &cb.haStore)
}

func parseCanoeConfig(conf map[string]string) (*canoe.NodeConfig, error) {
	canoeConfig := &canoe.NodeConfig{}
	if val, ok := conf["config_port"]; ok {
		configPort, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("canoe setup error: Invalid api_port value - %s", val)

		}
		canoeConfig.ConfigurationPort = configPort
	} else {
		return nil, errors.New("canoe setup error: missing api_port argument")
	}

	if val, ok := conf["raft_port"]; ok {
		raftPort, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("canoe setup error: Invalid raft_port value - %s", val)
		}
		canoeConfig.RaftPort = raftPort
	} else {
		return nil, errors.New("canoe setup error: missing raft_port argument")
	}

	if val, ok := conf["id"]; ok {
		id, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("canoe setup error: invalid id value - %s", id)
		}
		canoeConfig.ID = id
	}

	if val, ok := conf["cluster_id"]; ok {
		clusterID, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("canoe setup error: invalid cluster_id value - %s", clusterID)
		}
		canoeConfig.ClusterID = clusterID
	}

	if val, ok := conf["peers"]; ok {
		peers := strings.Split(val, ",")

		for _, peer := range peers {
			if _, err := url.Parse(peer); err != nil {
				return nil, fmt.Errorf("canoe setup error: invalid peer - %s", peer)
			}
		}
		canoeConfig.BootstrapPeers = peers
	}

	if val, ok := conf["bootstrap_node"]; ok {
		bootstrap, err := strconv.ParseBool(val)
		if err != nil {
			return nil, fmt.Errorf("canoe setup error: invalid value for bootstrap_node - %s", val)
		}
		canoeConfig.BootstrapNode = bootstrap
	}

	if val, ok := conf["data_dir"]; ok {
		//TODO: Check for valid dir?
		canoeConfig.DataDir = val
	}

	snapConfig := &canoe.SnapshotConfig{
		Interval:             30 * time.Second,
		MinCommittedLogs:     20,
		MaxRetainedSnapshots: 3,
	}

	if val, ok := conf["snapshot_interval"]; ok {
		interval, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("canoe setup error: Invalid snapshot_interval value - %s", val)
		}
		snapConfig.Interval = time.Duration(interval) * time.Second
	}

	if val, ok := conf["snapshot_min_committed_logs"]; ok {
		minLogs, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("canoe setup error: Invalid snapshot_min_committed_logs value - %s", val)
		}
		snapConfig.MinCommittedLogs = minLogs
	}

	if val, ok := conf["max_retained_snapshots"]; ok {
		maxRetained, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("canoe setup error: Invalid snapshot_min_committed_logs value - %s", val)
		}
		snapConfig.MaxRetainedSnapshots = uint(maxRetained)
	}
	return canoeConfig, nil
}

type canoeLogger struct {
	logger log.Logger
}

// Debug fulfills canoe.Logger
func (l *canoeLogger) Debug(v ...interface{}) {
	l.logger.Debug(fmt.Sprintf("physical/canoe: %s", fmt.Sprint(v)))
}

// Debugf fulfills canoe.Logger
func (l *canoeLogger) Debugf(format string, v ...interface{}) {
	l.logger.Debug(fmt.Sprintf("physical/canoe: %s", fmt.Sprintf(format, v...)))
}

// Error fulfills canoe.Logger
func (l *canoeLogger) Error(v ...interface{}) {
	l.logger.Error(fmt.Sprintf("physical/canoe: %s", fmt.Sprint(v)))
}

// Errorf fulfills canoe.Logger
func (l *canoeLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf("physical/canoe: %s", fmt.Sprintf(format, v...)))
}

// Info fulfills canoe.Logger
func (l *canoeLogger) Info(v ...interface{}) {
	l.logger.Info(fmt.Sprintf("physical/canoe: %s", fmt.Sprint(v)))
}

// Infof fulfills canoe.Logger
func (l *canoeLogger) Infof(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf("physical/canoe: %s", fmt.Sprintf(format, v...)))
}

// Warning fulfills canoe.Logger
func (l *canoeLogger) Warning(v ...interface{}) {
	l.logger.Warn(fmt.Sprintf("physical/canoe: %s", fmt.Sprint(v)))
}

// Warningf fulfills canoe.Logger
func (l *canoeLogger) Warningf(format string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf("physical/canoe: %s", fmt.Sprintf(format, v...)))
}

// Fatal fulfills canoe.Logger
func (l *canoeLogger) Fatal(v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf("physical/canoe: %s", fmt.Sprint(v)))
}

// Fatalf fulfills canoe.Logger
func (l *canoeLogger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal(fmt.Sprintf("physical/canoe: %s", fmt.Sprintf(format, v...)))
}

// Panic fulfills canoe.Logger
func (l *canoeLogger) Panic(v ...interface{}) {
	l.logger.Error(fmt.Sprintf("physical/canoe: %s", fmt.Sprint(v)))
	panic(v)
}

// Panicf fulfills canoe.Logger
func (l *canoeLogger) Panicf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf("physical/canoe: %s", fmt.Sprintf(format, v...)))
	panic(fmt.Sprintf(format, v))
}
