package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl"
)

const (
	PathPolicyDeny  = "deny"
	PathPolicyRead  = "read"
	PathPolicyWrite = "write"
	PathPolicySudo  = "sudo"
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
	Glob   bool
}

// TakesPrecedence is used when multiple policies
// collide on a path to determine which policy takes
// precendence.
func (p *PathPolicy) TakesPrecedence(other *PathPolicy) bool {
	// Handle the full merge matrix
	switch p.Policy {
	case PathPolicyDeny:
		// Deny always takes precendence
		return true

	case PathPolicyRead:
		// Read never takes precedence
		return false

	case PathPolicyWrite:
		switch other.Policy {
		case PathPolicyRead:
			return true
		case PathPolicyDeny, PathPolicyWrite, PathPolicySudo:
			return false
		default:
			panic("missing case")
		}

	case PathPolicySudo:
		switch other.Policy {
		case PathPolicyRead, PathPolicyWrite:
			return true
		case PathPolicyDeny, PathPolicySudo:
			return false
		default:
			panic("missing case")
		}

	default:
		panic("missing case")
	}
	return false
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
		// Strip the glob character if found
		if strings.HasSuffix(pp.Prefix, "*") {
			pp.Prefix = strings.TrimSuffix(pp.Prefix, "*")
			pp.Glob = true
		}

		// Check the policy is valid
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
