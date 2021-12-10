package github

import (
	"context"
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

	var verifyResp *verifyCredentialsResp
	if verifyResponse, resp, err := b.verifyCredentials(ctx, req, token); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		verifyResp = verifyResponse
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: *verifyResp.User.Login,
			},
		},
	}, nil
}

func (b *backend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	token := data.Get("token").(string)

	var verifyResp *verifyCredentialsResp
	if verifyResponse, resp, err := b.verifyCredentials(ctx, req, token); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		verifyResp = verifyResponse
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
		Auth: auth,
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

	var verifyResp *verifyCredentialsResp
	if verifyResponse, resp, err := b.verifyCredentials(ctx, req, token); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		verifyResp = verifyResponse
	}
	if !policyutil.EquivalentPolicies(verifyResp.Policies, req.Auth.TokenPolicies) {
		return nil, fmt.Errorf("policies do not match")
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.Period = verifyResp.Config.TokenPeriod
	resp.Auth.TTL = verifyResp.Config.TokenTTL
	resp.Auth.MaxTTL = verifyResp.Config.TokenMaxTTL

	// Remove old aliases
	resp.Auth.GroupAliases = nil

	for _, teamName := range verifyResp.TeamNames {
		resp.Auth.GroupAliases = append(resp.Auth.GroupAliases, &logical.Alias{
			Name: teamName,
		})
	}

	return resp, nil
}

func (b *backend) verifyCredentials(ctx context.Context, req *logical.Request, token string) (*verifyCredentialsResp, *logical.Response, error) {
	var resp logical.Response
	config, err := b.Config(ctx, req.Storage)
	if err != nil {
		return nil, nil, err
	}
	if config == nil {
		return nil, logical.ErrorResponse("configuration has not been set"), nil
	}

	// Check for a CIDR match.
	if len(config.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, config.TokenBoundCIDRs) {
			return nil, nil, logical.ErrPermissionDenied
		}
	}

	client, err := b.Client(token)
	if err != nil {
		return nil, nil, err
	}

	if config.BaseURL != "" {
		parsedURL, err := url.Parse(config.BaseURL)
		if err != nil {
			return nil, nil, fmt.Errorf("successfully parsed base_url when set but failing to parse now: %w", err)
		}
		client.BaseURL = parsedURL
	}

	if config.OrganizationID == 0 {
		// Previously we did not verify using the Org ID. So if the Org ID is
		// not set, we will trust-on-first-use and set it now.
		err = config.setOrganizationID(ctx, client)
		if err != nil {
			b.Logger().Error("failed to set the organization_id on login", "error", err)
			return nil, nil, err
		} else {
			entry, err := logical.StorageEntryJSON("config", config)
			if err != nil {
				return nil, nil, err
			}

			if err := req.Storage.Put(ctx, entry); err != nil {
				return nil, nil, err
			}
		}

	}

	// Get the user
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, nil, err
	}

	// Verify that the user is part of the organization
	var org *github.Organization

	orgOpt := &github.ListOptions{
		PerPage: 100,
	}

	var allOrgs []*github.Organization
	for {
		orgs, listResp, err := client.Organizations.List(ctx, "", orgOpt)
		if err != nil {
			return nil, nil, err
		}
		allOrgs = append(allOrgs, orgs...)
		if listResp.NextPage == 0 {
			break
		}
		orgOpt.Page = listResp.NextPage
	}

	orgLoginName := ""
	for _, o := range allOrgs {
		if o != nil && *o.ID == config.OrganizationID {
			org = o
			orgLoginName = *o.Login
			break
		}
	}
	if org == nil {
		return nil, logical.ErrorResponse("user is not part of required org"), nil
	}

	if orgLoginName != config.Organization {
		warningMsg := fmt.Sprintf(
			"the organization name has changed to %q. It is recommended to verify and update the organization name in the config: %s=%d",
			orgLoginName,
			"organization_id",
			config.OrganizationID,
		)
		b.Logger().Warn(warningMsg)
		resp.AddWarning(warningMsg)
	}

	// Get the teams that this user is part of to determine the policies
	var teamNames []string

	teamOpt := &github.ListOptions{
		PerPage: 100,
	}

	var allTeams []*github.Team
	for {
		teams, listResp, err := client.Teams.ListUserTeams(ctx, teamOpt)
		if err != nil {
			return nil, nil, err
		}
		allTeams = append(allTeams, teams...)
		if listResp.NextPage == 0 {
			break
		}
		teamOpt.Page = listResp.NextPage
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
		return nil, nil, err
	}

	userPoliciesList, err := b.UserMap.Policies(ctx, req.Storage, []string{*user.Login}...)
	if err != nil {
		return nil, nil, err
	}

	verifyResp := &verifyCredentialsResp{
		User:      user,
		Org:       org,
		Policies:  append(groupPoliciesList, userPoliciesList...),
		TeamNames: teamNames,
		Config:    config,
	}

	return verifyResp, &resp, nil
}

type verifyCredentialsResp struct {
	User      *github.User
	Org       *github.Organization
	Policies  []string
	TeamNames []string

	// This is just a cache to send back to the caller
	Config *config
}
