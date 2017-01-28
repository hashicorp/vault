package audit

import (
	"fmt"
	"sync"

	"github.com/hashicorp/vault/helper/salt"
)

type auditedHeaderSettings struct {
	HMAC bool
}

type AuditedHeadersConfig struct {
	Headers map[string]*auditedHeaderSettings

	sync.RWMutex
}

func NewAuditedHeadersConfig() *AuditedHeadersConfig {
	return &AuditedHeadersConfig{
		Headers: make(map[string]*auditedHeaderSettings),
	}
}

func (a *AuditedHeadersConfig) Add(header string, hmac bool) {
	a.Lock()
	a.Headers[header] = &auditedHeaderSettings{hmac}
	a.Unlock()
}

func (a *AuditedHeadersConfig) Remove(header string) {
	a.Lock()
	delete(a.Headers, header)
	a.Unlock()
}

func (a *AuditedHeadersConfig) ApplyConfig(headers map[string][]string, salt *salt.Salt) (result map[string][]string, err error) {
	a.RLock()
	defer a.RUnlock()

	fmt.Println(a.Headers)

	result = make(map[string][]string)
	for key, val := range headers {
		hVals := make([]string, len(val))
		copy(hVals, val)

		if settings, ok := a.Headers[key]; ok {
			if settings.HMAC {
				fmt.Println("HMAC'ING")
				if err := Hash(salt, hVals); err != nil {
					return nil, err
				}
			}

			result[key] = hVals
		}
	}

	return
}
