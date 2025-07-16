// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/permitpool"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/diagnose"
	"golang.org/x/net/http2"
)

const (
	// consistencyModeDefault is the configuration value used to tell
	// consul to use default consistency.
	consistencyModeDefault = "default"

	// consistencyModeStrong is the configuration value used to tell
	// consul to use strong consistency.
	consistencyModeStrong = "strong"

	// nonExistentKey is used as part of a capabilities check against Consul
	nonExistentKey = "F35C28E1-7035-40BB-B865-6BED9E3A1B28"
)

// Verify ConsulBackend satisfies the correct interfaces
var (
	_ physical.Backend             = (*ConsulBackend)(nil)
	_ physical.FencingHABackend    = (*ConsulBackend)(nil)
	_ physical.Lock                = (*ConsulLock)(nil)
	_ physical.Transactional       = (*ConsulBackend)(nil)
	_ physical.TransactionalLimits = (*ConsulBackend)(nil)

	GetInTxnDisabledError = errors.New("get operations inside transactions are disabled in consul backend")
)

// ConsulBackend is a physical backend that stores data at specific
// prefix within Consul. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
// failGetInTxn is only used in tests.
type ConsulBackend struct {
	logger          log.Logger
	client          *api.Client
	path            string
	kv              *api.KV
	txn             *api.Txn
	permitPool      *permitpool.Pool
	consistencyMode string
	sessionTTL      string
	lockWaitTime    time.Duration
	failGetInTxn    *uint32
	activeNodeLock  atomic.Pointer[ConsulLock]
}

// NewConsulBackend constructs a Consul backend using the given API client
// and the prefix in the KV store.
func NewConsulBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Get the path in Consul
	path, ok := conf["path"]
	if !ok {
		path = "vault/"
	}
	if logger.IsDebug() {
		logger.Debug("config path set", "path", path)
	}

	// Ensure path is suffixed but not prefixed
	if !strings.HasSuffix(path, "/") {
		logger.Warn("appending trailing forward slash to path")
		path += "/"
	}
	if strings.HasPrefix(path, "/") {
		logger.Warn("trimming path of its forward slash")
		path = strings.TrimPrefix(path, "/")
	}

	sessionTTL := api.DefaultLockSessionTTL
	sessionTTLStr, ok := conf["session_ttl"]
	if ok {
		_, err := parseutil.ParseDurationSecond(sessionTTLStr)
		if err != nil {
			return nil, fmt.Errorf("invalid session_ttl: %w", err)
		}
		sessionTTL = sessionTTLStr
		if logger.IsDebug() {
			logger.Debug("config session_ttl set", "session_ttl", sessionTTL)
		}
	}

	lockWaitTime := api.DefaultLockWaitTime
	lockWaitTimeRaw, ok := conf["lock_wait_time"]
	if ok {
		d, err := parseutil.ParseDurationSecond(lockWaitTimeRaw)
		if err != nil {
			return nil, fmt.Errorf("invalid lock_wait_time: %w", err)
		}
		lockWaitTime = d
		if logger.IsDebug() {
			logger.Debug("config lock_wait_time set", "lock_wait_time", d)
		}
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err := strconv.Atoi(maxParStr)
		if err != nil {
			return nil, fmt.Errorf("failed parsing max_parallel parameter: %w", err)
		}
		if logger.IsDebug() {
			logger.Debug("max_parallel set", "max_parallel", maxParInt)
		}
	}

	consistencyMode, ok := conf["consistency_mode"]
	if ok {
		switch consistencyMode {
		case consistencyModeDefault, consistencyModeStrong:
		default:
			return nil, fmt.Errorf("invalid consistency_mode value: %q", consistencyMode)
		}
	} else {
		consistencyMode = consistencyModeDefault
	}

	// Configure the client
	consulConf := api.DefaultConfig()
	// Set MaxIdleConnsPerHost to the number of processes used in expiration.Restore
	consulConf.Transport.MaxIdleConnsPerHost = consts.ExpirationRestoreWorkerCount

	if err := SetupSecureTLS(context.Background(), consulConf, conf, logger, false); err != nil {
		return nil, fmt.Errorf("client setup failed: %w", err)
	}

	consulConf.HttpClient = &http.Client{Transport: consulConf.Transport}
	client, err := api.NewClient(consulConf)
	if err != nil {
		return nil, fmt.Errorf("client setup failed: %w", err)
	}

	// Set up the backend
	c := &ConsulBackend{
		logger:          logger,
		path:            path,
		client:          client,
		kv:              client.KV(),
		txn:             client.Txn(),
		permitPool:      permitpool.New(maxParInt),
		consistencyMode: consistencyMode,
		sessionTTL:      sessionTTL,
		lockWaitTime:    lockWaitTime,
		failGetInTxn:    new(uint32),
	}

	return c, nil
}

func SetupSecureTLS(ctx context.Context, consulConf *api.Config, conf map[string]string, logger log.Logger, isDiagnose bool) error {
	if addr, ok := conf["address"]; ok {
		consulConf.Address = addr
		if logger.IsDebug() {
			logger.Debug("config address set", "address", addr)
		}

		// Copied from the Consul API module; set the Scheme based on
		// the protocol field if address looks ike a URL.
		// This can enable the TLS configuration below.
		parts := strings.SplitN(addr, "://", 2)
		if len(parts) == 2 {
			if parts[0] == "http" || parts[0] == "https" {
				consulConf.Scheme = parts[0]
				consulConf.Address = parts[1]
				if logger.IsDebug() {
					logger.Debug("config address parsed", "scheme", parts[0])
					logger.Debug("config scheme parsed", "address", parts[1])
				}
			} // allow "unix:" or whatever else consul supports in the future
		}
	}
	if scheme, ok := conf["scheme"]; ok {
		consulConf.Scheme = scheme
		if logger.IsDebug() {
			logger.Debug("config scheme set", "scheme", scheme)
		}
	}
	if token, ok := conf["token"]; ok {
		consulConf.Token = token
		logger.Debug("config token set")
	}

	if consulConf.Scheme == "https" {
		if isDiagnose {
			certPath, okCert := conf["tls_cert_file"]
			keyPath, okKey := conf["tls_key_file"]
			if okCert && okKey {
				warnings, err := diagnose.TLSFileChecks(certPath, keyPath)
				for _, warning := range warnings {
					diagnose.Warn(ctx, warning)
				}
				if err != nil {
					return err
				}
				return nil
			}
			return fmt.Errorf("key or cert path: %s, %s, cannot be loaded from consul config file", certPath, keyPath)
		}

		// Use the parsed Address instead of the raw conf['address']
		tlsClientConfig, err := tlsutil.SetupTLSConfig(conf, consulConf.Address)
		if err != nil {
			return err
		}

		consulConf.Transport.TLSClientConfig = tlsClientConfig
		if err := http2.ConfigureTransport(consulConf.Transport); err != nil {
			return err
		}
		logger.Debug("configured TLS")
	} else {
		if isDiagnose {
			diagnose.Skipped(ctx, "HTTPS is not used, Skipping TLS verification.")
		}
	}
	return nil
}

// ExpandedCapabilitiesAvailable tests to see if Consul has KVGetOrEmpty and 128 entries per transaction available
func (c *ConsulBackend) ExpandedCapabilitiesAvailable(ctx context.Context) bool {
	available := false

	maxEntries := 128
	ops := make([]*api.TxnOp, maxEntries)
	for i := 0; i < maxEntries; i++ {
		ops[i] = &api.TxnOp{KV: &api.KVTxnOp{
			Key:  c.path + nonExistentKey,
			Verb: api.KVGetOrEmpty,
		}}
	}

	if err := c.permitPool.Acquire(ctx); err != nil {
		return false
	}
	defer c.permitPool.Release()

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)

	ok, resp, _, err := c.txn.Txn(ops, queryOpts)
	if ok && len(resp.Errors) == 0 && err == nil {
		available = true
	}

	return available
}

func (c *ConsulBackend) writeTxnOps(ctx context.Context, len int) ([]*api.TxnOp, string) {
	if len < 1 {
		len = 1
	}
	ops := make([]*api.TxnOp, 0, len+1)

	// If we don't have a lock yet, return a transaction with no session check. We
	// need to do this to allow writes during cluster initialization before there
	// is an active node.
	lock := c.activeNodeLock.Load()
	if lock == nil {
		return ops, ""
	}

	lockKey, lockSession := lock.Info()
	if lockKey == "" || lockSession == "" {
		return ops, ""
	}

	// If the context used to write has been marked as a special case write that
	// happens outside of a lock then don't add the session check.
	if physical.IsUnfencedWrite(ctx) {
		return ops, ""
	}

	// Insert the session check operation at index 0. This will allow us later to
	// work out easily if a write failure is because of the session check.
	ops = append(ops, &api.TxnOp{
		KV: &api.KVTxnOp{
			Verb:    api.KVCheckSession,
			Key:     lockKey,
			Session: lockSession,
		},
	})
	return ops, lockSession
}

// Transaction is used to run multiple entries via a transaction.
func (c *ConsulBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	return c.txnInternal(ctx, txns, "transaction")
}

func (c *ConsulBackend) txnInternal(ctx context.Context, txns []*physical.TxnEntry, apiOpName string) error {
	if len(txns) == 0 {
		return nil
	}
	defer metrics.MeasureSince([]string{"consul", apiOpName}, time.Now())

	failGetInTxn := atomic.LoadUint32(c.failGetInTxn)
	for _, t := range txns {
		if t.Operation == physical.GetOperation && failGetInTxn != 0 {
			return GetInTxnDisabledError
		}
	}

	ops, sessionID := c.writeTxnOps(ctx, len(txns))
	for _, t := range txns {
		o, err := c.makeApiTxn(t)
		if err != nil {
			return fmt.Errorf("error converting physical transactions into api transactions: %w", err)
		}

		ops = append(ops, o)
	}

	if err := c.permitPool.Acquire(ctx); err != nil {
		return err
	}
	defer c.permitPool.Release()

	var retErr *multierror.Error
	kvMap := make(map[string][]byte, 0)

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)

	ok, resp, _, err := c.txn.Txn(ops, queryOpts)
	if err != nil {
		if strings.Contains(err.Error(), "is too large") {
			return fmt.Errorf("%s: %w", physical.ErrValueTooLarge, err)
		}
		return err
	}
	if ok && len(resp.Errors) == 0 {
		// Loop over results and cache them in a map. Note that we're only caching
		// the first time we see a key, which _should_ correspond to a Get
		// operation, since we expect those come first in our txns slice (though
		// after check-session).
		for _, txnr := range resp.Results {
			if len(txnr.KV.Value) > 0 {
				// We need to trim the Consul kv path (typically "vault/") from the key
				// otherwise it won't match the transaction entries we have.
				key := strings.TrimPrefix(txnr.KV.Key, c.path)
				if _, found := kvMap[key]; !found {
					kvMap[key] = txnr.KV.Value
				}
			}
		}
	}

	if len(resp.Errors) > 0 {
		for _, res := range resp.Errors {
			retErr = multierror.Append(retErr, errors.New(res.What))
			if res.OpIndex == 0 && sessionID != "" {
				// We added a session check (sessionID not empty) so an error at OpIndex
				// 0 means that we failed that session check. We don't attempt to string
				// match because Consul can return at least three different errors here
				// with no common string. In all cases though failing this check means
				// we no longer hold the lock because it was released, modified or
				// deleted. Rather than just continuing to try writing until the
				// blocking query manages to notice we're no longer the lock holder
				// (which can take 10s of seconds even in good network conditions in my
				// testing) we can now Unlock directly here. Our ConsulLock now has a
				// shortcut that will cause the lock to close the leaderCh immediately
				// when we call without waiting for the blocking query to return (unlike
				// Consul's current Lock implementation). But before we unlock, we
				// should re-load the lock and ensure it's still the same instance we
				// just tried to write with in case this goroutine is somehow really
				// delayed and we actually acquired a whole new lock in the meantime!
				lock := c.activeNodeLock.Load()
				if lock != nil {
					_, lockSessionID := lock.Info()
					if sessionID == lockSessionID {
						c.logger.Warn("session check failed on write, we lost active node lock, stepping down", "err", res.What)
						lock.Unlock()
					}
				}
			}
		}
	}

	if retErr != nil {
		return retErr
	}

	// Loop over our get transactions and populate any values found in our map cache.
	for _, t := range txns {
		if val, ok := kvMap[t.Entry.Key]; ok && t.Operation == physical.GetOperation {
			newVal := make([]byte, len(val))
			copy(newVal, val)
			t.Entry.Value = newVal
		}
	}

	return nil
}

func (c *ConsulBackend) makeApiTxn(txn *physical.TxnEntry) (*api.TxnOp, error) {
	op := &api.KVTxnOp{
		Key: c.path + txn.Entry.Key,
	}
	switch txn.Operation {
	case physical.GetOperation:
		op.Verb = api.KVGetOrEmpty
	case physical.DeleteOperation:
		op.Verb = api.KVDelete
	case physical.PutOperation:
		op.Verb = api.KVSet
		op.Value = txn.Entry.Value
	default:
		return nil, fmt.Errorf("%q is not a supported transaction operation", txn.Operation)
	}

	return &api.TxnOp{KV: op}, nil
}

func (c *ConsulBackend) TransactionLimits() (int, int) {
	// Note that even for modern Consul versions that support 128 entries per txn,
	// we have an effective limit of 64 write operations because the other 64 are
	// used for undo log read operations. We also reserve 1 for a check-session
	// operation to prevent split brain so the most we allow WAL to put in a batch
	// is 63.
	return 63, 128 * 1024
}

// Put is used to insert or update an entry
func (c *ConsulBackend) Put(ctx context.Context, entry *physical.Entry) error {
	txns := []*physical.TxnEntry{
		{
			Operation: physical.PutOperation,
			Entry:     entry,
		},
	}
	return c.txnInternal(ctx, txns, "put")
}

// Get is used to fetch an entry
func (c *ConsulBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"consul", "get"}, time.Now())

	if err := c.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer c.permitPool.Release()

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)

	if c.consistencyMode == consistencyModeStrong {
		queryOpts.RequireConsistent = true
	}

	pair, _, err := c.kv.Get(c.path+key, queryOpts)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, nil
	}
	ent := &physical.Entry{
		Key:   key,
		Value: pair.Value,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (c *ConsulBackend) Delete(ctx context.Context, key string) error {
	txns := []*physical.TxnEntry{
		{
			Operation: physical.DeleteOperation,
			Entry: &physical.Entry{
				Key: key,
			},
		},
	}
	return c.txnInternal(ctx, txns, "delete")
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (c *ConsulBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"consul", "list"}, time.Now())
	scan := c.path + prefix

	// The TrimPrefix call below will not work correctly if we have "//" at the
	// end. This can happen in cases where you are e.g. listing the root of a
	// prefix in a logical backend via "/" instead of ""
	if strings.HasSuffix(scan, "//") {
		scan = scan[:len(scan)-1]
	}

	if err := c.permitPool.Acquire(ctx); err != nil {
		return nil, err
	}
	defer c.permitPool.Release()

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)

	out, _, err := c.kv.Keys(scan, "/", queryOpts)
	for idx, val := range out {
		out[idx] = strings.TrimPrefix(val, scan)
	}

	return out, err
}

func (c *ConsulBackend) FailGetInTxn(fail bool) {
	var val uint32
	if fail {
		val = 1
	}
	atomic.StoreUint32(c.failGetInTxn, val)
}

// LockWith is used for mutual exclusion based on the given key.
func (c *ConsulBackend) LockWith(key, value string) (physical.Lock, error) {
	cl := &ConsulLock{
		logger:          c.logger,
		client:          c.client,
		key:             c.path + key,
		value:           value,
		consistencyMode: c.consistencyMode,
		sessionTTL:      c.sessionTTL,
		lockWaitTime:    c.lockWaitTime,
	}
	return cl, nil
}

// HAEnabled indicates whether the HA functionality should be exposed.
// Currently always returns true.
func (c *ConsulBackend) HAEnabled() bool {
	return true
}

// DetectHostAddr is used to detect the host address by asking the Consul agent
func (c *ConsulBackend) DetectHostAddr() (string, error) {
	agent := c.client.Agent()
	self, err := agent.Self()
	if err != nil {
		return "", err
	}
	addr, ok := self["Member"]["Addr"].(string)
	if !ok {
		return "", fmt.Errorf("unable to convert an address to string")
	}
	return addr, nil
}

// RegisterActiveNodeLock is called after active node lock is obtained to allow
// us to fence future writes.
func (c *ConsulBackend) RegisterActiveNodeLock(l physical.Lock) error {
	cl, ok := l.(*ConsulLock)
	if !ok {
		return fmt.Errorf("invalid Lock type")
	}
	c.activeNodeLock.Store(cl)
	key, sessionID := cl.Info()
	c.logger.Info("registered active node lock", "key", key, "sessionID", sessionID)
	return nil
}

// ConsulLock is used to provide the Lock interface backed by Consul. We work
// around some limitations of Consuls api.Lock noted in
// https://github.com/hashicorp/consul/issues/18271 by creating and managing the
// session ourselves, while using Consul's Lock to do the heavy lifting.
type ConsulLock struct {
	logger          log.Logger
	client          *api.Client
	key             string
	value           string
	consistencyMode string
	sessionTTL      string
	lockWaitTime    time.Duration

	mu      sync.Mutex // protects session state
	session *lockSession
	// sessionID is a copy of the value from session.id. We use a separate field
	// because `Info` needs to keep returning the same sessionID after Unlock has
	// cleaned up the session state so that we continue to fence any writes still
	// in flight after the lock is Unlocked. It's easier to reason about that as a
	// separate field rather than keeping an already-terminated session object
	// around. Once Lock is called again this will be replaced (while mu is
	// locked) with the new session ID. Must hold mu to read or write this.
	sessionID string
}

type lockSession struct {
	// id is immutable after the session is created so does not need mu held
	id string

	// mu protects the lock and unlockCh to ensure they are only cleaned up once
	mu       sync.Mutex
	lock     *api.Lock
	unlockCh chan struct{}
}

func (s *lockSession) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	lockHeld := false
	defer func() {
		if !lockHeld {
			s.cleanupLocked()
		}
	}()

	consulLeaderCh, err := s.lock.Lock(stopCh)
	if err != nil {
		return nil, err
	}
	if consulLeaderCh == nil {
		// If both leaderCh and err are nil from Consul's Lock then it means we
		// waited for the lockWait without grabbing it.
		return nil, nil
	}
	// We got the Lock, monitor it!
	lockHeld = true
	leaderCh := make(chan struct{})
	go s.monitorLock(leaderCh, s.unlockCh, consulLeaderCh)
	return leaderCh, nil
}

// monitorLock waits for either unlockCh or consulLeaderCh to close and then
// closes leaderCh. It's designed to be run in a separate goroutine. Note that
// we pass unlockCh rather than accessing it via the member variable because it
// is mutated under the lock during Unlock so reading it from c could be racy.
// We just need the chan created at the call site here so we pass it instead of
// locking and unlocking in here.
func (s *lockSession) monitorLock(leaderCh chan struct{}, unlockCh, consulLeaderCh <-chan struct{}) {
	select {
	case <-unlockCh:
	case <-consulLeaderCh:
	}
	// We lost the lock. Close the leaderCh
	close(leaderCh)

	// Whichever chan closed, cleanup to unwind all the state. If we were
	// triggered by a cleanup call this will be a no-op, but if not it ensures all
	// state is cleaned up correctly.
	s.cleanup()
}

func (s *lockSession) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupLocked()
}

func (s *lockSession) cleanupLocked() {
	if s.lock != nil {
		s.lock.Unlock()
		s.lock = nil
	}
	if s.unlockCh != nil {
		close(s.unlockCh)
		s.unlockCh = nil
	}
	// Don't bother destroying sessions as they will be destroyed after TTL
	// anyway.
}

func (c *ConsulLock) createSession() (*lockSession, error) {
	se := &api.SessionEntry{
		Name: "Vault Lock",
		TTL:  c.sessionTTL,
		// We use Consul's default LockDelay of 15s by not specifying it
	}
	session, _, err := c.client.Session().Create(se, nil)
	if err != nil {
		return nil, err
	}

	opts := &api.LockOptions{
		Key:            c.key,
		Value:          []byte(c.value),
		Session:        session,
		MonitorRetries: 5,
		LockWaitTime:   c.lockWaitTime,
		SessionTTL:     c.sessionTTL,
	}
	lock, err := c.client.LockOpts(opts)
	if err != nil {
		// Don't bother destroying sessions as they will be destroyed after TTL
		// anyway.
		return nil, fmt.Errorf("failed to create lock: %w", err)
	}

	unlockCh := make(chan struct{})

	s := &lockSession{
		id:       session,
		lock:     lock,
		unlockCh: unlockCh,
	}

	// Start renewals of the session
	go func() {
		// Note we capture unlockCh here rather than s.unlockCh because s.unlockCh
		// is mutated on cleanup which is racy since we don't hold a lock here.
		// unlockCh will never be mutated though.
		err := c.client.Session().RenewPeriodic(c.sessionTTL, session, nil, unlockCh)
		if err != nil {
			c.logger.Error("failed to renew consul session for more than the TTL, lock lost", "err", err)
		}
		// release other resources for this session only i.e. don't c.Unlock as that
		// might now be locked under a different session).
		s.cleanup()
	}()
	return s, nil
}

func (c *ConsulLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.session != nil {
		return nil, fmt.Errorf("lock instance already locked")
	}

	session, err := c.createSession()
	if err != nil {
		return nil, err
	}
	leaderCh, err := session.Lock(stopCh)
	if leaderCh != nil && err == nil {
		// We hold the lock, store the session
		c.session = session
		c.sessionID = session.id
	}
	return leaderCh, err
}

func (c *ConsulLock) Unlock() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.session != nil {
		c.session.cleanup()
		c.session = nil
		// Don't clear c.sessionID since we rely on returning the same old ID after
		// Unlock until the next Lock.
	}
	return nil
}

func (c *ConsulLock) Value() (bool, string, error) {
	kv := c.client.KV()

	var queryOptions *api.QueryOptions
	if c.consistencyMode == consistencyModeStrong {
		queryOptions = &api.QueryOptions{
			RequireConsistent: true,
		}
	}

	pair, _, err := kv.Get(c.key, queryOptions)
	if err != nil {
		return false, "", err
	}
	if pair == nil {
		return false, "", nil
	}
	// Note that held is expected to mean "does _any_ node hold the lock" not
	// "does this current instance hold the lock" so although we know what our own
	// session ID is, we don't check it matches here only that there is _some_
	// session in Consul holding the lock right now.
	held := pair.Session != ""
	value := string(pair.Value)
	return held, value, nil
}

func (c *ConsulLock) Info() (key, sessionid string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.key, c.sessionID
}
