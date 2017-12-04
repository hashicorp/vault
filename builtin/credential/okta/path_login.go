package okta

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-errors/errors"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login/(?P<username>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Username to be used for login.",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.pathLoginAliasLookahead,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLoginAliasLookahead(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	if username == "" {
		return nil, fmt.Errorf("missing username")
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: username,
			},
		},
	}, nil
}

func (b *backend) pathLogin(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	policies, resp, groupNames, err := b.Login(req, username, password)
	// Handle an internal error
	if err != nil {
		return nil, err
	}
	if resp != nil {
		// Handle a logical error
		if resp.IsError() {
			return resp, nil
		}
	} else {
		resp = &logical.Response{}
	}

	sort.Strings(policies)

	cfg, err := b.getConfig(req)
	if err != nil {
		return nil, err
	}

	resp.Auth = &logical.Auth{
		Policies: policies,
		Metadata: map[string]string{
			"username": username,
			"policies": strings.Join(policies, ","),
		},
		InternalData: map[string]interface{}{
			"password": password,
		},
		DisplayName: username,
		LeaseOptions: logical.LeaseOptions{
			TTL:       cfg.TTL,
			Renewable: true,
		},
		Alias: &logical.Alias{
			Name: username,
		},
	}

	for _, groupName := range groupNames {
		if groupName == "" {
			continue
		}
		resp.Auth.GroupAliases = append(resp.Auth.GroupAliases, &logical.Alias{
			Name: groupName,
		})
	}

	return resp, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	username := req.Auth.Metadata["username"]
	password := req.Auth.InternalData["password"].(string)

	loginPolicies, resp, groupNames, err := b.Login(req, username, password)
	if len(loginPolicies) == 0 {
		return resp, err
	}

	if !policyutil.EquivalentPolicies(loginPolicies, req.Auth.Policies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	cfg, err := b.getConfig(req)
	if err != nil {
		return nil, err
	}

	resp, err = framework.LeaseExtend(cfg.TTL, cfg.MaxTTL, b.System())(req, d)
	if err != nil {
		return nil, err
	}

	// Remove old aliases
	resp.Auth.GroupAliases = nil

	for _, groupName := range groupNames {
		resp.Auth.GroupAliases = append(resp.Auth.GroupAliases, &logical.Alias{
			Name: groupName,
		})
	}

	return resp, nil

}

func (b *backend) getConfig(req *logical.Request) (*ConfigEntry, error) {

	cfg, err := b.Config(req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("Okta backend not configured")
	}

	return cfg, nil
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
