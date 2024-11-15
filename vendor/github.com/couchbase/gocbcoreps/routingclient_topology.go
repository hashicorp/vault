package gocbcoreps

import (
	"context"
	"time"

	"github.com/couchbase/goprotostellar/genproto/routing_v1"
	"go.uber.org/zap"
)

func (p *RoutingClient) translateTopology(t *routing_v1.WatchRoutingResponse) *Topology {
	nodes := make([]*Node, len(t.Endpoints))
	for psEpIdx, psEp := range t.Endpoints {
		node := &Node{
			NodeID:      psEp.Id,
			ServerGroup: psEp.ServerGroup,
		}
		nodes[psEpIdx] = node
	}

	var vbRouting *VbucketRouting
	psVbRouting := t.GetVbucketDataRouting()
	if vbRouting != nil {
		dataNodes := make([]*DataNode, len(psVbRouting.Endpoints))
		for psDataEpIdx, psDataEp := range psVbRouting.Endpoints {
			dataNode := &DataNode{
				Node:          nodes[psDataEp.EndpointIdx],
				LocalVbuckets: psDataEp.LocalVbuckets,
				GroupVbuckets: psDataEp.GroupVbuckets,
			}
			dataNodes[psDataEpIdx] = dataNode
		}
		vbRouting = &VbucketRouting{
			NumVbuckets: uint(vbRouting.NumVbuckets),
			Nodes:       dataNodes,
		}
	}

	return &Topology{
		Revision:       t.Revision,
		Nodes:          nodes,
		VbucketRouting: vbRouting,
	}
}

func (p *RoutingClient) watchTopology(ctx context.Context, bucketName *string) (<-chan *Topology, error) {
	// TODO(brett19): PSProvider.watchTopology is something that probably belongs in the client.
	// Putting it in the client will be helpful because it's likely used in multipled places, and the
	// client itself needs to maintain knowledge of the topology anyways, we can perform coalescing
	// of those watchers in the client (to reduce the number of watchers).

	b := exponentialBackoff(0, 0, 0)
	var numRetries uint32

	routingStream, err := p.RoutingV1().WatchRouting(ctx, &routing_v1.WatchRoutingRequest{
		BucketName: bucketName,
	})
	if err != nil {
		return nil, err
	}

	routingResp, err := routingStream.Recv()
	if err != nil {
		return nil, err
	}

	outputCh := make(chan *Topology)
	go func() {
		outputCh <- p.translateTopology(routingResp)

	MainLoop:
		for {
			routingStream, err := p.RoutingV1().WatchRouting(ctx, &routing_v1.WatchRoutingRequest{
				BucketName: bucketName,
			})
			if err != nil {
				// TODO(brett19): Implement better error handling here...
				p.logger.Error("failed to watch routing", zap.Error(err))
				numRetries++

				select {
				case <-time.After(b(numRetries)):
					continue
				case <-ctx.Done():
					break MainLoop
				}
			}
			numRetries = 0

			for {
				routingResp, err := routingStream.Recv()
				if err != nil {
					p.logger.Error("failed to recv updated topology", zap.Error(err))
					break
				}

				outputCh <- p.translateTopology(routingResp)
			}
		}
	}()

	return outputCh, nil
}

func (p *RoutingClient) WatchTopology(ctx context.Context, bucketName string) (<-chan *Topology, error) {
	// TODO(brett19): Remove this pointer shenanigans
	var bucketNamePtr *string
	if bucketName != "" {
		bucketNamePtr = &bucketName
	}
	return p.watchTopology(ctx, bucketNamePtr)
}
