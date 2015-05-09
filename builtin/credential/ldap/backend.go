package ldap

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/vanackere/ldap"
)

func Factory(map[string]string) (logical.Backend, error) {
	return Backend(), nil
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Root: []string{
				"config",
				"groups/*",
			},

			Unauthenticated: []string{
				"login/*",
			},
		},

		Paths: append([]*framework.Path{
			pathLogin(&b),
			pathConfig(&b),
			pathGroups(&b),
		}),

		AuthRenew: b.pathLoginRenew,
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}

func (b *backend) Login(req *logical.Request, username string, password string) ([]string, *logical.Response, error) {

	cfg, err := b.Config(req)
	if err != nil {
		return nil, nil, err
	}
	if cfg == nil {
		return nil, logical.ErrorResponse("ldap backend not configured"), nil
	}

	c, err := cfg.DialLDAP()
	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}

	// Try to authenticate to the server using the provided credentials
	binddn := fmt.Sprintf("%s=%s,%s", cfg.UserAttr, username, cfg.UserDN)
	if err = c.Bind(binddn, password); err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("LDAP bind failed: %v", err)), nil
	}

	// Enumerate all groups the user is member of. The search filter should
	// work with both openldap and MS AD standard schemas.
	sresult, err := c.Search(&ldap.SearchRequest{
		BaseDN: cfg.GroupDN,
		Scope:  2, // subtree
		Filter: fmt.Sprintf("(|(memberUid=%s)(member=%s)(uniqueMember=%s))", username, binddn, binddn),
	})
	if err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("LDAP search failed: %v", err)), nil
	}

	var allgroups []string
	var policies []string
	for _, e := range sresult.Entries {
		// Expected syntax for group DN: cn=groupanem,ou=Group,dc=example,dc=com
		dn := strings.Split(e.DN, ",")
		gname := strings.SplitN(dn[0], "=", 2)[1]
		allgroups = append(allgroups, gname)
		group, err := b.Group(req.Storage, gname)
		if err == nil && group != nil {
			policies = append(policies, group.Policies...)
		}
	}

	if len(policies) == 0 {
		return nil, logical.ErrorResponse("user is not member of any authorized group"), nil
	}

	return policies, nil, nil
}

const backendHelp = `
The "ldap" credential provider allows authentication querying
a LDAP server, checking username and password, and associating groups
to set of policies.

Configuration of the server is done through the "config" and "groups"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
