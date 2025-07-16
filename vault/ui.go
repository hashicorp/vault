// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
)

const (
	uiConfigKey          = "config"
	uiConfigPlaintextKey = "config_plaintext"
)

// UIConfig contains UI configuration. This takes both a physical view and a barrier view
// because it is stored in both plaintext and encrypted to allow for getting the header
// values before the barrier is unsealed
type UIConfig struct {
	l               sync.RWMutex
	physicalStorage physical.Backend
	barrierStorage  logical.Storage

	enabled        bool
	defaultHeaders http.Header
}

// NewUIConfig creates a new UI config
func NewUIConfig(enabled bool, physicalStorage physical.Backend, barrierStorage logical.Storage) *UIConfig {
	defaultHeaders := http.Header{}
	defaultHeaders.Set("Service-Worker-Allowed", "/")
	defaultHeaders.Set("X-Content-Type-Options", "nosniff")
	defaultHeaders.Set("Content-Security-Policy", "default-src 'none'; connect-src 'self'; img-src 'self' data:; script-src 'self'; style-src 'unsafe-inline' 'self'; form-action  'none'; frame-ancestors 'none'; font-src 'self'")

	return &UIConfig{
		physicalStorage: physicalStorage,
		barrierStorage:  barrierStorage,
		enabled:         enabled,
		defaultHeaders:  defaultHeaders,
	}
}

// Enabled returns if the UI is enabled
func (c *UIConfig) Enabled() bool {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.enabled
}

// Headers returns the response headers that should be returned in the UI
func (c *UIConfig) Headers(ctx context.Context) (http.Header, error) {
	c.l.RLock()
	defer c.l.RUnlock()

	config, err := c.get(ctx)
	if err != nil {
		return nil, err
	}
	headers := make(http.Header)
	if config != nil {
		headers = config.Headers
	}

	for k := range c.defaultHeaders {
		if headers.Get(k) == "" {
			v := c.defaultHeaders.Get(k)
			headers.Set(k, v)
		}
	}
	return headers, nil
}

// HeaderKeys returns the list of the configured headers
func (c *UIConfig) HeaderKeys(ctx context.Context) ([]string, error) {
	c.l.RLock()
	defer c.l.RUnlock()

	config, err := c.get(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}
	var keys []string
	for k := range config.Headers {
		keys = append(keys, k)
	}
	return keys, nil
}

// GetHeader retrieves the configured values for the given header
func (c *UIConfig) GetHeader(ctx context.Context, header string) ([]string, error) {
	c.l.RLock()
	defer c.l.RUnlock()

	config, err := c.get(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, nil
	}

	value := config.Headers.Values(header)
	return value, nil
}

// SetHeader sets the values for the given header
func (c *UIConfig) SetHeader(ctx context.Context, header string, values []string) error {
	c.l.Lock()
	defer c.l.Unlock()

	config, err := c.get(ctx)
	if err != nil {
		return err
	}
	if config == nil {
		config = &uiConfigEntry{
			Headers: http.Header{},
		}
	}

	// Clear custom header values before setting new
	config.Headers.Del(header)

	// Set new values
	for _, value := range values {
		config.Headers.Add(header, value)
	}
	return c.save(ctx, config)
}

// DeleteHeader deletes the header configuration for the given header
func (c *UIConfig) DeleteHeader(ctx context.Context, header string) error {
	c.l.Lock()
	defer c.l.Unlock()

	config, err := c.get(ctx)
	if err != nil {
		return err
	}
	if config == nil {
		return nil
	}

	config.Headers.Del(header)
	return c.save(ctx, config)
}

func (c *UIConfig) get(ctx context.Context) (*uiConfigEntry, error) {
	// Read plaintext always to ensure in sync with barrier value
	plaintextConfigRaw, err := c.physicalStorage.Get(ctx, uiConfigPlaintextKey)
	if err != nil {
		return nil, err
	}
	configRaw, uiConfigGetErr := c.barrierStorage.Get(ctx, uiConfigKey)

	// Respond with error only if not sealed, otherwise do not throw the error
	if uiConfigGetErr != nil && !strings.Contains(uiConfigGetErr.Error(), ErrBarrierSealed.Error()) {
		return nil, uiConfigGetErr
	}
	if configRaw == nil {
		return nil, nil
	}

	config := new(uiConfigEntry)
	if config == nil {
		return nil, nil
	}
	if err := json.Unmarshal(configRaw.Value, config); err != nil {
		return nil, err
	}

	// Check that plaintext value matches barrier value, if not sync values
	if uiConfigGetErr == nil && (plaintextConfigRaw == nil ||
		!bytes.Equal(plaintextConfigRaw.Value, configRaw.Value)) {
		if err := c.save(ctx, config); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (c *UIConfig) save(ctx context.Context, config *uiConfigEntry) error {
	if len(config.Headers) == 0 {
		if err := c.physicalStorage.Delete(ctx, uiConfigPlaintextKey); err != nil {
			return err
		}
		return c.barrierStorage.Delete(ctx, uiConfigKey)
	}

	configRaw, err := json.Marshal(config)
	if err != nil {
		return err
	}

	entry := &physical.Entry{
		Key:   uiConfigPlaintextKey,
		Value: configRaw,
	}
	if err := c.physicalStorage.Put(ctx, entry); err != nil {
		return err
	}

	barrEntry := &logical.StorageEntry{
		Key:   uiConfigKey,
		Value: configRaw,
	}
	return c.barrierStorage.Put(ctx, barrEntry)
}

type uiConfigEntry struct {
	Headers http.Header `json:"headers"`
}
