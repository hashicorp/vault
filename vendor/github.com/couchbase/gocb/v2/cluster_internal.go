package gocb

import (
	"context"
	"time"
)

// InternalCluster is used for internal functionality.
// Internal: This should never be used and is not supported.
type InternalCluster struct {
	cluster *Cluster
}

// Internal returns an InternalCluster.
// Internal: This should never be used and is not supported.
func (c *Cluster) Internal() *InternalCluster {
	return &InternalCluster{
		cluster: c,
	}
}

// NodeMetadata contains information about a node in the cluster.
// Internal: This should never be used and is not supported.
type NodeMetadata struct {
	ClusterCompatibility int
	ClusterMembership    string
	CouchAPIBase         string
	Hostname             string
	InterestingStats     map[string]float64
	MCDMemoryAllocated   float64
	MCDMemoryReserved    float64
	MemoryFree           float64
	MemoryTotal          float64
	OS                   string
	Ports                map[string]int
	Status               string
	Uptime               int
	Version              string
	ThisNode             bool
}

type jsonClusterCfg struct {
	Nodes []jsonNodeMetadata `json:"nodes"`
}

type jsonNodeMetadata struct {
	ClusterCompatibility int                `json:"clusterCompatibility"`
	ClusterMembership    string             `json:"clusterMembership"`
	CouchAPIBase         string             `json:"couchApiBase"`
	Hostname             string             `json:"hostname"`
	InterestingStats     map[string]float64 `json:"interestingStats,omitempty"`
	MCDMemoryAllocated   float64            `json:"mcdMemoryAllocated"`
	MCDMemoryReserved    float64            `json:"mcdMemoryReserved"`
	MemoryFree           float64            `json:"memoryFree"`
	MemoryTotal          float64            `json:"memoryTotal"`
	OS                   string             `json:"os"`
	Ports                map[string]int     `json:"ports"`
	Status               string             `json:"status"`
	Uptime               int                `json:"uptime,string"`
	Version              string             `json:"version"`
	ThisNode             bool               `json:"thisNode,omitempty"`
}

// GetNodesMetadataOptions is the set of options available to the GetNodesMetadata operation.
// Internal: This should never be used and is not supported.
type GetNodesMetadataOptions struct {
	Timeout       time.Duration
	RetryStrategy RetryStrategy
	ParentSpan    RequestSpan

	// Using a deadlined Context alongside a Timeout will cause the shorter of the two to cause cancellation, this
	// also applies to global level timeouts.
	// UNCOMMITTED: This API may change in the future.
	Context context.Context
}

// GetNodesMetadata returns a list of information about nodes in the cluster.
func (ic *InternalCluster) GetNodesMetadata(opts *GetNodesMetadataOptions) ([]NodeMetadata, error) {
	return autoOpControl(ic.cluster.internalController(), func(provider internalProvider) ([]NodeMetadata, error) {
		if opts == nil {
			opts = &GetNodesMetadataOptions{}
		}

		return provider.GetNodesMetadata(opts)
	})
}
