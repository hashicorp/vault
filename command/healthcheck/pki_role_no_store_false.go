package healthcheck

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
)

type RoleNoStoreFalse struct {
	Enabled            bool
	UnsupportedVersion bool

	AllowedRoles map[string]bool

	CertCounts   int
	RoleEntryMap map[string]map[string]interface{}
	CRLConfig    *PathFetch
}

func NewRoleNoStoreFalseCheck() Check {
	return &RoleNoStoreFalse{
		AllowedRoles: make(map[string]bool),
		RoleEntryMap: make(map[string]map[string]interface{}),
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

	exit, _, leaves, err := pkiFetchLeaves(e, func() {
		h.UnsupportedVersion = true
	})
	if exit {
		return err
	}
	h.CertCounts = len(leaves)

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
			Endpoint: "/{{mount}}/role/" + role,
			Message:  "Role currently stores every issued certificate (no_store=false). Too many issued and/or revoked certificates can exceed Vault's storage limits and make operations slow. It is encouraged to enable auto-rebuild of CRLs to prevent every revocation from creating a new CRL, and to limit the number of certificates issued under roles with no_store=false: use shorter lifetimes and/or BYOC revocation instead.",
		}

		if crlAutoRebuild {
			ret.Status = ResultInformational
			ret.Message = "Role currently stores every issued certificate (no_store=false). Too many issued and/or revoked certificates can exceed Vault's storage limits and make operations slow. It is encouraged to enable auto-rebuild of CRLs to prevent every revocation from creating a new CRL, and to limit the number of certificates issued under roles with no_store=false: use shorter lifetimes and/or BYOC revocation instead."
		}

		results = append(results, &ret)
	}

	return
}
