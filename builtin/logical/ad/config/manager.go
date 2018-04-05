package config

import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const (
	BackendPath = "config"

	// This length is arbitrarily chosen but should work for
	// most Active Directory minimum and maximum length settings.
	// A bit tongue-in-cheek since programmers love their base-2 exponents.
	defaultPasswordLength = 64

	// The number of seconds in 32 days.
	defaultPasswordTTLs = 24 * 60 * 60 * 32
)

// NewManager creates a Manager, which manages all aspects of the config.
// Its only exported methods are the ones absolutely needed by the backend.
func NewManager(ctx context.Context, conf *logical.BackendConfig) (*Manager, error) {

	engineConf, err := readConfig(ctx, conf.StorageView)
	if err != nil {
		return nil, err
	}

	cache := newCache(engineConf)

	return &Manager{
		cache:  cache,
		logger: conf.Logger,
	}, nil
}

type Manager struct {
	cache  *cache
	reader *Reader
	logger hclog.Logger
}

func (m *Manager) ConfigReader() *Reader {
	return m.reader
}

func (m *Manager) Invalidate(ctx context.Context, key string) {
	if key == BackendPath {
		m.cache.Invalidate()
	}
}

func (m *Manager) Path() *framework.Path {
	opHandler := &operationHandler{
		logger: m.logger,
		cache:  m.cache,
	}
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

			"default_password_ttl": {
				Type:        framework.TypeDurationSecond,
				Default:     defaultPasswordTTLs,
				Description: "In seconds, the default password time-to-live.",
			},

			"max_password_ttl": {
				Type:        framework.TypeDurationSecond,
				Default:     defaultPasswordTTLs,
				Description: "In seconds, the maximum password time-to-live.",
			},

			"password_length": {
				Type:        framework.TypeInt,
				Default:     defaultPasswordLength,
				Description: "The desired length of passwords that Vault generates.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: opHandler.Delete,
			logical.ReadOperation:   opHandler.Read,
			logical.UpdateOperation: opHandler.Update,
		},
	}
}

// Reader pulls the current user conf from either the cache, or storage.
type Reader struct {
	cache *cache
}

func (r *Reader) Read(ctx context.Context, storage logical.Storage) (*EngineConf, error) {
	engineConf, ok := r.cache.Get()
	if ok {
		return engineConf, nil
	}
	return readConfig(ctx, storage)
}
