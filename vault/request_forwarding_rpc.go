// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"net/http"
	"os"
	"runtime/debug"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/forwarding"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/vault/replication"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type forwardedRequestRPCServer struct {
	UnimplementedRequestForwardingServer

	core                  *Core
	handler               http.Handler
	perfStandbySlots      chan struct{}
	perfStandbyRepCluster *replication.Cluster
	raftFollowerStates    *raft.FollowerStates
}

func (s *forwardedRequestRPCServer) ForwardRequest(ctx context.Context, freq *forwarding.Request) (*forwarding.Response, error) {
	// Parse an http.Request out of it
	req, err := forwarding.ParseForwardedRequest(freq)
	if err != nil {
		return nil, err
	}

	// A very dummy response writer that doesn't follow normal semantics, just
	// lets you write a status code (last written wins) and a body. But it
	// meets the interface requirements.
	w := forwarding.NewRPCResponseWriter()

	resp := &forwarding.Response{}

	runRequest := func() {
		defer func() {
			if err := recover(); err != nil {
				s.core.logger.Error("panic serving forwarded request", "path", req.URL.Path, "error", err, "stacktrace", string(debug.Stack()))
			}
		}()
		s.handler.ServeHTTP(w, req)
	}
	runRequest()
	resp.StatusCode = uint32(w.StatusCode())
	resp.Body = w.Body().Bytes()

	header := w.Header()
	if header != nil {
		resp.HeaderEntries = make(map[string]*forwarding.HeaderEntry, len(header))
		for k, v := range header {
			resp.HeaderEntries[k] = &forwarding.HeaderEntry{
				Values: v,
			}
		}
	}

	// Performance standby nodes will use this value to do wait for WALs to ship
	// in order to do a best-effort read after write guarantee
	resp.LastRemoteWal = s.core.EntLastWAL()

	return resp, nil
}

type nodeHAConnectionInfo struct {
	nodeInfo             *NodeInformation
	lastHeartbeat        time.Time
	version              string
	upgradeVersion       string
	redundancyZone       string
	localTime            time.Time
	echoDuration         time.Duration
	clockSkewMillis      int64
	replicationLagMillis int64
}

func (s *forwardedRequestRPCServer) Echo(ctx context.Context, in *EchoRequest) (*EchoReply, error) {
	incomingNodeConnectionInfo := nodeHAConnectionInfo{
		nodeInfo:             in.NodeInfo,
		lastHeartbeat:        time.Now(),
		version:              in.SdkVersion,
		upgradeVersion:       in.RaftUpgradeVersion,
		redundancyZone:       in.RaftRedundancyZone,
		localTime:            in.Now.AsTime(),
		echoDuration:         in.LastRoundtripTime.AsDuration(),
		clockSkewMillis:      in.ClockSkewMillis,
		replicationLagMillis: in.ReplicationPrimaryCanaryAgeMillis,
	}
	if in.ClusterAddr != "" {
		s.core.clusterPeerClusterAddrsCache.Set(in.ClusterAddr, incomingNodeConnectionInfo, 0)
	}

	if in.RaftAppliedIndex > 0 && len(in.RaftNodeID) > 0 && s.raftFollowerStates != nil {
		s.core.logger.Trace("forwarding RPC: echo received",
			"node_id", in.RaftNodeID,
			"applied_index", in.RaftAppliedIndex,
			"term", in.RaftTerm,
			"desired_suffrage", in.RaftDesiredSuffrage,
			"sdk_version", in.SdkVersion,
			"upgrade_version", in.RaftUpgradeVersion,
			"redundancy_zone", in.RaftRedundancyZone)

		s.raftFollowerStates.Update(&raft.EchoRequestUpdate{
			NodeID:          in.RaftNodeID,
			AppliedIndex:    in.RaftAppliedIndex,
			Term:            in.RaftTerm,
			DesiredSuffrage: in.RaftDesiredSuffrage,
			SDKVersion:      in.SdkVersion,
			UpgradeVersion:  in.RaftUpgradeVersion,
			RedundancyZone:  in.RaftRedundancyZone,
		})
	}

	reply := &EchoReply{
		Message:          "pong",
		ReplicationState: uint32(s.core.ReplicationState()),
		Now:              timestamppb.Now(),
	}

	if raftBackend := s.core.getRaftBackend(); raftBackend != nil {
		reply.RaftAppliedIndex = raftBackend.AppliedIndex()
		reply.RaftNodeID = raftBackend.NodeID()
	}

	return reply, nil
}

type forwardingClient struct {
	RequestForwardingClient
	core        *Core
	echoTicker  *time.Ticker
	echoContext context.Context
}

// NOTE: we also take advantage of gRPC's keepalive bits, but as we send data
// with these requests it's useful to keep this as well
func (c *forwardingClient) startHeartbeat() {
	go func() {
		clusterAddr := c.core.ClusterAddr()
		hostname, _ := os.Hostname()
		ni := NodeInformation{
			ApiAddr:  c.core.redirectAddr,
			Hostname: hostname,
			Mode:     "standby",
		}
		var echoDuration time.Duration
		var serverTimeDelta int64
		tick := func() {
			labels := make([]metrics.Label, 0, 1)
			defer metrics.MeasureSinceWithLabels([]string{"ha", "rpc", "client", "echo"}, time.Now(), labels)

			req := &EchoRequest{
				Message:                           "ping",
				ClusterAddr:                       clusterAddr,
				NodeInfo:                          &ni,
				SdkVersion:                        c.core.effectiveSDKVersion,
				LastRoundtripTime:                 durationpb.New(echoDuration),
				ClockSkewMillis:                   serverTimeDelta,
				ReplicationPrimaryCanaryAgeMillis: c.core.GetReplicationLagMillisIgnoreErrs(),
			}

			if raftBackend := c.core.getRaftBackend(); raftBackend != nil {
				req.RaftAppliedIndex = raftBackend.AppliedIndex()
				req.RaftNodeID = raftBackend.NodeID()
				req.RaftTerm = raftBackend.Term()
				req.RaftDesiredSuffrage = raftBackend.DesiredSuffrage()
				req.RaftRedundancyZone = raftBackend.RedundancyZone()
				req.RaftUpgradeVersion = raftBackend.UpgradeVersion()
				labels = append(labels, metrics.Label{Name: "peer_id", Value: raftBackend.NodeID()})
			}

			start := time.Now()
			req.Now = timestamppb.New(start)
			ctx, cancel := context.WithTimeout(c.echoContext, 2*time.Second)
			resp, err := c.RequestForwardingClient.Echo(ctx, req)
			cancel()

			now := time.Now()
			if err == nil {
				serverTimeDelta = resp.Now.AsTime().UnixMilli() - now.UnixMilli()
			} else {
				serverTimeDelta = 0
			}
			echoDuration = now.Sub(start)
			c.core.echoDuration.Store(echoDuration)
			c.core.activeNodeClockSkewMillis.Store(serverTimeDelta)

			if err != nil {
				metrics.IncrCounter([]string{"ha", "rpc", "client", "echo", "errors"}, 1)
				c.core.logger.Debug("forwarding: error sending echo request to active node", "error", err)
				return
			}
			c.core.rpcLastSuccessfulHeartbeat.Store(now)
			if resp == nil {
				c.core.logger.Debug("forwarding: empty echo response from active node")
				return
			}
			if resp.Message != "pong" {
				c.core.logger.Debug("forwarding: unexpected echo response from active node", "message", resp.Message)
				return
			}
			// Store the active node's replication state to display in
			// sys/health calls
			atomic.StoreUint32(c.core.activeNodeReplicationState, resp.ReplicationState)
		}

		// store a value before the first tick to indicate that we've started
		// sending heartbeats
		c.core.rpcLastSuccessfulHeartbeat.Store(time.Now())
		tick()

		for {
			select {
			case <-c.echoContext.Done():
				c.echoTicker.Stop()
				c.core.logger.Debug("forwarding: stopping heartbeating")
				atomic.StoreUint32(c.core.activeNodeReplicationState, uint32(consts.ReplicationUnknown))
				return
			case <-c.echoTicker.C:
				tick()
			}
		}
	}()
}
