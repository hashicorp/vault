// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (b *backend) secretAccessTokenResponse(ctx context.Context, s logical.Storage, tokenGen *TokenGenerator) (*logical.Response, error) {
	if tokenGen == nil || tokenGen.KeyName == "" {
		return logical.ErrorResponse("invalid token generator has no service account key"), nil
	}

	token, err := tokenGen.getAccessToken(ctx)
	if err != nil {
		return logical.ErrorResponse("unable to generate token - make sure your roleset service account and key are still valid: %v", err), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"token":              token.AccessToken,
			"token_ttl":          token.Expiry.UTC().Sub(time.Now().UTC()) / (time.Second),
			"expires_at_seconds": token.Expiry.Unix(),
		},
	}, nil
}

func (tg *TokenGenerator) getAccessToken(ctx context.Context) (*oauth2.Token, error) {
	jsonBytes, err := base64.StdEncoding.DecodeString(tg.B64KeyJSON)
	if err != nil {
		return nil, errwrap.Wrapf("could not b64-decode key data: {{err}}", err)
	}

	cfg, err := google.JWTConfigFromJSON(jsonBytes, tg.Scopes...)
	if err != nil {
		return nil, errwrap.Wrapf("could not generate token JWT config: {{err}}", err)
	}

	tkn, err := cfg.TokenSource(ctx).Token()
	if err != nil {
		return nil, errwrap.Wrapf("got error while creating OAuth2 token: {{err}}", err)
	}
	return tkn, err
}

const deprecationWarning = `
This endpoint no longer generates leases due to limitations of the GCP API, as OAuth2 tokens belonging to Service
Accounts cannot be revoked. This access_token and lease were created by a previous version of the GCP secrets
engine and will be cleaned up now. Note that there is the chance that this access_token, if not already expired,
will still be valid up to one hour.
`

const pathTokenHelpSyn = `Generate an OAuth2 access token secret.`
const pathTokenHelpDesc = `
This path will generate a new OAuth2 access token for accessing GCP APIs.

Either specify "roleset/my-roleset" or "static/my-account" to generate a key corresponding
to a roleset or static account respectively.

Please see backend documentation for more information:
https://www.vaultproject.io/docs/secrets/gcp/index.html
`

// THIS SECRET TYPE IS DEPRECATED - future secret requests returns a response with no framework.Secret
// We are keeping them as part of the created framework.Secret
// to allow for clean up of access_token secrets and leases
// from older versions of Vault.
const SecretTypeAccessToken = "access_token"

func secretAccessToken(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTypeAccessToken,
		Fields: map[string]*framework.FieldSchema{
			"token": {
				Type:        framework.TypeString,
				Description: "OAuth2 token",
			},
		},
		Renew:  b.secretAccessTokenRenew,
		Revoke: b.secretAccessTokenRevoke,
	}
}

// Renewal will still return an error, but return the warning in case as well.
func (b *backend) secretAccessTokenRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := logical.ErrorResponse("short-term access tokens cannot be renewed - request new access token instead")
	resp.AddWarning(deprecationWarning)
	return resp, nil
}

// Revoke will no-op and pass but warn the user. This is mostly to clean up old leases.
// Any associated secret (access_token) has already expired and thus doesn't need to
// actually be revoked,  or will expire within an hour and currently can't actually be revoked anyways.
func (b *backend) secretAccessTokenRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp := &logical.Response{}
	resp.AddWarning(deprecationWarning)
	return resp, nil
}
