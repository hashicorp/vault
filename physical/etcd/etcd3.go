package etcd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/physical"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	"go.etcd.io/etcd/pkg/transport"
)

// EtcdBackend is a physical backend that stores data at specific
// prefix within etcd. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type EtcdBackend struct {
	logger         log.Logger
	path           string
	haEnabled      bool
	lockTimeout    time.Duration
	requestTimeout time.Duration

	permitPool *physical.PermitPool

	etcd *clientv3.Client
}

// Verify EtcdBackend satisfies the correct interfaces
var _ physical.Backend = (*EtcdBackend)(nil)
var _ physical.HABackend = (*EtcdBackend)(nil)
var _ physical.Lock = (*EtcdLock)(nil)

// newEtcd3Backend constructs a etcd3 backend.
func newEtcd3Backend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Get the etcd path form the configuration.
	path, ok := conf["path"]
	if !ok {
		path = "/vault"
	}

	// Ensure path is prefixed.
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	endpoints, err := getEtcdEndpoints(conf)
	if err != nil {
		return nil, err
	}

	cfg := clientv3.Config{
		Endpoints: endpoints,
	}

	haEnabled := os.Getenv("ETCD_HA_ENABLED")
	if haEnabled == "" {
		haEnabled = conf["ha_enabled"]
	}
	if haEnabled == "" {
		haEnabled = "false"
	}
	haEnabledBool, err := strconv.ParseBool(haEnabled)
	if err != nil {
		return nil, fmt.Errorf("value [%v] of 'ha_enabled' could not be understood", haEnabled)
	}

	cert, hasCert := conf["tls_cert_file"]
	key, hasKey := conf["tls_key_file"]
	ca, hasCa := conf["tls_ca_file"]
	if (hasCert && hasKey) || hasCa {
		tls := transport.TLSInfo{
			TrustedCAFile: ca,
			CertFile:      cert,
			KeyFile:       key,
		}

		tlscfg, err := tls.ClientConfig()
		if err != nil {
			return nil, err
		}
		cfg.TLS = tlscfg
	}

	// Set credentials.
	username := os.Getenv("ETCD_USERNAME")
	if username == "" {
		username, _ = conf["username"]
	}

	password := os.Getenv("ETCD_PASSWORD")
	if password == "" {
		password, _ = conf["password"]
	}

	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}

	if maxReceive, ok := conf["max_receive_size"]; ok {
		// grpc converts this to uint32 internally, so parse as that to avoid passing invalid values
		val, err := strconv.ParseUint(maxReceive, 10, 32)
		if err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("value of 'max_receive_size' (%v) could not be understood: {{err}}", maxReceive), err)
		}
		cfg.MaxCallRecvMsgSize = int(val)
	}

	etcd, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	sReqTimeout := conf["request_timeout"]
	if sReqTimeout == "" {
		// etcd3 default request timeout is set to 5s. It should be long enough
		// for most cases, even with internal retry.
		sReqTimeout = "5s"
	}
	reqTimeout, err := parseutil.ParseDurationSecond(sReqTimeout)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("value [%v] of 'request_timeout' could not be understood: {{err}}", sReqTimeout), err)
	}

	ssync, ok := conf["sync"]
	if !ok {
		ssync = "true"
	}
	sync, err := strconv.ParseBool(ssync)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("value of 'sync' (%v) could not be understood: {{err}}", ssync), err)
	}

	if sync {
		ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
		err := etcd.Sync(ctx)
		cancel()
		if err != nil {
			return nil, err
		}
	}

	sLock := conf["lock_timeout"]
	if sLock == "" {
		// etcd3 default lease duration is 60s. set to 15s for faster recovery.
		sLock = "15s"
	}
	lock, err := parseutil.ParseDurationSecond(sLock)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("value [%v] of 'lock_timeout' could not be understood: {{err}}", sLock), err)
	}

	return &EtcdBackend{
		path:           path,
		etcd:           etcd,
		permitPool:     physical.NewPermitPool(physical.DefaultParallelOperations),
		logger:         logger,
		haEnabled:      haEnabledBool,
		lockTimeout:    lock,
		requestTimeout: reqTimeout,
	}, nil
}

func (c *EtcdBackend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"etcd", "put"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	_, err := c.etcd.Put(ctx, path.Join(c.path, entry.Key), string(entry.Value))
	return err
}

func (c *EtcdBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"etcd", "get"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	resp, err := c.etcd.Get(ctx, path.Join(c.path, key))
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	if len(resp.Kvs) > 1 {
		return nil, errors.New("unexpected number of keys from a get request")
	}
	return &physical.Entry{
		Key:   key,
		Value: resp.Kvs[0].Value,
	}, nil
}

func (c *EtcdBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"etcd", "delete"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	_, err := c.etcd.Delete(ctx, path.Join(c.path, key))
	if err != nil {
		return err
	}
	return nil
}

func (c *EtcdBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"etcd", "list"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	prefix = path.Join(c.path, prefix) + "/"
	resp, err := c.etcd.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, kv := range resp.Kvs {
		key := strings.TrimPrefix(string(kv.Key), prefix)
		key = strings.TrimPrefix(key, "/")

		if len(key) == 0 {
			continue
		}

		if i := strings.Index(key, "/"); i == -1 {
			keys = append(keys, key)
		} else if i != -1 {
			keys = strutil.AppendIfMissing(keys, key[:i+1])
		}
	}
	return keys, nil
}

func (e *EtcdBackend) HAEnabled() bool {
	return e.haEnabled
}

// EtcdLock implements a lock using and etcd backend.
type EtcdLock struct {
	lock           sync.Mutex
	held           bool
	timeout        time.Duration
	requestTimeout time.Duration

	etcdSession *concurrency.Session
	etcdMu      *concurrency.Mutex

	prefix string
	value  string

	etcd *clientv3.Client
}

// Lock is used for mutual exclusion based on the given key.
func (c *EtcdBackend) LockWith(key, value string) (physical.Lock, error) {
	p := path.Join(c.path, key)
	return &EtcdLock{
		prefix:         p,
		value:          value,
		etcd:           c.etcd,
		timeout:        c.lockTimeout,
		requestTimeout: c.requestTimeout,
	}, nil
}

func (c *EtcdLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.etcdMu == nil {
		if err := c.initMu(); err != nil {
			return nil, err
		}
	}

	if c.held {
		return nil, EtcdLockHeldError
	}

	select {
	case _, ok := <-c.etcdSession.Done():
		if !ok {
			// The session's done channel is closed, so the session is over,
			// and we need a new lock with a new session.
			if err := c.initMu(); err != nil {
				return nil, err
			}
		}
	default:
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stopCh
		cancel()
	}()
	if err := c.etcdMu.Lock(ctx); err != nil {
		if err == context.Canceled {
			return nil, nil
		}
		return nil, err
	}

	pctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	if _, err := c.etcd.Put(pctx, c.etcdMu.Key(), c.value, clientv3.WithLease(c.etcdSession.Lease())); err != nil {
		return nil, err
	}

	c.held = true

	return c.etcdSession.Done(), nil
}

func (c *EtcdLock) Unlock() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.held {
		return EtcdLockNotHeldError
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()
	return c.etcdMu.Unlock(ctx)
}

func (c *EtcdLock) Value() (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.requestTimeout)
	defer cancel()

	resp, err := c.etcd.Get(ctx,
		c.prefix, clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByCreateRevision, clientv3.SortAscend))

	if err != nil {
		return false, "", err
	}
	if len(resp.Kvs) == 0 {
		return false, "", nil
	}

	return true, string(resp.Kvs[0].Value), nil
}

func (c *EtcdLock) initMu() error {
	session, err := concurrency.NewSession(c.etcd, concurrency.WithTTL(int(c.timeout.Seconds())))
	if err != nil {
		return err
	}
	c.etcdSession = session
	c.etcdMu = concurrency.NewMutex(session, c.prefix)
	return nil
}
