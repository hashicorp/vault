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

type namedProvider struct {
	name             string
	Issuer           string   `json:"issuer"`
	AllowedClientIDs []string `json:"allowed_client_ids"`
	Scopes           []string `json:"scopes"`
}

const (
	namedAssignmentPath = oidcTokensPrefix + "named_assignments/"
	namedProviderPath   = oidcTokensPrefix + "named_providers/"
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

// pathOIDCCreateUpdateProvider is used to create a new named provider or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	var provider namedProvider
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, namedProviderPath+name)
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
		provider.AllowedClientIDs = d.Get("allowed_client_ids").([]string)
	}

	if scopesRaw, ok := d.GetOk("scopes"); ok {
		provider.Scopes = scopesRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		provider.Scopes = d.Get("scopes").([]string)
	}

	// store named provider
	entry, err := logical.StorageEntryJSON(namedProviderPath+name, provider)
	if err != nil {
		return nil, err
	}

	return nil, req.Storage.Put(ctx, entry)
}

// pathOIDCListProvider is used to list named providers
func (i *IdentityStore) pathOIDCListProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	providers, err := req.Storage.List(ctx, namedProviderPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(providers), nil
}

// pathOIDCReadProvider is used to read an existing provider
func (i *IdentityStore) pathOIDCReadProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, namedProviderPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var storedNameProvider namedProvider
	if err := entry.DecodeJSON(&storedNameProvider); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"issuer":             storedNameProvider.Issuer,
			"allowed_client_ids": storedNameProvider.AllowedClientIDs,
			"scopes":             storedNameProvider.Scopes,
		},
	}, nil
}

// pathOIDCDeleteProvider is used to delete an assignment
func (i *IdentityStore) pathOIDCDeleteProvider(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	return nil, req.Storage.Delete(ctx, namedProviderPath+name)
}

func (i *IdentityStore) pathOIDCProviderExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, namedProviderPath+name)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}
