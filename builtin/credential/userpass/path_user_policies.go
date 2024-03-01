// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package userpass

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathUserPolicies(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "users/" + framework.GenericNameRegex("username") + "/policies$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixUserpass,
			OperationVerb:   "update",
			OperationSuffix: "policies",
		},

		Fields: map[string]*framework.FieldSchema{
			"username": {
				Type:        framework.TypeString,
				Description: "Username for this user.",
			},
			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_policies"),
				Deprecated:  true,
			},
			"token_policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: "Comma-separated list of policies",
				DisplayAttrs: &framework.DisplayAttributes{
					Description: "A list of policies that will apply to the generated token for this user.",
				},
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathUserPoliciesUpdate,
		},

		HelpSynopsis:    pathUserPoliciesHelpSyn,
		HelpDescription: pathUserPoliciesHelpDesc,
	}
}

func (b *backend) pathUserPoliciesUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)

	userEntry, err := b.user(ctx, req.Storage, username)
	if err != nil {
		return nil, err
	}
	if userEntry == nil {
		return nil, fmt.Errorf("username does not exist")
	}

	policiesRaw, ok := d.GetOk("token_policies")
	if !ok {
		policiesRaw, ok = d.GetOk("policies")
		if ok {
			userEntry.Policies = policyutil.ParsePolicies(policiesRaw)
			userEntry.TokenPolicies = userEntry.Policies
		}
	} else {
		userEntry.TokenPolicies = policyutil.ParsePolicies(policiesRaw)
		_, ok = d.GetOk("policies")
		if ok {
			userEntry.Policies = userEntry.TokenPolicies
		} else {
			userEntry.Policies = nil
		}
	}

	return nil, b.setUser(ctx, req.Storage, username, userEntry)
}

const pathUserPoliciesHelpSyn = `
Update the policies associated with the username.
`

const pathUserPoliciesHelpDesc = `
This endpoint allows updating the policies associated with the username.
`
