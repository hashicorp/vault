package vault

import (
	"github.com/armon/go-radix"
	"github.com/hashicorp/vault/logical"
)

var (
	// Policy lists are used to restrict what is eligible for an operation
	anyPolicy     = []string{PathPolicyDeny}
	readWriteSudo = []string{PathPolicyRead, PathPolicyWrite, PathPolicySudo}
	writeSudo     = []string{PathPolicyWrite, PathPolicySudo}

	// permittedPolicyLevel is used to map each logical operation
	// into the set of policies that allow the operation.
	permittedPolicyLevels = map[logical.Operation][]string{
		logical.ReadOperation:     readWriteSudo,
		logical.WriteOperation:    writeSudo,
		logical.DeleteOperation:   writeSudo,
		logical.ListOperation:     readWriteSudo,
		logical.HelpOperation:     anyPolicy,
		logical.RevokeOperation:   writeSudo,
		logical.RenewOperation:    writeSudo,
		logical.RollbackOperation: writeSudo,
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
		for _, pp := range policy.Paths {
			// Check which tree to use
			tree := a.exactRules
			if pp.Glob {
				tree = a.globRules
			}

			// Check for an existing policy
			raw, ok := tree.Get(pp.Prefix)
			if !ok {
				tree.Insert(pp.Prefix, pp)
				continue
			}
			existing := raw.(*PathPolicy)

			// Check if this policy is takes precedence
			if pp.TakesPrecedence(existing) {
				tree.Insert(pp.Prefix, pp)
			}
		}
	}
	return a, nil
}

// AllowOperation is used to check if the given operation is permitted
func (a *ACL) AllowOperation(op logical.Operation, path string) bool {
	// Fast-path root
	if a.root {
		return true
	}

	// Check if any policy level allows this operation
	permitted := permittedPolicyLevels[op]
	if permitted[0] == PathPolicyDeny {
		return true
	}

	// Find an exact matching rule, look for glob if no match
	var policy *PathPolicy
	raw, ok := a.exactRules.Get(path)
	if ok {
		policy = raw.(*PathPolicy)
		goto CHECK
	}

	// Find a glob rule, default deny if no match
	_, raw, ok = a.globRules.LongestPrefix(path)
	if !ok {
		return false
	} else {
		policy = raw.(*PathPolicy)
	}

CHECK:
	// Check if the minimum permissions are met
	for _, allowed := range permitted {
		if allowed == policy.Policy {
			return true
		}
	}
	return false
}

// RootPrivilege checks if the user has root level permission
// to given path. This requires that the user be root, or that
// sudo privilege is available on that path.
func (a *ACL) RootPrivilege(path string) bool {
	// Fast-path root
	if a.root {
		return true
	}

	// Find an exact matching rule, look for glob if no match
	var policy *PathPolicy
	raw, ok := a.exactRules.Get(path)
	if ok {
		policy = raw.(*PathPolicy)
		goto CHECK
	}

	// Check the rules for a match, default deny if no match
	_, raw, ok = a.globRules.LongestPrefix(path)
	if !ok {
		return false
	} else {
		policy = raw.(*PathPolicy)
	}

CHECK:
	// Check the policy level
	return policy.Policy == PathPolicySudo
}
