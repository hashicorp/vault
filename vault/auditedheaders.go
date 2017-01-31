package vault

import (
	"fmt"
	"sync"

	"github.com/hashicorp/vault/logical"
)

const (
	auditedHeadersEntry   = "audited_headers"
	auditedHeadersSubPath = "auditedHeadersConfig/"
)

type auditedHeaderSettings struct {
	HMAC bool
}

type AuditedHeadersConfig struct {
	Headers map[string]*auditedHeaderSettings

	view *BarrierView
	sync.RWMutex
}

func NewAuditedHeadersConfig() *AuditedHeadersConfig {
	return &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
	}
}

func (a *AuditedHeadersConfig) add(header string, hmac bool) error {
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

func (a *AuditedHeadersConfig) remove(header string) error {
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

func (a *AuditedHeadersConfig) ApplyConfig(headers map[string][]string, hashFunc func(string) string) (result map[string][]string) {
	a.RLock()
	defer a.RUnlock()

	result = make(map[string][]string)
	for key, settings := range a.Headers {
		if val, ok := headers[key]; ok {
			hVals := make([]string, len(val))
			copy(hVals, val)

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
