// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package audit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/internal/observability/event"

	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

// newEvent should be used to create an audit event.
// subtype and format are needed for audit.
// It will use the supplied Options, generate an ID if required, and validate the event.
func newEvent(opt ...Option) (*auditEvent, error) {
	const op = "audit.newEvent"

	// Get the default options
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("%s: error applying options: %w", op, err)
	}

	if opts.withID == "" {
		var err error

		opts.withID, err = event.NewID(string(event.AuditType))
		if err != nil {
			return nil, fmt.Errorf("%s: error creating ID for event: %w", op, err)
		}
	}

	audit := &auditEvent{
		ID:             opts.withID,
		Timestamp:      opts.withNow,
		Version:        version,
		Subtype:        opts.withSubtype,
		RequiredFormat: opts.withFormat,
	}

	if err := audit.validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return audit, nil
}

// validate attempts to ensure the event has the basic requirements of the event type configured.
func (a *auditEvent) validate() error {
	const op = "audit.(auditEvent).validate"

	if a == nil {
		return fmt.Errorf("%s: event is nil: %w", op, event.ErrInvalidParameter)
	}

	if a.ID == "" {
		return fmt.Errorf("%s: missing ID: %w", op, event.ErrInvalidParameter)
	}

	if a.Version != version {
		return fmt.Errorf("%s: event version unsupported: %w", op, event.ErrInvalidParameter)
	}

	if a.Timestamp.IsZero() {
		return fmt.Errorf("%s: event timestamp cannot be the zero time instant: %w", op, event.ErrInvalidParameter)
	}

	err := a.Subtype.validate()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.RequiredFormat.validate()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// validate ensures that subtype is one of the set of allowed event subtypes.
func (t subtype) validate() error {
	const op = "audit.(subtype).validate"
	switch t {
	case RequestType, ResponseType:
		return nil
	default:
		return fmt.Errorf("%s: '%s' is not a valid event subtype: %w", op, t, event.ErrInvalidParameter)
	}
}

// validate ensures that format is one of the set of allowed event formats.
func (f format) validate() error {
	const op = "audit.(format).validate"
	switch f {
	case JSONFormat, JSONxFormat:
		return nil
	default:
		return fmt.Errorf("%s: '%s' is not a valid format: %w", op, f, event.ErrInvalidParameter)
	}
}

// String returns the string version of an format.
func (f format) String() string {
	return string(f)
}

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
