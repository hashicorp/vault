package proxy

import (
	"context"
	"fmt"
	"net/textproto"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
	glob "github.com/ryanuber/go-glob"
)

const (
	loginHelpSyn   = `Authenticate to vault using credentials supplied from a trusted proxy`
	loginRoleField = "role"
	loginNameField = "name"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login$`,
		Fields: map[string]*framework.FieldSchema{
			loginRoleField: {
				Type:        framework.TypeLowerCaseString,
				Description: "The role to login against.",
			},

			// support 'name' as an alias for 'role'.  This allows the login
			// method of the cert engine in the default vault client to work
			// with this backend
			loginNameField: {
				Type:        framework.TypeLowerCaseString,
				Description: fmt.Sprintf("alias for %q field", loginRoleField),
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathLogin,
				Summary:  loginHelpSyn,
			},
			logical.AliasLookaheadOperation: &framework.PathOperation{
				Callback: b.pathLogin,
			},
		},

		HelpSynopsis: loginHelpSyn,
	}
}

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("could not load configuration"), nil
	}

	rolename := data.Get(loginRoleField).(string)
	if rolename == "" {
		rolename = data.Get(loginNameField).(string)
		if rolename == "" {
			return logical.ErrorResponse(fmt.Sprintf("either %q or %q must be set", loginRoleField, loginNameField)), nil
		}
	}

	role, err := b.getRole(ctx, req.Storage, rolename)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("role %q could not be found", rolename), nil
	}

	if req.Headers == nil {
		// a likely cause of this error is when the proxy authentication
		// method was enabled the "passthrough_request_headers" option
		// was not specified.  This option should include the configured
		// user header
		return nil, fmt.Errorf("failed to get request headers")
	}

	username, err := getHeaderVal(config.UserHeader, req)
	if err != nil || username == "" {
		return nil, fmt.Errorf("could not identify remote user")
	}

	if req.Operation == logical.AliasLookaheadOperation {
		return &logical.Response{
			Auth: &logical.Auth{
				Alias: b.getAuthAlias(username, rolename),
			},
		}, nil
	}

	if allowed, resp, err := b.authorize(config, role, req, username); !allowed {
		return resp, err
	}

	// grant the token
	resp := logical.Response{
		Auth: &logical.Auth{
			DisplayName: username,
			Period:      role.Period,
			Policies:    role.Policies,

			Metadata: map[string]string{
				"username": username,
				"role":     rolename,
			},
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       role.TTL,
				MaxTTL:    role.MaxTTL,
			},
			Alias:      b.getAuthAlias(username, rolename),
			BoundCIDRs: config.BoundCIDRs,
		},
	}

	return &resp, nil
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("could not load configuration"), nil
	}

	rolename := req.Auth.Metadata["role"]
	role, err := b.getRole(ctx, req.Storage, rolename)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("role %q could not be found during renewal", rolename), nil
	}

	username := req.Auth.Metadata["username"]
	if allowed, resp, err := b.authorize(config, role, req, username); !allowed {
		return resp, err
	}

	if !policyutil.EquivalentPolicies(role.Policies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.TTL = role.TTL
	resp.Auth.MaxTTL = role.MaxTTL
	return resp, nil
}

func (b *backend) getAuthAlias(username, role string) *logical.Alias {
	return &logical.Alias{
		Name: username,
		Metadata: map[string]string{
			"role": role,
		},
	}
}

func (b *backend) hasRequiredUser(username string, role *proxyRole) bool {
	isAllowedUser := false
	for _, allowedUserGlob := range role.AllowedUsers {
		if glob.Glob(allowedUserGlob, username) {
			isAllowedUser = true
			break
		}
	}

	return isAllowedUser
}

func (b *backend) hasRequiredHeaders(role *proxyRole, req *logical.Request) bool {
	for reqHdrName, reqHdrVal := range role.RequiredHeaders {
		val, err := getHeaderVal(reqHdrName, req)
		if err != nil || val != reqHdrVal {
			return false
		}
	}

	return true
}

func (b *backend) hasRequiredRemoteAddr(config *proxyConfig, req *logical.Request) (bool, error) {
	if len(config.BoundCIDRs) != 0 {
		if req.Connection == nil || req.Connection.RemoteAddr == "" {
			return false, fmt.Errorf("failed to get connection information")
		}

		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, config.BoundCIDRs) {
			return false, nil
		}
	}

	return true, nil
}

func (b *backend) authorize(config *proxyConfig, role *proxyRole, req *logical.Request, username string) (allowed bool, resp *logical.Response, err error) {
	if !b.hasRequiredUser(username, role) {
		return false, logical.ErrorResponse("user not permitted to authenticate with this role"), nil
	}

	if !b.hasRequiredHeaders(role, req) {
		return false, logical.ErrorResponse("required header not present, or has incorrect value"), nil
	}

	if ok, err := b.hasRequiredRemoteAddr(config, req); err != nil {
		return false, nil, err
	} else if !ok {
		return false, logical.ErrorResponse("unauthorized source address"), nil
	}

	return true, nil, nil
}

func getHeaderVal(headerName string, req *logical.Request) (string, error) {
	canonHeaderName := textproto.CanonicalMIMEHeaderKey(headerName)
	for k, v := range req.Headers {
		if canonHeaderName == textproto.CanonicalMIMEHeaderKey(k) {
			if len(v) == 1 {
				return v[0], nil
			}

			return "", fmt.Errorf("header present %d times", len(v))
		}

	}

	return "", nil
}
