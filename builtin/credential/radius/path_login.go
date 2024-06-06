// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package radius

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
	"layeh.com/radius"
	. "layeh.com/radius/rfc2865"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login" + framework.OptionalParamRegex("urlusername"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixRadius,
			OperationVerb:   "login",
			OperationSuffix: "|with-username",
		},

		Fields: map[string]*framework.FieldSchema{
			"urlusername": {
				Type:        framework.TypeString,
				Description: "Username to be used for login. (URL parameter)",
			},

			"username": {
				Type:        framework.TypeString,
				Description: "Username to be used for login. (POST request body)",
			},

			"password": {
				Type:        framework.TypeString,
				Description: "Password for this user.",
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
	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("radius backend not configured"), nil
	}

	// Check for a CIDR match.
	if len(cfg.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, cfg.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}
	}

	username := d.Get("username").(string)
	password := d.Get("password").(string)

	if username == "" {
		username = d.Get("urlusername").(string)
		if username == "" {
			return logical.ErrorResponse("username cannot be empty"), nil
		}
	}

	if password == "" {
		return logical.ErrorResponse("password cannot be empty"), nil
	}

	policies, resp, err := b.RadiusLogin(ctx, req, username, password)
	// Handle an internal error
	if err != nil {
		return nil, err
	}
	if resp != nil {
		// Handle a logical error
		if resp.IsError() {
			return resp, nil
		}
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

	resp.Auth = auth
	if policies != nil {
		resp.Auth.Policies = append(resp.Auth.Policies, policies...)
	}

	return resp, nil
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("radius backend not configured"), nil
	}

	username := req.Auth.Metadata["username"]
	password := req.Auth.InternalData["password"].(string)

	var resp *logical.Response
	var loginPolicies []string

	loginPolicies, resp, err = b.RadiusLogin(ctx, req, username, password)
	if err != nil || (resp != nil && resp.IsError()) {
		return resp, err
	}
	finalPolicies := cfg.TokenPolicies
	if loginPolicies != nil {
		finalPolicies = append(finalPolicies, loginPolicies...)
	}

	if !policyutil.EquivalentPolicies(finalPolicies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	req.Auth.Period = cfg.TokenPeriod
	req.Auth.TTL = cfg.TokenTTL
	req.Auth.MaxTTL = cfg.TokenMaxTTL
	return &logical.Response{Auth: req.Auth}, nil
}

func (b *backend) RadiusLogin(ctx context.Context, req *logical.Request, username string, password string) ([]string, *logical.Response, error) {
	cfg, err := b.Config(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	if cfg == nil || cfg.Host == "" || cfg.Secret == "" {
		return nil, logical.ErrorResponse("radius backend not configured"), nil
	}

	hostport := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	packet := radius.New(radius.CodeAccessRequest, []byte(cfg.Secret))
	UserName_SetString(packet, username)
	UserPassword_SetString(packet, password)
	if cfg.NasIdentifier != "" {
		NASIdentifier_AddString(packet, cfg.NasIdentifier)
	}
	packet.Add(5, radius.NewInteger(uint32(cfg.NasPort)))

	client := radius.Client{
		Dialer: net.Dialer{
			Timeout: time.Duration(cfg.DialTimeout) * time.Second,
		},
	}
	clientCtx, cancelFunc := context.WithTimeout(ctx, time.Duration(cfg.ReadTimeout)*time.Second)
	received, err := client.Exchange(clientCtx, packet, hostport)
	cancelFunc()
	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}
	if received.Code != radius.CodeAccessAccept {
		return nil, logical.ErrorResponse("access denied by the authentication server"), nil
	}

	policies := cfg.UnregisteredUserPolicies

	// Retrieve user entry from storage
	user, err := b.user(ctx, req.Storage, username)
	if err != nil {
		return nil, logical.ErrorResponse("could not retrieve user entry from storage"), err
	}
	if user != nil {
		policies = user.Policies
	}

	return policies, &logical.Response{}, nil
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password. Please be sure to
read the note on escaping from the path-help for the 'config' endpoint.
`
