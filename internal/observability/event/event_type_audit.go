// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package event

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

// Audit subtypes.
const (
	AuditRequest  auditSubtype = "AuditRequest"
	AuditResponse auditSubtype = "AuditResponse"
)

// Audit formats.
const (
	AuditFormatJSON  auditFormat = "json"
	AuditFormatJSONx auditFormat = "jsonx"
)

// auditVersion defines the version of audit events.
const auditVersion = "v0.1"

// auditSubtype defines the type of audit event.
type auditSubtype string

// auditFormat defines types of format audit events support.
type auditFormat string

// audit is the audit event.
type audit struct {
	ID             string            `json:"id"`
	Version        string            `json:"version"`
	Subtype        auditSubtype      `json:"subtype"` // the subtype of the audit event.
	Timestamp      time.Time         `json:"timestamp"`
	Data           *logical.LogInput `json:"data"`
	RequiredFormat auditFormat       `json:"format"`
}

// newAudit should be used to create an audit event.
// auditSubtype and auditFormat are needed for audit.
// It will use the supplied options, generate an ID if required, and validate the event.
func newAudit(opt ...Option) (*audit, error) {
	const op = "event.newAudit"

	opts, err := getOpts(opt...)
	if err != nil {
		return nil, fmt.Errorf("%s: error applying options: %w", op, err)
	}

	if opts.withID == "" {
		var err error

		opts.withID, err = NewID(string(AuditType))
		if err != nil {
			return nil, fmt.Errorf("%s: error creating ID for event: %w", op, err)
		}
	}

	audit := &audit{
		ID:             opts.withID,
		Version:        auditVersion,
		Subtype:        auditSubtype(opts.withSubtype),
		Timestamp:      opts.withNow,
		RequiredFormat: opts.withFormat,
	}

	if err := audit.validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return audit, nil
}

// validate attempts to ensure the event has the basic requirements of the event type configured.
func (a *audit) validate() error {
	const op = "event.(audit).validate"
	if a == nil {
		return fmt.Errorf("%s: audit is nil: %w", op, ErrInvalidParameter)
	}

	if a.ID == "" {
		return fmt.Errorf("%s: missing ID: %w", op, ErrInvalidParameter)
	}

	if a.Version != auditVersion {
		return fmt.Errorf("%s: audit version unsupported: %w", op, ErrInvalidParameter)
	}

	if a.Timestamp.IsZero() {
		return fmt.Errorf("%s: audit timestamp cannot be the zero time instant: %w", op, ErrInvalidParameter)
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

// validate ensures that auditSubtype is one of the set of allowed event subtypes.
func (t auditSubtype) validate() error {
	const op = "event.(auditSubtype).validate"
	switch t {
	case AuditRequest, AuditResponse:
		return nil
	default:
		return fmt.Errorf("%s: '%s' is not a valid event subtype: %w", op, t, ErrInvalidParameter)
	}
}

// validate ensures that auditFormat is one of the set of allowed event formats.
func (f auditFormat) validate() error {
	const op = "event.(auditFormat).validate"
	switch f {
	case AuditFormatJSON, AuditFormatJSONx:
		return nil
	default:
		return fmt.Errorf("%s: '%s' is not a valid format: %w", op, f, ErrInvalidParameter)
	}
}

// String returns the string version of an auditFormat.
func (f auditFormat) String() string {
	return string(f)
}
