// Package duo provides a TOTP MFA handler to authenticate users
// with TOTP. This handler is registered as the "totp" type in
// mfa_config.
package totp

import (
	"time"

	"github.com/hashicorp/vault/helper/totputil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	cache "github.com/patrickmn/go-cache"
)

type backend struct {
	*totputil.Backend
}

// TotpPaths returns path functions to configure TOTP credentials.
func TotpPaths(fb *framework.Backend) []*framework.Path {
	var b totputil.Backend
	b.Backend = fb

	return []*framework.Path{
		b.PathListKeys("totp/"),
		b.PathKeys("totp/"),
		// We omit code generation / validation paths
	}
}

// FIXME?
// TotpRootPaths returns the paths that are used to configure TOTP.
func TotpRootPaths() []string {
	return []string{}
}

func GetTotpHandler(fb *framework.Backend) func(req *logical.Request, d *framework.FieldData, resp *logical.Response) (*logical.Response, error) {
	var b backend
	var tb totputil.Backend
	tb.Backend = fb
	tb.UsedCodes = cache.New(0, 30*time.Second)
	b.Backend = &tb

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
