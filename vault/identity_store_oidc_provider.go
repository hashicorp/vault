package vault

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type namedAssignment struct {
	name     string
	Groups   []string `json:"groups"`
	Entities []string `json:"entities"`
}

const (
	namedAssignmentPath = oidcTokensPrefix + "named_assignments/"
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
