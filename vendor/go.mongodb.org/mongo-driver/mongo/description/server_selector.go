// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/tag"
)

// ServerSelector is an interface implemented by types that can perform server selection given a topology description
// and list of candidate servers. The selector should filter the provided candidates list and return a subset that
// matches some criteria.
type ServerSelector interface {
	SelectServer(Topology, []Server) ([]Server, error)
}

// ServerSelectorFunc is a function that can be used as a ServerSelector.
type ServerSelectorFunc func(Topology, []Server) ([]Server, error)

// SelectServer implements the ServerSelector interface.
func (ssf ServerSelectorFunc) SelectServer(t Topology, s []Server) ([]Server, error) {
	return ssf(t, s)
}

// serverSelectorInfo contains metadata concerning the server selector for the
// purpose of publication.
type serverSelectorInfo struct {
	Type      string
	Data      string               `json:",omitempty"`
	Selectors []serverSelectorInfo `json:",omitempty"`
}

// String returns the JSON string representation of the serverSelectorInfo.
func (sss serverSelectorInfo) String() string {
	bytes, _ := json.Marshal(sss)

	return string(bytes)
}

// serverSelectorInfoGetter is an interface that defines an info() method to
// get the serverSelectorInfo.
type serverSelectorInfoGetter interface {
	info() serverSelectorInfo
}

type compositeSelector struct {
	selectors []ServerSelector
}

func (cs *compositeSelector) info() serverSelectorInfo {
	csInfo := serverSelectorInfo{Type: "compositeSelector"}

	for _, sel := range cs.selectors {
		if getter, ok := sel.(serverSelectorInfoGetter); ok {
			csInfo.Selectors = append(csInfo.Selectors, getter.info())
		}
	}

	return csInfo
}

// String returns the JSON string representation of the compositeSelector.
func (cs *compositeSelector) String() string {
	return cs.info().String()
}

// CompositeSelector combines multiple selectors into a single selector by applying them in order to the candidates
// list.
//
// For example, if the initial candidates list is [s0, s1, s2, s3] and two selectors are provided where the first
// matches s0 and s1 and the second matches s1 and s2, the following would occur during server selection:
//
// 1. firstSelector([s0, s1, s2, s3]) -> [s0, s1]
// 2. secondSelector([s0, s1]) -> [s1]
//
// The final list of candidates returned by the composite selector would be [s1].
func CompositeSelector(selectors []ServerSelector) ServerSelector {
	return &compositeSelector{selectors: selectors}
}

func (cs *compositeSelector) SelectServer(t Topology, candidates []Server) ([]Server, error) {
	var err error
	for _, sel := range cs.selectors {
		candidates, err = sel.SelectServer(t, candidates)
		if err != nil {
			return nil, err
		}
	}
	return candidates, nil
}

type latencySelector struct {
	latency time.Duration
}

// LatencySelector creates a ServerSelector which selects servers based on their average RTT values.
func LatencySelector(latency time.Duration) ServerSelector {
	return &latencySelector{latency: latency}
}

func (latencySelector) info() serverSelectorInfo {
	return serverSelectorInfo{Type: "latencySelector"}
}

func (selector latencySelector) String() string {
	return selector.info().String()
}

func (selector *latencySelector) SelectServer(t Topology, candidates []Server) ([]Server, error) {
	if selector.latency < 0 {
		return candidates, nil
	}
	if t.Kind == LoadBalanced {
		// In LoadBalanced mode, there should only be one server in the topology and it must be selected.
		return candidates, nil
	}

	switch len(candidates) {
	case 0, 1:
		return candidates, nil
	default:
		min := time.Duration(math.MaxInt64)
		for _, candidate := range candidates {
			if candidate.AverageRTTSet {
				if candidate.AverageRTT < min {
					min = candidate.AverageRTT
				}
			}
		}

		if min == math.MaxInt64 {
			return candidates, nil
		}

		max := min + selector.latency

		viableIndexes := make([]int, 0, len(candidates))
		for i, candidate := range candidates {
			if candidate.AverageRTTSet {
				if candidate.AverageRTT <= max {
					viableIndexes = append(viableIndexes, i)
				}
			}
		}
		if len(viableIndexes) == len(candidates) {
			return candidates, nil
		}
		result := make([]Server, len(viableIndexes))
		for i, idx := range viableIndexes {
			result[i] = candidates[idx]
		}
		return result, nil
	}
}

type writeServerSelector struct{}

// WriteSelector selects all the writable servers.
func WriteSelector() ServerSelector {
	return writeServerSelector{}
}

func (writeServerSelector) info() serverSelectorInfo {
	return serverSelectorInfo{Type: "writeSelector"}
}

func (selector writeServerSelector) String() string {
	return selector.info().String()
}

func (writeServerSelector) SelectServer(t Topology, candidates []Server) ([]Server, error) {
	switch t.Kind {
	case Single, LoadBalanced:
		return candidates, nil
	default:
		// Determine the capacity of the results slice.
		selected := 0
		for _, candidate := range candidates {
			switch candidate.Kind {
			case Mongos, RSPrimary, Standalone:
				selected++
			}
		}

		// Append candidates to the results slice.
		result := make([]Server, 0, selected)
		for _, candidate := range candidates {
			switch candidate.Kind {
			case Mongos, RSPrimary, Standalone:
				result = append(result, candidate)
			}
		}
		return result, nil
	}
}

type readPrefServerSelector struct {
	rp                *readpref.ReadPref
	isOutputAggregate bool
}

// ReadPrefSelector selects servers based on the provided read preference.
func ReadPrefSelector(rp *readpref.ReadPref) ServerSelector {
	return readPrefServerSelector{
		rp:                rp,
		isOutputAggregate: false,
	}
}

func (selector readPrefServerSelector) info() serverSelectorInfo {
	return serverSelectorInfo{
		Type: "readPrefSelector",
		Data: selector.rp.String(),
	}
}

func (selector readPrefServerSelector) String() string {
	return selector.info().String()
}

func (selector readPrefServerSelector) SelectServer(t Topology, candidates []Server) ([]Server, error) {
	if t.Kind == LoadBalanced {
		// In LoadBalanced mode, there should only be one server in the topology and it must be selected. We check
		// this before checking MaxStaleness support because there's no monitoring in this mode, so the candidate
		// server wouldn't have a wire version set, which would result in an error.
		return candidates, nil
	}

	switch t.Kind {
	case Single:
		return candidates, nil
	case ReplicaSetNoPrimary, ReplicaSetWithPrimary:
		return selectForReplicaSet(selector.rp, selector.isOutputAggregate, t, candidates)
	case Sharded:
		return selectByKind(candidates, Mongos), nil
	}

	return nil, nil
}

// OutputAggregateSelector selects servers based on the provided read preference
// given that the underlying operation is aggregate with an output stage.
func OutputAggregateSelector(rp *readpref.ReadPref) ServerSelector {
	return readPrefServerSelector{
		rp:                rp,
		isOutputAggregate: true,
	}
}

func selectForReplicaSet(rp *readpref.ReadPref, isOutputAggregate bool, t Topology, candidates []Server) ([]Server, error) {
	if err := verifyMaxStaleness(rp, t); err != nil {
		return nil, err
	}

	// If underlying operation is an aggregate with an output stage, only apply read preference
	// if all candidates are 5.0+. Otherwise, operate under primary read preference.
	if isOutputAggregate {
		for _, s := range candidates {
			if s.WireVersion.Max < 13 {
				return selectByKind(candidates, RSPrimary), nil
			}
		}
	}

	switch rp.Mode() {
	case readpref.PrimaryMode:
		return selectByKind(candidates, RSPrimary), nil
	case readpref.PrimaryPreferredMode:
		selected := selectByKind(candidates, RSPrimary)

		if len(selected) == 0 {
			selected = selectSecondaries(rp, candidates)
			return selectByTagSet(selected, rp.TagSets()), nil
		}

		return selected, nil
	case readpref.SecondaryPreferredMode:
		selected := selectSecondaries(rp, candidates)
		selected = selectByTagSet(selected, rp.TagSets())
		if len(selected) > 0 {
			return selected, nil
		}
		return selectByKind(candidates, RSPrimary), nil
	case readpref.SecondaryMode:
		selected := selectSecondaries(rp, candidates)
		return selectByTagSet(selected, rp.TagSets()), nil
	case readpref.NearestMode:
		selected := selectByKind(candidates, RSPrimary)
		selected = append(selected, selectSecondaries(rp, candidates)...)
		return selectByTagSet(selected, rp.TagSets()), nil
	}

	return nil, fmt.Errorf("unsupported mode: %d", rp.Mode())
}

func selectSecondaries(rp *readpref.ReadPref, candidates []Server) []Server {
	secondaries := selectByKind(candidates, RSSecondary)
	if len(secondaries) == 0 {
		return secondaries
	}
	if maxStaleness, set := rp.MaxStaleness(); set {
		primaries := selectByKind(candidates, RSPrimary)
		if len(primaries) == 0 {
			baseTime := secondaries[0].LastWriteTime
			for i := 1; i < len(secondaries); i++ {
				if secondaries[i].LastWriteTime.After(baseTime) {
					baseTime = secondaries[i].LastWriteTime
				}
			}

			var selected []Server
			for _, secondary := range secondaries {
				estimatedStaleness := baseTime.Sub(secondary.LastWriteTime) + secondary.HeartbeatInterval
				if estimatedStaleness <= maxStaleness {
					selected = append(selected, secondary)
				}
			}

			return selected
		}

		primary := primaries[0]

		var selected []Server
		for _, secondary := range secondaries {
			estimatedStaleness := secondary.LastUpdateTime.Sub(secondary.LastWriteTime) - primary.LastUpdateTime.Sub(primary.LastWriteTime) + secondary.HeartbeatInterval
			if estimatedStaleness <= maxStaleness {
				selected = append(selected, secondary)
			}
		}
		return selected
	}

	return secondaries
}

func selectByTagSet(candidates []Server, tagSets []tag.Set) []Server {
	if len(tagSets) == 0 {
		return candidates
	}

	for _, ts := range tagSets {
		// If this tag set is empty, we can take a fast path because the empty list is a subset of all tag sets, so
		// all candidate servers will be selected.
		if len(ts) == 0 {
			return candidates
		}

		var results []Server
		for _, s := range candidates {
			// ts is non-empty, so only servers with a non-empty set of tags need to be checked.
			if len(s.Tags) > 0 && s.Tags.ContainsAll(ts) {
				results = append(results, s)
			}
		}

		if len(results) > 0 {
			return results
		}
	}

	return []Server{}
}

func selectByKind(candidates []Server, kind ServerKind) []Server {
	// Record the indices of viable candidates first and then append those to the returned slice
	// to avoid appending costly Server structs directly as an optimization.
	viableIndexes := make([]int, 0, len(candidates))
	for i, s := range candidates {
		if s.Kind == kind {
			viableIndexes = append(viableIndexes, i)
		}
	}
	if len(viableIndexes) == len(candidates) {
		return candidates
	}
	result := make([]Server, len(viableIndexes))
	for i, idx := range viableIndexes {
		result[i] = candidates[idx]
	}
	return result
}

func verifyMaxStaleness(rp *readpref.ReadPref, t Topology) error {
	maxStaleness, set := rp.MaxStaleness()
	if !set {
		return nil
	}

	if maxStaleness < 90*time.Second {
		return fmt.Errorf("max staleness (%s) must be greater than or equal to 90s", maxStaleness)
	}

	if len(t.Servers) < 1 {
		// Maybe we should return an error here instead?
		return nil
	}

	// we'll assume all candidates have the same heartbeat interval.
	s := t.Servers[0]
	idleWritePeriod := 10 * time.Second

	if maxStaleness < s.HeartbeatInterval+idleWritePeriod {
		return fmt.Errorf(
			"max staleness (%s) must be greater than or equal to the heartbeat interval (%s) plus idle write period (%s)",
			maxStaleness, s.HeartbeatInterval, idleWritePeriod,
		)
	}

	return nil
}
