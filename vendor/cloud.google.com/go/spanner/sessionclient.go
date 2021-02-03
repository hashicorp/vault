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
	"cloud.google.com/go/internal/version"
	vkit "cloud.google.com/go/spanner/apiv1"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
	gtransport "google.golang.org/api/transport/grpc"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
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
	sessionReady(s *session)

	// sessionCreationFailed is called when the creation of a sub-batch of
	// sessions failed. The numSessions argument specifies the number of
	// sessions that could not be created as a result of this error. A
	// consumer may receive multiple errors per batch.
	sessionCreationFailed(err error, numSessions int32)
}

// sessionClient creates sessions for a database, either in batches or one at a
// time. Each session will be affiliated with a gRPC channel. sessionClient
// will ensure that the sessions that are created are evenly distributed over
// all available channels.
type sessionClient struct {
	mu     sync.Mutex
	closed bool

	connPool      gtransport.ConnPool
	database      string
	id            string
	sessionLabels map[string]string
	md            metadata.MD
	batchTimeout  time.Duration
	logger        *log.Logger
	callOptions   *vkit.CallOptions
}

// newSessionClient creates a session client to use for a database.
func newSessionClient(connPool gtransport.ConnPool, database string, sessionLabels map[string]string, md metadata.MD, logger *log.Logger, callOptions *vkit.CallOptions) *sessionClient {
	return &sessionClient{
		connPool:      connPool,
		database:      database,
		id:            cidGen.nextID(database),
		sessionLabels: sessionLabels,
		md:            md,
		batchTimeout:  time.Minute,
		logger:        logger,
		callOptions:   callOptions,
	}
}

func (sc *sessionClient) close() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.closed = true
	return sc.connPool.Close()
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
	ctx = contextWithOutgoingMetadata(ctx, sc.md)
	sid, err := client.CreateSession(ctx, &sppb.CreateSessionRequest{
		Database: sc.database,
		Session:  &sppb.Session{Labels: sc.sessionLabels},
	})
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
			go sc.executeBatchCreateSessions(client, createCountForChannel, sc.sessionLabels, sc.md, consumer)
			numBeingCreated += createCountForChannel
		}
	}
	return nil
}

// executeBatchCreateSessions executes the gRPC call for creating a batch of
// sessions.
func (sc *sessionClient) executeBatchCreateSessions(client *vkit.Client, createCount int32, labels map[string]string, md metadata.MD, consumer sessionConsumer) {
	ctx, cancel := context.WithTimeout(context.Background(), sc.batchTimeout)
	defer cancel()
	ctx = contextWithOutgoingMetadata(ctx, sc.md)

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
			consumer.sessionCreationFailed(err, remainingCreateCount)
			break
		}
		if ctx.Err() != nil {
			trace.TracePrintf(ctx, nil, "Context error while creating a batch of %d sessions: %v", createCount, ctx.Err())
			consumer.sessionCreationFailed(ToSpannerError(ctx.Err()), remainingCreateCount)
			break
		}
		response, err := client.BatchCreateSessions(ctx, &sppb.BatchCreateSessionsRequest{
			SessionCount:    remainingCreateCount,
			Database:        sc.database,
			SessionTemplate: &sppb.Session{Labels: labels},
		})
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error creating a batch of %d sessions: %v", remainingCreateCount, err)
			consumer.sessionCreationFailed(ToSpannerError(err), remainingCreateCount)
			break
		}
		actuallyCreated := int32(len(response.Session))
		trace.TracePrintf(ctx, nil, "Received a batch of %d sessions", actuallyCreated)
		for _, s := range response.Session {
			consumer.sessionReady(&session{valid: true, client: client, id: s.Name, createTime: time.Now(), md: md, logger: sc.logger})
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
func (sc *sessionClient) nextClient() (*vkit.Client, error) {
	// This call should never return an error as we are passing in an existing
	// connection, so we can safely ignore it.
	client, err := vkit.NewClient(context.Background(), option.WithGRPCConn(sc.connPool.Conn()))
	if err != nil {
		return nil, err
	}
	client.SetGoogleClientInfo("gccl", version.Repo)
	if sc.callOptions != nil {
		client.CallOptions = mergeCallOptions(client.CallOptions, sc.callOptions)
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
