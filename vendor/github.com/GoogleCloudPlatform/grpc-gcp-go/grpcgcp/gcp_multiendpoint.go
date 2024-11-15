/*
 *
 * Copyright 2023 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package grpcgcp

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp/multiendpoint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	pb "github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp/grpc_gcp"
)

var gmeCounter uint32

type contextMEKey int

var meKey contextMEKey

// NewMEContext returns a new Context that carries Multiendpoint name.
func NewMEContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, meKey, name)
}

// FromMEContext returns the MultiEndpoint name stored in ctx, if any.
func FromMEContext(ctx context.Context) (string, bool) {
	name, ok := ctx.Value(meKey).(string)
	return name, ok
}

// GCPMultiEndpoint holds the state of MultiEndpoints-enabled gRPC client connection.
//
// The purposes of GCPMultiEndpoint are:
//
//   - Fallback to an alternative endpoint (host:port) of a gRPC service when the original
//     endpoint is completely unavailable.
//   - Be able to route an RPC call to a specific group of endpoints.
//   - Be able to reconfigure endpoints in runtime.
//
// A group of endpoints is called a [multiendpoint.MultiEndpoint] and is essentially a list of endpoints
// where priority is defined by the position in the list with the first endpoint having top
// priority. A MultiEndpoint tracks endpoints' availability. When a MultiEndpoint is picked for an
// RPC call, it picks the top priority endpoint that is currently available. More information on the
// [multiendpoint.MultiEndpoint].
//
// GCPMultiEndpoint can have one or more MultiEndpoint identified by its name -- arbitrary
// string provided in the [GCPMultiEndpointOptions] when configuring MultiEndpoints. This name
// can be used to route an RPC call to this MultiEndpoint by using the [NewMEContext].
//
// GCPMultiEndpoint uses [GCPMultiEndpointOptions] for initial configuration.
// An updated configuration can be provided at any time later using [UpdateMultiEndpoints].
//
// Example:
//
// Let's assume we have a service with read and write operations and the following backends:
//
//   - service.example.com -- the main set of backends supporting all operations
//   - service-fallback.example.com -- read-write replica supporting all operations
//   - ro-service.example.com -- read-only replica supporting only read operations
//
// Example configuration:
//
//   - MultiEndpoint named "default" with endpoints:
//
//     1. service.example.com:443
//
//     2. service-fallback.example.com:443
//
//   - MultiEndpoint named "read" with endpoints:
//
//     1. ro-service.example.com:443
//
//     2. service-fallback.example.com:443
//
//     3. service.example.com:443
//
// With the configuration above GCPMultiEndpoint will use the "default" MultiEndpoint by
// default. It means that RPC calls by default will use the main endpoint and if it is not available
// then the read-write replica.
//
// To offload some read calls to the read-only replica we can specify "read" MultiEndpoint in the
// context. Then these calls will use the read-only replica endpoint and if it is not available
// then the read-write replica and if it is also not available then the main endpoint.
//
// GCPMultiEndpoint creates a [grpcgcp] connection pool for every unique
// endpoint. For the example above three connection pools will be created.
//
// [GCPMultiEndpoint] implements [grpc.ClientConnInterface] and can be used
// as a [grpc.ClientConn] when creating gRPC clients.
type GCPMultiEndpoint struct {
	mu sync.RWMutex

	defaultName string
	mes         map[string]multiendpoint.MultiEndpoint
	pools       map[string]*monitoredConn
	opts        []grpc.DialOption
	gcpConfig   *pb.ApiConfig
	dialFunc    func(ctx context.Context, target string, dopts ...grpc.DialOption) (*grpc.ClientConn, error)
	log         grpclog.LoggerV2

	grpc.ClientConnInterface
}

// Make sure GCPMultiEndpoint implements grpc.ClientConnInterface.
var _ grpc.ClientConnInterface = (*GCPMultiEndpoint)(nil)

func (gme *GCPMultiEndpoint) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	return gme.pickConn(ctx).Invoke(ctx, method, args, reply, opts...)
}

func (gme *GCPMultiEndpoint) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return gme.pickConn(ctx).NewStream(ctx, desc, method, opts...)
}

func (gme *GCPMultiEndpoint) pickConn(ctx context.Context) *grpc.ClientConn {
	name, ok := FromMEContext(ctx)
	me, ook := gme.mes[name]
	if !ok || !ook {
		me = gme.mes[gme.defaultName]
	}
	return gme.pools[me.Current()].conn
}

func (gme *GCPMultiEndpoint) Close() error {
	var errs multiError
	for e, mc := range gme.pools {
		mc.stopMonitoring()
		if err := mc.conn.Close(); err != nil {
			errs = append(errs, err)
			gme.log.Errorf("error while closing the pool for %q endpoint: %v", e, err)
		}
		if gme.log.V(FINE) {
			gme.log.Infof("closed channel pool for %q endpoint.", e)
		}
	}
	return errs.Combine()
}

func (gme *GCPMultiEndpoint) GCPConfig() *pb.ApiConfig {
	return proto.Clone(gme.gcpConfig).(*pb.ApiConfig)
}

// GCPMultiEndpointOptions holds options to construct a MultiEndpoints-enabled gRPC client
// connection.
type GCPMultiEndpointOptions struct {
	// Regular gRPC-GCP configuration to be applied to every endpoint.
	GRPCgcpConfig *pb.ApiConfig
	// Map of MultiEndpoints where key is the MultiEndpoint name.
	MultiEndpoints map[string]*multiendpoint.MultiEndpointOptions
	// Name of the default MultiEndpoint.
	Default string
	// Func to dial grpc ClientConn.
	DialFunc func(ctx context.Context, target string, dopts ...grpc.DialOption) (*grpc.ClientConn, error)
}

// NewGCPMultiEndpoint creates new [GCPMultiEndpoint] -- MultiEndpoints-enabled gRPC client
// connection.
//
// Deprecated: use NewGCPMultiEndpoint.
func NewGcpMultiEndpoint(meOpts *GCPMultiEndpointOptions, opts ...grpc.DialOption) (*GCPMultiEndpoint, error) {
	return NewGCPMultiEndpoint(meOpts, opts...)
}

// NewGCPMultiEndpoint creates new [GCPMultiEndpoint] -- MultiEndpoints-enabled gRPC client
// connection.
//
// [GCPMultiEndpoint] implements [grpc.ClientConnInterface] and can be used
// as a [grpc.ClientConn] when creating gRPC clients.
func NewGCPMultiEndpoint(meOpts *GCPMultiEndpointOptions, opts ...grpc.DialOption) (*GCPMultiEndpoint, error) {
	// Read config, create multiendpoints and pools.
	o, err := makeOpts(meOpts, opts)
	if err != nil {
		return nil, err
	}
	gme := &GCPMultiEndpoint{
		mes:         make(map[string]multiendpoint.MultiEndpoint),
		pools:       make(map[string]*monitoredConn),
		defaultName: meOpts.Default,
		opts:        o,
		gcpConfig:   proto.Clone(meOpts.GRPCgcpConfig).(*pb.ApiConfig),
		dialFunc:    meOpts.DialFunc,
		log:         NewGCPLogger(compLogger, fmt.Sprintf("[GCPMultiEndpoint #%d]", atomic.AddUint32(&gmeCounter, 1))),
	}
	if gme.dialFunc == nil {
		gme.dialFunc = func(_ context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
			return grpc.Dial(target, opts...)
		}
	}
	if err := gme.UpdateMultiEndpoints(meOpts); err != nil {
		return nil, err
	}
	return gme, nil
}

func makeOpts(meOpts *GCPMultiEndpointOptions, opts []grpc.DialOption) ([]grpc.DialOption, error) {
	grpcGCPjsonConfig, err := protojson.Marshal(meOpts.GRPCgcpConfig)
	if err != nil {
		return nil, err
	}
	o := append([]grpc.DialOption{}, opts...)
	o = append(o, []grpc.DialOption{
		grpc.WithDisableServiceConfig(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingConfig": [{"%s":%s}]}`, Name, string(grpcGCPjsonConfig))),
		grpc.WithChainUnaryInterceptor(GCPUnaryClientInterceptor),
		grpc.WithChainStreamInterceptor(GCPStreamClientInterceptor),
	}...)

	return o, nil
}

type monitoredConn struct {
	endpoint string
	conn     *grpc.ClientConn
	gme      *GCPMultiEndpoint
	cancel   context.CancelFunc
}

func newMonitoredConn(endpoint string, conn *grpc.ClientConn, gme *GCPMultiEndpoint) (mc *monitoredConn) {
	ctx, cancel := context.WithCancel(context.Background())
	mc = &monitoredConn{
		endpoint: endpoint,
		conn:     conn,
		gme:      gme,
		cancel:   cancel,
	}
	go mc.monitor(ctx)
	return
}

func (mc *monitoredConn) notify(state connectivity.State) {
	if mc.gme.log.V(FINE) {
		mc.gme.log.Infof("%q endpoint state changed to %v", mc.endpoint, state)
	}
	// Inform all multiendpoints.
	mc.gme.mu.RLock()
	for _, me := range mc.gme.mes {
		me.SetEndpointAvailability(mc.endpoint, state == connectivity.Ready)
	}
	mc.gme.mu.RUnlock()
}

func (mc *monitoredConn) monitor(ctx context.Context) {
	for {
		currentState := mc.conn.GetState()
		mc.notify(currentState)
		if !mc.conn.WaitForStateChange(ctx, currentState) {
			break
		}
	}
}

func (mc *monitoredConn) stopMonitoring() {
	mc.cancel()
}

// UpdateMultiEndpoints reconfigures MultiEndpoints.
//
// MultiEndpoints are matched with the current ones by name.
//
//   - If a current MultiEndpoint is missing in the updated list, the MultiEndpoint will be
//     removed.
//   - A new MultiEndpoint will be created for every new name in the list.
//   - For an existing MultiEndpoint only its endpoints will be updated (no recovery timeout
//     change).
//
// Endpoints are matched by the endpoint address (usually in the form of address:port).
//
//   - If an existing endpoint is not used by any MultiEndpoint in the updated list, then the
//     connection poll for this endpoint will be shutdown.
//   - A connection pool will be created for every new endpoint.
//   - For an existing endpoint nothing will change (the connection pool will not be re-created,
//     thus no connection credentials change, nor connection configuration change).
func (gme *GCPMultiEndpoint) UpdateMultiEndpoints(meOpts *GCPMultiEndpointOptions) error {
	gme.mu.Lock()
	defer gme.mu.Unlock()
	if _, ok := meOpts.MultiEndpoints[meOpts.Default]; !ok {
		return fmt.Errorf("default MultiEndpoint %q missing options", meOpts.Default)
	}

	validPools := make(map[string]bool)
	for _, meo := range meOpts.MultiEndpoints {
		for _, e := range meo.Endpoints {
			validPools[e] = true
		}
	}

	// Add missing pools.
	for e := range validPools {
		if _, ok := gme.pools[e]; !ok {
			// This creates a ClientConn with the gRPC-GCP balancer managing connection pool.
			conn, err := gme.dialFunc(context.Background(), e, gme.opts...)
			if err != nil {
				return err
			}
			if gme.log.V(FINE) {
				gme.log.Infof("created new channel pool for %q endpoint.", e)
			}
			gme.pools[e] = newMonitoredConn(e, conn, gme)
		}
	}

	// Add new multi-endpoints and update existing.
	for name, meo := range meOpts.MultiEndpoints {
		if me, ok := gme.mes[name]; ok {
			// Updating existing MultiEndpoint.
			me.SetEndpoints(meo.Endpoints)
			continue
		}

		// Add new MultiEndpoint.
		if gme.log.V(FINE) {
			gme.log.Infof("creating new %q multiendpoint.", name)
		}
		me, err := multiendpoint.NewMultiEndpoint(meo)
		if err != nil {
			return err
		}
		gme.mes[name] = me
	}
	gme.defaultName = meOpts.Default

	// Remove obsolete MultiEndpoints.
	for name := range gme.mes {
		if _, ok := meOpts.MultiEndpoints[name]; !ok {
			delete(gme.mes, name)
			if gme.log.V(FINE) {
				gme.log.Infof("removed obsolete %q multiendpoint.", name)
			}
		}
	}

	// Remove obsolete pools.
	for e, mc := range gme.pools {
		if _, ok := validPools[e]; !ok {
			if err := mc.conn.Close(); err != nil {
				gme.log.Errorf("error while closing the pool for %q endpoint: %v", e, err)
			}
			if gme.log.V(FINE) {
				gme.log.Infof("closed channel pool for %q endpoint.", e)
			}
			mc.stopMonitoring()
			delete(gme.pools, e)
		}
	}

	// Trigger status update.
	for e, mc := range gme.pools {
		s := mc.conn.GetState()
		for _, me := range gme.mes {
			me.SetEndpointAvailability(e, s == connectivity.Ready)
		}
	}
	return nil
}

type multiError []error

func (m multiError) Error() string {
	s, n := "", 0
	for _, e := range m {
		if e != nil {
			if n == 0 {
				s = e.Error()
			}
			n++
		}
	}
	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, n-1)
}

func (m multiError) Combine() error {
	if len(m) == 0 {
		return nil
	}

	return m
}
