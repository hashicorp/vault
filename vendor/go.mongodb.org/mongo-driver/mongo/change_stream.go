// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/internal/csot"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

var (
	// ErrMissingResumeToken indicates that a change stream notification from the server did not contain a resume token.
	ErrMissingResumeToken = errors.New("cannot provide resume functionality when the resume token is missing")
	// ErrNilCursor indicates that the underlying cursor for the change stream is nil.
	ErrNilCursor = errors.New("cursor is nil")

	minResumableLabelWireVersion int32 = 9 // Wire version at which the server includes the resumable error label
	networkErrorLabel                  = "NetworkError"
	resumableErrorLabel                = "ResumableChangeStreamError"
	errorCursorNotFound          int32 = 43 // CursorNotFound error code

	// Allowlist of error codes that are considered resumable.
	resumableChangeStreamErrors = map[int32]struct{}{
		6:     {}, // HostUnreachable
		7:     {}, // HostNotFound
		89:    {}, // NetworkTimeout
		91:    {}, // ShutdownInProgress
		189:   {}, // PrimarySteppedDown
		262:   {}, // ExceededTimeLimit
		9001:  {}, // SocketException
		10107: {}, // NotPrimary
		11600: {}, // InterruptedAtShutdown
		11602: {}, // InterruptedDueToReplStateChange
		13435: {}, // NotPrimaryNoSecondaryOK
		13436: {}, // NotPrimaryOrSecondary
		63:    {}, // StaleShardVersion
		150:   {}, // StaleEpoch
		13388: {}, // StaleConfig
		234:   {}, // RetryChangeStream
		133:   {}, // FailedToSatisfyReadPreference
	}
)

// ChangeStream is used to iterate over a stream of events. Each event can be decoded into a Go type via the Decode
// method or accessed as raw BSON via the Current field. This type is not goroutine safe and must not be used
// concurrently by multiple goroutines. For more information about change streams, see
// https://www.mongodb.com/docs/manual/changeStreams/.
type ChangeStream struct {
	// Current is the BSON bytes of the current event. This property is only valid until the next call to Next or
	// TryNext. If continued access is required, a copy must be made.
	Current bson.Raw

	aggregate       *operation.Aggregate
	pipelineSlice   []bsoncore.Document
	pipelineOptions map[string]bsoncore.Value
	cursor          changeStreamCursor
	cursorOptions   driver.CursorOptions
	batch           []bsoncore.Document
	resumeToken     bson.Raw
	err             error
	sess            *session.Client
	client          *Client
	bsonOpts        *options.BSONOptions
	registry        *bsoncodec.Registry
	streamType      StreamType
	options         *options.ChangeStreamOptions
	selector        description.ServerSelector
	operationTime   *primitive.Timestamp
	wireVersion     *description.VersionRange
}

type changeStreamConfig struct {
	readConcern    *readconcern.ReadConcern
	readPreference *readpref.ReadPref
	client         *Client
	bsonOpts       *options.BSONOptions
	registry       *bsoncodec.Registry
	streamType     StreamType
	collectionName string
	databaseName   string
	crypt          driver.Crypt
}

func newChangeStream(ctx context.Context, config changeStreamConfig, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	cursorOpts := config.client.createBaseCursorOptions()

	cursorOpts.MarshalValueEncoderFn = newEncoderFn(config.bsonOpts, config.registry)

	cs := &ChangeStream{
		client:     config.client,
		bsonOpts:   config.bsonOpts,
		registry:   config.registry,
		streamType: config.streamType,
		options:    options.MergeChangeStreamOptions(opts...),
		selector: description.CompositeSelector([]description.ServerSelector{
			description.ReadPrefSelector(config.readPreference),
			description.LatencySelector(config.client.localThreshold),
		}),
		cursorOptions: cursorOpts,
	}

	cs.sess = sessionFromContext(ctx)
	if cs.sess == nil && cs.client.sessionPool != nil {
		cs.sess = session.NewImplicitClientSession(cs.client.sessionPool, cs.client.id)
	}
	if cs.err = cs.client.validSession(cs.sess); cs.err != nil {
		closeImplicitSession(cs.sess)
		return nil, cs.Err()
	}

	cs.aggregate = operation.NewAggregate(nil).
		ReadPreference(config.readPreference).ReadConcern(config.readConcern).
		Deployment(cs.client.deployment).ClusterClock(cs.client.clock).
		CommandMonitor(cs.client.monitor).Session(cs.sess).ServerSelector(cs.selector).Retry(driver.RetryNone).
		ServerAPI(cs.client.serverAPI).Crypt(config.crypt).Timeout(cs.client.timeout).
		Authenticator(cs.client.authenticator)

	if cs.options.Collation != nil {
		cs.aggregate.Collation(bsoncore.Document(cs.options.Collation.ToDocument()))
	}
	if comment := cs.options.Comment; comment != nil {
		cs.aggregate.Comment(*comment)

		commentVal, err := marshalValue(comment, cs.bsonOpts, cs.registry)
		if err != nil {
			return nil, err
		}
		cs.cursorOptions.Comment = commentVal
	}
	if cs.options.BatchSize != nil {
		cs.aggregate.BatchSize(*cs.options.BatchSize)
		cs.cursorOptions.BatchSize = *cs.options.BatchSize
	}
	if cs.options.MaxAwaitTime != nil {
		cs.cursorOptions.MaxTimeMS = int64(*cs.options.MaxAwaitTime / time.Millisecond)
	}
	if cs.options.Custom != nil {
		// Marshal all custom options before passing to the initial aggregate. Return
		// any errors from Marshaling.
		customOptions := make(map[string]bsoncore.Value)
		for optionName, optionValue := range cs.options.Custom {
			bsonType, bsonData, err := bson.MarshalValueWithRegistry(cs.registry, optionValue)
			if err != nil {
				cs.err = err
				closeImplicitSession(cs.sess)
				return nil, cs.Err()
			}
			optionValueBSON := bsoncore.Value{Type: bsonType, Data: bsonData}
			customOptions[optionName] = optionValueBSON
		}
		cs.aggregate.CustomOptions(customOptions)
	}
	if cs.options.CustomPipeline != nil {
		// Marshal all custom pipeline options before building pipeline slice. Return
		// any errors from Marshaling.
		cs.pipelineOptions = make(map[string]bsoncore.Value)
		for optionName, optionValue := range cs.options.CustomPipeline {
			bsonType, bsonData, err := bson.MarshalValueWithRegistry(cs.registry, optionValue)
			if err != nil {
				cs.err = err
				closeImplicitSession(cs.sess)
				return nil, cs.Err()
			}
			optionValueBSON := bsoncore.Value{Type: bsonType, Data: bsonData}
			cs.pipelineOptions[optionName] = optionValueBSON
		}
	}

	switch cs.streamType {
	case ClientStream:
		cs.aggregate.Database("admin")
	case DatabaseStream:
		cs.aggregate.Database(config.databaseName)
	case CollectionStream:
		cs.aggregate.Collection(config.collectionName).Database(config.databaseName)
	default:
		closeImplicitSession(cs.sess)
		return nil, fmt.Errorf("must supply a valid StreamType in config, instead of %v", cs.streamType)
	}

	// When starting a change stream, cache startAfter as the first resume token if it is set. If not, cache
	// resumeAfter. If neither is set, do not cache a resume token.
	resumeToken := cs.options.StartAfter
	if resumeToken == nil {
		resumeToken = cs.options.ResumeAfter
	}
	var marshaledToken bson.Raw
	if resumeToken != nil {
		if marshaledToken, cs.err = bson.Marshal(resumeToken); cs.err != nil {
			closeImplicitSession(cs.sess)
			return nil, cs.Err()
		}
	}
	cs.resumeToken = marshaledToken

	if cs.err = cs.buildPipelineSlice(pipeline); cs.err != nil {
		closeImplicitSession(cs.sess)
		return nil, cs.Err()
	}
	var pipelineArr bsoncore.Document
	pipelineArr, cs.err = cs.pipelineToBSON()
	cs.aggregate.Pipeline(pipelineArr)

	if cs.err = cs.executeOperation(ctx, false); cs.err != nil {
		closeImplicitSession(cs.sess)
		return nil, cs.Err()
	}

	return cs, cs.Err()
}

func (cs *ChangeStream) createOperationDeployment(server driver.Server, connection driver.Connection) driver.Deployment {
	return &changeStreamDeployment{
		topologyKind: cs.client.deployment.Kind(),
		server:       server,
		conn:         connection,
	}
}

func (cs *ChangeStream) executeOperation(ctx context.Context, resuming bool) error {
	var server driver.Server
	var conn driver.Connection

	if server, cs.err = cs.client.deployment.SelectServer(ctx, cs.selector); cs.err != nil {
		return cs.Err()
	}
	if conn, cs.err = server.Connection(ctx); cs.err != nil {
		return cs.Err()
	}
	defer conn.Close()
	cs.wireVersion = conn.Description().WireVersion

	cs.aggregate.Deployment(cs.createOperationDeployment(server, conn))

	if resuming {
		cs.replaceOptions(cs.wireVersion)

		csOptDoc, err := cs.createPipelineOptionsDoc()
		if err != nil {
			return err
		}
		pipIdx, pipDoc := bsoncore.AppendDocumentStart(nil)
		pipDoc = bsoncore.AppendDocumentElement(pipDoc, "$changeStream", csOptDoc)
		if pipDoc, cs.err = bsoncore.AppendDocumentEnd(pipDoc, pipIdx); cs.err != nil {
			return cs.Err()
		}
		cs.pipelineSlice[0] = pipDoc

		var plArr bsoncore.Document
		if plArr, cs.err = cs.pipelineToBSON(); cs.err != nil {
			return cs.Err()
		}
		cs.aggregate.Pipeline(plArr)
	}

	// If cs.client.timeout is set and context is not already a Timeout context,
	// honor cs.client.timeout in new Timeout context for change stream
	// operation execution and potential retry.
	if cs.client.timeout != nil && !csot.IsTimeoutContext(ctx) {
		newCtx, cancelFunc := csot.MakeTimeoutContext(ctx, *cs.client.timeout)
		// Redefine ctx to be the new timeout-derived context.
		ctx = newCtx
		// Cancel the timeout-derived context at the end of executeOperation to avoid a context leak.
		defer cancelFunc()
	}

	// Execute the aggregate, retrying on retryable errors once (1) if retryable reads are enabled and
	// infinitely (-1) if context is a Timeout context.
	var retries int
	if cs.client.retryReads {
		retries = 1
	}
	if csot.IsTimeoutContext(ctx) {
		retries = -1
	}

	var err error
AggregateExecuteLoop:
	for {
		err = cs.aggregate.Execute(ctx)
		// If no error or no retries remain, do not retry.
		if err == nil || retries == 0 {
			break AggregateExecuteLoop
		}

		switch tt := err.(type) {
		case driver.Error:
			// If error is not retryable, do not retry.
			if !tt.RetryableRead() {
				break AggregateExecuteLoop
			}

			// If error is retryable: subtract 1 from retries, redo server selection, checkout
			// a connection, and restart loop.
			retries--
			server, err = cs.client.deployment.SelectServer(ctx, cs.selector)
			if err != nil {
				break AggregateExecuteLoop
			}

			conn.Close()
			conn, err = server.Connection(ctx)
			if err != nil {
				break AggregateExecuteLoop
			}
			defer conn.Close()

			// Update the wire version with data from the new connection.
			cs.wireVersion = conn.Description().WireVersion

			// Reset deployment.
			cs.aggregate.Deployment(cs.createOperationDeployment(server, conn))
		default:
			// Do not retry if error is not a driver error.
			break AggregateExecuteLoop
		}
	}
	if err != nil {
		cs.err = replaceErrors(err)
		return cs.err
	}

	cr := cs.aggregate.ResultCursorResponse()
	cr.Server = server

	cs.cursor, cs.err = driver.NewBatchCursor(cr, cs.sess, cs.client.clock, cs.cursorOptions)
	if cs.err = replaceErrors(cs.err); cs.err != nil {
		return cs.Err()
	}

	cs.updatePbrtFromCommand()
	if cs.options.StartAtOperationTime == nil && cs.options.ResumeAfter == nil &&
		cs.options.StartAfter == nil && cs.wireVersion.Max >= 7 &&
		cs.emptyBatch() && cs.resumeToken == nil {
		cs.operationTime = cs.sess.OperationTime
	}

	return cs.Err()
}

// Updates the post batch resume token after a successful aggregate or getMore operation.
func (cs *ChangeStream) updatePbrtFromCommand() {
	// Only cache the pbrt if an empty batch was returned and a pbrt was included
	if pbrt := cs.cursor.PostBatchResumeToken(); cs.emptyBatch() && pbrt != nil {
		cs.resumeToken = bson.Raw(pbrt)
	}
}

func (cs *ChangeStream) storeResumeToken() error {
	// If cs.Current is the last document in the batch and a pbrt is included, cache the pbrt
	// Otherwise, cache the _id of the document
	var tokenDoc bson.Raw
	if len(cs.batch) == 0 {
		if pbrt := cs.cursor.PostBatchResumeToken(); pbrt != nil {
			tokenDoc = bson.Raw(pbrt)
		}
	}

	if tokenDoc == nil {
		var ok bool
		tokenDoc, ok = cs.Current.Lookup("_id").DocumentOK()
		if !ok {
			_ = cs.Close(context.Background())
			return ErrMissingResumeToken
		}
	}

	cs.resumeToken = tokenDoc
	return nil
}

func (cs *ChangeStream) buildPipelineSlice(pipeline interface{}) error {
	val := reflect.ValueOf(pipeline)
	if !val.IsValid() || !(val.Kind() == reflect.Slice) {
		cs.err = errors.New("can only marshal slices and arrays into aggregation pipelines, but got invalid")
		return cs.err
	}

	cs.pipelineSlice = make([]bsoncore.Document, 0, val.Len()+1)

	csIdx, csDoc := bsoncore.AppendDocumentStart(nil)

	csDocTemp, err := cs.createPipelineOptionsDoc()
	if err != nil {
		return err
	}
	csDoc = bsoncore.AppendDocumentElement(csDoc, "$changeStream", csDocTemp)
	csDoc, cs.err = bsoncore.AppendDocumentEnd(csDoc, csIdx)
	if cs.err != nil {
		return cs.err
	}
	cs.pipelineSlice = append(cs.pipelineSlice, csDoc)

	for i := 0; i < val.Len(); i++ {
		var elem []byte
		elem, cs.err = marshal(val.Index(i).Interface(), cs.bsonOpts, cs.registry)
		if cs.err != nil {
			return cs.err
		}

		cs.pipelineSlice = append(cs.pipelineSlice, elem)
	}

	return cs.err
}

func (cs *ChangeStream) createPipelineOptionsDoc() (bsoncore.Document, error) {
	plDocIdx, plDoc := bsoncore.AppendDocumentStart(nil)

	if cs.streamType == ClientStream {
		plDoc = bsoncore.AppendBooleanElement(plDoc, "allChangesForCluster", true)
	}

	if cs.options.FullDocument != nil && *cs.options.FullDocument != options.Default {
		plDoc = bsoncore.AppendStringElement(plDoc, "fullDocument", string(*cs.options.FullDocument))
	}

	if cs.options.FullDocumentBeforeChange != nil {
		plDoc = bsoncore.AppendStringElement(plDoc, "fullDocumentBeforeChange", string(*cs.options.FullDocumentBeforeChange))
	}

	if cs.options.ResumeAfter != nil {
		var raDoc bsoncore.Document
		raDoc, cs.err = marshal(cs.options.ResumeAfter, cs.bsonOpts, cs.registry)
		if cs.err != nil {
			return nil, cs.err
		}

		plDoc = bsoncore.AppendDocumentElement(plDoc, "resumeAfter", raDoc)
	}

	if cs.options.ShowExpandedEvents != nil {
		plDoc = bsoncore.AppendBooleanElement(plDoc, "showExpandedEvents", *cs.options.ShowExpandedEvents)
	}

	if cs.options.StartAfter != nil {
		var saDoc bsoncore.Document
		saDoc, cs.err = marshal(cs.options.StartAfter, cs.bsonOpts, cs.registry)
		if cs.err != nil {
			return nil, cs.err
		}

		plDoc = bsoncore.AppendDocumentElement(plDoc, "startAfter", saDoc)
	}

	if cs.options.StartAtOperationTime != nil {
		plDoc = bsoncore.AppendTimestampElement(plDoc, "startAtOperationTime", cs.options.StartAtOperationTime.T, cs.options.StartAtOperationTime.I)
	}

	// Append custom pipeline options.
	for optionName, optionValue := range cs.pipelineOptions {
		plDoc = bsoncore.AppendValueElement(plDoc, optionName, optionValue)
	}

	if plDoc, cs.err = bsoncore.AppendDocumentEnd(plDoc, plDocIdx); cs.err != nil {
		return nil, cs.err
	}

	return plDoc, nil
}

func (cs *ChangeStream) pipelineToBSON() (bsoncore.Document, error) {
	pipelineDocIdx, pipelineArr := bsoncore.AppendArrayStart(nil)
	for i, doc := range cs.pipelineSlice {
		pipelineArr = bsoncore.AppendDocumentElement(pipelineArr, strconv.Itoa(i), doc)
	}
	if pipelineArr, cs.err = bsoncore.AppendArrayEnd(pipelineArr, pipelineDocIdx); cs.err != nil {
		return nil, cs.err
	}
	return pipelineArr, cs.err
}

func (cs *ChangeStream) replaceOptions(wireVersion *description.VersionRange) {
	// Cached resume token: use the resume token as the resumeAfter option and set no other resume options
	if cs.resumeToken != nil {
		cs.options.SetResumeAfter(cs.resumeToken)
		cs.options.SetStartAfter(nil)
		cs.options.SetStartAtOperationTime(nil)
		return
	}

	// No cached resume token but cached operation time: use the operation time as the startAtOperationTime option and
	// set no other resume options
	if (cs.sess.OperationTime != nil || cs.options.StartAtOperationTime != nil) && wireVersion.Max >= 7 {
		opTime := cs.options.StartAtOperationTime
		if cs.operationTime != nil {
			opTime = cs.sess.OperationTime
		}

		cs.options.SetStartAtOperationTime(opTime)
		cs.options.SetResumeAfter(nil)
		cs.options.SetStartAfter(nil)
		return
	}

	// No cached resume token or operation time: set none of the resume options
	cs.options.SetResumeAfter(nil)
	cs.options.SetStartAfter(nil)
	cs.options.SetStartAtOperationTime(nil)
}

// ID returns the ID for this change stream, or 0 if the cursor has been closed or exhausted.
func (cs *ChangeStream) ID() int64 {
	if cs.cursor == nil {
		return 0
	}
	return cs.cursor.ID()
}

// RemainingBatchLength returns the number of documents left in the current batch. If this returns zero, the subsequent
// call to Next or TryNext will do a network request to fetch the next batch.
func (cs *ChangeStream) RemainingBatchLength() int {
	return len(cs.batch)
}

// SetBatchSize sets the number of documents to fetch from the database with
// each iteration of the ChangeStream's "Next" or "TryNext" method. This setting
// only affects subsequent document batches fetched from the database.
func (cs *ChangeStream) SetBatchSize(size int32) {
	// Set batch size on the cursor options also so any "resumed" change stream
	// cursors will pick up the latest batch size setting.
	cs.cursorOptions.BatchSize = size
	cs.cursor.SetBatchSize(size)
}

// Decode will unmarshal the current event document into val and return any errors from the unmarshalling process
// without any modification. If val is nil or is a typed nil, an error will be returned.
func (cs *ChangeStream) Decode(val interface{}) error {
	if cs.cursor == nil {
		return ErrNilCursor
	}

	dec, err := getDecoder(cs.Current, cs.bsonOpts, cs.registry)
	if err != nil {
		return fmt.Errorf("error configuring BSON decoder: %w", err)
	}
	return dec.Decode(val)
}

// Err returns the last error seen by the change stream, or nil if no errors has occurred.
func (cs *ChangeStream) Err() error {
	if cs.err != nil {
		return replaceErrors(cs.err)
	}
	if cs.cursor == nil {
		return nil
	}

	return replaceErrors(cs.cursor.Err())
}

// Close closes this change stream and the underlying cursor. Next and TryNext must not be called after Close has been
// called. Close is idempotent. After the first call, any subsequent calls will not change the state.
func (cs *ChangeStream) Close(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	defer closeImplicitSession(cs.sess)

	if cs.cursor == nil {
		return nil // cursor is already closed
	}

	cs.err = replaceErrors(cs.cursor.Close(ctx))
	cs.cursor = nil
	return cs.Err()
}

// ResumeToken returns the last cached resume token for this change stream, or nil if a resume token has not been
// stored.
func (cs *ChangeStream) ResumeToken() bson.Raw {
	return cs.resumeToken
}

// Next gets the next event for this change stream. It returns true if there were no errors and the next event document
// is available.
//
// Next blocks until an event is available, an error occurs, or ctx expires. If ctx expires, the error
// will be set to ctx.Err(). In an error case, Next will return false.
//
// If Next returns false, subsequent calls will also return false.
func (cs *ChangeStream) Next(ctx context.Context) bool {
	return cs.next(ctx, false)
}

// TryNext attempts to get the next event for this change stream. It returns true if there were no errors and the next
// event document is available.
//
// TryNext returns false if the change stream is closed by the server, an error occurs when getting changes from the
// server, the next change is not yet available, or ctx expires. If ctx expires, the error will be set to ctx.Err().
//
// If TryNext returns false and an error occurred or the change stream was closed
// (i.e. cs.Err() != nil || cs.ID() == 0), subsequent attempts will also return false. Otherwise, it is safe to call
// TryNext again until a change is available.
//
// This method requires driver version >= 1.2.0.
func (cs *ChangeStream) TryNext(ctx context.Context) bool {
	return cs.next(ctx, true)
}

func (cs *ChangeStream) next(ctx context.Context, nonBlocking bool) bool {
	// return false right away if the change stream has already errored or if cursor is closed.
	if cs.err != nil {
		return false
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if len(cs.batch) == 0 {
		cs.loopNext(ctx, nonBlocking)
		if cs.err != nil {
			cs.err = replaceErrors(cs.err)
			return false
		}
		if len(cs.batch) == 0 {
			return false
		}
	}

	// successfully got non-empty batch
	cs.Current = bson.Raw(cs.batch[0])
	cs.batch = cs.batch[1:]
	if cs.err = cs.storeResumeToken(); cs.err != nil {
		return false
	}
	return true
}

func (cs *ChangeStream) loopNext(ctx context.Context, nonBlocking bool) {
	for {
		if cs.cursor == nil {
			return
		}

		if cs.cursor.Next(ctx) {
			// non-empty batch returned
			cs.batch, cs.err = cs.cursor.Batch().Documents()
			return
		}

		cs.err = replaceErrors(cs.cursor.Err())
		if cs.err == nil {
			// Check if cursor is alive
			if cs.ID() == 0 {
				return
			}

			// If a getMore was done but the batch was empty, the batch cursor will return false with no error.
			// Update the tracked resume token to catch the post batch resume token from the server response.
			cs.updatePbrtFromCommand()
			if nonBlocking {
				// stop after a successful getMore, even though the batch was empty
				return
			}
			continue // loop getMore until a non-empty batch is returned or an error occurs
		}

		if !cs.isResumableError() {
			return
		}

		// ignore error from cursor close because if the cursor is deleted or errors we tried to close it and will remake and try to get next batch
		_ = cs.cursor.Close(ctx)
		if cs.err = cs.executeOperation(ctx, true); cs.err != nil {
			return
		}
	}
}

func (cs *ChangeStream) isResumableError() bool {
	var commandErr CommandError
	if !errors.As(cs.err, &commandErr) || commandErr.HasErrorLabel(networkErrorLabel) {
		// All non-server errors or network errors are resumable.
		return true
	}

	if commandErr.Code == errorCursorNotFound {
		return true
	}

	// For wire versions 9 and above, a server error is resumable if it has the ResumableChangeStreamError label.
	if cs.wireVersion != nil && cs.wireVersion.Includes(minResumableLabelWireVersion) {
		return commandErr.HasErrorLabel(resumableErrorLabel)
	}

	// For wire versions below 9, a server error is resumable if its code is on the allowlist.
	_, resumable := resumableChangeStreamErrors[commandErr.Code]
	return resumable
}

// Returns true if the underlying cursor's batch is empty
func (cs *ChangeStream) emptyBatch() bool {
	return cs.cursor.Batch().Empty()
}

// StreamType represents the cluster type against which a ChangeStream was created.
type StreamType uint8

// These constants represent valid change stream types. A change stream can be initialized over a collection, all
// collections in a database, or over a cluster.
const (
	CollectionStream StreamType = iota
	DatabaseStream
	ClientStream
)
