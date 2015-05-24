package physical

import (
	"encoding/base64"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/coreos/go-etcd/etcd"
)

const (

	// Ideally, this prefix would match the "_" used in the file backend, but
	// that prefix has special meaining in etcd. Specifically, it excludes those
	// entries from directory listings.
	EtcdNodeFilePrefix = "."

	// The delimiter is the same as the `-C` flag of etcdctl.
	EtcdMachineDelimiter = ","
)

var (
	EtcdSyncClusterError = errors.New("client setup failed: unable to sync etcd cluster")
)

// errorIsMissingKey returns true if the given error is an etcd error with an
// error code corresponding to a missing key.
func errorIsMissingKey(err error) bool {
	etcdErr, ok := err.(*etcd.EtcdError)
	return ok && etcdErr.ErrorCode == 100
}

// EtcdBackend is a physical backend that stores data at specific
// prefix within Etcd. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type EtcdBackend struct {
	path   string
	client *etcd.Client
}

// newEtcdBackend constructs a etcd backend using a given machine address.
func newEtcdBackend(conf map[string]string) (Backend, error) {

	// Get the etcd path form the configuration.
	path, ok := conf["path"]
	if !ok {
		path = "/vault"
	}

	// Ensure path is prefixed.
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Set a default machines list and check for an overriding address value.
	machines := "http://128.0.0.1:4001"
	if address, ok := conf["address"]; ok {
		machines = address
	}

	// Create a new client from the supplied addres and attempt to sync with the
	// cluster.
	client := etcd.NewClient(strings.Split(machines, EtcdMachineDelimiter))
	if !client.SyncCluster() {
		return nil, EtcdSyncClusterError
	}

	// Setup the backend
	return &EtcdBackend{
		path:   path,
		client: client,
	}, nil
}

// Put is used to insert or update an entry.
func (c *EtcdBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"etcd", "put"}, time.Now())
	value := base64.StdEncoding.EncodeToString(entry.Value)
	_, err := c.client.Set(c.nodePath(entry.Key), value, 0)
	return err
}

// Get is used to fetch an entry.
func (c *EtcdBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"etcd", "get"}, time.Now())

	response, err := c.client.Get(c.nodePath(key), false, false)
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
	return &Entry{
		Key:   key,
		Value: value,
	}, nil
}

// Delete is used to permanently delete an entry.
func (c *EtcdBackend) Delete(key string) error {
	defer metrics.MeasureSince([]string{"etcd", "delete"}, time.Now())

	// Remove the key, non-recursively.
	_, err := c.client.Delete(c.nodePath(key), false)
	if err != nil && !errorIsMissingKey(err) {
		return err
	}

	return nil
}

// List is used to list all the keys under a given prefix, up to the next
// prefix.
func (c *EtcdBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"etcd", "list"}, time.Now())

	// Set a directory path from the given prefix.
	path := c.nodePathDir(prefix)

	// Get the directory, non-recursively, from etcd. If the directory is
	// missing, we just return an empty list of contents.
	response, err := c.client.Get(path, true, false)
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
func (b *EtcdBackend) nodePath(key string) string {
	return filepath.Join(b.path, filepath.Dir(key), EtcdNodeFilePrefix+filepath.Base(key))
}

// nodePathDir returns an etcd directory path based on the given key.
func (b *EtcdBackend) nodePathDir(key string) string {
	return filepath.Join(b.path, key) + "/"
}
