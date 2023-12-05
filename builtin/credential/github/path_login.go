// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGithub,
			OperationVerb:   "login",
		},

		Fields: map[string]*framework.FieldSchema{
			"token": {
				Type:        framework.TypeString,
				Description: "GitHub personal API token",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.pathLoginAliasLookahead,
		},
	}
}

func (b *backend) pathLoginAliasLookahead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	token := data.Get("token").(string)

	verifyResp, err := b.verifyCredentials(ctx, req, token)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Warnings: verifyResp.Warnings,
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: *verifyResp.User.Login,
			},
		},
	}, nil
}

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	token := data.Get("token").(string)

	verifyResp, err := b.verifyCredentials(ctx, req, token)
	if err != nil {
		return nil, err
	}

	auth := &logical.Auth{
		InternalData: map[string]interface{}{
			"token": token,
		},
		Metadata: map[string]string{
			"username": *verifyResp.User.Login,
			"org":      *verifyResp.Org.Login,
		},
		DisplayName: *verifyResp.User.Login,
		Alias: &logical.Alias{
			Name: *verifyResp.User.Login,
		},
	}
	verifyResp.Config.PopulateTokenAuth(auth)

	// Add in configured policies from user/group mapping
	if len(verifyResp.Policies) > 0 {
		auth.Policies = append(auth.Policies, verifyResp.Policies...)
	}

	resp := &logical.Response{
		Warnings: verifyResp.Warnings,
		Auth:     auth,
	}

	for _, teamName := range verifyResp.TeamNames {
		if teamName == "" {
			continue
		}
		resp.Auth.GroupAliases = append(resp.Auth.GroupAliases, &logical.Alias{
			Name: teamName,
		})
	}

	return resp, nil
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	if req.Auth == nil {
		return nil, fmt.Errorf("request auth was nil")
	}

	tokenRaw, ok := req.Auth.InternalData["token"]
	if !ok {
		return nil, fmt.Errorf("token created in previous version of Vault cannot be validated properly at renewal time")
	}
	token := tokenRaw.(string)

	verifyResp, err := b.verifyCredentials(ctx, req, token)
	if err != nil {
		return nil, err
	}

	if !policyutil.EquivalentPolicies(verifyResp.Policies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies do not match")
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.Period = verifyResp.Config.TokenPeriod
	resp.Auth.TTL = verifyResp.Config.TokenTTL
	resp.Auth.MaxTTL = verifyResp.Config.TokenMaxTTL
	resp.Warnings = verifyResp.Warnings

	// Remove old aliases
	resp.Auth.GroupAliases = nil

	for _, teamName := range verifyResp.TeamNames {
		resp.Auth.GroupAliases = append(resp.Auth.GroupAliases, &logical.Alias{
			Name: teamName,
		})
	}

	return resp, nil
}

func (b *backend) verifyCredentials(ctx context.Context, req *logical.Request, token string) (*verifyCredentialsResp, error) {
	var warnings []string
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("configuration has not been set")
	}

	// Check for a CIDR match.
	if len(config.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Error("token bound CIDRs found but no connection information available for validation")
			return nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, config.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}
	}

	client, err := b.Client(token)
	if err != nil {
		return nil, err
	}

	if config.BaseURL != "" {
		parsedURL, err := url.Parse(config.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("successfully parsed base_url when set but failing to parse now: %w", err)
		}
		client.BaseURL = parsedURL
	}

	if config.OrganizationID == 0 {
		// Previously we did not verify using the Org ID. So if the Org ID is
		// not set, we will trust-on-first-use and set it now.
		err = config.setOrganizationID(ctx, client)
		if err != nil {
			b.Logger().Error("failed to set the organization_id on login", "error", err)
			return nil, err
		}
		entry, err := logical.StorageEntryJSON("config", config)
		if err != nil {
			return nil, err
		}

		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		b.Logger().Info("set ID on a trust-on-first-use basis", "organization_id", config.OrganizationID)
	}

	// Get the user
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	// Verify that the user is part of the organization
	var org *github.Organization

	orgOpt := &github.ListOptions{
		PerPage: 100,
	}

	var allOrgs []*github.Organization
	for {
		orgs, resp, err := client.Organizations.List(ctx, "", orgOpt)
		if err != nil {
			return nil, err
		}
		allOrgs = append(allOrgs, orgs...)
		if resp.NextPage == 0 {
			break
		}
		orgOpt.Page = resp.NextPage
	}

	orgLoginName := ""
	for _, o := range allOrgs {
		if o.GetID() == config.OrganizationID {
			org = o
			orgLoginName = *o.Login
			break
		}
	}
	if org == nil {
		return nil, errors.New("user is not part of required org")
	}

	if orgLoginName != config.Organization {
		warningMsg := fmt.Sprintf(
			"the organization name has changed to %q. It is recommended to verify and update the organization name in the config: %s=%d",
			orgLoginName,
			"organization_id",
			config.OrganizationID,
		)
		b.Logger().Warn(warningMsg)
		warnings = append(warnings, warningMsg)
	}

	// Get the teams that this user is part of to determine the policies
	var teamNames []string

	teamOpt := &github.ListOptions{
		PerPage: 100,
	}

	var allTeams []*github.Team
	for {
		teams, resp, err := client.Teams.ListUserTeams(ctx, teamOpt)
		if err != nil {
			return nil, err
		}
		allTeams = append(allTeams, teams...)
		if resp.NextPage == 0 {
			break
		}
		teamOpt.Page = resp.NextPage
	}

	for _, t := range allTeams {
		// We only care about teams that are part of the organization we use
		if *t.Organization.ID != *org.ID {
			continue
		}

		// Append the names so we can get the policies
		teamNames = append(teamNames, *t.Name)
		if *t.Name != *t.Slug {
			teamNames = append(teamNames, *t.Slug)
		}
	}

	groupPoliciesList, err := b.TeamMap.Policies(ctx, req.Storage, teamNames...)
	if err != nil {
		return nil, err
	}

	userPoliciesList, err := b.UserMap.Policies(ctx, req.Storage, []string{*user.Login}...)
	if err != nil {
		return nil, err
	}

	verifyResp := &verifyCredentialsResp{
		User:      user,
		Org:       org,
		Policies:  append(groupPoliciesList, userPoliciesList...),
		TeamNames: teamNames,
		Config:    config,
		Warnings:  warnings,
	}

	return verifyResp, nil
}

type verifyCredentialsResp struct {
	User      *github.User
	Org       *github.Organization
	Policies  []string
	TeamNames []string

	// Warnings to send back to the caller
	Warnings []string

	// This is just a cache to send back to the caller
	Config *config
}
