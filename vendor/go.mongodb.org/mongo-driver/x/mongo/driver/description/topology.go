// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import (
	"go.mongodb.org/mongo-driver/x/mongo/driver/address"
)

// Topology represents a description of a mongodb topology
type Topology struct {
	Servers               []Server
	Kind                  TopologyKind
	SessionTimeoutMinutes uint32
}

// Server returns the server for the given address. Returns false if the server
// could not be found.
func (t Topology) Server(addr address.Address) (Server, bool) {
	for _, server := range t.Servers {
		if server.Addr.String() == addr.String() {
			return server, true
		}
	}
	return Server{}, false
}

// TopologyDiff is the difference between two different topology descriptions.
type TopologyDiff struct {
	Added   []Server
	Removed []Server
}

// DiffTopology compares the two topology descriptions and returns the difference.
func DiffTopology(old, new Topology) TopologyDiff {
	var diff TopologyDiff

	oldServers := make(map[string]bool)
	for _, s := range old.Servers {
		oldServers[s.Addr.String()] = true
	}

	for _, s := range new.Servers {
		addr := s.Addr.String()
		if oldServers[addr] {
			delete(oldServers, addr)
		} else {
			diff.Added = append(diff.Added, s)
		}
	}

	for _, s := range old.Servers {
		addr := s.Addr.String()
		if oldServers[addr] {
			diff.Removed = append(diff.Removed, s)
		}
	}

	return diff
}

// HostlistDiff is the difference between a topology and a host list.
type HostlistDiff struct {
	Added   []string
	Removed []string
}

// DiffHostlist compares the topology description and host list and returns the difference.
func (t Topology) DiffHostlist(hostlist []string) HostlistDiff {
	var diff HostlistDiff

	oldServers := make(map[string]bool)
	for _, s := range t.Servers {
		oldServers[s.Addr.String()] = true
	}

	for _, addr := range hostlist {
		if oldServers[addr] {
			delete(oldServers, addr)
		} else {
			diff.Added = append(diff.Added, addr)
		}
	}

	for addr := range oldServers {
		diff.Removed = append(diff.Removed, addr)
	}

	return diff
}
