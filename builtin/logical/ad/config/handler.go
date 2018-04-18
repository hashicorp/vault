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
	StorageKey  = "config"

	// This length is arbitrarily chosen but should work for
	// most Active Directory minimum and maximum length settings.
	// A bit tongue-in-cheek since programmers love their base-2 exponents.
	DefaultPasswordLength = 64

	// The number of seconds in 32 days.
	DefaultPasswordTTLs = 24 * 60 * 60 * 32
)

func Handler(logger hclog.Logger) *handler {
	return &handler{
		logger:  logger,
		rwMutex: &sync.RWMutex{},
		config:  nil,
	}
}

type handler struct {
	logger  hclog.Logger
	rwMutex *sync.RWMutex
	config  *EngineConf
}

func (h *handler) Read(ctx context.Context, storage logical.Storage) (*EngineConf, error) {

	// If the config is cached, just return it.
	h.rwMutex.RLock()
	if h.config != nil {
		defer h.rwMutex.RUnlock()
		return h.config, nil
	}

	// The config is nil, so let's try to get it from storage.
	h.rwMutex.RUnlock()
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	entry, err := storage.Get(ctx, StorageKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		// Provide an unset error for a consistent error message
		// and to reduce the lines of code needed to safely
		// use the conf.
		return nil, &Unset{}
	}

	config := &EngineConf{&PasswordConf{}, &activedirectory.Configuration{}}
	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	h.config = config

	return h.config, nil
}

func (h *handler) Invalidate(ctx context.Context, key string) {
	if key == BackendPath {
		h.rwMutex.Lock()
		h.config = nil
		h.rwMutex.Unlock()
	}
}

func (h *handler) Path() *framework.Path {
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
			logical.UpdateOperation: h.updateOperation,
			logical.ReadOperation:   h.readOperation,
			logical.DeleteOperation: h.deleteOperation,
		},
	}
}

func (h *handler) updateOperation(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {

	// Parse and validate the desired conf.
	activeDirectoryConf, err := activedirectory.NewConfiguration(fieldData)
	if err != nil {
		return nil, err
	}
	passwordConf, err := newPasswordConfig(fieldData)
	if err != nil {
		return nil, err
	}
	config := &EngineConf{passwordConf, activeDirectoryConf}

	// Write and cache it.
	entry, err := logical.StorageEntryJSON(StorageKey, config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	h.rwMutex.Lock()
	h.config = config
	h.rwMutex.Unlock()

	// Respond.
	resp := &logical.Response{
		Data: config.Map(),
	}
	h.addWarnings(resp)

	return resp, nil
}

func (h *handler) readOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	_, err := h.Read(ctx, req.Storage)
	if err != nil {
		_, ok := err.(*Unset)
		if ok {
			return nil, nil
		}
		return nil, err
	}

	resp := &logical.Response{
		Data: h.config.Map(),
	}
	h.addWarnings(resp)

	return resp, nil
}

func (h *handler) deleteOperation(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {

	if err := req.Storage.Delete(ctx, StorageKey); err != nil {
		return nil, err
	}

	h.rwMutex.Lock()
	h.config = nil
	h.rwMutex.Unlock()

	return nil, nil
}

func (h *handler) addWarnings(resp *logical.Response) {

	resp.AddWarning("Access to this endpoint should be controlled via ACLs as it will return the configuration information as-is, including any passwords.")

	if !h.config.ADConf.StartTLS {
		resp.AddWarning("Because \"starttls\" is false, Vault is currently unable to rotate passwords, though setup can otherwise be completed.")
	}
}

// Reader is the interface through which those outside the config package
// can access the current *EngineConf.
// *handler itself fulfills the interface.
type Reader interface {

	// Read returns the present *EngineConf, which should not be stored in memory in case of ongoing changes.
	// If error == nil, *EngineConf != nil.
	//
	// The returned error may be due to issues reaching storage,
	// or it may be because the *EngineConf is unset by the user.
	//
	// If knowing the error is useful to the caller, it may be inspected like so:
	//
	// 		engineConf, err := Read(ctx, storage)
	//		if err != nil {
	//			_, ok := err.(*Unset)
	// 			...
	// 		}
	Read(ctx context.Context, storage logical.Storage) (*EngineConf, error)
}

type Unset struct{}

func (e *Unset) Error() string {
	return "the config is currently unset"
}
