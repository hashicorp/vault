package etcd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/physical"
	log "github.com/mgutz/logxi/v1"
	"golang.org/x/net/context"
)

// EtcdBackend is a physical backend that stores data at specific
// prefix within etcd. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type EtcdBackend struct {
	logger    log.Logger
	path      string
	haEnabled bool

	permitPool *physical.PermitPool

	etcd *clientv3.Client
}

const (
	// etcd3 default lease duration is 60s. set to 15s for faster recovery.
	etcd3LockTimeoutInSeconds = 15
	// etcd3 default request timeout is set to 5s. It should be long enough
	// for most cases, even with internal retry.
	etcd3RequestTimeout = 5 * time.Second
)

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
			CAFile:   ca,
			CertFile: cert,
			KeyFile:  key,
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

	etcd, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	ssync, ok := conf["sync"]
	if !ok {
		ssync = "true"
	}
	sync, err := strconv.ParseBool(ssync)
	if err != nil {
		return nil, fmt.Errorf("value of 'sync' (%v) could not be understood", err)
	}

	if sync {
		ctx, cancel := context.WithTimeout(context.Background(), etcd3RequestTimeout)
		err := etcd.Sync(ctx)
		cancel()
		if err != nil {
			return nil, err
		}
	}

	return &EtcdBackend{
		path:       path,
		etcd:       etcd,
		permitPool: physical.NewPermitPool(physical.DefaultParallelOperations),
		logger:     logger,
		haEnabled:  haEnabledBool,
	}, nil
}

func (c *EtcdBackend) Put(entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"etcd", "put"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	ctx, cancel := context.WithTimeout(context.Background(), etcd3RequestTimeout)
	defer cancel()
	_, err := c.etcd.Put(ctx, path.Join(c.path, entry.Key), string(entry.Value))
	return err
}

func (c *EtcdBackend) Get(key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"etcd", "get"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	ctx, cancel := context.WithTimeout(context.Background(), etcd3RequestTimeout)
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

func (c *EtcdBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"etcd", "delete"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	ctx, cancel := context.WithTimeout(context.Background(), etcd3RequestTimeout)
	defer cancel()
	_, err := c.etcd.Delete(ctx, path.Join(c.path, key))
	if err != nil {
		return err
	}
	return nil
}

func (c *EtcdBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"etcd", "list"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	ctx, cancel := context.WithTimeout(context.Background(), etcd3RequestTimeout)
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

// EtcdLock emplements a lock using and etcd backend.
type EtcdLock struct {
	lock sync.Mutex
	held bool

	etcdSession *concurrency.Session
	etcdMu      *concurrency.Mutex

	prefix string
	value  string

	etcd *clientv3.Client
}

// Lock is used for mutual exclusion based on the given key.
func (c *EtcdBackend) LockWith(key, value string) (physical.Lock, error) {
	session, err := concurrency.NewSession(c.etcd, concurrency.WithTTL(etcd3LockTimeoutInSeconds))
	if err != nil {
		return nil, err
	}

	p := path.Join(c.path, key)
	return &EtcdLock{
		etcdSession: session,
		etcdMu:      concurrency.NewMutex(session, p),
		prefix:      p,
		value:       value,
		etcd:        c.etcd,
	}, nil
}

func (c *EtcdLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.held {
		return nil, EtcdLockHeldError
	}

	select {
	case _, ok := <-c.etcdSession.Done():
		if !ok {
			// The session's done channel is closed, so the session is over,
			// and we need a new one
			session, err := concurrency.NewSession(c.etcd, concurrency.WithTTL(etcd3LockTimeoutInSeconds))
			if err != nil {
				return nil, err
			}
			c.etcdSession = session
			c.etcdMu = concurrency.NewMutex(session, c.prefix)
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

	pctx, cancel := context.WithTimeout(context.Background(), etcd3RequestTimeout)
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

	ctx, cancel := context.WithTimeout(context.Background(), etcd3RequestTimeout)
	defer cancel()
	return c.etcdMu.Unlock(ctx)
}

func (c *EtcdLock) Value() (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), etcd3RequestTimeout)
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
