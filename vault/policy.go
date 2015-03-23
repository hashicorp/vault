package vault

import (
	"fmt"

	"github.com/hashicorp/hcl"
)

const (
	PathPolicyDeny  = "deny"
	PathPolicyRead  = "read"
	PathPolicyWrite = "write"
	PathPolicySudo  = "sudo"
)

var (
	pathPolicyLevel = map[string]int{
		PathPolicyDeny:  0,
		PathPolicyRead:  1,
		PathPolicyWrite: 2,
		PathPolicySudo:  3,
	}
)

// Policy is used to represent the policy specified by
// an ACL configuration.
type Policy struct {
	Name  string        `hcl:"name"`
	Paths []*PathPolicy `hcl:"path,expand"`
	Raw   string
}

// PathPolicy represents a policy for a path in the namespace
type PathPolicy struct {
	Prefix string `hcl:",key"`
	Policy string
}

// Parse is used to parse the specified ACL rules into an
// intermediary set of policies, before being compiled into
// the ACL
func Parse(rules string) (*Policy, error) {
	// Decode the rules
	p := &Policy{Raw: rules}
	if err := hcl.Decode(p, rules); err != nil {
		return nil, fmt.Errorf("Failed to parse ACL rules: %v", err)
	}

	// Validate the path policy
	for _, pp := range p.Paths {
		switch pp.Policy {
		case PathPolicyDeny:
		case PathPolicyRead:
		case PathPolicyWrite:
		case PathPolicySudo:
		default:
			return nil, fmt.Errorf("Invalid path policy: %#v", pp)
		}
	}
	return p, nil
}
