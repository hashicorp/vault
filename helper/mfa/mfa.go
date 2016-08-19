// Package mfa provides wrappers to add multi-factor authentication
// to any auth backend.
//
// To add MFA to a backend, replace its login path with the
// paths returned by MFAPaths and add the additional root
// paths returned by MFARootPaths. The backend provides
// the username to the MFA wrapper in Auth.Metadata['username'].
//
// To add an additional MFA type, create a subpackage that
// implements [Type]Paths, [Type]RootPaths, and [Type]Handler
// functions and add them to MFAPaths, MFARootPaths, and
// handlers respectively.
package mfa

import (
	"github.com/hashicorp/vault/helper/mfa/duo"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

// MFAPaths returns paths to wrap the original login path and configure MFA.
// When adding MFA to a backend, these paths should be included instead of
// the login path in Backend.Paths.
func MFAPaths(originalBackend *framework.Backend, loginPath *framework.Path) []*framework.Path {
	var b backend
	b.Backend = originalBackend
	return append(duo.DuoPaths(), pathMFAConfig(&b), wrapLoginPath(&b, loginPath))
}

// MFARootPaths returns path strings used to configure MFA. When adding MFA
// to a backend, these paths should be included in
// Backend.PathsSpecial.Root.
func MFARootPaths() []string {
	return append(duo.DuoRootPaths(), "mfa_config")
}

// HandlerFunc is the callback called to handle MFA for a login request.
type HandlerFunc func(*logical.Request, *framework.FieldData, *logical.Response) (*logical.Response, error)

// handlers maps each supported MFA type to its handler.
var handlers = map[string]HandlerFunc{
	"duo": duo.DuoHandler,
}

type backend struct {
	*framework.Backend
}

func wrapLoginPath(b *backend, loginPath *framework.Path) *framework.Path {
	loginPath.Fields["passcode"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "One time passcode (optional)",
	}
	loginPath.Fields["method"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "Multi-factor auth method to use (optional)",
	}
	// wrap write callback to do MFA after auth
	loginHandler := loginPath.Callbacks[logical.UpdateOperation]
	loginPath.Callbacks[logical.UpdateOperation] = b.wrapLoginHandler(loginHandler)
	return loginPath
}

func (b *backend) wrapLoginHandler(loginHandler framework.OperationFunc) framework.OperationFunc {
	return func(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		// login with original login function first
		resp, err := loginHandler(req, d)
		if err != nil || resp.Auth == nil {
			return resp, err
		}

		// check if multi-factor enabled
		mfa_config, err := b.MFAConfig(req)
		if err != nil || mfa_config == nil {
			return resp, nil
		}

		// perform multi-factor authentication if type supported
		handler, ok := handlers[mfa_config.Type]
		if ok {
			return handler(req, d, resp)
		} else {
			return resp, err
		}
	}
}
