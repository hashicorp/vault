package okta

import (
	"context"
	"fmt"

	"github.com/chrismalek/oktasdk-go/okta"
	"github.com/hashicorp/vault/helper/mfa"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Root: mfa.MFARootPaths(),

			Unauthenticated: []string{
				"login/*",
			},
			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathUsers(&b),
			pathGroups(&b),
			pathUsersList(&b),
			pathGroupsList(&b),
		},
			mfa.MFAPaths(b.Backend, pathLogin(&b))...,
		),

		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
	}

	return &b
}

type backend struct {
	*framework.Backend
}

func (b *backend) Login(ctx context.Context, req *logical.Request, username string, password string) ([]string, *logical.Response, []string, error) {
	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, nil, nil, err
	}
	if cfg == nil {
		return nil, logical.ErrorResponse("Okta auth method not configured"), nil, nil
	}

	client := cfg.OktaClient()

	type embeddedResult struct {
		User okta.User `json:"user"`
	}

	type authResult struct {
		Embedded embeddedResult `json:"_embedded"`
		Status   string         `json:"status"`
	}

	authReq, err := client.NewRequest("POST", "authn", map[string]interface{}{
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	var result authResult
	rsp, err := client.Do(authReq, &result)
	if err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("Okta auth failed: %v", err)), nil, nil
	}
	if rsp == nil {
		return nil, logical.ErrorResponse("okta auth method unexpected failure"), nil, nil
	}

	oktaResponse := &logical.Response{
		Data: map[string]interface{}{},
	}

	// If lockout failures are not configured to be hidden, the status needs to
	// be inspected for LOCKED_OUT status. Otherwise, it is handled above by an
	// error returned during the authentication request.
	switch result.Status {
	case "LOCKED_OUT":
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: user is locked out", "user", username)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil

	case "PASSWORD_EXPIRED":
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: password is expired", "user", username)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil

	case "PASSWORD_WARN":
		oktaResponse.AddWarning("Your Okta password is in warning state and needs to be changed soon.")

	case "MFA_REQUIRED", "MFA_ENROLL":
		if !cfg.BypassOktaMFA {
			return nil, logical.ErrorResponse("okta mfa required for this account but mfa bypass not set in config"), nil, nil
		}

	case "SUCCESS":
		// Do nothing here

	default:
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: unhandled result status", "status", result.Status)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil
	}

	// Verify result status again in case a switch case above modifies result
	switch {
	case result.Status == "SUCCESS",
		result.Status == "PASSWORD_WARN",
		result.Status == "MFA_REQUIRED" && cfg.BypassOktaMFA,
		result.Status == "MFA_ENROLL" && cfg.BypassOktaMFA:
		// Allowed
	default:
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: authentication returned a non-success status", "status", result.Status)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil
	}

	var allGroups []string
	// Only query the Okta API for group membership if we have a token
	if cfg.Token != "" {
		oktaGroups, err := b.getOktaGroups(client, &result.Embedded.User)
		if err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("okta failure retrieving groups: %v", err)), nil, nil
		}
		if len(oktaGroups) == 0 {
			errString := fmt.Sprintf(
				"no Okta groups found; only policies from locally-defined groups available")
			oktaResponse.AddWarning(errString)
		}
		allGroups = append(allGroups, oktaGroups...)
	}

	// Import the custom added groups from okta backend
	user, err := b.User(ctx, req.Storage, username)
	if err != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: error looking up user", "error", err)
		}
	}
	if err == nil && user != nil && user.Groups != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("auth/okta: adding local groups", "num_local_groups", len(user.Groups), "local_groups", user.Groups)
		}
		allGroups = append(allGroups, user.Groups...)
	}

	// Retrieve policies
	var policies []string
	for _, groupName := range allGroups {
		entry, _, err := b.Group(ctx, req.Storage, groupName)
		if err != nil {
			if b.Logger().IsDebug() {
				b.Logger().Debug("auth/okta: error looking up group policies", "error", err)
			}
		}
		if err == nil && entry != nil && entry.Policies != nil {
			policies = append(policies, entry.Policies...)
		}
	}

	// Merge local Policies into Okta Policies
	if user != nil && user.Policies != nil {
		policies = append(policies, user.Policies...)
	}

	if len(policies) == 0 {
		errStr := "user is not a member of any authorized policy"
		if len(oktaResponse.Warnings) > 0 {
			errStr = fmt.Sprintf("%s; additionally, %s", errStr, oktaResponse.Warnings[0])
		}

		oktaResponse.Data["error"] = errStr
		return nil, oktaResponse, nil, nil
	}

	return policies, oktaResponse, allGroups, nil
}

func (b *backend) getOktaGroups(client *okta.Client, user *okta.User) ([]string, error) {
	rsp, err := client.Users.PopulateGroups(user)
	if err != nil {
		return nil, err
	}
	if rsp == nil {
		return nil, fmt.Errorf("okta auth method unexpected failure")
	}
	oktaGroups := make([]string, 0, len(user.Groups))
	for _, group := range user.Groups {
		oktaGroups = append(oktaGroups, group.Profile.Name)
	}
	if b.Logger().IsDebug() {
		b.Logger().Debug("auth/okta: Groups fetched from Okta", "num_groups", len(oktaGroups), "groups", fmt.Sprintf("%#v", oktaGroups))
	}
	return oktaGroups, nil
}

const backendHelp = `
The Okta credential provider allows authentication querying,
checking username and password, and associating policies.  If an api token is configure
groups are pulled down from Okta.

Configuration of the connection is done through the "config" and "policies"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
