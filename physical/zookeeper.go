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

	// Configure the client
	var machines string
	machines, ok = conf["address"]
	if !ok {
		// Default to the localhost instance
		machines = "localhost:2181"
	}

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

// zookeeper requires nodes to be there before set and get
func (c *ZookeeperBackend) ensurePath(path string, value []byte) {

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
				c.client.Create(fullPath, nil, int32(0), acl)
			} else if isLastNode && !exists {
				c.client.Create(fullPath, value, int32(0), acl)
			} else if isLastNode && exists {
				c.client.Set(fullPath, value, int32(0))
			}

		}
	}
}

// Put is used to insert or update an entry
func (c *ZookeeperBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"zookeeper", "put"}, time.Now())

	fullPath := c.path + entry.Key

	_, err := c.client.Set(fullPath, entry.Value, 0)

	if err == zk.ErrNoNode {
		c.ensurePath(fullPath, entry.Value)
		return nil
	} else {
		return err
	}
}

// Get is used to fetch an entry
func (c *ZookeeperBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"zookeeper", "get"}, time.Now())

	fullPath := c.path + key

	value, _, err := c.client.Get(fullPath)

	// dont error if the node doesnt exist
	if err != nil && err != zk.ErrNoNode {
		return nil, err
	}
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

	fullPath := c.path + key

	err := c.client.Delete(fullPath, -1)

	if err == zk.ErrNoNode {
		return nil
	} else {
		return err
	}
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (c *ZookeeperBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"zookeeper", "list"}, time.Now())

	fullPath := strings.TrimSuffix(c.path+prefix, "/")

	result, _, err := c.client.Children(fullPath)

	if err == zk.ErrNoNode {
		return []string{}, nil
	}

	children := []string{}

	for _, key := range result {
		children = append(children, key)

		nodeChildren, _, _ := c.client.Children(fullPath + "/" + key)
		if nodeChildren != nil && len(nodeChildren) > 0 {
			children = append(children, key+"/")
		}
	}

	sort.Strings(children)

	return children, nil
}
