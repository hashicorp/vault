package gocbcore

// AgentInternal is a set of internal only functionality.
// Internal: This should never be used and is not supported.
type AgentInternal struct {
	agent *Agent
}

// Internal creates a new AgentInternal.
// Internal: This should never be used and is not supported.
func (agent *Agent) Internal() *AgentInternal {
	return &AgentInternal{
		agent: agent,
	}
}

// BucketCapabilityStatus returns the current status for a given bucket capability.
func (ai *AgentInternal) BucketCapabilityStatus(cap BucketCapability) CapabilityStatus {
	return ai.agent.kvMux.BucketCapabilityStatus(cap)
}
