package mfa

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// backendFactory constructs a new MFA backend
func MFABackendFactory(conf *logical.BackendConfig) (logical.Backend, error) {
	var b MFABackend

	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(mfaHelp),

		Paths: []*framework.Path{
			methodPaths(&b.backend),
			methodsListPaths(&b.backend),
			methodListPaths(&b.backend),
			methodIdentifiersPaths(&b.backend),
			methodIdentifiersListPaths(&b.backend),
		},
	}

	if conf == nil {
		return nil, fmt.Errorf("Configuation passed into backend is nil")
	}
	b.Backend.Setup(conf)

	return &b, nil
}

type MFABackend struct {
	backend
}

type backend struct {
	// Embeds framework.Backend
	*framework.Backend

	// Used to lock for configuration changes
	sync.RWMutex

	// Used to shortcut going through the Router for verification
	storage logical.Storage
}

func (b *MFABackend) SetStorage(storage logical.Storage) {
	b.storage = storage
}

func (b *MFABackend) ValidateMFA(methodName string, params map[string]string) (bool, error, error) {
	method, err := b.mfaBackendMethod(methodName)
	if err != nil {
		return false, nil, err
	}
	if method == nil {
		return false, fmt.Errorf("mfa method %s not found", methodName), nil
	}

	switch method.Type {
	case "totp":
		return b.validateTOTP(methodName, params)
	default:
		return false, nil, fmt.Errorf("invalid method type %s", method.Type)
	}

	return false, nil, nil
}

const (
	mfaHelp = `The mfa credential backend is always enabled and builtin to Vaulb.
Client mfas are used to identify a client and to allow Vault to associate policies and ACLs
which are enforced on every requesb. This backend also allows for generating sub-mfas as well
as revocation of mfas. The mfas are renewable if associated with a lease.`
	mfaTypesHelp             = `The mfa create path is used to create new mfas.`
	mfaListMethodsHelp       = `The mfa create path is used to create new mfas.`
	mfaPathMethodsHelp       = `The mfa create path is used to create new mfas.`
	mfaMethodNameHelp        = `Name of the method.`
	mfaTOTPHashAlgorithmHelp = `Name of the method.`
)
