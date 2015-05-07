package ldap

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/vanackere/ldap"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `login`,
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

	cfg, err := b.Config(req)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return logical.ErrorResponse("ldap backend not configured"), nil
	}

	c, err := cfg.DialLDAP()
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Try to authenticate to the server using the provided credentials
	binddn := fmt.Sprintf("%s=%s,%s", cfg.UserAttr, username, cfg.UserDN)
	if err = c.Bind(binddn, password); err != nil {
		return logical.ErrorResponse(fmt.Sprintf("LDAP bind failed: %v", err)), nil
	}

	// Enumerate all groups the user is member of. The search filter should
	// work with both openldap and MS AD standard schemas.
	sresult, err := c.Search(&ldap.SearchRequest{
		BaseDN: cfg.GroupDN,
		Scope:  2, // subtree
		Filter: fmt.Sprintf("(|(memberUid=%s)(member=%s)(uniqueMember=%s))", username, binddn, binddn),
	})
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("LDAP search failed: %v", err)), nil
	}

	var policies []string
	for _, e := range sresult.Entries {
		// Expected syntax for group DN: cn=groupanem,ou=Group,dc=example,dc=com
		dn := strings.Split(e.DN, ",")
		group, err := b.Group(req.Storage, dn[0][3:])
		if err == nil && group != nil {
			policies = append(policies, group.Policies...)
		}
	}

	if len(policies) == 0 {
		return logical.ErrorResponse("user is not member of any authorized group"), nil
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Policies: policies,
			Metadata: map[string]string{
				"username": username,
			},
			DisplayName: username,
		},
	}, nil
}

// func (b *backend) pathLoginRenew(
// 	req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
// 	// Get the user and validate auth
// 	user, err := b.User(req.Storage, req.Auth.Metadata["username"])
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user == nil {
// 		// User no longer exists, do not renew
// 		return nil, nil
// 	}

// 	return framework.LeaseExtend(1*time.Hour, 0)(req, d)
// }

const pathLoginSyn = `
Log in with a username and password.
`

const pathLoginDesc = `
This endpoint authenticates using a username and password.
`
