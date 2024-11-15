// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

func pathImpersonatedAccountSecretAccessToken(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s/%s/token", impersonatedAccountPathPrefix, framework.GenericNameRegex("name")),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationVerb:   "generate",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Required. Name of the impersonated account.",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathImpersonatedAccountAccessToken,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "impersonated-account-access-token",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathImpersonatedAccountAccessToken,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationSuffix: "impersonated-account-access-token2",
				},
			},
		},
		HelpSynopsis:    pathTokenHelpSyn,
		HelpDescription: pathTokenHelpDesc,
	}
}

func (b *backend) pathImpersonatedAccountAccessToken(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	acctName := d.Get("name").(string)

	acct, err := b.getImpersonatedAccount(acctName, ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if acct == nil {
		return logical.ErrorResponse("impersonated account %q does not exists", acctName), nil
	}

	creds, err := b.credentials(req.Storage)
	if err != nil {
		return nil, err
	}

	cfg, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &config{}
	}

	warnings := []string{}
	acctTtl := time.Duration(acct.Ttl) * time.Second
	if acctTtl > cfg.MaxTTL {
		warnings = append(warnings, fmt.Sprintf("using backend max ttl %q which is less than impersonated account ttl %q for token",
			cfg.MaxTTL.String(),
			acctTtl.String()))
		acctTtl = cfg.MaxTTL
	} else if acctTtl == 0 {
		warnings = append(warnings, fmt.Sprintf("using backend default ttl %q since impersonated account ttl not configured for token",
			cfg.TTL.String()))
		acctTtl = cfg.TTL
	}

	tokenSource, err := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
		TargetPrincipal: acct.EmailOrId,
		Scopes:          acct.TokenScopes,
		Lifetime:        acctTtl,
	}, option.WithCredentials(creds))
	if err != nil {
		return logical.ErrorResponse("unable to generate token source: %v", err), nil
	}
	token, err := tokenSource.Token()
	if err != nil {
		return logical.ErrorResponse("unable to generate token - make sure your service account and key are still valid: %v", err), nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"token":              token.AccessToken,
			"token_ttl":          token.Expiry.UTC().Sub(time.Now().UTC()) / (time.Second),
			"expires_at_seconds": token.Expiry.Unix(),
		},
		Warnings: warnings,
	}, nil
}
