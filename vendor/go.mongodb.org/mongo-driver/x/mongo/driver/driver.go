// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

// Package driver is intended for internal use only. It is made available to
// facilitate use cases that require access to internal MongoDB driver
// functionality and state. The API of this package is not stable and there is
// no backward compatibility guarantee.
//
// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package driver // import "go.mongodb.org/mongo-driver/x/mongo/driver"

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// AuthConfig holds the information necessary to perform an authentication attempt.
// this was moved from the auth package to avoid a circular dependency. The auth package
// reexports this under the old name to avoid breaking the public api.
type AuthConfig struct {
	Description   description.Server
	Connection    Connection
	ClusterClock  *session.ClusterClock
	HandshakeInfo HandshakeInformation
	ServerAPI     *ServerAPIOptions
}

// OIDCCallback is the type for both Human and Machine Callback flows. RefreshToken will always be
// nil in the OIDCArgs for the Machine flow.
type OIDCCallback func(context.Context, *OIDCArgs) (*OIDCCredential, error)

// OIDCArgs contains the arguments for the OIDC callback.
type OIDCArgs struct {
	Version      int
	IDPInfo      *IDPInfo
	RefreshToken *string
}

// OIDCCredential contains the access token and refresh token.
type OIDCCredential struct {
	AccessToken  string
	ExpiresAt    *time.Time
	RefreshToken *string
}

// IDPInfo contains the information needed to perform OIDC authentication with an Identity Provider.
type IDPInfo struct {
	Issuer        string   `bson:"issuer"`
	ClientID      string   `bson:"clientId"`
	RequestScopes []string `bson:"requestScopes"`
}

// Authenticator handles authenticating a connection. The implementers of this interface
// are all in the auth package. Most authentication mechanisms do not allow for Reauth,
// but this is included in the interface so that whenever a new mechanism is added, it
// must be explicitly considered.
type Authenticator interface {
	// Auth authenticates the connection.
	Auth(context.Context, *AuthConfig) error
	Reauth(context.Context, *AuthConfig) error
}

// Cred is a user's credential.
type Cred struct {
	Source              string
	Username            string
	Password            string
	PasswordSet         bool
	Props               map[string]string
	OIDCMachineCallback OIDCCallback
	OIDCHumanCallback   OIDCCallback
}

// Deployment is implemented by types that can select a server from a deployment.
type Deployment interface {
	SelectServer(context.Context, description.ServerSelector) (Server, error)
	Kind() description.TopologyKind
}

// Connector represents a type that can connect to a server.
type Connector interface {
	Connect() error
}

// Disconnector represents a type that can disconnect from a server.
type Disconnector interface {
	Disconnect(context.Context) error
}

// Subscription represents a subscription to topology updates. A subscriber can receive updates through the
// Updates field.
type Subscription struct {
	Updates <-chan description.Topology
	ID      uint64
}

// Subscriber represents a type to which another type can subscribe. A subscription contains a channel that
// is updated with topology descriptions.
type Subscriber interface {
	Subscribe() (*Subscription, error)
	Unsubscribe(*Subscription) error
}

// Server represents a MongoDB server. Implementations should pool connections and handle the
// retrieving and returning of connections.
type Server interface {
	Connection(context.Context) (Connection, error)

	// RTTMonitor returns the round-trip time monitor associated with this server.
	RTTMonitor() RTTMonitor
}

// Connection represents a connection to a MongoDB server.
type Connection interface {
	WriteWireMessage(context.Context, []byte) error
	ReadWireMessage(ctx context.Context) ([]byte, error)
	Description() description.Server

	// Close closes any underlying connection and returns or frees any resources held by the
	// connection. Close is idempotent and can be called multiple times, although subsequent calls
	// to Close may return an error. A connection cannot be used after it is closed.
	Close() error

	ID() string
	ServerConnectionID() *int64
	DriverConnectionID() uint64 // TODO(GODRIVER-2824): change type to int64.
	Address() address.Address
	Stale() bool
	OIDCTokenGenID() uint64
	SetOIDCTokenGenID(uint64)
}

// RTTMonitor represents a round-trip-time monitor.
type RTTMonitor interface {
	// EWMA returns the exponentially weighted moving average observed round-trip time.
	EWMA() time.Duration

	// Min returns the minimum observed round-trip time over the window period.
	Min() time.Duration

	// P90 returns the 90th percentile observed round-trip time over the window period.
	P90() time.Duration

	// Stats returns stringified stats of the current state of the monitor.
	Stats() string
}

var _ RTTMonitor = &csot.ZeroRTTMonitor{}

// PinnedConnection represents a Connection that can be pinned by one or more cursors or transactions. Implementations
// of this interface should maintain the following invariants:
//
// 1. Each Pin* call should increment the number of references for the connection.
// 2. Each Unpin* call should decrement the number of references for the connection.
// 3. Calls to Close() should be ignored until all resources have unpinned the connection.
type PinnedConnection interface {
	Connection
	PinToCursor() error
	PinToTransaction() error
	UnpinFromCursor() error
	UnpinFromTransaction() error
}

// The session.LoadBalancedTransactionConnection type is a copy of PinnedConnection that was introduced to avoid
// import cycles. This compile-time assertion ensures that these types remain in sync if the PinnedConnection interface
// is changed in the future.
var _ PinnedConnection = (session.LoadBalancedTransactionConnection)(nil)

// LocalAddresser is a type that is able to supply its local address
type LocalAddresser interface {
	LocalAddress() address.Address
}

// Expirable represents an expirable object.
type Expirable interface {
	Expire() error
	Alive() bool
}

// StreamerConnection represents a Connection that supports streaming wire protocol messages using the moreToCome and
// exhaustAllowed flags.
//
// The SetStreaming and CurrentlyStreaming functions correspond to the moreToCome flag on server responses. If a
// response has moreToCome set, SetStreaming(true) will be called and CurrentlyStreaming() should return true.
//
// CanStream corresponds to the exhaustAllowed flag. The operations layer will set exhaustAllowed on outgoing wire
// messages to inform the server that the driver supports streaming.
type StreamerConnection interface {
	Connection
	SetStreaming(bool)
	CurrentlyStreaming() bool
	SupportsStreaming() bool
}

// Compressor is an interface used to compress wire messages. If a Connection supports compression
// it should implement this interface as well. The CompressWireMessage method will be called during
// the execution of an operation if the wire message is allowed to be compressed.
type Compressor interface {
	CompressWireMessage(src, dst []byte) ([]byte, error)
}

// ProcessErrorResult represents the result of a ErrorProcessor.ProcessError() call. Exact values for this type can be
// checked directly (e.g. res == ServerMarkedUnknown), but it is recommended that applications use the ServerChanged()
// function instead.
type ProcessErrorResult int

const (
	// NoChange indicates that the error did not affect the state of the server.
	NoChange ProcessErrorResult = iota
	// ServerMarkedUnknown indicates that the error only resulted in the server being marked as Unknown.
	ServerMarkedUnknown
	// ConnectionPoolCleared indicates that the error resulted in the server being marked as Unknown and its connection
	// pool being cleared.
	ConnectionPoolCleared
)

// ErrorProcessor implementations can handle processing errors, which may modify their internal state.
// If this type is implemented by a Server, then Operation.Execute will call it's ProcessError
// method after it decodes a wire message.
type ErrorProcessor interface {
	ProcessError(err error, conn Connection) ProcessErrorResult
}

// HandshakeInformation contains information extracted from a MongoDB connection handshake. This is a helper type that
// augments description.Server by also tracking server connection ID and authentication-related fields. We use this type
// rather than adding authentication-related fields to description.Server to avoid retaining sensitive information in a
// user-facing type. The server connection ID is stored in this type because unlike description.Server, all handshakes are
// correlated with a single network connection.
type HandshakeInformation struct {
	Description             description.Server
	SpeculativeAuthenticate bsoncore.Document
	ServerConnectionID      *int64
	SaslSupportedMechs      []string
}

// Handshaker is the interface implemented by types that can perform a MongoDB
// handshake over a provided driver.Connection. This is used during connection
// initialization. Implementations must be goroutine safe.
type Handshaker interface {
	GetHandshakeInformation(context.Context, address.Address, Connection) (HandshakeInformation, error)
	FinishHandshake(context.Context, Connection) error
}

// SingleServerDeployment is an implementation of Deployment that always returns a single server.
type SingleServerDeployment struct{ Server }

var _ Deployment = SingleServerDeployment{}

// SelectServer implements the Deployment interface. This method does not use the
// description.SelectedServer provided and instead returns the embedded Server.
func (ssd SingleServerDeployment) SelectServer(context.Context, description.ServerSelector) (Server, error) {
	return ssd.Server, nil
}

// Kind implements the Deployment interface. It always returns description.Single.
func (SingleServerDeployment) Kind() description.TopologyKind { return description.Single }

// SingleConnectionDeployment is an implementation of Deployment that always returns the same Connection. This
// implementation should only be used for connection handshakes and server heartbeats as it does not implement
// ErrorProcessor, which is necessary for application operations.
type SingleConnectionDeployment struct{ C Connection }

var _ Deployment = SingleConnectionDeployment{}
var _ Server = SingleConnectionDeployment{}

// SelectServer implements the Deployment interface. This method does not use the
// description.SelectedServer provided and instead returns itself. The Connections returned from the
// Connection method have a no-op Close method.
func (scd SingleConnectionDeployment) SelectServer(context.Context, description.ServerSelector) (Server, error) {
	return scd, nil
}

// Kind implements the Deployment interface. It always returns description.Single.
func (SingleConnectionDeployment) Kind() description.TopologyKind { return description.Single }

// Connection implements the Server interface. It always returns the embedded connection.
func (scd SingleConnectionDeployment) Connection(context.Context) (Connection, error) {
	return scd.C, nil
}

// RTTMonitor implements the driver.Server interface.
func (scd SingleConnectionDeployment) RTTMonitor() RTTMonitor {
	return &csot.ZeroRTTMonitor{}
}

// TODO(GODRIVER-617): We can likely use 1 type for both the Type and the RetryMode by using 2 bits for the mode and 1
// TODO bit for the type. Although in the practical sense, we might not want to do that since the type of retryability
// TODO is tied to the operation itself and isn't going change, e.g. and insert operation will always be a write,
// TODO however some operations are both reads and  writes, for instance aggregate is a read but with a $out parameter
// TODO it's a write.

// Type specifies whether an operation is a read, write, or unknown.
type Type uint

// THese are the availables types of Type.
const (
	_ Type = iota
	Write
	Read
)

// RetryMode specifies the way that retries are handled for retryable operations.
type RetryMode uint

// These are the modes available for retrying. Note that if Timeout is specified on the Client, the
// operation will automatically retry as many times as possible within the context's deadline
// unless RetryNone is used.
const (
	// RetryNone disables retrying.
	RetryNone RetryMode = iota
	// RetryOnce will enable retrying the entire operation once if Timeout is not specified.
	RetryOnce
	// RetryOncePerCommand will enable retrying each command associated with an operation if Timeout
	// is not specified. For example, if an insert is batch split into 4 commands then each of
	// those commands is eligible for one retry.
	RetryOncePerCommand
	// RetryContext will enable retrying until the context.Context's deadline is exceeded or it is
	// cancelled.
	RetryContext
)

// Enabled returns if this RetryMode enables retrying.
func (rm RetryMode) Enabled() bool {
	return rm == RetryOnce || rm == RetryOncePerCommand || rm == RetryContext
}
