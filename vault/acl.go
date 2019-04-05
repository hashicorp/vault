package vault

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	radix "github.com/armon/go-radix"
	"github.com/hashicorp/errwrap"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
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
	ACLResults  *ACLResults
	Allowed     bool
	RootPrivs   bool
	DeniedError bool
	Error       *multierror.Error
}

type ACLResults struct {
	Allowed            bool
	RootPrivs          bool
	IsRoot             bool
	MFAMethods         []string
	ControlGroup       *ControlGroup
	CapabilitiesBitmap uint32
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
					return nil, errwrap.Wrapf("error cloning ACL permissions: {{err}}", err)
				}
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
					for _, method := range pc.Permissions.MFAMethods {
						existingPerms.MFAMethods = append(existingPerms.MFAMethods, method)
					}
				}
				existingPerms.MFAMethods = strutil.RemoveDuplicates(existingPerms.MFAMethods, false)
			}

			// No need to dedupe this list since any authorization can satisfy any factor
			if pc.Permissions.ControlGroup != nil {
				if len(pc.Permissions.ControlGroup.Factors) > 0 {
					if existingPerms.ControlGroup == nil {
						existingPerms.ControlGroup = pc.Permissions.ControlGroup
					} else {
						for _, authz := range pc.Permissions.ControlGroup.Factors {
							existingPerms.ControlGroup.Factors = append(existingPerms.ControlGroup.Factors, authz)
						}
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

	// Find a prefix rule, default deny if no match
	_, raw, ok = a.prefixRules.LongestPrefix(path)
	if ok {
		permissions = raw.(*ACLPermissions)
		capabilities = permissions.CapabilitiesBitmap
		goto CHECK
	}

	permissions = a.CheckAllowedFromSegmentWildcardPaths(path, false)
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

	operationAllowed := false
	switch op {
	case logical.ReadOperation:
		operationAllowed = capabilities&ReadCapabilityInt > 0
	case logical.ListOperation:
		operationAllowed = capabilities&ListCapabilityInt > 0
	case logical.UpdateOperation:
		operationAllowed = capabilities&UpdateCapabilityInt > 0
	case logical.DeleteOperation:
		operationAllowed = capabilities&DeleteCapabilityInt > 0
	case logical.CreateOperation:
		operationAllowed = capabilities&CreateCapabilityInt > 0

	// These three re-use UpdateCapabilityInt since that's the most appropriate
	// capability/operation mapping
	case logical.RevokeOperation, logical.RenewOperation, logical.RollbackOperation:
		operationAllowed = capabilities&UpdateCapabilityInt > 0

	default:
		return
	}

	if !operationAllowed {
		return
	}

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
	if op == logical.ReadOperation || op == logical.UpdateOperation || op == logical.CreateOperation {
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

// CheckAllowedFromSegmentWildcardPaths returns permissions corresponding to a
// matching path with wildcard segments. If bareMount is true, the path should
// correspond to a mount prefix, and what is returned is either a non-nil set
// of permissions from some allowed path underneath the mount (for use in mount
// access checks), or nil indicating no non-deny permissions were found.
func (a *ACL) CheckAllowedFromSegmentWildcardPaths(path string, bareMount bool) *ACLPermissions {
	if len(a.segmentWildcardPaths) == 0 {
		return nil
	}

	// We use these to find the most-specific matching path. In the case of
	// multiple matches, we use this priority order, which tries to most
	// closely match longest-prefix:
	//
	// * First wildcard position (prefer foo/bar/+/baz over foo/+/bar/baz)
	// * Number of wildcard segments (prefer foo/bar/+/baz over foo/+/+/baz)
	// * Total path segments (prefer foo/bar/+/baz/why over foo/bar/+/ba*)
	// * Whether it's a prefix (prefer foo/+/bar over foo/+/ba*)
	// * Length check (prefer foo/+/bar/ba* over foo/+/bar/b*)
	// * Lexicographical ordering (preferring less, arbitrarily)
	//
	// That final case (lexicographical) should never really come up. It's more
	// of a throwing-up-hands scenario akin to panic("should not be here")
	// statements, but less panicky.
	var totalPathSegments int
	var totalWildcardSegments int
	var currIsPrefix bool
	var permissions *ACLPermissions
	var currFoundPath string
	var currFirstWC int = -1
	constructedPath := make([]string, 0, 10)

	pathParts := strings.Split(path, "/")
	for currWCPath := range a.segmentWildcardPaths {
		if currWCPath == "" {
			continue
		}

		var isPrefix bool
		var invalid bool
		var wcSegments int
		var firstWC int = -1

		origCurrWCPath := currWCPath

		if currWCPath[len(currWCPath)-1] == '*' {
			isPrefix = true
			currWCPath = currWCPath[0 : len(currWCPath)-1]
		}
		splitCurrWCPath := strings.Split(currWCPath, "/")
		if !bareMount && len(pathParts) < len(splitCurrWCPath) {
			// If not a bare mount, check if the path coming in is shorter; if
			// so it can't match
			continue
		}
		if !isPrefix && !bareMount && len(splitCurrWCPath) != len(pathParts) {
			// If it's not a prefix or bare mount we expect the same number of
			// segments
			continue
		}

		//
		// Priority handling
		//
		// We don't do this in the bare mount check because we don't care about
		// priority, we only care about any capability at all.
		if !bareMount {
			for i, part := range splitCurrWCPath {
				if part == "+" {
					wcSegments++
					if firstWC == -1 {
						firstWC = i
					}
				}
			}
			switch {
			case currFirstWC == 0:
				// No current match so carry on
			case currFirstWC < firstWC:
				// Current one comes before so is less specific, continue on
			case currFirstWC > firstWC:
				// Current one is after so it's more specific, use that
				continue
			default:
				// Note: the various other global values will always be
				// populated if currFirstWC is
				//
				// Same first position of wildcard. Next switch on least
				// wildcard segments.
				switch {
				case totalWildcardSegments > wcSegments:
					// Continue on here; if it matches, this is less segments so we
					// want to use that as it's more specific
				case totalWildcardSegments < wcSegments:
					// What we have is more specific already, so use that
					continue
				default:
					// Same number of total wildcard segments, so next check
					// total path segments
					switch {
					// N.B.: We'd really only have a disparate number of segments
					// if a wildcard was also in play.
					case totalPathSegments < len(splitCurrWCPath):
						// We've found a match, but the current path has more overall
						// segments, hence is more specific, so keep going -- if it's not
						// invalid we'll replace the current match with it
					case totalPathSegments > len(splitCurrWCPath):
						// We already have a match and it's more overall segments (more
						// specific) so use that
						continue
					default:
						// Same number of path segments. Next switch on prefix.
						if isPrefix != currIsPrefix {
							if isPrefix {
								// This is a prefix but what we found earlier
								// isn't; use that
								continue
							}
							// Carry on checking
						} else {
							// At this point do a length check; failing that do
							// a lexicographical ordering. I actually can't
							// figure out when we'd get to the lexicographical
							// case, based on coming up with test cases, so
							// it's for safety more than anyting else.
							switch {
							case len(currFoundPath) > len(origCurrWCPath):
								// Use current one, it's longer
								continue
							case len(currFoundPath) < len(origCurrWCPath):
							// Try for the new one, it's longer
							default:
								if currFoundPath < origCurrWCPath {
									continue
								}
							}
						}
					}
				}
			}
		}

		constructedPath = constructedPath[:0]
		// We key off splitK here since it might be less than pathParts
		for i, aclPart := range splitCurrWCPath {
			var skipAfterConsPathCheck bool
			if aclPart == "+" {
				// Matches anything in the segment, so keep checking
				constructedPath = append(constructedPath, pathParts[i])
				skipAfterConsPathCheck = true
			} else {
				constructedPath = append(constructedPath, aclPart)
			}
			if bareMount {
				joinedConstructedPath := strings.Join(constructedPath, "/")
				// Check the current joined path so far. If we find a prefix,
				// check permissions. If they're defined but not deny, success.
				if strings.HasPrefix(joinedConstructedPath, path) {
					permissions = a.segmentWildcardPaths[origCurrWCPath].(*ACLPermissions)
					if permissions.CapabilitiesBitmap&DenyCapabilityInt == 0 && permissions.CapabilitiesBitmap > 0 {
						return permissions
					} else {
						// If we already found a match and the permissions
						// don't check out we're not going to do any better
						// looking at the rest of the path, so keep on with the
						// next one instead
						invalid = true
						break
					}
				}
			}
			if skipAfterConsPathCheck {
				continue
			}
			if i == len(splitCurrWCPath)-1 && isPrefix {
				// In this case we may have foo* or just * depending on if
				// originally it was foo* or foo/*.
				if aclPart == "" {
					// Ended in /*, so at this point we're at the final
					// glob which will match anything, so return success
					break
				}
				if !strings.HasPrefix(pathParts[i], aclPart) {
					// E.g., the final part of the acl is foo* and the
					// final part of the path is boofar
					invalid = true
					break
				}
				// Final prefixed matched and the rest is a wildcard,
				// matches
				break
			}
			if aclPart != pathParts[i] {
				// Mismatch, exit out
				invalid = true
				break
			}
		}
		// If invalid isn't set then we got through the full segmented path
		// without finding a mismatch, so it's valid
		if !invalid && !bareMount {
			permissions = a.segmentWildcardPaths[origCurrWCPath].(*ACLPermissions)
			totalPathSegments = len(splitCurrWCPath)
			totalWildcardSegments = wcSegments
			currIsPrefix = isPrefix
			currFoundPath = origCurrWCPath
			currFirstWC = firstWC
		}
	}

	return permissions
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
			//logger.Warn("token is root, skipping checks")
			ret.Allowed = true
			return ret
		}
		if !ret.ACLResults.Allowed {
			return ret
		}
		if !ret.RootPrivs && opts.RootPrivsRequired {
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
		if reflect.TypeOf(el).String() == "string" && reflect.TypeOf(v).String() == "string" {
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
