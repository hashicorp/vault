package okta

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-errors/errors"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"time"
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
			logical.UpdateOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLogin(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	policies, resp, err := b.Login(req, username, password)
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

	ttl, _, err := b.getTTLs(req)
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
			TTL:       ttl,
			Renewable: true,
		},
	}
	return resp, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	username := req.Auth.Metadata["username"]
	password := req.Auth.InternalData["password"].(string)

	loginPolicies, resp, err := b.Login(req, username, password)
	if len(loginPolicies) == 0 {
		return resp, err
	}

	if !policyutil.EquivalentPolicies(loginPolicies, req.Auth.Policies) {
		return nil, fmt.Errorf("policies have changed, not renewing")
	}

	ttl, maxTTL, err := b.getTTLs(req)
	if err != nil {
		return nil, err
	}

	return framework.LeaseExtend(ttl, maxTTL, b.System())(req, d)
}

func (b *backend) getTTLs(req *logical.Request) (ttl, maxTTL time.Duration, err error) {

	cfg, err := b.Config(req.Storage)
	if err != nil {
		return 0, 0, err
	}
	if cfg == nil {
		return 0, 0, errors.New("Okta backend not configured")
	}

	ttl, maxTTL, err = b.SanitizeTTLStr(cfg.TTL.String(), cfg.MaxTTL.String())

	return ttl, maxTTL, err
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
