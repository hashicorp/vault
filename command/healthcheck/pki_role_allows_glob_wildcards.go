package healthcheck

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type RoleAllowsGlobWildcards struct {
	Enabled            bool
	UnsupportedVersion bool

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
	exit, _, roles, err := pkiFetchRoles(e, func() {
		h.UnsupportedVersion = true
	})
	if exit {
		return err
	}

	for _, role := range roles {
		skip, _, entry, err := pkiFetchRole(e, role, func() {
			h.UnsupportedVersion = true
		})
		if skip || entry == nil {
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

	for role, entry := range h.RoleEntryMap {
		allowsWildcard, present := entry["allow_wildcard_certificates"].(bool)
		if !present || !allowsWildcard {
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
			Message:  fmt.Sprintf("Role currently allows wildcard issuance while allowing globs in allowed_domains (%v). Because globs can expand to one or more wildcard character, including wildcards under additional subdomains, these options are dangerous to enable together. If glob domains are required to be enabled, it is suggested to either disable wildcard issuance if not desired, or create two separate roles -- one with wildcard issuanced for specified domains, and one with glob matching enabled for concrete domain identifiers.", allowedDomains),
		}

		results = append(results, &ret)
	}

	return
}
