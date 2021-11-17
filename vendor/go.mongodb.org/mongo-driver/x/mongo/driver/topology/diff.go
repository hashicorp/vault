// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import "go.mongodb.org/mongo-driver/mongo/description"

// hostlistDiff is the difference between a topology and a host list.
type hostlistDiff struct {
	Added   []string
	Removed []string
}

// diffHostList compares the topology description and host list and returns the difference.
func diffHostList(t description.Topology, hostlist []string) hostlistDiff {
	var diff hostlistDiff

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

// topologyDiff is the difference between two different topology descriptions.
type topologyDiff struct {
	Added   []description.Server
	Removed []description.Server
}

// diffTopology compares the two topology descriptions and returns the difference.
func diffTopology(old, new description.Topology) topologyDiff {
	var diff topologyDiff

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
