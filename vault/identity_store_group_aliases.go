// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func groupAliasPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "group-alias$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "group",
				OperationVerb:   "create",
				OperationSuffix: "alias",
			},

			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the group alias.",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Alias of the group.",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "Mount accessor to which this alias belongs to.",
				},
				"canonical_id": {
					Type:        framework.TypeString,
					Description: "ID of the group to which this is an alias.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback:                    i.pathGroupAliasRegister(),
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},

			HelpSynopsis:    strings.TrimSpace(groupAliasHelp["group-alias"][0]),
			HelpDescription: strings.TrimSpace(groupAliasHelp["group-alias"][1]),
		},
		{
			Pattern: "group-alias/id/" + framework.GenericNameRegex("id"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "group",
				OperationSuffix: "alias-by-id",
			},

			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the group alias.",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Alias of the group.",
				},
				"mount_accessor": {
					Type:        framework.TypeString,
					Description: "Mount accessor to which this alias belongs to.",
				},
				"canonical_id": {
					Type:        framework.TypeString,
					Description: "ID of the group to which this is an alias.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathGroupAliasIDUpdate(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "update",
					},
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathGroupAliasIDRead(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "read",
					},
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: i.pathGroupAliasIDDelete(),
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "delete",
					},
					ForwardPerformanceStandby:   true,
					ForwardPerformanceSecondary: true,
				},
			},

			HelpSynopsis:    strings.TrimSpace(groupAliasHelp["group-alias-by-id"][0]),
			HelpDescription: strings.TrimSpace(groupAliasHelp["group-alias-by-id"][1]),
		},
		{
			Pattern: "group-alias/id/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "group",
				OperationSuffix: "aliases-by-id",
			},

			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathGroupAliasIDList(),
			},

			HelpSynopsis:    strings.TrimSpace(groupAliasHelp["group-alias-id-list"][0]),
			HelpDescription: strings.TrimSpace(groupAliasHelp["group-alias-id-list"][1]),
		},
	}
}

func (i *IdentityStore) pathGroupAliasRegister() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		_, ok := d.GetOk("id")
		if ok {
			return i.pathGroupAliasIDUpdate()(ctx, req, d)
		}

		i.groupLock.Lock()
		defer i.groupLock.Unlock()

		return i.handleGroupAliasUpdateCommon(ctx, req, d, nil)
	}
}

func (i *IdentityStore) pathGroupAliasIDUpdate() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupAliasID := d.Get("id").(string)
		if groupAliasID == "" {
			return logical.ErrorResponse("empty group alias ID"), nil
		}

		i.groupLock.Lock()
		defer i.groupLock.Unlock()

		groupAlias, err := i.MemDBAliasByID(groupAliasID, true, true)
		if err != nil {
			return nil, err
		}
		if groupAlias == nil {
			return logical.ErrorResponse("invalid group alias ID"), nil
		}

		return i.handleGroupAliasUpdateCommon(ctx, req, d, groupAlias)
	}
}

// NOTE: Currently we don't allow by-factors modification of group aliases the
// way we do with entities. As a result if a groupAlias is defined here we know
// that this is an update, where they provided an ID parameter.
func (i *IdentityStore) handleGroupAliasUpdateCommon(ctx context.Context, req *logical.Request, d *framework.FieldData, groupAlias *identity.Alias) (*logical.Response, error) {
	var newGroup, previousGroup *identity.Group

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	if groupAlias == nil {
		groupAlias = &identity.Alias{
			CreationTime: timestamppb.Now(),
			NamespaceID:  ns.ID,
		}
		groupAlias.LastUpdateTime = groupAlias.CreationTime
	} else {
		if ns.ID != groupAlias.NamespaceID {
			return logical.ErrorResponse("existing alias not in the same namespace as request"), logical.ErrPermissionDenied
		}
		groupAlias.LastUpdateTime = timestamppb.Now()
		if groupAlias.CreationTime == nil {
			groupAlias.CreationTime = groupAlias.LastUpdateTime
		}
	}

	// Get group alias name
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing alias name"), nil
	}

	mountAccessor := d.Get("mount_accessor").(string)
	if mountAccessor == "" {
		return logical.ErrorResponse("missing mount_accessor"), nil
	}

	canonicalID := d.Get("canonical_id").(string)

	if groupAlias.Name == name && groupAlias.MountAccessor == mountAccessor && (canonicalID == "" || groupAlias.CanonicalID == canonicalID) {
		// Nothing to do, be idempotent
		return nil, nil
	}

	// Explicitly correct for previous versions that persisted this
	groupAlias.MountType = ""

	// Canonical ID handling
	{
		if canonicalID != "" {
			newGroup, err = i.MemDBGroupByID(canonicalID, true)
			if err != nil {
				return nil, err
			}
			if newGroup == nil {
				return logical.ErrorResponse("invalid group ID given in 'canonical_id'"), nil
			}
			if newGroup.Type != groupTypeExternal {
				return logical.ErrorResponse("alias can't be set on an internal group"), nil
			}
			if newGroup.NamespaceID != groupAlias.NamespaceID {
				return logical.ErrorResponse("group referenced with 'canonical_id' not in the same namespace as alias"), logical.ErrPermissionDenied
			}
			groupAlias.CanonicalID = canonicalID
		}
	}

	// Validate name/accessor whether new or update
	{
		mountEntry := i.router.MatchingMountByAccessor(mountAccessor)
		if mountEntry == nil {
			return logical.ErrorResponse(fmt.Sprintf("invalid mount accessor %q", mountAccessor)), nil
		}
		if mountEntry.Local {
			return logical.ErrorResponse(fmt.Sprintf("mount accessor %q is a local mount", mountAccessor)), nil
		}
		if mountEntry.NamespaceID != groupAlias.NamespaceID {
			return logical.ErrorResponse("mount referenced via 'mount_accessor' not in the same namespace as alias"), logical.ErrPermissionDenied
		}

		groupAliasByFactors, err := i.MemDBAliasByFactors(mountEntry.Accessor, name, false, true)
		if err != nil {
			return nil, err
		}
		// This check will still work for the new case too since it won't have
		// an ID yet
		if groupAliasByFactors != nil && groupAliasByFactors.ID != groupAlias.ID {
			return logical.ErrorResponse("combination of mount and group alias name is already in use"), nil
		}

		groupAlias.Name = name
		groupAlias.MountAccessor = mountAccessor
	}

	switch groupAlias.ID {
	case "":
		// It's a new alias
		if newGroup == nil {
			// If this is a new alias being tied to a non-existent group,
			// create a new group for it
			newGroup = &identity.Group{
				Type: groupTypeExternal,
			}
		}

	default:
		// Fetch the group, if any, to which the alias is tied to
		previousGroup, err = i.MemDBGroupByAliasID(groupAlias.ID, true)
		if err != nil {
			return nil, err
		}
		if previousGroup == nil {
			return nil, fmt.Errorf("group alias is not associated with a group")
		}
		if previousGroup.NamespaceID != groupAlias.NamespaceID {
			return logical.ErrorResponse("previous group found for alias not in the same namespace as alias"), logical.ErrPermissionDenied
		}

		if newGroup == nil || newGroup.ID == previousGroup.ID {
			// If newGroup is nil they didn't specify a canonical ID, so they
			// aren't trying to update it; set the existing group as the "new"
			// one. If it's the same ID they specified the same canonical ID,
			// so follow the same behavior.
			newGroup = previousGroup
			previousGroup = nil
		} else {
			// The alias is moving, so nil out the previous group alias
			previousGroup.Alias = nil
		}
	}

	newGroup.Alias = groupAlias
	err = i.sanitizeAndUpsertGroup(ctx, newGroup, previousGroup, nil)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"id":           groupAlias.ID,
			"canonical_id": newGroup.ID,
		},
	}, nil
}

// pathGroupAliasIDRead returns the properties of an alias for a given
// alias ID
func (i *IdentityStore) pathGroupAliasIDRead() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupAliasID := d.Get("id").(string)
		if groupAliasID == "" {
			return logical.ErrorResponse("empty group alias id"), nil
		}

		groupAlias, err := i.MemDBAliasByID(groupAliasID, false, true)
		if err != nil {
			return nil, err
		}

		return i.handleAliasReadCommon(ctx, groupAlias)
	}
}

// pathGroupAliasIDDelete deletes the group's alias for a given group alias ID
func (i *IdentityStore) pathGroupAliasIDDelete() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		groupAliasID := d.Get("id").(string)
		if groupAliasID == "" {
			return logical.ErrorResponse("missing group alias ID"), nil
		}

		i.groupLock.Lock()
		defer i.groupLock.Unlock()

		txn := i.db.Txn(true)
		defer txn.Abort()

		alias, err := i.MemDBAliasByIDInTxn(txn, groupAliasID, false, true)
		if err != nil {
			return nil, err
		}

		if alias == nil {
			return nil, nil
		}

		ns, err := namespace.FromContext(ctx)
		if err != nil {
			return nil, err
		}
		if ns.ID != alias.NamespaceID {
			return logical.ErrorResponse("request namespace is not the same as the group alias namespace"), logical.ErrPermissionDenied
		}

		group, err := i.MemDBGroupByAliasIDInTxn(txn, alias.ID, true)
		if err != nil {
			return nil, err
		}

		// If there is no group tied to a valid alias, something is wrong
		if group == nil {
			return nil, fmt.Errorf("alias not associated to a group")
		}

		// Delete group alias in memdb
		err = i.MemDBDeleteAliasByIDInTxn(txn, group.Alias.ID, true)
		if err != nil {
			return nil, err
		}

		// Delete the alias
		group.Alias = nil

		err = i.UpsertGroupInTxn(ctx, txn, group, true)
		if err != nil {
			return nil, err
		}

		txn.Commit()

		return nil, nil
	}
}

// pathGroupAliasIDList lists the IDs of all the valid group aliases in the
// identity store
func (i *IdentityStore) pathGroupAliasIDList() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		return i.handleAliasListCommon(ctx, true)
	}
}

var groupAliasHelp = map[string][2]string{
	"group-alias": {
		"Creates a new group alias, or updates an existing one.",
		"",
	},
	"group-alias-id": {
		"Update, read or delete a group alias using ID.",
		"",
	},
	"group-alias-id-list": {
		"List all the group alias IDs.",
		"",
	},
}
