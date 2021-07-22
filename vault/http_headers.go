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
	httpHeadersKey          = "http_headers"
	httpHeadersPlaintextKey = "http_headers_plaintext"
)

// CommonResponseHeadersConfig contains common response headers configuration. This takes both a physical view and a barrier view
// because it is stored in both plaintext and encrypted to allow for getting the header
// values before the barrier is unsealed
type HttpHeadersConfig struct {
	l               sync.RWMutex
	physicalStorage physical.Backend
	barrierStorage  logical.Storage
	tlsDisabled     bool
	defaultHeaders  http.Header
}

// NewCommonResponseHeadersConfig creates a new CommonResponseHeadersConfig
func NewHttpHeadersConfig(tlsDisabled bool, physicalStorage physical.Backend, barrierStorage logical.Storage) *HttpHeadersConfig {

	defaultHeaders := setDefaultHeaders(tlsDisabled)

	return &HttpHeadersConfig{
		physicalStorage: physicalStorage,
		barrierStorage:  barrierStorage,
		tlsDisabled:     tlsDisabled,
		defaultHeaders:  defaultHeaders,
	}
}

func setDefaultHeaders(tlsDisabled bool) http.Header {
	defaultHeaders := http.Header{}
	defaultHeaders.Set("Content-Security-Policy", "default-src 'none'; connect-src 'self'; img-src 'self' data:; script-src 'self'; style-src 'unsafe-inline' 'self'; form-action 'none'; frame-ancestors 'none'; font-src 'self'")
	defaultHeaders.Set("X-XSS-Protection", "1; mode=block")
	defaultHeaders.Set("X-Frame-Options", "SAMEORIGIN")
	defaultHeaders.Set("X-Content-Type-Options", "nosniff")
	if !tlsDisabled {
		defaultHeaders.Set("Strict-Transport-Security", "max-age=63072000")
	}
	return defaultHeaders
}

// Headers returns the response headers
func (c *HttpHeadersConfig) Headers(ctx context.Context) (http.Header, error) {
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
func (c *HttpHeadersConfig) HeaderKeys(ctx context.Context) ([]string, error) {
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
func (c *HttpHeadersConfig) GetHeader(ctx context.Context, header string) ([]string, error) {
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
func (c *HttpHeadersConfig) SetHeader(ctx context.Context, header string, values []string) error {
	c.l.Lock()
	defer c.l.Unlock()

	config, err := c.get(ctx)
	if err != nil {
		return err
	}
	if config == nil {
		config = &httpHeadersConfigEntry{
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
func (c *HttpHeadersConfig) DeleteHeader(ctx context.Context, header string) error {
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

func (c *HttpHeadersConfig) get(ctx context.Context) (*httpHeadersConfigEntry, error) {
	// Read plaintext always to ensure in sync with barrier value
	plaintextConfigRaw, err := c.physicalStorage.Get(ctx, httpHeadersPlaintextKey)
	if err != nil {
		return nil, err
	}

	configRaw, err := c.barrierStorage.Get(ctx, httpHeadersKey)
	if err == nil {
		if configRaw == nil {
			return nil, nil
		}
		config := new(httpHeadersConfigEntry)
		if err := json.Unmarshal(configRaw.Value, config); err != nil {
			return nil, err
		}
		// Check that plaintext value matches barrier value, if not sync values
		if plaintextConfigRaw == nil || bytes.Compare(plaintextConfigRaw.Value, configRaw.Value) != 0 {
			if err := c.save(ctx, config); err != nil {
				return nil, err
			}
		}
		return config, nil
	}

	// Respond with error if not sealed
	if !strings.Contains(err.Error(), ErrBarrierSealed.Error()) {
		return nil, err
	}

	// Respond with plaintext value
	if configRaw == nil {
		return nil, nil
	}
	config := new(httpHeadersConfigEntry)
	if err := json.Unmarshal(plaintextConfigRaw.Value, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *HttpHeadersConfig) save(ctx context.Context, config *httpHeadersConfigEntry) error {
	if len(config.Headers) == 0 {
		if err := c.physicalStorage.Delete(ctx, httpHeadersPlaintextKey); err != nil {
			return err
		}
		return c.barrierStorage.Delete(ctx, httpHeadersKey)
	}

	configRaw, err := json.Marshal(config)
	if err != nil {
		return err
	}

	entry := &physical.Entry{
		Key:   httpHeadersPlaintextKey,
		Value: configRaw,
	}
	if err := c.physicalStorage.Put(ctx, entry); err != nil {
		return err
	}

	barrEntry := &logical.StorageEntry{
		Key:   httpHeadersKey,
		Value: configRaw,
	}
	return c.barrierStorage.Put(ctx, barrEntry)
}

type httpHeadersConfigEntry struct {
	Headers http.Header `json:"headers"`
}
