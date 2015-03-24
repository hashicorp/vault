package vault

import (
	"github.com/armon/go-radix"
	"github.com/hashicorp/vault/logical"
)

// operationPolicyLevel is used to map each logical operation
// into the minimum required permissions to allow the operation.
var operationPolicyLevel = map[logical.Operation]int{
	logical.ReadOperation:   pathPolicyLevel[PathPolicyRead],
	logical.WriteOperation:  pathPolicyLevel[PathPolicyWrite],
	logical.DeleteOperation: pathPolicyLevel[PathPolicyWrite],
	logical.ListOperation:   pathPolicyLevel[PathPolicyRead],
	logical.RevokeOperation: pathPolicyLevel[PathPolicyWrite],
	logical.RenewOperation:  pathPolicyLevel[PathPolicyRead],
	logical.HelpOperation:   pathPolicyLevel[PathPolicyDeny],
}

// ACL is used to wrap a set of policies to provide
// an efficient interface for access control.
type ACL struct {
	// pathRules contains the path policies
	pathRules *radix.Tree

	// root is enabled if the "root" named policy is present.
	root bool
}

// New is used to construct a policy based ACL from a set of policies.
func NewACL(policies []*Policy) (*ACL, error) {
	// Initialize
	a := &ACL{
		pathRules: radix.New(),
		root:      false,
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
			// Convert to a policy level
			policyLevel := pathPolicyLevel[pp.Policy]

			// Check for an existing policy
			raw, ok := a.pathRules.Get(pp.Prefix)
			if !ok {
				a.pathRules.Insert(pp.Prefix, policyLevel)
				continue
			}
			existing := raw.(int)

			// Check if this policy is a higher access level,
			// we want to store the highest permission permitted.
			if policyLevel > existing {
				a.pathRules.Insert(pp.Prefix, policyLevel)
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

	// Find a matching rule, default deny if no match
	policyLevel := 0
	_, rule, ok := a.pathRules.LongestPrefix(path)
	if ok {
		policyLevel = rule.(int)
	}

	// Convert the operation to a minimum required level
	requiredLevel := operationPolicyLevel[op]

	// Check if the minimum permissions are met
	return policyLevel >= requiredLevel
}

// RootPrivilege checks if the user has root level permission
// to given path. This requires that the user be root, or that
// sudo privilege is available on that path.
func (a *ACL) RootPrivilege(path string) bool {
	// Fast-path root
	if a.root {
		return true
	}

	// Check the rules for a match
	_, rule, ok := a.pathRules.LongestPrefix(path)
	if !ok {
		return false
	}

	// Check the policy level
	policyLevel := rule.(int)
	return policyLevel == pathPolicyLevel[PathPolicySudo]
}
