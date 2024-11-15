// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package raft

import "fmt"

// ServerSuffrage determines whether a Server in a Configuration gets a vote.
type ServerSuffrage int

// Note: Don't renumber these, since the numbers are written into the log.
const (
	// Voter is a server whose vote is counted in elections and whose match index
	// is used in advancing the leader's commit index.
	Voter ServerSuffrage = iota
	// Nonvoter is a server that receives log entries but is not considered for
	// elections or commitment purposes.
	Nonvoter
	// Staging is a server that acts like a Nonvoter. A configuration change
	// with a ConfigurationChangeCommand of Promote can change a Staging server
	// into a Voter.
	// Deprecated: use Nonvoter instead.
	Staging
)

func (s ServerSuffrage) String() string {
	switch s {
	case Voter:
		return "Voter"
	case Nonvoter:
		return "Nonvoter"
	case Staging:
		return "Staging"
	}
	return "ServerSuffrage"
}

// ConfigurationStore provides an interface that can optionally be implemented by FSMs
// to store configuration updates made in the replicated log. In general this is only
// necessary for FSMs that mutate durable state directly instead of applying changes
// in memory and snapshotting periodically. By storing configuration changes, the
// persistent FSM state can behave as a complete snapshot, and be able to recover
// without an external snapshot just for persisting the raft configuration.
type ConfigurationStore interface {
	// ConfigurationStore is a superset of the FSM functionality
	FSM

	// StoreConfiguration is invoked once a log entry containing a configuration
	// change is committed. It takes the index at which the configuration was
	// written and the configuration value.
	StoreConfiguration(index uint64, configuration Configuration)
}

type nopConfigurationStore struct{}

func (s nopConfigurationStore) StoreConfiguration(_ uint64, _ Configuration) {}

// ServerID is a unique string identifying a server for all time.
type ServerID string

// ServerAddress is a network address for a server that a transport can contact.
type ServerAddress string

// Server tracks the information about a single server in a configuration.
type Server struct {
	// Suffrage determines whether the server gets a vote.
	Suffrage ServerSuffrage
	// ID is a unique string identifying this server for all time.
	ID ServerID
	// Address is its network address that a transport can contact.
	Address ServerAddress
}

// Configuration tracks which servers are in the cluster, and whether they have
// votes. This should include the local server, if it's a member of the cluster.
// The servers are listed no particular order, but each should only appear once.
// These entries are appended to the log during membership changes.
type Configuration struct {
	Servers []Server
}

// Clone makes a deep copy of a Configuration.
func (c *Configuration) Clone() (copy Configuration) {
	copy.Servers = append(copy.Servers, c.Servers...)
	return
}

// ConfigurationChangeCommand is the different ways to change the cluster
// configuration.
type ConfigurationChangeCommand uint8

const (
	// AddVoter adds a server with Suffrage of Voter.
	AddVoter ConfigurationChangeCommand = iota
	// AddNonvoter makes a server Nonvoter unless its Staging or Voter.
	AddNonvoter
	// DemoteVoter makes a server Nonvoter unless its absent.
	DemoteVoter
	// RemoveServer removes a server entirely from the cluster membership.
	RemoveServer
	// Promote changes a server from Staging to Voter. The command will be a
	// no-op if the server is not Staging.
	// Deprecated: use AddVoter instead.
	Promote
	// AddStaging makes a server a Voter.
	// Deprecated: AddStaging was actually AddVoter. Use AddVoter instead.
	AddStaging = 0 // explicit 0 to preserve the old value.
)

func (c ConfigurationChangeCommand) String() string {
	switch c {
	case AddVoter:
		return "AddVoter"
	case AddNonvoter:
		return "AddNonvoter"
	case DemoteVoter:
		return "DemoteVoter"
	case RemoveServer:
		return "RemoveServer"
	case Promote:
		return "Promote"
	}
	return "ConfigurationChangeCommand"
}

// configurationChangeRequest describes a change that a leader would like to
// make to its current configuration. It's used only within a single server
// (never serialized into the log), as part of `configurationChangeFuture`.
type configurationChangeRequest struct {
	command       ConfigurationChangeCommand
	serverID      ServerID
	serverAddress ServerAddress // only present for AddVoter, AddNonvoter
	// prevIndex, if nonzero, is the index of the only configuration upon which
	// this change may be applied; if another configuration entry has been
	// added in the meantime, this request will fail.
	prevIndex uint64
}

// configurations is state tracked on every server about its Configurations.
// Note that, per Diego's dissertation, there can be at most one uncommitted
// configuration at a time (the next configuration may not be created until the
// prior one has been committed).
//
// One downside to storing just two configurations is that if you try to take a
// snapshot when your state machine hasn't yet applied the committedIndex, we
// have no record of the configuration that would logically fit into that
// snapshot. We disallow snapshots in that case now. An alternative approach,
// which LogCabin uses, is to track every configuration change in the
// log.
type configurations struct {
	// committed is the latest configuration in the log/snapshot that has been
	// committed (the one with the largest index).
	committed Configuration
	// committedIndex is the log index where 'committed' was written.
	committedIndex uint64
	// latest is the latest configuration in the log/snapshot (may be committed
	// or uncommitted)
	latest Configuration
	// latestIndex is the log index where 'latest' was written.
	latestIndex uint64
}

// Clone makes a deep copy of a configurations object.
func (c *configurations) Clone() (copy configurations) {
	copy.committed = c.committed.Clone()
	copy.committedIndex = c.committedIndex
	copy.latest = c.latest.Clone()
	copy.latestIndex = c.latestIndex
	return
}

// hasVote returns true if the server identified by 'id' is a Voter in the
// provided Configuration.
func hasVote(configuration Configuration, id ServerID) bool {
	for _, server := range configuration.Servers {
		if server.ID == id {
			return server.Suffrage == Voter
		}
	}
	return false
}

// inConfiguration returns true if the server identified by 'id' is in in the
// provided Configuration.
func inConfiguration(configuration Configuration, id ServerID) bool {
	for _, server := range configuration.Servers {
		if server.ID == id {
			return true
		}
	}
	return false
}

// checkConfiguration tests a cluster membership configuration for common
// errors.
func checkConfiguration(configuration Configuration) error {
	idSet := make(map[ServerID]bool)
	addressSet := make(map[ServerAddress]bool)
	var voters int
	for _, server := range configuration.Servers {
		if server.ID == "" {
			return fmt.Errorf("empty ID in configuration: %v", configuration)
		}
		if server.Address == "" {
			return fmt.Errorf("empty address in configuration: %v", server)
		}
		if idSet[server.ID] {
			return fmt.Errorf("found duplicate ID in configuration: %v", server.ID)
		}
		idSet[server.ID] = true
		if addressSet[server.Address] {
			return fmt.Errorf("found duplicate address in configuration: %v", server.Address)
		}
		addressSet[server.Address] = true
		if server.Suffrage == Voter {
			voters++
		}
	}
	if voters == 0 {
		return fmt.Errorf("need at least one voter in configuration: %v", configuration)
	}
	return nil
}

// nextConfiguration generates a new Configuration from the current one and a
// configuration change request. It's split from appendConfigurationEntry so
// that it can be unit tested easily.
func nextConfiguration(current Configuration, currentIndex uint64, change configurationChangeRequest) (Configuration, error) {
	if change.prevIndex > 0 && change.prevIndex != currentIndex {
		return Configuration{}, fmt.Errorf("configuration changed since %v (latest is %v)", change.prevIndex, currentIndex)
	}

	configuration := current.Clone()
	switch change.command {
	case AddVoter:
		newServer := Server{
			Suffrage: Voter,
			ID:       change.serverID,
			Address:  change.serverAddress,
		}
		found := false
		for i, server := range configuration.Servers {
			if server.ID == change.serverID {
				if server.Suffrage == Voter {
					configuration.Servers[i].Address = change.serverAddress
				} else {
					configuration.Servers[i] = newServer
				}
				found = true
				break
			}
		}
		if !found {
			configuration.Servers = append(configuration.Servers, newServer)
		}
	case AddNonvoter:
		newServer := Server{
			Suffrage: Nonvoter,
			ID:       change.serverID,
			Address:  change.serverAddress,
		}
		found := false
		for i, server := range configuration.Servers {
			if server.ID == change.serverID {
				if server.Suffrage != Nonvoter {
					configuration.Servers[i].Address = change.serverAddress
				} else {
					configuration.Servers[i] = newServer
				}
				found = true
				break
			}
		}
		if !found {
			configuration.Servers = append(configuration.Servers, newServer)
		}
	case DemoteVoter:
		for i, server := range configuration.Servers {
			if server.ID == change.serverID {
				configuration.Servers[i].Suffrage = Nonvoter
				break
			}
		}
	case RemoveServer:
		for i, server := range configuration.Servers {
			if server.ID == change.serverID {
				configuration.Servers = append(configuration.Servers[:i], configuration.Servers[i+1:]...)
				break
			}
		}
	case Promote:
		for i, server := range configuration.Servers {
			if server.ID == change.serverID && server.Suffrage == Staging {
				configuration.Servers[i].Suffrage = Voter
				break
			}
		}
	}

	// Make sure we didn't do something bad like remove the last voter
	if err := checkConfiguration(configuration); err != nil {
		return Configuration{}, err
	}

	return configuration, nil
}

// encodePeers is used to serialize a Configuration into the old peers format.
// This is here for backwards compatibility when operating with a mix of old
// servers and should be removed once we deprecate support for protocol version 1.
func encodePeers(configuration Configuration, trans Transport) []byte {
	// Gather up all the voters, other suffrage types are not supported by
	// this data format.
	var encPeers [][]byte
	for _, server := range configuration.Servers {
		if server.Suffrage == Voter {
			encPeers = append(encPeers, trans.EncodePeer(server.ID, server.Address))
		}
	}

	// Encode the entire array.
	buf, err := encodeMsgPack(encPeers)
	if err != nil {
		panic(fmt.Errorf("failed to encode peers: %v", err))
	}

	return buf.Bytes()
}

// decodePeers is used to deserialize an old list of peers into a Configuration.
// This is here for backwards compatibility with old log entries and snapshots;
// it should be removed eventually.
func decodePeers(buf []byte, trans Transport) (Configuration, error) {
	// Decode the buffer first.
	var encPeers [][]byte
	if err := decodeMsgPack(buf, &encPeers); err != nil {
		return Configuration{}, fmt.Errorf("failed to decode peers: %v", err)
	}

	// Deserialize each peer.
	var servers []Server
	for _, enc := range encPeers {
		p := trans.DecodePeer(enc)
		servers = append(servers, Server{
			Suffrage: Voter,
			ID:       ServerID(p),
			Address:  p,
		})
	}

	return Configuration{Servers: servers}, nil
}

// EncodeConfiguration serializes a Configuration using MsgPack, or panics on
// errors.
func EncodeConfiguration(configuration Configuration) []byte {
	buf, err := encodeMsgPack(configuration)
	if err != nil {
		panic(fmt.Errorf("failed to encode configuration: %v", err))
	}
	return buf.Bytes()
}

// DecodeConfiguration deserializes a Configuration using MsgPack, or panics on
// errors.
func DecodeConfiguration(buf []byte) Configuration {
	var configuration Configuration
	if err := decodeMsgPack(buf, &configuration); err != nil {
		panic(fmt.Errorf("failed to decode configuration: %v", err))
	}
	return configuration
}
