// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package healthcheck

import (
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type RoleNoStoreFalse struct {
	Enabled            bool
	UnsupportedVersion bool

	AllowedRoles map[string]bool

	RoleListFetchIssue *PathFetch
	RoleFetchIssues    map[string]*PathFetch
	RoleEntryMap       map[string]map[string]interface{}
	CRLConfig          *PathFetch
}

func NewRoleNoStoreFalseCheck() Check {
	return &RoleNoStoreFalse{
		RoleFetchIssues: make(map[string]*PathFetch),
		AllowedRoles:    make(map[string]bool),
		RoleEntryMap:    make(map[string]map[string]interface{}),
	}
}

func (h *RoleNoStoreFalse) Name() string {
	return "role_no_store_false"
}

func (h *RoleNoStoreFalse) IsEnabled() bool {
	return h.Enabled
}

func (h *RoleNoStoreFalse) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"allowed_roles": []string{},
	}
}

func (h *RoleNoStoreFalse) LoadConfig(config map[string]interface{}) error {
	value, present := config["allowed_roles"].([]interface{})
	if present {
		for _, rawValue := range value {
			h.AllowedRoles[rawValue.(string)] = true
		}
	}

	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *RoleNoStoreFalse) FetchResources(e *Executor) error {
	exit, f, roles, err := pkiFetchRolesList(e, func() {
		h.UnsupportedVersion = true
	})
	if exit || err != nil {
		if f != nil && f.IsSecretPermissionsError() {
			h.RoleListFetchIssue = f
		}
		return err
	}

	for _, role := range roles {
		skip, f, entry, err := pkiFetchRole(e, role, func() {
			h.UnsupportedVersion = true
		})
		if skip || err != nil || entry == nil {
			if f != nil && f.IsSecretPermissionsError() {
				h.RoleFetchIssues[role] = f
			}
			if err != nil {
				return err
			}
			continue
		}

		h.RoleEntryMap[role] = entry
	}

	// Check if the issuer is fetched yet.
	configRet, err := e.FetchIfNotFetched(logical.ReadOperation, "/{{mount}}/config/crl")
	if err != nil {
		return err
	}

	h.CRLConfig = configRet

	return nil
}

func (h *RoleNoStoreFalse) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		// Shouldn't happen; roles have been around forever.
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/{{mount}}/roles",
			Message:  "This health check requires Vault 1.11+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}

	if h.RoleListFetchIssue != nil && h.RoleListFetchIssue.IsSecretPermissionsError() {
		ret := Result{
			Status:   ResultInsufficientPermissions,
			Endpoint: h.RoleListFetchIssue.Path,
			Message:  "lacks permission either to list the roles. This restricts the ability to fully execute this health check.",
		}
		if e.Client.Token() == "" {
			ret.Message = "No token available and so this health check " + ret.Message
		} else {
			ret.Message = "This token " + ret.Message
		}
		return []*Result{&ret}, nil
	}

	for role, fetchPath := range h.RoleFetchIssues {
		if fetchPath != nil && fetchPath.IsSecretPermissionsError() {
			delete(h.RoleEntryMap, role)
			ret := Result{
				Status:   ResultInsufficientPermissions,
				Endpoint: fetchPath.Path,
				Message:  "Without this information, this health check is unable to function.",
			}

			if e.Client.Token() == "" {
				ret.Message = "No token available so unable for the endpoint for this mount. " + ret.Message
			} else {
				ret.Message = "This token lacks permission the endpoint for this mount. " + ret.Message
			}

			results = append(results, &ret)
		}
	}

	crlAutoRebuild := false
	if h.CRLConfig != nil {
		if h.CRLConfig.IsSecretPermissionsError() {
			ret := Result{
				Status:   ResultInsufficientPermissions,
				Endpoint: "/{{mount}}/config/crl",
				Message:  "This prevents the health check from seeing if the CRL is set to auto_rebuild=true and lowering the severity of check results appropriately.",
			}

			if e.Client.Token() == "" {
				ret.Message = "No token available so unable read authenticated CRL configuration for this mount. " + ret.Message
			} else {
				ret.Message = "This token lacks so permission to read the CRL configuration for this mount. " + ret.Message
			}

			results = append(results, &ret)
		} else if h.CRLConfig.Secret != nil && h.CRLConfig.Secret.Data["auto_rebuild"] != nil {
			crlAutoRebuild = h.CRLConfig.Secret.Data["auto_rebuild"].(bool)
		}
	}

	for role, entry := range h.RoleEntryMap {
		noStore := entry["no_store"].(bool)
		if noStore {
			continue
		}

		ret := Result{
			Status:   ResultWarning,
			Endpoint: "/{{mount}}/roles/" + role,
			Message:  "Role currently stores every issued certificate (no_store=false). Too many issued and/or revoked certificates can exceed Vault's storage limits and make operations slow. It is encouraged to enable auto-rebuild of CRLs to prevent every revocation from creating a new CRL, and to limit the number of certificates issued under roles with no_store=false: use shorter lifetimes and/or BYOC revocation instead.",
		}

		if crlAutoRebuild {
			ret.Status = ResultInformational
			ret.Message = "Role currently stores every issued certificate (no_store=false). With auto-rebuild CRL enabled, less performance impact occur on CRL rebuilding, but note that too many issued and/or revoked certificates can exceed Vault's storage limits and make operations slow. It is suggested to limit the number of certificates issued under roles with no_store=false: use shorter lifetimes to avoid revocation and/or BYOC revocation instead."
		}

		results = append(results, &ret)
	}

	if len(results) == 0 && len(h.RoleEntryMap) > 0 {
		ret := Result{
			Status:   ResultOK,
			Endpoint: "/{{mount}}/roles",
			Message:  "Roles follow best practices regarding certificate storage.",
		}

		results = append(results, &ret)
	}

	return
}
