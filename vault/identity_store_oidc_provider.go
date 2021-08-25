package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/identitytpl"
	"github.com/hashicorp/vault/sdk/logical"
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
	RedirectURIs   []string `json:"redirect_uris"`
	Assignments    []string `json:"assignments"`
	Key            string   `json:"key"`
	IDTokenTTL     int      `json:"id_token_ttl"`
	AccessTokenTTL int      `json:"access_token_ttl"`

	// used for OIDC endpoints
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

const (
	oidcProviderPrefix = "oidc_provider/"
	assignmentPath     = oidcProviderPrefix + "assignment/"
	scopePath          = oidcProviderPrefix + "scope/"
	clientPath         = oidcProviderPrefix + "client/"
	providerPath       = oidcProviderPrefix + "provider/"
)

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
					Description: "Name of the assignment",
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
	for client, _ := range clients {
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
	for client, _ := range clients {
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

	entry, err := req.Storage.Get(ctx, assignmentPath+name)
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
	return &logical.Response{
		Data: map[string]interface{}{
			"groups":   assignment.Groups,
			"entities": assignment.Entities,
		},
	}, nil
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
	if name == "openid" {
		return logical.ErrorResponse("the \"openid\" scope name is reserved"), nil
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

	var client client
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

	entry, err := req.Storage.Get(ctx, clientPath+name)
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
	err := req.Storage.Delete(ctx, clientPath+name)
	if err != nil {
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
		provider.effectiveIssuer = i.core.redirectAddr
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
