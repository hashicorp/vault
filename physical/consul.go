package physical

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"crypto/tls"
	"crypto/x509"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/lib"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/helper/tlsutil"
)

const (
	// checkJitterFactor specifies the jitter factor used to stagger checks
	checkJitterFactor = 16

	// checkMinBuffer specifies provides a guarantee that a check will not
	// be executed too close to the TTL check timeout
	checkMinBuffer = 100 * time.Millisecond

	// consulRetryInterval specifies the retry duration to use when an
	// API call to the Consul agent fails.
	consulRetryInterval = 1 * time.Second

	// defaultCheckTimeout changes the timeout of TTL checks
	defaultCheckTimeout = 5 * time.Second

	// defaultServiceName is the default Consul service name used when
	// advertising a Vault instance.
	defaultServiceName = "vault"

	// reconcileTimeout is how often Vault should query Consul to detect
	// and fix any state drift.
	reconcileTimeout = 60 * time.Second
)

type notifyEvent struct{}

// ConsulBackend is a physical backend that stores data at specific
// prefix within Consul. It is used for most production situations as
// it allows Vault to run on multiple machines in a highly-available manner.
type ConsulBackend struct {
	path                string
	logger              *log.Logger
	client              *api.Client
	kv                  *api.KV
	permitPool          *PermitPool
	serviceLock         sync.RWMutex
	advertiseHost       string
	advertisePort       int64
	serviceName         string
	disableRegistration bool
	checkTimeout        time.Duration

	notifyActiveCh chan notifyEvent
	notifySealedCh chan notifyEvent
}

// newConsulBackend constructs a Consul backend using the given API client
// and the prefix in the KV store.
func newConsulBackend(conf map[string]string, logger *log.Logger) (Backend, error) {
	// Get the path in Consul
	path, ok := conf["path"]
	if !ok {
		path = "vault/"
	}
	logger.Printf("[DEBUG]: consul: config path set to %v", path)

	// Ensure path is suffixed but not prefixed
	if !strings.HasSuffix(path, "/") {
		logger.Printf("[WARN]: consul: appending trailing forward slash to path")
		path += "/"
	}
	if strings.HasPrefix(path, "/") {
		logger.Printf("[WARN]: consul: trimming path of its forward slash")
		path = strings.TrimPrefix(path, "/")
	}

	// Allow admins to disable consul integration
	disableReg, ok := conf["disable_registration"]
	var disableRegistration bool
	if ok && disableReg != "" {
		b, err := strconv.ParseBool(disableReg)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing disable_registration parameter: {{err}}", err)
		}
		disableRegistration = b
	}
	logger.Printf("[DEBUG]: consul: config disable_registration set to %v", disableRegistration)

	// Get the service name to advertise in Consul
	service, ok := conf["service"]
	if !ok {
		service = defaultServiceName
	}
	logger.Printf("[DEBUG]: consul: config service set to %s", service)

	checkTimeout := defaultCheckTimeout
	checkTimeoutStr, ok := conf["check_timeout"]
	if ok {
		d, err := time.ParseDuration(checkTimeoutStr)
		if err != nil {
			return nil, err
		}

		min, _ := lib.DurationMinusBufferDomain(d, checkMinBuffer, checkJitterFactor)
		if min < checkMinBuffer {
			return nil, fmt.Errorf("Consul check_timeout must be greater than %v", min)
		}

		checkTimeout = d
		logger.Printf("[DEBUG]: consul: config check_timeout set to %v", d)
	}

	// Configure the client
	consulConf := api.DefaultConfig()

	if addr, ok := conf["address"]; ok {
		consulConf.Address = addr
		logger.Printf("[DEBUG]: consul: config address set to %d", addr)
	}
	if scheme, ok := conf["scheme"]; ok {
		consulConf.Scheme = scheme
		logger.Printf("[DEBUG]: consul: config scheme set to %d", scheme)
	}
	if token, ok := conf["token"]; ok {
		consulConf.Token = token
		logger.Printf("[DEBUG]: consul: config token set")
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
		logger.Printf("[DEBUG]: consul: configured TLS")
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
		logger.Printf("[DEBUG]: consul: max_parallel set to %d", maxParInt)
	}

	// Setup the backend
	c := &ConsulBackend{
		path:                path,
		logger:              logger,
		client:              client,
		kv:                  client.KV(),
		permitPool:          NewPermitPool(maxParInt),
		serviceName:         service,
		checkTimeout:        checkTimeout,
		disableRegistration: disableRegistration,
	}
	return c, nil
}

func setupTLSConfig(conf map[string]string) (*tls.Config, error) {
	serverName := strings.Split(conf["address"], ":")

	insecureSkipVerify := false
	if _, ok := conf["tls_skip_verify"]; ok {
		insecureSkipVerify = true
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

// HAEnabled indicates whether the HA functionality should be exposed.
// Currently always returns true.
func (c *ConsulBackend) HAEnabled() bool {
	return true
}

// DetectHostAddr is used to detect the host address by asking the Consul agent
func (c *ConsulBackend) DetectHostAddr() (string, error) {
	agent := c.client.Agent()
	self, err := agent.Self()
	if err != nil {
		return "", err
	}
	addr, ok := self["Member"]["Addr"].(string)
	if !ok {
		return "", fmt.Errorf("Unable to convert an address to string")
	}
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

func (c *ConsulBackend) NotifyActiveStateChange() error {
	select {
	case c.notifyActiveCh <- notifyEvent{}:
	default:
		// NOTE: If this occurs Vault's active status could be out of
		// sync with Consul until reconcileTimer expires.
		c.logger.Printf("[WARN]: consul: Concurrent state change notify dropped")
	}

	return nil
}

func (c *ConsulBackend) NotifySealedStateChange() error {
	select {
	case c.notifySealedCh <- notifyEvent{}:
	default:
		// NOTE: If this occurs Vault's sealed status could be out of
		// sync with Consul until checkTimer expires.
		c.logger.Printf("[WARN]: consul: Concurrent sealed state change notify dropped")
	}

	return nil
}

func (c *ConsulBackend) checkDuration() time.Duration {
	return lib.DurationMinusBuffer(c.checkTimeout, checkMinBuffer, checkJitterFactor)
}

func (c *ConsulBackend) RunServiceDiscovery(shutdownCh ShutdownChannel, advertiseAddr string, activeFunc activeFunction, sealedFunc sealedFunction) (err error) {
	if err := c.setAdvertiseAddr(advertiseAddr); err != nil {
		return err
	}

	go c.runEventDemuxer(shutdownCh, advertiseAddr, activeFunc, sealedFunc)

	return nil
}

func (c *ConsulBackend) runEventDemuxer(shutdownCh ShutdownChannel, advertiseAddr string, activeFunc activeFunction, sealedFunc sealedFunction) {
	// Fire the reconcileTimer immediately upon starting the event demuxer
	reconcileTimer := time.NewTimer(0)
	defer reconcileTimer.Stop()

	// Schedule the first check.  Consul TTL checks are passing by
	// default, checkTimer does not need to be run immediately.
	checkTimer := time.NewTimer(c.checkDuration())
	defer checkTimer.Stop()

	// Use a reactor pattern to handle and dispatch events to singleton
	// goroutine handlers for execution.  It is not acceptable to drop
	// inbound events from Notify*().
	//
	// goroutines are dispatched if the demuxer can acquire a lock (via
	// an atomic CAS incr) on the handler.  Handlers are responsible for
	// deregistering themselves (atomic CAS decr).  Handlers and the
	// demuxer share a lock to synchronize information at the beginning
	// and end of a handler's life (or after a handler wakes up from
	// sleeping during a back-off/retry).
	var shutdown bool
	var checkLock int64
	var registeredServiceID string
	var serviceRegLock int64
shutdown:
	for {
		select {
		case <-c.notifyActiveCh:
			// Run reconcile immediately upon active state change notification
			reconcileTimer.Reset(0)
		case <-c.notifySealedCh:
			// Run check timer immediately upon a seal state change notification
			checkTimer.Reset(0)
		case <-reconcileTimer.C:
			// Unconditionally rearm the reconcileTimer
			reconcileTimer.Reset(reconcileTimeout - lib.RandomStagger(reconcileTimeout/checkJitterFactor))

			// Abort if service discovery is disabled or a
			// reconcile handler is already active
			if !c.disableRegistration && atomic.CompareAndSwapInt64(&serviceRegLock, 0, 1) {
				// Enter handler with serviceRegLock held
				go func() {
					defer atomic.CompareAndSwapInt64(&serviceRegLock, 1, 0)
					for !shutdown {
						serviceID, err := c.reconcileConsul(registeredServiceID, activeFunc, sealedFunc)
						if err != nil {
							c.logger.Printf("[WARN]: consul: reconcile unable to talk with Consul backend: %v", err)
							time.Sleep(consulRetryInterval)
							continue
						}

						c.serviceLock.Lock()
						defer c.serviceLock.Unlock()

						registeredServiceID = serviceID
						return
					}
				}()
			}
		case <-checkTimer.C:
			checkTimer.Reset(c.checkDuration())
			// Abort if service discovery is disabled or a
			// reconcile handler is active
			if !c.disableRegistration && atomic.CompareAndSwapInt64(&checkLock, 0, 1) {
				// Enter handler with checkLock held
				go func() {
					defer atomic.CompareAndSwapInt64(&checkLock, 1, 0)
					for !shutdown {
						sealed := sealedFunc()
						if err := c.runCheck(sealed); err != nil {
							c.logger.Printf("[WARN]: consul: check unable to talk with Consul backend: %v", err)
							time.Sleep(consulRetryInterval)
							continue
						}
						return
					}
				}()
			}
		case <-shutdownCh:
			c.logger.Printf("[INFO]: consul: Shutting down consul backend")
			shutdown = true
			break shutdown
		}
	}

	c.serviceLock.RLock()
	defer c.serviceLock.RUnlock()
	if err := c.client.Agent().ServiceDeregister(registeredServiceID); err != nil {
		c.logger.Printf("[WARN]: consul: service deregistration failed: %v", err)
	}
}

// checkID returns the ID used for a Consul Check.  Assume at least a read
// lock is held.
func (c *ConsulBackend) checkID() string {
	return fmt.Sprintf("%s:vault-sealed-check", c.serviceID())
}

// serviceID returns the Vault ServiceID for use in Consul.  Assume at least
// a read lock is held.
func (c *ConsulBackend) serviceID() string {
	return fmt.Sprintf("%s:%s:%d", c.serviceName, c.advertiseHost, c.advertisePort)
}

// reconcileConsul queries the state of Vault Core and Consul and fixes up
// Consul's state according to what's in Vault.  reconcileConsul is called
// without any locks held and can be run concurrently, therefore no changes
// to ConsulBackend can be made in this method (i.e. wtb const receiver for
// compiler enforced safety).
func (c *ConsulBackend) reconcileConsul(registeredServiceID string, activeFunc activeFunction, sealedFunc sealedFunction) (serviceID string, err error) {
	// Query vault Core for its current state
	active := activeFunc()
	sealed := sealedFunc()

	agent := c.client.Agent()
	catalog := c.client.Catalog()

	serviceID = c.serviceID()

	// Get the current state of Vault from Consul
	var currentVaultService *api.CatalogService
	if services, _, err := catalog.Service(c.serviceName, "", &api.QueryOptions{AllowStale: true}); err == nil {
		for _, service := range services {
			if serviceID == service.ServiceID {
				currentVaultService = service
				break
			}
		}
	}

	tags := serviceTags(active)

	var reregister bool
	switch {
	case currentVaultService == nil,
		registeredServiceID == "":
		reregister = true
	default:
		switch {
		case len(currentVaultService.ServiceTags) != 1,
			currentVaultService.ServiceTags[0] != tags[0]:
			reregister = true
		}
	}

	if !reregister {
		// When re-registration is not required, return a valid serviceID
		// to avoid registration in the next cycle.
		return serviceID, nil
	}

	service := &api.AgentServiceRegistration{
		ID:                serviceID,
		Name:              c.serviceName,
		Tags:              tags,
		Port:              int(c.advertisePort),
		Address:           c.advertiseHost,
		EnableTagOverride: false,
	}

	checkStatus := api.HealthCritical
	if !sealed {
		checkStatus = api.HealthPassing
	}

	sealedCheck := &api.AgentCheckRegistration{
		ID:        c.checkID(),
		Name:      "Vault Sealed Status",
		Notes:     "Vault service is healthy when Vault is in an unsealed status and can become an active Vault server",
		ServiceID: serviceID,
		AgentServiceCheck: api.AgentServiceCheck{
			TTL:    c.checkTimeout.String(),
			Status: checkStatus,
		},
	}

	if err := agent.ServiceRegister(service); err != nil {
		return "", errwrap.Wrapf(`service registration failed: {{err}}`, err)
	}

	if err := agent.CheckRegister(sealedCheck); err != nil {
		return serviceID, errwrap.Wrapf(`service check registration failed: {{err}}`, err)
	}

	return serviceID, nil
}

// runCheck immediately pushes a TTL check.
func (c *ConsulBackend) runCheck(sealed bool) error {
	// Run a TTL check
	agent := c.client.Agent()
	if !sealed {
		return agent.PassTTL(c.checkID(), "Vault Unsealed")
	} else {
		return agent.FailTTL(c.checkID(), "Vault Sealed")
	}
}

// serviceTags returns all of the relevant tags for Consul.
func serviceTags(active bool) []string {
	activeTag := "standby"
	if active {
		activeTag = "active"
	}
	return []string{activeTag}
}

func (c *ConsulBackend) setAdvertiseAddr(addr string) (err error) {
	if addr == "" {
		return fmt.Errorf("advertise address must not be empty")
	}

	url, err := url.Parse(addr)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf(`failed to parse advertise URL "%v": {{err}}`, addr), err)
	}

	var portStr string
	c.advertiseHost, portStr, err = net.SplitHostPort(url.Host)
	if err != nil {
		if url.Scheme == "http" {
			portStr = "80"
		} else if url.Scheme == "https" {
			portStr = "443"
		} else if url.Scheme == "unix" {
			portStr = "-1"
			c.advertiseHost = url.Path
		} else {
			return errwrap.Wrapf(fmt.Sprintf(`failed to find a host:port in advertise address "%v": {{err}}`, url.Host), err)
		}
	}
	c.advertisePort, err = strconv.ParseInt(portStr, 10, 0)
	if err != nil || c.advertisePort < -1 || c.advertisePort > 65535 {
		return errwrap.Wrapf(fmt.Sprintf(`failed to parse valid port "%v": {{err}}`, portStr), err)
	}

	return nil
}
