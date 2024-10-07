// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package audit

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/logical"
)

// N.B.: While we could use textproto to get the canonical mime header, HTTP/2
// requires all headers to be converted to lower case, so we just do that.

const (
	// auditedHeadersEntry is the key used in storage to store and retrieve the header config
	auditedHeadersEntry = "audited-headers"

	// AuditedHeadersSubPath is the path used to create a sub view within storage.
	AuditedHeadersSubPath = "audited-headers-config/"
)

type durableStorer interface {
	Get(ctx context.Context, key string) (*logical.StorageEntry, error)
	Put(ctx context.Context, entry *logical.StorageEntry) error
}

// HeaderFormatter is an interface defining the methods of the
// vault.HeadersConfig structure needed in this package.
type HeaderFormatter interface {
	// ApplyConfig returns a map of header values that consists of the
	// intersection of the provided set of header values with a configured
	// set of headers and will hash headers that have been configured as such.
	ApplyConfig(context.Context, map[string][]string, Salter) (map[string][]string, error)
}

// AuditedHeadersKey returns the key at which audit header configuration is stored.
func AuditedHeadersKey() string {
	return AuditedHeadersSubPath + auditedHeadersEntry
}

type headerSettings struct {
	// HMAC is used to indicate whether the value of the header should be HMAC'd.
	HMAC bool `json:"hmac"`
}

// HeadersConfig is used by the Audit Broker to write only approved
// headers to the audit logs. It uses a BarrierView to persist the settings.
type HeadersConfig struct {
	// headerSettings stores the current headers that should be audited, and their settings.
	headerSettings map[string]*headerSettings

	// view is the barrier view which should be used to access underlying audit header config data.
	view durableStorer

	sync.RWMutex
}

// NewHeadersConfig should be used to create HeadersConfig.
func NewHeadersConfig(view durableStorer) (*HeadersConfig, error) {
	if view == nil {
		return nil, fmt.Errorf("barrier view cannot be nil")
	}

	// This should be the only place where the HeadersConfig struct is initialized.
	// Store the view so that we can reload headers when we 'Invalidate'.
	return &HeadersConfig{
		view:           view,
		headerSettings: make(map[string]*headerSettings),
	}, nil
}

// Header attempts to retrieve a copy of the settings associated with the specified header.
// The second boolean return parameter indicates whether the header existed in configuration,
// it should be checked as when 'false' the returned settings will have the default values.
func (a *HeadersConfig) Header(name string) (headerSettings, bool) {
	a.RLock()
	defer a.RUnlock()

	var s headerSettings
	v, ok := a.headerSettings[strings.ToLower(name)]

	if ok {
		s.HMAC = v.HMAC
	}

	return s, ok
}

// Headers returns all existing headers along with a copy of their current settings.
func (a *HeadersConfig) Headers() map[string]headerSettings {
	a.RLock()
	defer a.RUnlock()

	// We know how many entries the map should have.
	headers := make(map[string]headerSettings, len(a.headerSettings))

	// Clone the headers
	for name, setting := range a.headerSettings {
		headers[name] = headerSettings{HMAC: setting.HMAC}
	}

	return headers
}

// Add adds or overwrites a header in the config and updates the barrier view
// NOTE: Add will acquire a write lock in order to update the underlying headers.
func (a *HeadersConfig) Add(ctx context.Context, header string, hmac bool) error {
	if header == "" {
		return fmt.Errorf("header value cannot be empty")
	}

	// Grab a write lock
	a.Lock()
	defer a.Unlock()

	if a.headerSettings == nil {
		a.headerSettings = make(map[string]*headerSettings, 1)
	}

	a.headerSettings[strings.ToLower(header)] = &headerSettings{hmac}
	entry, err := logical.StorageEntryJSON(auditedHeadersEntry, a.headerSettings)
	if err != nil {
		return fmt.Errorf("failed to persist audited headers config: %w", err)
	}

	if err := a.view.Put(ctx, entry); err != nil {
		return fmt.Errorf("failed to persist audited headers config: %w", err)
	}

	return nil
}

// Remove deletes a header out of the header config and updates the barrier view
// NOTE: Remove will acquire a write lock in order to update the underlying headers.
func (a *HeadersConfig) Remove(ctx context.Context, header string) error {
	if header == "" {
		return fmt.Errorf("header value cannot be empty")
	}

	// Grab a write lock
	a.Lock()
	defer a.Unlock()

	// Nothing to delete
	if len(a.headerSettings) == 0 {
		return nil
	}

	delete(a.headerSettings, strings.ToLower(header))
	entry, err := logical.StorageEntryJSON(auditedHeadersEntry, a.headerSettings)
	if err != nil {
		return fmt.Errorf("failed to persist audited headers config: %w", err)
	}

	if err := a.view.Put(ctx, entry); err != nil {
		return fmt.Errorf("failed to persist audited headers config: %w", err)
	}

	return nil
}

// DefaultHeaders can be used to retrieve the set of default headers that will be
// added to HeadersConfig in order to allow them to appear in audit logs in a raw
// format. If the Vault Operator adds their own setting for any of the defaults,
// their setting will be honored.
func (a *HeadersConfig) DefaultHeaders() map[string]*headerSettings {
	// Support deprecated 'x-' prefix (https://datatracker.ietf.org/doc/html/rfc6648)
	const correlationID = "correlation-id"
	xCorrelationID := fmt.Sprintf("x-%s", correlationID)

	return map[string]*headerSettings{
		correlationID:  {},
		xCorrelationID: {},
		"user-agent":   {},
	}
}

// Invalidate attempts to refresh the allowed audit headers and their settings.
// NOTE: Invalidate will acquire a write lock in order to update the underlying headers.
func (a *HeadersConfig) Invalidate(ctx context.Context) error {
	a.Lock()
	defer a.Unlock()

	// Get the actual headers entries, e.g. sys/audited-headers-config/audited-headers
	out, err := a.view.Get(ctx, auditedHeadersEntry)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// If we cannot update the stored 'new' headers, we will clear the existing
	// ones as part of invalidation.
	headers := make(map[string]*headerSettings)
	if out != nil {
		err = out.DecodeJSON(&headers)
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}
	}

	// Ensure that we are able to case-sensitively access the headers;
	// necessary for the upgrade case
	lowerHeaders := make(map[string]*headerSettings, len(headers))
	for k, v := range headers {
		lowerHeaders[strings.ToLower(k)] = v
	}

	// Ensure that we have default headers configured to appear in the audit log.
	// Add them if they're missing.
	for header, setting := range a.DefaultHeaders() {
		if _, ok := lowerHeaders[header]; !ok {
			lowerHeaders[header] = setting
		}
	}

	a.headerSettings = lowerHeaders
	return nil
}

// ApplyConfig returns a map of approved headers and their values, either HMAC'd or plaintext.
// If the supplied headers are empty or nil, an empty set of headers will be returned.
func (a *HeadersConfig) ApplyConfig(ctx context.Context, headers map[string][]string, salter Salter) (result map[string][]string, retErr error) {
	// Return early if we don't have headers.
	if len(headers) < 1 {
		return map[string][]string{}, nil
	}

	// Grab a read lock
	a.RLock()
	defer a.RUnlock()

	// Make a copy of the incoming headers with everything lower so we can
	// case-insensitively compare
	lowerHeaders := make(map[string][]string, len(headers))
	for k, v := range headers {
		lowerHeaders[strings.ToLower(k)] = v
	}

	result = make(map[string][]string, len(a.headerSettings))
	for key, settings := range a.headerSettings {
		if val, ok := lowerHeaders[key]; ok {
			// copy the header values so we don't overwrite them
			hVals := make([]string, len(val))
			copy(hVals, val)

			// Optionally hmac the values
			if settings.HMAC {
				for i, el := range hVals {
					hVal, err := hashString(ctx, salter, el)
					if err != nil {
						return nil, err
					}
					hVals[i] = hVal
				}
			}

			result[key] = hVals
		}
	}

	return result, nil
}
