package vault

import (
	"github.com/armon/go-radix"
	"github.com/hashicorp/vault/logical"
)

// ACL is used to wrap a set of policies to provide
// an efficient interface for access control.
type ACL struct {
	// exactRules contains the path policies that are exact
	exactRules *radix.Tree

	// globRules contains the path policies that glob
	globRules *radix.Tree

	// root is enabled if the "root" named policy is present.
	root bool
}

type aclEntry struct {
	// The fast-path pure-bitmap capabilities set. If the MFA map is not nil, we use that instead (and this should be zero)
	capabilitiesBitmap uint32

	// A map of capabilities to MFA methods configured for them
	capabilitiesToMFAMap map[uint32][]string
}

// New is used to construct a policy based ACL from a set of policies.
func NewACL(policies []*Policy) (*ACL, error) {
	// Initialize
	a := &ACL{
		exactRules: radix.New(),
		globRules:  radix.New(),
		root:       false,
	}

	// Inject each policy
	for _, policy := range policies {
		// Ignore a nil policy object
		if policy == nil {
			continue
		}
		// Check if this is root
		if policy.Name == "root" {
			a.root = true
		}
		for _, pc := range policy.Paths {
			// Check which tree to use
			tree := a.exactRules
			if pc.Glob {
				tree = a.globRules
			}

			// Check for an existing policy
			raw, ok := tree.Get(pc.Prefix)
			if !ok {
				var entry aclEntry
				if pc.MFAMethods != nil && len(pc.MFAMethods) != 0 {
					// We have a MFA so we leave the bitmap at 0 and use the map
					entry.capabilitiesToMFAMap = convertCapBitmapToMFAMap(pc.CapabilitiesBitmap, pc.MFAMethods)
				} else {
					// We don't have MFA so we use the straight bitmap to be speedier
					entry.capabilitiesBitmap = pc.CapabilitiesBitmap
				}
				tree.Insert(pc.Prefix, entry)
				continue
			}
			existing := raw.(aclEntry)

			switch {
			case existing.capabilitiesBitmap&DenyCapabilityInt > 0:
				// If we are explicitly denied in the existing capability set,
				// don't save anything else. Explicit denial entries do not set
				// the MFA map, so checking the bitmap is sufficient.

			case pc.CapabilitiesBitmap&DenyCapabilityInt > 0:
				// If this new policy explicitly denies, only save the deny
				// value and leave the map at nil
				tree.Insert(pc.Prefix, aclEntry{
					capabilitiesBitmap: DenyCapabilityInt,
				})

			default:
				// Insert the capabilities in this new policy into the existing
				// value
				tree.Insert(pc.Prefix, *(mergeACLEntryCapabilities(&existing, pc)))
			}
		}
	}
	return a, nil
}

func mergeACLEntryCapabilities(existing *aclEntry, pc *PathCapabilities) (ret *aclEntry) {
	ret = &aclEntry{}

	// Ensure that we start out with a representative map based on existing
	switch {
	case existing.capabilitiesToMFAMap == nil && pc.MFAMethods == nil:
		// Neither are using MFA so simply merge and return
		ret.capabilitiesBitmap = existing.capabilitiesBitmap | pc.CapabilitiesBitmap
		return ret

	case existing.capabilitiesToMFAMap == nil && pc.MFAMethods != nil:
		// Convert to a map so that we can merge in the methods
		ret.capabilitiesToMFAMap = convertCapBitmapToMFAMap(existing.capabilitiesBitmap, nil)

	case existing.capabilitiesToMFAMap != nil && pc.MFAMethods == nil:
		// Start out using the existing map
		ret.capabilitiesToMFAMap = existing.capabilitiesToMFAMap
	}

	// Now, create the map for the new path capabilities object
	pathCapMFAMap := convertCapBitmapToMFAMap(pc.CapabilitiesBitmap, pc.MFAMethods)

	/*
		defer func() {
			panic(fmt.Sprintf("existing:\n%#v\n, pc:\n%#v\n, currRetMap:\n%#v\n, pathCapMFAMap: \n%#v\n, ret:\n%#v\n", *existing, *pc, ret.capabilitiesToMFAMap, pathCapMFAMap, *ret))
		}()
	*/

	switch {
	case ret.capabilitiesToMFAMap == nil && pathCapMFAMap == nil:
		// This really shouldn't happen, because in the first switch if the
		// existing maps were nil it should have merged bitmaps and
		// returned, but safety.
		return ret

	case pathCapMFAMap == nil:
		// We have nothing(?) from the new pathcapabilities so return what we do have
		return ret

	case ret.capabilitiesToMFAMap == nil:
		// In this case we had no existing capabilities, so use whatever has
		// come down the line
		ret.capabilitiesToMFAMap = pathCapMFAMap
		return ret
	}

	// If we are at this point, both MFA maps are not nil, so we merge them.
	for cap, mfas := range pathCapMFAMap {
		ret.capabilitiesToMFAMap[cap] = append(ret.capabilitiesToMFAMap[cap], mfas...)
	}

	return ret
}

func convertCapBitmapToMFAMap(bitmap uint32, mfaMethods []string) map[uint32][]string {
	if bitmap == 0 {
		return nil
	}

	ret := make(map[uint32][]string, 6)

	// We use nil as a baseline to save memory
	if mfaMethods != nil && len(mfaMethods) == 0 {
		mfaMethods = nil
	}

	if bitmap&SudoCapabilityInt > 0 {
		ret[SudoCapabilityInt] = mfaMethods
	}
	if bitmap&ReadCapabilityInt > 0 {
		ret[ReadCapabilityInt] = mfaMethods
	}
	if bitmap&ListCapabilityInt > 0 {
		ret[ListCapabilityInt] = mfaMethods
	}
	if bitmap&UpdateCapabilityInt > 0 {
		ret[UpdateCapabilityInt] = mfaMethods
	}
	if bitmap&DeleteCapabilityInt > 0 {
		ret[DeleteCapabilityInt] = mfaMethods
	}
	if bitmap&CreateCapabilityInt > 0 {
		ret[CreateCapabilityInt] = mfaMethods
	}

	return ret
}

func (a *ACL) Capabilities(path string) (pathCapabilities []string) {
	// Fast-path root
	if a.root {
		return []string{RootCapability}
	}

	// Find an exact matching rule, look for glob if no match
	var entry aclEntry
	raw, ok := a.exactRules.Get(path)
	if ok {
		entry = raw.(aclEntry)
		goto CHECK
	}

	// Find a glob rule, default deny if no match
	_, raw, ok = a.globRules.LongestPrefix(path)
	if !ok {
		return []string{DenyCapability}
	} else {
		entry = raw.(aclEntry)
	}

CHECK:
	if entry.capabilitiesToMFAMap == nil {
		if entry.capabilitiesBitmap&SudoCapabilityInt > 0 {
			pathCapabilities = append(pathCapabilities, SudoCapability)
		}
		if entry.capabilitiesBitmap&ReadCapabilityInt > 0 {
			pathCapabilities = append(pathCapabilities, ReadCapability)
		}
		if entry.capabilitiesBitmap&ListCapabilityInt > 0 {
			pathCapabilities = append(pathCapabilities, ListCapability)
		}
		if entry.capabilitiesBitmap&UpdateCapabilityInt > 0 {
			pathCapabilities = append(pathCapabilities, UpdateCapability)
		}
		if entry.capabilitiesBitmap&DeleteCapabilityInt > 0 {
			pathCapabilities = append(pathCapabilities, DeleteCapability)
		}
		if entry.capabilitiesBitmap&CreateCapabilityInt > 0 {
			pathCapabilities = append(pathCapabilities, CreateCapability)
		}

		// If "deny" is explicitly set or if the path has no capabilities at all,
		// set the path capabilities to "deny"
		if entry.capabilitiesBitmap&DenyCapabilityInt > 0 || len(pathCapabilities) == 0 {
			pathCapabilities = []string{DenyCapability}
		}

		return
	}

	// Potentially have some mfa methods to return, so look at the map
	if _, ok := entry.capabilitiesToMFAMap[SudoCapabilityInt]; ok {
		pathCapabilities = append(pathCapabilities, SudoCapability)
	}
	if _, ok := entry.capabilitiesToMFAMap[ReadCapabilityInt]; ok {
		pathCapabilities = append(pathCapabilities, ReadCapability)
	}
	if _, ok := entry.capabilitiesToMFAMap[ListCapabilityInt]; ok {
		pathCapabilities = append(pathCapabilities, ListCapability)
	}
	if _, ok := entry.capabilitiesToMFAMap[UpdateCapabilityInt]; ok {
		pathCapabilities = append(pathCapabilities, UpdateCapability)
	}
	if _, ok := entry.capabilitiesToMFAMap[DeleteCapabilityInt]; ok {
		pathCapabilities = append(pathCapabilities, DeleteCapability)
	}
	if _, ok := entry.capabilitiesToMFAMap[CreateCapabilityInt]; ok {
		pathCapabilities = append(pathCapabilities, CreateCapability)
	}

	// We don't store Deny in map; instead we clear the map and store just the bitmap, so we don't need to check deny here

	return
}

// AllowOperation is used to check if the given operation is permitted. The
// first bool indicates if an op is allowed, the second whether sudo priviliges
// exist for that op and path.
func (a *ACL) AllowOperation(op logical.Operation, path string) (allowed bool, sudo bool, mfaMethods, sudoMFAMethods []string) {
	// Fast-path root
	if a.root {
		return true, true, nil, nil
	}

	// Help is always allowed
	if op == logical.HelpOperation {
		return true, false, nil, nil
	}

	// Find an exact matching rule, look for glob if no match
	var entry aclEntry
	raw, ok := a.exactRules.Get(path)
	if ok {
		entry = raw.(aclEntry)
		goto CHECK
	}

	// Find a glob rule, default deny if no match
	_, raw, ok = a.globRules.LongestPrefix(path)
	if !ok {
		return false, false, nil, nil
	} else {
		entry = raw.(aclEntry)
	}

CHECK:
	// Check if the minimum permissions are met
	// If "deny" has been explicitly set, only deny will be in the map, so we
	// only need to check for the existence of other values
	if entry.capabilitiesToMFAMap == nil {
		sudo = entry.capabilitiesBitmap&SudoCapabilityInt > 0
		switch op {
		case logical.ReadOperation:
			allowed = entry.capabilitiesBitmap&ReadCapabilityInt > 0
		case logical.ListOperation:
			allowed = entry.capabilitiesBitmap&ListCapabilityInt > 0
		case logical.UpdateOperation:
			allowed = entry.capabilitiesBitmap&UpdateCapabilityInt > 0
		case logical.DeleteOperation:
			allowed = entry.capabilitiesBitmap&DeleteCapabilityInt > 0
		case logical.CreateOperation:
			allowed = entry.capabilitiesBitmap&CreateCapabilityInt > 0

		// These three re-use UpdateCapabilityInt since that's the most appropraite capability/operation mapping
		case logical.RevokeOperation, logical.RenewOperation, logical.RollbackOperation:
			allowed = entry.capabilitiesBitmap&UpdateCapabilityInt > 0

		default:
			return false, false, nil, nil
		}

		return
	}

	// Potentially have some mfa methods to return, so look at the map
	switch op {
	case logical.ReadOperation:
		if methods, ok := entry.capabilitiesToMFAMap[ReadCapabilityInt]; ok {
			allowed = true
			mfaMethods = append(mfaMethods, methods...)
		}
	case logical.ListOperation:
		if methods, ok := entry.capabilitiesToMFAMap[ListCapabilityInt]; ok {
			allowed = true
			mfaMethods = append(mfaMethods, methods...)
		}
	case logical.UpdateOperation:
		if methods, ok := entry.capabilitiesToMFAMap[UpdateCapabilityInt]; ok {
			allowed = true
			mfaMethods = append(mfaMethods, methods...)
		}
	case logical.DeleteOperation:
		if methods, ok := entry.capabilitiesToMFAMap[DeleteCapabilityInt]; ok {
			allowed = true
			mfaMethods = append(mfaMethods, methods...)
		}
	case logical.CreateOperation:
		if methods, ok := entry.capabilitiesToMFAMap[CreateCapabilityInt]; ok {
			allowed = true
			mfaMethods = append(mfaMethods, methods...)
		}

	// These three re-use UpdateCapabilityInt since that's the most appropriate capability/operation mapping
	case logical.RevokeOperation, logical.RenewOperation, logical.RollbackOperation:
		if _, ok := entry.capabilitiesToMFAMap[UpdateCapabilityInt]; ok {
			allowed = true
		}
	}

	// Sudo is additive, so we need to return a set of mfa methods specifically
	// for sudo use so as not to end up adding any defined on sudo capability
	// to all other methods
	if methods, ok := entry.capabilitiesToMFAMap[SudoCapabilityInt]; ok {
		sudo = true
		sudoMFAMethods = append(mfaMethods, methods...)
	}

	return
}
