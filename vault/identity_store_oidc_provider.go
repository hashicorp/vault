// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/identitytpl"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	// OIDC-related constants
	openIDScope              = "openid"
	scopesDelimiter          = " "
	accessTokenScopesMeta    = "scopes"
	accessTokenClientIDMeta  = "client_id"
	clientIDLength           = 32
	clientSecretLength       = 64
	clientSecretPrefix       = "hvo_secret_"
	codeChallengeMethodPlain = "plain"
	codeChallengeMethodS256  = "S256"
	defaultProviderName      = "default"
	defaultKeyName           = "default"
	allowAllAssignmentName   = "allow_all"

	// Storage path constants
	oidcProviderPrefix = "oidc_provider/"
	assignmentPath     = oidcProviderPrefix + "assignment/"
	scopePath          = oidcProviderPrefix + "scope/"
	clientPath         = oidcProviderPrefix + "client/"
	providerPath       = oidcProviderPrefix + "provider/"

	// Error constants used in the Authorization Endpoint. See details at
	// https://openid.net/specs/openid-connect-core-1_0.html#AuthError.
	ErrAuthUnsupportedResponseType = "unsupported_response_type"
	ErrAuthInvalidRequest          = "invalid_request"
	ErrAuthAccessDenied            = "access_denied"
	ErrAuthUnauthorizedClient      = "unauthorized_client"
	ErrAuthServerError             = "server_error"
	ErrAuthRequestNotSupported     = "request_not_supported"
	ErrAuthRequestURINotSupported  = "request_uri_not_supported"

	// Error constants used in the Token Endpoint. See details at
	// https://openid.net/specs/openid-connect-core-1_0.html#TokenErrorResponse
	ErrTokenInvalidRequest       = "invalid_request"
	ErrTokenInvalidClient        = "invalid_client"
	ErrTokenInvalidGrant         = "invalid_grant"
	ErrTokenUnsupportedGrantType = "unsupported_grant_type"
	ErrTokenServerError          = "server_error"

	// Error constants used in the UserInfo Endpoint. See details at
	// https://openid.net/specs/openid-connect-core-1_0.html#UserInfoError
	ErrUserInfoServerError    = "server_error"
	ErrUserInfoInvalidRequest = "invalid_request"
	ErrUserInfoInvalidToken   = "invalid_token"
	ErrUserInfoAccessDenied   = "access_denied"

	// The following errors are used by the UI for specific behavior of
	// the OIDC specification. Any changes to their values must come with
	// a corresponding change in the UI code.
	ErrAuthInvalidClientID      = "invalid_client_id"
	ErrAuthInvalidRedirectURI   = "invalid_redirect_uri"
	ErrAuthMaxAgeReAuthenticate = "max_age_violation"
)

type assignment struct {
	GroupIDs  []string `json:"group_ids"`
	EntityIDs []string `json:"entity_ids"`
}

type scope struct {
	Template    string `json:"template"`
	Description string `json:"description"`
}

type client struct {
	// Used for indexing in memdb
	Name        string `json:"name"`
	NamespaceID string `json:"namespace_id"`

	// User-supplied parameters
	RedirectURIs   []string      `json:"redirect_uris"`
	Assignments    []string      `json:"assignments"`
	Key            string        `json:"key"`
	IDTokenTTL     time.Duration `json:"id_token_ttl"`
	AccessTokenTTL time.Duration `json:"access_token_ttl"`
	Type           clientType    `json:"type"`

	// Generated values that are used in OIDC endpoints
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

//go:generate enumer -type=clientType -trimprefix=clientType -transform=snake
type clientType int

const (
	confidential clientType = iota
	public
)

type provider struct {
	Issuer           string   `json:"issuer"`
	AllowedClientIDs []string `json:"allowed_client_ids"`
	ScopesSupported  []string `json:"scopes_supported"`

	// effectiveIssuer is a calculated field and will be either Issuer (if
	// that's set) or the Vault instance's api_addr.
	effectiveIssuer string
}

// allowedClientID returns true if the given client ID is in
// the provider's set of allowed client IDs or its allowed client
// IDs contains the wildcard "*" char.
func (p *provider) allowedClientID(clientID string) bool {
	for _, allowedID := range p.AllowedClientIDs {
		switch allowedID {
		case "*", clientID:
			return true
		}
	}
	return false
}

type providerDiscovery struct {
	Issuer                string   `json:"issuer"`
	Keys                  string   `json:"jwks_uri"`
	AuthorizationEndpoint string   `json:"authorization_endpoint"`
	TokenEndpoint         string   `json:"token_endpoint"`
	UserinfoEndpoint      string   `json:"userinfo_endpoint"`
	RequestParameter      bool     `json:"request_parameter_supported"`
	RequestURIParameter   bool     `json:"request_uri_parameter_supported"`
	IDTokenAlgs           []string `json:"id_token_signing_alg_values_supported"`
	ResponseTypes         []string `json:"response_types_supported"`
	Scopes                []string `json:"scopes_supported"`
	Claims                []string `json:"claims_supported"`
	Subjects              []string `json:"subject_types_supported"`
	GrantTypes            []string `json:"grant_types_supported"`
	AuthMethods           []string `json:"token_endpoint_auth_methods_supported"`
	CodeChallengeMethods  []string `json:"code_challenge_methods_supported"`
}

type authCodeCacheEntry struct {
	provider            string
	clientID            string
	entityID            string
	redirectURI         string
	nonce               string
	scopes              []string
	authTime            time.Time
	codeChallenge       string
	codeChallengeMethod string
}

func oidcProviderPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "oidc/assignment/" + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "assignment",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the assignment",
				},
				"entity_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Comma separated string or array of identity entity IDs",
				},
				"group_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Comma separated string or array of identity group IDs",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathOIDCCreateUpdateAssignment,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: i.pathOIDCCreateUpdateAssignment,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathOIDCReadAssignment,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathOIDCDeleteAssignment,
				},
			},
			ExistenceCheck:  i.pathOIDCAssignmentExistenceCheck,
			HelpSynopsis:    "CRUD operations for OIDC assignments.",
			HelpDescription: "Create, Read, Update, and Delete OIDC assignments.",
		},
		{
			Pattern: "oidc/assignment/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "assignments",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.pathOIDCListAssignment,
				},
			},
			HelpSynopsis:    "List OIDC assignments",
			HelpDescription: "List all configured OIDC assignments in the identity backend.",
		},
		{
			Pattern: "oidc/scope/" + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "scope",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the scope",
				},
				"template": {
					Type:        framework.TypeString,
					Description: "The template string to use for the scope. This may be in string-ified JSON or base64 format.",
				},
				"description": {
					Type:        framework.TypeString,
					Description: "The description of the scope",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathOIDCCreateUpdateScope,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: i.pathOIDCCreateUpdateScope,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathOIDCReadScope,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathOIDCDeleteScope,
				},
			},
			ExistenceCheck:  i.pathOIDCScopeExistenceCheck,
			HelpSynopsis:    "CRUD operations for OIDC scopes.",
			HelpDescription: "Create, Read, Update, and Delete OIDC scopes.",
		},
		{
			Pattern: "oidc/scope/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "scopes",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.pathOIDCListScope,
				},
			},
			HelpSynopsis:    "List OIDC scopes",
			HelpDescription: "List all configured OIDC scopes in the identity backend.",
		},
		{
			Pattern: "oidc/client/" + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "client",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the client.",
				},
				"redirect_uris": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Comma separated string or array of redirect URIs used by the client. One of these values must exactly match the redirect_uri parameter value used in each authentication request.",
				},
				"assignments": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Comma separated string or array of assignment resources.",
				},
				"key": {
					Type:        framework.TypeString,
					Description: "A reference to a named key resource. Cannot be modified after creation. Defaults to the 'default' key.",
					Default:     "default",
				},
				"id_token_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "The time-to-live for ID tokens obtained by the client.",
					Default:     "24h",
				},
				"access_token_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "The time-to-live for access tokens obtained by the client.",
					Default:     "24h",
				},
				"client_type": {
					Type:        framework.TypeString,
					Description: "The client type based on its ability to maintain confidentiality of credentials. The following client types are supported: 'confidential', 'public'. Defaults to 'confidential'.",
					Default:     "confidential",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathOIDCCreateUpdateClient,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: i.pathOIDCCreateUpdateClient,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathOIDCReadClient,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathOIDCDeleteClient,
				},
			},
			ExistenceCheck:  i.pathOIDCClientExistenceCheck,
			HelpSynopsis:    "CRUD operations for OIDC clients.",
			HelpDescription: "Create, Read, Update, and Delete OIDC clients.",
		},
		{
			Pattern: "oidc/client/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "clients",
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.pathOIDCListClient,
				},
			},
			HelpSynopsis:    "List OIDC clients",
			HelpDescription: "List all configured OIDC clients in the identity backend.",
		},
		{
			Pattern: "oidc/provider/" + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "provider",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the provider",
				},
				"issuer": {
					Type:        framework.TypeString,
					Description: "Specifies what will be used for the iss claim of ID tokens.",
				},
				"allowed_client_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "The client IDs that are permitted to use the provider",
				},
				"scopes_supported": {
					Type:        framework.TypeCommaStringSlice,
					Description: "The scopes supported for requesting on the provider",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathOIDCCreateUpdateProvider,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: i.pathOIDCCreateUpdateProvider,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathOIDCReadProvider,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathOIDCDeleteProvider,
				},
			},
			ExistenceCheck:  i.pathOIDCProviderExistenceCheck,
			HelpSynopsis:    "CRUD operations for OIDC providers.",
			HelpDescription: "Create, Read, Update, and Delete OIDC named providers.",
		},
		{
			Pattern: "oidc/provider/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "providers",
			},
			Fields: map[string]*framework.FieldSchema{
				"allowed_client_id": {
					Type: framework.TypeString,
					Description: "Filters the list of OIDC providers to those " +
						"that allow the given client ID in their set of allowed_client_ids.",
					Default: "",
					Query:   true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.pathOIDCListProvider,
				},
			},
			HelpSynopsis:    "List OIDC providers",
			HelpDescription: "List all configured OIDC providers in the identity backend.",
		},
		{
			Pattern: "oidc/provider/" + framework.GenericNameRegex("name") + "/\\.well-known/openid-configuration",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "provider-open-id-configuration",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the provider",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathOIDCProviderDiscovery,
				},
			},
			HelpSynopsis:    "Query OIDC configurations",
			HelpDescription: "Query this path to retrieve the configured OIDC Issuer and Keys endpoints, response types, subject types, and signing algorithms used by the OIDC backend.",
		},
		{
			Pattern: "oidc/provider/" + framework.GenericNameRegex("name") + "/\\.well-known/keys",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "provider-public-keys",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the provider",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathOIDCReadProviderPublicKeys,
				},
			},
			HelpSynopsis:    "Retrieve public keys",
			HelpDescription: "Returns the public portion of keys for a named OIDC provider. Clients can use them to validate the authenticity of an ID token.",
		},
		{
			Pattern: "oidc/provider/" + framework.GenericNameRegex("name") + "/authorize",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc-provider",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the provider",
				},
				"client_id": {
					Type:        framework.TypeString,
					Description: "The ID of the requesting client.",
					Required:    true,
					Query:       true,
				},
				"scope": {
					Type:        framework.TypeString,
					Description: "A space-delimited, case-sensitive list of scopes to be requested. The 'openid' scope is required.",
					Required:    true,
					Query:       true,
				},
				"redirect_uri": {
					Type:        framework.TypeString,
					Description: "The redirection URI to which the response will be sent.",
					Required:    true,
					Query:       true,
				},
				"response_type": {
					Type:        framework.TypeString,
					Description: "The OIDC authentication flow to be used. The following response types are supported: 'code'",
					Required:    true,
					Query:       true,
				},
				"state": {
					Type:        framework.TypeString,
					Description: "The value used to maintain state between the authentication request and client.",
					Query:       true,
				},
				"nonce": {
					Type:        framework.TypeString,
					Description: "The value that will be returned in the ID token nonce claim after a token exchange.",
					Query:       true,
				},
				"max_age": {
					Type:        framework.TypeInt,
					Description: "The allowable elapsed time in seconds since the last time the end-user was actively authenticated.",
					Query:       true,
				},
				"code_challenge": {
					Type:        framework.TypeString,
					Description: "The code challenge derived from the code verifier.",
					Query:       true,
				},
				"code_challenge_method": {
					Type:        framework.TypeString,
					Description: "The method that was used to derive the code challenge. The following methods are supported: 'S256', 'plain'. Defaults to 'plain'.",
					Default:     codeChallengeMethodPlain,
					Query:       true,
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathOIDCAuthorize,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "authorize",
					},
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: false,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathOIDCAuthorize,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb:   "authorize",
						OperationSuffix: "with-parameters",
					},
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: false,
				},
			},
			HelpSynopsis:    "Provides the OIDC Authorization Endpoint.",
			HelpDescription: "The OIDC Authorization Endpoint performs authentication and authorization by using request parameters defined by OpenID Connect (OIDC).",
		},
		{
			Pattern: "oidc/provider/" + framework.GenericNameRegex("name") + "/token",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc-provider",
				OperationVerb:   "token",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the provider",
				},
				"code": {
					Type:        framework.TypeString,
					Description: "The authorization code received from the provider's authorization endpoint.",
					Required:    true,
				},
				"grant_type": {
					Type:        framework.TypeString,
					Description: "The authorization grant type. The following grant types are supported: 'authorization_code'.",
					Required:    true,
				},
				"redirect_uri": {
					Type:        framework.TypeString,
					Description: "The callback location where the authentication response was sent.",
					Required:    true,
				},
				"code_verifier": {
					Type:        framework.TypeString,
					Description: "The code verifier associated with the authorization code.",
				},
				// For confidential clients, the client_id and client_secret are provided to
				// the token endpoint via the 'client_secret_basic' or 'client_secret_post'
				// authentication methods. See the OIDC spec for details at:
				// https://openid.net/specs/openid-connect-core-1_0.html#ClientAuthentication

				// For public clients, the client_id is required and a client_secret does
				// not exist. This means that public clients use the 'none' authentication
				// method. However, public clients are required to use Proof Key for Code
				// Exchange (PKCE) when using the authorization code flow.
				"client_id": {
					Type:        framework.TypeString,
					Description: "The ID of the requesting client.",
				},
				"client_secret": {
					Type:        framework.TypeString,
					Description: "The secret of the requesting client.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    i.pathOIDCToken,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: false,
				},
			},
			HelpSynopsis:    "Provides the OIDC Token Endpoint.",
			HelpDescription: "The OIDC Token Endpoint allows a client to exchange its Authorization Grant for an Access Token and ID Token.",
		},
		{
			Pattern: "oidc/provider/" + framework.GenericNameRegex("name") + "/userinfo",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc-provider",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the provider",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathOIDCUserInfo,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "user-info",
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathOIDCUserInfo,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "user-info2",
					},
				},
			},
			HelpSynopsis:    "Provides the OIDC UserInfo Endpoint.",
			HelpDescription: "The OIDC UserInfo Endpoint returns claims about the authenticated end-user.",
		},
	}
}

// clientsReferencingTargetAssignmentName returns a map of client names to
// clients referencing targetAssignmentName.
func (i *IdentityStore) clientsReferencingTargetAssignmentName(ctx context.Context, req *logical.Request, targetAssignmentName string) (map[string]client, error) {
	clientNames, err := req.Storage.List(ctx, clientPath)
	if err != nil {
		return nil, err
	}

	var tempClient client
	clients := make(map[string]client)
	for _, clientName := range clientNames {
		entry, err := req.Storage.Get(ctx, clientPath+clientName)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&tempClient); err != nil {
				return nil, err
			}
			for _, a := range tempClient.Assignments {
				if a == targetAssignmentName {
					clients[clientName] = tempClient
				}
			}
		}
	}

	return clients, nil
}

// clientNamesReferencingTargetAssignmentName returns a slice of strings of client
// names referencing targetAssignmentName.
func (i *IdentityStore) clientNamesReferencingTargetAssignmentName(ctx context.Context, req *logical.Request, targetAssignmentName string) ([]string, error) {
	clients, err := i.clientsReferencingTargetAssignmentName(ctx, req, targetAssignmentName)
	if err != nil {
		return nil, err
	}

	var names []string
	for client := range clients {
		names = append(names, client)
	}
	sort.Strings(names)
	return names, nil
}

// clientsReferencingTargetKeyName returns a map of client names to
// clients referencing targetKeyName.
func (i *IdentityStore) clientsReferencingTargetKeyName(ctx context.Context, req *logical.Request, targetKeyName string) (map[string]client, error) {
	clientNames, err := req.Storage.List(ctx, clientPath)
	if err != nil {
		return nil, err
	}

	var tempClient client
	clients := make(map[string]client)
	for _, clientName := range clientNames {
		entry, err := req.Storage.Get(ctx, clientPath+clientName)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&tempClient); err != nil {
				return nil, err
			}
			if tempClient.Key == targetKeyName {
				clients[clientName] = tempClient
			}
		}
	}

	return clients, nil
}

// clientNamesReferencingTargetKeyName returns a slice of strings of client
// names referencing targetKeyName.
func (i *IdentityStore) clientNamesReferencingTargetKeyName(ctx context.Context, req *logical.Request, targetKeyName string) ([]string, error) {
	clients, err := i.clientsReferencingTargetKeyName(ctx, req, targetKeyName)
	if err != nil {
		return nil, err
	}

	var names []string
	for client := range clients {
		names = append(names, client)
	}
	sort.Strings(names)
	return names, nil
}

// providersReferencingTargetScopeName returns a list of provider names referencing targetScopeName.
// Not threadsafe. To be called with lock already held.
func (i *IdentityStore) providersReferencingTargetScopeName(ctx context.Context, req *logical.Request, targetScopeName string) ([]string, error) {
	providerNames, err := req.Storage.List(ctx, providerPath)
	if err != nil {
		return nil, err
	}

	var tempProvider provider
	var providers []string
	for _, providerName := range providerNames {
		entry, err := req.Storage.Get(ctx, providerPath+providerName)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&tempProvider); err != nil {
				return nil, err
			}
			for _, a := range tempProvider.ScopesSupported {
				if a == targetScopeName {
					providers = append(providers, providerName)
				}
			}
		}
	}

	return providers, nil
}

// pathOIDCCreateUpdateAssignment is used to create a new assignment or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateAssignment(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	if name == allowAllAssignmentName {
		return logical.ErrorResponse("modification of assignment %q not allowed",
			allowAllAssignmentName), nil
	}

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	var assignment assignment
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, assignmentPath+name)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&assignment); err != nil {
				return nil, err
			}
		}
	}

	if entitiesRaw, ok := d.GetOk("entity_ids"); ok {
		assignment.EntityIDs = entitiesRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		assignment.EntityIDs = d.Get("entity_ids").([]string)
	}

	if groupsRaw, ok := d.GetOk("group_ids"); ok {
		assignment.GroupIDs = groupsRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		assignment.GroupIDs = d.Get("group_ids").([]string)
	}

	// remove duplicates and lowercase entities and groups
	assignment.EntityIDs = strutil.RemoveDuplicates(assignment.EntityIDs, true)
	assignment.GroupIDs = strutil.RemoveDuplicates(assignment.GroupIDs, true)

	// store assignment
	entry, err := logical.StorageEntryJSON(assignmentPath+name, assignment)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathOIDCListAssignment is used to list assignments
func (i *IdentityStore) pathOIDCListAssignment(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	assignments, err := req.Storage.List(ctx, assignmentPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(assignments), nil
}

// pathOIDCReadAssignment is used to read an existing assignment
func (i *IdentityStore) pathOIDCReadAssignment(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	assignment, err := i.getOIDCAssignment(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"group_ids":  assignment.GroupIDs,
			"entity_ids": assignment.EntityIDs,
		},
	}, nil
}

func (i *IdentityStore) getOIDCAssignment(ctx context.Context, s logical.Storage, name string) (*assignment, error) {
	entry, err := s.Get(ctx, assignmentPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var assignment assignment
	if err := entry.DecodeJSON(&assignment); err != nil {
		return nil, err
	}

	return &assignment, nil
}

// pathOIDCDeleteAssignment is used to delete an assignment
func (i *IdentityStore) pathOIDCDeleteAssignment(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	if name == allowAllAssignmentName {
		return logical.ErrorResponse("deletion of assignment %q not allowed",
			allowAllAssignmentName), nil
	}

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	clientNames, err := i.clientNamesReferencingTargetAssignmentName(ctx, req, name)
	if err != nil {
		return nil, err
	}

	if len(clientNames) > 0 {
		errorMessage := fmt.Sprintf("unable to delete assignment %q because it is currently referenced by these clients: %s",
			name, strings.Join(clientNames, ", "))
		return logical.ErrorResponse(errorMessage), logical.ErrInvalidRequest
	}

	err = req.Storage.Delete(ctx, assignmentPath+name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *IdentityStore) pathOIDCAssignmentExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, assignmentPath+name)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

// pathOIDCCreateUpdateScope is used to create a new scope or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateScope(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == openIDScope {
		return logical.ErrorResponse("the %q scope name is reserved", openIDScope), nil
	}

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	var scope scope
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, scopePath+name)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&scope); err != nil {
				return nil, err
			}
		}
	}

	if descriptionRaw, ok := d.GetOk("description"); ok {
		scope.Description = descriptionRaw.(string)
	} else if req.Operation == logical.CreateOperation {
		scope.Description = d.Get("description").(string)
	}

	if templateRaw, ok := d.GetOk("template"); ok {
		scope.Template = templateRaw.(string)
	} else if req.Operation == logical.CreateOperation {
		scope.Template = d.Get("template").(string)
	}

	// Attempt to decode as base64 and use that if it works
	if decoded, err := base64.StdEncoding.DecodeString(scope.Template); err == nil {
		scope.Template = string(decoded)
	}

	// Validate that template can be parsed and results in valid JSON
	if scope.Template != "" {
		_, populatedTemplate, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
			Mode:   identitytpl.JSONTemplating,
			String: scope.Template,
			Entity: new(logical.Entity),
			Groups: make([]*logical.Group, 0),
		})
		if err != nil {
			return logical.ErrorResponse("error parsing template: %s", err.Error()), nil
		}

		var tmp map[string]interface{}
		if err := json.Unmarshal([]byte(populatedTemplate), &tmp); err != nil {
			return logical.ErrorResponse("error parsing template JSON: %s", err.Error()), nil
		}

		for key := range tmp {
			if strutil.StrListContains(reservedClaims, key) {
				return logical.ErrorResponse(`top level key %q not allowed. Restricted keys: %s`,
					key, strings.Join(reservedClaims, ", ")), nil
			}
		}
	}
	// store scope
	entry, err := logical.StorageEntryJSON(scopePath+name, scope)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathOIDCListScope is used to list scopes
func (i *IdentityStore) pathOIDCListScope(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	scopes, err := req.Storage.List(ctx, scopePath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(scopes), nil
}

// pathOIDCReadScope is used to read an existing scope
func (i *IdentityStore) pathOIDCReadScope(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	scope, err := i.getOIDCScope(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if scope == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"template":    scope.Template,
			"description": scope.Description,
		},
	}, nil
}

func (i *IdentityStore) getOIDCScope(ctx context.Context, s logical.Storage, name string) (*scope, error) {
	entry, err := s.Get(ctx, scopePath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var scope scope
	if err := entry.DecodeJSON(&scope); err != nil {
		return nil, err
	}

	return &scope, nil
}

// pathOIDCDeleteScope is used to delete a scope
func (i *IdentityStore) pathOIDCDeleteScope(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	providerNames, err := i.providersReferencingTargetScopeName(ctx, req, name)
	if err != nil {
		return nil, err
	}

	if len(providerNames) > 0 {
		errorMessage := fmt.Sprintf("unable to delete scope %q because it is currently referenced by these providers: %s",
			name, strings.Join(providerNames, ", "))
		return logical.ErrorResponse(errorMessage), logical.ErrInvalidRequest
	}
	err = req.Storage.Delete(ctx, scopePath+name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *IdentityStore) pathOIDCScopeExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, scopePath+name)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

// pathOIDCCreateUpdateClient is used to create a new client or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateClient(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	client := client{
		Name:        name,
		NamespaceID: ns.ID,
	}
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, clientPath+name)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&client); err != nil {
				return nil, err
			}
		}
	}

	if redirectURIsRaw, ok := d.GetOk("redirect_uris"); ok {
		client.RedirectURIs = redirectURIsRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		client.RedirectURIs = d.Get("redirect_uris").([]string)
	}

	if assignmentsRaw, ok := d.GetOk("assignments"); ok {
		client.Assignments = assignmentsRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		client.Assignments = d.Get("assignments").([]string)
	}

	// remove duplicate assignments and redirect URIs
	client.Assignments = strutil.RemoveDuplicates(client.Assignments, false)
	client.RedirectURIs = strutil.RemoveDuplicates(client.RedirectURIs, false)

	// enforce assignment existence
	for _, assignment := range client.Assignments {
		entry, err := req.Storage.Get(ctx, assignmentPath+assignment)
		if err != nil {
			return nil, err
		}
		if entry == nil {
			return logical.ErrorResponse("assignment %q does not exist", assignment), nil
		}
	}

	if keyRaw, ok := d.GetOk("key"); ok {
		key := keyRaw.(string)
		if req.Operation == logical.UpdateOperation && client.Key != key {
			return logical.ErrorResponse("key modification is not allowed"), nil
		}
		client.Key = key
	} else if req.Operation == logical.CreateOperation {
		client.Key = d.Get("key").(string)
	}

	// enforce key existence on client creation
	key, err := i.getNamedKey(ctx, req.Storage, client.Key)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return logical.ErrorResponse("key %q does not exist", client.Key), nil
	}

	if client.Key == defaultKeyName {
		if err := i.lazyGenerateDefaultKey(ctx, req.Storage); err != nil {
			return nil, fmt.Errorf("failed to generate default key: %w", err)
		}
	}

	if idTokenTTLRaw, ok := d.GetOk("id_token_ttl"); ok {
		client.IDTokenTTL = time.Duration(idTokenTTLRaw.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		client.IDTokenTTL = time.Duration(d.Get("id_token_ttl").(int)) * time.Second
	}

	if client.IDTokenTTL > key.VerificationTTL {
		return logical.ErrorResponse("a client's id_token_ttl cannot be greater than the verification_ttl of the key it references"), nil
	}

	if accessTokenTTLRaw, ok := d.GetOk("access_token_ttl"); ok {
		client.AccessTokenTTL = time.Duration(accessTokenTTLRaw.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		client.AccessTokenTTL = time.Duration(d.Get("access_token_ttl").(int)) * time.Second
	}

	if clientTypeRaw, ok := d.GetOk("client_type"); ok {
		clientType := clientTypeRaw.(string)
		if req.Operation == logical.UpdateOperation && client.Type.String() != clientType {
			return logical.ErrorResponse("client_type modification is not allowed"), nil
		}

		switch clientType {
		case confidential.String():
			client.Type = confidential
		case public.String():
			client.Type = public
		default:
			return logical.ErrorResponse("invalid client_type %q", clientType), nil
		}
	}

	if client.ClientID == "" {
		// generate client_id
		clientID, err := base62.Random(clientIDLength)
		if err != nil {
			return nil, err
		}
		client.ClientID = clientID
	}

	// client secrets are only generated for confidential clients
	if client.Type == confidential && client.ClientSecret == "" {
		// generate client_secret
		clientSecret, err := base62.Random(clientSecretLength)
		if err != nil {
			return nil, err
		}
		client.ClientSecret = clientSecretPrefix + clientSecret
	}

	// invalidate the cached client in memdb
	if err := i.memDBDeleteClientByName(ctx, name); err != nil {
		return nil, err
	}

	// store client
	entry, err := logical.StorageEntryJSON(clientPath+name, client)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathOIDCListClient is used to list clients
func (i *IdentityStore) pathOIDCListClient(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	clients, err := i.listClients(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(clients))
	keyInfo := make(map[string]interface{})
	for _, client := range clients {
		keys = append(keys, client.Name)
		keyInfo[client.Name] = map[string]interface{}{
			"redirect_uris":    client.RedirectURIs,
			"assignments":      client.Assignments,
			"key":              client.Key,
			"id_token_ttl":     int64(client.IDTokenTTL.Seconds()),
			"access_token_ttl": int64(client.AccessTokenTTL.Seconds()),
			"client_type":      client.Type.String(),
			"client_id":        client.ClientID,
			// client_secret is intentionally omitted
		}
	}

	return logical.ListResponseWithInfo(keys, keyInfo), nil
}

// pathOIDCReadClient is used to read an existing client
func (i *IdentityStore) pathOIDCReadClient(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	client, err := i.clientByName(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"redirect_uris":    client.RedirectURIs,
			"assignments":      client.Assignments,
			"key":              client.Key,
			"id_token_ttl":     int64(client.IDTokenTTL.Seconds()),
			"access_token_ttl": int64(client.AccessTokenTTL.Seconds()),
			"client_id":        client.ClientID,
			"client_type":      client.Type.String(),
		},
	}

	if client.Type == confidential {
		resp.Data["client_secret"] = client.ClientSecret
	}

	return resp, nil
}

// pathOIDCDeleteClient is used to delete a client
func (i *IdentityStore) pathOIDCDeleteClient(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	// Delete the client from memdb
	if err := i.memDBDeleteClientByName(ctx, name); err != nil {
		return nil, err
	}

	// Delete the client from storage
	if err := req.Storage.Delete(ctx, clientPath+name); err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *IdentityStore) pathOIDCClientExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, clientPath+name)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

// pathOIDCCreateUpdateProvider is used to create a new named provider or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var resp logical.Response
	name := d.Get("name").(string)

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	var provider provider
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, providerPath+name)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&provider); err != nil {
				return nil, err
			}
		}
	}

	if issuerRaw, ok := d.GetOk("issuer"); ok {
		provider.Issuer = issuerRaw.(string)
	} else if req.Operation == logical.CreateOperation {
		provider.Issuer = d.Get("issuer").(string)
	}

	if allowedClientIDsRaw, ok := d.GetOk("allowed_client_ids"); ok {
		provider.AllowedClientIDs = allowedClientIDsRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		provider.AllowedClientIDs = d.Get("allowed_client_ids").([]string)
	}

	if scopesRaw, ok := d.GetOk("scopes_supported"); ok {
		provider.ScopesSupported = scopesRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		provider.ScopesSupported = d.Get("scopes_supported").([]string)
	}

	// remove duplicate allowed client IDs and scopes
	provider.AllowedClientIDs = strutil.RemoveDuplicates(provider.AllowedClientIDs, false)
	provider.ScopesSupported = strutil.RemoveDuplicates(provider.ScopesSupported, false)

	if provider.Issuer != "" {
		// verify that issuer is the correct format:
		//   - http or https
		//   - host name
		//   - optional port
		//   - nothing more
		valid := false
		if u, err := url.Parse(provider.Issuer); err == nil {
			u2 := url.URL{
				Scheme: u.Scheme,
				Host:   u.Host,
			}
			valid = (*u == u2) &&
				(u.Scheme == "http" || u.Scheme == "https") &&
				u.Host != ""
		}

		if !valid {
			return logical.ErrorResponse(
				"invalid issuer, which must include only a scheme, host, " +
					"and optional port (e.g. https://example.com:8200)"), nil
		}

		resp.AddWarning(`If "issuer" is set explicitly, all tokens must be ` +
			`validated against that address, including those issued by secondary ` +
			`clusters. Setting issuer to "" will restore the default behavior of ` +
			`using the cluster's api_addr as the issuer.`)

	}

	scopeTemplateKeyNames := make(map[string]string)
	for _, scopeName := range provider.ScopesSupported {
		scope, err := i.getOIDCScope(ctx, req.Storage, scopeName)
		if err != nil {
			return nil, err
		}
		// enforce scope existence on provider create and update
		if scope == nil {
			return logical.ErrorResponse("scope %q does not exist", scopeName), nil
		}

		// ensure no two templates have the same top-level keys
		_, populatedTemplate, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
			Mode:   identitytpl.JSONTemplating,
			String: scope.Template,
			Entity: new(logical.Entity),
			Groups: make([]*logical.Group, 0),
		})
		if err != nil {
			return nil, fmt.Errorf("error parsing template for scope %q: %s", scopeName, err.Error())
		}

		jsonTemplate := make(map[string]interface{})
		if err = json.Unmarshal([]byte(populatedTemplate), &jsonTemplate); err != nil {
			return nil, err
		}

		for keyName := range jsonTemplate {
			val, ok := scopeTemplateKeyNames[keyName]
			if ok && val != scopeName {
				resp.AddWarning(fmt.Sprintf("Found scope templates with conflicting top-level keys: "+
					"conflict %q in scopes %q, %q. This may result in an error if the scopes are "+
					"requested in an OIDC Authentication Request.", keyName, scopeName, val))
			}

			scopeTemplateKeyNames[keyName] = scopeName
		}
	}

	// store named provider
	entry, err := logical.StorageEntryJSON(providerPath+name, provider)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	if len(resp.Warnings) == 0 {
		return nil, nil
	}

	return &resp, nil
}

// pathOIDCListProvider is used to list named providers
func (i *IdentityStore) pathOIDCListProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	providers, err := req.Storage.List(ctx, providerPath)
	if err != nil {
		return nil, err
	}

	// Build a map from provider name to provider struct
	providerMap := make(map[string]*provider)
	for _, name := range providers {
		provider, err := i.getOIDCProvider(ctx, req.Storage, name)
		if err != nil {
			return nil, err
		}
		if provider == nil {
			continue
		}
		providerMap[name] = provider
	}

	// If allowed_client_id is provided as a query parameter, filter the set of
	// returned OIDC providers to those that allow the given value in their set
	// of allowed_client_ids.
	if clientID := d.Get("allowed_client_id").(string); clientID != "" {
		for name, provider := range providerMap {
			if !provider.allowedClientID(clientID) {
				delete(providerMap, name)
			}
		}
	}

	keys := make([]string, 0, len(providerMap))
	keyInfo := make(map[string]interface{})
	for name, provider := range providerMap {
		keys = append(keys, name)
		keyInfo[name] = map[string]interface{}{
			"issuer":             provider.effectiveIssuer,
			"allowed_client_ids": provider.AllowedClientIDs,
			"scopes_supported":   provider.ScopesSupported,
		}
	}

	return logical.ListResponseWithInfo(keys, keyInfo), nil
}

// pathOIDCReadProvider is used to read an existing provider
func (i *IdentityStore) pathOIDCReadProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	provider, err := i.getOIDCProvider(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"issuer":             provider.effectiveIssuer,
			"allowed_client_ids": provider.AllowedClientIDs,
			"scopes_supported":   provider.ScopesSupported,
		},
	}, nil
}

func (i *IdentityStore) getOIDCProvider(ctx context.Context, s logical.Storage, name string) (*provider, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	entry, err := s.Get(ctx, providerPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var provider provider
	if err := entry.DecodeJSON(&provider); err != nil {
		return nil, err
	}

	provider.effectiveIssuer = provider.Issuer
	if provider.effectiveIssuer == "" {
		provider.effectiveIssuer = i.redirectAddr
	}

	provider.effectiveIssuer += "/v1/" + ns.Path + "identity/oidc/provider/" + name

	return &provider, nil
}

// pathOIDCDeleteProvider is used to delete a provider
func (i *IdentityStore) pathOIDCDeleteProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	if name == defaultProviderName {
		return logical.ErrorResponse("deletion of OIDC provider %q not allowed",
			defaultProviderName), nil
	}

	return nil, req.Storage.Delete(ctx, providerPath+name)
}

func (i *IdentityStore) pathOIDCProviderExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, providerPath+name)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

func (i *IdentityStore) pathOIDCProviderDiscovery(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	p, err := i.getOIDCProvider(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	// the "openid" scope is reserved and is included for every provider
	scopes := append(p.ScopesSupported, openIDScope)

	disc := providerDiscovery{
		Issuer:                p.effectiveIssuer,
		Keys:                  p.effectiveIssuer + "/.well-known/keys",
		AuthorizationEndpoint: strings.Replace(p.effectiveIssuer, "/v1/", "/ui/vault/", 1) + "/authorize",
		TokenEndpoint:         p.effectiveIssuer + "/token",
		UserinfoEndpoint:      p.effectiveIssuer + "/userinfo",
		IDTokenAlgs:           supportedAlgs,
		Scopes:                scopes,
		Claims:                []string{},
		RequestParameter:      false,
		RequestURIParameter:   false,
		ResponseTypes:         []string{"code"},
		Subjects:              []string{"public"},
		GrantTypes:            []string{"authorization_code"},
		AuthMethods: []string{
			// PKCE is required for auth method "none"
			"none",
			"client_secret_basic",
			"client_secret_post",
		},
		CodeChallengeMethods: []string{
			codeChallengeMethodPlain,
			codeChallengeMethodS256,
		},
	}

	data, err := json.Marshal(disc)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:         200,
			logical.HTTPRawBody:            data,
			logical.HTTPContentType:        "application/json",
			logical.HTTPCacheControlHeader: "max-age=3600",
		},
	}

	return resp, nil
}

// pathOIDCReadProviderPublicKeys is used to retrieve all public keys for a
// named provider so that clients can verify the validity of a signed OIDC token.
func (i *IdentityStore) pathOIDCReadProviderPublicKeys(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	providerName := d.Get("name").(string)

	var provider provider

	providerEntry, err := req.Storage.Get(ctx, providerPath+providerName)
	if err != nil {
		return nil, err
	}
	if providerEntry == nil {
		return nil, nil
	}
	if err := providerEntry.DecodeJSON(&provider); err != nil {
		return nil, err
	}

	keyIDs, err := i.keyIDsReferencedByTargetClientIDs(ctx, req.Storage, provider.AllowedClientIDs)
	if err != nil {
		return nil, err
	}

	jwks := &jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 0, len(keyIDs)),
	}

	for _, keyID := range keyIDs {
		key, err := loadOIDCPublicKey(ctx, req.Storage, keyID)
		if err != nil {
			return nil, err
		}
		jwks.Keys = append(jwks.Keys, *key)
	}

	data, err := json.Marshal(jwks)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:  200,
			logical.HTTPRawBody:     data,
			logical.HTTPContentType: "application/json",
		},
	}

	return resp, nil
}

// keyIDsReferencedByTargetClientIDs returns a slice of key IDs that are
// referenced by the clients' targetIDs.
// If targetIDs contains "*" then the IDs for all public keys are returned.
func (i *IdentityStore) keyIDsReferencedByTargetClientIDs(ctx context.Context, s logical.Storage, targetIDs []string) ([]string, error) {
	keyNames := make(map[string]bool)

	// Get all key names referenced by clients if wildcard "*" in target client IDs
	if strutil.StrListContains(targetIDs, "*") {
		clients, err := i.listClients(ctx, s)
		if err != nil {
			return nil, err
		}

		for _, client := range clients {
			keyNames[client.Key] = true
		}
	}

	// Otherwise, get the key names referenced by each target client ID
	if len(keyNames) == 0 {
		for _, clientID := range targetIDs {
			client, err := i.clientByID(ctx, s, clientID)
			if err != nil {
				return nil, err
			}

			if client != nil {
				keyNames[client.Key] = true
			}
		}
	}

	// Collect the key IDs
	var keyIDs []string
	for name := range keyNames {
		entry, err := s.Get(ctx, namedKeyConfigPath+name)
		if err != nil {
			return nil, err
		}

		if entry == nil {
			continue
		}

		var key namedKey
		if err := entry.DecodeJSON(&key); err != nil {
			return nil, err
		}
		for _, expirableKey := range key.KeyRing {
			keyIDs = append(keyIDs, expirableKey.KeyID)
		}
	}
	return keyIDs, nil
}

func (i *IdentityStore) pathOIDCAuthorize(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	state := d.Get("state").(string)

	// Validate the client ID
	clientID := d.Get("client_id").(string)
	if clientID == "" {
		return authResponse("", state, ErrAuthInvalidClientID, "client_id parameter is required")
	}
	client, err := i.clientByID(ctx, req.Storage, clientID)
	if err != nil {
		return authResponse("", state, ErrAuthServerError, err.Error())
	}
	if client == nil {
		return authResponse("", state, ErrAuthInvalidClientID, "client with client_id not found")
	}

	// Validate the redirect URI
	redirectURI := d.Get("redirect_uri").(string)
	if redirectURI == "" {
		return authResponse("", state, ErrAuthInvalidRequest, "redirect_uri parameter is required")
	}
	if !validRedirect(redirectURI, client.RedirectURIs) {
		return authResponse("", state, ErrAuthInvalidRedirectURI, "redirect_uri is not allowed for the client")
	}

	// Get the OIDC provider
	name := d.Get("name").(string)
	provider, err := i.getOIDCProvider(ctx, req.Storage, name)
	if err != nil {
		return authResponse("", state, ErrAuthServerError, err.Error())
	}
	if provider == nil {
		return authResponse("", state, ErrAuthInvalidRequest, "provider not found")
	}
	if !provider.allowedClientID(clientID) {
		return authResponse("", state, ErrAuthUnauthorizedClient, "client is not authorized to use the provider")
	}

	// We don't support the request or request_uri parameters. If they're provided,
	// the appropriate errors must be returned. For details, see the spec at:
	// https://openid.net/specs/openid-connect-core-1_0.html#RequestObject
	// https://openid.net/specs/openid-connect-core-1_0.html#RequestUriParameter
	if _, ok := d.Raw["request"]; ok {
		return authResponse("", "", ErrAuthRequestNotSupported, "request parameter is not supported")
	}
	if _, ok := d.Raw["request_uri"]; ok {
		return authResponse("", "", ErrAuthRequestURINotSupported, "request_uri parameter is not supported")
	}

	// Validate that a scope parameter is present and contains the openid scope value
	requestedScopes := strutil.ParseDedupAndSortStrings(d.Get("scope").(string), scopesDelimiter)
	if len(requestedScopes) == 0 || !strutil.StrListContains(requestedScopes, openIDScope) {
		return authResponse("", state, ErrAuthInvalidRequest,
			fmt.Sprintf("scope parameter must contain the %q value", openIDScope))
	}

	// Scope values that are not supported by the provider should be ignored
	scopes := make([]string, 0)
	for _, scope := range requestedScopes {
		if strutil.StrListContains(provider.ScopesSupported, scope) && scope != openIDScope {
			scopes = append(scopes, scope)
		}
	}

	// Validate the response type
	responseType := d.Get("response_type").(string)
	if responseType == "" {
		return authResponse("", state, ErrAuthInvalidRequest, "response_type parameter is required")
	}
	if responseType != "code" {
		return authResponse("", state, ErrAuthUnsupportedResponseType, "unsupported response_type value")
	}

	// Validate that there is an identity entity associated with the request
	if req.EntityID == "" {
		return authResponse("", state, ErrAuthAccessDenied, "identity entity must be associated with the request")
	}
	entity, err := i.MemDBEntityByID(req.EntityID, false)
	if err != nil {
		return authResponse("", state, ErrAuthServerError, err.Error())
	}
	if entity == nil {
		return authResponse("", state, ErrAuthAccessDenied, "identity entity associated with the request not found")
	}

	// Validate that the entity is a member of the client's assignments
	isMember, err := i.entityHasAssignment(ctx, req.Storage, entity, client.Assignments)
	if err != nil {
		return authResponse("", state, ErrAuthServerError, err.Error())
	}
	if !isMember {
		return authResponse("", state, ErrAuthAccessDenied, "identity entity not authorized by client assignment")
	}

	// A nonce is optional for the authorization code flow. If not
	// provided, the nonce claim will be omitted from the ID token.
	nonce := d.Get("nonce").(string)

	// Create the auth code cache entry
	authCodeEntry := &authCodeCacheEntry{
		provider:    name,
		clientID:    clientID,
		entityID:    entity.GetID(),
		redirectURI: redirectURI,
		nonce:       nonce,
		scopes:      scopes,
	}

	// Validate the Proof Key for Code Exchange (PKCE) code challenge and code challenge
	// method. PKCE is required for public clients and optional for confidential clients.
	// See details at https://datatracker.ietf.org/doc/html/rfc7636.
	codeChallengeRaw, okCodeChallenge := d.GetOk("code_challenge")
	if !okCodeChallenge && client.Type == public {
		return authResponse("", state, ErrAuthInvalidRequest, "PKCE is required for public clients")
	}
	if okCodeChallenge {
		codeChallenge := codeChallengeRaw.(string)

		// Validate the code challenge method
		codeChallengeMethod := d.Get("code_challenge_method").(string)
		switch codeChallengeMethod {
		case codeChallengeMethodPlain, codeChallengeMethodS256:
		case "":
			codeChallengeMethod = codeChallengeMethodPlain
		default:
			return authResponse("", state, ErrAuthInvalidRequest, "invalid code_challenge_method")
		}

		// Validate the code challenge
		if len(codeChallenge) < 43 || len(codeChallenge) > 128 {
			return authResponse("", state, ErrAuthInvalidRequest, "invalid code_challenge")
		}

		// Associate the code challenge and method with the authorization code.
		// This will be used to verify the code verifier in the token exchange.
		authCodeEntry.codeChallenge = codeChallenge
		authCodeEntry.codeChallengeMethod = codeChallengeMethod
	}

	// Validate the optional max_age parameter to check if an active re-authentication
	// of the user should occur. Re-authentication will be requested if the last time
	// the token actively authenticated exceeds the given max_age requirement. Returning
	// ErrAuthMaxAgeReAuthenticate will enforce the user to re-authenticate via the user agent.
	if maxAgeRaw, ok := d.GetOk("max_age"); ok {
		maxAge := maxAgeRaw.(int)
		if maxAge < 1 {
			return authResponse("", state, ErrAuthInvalidRequest, "max_age must be greater than zero")
		}

		// Look up the token associated with the request
		te, err := i.tokenStorer.LookupToken(ctx, req.ClientToken)
		if err != nil {
			return authResponse("", state, ErrAuthServerError, err.Error())
		}
		if te == nil {
			return authResponse("", state, ErrAuthAccessDenied, "token associated with request not found")
		}

		// Check if the token creation time violates the max age requirement
		now := time.Now().UTC()
		lastAuthTime := time.Unix(te.CreationTime, 0).UTC()
		secondsSince := int(now.Sub(lastAuthTime).Seconds())
		if secondsSince > maxAge {
			return authResponse("", state, ErrAuthMaxAgeReAuthenticate, "active re-authentication is required by max_age")
		}

		// Set the auth time to use for the auth_time claim in the token exchange
		authCodeEntry.authTime = lastAuthTime
	}

	// Generate the authorization code
	code, err := base62.Random(32)
	if err != nil {
		return authResponse("", state, ErrAuthServerError, err.Error())
	}

	// Get the namespace
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return authResponse("", state, ErrAuthServerError, err.Error())
	}

	// Cache the authorization code for a subsequent token exchange
	if err := i.oidcAuthCodeCache.SetDefault(ns, code, authCodeEntry); err != nil {
		return authResponse("", state, ErrAuthServerError, err.Error())
	}

	return authResponse(code, state, "", "")
}

// authResponse returns the OIDC Authentication Response. An error response is
// returned if the given error code is non-empty. For details, see spec at
//   - https://openid.net/specs/openid-connect-core-1_0.html#AuthResponse
//   - https://openid.net/specs/openid-connect-core-1_0.html#AuthError
func authResponse(code, state, errorCode, errorDescription string) (*logical.Response, error) {
	statusCode := http.StatusOK
	response := map[string]interface{}{
		"code":  code,
		"state": state,
	}

	// Set the error response and status code if error code isn't empty
	if errorCode != "" {
		statusCode = http.StatusBadRequest
		if errorCode == ErrAuthServerError {
			statusCode = http.StatusInternalServerError
		}

		response = map[string]interface{}{
			"error":             errorCode,
			"error_description": errorDescription,
			"state":             state,
		}
	}

	body, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:  statusCode,
			logical.HTTPRawBody:     body,
			logical.HTTPContentType: "application/json",
		},
	}, nil
}

func (i *IdentityStore) pathOIDCToken(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the namespace
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}

	// Get the OIDC provider
	name := d.Get("name").(string)
	provider, err := i.getOIDCProvider(ctx, req.Storage, name)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}
	if provider == nil {
		return tokenResponse(nil, ErrTokenInvalidRequest, "provider not found")
	}

	// client_secret_basic - Check for client credentials in the Authorization header
	clientID, clientSecret, okBasicAuth := basicAuth(req)
	if !okBasicAuth {
		// client_secret_post - Check for client credentials in the request body
		clientID = d.Get("client_id").(string)
		if clientID == "" {
			return tokenResponse(nil, ErrTokenInvalidRequest, "client_id parameter is required")
		}
		clientSecret = d.Get("client_secret").(string)
	}
	client, err := i.clientByID(ctx, req.Storage, clientID)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}
	if client == nil {
		i.Logger().Debug("client failed to authenticate with client not found", "client_id", clientID)
		return tokenResponse(nil, ErrTokenInvalidClient, "client failed to authenticate")
	}

	// Authenticate the client if it's a confidential client type.
	// Details at https://openid.net/specs/openid-connect-core-1_0.html#ClientAuthentication
	if client.Type == confidential &&
		subtle.ConstantTimeCompare([]byte(client.ClientSecret), []byte(clientSecret)) == 0 {
		i.Logger().Debug("client failed to authenticate with invalid client secret", "client_id", clientID)
		return tokenResponse(nil, ErrTokenInvalidClient, "client failed to authenticate")
	}

	// Validate that the client is authorized to use the provider
	if !provider.allowedClientID(clientID) {
		return tokenResponse(nil, ErrTokenInvalidClient, "client is not authorized to use the provider")
	}

	// Get the key that the client uses to sign ID tokens
	key, err := i.getNamedKey(ctx, req.Storage, client.Key)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}
	if key == nil {
		return tokenResponse(nil, ErrTokenServerError, fmt.Sprintf("client key %q not found", client.Key))
	}

	// Validate that the client is authorized to use the key
	if !strutil.StrListContains(key.AllowedClientIDs, "*") &&
		!strutil.StrListContains(key.AllowedClientIDs, clientID) {
		return tokenResponse(nil, ErrTokenInvalidClient, "client is not authorized to use the key")
	}

	// Validate the grant type
	grantType := d.Get("grant_type").(string)
	if grantType == "" {
		return tokenResponse(nil, ErrTokenInvalidRequest, "grant_type parameter is required")
	}
	if grantType != "authorization_code" {
		return tokenResponse(nil, ErrTokenUnsupportedGrantType, "unsupported grant_type value")
	}

	// Validate the authorization code
	code := d.Get("code").(string)
	if code == "" {
		return tokenResponse(nil, ErrTokenInvalidRequest, "code parameter is required")
	}

	// Get the authorization code entry and defer its deletion (single use)
	authCodeEntryRaw, ok, err := i.oidcAuthCodeCache.Get(ns, code)
	defer i.oidcAuthCodeCache.Delete(ns, code)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}
	if !ok {
		return tokenResponse(nil, ErrTokenInvalidGrant, "authorization grant is invalid or expired")
	}
	authCodeEntry, ok := authCodeEntryRaw.(*authCodeCacheEntry)
	if !ok {
		return tokenResponse(nil, ErrTokenServerError, "authorization grant is invalid or expired")
	}

	// Ensure the authorization code was issued to the authenticated client
	if authCodeEntry.clientID != clientID {
		return tokenResponse(nil, ErrTokenInvalidGrant, "authorization code was not issued to the client")
	}

	// Ensure the authorization code was issued by the provider
	if authCodeEntry.provider != name {
		return tokenResponse(nil, ErrTokenInvalidGrant, "authorization code was not issued by the provider")
	}

	// Ensure the redirect_uri parameter value is identical to the redirect_uri
	// parameter value that was included in the initial authorization request.
	redirectURI := d.Get("redirect_uri").(string)
	if redirectURI == "" {
		return tokenResponse(nil, ErrTokenInvalidRequest, "redirect_uri parameter is required")
	}
	if authCodeEntry.redirectURI != redirectURI {
		return tokenResponse(nil, ErrTokenInvalidGrant, "redirect_uri does not match the redirect_uri used in the authorization request")
	}

	// Get the entity associated with the initial authorization request
	entity, err := i.MemDBEntityByID(authCodeEntry.entityID, true)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}
	if entity == nil {
		return tokenResponse(nil, ErrTokenInvalidRequest, "identity entity associated with the request not found")
	}

	// Validate that the entity is a member of the client's assignments
	isMember, err := i.entityHasAssignment(ctx, req.Storage, entity, client.Assignments)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}
	if !isMember {
		return tokenResponse(nil, ErrTokenInvalidRequest, "identity entity not authorized by client assignment")
	}

	// Validate the PKCE code verifier. See details at
	// https://datatracker.ietf.org/doc/html/rfc7636#section-4.6.
	usedPKCE := authCodeUsedPKCE(authCodeEntry)
	codeVerifier := d.Get("code_verifier").(string)
	switch {
	case !usedPKCE && client.Type == public:
		return tokenResponse(nil, ErrTokenInvalidRequest, "PKCE is required for public clients")
	case !usedPKCE && codeVerifier != "":
		return tokenResponse(nil, ErrTokenInvalidRequest, "unexpected code_verifier for token exchange")
	case usedPKCE && codeVerifier == "":
		return tokenResponse(nil, ErrTokenInvalidRequest, "expected code_verifier for token exchange")
	case usedPKCE:
		codeChallenge, err := computeCodeChallenge(codeVerifier, authCodeEntry.codeChallengeMethod)
		if err != nil {
			return tokenResponse(nil, ErrTokenServerError, err.Error())
		}

		if subtle.ConstantTimeCompare([]byte(codeChallenge), []byte(authCodeEntry.codeChallenge)) == 0 {
			return tokenResponse(nil, ErrTokenInvalidGrant, "invalid code_verifier for token exchange")
		}
	}

	// The access token is a Vault batch token with a policy that only
	// provides access to the issuing provider's userinfo endpoint.
	accessTokenIssuedAt := time.Now()
	accessTokenExpiry := accessTokenIssuedAt.Add(client.AccessTokenTTL)
	accessToken := &logical.TokenEntry{
		Type:               logical.TokenTypeBatch,
		NamespaceID:        ns.ID,
		Path:               req.Path,
		TTL:                client.AccessTokenTTL,
		CreationTime:       accessTokenIssuedAt.Unix(),
		EntityID:           entity.ID,
		NoIdentityPolicies: true,
		Meta: map[string]string{
			"oidc_token_type": "access token",
		},
		InternalMeta: map[string]string{
			accessTokenClientIDMeta: client.ClientID,
			accessTokenScopesMeta:   strings.Join(authCodeEntry.scopes, scopesDelimiter),
		},
		InlinePolicy: fmt.Sprintf(`
			path "identity/oidc/provider/%s/userinfo" {
				capabilities = ["read", "update"]
			}
		`, name),
	}
	err = i.tokenStorer.CreateToken(ctx, accessToken)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}

	// Compute the access token hash claim (at_hash)
	atHash, err := computeHashClaim(key.Algorithm, accessToken.ID)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}

	// Compute the authorization code hash claim (c_hash)
	cHash, err := computeHashClaim(key.Algorithm, code)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}

	// Set the ID token claims
	idTokenIssuedAt := time.Now()
	idTokenExpiry := idTokenIssuedAt.Add(client.IDTokenTTL)
	idToken := idToken{
		Namespace:       ns.ID,
		Issuer:          provider.effectiveIssuer,
		Subject:         authCodeEntry.entityID,
		Audience:        authCodeEntry.clientID,
		Nonce:           authCodeEntry.nonce,
		Expiry:          idTokenExpiry.Unix(),
		IssuedAt:        idTokenIssuedAt.Unix(),
		AccessTokenHash: atHash,
		CodeHash:        cHash,
	}

	// Add the auth_time claim if it's not the zero time instant
	if !authCodeEntry.authTime.IsZero() {
		idToken.AuthTime = authCodeEntry.authTime.Unix()
	}

	// Populate each of the requested scope templates
	templates, conflict, err := i.populateScopeTemplates(ctx, req.Storage, ns, entity, authCodeEntry.scopes...)
	if !conflict && err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}
	if conflict && err != nil {
		return tokenResponse(nil, ErrTokenInvalidRequest, err.Error())
	}

	// Generate the ID token payload
	payload, err := idToken.generatePayload(i.Logger(), templates...)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}

	// Sign the ID token using the client's key
	signedIDToken, err := key.signPayload(payload)
	if err != nil {
		return tokenResponse(nil, ErrTokenServerError, err.Error())
	}

	return tokenResponse(map[string]interface{}{
		"token_type":   "Bearer",
		"access_token": accessToken.ID,
		"id_token":     signedIDToken,
		"expires_in":   int64(accessTokenExpiry.Sub(accessTokenIssuedAt).Seconds()),
	}, "", "")
}

// tokenResponse returns the OIDC Token Response. An error response is
// returned if the given error code is non-empty. For details, see spec at
//   - https://openid.net/specs/openid-connect-core-1_0.html#TokenResponse
//   - https://openid.net/specs/openid-connect-core-1_0.html#TokenErrorResponse
func tokenResponse(response map[string]interface{}, errorCode, errorDescription string) (*logical.Response, error) {
	statusCode := http.StatusOK

	// Set the error response and status code if error code isn't empty
	if errorCode != "" {
		switch errorCode {
		case ErrTokenInvalidClient:
			statusCode = http.StatusUnauthorized
		case ErrTokenServerError:
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusBadRequest
		}

		response = map[string]interface{}{
			"error":             errorCode,
			"error_description": errorDescription,
		}
	}

	body, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		logical.HTTPStatusCode:  statusCode,
		logical.HTTPRawBody:     body,
		logical.HTTPContentType: "application/json",

		// Token responses must include the following HTTP response headers
		// https://openid.net/specs/openid-connect-core-1_0.html#TokenResponse
		logical.HTTPCacheControlHeader: "no-store",
		logical.HTTPPragmaHeader:       "no-cache",
	}

	// Set the WWW-Authenticate response header when returning the
	// invalid_client error code per the OAuth 2.0 spec at
	// https://datatracker.ietf.org/doc/html/rfc6749#section-5.2
	if errorCode == ErrTokenInvalidClient {
		data[logical.HTTPWWWAuthenticateHeader] = "Basic"
	}

	return &logical.Response{
		Data: data,
	}, nil
}

func (i *IdentityStore) pathOIDCUserInfo(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// Get the namespace
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return userInfoResponse(nil, ErrUserInfoServerError, err.Error())
	}

	// Get the OIDC provider
	name := d.Get("name").(string)
	provider, err := i.getOIDCProvider(ctx, req.Storage, name)
	if err != nil {
		return userInfoResponse(nil, ErrUserInfoServerError, err.Error())
	}
	if provider == nil {
		return userInfoResponse(nil, ErrUserInfoInvalidRequest, "provider not found")
	}

	// Validate that the access token was sent as a Bearer token
	if req.ClientTokenSource != logical.ClientTokenFromAuthzHeader {
		return userInfoResponse(nil, ErrUserInfoInvalidToken, "access token must be sent as a Bearer token")
	}

	// Look up the access token
	te, err := i.tokenStorer.LookupToken(ctx, req.ClientToken)
	if err != nil {
		return userInfoResponse(nil, ErrUserInfoServerError, err.Error())
	}
	if te == nil {
		return userInfoResponse(nil, ErrUserInfoInvalidToken, "access token is expired")
	}
	if te.Type != logical.TokenTypeBatch {
		return userInfoResponse(nil, ErrUserInfoInvalidToken, "access token is malformed or invalid")
	}

	// Get the client ID that originated the request from the token metadata
	clientID, ok := te.InternalMeta[accessTokenClientIDMeta]
	if !ok {
		return userInfoResponse(nil, ErrUserInfoServerError, "expected client ID in token metadata")
	}
	client, err := i.clientByID(ctx, req.Storage, clientID)
	if err != nil {
		return userInfoResponse(nil, ErrUserInfoServerError, err.Error())
	}
	if client == nil {
		return userInfoResponse(nil, ErrUserInfoAccessDenied, "client not found")
	}

	// Validate that there is an identity entity associated with the request
	if req.EntityID == "" {
		return userInfoResponse(nil, ErrUserInfoAccessDenied, "identity entity must be associated with the request")
	}
	entity, err := i.MemDBEntityByID(req.EntityID, false)
	if err != nil {
		return userInfoResponse(nil, ErrUserInfoServerError, err.Error())
	}
	if entity == nil {
		return userInfoResponse(nil, ErrUserInfoAccessDenied, "identity entity associated with the request not found")
	}

	// Validate that the entity is a member of the client's assignments
	isMember, err := i.entityHasAssignment(ctx, req.Storage, entity, client.Assignments)
	if err != nil {
		return userInfoResponse(nil, ErrUserInfoServerError, err.Error())
	}
	if !isMember {
		return userInfoResponse(nil, ErrUserInfoAccessDenied, "identity entity not authorized by client assignment")
	}

	// Validate that the client is authorized to use the provider
	if !provider.allowedClientID(clientID) {
		return userInfoResponse(nil, ErrUserInfoAccessDenied, "client is not authorized to use the provider")
	}

	claims := map[string]interface{}{
		// The subject claim must always be in the response
		"sub": entity.ID,
	}

	// Get the scopes for the access token
	tokenScopes, ok := te.InternalMeta[accessTokenScopesMeta]
	if !ok || len(tokenScopes) == 0 {
		return userInfoResponse(claims, "", "")
	}
	parsedScopes := strutil.ParseStringSlice(tokenScopes, scopesDelimiter)

	// Scope values that are not supported by the provider should be ignored
	scopes := make([]string, 0)
	for _, scope := range parsedScopes {
		if strutil.StrListContains(provider.ScopesSupported, scope) {
			scopes = append(scopes, scope)
		}
	}

	// Populate each of the token's scope templates
	templates, conflict, err := i.populateScopeTemplates(ctx, req.Storage, ns, entity, scopes...)
	if !conflict && err != nil {
		return userInfoResponse(nil, ErrUserInfoServerError, err.Error())
	}
	if conflict && err != nil {
		return userInfoResponse(nil, ErrUserInfoInvalidRequest, err.Error())
	}

	// Merge all of the populated JSON scope templates into claims
	if err := mergeJSONTemplates(i.Logger(), claims, templates...); err != nil {
		return userInfoResponse(nil, ErrUserInfoServerError, err.Error())
	}

	return userInfoResponse(claims, "", "")
}

// userInfoResponse returns the OIDC UserInfo Response. An error response is
// returned if the given error code is non-empty. For details, see spec at
//   - https://openid.net/specs/openid-connect-core-1_0.html#UserInfoResponse
//   - https://openid.net/specs/openid-connect-core-1_0.html#UserInfoError
func userInfoResponse(response map[string]interface{}, errorCode, errorDescription string) (*logical.Response, error) {
	statusCode := http.StatusOK

	// Set the error response and status code if error code isn't empty
	if errorCode != "" {
		switch errorCode {
		case ErrUserInfoInvalidRequest:
			statusCode = http.StatusBadRequest
		case ErrUserInfoInvalidToken:
			statusCode = http.StatusUnauthorized
		case ErrUserInfoAccessDenied:
			statusCode = http.StatusForbidden
		case ErrUserInfoServerError:
			statusCode = http.StatusInternalServerError
		}

		response = map[string]interface{}{
			"error":             errorCode,
			"error_description": errorDescription,
		}
	}

	body, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{
		logical.HTTPStatusCode:  statusCode,
		logical.HTTPRawBody:     body,
		logical.HTTPContentType: "application/json",
	}

	// Set the WWW-Authenticate response header when returning error codes
	// defined in https://datatracker.ietf.org/doc/html/rfc6750#section-3
	if errorCode == ErrUserInfoInvalidRequest || errorCode == ErrUserInfoInvalidToken {
		data[logical.HTTPWWWAuthenticateHeader] = fmt.Sprintf("Bearer error=%q,error_description=%q",
			errorCode, errorDescription)
	}

	return &logical.Response{
		Data: data,
	}, nil
}

// getScopeTemplates returns a mapping from scope names to
// their templates for each of the given scopes.
func (i *IdentityStore) getScopeTemplates(ctx context.Context, s logical.Storage, scopes ...string) (map[string]string, error) {
	templates := make(map[string]string)
	for _, name := range scopes {
		if name == openIDScope {
			// No template for the openid scope
			continue
		}

		// Get the scope template
		scope, err := i.getOIDCScope(ctx, s, name)
		if err != nil {
			return nil, err
		}
		if scope == nil {
			// Scope values used that are not understood by an implementation should be ignored.
			// https://openid.net/specs/openid-connect-core-1_0.html#AuthRequest
			continue
		}
		templates[name] = scope.Template
	}

	return templates, nil
}

// populateScopeTemplates populates the templates for each of the passed scopes.
// Returns a slice of the populated JSON template strings and a bool to indicate
// if a conflict in scope template claims occurred.
func (i *IdentityStore) populateScopeTemplates(ctx context.Context, s logical.Storage, ns *namespace.Namespace, entity *identity.Entity, scopes ...string) ([]string, bool, error) {
	// Gather the templates for each scope
	templates, err := i.getScopeTemplates(ctx, s, scopes...)
	if err != nil {
		return nil, false, err
	}

	// Get the groups for the entity
	groups, inheritedGroups, err := i.groupsByEntityID(entity.ID)
	if err != nil {
		return nil, false, err
	}
	groups = append(groups, inheritedGroups...)

	claimsToScopes := make(map[string]string)
	populatedTemplates := make([]string, 0)
	for scope, template := range templates {
		// Parse and integrate the populated template. Structural errors with the template
		// should be caught during configuration. Errors found during runtime will be logged.
		_, populatedTemplate, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
			Mode:        identitytpl.JSONTemplating,
			String:      template,
			Entity:      identity.ToSDKEntity(entity),
			Groups:      identity.ToSDKGroups(groups),
			NamespaceID: ns.ID,
		})
		if err != nil {
			i.Logger().Warn("error populating OIDC token template", "scope", scope,
				"template", template, "error", err)
		}

		if populatedTemplate != "" {
			claimsMap := make(map[string]interface{})
			if err := json.Unmarshal([]byte(populatedTemplate), &claimsMap); err != nil {
				i.Logger().Warn("error parsing OIDC template", "template", template, "err", err)
			}

			// Check top-level claim keys for conflicts with other scopes
			for claimKey := range claimsMap {
				if conflictScope, ok := claimsToScopes[claimKey]; ok {
					return nil, true, fmt.Errorf("found scopes with conflicting top-level claim: claim %q in scopes %q, %q",
						claimKey, scope, conflictScope)
				}
				claimsToScopes[claimKey] = scope
			}

			populatedTemplates = append(populatedTemplates, populatedTemplate)
		}
	}

	return populatedTemplates, false, nil
}

// entityHasAssignment returns true if the entity is enabled and a member of any
// of the assignments' groups or entities. Otherwise, returns false or an error.
func (i *IdentityStore) entityHasAssignment(ctx context.Context, s logical.Storage, entity *identity.Entity, assignments []string) (bool, error) {
	if entity.GetDisabled() {
		return false, nil
	}

	if strutil.StrListContains(assignments, allowAllAssignmentName) {
		return true, nil
	}

	// Get the group IDs that the entity is a member of
	groups, inheritedGroups, err := i.groupsByEntityID(entity.GetID())
	if err != nil {
		return false, err
	}
	entityGroupIDs := make(map[string]bool)
	for _, group := range append(groups, inheritedGroups...) {
		entityGroupIDs[group.GetID()] = true
	}

	for _, a := range assignments {
		assignment, err := i.getOIDCAssignment(ctx, s, a)
		if err != nil {
			return false, err
		}
		if assignment == nil {
			return false, fmt.Errorf("client assignment %q not found", a)
		}

		// Check if the entity is a member of any groups in the assignment
		for _, id := range assignment.GroupIDs {
			if entityGroupIDs[id] {
				return true, nil
			}
		}

		// Check if the entity is a member of the assignment's entities
		if strutil.StrListContains(assignment.EntityIDs, entity.GetID()) {
			return true, nil
		}
	}

	return false, nil
}

func defaultOIDCProvider() provider {
	return provider{
		AllowedClientIDs: []string{"*"},
		ScopesSupported:  []string{},
	}
}

func defaultOIDCKey() namedKey {
	return namedKey{
		Algorithm:        "RS256",
		VerificationTTL:  24 * time.Hour,
		RotationPeriod:   24 * time.Hour,
		NextRotation:     time.Now().Add(24 * time.Hour),
		AllowedClientIDs: []string{"*"},
	}
}

func allowAllAssignment() assignment {
	return assignment{
		EntityIDs: []string{"*"},
		GroupIDs:  []string{"*"},
	}
}

func (i *IdentityStore) storeOIDCDefaultResources(ctx context.Context, view logical.Storage) error {
	// Store the default provider
	storageKey := providerPath + defaultProviderName
	entry, err := view.Get(ctx, storageKey)
	if err != nil {
		return err
	}
	if entry == nil {
		entry, err := logical.StorageEntryJSON(storageKey, defaultOIDCProvider())
		if err != nil {
			return err
		}
		if err := view.Put(ctx, entry); err != nil {
			return err
		}
		i.Logger().Debug("wrote OIDC default provider")
	}

	if _, err := i.ensureDefaultKey(ctx, view); err != nil {
		return fmt.Errorf("error writing default key to storage: %w", err)
	}

	// Store the allow all assignment
	storageKey = assignmentPath + allowAllAssignmentName
	entry, err = view.Get(ctx, storageKey)
	if err != nil {
		return err
	}
	if entry == nil {
		entry, err := logical.StorageEntryJSON(storageKey, allowAllAssignment())
		if err != nil {
			return err
		}
		if err := view.Put(ctx, entry); err != nil {
			return err
		}
		i.Logger().Debug("wrote OIDC allow_all assignment")
	}

	return nil
}

// ensureDefaultKey ensures that the OIDC default key is written to storage. If no
// error is returned, callers can be sure that it exists in storage. Note that it
// only writes the key's configuration to storage and does not generate key material
// for its current and next keys.
func (i *IdentityStore) ensureDefaultKey(ctx context.Context, storage logical.Storage) (*namedKey, error) {
	key, err := i.getNamedKey(ctx, storage, defaultKeyName)
	if err != nil {
		return nil, err
	}
	if key != nil {
		return key, nil
	}

	// The default key doesn't exist. Write it to storage.
	defaultKey := defaultOIDCKey()
	entry, err := logical.StorageEntryJSON(namedKeyConfigPath+defaultKeyName, defaultKey)
	if err != nil {
		return nil, err
	}
	if err := storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	i.Logger().Debug("wrote OIDC default key")
	return &defaultKey, nil
}

// lazyGenerateDefaultKey generates key material for the OIDC default key's current and
// next key if it hasn't already been generated. Must be called with the oidcLock write
// lock held.
func (i *IdentityStore) lazyGenerateDefaultKey(ctx context.Context, storage logical.Storage) error {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return err
	}

	defaultKey, err := i.ensureDefaultKey(ctx, storage)
	if err != nil {
		return err
	}

	if defaultKey.SigningKey == nil {
		if err := defaultKey.generateAndSetKey(ctx, i.Logger(), storage); err != nil {
			return err
		}
		if err := defaultKey.generateAndSetNextKey(ctx, i.Logger(), storage); err != nil {
			return err
		}

		entry, err := logical.StorageEntryJSON(namedKeyConfigPath+defaultKeyName, defaultKey)
		if err != nil {
			return err
		}
		if err := storage.Put(ctx, entry); err != nil {
			return err
		}

		if err := i.oidcCache.Flush(ns); err != nil {
			return err
		}
	}

	return nil
}

func (i *IdentityStore) loadOIDCClients(ctx context.Context) error {
	i.logger.Debug("identity loading OIDC clients")

	clients, err := i.view.List(ctx, clientPath)
	if err != nil {
		return err
	}

	txn := i.db.Txn(true)
	defer txn.Abort()
	for _, name := range clients {
		entry, err := i.view.Get(ctx, clientPath+name)
		if err != nil {
			return err
		}
		if entry == nil {
			continue
		}

		var client client
		if err := entry.DecodeJSON(&client); err != nil {
			return err
		}

		if err := i.memDBUpsertClientInTxn(txn, &client); err != nil {
			return err
		}
	}
	txn.Commit()

	return nil
}

// clientByID returns the client with the given ID.
func (i *IdentityStore) clientByID(ctx context.Context, s logical.Storage, id string) (*client, error) {
	// Read the client from memdb
	client, err := i.memDBClientByID(id)
	if err != nil {
		return nil, err
	}
	if client != nil {
		return client, nil
	}

	// Fall back to reading the client from storage
	client, err = i.storageClientByID(ctx, s, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, nil
	}

	// Upsert the client in memdb
	txn := i.db.Txn(true)
	defer txn.Abort()
	if err := i.memDBUpsertClientInTxn(txn, client); err != nil {
		i.logger.Debug("failed to upsert client in memdb", "error", err)
		return client, nil
	}
	txn.Commit()

	return client, nil
}

// clientByName returns the client with the given name.
func (i *IdentityStore) clientByName(ctx context.Context, s logical.Storage, name string) (*client, error) {
	// Read the client from memdb
	client, err := i.memDBClientByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if client != nil {
		return client, nil
	}

	// Fall back to reading the client from storage
	client, err = i.storageClientByName(ctx, s, name)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, nil
	}

	// Upsert the client in memdb
	txn := i.db.Txn(true)
	defer txn.Abort()
	if err := i.memDBUpsertClientInTxn(txn, client); err != nil {
		i.logger.Debug("failed to upsert client in memdb", "error", err)
		return client, nil
	}
	txn.Commit()

	return client, nil
}

// memDBClientByID returns the client with the given ID from memdb.
func (i *IdentityStore) memDBClientByID(id string) (*client, error) {
	if id == "" {
		return nil, errors.New("missing client ID")
	}

	txn := i.db.Txn(false)

	return i.memDBClientByIDInTxn(txn, id)
}

// memDBClientByIDInTxn returns the client with the given ID from memdb using the given txn.
func (i *IdentityStore) memDBClientByIDInTxn(txn *memdb.Txn, id string) (*client, error) {
	if id == "" {
		return nil, errors.New("missing client ID")
	}

	if txn == nil {
		return nil, errors.New("txn is nil")
	}

	clientRaw, err := txn.First(oidcClientsTable, "id", id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch client from memdb using ID: %w", err)
	}
	if clientRaw == nil {
		return nil, nil
	}

	client, ok := clientRaw.(*client)
	if !ok {
		return nil, errors.New("unexpected client type")
	}

	return client, nil
}

// memDBClientByName returns the client with the given name from memdb.
func (i *IdentityStore) memDBClientByName(ctx context.Context, name string) (*client, error) {
	if name == "" {
		return nil, errors.New("missing client name")
	}

	txn := i.db.Txn(false)

	return i.memDBClientByNameInTxn(ctx, txn, name)
}

// memDBClientByNameInTxn returns the client with the given ID from memdb using the given txn.
func (i *IdentityStore) memDBClientByNameInTxn(ctx context.Context, txn *memdb.Txn, name string) (*client, error) {
	if name == "" {
		return nil, errors.New("missing client name")
	}

	if txn == nil {
		return nil, errors.New("txn is nil")
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	clientRaw, err := txn.First(oidcClientsTable, "name", ns.ID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch client from memdb using name: %w", err)
	}
	if clientRaw == nil {
		return nil, nil
	}

	client, ok := clientRaw.(*client)
	if !ok {
		return nil, errors.New("unexpected client type")
	}

	return client, nil
}

// memDBDeleteClientByName deletes the client with the given name from memdb.
func (i *IdentityStore) memDBDeleteClientByName(ctx context.Context, name string) error {
	if name == "" {
		return errors.New("missing client name")
	}

	txn := i.db.Txn(true)
	defer txn.Abort()

	if err := i.memDBDeleteClientByNameInTxn(ctx, txn, name); err != nil {
		return err
	}

	txn.Commit()

	return nil
}

// memDBDeleteClientByNameInTxn deletes the client with name from memdb using the given txn.
func (i *IdentityStore) memDBDeleteClientByNameInTxn(ctx context.Context, txn *memdb.Txn, name string) error {
	if name == "" {
		return errors.New("missing client name")
	}

	if txn == nil {
		return errors.New("txn is nil")
	}

	client, err := i.memDBClientByNameInTxn(ctx, txn, name)
	if err != nil {
		return err
	}
	if client == nil {
		return nil
	}

	if err := txn.Delete(oidcClientsTable, client); err != nil {
		return fmt.Errorf("failed to delete client from memdb: %w", err)
	}

	return nil
}

// memDBUpsertClientInTxn creates or updates the given client in memdb using the given txn.
func (i *IdentityStore) memDBUpsertClientInTxn(txn *memdb.Txn, client *client) error {
	if client == nil {
		return errors.New("client is nil")
	}

	if txn == nil {
		return errors.New("nil txn")
	}

	clientRaw, err := txn.First(oidcClientsTable, "id", client.ClientID)
	if err != nil {
		return fmt.Errorf("failed to lookup client from memdb using ID: %w", err)
	}

	if clientRaw != nil {
		err = txn.Delete(oidcClientsTable, clientRaw)
		if err != nil {
			return fmt.Errorf("failed to delete client from memdb: %w", err)
		}
	}

	if err := txn.Insert(oidcClientsTable, client); err != nil {
		return fmt.Errorf("failed to update client in memdb: %w", err)
	}

	return nil
}

// storageClientByName returns the client with name from the given logical storage.
func (i *IdentityStore) storageClientByName(ctx context.Context, s logical.Storage, name string) (*client, error) {
	entry, err := s.Get(ctx, clientPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var client client
	if err := entry.DecodeJSON(&client); err != nil {
		return nil, err
	}

	return &client, nil
}

// storageClientByID returns the client with ID from the given logical storage.
func (i *IdentityStore) storageClientByID(ctx context.Context, s logical.Storage, id string) (*client, error) {
	clients, err := s.List(ctx, clientPath)
	if err != nil {
		return nil, err
	}

	for _, name := range clients {
		client, err := i.storageClientByName(ctx, s, name)
		if err != nil {
			return nil, err
		}
		if client == nil {
			continue
		}

		if client.ClientID == id {
			return client, nil
		}
	}

	return nil, nil
}

func (i *IdentityStore) listClients(ctx context.Context, s logical.Storage) ([]*client, error) {
	clientNames, err := s.List(ctx, clientPath)
	if err != nil {
		return nil, err
	}

	var clients []*client
	for _, name := range clientNames {
		entry, err := s.Get(ctx, clientPath+name)
		if err != nil {
			return nil, err
		}
		if entry == nil {
			continue
		}

		var client client
		if err := entry.DecodeJSON(&client); err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}

	return clients, nil
}
