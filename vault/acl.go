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
				tree.Insert(pc.Prefix, pc.CapabilitiesBitmap)
				continue
			}
			existing := raw.(uint32)

			switch {
			case existing&DenyCapabilityInt > 0:
				// If we are explicitly denied in the existing capability set,
				// don't save anything else

			case pc.CapabilitiesBitmap&DenyCapabilityInt > 0:
				// If this new policy explicitly denies, only save the deny value
				tree.Insert(pc.Prefix, DenyCapabilityInt)

			default:
				// Insert the capabilities in this new policy into the existing
				// value; since it's a pointer we can just modify the
				// underlying data
				tree.Insert(pc.Prefix, existing|pc.CapabilitiesBitmap)
			}
		}
	}
	return a, nil
}

// AllowOperation is used to check if the given operation is permitted. The
// first bool indicates if an op is allowed, the second whether sudo priviliges
// exist for that op and path.
func (a *ACL) AllowOperation(op logical.Operation, path string) (allowed bool, sudo bool) {
	// Fast-path root
	if a.root {
		return true, true
	}

	// Help is always allowed
	if op == logical.HelpOperation {
		return true, false
	}

	// Find an exact matching rule, look for glob if no match
	var capabilities uint32
	raw, ok := a.exactRules.Get(path)
	if ok {
		capabilities = raw.(uint32)
		goto CHECK
	}

	// Find a glob rule, default deny if no match
	_, raw, ok = a.globRules.LongestPrefix(path)
	if !ok {
		return false, false
	} else {
		capabilities = raw.(uint32)
	}

CHECK:
	// Check if the minimum permissions are met
	// If "deny" has been explicitly set, only deny will be in the map, so we
	// only need to check for the existence of other values
	sudo = capabilities&SudoCapabilityInt > 0
	switch op {
	case logical.ReadOperation:
		allowed = capabilities&ReadCapabilityInt > 0
	case logical.ListOperation:
		allowed = capabilities&ListCapabilityInt > 0
	case logical.UpdateOperation:
		allowed = capabilities&UpdateCapabilityInt > 0
	case logical.DeleteOperation:
		allowed = capabilities&DeleteCapabilityInt > 0
	case logical.CreateOperation:
		allowed = capabilities&CreateCapabilityInt > 0

	// These three re-use UpdateCapabilityInt since that's the most appropraite capability/operation mapping
	case logical.RevokeOperation, logical.RenewOperation, logical.RollbackOperation:
		allowed = capabilities&UpdateCapabilityInt > 0

	default:
		return false, false
	}
	return
}
