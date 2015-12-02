// Package duo provides a totp MFA handler.
// This handler is registered as the "totp" type in
// mfa_config.
package vaultTotp

import (
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	totp "github.com/ptollot/totp"
)

// TotpPaths returns path functions to configure Duo.
func TotpPaths() []*framework.Path {
	return []*framework.Path{
		pathTotpConfig(),
	}
}

// DuoRootPaths returns the paths that are used to configure Duo.
func TotpRootPaths() []string {
	return []string {
		"totp/config",
	}
}

// TotpHandler use the totp library to authenticate a user
// login request. If successful, the original response from the login
// backend is returned.
func TotpHandler(req *logical.Request, d *framework.FieldData, resp *logical.Response) (
	*logical.Response, error) {

	username, ok := resp.Auth.Metadata["username"]
	if !ok {
		return logical.ErrorResponse("Could not read username for MFA"), nil
	}

	var request *totpAuthRequest = &totpAuthRequest{}
	request.successResp = resp
	request.username = username
	request.method = d.Get("method").(string)
	request.passcode = d.Get("passcode").(string)

	return totpHandler(request)
}

type totpAuthRequest struct {
	successResp *logical.Response
	username string
	method string
	passcode string
}

func totpHandler(request *totpAuthRequest) (
	*logical.Response, error) {

	userPresent := totp.TotpReference.UserVerify(request.username)
	if userPresent == false {
		return nil, fmt.Errorf("unknown user")
	}

	correctPasscode, err := totp.TotpReference.Verify(request.username, request.passcode)
	if err != nil {
		return nil, err
	}

	if correctPasscode {
		return request.successResp, nil
	}

	return nil, fmt.Errorf("passcode not correct")
}
