// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/vault/internal/observability/event"
	"github.com/hashicorp/vault/sdk/logical"
)

// version defines the version of audit events.
const version = "v0.1"

// Audit subtypes.
const (
	RequestType  subtype = "AuditRequest"
	ResponseType subtype = "AuditResponse"
)

// Audit formats.
const (
	jsonFormat  format = "json"
	jsonxFormat format = "jsonx"
)

// Check AuditEvent implements the timeProvider at compile time.
var _ timeProvider = (*Event)(nil)

// Event is the audit event.
type Event struct {
	ID        string            `json:"id"`
	Version   string            `json:"version"`
	Subtype   subtype           `json:"subtype"` // the subtype of the audit event.
	Timestamp time.Time         `json:"timestamp"`
	Data      *logical.LogInput `json:"data"`
	prov      timeProvider
}

// setTimeProvider can be used to set a specific time provider which is used when
// creating an entry.
// NOTE: This is primarily used for testing to supply a known time value.
func (a *Event) setTimeProvider(t timeProvider) {
	a.prov = t
}

// timeProvider returns a configured time provider, or the default if not set.
func (a *Event) timeProvider() timeProvider {
	if a.prov == nil {
		return a
	}

	return a.prov
}

// format defines types of format audit events support.
type format string

// subtype defines the type of audit event.
type subtype string

// newEvent should be used to create an audit event. The subtype field is needed
// for audit events. It will generate an ID if no ID is supplied. Supported
// options: withID, withNow.
func newEvent(s subtype, opt ...option) (*Event, error) {
	// Get the default options
	opts, err := getOpts(opt...)
	if err != nil {
		return nil, err
	}

	if opts.withID == "" {
		var err error

		opts.withID, err = event.NewID(string(event.AuditType))
		if err != nil {
			return nil, fmt.Errorf("error creating ID for event: %w", err)
		}
	}

	audit := &Event{
		ID:        opts.withID,
		Timestamp: opts.withNow,
		Version:   version,
		Subtype:   s,
	}

	if err := audit.validate(); err != nil {
		return nil, err
	}
	return audit, nil
}

// validate attempts to ensure the audit event in its present state is valid.
func (a *Event) validate() error {
	if a == nil {
		return fmt.Errorf("event is nil: %w", ErrInvalidParameter)
	}

	if a.ID == "" {
		return fmt.Errorf("missing ID: %w", ErrInvalidParameter)
	}

	if a.Version != version {
		return fmt.Errorf("event version unsupported: %w", ErrInvalidParameter)
	}

	if a.Timestamp.IsZero() {
		return fmt.Errorf("event timestamp cannot be the zero time instant: %w", ErrInvalidParameter)
	}

	err := a.Subtype.validate()
	if err != nil {
		return err
	}

	return nil
}

// validate ensures that subtype is one of the set of allowed event subtypes.
func (t subtype) validate() error {
	switch t {
	case RequestType, ResponseType:
		return nil
	default:
		return fmt.Errorf("invalid event subtype %q: %w", t, ErrInvalidParameter)
	}
}

// validate ensures that format is one of the set of allowed event formats.
func (f format) validate() error {
	switch f {
	case jsonFormat, jsonxFormat:
		return nil
	default:
		return fmt.Errorf("invalid format %q: %w", f, ErrInvalidParameter)
	}
}

// String returns the string version of a format.
func (f format) String() string {
	return string(f)
}

// MetricTag returns a tag corresponding to this subtype to include in metrics.
// If a tag cannot be found the value is returned 'as-is' in string format.
func (t subtype) MetricTag() string {
	switch t {
	case RequestType:
		return "log_request"
	case ResponseType:
		return "log_response"
	}

	return t.String()
}

// String returns the subtype as a human-readable string.
func (t subtype) String() string {
	switch t {
	case RequestType:
		return "request"
	case ResponseType:
		return "response"
	}

	return string(t)
}

// formattedTime returns the UTC time the AuditEvent was created in the RFC3339Nano
// format (which removes trailing zeros from the seconds field).
func (a *Event) formattedTime() string {
	return a.Timestamp.UTC().Format(time.RFC3339Nano)
}

// isValidFormat provides a means to validate whether the supplied format is valid.
// Examples of valid formats are JSON and JSONx.
func isValidFormat(v string) bool {
	err := format(strings.TrimSpace(strings.ToLower(v))).validate()
	return err == nil
}
