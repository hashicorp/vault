package vault

import (
	"context"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/forwarding"
	"github.com/hashicorp/vault/vault/replication"
	cache "github.com/patrickmn/go-cache"
)

type forwardedRequestRPCServer struct {
	core                  *Core
	handler               http.Handler
	perfStandbySlots      chan struct{}
	perfStandbyRepCluster *replication.Cluster
	perfStandbyCache      *cache.Cache
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
			// Logic here comes mostly from the Go source code
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				s.core.logger.Error("panic serving forwarded request", "path", req.URL.Path, "error", err, "stacktrace", string(buf))
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
	// in order to do a best-effort read after write gurantee
	resp.LastRemoteWal = LastWAL(s.core)

	return resp, nil
}

func (s *forwardedRequestRPCServer) Echo(ctx context.Context, in *EchoRequest) (*EchoReply, error) {
	if in.ClusterAddr != "" {
		s.core.clusterPeerClusterAddrsCache.Set(in.ClusterAddr, nil, 0)
	}
	return &EchoReply{
		Message:          "pong",
		ReplicationState: uint32(s.core.ReplicationState()),
	}, nil
}

type forwardingClient struct {
	RequestForwardingClient

	core *Core

	echoTicker  *time.Ticker
	echoContext context.Context
}

// NOTE: we also take advantage of gRPC's keepalive bits, but as we send data
// with these requests it's useful to keep this as well
func (c *forwardingClient) startHeartbeat() {
	go func() {
		tick := func() {
			c.core.stateLock.RLock()
			clusterAddr := c.core.clusterAddr
			c.core.stateLock.RUnlock()

			ctx, cancel := context.WithTimeout(c.echoContext, 2*time.Second)
			resp, err := c.RequestForwardingClient.Echo(ctx, &EchoRequest{
				Message:     "ping",
				ClusterAddr: clusterAddr,
			})
			cancel()
			if err != nil {
				c.core.logger.Debug("forwarding: error sending echo request to active node", "error", err)
				return
			}
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
