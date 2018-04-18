package config

import (
	"context"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	BackendPath = "config"

	// This length is arbitrarily chosen but should work for
	// most Active Directory minimum and maximum length settings.
	// A bit tongue-in-cheek since programmers love their base-2 exponents.
	DefaultPasswordLength = 64

	// The number of seconds in 32 days.
	DefaultPasswordTTLs = 24 * 60 * 60 * 32
)

type Handler interface {
	// Config returns the current *EngineConf.
	// This shouldn't be retained in memory, but rather should be pulled on the fly every time,
	// in case the config changes.
	//
	// NOTE: If error is nil, *EngineConf still may be nil if it's unset by the user.
	Config(ctx context.Context, storage logical.Storage) (*EngineConf, error)
}

// NewManager creates a Manager, which manages all aspects of the config.
// Its only exported methods are the ones absolutely needed by the backend.
func NewManager(ctx context.Context, conf *logical.BackendConfig) (*Manager, error) {

	config, err := readConfig(ctx, conf.StorageView)
	if err != nil {
		return nil, err
	}
	rwMutex := &sync.RWMutex{}

	return &Manager{
		logger:  conf.Logger,
		rwMutex: rwMutex,
		config:  config,
	}, nil
}

type Manager struct {
	logger  hclog.Logger
	rwMutex *sync.RWMutex
	config  *EngineConf
}

type Unset struct{}

func (e *Unset) Error() string {
	return "config is currently unset"
}

func (m *Manager) Config(ctx context.Context, storage logical.Storage) (*EngineConf, error) {

	m.rwMutex.RLock()
	if m.config != nil {
		defer m.rwMutex.RUnlock()
		return m.config, nil
	}

	// upgrade the lock
	m.rwMutex.RUnlock()
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

	config, err := readConfig(ctx, storage)
	if err != nil {
		return nil, err
	}
	m.config = config

	if m.config == nil {
		// provide an unset error for a consistent error message
		// and to reduce the lines of code needed to safely
		// use the conf
		return nil, &Unset{}
	}
	return m.config, nil
}

func (m *Manager) Invalidate(ctx context.Context, key string) {
	if key == BackendPath {
		m.rwMutex.Lock()
		m.config = nil
		m.rwMutex.Unlock()
	}
}

func (m *Manager) Path() *framework.Path {
	return &framework.Path{
		Pattern: BackendPath,
		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "Username with sufficient permissions in Active Directory to administer passwords.",
			},

			"password": {
				Type:        framework.TypeString,
				Description: "Password for username with sufficient permissions in Active Directory to administer passwords.",
			},

			"urls": {
				Type:        framework.TypeCommaStringSlice,
				Default:     "ldap://127.0.0.1",
				Description: "LDAP URL to connect to (default: ldap://127.0.0.1). Multiple URLs can be specified by concatenating them with commas; they will be tried in-order.",
			},

			"certificate": {
				Type:        framework.TypeString,
				Description: "CA certificate to use when verifying LDAP server certificate, must be x509 PEM encoded.",
			},

			"dn": {
				Type:        framework.TypeString,
				Description: "The root distinguished name to bind to when managing service accounts.",
			},

			"insecure_tls": {
				Type:        framework.TypeBool,
				Description: "Skip LDAP server SSL Certificate verification - VERY insecure.",
			},

			"starttls": {
				Type:        framework.TypeBool,
				Default:     true,
				Description: "Issue a StartTLS command after establishing unencrypted connection.",
			},

			"tls_min_version": {
				Type:        framework.TypeString,
				Default:     "tls12",
				Description: "Minimum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'.",
			},

			"tls_max_version": {
				Type:        framework.TypeString,
				Default:     "tls12",
				Description: "Maximum TLS version to use. Accepted values are 'tls10', 'tls11' or 'tls12'. Defaults to 'tls12'.",
			},

			"ttl": {
				Type:        framework.TypeDurationSecond,
				Default:     DefaultPasswordTTLs,
				Description: "In seconds, the default password time-to-live.",
			},

			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Default:     DefaultPasswordTTLs,
				Description: "In seconds, the maximum password time-to-live.",
			},

			"password_length": {
				Type:        framework.TypeInt,
				Default:     DefaultPasswordLength,
				Description: "The desired length of passwords that Vault generates.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: m.delete,
			logical.ReadOperation:   m.read,
			logical.UpdateOperation: m.update,
		},
	}
}

func (m *Manager) delete(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	if err := deleteConfig(ctx, req.Storage); err != nil {
		return nil, err
	}

	m.rwMutex.Lock()
	m.config = nil
	m.rwMutex.Unlock()

	return nil, nil
}

func (m *Manager) read(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	_, err := m.Config(ctx, req.Storage)
	if err != nil {
		_, ok := err.(*Unset)
		if ok {
			return nil, nil
		}
		return nil, err
	}

	resp := &logical.Response{
		Data: m.config.Map(),
	}
	resp.AddWarning("read access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords.")
	return resp, nil
}

func (m *Manager) update(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {

	// Parse and validate the desired conf.
	activeDirectoryConf, err := activedirectory.NewConfiguration(m.logger, fieldData)
	if err != nil {
		return nil, err
	}
	passwordConf, err := newPasswordConfig(fieldData)
	if err != nil {
		return nil, err
	}
	config := &EngineConf{passwordConf, activeDirectoryConf}

	// Write and cache it.
	if err := writeConfig(ctx, req.Storage, config); err != nil {
		return nil, err
	}
	m.rwMutex.Lock()
	m.config = config
	m.rwMutex.Unlock()

	// Respond.
	resp := &logical.Response{
		Data: config.Map(),
	}
	resp.AddWarning("write access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords.")
	return resp, nil
}
