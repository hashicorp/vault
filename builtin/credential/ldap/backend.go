package ldap

import (
	"fmt"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/vault/helper/mfa"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	return Backend().Setup(conf)
}

func Backend() *framework.Backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Root: append([]string{
				"config",
				"groups/*",
				"users/*",
			},
				mfa.MFARootPaths()...,
			),

			Unauthenticated: []string{
				"login/*",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathGroups(&b),
			pathUsers(&b),
		},
			mfa.MFAPaths(b.Backend, pathLogin(&b))...,
		),

		AuthRenew: b.pathLoginRenew,
	}

	return b.Backend
}

type backend struct {
	*framework.Backend
}

func EscapeLDAPValue(input string) string {
	// RFC4514 forbids un-escaped:
	// - leading space or hash
	// - trailing space
	// - special characters '"', '+', ',', ';', '<', '>', '\\'
	// - null
	for i := 0; i < len(input); i++ {
		escaped := false
		if input[i] == '\\' {
			i++
			escaped = true
		}
		switch input[i] {
		case '"', '+', ',', ';', '<', '>', '\\':
			if !escaped {
				input = input[0:i] + "\\" + input[i:]
				i++
			}
			continue
		}
		if escaped {
			input = input[0:i] + "\\" + input[i:]
			i++
		}
	}
	if input[0] == ' ' || input[0] == '#' {
		input = "\\" + input
	}
	if input[len(input)-1] == ' ' {
		input = input[0:len(input)-1] + "\\ "
	}
	return input
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
	if c == nil {
		return nil, logical.ErrorResponse("invalid connection returned from LDAP dial"), nil
	}
	binddn := ""
	if cfg.DiscoverDN || (cfg.BindDN != "" && cfg.BindPassword != "") {
		if err = c.Bind(cfg.BindDN, cfg.BindPassword); err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("LDAP bind (service) failed: %v", err)), nil
		}
		sresult, err := c.Search(&ldap.SearchRequest{
			BaseDN: cfg.UserDN,
			Scope:  2, // subtree
			Filter: fmt.Sprintf("(%s=%s)", cfg.UserAttr, ldap.EscapeFilter(username)),
		})
		if err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("LDAP search for binddn failed: %v", err)), nil
		}
		if len(sresult.Entries) != 1 {
			return nil, logical.ErrorResponse("LDAP search for binddn 0 or not uniq"), nil
		}
		binddn = sresult.Entries[0].DN
	} else {
		if cfg.UPNDomain != "" {
			binddn = fmt.Sprintf("%s@%s", EscapeLDAPValue(username), cfg.UPNDomain)
		} else {
			binddn = fmt.Sprintf("%s=%s,%s", cfg.UserAttr, EscapeLDAPValue(username), cfg.UserDN)
		}
	}
	if err = c.Bind(binddn, password); err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("LDAP bind failed: %v", err)), nil
	}

	userdn := ""
	if cfg.UPNDomain != "" {
		// Find the distinguished name for the user if userPrincipalName used for login
		sresult, err := c.Search(&ldap.SearchRequest{
			BaseDN: cfg.UserDN,
			Scope:  2, // subtree
			Filter: fmt.Sprintf("(userPrincipalName=%s)", ldap.EscapeFilter(binddn)),
		})
		if err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("LDAP search failed: %v", err)), nil
		}
		for _, e := range sresult.Entries {
			userdn = e.DN
		}
	} else {
		userdn = binddn
	}

	var allgroups []string
	var policies []string
	resp := &logical.Response{
		Data: map[string]interface{}{},
	}

	// Fetch custom (local) groups the user has been added to
	user, err := b.User(req.Storage, username)
	if err == nil && user != nil {
		allgroups = append(allgroups, user.Groups...)
	}

	if cfg.GroupDN != "" {
		// Enumerate all groups the user is member of. The search filter should
		// work with both openldap and MS AD standard schemas.
		sresult, err := c.Search(&ldap.SearchRequest{
			BaseDN: cfg.GroupDN,
			Scope:  2, // subtree
			Filter: fmt.Sprintf("(|(memberUid=%s)(member=%s)(uniqueMember=%s))", ldap.EscapeFilter(username), ldap.EscapeFilter(userdn), ldap.EscapeFilter(userdn)),
		})
		if err != nil {
			return nil, logical.ErrorResponse(fmt.Sprintf("LDAP search failed: %v", err)), nil
		}

		for _, e := range sresult.Entries {
			dn, err := ldap.ParseDN(e.DN)
			if err != nil || len(dn.RDNs) == 0 || len(dn.RDNs[0].Attributes) == 0 {
				continue
			}
			gname := dn.RDNs[0].Attributes[0].Value
			allgroups = append(allgroups, gname)
		}
	} else {
		resp.AddWarning("no group DN configured; only policies from locally-defined groups available")
	}

	for _, gname := range allgroups {
		group, err := b.Group(req.Storage, gname)
		if err == nil && group != nil {
			policies = append(policies, group.Policies...)
		}
	}

	if len(policies) == 0 {
		errStr := "user is not a member of any authorized group"
		if len(resp.Warnings()) > 0 {
			errStr = fmt.Sprintf("%s; additionally, %s", errStr, resp.Warnings()[0])
		}

		resp.Data["error"] = errStr
		return nil, resp, nil
	}

	return policies, resp, nil
}

const backendHelp = `
The "ldap" credential provider allows authentication querying
a LDAP server, checking username and password, and associating groups
to set of policies.

Configuration of the server is done through the "config" and "groups"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
