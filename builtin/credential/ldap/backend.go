package ldap

import (
	"fmt"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/vault/helper/mfa"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"strings"
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

	// Format binddn
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
			return nil, logical.ErrorResponse("LDAP search for binddn 0 or not unique"), nil
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
	ldapGroups := make(map[string]bool)
	if cfg.UPNDomain != "" {
		// Find the distinguished name for the user if userPrincipalName used for login
		// and the groups from memberOf attributes
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
			// Find the groups the user is member of from the 'memberOf' attribute extracting the CN
			for _,dnAttr := range e.Attributes {
				if dnAttr.Name == "memberOf" {
					for _,value := range dnAttr.Values {
						memberOfDN, err := ldap.ParseDN(value)
						if err != nil {
							return nil, nil, err
						}

						for _, rdn := range memberOfDN.RDNs {
							for _, rdnTypeAndValue := range rdn.Attributes {
								if strings.EqualFold(rdnTypeAndValue.Type, "CN") {
									ldapGroups[rdnTypeAndValue.Value] = true
								}
							}
						}
					}
				}
			}
		}
	} else {
		userdn = binddn
	}

	resp := &logical.Response{
		Data: map[string]interface{}{},
	}
	// Find groups by searching in groupDN for any of the memberUid, member or uniqueMember attributes
	// and retrieving the CN in the DN result
	if cfg.GroupDN != "" {
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
			if err != nil {
				return nil, nil, err
			}
			for _, rdn := range dn.RDNs {
				for _, rdnTypeAndValue := range rdn.Attributes {
					if strings.EqualFold(rdnTypeAndValue.Type, "CN" ) {
						ldapGroups[rdnTypeAndValue.Value] = true
					}
				}
			}
		}
	}
	
	if len(ldapGroups) == 0 {
		errString := fmt.Sprintf(
			"no LDAP groups found in user DN '%s' or group DN '%s';only policies from locally-defined groups available",
			cfg.UserDN,
			cfg.GroupDN)
		resp.AddWarning(errString)
	}

	var allgroups []string
	// Import the custom added groups from ldap backend
	user, err := b.User(req.Storage, username)
	if err == nil && user != nil {
		allgroups = append(allgroups, user.Groups...)
	}
	// add the LDAP groups
	for key, _ := range ldapGroups {
		allgroups = append(allgroups, key)
	}

	// Retrieve policies
	var policies []string
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
