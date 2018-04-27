package advisor

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/salt"
)

// An Backend is an audit backend designed to
// communicate with an Advisor server over HTTP
type Backend struct {
	address   string
	formatter *audit.AuditFormatter
	config    audit.FormatterConfig
	salt      *salt.Salt
	saltMutex sync.RWMutex
}

var _ audit.Backend = &Backend{}

// LogRequest will serialize a Vault *Request into a JSON object
// and send that payload to Advisor over HTTP
func (backend *Backend) LogRequest(ctx context.Context, in *audit.LogInput) error {

	var reqBody bytes.Buffer

	err := backend.formatter.FormatRequest(ctx, &reqBody, backend.config, in)
	if err != nil {
		return err
	}

	httpResp, err := http.Post(backend.address, "application/json", &reqBody)
	if err != nil {
		return err
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode > 300 {
		return errors.New("request to backend failed")
	}
	return nil
}

// LogResponse will serialize a Vault *Request into a JSON object
// and send that payload to Advisor over HTTP
func (backend *Backend) LogResponse(ctx context.Context, in *audit.LogInput) error {

	var reqBody bytes.Buffer

	err := backend.formatter.FormatResponse(ctx, &reqBody, backend.config, in)
	if err != nil {
		return err
	}

	httpResp, err := http.Post(backend.address, "application/json", &reqBody)
	if err != nil {
		return err
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode > 300 {
		return errors.New("request to backend failed")
	}
	return nil
}

func (backend *Backend) setDefaultConfig() {
	backend.config = audit.FormatterConfig{}
}

func (backend *Backend) setDefaultFormatter() {
	backend.formatter = &audit.AuditFormatter{
		AuditFormatWriter: &audit.JSONFormatWriter{
			Prefix:   "",
			SaltFunc: backend.Salt,
		},
	}
}

// Pulls the address out of the configuration
func addressFromConfig(conf *audit.BackendConfig) (string, error) {
	if addr, ok := conf.Config["address"]; ok {
		return addr, nil
	}

	if addr, ok := conf.Config["a"]; ok {
		return addr, nil
	}

	return "", errors.New("no address provided")
}

// Factory returns a new Advisor audit backend
func Factory(ctx context.Context, a *audit.BackendConfig) (audit.Backend, error) {
	// First, create a new formatter for serializing
	// requests and responses into JSON
	advisor := &Backend{}
	advisor.setDefaultFormatter()
	advisor.setDefaultConfig()

	// Next, collect the address from the configuration
	addr, err := addressFromConfig(a)
	if err != nil {
		return nil, err
	}

	// Store the address on the new Backend struct.
	advisor.address = addr

	return advisor, nil
}

// Salt returns the salt singleton for this
// backend instance
func (backend *Backend) Salt(ctx context.Context) (*salt.Salt, error) {
	backend.saltMutex.RLock()
	if backend.salt != nil {
		defer backend.saltMutex.RUnlock()
		return backend.salt, nil
	}
	backend.saltMutex.RUnlock()
	backend.saltMutex.Lock()
	defer backend.saltMutex.Unlock()
	if backend.salt != nil {
		return backend.salt, nil
	}
	salt, err := salt.NewSalt(context.Background(), nil, nil)
	if err != nil {
		return nil, err
	}
	backend.salt = salt
	return salt, nil
}

// GetHash is needed to implement the AuditBackend interface
func (backend *Backend) GetHash(ctx context.Context, data string) (string, error) {
	salt, err := backend.Salt(ctx)
	if err != nil {
		return "", err
	}
	return salt.GetIdentifiedHMAC(data), nil
}

// Reload is required to implement the AuditBackend interface
func (backend *Backend) Reload(ctx context.Context) error {
	return nil
}

// Invalidate is needed to implement the AuditBackend interface
func (backend *Backend) Invalidate(ctx context.Context) {
	backend.saltMutex.Lock()
	defer backend.saltMutex.Unlock()
	backend.salt = nil
}
