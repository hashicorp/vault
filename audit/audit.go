package audit

import (
	"context"

	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

// Backend interface must be implemented for an audit
// mechanism to be made available. Audit backends can be enabled to
// sink information to different backends such as logs, file, databases,
// or other external services.
type Backend interface {
	// LogRequest is used to synchronously log a request. This is done after the
	// request is authorized but before the request is executed. The arguments
	// MUST not be modified in anyway. They should be deep copied if this is
	// a possibility.
	LogRequest(context.Context, *logical.LogInput) error

	// LogResponse is used to synchronously log a response. This is done after
	// the request is processed but before the response is sent. The arguments
	// MUST not be modified in anyway. They should be deep copied if this is
	// a possibility.
	LogResponse(context.Context, *logical.LogInput) error

	// LogTestMessage is used to check an audit backend before adding it
	// permanently. It should attempt to synchronously log the given test
	// message, WITHOUT using the normal Salt (which would require a storage
	// operation on creation, which is currently disallowed.)
	LogTestMessage(context.Context, *logical.LogInput, map[string]string) error

	// GetHash is used to return the given data with the backend's hash,
	// so that a caller can determine if a value in the audit log matches
	// an expected plaintext value
	GetHash(context.Context, string) (string, error)

	// Reload is called on SIGHUP for supporting backends.
	Reload(context.Context) error

	// Invalidate is called for path invalidation
	Invalidate(context.Context)
}

// BackendConfig contains configuration parameters used in the factory func to
// instantiate audit backends
type BackendConfig struct {
	// The view to store the salt
	SaltView logical.Storage

	// The salt config that should be used for any secret obfuscation
	SaltConfig *salt.Config

	// Config is the opaque user configuration provided when mounting
	Config map[string]string
}

// Factory is the factory function to create an audit backend.
type Factory func(context.Context, *BackendConfig) (Backend, error)
