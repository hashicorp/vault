// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/hclutil"
	"github.com/hashicorp/vault/sdk/helper/identitytpl"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
)

const (
	DenyCapability      = "deny"
	CreateCapability    = "create"
	ReadCapability      = "read"
	UpdateCapability    = "update"
	DeleteCapability    = "delete"
	ListCapability      = "list"
	SudoCapability      = "sudo"
	RootCapability      = "root"
	PatchCapability     = "patch"
	SubscribeCapability = "subscribe"

	// Backwards compatibility
	OldDenyPathPolicy  = "deny"
	OldReadPathPolicy  = "read"
	OldWritePathPolicy = "write"
	OldSudoPathPolicy  = "sudo"
)

const (
	DenyCapabilityInt uint32 = 1 << iota
	CreateCapabilityInt
	ReadCapabilityInt
	UpdateCapabilityInt
	DeleteCapabilityInt
	ListCapabilityInt
	SudoCapabilityInt
	PatchCapabilityInt
	SubscribeCapabilityInt
)

// Error constants for testing
const (
	// ControlledCapabilityPolicySubsetError is thrown when a control group's controlled capabilities
	// are not a subset of the policy's capabilities.
	ControlledCapabilityPolicySubsetError = "control group factor capabilities must be a subset of the policy's capabilities"
)

type PolicyType uint32

const (
	PolicyTypeACL PolicyType = iota
	PolicyTypeRGP
	PolicyTypeEGP

	// Triggers a lookup in the map to figure out if ACL or RGP
	PolicyTypeToken
)

func (p PolicyType) String() string {
	switch p {
	case PolicyTypeACL:
		return "acl"
	case PolicyTypeRGP:
		return "rgp"
	case PolicyTypeEGP:
		return "egp"
	}

	return ""
}

var cap2Int = map[string]uint32{
	DenyCapability:      DenyCapabilityInt,
	CreateCapability:    CreateCapabilityInt,
	ReadCapability:      ReadCapabilityInt,
	UpdateCapability:    UpdateCapabilityInt,
	DeleteCapability:    DeleteCapabilityInt,
	ListCapability:      ListCapabilityInt,
	SudoCapability:      SudoCapabilityInt,
	PatchCapability:     PatchCapabilityInt,
	SubscribeCapability: SubscribeCapabilityInt,
}

type egpPath struct {
	Path string `json:"path"`
	Glob bool   `json:"glob"`
}

// Policy is used to represent the policy specified by an ACL configuration.
type Policy struct {
	sentinelPolicy
	Name      string       `hcl:"name"`
	Paths     []*PathRules `hcl:"-"`
	Raw       string
	Type      PolicyType
	Templated bool
	namespace *namespace.Namespace
}

// ShallowClone returns a shallow clone of the policy. This should not be used
// if any of the reference-typed fields are going to be modified
func (p *Policy) ShallowClone() *Policy {
	return &Policy{
		sentinelPolicy: p.sentinelPolicy,
		Name:           p.Name,
		Paths:          p.Paths,
		Raw:            p.Raw,
		Type:           p.Type,
		Templated:      p.Templated,
		namespace:      p.namespace,
	}
}

// PathRules represents a policy for a path in the namespace.
type PathRules struct {
	Path                string
	Policy              string
	Permissions         *ACLPermissions
	IsPrefix            bool
	HasSegmentWildcards bool
	Capabilities        []string

	// These keys are used at the top level to make the HCL nicer; we store in
	// the ACLPermissions object though
	MinWrappingTTLHCL      interface{}              `hcl:"min_wrapping_ttl"`
	MaxWrappingTTLHCL      interface{}              `hcl:"max_wrapping_ttl"`
	AllowedParametersHCL   map[string][]interface{} `hcl:"allowed_parameters"`
	DeniedParametersHCL    map[string][]interface{} `hcl:"denied_parameters"`
	RequiredParametersHCL  []string                 `hcl:"required_parameters"`
	MFAMethodsHCL          []string                 `hcl:"mfa_methods"`
	ControlGroupHCL        *ControlGroupHCL         `hcl:"control_group"`
	SubscribeEventTypesHCL []string                 `hcl:"subscribe_event_types"`
}

type ControlGroupHCL struct {
	TTL     interface{}                    `hcl:"ttl"`
	Factors map[string]*ControlGroupFactor `hcl:"factor"`
}

type ControlGroup struct {
	TTL     time.Duration
	Factors []*ControlGroupFactor
}

func (c *ControlGroup) Clone() (*ControlGroup, error) {
	clonedControlGroup, err := copystructure.Copy(c)
	if err != nil {
		return nil, err
	}

	cg := clonedControlGroup.(*ControlGroup)

	return cg, nil
}

type ControlGroupFactor struct {
	Name                   string
	Identity               *IdentityFactor `hcl:"identity"`
	ControlledCapabilities []string        `hcl:"controlled_capabilities"`
}

type IdentityFactor struct {
	GroupIDs          []string `hcl:"group_ids"`
	GroupNames        []string `hcl:"group_names"`
	ApprovalsRequired int      `hcl:"approvals"`
}

type ACLPermissions struct {
	CapabilitiesBitmap  uint32
	MinWrappingTTL      time.Duration
	MaxWrappingTTL      time.Duration
	AllowedParameters   map[string][]interface{}
	DeniedParameters    map[string][]interface{}
	RequiredParameters  []string
	MFAMethods          []string
	ControlGroup        *ControlGroup
	GrantingPoliciesMap map[uint32][]logical.PolicyInfo
	SubscribeEventTypes []string
}

func (p *ACLPermissions) Clone() (*ACLPermissions, error) {
	ret := &ACLPermissions{
		CapabilitiesBitmap:  p.CapabilitiesBitmap,
		MinWrappingTTL:      p.MinWrappingTTL,
		MaxWrappingTTL:      p.MaxWrappingTTL,
		RequiredParameters:  p.RequiredParameters[:],
		SubscribeEventTypes: p.SubscribeEventTypes[:],
	}

	switch {
	case p.AllowedParameters == nil:
	case len(p.AllowedParameters) == 0:
		ret.AllowedParameters = make(map[string][]interface{})
	default:
		clonedAllowed, err := copystructure.Copy(p.AllowedParameters)
		if err != nil {
			return nil, err
		}
		ret.AllowedParameters = clonedAllowed.(map[string][]interface{})
	}

	switch {
	case p.DeniedParameters == nil:
	case len(p.DeniedParameters) == 0:
		ret.DeniedParameters = make(map[string][]interface{})
	default:
		clonedDenied, err := copystructure.Copy(p.DeniedParameters)
		if err != nil {
			return nil, err
		}
		ret.DeniedParameters = clonedDenied.(map[string][]interface{})
	}

	switch {
	case p.MFAMethods == nil:
	case len(p.MFAMethods) == 0:
		ret.MFAMethods = []string{}
	default:
		clonedMFAMethods, err := copystructure.Copy(p.MFAMethods)
		if err != nil {
			return nil, err
		}
		ret.MFAMethods = clonedMFAMethods.([]string)
	}

	switch {
	case p.ControlGroup == nil:
	default:
		clonedControlGroup, err := copystructure.Copy(p.ControlGroup)
		if err != nil {
			return nil, err
		}
		ret.ControlGroup = clonedControlGroup.(*ControlGroup)
	}

	switch {
	case p.GrantingPoliciesMap == nil:
	case len(p.GrantingPoliciesMap) == 0:
		ret.GrantingPoliciesMap = make(map[uint32][]logical.PolicyInfo)
	default:
		clonedGrantingPoliciesMap, err := copystructure.Copy(p.GrantingPoliciesMap)
		if err != nil {
			return nil, err
		}
		ret.GrantingPoliciesMap = clonedGrantingPoliciesMap.(map[uint32][]logical.PolicyInfo)
	}

	return ret, nil
}

func addGrantingPoliciesToMap(m map[uint32][]logical.PolicyInfo, policy *Policy, capabilitiesBitmap uint32) map[uint32][]logical.PolicyInfo {
	if m == nil {
		m = make(map[uint32][]logical.PolicyInfo)
	}

	// For all possible policies, check if the provided capabilities include
	// them
	for _, capability := range cap2Int {
		if capabilitiesBitmap&capability == 0 {
			continue
		}

		m[capability] = append(m[capability], logical.PolicyInfo{
			Name:          policy.Name,
			NamespaceId:   policy.namespace.ID,
			NamespacePath: policy.namespace.Path,
			Type:          "acl",
		})
	}

	return m
}

// ParseACLPolicy is used to parse the specified ACL rules into an
// intermediary set of policies, before being compiled into
// the ACL
func ParseACLPolicy(ns *namespace.Namespace, rules string) (*Policy, error) {
	return parseACLPolicyWithTemplating(ns, rules, false, nil, nil)
}

// parseACLPolicyWithTemplating performs the actual work and checks whether we
// should perform substitutions. If performTemplating is true we know that it
// is templated so we don't check again, otherwise we check to see if it's a
// templated policy.
func parseACLPolicyWithTemplating(ns *namespace.Namespace, rules string, performTemplating bool, entity *identity.Entity, groups []*identity.Group) (*Policy, error) {
	// Parse the rules
	root, err := hcl.Parse(rules)
	if err != nil {
		return nil, fmt.Errorf("failed to parse policy: %w", err)
	}

	// Top-level item should be the object list
	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("failed to parse policy: does not contain a root object")
	}

	// Check for invalid top-level keys
	valid := []string{
		"name",
		"path",
	}
	if err := hclutil.CheckHCLKeys(list, valid); err != nil {
		return nil, fmt.Errorf("failed to parse policy: %w", err)
	}

	// Create the initial policy and store the raw text of the rules
	p := Policy{
		Raw:       rules,
		Type:      PolicyTypeACL,
		namespace: ns,
	}
	if err := hcl.DecodeObject(&p, list); err != nil {
		return nil, fmt.Errorf("failed to parse policy: %w", err)
	}

	if o := list.Filter("path"); len(o.Items) > 0 {
		if err := parsePaths(&p, o, performTemplating, entity, groups); err != nil {
			return nil, fmt.Errorf("failed to parse policy: %w", err)
		}
	}

	return &p, nil
}

func parsePaths(result *Policy, list *ast.ObjectList, performTemplating bool, entity *identity.Entity, groups []*identity.Group) error {
	paths := make([]*PathRules, 0, len(list.Items))
	for _, item := range list.Items {
		key := "path"
		if len(item.Keys) > 0 {
			key = item.Keys[0].Token.Value().(string)
		}

		// Check the path
		if performTemplating {
			_, templated, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
				Mode:        identitytpl.ACLTemplating,
				String:      key,
				Entity:      identity.ToSDKEntity(entity),
				Groups:      identity.ToSDKGroups(groups),
				NamespaceID: result.namespace.ID,
			})
			if err != nil {
				continue
			}
			key = templated
		} else {
			hasTemplating, _, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
				Mode:              identitytpl.ACLTemplating,
				ValidityCheckOnly: true,
				String:            key,
			})
			if err != nil {
				return fmt.Errorf("failed to validate policy templating: %w", err)
			}
			if hasTemplating {
				result.Templated = true
			}
		}

		valid := []string{
			"comment",
			"policy",
			"capabilities",
			"allowed_parameters",
			"denied_parameters",
			"required_parameters",
			"min_wrapping_ttl",
			"max_wrapping_ttl",
			"mfa_methods",
			"control_group",
			"subscribe_event_types",
		}
		if err := hclutil.CheckHCLKeys(item.Val, valid); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("path %q:", key))
		}

		var pc PathRules

		// allocate memory so that DecodeObject can initialize the ACLPermissions struct
		pc.Permissions = new(ACLPermissions)

		pc.Path = key

		if err := hcl.DecodeObject(&pc, item.Val); err != nil {
			return multierror.Prefix(err, fmt.Sprintf("path %q:", key))
		}

		// Strip a leading '/' as paths in Vault start after the / in the API path
		if len(pc.Path) > 0 && pc.Path[0] == '/' {
			pc.Path = pc.Path[1:]
		}

		// Ensure we are using the full request path internally
		pc.Path = result.namespace.Path + pc.Path

		if strings.Contains(pc.Path, "+*") {
			return fmt.Errorf("path %q: invalid use of wildcards ('+*' is forbidden)", pc.Path)
		}

		if pc.Path == "+" || strings.Count(pc.Path, "/+") > 0 || strings.HasPrefix(pc.Path, "+/") {
			pc.HasSegmentWildcards = true
		}

		if strings.HasSuffix(pc.Path, "*") {
			// If there are segment wildcards, don't actually strip the
			// trailing asterisk, but don't want to hit the default case
			if !pc.HasSegmentWildcards {
				// Strip the glob character if found
				pc.Path = strings.TrimSuffix(pc.Path, "*")
				pc.IsPrefix = true
			}
		}

		// Map old-style policies into capabilities
		if len(pc.Policy) > 0 {
			switch pc.Policy {
			case OldDenyPathPolicy:
				pc.Capabilities = []string{DenyCapability}
			case OldReadPathPolicy:
				pc.Capabilities = append(pc.Capabilities, []string{ReadCapability, ListCapability}...)
			case OldWritePathPolicy:
				pc.Capabilities = append(pc.Capabilities, []string{CreateCapability, ReadCapability, UpdateCapability, DeleteCapability, ListCapability}...)
			case OldSudoPathPolicy:
				pc.Capabilities = append(pc.Capabilities, []string{CreateCapability, ReadCapability, UpdateCapability, DeleteCapability, ListCapability, SudoCapability}...)
			default:
				return fmt.Errorf("path %q: invalid policy %q", key, pc.Policy)
			}
		}

		// Initialize the map
		pc.Permissions.CapabilitiesBitmap = 0
		for _, cap := range pc.Capabilities {
			switch cap {
			// If it's deny, don't include any other capability
			case DenyCapability:
				pc.Capabilities = []string{DenyCapability}
				pc.Permissions.CapabilitiesBitmap = DenyCapabilityInt
				goto PathFinished
			case CreateCapability, ReadCapability, UpdateCapability, DeleteCapability, ListCapability, SudoCapability, PatchCapability, SubscribeCapability:
				pc.Permissions.CapabilitiesBitmap |= cap2Int[cap]
			default:
				return fmt.Errorf("path %q: invalid capability %q", key, cap)
			}
		}

		if pc.AllowedParametersHCL != nil {
			pc.Permissions.AllowedParameters = make(map[string][]interface{}, len(pc.AllowedParametersHCL))
			for k, v := range pc.AllowedParametersHCL {
				pc.Permissions.AllowedParameters[strings.ToLower(k)] = v
			}
		}
		if pc.DeniedParametersHCL != nil {
			pc.Permissions.DeniedParameters = make(map[string][]interface{}, len(pc.DeniedParametersHCL))

			for k, v := range pc.DeniedParametersHCL {
				pc.Permissions.DeniedParameters[strings.ToLower(k)] = v
			}
		}
		if pc.MinWrappingTTLHCL != nil {
			dur, err := parseutil.ParseDurationSecond(pc.MinWrappingTTLHCL)
			if err != nil {
				return fmt.Errorf("error parsing min_wrapping_ttl: %w", err)
			}
			pc.Permissions.MinWrappingTTL = dur
		}
		if pc.MaxWrappingTTLHCL != nil {
			dur, err := parseutil.ParseDurationSecond(pc.MaxWrappingTTLHCL)
			if err != nil {
				return fmt.Errorf("error parsing max_wrapping_ttl: %w", err)
			}
			pc.Permissions.MaxWrappingTTL = dur
		}
		if pc.MFAMethodsHCL != nil {
			pc.Permissions.MFAMethods = make([]string, len(pc.MFAMethodsHCL))
			copy(pc.Permissions.MFAMethods, pc.MFAMethodsHCL)
		}
		if pc.ControlGroupHCL != nil {
			pc.Permissions.ControlGroup = new(ControlGroup)
			if pc.ControlGroupHCL.TTL != nil {
				dur, err := parseutil.ParseDurationSecond(pc.ControlGroupHCL.TTL)
				if err != nil {
					return fmt.Errorf("error parsing control group max ttl: %w", err)
				}
				pc.Permissions.ControlGroup.TTL = dur
			}
			var factors []*ControlGroupFactor
			if pc.ControlGroupHCL.Factors != nil {
				for key, factor := range pc.ControlGroupHCL.Factors {
					// Although we only have one factor here, we need to check to make sure there is at least
					// one factor defined in this factor block.
					if factor.Identity == nil {
						return errors.New("no control_group factor provided")
					}

					if factor.Identity.ApprovalsRequired <= 0 ||
						(len(factor.Identity.GroupIDs) == 0 && len(factor.Identity.GroupNames) == 0) {
						return errors.New("must provide more than one identity group and approvals > 0")
					}

					// Ensure that configured ControlledCapabilities for factor are a subset of the
					// Capabilities of the policy.
					if len(factor.ControlledCapabilities) > 0 {
						var found bool
						for _, controlledCapability := range factor.ControlledCapabilities {
							found = false
							for _, policyCap := range pc.Capabilities {
								if controlledCapability == policyCap {
									found = true
								}
							}
							if !found {
								return errors.New(ControlledCapabilityPolicySubsetError)
							}
						}
					}

					factors = append(factors, &ControlGroupFactor{
						Name:                   key,
						Identity:               factor.Identity,
						ControlledCapabilities: factor.ControlledCapabilities,
					})
				}
			}
			if len(factors) == 0 {
				return errors.New("no control group factors provided")
			}
			pc.Permissions.ControlGroup.Factors = factors
		}
		if pc.Permissions.MinWrappingTTL != 0 &&
			pc.Permissions.MaxWrappingTTL != 0 &&
			pc.Permissions.MaxWrappingTTL < pc.Permissions.MinWrappingTTL {
			return errors.New("max_wrapping_ttl cannot be less than min_wrapping_ttl")
		}
		if len(pc.RequiredParametersHCL) > 0 {
			pc.Permissions.RequiredParameters = pc.RequiredParametersHCL[:]
		}
		if len(pc.SubscribeEventTypesHCL) > 0 {
			pc.Permissions.SubscribeEventTypes = pc.SubscribeEventTypesHCL[:]
		}

	PathFinished:
		paths = append(paths, &pc)
	}

	result.Paths = paths
	return nil
}
