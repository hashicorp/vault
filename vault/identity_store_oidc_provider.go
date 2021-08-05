package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/identitytpl"
	"github.com/hashicorp/vault/sdk/logical"
)

type namedAssignment struct {
	name     string
	Groups   []string `json:"groups"`
	Entities []string `json:"entities"`
}

type namedScope struct {
	name        string
	Template    string `json:"template"`
	Description string `json:"description"`
}

type namedClient struct {
	name           string
	RedirectURIs   []string `json:"redirect_uris"`
	Assignments    []string `json:"assignments"`
	Key            string   `json:"key"`
	IDTokenTTL     int      `json:"id_token_ttl"`
	AccessTokenTTL int      `json:"access_token_ttl"`
}

const (
	namedAssignmentPath = oidcTokensPrefix + "named_assignments/"
	namedScopePath      = oidcTokensPrefix + "named_scopes/"
	namedClientPath     = oidcTokensPrefix + "named_clients/"
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
			HelpDescription: "Create, Read, Update, and Delete OIDC named assignments.",
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
					Description: "The template string for the scope",
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
			HelpDescription: "Create, Read, Update, and Delete OIDC named scopes.",
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
					Description: "A reference to a named key resource. Cannot be modified after creation. If not provided, the default key will be used.",
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
			HelpDescription: "Create, Read, Update, and Delete OIDC named clients.",
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
	}
}

// pathOIDCCreateUpdateAssignment is used to create a new named assignment or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateAssignment(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	var assignment namedAssignment
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, namedAssignmentPath+name)
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
		assignment.Entities = d.Get("entities").([]string)
	}

	if groupsRaw, ok := d.GetOk("groups"); ok {
		assignment.Groups = groupsRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		assignment.Groups = d.Get("groups").([]string)
	}

	// store named assignment
	entry, err := logical.StorageEntryJSON(namedAssignmentPath+name, assignment)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathOIDCListAssignment is used to list named assignments
func (i *IdentityStore) pathOIDCListAssignment(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	assignments, err := req.Storage.List(ctx, namedAssignmentPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(assignments), nil
}

// pathOIDCReadAssignment is used to read an existing assignment
func (i *IdentityStore) pathOIDCReadAssignment(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, namedAssignmentPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var storedNamedAssignment namedAssignment
	if err := entry.DecodeJSON(&storedNamedAssignment); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"groups":   storedNamedAssignment.Groups,
			"entities": storedNamedAssignment.Entities,
		},
	}, nil
}

// pathOIDCDeleteAssignment is used to delete an assignment
func (i *IdentityStore) pathOIDCDeleteAssignment(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	err := req.Storage.Delete(ctx, namedAssignmentPath+name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *IdentityStore) pathOIDCAssignmentExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, namedAssignmentPath+name)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

// pathOIDCCreateUpdateScope is used to create a new named scope or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateScope(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	var scope namedScope
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, namedScopePath+name)
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
			// namespace?
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
	// store named scope
	entry, err := logical.StorageEntryJSON(namedScopePath+name, scope)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathOIDCListScope is used to list named scopes
func (i *IdentityStore) pathOIDCListScope(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	scopes, err := req.Storage.List(ctx, namedScopePath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(scopes), nil
}

// pathOIDCReadScope is used to read an existing scope
func (i *IdentityStore) pathOIDCReadScope(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, namedScopePath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var storedNamedScope namedScope
	if err := entry.DecodeJSON(&storedNamedScope); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"template":    storedNamedScope.Template,
			"description": storedNamedScope.Description,
		},
	}, nil
}

// pathOIDCDeleteScope is used to delete an scope
func (i *IdentityStore) pathOIDCDeleteScope(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	err := req.Storage.Delete(ctx, namedScopePath+name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *IdentityStore) pathOIDCScopeExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, namedScopePath+name)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

// pathOIDCCreateUpdateClient is used to create a new named client or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateClient(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	name := d.Get("name").(string)

	var client namedClient
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, namedClientPath+name)
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

	if keyRaw, ok := d.GetOk("key"); ok {
		client.Key = keyRaw.(string)
	} else if req.Operation == logical.CreateOperation {
		client.Key = d.Get("key").(string)
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

	if err := i.oidcCache.Flush(ns); err != nil {
		return nil, err
	}

	// store named client
	entry, err := logical.StorageEntryJSON(namedClientPath+name, client)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathOIDCListClient is used to list named clients
func (i *IdentityStore) pathOIDCListClient(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	clients, err := req.Storage.List(ctx, namedClientPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(clients), nil
}

// pathOIDCReadClient is used to read an existing client
func (i *IdentityStore) pathOIDCReadClient(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, namedClientPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var storedNamedClient namedClient
	if err := entry.DecodeJSON(&storedNamedClient); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"redirect_uris":    storedNamedClient.RedirectURIs,
			"assignments":      storedNamedClient.Assignments,
			"key":              storedNamedClient.Key,
			"id_token_ttl":     storedNamedClient.IDTokenTTL,
			"access_token_ttl": storedNamedClient.AccessTokenTTL,
		},
	}, nil
}

// pathOIDCDeleteClient is used to delete an client
func (i *IdentityStore) pathOIDCDeleteClient(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	err := req.Storage.Delete(ctx, namedClientPath+name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *IdentityStore) pathOIDCClientExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, namedClientPath+name)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}
