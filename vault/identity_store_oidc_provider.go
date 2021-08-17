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

type assignment struct {
	Groups   []string `json:"groups"`
	Entities []string `json:"entities"`
}

type scope struct {
	Template    string `json:"template"`
	Description string `json:"description"`
}

const (
	oidcProviderPrefix = "oidc_provider/"
	assignmentPath     = oidcProviderPrefix + "assignment/"
	scopePath          = oidcProviderPrefix + "scope/"
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
	}
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
	err := req.Storage.Delete(ctx, assignmentPath+name)
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
	err := req.Storage.Delete(ctx, scopePath+name)
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
