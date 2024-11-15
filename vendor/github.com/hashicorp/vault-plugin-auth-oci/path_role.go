// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// Constants for role specific data
const (
	// Increasing this above this limit might require implementing
	// client-side paging in the filterGroupMembership API
	MaxOCIDsPerRole = 100
)

func pathRole(b *backend) *framework.Path {
	p := &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("role"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixOCI,
			OperationSuffix: "role",
		},

		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the role.",
			},
			"ocid_list": {
				Type:        framework.TypeCommaStringSlice,
				Description: `A comma separated list of Group or Dynamic Group OCIDs that are allowed to take this role.`,
			},
		},

		ExistenceCheck: b.pathRoleExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathRoleCreateUpdate,
			logical.UpdateOperation: b.pathRoleCreateUpdate,
			logical.ReadOperation:   b.pathRoleRead,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleSyn,
		HelpDescription: pathRoleDesc,
	}

	tokenutil.AddTokenFields(p.Fields)

	return p
}

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/?",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixOCI,
			OperationVerb:   "list",
			OperationSuffix: "roles",
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

// Establishes dichotomy of request operation between CreateOperation and UpdateOperation.
// Returning 'true' forces an UpdateOperation, CreateOperation otherwise.
func (b *backend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	entry, err := b.getOCIRole(ctx, req.Storage, data.Get("role").(string))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// setOciRole creates or updates a role in the storage.
func (b *backend) setOCIRole(ctx context.Context, s logical.Storage, roleName string,
	roleEntry *OCIRoleEntry) error {
	if roleName == "" {
		return fmt.Errorf("missing role name")
	}

	if roleEntry == nil {
		return fmt.Errorf("nil role entry")
	}

	entry, err := logical.StorageEntryJSON("role/"+roleName, roleEntry)
	if err != nil {
		return err
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

// getOCIRole returns the properties set on the given role.
// This method does NOT check to see if a role upgrade is required. It is
// the responsibility of the caller to check if a role upgrade is required and,
// if so, to upgrade the role
func (b *backend) getOCIRole(ctx context.Context, s logical.Storage, roleName string) (*OCIRoleEntry, error) {
	if roleName == "" {
		return nil, fmt.Errorf("missing role name")
	}

	entry, err := s.Get(ctx, "role/"+roleName)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result OCIRoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (b *backend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	return nil, req.Storage.Delete(ctx, "role/"+roleName)
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

func (b *backend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleEntry, err := b.getOCIRole(ctx, req.Storage, data.Get("role").(string))
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, nil
	}

	responseData := map[string]interface{}{
		"ocid_list": append([]string{}, roleEntry.OcidList...),
	}

	roleEntry.PopulateTokenData(responseData)

	return &logical.Response{
		Data: responseData,
	}, nil
}

// create a Role
func (b *backend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleName := data.Get("role").(string)

	roleEntry, err := b.getOCIRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if roleEntry == nil && req.Operation == logical.CreateOperation {
		roleEntry = &OCIRoleEntry{}
	} else if roleEntry == nil {
		return logical.ErrorResponse("The specified role does not exist"), nil
	}

	if ocidList, ok := data.GetOk("ocid_list"); ok {
		roleEntry.OcidList = ocidList.([]string)
		if len(roleEntry.OcidList) > MaxOCIDsPerRole {
			return logical.ErrorResponse("Number of OCIDs for this role exceeds the limit"), nil
		}
	}

	if err := roleEntry.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	var resp *logical.Response

	if err := b.setOCIRole(ctx, req.Storage, roleName, roleEntry); err != nil {
		return nil, err
	}

	return resp, nil
}

// Struct to hold the information associated with an OCI role
type OCIRoleEntry struct {
	tokenutil.TokenParams

	OcidList []string `json:"ocid_list"`
}

const pathRoleSyn = `
Create a role and associate policies to it.
`

const pathRoleDesc = `
Create a role and associate policies to it.
`

const pathListRolesHelpSyn = `
Lists all the roles that are registered with Vault.
`

const pathListRolesHelpDesc = `
Roles will be listed by their respective role names.
`
