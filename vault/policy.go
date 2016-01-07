package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl"
)

const (
	CreateCapability = "create"
	ReadCapability   = "read"
	UpdateCapability = "update"
	DeleteCapability = "delete"
	ListCapability   = "list"
	DenyCapability   = "deny"
	SudoCapability   = "sudo"

	// Backwards compatibility
	OldDenyPathPolicy  = "deny"
	OldReadPathPolicy  = "read"
	OldWritePathPolicy = "write"
	OldSudoPathPolicy  = "sudo"
)

// Policy is used to represent the policy specified by
// an ACL configuration.
type Policy struct {
	Name  string              `hcl:"name"`
	Paths []*PathCapabilities `hcl:"path,expand"`
	Raw   string
}

// Capability represents a policy for a path in the namespace
type PathCapabilities struct {
	Prefix          string `hcl:",key"`
	Policy          string
	Capabilities    []string
	CapabilitiesMap map[string]bool `hcl:"-"`
	Glob            bool
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
	for _, pc := range p.Paths {
		// Strip the glob character if found
		if strings.HasSuffix(pc.Prefix, "*") {
			pc.Prefix = strings.TrimSuffix(pc.Prefix, "*")
			pc.Glob = true
		}

		// Map old-style policies into capabilities
		switch pc.Policy {
		case OldDenyPathPolicy:
			pc.Capabilities = append(pc.Capabilities, DenyCapability)
		case OldReadPathPolicy:
			pc.Capabilities = append(pc.Capabilities, []string{ReadCapability, ListCapability}...)
		case OldWritePathPolicy:
			pc.Capabilities = append(pc.Capabilities, []string{CreateCapability, ReadCapability, UpdateCapability, DeleteCapability, ListCapability}...)
		case OldSudoPathPolicy:
			pc.Capabilities = append(pc.Capabilities, []string{CreateCapability, ReadCapability, UpdateCapability, DeleteCapability, ListCapability, SudoCapability}...)
		}

		// Initialize the map
		pc.CapabilitiesMap = make(map[string]bool, len(pc.Capabilities))
		for _, cap := range pc.Capabilities {
			switch cap {
			// If it's deny, don't include any other capability
			case DenyCapability:
				pc.Capabilities = []string{DenyCapability}
				pc.CapabilitiesMap = map[string]bool{
					DenyCapability: true,
				}
				goto PathFinished
			case CreateCapability, ReadCapability, UpdateCapability, DeleteCapability, ListCapability, SudoCapability:
				pc.CapabilitiesMap[cap] = true
			default:
				return nil, fmt.Errorf("Invalid capability: %#v", pc)
			}
		}

	PathFinished:
	}
	return p, nil
}
