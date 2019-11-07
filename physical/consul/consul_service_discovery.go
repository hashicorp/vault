package consul

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	"github.com/hashicorp/vault/sdk/physical"
	"golang.org/x/net/http2"
)

var _ physical.ServiceDiscovery = (*ConsulServiceDiscovery)(nil)

// ConsulServiceDiscovery is a ServiceDiscovery that advertises the state of
// Vault to Consul.
type ConsulServiceDiscovery struct {
	logger              log.Logger
	client              *api.Client
	serviceLock         sync.RWMutex
	redirectHost        string
	redirectPort        int64
	serviceName         string
	serviceTags         []string
	serviceAddress      *string
	disableRegistration bool
	checkTimeout        time.Duration

	notifyActiveCh      chan notifyEvent
	notifySealedCh      chan notifyEvent
	notifyPerfStandbyCh chan notifyEvent
}

// NewConsulServiceDiscovery constructs a Consul-based ServiceDiscovery.
func NewConsulServiceDiscovery(conf map[string]string, logger log.Logger) (physical.ServiceDiscovery, error) {

	// Allow admins to disable consul integration
	disableReg, ok := conf["disable_registration"]
	var disableRegistration bool
	if ok && disableReg != "" {
		b, err := parseutil.ParseBool(disableReg)
		if err != nil {
			return nil, errwrap.Wrapf("failed parsing disable_registration parameter: {{err}}", err)
		}
		disableRegistration = b
	}
	if logger.IsDebug() {
		logger.Debug("config disable_registration set", "disable_registration", disableRegistration)
	}

	// Get the service name to advertise in Consul
	service, ok := conf["service"]
	if !ok {
		service = DefaultServiceName
	}
	if !hostnameRegex.MatchString(service) {
		return nil, errors.New("service name must be valid per RFC 1123 and can contain only alphanumeric characters or dashes")
	}
	if logger.IsDebug() {
		logger.Debug("config service set", "service", service)
	}

	// Get the additional tags to attach to the registered service name
	tags := conf["service_tags"]
	if logger.IsDebug() {
		logger.Debug("config service_tags set", "service_tags", tags)
	}

	// Get the service-specific address to override the use of the HA redirect address
	var serviceAddr *string
	serviceAddrStr, ok := conf["service_address"]
	if ok {
		serviceAddr = &serviceAddrStr
	}
	if logger.IsDebug() {
		logger.Debug("config service_address set", "service_address", serviceAddr)
	}

	checkTimeout := defaultCheckTimeout
	checkTimeoutStr, ok := conf["check_timeout"]
	if ok {
		d, err := parseutil.ParseDurationSecond(checkTimeoutStr)
		if err != nil {
			return nil, err
		}

		min, _ := DurationMinusBufferDomain(d, checkMinBuffer, checkJitterFactor)
		if min < checkMinBuffer {
			return nil, fmt.Errorf("consul check_timeout must be greater than %v", min)
		}

		checkTimeout = d
		if logger.IsDebug() {
			logger.Debug("config check_timeout set", "check_timeout", d)
		}
	}

	// Configure the client
	consulConf := api.DefaultConfig()
	// Set MaxIdleConnsPerHost to the number of processes used in expiration.Restore
	consulConf.Transport.MaxIdleConnsPerHost = consts.ExpirationRestoreWorkerCount

	if addr, ok := conf["address"]; ok {
		consulConf.Address = addr
		if logger.IsDebug() {
			logger.Debug("config address set", "address", addr)
		}

		// Copied from the Consul API module; set the Scheme based on
		// the protocol field if address looks ike a URL.
		// This can enable the TLS configuration below.
		parts := strings.SplitN(addr, "://", 2)
		if len(parts) == 2 {
			if parts[0] == "http" || parts[0] == "https" {
				consulConf.Scheme = parts[0]
				consulConf.Address = parts[1]
				if logger.IsDebug() {
					logger.Debug("config address parsed", "scheme", parts[0])
					logger.Debug("config scheme parsed", "address", parts[1])
				}
			} // allow "unix:" or whatever else consul supports in the future
		}
	}
	if scheme, ok := conf["scheme"]; ok {
		consulConf.Scheme = scheme
		if logger.IsDebug() {
			logger.Debug("config scheme set", "scheme", scheme)
		}
	}
	if token, ok := conf["token"]; ok {
		consulConf.Token = token
		logger.Debug("config token set")
	}

	if consulConf.Scheme == "https" {
		// Use the parsed Address instead of the raw conf['address']
		tlsClientConfig, err := setupTLSConfig(conf, consulConf.Address)
		if err != nil {
			return nil, err
		}

		consulConf.Transport.TLSClientConfig = tlsClientConfig
		if err := http2.ConfigureTransport(consulConf.Transport); err != nil {
			return nil, err
		}
		logger.Debug("configured TLS")
	}

	consulConf.HttpClient = &http.Client{Transport: consulConf.Transport}
	client, err := api.NewClient(consulConf)
	if err != nil {
		return nil, errwrap.Wrapf("client setup failed: {{err}}", err)
	}

	// Setup the backend
	c := &ConsulServiceDiscovery{
		logger:              logger,
		client:              client,
		serviceName:         service,
		serviceTags:         strutil.ParseDedupLowercaseAndSortStrings(tags, ","),
		serviceAddress:      serviceAddr,
		checkTimeout:        checkTimeout,
		disableRegistration: disableRegistration,

		notifyActiveCh:      make(chan notifyEvent),
		notifySealedCh:      make(chan notifyEvent),
		notifyPerfStandbyCh: make(chan notifyEvent),
	}
	return c, nil
}

func setupTLSConfig(conf map[string]string, address string) (*tls.Config, error) {
	serverName, _, err := net.SplitHostPort(address)
	switch {
	case err == nil:
	case strings.Contains(err.Error(), "missing port"):
		serverName = conf["address"]
	default:
		return nil, err
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
			return nil, errwrap.Wrapf("client tls setup failed: {{err}}", err)
		}

		tlsClientConfig.Certificates = []tls.Certificate{tlsCert}
	}

	if tlsCaFile, ok := conf["tls_ca_file"]; ok {
		caPool := x509.NewCertPool()

		data, err := ioutil.ReadFile(tlsCaFile)
		if err != nil {
			return nil, errwrap.Wrapf("failed to read CA file: {{err}}", err)
		}

		if !caPool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}

		tlsClientConfig.RootCAs = caPool
	}

	return tlsClientConfig, nil
}

func (c *ConsulServiceDiscovery) NotifyActiveStateChange() error {
	select {
	case c.notifyActiveCh <- notifyEvent{}:
	default:
		// NOTE: If this occurs Vault's active status could be out of
		// sync with Consul until reconcileTimer expires.
		c.logger.Warn("concurrent state change notify dropped")
	}

	return nil
}

func (c *ConsulServiceDiscovery) NotifyPerformanceStandbyStateChange() error {
	select {
	case c.notifyPerfStandbyCh <- notifyEvent{}:
	default:
		// NOTE: If this occurs Vault's active status could be out of
		// sync with Consul until reconcileTimer expires.
		c.logger.Warn("concurrent state change notify dropped")
	}

	return nil
}

func (c *ConsulServiceDiscovery) NotifySealedStateChange() error {
	select {
	case c.notifySealedCh <- notifyEvent{}:
	default:
		// NOTE: If this occurs Vault's sealed status could be out of
		// sync with Consul until checkTimer expires.
		c.logger.Warn("concurrent sealed state change notify dropped")
	}

	return nil
}

func (c *ConsulServiceDiscovery) checkDuration() time.Duration {
	return DurationMinusBuffer(c.checkTimeout, checkMinBuffer, checkJitterFactor)
}

func (c *ConsulServiceDiscovery) RunServiceDiscovery(waitGroup *sync.WaitGroup, shutdownCh physical.ShutdownChannel, redirectAddr string, activeFunc physical.ActiveFunction, sealedFunc physical.SealedFunction, perfStandbyFunc physical.PerformanceStandbyFunction) (err error) {
	if err := c.setRedirectAddr(redirectAddr); err != nil {
		return err
	}

	// 'server' command will wait for the below goroutine to complete
	waitGroup.Add(1)

	go c.runEventDemuxer(waitGroup, shutdownCh, redirectAddr, activeFunc, sealedFunc, perfStandbyFunc)

	return nil
}

func (c *ConsulServiceDiscovery) runEventDemuxer(waitGroup *sync.WaitGroup, shutdownCh physical.ShutdownChannel, redirectAddr string, activeFunc physical.ActiveFunction, sealedFunc physical.SealedFunction, perfStandbyFunc physical.PerformanceStandbyFunction) {
	// This defer statement should be executed last. So push it first.
	defer waitGroup.Done()

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
	var registeredServiceID string
	checkLock := new(int32)
	serviceRegLock := new(int32)

	for !shutdown {
		select {
		case <-c.notifyActiveCh:
			// Run reconcile immediately upon active state change notification
			reconcileTimer.Reset(0)
		case <-c.notifySealedCh:
			// Run check timer immediately upon a seal state change notification
			checkTimer.Reset(0)
		case <-c.notifyPerfStandbyCh:
			// Run check timer immediately upon a seal state change notification
			checkTimer.Reset(0)
		case <-reconcileTimer.C:
			// Unconditionally rearm the reconcileTimer
			reconcileTimer.Reset(reconcileTimeout - RandomStagger(reconcileTimeout/checkJitterFactor))

			// Abort if service discovery is disabled or a
			// reconcile handler is already active
			if !c.disableRegistration && atomic.CompareAndSwapInt32(serviceRegLock, 0, 1) {
				// Enter handler with serviceRegLock held
				go func() {
					defer atomic.CompareAndSwapInt32(serviceRegLock, 1, 0)
					for !shutdown {
						serviceID, err := c.reconcileConsul(registeredServiceID, activeFunc, sealedFunc, perfStandbyFunc)
						if err != nil {
							if c.logger.IsWarn() {
								c.logger.Warn("reconcile unable to talk with Consul backend", "error", err)
							}
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
			if !c.disableRegistration && atomic.CompareAndSwapInt32(checkLock, 0, 1) {
				// Enter handler with checkLock held
				go func() {
					defer atomic.CompareAndSwapInt32(checkLock, 1, 0)
					for !shutdown {
						sealed := sealedFunc()
						if err := c.runCheck(sealed); err != nil {
							if c.logger.IsWarn() {
								c.logger.Warn("check unable to talk with Consul backend", "error", err)
							}
							time.Sleep(consulRetryInterval)
							continue
						}
						return
					}
				}()
			}
		case <-shutdownCh:
			c.logger.Info("shutting down consul backend")
			shutdown = true
		}
	}

	c.serviceLock.RLock()
	defer c.serviceLock.RUnlock()
	if err := c.client.Agent().ServiceDeregister(registeredServiceID); err != nil {
		if c.logger.IsWarn() {
			c.logger.Warn("service deregistration failed", "error", err)
		}
	}
}

// checkID returns the ID used for a Consul Check.  Assume at least a read
// lock is held.
func (c *ConsulServiceDiscovery) checkID() string {
	return fmt.Sprintf("%s:vault-sealed-check", c.serviceID())
}

// serviceID returns the Vault ServiceID for use in Consul.  Assume at least
// a read lock is held.
func (c *ConsulServiceDiscovery) serviceID() string {
	return fmt.Sprintf("%s:%s:%d", c.serviceName, c.redirectHost, c.redirectPort)
}

// reconcileConsul queries the state of Vault Core and Consul and fixes up
// Consul's state according to what's in Vault.  reconcileConsul is called
// without any locks held and can be run concurrently, therefore no changes
// to ConsulServiceDiscovery can be made in this method (i.e. wtb const receiver for
// compiler enforced safety).
func (c *ConsulServiceDiscovery) reconcileConsul(registeredServiceID string, activeFunc physical.ActiveFunction, sealedFunc physical.SealedFunction, perfStandbyFunc physical.PerformanceStandbyFunction) (serviceID string, err error) {
	// Query vault Core for its current state
	active := activeFunc()
	sealed := sealedFunc()
	perfStandby := perfStandbyFunc()

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

	tags := c.fetchServiceTags(active, perfStandby)

	var reregister bool

	switch {
	case currentVaultService == nil, registeredServiceID == "":
		reregister = true
	default:
		switch {
		case !strutil.EquivalentSlices(currentVaultService.ServiceTags, tags):
			reregister = true
		}
	}

	if !reregister {
		// When re-registration is not required, return a valid serviceID
		// to avoid registration in the next cycle.
		return serviceID, nil
	}

	// If service address was set explicitly in configuration, use that
	// as the service-specific address instead of the HA redirect address.
	var serviceAddress string
	if c.serviceAddress == nil {
		serviceAddress = c.redirectHost
	} else {
		serviceAddress = *c.serviceAddress
	}

	service := &api.AgentServiceRegistration{
		ID:                serviceID,
		Name:              c.serviceName,
		Tags:              tags,
		Port:              int(c.redirectPort),
		Address:           serviceAddress,
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
func (c *ConsulServiceDiscovery) runCheck(sealed bool) error {
	// Run a TTL check
	agent := c.client.Agent()
	if !sealed {
		return agent.PassTTL(c.checkID(), "Vault Unsealed")
	} else {
		return agent.FailTTL(c.checkID(), "Vault Sealed")
	}
}

// fetchServiceTags returns all of the relevant tags for Consul.
func (c *ConsulServiceDiscovery) fetchServiceTags(active bool, perfStandby bool) []string {
	activeTag := "standby"
	if active {
		activeTag = "active"
	}

	result := append(c.serviceTags, activeTag)

	if perfStandby {
		result = append(c.serviceTags, "performance-standby")
	}

	return result
}

func (c *ConsulServiceDiscovery) setRedirectAddr(addr string) (err error) {
	if addr == "" {
		return fmt.Errorf("redirect address must not be empty")
	}

	url, err := url.Parse(addr)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("failed to parse redirect URL %q: {{err}}", addr), err)
	}

	var portStr string
	c.redirectHost, portStr, err = net.SplitHostPort(url.Host)
	if err != nil {
		if url.Scheme == "http" {
			portStr = "80"
		} else if url.Scheme == "https" {
			portStr = "443"
		} else if url.Scheme == "unix" {
			portStr = "-1"
			c.redirectHost = url.Path
		} else {
			return errwrap.Wrapf(fmt.Sprintf(`failed to find a host:port in redirect address "%v": {{err}}`, url.Host), err)
		}
	}
	c.redirectPort, err = strconv.ParseInt(portStr, 10, 0)
	if err != nil || c.redirectPort < -1 || c.redirectPort > 65535 {
		return errwrap.Wrapf(fmt.Sprintf(`failed to parse valid port "%v": {{err}}`, portStr), err)
	}

	return nil
}
