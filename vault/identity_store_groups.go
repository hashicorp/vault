// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	groupTypeInternal = "internal"
	groupTypeExternal = "external"
)

func groupPathFields() map[string]*framework.FieldSchema {
	return map[string]*framework.FieldSchema{
		"id": {
			Type:        framework.TypeString,
			Description: "ID of the group. If set, updates the corresponding existing group.",
		},
		"type": {
			Type:        framework.TypeString,
			Description: "Type of the group, 'internal' or 'external'. Defaults to 'internal'",
		},
		"name": {
			Type:        framework.TypeString,
			Description: "Name of the group.",
		},
		"metadata": {
			Type: framework.TypeKVPairs,
			Description: `Metadata to be associated with the group.
In CLI, this parameter can be repeated multiple times, and it all gets merged together.
For example:
vault <command> <path> metadata=key1=value1 metadata=key2=value2
					`,
		},
		"policies": {
			Type:        framework.TypeCommaStringSlice,
			Description: "Policies to be tied to the group.",
		},
		"member_group_ids": {
			Type:        framework.TypeCommaStringSlice,
			Description: "Group IDs to be assigned as group members.",
		},
		"member_entity_ids": {
			Type:        framework.TypeCommaStringSlice,
			Description: "Entity IDs to be assigned as group members.",
		},
	}
}

func groupPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "group$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "group",
				OperationVerb:   "create",
			},

			Fields: groupPathFields(),
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    i.pathGroupRegister(),
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["register"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["register"][1]),
		},
		{
			Pattern: "group/id/" + framework.GenericNameRegex("id"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "group",
				OperationSuffix: "by-id",
			},

			Fields: groupPathFields(),

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathGroupIDUpdate(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "update",
					},
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathGroupIDRead(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathGroupIDDelete(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["group-by-id"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-by-id"][1]),
		},
		{
			Pattern: "group/id/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "group",
				OperationSuffix: "by-id",
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathGroupIDList(),
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["group-id-list"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-id-list"][1]),
		},
		{
			Pattern: "group/name/(?P<name>.+)",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "group",
				OperationSuffix: "by-name",
			},

			Fields: groupPathFields(),

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathGroupNameUpdate(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "update",
					},
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathGroupNameRead(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathGroupNameDelete(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["group-by-name"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-by-name"][1]),
		},
		{
			Pattern: "group/name/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "group",
				OperationSuffix: "by-name",
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathGroupNameList(),
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["group-name-list"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-name-list"][1]),
		},
	}
}

// pathGroupRegister is always called by the active primary node of the cluster.
func (i *IdentityStore) pathGroupRegister() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		_, ok := d.GetOk("id")
		if ok {
			return i.pathGroupIDUpdate()(ctx, req, d)
		}

		_, ok = d.GetOk("name")
		if ok {
			return i.pathGroupNameUpdate()(ctx, req, d)
		}

		i.groupLock.Lock()
		defer i.groupLock.Unlock()

		return i.handleGroupUpdateCommon(ctx, req, d, nil)
	}
}

// pathGroupIDUpdate is always called by the active primary node of the cluster.
func (i *IdentityStore) pathGroupIDUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupID := d.Get("id").(string)
		if groupID == "" {
			return logical.ErrorResponse("empty group ID"), nil
		}

		i.groupLock.Lock()
		defer i.groupLock.Unlock()

		group, err := i.MemDBGroupByID(groupID, true)
		if err != nil {
			return nil, err
		}
		if group == nil {
			return logical.ErrorResponse("invalid group ID"), nil
		}

		return i.handleGroupUpdateCommon(ctx, req, d, group)
	}
}

// pathGroupNameUpdate is always called by the active primary node of the cluster.
func (i *IdentityStore) pathGroupNameUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupName := d.Get("name").(string)
		if groupName == "" {
			return logical.ErrorResponse("empty group name"), nil
		}

		i.groupLock.Lock()
		defer i.groupLock.Unlock()

		group, err := i.MemDBGroupByName(ctx, groupName, true)
		if err != nil {
			return nil, err
		}
		return i.handleGroupUpdateCommon(ctx, req, d, group)
	}
}

// handleGroupUpdateCommon is always handled by the active primary node of the cluster.
func (i *IdentityStore) handleGroupUpdateCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, group *identity.Group) (*logical.Response, error) {
	var newGroup bool
	if group == nil {
		group = new(identity.Group)
		newGroup = true
	}

	// Update the policies if supplied
	policiesRaw, ok := d.GetOk("policies")
	if ok {
		group.Policies = strutil.RemoveDuplicatesStable(policiesRaw.([]string), true)
	}

	if strutil.StrListContains(group.Policies, "root") {
		return logical.ErrorResponse("policies cannot contain root"), nil
	}

	groupTypeRaw, ok := d.GetOk("type")
	if ok {
		groupType := groupTypeRaw.(string)
		if group.Type != "" && groupType != group.Type {
			return logical.ErrorResponse(fmt.Sprintf("group type cannot be changed")), nil
		}

		group.Type = groupType
	}

	// If group type is not set, default to internal type
	if group.Type == "" {
		group.Type = groupTypeInternal
	}

	if group.Type != groupTypeInternal && group.Type != groupTypeExternal {
		return logical.ErrorResponse(fmt.Sprintf("invalid group type %q", group.Type)), nil
	}

	// Get the name
	groupName := d.Get("name").(string)
	if groupName != "" {
		// Check if there is a group already existing for the given name
		groupByName, err := i.MemDBGroupByName(ctx, groupName, false)
		if err != nil {
			return nil, err
		}

		// If no existing group has this name, go ahead with the creation or rename.
		// If there is a group, it must match the group passed in; groupByName
		// should not be modified as it's in memdb.
		switch {
		case groupByName == nil:
			// Allowed
		case groupByName.ID != group.ID:
			return logical.ErrorResponse("group name is already in use"), nil
		}
		group.Name = groupName
	}

	metadata, ok, err := d.GetOkErr("metadata")
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to parse metadata: %v", err)), nil
	}
	if ok {
		group.Metadata = metadata.(map[string]string)
	}

	memberEntityIDsRaw, ok := d.GetOk("member_entity_ids")
	if ok {
		if group.Type == groupTypeExternal {
			return logical.ErrorResponse("member entities can't be set manually for external groups"), nil
		}
		group.MemberEntityIDs = memberEntityIDsRaw.([]string)
	}

	memberGroupIDsRaw, ok := d.GetOk("member_group_ids")
	var memberGroupIDs []string
	if ok {
		if group.Type == groupTypeExternal {
			return logical.ErrorResponse("member groups can't be set for external groups"), nil
		}
		memberGroupIDs = memberGroupIDsRaw.([]string)
	}

	err = i.sanitizeAndUpsertGroup(ctx, group, nil, memberGroupIDs)
	if err != nil {
		if errStr := err.Error(); strings.HasPrefix(errStr, errCycleDetectedPrefix) {
			return logical.ErrorResponse(errStr), nil
		}

		return nil, err
	}

	if !newGroup {
		return nil, nil
	}

	respData := map[string]interface{}{
		"id":   group.ID,
		"name": group.Name,
	}
	return &logical.Response{
		Data: respData,
	}, nil
}

func (i *IdentityStore) pathGroupIDRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupID := d.Get("id").(string)
		if groupID == "" {
			return logical.ErrorResponse("empty group id"), nil
		}

		group, err := i.MemDBGroupByID(groupID, false)
		if err != nil {
			return nil, err
		}
		if group == nil {
			return nil, nil
		}

		return i.handleGroupReadCommon(ctx, group)
	}
}

func (i *IdentityStore) pathGroupNameRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupName := d.Get("name").(string)
		if groupName == "" {
			return logical.ErrorResponse("empty group name"), nil
		}

		group, err := i.MemDBGroupByName(ctx, groupName, false)
		if err != nil {
			return nil, err
		}
		if group == nil {
			return nil, nil
		}

		return i.handleGroupReadCommon(ctx, group)
	}
}

func (i *IdentityStore) handleGroupReadCommon(ctx context.Context, group *identity.Group) (*logical.Response, error) {
	if group == nil {
		return nil, nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns.ID != group.NamespaceID {
		return logical.ErrorResponse("request namespace is not the same as the group namespace"), logical.ErrPermissionDenied
	}

	respData := map[string]interface{}{}
	respData["id"] = group.ID
	respData["name"] = group.Name
	respData["policies"] = group.Policies
	respData["member_entity_ids"] = group.MemberEntityIDs
	respData["parent_group_ids"] = group.ParentGroupIDs
	respData["metadata"] = group.Metadata
	respData["creation_time"] = ptypes.TimestampString(group.CreationTime)
	respData["last_update_time"] = ptypes.TimestampString(group.LastUpdateTime)
	respData["modify_index"] = group.ModifyIndex
	respData["type"] = group.Type
	respData["namespace_id"] = group.NamespaceID

	aliasMap := map[string]interface{}{}
	if group.Alias != nil {
		aliasMap["id"] = group.Alias.ID
		aliasMap["canonical_id"] = group.Alias.CanonicalID
		aliasMap["mount_accessor"] = group.Alias.MountAccessor
		aliasMap["metadata"] = group.Alias.Metadata
		aliasMap["name"] = group.Alias.Name
		aliasMap["merged_from_canonical_ids"] = group.Alias.MergedFromCanonicalIDs
		aliasMap["creation_time"] = ptypes.TimestampString(group.Alias.CreationTime)
		aliasMap["last_update_time"] = ptypes.TimestampString(group.Alias.LastUpdateTime)

		if mountValidationResp := i.router.ValidateMountByAccessor(group.Alias.MountAccessor); mountValidationResp != nil {
			aliasMap["mount_path"] = mountValidationResp.MountPath
			aliasMap["mount_type"] = mountValidationResp.MountType
		}
	}

	respData["alias"] = aliasMap

	var memberGroupIDs []string
	memberGroups, err := i.MemDBGroupsByParentGroupID(group.ID, false)
	if err != nil {
		return nil, err
	}
	for _, memberGroup := range memberGroups {
		memberGroupIDs = append(memberGroupIDs, memberGroup.ID)
	}

	respData["member_group_ids"] = memberGroupIDs

	return &logical.Response{
		Data: respData,
	}, nil
}

// pathGroupIDDelete is always called by the active primary node of the cluster.
func (i *IdentityStore) pathGroupIDDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupID := d.Get("id").(string)
		if groupID == "" {
			return logical.ErrorResponse("empty group ID"), nil
		}

		return i.handleGroupDeleteCommon(ctx, groupID, true)
	}
}

// pathGroupNameDelete is always called by the active primary node of the cluster.
func (i *IdentityStore) pathGroupNameDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupName := d.Get("name").(string)
		if groupName == "" {
			return logical.ErrorResponse("empty group name"), nil
		}

		return i.handleGroupDeleteCommon(ctx, groupName, false)
	}
}

// handleGroupDeleteCommon is always handled by the active primary node of the cluster.
func (i *IdentityStore) handleGroupDeleteCommon(ctx context.Context, key string, byID bool) (*logical.Response, error) {
	// Acquire the lock to modify the group storage entry
	i.groupLock.Lock()
	defer i.groupLock.Unlock()

	// Create a MemDB transaction to delete group
	txn := i.db.Txn(true)
	defer txn.Abort()

	var group *identity.Group
	var err error
	switch byID {
	case true:
		group, err = i.MemDBGroupByIDInTxn(txn, key, false)
		if err != nil {
			return nil, err
		}
	default:
		group, err = i.MemDBGroupByNameInTxn(ctx, txn, key, false)
		if err != nil {
			return nil, err
		}
	}
	if group == nil {
		return nil, nil
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if group.NamespaceID != ns.ID {
		return logical.ErrorResponse("request namespace is not the same as the group namespace"), logical.ErrPermissionDenied
	}

	// Delete group alias from memdb
	if group.Type == groupTypeExternal && group.Alias != nil {
		err = i.MemDBDeleteAliasByIDInTxn(txn, group.Alias.ID, true)
		if err != nil {
			return nil, err
		}
	}

	// Delete the group using the same transaction
	err = i.MemDBDeleteGroupByIDInTxn(txn, group.ID)
	if err != nil {
		return nil, err
	}

	// Delete the group from storage
	err = i.groupPacker.DeleteItem(ctx, group.ID)
	if err != nil {
		return nil, err
	}

	// Committing the transaction *after* successfully deleting group
	txn.Commit()

	return nil, nil
}

// pathGroupIDList lists the IDs of all the groups in the identity store
func (i *IdentityStore) pathGroupIDList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		return i.handleGroupListCommon(ctx, true)
	}
}

// pathGroupNameList lists the names of all the groups in the identity store
func (i *IdentityStore) pathGroupNameList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		return i.handleGroupListCommon(ctx, false)
	}
}

func (i *IdentityStore) handleGroupListCommon(ctx context.Context, byID bool) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	txn := i.db.Txn(false)

	iter, err := txn.Get(groupsTable, "namespace_id", ns.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup groups using namespace ID: %w", err)
	}

	var keys []string
	groupInfo := map[string]interface{}{}

	type mountInfo struct {
		MountType string
		MountPath string
	}
	mountAccessorMap := map[string]mountInfo{}

	for entry := iter.Next(); entry != nil; entry = iter.Next() {
		group := entry.(*identity.Group)

		if byID {
			keys = append(keys, group.ID)
		} else {
			keys = append(keys, group.Name)
		}

		groupInfoEntry := map[string]interface{}{
			"name":                group.Name,
			"num_member_entities": len(group.MemberEntityIDs),
			"num_parent_groups":   len(group.ParentGroupIDs),
		}
		if group.Alias != nil {
			entry := map[string]interface{}{
				"id":             group.Alias.ID,
				"name":           group.Alias.Name,
				"mount_accessor": group.Alias.MountAccessor,
			}

			mi, ok := mountAccessorMap[group.Alias.MountAccessor]
			if ok {
				entry["mount_type"] = mi.MountType
				entry["mount_path"] = mi.MountPath
			} else {
				mi = mountInfo{}
				if mountValidationResp := i.router.ValidateMountByAccessor(group.Alias.MountAccessor); mountValidationResp != nil {
					mi.MountType = mountValidationResp.MountType
					mi.MountPath = mountValidationResp.MountPath
					entry["mount_type"] = mi.MountType
					entry["mount_path"] = mi.MountPath
				}
				mountAccessorMap[group.Alias.MountAccessor] = mi
			}

			groupInfoEntry["alias"] = entry
		}
		groupInfo[group.ID] = groupInfoEntry
	}

	return logical.ListResponseWithInfo(keys, groupInfo), nil
}

var groupHelp = map[string][2]string{
	"register": {
		"Create a new group.",
		"",
	},
	"group-by-id": {
		"Update or delete an existing group using its ID.",
		"",
	},
	"group-id-list": {
		"List all the group IDs.",
		"",
	},
}
