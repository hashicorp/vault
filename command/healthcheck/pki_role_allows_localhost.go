package healthcheck

import (
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type RoleAllowsLocalhost struct {
	Enabled            bool
	UnsupportedVersion bool
	NoPerms            bool

	RoleEntryMap map[string]map[string]interface{}
}

func NewRoleAllowsLocalhostCheck() Check {
	return &RoleAllowsLocalhost{
		RoleEntryMap: make(map[string]map[string]interface{}),
	}
}

func (h *RoleAllowsLocalhost) Name() string {
	return "role_allows_localhost"
}

func (h *RoleAllowsLocalhost) IsEnabled() bool {
	return h.Enabled
}

func (h *RoleAllowsLocalhost) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{}
}

func (h *RoleAllowsLocalhost) LoadConfig(config map[string]interface{}) error {
	enabled, err := parseutil.ParseBool(config["enabled"])
	if err != nil {
		return fmt.Errorf("error parsing %v.enabled: %w", h.Name(), err)
	}
	h.Enabled = enabled

	return nil
}

func (h *RoleAllowsLocalhost) FetchResources(e *Executor) error {
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

func (h *RoleAllowsLocalhost) Evaluate(e *Executor) (results []*Result, err error) {
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
			Message:  "lacks permission either to list the roles or to read a specific role. This may restrict the ability to fully execute this health check",
		}
		if e.Client.Token() == "" {
			ret.Message = "No token available and so this health check " + ret.Message
		} else {
			ret.Message = "This token " + ret.Message
		}
		results = append(results, &ret)
	}

	for role, entry := range h.RoleEntryMap {
		allowsLocalhost := entry["allow_localhost"].(bool)
		if !allowsLocalhost {
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

		ret := Result{
			Status:   ResultWarning,
			Endpoint: "/{{mount}}/role/" + role,
			Message:  fmt.Sprintf("Role currently allows localhost issuance with a non-empty allowed_domains (%v): this role is intended for issuing other hostnames and the allow_localhost=true option may be overlooked by operators. If this role is intended to issue certificates valid for localhost, consider setting allow_localhost=false and explicitly adding localhost to the list of allowed domains.", allowedDomains),
		}

		results = append(results, &ret)
	}

	if len(results) == 0 && len(h.RoleEntryMap) > 0 {
		ret := Result{
			Status:   ResultOK,
			Endpoint: "/{{mount}}/roles",
			Message:  "Roles follow best practices regarding allowing issuance for localhost domains.",
		}

		results = append(results, &ret)
	}

	return
}
