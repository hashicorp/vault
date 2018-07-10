package acl

import (
	"fmt"

	iradix "github.com/hashicorp/go-immutable-radix"
)

// ManagementACL is a singleton used for management tokens
var ManagementACL *ACL

func init() {
	var err error
	ManagementACL, err = NewACL(true, nil)
	if err != nil {
		panic(fmt.Errorf("failed to setup management ACL: %v", err))
	}
}

// capabilitySet is a type wrapper to help managing a set of capabilities
type capabilitySet map[string]struct{}

func (c capabilitySet) Check(k string) bool {
	_, ok := c[k]
	return ok
}

func (c capabilitySet) Set(k string) {
	c[k] = struct{}{}
}

func (c capabilitySet) Clear() {
	for cap := range c {
		delete(c, cap)
	}
}

// ACL object is used to convert a set of policies into a structure that
// can be efficiently evaluated to determine if an action is allowed.
type ACL struct {
	// management tokens are allowed to do anything
	management bool

	// namespaces maps a namespace to a capabilitySet
	namespaces *iradix.Tree

	agent    string
	node     string
	operator string
	quota    string
}

// maxPrivilege returns the policy which grants the most privilege
// This handles the case of Deny always taking maximum precedence.
func maxPrivilege(a, b string) string {
	switch {
	case a == PolicyDeny || b == PolicyDeny:
		return PolicyDeny
	case a == PolicyWrite || b == PolicyWrite:
		return PolicyWrite
	case a == PolicyRead || b == PolicyRead:
		return PolicyRead
	default:
		return ""
	}
}

// NewACL compiles a set of policies into an ACL object
func NewACL(management bool, policies []*Policy) (*ACL, error) {
	// Hot-path management tokens
	if management {
		return &ACL{management: true}, nil
	}

	// Create the ACL object
	acl := &ACL{}
	nsTxn := iradix.New().Txn()

	for _, policy := range policies {
	NAMESPACES:
		for _, ns := range policy.Namespaces {
			// Check for existing capabilities
			var capabilities capabilitySet
			raw, ok := nsTxn.Get([]byte(ns.Name))
			if ok {
				capabilities = raw.(capabilitySet)
			} else {
				capabilities = make(capabilitySet)
				nsTxn.Insert([]byte(ns.Name), capabilities)
			}

			// Deny always takes precedence
			if capabilities.Check(NamespaceCapabilityDeny) {
				continue NAMESPACES
			}

			// Add in all the capabilities
			for _, cap := range ns.Capabilities {
				if cap == NamespaceCapabilityDeny {
					// Overwrite any existing capabilities
					capabilities.Clear()
					capabilities.Set(NamespaceCapabilityDeny)
					continue NAMESPACES
				}
				capabilities.Set(cap)
			}
		}

		// Take the maximum privilege for agent, node, and operator
		if policy.Agent != nil {
			acl.agent = maxPrivilege(acl.agent, policy.Agent.Policy)
		}
		if policy.Node != nil {
			acl.node = maxPrivilege(acl.node, policy.Node.Policy)
		}
		if policy.Operator != nil {
			acl.operator = maxPrivilege(acl.operator, policy.Operator.Policy)
		}
		if policy.Quota != nil {
			acl.quota = maxPrivilege(acl.quota, policy.Quota.Policy)
		}
	}

	// Finalize the namespaces
	acl.namespaces = nsTxn.Commit()
	return acl, nil
}

// AllowNsOp is shorthand for AllowNamespaceOperation
func (a *ACL) AllowNsOp(ns string, op string) bool {
	return a.AllowNamespaceOperation(ns, op)
}

// AllowNamespaceOperation checks if a given operation is allowed for a namespace
func (a *ACL) AllowNamespaceOperation(ns string, op string) bool {
	// Hot path management tokens
	if a.management {
		return true
	}

	// Check for a matching capability set
	raw, ok := a.namespaces.Get([]byte(ns))
	if !ok {
		return false
	}

	// Check if the capability has been granted
	capabilities := raw.(capabilitySet)
	return capabilities.Check(op)
}

// AllowNamespace checks if any operations are allowed for a namespace
func (a *ACL) AllowNamespace(ns string) bool {
	// Hot path management tokens
	if a.management {
		return true
	}

	// Check for a matching capability set
	raw, ok := a.namespaces.Get([]byte(ns))
	if !ok {
		return false
	}

	// Check if the capability has been granted
	capabilities := raw.(capabilitySet)
	if len(capabilities) == 0 {
		return false
	}

	return !capabilities.Check(PolicyDeny)
}

// AllowAgentRead checks if read operations are allowed for an agent
func (a *ACL) AllowAgentRead() bool {
	switch {
	case a.management:
		return true
	case a.agent == PolicyWrite:
		return true
	case a.agent == PolicyRead:
		return true
	default:
		return false
	}
}

// AllowAgentWrite checks if write operations are allowed for an agent
func (a *ACL) AllowAgentWrite() bool {
	switch {
	case a.management:
		return true
	case a.agent == PolicyWrite:
		return true
	default:
		return false
	}
}

// AllowNodeRead checks if read operations are allowed for a node
func (a *ACL) AllowNodeRead() bool {
	switch {
	case a.management:
		return true
	case a.node == PolicyWrite:
		return true
	case a.node == PolicyRead:
		return true
	default:
		return false
	}
}

// AllowNodeWrite checks if write operations are allowed for a node
func (a *ACL) AllowNodeWrite() bool {
	switch {
	case a.management:
		return true
	case a.node == PolicyWrite:
		return true
	default:
		return false
	}
}

// AllowOperatorRead checks if read operations are allowed for a operator
func (a *ACL) AllowOperatorRead() bool {
	switch {
	case a.management:
		return true
	case a.operator == PolicyWrite:
		return true
	case a.operator == PolicyRead:
		return true
	default:
		return false
	}
}

// AllowOperatorWrite checks if write operations are allowed for a operator
func (a *ACL) AllowOperatorWrite() bool {
	switch {
	case a.management:
		return true
	case a.operator == PolicyWrite:
		return true
	default:
		return false
	}
}

// AllowQuotaRead checks if read operations are allowed for all quotas
func (a *ACL) AllowQuotaRead() bool {
	switch {
	case a.management:
		return true
	case a.quota == PolicyWrite:
		return true
	case a.quota == PolicyRead:
		return true
	default:
		return false
	}
}

// AllowQuotaWrite checks if write operations are allowed for quotas
func (a *ACL) AllowQuotaWrite() bool {
	switch {
	case a.management:
		return true
	case a.quota == PolicyWrite:
		return true
	default:
		return false
	}
}

// IsManagement checks if this represents a management token
func (a *ACL) IsManagement() bool {
	return a.management
}
