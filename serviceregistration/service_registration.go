// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package serviceregistration

/*
ServiceRegistration is an interface that can be fulfilled to use
varying applications for service discovery, regardless of the physical
back-end used.

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
// The config is the key/value pairs set _inside_ the service registration config stanza.
// The state is the initial state.
// The redirectAddr is Vault core's RedirectAddr.
type Factory func(config map[string]string, logger log.Logger, state State) (ServiceRegistration, error)

// ServiceRegistration is an interface that advertises the state of Vault to a
// service discovery network.
type ServiceRegistration interface {
	// Run provides a shutdownCh, wait WaitGroup, and redirectAddr. The
	// shutdownCh is for monitoring when a shutdown occurs and initiating any
	// actions needed to leave service registration in a final state. When
	// finished, signalling that with wait means that Vault will wait until
	// complete. The redirectAddr is an optional parameter for implementations
	// that might need to communicate with Vault's listener via this address.
	//
	// Run is called just after Factory instantiation so can be relied upon
	// for controlling shutdown behavior.
	// Here is an example of its intended use:
	//	func Run(shutdownCh <-chan struct{}, wait sync.WaitGroup, redirectAddr string) error {
	//
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
	//			// Now do whatever we need to clean up.
	//			if err := someService.SetFinalState(); err != nil {
	//				// Log it at error level.
	//			}
	//		}()
	//		return nil
	//	}
	Run(shutdownCh <-chan struct{}, wait *sync.WaitGroup, redirectAddr string) error

	// NotifyActiveStateChange is used by Core to notify that this Vault
	// instance has changed its status on whether it's active or is
	// a standby.
	// If errors are returned, Vault only logs a warning, so it is
	// the implementation's responsibility to retry updating state
	// in the face of errors.
	NotifyActiveStateChange(isActive bool) error

	// NotifySealedStateChange is used by Core to notify that Vault has changed
	// its Sealed status to sealed or unsealed.
	// If errors are returned, Vault only logs a warning, so it is
	// the implementation's responsibility to retry updating state
	// in the face of errors.
	NotifySealedStateChange(isSealed bool) error

	// NotifyPerformanceStandbyStateChange is used by Core to notify that this
	// Vault instance has changed its performance standby status.
	// If errors are returned, Vault only logs a warning, so it is
	// the implementation's responsibility to retry updating state
	// in the face of errors.
	NotifyPerformanceStandbyStateChange(isStandby bool) error

	// NotifyInitializedStateChange is used by Core to notify that storage
	// has been initialized.  An unsealed core will always also be initialized.
	// If errors are returned, Vault only logs a warning, so it is
	// the implementation's responsibility to retry updating state
	// in the face of errors.
	NotifyInitializedStateChange(isInitialized bool) error

	// NotifyConfigurationReload is used by Core to notify that the Vault
	// configuration has been reloaded.
	// If errors are returned, Vault only logs a warning, so it is
	// the implementation's responsibility to retry updating state
	// in the face of errors.
	//
	// If the passed in conf is nil, it is assumed that the service registration
	// configuration no longer exits and should be deregistered.
	NotifyConfigurationReload(conf *map[string]string) error
}
