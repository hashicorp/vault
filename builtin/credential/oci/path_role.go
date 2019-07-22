// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/framework"
	"regexp"
	"strings"
)

var (
	currentRoleStorageVersion = 1
)

const DEFAULT_TTL_SECONDS int = 1800
const MAX_ROLE_NAME_LENGTH int = 50
const MAX_OCIDS_PER_ROLE = 100; //increasing this above this limit might require implementing client-side paging in the filterGroupMembership API

func pathRole(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: "Name of the role.",
			},
			"description": {
				Type:        framework.TypeString,
				Description: `A description of the role.`,
			},
			"add_ocid_list": {
				Type:        framework.TypeCommaStringSlice,
				Description: `A comma separated list of OCIDs to add.`,
			},
			"remove_ocid_list": {
				Type:        framework.TypeCommaStringSlice,
				Description: `A comma separated list of OCIDs to remove. Is applicable only for the UPDATE operation.`,
			},
			"ttl": {
				Type:        framework.TypeInt,
				Default:     0,
				Description: `Duration in seconds after which the issued token should expire. Defaults to 1800 seconds.`,
			},
			"add_policy_list": {
				Type:        framework.TypeCommaStringSlice,
				Default:     "default",
				Description: "A list of Policies to be set on tokens issued using this role.",
			},
			"remove_policy_list": {
				Type:        framework.TypeCommaStringSlice,
				Default:     "default",
				Description: "A list of Policies to be remove from the role. Is applicable only for the UPDATE operation.",
			},
		},

		ExistenceCheck: b.pathRoleExistenceCheck,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathRoleCreate,
			logical.UpdateOperation: b.pathRoleUpdate,
			logical.ReadOperation:   b.pathRoleRead,
			logical.DeleteOperation: b.pathRoleDelete,
		},

		HelpSynopsis:    pathRoleSyn,
		HelpDescription: pathRoleDesc,
	}
}

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "role/?",

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
	entry, err := b.lockedOCIRole(ctx, req.Storage, strings.ToLower(data.Get("role").(string)))
	if err != nil {
		return false, err
	}
	return entry != nil, nil
}

// lockedOCIRole returns the properties set on the given role. This method
// acquires the read lock before reading the role from the storage.
func (b *backend) lockedOCIRole(ctx context.Context, s logical.Storage, roleName string) (*OCIRoleEntry, error) {
	if roleName == "" {
		return nil, fmt.Errorf("missing role name")
	}

	b.roleMutex.RLock()
	roleEntry, err := b.nonLockedOCIRole(ctx, s, roleName)
	// we manually unlock rather than defer the unlock because we might need to grab
	// a read/write lock in the upgrade path
	b.roleMutex.RUnlock()
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, nil
	}
	return roleEntry, nil
}

// lockedSetOCIRole creates or updates a role in the storage. This method
// acquires the write lock before creating or updating the role at the storage.
func (b *backend) lockedSetOCIRole(ctx context.Context, s logical.Storage, roleName string, roleEntry *OCIRoleEntry) error {
	if roleName == "" {
		return fmt.Errorf("missing role name")
	}

	if roleEntry == nil {
		return fmt.Errorf("nil role entry")
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	return b.nonLockedSetOCIRole(ctx, s, roleName, roleEntry)
}

// nonLockedSetOCIRole creates or updates a role in the storage. This method
// does not acquire the write lock before reading the role from the storage. If
// locking is desired, use lockedSetOCIRole instead.
func (b *backend) nonLockedSetOCIRole(ctx context.Context, s logical.Storage, roleName string,
	roleEntry *OCIRoleEntry) error {
	if roleName == "" {
		return fmt.Errorf("missing role name")
	}

	if roleEntry == nil {
		return fmt.Errorf("nil role entry")
	}

	entry, err := logical.StorageEntryJSON("role/"+strings.ToLower(roleName), roleEntry)
	if err != nil {
		return err
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

// nonLockedOCIRole returns the properties set on the given role. This method
// does not acquire the read lock before reading the role from the storage. If
// locking is desired, use lockedOCIRole instead.
// This method also does NOT check to see if a role upgrade is required. It is
// the responsibility of the caller to check if a role upgrade is required and,
// if so, to upgrade the role
func (b *backend) nonLockedOCIRole(ctx context.Context, s logical.Storage, roleName string) (*OCIRoleEntry, error) {
	if roleName == "" {
		return nil, fmt.Errorf("missing role name")
	}

	entry, err := s.Get(ctx, "role/"+strings.ToLower(roleName))
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
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	return nil, req.Storage.Delete(ctx, "role/"+strings.ToLower(roleName))
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b.roleMutex.RLock()
	defer b.roleMutex.RUnlock()

	roles, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

func (b *backend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleEntry, err := b.lockedOCIRole(ctx, req.Storage, strings.ToLower(data.Get("role").(string)))
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: roleEntry.ToResponseData(),
	}, nil
}

// create a Role
func (b *backend) pathRoleCreate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleName := strings.ToLower(data.Get("role").(string))
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	if len(roleName) > MAX_ROLE_NAME_LENGTH {
		return logical.ErrorResponse("role length exceeds the limit"), nil
	}

	validateRoleRegEx := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !validateRoleRegEx(roleName) {
		return logical.ErrorResponse("role is invalid"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	roleEntry, err := b.nonLockedOCIRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if roleEntry != nil {
		return logical.ErrorResponse("The specified role already exists"), nil
	}

	roleEntry = &OCIRoleEntry{
		Role:    roleName,
		Version: 1,
	}

	if description, ok := data.GetOk("description"); ok {
		roleEntry.Description = description.(string)
	}

	if add_ocid_list, ok := data.GetOk("add_ocid_list"); ok {
		roleEntry.OcidList = add_ocid_list.([]string)
		if(len(roleEntry.OcidList) > MAX_OCIDS_PER_ROLE) {
			return logical.ErrorResponse("Number of OCIDs for this role exceeds the limit"), nil
		}
	}

	var resp logical.Response

	ttl, ok := data.GetOk("ttl")
	if ok {
		ttlVal := ttl.(int)

		if ttlVal > DEFAULT_TTL_SECONDS {
			return logical.ErrorResponse(fmt.Sprintf("Given ttl of %d seconds should be lesser than %d seconds;", ttlVal, DEFAULT_TTL_SECONDS)), nil
		}

		if ttlVal < 0 {
			return logical.ErrorResponse("ttl cannot be negative"), nil
		}

		roleEntry.TTL = ttlVal
	} else {
		roleEntry.TTL = DEFAULT_TTL_SECONDS
	}

	if add_policy_list, ok := data.GetOk("add_policy_list"); ok {
		roleEntry.PolicyList = add_policy_list.([]string)
	}

	if err := b.nonLockedSetOCIRole(ctx, req.Storage, roleName, roleEntry); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Update if it already exists
func (b *backend) pathRoleUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	roleName := strings.ToLower(data.Get("role").(string))
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	b.roleMutex.Lock()
	defer b.roleMutex.Unlock()

	roleEntry, err := b.nonLockedOCIRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if roleEntry == nil {
		return logical.ErrorResponse("The specified role does not exist"), nil
	} else {
		roleEntry.Version = roleEntry.Version + 1
	}

	if description, ok := data.GetOk("description"); ok {
		roleEntry.Description = description.(string)
	}

	//Add and Remove the OCIDs
	ocidMap := sliceToMap(roleEntry.OcidList)

	if add_ocid_list, ok := data.GetOk("add_ocid_list"); ok {
		addOcidSlice := add_ocid_list.([]string)
		ocidMap = addSliceToMap(addOcidSlice, ocidMap)
	}

	if remove_ocid_list, ok := data.GetOk("remove_ocid_list"); ok {
		removeOcidSlice := remove_ocid_list.([]string)
		ocidMap = removeSliceFromMap(removeOcidSlice, ocidMap)
	}

	roleEntry.OcidList = mapToSlice(ocidMap)

	if(len(roleEntry.OcidList) > MAX_OCIDS_PER_ROLE) {
		return logical.ErrorResponse("Number of OCIDs for this role exceeds the limit"), nil
	}

	var resp logical.Response

	ttl, ok := data.GetOk("ttl")
	if ok {
		ttlVal := ttl.(int)

		if ttlVal > DEFAULT_TTL_SECONDS {
			return logical.ErrorResponse(fmt.Sprintf("Given ttl of %d seconds should be lesser than %d seconds;", ttlVal, DEFAULT_TTL_SECONDS)), nil
		}

		if ttlVal < 0 {
			return logical.ErrorResponse("ttl cannot be negative"), nil
		}

		roleEntry.TTL = ttlVal
	}

	//Add and Remove the Policies
	policyMap := sliceToMap(roleEntry.PolicyList)

	if add_policy_list, ok := data.GetOk("add_policy_list"); ok {
		addPolicySlice := add_policy_list.([]string)
		policyMap = addSliceToMap(addPolicySlice, policyMap)
	}

	if remove_policy_list, ok := data.GetOk("remove_policy_list"); ok {
		removePolicySlice := remove_policy_list.([]string)
		policyMap = removeSliceFromMap(removePolicySlice, policyMap)
	}

	roleEntry.PolicyList = mapToSlice(policyMap)

	if err := b.nonLockedSetOCIRole(ctx, req.Storage, roleName, roleEntry); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Struct to hold the information associated with an OCI role
type OCIRoleEntry struct {
	Role        string   `json:"role" `
	Description string   `json:"description" `
	OcidList    []string `json:"ocid_list"`
	TTL         int      `json:"ttl"`
	PolicyList  []string `json:"policy_list"`
	Version     int      `json:"version"`
}

func (r *OCIRoleEntry) ToResponseData() map[string]interface{} {
	responseData := map[string]interface{}{
		"role":        r.Role,
		"description": r.Description,
		"ocid_list":   r.OcidList,
		"ttl":         r.TTL,
		"policy_list": r.PolicyList,
		"version":     r.Version,
	}

	convertNilToEmptySlice := func(data map[string]interface{}, field string) {
		if data[field] == nil || len(data[field].([]string)) == 0 {
			data[field] = []string{}
		}
	}

	convertNilToEmptySlice(responseData, "ocid_list")
	convertNilToEmptySlice(responseData, "policy_list")

	return responseData
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
