package okta

import (
	"context"
	"fmt"
	"time"

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

	type mfaFactor struct {
		Id       string `json:"id"`
		Type     string `json:"factorType"`
		Provider string `json:"provider"`
	}

	type embeddedResult struct {
		User    okta.User   `json:"user"`
		Factors []mfaFactor `json:"factors"`
	}

	type authResult struct {
		Embedded     embeddedResult `json:"_embedded"`
		Status       string         `json:"status"`
		FactorResult string         `json:"factorResult"`
		StateToken   string         `json:"stateToken"`
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

	// More about Okta's Auth transaction state here:
	// https://developer.okta.com/docs/api/resources/authn#transaction-state

	// If lockout failures are not configured to be hidden, the status needs to
	// be inspected for LOCKED_OUT status. Otherwise, it is handled above by an
	// error returned during the authentication request.
	switch result.Status {
	case "LOCKED_OUT":
		if b.Logger().IsDebug() {
			b.Logger().Debug("user is locked out", "user", username)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil

	case "PASSWORD_EXPIRED":
		if b.Logger().IsDebug() {
			b.Logger().Debug("password is expired", "user", username)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil

	case "PASSWORD_WARN":
		oktaResponse.AddWarning("Your Okta password is in warning state and needs to be changed soon.")

	case "MFA_ENROLL", "MFA_ENROLL_ACTIVATE":
		if !cfg.BypassOktaMFA {
			if b.Logger().IsDebug() {
				b.Logger().Debug("user must enroll or complete mfa enrollment", "user", username)
			}
			return nil, logical.ErrorResponse("okta authentication failed: you must complete MFA enrollment to continue"), nil, nil
		}

	case "MFA_REQUIRED":
		// Per Okta documentation: Users are challenged for MFA (MFA_REQUIRED)
		// before the Status of PASSWORD_EXPIRED is exposed (if they have an
		// active factor enrollment). This bypass removes visibility
		// into the authenticating user's password expiry, but still ensures the
		// credentials are valid and the user is not locked out.
		if cfg.BypassOktaMFA {
			result.Status = "SUCCESS"
			break
		}

		factorAvailable := false

		var selectedFactor mfaFactor
		// only okta push is currently supported
		for _, v := range result.Embedded.Factors {
			if v.Type == "push" && v.Provider == "OKTA" {
				factorAvailable = true
				selectedFactor = v
			}
		}

		if !factorAvailable {
			return nil, logical.ErrorResponse("Okta Verify Push factor is required in order to perform MFA"), nil, nil
		}

		requestPath := fmt.Sprintf("authn/factors/%s/verify", selectedFactor.Id)
		payload := map[string]interface{}{
			"stateToken": result.StateToken,
		}
		verifyReq, err := client.NewRequest("POST", requestPath, payload)
		if err != nil {
			return nil, nil, nil, err
		}

		rsp, err := client.Do(verifyReq, &result)
		if err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("Okta auth failed: %v", err)), nil, nil
		}
		if rsp == nil {
			return nil, logical.ErrorResponse("okta auth backend unexpected failure"), nil, nil
		}
		for result.Status == "MFA_CHALLENGE" {
			switch result.FactorResult {
			case "WAITING":
				verifyReq, err := client.NewRequest("POST", requestPath, payload)
				rsp, err := client.Do(verifyReq, &result)
				if err != nil {
					return nil, logical.ErrorResponse(fmt.Sprintf("Okta auth failed checking loop: %v", err)), nil, nil
				}
				if rsp == nil {
					return nil, logical.ErrorResponse("okta auth backend unexpected failure"), nil, nil
				}

				select {
				case <-time.After(500 * time.Millisecond):
					// Continue
				case <-ctx.Done():
					return nil, logical.ErrorResponse("exiting pending mfa challenge"), nil, nil
				}
			case "REJECTED":
				return nil, logical.ErrorResponse("multi-factor authentication denied"), nil, nil
			case "TIMEOUT":
				return nil, logical.ErrorResponse("failed to complete multi-factor authentication"), nil, nil
			case "SUCCESS":
				// Allowed
			default:
				if b.Logger().IsDebug() {
					b.Logger().Debug("unhandled result status", "status", result.Status, "factorstatus", result.FactorResult)
				}
				return nil, logical.ErrorResponse("okta authentication failed"), nil, nil
			}
		}

	case "SUCCESS":
		// Do nothing here

	default:
		if b.Logger().IsDebug() {
			b.Logger().Debug("unhandled result status", "status", result.Status)
		}
		return nil, logical.ErrorResponse("okta authentication failed"), nil, nil
	}

	// Verify result status again in case a switch case above modifies result
	switch {
	case result.Status == "SUCCESS",
		result.Status == "PASSWORD_WARN",
		result.Status == "MFA_REQUIRED" && cfg.BypassOktaMFA,
		result.Status == "MFA_ENROLL" && cfg.BypassOktaMFA,
		result.Status == "MFA_ENROLL_ACTIVATE" && cfg.BypassOktaMFA:
		// Allowed
	default:
		if b.Logger().IsDebug() {
			b.Logger().Debug("authentication returned a non-success status", "status", result.Status)
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
			b.Logger().Debug("error looking up user", "error", err)
		}
	}
	if err == nil && user != nil && user.Groups != nil {
		if b.Logger().IsDebug() {
			b.Logger().Debug("adding local groups", "num_local_groups", len(user.Groups), "local_groups", user.Groups)
		}
		allGroups = append(allGroups, user.Groups...)
	}

	// Retrieve policies
	var policies []string
	for _, groupName := range allGroups {
		entry, _, err := b.Group(ctx, req.Storage, groupName)
		if err != nil {
			if b.Logger().IsDebug() {
				b.Logger().Debug("error looking up group policies", "error", err)
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
		b.Logger().Debug("Groups fetched from Okta", "num_groups", len(oktaGroups), "groups", fmt.Sprintf("%#v", oktaGroups))
	}
	return oktaGroups, nil
}

const backendHelp = `
The Okta credential provider allows authentication querying,
checking username and password, and associating policies.  If an api token is
configured groups are pulled down from Okta.

Configuration of the connection is done through the "config" and "policies"
endpoints by a user with root access. Authentication is then done
by supplying the two fields for "login".
`
