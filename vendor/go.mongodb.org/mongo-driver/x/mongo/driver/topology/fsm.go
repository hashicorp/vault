// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"bytes"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/mongo/driver/address"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
)

var supportedWireVersions = description.NewVersionRange(2, 8)
var minSupportedMongoDBVersion = "2.6"

type fsm struct {
	description.Topology
	SetName       string
	maxElectionID primitive.ObjectID
	maxSetVersion uint32
}

func newFSM() *fsm {
	return new(fsm)
}

// apply should operate on immutable TopologyDescriptions and Descriptions. This way we don't have to
// lock for the entire time we're applying server description.
func (f *fsm) apply(s description.Server) (description.Topology, error) {

	newServers := make([]description.Server, len(f.Servers))
	copy(newServers, f.Servers)

	oldMinutes := f.SessionTimeoutMinutes
	f.Topology = description.Topology{
		Kind:    f.Kind,
		Servers: newServers,
	}

	// For data bearing servers, set SessionTimeoutMinutes to the lowest among them
	if oldMinutes == 0 {
		// If timeout currently 0, check all servers to see if any still don't have a timeout
		// If they all have timeout, pick the lowest.
		timeout := s.SessionTimeoutMinutes
		for _, server := range f.Servers {
			if server.DataBearing() && server.SessionTimeoutMinutes < timeout {
				timeout = server.SessionTimeoutMinutes
			}
		}
		f.SessionTimeoutMinutes = timeout
	} else {
		if s.DataBearing() && oldMinutes > s.SessionTimeoutMinutes {
			f.SessionTimeoutMinutes = s.SessionTimeoutMinutes
		} else {
			f.SessionTimeoutMinutes = oldMinutes
		}
	}

	if _, ok := f.findServer(s.Addr); !ok {
		return f.Topology, nil
	}

	if s.WireVersion != nil {
		if s.WireVersion.Max < supportedWireVersions.Min {
			return description.Topology{}, fmt.Errorf(
				"server at %s reports wire version %d, but this version of the Go driver requires "+
					"at least %d (MongoDB %s)",
				s.Addr.String(),
				s.WireVersion.Max,
				supportedWireVersions.Min,
				minSupportedMongoDBVersion,
			)
		}

		if s.WireVersion.Min > supportedWireVersions.Max {
			return description.Topology{}, fmt.Errorf(
				"server at %s requires wire version %d, but this version of the Go driver only "+
					"supports up to %d",
				s.Addr.String(),
				s.WireVersion.Min,
				supportedWireVersions.Max,
			)
		}
	}

	switch f.Kind {
	case description.Unknown:
		f.applyToUnknown(s)
	case description.Sharded:
		f.applyToSharded(s)
	case description.ReplicaSetNoPrimary:
		f.applyToReplicaSetNoPrimary(s)
	case description.ReplicaSetWithPrimary:
		f.applyToReplicaSetWithPrimary(s)
	case description.Single:
		f.applyToSingle(s)
	}

	return f.Topology, nil
}

func (f *fsm) applyToReplicaSetNoPrimary(s description.Server) {
	switch s.Kind {
	case description.Standalone, description.Mongos:
		f.removeServerByAddr(s.Addr)
	case description.RSPrimary:
		f.updateRSFromPrimary(s)
	case description.RSSecondary, description.RSArbiter, description.RSMember:
		f.updateRSWithoutPrimary(s)
	case description.Unknown, description.RSGhost:
		f.replaceServer(s)
	}
}

func (f *fsm) applyToReplicaSetWithPrimary(s description.Server) {
	switch s.Kind {
	case description.Standalone, description.Mongos:
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
	case description.RSPrimary:
		f.updateRSFromPrimary(s)
	case description.RSSecondary, description.RSArbiter, description.RSMember:
		f.updateRSWithPrimaryFromMember(s)
	case description.Unknown, description.RSGhost:
		f.replaceServer(s)
		f.checkIfHasPrimary()
	}
}

func (f *fsm) applyToSharded(s description.Server) {
	switch s.Kind {
	case description.Mongos, description.Unknown:
		f.replaceServer(s)
	case description.Standalone, description.RSPrimary, description.RSSecondary, description.RSArbiter, description.RSMember, description.RSGhost:
		f.removeServerByAddr(s.Addr)
	}
}

func (f *fsm) applyToSingle(s description.Server) {
	switch s.Kind {
	case description.Unknown:
		f.replaceServer(s)
	case description.Standalone, description.Mongos:
		if f.SetName != "" {
			f.removeServerByAddr(s.Addr)
			return
		}

		f.replaceServer(s)
	case description.RSPrimary, description.RSSecondary, description.RSArbiter, description.RSMember, description.RSGhost:
		if f.SetName != "" && f.SetName != s.SetName {
			f.removeServerByAddr(s.Addr)
			return
		}

		f.replaceServer(s)
	}
}

func (f *fsm) applyToUnknown(s description.Server) {
	switch s.Kind {
	case description.Mongos:
		f.setKind(description.Sharded)
		f.replaceServer(s)
	case description.RSPrimary:
		f.updateRSFromPrimary(s)
	case description.RSSecondary, description.RSArbiter, description.RSMember:
		f.setKind(description.ReplicaSetNoPrimary)
		f.updateRSWithoutPrimary(s)
	case description.Standalone:
		f.updateUnknownWithStandalone(s)
	case description.Unknown, description.RSGhost:
		f.replaceServer(s)
	}
}

func (f *fsm) checkIfHasPrimary() {
	if _, ok := f.findPrimary(); ok {
		f.setKind(description.ReplicaSetWithPrimary)
	} else {
		f.setKind(description.ReplicaSetNoPrimary)
	}
}

func (f *fsm) updateRSFromPrimary(s description.Server) {
	if f.SetName == "" {
		f.SetName = s.SetName
	} else if f.SetName != s.SetName {
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
		return
	}

	if s.SetVersion != 0 && !bytes.Equal(s.ElectionID[:], primitive.NilObjectID[:]) {
		if f.maxSetVersion > s.SetVersion || bytes.Compare(f.maxElectionID[:], s.ElectionID[:]) == 1 {
			f.replaceServer(description.Server{
				Addr:      s.Addr,
				LastError: fmt.Errorf("was a primary, but its set version or election id is stale"),
			})
			f.checkIfHasPrimary()
			return
		}

		f.maxElectionID = s.ElectionID
	}

	if s.SetVersion > f.maxSetVersion {
		f.maxSetVersion = s.SetVersion
	}

	if j, ok := f.findPrimary(); ok {
		f.setServer(j, description.Server{
			Addr:      f.Servers[j].Addr,
			LastError: fmt.Errorf("was a primary, but a new primary was discovered"),
		})
	}

	f.replaceServer(s)

	for j := len(f.Servers) - 1; j >= 0; j-- {
		found := false
		for _, member := range s.Members {
			if member == f.Servers[j].Addr {
				found = true
				break
			}
		}
		if !found {
			f.removeServer(j)
		}
	}

	for _, member := range s.Members {
		if _, ok := f.findServer(member); !ok {
			f.addServer(member)
		}
	}

	f.checkIfHasPrimary()
}

func (f *fsm) updateRSWithPrimaryFromMember(s description.Server) {
	if f.SetName != s.SetName {
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
		return
	}

	if s.Addr != s.CanonicalAddr {
		f.removeServerByAddr(s.Addr)
		f.checkIfHasPrimary()
		return
	}

	f.replaceServer(s)

	if _, ok := f.findPrimary(); !ok {
		f.setKind(description.ReplicaSetNoPrimary)
	}
}

func (f *fsm) updateRSWithoutPrimary(s description.Server) {
	if f.SetName == "" {
		f.SetName = s.SetName
	} else if f.SetName != s.SetName {
		f.removeServerByAddr(s.Addr)
		return
	}

	for _, member := range s.Members {
		if _, ok := f.findServer(member); !ok {
			f.addServer(member)
		}
	}

	if s.Addr != s.CanonicalAddr {
		f.removeServerByAddr(s.Addr)
		return
	}

	f.replaceServer(s)
}

func (f *fsm) updateUnknownWithStandalone(s description.Server) {
	if len(f.Servers) > 1 {
		f.removeServerByAddr(s.Addr)
		return
	}

	f.setKind(description.Single)
	f.replaceServer(s)
}

func (f *fsm) addServer(addr address.Address) {
	f.Servers = append(f.Servers, description.Server{
		Addr: addr.Canonicalize(),
	})
}

func (f *fsm) findPrimary() (int, bool) {
	for i, s := range f.Servers {
		if s.Kind == description.RSPrimary {
			return i, true
		}
	}

	return 0, false
}

func (f *fsm) findServer(addr address.Address) (int, bool) {
	canon := addr.Canonicalize()
	for i, s := range f.Servers {
		if canon == s.Addr {
			return i, true
		}
	}

	return 0, false
}

func (f *fsm) removeServer(i int) {
	f.Servers = append(f.Servers[:i], f.Servers[i+1:]...)
}

func (f *fsm) removeServerByAddr(addr address.Address) {
	if i, ok := f.findServer(addr); ok {
		f.removeServer(i)
	}
}

func (f *fsm) replaceServer(s description.Server) bool {
	if i, ok := f.findServer(s.Addr); ok {
		f.setServer(i, s)
		return true
	}
	return false
}

func (f *fsm) setServer(i int, s description.Server) {
	f.Servers[i] = s
}

func (f *fsm) setKind(k description.TopologyKind) {
	f.Kind = k
}
