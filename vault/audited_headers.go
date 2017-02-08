package vault

import (
	"fmt"
	"sync"

	"github.com/hashicorp/vault/logical"
)

const (
	// Key used in the BarrierView to store and retrieve the header config
	auditedHeadersEntry = "audited-headers"
	// Path used to create a sub view off of BarrierView
	auditedHeadersSubPath = "audited-headers-config/"
)

type auditedHeaderSettings struct {
	HMAC bool `json:"hmac"`
}

// AuditedHeadersConfig is used by the Audit Broker to write only approved
// headers to the audit logs. It uses a BarrierView to persist the settings.
type AuditedHeadersConfig struct {
	Headers map[string]*auditedHeaderSettings

	view *BarrierView
	sync.RWMutex
}

// add adds or overwrites a header in the config and updates the barrier view
func (a *AuditedHeadersConfig) add(header string, hmac bool) error {
	if header == "" {
		return fmt.Errorf("header value cannot be empty")
	}

	// Grab a write lock
	a.Lock()
	defer a.Unlock()

	a.Headers[header] = &auditedHeaderSettings{hmac}
	entry, err := logical.StorageEntryJSON(auditedHeadersEntry, a.Headers)
	if err != nil {
		return fmt.Errorf("failed to persist audited headers config: %v", err)
	}

	if err := a.view.Put(entry); err != nil {
		return fmt.Errorf("failed to persist audited headers config: %v", err)
	}

	return nil
}

// remove deletes a header out of the header config and updates the barrier view
func (a *AuditedHeadersConfig) remove(header string) error {
	if header == "" {
		return fmt.Errorf("header value cannot be empty")
	}

	// Grab a write lock
	a.Lock()
	defer a.Unlock()

	delete(a.Headers, header)
	entry, err := logical.StorageEntryJSON(auditedHeadersEntry, a.Headers)
	if err != nil {
		return fmt.Errorf("failed to persist audited headers config: %v", err)
	}

	if err := a.view.Put(entry); err != nil {
		return fmt.Errorf("failed to persist audited headers config: %v", err)
	}

	return nil
}

// ApplyConfig returns a map of approved headers and their values, either
// hmac'ed or plaintext
func (a *AuditedHeadersConfig) ApplyConfig(headers map[string][]string, hashFunc func(string) string) (result map[string][]string) {
	// Grab a read lock
	a.RLock()
	defer a.RUnlock()

	result = make(map[string][]string, len(a.Headers))
	for key, settings := range a.Headers {
		if val, ok := headers[key]; ok {
			// copy the header values so we don't overwrite them
			hVals := make([]string, len(val))
			copy(hVals, val)

			// Optionally hmac the values
			if settings.HMAC {
				for i, el := range hVals {
					hVals[i] = hashFunc(el)
				}
			}

			result[key] = hVals
		}
	}

	return
}

// Initalize the headers config by loading from the barrier view
func (c *Core) setupAuditedHeadersConfig() error {
	// Create a sub-view
	view := c.systemBarrierView.SubView(auditedHeadersSubPath)

	// Create the config
	out, err := view.Get(auditedHeadersEntry)
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	headers := make(map[string]*auditedHeaderSettings)
	if out != nil {
		err = out.DecodeJSON(&headers)
		if err != nil {
			return err
		}
	}

	c.auditedHeaders = &AuditedHeadersConfig{
		Headers: headers,
		view:    view,
	}

	return nil
}
