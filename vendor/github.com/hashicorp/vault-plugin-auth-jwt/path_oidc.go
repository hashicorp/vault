package jwtauth

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/vault/helper/strutil"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"golang.org/x/oauth2"
)

var oidcStateTimeout = 10 * time.Minute

// OIDC error prefixes. These are searched for specifically by the UI, so any
// changes to them must be aligned with a UI change.
const errLoginFailed = "Vault login failed."
const errNoResponse = "No response from provider."
const errTokenVerification = "Token verification failed."

// oidcState is created when an authURL is requested. The state identifier is
// passed throughout the OAuth process.
type oidcState struct {
	rolename    string
	nonce       string
	redirectURI string
}

func pathOIDC(b *jwtAuthBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: `oidc/callback`,
			Fields: map[string]*framework.FieldSchema{
				"state": {
					Type: framework.TypeString,
				},
				"code": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathCallback,
					Summary:  "Callback endpoint to complete an OIDC login.",
				},
			},
		},
		{
			Pattern: `oidc/auth_url`,
			Fields: map[string]*framework.FieldSchema{
				"role": {
					Type:        framework.TypeLowerCaseString,
					Description: "The role to issue an OIDC authorization URL against.",
				},
				"redirect_uri": {
					Type:        framework.TypeString,
					Description: "The OAuth redirect_uri to use in the authorization URL.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.authURL,
					Summary:  "Request an authorization URL to start an OIDC login flow.",
				},
			},
		},
	}
}

func (b *jwtAuthBackend) pathCallback(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	state := b.verifyState(d.Get("state").(string))
	if state == nil {
		return logical.ErrorResponse(errLoginFailed + " Expired or missing OAuth state."), nil
	}

	roleName := state.rolename
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(errLoginFailed + " Role could not be found"), nil
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse(errLoginFailed + " Could not load configuration"), nil
	}

	provider, err := b.getProvider(ctx, config)
	if err != nil {
		return nil, errwrap.Wrapf(errLoginFailed+" Error getting provider for login operation: {{err}}", err)
	}

	var oauth2Config = oauth2.Config{
		ClientID:     config.OIDCClientID,
		ClientSecret: config.OIDCClientSecret,
		RedirectURL:  state.redirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}

	code := d.Get("code").(string)
	if code == "" {
		return logical.ErrorResponse(errLoginFailed + " OAuth code parameter not provided"), nil
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return logical.ErrorResponse(errLoginFailed+" Error exchanging oidc code: %q.", err.Error()), nil
	}

	// Extract the ID Token from OAuth2 token.
	rawToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return logical.ErrorResponse(errTokenVerification + " No id_token found in response."), nil
	}

	// Parse and verify ID Token payload.
	allClaims, err := b.verifyOIDCToken(ctx, config, role, rawToken)
	if err != nil {
		return logical.ErrorResponse("%s %s", errTokenVerification, err.Error()), nil
	}

	if allClaims["nonce"] != state.nonce {
		return logical.ErrorResponse(errTokenVerification + " Invalid ID token nonce."), nil
	}
	delete(allClaims, "nonce")

	// Attempt to fetch information from the /userinfo endpoint and merge it with
	// the existing claims data. A failure to fetch additional information from this
	// endpoint will not invalidate the authorization flow.
	if userinfo, err := provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token)); err == nil {
		_ = userinfo.Claims(&allClaims)
	} else {
		logFunc := b.Logger().Warn
		if strings.Contains(err.Error(), "user info endpoint is not supported") {
			logFunc = b.Logger().Info
		}
		logFunc("error reading /userinfo endpoint", "error", err)
	}

	if err := validateBoundClaims(b.Logger(), role.BoundClaims, allClaims); err != nil {
		return logical.ErrorResponse("error validating claims: %s", err.Error()), nil
	}

	alias, groupAliases, err := b.createIdentity(allClaims, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	tokenMetadata := map[string]string{"role": roleName}
	for k, v := range alias.Metadata {
		tokenMetadata[k] = v
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			Policies:     role.Policies,
			DisplayName:  alias.Name,
			Period:       role.Period,
			NumUses:      role.NumUses,
			Alias:        alias,
			GroupAliases: groupAliases,
			InternalData: map[string]interface{}{
				"role": roleName,
			},
			Metadata: tokenMetadata,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       role.TTL,
				MaxTTL:    role.MaxTTL,
			},
			BoundCIDRs: role.BoundCIDRs,
		},
	}

	return resp, nil
}

// authURL returns a URL used for redirection to receive an authorization code.
// This path requires a role name, or that a default_role has been configured.
// Because this endpoint is unauthenticated, the response to invalid or non-OIDC
// roles is intentionally non-descriptive and will simply be an empty string.
func (b *jwtAuthBackend) authURL(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	logger := b.Logger()

	// default response for most error/invalid conditions
	resp := &logical.Response{
		Data: map[string]interface{}{
			"auth_url": "",
		},
	}

	config, err := b.config(ctx, req.Storage)
	if err != nil {
		logger.Warn("error loading configuration", "error", err)
		return resp, nil
	}

	if config == nil {
		logger.Warn("nil configuration")
		return resp, nil
	}

	roleName := d.Get("role").(string)
	if roleName == "" {
		roleName = config.DefaultRole
	}
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	redirectURI := d.Get("redirect_uri").(string)
	if redirectURI == "" {
		return logical.ErrorResponse("missing redirect_uri"), nil
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return resp, nil
	}

	if role == nil || role.RoleType != "oidc" {
		return resp, nil
	}

	if !validRedirect(redirectURI, role.AllowedRedirectURIs) {
		logger.Warn("unauthorized redirect_uri", "redirect_uri", redirectURI)
		return resp, nil
	}

	provider, err := b.getProvider(ctx, config)
	if err != nil {
		logger.Warn("error getting provider for login operation", "error", err)
		return resp, nil
	}

	// "openid" is a required scope for OpenID Connect flows
	scopes := append([]string{oidc.ScopeOpenID}, role.OIDCScopes...)

	// Configure an OpenID Connect aware OAuth2 client
	oauth2Config := oauth2.Config{
		ClientID:     config.OIDCClientID,
		ClientSecret: config.OIDCClientSecret,
		RedirectURL:  redirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	stateID, nonce, err := b.createState(roleName, redirectURI)
	if err != nil {
		logger.Warn("error generating OAuth state", "error", err)
		return resp, nil
	}

	resp.Data["auth_url"] = oauth2Config.AuthCodeURL(stateID, oidc.Nonce(nonce))

	return resp, nil
}

// createState make an expiring state object, associated with a random state ID
// that is passed throughout the OAuth process. A nonce is also included in the
// auth process, and for simplicity will be identical in length/format as the state ID.
func (b *jwtAuthBackend) createState(rolename, redirectURI string) (string, string, error) {
	// Get enough bytes for 2 160-bit IDs (per rfc6749#section-10.10)
	bytes, err := uuid.GenerateRandomBytes(2 * 20)
	if err != nil {
		return "", "", err
	}

	stateID := fmt.Sprintf("%x", bytes[:20])
	nonce := fmt.Sprintf("%x", bytes[20:])

	b.oidcStates.SetDefault(stateID, &oidcState{
		rolename:    rolename,
		nonce:       nonce,
		redirectURI: redirectURI,
	})

	return stateID, nonce, nil
}

// verifyState tests whether the provided state ID is valid and returns the
// associated state object if so. A nil state is returned if the ID is not found
// or expired. The state should only ever be retrieved once and is deleted as
// part of this request.
func (b *jwtAuthBackend) verifyState(stateID string) *oidcState {
	defer b.oidcStates.Delete(stateID)

	if stateRaw, ok := b.oidcStates.Get(stateID); ok {
		return stateRaw.(*oidcState)
	}

	return nil
}

// validRedirect checks whether uri is in allowed using special handling for loopback uris.
// Ref: https://tools.ietf.org/html/rfc8252#section-7.3
func validRedirect(uri string, allowed []string) bool {
	inputURI, err := url.Parse(uri)
	if err != nil {
		return false
	}

	// if uri isn't a loopback, just string search the allowed list
	if !strutil.StrListContains([]string{"localhost", "127.0.0.1", "::1"}, inputURI.Hostname()) {
		return strutil.StrListContains(allowed, uri)
	}

	// otherwise, search for a match in a port-agnostic manner, per the OAuth RFC.
	inputURI.Host = inputURI.Hostname()

	for _, a := range allowed {
		allowedURI, err := url.Parse(a)
		if err != nil {
			return false
		}
		allowedURI.Host = allowedURI.Hostname()

		if inputURI.String() == allowedURI.String() {
			return true
		}
	}

	return false
}
