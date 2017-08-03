// Package duo provides a TOTP MFA handler to authenticate users
// with TOTP. This handler is registered as the "totp" type in
// mfa_config.
package totp

import (
	"time"

	totpbackend "github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	cache "github.com/patrickmn/go-cache"
)

type backend struct {
	*totpbackend.Backend
}

// TotpPaths returns path functions to configure TOTP credentials.
func TotpPaths(inb *framework.Backend) []*framework.Path {
	var b totpbackend.Backend
	b.Backend = inb

	return []*framework.Path{
		totpbackend.PrefixedPathListKeys("totp/", &b),
		totpbackend.PrefixedPathKeys("totp/", &b),
		//totpbackend.PrefixedPathCode("totp/", &b),
	}
}

// FIXME?
// TotpRootPaths returns the paths that are used to configure TOTP.
func TotpRootPaths() []string {
	return []string{}
}

func GetTotpHandler(inb *framework.Backend) func(req *logical.Request, d *framework.FieldData, resp *logical.Response) (*logical.Response, error) {
	var b backend
	var bb totpbackend.Backend
	bb.Backend = inb
	bb.UsedCodes = cache.New(0, 30*time.Second)
	b.Backend = &bb

	return b.TotpHandler
}

// TotpHandler interacts with the builtin totp backend to authenticate a user
// login request. If successful, the original response from the login
// backend is returned.
func (b *backend) TotpHandler(req *logical.Request, d *framework.FieldData, resp *logical.Response) (
	*logical.Response, error) {
	username, ok := resp.Auth.Metadata["username"]
	if !ok {
		return logical.ErrorResponse("Could not read username for MFA"), nil
	}

	passcode := d.Get("passcode").(string)

	result, err := b.Backend.ValidateCode(req, username, passcode)

	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return result, nil
	}

	if !result.Data["valid"].(bool) {
		return logical.ErrorResponse("The specified passcode is not valid"), nil
	}
	return resp, nil
}
