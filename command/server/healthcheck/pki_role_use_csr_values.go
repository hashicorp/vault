// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package healthcheck

import (
	"fmt"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type RoleUseCsrValues struct {
	Enabled            bool
	UnsupportedVersion bool

	AllowedRoles map[string]bool

	RoleListFetchIssue *PathFetch
	RoleFetchIssues    map[string]*PathFetch
	RoleEntryMap       map[string]map[string]interface{}
}

func NewRoleUseCsrValuesCheck() Check {
	return &RoleUseCsrValues{
		RoleFetchIssues: make(map[string]*PathFetch),
		AllowedRoles:    make(map[string]bool),
		RoleEntryMap:    make(map[string]map[string]interface{}),
	}
}

func (h *RoleUseCsrValues) Name() string {
	return "role_use_csr_values"
}

func (h *RoleUseCsrValues) IsEnabled() bool {
	return h.Enabled
}

func (h *RoleUseCsrValues) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"allowed_roles": []string{},
	}
}

func (h *RoleUseCsrValues) LoadConfig(config map[string]interface{}) error {
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

func (h *RoleUseCsrValues) FetchResources(e *Executor) error {
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

	return nil
}

func (h *RoleUseCsrValues) Evaluate(e *Executor) (results []*Result, err error) {
	if h.UnsupportedVersion {
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

	for role, entry := range h.RoleEntryMap {
		if h.AllowedRoles[role] {
			continue
		}

		useCsrValues, _ := entry["use_csr_values"].(bool)
		if !useCsrValues {
			continue
		}

		ret := Result{
			Status:   ResultWarning,
			Endpoint: "/{{mount}}/roles/" + role,
			Message: "Role has use_csr_values=true. This causes the entire CSR subject, SANs, and " +
				"extensions to be taken verbatim, bypassing role-level subject validation. " +
				"Sentinel policies referencing subject fields (common_name, organization, etc.) " +
				"will not be evaluated against the actual values in the issued certificate. " +
				"Consider disabling use_csr_values and using use_csr_common_name and use_csr_sans " +
				"instead, which retain role-level validation.",
		}
		results = append(results, &ret)
	}

	if len(results) == 0 && len(h.RoleEntryMap) > 0 {
		ret := Result{
			Status:   ResultOK,
			Endpoint: "/{{mount}}/roles",
			Message:  "No roles use use_csr_values=true. Role subject validation is enforced for all roles.",
		}
		results = append(results, &ret)
	}

	return
}
