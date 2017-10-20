package vault

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/ptypes"
	memdb "github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func groupPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "group$",
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the group.",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the group.",
				},
				"metadata": {
					Type:        framework.TypeStringSlice,
					Description: "Metadata to be associated with the group. Format should be a list of `key=value` pairs.",
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
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathGroupRegister,
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["register"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["register"][1]),
		},
		{
			Pattern: "group/id/" + framework.GenericNameRegex("id"),
			Fields: map[string]*framework.FieldSchema{
				"id": {
					Type:        framework.TypeString,
					Description: "ID of the group.",
				},
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the group.",
				},
				"metadata": {
					Type:        framework.TypeStringSlice,
					Description: "Metadata to be associated with the group. Format should be a list of `key=value` pairs.",
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
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathGroupIDUpdate,
				logical.ReadOperation:   i.pathGroupIDRead,
				logical.DeleteOperation: i.pathGroupIDDelete,
			},

			HelpSynopsis:    strings.TrimSpace(groupHelp["group-by-id"][0]),
			HelpDescription: strings.TrimSpace(groupHelp["group-by-id"][1]),
		},
		{
			Pattern: "group/id/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathGroupIDList,
			},

			HelpSynopsis:    strings.TrimSpace(entityHelp["group-id-list"][0]),
			HelpDescription: strings.TrimSpace(entityHelp["group-id-list"][1]),
		},
	}
}

func (i *IdentityStore) pathGroupRegister(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	_, ok := d.GetOk("id")
	if ok {
		return i.pathGroupIDUpdate(req, d)
	}

	i.groupLock.Lock()
	defer i.groupLock.Unlock()

	return i.handleGroupUpdateCommon(req, d, nil)
}

func (i *IdentityStore) pathGroupIDUpdate(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	groupID := d.Get("id").(string)
	if groupID == "" {
		return logical.ErrorResponse("empty group ID"), nil
	}

	i.groupLock.Lock()
	defer i.groupLock.Unlock()

	group, err := i.memDBGroupByID(groupID, true)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return logical.ErrorResponse("invalid group ID"), nil
	}

	return i.handleGroupUpdateCommon(req, d, group)
}

func (i *IdentityStore) handleGroupUpdateCommon(req *logical.Request, d *framework.FieldData, group *identity.Group) (*logical.Response, error) {
	var err error
	var newGroup bool
	if group == nil {
		group = &identity.Group{}
		newGroup = true
	}

	// Update the policies if supplied
	policiesRaw, ok := d.GetOk("policies")
	if ok {
		group.Policies = policiesRaw.([]string)
	}

	// Get the name
	groupName := d.Get("name").(string)
	if groupName != "" {
		// Check if there is a group already existing for the given name
		groupByName, err := i.memDBGroupByName(groupName, false)
		if err != nil {
			return nil, err
		}

		// If this is a new group and if there already exists a group by this
		// name, error out. If the name of an existing group is about to be
		// modified into something which is already tied to a different group,
		// error out.
		switch {
		case (newGroup && groupByName != nil), (groupByName != nil && group.ID != "" && groupByName.ID != group.ID):
			return logical.ErrorResponse("group name is already in use"), nil
		}
		group.Name = groupName
	}

	metadataRaw, ok := d.GetOk("metadata")
	if ok {
		group.Metadata, err = parseMetadata(metadataRaw.([]string))
		if err != nil {
			return logical.ErrorResponse(fmt.Sprintf("failed to parse group metadata: %v", err)), nil
		}
	}

	memberEntityIDsRaw, ok := d.GetOk("member_entity_ids")
	if ok {
		group.MemberEntityIDs = memberEntityIDsRaw.([]string)
		if len(group.MemberEntityIDs) > 512 {
			return logical.ErrorResponse("member entity IDs exceeding the limit of 512"), nil
		}
	}

	memberGroupIDsRaw, ok := d.GetOk("member_group_ids")
	var memberGroupIDs []string
	if ok {
		memberGroupIDs = memberGroupIDsRaw.([]string)
	}

	err = i.sanitizeAndUpsertGroup(group, memberGroupIDs)
	if err != nil {
		return nil, err
	}

	respData := map[string]interface{}{
		"id":   group.ID,
		"name": group.Name,
	}
	return &logical.Response{
		Data: respData,
	}, nil
}

func (i *IdentityStore) pathGroupIDRead(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	groupID := d.Get("id").(string)
	if groupID == "" {
		return logical.ErrorResponse("empty group id"), nil
	}

	group, err := i.memDBGroupByID(groupID, false)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	return i.handleGroupReadCommon(group)
}

func (i *IdentityStore) handleGroupReadCommon(group *identity.Group) (*logical.Response, error) {
	if group == nil {
		return nil, fmt.Errorf("nil group")
	}

	respData := map[string]interface{}{}
	respData["id"] = group.ID
	respData["name"] = group.Name
	respData["policies"] = group.Policies
	respData["member_entity_ids"] = group.MemberEntityIDs
	respData["metadata"] = group.Metadata
	respData["creation_time"] = ptypes.TimestampString(group.CreationTime)
	respData["last_update_time"] = ptypes.TimestampString(group.LastUpdateTime)
	respData["modify_index"] = group.ModifyIndex

	memberGroupIDs, err := i.memberGroupIDsByID(group.ID)
	if err != nil {
		return nil, err
	}
	respData["member_group_ids"] = memberGroupIDs

	return &logical.Response{
		Data: respData,
	}, nil
}

func (i *IdentityStore) pathGroupIDDelete(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	groupID := d.Get("id").(string)
	if groupID == "" {
		return logical.ErrorResponse("empty group ID"), nil
	}
	return nil, i.deleteGroupByID(groupID)
}

// pathGroupIDList lists the IDs of all the groups in the identity store
func (i *IdentityStore) pathGroupIDList(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ws := memdb.NewWatchSet()
	iter, err := i.memDBGroupIterator(ws)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch iterator for group in memdb: %v", err)
	}

	var groupIDs []string
	for {
		raw := iter.Next()
		if raw == nil {
			break
		}
		groupIDs = append(groupIDs, raw.(*identity.Group).ID)
	}

	return logical.ListResponse(groupIDs), nil
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
