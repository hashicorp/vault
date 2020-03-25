package jwtauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
)

var oidcStateTimeout = 10 * time.Minute

const (
	// OIDC error prefixes. These are searched for specifically by the UI, so any
	// changes to them must be aligned with a UI change.
	errLoginFailed       = "Vault login failed."
	errNoResponse        = "No response from provider."
	errTokenVerification = "Token verification failed."

	noCode = "no_code"
)

// oidcState is created when an authURL is requested. The state identifier is
// passed throughout the OAuth process.
type oidcState struct {
	rolename    string
	nonce       string
	redirectURI string
	code        string
	idToken     string
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
				"id_token": {
					Type: framework.TypeString,
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathCallback,
					Summary:  "Callback endpoint to complete an OIDC login.",

					// state is cached so don't process OIDC logins on perf standbys
					ForwardPerformanceStandby: true,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathCallbackPost,
					Summary:  "Callback endpoint to handle form_posts.",

					// state is cached so don't process OIDC logins on perf standbys
					ForwardPerformanceStandby: true,
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

					// state is cached so don't process OIDC logins on perf standbys
					ForwardPerformanceStandby: true,
				},
			},
		},
	}
}

func (b *jwtAuthBackend) pathCallbackPost(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse(errLoginFailed + " Could not load configuration."), nil
	}

	if config.OIDCResponseMode != responseModeFormPost {
		return logical.RespondWithStatusCode(nil, req, http.StatusMethodNotAllowed)
	}

	stateID := d.Get("state").(string)
	code := d.Get("code").(string)
	idToken := d.Get("id_token").(string)

	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: "text/html",
			logical.HTTPStatusCode:  http.StatusOK,
		},
	}

	// Store the provided code and/or token into state, which must already exist.
	state, err := b.amendState(stateID, code, idToken)
	if err != nil {
		resp.Data[logical.HTTPRawBody] = []byte(errorHTML(errLoginFailed, "Expired or missing OAuth state."))
		resp.Data[logical.HTTPStatusCode] = http.StatusBadRequest
	} else {
		mount := parseMount(state.redirectURI)
		if mount == "" {
			resp.Data[logical.HTTPRawBody] = []byte(errorHTML(errLoginFailed, "Invalid redirect path."))
			resp.Data[logical.HTTPStatusCode] = http.StatusBadRequest
		} else {
			resp.Data[logical.HTTPRawBody] = []byte(formpostHTML(mount, noCode, stateID))
		}
	}

	return resp, nil
}

func (b *jwtAuthBackend) pathCallback(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse(errLoginFailed + " Could not load configuration"), nil
	}

	stateID := d.Get("state").(string)

	state := b.verifyState(stateID)
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

	if len(role.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, role.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}
	}

	provider, err := b.getProvider(config)
	if err != nil {
		return nil, errwrap.Wrapf("error getting provider for login operation: {{err}}", err)
	}

	oidcCtx, err := b.createCAContext(ctx, config.OIDCDiscoveryCAPEM)
	if err != nil {
		return nil, errwrap.Wrapf("error preparing context for login operation: {{err}}", err)
	}

	var oauth2Config = oauth2.Config{
		ClientID:     config.OIDCClientID,
		ClientSecret: config.OIDCClientSecret,
		RedirectURL:  state.redirectURI,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}

	var rawToken string
	var oauth2Token *oauth2.Token

	code := d.Get("code").(string)
	if code == noCode {
		code = state.code
	}

	if code == "" {
		if state.idToken == "" {
			return logical.ErrorResponse(errLoginFailed + " No code or id_token received."), nil
		}
		rawToken = state.idToken
	} else {
		oauth2Token, err = oauth2Config.Exchange(oidcCtx, code)
		if err != nil {
			return logical.ErrorResponse(errLoginFailed+" Error exchanging oidc code: %q.", err.Error()), nil
		}

		// Extract the ID Token from OAuth2 token.
		var ok bool
		rawToken, ok = oauth2Token.Extra("id_token").(string)
		if !ok {
			return logical.ErrorResponse(errTokenVerification + " No id_token found in response."), nil
		}
	}

	if role.VerboseOIDCLogging {
		b.Logger().Debug("OIDC provider response", "ID token", rawToken)
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

	// If we have a token, attempt to fetch information from the /userinfo endpoint
	// and merge it with the existing claims data. A failure to fetch additional information
	// from this endpoint will not invalidate the authorization flow.
	if oauth2Token != nil {
		if userinfo, err := provider.UserInfo(oidcCtx, oauth2.StaticTokenSource(oauth2Token)); err == nil {
			_ = userinfo.Claims(&allClaims)
		} else {
			logFunc := b.Logger().Warn
			if strings.Contains(err.Error(), "user info endpoint is not supported") {
				logFunc = b.Logger().Info
			}
			logFunc("error reading /userinfo endpoint", "error", err)
		}
	}

	if role.VerboseOIDCLogging {
		if c, err := json.Marshal(allClaims); err == nil {
			b.Logger().Debug("OIDC provider response", "claims", string(c))
		} else {
			b.Logger().Debug("OIDC provider response", "marshalling error", err.Error())
		}
	}

	if err := validateBoundClaims(b.Logger(), role.BoundClaimsType, role.BoundClaims, allClaims); err != nil {
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

	auth := &logical.Auth{
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
	}

	role.PopulateTokenAuth(auth)

	resp := &logical.Response{
		Auth: auth,
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
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("could not load configuration"), nil
	}

	if config.authType() != OIDCFlow {
		return logical.ErrorResponse("OIDC login is not configured for this mount"), nil
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
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("role %q could not be found", roleName), nil
	}

	if !validRedirect(redirectURI, role.AllowedRedirectURIs) {
		logger.Warn("unauthorized redirect_uri", "redirect_uri", redirectURI)
		return resp, nil
	}

	// If configured for form_post, redirect directly to Vault instead of the UI,
	// if this was initiated by the UI (which currently has no knowledge of mode).
	///
	// TODO: it would be better to convey this to the UI and have it send the
	// correct URL directly.
	if config.OIDCResponseMode == responseModeFormPost {
		redirectURI = strings.Replace(redirectURI, "ui/vault", "v1", 1)
	}

	provider, err := b.getProvider(config)
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

	authCodeOpts := []oauth2.AuthCodeOption{
		oidc.Nonce(nonce),
	}

	// Add "form_post" param if requested. Note: the operator is allowed to configure "query"
	// as well, but that is the default for the AuthCode method and needn't be explicitly added.
	if config.OIDCResponseMode == responseModeFormPost {
		authCodeOpts = append(authCodeOpts, oauth2.SetAuthURLParam("response_mode", responseModeFormPost))
	}

	// Build the final authorization URL. oauth2Config doesn't support response types other than
	// code, so some manual tweaking is required.
	urlStr := oauth2Config.AuthCodeURL(stateID, authCodeOpts...)

	var rt string
	if config.hasType(responseTypeCode) {
		rt += responseTypeCode + " "
	}
	if config.hasType(responseTypeIDToken) {
		rt += responseTypeIDToken + " "
	}

	rt = strings.TrimSpace(rt)
	urlStr = strings.Replace(urlStr, "response_type=code",
		fmt.Sprintf("response_type=%s", url.QueryEscape(rt)), 1)

	resp.Data["auth_url"] = urlStr

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

func (b *jwtAuthBackend) amendState(stateID, code, idToken string) (*oidcState, error) {
	stateRaw, ok := b.oidcStates.Get(stateID)
	if !ok {
		return nil, errors.New("OIDC state not found")
	}

	state := stateRaw.(*oidcState)
	state.code = code
	state.idToken = idToken

	b.oidcStates.SetDefault(stateID, state)

	return state, nil
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

// parseMount attempts to extract the mount path from a redirect URI.
func parseMount(redirectURI string) string {
	parts := strings.Split(redirectURI, "/")

	for i := 0; i+2 < len(parts); i++ {
		if parts[i] == "v1" && parts[i+1] == "auth" {
			return parts[i+2]
		}
	}
	return ""
}
