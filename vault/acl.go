// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/armon/go-radix"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault/observations"
	"github.com/mitchellh/copystructure"
)

// ACL is used to wrap a set of policies to provide
// an efficient interface for access control.
type ACL struct {
	entAcl

	// exactRules contains the path policies that are exact
	exactRules *radix.Tree

	// prefixRules contains the path policies that are a prefix
	prefixRules *radix.Tree

	segmentWildcardPaths map[string]interface{}

	// root is enabled if the "root" named policy is present
	root bool

	// Stores policies that are actually RGPs for later fetching
	rgpPolicies []*Policy
}

type PolicyCheckOpts struct {
	RootPrivsRequired          bool
	Unauth                     bool
	CheckSourcePath            bool
	RecoverAlternateCapability *logical.Operation
}

type AuthResults struct {
	entAuthResults
	ACLResults      *ACLResults
	SentinelResults *SentinelResults
	Allowed         bool
	RootPrivs       bool
	DeniedError     bool
	Error           *multierror.Error
}

type ACLResults struct {
	Allowed             bool
	RootPrivs           bool
	IsRoot              bool
	MFAMethods          []string
	ControlGroup        *ControlGroup
	CapabilitiesBitmap  uint32
	GrantingPolicies    []logical.PolicyInfo
	SubscribeEventTypes []string
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

			if len(pc.Permissions.SubscribeEventTypes) > 0 {
				if len(existingPerms.SubscribeEventTypes) > 0 {
					existingPerms.SubscribeEventTypes = strutil.RemoveDuplicates(append(existingPerms.SubscribeEventTypes, pc.Permissions.SubscribeEventTypes...), false)
				} else {
					existingPerms.SubscribeEventTypes = slices.Clone(pc.Permissions.SubscribeEventTypes)
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

func (a *ACL) CapabilitiesAndSubscribeEventTypes(ctx context.Context, path string) (pathCapabilities []string, subscribeEventTypes []string) {
	req := &logical.Request{
		Path: path,
		// doesn't matter, but use List to trigger fallback behavior so we can
		// model real behavior
		Operation: logical.ListOperation,
	}

	res := a.AllowOperation(ctx, req, true)
	if res.IsRoot {
		return []string{RootCapability}, []string{"*"}
	}
	subscribeEventTypes = res.SubscribeEventTypes
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
	if capabilities&RecoverCapabilityInt > 0 {
		pathCapabilities = append(pathCapabilities, RecoverCapability)
	}

	// If "deny" is explicitly set or if the path has no capabilities at all,
	// set the path capabilities to "deny"
	if capabilities&DenyCapabilityInt > 0 || len(pathCapabilities) == 0 {
		pathCapabilities = []string{DenyCapability}
	}

	return
}

func (a *ACL) Capabilities(ctx context.Context, path string) []string {
	pathCapabilities, _ := a.CapabilitiesAndSubscribeEventTypes(ctx, path)
	return pathCapabilities
}

// AllowOperation is used to check if the given operation is permitted.
func (a *ACL) AllowOperation(ctx context.Context, req *logical.Request, capCheckOnly bool) (ret *ACLResults) {
	ret = a.performEnterpriseAclChecks(ctx, req, capCheckOnly)
	if ret != nil {
		return ret
	}

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
		ret.SubscribeEventTypes = []string{"*"}
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

	// For LIST operations with a trailing slash we must consider candidates
	// from *both* the slash-stripped path ("kv1/private") and the full path
	// ("kv1/private/") before choosing a winner.
	//
	// Why both forms are needed:
	//   - Prefix rules (e.g. "kv1/private/*" stored as key "kv1/private/")
	//     are only reachable via LongestPrefix when the query string starts
	//     with that key.  "kv1/private/" starts with "kv1/private/" but
	//     "kv1/private" (trimmed) does NOT, so a deny on "kv1/private/*"
	//     would be silently skipped if we only use the trimmed form.
	//   - Segment-wildcard rules (e.g. "kv1/+") are matched by splitting on
	//     "/".  Split("kv1/private/") yields ["kv1","private",""] (3 parts)
	//     which does NOT match the 2-segment rule "kv1/+", while
	//     Split("kv1/private") yields ["kv1","private"] which does.  So the
	//     trimmed form is essential to correctly honour segment-wildcard
	//     policies on LIST paths (original VAULT-3825 intent).
	//
	// Selection logic:
	//   1. Collect candidates from both the full path and the trimmed path.
	//   2. Among all candidates, find the most-specific deny (if any).
	//   3. Find the most-specific overall winner (deny or allow) via the
	//      existing specificity comparator.
	//   4. If the most-specific deny is at least as specific as the overall
	//      winner, use the deny — this closes the CVE.
	//   5. Otherwise use the overall winner — this preserves the VAULT-3825
	//      behaviour where a trimmed-match grants the LIST capability even
	//      when a different (and more-specific) full-path rule exists but
	//      does not deny.
	if op == logical.ListOperation && strings.HasSuffix(path, "/") {
		trimmedDescrs := a.prefixAndWCCandidatesForPath(strings.TrimSuffix(path, "/"), nil)
		fullDescrs := a.prefixAndWCCandidatesForPath(path, nil)

		if len(trimmedDescrs) > 0 || len(fullDescrs) > 0 {
			permissions = a.resolveACLPermsForListOp(trimmedDescrs, fullDescrs)
			capabilities = permissions.CapabilitiesBitmap
			goto CHECK
		}
		// No prefix or segment-wildcard rule matched either form; fall through
		// to the "no match" return below.
		return
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
		ret.SubscribeEventTypes = slices.Clone(permissions.SubscribeEventTypes)
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
	case logical.RecoverOperation:
		operationAllowed = capabilities&RecoverCapabilityInt > 0
		grantingPolicies = permissions.GrantingPoliciesMap[RecoverCapabilityInt]

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
	if op == logical.ReadOperation || op == logical.UpdateOperation || op == logical.CreateOperation || op == logical.PatchOperation || op == logical.RecoverOperation {
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

		useLegacyMatching := os.Getenv("VAULT_LEGACY_EXACT_MATCHING_ON_LIST") != ""

		if len(permissions.DeniedParameters) > 0 {
			// Check if all parameters have been denied
			if _, ok := permissions.DeniedParameters["*"]; ok {
				return
			}

			for parameter, value := range req.Data {
				// Check if parameter has been explicitly denied
				if valueSlice, ok := permissions.DeniedParameters[strings.ToLower(parameter)]; ok {
					normalizedValue := normalizePolicyParameterValue(parameter, value)
					if valueInDeniedParameterList(normalizedValue, valueSlice, useLegacyMatching) {
						return
					}
				}
			}
		}

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

			normalizedValue := normalizePolicyParameterValue(parameter, value)
			if ok && !valueInAllowedParameterList(normalizedValue, valueSlice, useLegacyMatching) {
				return
			}
		}
	}

	ret.Allowed = true
	return
}

// wcPathDescr is a single candidate ACL rule matched via a prefix-glob or
// segment-wildcard pattern. All fields are derived from the rule itself, not
// from the request path, so candidates from different query forms share one
// comparison space and can be ranked together by wildcardPathDescriptorComparePriority.
type wcPathDescr struct {
	// firstWCOrGlob: byte position of the first '+' or '*' in the rule.
	// For prefix rules this equals len(prefix key); for segment-wildcard
	// rules it is strings.Index(rule, "+"). Larger = more literal leading
	// text = higher specificity.
	firstWCOrGlob int

	// wildcards: number of '+' path segments consumed during matching.
	// Fewer wildcards = more precise = higher specificity.
	wildcards int

	// isPrefix: true when the rule ends with a trailing glob ('*'), making
	// it open-ended. A non-prefix ('+'-only) rule at the same firstWCOrGlob
	// is considered more specific than a prefix rule.
	isPrefix bool

	// wcPath: rule path with the trailing '*' stripped for prefix rules.
	// Used as a length and lexicographic tie-break when other tiers are equal.
	wcPath string

	// perms: permissions granted (or denied) by this rule. A set
	// DenyCapabilityInt bit signals an unconditional deny.
	perms *ACLPermissions
}

// wildcardPathDescriptorComparePriority returns true if a is lower priority than b.
// It is the single source of truth for the wcPathDescr specificity ordering,
// centralised here so that wildcardPathSpecificityLess (sort.Slice) and
// resolveACLPermsForListOp (direct candidate comparison) both delegate to it —
// keeping the ranking logic in one place prevents the two call-sites from
// silently diverging over time.
//
//   - Later first-wildcard/glob position wins (earlier occurrence = lower priority)
//   - Non-prefix beats prefix (prefix/glob-terminated = lower priority)
//   - Fewer wildcard (+) segments wins
//   - Longer wcPath wins
//   - Lexicographically larger wcPath (tie-break; in practice unreachable)
func wildcardPathDescriptorComparePriority(a, b wcPathDescr) bool {
	if a.firstWCOrGlob != b.firstWCOrGlob {
		return a.firstWCOrGlob < b.firstWCOrGlob
	}
	if a.isPrefix != b.isPrefix {
		return a.isPrefix // prefix is lower priority than non-prefix
	}
	if a.wildcards != b.wildcards {
		return a.wildcards > b.wildcards // more wildcards = lower priority
	}
	if len(a.wcPath) != len(b.wcPath) {
		return len(a.wcPath) < len(b.wcPath)
	}
	// Lexicographical tie-break. This should never really come up. It's more
	// of a throwing-up-hands scenario akin to panic("should not be here")
	// statements, but less panicky.
	return a.wcPath < b.wcPath
}

// wildcardPathSpecificityLess returns a sort.Slice less function that orders
// wcPathDescr candidates by specificity (ascending), so that the
// highest-priority candidate ends up last and can be picked with descrs[len-1].
// It is a thin index-based wrapper around wildcardPathDescriptorComparePriority.
func wildcardPathSpecificityLess(descrs []wcPathDescr) func(i, j int) bool {
	return func(i, j int) bool {
		return wildcardPathDescriptorComparePriority(descrs[i], descrs[j])
	}
}

// tryMatchWildcardPath attempts to match a single segment-wildcard rule (fullWCPath)
// against pathParts and returns the populated wcPathDescr candidate.
//
// bareMount contract: when bareMount is true, the caller is performing a
// mount-access check (does any policy grant non-deny access under this mount
// prefix?) rather than a precise path lookup. In that mode the length guards
// are relaxed and the caller needs to inspect each partial segment match
// itself (at i == len(pathParts)-2) to decide whether to accept or continue.
// To signal this boundary, tryMatchWildcardPath returns (wcPathDescr{}, false) at
// that point rather than completing the match — the caller then does the
// join-and-check inline. This means tryMatchWildcardPath never returns a non-zero
// candidate when bareMount is true; callers must handle bareMount logic
// themselves and should not pass bareMount=true unless they implement that
// inline loop.
//
// Returns (candidate, true) on a successful non-bareMount match, or
// (wcPathDescr{}, false) when the rule does not match or when bareMount is
// true (the caller drives bareMount logic inline).
func (a *ACL) tryMatchWildcardPath(fullWCPath string, pathParts []string, bareMount bool) (wcPathDescr, bool) {
	if fullWCPath == "" {
		return wcPathDescr{}, false
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
		return wcPathDescr{}, false
	}
	if !bareMount && !pd.isPrefix && len(splitCurrWCPath) != len(pathParts) {
		return wcPathDescr{}, false
	}

	segments := make([]string, 0, len(splitCurrWCPath))
	for i, aclPart := range splitCurrWCPath {
		switch {
		case aclPart == "+":
			pd.wildcards++ // each '+' consumed reduces specificity
			segments = append(segments, pathParts[i])

		case aclPart == pathParts[i]:
			segments = append(segments, pathParts[i])

		case pd.isPrefix && i == len(splitCurrWCPath)-1 && strings.HasPrefix(pathParts[i], aclPart):
			segments = append(segments, pathParts[i:]...)

		default:
			// Found a mismatch; this rule does not apply.
			return wcPathDescr{}, false
		}

		// bareMount early-return: the caller checks whether this rule provides
		// any non-deny permission under the mount prefix; we signal that by
		// returning false so the caller can handle it inline.
		// -2 because we're always invoked with a trailing "/" in bareMount mode.
		if bareMount && i == len(pathParts)-2 {
			return wcPathDescr{}, false
		}
	}
	pd.perms = a.segmentWildcardPaths[fullWCPath].(*ACLPermissions)
	return pd, true
}

// prefixAndWCCandidatesForPath collects all prefix-rule and segment-wildcard
// candidates matching path (bareMount=false semantics only) and appends them
// to descrs. It does NOT sort or pick a winner — callers are responsible for
// that. Accepting an existing slice lets a caller invoke this function twice
// with different query forms (e.g. slash-stripped and slash-retained) and
// obtain a single merged pool that can be ranked in one pass.
func (a *ACL) prefixAndWCCandidatesForPath(path string, descrs []wcPathDescr) []wcPathDescr {
	// Collect prefix rule candidate.
	if prefix, raw, ok := a.prefixRules.LongestPrefix(path); ok {
		descrs = append(descrs, wcPathDescr{
			firstWCOrGlob: len(prefix),
			wcPath:        prefix,
			isPrefix:      true,
			perms:         raw.(*ACLPermissions),
		})
	}

	if len(a.segmentWildcardPaths) == 0 {
		return descrs
	}

	pathParts := strings.Split(path, "/")
	for fullWCPath := range a.segmentWildcardPaths {
		if pd, ok := a.tryMatchWildcardPath(fullWCPath, pathParts, false); ok {
			descrs = append(descrs, pd)
		}
	}
	return descrs
}

// resolveACLPermsForListOp selects the winning ACLPermissions for a LIST
// operation that carries a trailing slash, given candidate sets from both the
// trimmed (slash-stripped) and full (slash-retained) forms of the request path.
//
// Selection logic:
//  1. Find the most-specific DENY candidate across all candidates (both sets).
//     A deny from a more-specific rule must always win regardless of operation
//     capability: a more-specific deny policy must not be bypassed by a
//     broader allow reachable only via the trimmed path.
//  2. Find the most-specific candidate that explicitly grants LIST.
//  3. Decision:
//     a. If the most-specific deny outranks the most-specific LIST-granting
//     candidate (or if no LIST-granting candidate exists), return the deny.
//     b. Otherwise return the most-specific LIST-granting candidate.
//     c. If there is no deny and no LIST-granting candidate (e.g. the only
//     matching rule grants read/write but not list), return the overall
//     most-specific candidate so the CHECK phase can evaluate its capabilities
//     normally — including returning "not allowed" when the operation is not
//     present in the bitmap.
//
// The trimmed-form candidates ensure that rules written without a trailing
// slash (e.g. "kv1/+") still match trailing-slash LIST requests. The
// full-form candidates ensure that more-specific rules keyed with a trailing
// slash (e.g. "kv1/private/") are visible to the comparator.
//
// Callers MUST only invoke this function when len(trimmedDescrs)+len(fullDescrs) > 0.
func (a *ACL) resolveACLPermsForListOp(trimmedDescrs, fullDescrs []wcPathDescr) *ACLPermissions {
	all := make([]wcPathDescr, 0, len(trimmedDescrs)+len(fullDescrs))
	all = append(all, trimmedDescrs...)
	all = append(all, fullDescrs...)

	// Find the most-specific deny and the most-specific LIST-granting candidate.
	var bestDeny *wcPathDescr
	var bestList *wcPathDescr
	for i := range all {
		pd := &all[i]
		if pd.perms.CapabilitiesBitmap&DenyCapabilityInt != 0 {
			if bestDeny == nil || wildcardPathDescriptorComparePriority(*bestDeny, *pd) {
				bestDeny = pd
			}
		}
		if pd.perms.CapabilitiesBitmap&ListCapabilityInt != 0 {
			if bestList == nil || wildcardPathDescriptorComparePriority(*bestList, *pd) {
				bestList = pd
			}
		}
	}

	// Case 3a: deny is more specific than any LIST-granting candidate → deny wins.
	if bestDeny != nil {
		if bestList == nil || wildcardPathDescriptorComparePriority(*bestList, *bestDeny) {
			return bestDeny.perms
		}
	}

	// Case 3b: return the most-specific LIST-granting candidate.
	if bestList != nil {
		return bestList.perms
	}

	// Case 3c: no deny, no LIST grant — return the overall most-specific candidate.
	sort.Slice(all, wildcardPathSpecificityLess(all))
	return all[len(all)-1].perms
}

// CheckAllowedFromNonExactPaths returns permissions corresponding to a
// matching path with wildcards/globs. If bareMount is true, the path should
// correspond to a mount prefix, and what is returned is either a non-nil set
// of permissions from some allowed path underneath the mount (for use in mount
// access checks), or nil indicating no non-deny permissions were found.
//
// bareMount=false delegates to prefixAndWCCandidatesForPath + wildcardPathSpecificityLess
// so the ranking logic is centralised and the two non-LIST call-sites share it.
//
// bareMount=true intentionally keeps its own inline loop rather than delegating
// to tryMatchWildcardPath: in mount-access mode the caller needs an early-return at
// i == len(pathParts)-2 (one segment before the trailing slash of the mount
// prefix) to check whether the partial match covers the mount — a contract that
// tryMatchWildcardPath signals with (wcPathDescr{}, false) rather than completing.
// Keeping the loop inline makes the early-return and the joinedPath check
// co-located and easy to audit together.
func (a *ACL) CheckAllowedFromNonExactPaths(path string, bareMount bool) *ACLPermissions {
	// bareMount=false: delegate to prefixAndWCCandidatesForPath + sort.
	if !bareMount {
		descrs := a.prefixAndWCCandidatesForPath(path, make([]wcPathDescr, 0, len(a.segmentWildcardPaths)+1))
		if len(descrs) == 0 {
			return nil
		}
		sort.Slice(descrs, wildcardPathSpecificityLess(descrs))
		return descrs[len(descrs)-1].perms
	}

	// bareMount=true: preserved byte-for-byte — we need the early-return
	// logic that checks non-deny access under a mount prefix.
	pathParts := strings.Split(path, "/")

SWCPATH:
	for fullWCPath := range a.segmentWildcardPaths {
		if fullWCPath == "" {
			continue
		}

		currWCPath := fullWCPath
		isPrefix := false
		if currWCPath[len(currWCPath)-1] == '*' {
			isPrefix = true
			currWCPath = currWCPath[0 : len(currWCPath)-1]
		}

		splitCurrWCPath := strings.Split(currWCPath, "/")

		// In bareMount mode len(pathParts) < len(splitCurrWCPath) is allowed
		// (we're checking a prefix of the mount).

		segments := make([]string, 0, len(splitCurrWCPath))
		for i, aclPart := range splitCurrWCPath {
			switch {
			case aclPart == "+":
				segments = append(segments, pathParts[i])

			case aclPart == pathParts[i]:
				segments = append(segments, pathParts[i])

			case isPrefix && i == len(splitCurrWCPath)-1 && strings.HasPrefix(pathParts[i], aclPart):
				segments = append(segments, pathParts[i:]...)

			default:
				continue SWCPATH
			}

			// -2 because we're always invoked with a trailing "/" in case bareMount.
			if i == len(pathParts)-2 {
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
	}

	return nil
}

func (c *Core) recordPolicyEvaluationObservation(ctx context.Context, te *logical.TokenEntry, req *logical.Request, results *AuthResults) {
	observation := map[string]interface{}{
		"request_id": req.ID,
		"path":       req.Path,
		"entity_id":  req.EntityID,
		"client_id":  req.ClientID,
	}
	if te != nil {
		observation["policies"] = te.Policies
		observation["is_root"] = te.IsRoot()
	}

	if results != nil {
		if results.ACLResults != nil {
			observation["request_allowed"] = results.Allowed
			observation["request_acl_allowed"] = results.ACLResults.Allowed

			grantingPolicies := make([]logical.PolicyInfo, 0)
			if len(results.ACLResults.GrantingPolicies) > 0 {
				grantingPolicies = append(grantingPolicies, results.ACLResults.GrantingPolicies...)
			}
			if results.SentinelResults != nil && len(results.SentinelResults.GrantingPolicies) > 0 {
				grantingPolicies = append(grantingPolicies, results.SentinelResults.GrantingPolicies...)
			}
			if len(grantingPolicies) > 0 {
				observation["granting_policies"] = grantingPolicies
			}

			err := c.Observations().RecordObservationToLedger(ctx, observations.ObservationTypePolicyACLEvaluation, nil, observation)
			if err != nil {
				c.logger.Error("error recording observation for policy checks", "error", err)
			}
		}
	}
}

func (c *Core) performPolicyChecksSinglePath(ctx context.Context, acl *ACL, te *logical.TokenEntry, req *logical.Request, inEntity *identity.Entity, opts *PolicyCheckOpts) *AuthResults {
	ret := new(AuthResults)

	// First, perform normal ACL checks if requested. The only time no ACL
	// should be applied is if we are only processing EGPs against a login
	// path in which case opts.Unauth will be set.
	if acl != nil && !opts.Unauth {
		ret.ACLResults = acl.AllowOperation(ctx, req, false)
		ret.RootPrivs = ret.ACLResults.RootPrivs
		// Root is always allowed; skip Sentinel/MFA checks
		if ret.ACLResults.IsRoot {
			ret.Allowed = true
			c.recordPolicyEvaluationObservation(ctx, te, req, ret)
			return ret
		}
		if !ret.ACLResults.Allowed {
			c.recordPolicyEvaluationObservation(ctx, te, req, ret)
			return ret
		}
		// Since HelpOperation was fast-pathed inside AllowOperation, RootPrivs will not have been populated in this
		// case, so we need to special-case that here as well, or we'll block HelpOperation on all sudo-protected paths.
		if !ret.RootPrivs && opts.RootPrivsRequired && req.Operation != logical.HelpOperation {
			c.recordPolicyEvaluationObservation(ctx, te, req, ret)
			return ret
		}
	}

	c.performEntPolicyChecks(ctx, acl, te, req, inEntity, opts, ret)

	c.recordPolicyEvaluationObservation(ctx, te, req, ret)
	return ret
}

// normalizePolicyParameterValue returns a lowercased copy of value when
// parameter is one that Vault canonicalises to lowercase internally before
// storing or enforcing policy names. Without this normalisation a caller can
// bypass a denied_parameters (or allowed_parameters) constraint by submitting
// a mixed-case variant.
func normalizePolicyParameterValue(parameter string, value interface{}) interface{} {
	switch strings.ToLower(parameter) {
	case "policies", "token_policies":
		// fall through to normalisation below
	default:
		return value
	}
	switch v := value.(type) {
	case string:
		return strings.ToLower(v)
	case []string:
		lowered := make([]interface{}, len(v))
		for i, s := range v {
			lowered[i] = strings.ToLower(s)
		}
		return lowered
	case []interface{}:
		lowered := make([]interface{}, len(v))
		for i, el := range v {
			if s, ok := el.(string); ok {
				lowered[i] = strings.ToLower(s)
			} else {
				lowered[i] = el
			}
		}
		return lowered
	default:
		return value
	}
}

func valueInAllowedParameterList(v interface{}, list []interface{}, useLegacyMatching bool) bool {
	// Empty list is equivalent to the item always existing in the list
	if len(list) == 0 {
		return true
	}

	if valueInParameterList(v, list) {
		return true
	}

	if useLegacyMatching {
		// prevent execution of the new behaviour if we're in legacy mode
		return false
	}

	if vSlice, ok := v.([]interface{}); ok {
		// when not running in legacy mode, we run a relaxed check for slices that verifies if all
		// elements in the slice exist in the allowed list, as opposed to checking if the allowed
		// list contains a single element that matches the entire slice (but this whole-slice match
		// is still supported)
		for _, v := range vSlice {
			if !valueInParameterList(v, list) {
				return false
			}
		}

		return true
	} else if vString, ok := v.(string); ok {
		// At this point we don't know if the field is of framework.TypeCommaStringSlice, but we assume it is
		// because failing to match a value because of it being in a comma-separated string is way more likely
		// and worse than accidentally matching a substring of a string value.
		if vSlice, err := parseutil.ParseCommaStringSlice(vString); err == nil {
			for _, v := range vSlice {
				if !valueInParameterList(v, list) {
					return false
				}
			}
			return true
		}
	}

	return false
}

func valueInDeniedParameterList(v interface{}, list []interface{}, useLegacyMatching bool) bool {
	// Empty list is equivalent to the item always existing in the list
	if len(list) == 0 {
		return true
	}

	if valueInParameterList(v, list) {
		return true
	}

	if useLegacyMatching {
		// prevent execution of the new behaviour if we're in legacy mode
		return false
	}

	// The new behaviour is that if any value in the slice is in the denied list, we deny.
	if vSlice, ok := v.([]interface{}); ok {
		for _, v := range vSlice {
			if valueInParameterList(v, list) {
				return true
			}
		}
	} else if vString, ok := v.(string); ok {
		// At this point we don't know if the field is of framework.TypeCommaStringSlice, but we assume it is
		// because failing to match a value because of it being in a comma-separated string is way more likely
		// and worse than accidentally matching a substring of a string value.
		if vSlice, err := parseutil.ParseCommaStringSlice(vString); err == nil {
			for _, v := range vSlice {
				if valueInParameterList(v, list) {
					return true
				}
			}
		}
	}

	return false
}

func valueInParameterList(v interface{}, list []interface{}) bool {
	for _, el := range list {
		if el == nil || v == nil {
			// It doesn't seem possible to set up a nil entry in the list, but it is possible
			// to pass in a null entry in the API request being checked. Just in case,
			// nil will match nil.
			if el == v {
				return true
			}
		} else if elStr, ok := el.(string); ok {
			if vStr, ok := v.(string); ok && strutil.GlobbedStringsMatch(elStr, vStr) {
				return true
			}
		} else if reflect.DeepEqual(el, v) {
			return true
		}
	}

	return false
}
