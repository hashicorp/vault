// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dependency

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"log"
	"net/url"
	"regexp"
	"sort"
	"time"
)

var (
	// Ensure implements
	_ Dependency = (*ListPeeringQuery)(nil)

	// ListPeeringQueryRe is the regular expression to use.
	ListPeeringQueryRe = regexp.MustCompile(`\A` + queryRe + `\z`)
)

func init() {
	gob.Register([]*Peering{})
	gob.Register([]*PeeringStreamStatus{})
	gob.Register([]*PeeringRemoteInfo{})
}

// ListPeeringQuery fetches all peering for a Consul cluster.
// https://developer.hashicorp.com/consul/api-docs/peering#list-all-peerings
type ListPeeringQuery struct {
	stopCh chan struct{}

	partition string
}

// Peering represent the response of the Consul peering API.
type Peering struct {
	ID                  string
	Name                string
	Partition           string
	Meta                map[string]string
	PeeringState        string
	PeerID              string
	PeerServerName      string
	PeerServerAddresses []string
	StreamStatus        PeeringStreamStatus
	Remote              PeeringRemoteInfo
}

type PeeringStreamStatus struct {
	ImportedServices []string
	ExportedServices []string
	LastHeartbeat    *time.Time
	LastReceive      *time.Time
	LastSend         *time.Time
}

type PeeringRemoteInfo struct {
	Partition  string
	Datacenter string
}

func NewListPeeringQuery(s string) (*ListPeeringQuery, error) {
	if s != "" && !ListPeeringQueryRe.MatchString(s) {
		return nil, fmt.Errorf("list.peering: invalid format: %q", s)
	}

	m := regexpMatch(ListPeeringQueryRe, s)

	queryParams, err := GetConsulQueryOpts(m, "list.peering")
	if err != nil {
		return nil, err
	}

	return &ListPeeringQuery{
		stopCh:    make(chan struct{}, 1),
		partition: queryParams.Get(QueryPartition),
	}, nil
}

func (l *ListPeeringQuery) Fetch(clients *ClientSet, opts *QueryOptions) (interface{}, *ResponseMetadata, error) {
	select {
	case <-l.stopCh:
		return nil, nil, ErrStopped
	default:
	}

	opts = opts.Merge(&QueryOptions{
		ConsulPartition: l.partition,
	})

	log.Printf("[TRACE] %s: GET %s", l, &url.URL{
		Path:     "/v1/peerings",
		RawQuery: opts.String(),
	})

	// list peering is a blocking API, so making sure the ctx passed while calling it
	// times out after the default wait time.
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout)
	defer cancel()

	p, meta, err := clients.Consul().Peerings().List(ctx, opts.ToConsulOpts())
	if err != nil {
		return nil, nil, errors.Wrap(err, l.String())
	}

	log.Printf("[TRACE] %s: returned %d results", l, len(p))

	peers := make([]*Peering, 0, len(p))
	for _, peering := range p {
		peers = append(peers, toPeering(peering))
	}

	// sort so that the result is deterministic
	sort.Stable(ByPeer(peers))

	rm := &ResponseMetadata{
		LastIndex:   meta.LastIndex,
		LastContact: meta.LastContact,
	}

	return peers, rm, nil
}

func toPeering(p *api.Peering) *Peering {
	return &Peering{
		ID:                  p.ID,
		Name:                p.Name,
		Partition:           p.Partition,
		Meta:                p.Meta,
		PeeringState:        string(p.State),
		PeerID:              p.PeerID,
		PeerServerName:      p.PeerServerName,
		PeerServerAddresses: p.PeerServerAddresses,
		StreamStatus: PeeringStreamStatus{
			ImportedServices: p.StreamStatus.ImportedServices,
			ExportedServices: p.StreamStatus.ExportedServices,
			LastHeartbeat:    p.StreamStatus.LastHeartbeat,
			LastReceive:      p.StreamStatus.LastReceive,
			LastSend:         p.StreamStatus.LastSend,
		},
		Remote: PeeringRemoteInfo{
			Partition:  p.Remote.Partition,
			Datacenter: p.Remote.Datacenter,
		},
	}
}

func (l *ListPeeringQuery) String() string {
	partitionStr := l.partition

	if len(partitionStr) > 0 {
		partitionStr = fmt.Sprintf("?partition=%s", partitionStr)
	} else {
		return "list.peerings"
	}

	return fmt.Sprintf("list.peerings%s", partitionStr)
}

func (l *ListPeeringQuery) Stop() {
	close(l.stopCh)
}

func (l *ListPeeringQuery) Type() Type {
	return TypeConsul
}

func (l *ListPeeringQuery) CanShare() bool {
	return false
}

// ByPeer is a sortable list of peerings in this order:
// 1. State
// 2. Partition
// 3. Name
type ByPeer []*Peering

func (p ByPeer) Len() int      { return len(p) }
func (p ByPeer) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Less if peer names are cluster-2, cluster-12, cluster-1
// our sorting will be cluster-1, cluster-12, cluster-2
func (p ByPeer) Less(i, j int) bool {
	if p[i].PeeringState == p[j].PeeringState {
		if p[i].Partition == p[j].Partition {
			return p[i].Name < p[j].Name
		}
		return p[i].Partition < p[j].Partition
	}
	return p[i].PeeringState < p[j].PeeringState
}
