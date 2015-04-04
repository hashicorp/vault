package physical

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/consul/api"
)

// ConsulBackend is a physical backend that stores data at specific
// prefix within Consul. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type ConsulBackend struct {
	path   string
	client *api.Client
	kv     *api.KV
}

// newConsulBackend constructs a Consul backend using the given API client
// and the prefix in the KV store.
func newConsulBackend(conf map[string]string) (Backend, error) {
	// Get the path in Consul
	path, ok := conf["path"]
	if !ok {
		path = "vault/"
	}

	// Ensure path is suffixed but not prefixed
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}

	// Configure the client
	consulConf := api.DefaultConfig()
	if addr, ok := conf["address"]; ok {
		consulConf.Address = addr
	}
	if scheme, ok := conf["scheme"]; ok {
		consulConf.Scheme = scheme
	}
	if dc, ok := conf["datacenter"]; ok {
		consulConf.Datacenter = dc
	}
	if token, ok := conf["token"]; ok {
		consulConf.Token = token
	}
	client, err := api.NewClient(consulConf)
	if err != nil {
		return nil, fmt.Errorf("client setup failed: %v", err)
	}

	// Setup the backend
	c := &ConsulBackend{
		path:   path,
		client: client,
		kv:     client.KV(),
	}
	return c, nil
}

// Put is used to insert or update an entry
func (c *ConsulBackend) Put(entry *Entry) error {
	pair := &api.KVPair{
		Key:   c.path + entry.Key,
		Value: entry.Value,
	}
	_, err := c.kv.Put(pair, nil)
	return err
}

// Get is used to fetch an entry
func (c *ConsulBackend) Get(key string) (*Entry, error) {
	pair, _, err := c.kv.Get(c.path+key, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, nil
	}
	ent := &Entry{
		Key:   key,
		Value: pair.Value,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (c *ConsulBackend) Delete(key string) error {
	_, err := c.kv.Delete(c.path+key, nil)
	return err
}

// List is used ot list all the keys under a given
// prefix, up to the next prefix.
func (c *ConsulBackend) List(prefix string) ([]string, error) {
	scan := c.path + prefix
	out, _, err := c.kv.Keys(scan, "/", nil)
	for idx, val := range out {
		out[idx] = strings.TrimPrefix(val, scan)
	}
	sort.Strings(out)
	return out, err
}
