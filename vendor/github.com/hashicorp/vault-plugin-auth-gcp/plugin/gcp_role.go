// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
)

type gcpRole struct {
	tokenutil.TokenParams

	// RoleID is a unique identifier for this role.
	RoleID string `json:"role_id"`

	// Type of this role. See path_role constants for currently supported types.
	RoleType string `json:"role_type,omitempty"`

	// Policies for Vault to assign to authorized entities.
	Policies []string `json:"policies,omitempty"`

	// TTL of Vault auth leases under this role.
	TTL time.Duration `json:"ttl,omitempty"`

	// Max total TTL including renewals, of Vault auth leases under this role.
	MaxTTL time.Duration `json:"max_ttl,omitempty"`

	// Period, If set, indicates that this token should not expire and
	// should be automatically renewed within this time period
	// with TTL equal to this value.
	Period time.Duration `json:"period,omitempty"`

	// Projects that entities must belong to
	BoundProjects []string `json:"bound_projects,omitempty"`

	// Service accounts allowed to login under this role.
	BoundServiceAccounts []string `json:"bound_service_accounts,omitempty"`

	// AddGroupAliases adds Vault group aliases to the response.
	AddGroupAliases bool `json:"add_group_aliases,omitempty"`

	// --| IAM-only attributes |--
	// MaxJwtExp is the duration from time of authentication that a JWT used to authenticate to role must expire within.
	// TODO(emilymye): Allow this to be updated for GCE roles once 'exp' parameter has been allowed for GCE metadata.
	MaxJwtExp time.Duration `json:"max_jwt_exp,omitempty"`

	// AllowGCEInference, if false, does not allow a GCE instance to login under this 'iam' role. If true (default),
	// a service account is inferred from the instance metadata and used as the authenticating instance.
	AllowGCEInference bool `json:"allow_gce_inference,omitempty"`

	// --| GCE-only attributes |--
	// BoundRegions that instances must belong to in order to login under this role.
	BoundRegions []string `json:"bound_regions,omitempty"`

	// BoundZones that instances must belong to in order to login under this role.
	BoundZones []string `json:"bound_zones,omitempty"`

	// BoundInstanceGroups are the instance group that instances must belong to in order to login under this role.
	BoundInstanceGroups []string `json:"bound_instance_groups,omitempty"`

	// BoundLabels that instances must currently have set in order to login under this role.
	BoundLabels map[string]string `json:"bound_labels,omitempty"`

	// Deprecated fields
	// TODO: Remove in 0.5.0+
	ProjectId          string `json:"project_id,omitempty"`
	BoundRegion        string `json:"bound_region,omitempty"`
	BoundZone          string `json:"bound_zone,omitempty"`
	BoundInstanceGroup string `json:"bound_instance_group,omitempty"`
}

// updateRole updates the given role with values parsed/validated from given FieldData.
// Exactly one of the response and error will be nil. The response is only used to pass back warnings.
// This method does not validate the role. Validation is done before storage.
func (role *gcpRole) updateRole(sys logical.SystemView, req *logical.Request, data *framework.FieldData) (warnings []string, err error) {
	if e := role.ParseTokenFields(req, data); e != nil {
		return nil, e
	}

	// Handle token field upgrades
	{
		if e := tokenutil.UpgradeValue(data, "policies", "token_policies", &role.Policies, &role.TokenPolicies); e != nil {
			return nil, e
		}

		if e := tokenutil.UpgradeValue(data, "ttl", "token_ttl", &role.TTL, &role.TokenTTL); e != nil {
			return nil, e
		}

		if e := tokenutil.UpgradeValue(data, "max_ttl", "token_max_ttl", &role.MaxTTL, &role.TokenMaxTTL); e != nil {
			return nil, e
		}

		if e := tokenutil.UpgradeValue(data, "period", "token_period", &role.Period, &role.TokenPeriod); e != nil {
			return nil, e
		}
	}

	// Set role type
	if rt, ok := data.GetOk("type"); ok {
		roleType := rt.(string)
		if role.RoleType != roleType && req.Operation == logical.UpdateOperation {
			return nil, fmt.Errorf("role type cannot be changed for an existing role")
		}
		role.RoleType = roleType
	} else if req.Operation == logical.CreateOperation {
		return nil, fmt.Errorf(errEmptyRoleType)
	}

	def := sys.DefaultLeaseTTL()
	if role.TokenTTL > def {
		warnings = append(warnings, fmt.Sprintf(`Given token ttl of %q is greater `+
			`than the maximum system/mount TTL of %q. The TTL will be capped at `+
			`%q during login.`, role.TokenTTL, def, def))
	}

	// Update token Max TTL.
	def = sys.MaxLeaseTTL()
	if role.TokenMaxTTL > def {
		warnings = append(warnings, fmt.Sprintf(`Given token max ttl of %q is greater `+
			`than the maximum system/mount MaxTTL of %q. The MaxTTL will be `+
			`capped at %q during login.`, role.TokenMaxTTL, def, def))
	}
	if role.TokenPeriod > def {
		warnings = append(warnings, fmt.Sprintf(`Given token period of %q is greater `+
			`than the maximum system/mount period of %q. The period will be `+
			`capped at %q during login.`, role.TokenPeriod, def, def))
	}

	// Update bound GCP service accounts.
	if sa, ok := data.GetOk("bound_service_accounts"); ok {
		role.BoundServiceAccounts = sa.([]string)
	} else {
		// Check for older version of param name
		if sa, ok := data.GetOk("service_accounts"); ok {
			warnings = append(warnings, `The "service_accounts" field is deprecated. `+
				`Please use "bound_service_accounts" instead. The "service_accounts" `+
				`field will be removed in a later release, so please update accordingly.`)
			role.BoundServiceAccounts = sa.([]string)
		}
	}
	if len(role.BoundServiceAccounts) > 0 {
		role.BoundServiceAccounts = strutil.TrimStrings(role.BoundServiceAccounts)
		role.BoundServiceAccounts = strutil.RemoveDuplicates(role.BoundServiceAccounts, false)
	}

	// Update bound GCP projects.
	boundProjects, givenBoundProj := data.GetOk("bound_projects")
	if givenBoundProj {
		role.BoundProjects = boundProjects.([]string)
	}
	if projectId, ok := data.GetOk("project_id"); ok {
		if givenBoundProj {
			return warnings, errors.New("only one of 'bound_projects' or 'project_id' can be given")
		}
		warnings = append(warnings,
			`The "project_id" (singular) field is deprecated. `+
				`Please use plural "bound_projects" instead to bind required GCP projects. `+
				`The "project_id" field will be removed in a later release, so please update accordingly.`)
		role.BoundProjects = []string{projectId.(string)}
	}
	if len(role.BoundProjects) > 0 {
		role.BoundProjects = strutil.TrimStrings(role.BoundProjects)
		role.BoundProjects = strutil.RemoveDuplicates(role.BoundProjects, false)
	}

	// Update bound GCP projects.
	addGroupAliases, ok := data.GetOk("add_group_aliases")
	if ok {
		role.AddGroupAliases = addGroupAliases.(bool)
	}

	// Update fields specific to this type
	switch role.RoleType {
	case iamRoleType:
		if err = checkInvalidRoleTypeArgs(data, gceOnlyFieldSchema); err != nil {
			return warnings, err
		}
		if warnings, err = role.updateIamFields(data, req.Operation); err != nil {
			return warnings, err
		}
	case gceRoleType:
		if err = checkInvalidRoleTypeArgs(data, iamOnlyFieldSchema); err != nil {
			return warnings, err
		}
		if warnings, err = role.updateGceFields(data, req.Operation); err != nil {
			return warnings, err
		}
	}

	return warnings, nil
}

func (role *gcpRole) validate(sys logical.SystemView) (warnings []string, err error) {
	warnings = []string{}

	switch role.RoleType {
	case iamRoleType:
		if warnings, err = role.validateForIAM(); err != nil {
			return warnings, err
		}
	case gceRoleType:
		if warnings, err = role.validateForGCE(); err != nil {
			return warnings, err
		}
	case "":
		return warnings, errors.New(errEmptyRoleType)
	default:
		return warnings, fmt.Errorf("role type '%s' is invalid", role.RoleType)
	}

	defaultLeaseTTL := sys.DefaultLeaseTTL()
	if role.TokenTTL > defaultLeaseTTL {
		warnings = append(warnings, fmt.Sprintf(
			"Given ttl of %d seconds greater than current mount/system default of %d seconds; ttl will be capped at login time",
			role.TokenTTL/time.Second, defaultLeaseTTL/time.Second))
	}

	defaultMaxTTL := sys.MaxLeaseTTL()
	if role.TokenMaxTTL > defaultMaxTTL {
		warnings = append(warnings, fmt.Sprintf(
			"Given max_ttl of %d seconds greater than current mount/system default of %d seconds; max_ttl will be capped at login time",
			role.TokenMaxTTL/time.Second, defaultMaxTTL/time.Second))
	}
	if role.TokenMaxTTL < time.Duration(0) {
		return warnings, errors.New("max_ttl cannot be negative")
	}
	if role.TokenMaxTTL != 0 && role.TokenMaxTTL < role.TokenTTL {
		return warnings, errors.New("ttl should be shorter than max_ttl")
	}

	if role.TokenPeriod > sys.MaxLeaseTTL() {
		return warnings, fmt.Errorf("'period' of '%s' is greater than the backend's maximum lease TTL of '%s'", role.TokenPeriod.String(), sys.MaxLeaseTTL().String())
	}

	return warnings, nil
}

// updateIamFields updates IAM-only fields for a role.
func (role *gcpRole) updateIamFields(data *framework.FieldData, op logical.Operation) (warnings []string, err error) {
	if allowGCEInference, ok := data.GetOk("allow_gce_inference"); ok {
		role.AllowGCEInference = allowGCEInference.(bool)
	} else if op == logical.CreateOperation {
		role.AllowGCEInference = data.Get("allow_gce_inference").(bool)
	}

	if maxJwtExp, ok := data.GetOk("max_jwt_exp"); ok {
		role.MaxJwtExp = time.Duration(maxJwtExp.(int)) * time.Second
	} else if op == logical.CreateOperation {
		role.MaxJwtExp = time.Duration(defaultIamMaxJwtExpMinutes) * time.Minute
	}

	return warnings, nil
}

// updateGceFields updates GCE-only fields for a role.
func (role *gcpRole) updateGceFields(data *framework.FieldData, op logical.Operation) (warnings []string, err error) {
	if regions, ok := data.GetOk("bound_regions"); ok {
		role.BoundRegions = regions.([]string)
	} else if op == logical.CreateOperation {
		role.BoundRegions = data.Get("bound_regions").([]string)
	}

	if zones, ok := data.GetOk("bound_zones"); ok {
		role.BoundZones = zones.([]string)
	} else if op == logical.CreateOperation {
		role.BoundZones = data.Get("bound_zones").([]string)
	}

	if instanceGroups, ok := data.GetOk("bound_instance_groups"); ok {
		role.BoundInstanceGroups = instanceGroups.([]string)
	} else if op == logical.CreateOperation {
		role.BoundInstanceGroups = data.Get("bound_instance_groups").([]string)
	}

	if boundRegion, ok := data.GetOk("bound_region"); ok {
		if _, ok := data.GetOk("bound_regions"); ok {
			return warnings, fmt.Errorf(`cannot specify both "bound_region" and "bound_regions"`)
		}

		warnings = append(warnings, `The "bound_region" field is deprecated. `+
			`Please use "bound_regions" (plural) instead. You can still specify a `+
			`single region, but multiple regions are also now supported. The `+
			`"bound_region" field will be removed in a later release, so please `+
			`update accordingly.`)
		role.BoundRegions = append(role.BoundRegions, boundRegion.(string))
	}

	if boundZone, ok := data.GetOk("bound_zone"); ok {
		if _, ok := data.GetOk("bound_zones"); ok {
			return warnings, fmt.Errorf(`cannot specify both "bound_zone" and "bound_zones"`)
		}

		warnings = append(warnings, `The "bound_zone" field is deprecated. `+
			`Please use "bound_zones" (plural) instead. You can still specify a `+
			`single zone, but multiple zones are also now supported. The `+
			`"bound_zone" field will be removed in a later release, so please `+
			`update accordingly.`)
		role.BoundZones = append(role.BoundZones, boundZone.(string))
	}

	if boundInstanceGroup, ok := data.GetOk("bound_instance_group"); ok {
		if _, ok := data.GetOk("bound_instance_groups"); ok {
			return warnings, fmt.Errorf(`cannot specify both "bound_instance_group" and "bound_instance_groups"`)
		}

		warnings = append(warnings, `The "bound_instance_group" field is deprecated. `+
			`Please use "bound_instance_groups" (plural) instead. You can still specify a `+
			`single instance group, but multiple instance groups are also now supported. The `+
			`"bound_instance_group" field will be removed in a later release, so please `+
			`update accordingly.`)
		role.BoundInstanceGroups = append(role.BoundInstanceGroups, boundInstanceGroup.(string))
	}

	if labelsRaw, ok := data.GetOk("bound_labels"); ok {
		labels, invalidLabels := gcputil.ParseGcpLabels(labelsRaw.([]string))
		if len(invalidLabels) > 0 {
			return warnings, fmt.Errorf("invalid labels given: %q", invalidLabels)
		}
		role.BoundLabels = labels
	}

	if len(role.Policies) > 0 {
		role.Policies = strutil.TrimStrings(role.Policies)
		role.Policies = strutil.RemoveDuplicates(role.Policies, false)
	}

	if len(role.BoundRegions) > 0 {
		role.BoundRegions = strutil.TrimStrings(role.BoundRegions)
		role.BoundRegions = strutil.RemoveDuplicates(role.BoundRegions, false)
	}

	if len(role.BoundZones) > 0 {
		role.BoundZones = strutil.TrimStrings(role.BoundZones)
		role.BoundZones = strutil.RemoveDuplicates(role.BoundZones, false)
	}

	if len(role.BoundInstanceGroups) > 0 {
		role.BoundInstanceGroups = strutil.TrimStrings(role.BoundInstanceGroups)
		role.BoundInstanceGroups = strutil.RemoveDuplicates(role.BoundInstanceGroups, false)
	}

	return warnings, nil
}

// validateIamFields validates the IAM-only fields for a role.
func (role *gcpRole) validateForIAM() (warnings []string, err error) {
	if len(role.BoundServiceAccounts) == 0 {
		return []string{}, errors.New(errEmptyIamServiceAccounts)
	}

	if len(role.BoundServiceAccounts) > 1 && strutil.StrListContains(role.BoundServiceAccounts, serviceAccountsWildcard) {
		return []string{}, fmt.Errorf("cannot provide IAM service account wildcard '%s' (for all service accounts) with other service accounts", serviceAccountsWildcard)
	}

	maxMaxJwtExp := time.Duration(maxJwtExpMaxMinutes) * time.Minute
	if role.MaxJwtExp > maxMaxJwtExp {
		return warnings, fmt.Errorf("max_jwt_exp cannot be more than %d minutes", maxJwtExpMaxMinutes)
	}

	return []string{}, nil
}

// validateGceFields validates the GCE-only fields for a role.
func (role *gcpRole) validateForGCE() (warnings []string, err error) {
	warnings = []string{}

	hasRegion := len(role.BoundRegions) > 0
	hasZone := len(role.BoundZones) > 0
	hasRegionOrZone := hasRegion || hasZone

	hasInstanceGroup := len(role.BoundInstanceGroups) > 0

	if hasInstanceGroup && !hasRegionOrZone {
		return warnings, errors.New(`region or zone information must be specified if an instance group is given`)
	}

	if hasRegion && hasZone {
		warnings = append(warnings, `Given both "bound_regions" and "bound_zones" `+
			`fields for role type "gce", "bound_regions" will be ignored in favor `+
			`of the more specific "bound_zones" field. To fix this warning, update `+
			`the role to remove either the "bound_regions" or "bound_zones" field.`)
	}

	return warnings, nil
}
