package gocql

import (
	"fmt"
	"strconv"
	"strings"
)

type placementStrategy interface {
	replicaMap(hosts []*HostInfo, tokens []hostToken) map[token][]*HostInfo
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

func (s *simpleStrategy) replicaMap(_ []*HostInfo, tokens []hostToken) map[token][]*HostInfo {
	tokenRing := make(map[token][]*HostInfo, len(tokens))

	for i, th := range tokens {
		replicas := make([]*HostInfo, 0, s.rf)
		for j := 0; j < len(tokens) && len(replicas) < s.rf; j++ {
			// TODO: need to ensure we dont add the same hosts twice
			h := tokens[(i+j)%len(tokens)]
			replicas = append(replicas, h.host)
		}
		tokenRing[th.token] = replicas
	}

	return tokenRing
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

func (n *networkTopology) replicaMap(hosts []*HostInfo, tokens []hostToken) map[token][]*HostInfo {
	dcRacks := make(map[string]map[string]struct{})

	for _, h := range hosts {
		dc := h.DataCenter()
		rack := h.Rack()

		racks, ok := dcRacks[dc]
		if !ok {
			racks = make(map[string]struct{})
			dcRacks[dc] = racks
		}
		racks[rack] = struct{}{}
	}

	tokenRing := make(map[token][]*HostInfo, len(tokens))

	var totalRF int
	for _, rf := range n.dcs {
		totalRF += rf
	}

	for i, th := range tokens {
		// number of replicas per dc
		// TODO: recycle these
		replicasInDC := make(map[string]int, len(n.dcs))
		// dc -> racks
		seenDCRacks := make(map[string]map[string]struct{}, len(n.dcs))
		// skipped hosts in a dc
		skipped := make(map[string][]*HostInfo, len(n.dcs))

		replicas := make([]*HostInfo, 0, totalRF)
		for j := 0; j < len(tokens) && !n.haveRF(replicasInDC); j++ {
			// TODO: ensure we dont add the same host twice
			h := tokens[(i+j)%len(tokens)].host

			dc := h.DataCenter()
			rack := h.Rack()

			rf, ok := n.dcs[dc]
			if !ok {
				// skip this DC, dont know about it
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
			} else if len(replicas) >= totalRF {
				if replicasInDC[dc] > rf {
					panic(fmt.Sprintf("replica overflow. total rf=%d have=%d", totalRF, len(replicas)))
				}

				// we now have enough replicas
				break
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
				replicasInDC[dc]++

				if len(racks) == len(dcRacks[dc]) {
					// if we have been through all the racks, drain the rest of the skipped
					// hosts until we have RF. The next iteration will skip in the block
					// above
					skippedHosts := skipped[dc]
					var k int
					for ; k < len(skippedHosts) && replicasInDC[dc] < rf; k++ {
						sh := skippedHosts[k]
						replicas = append(replicas, sh)
						replicasInDC[dc]++
					}
					skipped[dc] = skippedHosts[k:]
				}
			} else {
				// already seen this rack, keep hold of this host incase
				// we dont get enough for rf
				skipped[dc] = append(skipped[dc], h)
			}
		}

		if len(replicas) == 0 || replicas[0] != th.host {
			panic("first replica is not the primary replica for the token")
		}

		tokenRing[th.token] = replicas
	}

	if len(tokenRing) != len(tokens) {
		panic(fmt.Sprintf("token map different size to token ring: got %d expected %d", len(tokenRing), len(tokens)))
	}

	return tokenRing
}
