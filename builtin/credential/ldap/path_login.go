package ldap

import (
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login/(?P<username>.+)`,
		Fields: map[string]*framework.FieldSchema{
			"username": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "DN (distinguished name) to be used for login.",
			},

			"password": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Password for this user.",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.WriteOperation: b.pathLogin,
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
	if len(policies) == 0 {
		return resp, err
	}

	sort.Strings(policies)

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: policies,
			Metadata: map[string]string{
				"username": username,
				"policies": strings.Join(policies, ","),
			},
			InternalData: map[string]interface{}{
				"password": password,
			},
			DisplayName: username,
		},
	}, nil
}

func (b *backend) pathLoginRenew(
	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	username := req.Auth.Metadata["username"]
	password := req.Auth.InternalData["password"].(string)
	prevpolicies := req.Auth.Metadata["policies"]

	policies, resp, err := b.Login(req, username, password)
	if len(policies) == 0 {
		return resp, err
	}

	sort.Strings(policies)
	if strings.Join(policies, ",") != prevpolicies {
		return logical.ErrorResponse("policies have changed, revoking login"), nil
	}

	return framework.LeaseExtend(1*time.Hour, 0, false)(req, d)
}

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
