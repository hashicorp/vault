package zookeeper

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/parseutil"
	"github.com/hashicorp/vault/physical"

	metrics "github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	// ZKNodeFilePrefix is prefixed to any "files" in ZooKeeper,
	// so that they do not collide with directory entries. Otherwise,
	// we cannot delete a file if the path is a full-prefix of another
	// key.
	ZKNodeFilePrefix = "_"
)

// Verify ZooKeeperBackend satisfies the correct interfaces
var _ physical.Backend = (*ZooKeeperBackend)(nil)
var _ physical.HABackend = (*ZooKeeperBackend)(nil)
var _ physical.Lock = (*ZooKeeperHALock)(nil)

// ZooKeeperBackend is a physical backend that stores data at specific
// prefix within ZooKeeper. It is used in production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type ZooKeeperBackend struct {
	path   string
	client *zk.Conn
	acl    []zk.ACL
	logger log.Logger
}

// NewZooKeeperBackend constructs a ZooKeeper backend using the given API client
// and the prefix in the KV store.
func NewZooKeeperBackend(conf map[string]string, logger log.Logger) (physical.Backend, error) {
	// Get the path in ZooKeeper
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

	acl := []zk.ACL{
		{
			Perms:  zk.PermAll,
			Scheme: schema,
			ID:     owner,
		},
	}

	// Authentication info
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
	client, _, err := createClient(conf, machines, time.Second)
	if err != nil {
		return nil, errwrap.Wrapf("client setup failed: {{err}}", err)
	}

	// ZK AddAuth API if the user asked for it
	if useAddAuth {
		err = client.AddAuth(schema, []byte(owner))
		if err != nil {
			return nil, errwrap.Wrapf("ZooKeeper rejected authentication information provided at auth_info: {{err}}", err)
		}
	}

	// Setup the backend
	c := &ZooKeeperBackend{
		path:   path,
		client: client,
		acl:    acl,
		logger: logger,
	}
	return c, nil
}

func caseInsenstiveContains(superset, val string) bool {
	return strings.Contains(strings.ToUpper(superset), strings.ToUpper(val))
}

// Returns a client for ZK connection. Config value 'tls_enabled' determines if TLS is enabled or not.
func createClient(conf map[string]string, machines string, timeout time.Duration) (*zk.Conn, <-chan zk.Event, error) {
	// 'tls_enabled' defaults to false
	isTlsEnabled := false
	isTlsEnabledStr, ok := conf["tls_enabled"]

	if ok && isTlsEnabledStr != "" {
		parsedBoolval, err := parseutil.ParseBool(isTlsEnabledStr)
		if err != nil {
			return nil, nil, errwrap.Wrapf("failed parsing tls_enabled parameter: {{err}}", err)
		}
		isTlsEnabled = parsedBoolval
	}

	if isTlsEnabled {
		// Create a custom Dialer with cert configuration for TLS handshake.
		tlsDialer := customTLSDial(conf, machines)
		options := zk.WithDialer(tlsDialer)
		return zk.Connect(strings.Split(machines, ","), timeout, options)
	} else {
		return zk.Connect(strings.Split(machines, ","), timeout)
	}
}

// Vault config file properties:
// 1. tls_skip_verify: skip host name verification.
// 2. tls_min_version: minimum supported/acceptable tls version
// 3. tls_cert_file: Cert file Absolute path
// 4. tls_key_file: Key file Absolute path
// 5. tls_ca_file: ca file absolute path
// 6. tls_verify_ip: If set to true, server's IP is verified in certificate if tls_skip_verify is false.
func customTLSDial(conf map[string]string, machines string) zk.Dialer {
	return func(network, addr string, timeout time.Duration) (net.Conn, error) {
		// Sets the serverName. *Note* the addr field comes in as an IP address
		serverName, _, sParseErr := net.SplitHostPort(addr)
		if sParseErr != nil {
			// If the address is only missing port, assign the full address anyway
			if strings.Contains(sParseErr.Error(), "missing port") {
				serverName = addr
			} else {
				return nil, errwrap.Wrapf("failed parsing the server address for 'serverName' setting {{err}}", sParseErr)
			}
		}

		insecureSkipVerify := false
		tlsSkipVerify, ok := conf["tls_skip_verify"]

		if ok && tlsSkipVerify != "" {
			b, err := parseutil.ParseBool(tlsSkipVerify)
			if err != nil {
				return nil, errwrap.Wrapf("failed parsing tls_skip_verify parameter: {{err}}", err)
			}
			insecureSkipVerify = b
		}

		if !insecureSkipVerify {
			// If tls_verify_ip is set to false, Server's DNS name is verified in the CN/SAN of the certificate.
			// if tls_verify_ip is true, Server's IP is verified in the CN/SAN of the certificate.
			// These checks happen only when tls_skip_verify is set to false.
			// This value defaults to false
			ipSanCheck := false
			configVal, lookupOk := conf["tls_verify_ip"]

			if lookupOk && configVal != "" {
				parsedIpSanCheck, ipSanErr := parseutil.ParseBool(configVal)
				if ipSanErr != nil {
					return nil, errwrap.Wrapf("failed parsing tls_verify_ip parameter: {{err}}", ipSanErr)
				}
				ipSanCheck = parsedIpSanCheck
			}
			// The addr/serverName parameter to this method comes in as an IP address.
			// Here we lookup the DNS name and assign it to serverName if ipSanCheck is set to false
			if !ipSanCheck {
				lookupAddressMany, lookupErr := net.LookupAddr(serverName)
				if lookupErr == nil {
					for _, lookupAddress := range lookupAddressMany {
						// strip the trailing '.' from lookupAddr
						if lookupAddress[len(lookupAddress)-1] == '.' {
							lookupAddress = lookupAddress[:len(lookupAddress)-1]
						}
						// Allow serverName to be replaced only if the lookupname is part of the
						// supplied machine names
						// If there is no match, the serverName will continue to be an IP value.
						if caseInsenstiveContains(machines, lookupAddress) {
							serverName = lookupAddress
							break
						}
					}
				}
			}

		}

		tlsMinVersionStr, ok := conf["tls_min_version"]
		if !ok {
			// Set the default value
			tlsMinVersionStr = "tls12"
		}

		tlsMinVersion, ok := tlsutil.TLSLookup[tlsMinVersionStr]
		if !ok {
			return nil, fmt.Errorf("invalid 'tls_min_version'")
		}

		tlsClientConfig := &tls.Config{
			MinVersion:         tlsMinVersion,
			InsecureSkipVerify: insecureSkipVerify,
			ServerName:         serverName,
		}

		_, okCert := conf["tls_cert_file"]
		_, okKey := conf["tls_key_file"]

		if okCert && okKey {
			tlsCert, err := tls.LoadX509KeyPair(conf["tls_cert_file"], conf["tls_key_file"])
			if err != nil {
				return nil, errwrap.Wrapf("client tls setup failed for ZK: {{err}}", err)
			}

			tlsClientConfig.Certificates = []tls.Certificate{tlsCert}
		}

		if tlsCaFile, ok := conf["tls_ca_file"]; ok {
			caPool := x509.NewCertPool()

			data, err := ioutil.ReadFile(tlsCaFile)
			if err != nil {
				return nil, errwrap.Wrapf("failed to read ZK CA file: {{err}}", err)
			}

			if !caPool.AppendCertsFromPEM(data) {
				return nil, fmt.Errorf("failed to parse ZK CA certificate")
			}
			tlsClientConfig.RootCAs = caPool
		}

		if network != "tcp" {
			return nil, fmt.Errorf("unsupported network %q", network)
		}

		tcpConn, err := net.DialTimeout("tcp", addr, timeout)
		if err != nil {
			return nil, err
		}
		conn := tls.Client(tcpConn, tlsClientConfig)
		if err := conn.Handshake(); err != nil {
			return nil, fmt.Errorf("Handshake failed with Zookeeper : %v", err)
		}
		return conn, nil
	}
}

// ensurePath is used to create each node in the path hierarchy.
// We avoid calling this optimistically, and invoke it when we get
// an error during an operation
func (c *ZooKeeperBackend) ensurePath(path string, value []byte) error {
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

// cleanupLogicalPath is used to remove all empty nodes, beginning with deepest one,
// aborting on first non-empty one, up to top-level node.
func (c *ZooKeeperBackend) cleanupLogicalPath(path string) error {
	nodes := strings.Split(path, "/")
	for i := len(nodes) - 1; i > 0; i-- {
		fullPath := c.path + strings.Join(nodes[:i], "/")

		_, stat, err := c.client.Exists(fullPath)
		if err != nil {
			return errwrap.Wrapf("failed to acquire node data: {{err}}", err)
		}

		if stat.DataLength > 0 && stat.NumChildren > 0 {
			panic(fmt.Sprintf("node %q is both of data and leaf type", fullPath))
		} else if stat.DataLength > 0 {
			panic(fmt.Sprintf("node %q is a data node, this is either a bug or backend data is corrupted", fullPath))
		} else if stat.NumChildren > 0 {
			return nil
		} else {
			// Empty node, lets clean it up!
			if err := c.client.Delete(fullPath, -1); err != nil && err != zk.ErrNoNode {
				return errwrap.Wrapf(fmt.Sprintf("removal of node %q failed: {{err}}", fullPath), err)
			}
		}
	}
	return nil
}

// nodePath returns an zk path based on the given key.
func (c *ZooKeeperBackend) nodePath(key string) string {
	return filepath.Join(c.path, filepath.Dir(key), ZKNodeFilePrefix+filepath.Base(key))
}

// Put is used to insert or update an entry
func (c *ZooKeeperBackend) Put(ctx context.Context, entry *physical.Entry) error {
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
func (c *ZooKeeperBackend) Get(ctx context.Context, key string) (*physical.Entry, error) {
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
	ent := &physical.Entry{
		Key:   key,
		Value: value,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (c *ZooKeeperBackend) Delete(ctx context.Context, key string) error {
	defer metrics.MeasureSince([]string{"zookeeper", "delete"}, time.Now())

	if key == "" {
		return nil
	}

	// Delete the full path
	fullPath := c.nodePath(key)
	err := c.client.Delete(fullPath, -1)

	// Mask if the node does not exist
	if err != nil && err != zk.ErrNoNode {
		return errwrap.Wrapf(fmt.Sprintf("failed to remove %q: {{err}}", fullPath), err)
	}

	err = c.cleanupLogicalPath(key)

	return err
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (c *ZooKeeperBackend) List(ctx context.Context, prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"zookeeper", "list"}, time.Now())

	// Query the children at the full path
	fullPath := strings.TrimSuffix(c.path+prefix, "/")
	result, _, err := c.client.Children(fullPath)

	// If the path nodes are missing, no children!
	if err == zk.ErrNoNode {
		return []string{}, nil
	} else if err != nil {
		return []string{}, err
	}

	children := []string{}
	for _, key := range result {
		childPath := fullPath + "/" + key
		_, stat, err := c.client.Exists(childPath)
		if err != nil {
			// Node is ought to exists, so it must be something different
			return []string{}, err
		}

		// Check if this entry is a leaf of a node,
		// and append the slash which is what Vault depends on
		// for iteration
		if stat.DataLength > 0 && stat.NumChildren > 0 {
			if childPath == c.nodePath("core/lock") {
				// go-zookeeper Lock() breaks Vault semantics and creates a directory
				// under the lock file; just treat it like the file Vault expects
				children = append(children, key[1:])
			} else {
				panic(fmt.Sprintf("node %q is both of data and leaf type", childPath))
			}
		} else if stat.DataLength == 0 {
			// No, we cannot differentiate here on number of children as node
			// can have all it leafs removed, and it still is a node.
			children = append(children, key+"/")
		} else {
			children = append(children, key[1:])
		}
	}
	sort.Strings(children)
	return children, nil
}

// LockWith is used for mutual exclusion based on the given key.
func (c *ZooKeeperBackend) LockWith(key, value string) (physical.Lock, error) {
	l := &ZooKeeperHALock{
		in:     c,
		key:    key,
		value:  value,
		logger: c.logger,
	}
	return l, nil
}

// HAEnabled indicates whether the HA functionality should be exposed.
// Currently always returns true.
func (c *ZooKeeperBackend) HAEnabled() bool {
	return true
}

// ZooKeeperHALock is a ZooKeeper Lock implementation for the HABackend
type ZooKeeperHALock struct {
	in     *ZooKeeperBackend
	key    string
	value  string
	logger log.Logger

	held      bool
	localLock sync.Mutex
	leaderCh  chan struct{}
	stopCh    <-chan struct{}
	zkLock    *zk.Lock
}

func (i *ZooKeeperHALock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
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
		return nil, errwrap.Wrapf("unable to watch HA lock: {{err}}", err)
	}
	if i.value != string(currentVal) {
		return nil, fmt.Errorf("lost HA lock immediately before watch")
	}
	go i.monitorLock(lockeventCh, i.leaderCh)

	i.stopCh = stopCh

	return i.leaderCh, nil
}

func (i *ZooKeeperHALock) attemptLock(lockpath string, didLock chan struct{}, failLock chan error, releaseCh chan bool) {
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

func (i *ZooKeeperHALock) monitorLock(lockeventCh <-chan zk.Event, leaderCh chan struct{}) {
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

func (i *ZooKeeperHALock) unlockInternal() error {
	i.localLock.Lock()
	defer i.localLock.Unlock()
	if !i.held {
		return nil
	}

	err := i.zkLock.Unlock()

	if err == nil {
		i.held = false
		return nil
	}

	return err
}

func (i *ZooKeeperHALock) Unlock() error {
	var err error

	if err = i.unlockInternal(); err != nil {
		i.logger.Error("failed to release distributed lock", "error", err)

		go func(i *ZooKeeperHALock) {
			attempts := 0
			i.logger.Info("launching automated distributed lock release")

			for {
				if err := i.unlockInternal(); err == nil {
					i.logger.Info("distributed lock released")
					return
				}

				select {
				case <-time.After(time.Second):
					attempts := attempts + 1
					if attempts >= 10 {
						i.logger.Error("release lock max attempts reached. Lock may not be released", "error", err)
						return
					}
					continue
				case <-i.stopCh:
					return
				}
			}
		}(i)
	}

	return err
}

func (i *ZooKeeperHALock) Value() (bool, string, error) {
	lockpath := i.in.nodePath(i.key)
	value, _, err := i.in.client.Get(lockpath)
	return (value != nil), string(value), err
}
