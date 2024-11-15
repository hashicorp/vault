// Copyright (C) MongoDB, Inc. 2019-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package operation

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/driverutil"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Find performs a find operation.
type Find struct {
	authenticator       driver.Authenticator
	allowDiskUse        *bool
	allowPartialResults *bool
	awaitData           *bool
	batchSize           *int32
	collation           bsoncore.Document
	comment             *string
	filter              bsoncore.Document
	hint                bsoncore.Value
	let                 bsoncore.Document
	limit               *int64
	max                 bsoncore.Document
	maxTime             *time.Duration
	min                 bsoncore.Document
	noCursorTimeout     *bool
	oplogReplay         *bool
	projection          bsoncore.Document
	returnKey           *bool
	showRecordID        *bool
	singleBatch         *bool
	skip                *int64
	snapshot            *bool
	sort                bsoncore.Document
	tailable            *bool
	session             *session.Client
	clock               *session.ClusterClock
	collection          string
	monitor             *event.CommandMonitor
	crypt               driver.Crypt
	database            string
	deployment          driver.Deployment
	readConcern         *readconcern.ReadConcern
	readPreference      *readpref.ReadPref
	selector            description.ServerSelector
	retry               *driver.RetryMode
	result              driver.CursorResponse
	serverAPI           *driver.ServerAPIOptions
	timeout             *time.Duration
	omitCSOTMaxTimeMS   bool
	logger              *logger.Logger
}

// NewFind constructs and returns a new Find.
func NewFind(filter bsoncore.Document) *Find {
	return &Find{
		filter: filter,
	}
}

// Result returns the result of executing this operation.
func (f *Find) Result(opts driver.CursorOptions) (*driver.BatchCursor, error) {
	opts.ServerAPI = f.serverAPI
	return driver.NewBatchCursor(f.result, f.session, f.clock, opts)
}

func (f *Find) processResponse(info driver.ResponseInfo) error {
	var err error
	f.result, err = driver.NewCursorResponse(info)
	return err
}

// Execute runs this operations and returns an error if the operation did not execute successfully.
func (f *Find) Execute(ctx context.Context) error {
	if f.deployment == nil {
		return errors.New("the Find operation must have a Deployment set before Execute can be called")
	}

	return driver.Operation{
		CommandFn:         f.command,
		ProcessResponseFn: f.processResponse,
		RetryMode:         f.retry,
		Type:              driver.Read,
		Client:            f.session,
		Clock:             f.clock,
		CommandMonitor:    f.monitor,
		Crypt:             f.crypt,
		Database:          f.database,
		Deployment:        f.deployment,
		MaxTime:           f.maxTime,
		ReadConcern:       f.readConcern,
		ReadPreference:    f.readPreference,
		Selector:          f.selector,
		Legacy:            driver.LegacyFind,
		ServerAPI:         f.serverAPI,
		Timeout:           f.timeout,
		Logger:            f.logger,
		Name:              driverutil.FindOp,
		OmitCSOTMaxTimeMS: f.omitCSOTMaxTimeMS,
		Authenticator:     f.authenticator,
	}.Execute(ctx)

}

func (f *Find) command(dst []byte, desc description.SelectedServer) ([]byte, error) {
	dst = bsoncore.AppendStringElement(dst, "find", f.collection)
	if f.allowDiskUse != nil {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(4) {
			return nil, errors.New("the 'allowDiskUse' command parameter requires a minimum server wire version of 4")
		}
		dst = bsoncore.AppendBooleanElement(dst, "allowDiskUse", *f.allowDiskUse)
	}
	if f.allowPartialResults != nil {
		dst = bsoncore.AppendBooleanElement(dst, "allowPartialResults", *f.allowPartialResults)
	}
	if f.awaitData != nil {
		dst = bsoncore.AppendBooleanElement(dst, "awaitData", *f.awaitData)
	}
	if f.batchSize != nil {
		dst = bsoncore.AppendInt32Element(dst, "batchSize", *f.batchSize)
	}
	if f.collation != nil {
		if desc.WireVersion == nil || !desc.WireVersion.Includes(5) {
			return nil, errors.New("the 'collation' command parameter requires a minimum server wire version of 5")
		}
		dst = bsoncore.AppendDocumentElement(dst, "collation", f.collation)
	}
	if f.comment != nil {
		dst = bsoncore.AppendStringElement(dst, "comment", *f.comment)
	}
	if f.filter != nil {
		dst = bsoncore.AppendDocumentElement(dst, "filter", f.filter)
	}
	if f.hint.Type != bsontype.Type(0) {
		dst = bsoncore.AppendValueElement(dst, "hint", f.hint)
	}
	if f.let != nil {
		dst = bsoncore.AppendDocumentElement(dst, "let", f.let)
	}
	if f.limit != nil {
		dst = bsoncore.AppendInt64Element(dst, "limit", *f.limit)
	}
	if f.max != nil {
		dst = bsoncore.AppendDocumentElement(dst, "max", f.max)
	}
	if f.min != nil {
		dst = bsoncore.AppendDocumentElement(dst, "min", f.min)
	}
	if f.noCursorTimeout != nil {
		dst = bsoncore.AppendBooleanElement(dst, "noCursorTimeout", *f.noCursorTimeout)
	}
	if f.oplogReplay != nil {
		dst = bsoncore.AppendBooleanElement(dst, "oplogReplay", *f.oplogReplay)
	}
	if f.projection != nil {
		dst = bsoncore.AppendDocumentElement(dst, "projection", f.projection)
	}
	if f.returnKey != nil {
		dst = bsoncore.AppendBooleanElement(dst, "returnKey", *f.returnKey)
	}
	if f.showRecordID != nil {
		dst = bsoncore.AppendBooleanElement(dst, "showRecordId", *f.showRecordID)
	}
	if f.singleBatch != nil {
		dst = bsoncore.AppendBooleanElement(dst, "singleBatch", *f.singleBatch)
	}
	if f.skip != nil {
		dst = bsoncore.AppendInt64Element(dst, "skip", *f.skip)
	}
	if f.snapshot != nil {
		dst = bsoncore.AppendBooleanElement(dst, "snapshot", *f.snapshot)
	}
	if f.sort != nil {
		dst = bsoncore.AppendDocumentElement(dst, "sort", f.sort)
	}
	if f.tailable != nil {
		dst = bsoncore.AppendBooleanElement(dst, "tailable", *f.tailable)
	}
	return dst, nil
}

// AllowDiskUse when true allows temporary data to be written to disk during the find command."
func (f *Find) AllowDiskUse(allowDiskUse bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.allowDiskUse = &allowDiskUse
	return f
}

// AllowPartialResults when true allows partial results to be returned if some shards are down.
func (f *Find) AllowPartialResults(allowPartialResults bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.allowPartialResults = &allowPartialResults
	return f
}

// AwaitData when true makes a cursor block before returning when no data is available.
func (f *Find) AwaitData(awaitData bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.awaitData = &awaitData
	return f
}

// BatchSize specifies the number of documents to return in every batch.
func (f *Find) BatchSize(batchSize int32) *Find {
	if f == nil {
		f = new(Find)
	}

	f.batchSize = &batchSize
	return f
}

// Collation specifies a collation to be used.
func (f *Find) Collation(collation bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.collation = collation
	return f
}

// Comment sets a string to help trace an operation.
func (f *Find) Comment(comment string) *Find {
	if f == nil {
		f = new(Find)
	}

	f.comment = &comment
	return f
}

// Filter determines what results are returned from find.
func (f *Find) Filter(filter bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.filter = filter
	return f
}

// Hint specifies the index to use.
func (f *Find) Hint(hint bsoncore.Value) *Find {
	if f == nil {
		f = new(Find)
	}

	f.hint = hint
	return f
}

// Let specifies the let document to use. This option is only valid for server versions 5.0 and above.
func (f *Find) Let(let bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.let = let
	return f
}

// Limit sets a limit on the number of documents to return.
func (f *Find) Limit(limit int64) *Find {
	if f == nil {
		f = new(Find)
	}

	f.limit = &limit
	return f
}

// Max sets an exclusive upper bound for a specific index.
func (f *Find) Max(max bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.max = max
	return f
}

// MaxTime specifies the maximum amount of time to allow the query to run on the server.
func (f *Find) MaxTime(maxTime *time.Duration) *Find {
	if f == nil {
		f = new(Find)
	}

	f.maxTime = maxTime
	return f
}

// Min sets an inclusive lower bound for a specific index.
func (f *Find) Min(min bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.min = min
	return f
}

// NoCursorTimeout when true prevents cursor from timing out after an inactivity period.
func (f *Find) NoCursorTimeout(noCursorTimeout bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.noCursorTimeout = &noCursorTimeout
	return f
}

// OplogReplay when true replays a replica set's oplog.
func (f *Find) OplogReplay(oplogReplay bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.oplogReplay = &oplogReplay
	return f
}

// Projection limits the fields returned for all documents.
func (f *Find) Projection(projection bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.projection = projection
	return f
}

// ReturnKey when true returns index keys for all result documents.
func (f *Find) ReturnKey(returnKey bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.returnKey = &returnKey
	return f
}

// ShowRecordID when true adds a $recordId field with the record identifier to returned documents.
func (f *Find) ShowRecordID(showRecordID bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.showRecordID = &showRecordID
	return f
}

// SingleBatch specifies whether the results should be returned in a single batch.
func (f *Find) SingleBatch(singleBatch bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.singleBatch = &singleBatch
	return f
}

// Skip specifies the number of documents to skip before returning.
func (f *Find) Skip(skip int64) *Find {
	if f == nil {
		f = new(Find)
	}

	f.skip = &skip
	return f
}

// Snapshot prevents the cursor from returning a document more than once because of an intervening write operation.
func (f *Find) Snapshot(snapshot bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.snapshot = &snapshot
	return f
}

// Sort specifies the order in which to return results.
func (f *Find) Sort(sort bsoncore.Document) *Find {
	if f == nil {
		f = new(Find)
	}

	f.sort = sort
	return f
}

// Tailable keeps a cursor open and resumable after the last data has been retrieved.
func (f *Find) Tailable(tailable bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.tailable = &tailable
	return f
}

// Session sets the session for this operation.
func (f *Find) Session(session *session.Client) *Find {
	if f == nil {
		f = new(Find)
	}

	f.session = session
	return f
}

// ClusterClock sets the cluster clock for this operation.
func (f *Find) ClusterClock(clock *session.ClusterClock) *Find {
	if f == nil {
		f = new(Find)
	}

	f.clock = clock
	return f
}

// Collection sets the collection that this command will run against.
func (f *Find) Collection(collection string) *Find {
	if f == nil {
		f = new(Find)
	}

	f.collection = collection
	return f
}

// CommandMonitor sets the monitor to use for APM events.
func (f *Find) CommandMonitor(monitor *event.CommandMonitor) *Find {
	if f == nil {
		f = new(Find)
	}

	f.monitor = monitor
	return f
}

// Crypt sets the Crypt object to use for automatic encryption and decryption.
func (f *Find) Crypt(crypt driver.Crypt) *Find {
	if f == nil {
		f = new(Find)
	}

	f.crypt = crypt
	return f
}

// Database sets the database to run this operation against.
func (f *Find) Database(database string) *Find {
	if f == nil {
		f = new(Find)
	}

	f.database = database
	return f
}

// Deployment sets the deployment to use for this operation.
func (f *Find) Deployment(deployment driver.Deployment) *Find {
	if f == nil {
		f = new(Find)
	}

	f.deployment = deployment
	return f
}

// ReadConcern specifies the read concern for this operation.
func (f *Find) ReadConcern(readConcern *readconcern.ReadConcern) *Find {
	if f == nil {
		f = new(Find)
	}

	f.readConcern = readConcern
	return f
}

// ReadPreference set the read preference used with this operation.
func (f *Find) ReadPreference(readPreference *readpref.ReadPref) *Find {
	if f == nil {
		f = new(Find)
	}

	f.readPreference = readPreference
	return f
}

// ServerSelector sets the selector used to retrieve a server.
func (f *Find) ServerSelector(selector description.ServerSelector) *Find {
	if f == nil {
		f = new(Find)
	}

	f.selector = selector
	return f
}

// Retry enables retryable mode for this operation. Retries are handled automatically in driver.Operation.Execute based
// on how the operation is set.
func (f *Find) Retry(retry driver.RetryMode) *Find {
	if f == nil {
		f = new(Find)
	}

	f.retry = &retry
	return f
}

// ServerAPI sets the server API version for this operation.
func (f *Find) ServerAPI(serverAPI *driver.ServerAPIOptions) *Find {
	if f == nil {
		f = new(Find)
	}

	f.serverAPI = serverAPI
	return f
}

// Timeout sets the timeout for this operation.
func (f *Find) Timeout(timeout *time.Duration) *Find {
	if f == nil {
		f = new(Find)
	}

	f.timeout = timeout
	return f
}

// OmitCSOTMaxTimeMS omits the automatically-calculated "maxTimeMS" from the
// command when CSOT is enabled. It does not effect "maxTimeMS" set by
// [Find.MaxTime].
func (f *Find) OmitCSOTMaxTimeMS(omit bool) *Find {
	if f == nil {
		f = new(Find)
	}

	f.omitCSOTMaxTimeMS = omit
	return f
}

// Logger sets the logger for this operation.
func (f *Find) Logger(logger *logger.Logger) *Find {
	if f == nil {
		f = new(Find)
	}

	f.logger = logger
	return f
}

// Authenticator sets the authenticator to use for this operation.
func (f *Find) Authenticator(authenticator driver.Authenticator) *Find {
	if f == nil {
		f = new(Find)
	}

	f.authenticator = authenticator
	return f
}
