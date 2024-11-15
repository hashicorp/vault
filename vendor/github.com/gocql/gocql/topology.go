/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*
 * Content before git sha 34fdeebefcbf183ed7f916f931aa0586fdaa1b40
 * Copyright (c) 2016, The Gocql authors,
 * provided under the BSD-3-Clause License.
 * See the NOTICE file distributed with this work for additional information.
 */

package gocql

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type hostTokens struct {
	// token is end (inclusive) of token range these hosts belong to
	token token
	hosts []*HostInfo
}

// tokenRingReplicas maps token ranges to list of replicas.
// The elements in tokenRingReplicas are sorted by token ascending.
// The range for a given item in tokenRingReplicas starts after preceding range and ends with the token specified in
// token. The end token is part of the range.
// The lowest (i.e. index 0) range wraps around the ring (its preceding range is the one with largest index).
type tokenRingReplicas []hostTokens

func (h tokenRingReplicas) Less(i, j int) bool { return h[i].token.Less(h[j].token) }
func (h tokenRingReplicas) Len() int           { return len(h) }
func (h tokenRingReplicas) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h tokenRingReplicas) replicasFor(t token) *hostTokens {
	if len(h) == 0 {
		return nil
	}

	p := sort.Search(len(h), func(i int) bool {
		return !h[i].token.Less(t)
	})

	if p >= len(h) {
		// rollover
		p = 0
	}

	return &h[p]
}

type placementStrategy interface {
	replicaMap(tokenRing *tokenRing) tokenRingReplicas
	replicationFactor(dc string) int
}

func getReplicationFactorFromOpts(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		if v < 0 {
			return 0, fmt.Errorf("invalid replication_factor %d", v)
		}
		return v, nil
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("invalid replication_factor %q: %v", v, err)
		} else if n < 0 {
			return 0, fmt.Errorf("invalid replication_factor %d", n)
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unknown replication_factor type %T", v)
	}
}

func getStrategy(ks *KeyspaceMetadata, logger StdLogger) placementStrategy {
	switch {
	case strings.Contains(ks.StrategyClass, "SimpleStrategy"):
		rf, err := getReplicationFactorFromOpts(ks.StrategyOptions["replication_factor"])
		if err != nil {
			logger.Printf("parse rf for keyspace %q: %v", ks.Name, err)
			return nil
		}
		return &simpleStrategy{rf: rf}
	case strings.Contains(ks.StrategyClass, "NetworkTopologyStrategy"):
		dcs := make(map[string]int)
		for dc, rf := range ks.StrategyOptions {
			if dc == "class" {
				continue
			}

			rf, err := getReplicationFactorFromOpts(rf)
			if err != nil {
				logger.Println("parse rf for keyspace %q, dc %q: %v", err)
				// skip DC if the rf is invalid/unsupported, so that we can at least work with other working DCs.
				continue
			}

			dcs[dc] = rf
		}
		return &networkTopology{dcs: dcs}
	case strings.Contains(ks.StrategyClass, "LocalStrategy"):
		return nil
	default:
		logger.Printf("parse rf for keyspace %q: unsupported strategy class: %v", ks.StrategyClass)
		return nil
	}
}

type simpleStrategy struct {
	rf int
}

func (s *simpleStrategy) replicationFactor(dc string) int {
	return s.rf
}

func (s *simpleStrategy) replicaMap(tokenRing *tokenRing) tokenRingReplicas {
	tokens := tokenRing.tokens
	ring := make(tokenRingReplicas, len(tokens))

	for i, th := range tokens {
		replicas := make([]*HostInfo, 0, s.rf)
		seen := make(map[*HostInfo]bool)

		for j := 0; j < len(tokens) && len(replicas) < s.rf; j++ {
			h := tokens[(i+j)%len(tokens)]
			if !seen[h.host] {
				replicas = append(replicas, h.host)
				seen[h.host] = true
			}
		}

		ring[i] = hostTokens{th.token, replicas}
	}

	sort.Sort(ring)

	return ring
}

type networkTopology struct {
	dcs map[string]int
}

func (n *networkTopology) replicationFactor(dc string) int {
	return n.dcs[dc]
}

func (n *networkTopology) haveRF(replicaCounts map[string]int) bool {
	if len(replicaCounts) != len(n.dcs) {
		return false
	}

	for dc, rf := range n.dcs {
		if rf != replicaCounts[dc] {
			return false
		}
	}

	return true
}

func (n *networkTopology) replicaMap(tokenRing *tokenRing) tokenRingReplicas {
	dcRacks := make(map[string]map[string]struct{}, len(n.dcs))
	// skipped hosts in a dc
	skipped := make(map[string][]*HostInfo, len(n.dcs))
	// number of replicas per dc
	replicasInDC := make(map[string]int, len(n.dcs))
	// dc -> racks
	seenDCRacks := make(map[string]map[string]struct{}, len(n.dcs))

	for _, h := range tokenRing.hosts {
		dc := h.DataCenter()
		rack := h.Rack()

		racks, ok := dcRacks[dc]
		if !ok {
			racks = make(map[string]struct{})
			dcRacks[dc] = racks
		}
		racks[rack] = struct{}{}
	}

	for dc, racks := range dcRacks {
		replicasInDC[dc] = 0
		seenDCRacks[dc] = make(map[string]struct{}, len(racks))
	}

	tokens := tokenRing.tokens
	replicaRing := make(tokenRingReplicas, 0, len(tokens))

	var totalRF int
	for _, rf := range n.dcs {
		totalRF += rf
	}

	for i, th := range tokenRing.tokens {
		if rf := n.dcs[th.host.DataCenter()]; rf == 0 {
			// skip this token since no replica in this datacenter.
			continue
		}

		for k, v := range skipped {
			skipped[k] = v[:0]
		}

		for dc := range n.dcs {
			replicasInDC[dc] = 0
			for rack := range seenDCRacks[dc] {
				delete(seenDCRacks[dc], rack)
			}
		}

		replicas := make([]*HostInfo, 0, totalRF)
		for j := 0; j < len(tokens) && (len(replicas) < totalRF && !n.haveRF(replicasInDC)); j++ {
			// TODO: ensure we dont add the same host twice
			p := i + j
			if p >= len(tokens) {
				p -= len(tokens)
			}
			h := tokens[p].host

			dc := h.DataCenter()
			rack := h.Rack()

			rf := n.dcs[dc]
			if rf == 0 {
				// skip this DC, dont know about it or replication factor is zero
				continue
			} else if replicasInDC[dc] >= rf {
				if replicasInDC[dc] > rf {
					panic(fmt.Sprintf("replica overflow. rf=%d have=%d in dc %q", rf, replicasInDC[dc], dc))
				}

				// have enough replicas in this DC
				continue
			} else if _, ok := dcRacks[dc][rack]; !ok {
				// dont know about this rack
				continue
			}

			racks := seenDCRacks[dc]
			if _, ok := racks[rack]; ok && len(racks) == len(dcRacks[dc]) {
				// we have been through all the racks and dont have RF yet, add this
				replicas = append(replicas, h)
				replicasInDC[dc]++
			} else if !ok {
				if racks == nil {
					racks = make(map[string]struct{}, 1)
					seenDCRacks[dc] = racks
				}

				// new rack
				racks[rack] = struct{}{}
				replicas = append(replicas, h)
				r := replicasInDC[dc] + 1

				if len(racks) == len(dcRacks[dc]) {
					// if we have been through all the racks, drain the rest of the skipped
					// hosts until we have RF. The next iteration will skip in the block
					// above
					skippedHosts := skipped[dc]
					var k int
					for ; k < len(skippedHosts) && r+k < rf; k++ {
						sh := skippedHosts[k]
						replicas = append(replicas, sh)
					}
					r += k
					skipped[dc] = skippedHosts[k:]
				}
				replicasInDC[dc] = r
			} else {
				// already seen this rack, keep hold of this host incase
				// we dont get enough for rf
				skipped[dc] = append(skipped[dc], h)
			}
		}

		if len(replicas) == 0 {
			panic(fmt.Sprintf("no replicas for token: %v", th.token))
		} else if !replicas[0].Equal(th.host) {
			panic(fmt.Sprintf("first replica is not the primary replica for the token: expected %v got %v", replicas[0].ConnectAddress(), th.host.ConnectAddress()))
		}

		replicaRing = append(replicaRing, hostTokens{th.token, replicas})
	}

	dcsWithReplicas := 0
	for _, dc := range n.dcs {
		if dc > 0 {
			dcsWithReplicas++
		}
	}

	if dcsWithReplicas == len(dcRacks) && len(replicaRing) != len(tokens) {
		panic(fmt.Sprintf("token map different size to token ring: got %d expected %d", len(replicaRing), len(tokens)))
	}

	return replicaRing
}
