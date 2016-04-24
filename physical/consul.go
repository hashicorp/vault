package physical

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"crypto/tls"
	"crypto/x509"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/lib"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
)

const (
	// checkJitterFactor specifies the jitter factor used to stagger checks
	checkJitterFactor = 16

	// checkMinBuffer specifies provides a guarantee that a check will not
	// be executed too close to the TTL check timeout
	checkMinBuffer = 100 * time.Millisecond

	// defaultCheckTimeout changes the timeout of TTL checks
	defaultCheckTimeout = 5 * time.Second

	// defaultCheckInterval specifies the default interval used to send
	// checks
	defaultCheckInterval = 4 * time.Second

	// defaultServiceName is the default Consul service name used when
	// advertising a Vault instance.
	defaultServiceName = "vault"
)

// ConsulBackend is a physical backend that stores data at specific
// prefix within Consul. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type ConsulBackend struct {
	path             string
	client           *api.Client
	kv               *api.KV
	permitPool       *PermitPool
	serviceLock      sync.RWMutex
	service          *api.AgentServiceRegistration
	sealedCheck      *api.AgentCheckRegistration
	advertiseAddr    string
	consulClientConf *api.Config
	serviceName      string
	running          bool
	active           bool
	sealed           bool
	checkTimeout     time.Duration
	checkTimer       *time.Timer
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

	// Get the service name to advertise in Consul
	service, ok := conf["service"]
	if !ok {
		service = defaultServiceName
	}

	checkTimeout := defaultCheckTimeout
	checkTimeoutStr, ok := conf["check_timeout"]
	if ok {
		d, err := time.ParseDuration(checkTimeoutStr)
		if err != nil {
			return nil, err
		}

		min, _ := lib.DurationMinusBufferDomain(d, checkMinBuffer, checkJitterFactor)
		if d < min {
			return nil, fmt.Errorf("Consul check_timeout must be greater than %v", min)
		}

		checkTimeout = d
	}

	// Configure the client
	consulConf := api.DefaultConfig()

	if addr, ok := conf["address"]; ok {
		consulConf.Address = addr
	}
	if scheme, ok := conf["scheme"]; ok {
		consulConf.Scheme = scheme
	}
	if token, ok := conf["token"]; ok {
		consulConf.Token = token
	}

	if consulConf.Scheme == "https" {
		tlsClientConfig, err := setupTLSConfig(conf)
		if err != nil {
			return nil, err
		}

		transport := cleanhttp.DefaultPooledTransport()
		transport.MaxIdleConnsPerHost = 4
		transport.TLSClientConfig = tlsClientConfig
		consulConf.HttpClient.Transport = transport
	}

	client, err := api.NewClient(consulConf)
	if err != nil {
		return nil, errwrap.Wrapf("client setup failed: {{err}}", err)
	}

	maxParStr, ok := conf["max_parallel"]
	var maxParInt int
	if ok {
		maxParInt, err = strconv.Atoi(maxParStr)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing max_parallel parameter: {{err}}", err)
		}
	}

	// Setup the backend
	c := &ConsulBackend{
		path:             path,
		client:           client,
		kv:               client.KV(),
		permitPool:       NewPermitPool(maxParInt),
		consulClientConf: consulConf,
		serviceName:      service,
		checkTimeout:     checkTimeout,
		checkTimer:       time.NewTimer(checkTimeout),
	}
	return c, nil
}

// UpdateAdvertiseAddr provides a pre-initialization hook for updating
// Consul's advertise address.
func (c *ConsulBackend) UpdateAdvertiseAddr(addr string) error {
	if c.running {
		return fmt.Errorf("service registration unable to update advertise address, backend already running")
	}

	url, err := url.Parse(addr)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf(`updating advertise address failed to parse URL "%v": {{err}}`, addr), err)
	}

	_, portStr, err := net.SplitHostPort(url.Host)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf(`updating advertise address failed to find a host:port in advertise address "%v": {{err}}`, url.Host), err)
	}
	_, err = strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf(`updating advertise address failed to parse port "%v": {{err}}`, portStr), err)
	}

	c.advertiseAddr = addr
	return nil
}

// serviceTags returns all of the relevant tags for Consul.  Assumes
// c.serviceLock held for writing.
func serviceTags(active bool) []string {
	activeTag := "standby"
	if active {
		activeTag = "active"
	}
	return []string{activeTag}
}

func (c *ConsulBackend) AdvertiseActive(active bool) error {
	c.serviceLock.Lock()
	defer c.serviceLock.Unlock()

	// Vault is still bootstrapping
	if c.service == nil {
		return nil
	}

	c.service.Tags = serviceTags(active)
	agent := c.client.Agent()
	if err := agent.ServiceRegister(c.service); err != nil {
		return errwrap.Wrapf("service registration failed: {{err}}", err)
	}

	// Save a cached copy of the active state: no way to query Core
	c.active = active

	return nil
}

func (c *ConsulBackend) AdvertiseSealed(sealed bool) error {
	c.serviceLock.Lock()
	defer c.serviceLock.Unlock()
	c.sealed = sealed

	// Vault is still bootstrapping
	if c.service == nil {
		return nil
	}

	// Push a TTL check immediately to update the state
	c.runCheck()

	return nil
}

func (c *ConsulBackend) RunServiceDiscovery(shutdownCh ShutdownChannel) (err error) {
	c.serviceLock.Lock()
	defer c.serviceLock.Unlock()

	if c.running {
		return fmt.Errorf("service registration routine already running")
	}

	url, err := url.Parse(c.advertiseAddr)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf(`service registration failed to parse URL "%v": {{err}}`, c.advertiseAddr), err)
	}

	host, portStr, err := net.SplitHostPort(url.Host)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf(`service registration failed to find a host:port in advertise address "%v": {{err}}`, url.Host), err)
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf(`service registration failed to parse port "%v": {{err}}`, portStr), err)
	}

	serviceID, err := c.serviceID()
	if err != nil {
		return err
	}

	c.service = &api.AgentServiceRegistration{
		ID:                serviceID,
		Name:              c.serviceName,
		Tags:              serviceTags(c.active),
		Port:              int(port),
		Address:           host,
		EnableTagOverride: false,
	}

	checkStatus := "failing"
	if !c.sealed {
		checkStatus = "passing"
	}

	c.sealedCheck = &api.AgentCheckRegistration{
		ID:        c.checkID(),
		Name:      "Vault Sealed Status",
		Notes:     "Vault service is healthy when Vault is in an unsealed status and can become an active Vault server",
		ServiceID: serviceID,
		AgentServiceCheck: api.AgentServiceCheck{
			TTL:    c.checkTimeout.String(),
			Status: checkStatus,
		},
	}

	agent := c.client.Agent()
	if err := agent.ServiceRegister(c.service); err != nil {
		return errwrap.Wrapf("service registration failed: {{err}}", err)
	}

	if err := agent.CheckRegister(c.sealedCheck); err != nil {
		return errwrap.Wrapf("service registration check registration failed: {{err}}", err)
	}

	go c.checkRunner(shutdownCh)
	c.running = true

	// Deregister upon shutdown
	go func() {
	shutdown:
		for {
			select {
			case <-shutdownCh:
				// wtb logger: log.Printf("[DEBUG]: Shutting down consul backend")
				break shutdown
			}
		}

		if err := agent.ServiceDeregister(serviceID); err != nil {
			// wtb logger: log.Printf("[WARNING]: service deregistration failed: {{err}}", err)
		}
		c.running = false
	}()

	return nil
}

// checkRunner periodically runs TTL checks
func (c *ConsulBackend) checkRunner(shutdownCh ShutdownChannel) {
	defer c.checkTimer.Stop()

	for {
		select {
		case <-c.checkTimer.C:
			go func() {
				c.serviceLock.Lock()
				defer c.serviceLock.Unlock()
				c.runCheck()
			}()
		case <-shutdownCh:
			return
		}
	}
}

// runCheck immediately pushes a TTL check.  Assumes c.serviceLock is held
// exclusively.
func (c *ConsulBackend) runCheck() {
	// Reset timer before calling run check in order to not slide the
	// window of the next check.
	c.checkTimer.Reset(lib.DurationMinusBuffer(c.checkTimeout, checkMinBuffer, checkJitterFactor))

	// Run a TTL check
	agent := c.client.Agent()
	if !c.sealed {
		agent.UpdateTTL(c.checkID(), "Vault Unsealed", api.HealthPassing)
	} else {
		agent.UpdateTTL(c.checkID(), "Vault Sealed", api.HealthCritical)
	}
}

// checkID returns the ID used for a Consul Check.  Assume at least a read
// lock is held.
func (c *ConsulBackend) checkID() string {
	return "vault-sealed-check"
}

// serviceID returns the Vault ServiceID for use in Consul.  Assume at least
// a read lock is held.
func (c *ConsulBackend) serviceID() (string, error) {
	url, err := url.Parse(c.advertiseAddr)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf(`service registration failed to parse URL "%v": {{err}}`, c.advertiseAddr), err)
	}

	host, portStr, err := net.SplitHostPort(url.Host)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf(`service registration failed to find a host:port in advertise address "%v": {{err}}`, url.Host), err)
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf(`service registration failed to parse port "%v": {{err}}`, portStr), err)
	}

	return fmt.Sprintf("%s:%s:%d", c.serviceName, host, int(port)), nil
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
	_, okKey := conf["tls_key_file"]

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

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	_, err := c.kv.Put(pair, nil)
	return err
}

// Get is used to fetch an entry
func (c *ConsulBackend) Get(key string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"consul", "get"}, time.Now())

	c.permitPool.Acquire()
	defer c.permitPool.Release()

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

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	_, err := c.kv.Delete(c.path+key, nil)
	return err
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (c *ConsulBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"consul", "list"}, time.Now())
	scan := c.path + prefix

	// The TrimPrefix call below will not work correctly if we have "//" at the
	// end. This can happen in cases where you are e.g. listing the root of a
	// prefix in a logical backend via "/" instead of ""
	if strings.HasSuffix(scan, "//") {
		scan = scan[:len(scan)-1]
	}

	c.permitPool.Acquire()
	defer c.permitPool.Release()

	out, _, err := c.kv.Keys(scan, "/", nil)
	for idx, val := range out {
		out[idx] = strings.TrimPrefix(val, scan)
	}

	return out, err
}

// Lock is used for mutual exclusion based on the given key.
func (c *ConsulBackend) LockWith(key, value string) (Lock, error) {
	// Create the lock
	opts := &api.LockOptions{
		Key:            c.path + key,
		Value:          []byte(value),
		SessionName:    "Vault Lock",
		MonitorRetries: 5,
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
