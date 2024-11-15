// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kerberos

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/jcmturner/gokrb5/v8/keytab"
	"github.com/jcmturner/gokrb5/v8/service"
	"github.com/jcmturner/gokrb5/v8/spnego"
	"gopkg.in/jcmturner/goidentity.v3"
)

func (b *backend) pathLogin() *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixKerberos,
		},
		Fields: map[string]*framework.FieldSchema{
			"authorization": {
				Type:        framework.TypeString,
				Description: `SPNEGO Authorization header. Required.`,
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathLoginGet,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "login2",
				},
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathLoginUpdate,
				DisplayAttrs: &framework.DisplayAttributes{
					OperationVerb: "login",
				},
			},
		},
	}
}

func parseKeytab(b64EncodedKt string) (*keytab.Keytab, error) {
	decodedKt, err := base64.StdEncoding.DecodeString(b64EncodedKt)
	if err != nil {
		return nil, err
	}
	parsedKt := new(keytab.Keytab)
	if err := parsedKt.Unmarshal(decodedKt); err != nil {
		return nil, err
	}
	return parsedKt, nil
}

func (b *backend) pathLoginGet(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	return &logical.Response{
		Auth: &logical.Auth{},
		Headers: map[string][]string{
			"www-authenticate": {"Negotiate"},
		},
	}, logical.CodedError(401, "authentication required")
}

func (b *backend) pathLoginUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	kerbCfg, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, errwrap.Wrapf("unable to get kerberos config: {{err}}", err)
	}
	if kerbCfg == nil {
		return nil, errors.New("backend kerberos not configured")
	}

	ldapCfg, err := b.ConfigLdap(ctx, req)
	if err != nil {
		return nil, errwrap.Wrapf("unable to get ldap config: {{err}}", err)
	}
	if ldapCfg == nil {
		return nil, errors.New("ldap backend not configured")
	}

	// Check for a CIDR match.
	if len(ldapCfg.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, ldapCfg.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}
	}

	authorizationString := ""
	authorizationHeaders := req.Headers["Authorization"]
	if len(authorizationHeaders) > 0 {
		authorizationString = authorizationHeaders[0]
	} else {
		authorizationString = d.Get("authorization").(string)
	}

	kt, err := parseKeytab(kerbCfg.Keytab)
	if err != nil {
		return nil, errwrap.Wrapf("could not parse keytab: {{err}}", err)
	}

	if kerbCfg.RemoveInstanceName {
		removeInstanceNameFromKeytab(kt)
	}

	s := strings.SplitN(authorizationString, " ", 2)
	if len(s) != 2 || s[0] != "Negotiate" {
		return b.pathLoginGet(ctx, req, d)
	}

	// The SPNEGOKRB5Authenticate method only calls an inner function if it's
	// successful. Let's use it to record success, and to retrieve the caller's
	// identity.
	username := ""
	authenticated := false
	var identity goidentity.Identity
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Context().Value(goidentity.CTXKey)
		if raw == nil {
			w.WriteHeader(400)
			_, _ = w.Write([]byte("identity credentials are not included"))
			return
		}
		ok := false
		identity, ok = raw.(goidentity.Identity)
		if !ok {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(fmt.Sprintf("identity credentials are malformed: %+v", raw)))
			return
		}
		b.Logger().Debug(fmt.Sprintf("identity: %+v", identity))
		username = identity.UserName()

		if kerbCfg.RemoveInstanceName {
			user := splitUsername(identity.UserName())
			if len(user) > 1 {
				username = user[0]
			}
		}

		// Verify that the realm on the LDAP config (if set) is the same as the identity's
		// realm. The UPNDomain denotes the realm on the LDAP config, and the identity
		// domain likewise identifies the realm. This is a case sensitive check.
		// This covers an edge case where, potentially, there has been drift between the LDAP
		// config's realm and the Kerberos realm. In such a case, it prevents a user from
		// passing Kerberos authentication, and then extracting group membership, and
		// therefore policies, from a separate directory.
		if ldapCfg.ConfigEntry.UPNDomain != "" && identity.Domain() != ldapCfg.ConfigEntry.UPNDomain {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(fmt.Sprintf("identity domain of %q doesn't match LDAP upndomain of %q", identity.Domain(), ldapCfg.ConfigEntry.UPNDomain)))
			return
		}
		authenticated = true
	})

	// Let's pass in a logger so we can get debugging information if anything
	// goes wrong.
	l := b.Logger().StandardLogger(&hclog.StandardLoggerOptions{
		InferLevels: true,
	})

	// Now let's use our inner handler to compose the overall function.
	authHTTPHandler := spnego.SPNEGOKRB5Authenticate(inner, kt, service.Logger(l), service.KeytabPrincipal(kerbCfg.ServiceAccount))

	// Because the outer application strips off the raw request, we need to
	// re-compose it to use this authentication handler. Only the request
	// remote addr and headers are used anyways. We use an arbitrary port
	// of 8080 because it's not used for anything but logging, but is required
	// by an underlying parser.
	rebuiltReq := &http.Request{
		Header:     req.Headers,
		RemoteAddr: req.Connection.RemoteAddr + ":8080",
	}

	// Finally, execute the SPNEGO authentication check.
	w := &simpleResponseWriter{}
	authHTTPHandler.ServeHTTP(w, rebuiltReq)
	if !authenticated {
		resp := &logical.Response{
			Warnings: []string{string(w.body)},
		}
		return logical.RespondWithStatusCode(resp, req, w.statusCode)
	}

	// Now that they've passed the Kerb authentication, begin checking if
	// they're a member of an LDAP group that should have additional policies
	// attached.
	ldapClient := ldaputil.Client{
		Logger: b.Logger(),
		LDAP:   ldaputil.NewLDAP(),
	}

	ldapConnection, err := ldapClient.DialLDAP(ldapCfg.ConfigEntry)
	if err != nil {
		return nil, errwrap.Wrapf("could not connect to LDAP: {{err}}", err)
	}
	if ldapConnection == nil {
		return nil, errors.New("invalid connection returned from LDAP dial")
	}

	// Clean ldap connection
	defer ldapConnection.Close()

	if len(ldapCfg.BindPassword) > 0 {
		err = ldapConnection.Bind(ldapCfg.BindDN, ldapCfg.BindPassword)
	} else {
		err = ldapConnection.UnauthenticatedBind(ldapCfg.BindDN)
	}
	if err != nil {
		return nil, fmt.Errorf("LDAP bind failed: %v", err)
	}

	userBindDN, err := ldapClient.GetUserBindDN(ldapCfg.ConfigEntry, ldapConnection, username)
	if err != nil {
		return nil, errwrap.Wrapf("unable to get user binddn: {{err}}", err)
	}
	b.Logger().Debug("auth/ldap: User BindDN fetched", "username", identity.UserName(), "binddn", userBindDN)

	userDN, err := ldapClient.GetUserDN(ldapCfg.ConfigEntry, ldapConnection, userBindDN, username)
	if err != nil {
		return nil, errwrap.Wrapf("unable to get user dn: {{err}}", err)
	}

	ldapGroups, err := ldapClient.GetLdapGroups(ldapCfg.ConfigEntry, ldapConnection, userDN, username)
	if err != nil {
		return nil, errwrap.Wrapf("unable to get ldap groups: {{err}}", err)
	}
	b.Logger().Debug("auth/ldap: Groups fetched from server", "num_server_groups", len(ldapGroups), "server_groups", ldapGroups)

	var allGroups []string
	// Merge local and LDAP groups
	allGroups = append(allGroups, ldapGroups...)

	// Retrieve policies
	var policies []string
	for _, groupName := range allGroups {
		group, err := b.Group(ctx, req.Storage, groupName)
		if err != nil {
			b.Logger().Debug(fmt.Sprintf("unable to retrieve %s: %s", groupName, err.Error()))
			continue
		}
		if group == nil {
			b.Logger().Debug(fmt.Sprintf("unable to find %s, does not currently exist", groupName))
			continue
		}
		policies = append(policies, group.Policies...)
	}

	// Policies from each group may overlap
	policies = strutil.RemoveDuplicates(policies, true)
	auth := &logical.Auth{
		InternalData: map[string]interface{}{},
		Metadata: map[string]string{
			"user":   identity.UserName(),
			"domain": identity.Domain(),
		},
		DisplayName: identity.UserName(),
		Alias:       &logical.Alias{Name: identity.UserName()},
	}

	ldapCfg.PopulateTokenAuth(auth)

	// This is done after PopulateTokenAuth because it forces Renewable to be true.
	// Renewable was always false at the time of the code's introduction, and we would
	// like to keep it the same until we have a concrete reason to change its behavior.
	auth.LeaseOptions = logical.LeaseOptions{
		Renewable: false,
	}

	// Combine our policies with the ones parsed from PopulateTokenAuth.
	if len(policies) > 0 {
		auth.Policies = append(auth.Policies, policies...)
	}

	// Add the LDAP groups so the Identity system can use them
	if kerbCfg.AddGroupAliases {
		for _, groupName := range allGroups {
			if groupName == "" {
				continue
			}
			auth.GroupAliases = append(auth.GroupAliases, &logical.Alias{
				Name: groupName,
			})
		}
	}

	return &logical.Response{
		Auth: auth,
	}, nil
}

type simpleResponseWriter struct {
	body       []byte
	statusCode int
}

func (w *simpleResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (w *simpleResponseWriter) Write(b []byte) (int, error) {
	w.body = b
	return 0, nil
}

func (w *simpleResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}
