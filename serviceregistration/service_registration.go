package serviceregistration

import (
	"sync"

	log "github.com/hashicorp/go-hclog"
)

// Factory is the factory function to create a ServiceRegistration.
type Factory func(config map[string]string, logger log.Logger) (ServiceRegistration, error)

// ServiceRegistration is an interface that advertises the state of Vault to a
// service discovery network.
type ServiceRegistration interface {
	// NotifyActiveStateChange is used by Core to notify that this Vault
	// instance has changed its status to active or standby.
	NotifyActiveStateChange() error

	// NotifySealedStateChange is used by Core to notify that Vault has changed
	// its Sealed status to sealed or unsealed.
	NotifySealedStateChange() error

	// NotifyPerformanceStandbyStateChange is used by Core to notify that this
	// Vault instance has changed it status to performance standby or standby.
	NotifyPerformanceStandbyStateChange() error

	// Run executes any background service discovery tasks until the
	// shutdown channel is closed.
	RunServiceRegistration(
		waitGroup *sync.WaitGroup, shutdownCh ShutdownChannel, redirectAddr string,
		activeFunc ActiveFunction, sealedFunc SealedFunction, perfStandbyFunc PerformanceStandbyFunction) error
}

// Callback signatures for RunServiceRegistration
type ActiveFunction func() bool
type SealedFunction func() bool
type PerformanceStandbyFunction func() bool

// ShutdownChannel is the shutdown signal for RunServiceRegistration
type ShutdownChannel chan struct{}
