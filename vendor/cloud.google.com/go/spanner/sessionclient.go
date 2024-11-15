/*
Copyright 2019 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package spanner

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"cloud.google.com/go/internal/trace"
	vkit "cloud.google.com/go/spanner/apiv1"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"cloud.google.com/go/spanner/internal"
	"github.com/googleapis/gax-go/v2"
	"go.opencensus.io/tag"
	"google.golang.org/api/option"
	gtransport "google.golang.org/api/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var cidGen = newClientIDGenerator()

type clientIDGenerator struct {
	mu  sync.Mutex
	ids map[string]int
}

func newClientIDGenerator() *clientIDGenerator {
	return &clientIDGenerator{ids: make(map[string]int)}
}

func (cg *clientIDGenerator) nextID(database string) string {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	var id int
	if val, ok := cg.ids[database]; ok {
		id = val + 1
	} else {
		id = 1
	}
	cg.ids[database] = id
	return fmt.Sprintf("client-%d", id)
}

// sessionConsumer is passed to the batchCreateSessions method and will receive
// the sessions that are created as they become available. A sessionConsumer
// implementation must be safe for concurrent use.
//
// The interface is implemented by sessionPool and is used for testing the
// sessionClient.
type sessionConsumer interface {
	// sessionReady is called when a session has been created and is ready for
	// use.
	sessionReady(ctx context.Context, s *session)

	// sessionCreationFailed is called when the creation of a sub-batch of
	// sessions failed. The numSessions argument specifies the number of
	// sessions that could not be created as a result of this error. A
	// consumer may receive multiple errors per batch.
	sessionCreationFailed(ctx context.Context, err error, numSessions int32, isMultiplexed bool)
}

// sessionClient creates sessions for a database, either in batches or one at a
// time. Each session will be affiliated with a gRPC channel. sessionClient
// will ensure that the sessions that are created are evenly distributed over
// all available channels.
type sessionClient struct {
	waitWorkers          sync.WaitGroup
	mu                   sync.Mutex
	closed               bool
	disableRouteToLeader bool

	connPool             gtransport.ConnPool
	database             string
	id                   string
	userAgent            string
	sessionLabels        map[string]string
	databaseRole         string
	md                   metadata.MD
	batchTimeout         time.Duration
	logger               *log.Logger
	callOptions          *vkit.CallOptions
	otConfig             *openTelemetryConfig
	metricsTracerFactory *builtinMetricsTracerFactory
}

// newSessionClient creates a session client to use for a database.
func newSessionClient(connPool gtransport.ConnPool, database, userAgent string, sessionLabels map[string]string, databaseRole string, disableRouteToLeader bool, md metadata.MD, batchTimeout time.Duration, logger *log.Logger, callOptions *vkit.CallOptions) *sessionClient {
	return &sessionClient{
		connPool:             connPool,
		database:             database,
		userAgent:            userAgent,
		id:                   cidGen.nextID(database),
		sessionLabels:        sessionLabels,
		databaseRole:         databaseRole,
		disableRouteToLeader: disableRouteToLeader,
		md:                   md,
		batchTimeout:         batchTimeout,
		logger:               logger,
		callOptions:          callOptions,
	}
}

func (sc *sessionClient) close() error {
	defer sc.waitWorkers.Wait()

	var err error
	func() {
		sc.mu.Lock()
		defer sc.mu.Unlock()

		sc.closed = true
		err = sc.connPool.Close()
	}()
	return err
}

// createSession creates one session for the database of the sessionClient. The
// session is created using one synchronous RPC.
func (sc *sessionClient) createSession(ctx context.Context) (*session, error) {
	sc.mu.Lock()
	if sc.closed {
		sc.mu.Unlock()
		return nil, spannerErrorf(codes.FailedPrecondition, "SessionClient is closed")
	}
	sc.mu.Unlock()
	client, err := sc.nextClient()
	if err != nil {
		return nil, err
	}

	var md metadata.MD
	sid, err := client.CreateSession(contextWithOutgoingMetadata(ctx, sc.md, sc.disableRouteToLeader), &sppb.CreateSessionRequest{
		Database: sc.database,
		Session:  &sppb.Session{Labels: sc.sessionLabels, CreatorRole: sc.databaseRole},
	}, gax.WithGRPCOptions(grpc.Header(&md)))

	if getGFELatencyMetricsFlag() && md != nil {
		_, instance, database, err := parseDatabaseName(sc.database)
		if err != nil {
			return nil, ToSpannerError(err)
		}
		ctxGFE, err := tag.New(ctx,
			tag.Upsert(tagKeyClientID, sc.id),
			tag.Upsert(tagKeyDatabase, database),
			tag.Upsert(tagKeyInstance, instance),
			tag.Upsert(tagKeyLibVersion, internal.Version),
		)
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", ToSpannerError(err))
		}
		err = captureGFELatencyStats(ctxGFE, md, "createSession")
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", ToSpannerError(err))
		}
	}
	if metricErr := recordGFELatencyMetricsOT(ctx, md, "createSession", sc.otConfig); metricErr != nil {
		trace.TracePrintf(ctx, nil, "Error in recording GFE Latency through OpenTelemetry. Error: %v", metricErr)
	}
	if err != nil {
		return nil, ToSpannerError(err)
	}
	return &session{valid: true, client: client, id: sid.Name, createTime: time.Now(), md: sc.md, logger: sc.logger}, nil
}

// batchCreateSessions creates a batch of sessions for the database of the
// sessionClient and returns these to the given sessionConsumer.
//
// createSessionCount is the number of sessions that should be created. The
// sessionConsumer is guaranteed to receive the requested number of sessions if
// no error occurs. If one or more errors occur, the sessionConsumer will
// receive any number of sessions + any number of errors, where each error will
// include the number of sessions that could not be created as a result of the
// error. The sum of returned sessions and errored sessions will be equal to
// the number of requested sessions.
// If distributeOverChannels is true, the sessions will be equally distributed
// over all the channels that are in use by the client.
func (sc *sessionClient) batchCreateSessions(createSessionCount int32, distributeOverChannels bool, consumer sessionConsumer) error {
	var sessionCountPerChannel int32
	var remainder int32
	if distributeOverChannels {
		// The sessions that we create should be evenly distributed over all the
		// channels (gapic clients) that are used by the client. Each gapic client
		// will do a request for a fraction of the total.
		sessionCountPerChannel = createSessionCount / int32(sc.connPool.Num())
		// The remainder of the calculation will be added to the number of sessions
		// that will be created for the first channel, to ensure that we create the
		// exact number of requested sessions.
		remainder = createSessionCount % int32(sc.connPool.Num())
	} else {
		sessionCountPerChannel = createSessionCount
	}
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed {
		return spannerErrorf(codes.FailedPrecondition, "SessionClient is closed")
	}
	// Spread the session creation over all available gRPC channels. Spanner
	// will maintain server side caches for a session on the gRPC channel that
	// is used by the session. A session should therefore always use the same
	// channel, and the sessions should be as evenly distributed as possible
	// over the channels.
	var numBeingCreated int32
	for i := 0; i < sc.connPool.Num() && numBeingCreated < createSessionCount; i++ {
		client, err := sc.nextClient()
		if err != nil {
			return err
		}
		// Determine the number of sessions that should be created for this
		// channel. The createCount for the first channel will be increased
		// with the remainder of the division of the total number of sessions
		// with the number of channels. All other channels will just use the
		// result of the division over all channels.
		createCountForChannel := sessionCountPerChannel
		if i == 0 {
			// We add the remainder to the first gRPC channel we use. We could
			// also spread the remainder over all channels, but this ensures
			// that small batches of sessions (i.e. less than numChannels) are
			// created in one RPC.
			createCountForChannel += remainder
		}
		if createCountForChannel > 0 {
			sc.waitWorkers.Add(1)
			go sc.executeBatchCreateSessions(client, createCountForChannel, sc.sessionLabels, sc.md, consumer)
			numBeingCreated += createCountForChannel
		}
	}
	return nil
}

// executeBatchCreateSessions executes the gRPC call for creating a batch of
// sessions.
func (sc *sessionClient) executeBatchCreateSessions(client spannerClient, createCount int32, labels map[string]string, md metadata.MD, consumer sessionConsumer) {
	defer sc.waitWorkers.Done()
	ctx, cancel := context.WithTimeout(context.Background(), sc.batchTimeout)
	defer cancel()
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.BatchCreateSessions")
	defer func() { trace.EndSpan(ctx, nil) }()
	trace.TracePrintf(ctx, nil, "Creating a batch of %d sessions", createCount)

	remainingCreateCount := createCount
	for {
		sc.mu.Lock()
		closed := sc.closed
		sc.mu.Unlock()
		if closed {
			err := spannerErrorf(codes.Canceled, "Session client closed")
			trace.TracePrintf(ctx, nil, "Session client closed while creating a batch of %d sessions: %v", createCount, err)
			consumer.sessionCreationFailed(ctx, err, remainingCreateCount, false)
			break
		}
		if ctx.Err() != nil {
			trace.TracePrintf(ctx, nil, "Context error while creating a batch of %d sessions: %v", createCount, ctx.Err())
			consumer.sessionCreationFailed(ctx, ToSpannerError(ctx.Err()), remainingCreateCount, false)
			break
		}
		var mdForGFELatency metadata.MD
		response, err := client.BatchCreateSessions(contextWithOutgoingMetadata(ctx, sc.md, sc.disableRouteToLeader), &sppb.BatchCreateSessionsRequest{
			SessionCount:    remainingCreateCount,
			Database:        sc.database,
			SessionTemplate: &sppb.Session{Labels: labels, CreatorRole: sc.databaseRole},
		}, gax.WithGRPCOptions(grpc.Header(&mdForGFELatency)))

		if getGFELatencyMetricsFlag() && mdForGFELatency != nil {
			_, instance, database, err := parseDatabaseName(sc.database)
			if err != nil {
				trace.TracePrintf(ctx, nil, "Error getting instance and database name: %v", err)
			}
			// Errors should not prevent initializing the session pool.
			ctxGFE, err := tag.New(ctx,
				tag.Upsert(tagKeyClientID, sc.id),
				tag.Upsert(tagKeyDatabase, database),
				tag.Upsert(tagKeyInstance, instance),
				tag.Upsert(tagKeyLibVersion, internal.Version),
			)
			if err != nil {
				trace.TracePrintf(ctx, nil, "Error in adding tags in BatchCreateSessions for GFE Latency: %v", err)
			}
			err = captureGFELatencyStats(ctxGFE, mdForGFELatency, "executeBatchCreateSessions")
			if err != nil {
				trace.TracePrintf(ctx, nil, "Error in Capturing GFE Latency and Header Missing count. Try disabling and rerunning. Error: %v", err)
			}
		}
		if metricErr := recordGFELatencyMetricsOT(ctx, mdForGFELatency, "executeBatchCreateSessions", sc.otConfig); metricErr != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency through OpenTelemetry. Error: %v", metricErr)
		}
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error creating a batch of %d sessions: %v", remainingCreateCount, err)
			consumer.sessionCreationFailed(ctx, ToSpannerError(err), remainingCreateCount, false)
			break
		}
		actuallyCreated := int32(len(response.Session))
		trace.TracePrintf(ctx, nil, "Received a batch of %d sessions", actuallyCreated)
		for _, s := range response.Session {
			consumer.sessionReady(ctx, &session{valid: true, client: client, id: s.Name, createTime: time.Now(), md: md, logger: sc.logger})
		}
		if actuallyCreated < remainingCreateCount {
			// Spanner could return less sessions than requested. In that case, we
			// should do another call using the same gRPC channel.
			remainingCreateCount -= actuallyCreated
		} else {
			trace.TracePrintf(ctx, nil, "Finished creating %d sessions", createCount)
			break
		}
	}
}

func (sc *sessionClient) executeCreateMultiplexedSession(ctx context.Context, client spannerClient, md metadata.MD, consumer sessionConsumer) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.CreateSession")
	defer func() { trace.EndSpan(ctx, nil) }()
	trace.TracePrintf(ctx, nil, "Creating a multiplexed session")
	sc.mu.Lock()
	closed := sc.closed
	sc.mu.Unlock()
	if closed {
		err := spannerErrorf(codes.Canceled, "Session client closed")
		trace.TracePrintf(ctx, nil, "Session client closed while creating a multiplexed session: %v", err)
		return
	}
	if ctx.Err() != nil {
		trace.TracePrintf(ctx, nil, "Context error while creating a multiplexed session: %v", ctx.Err())
		consumer.sessionCreationFailed(ctx, ToSpannerError(ctx.Err()), 1, true)
		return
	}
	var mdForGFELatency metadata.MD
	response, err := client.CreateSession(contextWithOutgoingMetadata(ctx, sc.md, sc.disableRouteToLeader), &sppb.CreateSessionRequest{
		Database: sc.database,
		// Multiplexed sessions do not support labels.
		Session: &sppb.Session{CreatorRole: sc.databaseRole, Multiplexed: true},
	}, gax.WithGRPCOptions(grpc.Header(&mdForGFELatency)))

	if getGFELatencyMetricsFlag() && mdForGFELatency != nil {
		_, instance, database, err := parseDatabaseName(sc.database)
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error getting instance and database name: %v", err)
		}
		// Errors should not prevent initializing the session pool.
		ctxGFE, err := tag.New(ctx,
			tag.Upsert(tagKeyClientID, sc.id),
			tag.Upsert(tagKeyDatabase, database),
			tag.Upsert(tagKeyInstance, instance),
			tag.Upsert(tagKeyLibVersion, internal.Version),
		)
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error in adding tags in CreateSession for GFE Latency: %v", err)
		}
		err = captureGFELatencyStats(ctxGFE, mdForGFELatency, "executeCreateSession")
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error in Capturing GFE Latency and Header Missing count. Try disabling and rerunning. Error: %v", err)
		}
	}
	if metricErr := recordGFELatencyMetricsOT(ctx, mdForGFELatency, "executeCreateSession", sc.otConfig); metricErr != nil {
		trace.TracePrintf(ctx, nil, "Error in recording GFE Latency through OpenTelemetry. Error: %v", metricErr)
	}
	if err != nil {
		trace.TracePrintf(ctx, nil, "Error creating a multiplexed sessions: %v", err)
		consumer.sessionCreationFailed(ctx, ToSpannerError(err), 1, true)
		return
	}
	consumer.sessionReady(ctx, &session{valid: true, client: client, id: response.Name, createTime: time.Now(), md: md, logger: sc.logger, isMultiplexed: response.Multiplexed})
	trace.TracePrintf(ctx, nil, "Finished creating multiplexed sessions")
}

func (sc *sessionClient) sessionWithID(id string) (*session, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	client, err := sc.nextClient()
	if err != nil {
		return nil, err
	}
	return &session{valid: true, client: client, id: id, createTime: time.Now(), md: sc.md, logger: sc.logger}, nil
}

// nextClient returns the next gRPC client to use for session creation. The
// client is set on the session, and used by all subsequent gRPC calls on the
// session. Using the same channel for all gRPC calls for a session ensures the
// optimal usage of server side caches.
func (sc *sessionClient) nextClient() (spannerClient, error) {
	var clientOpt option.ClientOption
	if _, ok := sc.connPool.(*gmeWrapper); ok {
		// Pass GCPMultiEndpoint as a pool.
		clientOpt = gtransport.WithConnPool(sc.connPool)
	} else {
		// Pick a grpc.ClientConn from a regular pool.
		clientOpt = option.WithGRPCConn(sc.connPool.Conn())
	}
	client, err := newGRPCSpannerClient(context.Background(), sc, clientOpt)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// mergeCallOptions merges two CallOptions into one and the first argument has
// a lower order of precedence than the second one.
func mergeCallOptions(a *vkit.CallOptions, b *vkit.CallOptions) *vkit.CallOptions {
	res := &vkit.CallOptions{}
	resVal := reflect.ValueOf(res).Elem()
	aVal := reflect.ValueOf(a).Elem()
	bVal := reflect.ValueOf(b).Elem()

	t := aVal.Type()

	for i := 0; i < aVal.NumField(); i++ {
		fieldName := t.Field(i).Name

		aFieldVal := aVal.Field(i).Interface().([]gax.CallOption)
		bFieldVal := bVal.Field(i).Interface().([]gax.CallOption)

		merged := append(aFieldVal, bFieldVal...)
		resVal.FieldByName(fieldName).Set(reflect.ValueOf(merged))
	}
	return res
}
