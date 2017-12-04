package vault

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/logical"
)

// N.B.: While we could use textproto to get the canonical mime header, HTTP/2
// requires all headers to be converted to lower case, so we just do that.

const (
	// Key used in the BarrierView to store and retrieve the header config
	passthroughHeadersEntry = "passthrough-headers"
	// Path used to create a sub view off of BarrierView
	passthroughHeadersSubPath = "passthrough-headers-config/"
)

type passthroughHeaderSettings struct {
	Backends []string `json:"backends"`
}

// PassthroughHeadersConfig is used by the router and core to
// determine which headers can be passed to which backends.
// It uses a BarrierView to persist the settings.
type PassthroughHeadersConfig struct {
	Headers map[string]*passthroughHeaderSettings

	view *BarrierView
	sync.RWMutex
}

// add adds or overwrites a header in the config and updates the barrier view
func (pt *PassthroughHeadersConfig) add(header string, backends []string) error {
	if header == "" {
		return fmt.Errorf("header value cannot be empty")
	}

	if len(backends) == 0 {
		return fmt.Errorf("backends value cannot be empty")
	}

	// Grab a write lock
	pt.Lock()
	defer pt.Unlock()

	if pt.Headers == nil {
		pt.Headers = make(map[string]*passthroughHeaderSettings, 1)
	}

	pt.Headers[strings.ToLower(header)] = &passthroughHeaderSettings{backends}
	entry, err := logical.StorageEntryJSON(passthroughHeadersEntry, pt.Headers)
	if err != nil {
		return fmt.Errorf("failed to persist passthrough headers config: %v", err)
	}

	if err := pt.view.Put(entry); err != nil {
		return fmt.Errorf("failed to persist passthrough headers config: %v", err)
	}

	return nil
}

// remove deletes a header out of the header config and updates the barrier view
func (pt *PassthroughHeadersConfig) remove(header string) error {
	if header == "" {
		return fmt.Errorf("header value cannot be empty")
	}

	// Grab a write lock
	pt.Lock()
	defer pt.Unlock()

	// Nothing to delete
	if len(pt.Headers) == 0 {
		return nil
	}

	delete(pt.Headers, strings.ToLower(header))
	entry, err := logical.StorageEntryJSON(passthroughHeadersEntry, pt.Headers)
	if err != nil {
		return fmt.Errorf("failed to persist passthrough headers config: %v", err)
	}

	if err := pt.view.Put(entry); err != nil {
		return fmt.Errorf("failed to persist passthrough headers config: %v", err)
	}

	return nil
}

// ApplyConfig returns a map of approved headers and their values, either
// hmac'ed or plaintext
func (pt *PassthroughHeadersConfig) ApplyConfig(headers map[string][]string, originalPath string) (result map[string][]string, retErr error) {
	if pt == nil {
		return nil, nil
	}

	// Grab a read lock
	pt.RLock()
	defer pt.RUnlock()

	// Make a copy of the incoming headers with everything lower so we can
	// case-insensitively compare
	lowerHeaders := make(map[string][]string, len(headers))
	for k, v := range headers {
		lowerHeaders[strings.ToLower(k)] = v
	}

	result = make(map[string][]string, len(pt.Headers))
	for key, settings := range pt.Headers {
		if val, ok := lowerHeaders[key]; ok {
			for _, allowedBackend := range settings.Backends {
				if strings.HasPrefix(originalPath, allowedBackend) {

					// copy the header values so we don't overwrite them
					hVals := make([]string, len(val))
					copy(hVals, val)

					result[key] = hVals
					break
				}
			}
		}
	}

	return result, nil
}

// Initialize the headers config by loading from the barrier view
func (c *Core) setupPassthroughHeadersConfig() error {
	// Create a sub-view
	view := c.systemBarrierView.SubView(passthroughHeadersSubPath)

	// Create the config
	out, err := view.Get(passthroughHeadersEntry)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	headers := make(map[string]*passthroughHeaderSettings)
	if out != nil {
		err = out.DecodeJSON(&headers)
		if err != nil {
			return err
		}
	}

	// Ensure that we are able to case-sensitively access the headers;
	// necessary for the upgrade case
	lowerHeaders := make(map[string]*passthroughHeaderSettings, len(headers))
	for k, v := range headers {
		lowerHeaders[strings.ToLower(k)] = v
	}

	c.passthroughHeaders = &PassthroughHeadersConfig{
		Headers: lowerHeaders,
		view:    view,
	}

	c.router.passthroughHeaders = c.passthroughHeaders

	return nil
}
