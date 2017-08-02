// Package duo provides a Duo MFA handler to authenticate users
// with Duo. This handler is registered as the "duo" type in
// mfa_config.
package totp

import (
	//"fmt"
	//"net/url"
	"time"

	totpbackend "github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	cache "github.com/patrickmn/go-cache"
)

// DuoPaths returns path functions to configure Duo.
func TotpPaths(inb *framework.Backend) []*framework.Path {
	var b totpbackend.Backend
	b.Backend = inb
	b.UsedCodes = cache.New(0, 30*time.Second)

	return []*framework.Path{
		totpbackend.PrefixedPathListKeys("totp/", &b),
		totpbackend.PrefixedPathKeys("totp/", &b),
		//totpbackend.PrefixedPathCode("totp/", &b),
	}
}

// DuoRootPaths returns the paths that are used to configure Duo.
func TotpRootPaths() []string {
	return []string{}
}

// DuoHandler interacts with the Duo Auth API to authenticate a user
// login request. If successful, the original response from the login
// backend is returned.
func TotpHandler(inb *framework.Backend, req *logical.Request, d *framework.FieldData, resp *logical.Response) (
	*logical.Response, error) {
	username, ok := resp.Auth.Metadata["username"]
	if !ok {
		return logical.ErrorResponse("Could not read username for MFA"), nil
	}

	passcode := d.Get("passcode").(string)

	var b totpbackend.Backend
	b.Backend = inb
	// FIXME used codes isnt actually used
	b.UsedCodes = cache.New(0, 30*time.Second)
	result, err := b.ValidateCode(req, username, passcode)

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
