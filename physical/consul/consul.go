package consul

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	"github.com/hashicorp/vault/sdk/physical"
	"golang.org/x/net/http2"
)

const (
	// consistencyModeDefault is the configuration value used to tell
	// consul to use default consistency.
	consistencyModeDefault = "default"

	// consistencyModeStrong is the configuration value used to tell
	// consul to use strong consistency.
	consistencyModeStrong = "strong"
)

// Verify ConsulBackend satisfies the correct interfaces
var _ physical.Backend = (*ConsulBackend)(nil)
var _ physical.HABackend = (*ConsulBackend)(nil)
var _ physical.Lock = (*ConsulLock)(nil)
var _ physical.Transactional = (*ConsulBackend)(nil)

// ConsulBackend is a physical backend that stores data at specific
// prefix within Consul. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type ConsulBackend struct {
	client          *api.Client
	path            string
	kv              *api.KV
	permitPool      *physical.PermitPool
	consistencyMode string

	sessionTTL   string
	lockWaitTime time.Duration
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
			return nil, errwrap.Wrapf("invalid session_ttl: {{err}}", err)
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
			return nil, errwrap.Wrapf("invalid lock_wait_time: {{err}}", err)
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
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
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
		// Use the parsed Address instead of the raw conf['address']
		tlsClientConfig, err := tlsutil.SetupTLSConfig(conf, consulConf.Address)
		if err != nil {
			return nil, err
		}

		consulConf.Transport.TLSClientConfig = tlsClientConfig
		if err := http2.ConfigureTransport(consulConf.Transport); err != nil {
			return nil, err
		}
		logger.Debug("configured TLS")
	}

	consulConf.HttpClient = &http.Client{Transport: consulConf.Transport}
	client, err := api.NewClient(consulConf)
	if err != nil {
		return nil, errwrap.Wrapf("client setup failed: {{err}}", err)
	}

	// Setup the backend
	c := &ConsulBackend{
		path:            path,
		client:          client,
		kv:              client.KV(),
		permitPool:      physical.NewPermitPool(maxParInt),
		consistencyMode: consistencyMode,

		sessionTTL:   sessionTTL,
		lockWaitTime: lockWaitTime,
	}
	return c, nil
}

// Used to run multiple entries via a transaction
func (c *ConsulBackend) Transaction(ctx context.Context, txns []*physical.TxnEntry) error {
	if len(txns) == 0 {
		return nil
	}

	ops := make([]*api.KVTxnOp, 0, len(txns))

	for _, op := range txns {
		cop := &api.KVTxnOp{
			Key: c.path + op.Entry.Key,
		}
		switch op.Operation {
		case physical.DeleteOperation:
			cop.Verb = api.KVDelete
		case physical.PutOperation:
			cop.Verb = api.KVSet
			cop.Value = op.Entry.Value
		default:
			return fmt.Errorf("%q is not a supported transaction operation", op.Operation)
		}

		ops = append(ops, cop)
	}

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	queryOpts := &api.QueryOptions{}
	queryOpts = queryOpts.WithContext(ctx)

	ok, resp, _, err := c.kv.Txn(ops, queryOpts)
	if err != nil {
		if strings.Contains(err.Error(), "is too large") {
			return errwrap.Wrapf(fmt.Sprintf("%s: {{err}}", physical.ErrValueTooLarge), err)
		}
		return err
	}
	if ok && len(resp.Errors) == 0 {
		return nil
	}

	var retErr *multierror.Error
	for _, res := range resp.Errors {
		retErr = multierror.Append(retErr, errors.New(res.What))
	}

	return retErr
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
			return errwrap.Wrapf(fmt.Sprintf("%s: {{err}}", physical.ErrValueTooLarge), err)
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

// Lock is used for mutual exclusion based on the given key.
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
		return nil, errwrap.Wrapf("failed to create lock: {{err}}", err)
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
