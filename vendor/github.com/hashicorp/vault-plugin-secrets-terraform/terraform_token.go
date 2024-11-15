// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfc

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	terraformTokenType = "terraform_token"
)

func isOrgToken(organization string, team string) bool {
	return organization != "" && team == ""
}

func isTeamToken(team string) bool {
	return team != ""
}

func createOrgToken(ctx context.Context, c *client, organization string) (*terraformToken, error) {
	if _, err := c.Organizations.Read(ctx, organization); err != nil {
		return nil, err
	}

	token, err := c.OrganizationTokens.Create(ctx, organization)
	if err != nil {
		return nil, err
	}

	return &terraformToken{
		ID:          token.ID,
		Description: token.Description,
		Token:       token.Token,
	}, nil
}

func createTeamToken(ctx context.Context, c *client, teamID string) (*terraformToken, error) {
	if _, err := c.Teams.Read(ctx, teamID); err != nil {
		return nil, err
	}

	token, err := c.TeamTokens.Create(ctx, teamID)
	if err != nil {
		return nil, err
	}

	return &terraformToken{
		ID:          token.ID,
		Description: token.Description,
		Token:       token.Token,
	}, nil
}

func createUserToken(ctx context.Context, c *client, userID string) (*terraformToken, error) {
	token, err := c.UserTokens.Create(ctx, userID, tfe.UserTokenCreateOptions{})
	if err != nil {
		return nil, err
	}

	return &terraformToken{
		ID:          token.ID,
		Description: token.Description,
		Token:       token.Token,
	}, nil
}

func (b *tfBackend) terraformTokenRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	organization := ""
	organizationRaw, ok := req.Secret.InternalData["organization"]
	if ok {
		organization, ok = organizationRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid value for organization in secret internal data")
		}
	}

	teamID := ""
	teamIDRaw, ok := req.Secret.InternalData["team_id"]
	if ok {
		teamID, ok = teamIDRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid value for team_id in secret internal data")
		}

	}

	if isOrgToken(organization, teamID) {
		// revoke org API token
		if err := client.OrganizationTokens.Delete(ctx, organization); err != nil {
			return nil, fmt.Errorf("error revoking organization token: %w", err)
		}
		return nil, nil
	}

	if isTeamToken(teamID) {
		// revoke team API token
		if err := client.TeamTokens.Delete(ctx, teamID); err != nil {
			return nil, fmt.Errorf("error revoking team token: %w", err)
		}
		return nil, nil
	}

	// if we haven't returned yet, then the token is a user API token
	tokenID := ""
	tokenIDRaw, ok := req.Secret.InternalData["token_id"]
	if ok {
		tokenID, ok = tokenIDRaw.(string)
		if !ok {
			return nil, fmt.Errorf("secret is missing tokenID internal data")
		}
	}

	if tokenID == "" {
		return nil, fmt.Errorf("secret is missing tokenID internal data")
	}

	if err := client.UserTokens.Delete(ctx, tokenID); err != nil {
		return nil, fmt.Errorf("error revoking user token: %w", err)
	}
	return nil, nil
}

func (b *tfBackend) terraformTokenRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleRaw, ok := req.Secret.InternalData["role"]
	if !ok {
		return nil, fmt.Errorf("secret is missing role internal data")
	}

	// get the role entry
	role := roleRaw.(string)
	roleEntry, err := b.getRole(ctx, req.Storage, role)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %w", err)
	}

	if roleEntry == nil {
		return nil, errors.New("error retrieving role: role is nil")
	}

	resp := &logical.Response{Secret: req.Secret}

	if roleEntry.TTL > 0 {
		resp.Secret.TTL = roleEntry.TTL
	}
	if roleEntry.MaxTTL > 0 {
		resp.Secret.MaxTTL = roleEntry.MaxTTL
	}

	return resp, nil
}
