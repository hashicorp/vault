package structs

import (
	"time"

	"github.com/hashicorp/raft"
)

// RaftServer has information about a server in the Raft configuration.
type RaftServer struct {
	// ID is the unique ID for the server. These are currently the same
	// as the address, but they will be changed to a real GUID in a future
	// release of Nomad.
	ID raft.ServerID

	// Node is the node name of the server, as known by Nomad, or this
	// will be set to "(unknown)" otherwise.
	Node string

	// Address is the IP:port of the server, used for Raft communications.
	Address raft.ServerAddress

	// Leader is true if this server is the current cluster leader.
	Leader bool

	// Voter is true if this server has a vote in the cluster. This might
	// be false if the server is staging and still coming online, or if
	// it's a non-voting server, which will be added in a future release of
	// Nomad.
	Voter bool

	// RaftProtocol is the version of the Raft protocol spoken by this server.
	RaftProtocol string
}

// RaftConfigurationResponse is returned when querying for the current Raft
// configuration.
type RaftConfigurationResponse struct {
	// Servers has the list of servers in the Raft configuration.
	Servers []*RaftServer

	// Index has the Raft index of this configuration.
	Index uint64
}

// RaftPeerByAddressRequest is used by the Operator endpoint to apply a Raft
// operation on a specific Raft peer by address in the form of "IP:port".
type RaftPeerByAddressRequest struct {
	// Address is the peer to remove, in the form "IP:port".
	Address raft.ServerAddress

	// WriteRequest holds the Region for this request.
	WriteRequest
}

// RaftPeerByIDRequest is used by the Operator endpoint to apply a Raft
// operation on a specific Raft peer by ID.
type RaftPeerByIDRequest struct {
	// ID is the peer ID to remove.
	ID raft.ServerID

	// WriteRequest holds the Region for this request.
	WriteRequest
}

// AutopilotSetConfigRequest is used by the Operator endpoint to update the
// current Autopilot configuration of the cluster.
type AutopilotSetConfigRequest struct {
	// Datacenter is the target this request is intended for.
	Datacenter string

	// Config is the new Autopilot configuration to use.
	Config AutopilotConfig

	// CAS controls whether to use check-and-set semantics for this request.
	CAS bool

	// WriteRequest holds the ACL token to go along with this request.
	WriteRequest
}

// RequestDatacenter returns the datacenter for a given request.
func (op *AutopilotSetConfigRequest) RequestDatacenter() string {
	return op.Datacenter
}

// AutopilotConfig is the internal config for the Autopilot mechanism.
type AutopilotConfig struct {
	// CleanupDeadServers controls whether to remove dead servers when a new
	// server is added to the Raft peers.
	CleanupDeadServers bool

	// ServerStabilizationTime is the minimum amount of time a server must be
	// in a stable, healthy state before it can be added to the cluster. Only
	// applicable with Raft protocol version 3 or higher.
	ServerStabilizationTime time.Duration

	// LastContactThreshold is the limit on the amount of time a server can go
	// without leader contact before being considered unhealthy.
	LastContactThreshold time.Duration

	// MaxTrailingLogs is the amount of entries in the Raft Log that a server can
	// be behind before being considered unhealthy.
	MaxTrailingLogs uint64

	// (Enterprise-only) EnableRedundancyZones specifies whether to enable redundancy zones.
	EnableRedundancyZones bool

	// (Enterprise-only) DisableUpgradeMigration will disable Autopilot's upgrade migration
	// strategy of waiting until enough newer-versioned servers have been added to the
	// cluster before promoting them to voters.
	DisableUpgradeMigration bool

	// (Enterprise-only) EnableCustomUpgrades specifies whether to enable using custom
	// upgrade versions when performing migrations.
	EnableCustomUpgrades bool

	// CreateIndex/ModifyIndex store the create/modify indexes of this configuration.
	CreateIndex uint64
	ModifyIndex uint64
}
