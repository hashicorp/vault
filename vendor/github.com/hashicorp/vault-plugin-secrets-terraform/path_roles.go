// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfc

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// terraformRoleEntry is a Vault role construct that maps to TFC/TFE
type terraformRoleEntry struct {
	Name         string        `json:"name"`
	Organization string        `json:"organization,omitempty"`
	TeamID       string        `json:"team_id,omitempty"`
	UserID       string        `json:"user_id,omitempty"`
	TTL          time.Duration `json:"ttl"`
	MaxTTL       time.Duration `json:"max_ttl"`
	Token        string        `json:"token,omitempty"`
	TokenID      string        `json:"token_id,omitempty"`
}

func (r *terraformRoleEntry) toResponseData() map[string]interface{} {
	respData := map[string]interface{}{
		"name":    r.Name,
		"ttl":     r.TTL.Seconds(),
		"max_ttl": r.MaxTTL.Seconds(),
	}
	if r.Organization != "" {
		respData["organization"] = r.Organization
	}
	if r.TeamID != "" {
		respData["team_id"] = r.TeamID
	}
	if r.UserID != "" {
		respData["user_id"] = r.UserID
	}
	return respData
}

func pathRole(b *tfBackend) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "role/" + framework.GenericNameRegex("name"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixTerraformCloud,
				OperationSuffix: "role",
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Name of the role",
					Required:    true,
				},
				"organization": {
					Type:        framework.TypeString,
					Description: "Name of the Terraform Cloud or Enterprise organization",
				},
				"team_id": {
					Type:        framework.TypeString,
					Description: "ID of the Terraform Cloud or Enterprise team under organization (e.g., settings/teams/team-xxxxxxxxxxxxx)",
				},
				"user_id": {
					Type:        framework.TypeString,
					Description: "ID of the Terraform Cloud or Enterprise user (e.g., user-xxxxxxxxxxxxxxxx)",
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Default lease for generated credentials. If not set or set to 0, will use system default.",
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Maximum time for role. If not set or set to 0, will use system default.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathRolesRead,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathRolesWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathRolesDelete,
				},
			},
			HelpSynopsis:    pathRoleHelpSynopsis,
			HelpDescription: pathRoleHelpDescription,
		},
		{
			Pattern: "role/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: operationPrefixTerraformCloud,
				OperationVerb:   "list",
				OperationSuffix: "roles",
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathRolesList,
				},
			},

			HelpSynopsis:    pathRoleListHelpSynopsis,
			HelpDescription: pathRoleListHelpDescription,
		},
	}
}

func (b *tfBackend) pathRolesList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

func (b *tfBackend) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entry, err := b.getRole(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: entry.toResponseData(),
	}, nil
}

func (b *tfBackend) pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	roleEntry, err := b.getRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}

	if roleEntry == nil {
		roleEntry = &terraformRoleEntry{}
	}

	roleEntry.Name = name
	if organization, ok := d.GetOk("organization"); ok {
		roleEntry.Organization = organization.(string)
	}

	if teamID, ok := d.GetOk("team_id"); ok {
		roleEntry.TeamID = teamID.(string)
	}

	if userID, ok := d.GetOk("user_id"); ok {
		roleEntry.UserID = userID.(string)
	}

	if roleEntry.UserID != "" && (roleEntry.Organization != "" || roleEntry.TeamID != "") {
		return logical.ErrorResponse("cannot provide a user_id in combination with organization or team_id"), nil
	}

	if roleEntry.UserID == "" && roleEntry.Organization == "" && roleEntry.TeamID == "" {
		return logical.ErrorResponse("must provide an organization name, team id, or user id"), nil
	}

	if ttlRaw, ok := d.GetOk("ttl"); ok {
		roleEntry.TTL = time.Duration(ttlRaw.(int)) * time.Second
	}

	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		roleEntry.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	}

	if roleEntry.MaxTTL != 0 && roleEntry.TTL > roleEntry.MaxTTL {
		return logical.ErrorResponse("ttl cannot be greater than max_ttl"), nil
	}

	// if we're creating a role to manage a Team or Organization, we need to
	// create the token now. User tokens will be created when credentials are
	// read.
	if roleEntry.Organization != "" || roleEntry.TeamID != "" {
		token, err := b.createToken(ctx, req.Storage, roleEntry)
		if err != nil {
			return nil, err
		}

		roleEntry.Token = token.Token
		roleEntry.TokenID = token.ID
	}

	if err := setRole(ctx, req.Storage, name, roleEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *tfBackend) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "role/"+d.Get("name").(string))
	if err != nil {
		return nil, fmt.Errorf("error deleting terraform role: %w", err)
	}

	return nil, nil
}

func setRole(ctx context.Context, s logical.Storage, name string, roleEntry *terraformRoleEntry) error {
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

func (b *tfBackend) getRole(ctx context.Context, s logical.Storage, name string) (*terraformRoleEntry, error) {
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

	var role terraformRoleEntry

	if err := entry.DecodeJSON(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

const (
	pathRoleHelpSynopsis    = `Manages the Vault role for generating Terraform Cloud / Enterprise tokens.`
	pathRoleHelpDescription = `
This path allows you to read and write roles used to generate Terraform Cloud /
Enterprise tokens. You can configure a role to manage an organization's token, a
team's token, or a user's dynamic tokens.

A Terraform Cloud/Enterprise Organization can only have one active token at a
time. To manage an Organization's token, set the organization field.

A Terraform Cloud/Enterprise Team can only have one active token at a time. To
manage a Teams's token, set the team_id field.

A Terraform Cloud/Enterprise User can have multiple API tokens. To manage a
User's token, set the user_id field.
`

	pathRoleListHelpSynopsis    = `List the existing roles in Terraform Cloud / Enterprise backend`
	pathRoleListHelpDescription = `Roles will be listed by the role name.`
)
