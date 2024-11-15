package gocbcore

// AgentGroupInternal is a set of internal-only functionality.
// Internal: This should never be used and is not supported.
type AgentGroupInternal struct {
	agentGroup *AgentGroup
}

// Internal creates a new AgentGroupInternal.
// Internal: This should never be used and is not supported.
func (ag *AgentGroup) Internal() *AgentGroupInternal {
	return &AgentGroupInternal{
		agentGroup: ag,
	}
}

// SearchCapabilityStatus returns the current status for a given search capability.
// Internal: This should never be used and is not supported.
func (agi *AgentGroupInternal) SearchCapabilityStatus(cap SearchCapability) CapabilityStatus {
	return agi.agentGroup.clusterAgent.search.capabilityStatus(cap)
}
