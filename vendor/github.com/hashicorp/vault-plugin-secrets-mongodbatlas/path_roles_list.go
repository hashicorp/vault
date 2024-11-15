// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodbatlas

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *Backend) pathRolesList() *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixMongoDBAtlas,
			OperationVerb:   "list",
			OperationSuffix: "roles",
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.operationListRoles,
		},

		HelpSynopsis:    pathRolesListHelpSyn,
		HelpDescription: pathRolesListHelpDesc,
	}
}

func (b *Backend) operationListRoles(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "roles/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

const pathRolesListHelpSyn = `List the existing roles in this backend`
const pathRolesListHelpDesc = `Roles will be listed by the role name.`
