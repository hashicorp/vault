package vault

import (
	"github.com/armon/go-radix"
	"github.com/hashicorp/vault/logical"
)

var (
	// permittedPolicyLevel is used to map each logical operation
	// into the set of capabilities that allow the operation.
	permittedPolicyLevels = map[logical.Operation][]string{
		logical.CreateOperation:   []string{CreateCapability},
		logical.ReadOperation:     []string{ReadCapability},
		logical.UpdateOperation:   []string{UpdateCapability},
		logical.DeleteOperation:   []string{DeleteCapability},
		logical.ListOperation:     []string{ListCapability},
		logical.RevokeOperation:   []string{UpdateCapability},
		logical.RenewOperation:    []string{UpdateCapability},
		logical.RollbackOperation: []string{UpdateCapability},

		// Help is special-cased to always be allowed, so we don't need anything in this list
		logical.HelpOperation: []string{},
	}
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
				tree.Insert(pc.Prefix, pc)
				continue
			}
			existing := raw.(*PathCapabilities)

			switch {
			case existing.CapabilitiesMap[DenyCapability]:
				// If we are explicitly denied in the existing capability set,
				// don't save anything else

			case pc.CapabilitiesMap[DenyCapability]:
				// If this new policy explicitly denies, only save the deny value
				tree.Insert(pc.Prefix, pc)

			default:
				// Insert the capabilities in this new policy into the existing
				// value; since it's a pointer we can just modify the
				// underlying data

				for k, _ := range pc.CapabilitiesMap {
					existing.CapabilitiesMap[k] = true
				}
			}
		}
	}
	return a, nil
}

// AllowOperation is used to check if the given operation is permitted
func (a *ACL) AllowOperation(op logical.Operation, path string) (opAllowed bool, sudoPriv bool) {
	// Fast-path root
	if a.root {
		return true, true
	}

	// Help is always allowed
	if op == logical.HelpOperation {
		return true, false
	}

	// Find an exact matching rule, look for glob if no match
	var policy *PathCapabilities
	raw, ok := a.exactRules.Get(path)
	if ok {
		policy = raw.(*PathCapabilities)
		goto CHECK
	}

	// Find a glob rule, default deny if no match
	_, raw, ok = a.globRules.LongestPrefix(path)
	if !ok {
		return false, false
	} else {
		policy = raw.(*PathCapabilities)
	}

CHECK:
	// Check if the minimum permissions are met
	// If "deny" has been explicitly set, only deny will be in the map, so we
	// only need to check for the existence of other values
	return policy.CapabilitiesMap[op.String()], policy.CapabilitiesMap[SudoCapability]
}
