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
type Factory func(shutdownCh <-chan struct{}, config map[string]string, logger log.Logger, state *State, redirectAddr string) (ServiceRegistration, error)

// ServiceRegistration is an interface that advertises the state of Vault to a
// service discovery network.
type ServiceRegistration interface {
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
