// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Topology contains information about a MongoDB cluster.
type Topology struct {
	Servers               []Server
	SetName               string
	Kind                  TopologyKind
	SessionTimeoutMinutes uint32
	CompatibilityErr      error
}

// String implements the Stringer interface.
func (t Topology) String() string {
	var serversStr string
	for _, s := range t.Servers {
		serversStr += "{ " + s.String() + " }, "
	}
	return fmt.Sprintf("Type: %s, Servers: [%s]", t.Kind, serversStr)
}

// Equal compares two topology descriptions and returns true if they are equal.
func (t Topology) Equal(other Topology) bool {
	if t.Kind != other.Kind {
		return false
	}

	topoServers := make(map[string]Server)
	for _, s := range t.Servers {
		topoServers[s.Addr.String()] = s
	}

	otherServers := make(map[string]Server)
	for _, s := range other.Servers {
		otherServers[s.Addr.String()] = s
	}

	if len(topoServers) != len(otherServers) {
		return false
	}

	for _, server := range topoServers {
		otherServer := otherServers[server.Addr.String()]

		if !server.Equal(otherServer) {
			return false
		}
	}

	return true
}

// HasReadableServer returns true if the topology contains a server suitable for reading.
//
// If the Topology's kind is Single or Sharded, the mode parameter is ignored and the function contains true if any of
// the servers in the Topology are of a known type.
//
// For replica sets, the function returns true if the cluster contains a server that matches the provided read
// preference mode.
func (t Topology) HasReadableServer(mode readpref.Mode) bool {
	switch t.Kind {
	case Single, Sharded:
		return hasAvailableServer(t.Servers, 0)
	case ReplicaSetWithPrimary:
		return hasAvailableServer(t.Servers, mode)
	case ReplicaSetNoPrimary, ReplicaSet:
		if mode == readpref.PrimaryMode {
			return false
		}
		// invalid read preference
		if !mode.IsValid() {
			return false
		}

		return hasAvailableServer(t.Servers, mode)
	}
	return false
}

// HasWritableServer returns true if a topology has a server available for writing.
//
// If the Topology's kind is Single or Sharded, this function returns true if any of the servers in the Topology are of
// a known type.
//
// For replica sets, the function returns true if the replica set contains a primary.
func (t Topology) HasWritableServer() bool {
	return t.HasReadableServer(readpref.PrimaryMode)
}

// hasAvailableServer returns true if any servers are available based on the read preference.
func hasAvailableServer(servers []Server, mode readpref.Mode) bool {
	switch mode {
	case readpref.PrimaryMode:
		for _, s := range servers {
			if s.Kind == RSPrimary {
				return true
			}
		}
		return false
	case readpref.PrimaryPreferredMode, readpref.SecondaryPreferredMode, readpref.NearestMode:
		for _, s := range servers {
			if s.Kind == RSPrimary || s.Kind == RSSecondary {
				return true
			}
		}
		return false
	case readpref.SecondaryMode:
		for _, s := range servers {
			if s.Kind == RSSecondary {
				return true
			}
		}
		return false
	}

	// read preference is not specified
	for _, s := range servers {
		switch s.Kind {
		case Standalone,
			RSMember,
			RSPrimary,
			RSSecondary,
			RSArbiter,
			RSGhost,
			Mongos:
			return true
		}
	}

	return false
}
