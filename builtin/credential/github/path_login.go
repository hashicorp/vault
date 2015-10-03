package github

import (
	"fmt"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login",
		Fields: map[string]*framework.FieldSchema{
			"token": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "GitHub personal API token",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathLogin,
		},
	}
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get all our stored state
	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	if config.Org == "" {
		return logical.ErrorResponse(
			"configure the github credential backend first"), nil
	}

	client, err := b.Client(data.Get("token").(string))
	if err != nil {
		return nil, err
	}

	if config.BaseURL != "" {
		parsedURL, err := url.Parse(config.BaseURL)
		if err != nil {
			return nil, fmt.Errorf("Successfully parsed base_url when set but failing to parse now: %s", err)
		}
		client.BaseURL = parsedURL
	}

	// Get the user
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, err
	}

	// Verify that the user is part of the organization
	var org *github.Organization

	orgOpt := &github.ListOptions{
		PerPage: 100,
	}

	var allOrgs []github.Organization
	for {
		orgs, resp, err := client.Organizations.List("", orgOpt)
		if err != nil {
			return nil, err
		}
		allOrgs = append(allOrgs, orgs...)
		if resp.NextPage == 0 {
			break
		}
		orgOpt.Page = resp.NextPage
	}

	for _, o := range allOrgs {
		if *o.Login == config.Org {
			org = &o
			break
		}
	}
	if org == nil {
		return logical.ErrorResponse("user is not part of required org"), nil
	}

	// Get the teams that this user is part of to determine the policies
	var teamNames []string

	teamOpt := &github.ListOptions{
		PerPage: 100,
	}

	var allTeams []github.Team
	for {
		teams, resp, err := client.Organizations.ListUserTeams(teamOpt)
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

	policiesList, err := b.Map.Policies(req.Storage, teamNames...)
	if err != nil {
		return nil, err
	}

	ttl, _, err := b.SanitizeTTL(config.TTL.String(), config.MaxTTL.String())
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("[ERR]:%s", err)), nil
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: policiesList,
			Metadata: map[string]string{
				"username": *user.Login,
				"org":      *org.Login,
			},
			DisplayName: *user.Login,
			LeaseOptions: logical.LeaseOptions{
				TTL:         ttl,
				GracePeriod: ttl / 10,
				Renewable:   ttl > 0,
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	return framework.LeaseExtend(config.MaxTTL, 0, false)(req, d)
}
