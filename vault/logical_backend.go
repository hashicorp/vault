package vault

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrUnsupportedOperation is returned if the operation is not supported
	// by the logical backend.
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrUnsupportedPath is returned if the path is not supported
	// by the logical backend.
	ErrUnsupportedPath = errors.New("unsupported path")
)

// LogicalBackend interface must be implemented to be "mountable" at
// a given path. Requests flow through a router which has various mount
// points that flow to a logical backend. The logic of each backend is flexible,
// and this is what allows materialized keys to function. There can be specialized
// logical backends for various upstreams (Consul, PostgreSQL, MySQL, etc) that can
// interact with remote APIs to generate keys dynamically. This interface also
// allows for a "procfs" like interaction, as internal state can be exposed by
// acting like a logical backend and being mounted.
type LogicalBackend interface {
	// HandleRequest is used to handle a request and generate a response.
	// The backends must check the operation type and handle appropriately.
	HandleRequest(*Request) (*Response, error)

	// RootPaths is a list of paths that require root level privileges.
	// These paths will be enforced by the router so that backends do
	// not need to handle the authorization. Paths are enforced exactly
	// or using a prefix match if they end in '*'
	RootPaths() []string
}

// Operation is an enum that is used to specify the type
// of request being made
type Operation string

const (
	ReadOperation   Operation = "read"
	WriteOperation            = "write"
	DeleteOperation           = "delete"
	ListOperation             = "list"
	RevokeOperation           = "revoke"
	HelpOperation             = "help"
)

// Request is a struct that stores the parameters and context
// of a request being made to Vault. It is used to abstract
// the details of the higher level request protocol from the handlers.
type Request struct {
	// Operation is the requested operation type
	Operation Operation

	// Path is the part of the request path not consumed by the
	// routing. As an example, if the original request path is "prod/aws/foo"
	// and the AWS logical backend is mounted at "prod/aws/", then the
	// final path is "foo" since the mount prefix is trimmed.
	Path string

	// Request data is an opaque map that must have string keys.
	Data map[string]interface{}

	// View is the storage view of this logical backend. It can be used
	// to durably store and retrieve state from the backend.
	View *BarrierView
}

// Response is a struct that stores the response of a request.
// It is used to abstract the details of the higher level request protocol.
type Response struct {
	// IsSecret is used to indicate this is secret material instead of policy or configuration.
	// Non-secrets never have a VaultID or renewable properties.
	IsSecret bool

	// The lease settings if applicable.
	Lease *Lease

	// Response data is an opaque map that must have string keys.
	Data map[string]interface{}
}

// Lease is used to provide more information about the lease
type Lease struct {
	VaultID      string        // VaultID is the unique identifier used for renewal and revocation
	Renewable    bool          // Is the VaultID renewable
	Revokable    bool          // Is the secret revokable. Must support 'Revoke' operation.
	Duration     time.Duration // Current lease duration
	MaxDuration  time.Duration // Maximum lease duration
	MaxIncrement time.Duration // Maximum increment to lease duration
}

// Factory is the factory function to create a logical backend.
type Factory func(map[string]string) (LogicalBackend, error)

// BuiltinBackends contains all of the available backends
var BuiltinBackends = map[string]Factory{
	"generic": newGenericBackend,
}

// NewBackend returns a new logical Backend with the given type and configuration.
// The backend is looked up in the BuiltinBackends variable.
func NewBackend(t string, conf map[string]string) (LogicalBackend, error) {
	f, ok := BuiltinBackends[t]
	if !ok {
		return nil, fmt.Errorf("unknown logical backend type: %s", t)
	}
	return f(conf)
}
