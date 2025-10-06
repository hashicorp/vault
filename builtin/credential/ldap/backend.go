// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ldap

import (
	"context"
	"fmt"
	"strings"
	"sync"

	goldap "github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/cap/ldap"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	operationPrefixLDAP   = "ldap"
	errUserBindFailed     = "ldap operation failed: failed to bind as user"
	defaultPasswordLength = 64 // length to use for configured root password on rotations by default
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: backendHelp,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login/*",
			},

			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: []*framework.Path{
			pathConfig(&b),
			pathGroups(&b),
			pathGroupsList(&b),
			pathUsers(&b),
			pathUsersList(&b),
			pathLogin(&b),
			pathConfigRotateRoot(&b),
		},

		AuthRenew:        b.pathLoginRenew,
		BackendType:      logical.TypeCredential,
		RotateCredential: b.rotateRootCredential,
	}

	return &b
}

type backend struct {
	*framework.Backend

	mu sync.RWMutex
}

func (b *backend) maybeLogDebug(msg string, args ...interface{}) {
	if b.Logger().IsDebug() {
		b.Logger().Debug(msg, args...)
	}
}

// Login authenticates a user against the LDAP server using the provided username and password.
// It returns the effective username, policies associated with the user, a response containing
func (b *backend) Login(ctx context.Context, req *logical.Request, username string, password string, usernameAsAlias bool) (string, []string, *logical.Response, []string, error) {
	cfg, err := b.Config(ctx, req)
	if err != nil {
		return "", nil, nil, nil, err
	}
	if cfg == nil {
		return "", nil, logical.ErrorResponse("ldap backend not configured"), nil, nil
	}

	if cfg.DenyNullBind && len(password) == 0 {
		return "", nil, logical.ErrorResponse("password cannot be of zero length when passwordless binds are being denied"), nil, nil
	}

	ldapClient, err := ldap.NewClient(ctx, ldaputil.ConvertConfig(cfg.ConfigEntry))
	if err != nil {
		b.maybeLogDebug("error creating client", "error", err)
		return "", nil, logical.ErrorResponse(err.Error()), nil, nil
	}

	// Clean connection
	defer ldapClient.Close(ctx)

	c, err := ldapClient.Authenticate(ctx, username, password, ldap.WithGroups(), ldap.WithUserAttributes())
	if err != nil {
		if strings.Contains(err.Error(), "discovery of user bind DN failed") ||
			strings.Contains(err.Error(), "unable to bind user") {
			b.maybeLogDebug("error getting user bind DN", username, "error", err)
			return "", nil, logical.ErrorResponse(errUserBindFailed), nil, logical.ErrInvalidCredentials
		}

		return "", nil, logical.ErrorResponse(err.Error()), nil, nil
	}

	b.maybeLogDebug("user binddn fetched", "username", username, "binddn", c.UserDN)
	if c.UserDN == "" {
		return "", nil, logical.ErrorResponse("user bind DN is empty after authentication"), nil, nil
	}

	ldapGroups := c.Groups
	ldapResponse := &logical.Response{
		Data: map[string]interface{}{},
	}
	if len(ldapGroups) == 0 {
		errString := fmt.Sprintf(
			"no LDAP groups found in groupDN %q; only policies from locally-defined groups available",
			cfg.GroupDN)

		b.maybeLogDebug(errString)
	}

	for _, warning := range c.Warnings {
		b.maybeLogDebug(string(warning))
	}

	canonicalUsername := username
	if usernameAsAlias {
		// expect to get the username from UserDN in the case where we are setting the
		// entity alias to be the username.
		parsed, err := goldap.ParseDN(c.UserDN)
		if err != nil {
			b.Logger().Error("Invalid DN after authentication", "user_dn", c.UserDN, "error", err)
			return "", nil, logical.ErrorResponse("invalid DN after authentication"), nil, nil
		}

		b.maybeLogDebug("User DN parsed", "parsed", parsed, "username", username)
		var found bool
		for _, rdn := range parsed.RDNs {
			if found {
				b.maybeLogDebug("Found canonical username", "canonicalUsername", canonicalUsername, "cfg.userAttr", cfg.UserAttr)
				break
			}
			for _, attr := range rdn.Attributes {
				b.maybeLogDebug("Handling RDN", "attr.Type", attr.Type, "attr.Value", attr.Value, "cfg.UserAttr", cfg.UserAttr)
				if strings.EqualFold(attr.Type, cfg.UserAttr) {
					b.maybeLogDebug("Found user attribute in RDN", "attr.Type", attr.Type, "attr.Value", attr.Value)
					canonicalUsername = attr.Value
					found = true
					break
				}
			}
		}
	}

	var allGroups []string

	cs := cfg.CaseSensitiveNames != nil && *cfg.CaseSensitiveNames
	if !cs {
		canonicalUsername = strings.ToLower(canonicalUsername)
	}

	// Import the custom added groups from ldap backend
	user, err := b.User(ctx, req.Storage, canonicalUsername)
	if err == nil && user != nil && user.Groups != nil {
		allGroups = append(allGroups, user.Groups...)
	}
	// Merge local and LDAP groups
	allGroups = append(allGroups, ldapGroups...)

	canonicalGroups := allGroups
	// If not case-sensitive, lowercase all
	if !cs {
		canonicalGroups = make([]string, len(allGroups))
		for i, v := range allGroups {
			canonicalGroups[i] = strings.ToLower(v)
		}
	}

	// Retrieve policies
	var policies []string
	for _, groupName := range canonicalGroups {
		group, err := b.Group(ctx, req.Storage, groupName)
		if err == nil && group != nil {
			policies = append(policies, group.Policies...)
		}
	}
	if user != nil && user.Policies != nil {
		policies = append(policies, user.Policies...)
	}
	// Policies from each group may overlap
	policies = strutil.RemoveDuplicates(policies, true)

	if usernameAsAlias {
		b.maybeLogDebug("UsernameAlias is set", "canonicalUsername", canonicalUsername, "policies", policies, "groups", allGroups)
		return canonicalUsername, policies, ldapResponse, allGroups, nil
	}

	userAttrValues := c.UserAttributes[cfg.UserAttr]
	if len(userAttrValues) == 0 {
		b.Logger().Error("missing entity alias attribute value")
		return "", nil, logical.ErrorResponse("missing entity alias attribute value"), nil, nil
	}

	effectiveUsername := userAttrValues[0]
	if effectiveUsername == "" {
		b.Logger().Error("empty entity alias attribute value")
		return "", nil, logical.ErrorResponse("empty entity alias attribute value"), nil, nil
	}

	return effectiveUsername, policies, ldapResponse, allGroups, nil
}

const backendHelp = `
The "ldap" credential provider allows authentication querying
a LDAP server, checking username and password, and associating groups
to set of policies.

Configuration of the server is done through the "config" and "groups"
endpoints by a user with root access. Authentication is then done
by supplying the two fields for "login".
`
