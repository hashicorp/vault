package mfa

import (
	"github.com/hashicorp/vault/helper/mfa/duo"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func MFAPaths(originalBackend *framework.Backend, loginPath *framework.Path) []*framework.Path {
	var b backend
	b.Backend = originalBackend
	return append(duo.DuoPaths(), pathMFAConfig(&b), wrapLoginPath(&b, loginPath))
}

func MFAPathsSpecial() []string {
	return append(duo.DuoPathsSpecial(), "mfa_config")
}

type backend struct {
	*framework.Backend
}

func wrapLoginPath(b *backend, loginPath *framework.Path) *framework.Path {
	(*loginPath).Fields["passcode"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "One time passcode (optional)",
	}
	(*loginPath).Fields["method"] = &framework.FieldSchema{
		Type:      framework.TypeString,
		Description: "Multi-factor auth method to use (optional)",
	}
	// wrap write callback to do duo two factor after auth
	loginHandler := loginPath.Callbacks[logical.WriteOperation]
	loginPath.Callbacks[logical.WriteOperation] = b.wrapLoginHandler(loginHandler)
	return loginPath
}

func (b *backend) wrapLoginHandler(loginHandler framework.OperationFunc) framework.OperationFunc {
	return func (req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		// login with original login function first
		resp, err := loginHandler(req, d);
		if err != nil || resp.Auth == nil {
			return resp, err
		}

		// check if multi-factor enabled
		mfa_config, err := b.MFAConfig(req)
		if err != nil || mfa_config == nil {
			return resp, nil
		}

		switch (mfa_config.Type) {
		case "duo":
			return duo.DuoHandler(req, d, resp)
		default:
			return resp, err
		}
	}
}
