package consul

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/tlsutil"
	sr "github.com/hashicorp/vault/serviceregistration"
	"golang.org/x/net/http2"
)

const (
	// checkJitterFactor specifies the jitter factor used to stagger checks
	checkJitterFactor = 16

	// checkMinBuffer specifies provides a guarantee that a check will not
	// be executed too close to the TTL check timeout
	checkMinBuffer = 100 * time.Millisecond

	// defaultCheckTimeout changes the timeout of TTL checks
	defaultCheckTimeout = 5 * time.Second

	// DefaultServiceName is the default Consul service name used when
	// advertising a Vault instance.
	DefaultServiceName = "vault"

	// reconcileTimeout is how often Vault should query Consul to detect
	// and fix any state drift.
	reconcileTimeout = 60 * time.Second

	// These tags reflect Vault's state at a point in time.
	tagPerfStandby    = "performance-standby"
	tagNotPerfStandby = "not-performance-standby"
	tagIsActive       = "active"
	tagNotActive      = "inactive"
	tagInitialized    = "initialized"
	tagUninitialized  = "uninitialized"
	tagSealed         = "sealed"
	tagUnsealed       = "unsealed"
)

var (
	hostnameRegex = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
)

// ServiceRegistration is a ServiceRegistration that advertises the state of
// Vault to Consul.
type ServiceRegistration struct {
	Client *api.Client

	logger    log.Logger
	state     *sr.State
	stateLock sync.RWMutex

	redirectHost        string
	redirectPort        int64
	serviceName         string
	serviceAddress      *string
	usersTags           []string
	disableRegistration bool
	checkTimeout        time.Duration
}

// NewServiceRegistration constructs a Consul-based ServiceRegistration.
func NewServiceRegistration(shutdownCh <-chan struct{}, conf map[string]string, logger log.Logger, state *sr.State, redirectAddr string) (sr.ServiceRegistration, error) {

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

		min, _ := durationMinusBufferDomain(d, checkMinBuffer, checkJitterFactor)
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
		tlsClientConfig, err := tlsutil.SetupTLSConfig(conf, consulConf.Address)
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
	c := &ServiceRegistration{
		Client: client,

		logger:              logger,
		state:               state,
		serviceName:         service,
		usersTags:           strutil.ParseDedupLowercaseAndSortStrings(tags, ","),
		serviceAddress:      serviceAddr,
		disableRegistration: disableRegistration,
		checkTimeout:        checkTimeout,
	}

	if err := c.setRedirectAddr(redirectAddr); err != nil {
		return nil, err
	}

	// Do an initial reconciliation to register Vault with Consul.
	if err := c.reconcileConsul(); err != nil {
		return nil, err
	}

	go c.RunOngoingReconciliations(shutdownCh)
	go c.RunOngoingChecks(shutdownCh, checkTimeout)
	go c.GoToFinalState(shutdownCh)

	return c, nil
}

func (c *ServiceRegistration) RunOngoingReconciliations(shutdownCh <-chan struct{}) {
	for {
		select {
		case <-shutdownCh:
			return
		case <-time.After(reconcileTimeout):
			if c.logger.IsDebug() {
				c.logger.Debug("reconciling")
			}
			if err := c.reconcileConsul(); err != nil {
				if c.logger.IsWarn() {
					c.logger.Warn(fmt.Sprintf("unable to reconcile consul: %s", err))
				}
			}
		}
	}
}

func (c *ServiceRegistration) RunOngoingChecks(shutdownCh <-chan struct{}, checkTimeout time.Duration) {
	for {
		select {
		case <-shutdownCh:
			return
		case <-time.After(addJitter(checkTimeout)):
			if c.logger.IsDebug() {
				c.logger.Debug("checking")
			}
			if err := c.runCheck(); err != nil {
				if c.logger.IsWarn() {
					c.logger.Warn(fmt.Sprintf("unable to check consul: %s", err))
				}
			}
		}
	}
}

func addJitter(checkTimeout time.Duration) time.Duration {
	d := checkTimeout - checkMinBuffer
	intv := time.Duration(int64(d) / checkJitterFactor)
	if intv == 0 {
		return 0
	}
	d -= time.Duration(uint64(rand.Int63()) % uint64(intv))
	return d
}

func (c *ServiceRegistration) NotifyActiveStateChange(isActive bool) error {
	if c.logger.IsDebug() {
		c.logger.Debug(fmt.Sprintf("received NotifyActiveStateChange isActive: %v", isActive))
	}
	c.stateLock.Lock()
	c.state.IsActive = isActive
	c.stateLock.Unlock()
	return c.reconcileConsul()
}

func (c *ServiceRegistration) NotifyPerformanceStandbyStateChange(isStandby bool) error {
	if c.logger.IsDebug() {
		c.logger.Debug(fmt.Sprintf("received NotifyPerformanceStandbyStateChange isStandby: %v", isStandby))
	}
	c.stateLock.Lock()
	c.state.IsPerformanceStandby = isStandby
	c.stateLock.Unlock()
	return c.runCheck()
}

func (c *ServiceRegistration) NotifySealedStateChange(isSealed bool) error {
	if c.logger.IsDebug() {
		c.logger.Debug(fmt.Sprintf("received NotifySealedStateChange isSealed: %v", isSealed))
	}
	c.stateLock.Lock()
	c.state.IsSealed = isSealed
	c.stateLock.Unlock()
	return c.runCheck()
}

func (c *ServiceRegistration) NotifyInitializedStateChange(isInitialized bool) error {
	if c.logger.IsDebug() {
		c.logger.Debug(fmt.Sprintf("received NotifyInitializedStateChange isSealed: %v", isInitialized))
	}
	c.stateLock.Lock()
	c.state.IsInitialized = isInitialized
	c.stateLock.Unlock()
	return c.runCheck()
}

func (c *ServiceRegistration) GoToFinalState(shutdownCh <-chan struct{}) {
	if c.disableRegistration {
		// Nothing further to do here.
		return
	}
	<-shutdownCh
	c.stateLock.Lock()
	c.state.IsSealed = true
	c.state.IsInitialized = false
	c.state.IsPerformanceStandby = false
	c.state.IsActive = false
	c.stateLock.Unlock()

	// Try to register one last time to make sure final tags look correct.
	if err := c.reconcileConsul(); err != nil {
		if c.logger.IsWarn() {
			c.logger.Warn("final reconciliation failed", "error", err)
		}
	}

	// Deregister the service.
	if err := c.Client.Agent().ServiceDeregister(c.serviceID()); err != nil {
		if c.logger.IsWarn() {
			c.logger.Warn("service deregistration failed", "error", err)
		}
	}
}

// checkID returns the ID used for a Consul Check.  Assume at least a read
// lock is held.
func (c *ServiceRegistration) checkID() string {
	return fmt.Sprintf("%s:vault-sealed-check", c.serviceID())
}

// serviceID returns the Vault ServiceID for use in Consul.  Assume at least
// a read lock is held.
func (c *ServiceRegistration) serviceID() string {
	return fmt.Sprintf("%s:%s:%d", c.serviceName, c.redirectHost, c.redirectPort)
}

// reconcileConsul queries the state of Vault Core and Consul and fixes up
// Consul's state according to what's in Vault.  reconcileConsul is called
// without any locks held and can be run concurrently, therefore no changes
// to ServiceRegistration can be made in this method (i.e. wtb const receiver for
// compiler enforced safety).
func (c *ServiceRegistration) reconcileConsul() error {
	if c.disableRegistration {
		// Nothing further to do here.
		return nil
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()

	agent := c.Client.Agent()
	catalog := c.Client.Catalog()

	serviceID := c.serviceID()

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

	tags := buildTags(c.usersTags, c.state)

	var reregister bool

	switch {
	case currentVaultService == nil:
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
		return nil
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
	if !c.state.IsSealed {
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
		return errwrap.Wrapf(`service registration failed: {{err}}`, err)
	}

	if err := agent.CheckRegister(sealedCheck); err != nil {
		return errwrap.Wrapf(`service check registration failed: {{err}}`, err)
	}

	return nil
}

// runCheck immediately pushes a TTL check.
func (c *ServiceRegistration) runCheck() error {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()

	// Run a TTL check
	agent := c.Client.Agent()
	if !c.state.IsSealed {
		return agent.PassTTL(c.checkID(), "Vault Unsealed")
	} else {
		return agent.FailTTL(c.checkID(), "Vault Sealed")
	}
}

func (c *ServiceRegistration) setRedirectAddr(addr string) (err error) {
	if addr == "" {
		return fmt.Errorf("redirect address must not be empty")
	}

	u, err := url.Parse(addr)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("failed to parse redirect URL %q: {{err}}", addr), err)
	}

	var portStr string
	c.redirectHost, portStr, err = net.SplitHostPort(u.Host)
	if err != nil {
		if u.Scheme == "http" {
			portStr = "80"
		} else if u.Scheme == "https" {
			portStr = "443"
		} else if u.Scheme == "unix" {
			portStr = "-1"
			c.redirectHost = u.Path
		} else {
			return errwrap.Wrapf(fmt.Sprintf(`failed to find a host:port in redirect address "%v": {{err}}`, u.Host), err)
		}
	}
	c.redirectPort, err = strconv.ParseInt(portStr, 10, 0)
	if err != nil || c.redirectPort < -1 || c.redirectPort > 65535 {
		return errwrap.Wrapf(fmt.Sprintf(`failed to parse valid port "%v": {{err}}`, portStr), err)
	}

	return nil
}

// durationMinusBufferDomain returns the domain of valid durations from a
// call to durationMinusBuffer.  This function is used to check user
// specified input values to durationMinusBuffer.
func durationMinusBufferDomain(intv time.Duration, buffer time.Duration, jitter int64) (min time.Duration, max time.Duration) {
	max = intv - buffer
	if jitter == 0 {
		min = max
	} else {
		min = max - time.Duration(int64(max)/jitter)
	}
	return min, max
}

func buildTags(usersTags []string, state *sr.State) []string {
	result := usersTags
	if state.IsPerformanceStandby {
		result = append(result, tagPerfStandby)
	} else {
		result = append(result, tagNotPerfStandby)
	}
	if state.IsActive {
		result = append(result, tagIsActive)
	} else {
		result = append(result, tagNotActive)
	}
	if state.IsInitialized {
		result = append(result, tagInitialized)
	} else {
		result = append(result, tagUninitialized)
	}
	if state.IsSealed {
		result = append(result, tagSealed)
	} else {
		result = append(result, tagUnsealed)
	}
	return result
}
