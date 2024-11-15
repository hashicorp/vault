package radix

import (
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"

	"github.com/mediocregopher/radix/v4/resp"
	"github.com/mediocregopher/radix/v4/resp/resp3"
)

// ClusterNode describes a single node in a redis cluster at a moment in time.
type ClusterNode struct {
	// older versions of redis might not actually send back the id, so it may be
	// blank
	Addr, ID string
	// start is inclusive, end is exclusive
	Slots [][2]uint16
	// address and id this node is the secondary of, if it's a secondary
	SecondaryOfAddr, SecondaryOfID string
}

// ClusterTopo describes the topology of a redis cluster at a given moment. It
// will be sorted first by slot number of each node and then by secondary
// status, so primaries will come before secondaries.
type ClusterTopo []ClusterNode

// MarshalRESP implements the resp.Marshaler interface, and will marshal the
// ClusterTopo in the same format as the return from CLUSTER SLOTS.
func (tt ClusterTopo) MarshalRESP(w io.Writer, o *resp.Opts) error {
	m := map[[2]uint16]topoSlotSet{}
	for _, t := range tt {
		for _, slots := range t.Slots {
			tss := m[slots]
			tss.slots = slots
			tss.nodes = append(tss.nodes, t)
			m[slots] = tss
		}
	}

	// we sort the topoSlotSets by their slot number so that the order is
	// deterministic, mostly so tests pass consistently, I'm not sure if actual
	// redis has any contract on the order
	allTSS := make([]topoSlotSet, 0, len(m))
	for _, tss := range m {
		allTSS = append(allTSS, tss)
	}
	sort.Slice(allTSS, func(i, j int) bool {
		return allTSS[i].slots[0] < allTSS[j].slots[0]
	})

	if err := (resp3.ArrayHeader{NumElems: len(allTSS)}).MarshalRESP(w, o); err != nil {
		return err
	}
	for _, tss := range allTSS {
		if err := tss.MarshalRESP(w, o); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalRESP implements the resp.Unmarshaler interface, but only supports
// unmarshaling the return from CLUSTER SLOTS. The unmarshaled nodes will be
// sorted before they are returned.
func (tt *ClusterTopo) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	var arrHead resp3.ArrayHeader
	if err := arrHead.UnmarshalRESP(br, o); err != nil {
		return err
	}
	slotSets := make([]topoSlotSet, arrHead.NumElems)
	for i := range slotSets {
		if err := (&(slotSets[i])).UnmarshalRESP(br, o); err != nil {
			return err
		}
	}

	nodeAddrM := map[string]ClusterNode{}
	for _, tss := range slotSets {
		for _, n := range tss.nodes {
			if existingN, ok := nodeAddrM[n.Addr]; ok {
				existingN.Slots = append(existingN.Slots, n.Slots...)
				nodeAddrM[n.Addr] = existingN
			} else {
				nodeAddrM[n.Addr] = n
			}
		}
	}

	for _, n := range nodeAddrM {
		*tt = append(*tt, n)
	}
	tt.sort()
	return nil
}

func (tt ClusterTopo) sort() {
	// first go through each node and make sure the individual slot sets are
	// sorted
	for _, node := range tt {
		sort.Slice(node.Slots, func(i, j int) bool {
			return node.Slots[i][0] < node.Slots[j][0]
		})
	}

	sort.Slice(tt, func(i, j int) bool {
		if tt[i].Slots[0] != tt[j].Slots[0] {
			return tt[i].Slots[0][0] < tt[j].Slots[0][0]
		}
		// we want secondaries to come after, which actually means they should
		// be sorted as greater
		return tt[i].SecondaryOfAddr == ""
	})

}

// Map returns the topology as a mapping of node address to its ClusterNode.
func (tt ClusterTopo) Map() map[string]ClusterNode {
	m := make(map[string]ClusterNode, len(tt))
	for _, t := range tt {
		m[t.Addr] = t
	}
	return m
}

// Primaries returns a ClusterTopo instance containing only the primary nodes
// from the ClusterTopo being called on.
func (tt ClusterTopo) Primaries() ClusterTopo {
	mtt := make(ClusterTopo, 0, len(tt))
	for _, node := range tt {
		if node.SecondaryOfAddr == "" {
			mtt = append(mtt, node)
		}
	}
	return mtt
}

// we only use this type during unmarshalling, the topo Unmarshal method will
// convert these into ClusterNodes.
type topoSlotSet struct {
	slots [2]uint16
	nodes []ClusterNode
}

func (tss topoSlotSet) MarshalRESP(w io.Writer, o *resp.Opts) error {
	var err error
	marshal := func(i interface{}) {
		if err == nil {
			err = resp3.Marshal(w, i, o)
		}
	}

	marshal(resp3.ArrayHeader{NumElems: 2 + len(tss.nodes)})
	marshal(tss.slots[0])
	marshal(tss.slots[1] - 1)

	for _, n := range tss.nodes {
		host, portStr, _ := net.SplitHostPort(n.Addr)

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}

		node := []interface{}{host, port}
		if n.ID != "" {
			node = append(node, n.ID)
		}
		marshal(node)
	}

	return resp.ErrConnUnusable(err)
}

func (tss *topoSlotSet) UnmarshalRESP(br resp.BufferedReader, o *resp.Opts) error {
	var arrHead resp3.ArrayHeader
	if err := arrHead.UnmarshalRESP(br, o); err != nil {
		return err
	}

	// first two array elements are the slot numbers. We increment the second to
	// preserve inclusive start/exclusive end, which redis doesn't
	for i := range tss.slots {
		if err := resp3.Unmarshal(br, &tss.slots[i], o); err != nil {
			return err
		}
	}
	tss.slots[1]++
	arrHead.NumElems -= len(tss.slots)

	var primaryNode ClusterNode
	for i := 0; i < arrHead.NumElems; i++ {

		var nodeArrHead resp3.ArrayHeader
		if err := nodeArrHead.UnmarshalRESP(br, o); err != nil {
			return err
		} else if nodeArrHead.NumElems < 2 {
			return fmt.Errorf("expected at least 2 array elements, got %d", nodeArrHead.NumElems)
		}

		var ip resp3.BlobString
		if err := ip.UnmarshalRESP(br, o); err != nil {
			return err
		}

		var port resp3.Number
		if err := port.UnmarshalRESP(br, o); err != nil {
			return err
		}

		nodeArrHead.NumElems -= 2

		var id resp3.BlobString
		if nodeArrHead.NumElems > 0 {
			if err := id.UnmarshalRESP(br, o); err != nil {
				return err
			}
			nodeArrHead.NumElems--
		}

		// discard anything after
		for i := 0; i < nodeArrHead.NumElems; i++ {
			if err := resp3.Unmarshal(br, nil, o); err != nil {
				return err
			}
		}

		node := ClusterNode{
			Addr:  net.JoinHostPort(ip.S, strconv.FormatInt(port.N, 10)),
			ID:    id.S,
			Slots: [][2]uint16{tss.slots},
		}

		if i == 0 {
			primaryNode = node
		} else {
			node.SecondaryOfAddr = primaryNode.Addr
			node.SecondaryOfID = primaryNode.ID
		}

		tss.nodes = append(tss.nodes, node)
	}

	return nil
}
