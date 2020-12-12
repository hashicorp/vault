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

func getReplicationFactorFromOpts(keyspace string, val interface{}) int {
	// TODO: dont really want to panic here, but is better
	// than spamming
	switch v := val.(type) {
	case int:
		if v <= 0 {
			panic(fmt.Sprintf("invalid replication_factor %d. Is the %q keyspace configured correctly?", v, keyspace))
		}
		return v
	case string:
		n, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("invalid replication_factor. Is the %q keyspace configured correctly? %v", keyspace, err))
		} else if n <= 0 {
			panic(fmt.Sprintf("invalid replication_factor %d. Is the %q keyspace configured correctly?", n, keyspace))
		}
		return n
	default:
		panic(fmt.Sprintf("unkown replication_factor type %T", v))
	}
}

func getStrategy(ks *KeyspaceMetadata) placementStrategy {
	switch {
	case strings.Contains(ks.StrategyClass, "SimpleStrategy"):
		return &simpleStrategy{rf: getReplicationFactorFromOpts(ks.Name, ks.StrategyOptions["replication_factor"])}
	case strings.Contains(ks.StrategyClass, "NetworkTopologyStrategy"):
		dcs := make(map[string]int)
		for dc, rf := range ks.StrategyOptions {
			if dc == "class" {
				continue
			}

			dcs[dc] = getReplicationFactorFromOpts(ks.Name+":dc="+dc, rf)
		}
		return &networkTopology{dcs: dcs}
	case strings.Contains(ks.StrategyClass, "LocalStrategy"):
		return nil
	default:
		// TODO: handle unknown replicas and just return the primary host for a token
		panic(fmt.Sprintf("unsupported strategy class: %v", ks.StrategyClass))
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

	if len(n.dcs) == len(dcRacks) && len(replicaRing) != len(tokens) {
		panic(fmt.Sprintf("token map different size to token ring: got %d expected %d", len(replicaRing), len(tokens)))
	}

	return replicaRing
}
