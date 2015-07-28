package physical

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"net/http"
	"io/ioutil"

	"crypto/tls"
	"crypto/x509"

	"github.com/armon/go-metrics"
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

	if consulConf.Scheme == "https" {
		tlsClientConfig, err := setupTLSConfig(conf)
		if err != nil {
			return nil, err
		}

		consulConf.HttpClient.Transport = &http.Transport{
			TLSClientConfig: tlsClientConfig,
		}
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

func setupTLSConfig(conf map[string]string) (*tls.Config, error) {
	serverName := strings.Split(conf["address"], ":")

	insecureSkipVerify := false
	if _, ok := conf["tls_skip_verify"]; ok {
		insecureSkipVerify = true
	}

	tlsClientConfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
		ServerName:         serverName[0],
	}

	_, okCert := conf["tls_cert_file"]
	_, okKey  := conf["tls_key_file"]

	if okCert && okKey {
		tlsCert, err := tls.LoadX509KeyPair(conf["tls_cert_file"], conf["tls_key_file"])
		if err != nil {
			return nil, fmt.Errorf("client tls setup failed: %v", err)
		}

		tlsClientConfig.Certificates = []tls.Certificate{tlsCert}
	}

	if tlsCaFile, ok := conf["tls_ca_file"]; ok {
		caPool := x509.NewCertPool()

		data, err := ioutil.ReadFile(tlsCaFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file: %v", err)
		}

		if !caPool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}

		tlsClientConfig.RootCAs = caPool
	}

	return tlsClientConfig, nil
}

// Put is used to insert or update an entry
func (c *ConsulBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"consul", "put"}, time.Now())
	pair := &api.KVPair{
		Key:   c.path + entry.Key,
		Value: entry.Value,
	}
	_, err := c.kv.Put(pair, nil)
	return err
}

// Get is used to fetch an entry
func (c *ConsulBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"consul", "get"}, time.Now())
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
	defer metrics.MeasureSince([]string{"consul", "delete"}, time.Now())
	_, err := c.kv.Delete(c.path+key, nil)
	return err
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (c *ConsulBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"consul", "list"}, time.Now())
	scan := c.path + prefix
	out, _, err := c.kv.Keys(scan, "/", nil)
	for idx, val := range out {
		out[idx] = strings.TrimPrefix(val, scan)
	}
	sort.Strings(out)
	return out, err
}

// Lock is used for mutual exclusion based on the given key.
func (c *ConsulBackend) LockWith(key, value string) (Lock, error) {
	// Create the lock
	opts := &api.LockOptions{
		Key:         c.path + key,
		Value:       []byte(value),
		SessionName: "Vault Lock",
	}
	lock, err := c.client.LockOpts(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create lock: %v", err)
	}
	cl := &ConsulLock{
		client: c.client,
		key:    c.path + key,
		lock:   lock,
	}
	return cl, nil
}

// DetectHostAddr is used to detect the host address by asking the Consul agent
func (c *ConsulBackend) DetectHostAddr() (string, error) {
	agent := c.client.Agent()
	self, err := agent.Self()
	if err != nil {
		return "", err
	}
	addr := self["Member"]["Addr"].(string)
	return addr, nil
}

// ConsulLock is used to provide the Lock interface backed by Consul
type ConsulLock struct {
	client *api.Client
	key    string
	lock   *api.Lock
}

func (c *ConsulLock) Lock(stopCh <-chan struct{}) (<-chan struct{}, error) {
	return c.lock.Lock(stopCh)
}

func (c *ConsulLock) Unlock() error {
	return c.lock.Unlock()
}

func (c *ConsulLock) Value() (bool, string, error) {
	kv := c.client.KV()
	pair, _, err := kv.Get(c.key, nil)
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
