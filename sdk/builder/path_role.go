package builder

import (
	"context"
	"errors"
	"fmt"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

var role any

// pathRole extends the Vault API with a `/role`
// endpoint for the backend. You can choose whether
// or not certain attributes should be displayed,
// required, and named. You can also define different
// path patterns to list all roles.
func (gb *GenericBackend[CC, C, R]) pathRole(inputRole *Role[R, C]) []*framework.Path {
	role = inputRole

	return []*framework.Path{
		{
			Pattern: "role/" + framework.GenericNameRegex("name"),
			Fields:  inputRole.Fields,
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: gb.pathRolesRead,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: gb.pathRolesWrite,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: gb.pathRolesWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: gb.pathRolesDelete,
				},
			},
			HelpSynopsis:    pathRoleHelpSynopsis,
			HelpDescription: pathRoleHelpDescription,
		},
		{
			Pattern: "role/?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: gb.pathRolesList,
				},
			},
			HelpSynopsis:    pathRoleListHelpSynopsis,
			HelpDescription: pathRoleListHelpDescription,
		},
	}
}

// pathRolesList makes a request to Vault storage to retrieve a list of roles for the backend
func (gb *GenericBackend[CC, C, R]) pathRolesList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

// pathRolesRead makes a request to Vault storage to read a role and return response data
func (gb *GenericBackend[CC, C, R]) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role, err := gb.getRole(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	responseData := structs.Map(role)

	return &logical.Response{
		Data: responseData,
	}, nil
}

// pathRolesWrite makes a request to Vault storage to update a role based on the attributes passed to the role configuration
func (gb *GenericBackend[CC, C, R]) pathRolesWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	name, ok := data.GetOk("name")
	if !ok {
		return logical.ErrorResponse("missing role name"), nil
	}

	roleEntry, err := gb.getRole(ctx, req.Storage, name.(string))
	if err != nil {
		return nil, err
	}

	createOperation := (req.Operation == logical.CreateOperation)
	if roleEntry == nil {
		if !createOperation {
			return nil, errors.New("role not found during update operation")
		}
		roleEntry = new(R)
	}

	writeData := structs.Map(roleEntry)

	for k := range writeData {
		if userInput, ok := data.GetOk(k); ok {
			writeData[k] = userInput
		} else if createOperation {
			return nil, fmt.Errorf("missing %s in configuration", k)
		}
	}

	if err = gb.validateRole(writeData); err != nil {
		return nil, err
	}

	if err = gb.setRole(ctx, req.Storage, name.(string), writeData); err != nil {
		return nil, err
	}

	return nil, nil
}

func (gb *GenericBackend[CC, C, R]) validateRole(writeData map[string]any) error {
	inputRole := role.(Role[R, C])
	result := new(R)
	err := mapstructure.Decode(writeData, result)
	if err != nil {
		return err
	}

	return inputRole.ValidateFunc(result)
}

// pathRolesDelete makes a request to Vault storage to delete a role
func (gb *GenericBackend[CC, C, R]) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "role/"+d.Get("name").(string))
	if err != nil {
		return nil, fmt.Errorf("error deleting hashiCups role: %w", err)
	}

	return nil, nil
}

// setRole adds the role to the Vault storage API
func (gb *GenericBackend[CC, C, R]) setRole(ctx context.Context, s logical.Storage, name string, roleEntry map[string]interface{}) error {
	entry, err := logical.StorageEntryJSON("role/"+name, roleEntry)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("failed to create storage entry for role")
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

// getRole gets the role from the Vault storage API
func (gb *GenericBackend[CC, C, R]) getRole(ctx context.Context, s logical.Storage, name string) (*R, error) {
	if name == "" {
		return nil, fmt.Errorf("missing role name")
	}

	entry, err := s.Get(ctx, "role/"+name)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	role := new(R)

	if err := entry.DecodeJSON(&role); err != nil {
		return nil, err
	}
	return role, nil
}

const (
	pathRoleHelpSynopsis    = `Manages the Vault role for generating HashiCups tokens.`
	pathRoleHelpDescription = `
This path allows you to read and write roles used to generate HashiCups tokens.
You can configure a role to manage a user's token by setting the username field.
`

	pathRoleListHelpSynopsis    = `List the existing roles in HashiCups backend`
	pathRoleListHelpDescription = `Roles will be listed by the role name.`
)
