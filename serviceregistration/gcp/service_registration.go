package gcp

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	sr "github.com/hashicorp/vault/serviceregistration"
	atomicB "go.uber.org/atomic"
	"google.golang.org/api/dns/v1"
)

const (
	// Environment variable for google credentials
	EnvVarGoogleApplicationCredentials = "GOOGLE_APPLICATION_CREDENTIALS"

	// jitterFactor specifies the jitter factor used to stagger checks
	jitterFactor = 16

	retryAttempts = 30

	// retryInterval specifies the retry duration to use when an API call to the GCP fails.
	retryInterval = 3 * time.Second

	// reconcileTimeout is how often Vault should query GCP to detect and fix any state drift.
	reconcileTimeout = 60 * time.Second
)

type serviceRegistration struct {
	logger              hclog.Logger
	client              *dns.Service
	credentials         string
	project             string
	managedZone         string
	zoneDNSSuffix       string
	activeDNSName       string
	standbyDNSName      string
	allUnsealedDNSName  string
	redirectHost        string
	redirectPort        int64
	serviceLock         sync.RWMutex
	notifyActiveCh      chan struct{}
	notifySealedCh      chan struct{}
	notifyPerfStandbyCh chan struct{}
	notifyInitializedCh chan struct{}
	isActive            *atomicB.Bool
	isSealed            *atomicB.Bool
	isPerfStandby       *atomicB.Bool
	isInitialized       *atomicB.Bool
}

func NewServiceRegistration(config map[string]string, logger hclog.Logger, state sr.State) (sr.ServiceRegistration, error) {
	// Parse and validate config
	credentials, err := getRequiredField(logger, config, EnvVarGoogleApplicationCredentials, "credentials")
	if err != nil {
		return nil, err
	}
	project, err := getRequiredField(logger, config, "", "zone_project")
	if err != nil {
		return nil, err
	}
	zoneName, err := getRequiredField(logger, config, "", "zone_name")
	if err != nil {
		return nil, err
	}
	zoneDNSSuffix, err := getRequiredField(logger, config, "", "zone_dns_suffix")
	if err != nil {
		return nil, err
	}

	// Create the client to talk to the cloud DNS service
	client, err := dns.NewService(context.Background())
	if err != nil {
		return nil, err
	}

	return &serviceRegistration{
		logger: logger,

		client:             client,
		credentials:        credentials,
		project:            project,
		managedZone:        zoneName,
		zoneDNSSuffix:      zoneDNSSuffix,
		activeDNSName:      fmt.Sprintf("active.vault.%s", zoneDNSSuffix),
		standbyDNSName:     fmt.Sprintf("standby.vault.%s", zoneDNSSuffix),
		allUnsealedDNSName: fmt.Sprintf("vault.%s", zoneDNSSuffix),

		notifyActiveCh:      make(chan struct{}),
		notifySealedCh:      make(chan struct{}),
		notifyPerfStandbyCh: make(chan struct{}),
		notifyInitializedCh: make(chan struct{}),

		isActive:      atomicB.NewBool(state.IsActive),
		isSealed:      atomicB.NewBool(state.IsSealed),
		isPerfStandby: atomicB.NewBool(state.IsPerformanceStandby),
		isInitialized: atomicB.NewBool(state.IsInitialized),
	}, nil
}

func (r *serviceRegistration) Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup, redirectAddr string) error {
	if err := r.setRedirectAddr(redirectAddr); err != nil {
		return err
	}

	// Since we are going to want Vault to wait to shutdown until after we do cleanup
	wait.Add(1)

	// Run shutdown code in a goroutine so Run doesn't block.
	go func() {
		r.runEventDemuxer(wait, shutdownCh)
	}()

	return nil
}

func (r *serviceRegistration) runEventDemuxer(wait *sync.WaitGroup, shutdownCh <-chan struct{}) {
	// This defer statement should be executed last. So push it first.
	defer wait.Done()

	// Fire the reconcileTimer immediately upon starting the event demuxer
	reconcileTimer := time.NewTimer(0)
	defer reconcileTimer.Stop()

	var shutdown bool
	serviceRegLock := new(int32)

	for !shutdown {
		select {
		case <-r.notifyActiveCh:
			// Run reconcile immediately upon active state change notification
			reconcileTimer.Reset(0)
		case <-r.notifySealedCh:
			// Run check timer immediately upon a seal state change notification
			reconcileTimer.Reset(0)
		case <-r.notifyPerfStandbyCh:
			// Run check timer immediately upon a perfstandby state change notification
			reconcileTimer.Reset(0)
		case <-r.notifyInitializedCh:
			// Run check timer immediately upon an initialized state change notification
			reconcileTimer.Reset(0)
		case <-reconcileTimer.C:
			// Unconditionally rearm the reconcileTimer
			reconcileTimer.Reset(reconcileTimeout - randomStagger(reconcileTimeout/jitterFactor))

			// Abort if reconcile handler is already active
			if atomic.CompareAndSwapInt32(serviceRegLock, 0, 1) {
				go func() {
					defer atomic.CompareAndSwapInt32(serviceRegLock, 1, 0)
					attempts := 0
					for !shutdown && attempts < retryAttempts {
						if err := r.reconcileCloudDNS(); err != nil {
							if r.logger.IsWarn() {
								r.logger.Warn("reconcile unable to talk with GCP cloud DNS", "error", err)
							}
							time.Sleep(retryInterval)
							attempts++
							continue
						}

						r.serviceLock.Lock()
						defer r.serviceLock.Unlock()

						return
					}
				}()
			}
		case <-shutdownCh:
			r.logger.Info("shutting down")
			shutdown = true
		}
	}

	r.serviceLock.RLock()
	defer r.serviceLock.RUnlock()

	// TODO: Remove the service record from GCP
	// if err := r.Client.Agent().ServiceDeregister(registeredServiceID); err != nil {
	// 	r.logger.Warn("service deregistration failed", "error", err)
	// }
}

func (r *serviceRegistration) reconcileCloudDNS() error {
	if !r.isInitialized.Load() {
		return nil
	}

	// Ensure that the managed zone exists
	managedZone, err := r.client.ManagedZones.Get(r.project, r.managedZone).Do()
	if err != nil || managedZone == nil {
		return fmt.Errorf("managed zone must exist: %w", err)
	}

	// Get the current record set for the managed zone
	recordSet, err := r.client.ResourceRecordSets.List(r.project, r.managedZone).Do()
	if err != nil {
		return err
	}

	// Get the current active and standby addresses
	actives := make(map[string]bool)
	standbys := make(map[string]bool)
	for _, rs := range recordSet.Rrsets {
		switch rs.Name {
		case r.activeDNSName:
			for _, d := range rs.Rrdatas {
				actives[d] = true
			}
		case r.standbyDNSName:
			for _, d := range rs.Rrdatas {
				standbys[d] = true
			}
		}
	}

	change := &dns.Change{}

	// If this instance is active and unsealed, all current active records
	// must be deleted and this instance must become the active.
	if r.isActive.Load() && !r.isSealed.Load() {
		r.logger.Info("ACTIVE: Current standbys", "standbys", standbys)
		r.logger.Info("ACTIVE: Current actives", "actives", actives)

		change.Additions = []*dns.ResourceRecordSet{
			{
				Name:    r.activeDNSName,
				Type:    "A",
				Ttl:     5,
				Rrdatas: []string{r.redirectHost},
			},
		}
		if len(actives) > 0 {
			change.Deletions = []*dns.ResourceRecordSet{
				{
					Name:    r.activeDNSName,
					Type:    "A",
					Ttl:     5,
					Rrdatas: mapKeysToSlice(actives),
				},
			}
		}
		if _, err := r.client.Changes.Create(r.project, r.managedZone, change).Do(); err != nil {
			return err
		}

		// If this instance was a standby, remove it from standbys
		if _, ok := standbys[r.redirectHost]; ok {
			delete(standbys, r.redirectHost)
			change := &dns.Change{}

			if len(standbys) > 0 {
				change.Additions = []*dns.ResourceRecordSet{
					{
						Name:    r.standbyDNSName,
						Type:    "A",
						Ttl:     5,
						Rrdatas: mapKeysToSlice(standbys),
					},
				}
			}
			change.Deletions = []*dns.ResourceRecordSet{
				{
					Name:    r.standbyDNSName,
					Type:    "A",
					Ttl:     5,
					Rrdatas: []string{r.redirectHost},
				},
			}
			if _, err := r.client.Changes.Create(r.project, r.managedZone, change).Do(); err != nil {
				return err
			}
		}

		return nil
	}

	// Otherwise, it is a standby. We need to add it in addition to all others that currently exist.
	if !r.isActive.Load() && !r.isSealed.Load() {
		r.logger.Info("STANDBY: Current standbys", "standbys", standbys)
		r.logger.Info("STANDBY: Current actives", "actives", actives)

		// If the instance is not a standby, then add and delete. Otherwise, there is nothing to do.
		if _, ok := standbys[r.redirectHost]; !ok {
			// attempt to delete the active from the standby
			for k := range actives {
				delete(standbys, k)
			}

			change.Additions = []*dns.ResourceRecordSet{
				{
					Name:    r.standbyDNSName,
					Type:    "A",
					Ttl:     5,
					Rrdatas: append(mapKeysToSlice(standbys), r.redirectHost),
				},
			}

			if len(standbys) > 0 {
				change.Deletions = []*dns.ResourceRecordSet{
					{
						Name:    r.standbyDNSName,
						Type:    "A",
						Ttl:     5,
						Rrdatas: mapKeysToSlice(standbys),
					},
				}
			}

			if _, err := r.client.Changes.Create(r.project, r.managedZone, change).Do(); err != nil {
				return err
			}
		}

		return nil
	}

	return nil
}

func mapKeysToSlice(m map[string]bool) []string {
	s := make([]string, 0)
	for k := range m {
		s = append(s, k)
	}
	return s
}

// randomStagger returns an interval between 0 and the duration
func randomStagger(intv time.Duration) time.Duration {
	if intv == 0 {
		return 0
	}
	return time.Duration(uint64(rand.Int63()) % uint64(intv))
}

func (r *serviceRegistration) NotifyActiveStateChange(isActive bool) error {
	r.isActive.Store(isActive)
	select {
	case r.notifyActiveCh <- struct{}{}:
	default:
		r.logger.Warn("concurrent state change notify dropped")
	}

	return nil
}

func (r *serviceRegistration) NotifyPerformanceStandbyStateChange(isStandby bool) error {
	r.isPerfStandby.Store(isStandby)
	select {
	case r.notifyPerfStandbyCh <- struct{}{}:
	default:
		r.logger.Warn("concurrent state change notify dropped")
	}

	return nil
}

func (r *serviceRegistration) NotifySealedStateChange(isSealed bool) error {
	r.isSealed.Store(isSealed)
	select {
	case r.notifySealedCh <- struct{}{}:
	default:
		r.logger.Warn("concurrent sealed state change notify dropped")
	}

	return nil
}

func (r *serviceRegistration) NotifyInitializedStateChange(isInitialized bool) error {
	r.isInitialized.Store(isInitialized)
	select {
	case r.notifyInitializedCh <- struct{}{}:
	default:
		r.logger.Warn("concurrent initalize state change notify dropped")
	}

	return nil
}

func (r *serviceRegistration) setRedirectAddr(addr string) (err error) {
	if addr == "" {
		return fmt.Errorf("redirect address must not be empty")
	}

	url, err := url.Parse(addr)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("failed to parse redirect URL %q: {{err}}", addr), err)
	}

	var portStr string
	r.redirectHost, portStr, err = net.SplitHostPort(url.Host)
	if err != nil {
		if url.Scheme == "http" {
			portStr = "80"
		} else if url.Scheme == "https" {
			portStr = "443"
		} else if url.Scheme == "unix" {
			portStr = "-1"
			r.redirectHost = url.Path
		} else {
			return errwrap.Wrapf(fmt.Sprintf(`failed to find a host:port in redirect address "%v": {{err}}`, url.Host), err)
		}
	}
	r.redirectPort, err = strconv.ParseInt(portStr, 10, 0)
	if err != nil || r.redirectPort < -1 || r.redirectPort > 65535 {
		return errwrap.Wrapf(fmt.Sprintf(`failed to parse valid port "%v": {{err}}`, portStr), err)
	}

	return nil
}

func getRequiredField(logger hclog.Logger, config map[string]string, envVar, configParam string) (string, error) {
	value := ""
	switch {
	case os.Getenv(envVar) != "":
		value = os.Getenv(envVar)
	case config[configParam] != "":
		value = config[configParam]
	default:
		return "", fmt.Errorf(`%s must be provided via %q or the %q config parameter`, configParam, envVar, configParam)
	}
	if logger.IsDebug() {
		logger.Debug(fmt.Sprintf("%q: %q", configParam, value))
	}
	return value, nil
}
