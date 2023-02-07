// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcheck

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type RoleAllowsGlobWildcards struct {
	Enabled            bool
	UnsupportedVersion bool
	NoPerms            bool

	RoleEntryMap map[string]map[string]interface{}
}

func NewRoleAllowsGlobWildcardsCheck() Check {
	return &RoleAllowsGlobWildcards{
		RoleEntryMap: make(map[string]map[string]interface{}),
	}
}

func (h *RoleAllowsGlobWildcards) Name() string {
	return "role_allows_glob_wildcards"
}

func (h *RoleAllowsGlobWildcards) IsEnabled() bool {
	return h.Enabled
}

func (h *RoleAllowsGlobWildcards) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (h *RoleAllowsGlobWildcards) LoadConfig(config map[string]interface{}) error {
	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *RoleAllowsGlobWildcards) FetchResources(e *Executor) error {
	exit, f, roles, err := pkiFetchRolesList(e, func() {
		h.UnsupportedVersion = true
	})
	if exit || err != nil {
		if f != nil && f.IsSecretPermissionsError() {
			h.NoPerms = true
		}
		return err
	}

	for _, role := range roles {
		skip, f, entry, err := pkiFetchRole(e, role, func() {
			h.UnsupportedVersion = true
		})
		if skip || err != nil || entry == nil {
			if f != nil && f.IsSecretPermissionsError() {
				h.NoPerms = true
			}
			if err != nil {
				return err
			}
			continue
		}

		h.RoleEntryMap[role] = entry
	}

	return nil
}

func (h *RoleAllowsGlobWildcards) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
		// Shouldn't happen; roles have been around forever.
		ret := Result{
			Status:   ResultInvalidVersion,
			Endpoint: "/{{mount}}/roles",
			Message:  "This health check requires Vault 1.11+ but an earlier version of Vault Server was contacted, preventing this health check from running.",
		}
		return []*Result{&ret}, nil
	}
	if h.NoPerms {
		ret := Result{
			Status:   ResultInsufficientPermissions,
			Endpoint: "/{{mount}}/roles",
			Message:  "lacks permission either to list the roles or to read a specific role. This may restrict the ability to fully execute this health check.",
		}
		if e.Client.Token() == "" {
			ret.Message = "No token available and so this health check " + ret.Message
		} else {
			ret.Message = "This token " + ret.Message
		}
		results = append(results, &ret)
	}

	for role, entry := range h.RoleEntryMap {
		allowsWildcards, present := entry["allow_wildcard_certificates"]
		if !present {
			ret := Result{
				Status:   ResultInvalidVersion,
				Endpoint: "/{{mount}}/roles",
				Message:  "This health check requires a version of Vault with allow_wildcard_certificates (Vault 1.8.9+, 1.9.4+, or 1.10.0+), but an earlier version of Vault Server was contacted, preventing this health check from running.",
			}
			return []*Result{&ret}, nil
		}
		if !allowsWildcards.(bool) {
			continue
		}

		allowsGlobs := entry["allow_glob_domains"].(bool)
		if !allowsGlobs {
			continue
		}

		rawAllowedDomains := entry["allowed_domains"].([]interface{})
		var allowedDomains []string
		for _, rawDomain := range rawAllowedDomains {
			allowedDomains = append(allowedDomains, rawDomain.(string))
		}

		if len(allowedDomains) == 0 {
			continue
		}

		hasGlobbedDomain := false
		for _, domain := range allowedDomains {
			if strings.Contains(domain, "*") {
				hasGlobbedDomain = true
				break
			}
		}

		if !hasGlobbedDomain {
			continue
		}

		ret := Result{
			Status:   ResultWarning,
			Endpoint: "/{{mount}}/role/" + role,
			Message:  fmt.Sprintf("Role currently allows wildcard issuance while allowing globs in allowed_domains (%v). Because globs can expand to one or more wildcard character, including wildcards under additional subdomains, these options are dangerous to enable together. If glob domains are required to be enabled, it is suggested to either disable wildcard issuance if not desired, or create two separate roles -- one with wildcard issuance for specified domains and one with glob matching enabled for concrete domain identifiers.", allowedDomains),
		}

		results = append(results, &ret)
	}

	if len(results) == 0 && len(h.RoleEntryMap) > 0 {
		ret := Result{
			Status:   ResultOK,
			Endpoint: "/{{mount}}/roles",
			Message:  "Roles follow best practices regarding restricting wildcard certificate issuance in roles.",
		}

		results = append(results, &ret)
	}

	return
}
