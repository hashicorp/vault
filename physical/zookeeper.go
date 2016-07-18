package physical

import (
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	// ZKNodeFilePrefix is prefixed to any "files" in ZooKeeper,
	// so that they do not collide with directory entries. Otherwise,
	// we cannot delete a file if the path is a full-prefix of another
	// key.
	ZKNodeFilePrefix = "_"
)

// ZookeeperBackend is a physical backend that stores data at specific
// prefix within Zookeeper. It is used in production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type ZookeeperBackend struct {
	path   string
	client *zk.Conn
	acl    []zk.ACL
	logger *log.Logger
}

// newZookeeperBackend constructs a Zookeeper backend using the given API client
// and the prefix in the KV store.
func newZookeeperBackend(conf map[string]string, logger *log.Logger) (Backend, error) {
	// Get the path in Zookeeper
	path, ok := conf["path"]
	if !ok {
		path = "vault/"
	}

	// Ensure path is suffixed and prefixed (zk requires prefix /)
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Configure the client, default to localhost instance
	var machines string
	machines, ok = conf["address"]
	if !ok {
		machines = "localhost:2181"
	}

	// zNode owner and schema.
	var owner string
	var schema string
	var schemaAndOwner string
	schemaAndOwner, ok = conf["znode_owner"]
	if !ok {
		owner = "anyone"
		schema = "world"
	} else {
		parsedSchemaAndOwner := strings.SplitN(schemaAndOwner, ":", 2)
		if len(parsedSchemaAndOwner) != 2 {
			return nil, fmt.Errorf("znode_owner expected format is 'schema:owner'")
		} else {
			schema = parsedSchemaAndOwner[0]
			owner = parsedSchemaAndOwner[1]

			// znode_owner is in config and structured correctly - but does it make any sense?
			// Either 'owner' or 'schema' was set but not both - this seems like a failed attempt
			// (e.g. ':MyUser' which omit the schema, or ':' omitting both)
			if owner == "" || schema == "" {
				return nil, fmt.Errorf("znode_owner expected format is 'schema:auth'")
			}
		}
	}

	acl := []zk.ACL{{zk.PermAll, schema, owner}}

	// Authnetication info
	var schemaAndUser string
	var useAddAuth bool
	schemaAndUser, useAddAuth = conf["auth_info"]
	if useAddAuth {
		parsedSchemaAndUser := strings.SplitN(schemaAndUser, ":", 2)
		if len(parsedSchemaAndUser) != 2 {
			return nil, fmt.Errorf("auth_info expected format is 'schema:auth'")
		} else {
			schema = parsedSchemaAndUser[0]
			owner = parsedSchemaAndUser[1]

			// auth_info is in config and structured correctly - but does it make any sense?
			// Either 'owner' or 'schema' was set but not both - this seems like a failed attempt
			// (e.g. ':MyUser' which omit the schema, or ':' omitting both)
			if owner == "" || schema == "" {
				return nil, fmt.Errorf("auth_info expected format is 'schema:auth'")
			}
		}
	}

	// We have all of the configuration in hand - let's try and connect to ZK
	client, _, err := zk.Connect(strings.Split(machines, ","), time.Second)
	if err != nil {
		return nil, fmt.Errorf("client setup failed: %v", err)
	}

	// ZK AddAuth API if the user asked for it
	if useAddAuth {
		err = client.AddAuth(schema, []byte(owner))
		if err != nil {
			return nil, fmt.Errorf("Zookeeper rejected authentication information provided at auth_info: %v", err)
		}
	}

	// Setup the backend
	c := &ZookeeperBackend{
		path:   path,
		client: client,
		acl:    acl,
		logger: logger,
	}
	return c, nil
}

// ensurePath is used to create each node in the path hierarchy.
// We avoid calling this optimistically, and invoke it when we get
// an error during an operation
func (c *ZookeeperBackend) ensurePath(path string, value []byte) error {
	nodes := strings.Split(path, "/")
	fullPath := ""
	for index, node := range nodes {
		if strings.TrimSpace(node) != "" {
			fullPath += "/" + node
			isLastNode := index+1 == len(nodes)

			// set parent nodes to nil, leaf to value
			// this block reduces round trips by being smart on the leaf create/set
			if exists, _, _ := c.client.Exists(fullPath); !isLastNode && !exists {
				if _, err := c.client.Create(fullPath, nil, int32(0), c.acl); err != nil {
					return err
				}
			} else if isLastNode && !exists {
				if _, err := c.client.Create(fullPath, value, int32(0), c.acl); err != nil {
					return err
				}
			} else if isLastNode && exists {
				if _, err := c.client.Set(fullPath, value, int32(-1)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// nodePath returns an etcd filepath based on the given key.
func (c *ZookeeperBackend) nodePath(key string) string {
	return filepath.Join(c.path, filepath.Dir(key), ZKNodeFilePrefix+filepath.Base(key))
}

// Put is used to insert or update an entry
func (c *ZookeeperBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"zookeeper", "put"}, time.Now())

	// Attempt to set the full path
	fullPath := c.nodePath(entry.Key)
	_, err := c.client.Set(fullPath, entry.Value, -1)

	// If we get ErrNoNode, we need to construct the path hierarchy
	if err == zk.ErrNoNode {
		return c.ensurePath(fullPath, entry.Value)
	}
	return err
}

// Get is used to fetch an entry
func (c *ZookeeperBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"zookeeper", "get"}, time.Now())

	// Attempt to read the full path
	fullPath := c.nodePath(key)
	value, _, err := c.client.Get(fullPath)

	// Ignore if the node does not exist
	if err == zk.ErrNoNode {
		err = nil
	}
	if err != nil {
		return nil, err
	}

	// Handle a non-existing value
	if value == nil {
		return nil, nil
	}
	ent := &Entry{
		Key:   key,
		Value: value,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (c *ZookeeperBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"zookeeper", "delete"}, time.Now())

	// Delete the full path
	fullPath := c.nodePath(key)
	err := c.client.Delete(fullPath, -1)

	// Mask if the node does not exist
	if err == zk.ErrNoNode {
		err = nil
	}
	return err
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (c *ZookeeperBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"zookeeper", "list"}, time.Now())

	// Query the children at the full path
	fullPath := strings.TrimSuffix(c.path+prefix, "/")
	result, _, err := c.client.Children(fullPath)

	// If the path nodes are missing, no children!
	if err == zk.ErrNoNode {
		return []string{}, nil
	}

	children := []string{}
	for _, key := range result {
		// Check if this entry has any child entries,
		// and append the slash which is what Vault depends on
		// for iteration
		nodeChildren, _, _ := c.client.Children(fullPath + "/" + key)
		if nodeChildren != nil && len(nodeChildren) > 0 {
			children = append(children, key+"/")
		} else {
			children = append(children, key[1:])
		}
	}
	sort.Strings(children)
	return children, nil
}

// LockWith is used for mutual exclusion based on the given key.
func (c *ZookeeperBackend) LockWith(key, value string) (Lock, error) {
	l := &ZookeeperHALock{
		in:    c,
		key:   key,
		value: value,
	}
	return l, nil
}

// HAEnabled indicates whether the HA functionality should be exposed.
// Currently always returns true.
func (c *ZookeeperBackend) HAEnabled() bool {
	return true
}

// ZookeeperHALock is a Zookeeper Lock implementation for the HABackend
type ZookeeperHALock struct {
	in    *ZookeeperBackend
	key   string
	value string

	held      bool
	localLock sync.Mutex
	leaderCh  chan struct{}
	zkLock    *zk.Lock
}

func (i *ZookeeperHALock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	i.localLock.Lock()
	defer i.localLock.Unlock()
	if i.held {
		return nil, fmt.Errorf("lock already held")
	}

	// Attempt an async acquisition
	didLock := make(chan struct{})
	failLock := make(chan error, 1)
	releaseCh := make(chan bool, 1)
	lockpath := i.in.nodePath(i.key)
	go i.attemptLock(lockpath, didLock, failLock, releaseCh)

	// Wait for lock acquisition, failure, or shutdown
	select {
	case <-didLock:
		releaseCh <- false
	case err := <-failLock:
		return nil, err
	case <-stopCh:
		releaseCh <- true
		return nil, nil
	}

	// Create the leader channel
	i.held = true
	i.leaderCh = make(chan struct{})

	// Watch for Events which could result in loss of our zkLock and close(i.leaderCh)
	currentVal, _, lockeventCh, err := i.in.client.GetW(lockpath)
	if err != nil {
		return nil, fmt.Errorf("unable to watch HA lock: %v", err)
	}
	if i.value != string(currentVal) {
		return nil, fmt.Errorf("lost HA lock immediately before watch")
	}
	go i.monitorLock(lockeventCh, i.leaderCh)

	return i.leaderCh, nil
}

func (i *ZookeeperHALock) attemptLock(lockpath string, didLock chan struct{}, failLock chan error, releaseCh chan bool) {
	// Wait to acquire the lock in ZK
	lock := zk.NewLock(i.in.client, lockpath, i.in.acl)
	err := lock.Lock()
	if err != nil {
		failLock <- err
		return
	}
	// Set node value
	data := []byte(i.value)
	err = i.in.ensurePath(lockpath, data)
	if err != nil {
		failLock <- err
		lock.Unlock()
		return
	}
	i.zkLock = lock

	// Signal that lock is held
	close(didLock)

	// Handle an early abort
	release := <-releaseCh
	if release {
		lock.Unlock()
	}
}

func (i *ZookeeperHALock) monitorLock(lockeventCh <-chan zk.Event, leaderCh chan struct{}) {
	for {
		select {
		case event := <-lockeventCh:
			// Lost connection?
			switch event.State {
			case zk.StateConnected:
			case zk.StateHasSession:
			default:
				close(leaderCh)
				return
			}

			// Lost lock?
			switch event.Type {
			case zk.EventNodeChildrenChanged:
			case zk.EventSession:
			default:
				close(leaderCh)
				return
			}
		}
	}
}

func (i *ZookeeperHALock) Unlock() error {
	i.localLock.Lock()
	defer i.localLock.Unlock()
	if !i.held {
		return nil
	}

	i.held = false
	i.zkLock.Unlock()
	return nil
}

func (i *ZookeeperHALock) Value() (bool, string, error) {
	lockpath := i.in.nodePath(i.key)
	value, _, err := i.in.client.Get(lockpath)
	return (value != nil), string(value), err
}
