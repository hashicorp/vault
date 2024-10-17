// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package okta

import (
	"context"
	"fmt"
	"net/textproto"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/okta/okta-sdk-golang/v4/okta"
	"github.com/patrickmn/go-cache"
)

const (
	operationPrefixOkta = "okta"
	mfaPushMethod       = "push"
	mfaTOTPMethod       = "token:software:totp"
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
			Unauthenticated: []string{
				"login/*",
				"verify/*",
			},
			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: []*framework.Path{
			pathConfig(&b),
			pathUsers(&b),
			pathGroups(&b),
			pathUsersList(&b),
			pathGroupsList(&b),
			pathLogin(&b),
			pathVerify(&b),
		},

		AuthRenew:   b.pathLoginRenew,
		BackendType: logical.TypeCredential,
	}
	b.verifyCache = cache.New(5*time.Minute, time.Minute)

	return &b
}

type backend struct {
	*framework.Backend
	verifyCache *cache.Cache
}

func (b *backend) Login(ctx context.Context, req *logical.Request, username, password, totp, nonce, preferredProvider string) ([]string, *logical.Response, []string, error) {
	cfg, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, nil, nil, err
	}
	if cfg == nil {
		return nil, logical.ErrorResponse("Okta auth method not configured"), nil, nil
	}

	// Check for a CIDR match.
	if len(cfg.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, nil, nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, cfg.TokenBoundCIDRs) {
			return nil, nil, nil, logical.ErrPermissionDenied
		}
	}

	shim, err := cfg.OktaClient(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	type mfaFactor struct {
		Id       string `json:"id"`
		Type     string `json:"factorType"`
		Provider string `json:"provider"`
		Embedded struct {
			Challenge struct {
				CorrectAnswer *int `json:"correctAnswer"`
			} `json:"challenge"`
		} `json:"_embedded"`
	}

	type embeddedResult struct {
		User    okta.User   `json:"user"`
		Factors []mfaFactor `json:"factors"`
		Factor  *mfaFactor  `json:"factor"`
	}

	type authResult struct {
		Embedded     embeddedResult `json:"_embedded"`
		Status       string         `json:"status"`
		FactorResult string         `json:"factorResult"`
		StateToken   string         `json:"stateToken"`
	}

	authReq, err := shim.NewRequest("POST", "authn", map[string]interface{}{
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	var result authResult
	rsp, err := shim.Do(authReq, &result)
	if err != nil {
		if oe, ok := err.(*okta.Error); ok {
			return nil, logical.ErrorResponse("Okta auth failed: %v (code=%v)", err, oe.ErrorCode), nil, nil
		}
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
		//
		// API reference: https://developer.okta.com/docs/reference/api/authn/#verify-factor
		if cfg.BypassOktaMFA {
			result.Status = "SUCCESS"
			break
		}

		var selectedFactor, totpFactor, pushFactor *mfaFactor

		// Scan for available factors
		for _, v := range result.Embedded.Factors {
			v := v // create a new copy since we'll be taking the address later

			if preferredProvider != "" && preferredProvider != v.Provider {
				continue
			}

			if !strutil.StrListContains(b.getSupportedProviders(), v.Provider) {
				continue
			}

			switch v.Type {
			case mfaTOTPMethod:
				totpFactor = &v
			case mfaPushMethod:
				pushFactor = &v
			}
		}

		// Okta push and totp, and Google totp are currently supported.
		// If a totp passcode is provided during login and is supported,
		// that will be the preferred method.
		switch {
		case totpFactor != nil && totp != "":
			selectedFactor = totpFactor
		case pushFactor != nil && pushFactor.Provider == oktaProvider:
			selectedFactor = pushFactor
		case totpFactor != nil && totp == "":
			return nil, logical.ErrorResponse("'totp' passcode parameter is required to perform MFA"), nil, nil
		default:
			return nil, logical.ErrorResponse("Okta Verify Push or TOTP or Google TOTP factor is required in order to perform MFA"), nil, nil
		}

		requestPath := fmt.Sprintf("authn/factors/%s/verify", selectedFactor.Id)

		payload := map[string]interface{}{
			"stateToken": result.StateToken,
		}
		if selectedFactor.Type == mfaTOTPMethod {
			payload["passCode"] = totp
		}

		verifyReq, err := shim.NewRequest("POST", requestPath, payload)
		if err != nil {
			return nil, nil, nil, err
		}
		if len(req.Headers["X-Forwarded-For"]) > 0 {
			verifyReq.Header.Set("X-Forwarded-For", req.Headers[textproto.CanonicalMIMEHeaderKey("X-Forwarded-For")][0])
		}

		rsp, err := shim.Do(verifyReq, &result)
		if err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("Okta auth failed: %v", err)), nil, nil
		}
		if rsp == nil {
			return nil, logical.ErrorResponse("okta auth backend unexpected failure"), nil, nil
		}
		for result.Status == "MFA_CHALLENGE" {
			switch result.FactorResult {
			case "WAITING":
				verifyReq, err := shim.NewRequest("POST", requestPath, payload)
				if err != nil {
					return nil, logical.ErrorResponse(fmt.Sprintf("okta auth failed creating verify request: %v", err)), nil, nil
				}
				rsp, err := shim.Do(verifyReq, &result)

				// Store number challenge if found
				numberChallenge := result.Embedded.Factor.Embedded.Challenge.CorrectAnswer
				if numberChallenge != nil {
					if nonce == "" {
						return nil, logical.ErrorResponse("nonce must be provided during login request when presented with number challenge"), nil, nil
					}

					b.verifyCache.SetDefault(nonce, *numberChallenge)
				}

				if err != nil {
					return nil, logical.ErrorResponse(fmt.Sprintf("Okta auth failed checking loop: %v", err)), nil, nil
				}
				if rsp == nil {
					return nil, logical.ErrorResponse("okta auth backend unexpected failure"), nil, nil
				}

				timer := time.NewTimer(1 * time.Second)
				select {
				case <-timer.C:
					// Continue
				case <-ctx.Done():
					timer.Stop()
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
	client, oktactx := shim.Client()
	if client != nil {
		oktaGroups, err := b.getOktaGroups(oktactx, client, &result.Embedded.User)
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

func (b *backend) getOktaGroups(ctx context.Context, client *okta.Client, user *okta.User) ([]string, error) {
	groups, resp, err := client.User.ListUserGroups(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	oktaGroups := make([]string, 0, len(groups))
	for _, group := range groups {
		oktaGroups = append(oktaGroups, group.Profile.Name)
	}
	for resp.HasNextPage() {
		var nextGroups []*okta.Group
		resp, err = resp.Next(ctx, &nextGroups)
		if err != nil {
			return nil, err
		}
		for _, group := range nextGroups {
			oktaGroups = append(oktaGroups, group.Profile.Name)
		}
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
