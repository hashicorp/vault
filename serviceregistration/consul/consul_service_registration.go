// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package consul

import (
	"context"
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
	"sync/atomic"
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/vault/diagnose"
	atomicB "go.uber.org/atomic"
	"golang.org/x/net/http2"
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

	// DefaultServiceName is the default Consul service name used when
	// advertising a Vault instance.
	DefaultServiceName = "vault"

	// reconcileTimeout is how often Vault should query Consul to detect
	// and fix any state drift.
	DefaultReconcileTimeout = 60 * time.Second

	// metaExternalSource is a metadata value for external-source that can be
	// used by the Consul UI.
	metaExternalSource = "vault"
)

var hostnameRegex = regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)

// serviceRegistration is a ServiceRegistration that advertises the state of
// Vault to Consul.
type serviceRegistration struct {
	Client *api.Client
	config *api.Config

	logger              log.Logger
	serviceLock         sync.RWMutex
	registeredServiceID string
	redirectHost        string
	redirectPort        int64
	serviceName         string
	serviceTags         []string
	serviceAddress      *string
	disableRegistration bool
	checkTimeout        time.Duration
	reconcileTimeout    time.Duration

	notifyActiveCh      chan struct{}
	notifySealedCh      chan struct{}
	notifyPerfStandbyCh chan struct{}
	notifyInitializedCh chan struct{}

	isActive      *atomicB.Bool
	isSealed      *atomicB.Bool
	isPerfStandby *atomicB.Bool
	isInitialized *atomicB.Bool
}

// NewConsulServiceRegistration constructs a Consul-based ServiceRegistration.
func NewServiceRegistration(conf map[string]string, logger log.Logger, state sr.State) (sr.ServiceRegistration, error) {
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	// Setup the backend
	c := &serviceRegistration{
		logger: logger,

		notifyActiveCh:      make(chan struct{}),
		notifySealedCh:      make(chan struct{}),
		notifyPerfStandbyCh: make(chan struct{}),
		notifyInitializedCh: make(chan struct{}),

		isActive:      atomicB.NewBool(state.IsActive),
		isSealed:      atomicB.NewBool(state.IsSealed),
		isPerfStandby: atomicB.NewBool(state.IsPerformanceStandby),
		isInitialized: atomicB.NewBool(state.IsInitialized),
	}

	c.serviceLock.Lock()
	defer c.serviceLock.Unlock()
	err := c.merge(conf)
	return c, err
}

func SetupSecureTLS(ctx context.Context, consulConf *api.Config, conf map[string]string, logger log.Logger, isDiagnose bool) error {
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
		if isDiagnose {
			certPath, okCert := conf["tls_cert_file"]
			keyPath, okKey := conf["tls_key_file"]
			if okCert && okKey {
				warnings, err := diagnose.TLSFileChecks(certPath, keyPath)
				for _, warning := range warnings {
					diagnose.Warn(ctx, warning)
				}
				if err != nil {
					return err
				}
				return nil
			}
			return fmt.Errorf("key or cert path: %s, %s, cannot be loaded from consul config file", certPath, keyPath)
		}

		// Use the parsed Address instead of the raw conf['address']
		tlsClientConfig, err := tlsutil.SetupTLSConfig(conf, consulConf.Address)
		if err != nil {
			return err
		}

		consulConf.Transport.TLSClientConfig = tlsClientConfig
		if err := http2.ConfigureTransport(consulConf.Transport); err != nil {
			return err
		}
		logger.Debug("configured TLS")
	} else {
		if isDiagnose {
			diagnose.Skipped(ctx, "HTTPS is not used, Skipping TLS verification.")
		}
	}
	return nil
}

func (c *serviceRegistration) Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup, redirectAddr string) error {
	go func() {
		if err := c.runServiceRegistration(wait, shutdownCh, redirectAddr); err != nil {
			if c.logger.IsError() {
				c.logger.Error(fmt.Sprintf("error running service registration: %s", err))
			}
		}
	}()
	return nil
}

func (c *serviceRegistration) merge(conf map[string]string) error {
	// Allow admins to disable consul integration
	disableReg, ok := conf["disable_registration"]
	var disableRegistration bool
	if ok && disableReg != "" {
		b, err := parseutil.ParseBool(disableReg)
		if err != nil {
			return fmt.Errorf("failed parsing disable_registration parameter: %w", err)
		}
		disableRegistration = b
	}
	if c.logger.IsDebug() {
		c.logger.Debug("config disable_registration set", "disable_registration", disableRegistration)
	}

	// Get the service name to advertise in Consul
	service, ok := conf["service"]
	if !ok {
		service = DefaultServiceName
	}
	if !hostnameRegex.MatchString(service) {
		return errors.New("service name must be valid per RFC 1123 and can contain only alphanumeric characters or dashes")
	}
	if c.logger.IsDebug() {
		c.logger.Debug("config service set", "service", service)
	}

	// Get the additional tags to attach to the registered service name
	tags := conf["service_tags"]
	if c.logger.IsDebug() {
		c.logger.Debug("config service_tags set", "service_tags", tags)
	}

	// Get the service-specific address to override the use of the HA redirect address
	var serviceAddr *string
	serviceAddrStr, ok := conf["service_address"]
	if ok {
		serviceAddr = &serviceAddrStr
	}
	if c.logger.IsDebug() {
		c.logger.Debug("config service_address set", "service_address", serviceAddrStr)
	}

	checkTimeout := defaultCheckTimeout
	checkTimeoutStr, ok := conf["check_timeout"]
	if ok {
		d, err := parseutil.ParseDurationSecond(checkTimeoutStr)
		if err != nil {
			return err
		}

		min, _ := durationMinusBufferDomain(d, checkMinBuffer, checkJitterFactor)
		if min < checkMinBuffer {
			return fmt.Errorf("consul check_timeout must be greater than %v", min)
		}

		checkTimeout = d
		if c.logger.IsDebug() {
			c.logger.Debug("config check_timeout set", "check_timeout", d)
		}
	}

	reconcileTimeout := DefaultReconcileTimeout
	reconcileTimeoutStr, ok := conf["reconcile_timeout"]
	if ok {
		d, err := parseutil.ParseDurationSecond(reconcileTimeoutStr)
		if err != nil {
			return err
		}

		min, _ := durationMinusBufferDomain(d, checkMinBuffer, checkJitterFactor)
		if min < checkMinBuffer {
			return fmt.Errorf("consul reconcile_timeout must be greater than %v", min)
		}

		reconcileTimeout = d
		if c.logger.IsDebug() {
			c.logger.Debug("config reconcile_timeout set", "reconcile_timeout", d)
		}
	}

	// Configure the client
	consulConf := api.DefaultConfig()
	// Set MaxIdleConnsPerHost to the number of processes used in expiration.Restore
	consulConf.Transport.MaxIdleConnsPerHost = consts.ExpirationRestoreWorkerCount

	SetupSecureTLS(context.Background(), consulConf, conf, c.logger, false)

	consulConf.HttpClient = &http.Client{Transport: consulConf.Transport}
	client, err := api.NewClient(consulConf)
	if err != nil {
		return fmt.Errorf("client setup failed: %w", err)
	}

	c.Client = client
	c.config = consulConf
	c.serviceName = service
	c.serviceTags = strutil.ParseDedupAndSortStrings(tags, ",")
	c.serviceAddress = serviceAddr
	c.checkTimeout = checkTimeout
	c.disableRegistration = disableRegistration
	c.reconcileTimeout = reconcileTimeout

	return nil
}

func (c *serviceRegistration) NotifyActiveStateChange(isActive bool) error {
	c.isActive.Store(isActive)
	select {
	case c.notifyActiveCh <- struct{}{}:
	default:
		// NOTE: If this occurs Vault's active status could be out of
		// sync with Consul until reconcileTimer expires.
		c.logger.Warn("concurrent state change notify dropped")
	}

	return nil
}

func (c *serviceRegistration) NotifyPerformanceStandbyStateChange(isStandby bool) error {
	c.isPerfStandby.Store(isStandby)
	select {
	case c.notifyPerfStandbyCh <- struct{}{}:
	default:
		// NOTE: If this occurs Vault's active status could be out of
		// sync with Consul until reconcileTimer expires.
		c.logger.Warn("concurrent state change notify dropped")
	}

	return nil
}

func (c *serviceRegistration) NotifySealedStateChange(isSealed bool) error {
	c.isSealed.Store(isSealed)
	select {
	case c.notifySealedCh <- struct{}{}:
	default:
		// NOTE: If this occurs Vault's sealed status could be out of
		// sync with Consul until checkTimer expires.
		c.logger.Warn("concurrent sealed state change notify dropped")
	}

	return nil
}

func (c *serviceRegistration) NotifyInitializedStateChange(isInitialized bool) error {
	c.isInitialized.Store(isInitialized)
	select {
	case c.notifyInitializedCh <- struct{}{}:
	default:
		// NOTE: If this occurs Vault's initialized status could be out of
		// sync with Consul until checkTimer expires.
		c.logger.Warn("concurrent initialize state change notify dropped")
	}

	return nil
}

func (c *serviceRegistration) NotifyConfigurationReload(conf *map[string]string) error {
	c.serviceLock.Lock()
	defer c.serviceLock.Unlock()
	if conf == nil {
		if c.logger.IsDebug() {
			c.logger.Debug("registration is now empty, deregistering service from consul")
		}
		c.disableRegistration = true
		err := c.deregisterService()
		c.Client = nil
		return err
	} else {
		if c.logger.IsDebug() {
			c.logger.Debug("service registration configuration received, merging with existing configuation")
		}
		return c.merge(*conf)
	}
}

func (c *serviceRegistration) checkDuration() time.Duration {
	return durationMinusBuffer(c.checkTimeout, checkMinBuffer, checkJitterFactor)
}

func (c *serviceRegistration) runServiceRegistration(waitGroup *sync.WaitGroup, shutdownCh <-chan struct{}, redirectAddr string) (err error) {
	if err := c.setRedirectAddr(redirectAddr); err != nil {
		return err
	}

	// 'server' command will wait for the below goroutine to complete
	waitGroup.Add(1)

	go c.runEventDemuxer(waitGroup, shutdownCh)

	return nil
}

func (c *serviceRegistration) runEventDemuxer(waitGroup *sync.WaitGroup, shutdownCh <-chan struct{}) {
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
	var shutdown atomicB.Bool
	checkLock := new(int32)
	serviceRegLock := new(int32)

	for !shutdown.Load() {
		select {
		case <-c.notifyActiveCh:
			// Run reconcile immediately upon active state change notification
			reconcileTimer.Reset(0)
		case <-c.notifySealedCh:
			// Run check timer immediately upon a seal state change notification
			checkTimer.Reset(0)
		case <-c.notifyPerfStandbyCh:
			// Run check timer immediately upon a perfstandby state change notification
			checkTimer.Reset(0)
		case <-c.notifyInitializedCh:
			// Run check timer immediately upon an initialized state change notification
			checkTimer.Reset(0)
		case <-reconcileTimer.C:
			// Unconditionally rearm the reconcileTimer
			c.serviceLock.RLock()
			reconcileTimer.Reset(c.reconcileTimeout - randomStagger(c.reconcileTimeout/checkJitterFactor))
			disableRegistration := c.disableRegistration
			c.serviceLock.RUnlock()

			// Abort if service discovery is disabled or a
			// reconcile handler is already active
			if !disableRegistration && atomic.CompareAndSwapInt32(serviceRegLock, 0, 1) {
				// Enter handler with serviceRegLock held
				go func() {
					defer atomic.CompareAndSwapInt32(serviceRegLock, 1, 0)
					for !shutdown.Load() {
						serviceID, err := c.reconcileConsul()
						if err != nil {
							if c.logger.IsWarn() {
								c.logger.Warn("reconcile unable to talk with Consul backend", "error", err)
							}
							time.Sleep(consulRetryInterval)
							continue
						}

						c.serviceLock.Lock()
						c.registeredServiceID = serviceID
						c.serviceLock.Unlock()

						return
					}
				}()
			}
		case <-checkTimer.C:
			checkTimer.Reset(c.checkDuration())
			c.serviceLock.RLock()
			disableRegistration := c.disableRegistration
			c.serviceLock.RUnlock()

			// Abort if service discovery is disabled or a
			// reconcile handler is active
			if !disableRegistration && atomic.CompareAndSwapInt32(checkLock, 0, 1) {
				// Enter handler with checkLock held
				go func() {
					defer atomic.CompareAndSwapInt32(checkLock, 1, 0)
					for !shutdown.Load() {
						c.serviceLock.RLock()
						registeredServiceID := c.registeredServiceID
						c.serviceLock.RUnlock()

						if registeredServiceID != "" {
							if err := c.runCheck(c.isSealed.Load()); err != nil {
								if c.logger.IsWarn() {
									c.logger.Warn("check unable to talk with Consul backend", "error", err)
								}
								time.Sleep(consulRetryInterval)
								continue
							}
						}
						return
					}
				}()
			}
		case <-shutdownCh:
			c.logger.Info("shutting down consul backend")
			shutdown.Store(true)
		}
	}

	c.serviceLock.Lock()
	defer c.serviceLock.Unlock()
	c.deregisterService()
}

func (c *serviceRegistration) deregisterService() error {
	if c.registeredServiceID != "" {
		if err := c.Client.Agent().ServiceDeregister(c.registeredServiceID); err != nil {
			if c.logger.IsWarn() {
				c.logger.Warn("service deregistration failed", "error", err)
			}
			return err
		}
		c.registeredServiceID = ""
	}

	return nil
}

// checkID returns the ID used for a Consul Check.  Assume at least a read
// lock is held.
func (c *serviceRegistration) checkID() string {
	return fmt.Sprintf("%s:vault-sealed-check", c.serviceID())
}

// serviceID returns the Vault ServiceID for use in Consul.  Assume at least
// a read lock is held.
func (c *serviceRegistration) serviceID() string {
	return fmt.Sprintf("%s:%s:%d", c.serviceName, c.redirectHost, c.redirectPort)
}

// reconcileConsul queries the state of Vault Core and Consul and fixes up
// Consul's state according to what's in Vault.  reconcileConsul is called
// with a read lock and can be run concurrently, therefore no changes
// to serviceRegistration can be made in this method (i.e. wtb const receiver for
// compiler enforced safety).
func (c *serviceRegistration) reconcileConsul() (serviceID string, err error) {
	c.serviceLock.RLock()
	defer c.serviceLock.RUnlock()
	agent := c.Client.Agent()
	catalog := c.Client.Catalog()

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

	tags := c.fetchServiceTags(c.isActive.Load(), c.isPerfStandby.Load(), c.isInitialized.Load())

	var reregister bool

	switch {
	case currentVaultService == nil, c.registeredServiceID == "":
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
		Meta: map[string]string{
			"external-source": metaExternalSource,
		},
	}

	checkStatus := api.HealthCritical
	if !c.isSealed.Load() {
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
		return "", fmt.Errorf(`service registration failed: %w`, err)
	}

	if err := agent.CheckRegister(sealedCheck); err != nil {
		return serviceID, fmt.Errorf(`service check registration failed: %w`, err)
	}

	return serviceID, nil
}

// runCheck immediately pushes a TTL check.
func (c *serviceRegistration) runCheck(sealed bool) error {
	// Run a TTL check
	agent := c.Client.Agent()
	if !sealed {
		return agent.PassTTL(c.checkID(), "Vault Unsealed")
	} else {
		return agent.FailTTL(c.checkID(), "Vault Sealed")
	}
}

// fetchServiceTags returns all of the relevant tags for Consul.
func (c *serviceRegistration) fetchServiceTags(active, perfStandby, initialized bool) []string {
	activeTag := "standby"
	if active {
		activeTag = "active"
	}

	result := append(c.serviceTags, activeTag)

	if perfStandby {
		result = append(c.serviceTags, "performance-standby")
	}

	if initialized {
		result = append(result, "initialized")
	}

	return result
}

func (c *serviceRegistration) setRedirectAddr(addr string) (err error) {
	if addr == "" {
		return fmt.Errorf("redirect address must not be empty")
	}

	url, err := url.Parse(addr)
	if err != nil {
		return fmt.Errorf("failed to parse redirect URL %q: %w", addr, err)
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
			return fmt.Errorf("failed to find a host:port in redirect address %q: %w", url.Host, err)
		}
	}
	c.redirectPort, err = strconv.ParseInt(portStr, 10, 0)
	if err != nil || c.redirectPort < -1 || c.redirectPort > 65535 {
		return fmt.Errorf("failed to parse valid port %q: %w", portStr, err)
	}

	return nil
}

// durationMinusBuffer returns a duration, minus a buffer and jitter
// subtracted from the duration.  This function is used primarily for
// servicing Consul TTL Checks in advance of the TTL.
func durationMinusBuffer(intv time.Duration, buffer time.Duration, jitter int64) time.Duration {
	d := intv - buffer
	if jitter == 0 {
		d -= randomStagger(d)
	} else {
		d -= randomStagger(time.Duration(int64(d) / jitter))
	}
	return d
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

// randomStagger returns an interval between 0 and the duration
func randomStagger(intv time.Duration) time.Duration {
	if intv == 0 {
		return 0
	}
	return time.Duration(uint64(rand.Int63()) % uint64(intv))
}
