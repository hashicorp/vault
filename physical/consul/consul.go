// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consul

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
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
	_ physical.Backend       = (*ConsulBackend)(nil)
	_ physical.HABackend     = (*ConsulBackend)(nil)
	_ physical.Lock          = (*ConsulLock)(nil)
	_ physical.Transactional = (*ConsulBackend)(nil)

	GetInTxnDisabledError = errors.New("get operations inside transactions are disabled in consul backend")
)

// ConsulBackend is a physical backend that stores data at specific
// prefix within Consul. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
// failGetInTxn is only used in tests.
type ConsulBackend struct {
	client          *api.Client
	path            string
	kv              *api.KV
	txn             *api.Txn
	permitPool      *physical.PermitPool
	consistencyMode string
	sessionTTL      string
	lockWaitTime    time.Duration
	failGetInTxn    *uint32
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
		path:            path,
		client:          client,
		kv:              client.KV(),
		txn:             client.Txn(),
		permitPool:      physical.NewPermitPool(maxParInt),
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

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)

	ok, resp, _, err := c.txn.Txn(ops, queryOpts)
	if ok && len(resp.Errors) == 0 && err == nil {
		available = true
	}

	return available
}

// Transaction is used to run multiple entries via a transaction.
func (c *ConsulBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	if len(txns) == 0 {
		return nil
	}
	defer metrics.MeasureSince([]string{"consul", "transaction"}, time.Now())

	failGetInTxn := atomic.LoadUint32(c.failGetInTxn)
	for _, t := range txns {
		if t.Operation == physical.GetOperation && failGetInTxn != 0 {
			return GetInTxnDisabledError
		}
	}

	ops := make([]*api.TxnOp, 0, len(txns))
	for _, t := range txns {
		o, err := c.makeApiTxn(t)
		if err != nil {
			return fmt.Errorf("error converting physical transactions into api transactions: %w", err)
		}

		ops = append(ops, o)
	}

	c.permitPool.Acquire()
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
		// Loop over results and cache them in a map. Note that we're only caching the first time we see a key,
		// which _should_ correspond to a Get operation, since we expect those come first in our txns slice.
		for _, txnr := range resp.Results {
			if len(txnr.KV.Value) > 0 {
				// We need to trim the Consul kv path (typically "vault/") from the key otherwise it won't
				// match the transaction entries we have.
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

// Put is used to insert or update an entry
func (c *ConsulBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"consul", "put"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	pair := &api.KVPair{
		Key:   c.path + entry.Key,
		Value: entry.Value,
	}

	writeOpts := &api.WriteOptions{}
	writeOpts = writeOpts.WithContext(ctx)

	_, err := c.kv.Put(pair, writeOpts)
	if err != nil {
		if strings.Contains(err.Error(), "Value exceeds") {
			return fmt.Errorf("%s: %w", physical.ErrValueTooLarge, err)
		}
		return err
	}
	return nil
}

// Get is used to fetch an entry
func (c *ConsulBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"consul", "get"}, time.Now())

	c.permitPool.Acquire()
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
	defer metrics.MeasureSince([]string{"consul", "delete"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	writeOpts := &api.WriteOptions{}
	writeOpts = writeOpts.WithContext(ctx)

	_, err := c.kv.Delete(c.path+key, writeOpts)
	return err
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

	c.permitPool.Acquire()
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
	// Create the lock
	opts := &api.LockOptions{
		Key:            c.path + key,
		Value:          []byte(value),
		SessionName:    "Vault Lock",
		MonitorRetries: 5,
		SessionTTL:     c.sessionTTL,
		LockWaitTime:   c.lockWaitTime,
	}
	lock, err := c.client.LockOpts(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create lock: %w", err)
	}
	cl := &ConsulLock{
		client:          c.client,
		key:             c.path + key,
		lock:            lock,
		consistencyMode: c.consistencyMode,
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

// ConsulLock is used to provide the Lock interface backed by Consul
type ConsulLock struct {
	client          *api.Client
	key             string
	lock            *api.Lock
	consistencyMode string
}

func (c *ConsulLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	return c.lock.Lock(stopCh)
}

func (c *ConsulLock) Unlock() error {
	return c.lock.Unlock()
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
	held := pair.Session != ""
	value := string(pair.Value)
	return held, value, nil
}
