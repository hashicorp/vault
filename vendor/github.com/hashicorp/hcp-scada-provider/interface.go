// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"net"
	"time"
)

// SCADAProvider allows to expose services via SCADA capabilities.
type SCADAProvider interface {
	// UpdateMeta overwrites the internal map of meta-data values
	// and performs a re-handshake to update the remote broker.
	UpdateMeta(map[string]string)

	// AddMeta upserts keys and values in the internal map of meta-data
	// and performs a re-handshake to update the remote broker.
	AddMeta(...Meta)

	// DeleteMeta delete keys from the meta-date values and then perform a
	// re-handshake to update the remote broker.
	DeleteMeta(...string)

	// GetMeta returns the provider's current meta-data.
	GetMeta() map[string]string

	// Listen will expose the provided capability and make new connections
	// available through the returned listener. Closing the listener will stop
	// exposing the provided capability.
	//
	// The method will return an existing listener if the capability already existed.
	// Listeners can be retrieved even when the provider is stopped (e.g. before it is
	// started). New capabilities and new-meta data can be added at any time.
	//
	// The listener will only be closed, if it is closed explicitly by calling Close().
	// The listener will not be closed due to errors or when the provider is stopped.
	// The listener can hence be used after a restart of the provider.
	Listen(capability string) (net.Listener, error)

	// Start will register the provider on the SCADA broker and expose the
	// registered capabilities.
	Start() error

	// Stop will try to gracefully close the currently active SCADA session. This will
	// not close the capability listeners.
	Stop() error

	// SessionStatus returns the status of the SCADA connection.
	//
	// The possibles statuses are:
	//   - SessionStatusDisconnected: the provider is stopped
	//   - SessionStatusConnecting:   in the connect/handshake cycle
	//   - SessionStatusConnected:    connected and serving scada consumers
	//   - SessionStatusWaiting:      disconnected and waiting to retry a connection to the broker
	//
	// The full lifecycle is: connecting -> connected -> waiting -> connecting -> ... -> disconnected.
	SessionStatus() SessionStatus

	// LastError returns the last error recorded in the provider
	// connection state engine as well as the time at which the error occured.
	// That record is erased at each occasion when the provider achieves a new connection.
	//
	// A few common internal error will return a known type:
	//   - ErrProviderNotStarted: the provider is not started
	//   - ErrInvalidCredentials: could not obtain a token with the supplied credentials
	//   - ErrPermissionDenied:   principal does not have the permision to register as a provider
	//
	// Any other internal error will be returned directly and unchanged.
	LastError() (time.Time, error)

	// UpdateConfig overwrites the provider's configuration with the given configuration.
	UpdateConfig(config *Config) error
}

// SessionStatus is used to express the current status of the SCADA session.
type SessionStatus = string

const (
	// SessionStatusDisconnected is the state of the SCADA session if the
	// provider has not been started or has been stopped.
	SessionStatusDisconnected = SessionStatus("disconnected")

	// SessionStatusConnecting is the initial state of the SCADA connection
	// as well as the state it will be in if the connection got disrupted and
	// the library is trying to reconnect.
	//
	// The connection will transition to connected once the SCADA session is
	// established.
	SessionStatusConnecting = SessionStatus("connecting")

	// SessionStatusConnected is the state of the SCADA session if the
	// session is established and active.
	SessionStatusConnected = SessionStatus("connected")

	// SessionStatusRetrying is the state of a SCADA session that was
	// previous connected and is now in a wait-connect cycle
	SessionStatusWaiting = SessionStatus("waiting")
)

type Meta struct {
	Key, Value string
}
