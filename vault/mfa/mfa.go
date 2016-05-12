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

type backend struct {
	// Embeds framework.Backend
	*framework.Backend

	// Used to lock for configuration changes
	sync.RWMutex

	// Used to avoid going through the Router for verification
	storage logical.Storage
}

// MFABackend wraps the internal backend object to allow us to make public
// methods on it
type MFABackend struct {
	backend
}

// Sets the storage for the backend. Since this backend is a singleton, the
// storage is constant and we do not need to rely on the router giving the
// right storage for a request.
func (b *MFABackend) SetStorage(storage logical.Storage) {
	b.storage = storage
}

// ValidateMFA looks for the given method name and if found, attempts vaidation
// with the given parameters
func (b *MFABackend) ValidateMFA(methodName string, mfaInfo *logical.MFAInfo) (bool, error, error) {
	if mfaInfo == nil {
		return false, nil, fmt.Errorf("nil mfa information supplied for validation")
	}

	method, err := b.mfaBackendMethod(methodName)
	if err != nil {
		return false, nil, err
	}
	if method == nil {
		return false, fmt.Errorf("mfa method %s not found", methodName), nil
	}

	switch method.Type {
	case "totp":
		return b.validateTOTP(methodName, mfaInfo)
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
