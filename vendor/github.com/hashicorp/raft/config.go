package raft

import (
	"fmt"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
)

// ProtocolVersion is the version of the protocol (which includes RPC messages
// as well as Raft-specific log entries) that this server can _understand_. Use
// the ProtocolVersion member of the Config object to control the version of
// the protocol to use when _speaking_ to other servers. Note that depending on
// the protocol version being spoken, some otherwise understood RPC messages
// may be refused. See dispositionRPC for details of this logic.
//
// There are notes about the upgrade path in the description of the versions
// below. If you are starting a fresh cluster then there's no reason not to
// jump right to the latest protocol version. If you need to interoperate with
// older, version 0 Raft servers you'll need to drive the cluster through the
// different versions in order.
//
// The version details are complicated, but here's a summary of what's required
// to get from a version 0 cluster to version 3:
//
// 1. In version N of your app that starts using the new Raft library with
//    versioning, set ProtocolVersion to 1.
// 2. Make version N+1 of your app require version N as a prerequisite (all
//    servers must be upgraded). For version N+1 of your app set ProtocolVersion
//    to 2.
// 3. Similarly, make version N+2 of your app require version N+1 as a
//    prerequisite. For version N+2 of your app, set ProtocolVersion to 3.
//
// During this upgrade, older cluster members will still have Server IDs equal
// to their network addresses. To upgrade an older member and give it an ID, it
// needs to leave the cluster and re-enter:
//
// 1. Remove the server from the cluster with RemoveServer, using its network
//    address as its ServerID.
// 2. Update the server's config to use a UUID or something else that is
//	  not tied to the machine as the ServerID (restarting the server).
// 3. Add the server back to the cluster with AddVoter, using its new ID.
//
// You can do this during the rolling upgrade from N+1 to N+2 of your app, or
// as a rolling change at any time after the upgrade.
//
// Version History
//
// 0: Original Raft library before versioning was added. Servers running this
//    version of the Raft library use AddPeerDeprecated/RemovePeerDeprecated
//    for all configuration changes, and have no support for LogConfiguration.
// 1: First versioned protocol, used to interoperate with old servers, and begin
//    the migration path to newer versions of the protocol. Under this version
//    all configuration changes are propagated using the now-deprecated
//    RemovePeerDeprecated Raft log entry. This means that server IDs are always
//    set to be the same as the server addresses (since the old log entry type
//    cannot transmit an ID), and only AddPeer/RemovePeer APIs are supported.
//    Servers running this version of the protocol can understand the new
//    LogConfiguration Raft log entry but will never generate one so they can
//    remain compatible with version 0 Raft servers in the cluster.
// 2: Transitional protocol used when migrating an existing cluster to the new
//    server ID system. Server IDs are still set to be the same as server
//    addresses, but all configuration changes are propagated using the new
//    LogConfiguration Raft log entry type, which can carry full ID information.
//    This version supports the old AddPeer/RemovePeer APIs as well as the new
//    ID-based AddVoter/RemoveServer APIs which should be used when adding
//    version 3 servers to the cluster later. This version sheds all
//    interoperability with version 0 servers, but can interoperate with newer
//    Raft servers running with protocol version 1 since they can understand the
//    new LogConfiguration Raft log entry, and this version can still understand
//    their RemovePeerDeprecated Raft log entries. We need this protocol version
//    as an intermediate step between 1 and 3 so that servers will propagate the
//    ID information that will come from newly-added (or -rolled) servers using
//    protocol version 3, but since they are still using their address-based IDs
//    from the previous step they will still be able to track commitments and
//    their own voting status properly. If we skipped this step, servers would
//    be started with their new IDs, but they wouldn't see themselves in the old
//    address-based configuration, so none of the servers would think they had a
//    vote.
// 3: Protocol adding full support for server IDs and new ID-based server APIs
//    (AddVoter, AddNonvoter, etc.), old AddPeer/RemovePeer APIs are no longer
//    supported. Version 2 servers should be swapped out by removing them from
//    the cluster one-by-one and re-adding them with updated configuration for
//    this protocol version, along with their server ID. The remove/add cycle
//    is required to populate their server ID. Note that removing must be done
//    by ID, which will be the old server's address.
type ProtocolVersion int

const (
	// ProtocolVersionMin is the minimum protocol version
	ProtocolVersionMin ProtocolVersion = 0
	// ProtocolVersionMax is the maximum protocol version
	ProtocolVersionMax = 3
)

// SnapshotVersion is the version of snapshots that this server can understand.
// Currently, it is always assumed that the server generates the latest version,
// though this may be changed in the future to include a configurable version.
//
// Version History
//
// 0: Original Raft library before versioning was added. The peers portion of
//    these snapshots is encoded in the legacy format which requires decodePeers
//    to parse. This version of snapshots should only be produced by the
//    unversioned Raft library.
// 1: New format which adds support for a full configuration structure and its
//    associated log index, with support for server IDs and non-voting server
//    modes. To ease upgrades, this also includes the legacy peers structure but
//    that will never be used by servers that understand version 1 snapshots.
//    Since the original Raft library didn't enforce any versioning, we must
//    include the legacy peers structure for this version, but we can deprecate
//    it in the next snapshot version.
type SnapshotVersion int

const (
	// SnapshotVersionMin is the minimum snapshot version
	SnapshotVersionMin SnapshotVersion = 0
	// SnapshotVersionMax is the maximum snapshot version
	SnapshotVersionMax = 1
)

// Config provides any necessary configuration for the Raft server.
type Config struct {
	// ProtocolVersion allows a Raft server to inter-operate with older
	// Raft servers running an older version of the code. This is used to
	// version the wire protocol as well as Raft-specific log entries that
	// the server uses when _speaking_ to other servers. There is currently
	// no auto-negotiation of versions so all servers must be manually
	// configured with compatible versions. See ProtocolVersionMin and
	// ProtocolVersionMax for the versions of the protocol that this server
	// can _understand_.
	ProtocolVersion ProtocolVersion

	// HeartbeatTimeout specifies the time in follower state without
	// a leader before we attempt an election.
	HeartbeatTimeout time.Duration

	// ElectionTimeout specifies the time in candidate state without
	// a leader before we attempt an election.
	ElectionTimeout time.Duration

	// CommitTimeout controls the time without an Apply() operation
	// before we heartbeat to ensure a timely commit. Due to random
	// staggering, may be delayed as much as 2x this value.
	CommitTimeout time.Duration

	// MaxAppendEntries controls the maximum number of append entries
	// to send at once. We want to strike a balance between efficiency
	// and avoiding waste if the follower is going to reject because of
	// an inconsistent log.
	MaxAppendEntries int

	// BatchApplyCh indicates whether we should buffer applyCh
	// to size MaxAppendEntries. This enables batch log commitment,
	// but breaks the timeout guarantee on Apply. Specifically,
	// a log can be added to the applyCh buffer but not actually be
	// processed until after the specified timeout.
	BatchApplyCh bool

	// If we are a member of a cluster, and RemovePeer is invoked for the
	// local node, then we forget all peers and transition into the follower state.
	// If ShutdownOnRemove is set, we additional shutdown Raft. Otherwise,
	// we can become a leader of a cluster containing only this node.
	ShutdownOnRemove bool

	// TrailingLogs controls how many logs we leave after a snapshot. This is used
	// so that we can quickly replay logs on a follower instead of being forced to
	// send an entire snapshot. The value passed here is the initial setting used.
	// This can be tuned during operation using ReloadConfig.
	TrailingLogs uint64

	// SnapshotInterval controls how often we check if we should perform a
	// snapshot. We randomly stagger between this value and 2x this value to avoid
	// the entire cluster from performing a snapshot at once. The value passed
	// here is the initial setting used. This can be tuned during operation using
	// ReloadConfig.
	SnapshotInterval time.Duration

	// SnapshotThreshold controls how many outstanding logs there must be before
	// we perform a snapshot. This is to prevent excessive snapshotting by
	// replaying a small set of logs instead. The value passed here is the initial
	// setting used. This can be tuned during operation using ReloadConfig.
	SnapshotThreshold uint64

	// LeaderLeaseTimeout is used to control how long the "lease" lasts
	// for being the leader without being able to contact a quorum
	// of nodes. If we reach this interval without contact, we will
	// step down as leader.
	LeaderLeaseTimeout time.Duration

	// LocalID is a unique ID for this server across all time. When running with
	// ProtocolVersion < 3, you must set this to be the same as the network
	// address of your transport.
	LocalID ServerID

	// NotifyCh is used to provide a channel that will be notified of leadership
	// changes. Raft will block writing to this channel, so it should either be
	// buffered or aggressively consumed.
	NotifyCh chan<- bool

	// LogOutput is used as a sink for logs, unless Logger is specified.
	// Defaults to os.Stderr.
	LogOutput io.Writer

	// LogLevel represents a log level. If the value does not match a known
	// logging level hclog.NoLevel is used.
	LogLevel string

	// Logger is a user-provided logger. If nil, a logger writing to
	// LogOutput with LogLevel is used.
	Logger hclog.Logger

	// NoSnapshotRestoreOnStart controls if raft will restore a snapshot to the
	// FSM on start. This is useful if your FSM recovers from other mechanisms
	// than raft snapshotting. Snapshot metadata will still be used to initialize
	// raft's configuration and index values.
	NoSnapshotRestoreOnStart bool

	// skipStartup allows NewRaft() to bypass all background work goroutines
	skipStartup bool
}

// ReloadableConfig is the subset of Config that may be reconfigured during
// runtime using raft.ReloadConfig. We choose to duplicate fields over embedding
// or accepting a Config but only using specific fields to keep the API clear.
// Reconfiguring some fields is potentially dangerous so we should only
// selectively enable it for fields where that is allowed.
type ReloadableConfig struct {
	// TrailingLogs controls how many logs we leave after a snapshot. This is used
	// so that we can quickly replay logs on a follower instead of being forced to
	// send an entire snapshot. The value passed here updates the setting at runtime
	// which will take effect as soon as the next snapshot completes and truncation
	// occurs.
	TrailingLogs uint64

	// SnapshotInterval controls how often we check if we should perform a snapshot.
	// We randomly stagger between this value and 2x this value to avoid the entire
	// cluster from performing a snapshot at once.
	SnapshotInterval time.Duration

	// SnapshotThreshold controls how many outstanding logs there must be before
	// we perform a snapshot. This is to prevent excessive snapshots when we can
	// just replay a small set of logs.
	SnapshotThreshold uint64
}

// apply sets the reloadable fields on the passed Config to the values in
// `ReloadableConfig`. It returns a copy of Config with the fields from this
// ReloadableConfig set.
func (rc *ReloadableConfig) apply(to Config) Config {
	to.TrailingLogs = rc.TrailingLogs
	to.SnapshotInterval = rc.SnapshotInterval
	to.SnapshotThreshold = rc.SnapshotThreshold
	return to
}

// fromConfig copies the reloadable fields from the passed Config.
func (rc *ReloadableConfig) fromConfig(from Config) {
	rc.TrailingLogs = from.TrailingLogs
	rc.SnapshotInterval = from.SnapshotInterval
	rc.SnapshotThreshold = from.SnapshotThreshold
}

// DefaultConfig returns a Config with usable defaults.
func DefaultConfig() *Config {
	return &Config{
		ProtocolVersion:    ProtocolVersionMax,
		HeartbeatTimeout:   1000 * time.Millisecond,
		ElectionTimeout:    1000 * time.Millisecond,
		CommitTimeout:      50 * time.Millisecond,
		MaxAppendEntries:   64,
		ShutdownOnRemove:   true,
		TrailingLogs:       10240,
		SnapshotInterval:   120 * time.Second,
		SnapshotThreshold:  8192,
		LeaderLeaseTimeout: 500 * time.Millisecond,
		LogLevel:           "DEBUG",
	}
}

// ValidateConfig is used to validate a sane configuration
func ValidateConfig(config *Config) error {
	// We don't actually support running as 0 in the library any more, but
	// we do understand it.
	protocolMin := ProtocolVersionMin
	if protocolMin == 0 {
		protocolMin = 1
	}
	if config.ProtocolVersion < protocolMin ||
		config.ProtocolVersion > ProtocolVersionMax {
		return fmt.Errorf("ProtocolVersion %d must be >= %d and <= %d",
			config.ProtocolVersion, protocolMin, ProtocolVersionMax)
	}
	if len(config.LocalID) == 0 {
		return fmt.Errorf("LocalID cannot be empty")
	}
	if config.HeartbeatTimeout < 5*time.Millisecond {
		return fmt.Errorf("HeartbeatTimeout is too low")
	}
	if config.ElectionTimeout < 5*time.Millisecond {
		return fmt.Errorf("ElectionTimeout is too low")
	}
	if config.CommitTimeout < time.Millisecond {
		return fmt.Errorf("CommitTimeout is too low")
	}
	if config.MaxAppendEntries <= 0 {
		return fmt.Errorf("MaxAppendEntries must be positive")
	}
	if config.MaxAppendEntries > 1024 {
		return fmt.Errorf("MaxAppendEntries is too large")
	}
	if config.SnapshotInterval < 5*time.Millisecond {
		return fmt.Errorf("SnapshotInterval is too low")
	}
	if config.LeaderLeaseTimeout < 5*time.Millisecond {
		return fmt.Errorf("LeaderLeaseTimeout is too low")
	}
	if config.LeaderLeaseTimeout > config.HeartbeatTimeout {
		return fmt.Errorf("LeaderLeaseTimeout cannot be larger than heartbeat timeout")
	}
	if config.ElectionTimeout < config.HeartbeatTimeout {
		return fmt.Errorf("ElectionTimeout must be equal or greater than Heartbeat Timeout")
	}
	return nil
}
