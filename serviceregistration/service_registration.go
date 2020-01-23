package serviceregistration

/*
ServiceRegistration is an interface that can be fulfilled to use
varying applications for service discovery, regardless of the physical
back-end used. It uses [the observer pattern](https://refactoring.guru/design-patterns/observer).

Implementing the Factory and adding your factory to the list in commands.go
is essentially how you register your implementation as an observer. There is
no deregistration because service discovery stops on its own if Vault stops.

Service registration implements notifications for changes in _dynamic_
properties regarding Vault's health. Vault's version is the only static
property given in state for now, but if there's a need for more in the future,
we could add them on.
*/

import (
	"sync"

	log "github.com/hashicorp/go-hclog"
)

type State struct {
	VaultVersion                                            string
	IsInitialized, IsSealed, IsActive, IsPerformanceStandby bool
}

// Factory is the factory function to create a ServiceRegistration.
//
// The shutdownCh is the channel to watch for graceful shutdowns _of Vault_,
// and is great to use for creating background cleanup processes, or for stopping
// any ongoing goroutines.
//
// The config is the key/value pairs set _inside_ the service registration config stanza.
//
// The state is only the initial state. The pointer won't be updated over time. All
// state notifications will come through Notify methods.
//
// The redirectAddr is Vault core's RedirectAddr.
type Factory func(config map[string]string, logger log.Logger, state *State, redirectAddr string) (ServiceRegistration, error)

// ServiceRegistration is an interface that advertises the state of Vault to a
// service discovery network.
type ServiceRegistration interface {
	// Run provides a shutdownCh and wait WaitGroup. The shutdownCh
	// is for monitoring when a shutdown occurs and initiating any actions needed
	// to leave service registration in a final state. When finished, signalling
	// that with wait means that Vault will wait until complete.
	// Run is called just after Factory instantiation so can be relied upon
	// for controlling shutdown behavior.
	// Here is an example of its intended use:
	//	func Run(shutdownCh <-chan struct{}, wait sync.WaitGroup) error {
	//		// Since we are going to want Vault to wait to shutdown
	//		// until after we do cleanup...
	//		wait.Add(1)
	//
	//		// Run shutdown code in a goroutine so Run doesn't block.
	//		go func(){
	//			// Ensure that when this ends, no matter how it ends,
	//			// we don't cause Vault to hang on shutdown.
	//			defer wait.Done()
	//
	//			// Now wait until we're actually receiving a shutdown.
	//			<-shutdownCh
	//
	//			// Now do whatever we need to clean up. This is essentially
	//			// an OnStop method, and we may wish someday to replace Run
	//			// with OnStop to further simplify the interface.
	//			if err := someService.SetFinalState(); err != nil {
	//				// Log it at error level.
	//			}
	//		}()
	//		return nil
	//	}
	Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup) error

	// NotifyActiveStateChange is used by Core to notify that this Vault
	// instance has changed its status on whether it's active or is
	// a standby.
	NotifyActiveStateChange(isActive bool) error

	// NotifySealedStateChange is used by Core to notify that Vault has changed
	// its Sealed status to sealed or unsealed.
	NotifySealedStateChange(isSealed bool) error

	// NotifyPerformanceStandbyStateChange is used by Core to notify that this
	// Vault instance has changed its performance standby status.
	NotifyPerformanceStandbyStateChange(isStandby bool) error

	// NotifyInitializedStateChange is used by Core to notify that the core is
	// initialized.
	NotifyInitializedStateChange(isInitialized bool) error
}
