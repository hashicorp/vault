package acl

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/hcl"
)

const (
	// The following levels are the only valid values for the `policy = "read"` stanza.
	// When policies are merged together, the most privilege is granted, except for deny
	// which always takes precedence and supercedes.
	PolicyDeny  = "deny"
	PolicyRead  = "read"
	PolicyWrite = "write"
)

const (
	// The following are the fine-grained capabilities that can be granted within a namespace.
	// The Policy stanza is a short hand for granting several of these. When capabilities are
	// combined we take the union of all capabilities. If the deny capability is present, it
	// takes precedence and overwrites all other capabilities.
	NamespaceCapabilityDeny             = "deny"
	NamespaceCapabilityListJobs         = "list-jobs"
	NamespaceCapabilityReadJob          = "read-job"
	NamespaceCapabilitySubmitJob        = "submit-job"
	NamespaceCapabilityDispatchJob      = "dispatch-job"
	NamespaceCapabilityReadLogs         = "read-logs"
	NamespaceCapabilityReadFS           = "read-fs"
	NamespaceCapabilitySentinelOverride = "sentinel-override"
)

var (
	validNamespace = regexp.MustCompile("^[a-zA-Z0-9-]{1,128}$")
)

// Policy represents a parsed HCL or JSON policy.
type Policy struct {
	Namespaces []*NamespacePolicy `hcl:"namespace,expand"`
	Agent      *AgentPolicy       `hcl:"agent"`
	Node       *NodePolicy        `hcl:"node"`
	Operator   *OperatorPolicy    `hcl:"operator"`
	Quota      *QuotaPolicy       `hcl:"quota"`
	Raw        string             `hcl:"-"`
}

// IsEmpty checks to make sure that at least one policy has been set and is not
// comprised of only a raw policy.
func (p *Policy) IsEmpty() bool {
	return len(p.Namespaces) == 0 &&
		p.Agent == nil &&
		p.Node == nil &&
		p.Operator == nil &&
		p.Quota == nil
}

// NamespacePolicy is the policy for a specific namespace
type NamespacePolicy struct {
	Name         string `hcl:",key"`
	Policy       string
	Capabilities []string
}

type AgentPolicy struct {
	Policy string
}

type NodePolicy struct {
	Policy string
}

type OperatorPolicy struct {
	Policy string
}

type QuotaPolicy struct {
	Policy string
}

// isPolicyValid makes sure the given string matches one of the valid policies.
func isPolicyValid(policy string) bool {
	switch policy {
	case PolicyDeny, PolicyRead, PolicyWrite:
		return true
	default:
		return false
	}
}

// isNamespaceCapabilityValid ensures the given capability is valid for a namespace policy
func isNamespaceCapabilityValid(cap string) bool {
	switch cap {
	case NamespaceCapabilityDeny, NamespaceCapabilityListJobs, NamespaceCapabilityReadJob,
		NamespaceCapabilitySubmitJob, NamespaceCapabilityDispatchJob, NamespaceCapabilityReadLogs,
		NamespaceCapabilityReadFS:
		return true
	// Separate the enterprise-only capabilities
	case NamespaceCapabilitySentinelOverride:
		return true
	default:
		return false
	}
}

// expandNamespacePolicy provides the equivalent set of capabilities for
// a namespace policy
func expandNamespacePolicy(policy string) []string {
	switch policy {
	case PolicyDeny:
		return []string{NamespaceCapabilityDeny}
	case PolicyRead:
		return []string{
			NamespaceCapabilityListJobs,
			NamespaceCapabilityReadJob,
		}
	case PolicyWrite:
		return []string{
			NamespaceCapabilityListJobs,
			NamespaceCapabilityReadJob,
			NamespaceCapabilitySubmitJob,
			NamespaceCapabilityDispatchJob,
			NamespaceCapabilityReadLogs,
			NamespaceCapabilityReadFS,
		}
	default:
		return nil
	}
}

// Parse is used to parse the specified ACL rules into an
// intermediary set of policies, before being compiled into
// the ACL
func Parse(rules string) (*Policy, error) {
	// Decode the rules
	p := &Policy{Raw: rules}
	if rules == "" {
		// Hot path for empty rules
		return p, nil
	}

	// Attempt to parse
	if err := hcl.Decode(p, rules); err != nil {
		return nil, fmt.Errorf("Failed to parse ACL Policy: %v", err)
	}

	// At least one valid policy must be specified, we don't want to store only
	// raw data
	if p.IsEmpty() {
		return nil, fmt.Errorf("Invalid policy: %s", p.Raw)
	}

	// Validate the policy
	for _, ns := range p.Namespaces {
		if !validNamespace.MatchString(ns.Name) {
			return nil, fmt.Errorf("Invalid namespace name: %#v", ns)
		}
		if ns.Policy != "" && !isPolicyValid(ns.Policy) {
			return nil, fmt.Errorf("Invalid namespace policy: %#v", ns)
		}
		for _, cap := range ns.Capabilities {
			if !isNamespaceCapabilityValid(cap) {
				return nil, fmt.Errorf("Invalid namespace capability '%s': %#v", cap, ns)
			}
		}

		// Expand the short hand policy to the capabilities and
		// add to any existing capabilities
		if ns.Policy != "" {
			extraCap := expandNamespacePolicy(ns.Policy)
			ns.Capabilities = append(ns.Capabilities, extraCap...)
		}
	}

	if p.Agent != nil && !isPolicyValid(p.Agent.Policy) {
		return nil, fmt.Errorf("Invalid agent policy: %#v", p.Agent)
	}

	if p.Node != nil && !isPolicyValid(p.Node.Policy) {
		return nil, fmt.Errorf("Invalid node policy: %#v", p.Node)
	}

	if p.Operator != nil && !isPolicyValid(p.Operator.Policy) {
		return nil, fmt.Errorf("Invalid operator policy: %#v", p.Operator)
	}

	if p.Quota != nil && !isPolicyValid(p.Quota.Policy) {
		return nil, fmt.Errorf("Invalid quota policy: %#v", p.Quota)
	}
	return p, nil
}
