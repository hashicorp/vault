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

// HasBucketCapabilityStatus verifies whether the specified capability has the specified status.
func (ai *AgentInternal) HasBucketCapabilityStatus(cap BucketCapability, status BucketCapabilityStatus) bool {
	return ai.agent.kvMux.HasBucketCapabilityStatus(cap, status)
}
