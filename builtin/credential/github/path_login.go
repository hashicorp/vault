package github

import (
	"fmt"
	"net/url"

	"github.com/google/go-github/github"
	"github.com/hashicorp/vault/helper/policyutil"
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
			logical.UpdateOperation: b.pathLogin,
		},
	}
}

func (b *backend) pathLogin(
	req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	token := data.Get("token").(string)

	var verifyResp *verifyCredentialsResp
	if verifyResponse, resp, err := b.verifyCredentials(req, token); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		verifyResp = verifyResponse
	}

	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}

	ttl, _, err := b.SanitizeTTL(config.TTL.String(), config.MaxTTL.String())
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("[ERR]:%s", err)), nil
	}

	return &logical.Response{
		Auth: &logical.Auth{
			InternalData: map[string]interface{}{
				"token": token,
			},
			Policies: verifyResp.Policies,
			Metadata: map[string]string{
				"username": *verifyResp.User.Login,
				"org":      *verifyResp.Org.Login,
			},
			DisplayName: *verifyResp.User.Login,
			LeaseOptions: logical.LeaseOptions{
				TTL:       ttl,
				Renewable: true,
			},
		},
	}, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	token := req.Auth.InternalData["token"].(string)

	var verifyResp *verifyCredentialsResp
	if verifyResponse, resp, err := b.verifyCredentials(req, token); err != nil {
		return nil, err
	} else if resp != nil {
		return resp, nil
	} else {
		verifyResp = verifyResponse
	}
	if !policyutil.EquivalentPolicies(verifyResp.Policies, req.Auth.Policies) {
		return logical.ErrorResponse("policies do not match"), nil
	}

	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	return framework.LeaseExtend(config.TTL, config.MaxTTL, b.System())(req, d)
}

func (b *backend) verifyCredentials(req *logical.Request, token string) (*verifyCredentialsResp, *logical.Response, error) {
	config, err := b.Config(req.Storage)
	if err != nil {
		return nil, nil, err
	}
	if config.Org == "" {
		return nil, logical.ErrorResponse(
			"configure the github credential backend first"), nil
	}

	client, err := b.Client(token)
	if err != nil {
		return nil, nil, err
	}

	if config.BaseURL != "" {
		parsedURL, err := url.Parse(config.BaseURL)
		if err != nil {
			return nil, nil, fmt.Errorf("Successfully parsed base_url when set but failing to parse now: %s", err)
		}
		client.BaseURL = parsedURL
	}

	// Get the user
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, nil, err
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
			return nil, nil, err
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
		return nil, logical.ErrorResponse("user is not part of required org"), nil
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
			return nil, nil, err
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
		return nil, nil, err
	}
	return &verifyCredentialsResp{
		User:     user,
		Org:      org,
		Policies: policiesList,
	}, nil, nil
}

type verifyCredentialsResp struct {
	User     *github.User
	Org      *github.Organization
	Policies []string
}
