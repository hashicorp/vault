package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/identitytpl"
	"github.com/hashicorp/vault/sdk/logical"
	"gopkg.in/square/go-jose.v2"
)

const (
	// OIDC-related constants
	openIDScope = "openid"

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

	// The following errors are used by the UI for specific behavior of
	// the OIDC specification. Any changes to their values must come with
	// a corresponding change in the UI code.
	ErrAuthInvalidClientID      = "invalid_client_id"
	ErrAuthInvalidRedirectURI   = "invalid_redirect_uri"
	ErrAuthMaxAgeReAuthenticate = "max_age_violation"
)

type assignment struct {
	Groups   []string `json:"groups"`
	Entities []string `json:"entities"`
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
	RedirectURIs   []string `json:"redirect_uris"`
	Assignments    []string `json:"assignments"`
	Key            string   `json:"key"`
	IDTokenTTL     int      `json:"id_token_ttl"`
	AccessTokenTTL int      `json:"access_token_ttl"`

	// Generated values that are used in OIDC endpoints
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type provider struct {
	Issuer           string   `json:"issuer"`
	AllowedClientIDs []string `json:"allowed_client_ids"`
	Scopes           []string `json:"scopes"`

	// effectiveIssuer is a calculated field and will be either Issuer (if
	// that's set) or the Vault instance's api_addr.
	effectiveIssuer string
}

type providerDiscovery struct {
	AuthorizationEndpoint string   `json:"authorization_endpoint"`
	IDTokenAlgs           []string `json:"id_token_signing_alg_values_supported"`
	Issuer                string   `json:"issuer"`
	Keys                  string   `json:"jwks_uri"`
	ResponseTypes         []string `json:"response_types_supported"`
	Scopes                []string `json:"scopes_supported"`
	Subjects              []string `json:"subject_types_supported"`
	TokenEndpoint         string   `json:"token_endpoint"`
	UserinfoEndpoint      string   `json:"userinfo_endpoint"`
}

type authCodeCacheEntry struct {
	entityID string
	nonce    string
	scopes   []string
	authTime time.Time
}

func oidcProviderPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "oidc/assignment/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the assignment",
				},
				"entities": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Comma separated string or array of identity entity names",
				},
				"groups": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Comma separated string or array of identity group names",
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
					Description: "A reference to a named key resource. Cannot be modified after creation.",
					Required:    true,
				},
				"id_token_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "The time-to-live for ID tokens obtained by the client.",
				},
				"access_token_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "The time-to-live for access tokens obtained by the client.",
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
				"scopes": {
					Type:        framework.TypeCommaStringSlice,
					Description: "The scopes available for requesting on the provider",
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
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: i.pathOIDCListProvider,
				},
			},
			HelpSynopsis:    "List OIDC providers",
			HelpDescription: "List all configured OIDC providers in the identity backend.",
		},
		{
			Pattern: "oidc/provider/" + framework.GenericNameRegex("name") + "/.well-known/openid-configuration",
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the provider",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.pathOIDCProviderDiscovery,
			},
			HelpSynopsis:    "Query OIDC configurations",
			HelpDescription: "Query this path to retrieve the configured OIDC Issuer and Keys endpoints, response types, subject types, and signing algorithms used by the OIDC backend.",
		},
		{
			Pattern: "oidc/provider/" + framework.GenericNameRegex("name") + "/.well-known/keys",
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the provider",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.pathOIDCReadProviderPublicKeys,
			},
			HelpSynopsis:    "Retrieve public keys",
			HelpDescription: "Returns the public portion of keys for a named OIDC provider. Clients can use them to validate the authenticity of an ID token.",
		},
		{
			Pattern: "oidc/provider/" + framework.GenericNameRegex("name") + "/authorize",
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the provider",
				},
				"client_id": {
					Type:        framework.TypeString,
					Description: "The ID of the requesting client.",
					Required:    true,
				},
				"scope": {
					Type:        framework.TypeString,
					Description: "A space-delimited, case-sensitive list of scopes to be requested. The 'openid' scope is required.",
					Required:    true,
				},
				"redirect_uri": {
					Type:        framework.TypeString,
					Description: "The redirection URI to which the response will be sent.",
					Required:    true,
				},
				"response_type": {
					Type:        framework.TypeString,
					Description: "The OIDC authentication flow to be used. The following response types are supported: 'code'",
					Required:    true,
				},
				"state": {
					Type:        framework.TypeString,
					Description: "The value used to maintain state between the authentication request and client.",
					Required:    true,
				},
				"nonce": {
					Type:        framework.TypeString,
					Description: "The value that will be returned in the ID token nonce claim after a token exchange.",
					Required:    true,
				},
				"max_age": {
					Type:        framework.TypeInt,
					Description: "The allowable elapsed time in seconds since the last time the end-user was actively authenticated.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback:                    i.pathOIDCAuthorize,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: false,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    i.pathOIDCAuthorize,
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: false,
				},
			},
			HelpSynopsis:    "Provides the OIDC Authorization Endpoint.",
			HelpDescription: "The OIDC Authorization Endpoint performs authentication and authorization by using request parameters defined by OpenID Connect (OIDC).",
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
			for _, a := range tempProvider.Scopes {
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

	if entitiesRaw, ok := d.GetOk("entities"); ok {
		assignment.Entities = entitiesRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		assignment.Entities = d.GetDefaultOrZero("entities").([]string)
	}

	if groupsRaw, ok := d.GetOk("groups"); ok {
		assignment.Groups = groupsRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		assignment.Groups = d.GetDefaultOrZero("groups").([]string)
	}

	// remove duplicates and lowercase entities and groups
	assignment.Entities = strutil.RemoveDuplicates(assignment.Entities, true)
	assignment.Groups = strutil.RemoveDuplicates(assignment.Groups, true)

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
			"groups":   assignment.Groups,
			"entities": assignment.Entities,
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
		scope.Description = d.GetDefaultOrZero("description").(string)
	}

	if templateRaw, ok := d.GetOk("template"); ok {
		scope.Template = templateRaw.(string)
	} else if req.Operation == logical.CreateOperation {
		scope.Template = d.GetDefaultOrZero("template").(string)
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
			if strutil.StrListContains(requiredClaims, key) {
				return logical.ErrorResponse(`top level key %q not allowed. Restricted keys: %s`,
					key, strings.Join(requiredClaims, ", ")), nil
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

	entry, err := req.Storage.Get(ctx, scopePath+name)
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
	return &logical.Response{
		Data: map[string]interface{}{
			"template":    scope.Template,
			"description": scope.Description,
		},
	}, nil
}

// pathOIDCDeleteScope is used to delete an scope
func (i *IdentityStore) pathOIDCDeleteScope(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	targetScopeName := d.Get("name").(string)

	providerNames, err := i.providersReferencingTargetScopeName(ctx, req, targetScopeName)
	if err != nil {
		return nil, err
	}

	if len(providerNames) > 0 {
		errorMessage := fmt.Sprintf("unable to delete scope %q because it is currently referenced by these providers: %s",
			targetScopeName, strings.Join(providerNames, ", "))
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

	if client.Key == "" {
		return logical.ErrorResponse("the key parameter is required"), nil
	}

	// enforce key existence on client creation
	entry, err := req.Storage.Get(ctx, namedKeyConfigPath+client.Key)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse("key %q does not exist", client.Key), nil
	}

	if idTokenTTLRaw, ok := d.GetOk("id_token_ttl"); ok {
		client.IDTokenTTL = idTokenTTLRaw.(int)
	} else if req.Operation == logical.CreateOperation {
		client.IDTokenTTL = d.Get("id_token_ttl").(int)
	}

	if accessTokenTTLRaw, ok := d.GetOk("access_token_ttl"); ok {
		client.AccessTokenTTL = accessTokenTTLRaw.(int)
	} else if req.Operation == logical.CreateOperation {
		client.AccessTokenTTL = d.Get("access_token_ttl").(int)
	}

	if client.ClientID == "" {
		// generate client_id
		clientID, err := base62.Random(32)
		if err != nil {
			return nil, err
		}
		client.ClientID = clientID
	}

	if client.ClientSecret == "" {
		// generate client_secret
		clientSecret, err := base62.Random(64)
		if err != nil {
			return nil, err
		}
		client.ClientSecret = clientSecret
	}

	// invalidate the cached client in memdb
	if err := i.memDBDeleteClientByName(ctx, name); err != nil {
		return nil, err
	}

	// store client
	entry, err = logical.StorageEntryJSON(clientPath+name, client)
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
	clients, err := req.Storage.List(ctx, clientPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(clients), nil
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

	return &logical.Response{
		Data: map[string]interface{}{
			"redirect_uris":    client.RedirectURIs,
			"assignments":      client.Assignments,
			"key":              client.Key,
			"id_token_ttl":     client.IDTokenTTL,
			"access_token_ttl": client.AccessTokenTTL,
			"client_id":        client.ClientID,
			"client_secret":    client.ClientSecret,
		},
	}, nil
}

// pathOIDCDeleteClient is used to delete an client
func (i *IdentityStore) pathOIDCDeleteClient(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

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
	resp := &logical.Response{}
	name := d.Get("name").(string)

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
		provider.Issuer = d.GetDefaultOrZero("issuer").(string)
	}

	if allowedClientIDsRaw, ok := d.GetOk("allowed_client_ids"); ok {
		provider.AllowedClientIDs = allowedClientIDsRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		provider.AllowedClientIDs = d.GetDefaultOrZero("allowed_client_ids").([]string)
	}

	if scopesRaw, ok := d.GetOk("scopes"); ok {
		provider.Scopes = scopesRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		provider.Scopes = d.GetDefaultOrZero("scopes").([]string)
	}

	// remove duplicate allowed client IDs and scopes
	provider.AllowedClientIDs = strutil.RemoveDuplicates(provider.AllowedClientIDs, false)
	provider.Scopes = strutil.RemoveDuplicates(provider.Scopes, false)

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
	for _, scopeName := range provider.Scopes {
		entry, err := req.Storage.Get(ctx, scopePath+scopeName)
		if err != nil {
			return nil, err
		}
		// enforce scope existence on provider create and update
		if entry == nil {
			return logical.ErrorResponse("scope %q does not exist", scopeName), nil
		}

		// ensure no two templates have the same top-level keys
		var storedScope scope
		if err := entry.DecodeJSON(&storedScope); err != nil {
			return nil, err
		}

		_, populatedTemplate, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
			Mode:   identitytpl.JSONTemplating,
			String: storedScope.Template,
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

	return resp, req.Storage.Put(ctx, entry)
}

// pathOIDCListProvider is used to list named providers
func (i *IdentityStore) pathOIDCListProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	providers, err := req.Storage.List(ctx, providerPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(providers), nil
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
			"issuer":             provider.Issuer,
			"allowed_client_ids": provider.AllowedClientIDs,
			"scopes":             provider.Scopes,
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

// pathOIDCDeleteProvider is used to delete an assignment
func (i *IdentityStore) pathOIDCDeleteProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
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
	scopes := append(p.Scopes, openIDScope)

	disc := providerDiscovery{
		AuthorizationEndpoint: strings.Replace(p.effectiveIssuer, "/v1/", "/ui/vault/", 1) + "/authorize",
		IDTokenAlgs:           supportedAlgs,
		Issuer:                p.effectiveIssuer,
		Keys:                  p.effectiveIssuer + "/.well-known/keys",
		ResponseTypes:         []string{"code"},
		Scopes:                scopes,
		Subjects:              []string{"public"},
		TokenEndpoint:         p.effectiveIssuer + "/token",
		UserinfoEndpoint:      p.effectiveIssuer + "/userinfo",
	}

	data, err := json.Marshal(disc)
	if err != nil {
		return nil, err
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:      200,
			logical.HTTPRawBody:         data,
			logical.HTTPContentType:     "application/json",
			logical.HTTPRawCacheControl: "max-age=3600",
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
	for name, _ := range keyNames {
		entry, err := s.Get(ctx, namedKeyConfigPath+name)
		if err != nil {
			return nil, err
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
	// Helper for preparing the non-standard OIDC error response:
	// https://openid.net/specs/openid-connect-core-1_0.html#AuthError
	errorResp := func(errorCode, description, state string) (*logical.Response, error) {
		response := map[string]interface{}{
			"error":             errorCode,
			"error_description": description,
			"state":             state,
		}

		statusCode := http.StatusBadRequest
		if errorCode == ErrAuthServerError {
			statusCode = http.StatusInternalServerError
		}

		data, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				logical.HTTPStatusCode:  statusCode,
				logical.HTTPRawBody:     data,
				logical.HTTPContentType: "application/json",
			},
		}, nil
	}

	// Validate the state
	state := d.Get("state").(string)
	if state == "" {
		return errorResp(ErrAuthInvalidRequest, "state parameter is required", "")
	}

	// Get the namespace
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return errorResp(ErrAuthServerError, err.Error(), state)
	}

	// Get the OIDC provider
	name := d.Get("name").(string)
	provider, err := i.getOIDCProvider(ctx, req.Storage, name)
	if err != nil {
		return errorResp(ErrAuthServerError, err.Error(), state)
	}
	if provider == nil {
		return errorResp(ErrAuthInvalidRequest, "provider not found", state)
	}

	// Validate that a scope parameter is present and contains the openid scope value
	scopes := strutil.ParseStringSlice(d.Get("scope").(string), " ")
	if len(scopes) == 0 || !strutil.StrListContains(scopes, openIDScope) {
		return errorResp(ErrAuthInvalidRequest,
			fmt.Sprintf("scope parameter must contain the %q value", openIDScope), state)
	}

	// Validate the response type
	responseType := d.Get("response_type").(string)
	if responseType == "" {
		return errorResp(ErrAuthInvalidRequest, "response_type parameter is required", state)
	}
	if responseType != "code" {
		return errorResp(ErrAuthUnsupportedResponseType, "unsupported response_type value", state)
	}

	// Validate the client ID
	clientID := d.Get("client_id").(string)
	if clientID == "" {
		return errorResp(ErrAuthInvalidClientID, "client_id parameter is required", state)
	}
	client, err := i.clientByID(ctx, req.Storage, clientID)
	if err != nil {
		return errorResp(ErrAuthServerError, err.Error(), state)
	}
	if client == nil {
		return errorResp(ErrAuthInvalidClientID, "client with client_id not found", state)
	}
	if !strutil.StrListContains(provider.AllowedClientIDs, "*") &&
		!strutil.StrListContains(provider.AllowedClientIDs, clientID) {
		return errorResp(ErrAuthUnauthorizedClient, "client is not authorized to use the provider", state)
	}

	// Validate the redirect URI
	redirectURI := d.Get("redirect_uri").(string)
	if redirectURI == "" {
		return errorResp(ErrAuthInvalidRequest, "redirect_uri parameter is required", state)
	}
	if !strutil.StrListContains(client.RedirectURIs, redirectURI) {
		return errorResp(ErrAuthInvalidRedirectURI, "redirect_uri is not allowed for the client", state)
	}

	// Validate the nonce
	nonce := d.Get("nonce").(string)
	if nonce == "" {
		return errorResp(ErrAuthInvalidRequest, "nonce parameter is required", state)
	}

	// Validate that there is an identity entity associated with the request
	entity, err := i.MemDBEntityByID(req.EntityID, false)
	if err != nil {
		return errorResp(ErrAuthServerError, err.Error(), state)
	}
	if entity == nil {
		return errorResp(ErrAuthAccessDenied, "identity entity must be associated with the request", state)
	}

	// Validate that the identity entity associated with the request
	// is a member of the client assignments' groups or entities
	isMember, err := i.entityHasAssignment(ctx, req.Storage, entity, client.Assignments)
	if err != nil {
		return errorResp(ErrAuthServerError, err.Error(), state)
	}
	if !isMember {
		return errorResp(ErrAuthAccessDenied, "identity entity not authorized by client assignment", state)
	}

	// Create the auth code cache entry
	authCodeEntry := &authCodeCacheEntry{
		entityID: entity.GetID(),
		nonce:    nonce,
		scopes:   scopes,
	}

	// Validate the optional max_age parameter to check if an active re-authentication
	// of the user should occur. Re-authentication will be requested if max_age=0 or the
	// last time the token actively authenticated exceeds the given max_age requirement.
	// Returning ErrAuthMaxAgeReAuthenticate will enforce the user to re-authenticate via
	// the user agent.
	if maxAgeRaw, ok := d.GetOk("max_age"); ok {
		maxAge := maxAgeRaw.(int)
		if maxAge < 0 {
			return errorResp(ErrAuthInvalidRequest, "max_age must be greater than zero", state)
		}
		if maxAge == 0 {
			// TODO: solve for the potential UI loop here or make max_age=0 invalid
			return errorResp(ErrAuthMaxAgeReAuthenticate, "active re-authentication is required by max_age", state)
		}

		// Look up the token associated with the request
		te, err := i.tokenStorer.LookupToken(ctx, req.ClientToken)
		if err != nil {
			return errorResp(ErrAuthServerError, err.Error(), state)
		}
		if te == nil {
			return errorResp(ErrAuthAccessDenied, "token associated with request not found", state)
		}

		// Check if the token creation time violates the max age requirement
		now := time.Now().UTC()
		lastAuthTime := time.Unix(te.CreationTime, 0).UTC()
		secondsSince := int(now.Sub(lastAuthTime).Seconds())
		if secondsSince > maxAge {
			return errorResp(ErrAuthMaxAgeReAuthenticate, "active re-authentication is required by max_age", state)
		}

		// Set the auth time to use for the auth_time claim in the token exchange
		authCodeEntry.authTime = lastAuthTime
	}

	// Generate the authorization code
	code, err := base62.Random(32)
	if err != nil {
		return errorResp(ErrAuthServerError, err.Error(), state)
	}

	// Cache the authorization code for a subsequent token exchange
	if err := i.oidcAuthCodeCache.SetDefault(ns, code, authCodeEntry); err != nil {
		return errorResp(ErrAuthServerError, err.Error(), state)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"code":  code,
			"state": state,
		},
	}, nil
}

// entityHasAssignment returns true if the entity is a member of any of the
// assignments' groups or entities. Otherwise, returns false or an error.
func (i *IdentityStore) entityHasAssignment(ctx context.Context, s logical.Storage, entity *identity.Entity, assignments []string) (bool, error) {
	for _, a := range assignments {
		assignment, err := i.getOIDCAssignment(ctx, s, a)
		if err != nil {
			return false, err
		}
		if assignment == nil {
			return false, fmt.Errorf("client assignment %q not found", a)
		}

		// Get the group names that the entity is a member of
		entityGroups, err := i.MemDBGroupsByMemberEntityID(entity.GetID(), true, false)
		if err != nil {
			return false, err
		}
		entityGroupNames := make(map[string]bool)
		for _, group := range entityGroups {
			entityGroupNames[group.Name] = true
		}

		// Check if the entity is a member of any groups in the assignment
		for _, group := range assignment.Groups {
			if entityGroupNames[group] {
				return true, nil
			}
		}

		// Check if the entity is a member of the assignment's entities
		if strutil.StrListContains(assignment.Entities, entity.GetName()) {
			return true, nil
		}
	}

	return false, nil
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
