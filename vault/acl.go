// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/armon/go-radix"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/copystructure"
)

// ACL is used to wrap a set of policies to provide
// an efficient interface for access control.
type ACL struct {
	// exactRules contains the path policies that are exact
	exactRules *radix.Tree

	// prefixRules contains the path policies that are a prefix
	prefixRules *radix.Tree

	segmentWildcardPaths map[string]interface{}

	// root is enabled if the "root" named policy is present.
	root bool

	// Stores policies that are actually RGPs for later fetching
	rgpPolicies []*Policy
}

type PolicyCheckOpts struct {
	RootPrivsRequired bool
	Unauth            bool
}

type AuthResults struct {
	ACLResults      *ACLResults
	SentinelResults *SentinelResults
	Allowed         bool
	RootPrivs       bool
	DeniedError     bool
	Error           *multierror.Error
}

type ACLResults struct {
	Allowed            bool
	RootPrivs          bool
	IsRoot             bool
	MFAMethods         []string
	ControlGroup       *ControlGroup
	CapabilitiesBitmap uint32
	GrantingPolicies   []logical.PolicyInfo
}

type SentinelResults struct {
	GrantingPolicies []logical.PolicyInfo
}

// NewACL is used to construct a policy based ACL from a set of policies.
func NewACL(ctx context.Context, policies []*Policy) (*ACL, error) {
	// Initialize
	a := &ACL{
		exactRules:           radix.New(),
		prefixRules:          radix.New(),
		segmentWildcardPaths: make(map[string]interface{}, len(policies)),
		root:                 false,
	}

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if ns == nil {
		return nil, namespace.ErrNoNamespace
	}

	// Inject each policy
	for _, policy := range policies {
		// Ignore a nil policy object
		if policy == nil {
			continue
		}

		switch policy.Type {
		case PolicyTypeACL:
		case PolicyTypeRGP:
			a.rgpPolicies = append(a.rgpPolicies, policy)
			continue
		default:
			return nil, fmt.Errorf("unable to parse policy (wrong type)")
		}

		// Check if this is root
		if policy.Name == "root" {
			if ns.ID != namespace.RootNamespaceID {
				return nil, fmt.Errorf("root policy is only allowed in root namespace")
			}

			if len(policies) != 1 {
				return nil, fmt.Errorf("other policies present along with root")
			}
			a.root = true
		}

		for _, pc := range policy.Paths {
			var raw interface{}
			var ok bool
			var tree *radix.Tree

			switch {
			case pc.HasSegmentWildcards:
				raw, ok = a.segmentWildcardPaths[pc.Path]
			default:
				// Check which tree to use
				tree = a.exactRules
				if pc.IsPrefix {
					tree = a.prefixRules
				}

				// Check for an existing policy
				raw, ok = tree.Get(pc.Path)
			}

			if !ok {
				clonedPerms, err := pc.Permissions.Clone()
				if err != nil {
					return nil, fmt.Errorf("error cloning ACL permissions: %w", err)
				}

				// Store this policy name as the policy that permits these
				// capabilities
				clonedPerms.GrantingPoliciesMap = addGrantingPoliciesToMap(nil, policy, clonedPerms.CapabilitiesBitmap)
				switch {
				case pc.HasSegmentWildcards:
					a.segmentWildcardPaths[pc.Path] = clonedPerms
				default:
					tree.Insert(pc.Path, clonedPerms)
				}
				continue
			}

			// these are the ones already in the tree
			existingPerms := raw.(*ACLPermissions)

			switch {
			case existingPerms.CapabilitiesBitmap&DenyCapabilityInt > 0:
				// If we are explicitly denied in the existing capability set,
				// don't save anything else
				continue

			case pc.Permissions.CapabilitiesBitmap&DenyCapabilityInt > 0:
				// If this new policy explicitly denies, only save the deny value
				existingPerms.CapabilitiesBitmap = DenyCapabilityInt
				existingPerms.AllowedParameters = nil
				existingPerms.DeniedParameters = nil
				goto INSERT

			default:
				// Insert the capabilities in this new policy into the existing
				// value
				existingPerms.CapabilitiesBitmap = existingPerms.CapabilitiesBitmap | pc.Permissions.CapabilitiesBitmap
				existingPerms.GrantingPoliciesMap = addGrantingPoliciesToMap(existingPerms.GrantingPoliciesMap, policy, pc.Permissions.CapabilitiesBitmap)
			}

			// Note: In these stanzas, we're preferring minimum lifetimes. So
			// we take the lesser of two specified max values, or we take the
			// lesser of two specified min values, the idea being, allowing
			// token lifetime to be minimum possible.
			//
			// If we have an existing max, and we either don't have a current
			// max, or the current is greater than the previous, use the
			// existing.
			if pc.Permissions.MaxWrappingTTL > 0 &&
				(existingPerms.MaxWrappingTTL == 0 ||
					pc.Permissions.MaxWrappingTTL < existingPerms.MaxWrappingTTL) {
				existingPerms.MaxWrappingTTL = pc.Permissions.MaxWrappingTTL
			}
			// If we have an existing min, and we either don't have a current
			// min, or the current is greater than the previous, use the
			// existing
			if pc.Permissions.MinWrappingTTL > 0 &&
				(existingPerms.MinWrappingTTL == 0 ||
					pc.Permissions.MinWrappingTTL < existingPerms.MinWrappingTTL) {
				existingPerms.MinWrappingTTL = pc.Permissions.MinWrappingTTL
			}

			if len(pc.Permissions.AllowedParameters) > 0 {
				if existingPerms.AllowedParameters == nil {
					clonedAllowed, err := copystructure.Copy(pc.Permissions.AllowedParameters)
					if err != nil {
						return nil, err
					}
					existingPerms.AllowedParameters = clonedAllowed.(map[string][]interface{})
				} else {
					for key, value := range pc.Permissions.AllowedParameters {
						pcValue, ok := existingPerms.AllowedParameters[key]
						// If an empty array exist it should overwrite any other
						// value.
						if len(value) == 0 || (ok && len(pcValue) == 0) {
							existingPerms.AllowedParameters[key] = []interface{}{}
						} else {
							// Merge the two maps, appending values on key conflict.
							existingPerms.AllowedParameters[key] = append(value, existingPerms.AllowedParameters[key]...)
						}
					}
				}
			}

			if len(pc.Permissions.DeniedParameters) > 0 {
				if existingPerms.DeniedParameters == nil {
					clonedDenied, err := copystructure.Copy(pc.Permissions.DeniedParameters)
					if err != nil {
						return nil, err
					}
					existingPerms.DeniedParameters = clonedDenied.(map[string][]interface{})
				} else {
					for key, value := range pc.Permissions.DeniedParameters {
						pcValue, ok := existingPerms.DeniedParameters[key]
						// If an empty array exist it should overwrite any other
						// value.
						if len(value) == 0 || (ok && len(pcValue) == 0) {
							existingPerms.DeniedParameters[key] = []interface{}{}
						} else {
							// Merge the two maps, appending values on key conflict.
							existingPerms.DeniedParameters[key] = append(value, existingPerms.DeniedParameters[key]...)
						}
					}
				}
			}

			if len(pc.Permissions.RequiredParameters) > 0 {
				if len(existingPerms.RequiredParameters) == 0 {
					existingPerms.RequiredParameters = pc.Permissions.RequiredParameters
				} else {
					for _, v := range pc.Permissions.RequiredParameters {
						if !strutil.StrListContains(existingPerms.RequiredParameters, v) {
							existingPerms.RequiredParameters = append(existingPerms.RequiredParameters, v)
						}
					}
				}
			}

			if len(pc.Permissions.MFAMethods) > 0 {
				if existingPerms.MFAMethods == nil {
					existingPerms.MFAMethods = pc.Permissions.MFAMethods
				} else {
					existingPerms.MFAMethods = append(existingPerms.MFAMethods, pc.Permissions.MFAMethods...)
				}
				existingPerms.MFAMethods = strutil.RemoveDuplicates(existingPerms.MFAMethods, false)
			}

			// No need to dedupe this list since any authorization can satisfy any factor, so long as
			// the factor matches the specified permission requested.
			if pc.Permissions.ControlGroup != nil {
				if len(pc.Permissions.ControlGroup.Factors) > 0 {
					if existingPerms.ControlGroup == nil {
						cg, err := pc.Permissions.ControlGroup.Clone()
						if err != nil {
							return nil, err
						}
						existingPerms.ControlGroup = cg
					} else {
						existingPerms.ControlGroup.Factors = append(existingPerms.ControlGroup.Factors, pc.Permissions.ControlGroup.Factors...)
					}
				}
			}

		INSERT:
			switch {
			case pc.HasSegmentWildcards:
				a.segmentWildcardPaths[pc.Path] = existingPerms
			default:
				tree.Insert(pc.Path, existingPerms)
			}
		}
	}
	return a, nil
}

func (a *ACL) Capabilities(ctx context.Context, path string) (pathCapabilities []string) {
	req := &logical.Request{
		Path: path,
		// doesn't matter, but use List to trigger fallback behavior so we can
		// model real behavior
		Operation: logical.ListOperation,
	}

	res := a.AllowOperation(ctx, req, true)
	if res.IsRoot {
		return []string{RootCapability}
	}

	capabilities := res.CapabilitiesBitmap

	if capabilities&SudoCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, SudoCapability)
	}
	if capabilities&ReadCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, ReadCapability)
	}
	if capabilities&ListCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, ListCapability)
	}
	if capabilities&UpdateCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, UpdateCapability)
	}
	if capabilities&DeleteCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, DeleteCapability)
	}
	if capabilities&CreateCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, CreateCapability)
	}
	if capabilities&PatchCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, PatchCapability)
	}
	if capabilities&SubscribeCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, SubscribeCapability)
	}

	// If "deny" is explicitly set or if the path has no capabilities at all,
	// set the path capabilities to "deny"
	if capabilities&DenyCapabilityInt > 0 || len(pathCapabilities) == 0 {
		pathCapabilities = []string{DenyCapability}
	}
	return
}

// AllowOperation is used to check if the given operation is permitted.
func (a *ACL) AllowOperation(ctx context.Context, req *logical.Request, capCheckOnly bool) (ret *ACLResults) {
	ret = new(ACLResults)

	// Fast-path root
	if a.root {
		ret.Allowed = true
		ret.RootPrivs = true
		ret.IsRoot = true
		ret.GrantingPolicies = []logical.PolicyInfo{{
			Name:          "root",
			NamespaceId:   "root",
			NamespacePath: "",
			Type:          "acl",
		}}
		return
	}
	op := req.Operation

	// Help is always allowed
	if op == logical.HelpOperation {
		ret.Allowed = true
		return
	}

	var permissions *ACLPermissions

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return
	}
	path := ns.Path + req.Path

	// The request path should take care of this already but this is useful for
	// tests and as defense in depth
	for {
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		} else {
			break
		}
	}

	// Find an exact matching rule, look for prefix if no match
	var capabilities uint32
	raw, ok := a.exactRules.Get(path)
	if ok {
		permissions = raw.(*ACLPermissions)
		capabilities = permissions.CapabilitiesBitmap
		goto CHECK
	}
	if op == logical.ListOperation {
		raw, ok = a.exactRules.Get(strings.TrimSuffix(path, "/"))
		if ok {
			permissions = raw.(*ACLPermissions)
			capabilities = permissions.CapabilitiesBitmap
			goto CHECK
		}
	}

	permissions = a.CheckAllowedFromNonExactPaths(path, false)
	if permissions != nil {
		capabilities = permissions.CapabilitiesBitmap
		goto CHECK
	}

	// No exact, prefix, or segment wildcard paths found, return without
	// setting allowed
	return

CHECK:
	// Check if the minimum permissions are met
	// If "deny" has been explicitly set, only deny will be in the map, so we
	// only need to check for the existence of other values
	ret.RootPrivs = capabilities&SudoCapabilityInt > 0

	// This is after the RootPrivs check so we can gate on it being from sudo
	// rather than policy root
	if capCheckOnly {
		ret.CapabilitiesBitmap = capabilities
		return ret
	}

	ret.MFAMethods = permissions.MFAMethods
	ret.ControlGroup = permissions.ControlGroup

	var grantingPolicies []logical.PolicyInfo
	operationAllowed := false
	switch op {
	case logical.ReadOperation:
		operationAllowed = capabilities&ReadCapabilityInt > 0
		grantingPolicies = permissions.GrantingPoliciesMap[ReadCapabilityInt]
	case logical.ListOperation:
		operationAllowed = capabilities&ListCapabilityInt > 0
		grantingPolicies = permissions.GrantingPoliciesMap[ListCapabilityInt]
	case logical.UpdateOperation:
		operationAllowed = capabilities&UpdateCapabilityInt > 0
		grantingPolicies = permissions.GrantingPoliciesMap[UpdateCapabilityInt]
	case logical.DeleteOperation:
		operationAllowed = capabilities&DeleteCapabilityInt > 0
		grantingPolicies = permissions.GrantingPoliciesMap[DeleteCapabilityInt]
	case logical.CreateOperation:
		operationAllowed = capabilities&CreateCapabilityInt > 0
		grantingPolicies = permissions.GrantingPoliciesMap[CreateCapabilityInt]
	case logical.PatchOperation:
		operationAllowed = capabilities&PatchCapabilityInt > 0
		grantingPolicies = permissions.GrantingPoliciesMap[PatchCapabilityInt]

	// These three re-use UpdateCapabilityInt since that's the most appropriate
	// capability/operation mapping
	case logical.RevokeOperation, logical.RenewOperation, logical.RollbackOperation:
		operationAllowed = capabilities&UpdateCapabilityInt > 0
		grantingPolicies = permissions.GrantingPoliciesMap[UpdateCapabilityInt]

	default:
		return
	}

	if !operationAllowed {
		return
	}

	ret.GrantingPolicies = grantingPolicies

	if permissions.MaxWrappingTTL > 0 {
		if req.WrapInfo == nil || req.WrapInfo.TTL > permissions.MaxWrappingTTL {
			return
		}
	}
	if permissions.MinWrappingTTL > 0 {
		if req.WrapInfo == nil || req.WrapInfo.TTL < permissions.MinWrappingTTL {
			return
		}
	}
	// This situation can happen because of merging, even though in a single
	// path statement we check on ingress
	if permissions.MinWrappingTTL != 0 &&
		permissions.MaxWrappingTTL != 0 &&
		permissions.MaxWrappingTTL < permissions.MinWrappingTTL {
		return
	}

	// Only check parameter permissions for operations that can modify
	// parameters.
	if op == logical.ReadOperation || op == logical.UpdateOperation || op == logical.CreateOperation || op == logical.PatchOperation {
		for _, parameter := range permissions.RequiredParameters {
			if _, ok := req.Data[strings.ToLower(parameter)]; !ok {
				return
			}
		}

		// If there are no data fields, allow
		if len(req.Data) == 0 {
			ret.Allowed = true
			return
		}

		if len(permissions.DeniedParameters) == 0 {
			goto ALLOWED_PARAMETERS
		}

		// Check if all parameters have been denied
		if _, ok := permissions.DeniedParameters["*"]; ok {
			return
		}

		for parameter, value := range req.Data {
			// Check if parameter has been explicitly denied
			if valueSlice, ok := permissions.DeniedParameters[strings.ToLower(parameter)]; ok {
				// If the value exists in denied values slice, deny
				if valueInParameterList(value, valueSlice) {
					return
				}
			}
		}

	ALLOWED_PARAMETERS:
		// If we don't have any allowed parameters set, allow
		if len(permissions.AllowedParameters) == 0 {
			ret.Allowed = true
			return
		}

		_, allowedAll := permissions.AllowedParameters["*"]
		if len(permissions.AllowedParameters) == 1 && allowedAll {
			ret.Allowed = true
			return
		}

		for parameter, value := range req.Data {
			valueSlice, ok := permissions.AllowedParameters[strings.ToLower(parameter)]
			// Requested parameter is not in allowed list
			if !ok && !allowedAll {
				return
			}

			// If the value doesn't exists in the allowed values slice,
			// deny
			if ok && !valueInParameterList(value, valueSlice) {
				return
			}
		}
	}

	ret.Allowed = true
	return
}

type wcPathDescr struct {
	firstWCOrGlob int
	wildcards     int
	isPrefix      bool
	wcPath        string
	perms         *ACLPermissions
}

// CheckAllowedFromNonExactPaths returns permissions corresponding to a
// matching path with wildcards/globs. If bareMount is true, the path should
// correspond to a mount prefix, and what is returned is either a non-nil set
// of permissions from some allowed path underneath the mount (for use in mount
// access checks), or nil indicating no non-deny permissions were found.
func (a *ACL) CheckAllowedFromNonExactPaths(path string, bareMount bool) *ACLPermissions {
	wcPathDescrs := make([]wcPathDescr, 0, len(a.segmentWildcardPaths)+1)

	less := func(i, j int) bool {
		// In the case of multiple matches, we use this priority order,
		// which tries to most closely match longest-prefix:
		//
		// * First glob or wildcard position (prefer foo/a* over foo/+,
		//   foo/bar/+/baz over foo/+/bar/baz)
		// * Whether it's a prefix (prefer foo/+/bar over foo/+/ba*,
		//   foo/+ over foo/*)
		// * Number of wildcard segments (prefer foo/bar/+/baz over foo/+/+/baz)
		// * Length check (prefer foo/+/bar/ba* over foo/+/bar/b*)
		// * Lexicographical ordering (preferring less, arbitrarily)
		//
		// That final case (lexigraphical) should never really come up. It's more
		// of a throwing-up-hands scenario akin to panic("should not be here")
		// statements, but less panicky.

		pdi, pdj := wcPathDescrs[i], wcPathDescrs[j]

		// If the first wildcard (+) or glob (*) occurs earlier in pdi,
		// pdi is lower priority
		if pdi.firstWCOrGlob < pdj.firstWCOrGlob {
			return true
		} else if pdi.firstWCOrGlob > pdj.firstWCOrGlob {
			return false
		}

		// If pdi ends in * and pdj doesn't, pdi is lower priority
		if pdi.isPrefix && !pdj.isPrefix {
			return true
		} else if !pdi.isPrefix && pdj.isPrefix {
			return false
		}

		// If pdi has more wc segs, pdi is lower priority
		if pdi.wildcards > pdj.wildcards {
			return true
		} else if pdi.wildcards < pdj.wildcards {
			return false
		}

		// If pdi is shorter, it is lower priority
		if len(pdi.wcPath) < len(pdj.wcPath) {
			return true
		} else if len(pdi.wcPath) > len(pdj.wcPath) {
			return false
		}

		// If pdi is smaller lexicographically, it is lower priority
		if pdi.wcPath < pdj.wcPath {
			return true
		} else if pdi.wcPath > pdj.wcPath {
			return false
		}
		return false
	}

	// Find a prefix rule if any.
	{
		prefix, raw, ok := a.prefixRules.LongestPrefix(path)
		if ok {
			if len(a.segmentWildcardPaths) == 0 {
				return raw.(*ACLPermissions)
			}
			wcPathDescrs = append(wcPathDescrs, wcPathDescr{
				firstWCOrGlob: len(prefix),
				wcPath:        prefix,
				isPrefix:      true,
				perms:         raw.(*ACLPermissions),
			})
		}
	}

	if len(a.segmentWildcardPaths) == 0 {
		return nil
	}

	pathParts := strings.Split(path, "/")

SWCPATH:
	for fullWCPath := range a.segmentWildcardPaths {
		if fullWCPath == "" {
			continue
		}
		pd := wcPathDescr{firstWCOrGlob: strings.Index(fullWCPath, "+")}

		currWCPath := fullWCPath
		if currWCPath[len(currWCPath)-1] == '*' {
			pd.isPrefix = true
			currWCPath = currWCPath[0 : len(currWCPath)-1]
		}
		pd.wcPath = currWCPath

		splitCurrWCPath := strings.Split(currWCPath, "/")

		if !bareMount && len(pathParts) < len(splitCurrWCPath) {
			// check if the path coming in is shorter; if so it can't match
			continue
		}
		if !bareMount && !pd.isPrefix && len(splitCurrWCPath) != len(pathParts) {
			// If it's not a prefix we expect the same number of segments
			continue
		}

		segments := make([]string, 0, len(splitCurrWCPath))
		for i, aclPart := range splitCurrWCPath {
			switch {
			case aclPart == "+":
				pd.wildcards++
				segments = append(segments, pathParts[i])

			case aclPart == pathParts[i]:
				segments = append(segments, pathParts[i])

			case pd.isPrefix && i == len(splitCurrWCPath)-1 && strings.HasPrefix(pathParts[i], aclPart):
				segments = append(segments, pathParts[i:]...)

			case !bareMount:
				// Found a mismatch, give up on this segmentWildcardPath
				continue SWCPATH
			}

			// -2 because we're always invoked with a trailing "/" in case bareMount.
			if bareMount && i == len(pathParts)-2 {
				joinedPath := strings.Join(segments, "/") + "/"
				// Check the current joined path so far. If we find a prefix,
				// check permissions. If they're defined but not deny, success.
				if strings.HasPrefix(joinedPath, path) {
					permissions := a.segmentWildcardPaths[fullWCPath].(*ACLPermissions)
					if permissions.CapabilitiesBitmap&DenyCapabilityInt == 0 && permissions.CapabilitiesBitmap > 0 {
						return permissions
					}
				}
				continue SWCPATH
			}
		}
		pd.perms = a.segmentWildcardPaths[fullWCPath].(*ACLPermissions)
		wcPathDescrs = append(wcPathDescrs, pd)
	}

	if bareMount || len(wcPathDescrs) == 0 {
		return nil
	}

	// We don't do this in the bare mount check because we don't care about
	// priority, we only care about any capability at all.
	sort.Slice(wcPathDescrs, less)

	return wcPathDescrs[len(wcPathDescrs)-1].perms
}

func (c *Core) performPolicyChecks(ctx context.Context, acl *ACL, te *logical.TokenEntry, req *logical.Request, inEntity *identity.Entity, opts *PolicyCheckOpts) *AuthResults {
	ret := new(AuthResults)

	// First, perform normal ACL checks if requested. The only time no ACL
	// should be applied is if we are only processing EGPs against a login
	// path in which case opts.Unauth will be set.
	if acl != nil && !opts.Unauth {
		ret.ACLResults = acl.AllowOperation(ctx, req, false)
		ret.RootPrivs = ret.ACLResults.RootPrivs
		// Root is always allowed; skip Sentinel/MFA checks
		if ret.ACLResults.IsRoot {
			// logger.Warn("token is root, skipping checks")
			ret.Allowed = true
			return ret
		}
		if !ret.ACLResults.Allowed {
			return ret
		}
		// Since HelpOperation was fast-pathed inside AllowOperation, RootPrivs will not have been populated in this
		// case, so we need to special-case that here as well, or we'll block HelpOperation on all sudo-protected paths.
		if !ret.RootPrivs && opts.RootPrivsRequired && req.Operation != logical.HelpOperation {
			return ret
		}
	}

	c.performEntPolicyChecks(ctx, acl, te, req, inEntity, opts, ret)

	return ret
}

func valueInParameterList(v interface{}, list []interface{}) bool {
	// Empty list is equivalent to the item always existing in the list
	if len(list) == 0 {
		return true
	}

	return valueInSlice(v, list)
}

func valueInSlice(v interface{}, list []interface{}) bool {
	for _, el := range list {
		if el == nil || v == nil {
			// It doesn't seem possible to set up a nil entry in the list, but it is possible
			// to pass in a null entry in the API request being checked. Just in case,
			// nil will match nil.
			if el == v {
				return true
			}
		} else if reflect.TypeOf(el).String() == "string" && reflect.TypeOf(v).String() == "string" {
			item := el.(string)
			val := v.(string)

			if strutil.GlobbedStringsMatch(item, val) {
				return true
			}
		} else if reflect.DeepEqual(el, v) {
			return true
		}
	}

	return false
}
