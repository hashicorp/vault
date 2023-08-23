// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package okta

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	googleProvider = "GOOGLE"
	oktaProvider   = "OKTA"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login/(?P<username>.+)`,

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixOkta,
			OperationVerb:   "login",
		},

		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "Username to be used for login.",
			},

			"password": {
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},
			"totp": {
				Type:        framework.TypeString,
				Description: "TOTP passcode.",
			},
			"nonce": {
				Type: framework.TypeString,
				Description: `Nonce provided if performing login that requires 
number verification challenge. Logins through the vault login CLI command will 
automatically generate a nonce.`,
			},
			"provider": {
				Type:        framework.TypeString,
				Description: "Preferred factor provider.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.pathLoginAliasLookahead,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) getSupportedProviders() []string {
	return []string{googleProvider, oktaProvider}
}

func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: username,
			},
		},
	}, nil
}

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	totp := d.Get("totp").(string)
	nonce := d.Get("nonce").(string)
	preferredProvider := strings.ToUpper(d.Get("provider").(string))
	if preferredProvider != "" && !strutil.StrListContains(b.getSupportedProviders(), preferredProvider) {
		return logical.ErrorResponse(fmt.Sprintf("provider %s is not among the supported ones %v", preferredProvider, b.getSupportedProviders())), nil
	}

	defer b.verifyCache.Delete(nonce)

	policies, resp, groupNames, err := b.Login(ctx, req, username, password, totp, nonce, preferredProvider)
	// Handle an internal error
	if err != nil {
		return nil, err
	}
	if resp != nil {
		// Handle a logical error
		if resp.IsError() {
			return resp, nil
		}
	} else {
		resp = &logical.Response{}
	}

	cfg, err := b.getConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	auth := &logical.Auth{
		Metadata: map[string]string{
			"username": username,
			"policies": strings.Join(policies, ","),
		},
		InternalData: map[string]interface{}{
			"password": password,
		},
		DisplayName: username,
		Alias: &logical.Alias{
			Name: username,
		},
	}
	cfg.PopulateTokenAuth(auth)

	// Add in configured policies from mappings
	if len(policies) > 0 {
		auth.Policies = append(auth.Policies, policies...)
	}

	resp.Auth = auth

	for _, groupName := range groupNames {
		if groupName == "" {
			continue
		}
		resp.Auth.GroupAliases = append(resp.Auth.GroupAliases, &logical.Alias{
			Name: groupName,
		})
	}

	return resp, nil
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := req.Auth.Metadata["username"]
	password := req.Auth.InternalData["password"].(string)

	var nonce string
	if d != nil {
		nonce = d.Get("nonce").(string)
	}

	cfg, err := b.getConfig(ctx, req)
	if err != nil {
		return nil, err
	}

	// No TOTP entry is possible on renew. If push MFA is enabled it will still be triggered, however.
	// Sending "" as the totp will prompt the push action if it is configured.
	loginPolicies, resp, groupNames, err := b.Login(ctx, req, username, password, "", nonce, "")
	if err != nil || (resp != nil && resp.IsError()) {
		return resp, err
	}

	finalPolicies := cfg.TokenPolicies
	if len(loginPolicies) > 0 {
		finalPolicies = append(finalPolicies, loginPolicies...)
	}
	if !policyutil.EquivalentPolicies(finalPolicies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	resp.Auth = req.Auth
	resp.Auth.Period = cfg.TokenPeriod
	resp.Auth.TTL = cfg.TokenTTL
	resp.Auth.MaxTTL = cfg.TokenMaxTTL

	// Remove old aliases
	resp.Auth.GroupAliases = nil

	for _, groupName := range groupNames {
		resp.Auth.GroupAliases = append(resp.Auth.GroupAliases, &logical.Alias{
			Name: groupName,
		})
	}

	return resp, nil
}

func pathVerify(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `verify/(?P<nonce>.+)`,
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixOkta,
			OperationVerb:   "verify",
		},
		Fields: map[string]*framework.FieldSchema{
			"nonce": {
				Type: framework.TypeString,
				Description: `Nonce provided during a login request to
retrieve the number verification challenge for the matching request.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathVerify,
			},
		},
	}
}

func (b *backend) pathVerify(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	nonce := d.Get("nonce").(string)

	correctRaw, ok := b.verifyCache.Get(nonce)
	if !ok {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"correct_answer": correctRaw.(int),
		},
	}

	return resp, nil
}

func (b *backend) getConfig(ctx context.Context, req *logical.Request) (*ConfigEntry, error) {
	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("Okta backend not configured")
	}

	return cfg, nil
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
