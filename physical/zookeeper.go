package physical

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/samuel/go-zookeeper/zk"
)

// ZookeeperBackend is a physical backend that stores data at specific
// prefix within Zookeeper. It is used in production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type ZookeeperBackend struct {
	path   string
	client *zk.Conn
}

// newZookeeperBackend constructs a Zookeeper backend using the given API client
// and the prefix in the KV store.
func newZookeeperBackend(conf map[string]string) (Backend, error) {
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

	// Attempt to create the ZK client
	client, _, err := zk.Connect(strings.Split(machines, ","), time.Second)
	if err != nil {
		return nil, fmt.Errorf("client setup failed: %v", err)
	}

	// Setup the backend
	c := &ZookeeperBackend{
		path:   path,
		client: client,
	}
	return c, nil
}

// ensurePath is used to create each node in the path hierarchy.
// We avoid calling this optimistically, and invoke it when we get
// an error during an operation
func (c *ZookeeperBackend) ensurePath(path string, value []byte) error {
	nodes := strings.Split(path, "/")
	acl := zk.WorldACL(zk.PermAll)
	fullPath := ""
	for index, node := range nodes {
		if strings.TrimSpace(node) != "" {
			fullPath += "/" + node
			isLastNode := index+1 == len(nodes)

			// set parent nodes to nil, leaf to value
			// this block reduces round trips by being smart on the leaf create/set
			if exists, _, _ := c.client.Exists(fullPath); !isLastNode && !exists {
				if _, err := c.client.Create(fullPath, nil, int32(0), acl); err != nil {
					return err
				}
			} else if isLastNode && !exists {
				if _, err := c.client.Create(fullPath, value, int32(0), acl); err != nil {
					return err
				}
			} else if isLastNode && exists {
				if _, err := c.client.Set(fullPath, value, int32(0)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Put is used to insert or update an entry
func (c *ZookeeperBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"zookeeper", "put"}, time.Now())

	// Attempt to set the full path
	fullPath := c.path + entry.Key
	_, err := c.client.Set(fullPath, entry.Value, 0)

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
	fullPath := c.path + key
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
	fullPath := c.path + key
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
		children = append(children, key)

		// Check if this entry has any child entries,
		// and append the slash which is what Vault depends on
		// for iteration
		nodeChildren, _, _ := c.client.Children(fullPath + "/" + key)
		if nodeChildren != nil && len(nodeChildren) > 0 {
			children = append(children, key+"/")
		}
	}
	sort.Strings(children)
	return children, nil
}
