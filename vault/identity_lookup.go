package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func lookupPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "lookup/group$",
			Fields: map[string]*framework.FieldSchema{
				"type": {
					Type:        framework.TypeString,
					Description: "Type of lookup. Current supported values are 'by_id' and 'by_name'",
				},
				"group_name": {
					Type:        framework.TypeString,
					Description: "Name of the group.",
				},
				"group_id": {
					Type:        framework.TypeString,
					Description: "ID of the group.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathLookupGroupUpdate,
			},

			HelpSynopsis:    strings.TrimSpace(lookupHelp["lookup-group"][0]),
			HelpDescription: strings.TrimSpace(lookupHelp["lookup-group"][1]),
		},
	}
}

func (i *IdentityStore) pathLookupGroupUpdate(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	lookupType := d.Get("type").(string)
	if lookupType == "" {
		return logical.ErrorResponse("empty type"), nil
	}

	switch lookupType {
	case "by_id":
		groupID := d.Get("group_id").(string)
		if groupID == "" {
			return logical.ErrorResponse("empty group_id"), nil
		}
		group, err := i.memDBGroupByID(groupID, false)
		if err != nil {
			return nil, err
		}
		return i.handleGroupReadCommon(group)
	case "by_name":
		groupName := d.Get("group_name").(string)
		if groupName == "" {
			return logical.ErrorResponse("empty group_name"), nil
		}
		group, err := i.memDBGroupByName(groupName, false)
		if err != nil {
			return nil, err
		}
		return i.handleGroupReadCommon(group)
	default:
		return logical.ErrorResponse(fmt.Sprintf("unrecognized type %q", lookupType)), nil
	}

	return nil, nil
}

var lookupHelp = map[string][2]string{
	"lookup-group": {
		"Query groups based on factors.",
		"Currently this supports querying groups by its name or ID.",
	},
}
