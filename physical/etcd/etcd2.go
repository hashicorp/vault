package etcd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	metrics "github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/physical"
	"go.etcd.io/etcd/client"
	"go.etcd.io/etcd/pkg/transport"
)

const (
	// Ideally, this prefix would match the "_" used in the file backend, but
	// that prefix has special meaning in etcd. Specifically, it excludes those
	// entries from directory listings.
	Etcd2NodeFilePrefix = "."

	// The lock prefix can (and probably should) cause an entry to be excluded
	// from directory listings, so "_" works here.
	Etcd2NodeLockPrefix = "_"

	// The delimiter is the same as the `-C` flag of etcdctl.
	Etcd2MachineDelimiter = ","

	// The lock TTL matches the default that Consul API uses, 15 seconds.
	Etcd2LockTTL = 15 * time.Second

	// The amount of time to wait between the semaphore key renewals
	Etcd2LockRenewInterval = 5 * time.Second

	// The amount of time to wait if a watch fails before trying again.
	Etcd2WatchRetryInterval = time.Second

	// The number of times to re-try a failed watch before signaling that leadership is lost.
	Etcd2WatchRetryMax = 5
)

// Etcd2Backend is a physical backend that stores data at specific
// prefix within etcd. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type Etcd2Backend struct {
	path       string
	kAPI       client.KeysAPI
	permitPool *physical.PermitPool
	logger     log.Logger
	haEnabled  bool
}

// Verify Etcd2Backend satisfies the correct interfaces
var _ physical.Backend = (*Etcd2Backend)(nil)
var _ physical.HABackend = (*Etcd2Backend)(nil)
var _ physical.Lock = (*Etcd2Lock)(nil)

func newEtcd2Backend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Get the etcd path form the configuration.
	path, ok := conf["path"]
	if !ok {
		path = "/vault"
	}

	// Ensure path is prefixed.
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	c, err := newEtcdV2Client(conf)
	if err != nil {
		return nil, err
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

	// Should we sync the cluster state? There are three available options
	// for our client library: don't sync (required for some proxies), sync
	// once, or sync periodically with AutoSync.  We currently support the
	// first two.
	sync, ok := conf["sync"]
	if !ok {
		sync = "yes"
	}
	switch sync {
	case "yes", "true", "y", "1":
		ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
		syncErr := c.Sync(ctx)
		cancel()
		if syncErr != nil {
			return nil, multierror.Append(EtcdSyncClusterError, syncErr)
		}
	case "no", "false", "n", "0":
	default:
		return nil, fmt.Errorf("value of 'sync' could not be understood")
	}

	kAPI := client.NewKeysAPI(c)

	// Setup the backend.
	return &Etcd2Backend{
		path:       path,
		kAPI:       kAPI,
		permitPool: physical.NewPermitPool(physical.DefaultParallelOperations),
		logger:     logger,
		haEnabled:  haEnabledBool,
	}, nil
}

func newEtcdV2Client(conf map[string]string) (client.Client, error) {
	endpoints, err := getEtcdEndpoints(conf)
	if err != nil {
		return nil, err
	}

	// Create a new client from the supplied address and attempt to sync with the
	// cluster.
	var cTransport client.CancelableTransport
	cert, hasCert := conf["tls_cert_file"]
	key, hasKey := conf["tls_key_file"]
	ca, hasCa := conf["tls_ca_file"]
	if (hasCert && hasKey) || hasCa {
		var transportErr error
		tls := transport.TLSInfo{
			TrustedCAFile: ca,
			CertFile:      cert,
			KeyFile:       key,
		}
		cTransport, transportErr = transport.NewTransport(tls, 30*time.Second)

		if transportErr != nil {
			return nil, transportErr
		}
	} else {
		cTransport = client.DefaultTransport
	}

	cfg := client.Config{
		Endpoints: endpoints,
		Transport: cTransport,
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

	return client.New(cfg)
}

// Put is used to insert or update an entry.
func (c *Etcd2Backend) Put(ctx context.Context, entry *physical.Entry) error {
	defer metrics.MeasureSince([]string{"etcd", "put"}, time.Now())
	value := base64.StdEncoding.EncodeToString(entry.Value)

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	_, err := c.kAPI.Set(context.Background(), c.nodePath(entry.Key), value, nil)
	return err
}

// Get is used to fetch an entry.
func (c *Etcd2Backend) Get(ctx context.Context, key string) (*physical.Entry, error) {
	defer metrics.MeasureSince([]string{"etcd", "get"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	getOpts := &client.GetOptions{
		Recursive: false,
		Sort:      false,
	}
	response, err := c.kAPI.Get(context.Background(), c.nodePath(key), getOpts)
	if err != nil {
		if errorIsMissingKey(err) {
			return nil, nil
		}
		return nil, err
	}

	// Decode the stored value from base-64.
	value, err := base64.StdEncoding.DecodeString(response.Node.Value)
	if err != nil {
		return nil, err
	}

	// Construct and return a new entry.
	return &physical.Entry{
		Key:   key,
		Value: value,
	}, nil
}

// Delete is used to permanently delete an entry.
func (c *Etcd2Backend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"etcd", "delete"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	// Remove the key, non-recursively.
	delOpts := &client.DeleteOptions{
		Recursive: false,
	}
	_, err := c.kAPI.Delete(context.Background(), c.nodePath(key), delOpts)
	if err != nil && !errorIsMissingKey(err) {
		return err
	}
	return nil
}

// List is used to list all the keys under a given prefix, up to the next
// prefix.
func (c *Etcd2Backend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"etcd", "list"}, time.Now())

	// Set a directory path from the given prefix.
	path := c.nodePathDir(prefix)

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	// Get the directory, non-recursively, from etcd. If the directory is
	// missing, we just return an empty list of contents.
	getOpts := &client.GetOptions{
		Recursive: false,
		Sort:      true,
	}
	response, err := c.kAPI.Get(context.Background(), path, getOpts)
	if err != nil {
		if errorIsMissingKey(err) {
			return []string{}, nil
		}
		return nil, err
	}

	out := make([]string, len(response.Node.Nodes))
	for i, node := range response.Node.Nodes {

		// etcd keys include the full path, so let's trim the prefix directory
		// path.
		name := strings.TrimPrefix(node.Key, path)

		// Check if this node is itself a directory. If it is, add a trailing
		// slash; if it isn't remove the node file prefix.
		if node.Dir {
			out[i] = name + "/"
		} else {
			out[i] = name[1:]
		}
	}
	return out, nil
}

// nodePath returns an etcd filepath based on the given key.
func (b *Etcd2Backend) nodePath(key string) string {
	return filepath.Join(b.path, filepath.Dir(key), Etcd2NodeFilePrefix+filepath.Base(key))
}

// nodePathDir returns an etcd directory path based on the given key.
func (b *Etcd2Backend) nodePathDir(key string) string {
	return filepath.Join(b.path, key) + "/"
}

// nodePathLock returns an etcd directory path used specifically for semaphore
// indices based on the given key.
func (b *Etcd2Backend) nodePathLock(key string) string {
	return filepath.Join(b.path, filepath.Dir(key), Etcd2NodeLockPrefix+filepath.Base(key)+"/")
}

// Lock is used for mutual exclusion based on the given key.
func (c *Etcd2Backend) LockWith(key, value string) (physical.Lock, error) {
	return &Etcd2Lock{
		kAPI:            c.kAPI,
		value:           value,
		semaphoreDirKey: c.nodePathLock(key),
	}, nil
}

// HAEnabled indicates whether the HA functionality should be exposed.
// Currently always returns true.
func (e *Etcd2Backend) HAEnabled() bool {
	return e.haEnabled
}

// Etcd2Lock implements a lock using and Etcd2 backend.
type Etcd2Lock struct {
	kAPI                                 client.KeysAPI
	value, semaphoreDirKey, semaphoreKey string
	lock                                 sync.Mutex
}

// addSemaphoreKey acquires a new ordered semaphore key.
func (c *Etcd2Lock) addSemaphoreKey() (string, uint64, error) {
	// CreateInOrder is an atomic operation that can be used to enqueue a
	// request onto a semaphore. In the rest of the comments, we refer to the
	// resulting key as a "semaphore key".
	// https://coreos.com/etcd/docs/2.0.8/api.html#atomically-creating-in-order-keys
	opts := &client.CreateInOrderOptions{
		TTL: Etcd2LockTTL,
	}
	response, err := c.kAPI.CreateInOrder(context.Background(), c.semaphoreDirKey, c.value, opts)
	if err != nil {
		return "", 0, err
	}
	return response.Node.Key, response.Index, nil
}

// renewSemaphoreKey renews an existing semaphore key.
func (c *Etcd2Lock) renewSemaphoreKey() (string, uint64, error) {
	setOpts := &client.SetOptions{
		TTL:       Etcd2LockTTL,
		PrevExist: client.PrevExist,
	}
	response, err := c.kAPI.Set(context.Background(), c.semaphoreKey, c.value, setOpts)
	if err != nil {
		return "", 0, err
	}
	return response.Node.Key, response.Index, nil
}

// getSemaphoreKey determines which semaphore key holder has acquired the lock
// and its value.
func (c *Etcd2Lock) getSemaphoreKey() (string, string, uint64, error) {
	// Get the list of waiters in order to see if we are next.
	getOpts := &client.GetOptions{
		Recursive: false,
		Sort:      true,
	}
	response, err := c.kAPI.Get(context.Background(), c.semaphoreDirKey, getOpts)
	if err != nil {
		return "", "", 0, err
	}

	// Make sure the list isn't empty.
	if response.Node.Nodes.Len() == 0 {
		return "", "", response.Index, nil
	}
	return response.Node.Nodes[0].Key, response.Node.Nodes[0].Value, response.Index, nil
}

// isHeld determines if we are the current holders of the lock.
func (c *Etcd2Lock) isHeld() (bool, error) {
	if c.semaphoreKey == "" {
		return false, nil
	}

	// Get the key of the current holder of the lock.
	currentSemaphoreKey, _, _, err := c.getSemaphoreKey()
	if err != nil {
		return false, err
	}
	return c.semaphoreKey == currentSemaphoreKey, nil
}

// assertHeld determines whether or not we are the current holders of the lock
// and returns an Etcd2LockNotHeldError if we are not.
func (c *Etcd2Lock) assertHeld() error {
	held, err := c.isHeld()
	if err != nil {
		return err
	}

	// Check if we don't hold the lock.
	if !held {
		return EtcdLockNotHeldError
	}
	return nil
}

// assertNotHeld determines whether or not we are the current holders of the
// lock and returns an Etcd2LockHeldError if we are.
func (c *Etcd2Lock) assertNotHeld() error {
	held, err := c.isHeld()
	if err != nil {
		return err
	}

	// Check if we hold the lock.
	if held {
		return EtcdLockHeldError
	}
	return nil
}

// periodically renew our semaphore key so that it doesn't expire
func (c *Etcd2Lock) periodicallyRenewSemaphoreKey(stopCh chan struct{}) {
	for {
		select {
		case <-time.After(Etcd2LockRenewInterval):
			c.renewSemaphoreKey()
		case <-stopCh:
			return
		}
	}
}

// watchForKeyRemoval continuously watches a single non-directory key starting
// from the provided etcd index and closes the provided channel when it's
// deleted, expires, or appears to be missing.
func (c *Etcd2Lock) watchForKeyRemoval(key string, etcdIndex uint64, closeCh chan struct{}) {
	retries := Etcd2WatchRetryMax

	for {
		// Start a non-recursive watch of the given key.
		w := c.kAPI.Watcher(key, &client.WatcherOptions{AfterIndex: etcdIndex, Recursive: false})
		response, err := w.Next(context.TODO())
		if err != nil {

			// If the key is just missing, we can exit the loop.
			if errorIsMissingKey(err) {
				break
			}

			// If the error is something else, there's nothing we can do but retry
			// the watch. Check that we still have retries left.
			retries -= 1
			if retries == 0 {
				break
			}

			// Sleep for a period of time to avoid slamming etcd.
			time.Sleep(Etcd2WatchRetryInterval)
			continue
		}

		// Check if the key we are concerned with has been removed. If it has, we
		// can exit the loop.
		if response.Node.Key == key &&
			(response.Action == "delete" || response.Action == "expire") {
			break
		}

		// Update the etcd index.
		etcdIndex = response.Index + 1
	}

	// Regardless of what happened, we need to close the close channel.
	close(closeCh)
}

// Lock attempts to acquire the lock by waiting for a new semaphore key in etcd
// to become the first in the queue and will block until it is successful or
// it receives a signal on the provided channel. The returned channel will be
// closed when the lock is lost, either by an explicit call to Unlock or by
// the associated semaphore key in etcd otherwise being deleted or expiring.
//
// If the lock is currently held by this instance of Etcd2Lock, Lock will
// return an Etcd2LockHeldError error.
func (c *Etcd2Lock) Lock(stopCh <-chan struct{}) (doneCh <-chan struct{}, retErr error) {
	// Get the local lock before interacting with etcd.
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if the lock is already held.
	if err := c.assertNotHeld(); err != nil {
		return nil, err
	}

	// Add a new semaphore key that we will track.
	semaphoreKey, _, err := c.addSemaphoreKey()
	if err != nil {
		return nil, err
	}
	c.semaphoreKey = semaphoreKey

	// Get the current semaphore key.
	currentSemaphoreKey, _, currentEtcdIndex, err := c.getSemaphoreKey()
	if err != nil {
		return nil, err
	}

	// Create an etcd-compatible boolean stop channel from the provided
	// interface stop channel.
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-stopCh
		cancel()
	}()
	defer cancel()

	// Create a channel to signal when we lose the semaphore key.
	done := make(chan struct{})
	defer func() {
		if retErr != nil {
			close(done)
		}
	}()

	go c.periodicallyRenewSemaphoreKey(done)

	// Loop until the we current semaphore key matches ours.
	for semaphoreKey != currentSemaphoreKey {
		var err error

		// Start a watch of the entire lock directory
		w := c.kAPI.Watcher(c.semaphoreDirKey, &client.WatcherOptions{AfterIndex: currentEtcdIndex, Recursive: true})
		response, err := w.Next(ctx)
		if err != nil {

			// If the error is not an etcd error, we can assume it's a notification
			// of the stop channel having closed. In this scenario, we also want to
			// remove our semaphore key as we are no longer waiting to acquire the
			// lock.
			if _, ok := err.(*client.Error); !ok {
				delOpts := &client.DeleteOptions{
					Recursive: false,
				}
				_, err = c.kAPI.Delete(context.Background(), c.semaphoreKey, delOpts)
			}
			return nil, err
		}

		// Make sure the index we are waiting for has not been removed. If it has,
		// this is an error and nothing else needs to be done.
		if response.Node.Key == semaphoreKey &&
			(response.Action == "delete" || response.Action == "expire") {
			return nil, EtcdSemaphoreKeyRemovedError
		}

		// Get the current semaphore key and etcd index.
		currentSemaphoreKey, _, currentEtcdIndex, err = c.getSemaphoreKey()
		if err != nil {
			return nil, err
		}
	}

	go c.watchForKeyRemoval(c.semaphoreKey, currentEtcdIndex, done)
	return done, nil
}

// Unlock releases the lock by deleting the associated semaphore key in etcd.
//
// If the lock is not currently held by this instance of Etcd2Lock, Unlock will
// return an Etcd2LockNotHeldError error.
func (c *Etcd2Lock) Unlock() error {
	// Get the local lock before interacting with etcd.
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check that the lock is held.
	if err := c.assertHeld(); err != nil {
		return err
	}

	// Delete our semaphore key.
	delOpts := &client.DeleteOptions{
		Recursive: false,
	}
	if _, err := c.kAPI.Delete(context.Background(), c.semaphoreKey, delOpts); err != nil {
		return err
	}
	return nil
}

// Value checks whether or not the lock is held by any instance of Etcd2Lock,
// including this one, and returns the current value.
func (c *Etcd2Lock) Value() (bool, string, error) {
	semaphoreKey, semaphoreValue, _, err := c.getSemaphoreKey()
	if err != nil {
		return false, "", err
	}

	if semaphoreKey == "" {
		return false, "", nil
	}
	return true, semaphoreValue, nil
}

// errorIsMissingKey returns true if the given error is an etcd error with an
// error code corresponding to a missing key.
func errorIsMissingKey(err error) bool {
	etcdErr, ok := err.(client.Error)
	return ok && etcdErr.Code == client.ErrorCodeKeyNotFound
}
