// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package topology

import (
	"bytes"
	"fmt"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal/ptrutil"
	"go.mongodb.org/mongo-driver/mongo/address"
	"go.mongodb.org/mongo-driver/mongo/description"
)

var (
	// MinSupportedMongoDBVersion is the version string for the lowest MongoDB version supported by the driver.
	MinSupportedMongoDBVersion = "3.6"

	// SupportedWireVersions is the range of wire versions supported by the driver.
	SupportedWireVersions = description.NewVersionRange(6, 25)
)

type fsm struct {
	description.Topology
	maxElectionID    primitive.ObjectID
	maxSetVersion    uint32
	compatible       atomic.Value
	compatibilityErr error
}

func newFSM() *fsm {
	f := fsm{}
	f.compatible.Store(true)
	return &f
}

// selectFSMSessionTimeout selects the timeout to return for the topology's
// finite state machine. If the logicalSessionTimeoutMinutes on the FSM exists
// and the server is data-bearing, then we determine this value by returning
//
//	min{server timeout, FSM timeout}
//
// where a "nil" value is considered less than 0.
//
// Otherwise, if the FSM's logicalSessionTimeoutMinutes exist, then this
// function returns the FSM timeout.
//
// In the case where the FSM timeout DNE, we check all servers to see if any
// still do not have a timeout. This function chooses the lowest of the existing
// timeouts.
func selectFSMSessionTimeout(f *fsm, s description.Server) *int64 {
	oldMinutes := f.SessionTimeoutMinutesPtr
	comp := ptrutil.CompareInt64(oldMinutes, s.SessionTimeoutMinutesPtr)

	// If the server is data-bearing and the current timeout exists and is
	// either:
	//
	// 1. larger than the server timeout, or
	// 2. non-nil while the server timeout is nil
	//
	// then return the server timeout.
	if s.DataBearing() && (comp == 1 || comp == 2) {
		return s.SessionTimeoutMinutesPtr
	}

	// If the current timeout exists and the server is not data-bearing OR
	// min{server timeout, current timeout} = current timeout, then return
	// the current timeout.
	if oldMinutes != nil {
		return oldMinutes
	}

	timeout := s.SessionTimeoutMinutesPtr
	for _, server := range f.Servers {
		// If the server is not data-bearing, then we do not consider
		// it's timeout whether set or not.
		if !server.DataBearing() {
			continue
		}

		srvTimeout := server.SessionTimeoutMinutesPtr
		comp := ptrutil.CompareInt64(timeout, srvTimeout)

		if comp <= 0 { // timeout <= srvTimout
			continue
		}

		timeout = server.SessionTimeoutMinutesPtr
	}

	return timeout
}

// apply takes a new server description and modifies the FSM's topology description based on it. It returns the
// updated topology description as well as a server description. The returned server description is either the same
// one that was passed in, or a new one in the case that it had to be changed.
//
// apply should operation on immutable descriptions so we don't have to lock for the entire time we're applying the
// server description.
func (f *fsm) apply(s description.Server) (description.Topology, description.Server) {
	newServers := make([]description.Server, len(f.Servers))
	copy(newServers, f.Servers)

	// Reset the logicalSessionTimeoutMinutes to the minimum of the FSM
	// and the description.server/f.servers.
	serverTimeoutMinutes := selectFSMSessionTimeout(f, s)

	f.Topology = description.Topology{
		Kind:    f.Kind,
		Servers: newServers,
		SetName: f.SetName,
	}

	f.Topology.SessionTimeoutMinutesPtr = serverTimeoutMinutes

	if serverTimeoutMinutes != nil {
		f.SessionTimeoutMinutes = uint32(*serverTimeoutMinutes)
	}

	if _, ok := f.findServer(s.Addr); !ok {
		return f.Topology, s
	}

	updatedDesc := s
	switch f.Kind {
	case description.Unknown:
		updatedDesc = f.applyToUnknown(s)
	case description.Sharded:
		updatedDesc = f.applyToSharded(s)
	case description.ReplicaSetNoPrimary:
		updatedDesc = f.applyToReplicaSetNoPrimary(s)
	case description.ReplicaSetWithPrimary:
		updatedDesc = f.applyToReplicaSetWithPrimary(s)
	case description.Single:
		updatedDesc = f.applyToSingle(s)
	}

	for _, server := range f.Servers {
		if server.WireVersion != nil {
			if server.WireVersion.Max < SupportedWireVersions.Min {
				f.compatible.Store(false)
				f.compatibilityErr = fmt.Errorf(
					"server at %s reports wire version %d, but this version of the Go driver requires "+
						"at least %d (MongoDB %s)",
					server.Addr.String(),
					server.WireVersion.Max,
					SupportedWireVersions.Min,
					MinSupportedMongoDBVersion,
				)
				f.Topology.CompatibilityErr = f.compatibilityErr
				return f.Topology, s
			}

			if server.WireVersion.Min > SupportedWireVersions.Max {
				f.compatible.Store(false)
				f.compatibilityErr = fmt.Errorf(
					"server at %s requires wire version %d, but this version of the Go driver only supports up to %d",
					server.Addr.String(),
					server.WireVersion.Min,
					SupportedWireVersions.Max,
				)
				f.Topology.CompatibilityErr = f.compatibilityErr
				return f.Topology, s
			}
		}
	}

	f.compatible.Store(true)
	f.compatibilityErr = nil

	return f.Topology, updatedDesc
}

func (f *fsm) applyToReplicaSetNoPrimary(s description.Server) description.Server {
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

	return s
}

func (f *fsm) applyToReplicaSetWithPrimary(s description.Server) description.Server {
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

	return s
}

func (f *fsm) applyToSharded(s description.Server) description.Server {
	switch s.Kind {
	case description.Mongos, description.Unknown:
		f.replaceServer(s)
	case description.Standalone, description.RSPrimary, description.RSSecondary, description.RSArbiter, description.RSMember, description.RSGhost:
		f.removeServerByAddr(s.Addr)
	}

	return s
}

func (f *fsm) applyToSingle(s description.Server) description.Server {
	switch s.Kind {
	case description.Unknown:
		f.replaceServer(s)
	case description.Standalone, description.Mongos:
		if f.SetName != "" {
			f.removeServerByAddr(s.Addr)
			return s
		}

		f.replaceServer(s)
	case description.RSPrimary, description.RSSecondary, description.RSArbiter, description.RSMember, description.RSGhost:
		// A replica set name can be provided when creating a direct connection. In this case, if the set name returned
		// by the hello response doesn't match up with the one provided during configuration, the server description
		// is replaced with a default Unknown description.
		//
		// We create a new server description rather than doing s.Kind = description.Unknown because the other fields,
		// such as RTT, need to be cleared for Unknown descriptions as well.
		if f.SetName != "" && f.SetName != s.SetName {
			s = description.Server{
				Addr: s.Addr,
				Kind: description.Unknown,
			}
		}

		f.replaceServer(s)
	}

	return s
}

func (f *fsm) applyToUnknown(s description.Server) description.Server {
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

	return s
}

func (f *fsm) checkIfHasPrimary() {
	if _, ok := f.findPrimary(); ok {
		f.setKind(description.ReplicaSetWithPrimary)
	} else {
		f.setKind(description.ReplicaSetNoPrimary)
	}
}

// hasStalePrimary returns true if the topology has a primary that is "stale".
func hasStalePrimary(fsm fsm, srv description.Server) bool {
	// Compare the election ID values of the server and the topology lexicographically.
	compRes := bytes.Compare(srv.ElectionID[:], fsm.maxElectionID[:])

	if wireVersion := srv.WireVersion; wireVersion != nil && wireVersion.Max >= 17 {
		// In the Post-6.0 case, a primary is considered "stale" if the server's election ID is greater than the
		// topology's max election ID. In these versions, the primary is also considered "stale" if the server's
		// election ID is LTE to the topologies election ID and the server's "setVersion" is less than the topology's
		// max "setVersion".
		return compRes == -1 || (compRes != 1 && srv.SetVersion < fsm.maxSetVersion)
	}

	// If the server's election ID is less than the topology's max election ID, the primary is considered
	// "stale". Similarly, if the server's "setVersion" is less than the topology's max "setVersion", the
	// primary is considered stale.
	return compRes == -1 || fsm.maxSetVersion > srv.SetVersion
}

// transferEVTuple will transfer the ("ElectionID", "SetVersion") tuple from the description server to the topology.
// If the primary is stale, the tuple will not be transferred, the topology will update it's "Kind" value, and this
// routine will return "false".
func transferEVTuple(srv description.Server, fsm *fsm) bool {
	stalePrimary := hasStalePrimary(*fsm, srv)

	if wireVersion := srv.WireVersion; wireVersion != nil && wireVersion.Max >= 17 {
		if stalePrimary {
			fsm.checkIfHasPrimary()
			return false
		}

		fsm.maxElectionID = srv.ElectionID
		fsm.maxSetVersion = srv.SetVersion

		return true
	}

	if srv.SetVersion != 0 && !srv.ElectionID.IsZero() {
		if stalePrimary {
			fsm.replaceServer(description.Server{
				Addr: srv.Addr,
				LastError: fmt.Errorf(
					"was a primary, but its set version or election id is stale"),
			})

			fsm.checkIfHasPrimary()

			return false
		}

		fsm.maxElectionID = srv.ElectionID
	}

	if srv.SetVersion > fsm.maxSetVersion {
		fsm.maxSetVersion = srv.SetVersion
	}

	return true
}

func (f *fsm) updateRSFromPrimary(srv description.Server) {
	if f.SetName == "" {
		f.SetName = srv.SetName
	} else if f.SetName != srv.SetName {
		f.removeServerByAddr(srv.Addr)
		f.checkIfHasPrimary()

		return
	}

	if ok := transferEVTuple(srv, f); !ok {
		return
	}

	if j, ok := f.findPrimary(); ok {
		f.setServer(j, description.Server{
			Addr:      f.Servers[j].Addr,
			LastError: fmt.Errorf("was a primary, but a new primary was discovered"),
		})
	}

	f.replaceServer(srv)

	for j := len(f.Servers) - 1; j >= 0; j-- {
		found := false
		for _, member := range srv.Members {
			if member == f.Servers[j].Addr {
				found = true
				break
			}
		}

		if !found {
			f.removeServer(j)
		}
	}

	for _, member := range srv.Members {
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

func (f *fsm) replaceServer(s description.Server) {
	if i, ok := f.findServer(s.Addr); ok {
		f.setServer(i, s)
	}
}

func (f *fsm) setServer(i int, s description.Server) {
	f.Servers[i] = s
}

func (f *fsm) setKind(k description.TopologyKind) {
	f.Kind = k
}
