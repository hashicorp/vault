package ldap

import (
	"fmt"

	"strings"

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
			Root: mfa.MFARootPaths(),

			Unauthenticated: []string{
				"login/*",
			},
		},

		Paths: append([]*framework.Path{
			pathConfig(&b),
			pathGroups(&b),
			pathGroupsList(&b),
			pathUsers(&b),
			pathUsersList(&b),
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

	bindDN, err := getBindDN(cfg, c, username)
	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}

	if err = c.Bind(bindDN, password); err != nil {
		return nil, logical.ErrorResponse(fmt.Sprintf("LDAP bind failed: %v", err)), nil
	}

	userDN, err := getUserDN(cfg, c, bindDN)
	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}

	ldapGroups, err := getLdapGroups(cfg, c, userDN, username)
	if err != nil {
		return nil, logical.ErrorResponse(err.Error()), nil
	}

	ldapResponse := &logical.Response{
		Data: map[string]interface{}{},
	}
	if len(ldapGroups) == 0 {
		errString := fmt.Sprintf(
			"no LDAP groups found in userDN '%s' or groupDN '%s';only policies from locally-defined groups available",
			cfg.UserDN,
			cfg.GroupDN)
		ldapResponse.AddWarning(errString)
	}

	var allGroups []string
	// Import the custom added groups from ldap backend
	user, err := b.User(req.Storage, username)
	if err == nil && user != nil {
		allGroups = append(allGroups, user.Groups...)
	}
	// add the LDAP groups
	allGroups = append(allGroups, ldapGroups...)

	// Retrieve policies
	var policies []string
	for _, groupName := range allGroups {
		group, err := b.Group(req.Storage, groupName)
		if err == nil && group != nil {
			policies = append(policies, group.Policies...)
		}
	}

	if len(policies) == 0 {
		errStr := "user is not a member of any authorized group"
		if len(ldapResponse.Warnings()) > 0 {
			errStr = fmt.Sprintf("%s; additionally, %s", errStr, ldapResponse.Warnings()[0])
		}

		ldapResponse.Data["error"] = errStr
		return nil, ldapResponse, nil
	}

	return policies, ldapResponse, nil
}

func getBindDN(cfg *ConfigEntry, c *ldap.Conn, username string) (string, error) {
	bindDN := ""
	if cfg.DiscoverDN || (cfg.BindDN != "" && cfg.BindPassword != "") {
		if err := c.Bind(cfg.BindDN, cfg.BindPassword); err != nil {
			return bindDN, fmt.Errorf("LDAP bind (service) failed: %v", err)
		}
		result, err := c.Search(&ldap.SearchRequest{
			BaseDN: cfg.UserDN,
			Scope:  2, // subtree
			Filter: fmt.Sprintf("(%s=%s)", cfg.UserAttr, ldap.EscapeFilter(username)),
		})
		if err != nil {
			return bindDN, fmt.Errorf("LDAP search for binddn failed: %v", err)
		}
		if len(result.Entries) != 1 {
			return bindDN, fmt.Errorf("LDAP search for binddn 0 or not unique")
		}
		bindDN = result.Entries[0].DN
	} else {
		if cfg.UPNDomain != "" {
			bindDN = fmt.Sprintf("%s@%s", EscapeLDAPValue(username), cfg.UPNDomain)
		} else {
			bindDN = fmt.Sprintf("%s=%s,%s", cfg.UserAttr, EscapeLDAPValue(username), cfg.UserDN)
		}
	}

	return bindDN, nil
}

func getUserDN(cfg *ConfigEntry, c *ldap.Conn, bindDN string) (string, error) {
	userDN := ""
	if cfg.UPNDomain != "" {
		// Find the distinguished name for the user if userPrincipalName used for login
		result, err := c.Search(&ldap.SearchRequest{
			BaseDN: cfg.UserDN,
			Scope:  2, // subtree
			Filter: fmt.Sprintf("(userPrincipalName=%s)", ldap.EscapeFilter(bindDN)),
		})
		if err != nil {
			return userDN, fmt.Errorf("LDAP search failed for detecting user: %v", err)
		}
		for _, e := range result.Entries {
			userDN = e.DN
		}
	} else {
		userDN = bindDN
	}

	return userDN, nil
}

func getLdapGroups(cfg *ConfigEntry, c *ldap.Conn, userDN string, username string) ([]string, error) {
	// retrieve the groups in a string/bool map as a structure to avoid duplicates inside
	ldapMap := make(map[string]bool)
	// Fetch the optional memberOf property values on the user object
	// This is the most common method used in Active Directory setup to retrieve the groups
	result, err := c.Search(&ldap.SearchRequest{
		BaseDN: userDN,
		Scope:  0,        // base scope to fetch only the userDN
		Filter: "(cn=*)", // bogus filter, required to fetch the CN from userDN
		Attributes: []string{
			"memberOf",
		},
	})
	// this check remains in case something happens with the ldap query or connection
	if err != nil {
		return nil, fmt.Errorf("LDAP fetch of distinguishedName=%s failed: %v", userDN, err)
	}
	// if there are more than one entry, we consider the results irrelevant and ignore them
	if len(result.Entries) == 1 {
		for _, attr := range result.Entries[0].Attributes {
			// Find the groups the user is member of from the 'memberOf' attribute extracting the CN
			if attr.Name == "memberOf" {
				for _, value := range attr.Values {
					memberOfDN, err := ldap.ParseDN(value)
					if err != nil || len(memberOfDN.RDNs) == 0 {
						continue
					}

					for _, rdn := range memberOfDN.RDNs {
						for _, rdnTypeAndValue := range rdn.Attributes {
							if strings.EqualFold(rdnTypeAndValue.Type, "CN") {
								ldapMap[rdnTypeAndValue.Value] = true
							}
						}
					}
				}
			}
		}
	}

	// Find groups by searching in groupDN for any of the memberUid, member or uniqueMember attributes
	// and retrieving the CN in the DN result
	if cfg.GroupDN != "" {
		result, err := c.Search(&ldap.SearchRequest{
			BaseDN: cfg.GroupDN,
			Scope:  2, // subtree
			Filter: fmt.Sprintf("(|(memberUid=%s)(member=%s)(uniqueMember=%s))", ldap.EscapeFilter(username), ldap.EscapeFilter(userDN), ldap.EscapeFilter(userDN)),
		})
		if err != nil {
			return nil, fmt.Errorf("LDAP search failed: %v", err)
		}

		for _, e := range result.Entries {
			dn, err := ldap.ParseDN(e.DN)
			if err != nil || len(dn.RDNs) == 0 {
				continue
			}
			for _, rdn := range dn.RDNs {
				for _, rdnTypeAndValue := range rdn.Attributes {
					if strings.EqualFold(rdnTypeAndValue.Type, "CN") {
						ldapMap[rdnTypeAndValue.Value] = true
					}
				}
			}
		}
	}

	ldapGroups := make([]string, len(ldapMap))
	for key, _ := range ldapMap {
		ldapGroups = append(ldapGroups, key)
	}
	return ldapGroups, nil
}

const backendHelp = `
The "ldap" credential provider allows authentication querying
a LDAP server, checking username and password, and associating groups
to set of policies.

Configuration of the server is done through the "config" and "groups"
endpoints by a user with root access. Authentication is then done
by suppying the two fields for "login".
`
