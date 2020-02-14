// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package description

import (
	"fmt"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/tag"
)

// ServerSelector is an interface implemented by types that can select a server given a
// topology description.
type ServerSelector interface {
	SelectServer(Topology, []Server) ([]Server, error)
}

// ServerSelectorFunc is a function that can be used as a ServerSelector.
type ServerSelectorFunc func(Topology, []Server) ([]Server, error)

// SelectServer implements the ServerSelector interface.
func (ssf ServerSelectorFunc) SelectServer(t Topology, s []Server) ([]Server, error) {
	return ssf(t, s)
}

type compositeSelector struct {
	selectors []ServerSelector
}

// CompositeSelector combines multiple selectors into a single selector.
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

// LatencySelector creates a ServerSelector which selects servers based on their latency.
func LatencySelector(latency time.Duration) ServerSelector {
	return &latencySelector{latency: latency}
}

func (ls *latencySelector) SelectServer(t Topology, candidates []Server) ([]Server, error) {
	if ls.latency < 0 {
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

		max := min + ls.latency

		var result []Server
		for _, candidate := range candidates {
			if candidate.AverageRTTSet {
				if candidate.AverageRTT <= max {
					result = append(result, candidate)
				}
			}
		}

		return result, nil
	}
}

// WriteSelector selects all the writable servers.
func WriteSelector() ServerSelector {
	return ServerSelectorFunc(func(t Topology, candidates []Server) ([]Server, error) {
		switch t.Kind {
		case Single:
			return candidates, nil
		default:
			result := []Server{}
			for _, candidate := range candidates {
				switch candidate.Kind {
				case Mongos, RSPrimary, Standalone:
					result = append(result, candidate)
				}
			}
			return result, nil
		}
	})
}

// ReadPrefSelector selects servers based on the provided read preference.
func ReadPrefSelector(rp *readpref.ReadPref) ServerSelector {
	return ServerSelectorFunc(func(t Topology, candidates []Server) ([]Server, error) {
		if _, set := rp.MaxStaleness(); set {
			for _, s := range candidates {
				if s.Kind != Unknown {
					if err := MaxStalenessSupported(s.WireVersion); err != nil {
						return nil, err
					}
				}
			}
		}

		switch t.Kind {
		case Single:
			return candidates, nil
		case ReplicaSetNoPrimary, ReplicaSetWithPrimary:
			return selectForReplicaSet(rp, t, candidates)
		case Sharded:
			return selectByKind(candidates, Mongos), nil
		}

		return nil, nil
	})
}

func selectForReplicaSet(rp *readpref.ReadPref, t Topology, candidates []Server) ([]Server, error) {
	if err := verifyMaxStaleness(rp, t); err != nil {
		return nil, err
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
		var results []Server
		for _, s := range candidates {
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
	var result []Server
	for _, s := range candidates {
		if s.Kind == kind {
			result = append(result, s)
		}
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
