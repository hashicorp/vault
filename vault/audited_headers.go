// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/logical"
)

// N.B.: While we could use textproto to get the canonical mime header, HTTP/2
// requires all headers to be converted to lower case, so we just do that.

const (
	// Key used in the BarrierView to store and retrieve the header config
	auditedHeadersEntry = "audited-headers"
	// Path used to create a sub view off of BarrierView
	auditedHeadersSubPath = "audited-headers-config/"
)

// auditedHeadersKey returns the key at which audit header configuration is stored.
func auditedHeadersKey() string {
	return auditedHeadersSubPath + auditedHeadersEntry
}

type auditedHeaderSettings struct {
	// HMAC is used to indicate whether the value of the header should be HMAC'd.
	HMAC bool `json:"hmac"`
}

// AuditedHeadersConfig is used by the Audit Broker to write only approved
// headers to the audit logs. It uses a BarrierView to persist the settings.
type AuditedHeadersConfig struct {
	// Headers stores the current headers that should be audited, and their settings.
	Headers map[string]*auditedHeaderSettings

	// view is the barrier view which should be used to access underlying audit header config data.
	view *BarrierView

	sync.RWMutex
}

// NewAuditedHeadersConfig should be used to create AuditedHeadersConfig.
func NewAuditedHeadersConfig(view *BarrierView) (*AuditedHeadersConfig, error) {
	if view == nil {
		return nil, fmt.Errorf("barrier view cannot be nil")
	}

	// This should be the only place where the AuditedHeadersConfig struct is initialized.
	// Store the view so that we can reload headers when we 'invalidate'.
	return &AuditedHeadersConfig{view: view}, nil
}

// add adds or overwrites a header in the config and updates the barrier view
// NOTE: add will acquire a write lock in order to update the underlying headers.
func (a *AuditedHeadersConfig) add(ctx context.Context, header string, hmac bool) error {
	if header == "" {
		return fmt.Errorf("header value cannot be empty")
	}

	// Grab a write lock
	a.Lock()
	defer a.Unlock()

	if a.Headers == nil {
		a.Headers = make(map[string]*auditedHeaderSettings, 1)
	}

	a.Headers[strings.ToLower(header)] = &auditedHeaderSettings{hmac}
	entry, err := logical.StorageEntryJSON(auditedHeadersEntry, a.Headers)
	if err != nil {
		return fmt.Errorf("failed to persist audited headers config: %w", err)
	}

	if err := a.view.Put(ctx, entry); err != nil {
		return fmt.Errorf("failed to persist audited headers config: %w", err)
	}

	return nil
}

// remove deletes a header out of the header config and updates the barrier view
// NOTE: remove will acquire a write lock in order to update the underlying headers.
func (a *AuditedHeadersConfig) remove(ctx context.Context, header string) error {
	if header == "" {
		return fmt.Errorf("header value cannot be empty")
	}

	// Grab a write lock
	a.Lock()
	defer a.Unlock()

	// Nothing to delete
	if len(a.Headers) == 0 {
		return nil
	}

	delete(a.Headers, strings.ToLower(header))
	entry, err := logical.StorageEntryJSON(auditedHeadersEntry, a.Headers)
	if err != nil {
		return fmt.Errorf("failed to persist audited headers config: %w", err)
	}

	if err := a.view.Put(ctx, entry); err != nil {
		return fmt.Errorf("failed to persist audited headers config: %w", err)
	}

	return nil
}

// invalidate attempts to refresh the allowed audit headers and their settings.
// NOTE: invalidate will acquire a write lock in order to update the underlying headers.
func (a *AuditedHeadersConfig) invalidate(ctx context.Context) error {
	a.Lock()
	defer a.Unlock()

	// Get the actual headers entries, e.g. sys/audited-headers-config/audited-headers
	out, err := a.view.Get(ctx, auditedHeadersEntry)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	if out == nil {
		// We didn't get any data from this view.
		return nil
	}

	headers := make(map[string]*auditedHeaderSettings)
	err = out.DecodeJSON(&headers)
	if err != nil {
		return err
	}

	// Ensure that we are able to case-sensitively access the headers;
	// necessary for the upgrade case
	lowerHeaders := make(map[string]*auditedHeaderSettings, len(headers))
	for k, v := range headers {
		lowerHeaders[strings.ToLower(k)] = v
	}

	a.Headers = headers
	return nil
}

// ApplyConfig returns a map of approved headers and their values, either hmac'ed or plaintext.
func (a *AuditedHeadersConfig) ApplyConfig(ctx context.Context, headers map[string][]string, salter audit.Salter) (result map[string][]string, retErr error) {
	// Grab a read lock
	a.RLock()
	defer a.RUnlock()

	// Make a copy of the incoming headers with everything lower so we can
	// case-insensitively compare
	lowerHeaders := make(map[string][]string, len(headers))
	for k, v := range headers {
		lowerHeaders[strings.ToLower(k)] = v
	}

	result = make(map[string][]string, len(a.Headers))
	for key, settings := range a.Headers {
		if val, ok := lowerHeaders[key]; ok {
			// copy the header values so we don't overwrite them
			hVals := make([]string, len(val))
			copy(hVals, val)

			// Optionally hmac the values
			if settings.HMAC {
				for i, el := range hVals {
					hVal, err := audit.HashString(ctx, salter, el)
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

// setupAuditedHeadersConfig will initialize new audited headers configuration on
// the Core by loading data from the barrier view.
func (c *Core) setupAuditedHeadersConfig(ctx context.Context) error {
	// Create a sub-view, e.g. sys/audited-headers-config/
	view := c.systemBarrierView.SubView(auditedHeadersSubPath)

	headers, err := NewAuditedHeadersConfig(view)
	if err != nil {
		return err
	}

	// Invalidate the headers now in order to load them for the first time.
	err = headers.invalidate(ctx)
	if err != nil {
		return err
	}

	// Update the Core.
	c.auditedHeaders = headers

	return nil
}
