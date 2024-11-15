/*
 *
 * Copyright 2018 gRPC authors.
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
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	pb "github.com/GoogleCloudPlatform/grpc-gcp-go/grpcgcp/grpc_gcp"
)

var _ balancer.Balancer = (*gcpBalancer)(nil) // Ensure gcpBalancer implements Balancer

const (
	// Name is the name of grpc_gcp balancer.
	Name = "grpc_gcp"

	healthCheckEnabled = true
	defaultMinSize     = 1
	defaultMaxSize     = 4
	defaultMaxStreams  = 100
)

func init() {
	balancer.Register(newBuilder())
}

type gcpBalancerBuilder struct {
	balancer.ConfigParser
}

type GCPBalancerConfig struct {
	serviceconfig.LoadBalancingConfig
	*pb.ApiConfig
}

func (bb *gcpBalancerBuilder) Build(
	cc balancer.ClientConn,
	opt balancer.BuildOptions,
) balancer.Balancer {
	gb := &gcpBalancer{
		cc:               cc,
		methodCfg:        make(map[string]*pb.AffinityConfig),
		affinityMap:      make(map[string]balancer.SubConn),
		fallbackMap:      make(map[string]balancer.SubConn),
		scRefs:           make(map[balancer.SubConn]*subConnRef),
		scStates:         make(map[balancer.SubConn]connectivity.State),
		refreshingScRefs: make(map[balancer.SubConn]*subConnRef),
		scRefList:        []*subConnRef{},
		rrRefId:          ^uint32(0),
		csEvltr:          &connectivityStateEvaluator{},
		// Initialize picker to a picker that always return
		// ErrNoSubConnAvailable, because when state of a SubConn changes, we
		// may call UpdateBalancerState with this picker.
		picker: newErrPicker(balancer.ErrNoSubConnAvailable),
	}
	gb.log = NewGCPLogger(compLogger, fmt.Sprintf("[gcpBalancer %p]", gb))
	return gb
}

func (*gcpBalancerBuilder) Name() string {
	return Name
}

// ParseConfig converts raw json config into GCPBalancerConfig.
// This is called by ClientConn on any load balancer config update.
// After parsing the config, ClientConn calls UpdateClientConnState passing the config.
func (*gcpBalancerBuilder) ParseConfig(j json.RawMessage) (serviceconfig.LoadBalancingConfig, error) {
	c := &GCPBalancerConfig{
		ApiConfig: &pb.ApiConfig{},
	}
	err := protojson.Unmarshal(j, c)
	return c, err
}

// newBuilder creates a new grpcgcp balancer builder.
func newBuilder() balancer.Builder {
	return &gcpBalancerBuilder{}
}

// connectivityStateEvaluator gets updated by addrConns when their
// states transition, based on which it evaluates the state of
// ClientConn.
type connectivityStateEvaluator struct {
	numReady            uint64 // Number of addrConns in ready state.
	numConnecting       uint64 // Number of addrConns in connecting state.
	numTransientFailure uint64 // Number of addrConns in transientFailure.
}

// recordTransition records state change happening in every subConn and based on
// that it evaluates what aggregated state should be.
// It can only transition between Ready, Connecting and TransientFailure. Other states,
// Idle and Shutdown are transitioned into by ClientConn; in the beginning of the connection
// before any subConn is created ClientConn is in idle state. In the end when ClientConn
// closes it is in Shutdown state.
//
// recordTransition should only be called synchronously from the same goroutine.
func (cse *connectivityStateEvaluator) recordTransition(
	oldState,
	newState connectivity.State,
) connectivity.State {
	// Update counters.
	for idx, state := range []connectivity.State{oldState, newState} {
		updateVal := 2*uint64(idx) - 1 // -1 for oldState and +1 for new.
		switch state {
		case connectivity.Ready:
			cse.numReady += updateVal
		case connectivity.Connecting:
			cse.numConnecting += updateVal
		case connectivity.TransientFailure:
			cse.numTransientFailure += updateVal
		}
	}

	// Evaluate.
	if cse.numReady > 0 {
		return connectivity.Ready
	}
	if cse.numConnecting > 0 {
		return connectivity.Connecting
	}
	return connectivity.TransientFailure
}

// subConnRef keeps reference to the real SubConn with its
// connectivity state, affinity count and streams count.
type subConnRef struct {
	subConn     balancer.SubConn
	stateSignal chan struct{} // This channel is closed and re-created when subConn or its state changes.
	affinityCnt int32         // Keeps track of the number of keys bound to the subConn.
	streamsCnt  int32         // Keeps track of the number of streams opened on the subConn.
	lastResp    time.Time     // Timestamp of the last response from the server.
	deCalls     uint32        // Keeps track of deadline exceeded calls since last response.
	refreshing  bool          // If this subconn is in the process of refreshing.
	refreshCnt  uint32        // Number of refreshes since last response.
}

func (ref *subConnRef) getAffinityCnt() int32 {
	return atomic.LoadInt32(&ref.affinityCnt)
}

func (ref *subConnRef) getStreamsCnt() int32 {
	return atomic.LoadInt32(&ref.streamsCnt)
}

func (ref *subConnRef) affinityIncr() {
	atomic.AddInt32(&ref.affinityCnt, 1)
}

func (ref *subConnRef) affinityDecr() {
	atomic.AddInt32(&ref.affinityCnt, -1)
}

func (ref *subConnRef) streamsIncr() {
	atomic.AddInt32(&ref.streamsCnt, 1)
}

func (ref *subConnRef) streamsDecr() {
	atomic.AddInt32(&ref.streamsCnt, -1)
}

func (ref *subConnRef) deCallsInc() uint32 {
	return atomic.AddUint32(&ref.deCalls, 1)
}

func (ref *subConnRef) gotResp() {
	ref.lastResp = time.Now()
	atomic.StoreUint32(&ref.deCalls, 0)
	ref.refreshCnt = 0
}

type gcpBalancer struct {
	cfg       *GCPBalancerConfig
	methodCfg map[string]*pb.AffinityConfig

	addrs   []resolver.Address
	cc      balancer.ClientConn
	csEvltr *connectivityStateEvaluator
	state   connectivity.State

	mu          sync.RWMutex
	affinityMap map[string]balancer.SubConn
	fallbackMap map[string]balancer.SubConn
	scStates    map[balancer.SubConn]connectivity.State
	scRefs      map[balancer.SubConn]*subConnRef
	scRefList   []*subConnRef
	rrRefId     uint32

	// Map from a fresh SubConn to the subConnRef where we want to refresh subConn.
	refreshingScRefs map[balancer.SubConn]*subConnRef
	// Unresponsive detection enabled flag.
	unresponsiveDetection bool

	picker balancer.Picker
	log    grpclog.LoggerV2
}

func (gb *gcpBalancer) initializeConfig(cfg *GCPBalancerConfig) {
	gb.cfg = &GCPBalancerConfig{
		ApiConfig: &pb.ApiConfig{
			ChannelPool: &pb.ChannelPoolConfig{},
		},
	}
	if cfg != nil && cfg.ApiConfig != nil {
		gb.cfg = &GCPBalancerConfig{
			ApiConfig: proto.Clone(cfg.ApiConfig).(*pb.ApiConfig),
		}
	}

	if gb.cfg.GetChannelPool() == nil {
		gb.cfg.ChannelPool = &pb.ChannelPoolConfig{}
	}
	cp := gb.cfg.GetChannelPool()
	if cp.GetMinSize() == 0 {
		cp.MinSize = defaultMinSize
	}
	if cp.GetMaxSize() == 0 {
		cp.MaxSize = defaultMaxSize
	}
	if cp.GetMaxConcurrentStreamsLowWatermark() == 0 {
		cp.MaxConcurrentStreamsLowWatermark = defaultMaxStreams
	}
	mp := make(map[string]*pb.AffinityConfig)
	methodCfgs := gb.cfg.GetMethod()
	for _, methodCfg := range methodCfgs {
		methodNames := methodCfg.GetName()
		affinityCfg := methodCfg.GetAffinity()
		if methodNames != nil && affinityCfg != nil {
			for _, method := range methodNames {
				mp[method] = affinityCfg
			}
		}
	}
	gb.methodCfg = mp
	gb.unresponsiveDetection = cp.GetUnresponsiveCalls() > 0 && cp.GetUnresponsiveDetectionMs() > 0
	gb.enforceMinSize()
}

func (gb *gcpBalancer) enforceMinSize() {
	for len(gb.scRefs) < int(gb.cfg.GetChannelPool().GetMinSize()) {
		gb.addSubConn()
	}
}

func (gb *gcpBalancer) UpdateClientConnState(ccs balancer.ClientConnState) error {
	gb.mu.Lock()
	defer gb.mu.Unlock()
	addrs := ccs.ResolverState.Addresses
	if gb.log.V(FINE) {
		gb.log.Infoln("got new resolved addresses: ", addrs, " and balancer config: ", ccs.BalancerConfig)
	}
	gb.addrs = addrs
	if gb.cfg == nil {
		cfg, ok := ccs.BalancerConfig.(*GCPBalancerConfig)
		if !ok && ccs.BalancerConfig != nil {
			return fmt.Errorf("provided config is not GCPBalancerConfig: %v", ccs.BalancerConfig)
		}
		gb.initializeConfig(cfg)
	}

	if len(gb.scRefs) == 0 {
		gb.newSubConn()
		return nil
	}

	for _, scRef := range gb.scRefs {
		// TODO(weiranf): update streams count when new addrs resolved?
		scRef.subConn.UpdateAddresses(addrs)
		scRef.subConn.Connect()
	}

	return nil
}

func (gb *gcpBalancer) ResolverError(err error) {
	gb.log.Warningf("ResolverError: %v", err)
}

// check current connection pool size
func (gb *gcpBalancer) getConnectionPoolSize() int {
	// TODO(golobokov): replace this with locked increase of subconns.
	gb.mu.Lock()
	defer gb.mu.Unlock()
	return len(gb.scRefs)
}

// newSubConn creates a new SubConn using cc.NewSubConn and initialize the subConnRef
// if none of the subconns are in the Connecting state.
func (gb *gcpBalancer) newSubConn() {
	gb.mu.Lock()
	defer gb.mu.Unlock()

	// there are chances the newly created subconns are still connecting,
	// we can wait on those new subconns.
	for _, scState := range gb.scStates {
		if scState == connectivity.Connecting || scState == connectivity.Idle {
			return
		}
	}
	gb.addSubConn()
}

// addSubConn creates a new SubConn using cc.NewSubConn and initialize the subConnRef.
// Must be called holding the mutex lock.
func (gb *gcpBalancer) addSubConn() {
	sc, err := gb.cc.NewSubConn(
		gb.addrs,
		balancer.NewSubConnOptions{HealthCheckEnabled: healthCheckEnabled},
	)
	if err != nil {
		gb.log.Errorf("failed to NewSubConn: %v", err)
		return
	}
	gb.scRefs[sc] = &subConnRef{
		subConn:     sc,
		stateSignal: make(chan struct{}),
		lastResp:    time.Now(),
	}
	gb.scStates[sc] = connectivity.Idle
	gb.scRefList = append(gb.scRefList, gb.scRefs[sc])
	sc.Connect()
}

// getReadySubConnRef returns a subConnRef and a bool. The bool indicates whether
// the boundKey exists in the affinityMap. If returned subConnRef is a nil, it
// means the underlying subconn is not READY yet.
func (gb *gcpBalancer) getReadySubConnRef(boundKey string) (*subConnRef, bool) {
	gb.mu.Lock()
	defer gb.mu.Unlock()

	if sc, ok := gb.affinityMap[boundKey]; ok {
		if gb.scStates[sc] != connectivity.Ready {
			// It's possible that the bound subconn is not in the readySubConns list,
			// If it's not ready, we throw ErrNoSubConnAvailable or
			// fallback to a previously mapped ready subconn or the least busy.
			if gb.cfg.GetChannelPool().GetFallbackToReady() {
				if sc, ok := gb.fallbackMap[boundKey]; ok {
					return gb.scRefs[sc], true
				}
				// Try to create fallback mapping.
				if scRef, err := gb.picker.(*gcpPicker).getLeastBusySubConnRef(); err == nil {
					gb.fallbackMap[boundKey] = scRef.subConn
					return scRef, true
				}
			}
			return nil, true
		}
		return gb.scRefs[sc], true
	}
	return nil, false
}

func (gb *gcpBalancer) getSubConnRoundRobin(ctx context.Context) *subConnRef {
	if len(gb.scRefList) == 0 {
		gb.newSubConn()
	}
	scRef := gb.scRefList[atomic.AddUint32(&gb.rrRefId, 1)%uint32(len(gb.scRefList))]

	gb.mu.RLock()
	if state := gb.scStates[scRef.subConn]; state == connectivity.Ready {
		gb.mu.RUnlock()
		return scRef
	} else {
		grpclog.Infof("grpcgcp.gcpBalancer: scRef is not ready: %v", state)
	}

	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()

	// Wait until SubConn is ready or call context is done.
	for gb.scStates[scRef.subConn] != connectivity.Ready {
		sigChan := scRef.stateSignal
		gb.mu.RUnlock()
		select {
		case <-ctx.Done():
			return scRef
		case <-ticker.C:
		case <-sigChan:
		}
		gb.mu.RLock()
	}
	gb.mu.RUnlock()

	return scRef
}

// bindSubConn binds the given affinity key to an existing subConnRef.
func (gb *gcpBalancer) bindSubConn(bindKey string, sc balancer.SubConn) {
	gb.mu.Lock()
	defer gb.mu.Unlock()
	_, ok := gb.affinityMap[bindKey]
	if !ok {
		gb.affinityMap[bindKey] = sc
	}
	gb.scRefs[sc].affinityIncr()
}

// unbindSubConn removes the existing binding associated with the key.
func (gb *gcpBalancer) unbindSubConn(boundKey string) {
	gb.mu.Lock()
	defer gb.mu.Unlock()
	boundSC, ok := gb.affinityMap[boundKey]
	if ok {
		gb.scRefs[boundSC].affinityDecr()
		delete(gb.affinityMap, boundKey)
	}
}

// regeneratePicker takes a snapshot of the balancer, and generates a picker
// from it. The picker is
//   - errPicker with ErrTransientFailure if the balancer is in TransientFailure,
//   - built by the pickerBuilder with all READY SubConns otherwise.
func (gb *gcpBalancer) regeneratePicker() {
	if gb.state == connectivity.TransientFailure {
		gb.picker = newErrPicker(balancer.ErrTransientFailure)
		return
	}
	readyRefs := []*subConnRef{}

	// Select ready subConns from subConn map.
	for sc, scState := range gb.scStates {
		if scState == connectivity.Ready {
			readyRefs = append(readyRefs, gb.scRefs[sc])
		}
	}
	gb.picker = newGCPPicker(readyRefs, gb)
}

func (gb *gcpBalancer) UpdateSubConnState(sc balancer.SubConn, scs balancer.SubConnState) {
	gb.mu.Lock()
	defer gb.mu.Unlock()
	s := scs.ConnectivityState

	if scRef, found := gb.refreshingScRefs[sc]; found {
		if gb.log.V(FINE) {
			gb.log.Infof("handle replacement SubConn state change: %p, %v", sc, s)
		}
		if s != connectivity.Ready {
			// Ignore the replacement sc until it's ready.
			return
		}

		// Replace SubConn of the scRef with the fresh SubConn (sc) concluding
		// the refresh process initiated by refresh(*subConnRef).
		oldSc := scRef.subConn
		gb.scStates[sc] = gb.scStates[oldSc]
		delete(gb.refreshingScRefs, sc)
		delete(gb.scRefs, oldSc)
		delete(gb.scStates, oldSc)
		gb.scRefs[sc] = scRef
		scRef.subConn = sc
		scRef.deCalls = 0
		scRef.lastResp = time.Now()
		scRef.refreshing = false
		scRef.refreshCnt++
		gb.cc.RemoveSubConn(oldSc)
	}

	if gb.log.V(FINE) {
		gb.log.Infof("handle SubConn state change: %p, %v", sc, s)
	}

	oldS, ok := gb.scStates[sc]
	if !ok {
		if gb.log.V(FINE) {
			gb.log.Infof(
				"got state changes for an unknown/replaced SubConn: %p, %v",
				sc,
				s,
			)
		}
		return
	}
	gb.scStates[sc] = s
	switch s {
	case connectivity.Idle:
		sc.Connect()
	case connectivity.Shutdown:
		delete(gb.scRefs, sc)
		delete(gb.scStates, sc)
	}
	if oldS == connectivity.Ready && s != oldS {
		// Subconn is broken. Remove fallback mapping to this subconn.
		for k, v := range gb.fallbackMap {
			if v == sc {
				delete(gb.fallbackMap, k)
			}
		}
	}
	if oldS != connectivity.Ready && s == connectivity.Ready {
		// Remove fallback mapping for the keys of recovered subconn.
		for k := range gb.fallbackMap {
			if gb.affinityMap[k] == sc {
				delete(gb.fallbackMap, k)
			}
		}
	}

	oldAggrState := gb.state
	gb.state = gb.csEvltr.recordTransition(oldS, s)

	// Regenerate picker when one of the following happens:
	//  - this sc became ready from not-ready
	//  - this sc became not-ready from ready
	//  - the aggregated state of balancer became TransientFailure from non-TransientFailure
	//  - the aggregated state of balancer became non-TransientFailure from TransientFailure
	if (s == connectivity.Ready) != (oldS == connectivity.Ready) ||
		(gb.state == connectivity.TransientFailure) != (oldAggrState == connectivity.TransientFailure) {
		gb.regeneratePicker()
		gb.cc.UpdateState(balancer.State{
			ConnectivityState: gb.state,
			Picker:            gb.picker,
		})
	}

	if scRef := gb.scRefs[sc]; scRef != nil {
		// Inform of the state change.
		close(scRef.stateSignal)
		scRef.stateSignal = make(chan struct{})
	}
}

// refresh initiates a new SubConn for a specific subConnRef and starts connecting.
// If the refresh is already initiated for the ref, then this is a no-op.
func (gb *gcpBalancer) refresh(ref *subConnRef) {
	if ref.refreshing {
		return
	}
	gb.mu.Lock()
	defer gb.mu.Unlock()
	if ref.refreshing {
		return
	}
	ref.refreshing = true
	sc, err := gb.cc.NewSubConn(
		gb.addrs,
		balancer.NewSubConnOptions{HealthCheckEnabled: healthCheckEnabled},
	)
	if err != nil {
		gb.log.Errorf("failed to create a replacement SubConn with NewSubConn: %v", err)
		return
	}
	gb.refreshingScRefs[sc] = ref
	sc.Connect()
}

func (gb *gcpBalancer) Close() {
}
