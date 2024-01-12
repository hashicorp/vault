// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"fmt"

	"github.com/hashicorp/vault/internal/observability/event"
)

// NewEvent should be used to create an audit event. The subtype field is needed
// for audit events. It will generate an ID if no ID is supplied. Supported
// options: WithID, WithNow.
func NewEvent(s subtype, opt ...Option) (*AuditEvent, error) {
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

	audit := &AuditEvent{
		ID:        opts.withID,
		Timestamp: opts.withNow,
		Version:   version,
		Subtype:   s,
	}

	if err := audit.validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return audit, nil
}

// validate attempts to ensure the audit event in its present state is valid.
func (a *AuditEvent) validate() error {
	const op = "audit.(AuditEvent).validate"

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

// String returns the string version of a format.
func (f format) String() string {
	return string(f)
}

// MetricTag returns a tag corresponding to this subtype to include in metrics.
func (st subtype) MetricTag() string {
	switch st {
	case RequestType:
		return "log_request"
	case ResponseType:
		return "log_response"
	}

	return ""
}
