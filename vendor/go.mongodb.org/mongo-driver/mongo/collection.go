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
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

// Collection is a handle to a MongoDB collection. It is safe for concurrent use by multiple goroutines.
type Collection struct {
	client         *Client
	db             *Database
	name           string
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	readPreference *readpref.ReadPref
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
	registry       *bsoncodec.Registry
}

// aggregateParams is used to store information to configure an Aggregate operation.
type aggregateParams struct {
	ctx            context.Context
	pipeline       interface{}
	client         *Client
	registry       *bsoncodec.Registry
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	retryRead      bool
	db             string
	col            string
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
	readPreference *readpref.ReadPref
	opts           []*options.AggregateOptions
}

func closeImplicitSession(sess *session.Client) {
	if sess != nil && sess.SessionType == session.Implicit {
		sess.EndSession()
	}
}

func newCollection(db *Database, name string, opts ...*options.CollectionOptions) *Collection {
	collOpt := options.MergeCollectionOptions(opts...)

	rc := db.readConcern
	if collOpt.ReadConcern != nil {
		rc = collOpt.ReadConcern
	}

	wc := db.writeConcern
	if collOpt.WriteConcern != nil {
		wc = collOpt.WriteConcern
	}

	rp := db.readPreference
	if collOpt.ReadPreference != nil {
		rp = collOpt.ReadPreference
	}

	reg := db.registry
	if collOpt.Registry != nil {
		reg = collOpt.Registry
	}

	readSelector := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(rp),
		description.LatencySelector(db.client.localThreshold),
	})

	writeSelector := description.CompositeSelector([]description.ServerSelector{
		description.WriteSelector(),
		description.LatencySelector(db.client.localThreshold),
	})

	coll := &Collection{
		client:         db.client,
		db:             db,
		name:           name,
		readPreference: rp,
		readConcern:    rc,
		writeConcern:   wc,
		readSelector:   readSelector,
		writeSelector:  writeSelector,
		registry:       reg,
	}

	return coll
}

func (coll *Collection) copy() *Collection {
	return &Collection{
		client:         coll.client,
		db:             coll.db,
		name:           coll.name,
		readConcern:    coll.readConcern,
		writeConcern:   coll.writeConcern,
		readPreference: coll.readPreference,
		readSelector:   coll.readSelector,
		writeSelector:  coll.writeSelector,
		registry:       coll.registry,
	}
}

// Clone creates a copy of the Collection configured with the given CollectionOptions.
// The specified options are merged with the existing options on the collection, with the specified options taking
// precedence.
func (coll *Collection) Clone(opts ...*options.CollectionOptions) (*Collection, error) {
	copyColl := coll.copy()
	optsColl := options.MergeCollectionOptions(opts...)

	if optsColl.ReadConcern != nil {
		copyColl.readConcern = optsColl.ReadConcern
	}

	if optsColl.WriteConcern != nil {
		copyColl.writeConcern = optsColl.WriteConcern
	}

	if optsColl.ReadPreference != nil {
		copyColl.readPreference = optsColl.ReadPreference
	}

	if optsColl.Registry != nil {
		copyColl.registry = optsColl.Registry
	}

	copyColl.readSelector = description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(copyColl.readPreference),
		description.LatencySelector(copyColl.client.localThreshold),
	})

	return copyColl, nil
}

// Name returns the name of the collection.
func (coll *Collection) Name() string {
	return coll.name
}

// Database returns the Database that was used to create the Collection.
func (coll *Collection) Database() *Database {
	return coll.db
}

// BulkWrite performs a bulk write operation (https://docs.mongodb.com/manual/core/bulk-write-operations/).
//
// The models parameter must be a slice of operations to be executed in this bulk write. It cannot be nil or empty.
// All of the models must be non-nil. See the mongo.WriteModel documentation for a list of valid model types and
// examples of how they should be used.
//
// The opts parameter can be used to specify options for the operation (see the options.BulkWriteOptions documentation.)
func (coll *Collection) BulkWrite(ctx context.Context, models []WriteModel,
	opts ...*options.BulkWriteOptions) (*BulkWriteResult, error) {

	if len(models) == 0 {
		return nil, ErrEmptySlice
	}

	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return nil, err
		}
		defer sess.EndSession()
	}

	err := coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	for _, model := range models {
		if model == nil {
			return nil, ErrNilDocument
		}
	}

	bwo := options.MergeBulkWriteOptions(opts...)

	op := bulkWrite{
		ordered:                  bwo.Ordered,
		bypassDocumentValidation: bwo.BypassDocumentValidation,
		models:                   models,
		session:                  sess,
		collection:               coll,
		selector:                 selector,
		writeConcern:             wc,
	}

	err = op.execute(ctx)

	return &op.result, replaceErrors(err)
}

func (coll *Collection) insert(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) ([]interface{}, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	result := make([]interface{}, len(documents))
	docs := make([]bsoncore.Document, len(documents))

	for i, doc := range documents {
		var err error
		docs[i], result[i], err = transformAndEnsureIDv2(coll.registry, doc)
		if err != nil {
			return nil, err
		}
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return nil, err
		}
		defer sess.EndSession()
	}

	err := coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	op := operation.NewInsert(docs...).
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.crypt)
	imo := options.MergeInsertManyOptions(opts...)
	if imo.BypassDocumentValidation != nil && *imo.BypassDocumentValidation {
		op = op.BypassDocumentValidation(*imo.BypassDocumentValidation)
	}
	if imo.Ordered != nil {
		op = op.Ordered(*imo.Ordered)
	}
	retry := driver.RetryNone
	if coll.client.retryWrites {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(ctx)
	wce, ok := err.(driver.WriteCommandError)
	if !ok {
		return result, err
	}

	// remove the ids that had writeErrors from result
	for i, we := range wce.WriteErrors {
		// i indexes have been removed before the current error, so the index is we.Index-i
		idIndex := int(we.Index) - i
		// if the insert is ordered, nothing after the error was inserted
		if imo.Ordered == nil || *imo.Ordered {
			result = result[:idIndex]
			break
		}
		result = append(result[:idIndex], result[idIndex+1:]...)
	}

	return result, err
}

// InsertOne executes an insert command to insert a single document into the collection.
//
// The document parameter must be the document to be inserted. It cannot be nil. If the document does not have an _id
// field when transformed into BSON, one will be added automatically to the marshalled document. The original document
// will not be modified. The _id can be retrieved from the InsertedID field of the returned InsertOneResult.
//
// The opts parameter can be used to specify options for the operation (see the options.InsertOneOptions documentation.)
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/insert/.
func (coll *Collection) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*InsertOneResult, error) {

	ioOpts := options.MergeInsertOneOptions(opts...)
	imOpts := options.InsertMany()

	if ioOpts.BypassDocumentValidation != nil && *ioOpts.BypassDocumentValidation {
		imOpts.SetBypassDocumentValidation(*ioOpts.BypassDocumentValidation)
	}
	res, err := coll.insert(ctx, []interface{}{document}, imOpts)

	rr, err := processWriteError(err)
	if rr&rrOne == 0 {
		return nil, err
	}
	return &InsertOneResult{InsertedID: res[0]}, err
}

// InsertMany executes an insert command to insert multiple documents into the collection. If write errors occur
// during the operation (e.g. duplicate key error), this method returns a BulkWriteException error.
//
// The documents parameter must be a slice of documents to insert. The slice cannot be nil or empty. The elements must
// all be non-nil. For any document that does not have an _id field when transformed into BSON, one will be added
// automatically to the marshalled document. The original document will not be modified. The _id values for the inserted
// documents can be retrieved from the InsertedIDs field of the returnd InsertManyResult.
//
// The opts parameter can be used to specify options for the operation (see the options.InsertManyOptions documentation.)
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/insert/.
func (coll *Collection) InsertMany(ctx context.Context, documents []interface{},
	opts ...*options.InsertManyOptions) (*InsertManyResult, error) {

	if len(documents) == 0 {
		return nil, ErrEmptySlice
	}

	result, err := coll.insert(ctx, documents, opts...)
	rr, err := processWriteError(err)
	if rr&rrMany == 0 {
		return nil, err
	}

	imResult := &InsertManyResult{InsertedIDs: result}
	writeException, ok := err.(WriteException)
	if !ok {
		return imResult, err
	}

	// create and return a BulkWriteException
	bwErrors := make([]BulkWriteError, 0, len(writeException.WriteErrors))
	for _, we := range writeException.WriteErrors {
		bwErrors = append(bwErrors, BulkWriteError{
			WriteError{
				Index:   we.Index,
				Code:    we.Code,
				Message: we.Message,
			},
			nil,
		})
	}

	return imResult, BulkWriteException{
		WriteErrors:       bwErrors,
		WriteConcernError: writeException.WriteConcernError,
		Labels:            writeException.Labels,
	}
}

func (coll *Collection) delete(ctx context.Context, filter interface{}, deleteOne bool, expectedRr returnResult,
	opts ...*options.DeleteOptions) (*DeleteResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := transformBsoncoreDocument(coll.registry, filter)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return nil, err
		}
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	var limit int32
	if deleteOne {
		limit = 1
	}
	do := options.MergeDeleteOptions(opts...)
	didx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendDocumentElement(doc, "q", f)
	doc = bsoncore.AppendInt32Element(doc, "limit", limit)
	if do.Collation != nil {
		doc = bsoncore.AppendDocumentElement(doc, "collation", do.Collation.ToDocument())
	}
	if do.Hint != nil {
		hint, err := transformValue(coll.registry, do.Hint)
		if err != nil {
			return nil, err
		}

		doc = bsoncore.AppendValueElement(doc, "hint", hint)
	}
	doc, _ = bsoncore.AppendDocumentEnd(doc, didx)

	op := operation.NewDelete(doc).
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.crypt)
	if do.Hint != nil {
		op = op.Hint(true)
	}

	// deleteMany cannot be retried
	retryMode := driver.RetryNone
	if deleteOne && coll.client.retryWrites {
		retryMode = driver.RetryOncePerCommand
	}
	op = op.Retry(retryMode)
	rr, err := processWriteError(op.Execute(ctx))
	if rr&expectedRr == 0 {
		return nil, err
	}
	return &DeleteResult{DeletedCount: int64(op.Result().N)}, err
}

// DeleteOne executes a delete command to delete at most one document from the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// deleted. It cannot be nil. If the filter does not match any documents, the operation will succeed and a DeleteResult
// with a DeletedCount of 0 will be returned. If the filter matches multiple documents, one will be selected from the
// matched set.
//
// The opts parameter can be used to specify options for the operation (see the options.DeleteOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/delete/.
func (coll *Collection) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*DeleteResult, error) {

	return coll.delete(ctx, filter, true, rrOne, opts...)
}

// DeleteMany executes a delete command to delete documents from the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the documents to
// be deleted. It cannot be nil. An empty document (e.g. bson.D{}) should be used to delete all documents in the
// collection. If the filter does not match any documents, the operation will succeed and a DeleteResult with a
// DeletedCount of 0 will be returned.
//
// The opts parameter can be used to specify options for the operation (see the options.DeleteOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/delete/.
func (coll *Collection) DeleteMany(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*DeleteResult, error) {

	return coll.delete(ctx, filter, false, rrMany, opts...)
}

func (coll *Collection) updateOrReplace(ctx context.Context, filter bsoncore.Document, update interface{}, multi bool,
	expectedRr returnResult, checkDollarKey bool, opts ...*options.UpdateOptions) (*UpdateResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	uo := options.MergeUpdateOptions(opts...)

	// collation, arrayFilters, upsert, and hint are included on the individual update documents rather than as part of the
	// command
	updateDoc, err := createUpdateDoc(filter, update, uo.Hint, uo.ArrayFilters, uo.Collation, uo.Upsert, multi,
		checkDollarKey, coll.registry)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return nil, err
		}
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	op := operation.NewUpdate(updateDoc).
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.crypt).Hint(uo.Hint != nil).
		ArrayFilters(uo.ArrayFilters != nil)

	if uo.BypassDocumentValidation != nil && *uo.BypassDocumentValidation {
		op = op.BypassDocumentValidation(*uo.BypassDocumentValidation)
	}
	retry := driver.RetryNone
	// retryable writes are only enabled updateOne/replaceOne operations
	if !multi && coll.client.retryWrites {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)
	err = op.Execute(ctx)

	rr, err := processWriteError(err)
	if rr&expectedRr == 0 {
		return nil, err
	}

	opRes := op.Result()
	res := &UpdateResult{
		MatchedCount:  int64(opRes.N),
		ModifiedCount: int64(opRes.NModified),
		UpsertedCount: int64(len(opRes.Upserted)),
	}
	if len(opRes.Upserted) > 0 {
		res.UpsertedID = opRes.Upserted[0].ID
		res.MatchedCount--
	}

	return res, err
}

// UpdateOne executes an update command to update at most one document in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// updated. It cannot be nil. If the filter does not match any documents, the operation will succeed and an UpdateResult
// with a MatchedCount of 0 will be returned. If the filter matches multiple documents, one will be selected from the
// matched set and MatchedCount will equal 1.
//
// The update parameter must be a document containing update operators
// (https://docs.mongodb.com/manual/reference/operator/update/) and can be used to specify the modifications to be
// made to the selected document. It cannot be nil or empty.
//
// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/update/.
func (coll *Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*UpdateResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := transformBsoncoreDocument(coll.registry, filter)
	if err != nil {
		return nil, err
	}

	return coll.updateOrReplace(ctx, f, update, false, rrOne, true, opts...)
}

// UpdateMany executes an update command to update documents in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the documents to be
// updated. It cannot be nil. If the filter does not match any documents, the operation will succeed and an UpdateResult
// with a MatchedCount of 0 will be returned.
//
// The update parameter must be a document containing update operators
// (https://docs.mongodb.com/manual/reference/operator/update/) and can be used to specify the modifications to be made
// to the selected documents. It cannot be nil or empty.
//
// The opts parameter can be used to specify options for the operation (see the options.UpdateOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/update/.
func (coll *Collection) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*UpdateResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := transformBsoncoreDocument(coll.registry, filter)
	if err != nil {
		return nil, err
	}

	return coll.updateOrReplace(ctx, f, update, true, rrMany, true, opts...)
}

// ReplaceOne executes an update command to replace at most one document in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// replaced. It cannot be nil. If the filter does not match any documents, the operation will succeed and an
// UpdateResult with a MatchedCount of 0 will be returned. If the filter matches multiple documents, one will be
// selected from the matched set and MatchedCount will equal 1.
//
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators (https://docs.mongodb.com/manual/reference/operator/update/).
//
// The opts parameter can be used to specify options for the operation (see the options.ReplaceOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/update/.
func (coll *Collection) ReplaceOne(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.ReplaceOptions) (*UpdateResult, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := transformBsoncoreDocument(coll.registry, filter)
	if err != nil {
		return nil, err
	}

	r, err := transformBsoncoreDocument(coll.registry, replacement)
	if err != nil {
		return nil, err
	}

	if err := ensureNoDollarKey(r); err != nil {
		return nil, err
	}

	updateOptions := make([]*options.UpdateOptions, 0, len(opts))
	for _, opt := range opts {
		uOpts := options.Update()
		uOpts.BypassDocumentValidation = opt.BypassDocumentValidation
		uOpts.Collation = opt.Collation
		uOpts.Upsert = opt.Upsert
		uOpts.Hint = opt.Hint
		updateOptions = append(updateOptions, uOpts)
	}

	return coll.updateOrReplace(ctx, f, r, false, rrOne, false, updateOptions...)
}

// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
//
// The pipeline parameter must be an array of documents, each representing an aggregation stage. The pipeline cannot
// be nil but can be empty. The stage documents must all be non-nil. For a pipeline of bson.D documents, the
// mongo.Pipeline type can be used. See
// https://docs.mongodb.com/manual/reference/operator/aggregation-pipeline/#db-collection-aggregate-stages for a list of
// valid stages in aggregations.
//
// The opts parameter can be used to specify options for the operation (see the options.AggregateOptions documentation.)
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/aggregate/.
func (coll *Collection) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (*Cursor, error) {
	a := aggregateParams{
		ctx:            ctx,
		pipeline:       pipeline,
		client:         coll.client,
		registry:       coll.registry,
		readConcern:    coll.readConcern,
		writeConcern:   coll.writeConcern,
		retryRead:      coll.client.retryReads,
		db:             coll.db.name,
		col:            coll.name,
		readSelector:   coll.readSelector,
		writeSelector:  coll.writeSelector,
		readPreference: coll.readPreference,
		opts:           opts,
	}
	return aggregate(a)
}

// aggreate is the helper method for Aggregate
func aggregate(a aggregateParams) (*Cursor, error) {

	if a.ctx == nil {
		a.ctx = context.Background()
	}

	pipelineArr, hasOutputStage, err := transformAggregatePipelinev2(a.registry, a.pipeline)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(a.ctx)
	if sess == nil && a.client.sessionPool != nil {
		sess, err = session.NewClientSession(a.client.sessionPool, a.client.id, session.Implicit)
		if err != nil {
			return nil, err
		}
	}
	if err = a.client.validSession(sess); err != nil {
		return nil, err
	}

	var wc *writeconcern.WriteConcern
	if hasOutputStage {
		wc = a.writeConcern
	}
	rc := a.readConcern
	if sess.TransactionRunning() {
		wc = nil
		rc = nil
	}
	if !writeconcern.AckWrite(wc) {
		closeImplicitSession(sess)
		sess = nil
	}

	selector := makePinnedSelector(sess, a.writeSelector)
	if !hasOutputStage {
		selector = makeReadPrefSelector(sess, a.readSelector, a.client.localThreshold)
	}

	ao := options.MergeAggregateOptions(a.opts...)
	cursorOpts := driver.CursorOptions{
		CommandMonitor: a.client.monitor,
		Crypt:          a.client.crypt,
	}

	op := operation.NewAggregate(pipelineArr).
		Session(sess).
		WriteConcern(wc).
		ReadConcern(rc).
		CommandMonitor(a.client.monitor).
		ServerSelector(selector).
		ClusterClock(a.client.clock).
		Database(a.db).
		Collection(a.col).
		Deployment(a.client.deployment).
		Crypt(a.client.crypt)
	if !hasOutputStage {
		// Only pass the user-specified read preference if the aggregation doesn't have a $out or $merge stage.
		// Otherwise, the read preference could be forwarded to a mongos, which would error if the aggregation were
		// executed against a non-primary node.
		op.ReadPreference(a.readPreference)
	}

	if ao.AllowDiskUse != nil {
		op.AllowDiskUse(*ao.AllowDiskUse)
	}
	// ignore batchSize of 0 with $out
	if ao.BatchSize != nil && !(*ao.BatchSize == 0 && hasOutputStage) {
		op.BatchSize(*ao.BatchSize)
		cursorOpts.BatchSize = *ao.BatchSize
	}
	if ao.BypassDocumentValidation != nil && *ao.BypassDocumentValidation {
		op.BypassDocumentValidation(*ao.BypassDocumentValidation)
	}
	if ao.Collation != nil {
		op.Collation(bsoncore.Document(ao.Collation.ToDocument()))
	}
	if ao.MaxTime != nil {
		op.MaxTimeMS(int64(*ao.MaxTime / time.Millisecond))
	}
	if ao.MaxAwaitTime != nil {
		cursorOpts.MaxTimeMS = int64(*ao.MaxAwaitTime / time.Millisecond)
	}
	if ao.Comment != nil {
		op.Comment(*ao.Comment)
	}
	if ao.Hint != nil {
		hintVal, err := transformValue(a.registry, ao.Hint)
		if err != nil {
			closeImplicitSession(sess)
			return nil, err
		}
		op.Hint(hintVal)
	}

	retry := driver.RetryNone
	if a.retryRead && !hasOutputStage {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(a.ctx)
	if err != nil {
		closeImplicitSession(sess)
		if wce, ok := err.(driver.WriteCommandError); ok && wce.WriteConcernError != nil {
			return nil, *convertDriverWriteConcernError(wce.WriteConcernError)
		}
		return nil, replaceErrors(err)
	}

	bc, err := op.Result(cursorOpts)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}
	cursor, err := newCursorWithSession(bc, a.registry, sess)
	return cursor, replaceErrors(err)
}

// CountDocuments returns the number of documents in the collection. For a fast count of the documents in the
// collection, see the EstimatedDocumentCount method.
//
// The filter parameter must be a document and can be used to select which documents contribute to the count. It
// cannot be nil. An empty document (e.g. bson.D{}) should be used to count all documents in the collection. This will
// result in a full collection scan.
//
// The opts parameter can be used to specify options for the operation (see the options.CountOptions documentation).
func (coll *Collection) CountDocuments(ctx context.Context, filter interface{},
	opts ...*options.CountOptions) (int64, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	countOpts := options.MergeCountOptions(opts...)

	pipelineArr, err := countDocumentsAggregatePipeline(coll.registry, filter, countOpts)
	if err != nil {
		return 0, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return 0, err
		}
		defer sess.EndSession()
	}
	if err = coll.client.validSession(sess); err != nil {
		return 0, err
	}

	rc := coll.readConcern
	if sess.TransactionRunning() {
		rc = nil
	}

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	op := operation.NewAggregate(pipelineArr).Session(sess).ReadConcern(rc).ReadPreference(coll.readPreference).
		CommandMonitor(coll.client.monitor).ServerSelector(selector).ClusterClock(coll.client.clock).Database(coll.db.name).
		Collection(coll.name).Deployment(coll.client.deployment).Crypt(coll.client.crypt)
	if countOpts.Collation != nil {
		op.Collation(bsoncore.Document(countOpts.Collation.ToDocument()))
	}
	if countOpts.MaxTime != nil {
		op.MaxTimeMS(int64(*countOpts.MaxTime / time.Millisecond))
	}
	if countOpts.Hint != nil {
		hintVal, err := transformValue(coll.registry, countOpts.Hint)
		if err != nil {
			return 0, err
		}
		op.Hint(hintVal)
	}
	retry := driver.RetryNone
	if coll.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil {
		return 0, replaceErrors(err)
	}

	batch := op.ResultCursorResponse().FirstBatch
	if batch == nil {
		return 0, errors.New("invalid response from server, no 'firstBatch' field")
	}

	docs, err := batch.Documents()
	if err != nil || len(docs) == 0 {
		return 0, nil
	}

	val, ok := docs[0].Lookup("n").AsInt64OK()
	if !ok {
		return 0, errors.New("invalid response from server, no 'n' field")
	}

	return val, nil
}

// EstimatedDocumentCount executes a count command and returns an estimate of the number of documents in the collection
// using collection metadata.
//
// The opts parameter can be used to specify options for the operation (see the options.EstimatedDocumentCountOptions
// documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/count/.
func (coll *Collection) EstimatedDocumentCount(ctx context.Context,
	opts ...*options.EstimatedDocumentCountOptions) (int64, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)

	var err error
	if sess == nil && coll.client.sessionPool != nil {
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return 0, err
		}
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return 0, err
	}

	rc := coll.readConcern
	if sess.TransactionRunning() {
		rc = nil
	}

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	op := operation.NewCount().Session(sess).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).CommandMonitor(coll.client.monitor).
		Deployment(coll.client.deployment).ReadConcern(rc).ReadPreference(coll.readPreference).
		ServerSelector(selector).Crypt(coll.client.crypt)

	co := options.MergeEstimatedDocumentCountOptions(opts...)
	if co.MaxTime != nil {
		op = op.MaxTimeMS(int64(*co.MaxTime / time.Millisecond))
	}
	retry := driver.RetryNone
	if coll.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op.Retry(retry)

	err = op.Execute(ctx)

	return op.Result().N, replaceErrors(err)
}

// Distinct executes a distinct command to find the unique values for a specified field in the collection.
//
// The fieldName parameter specifies the field name for which distinct values should be returned.
//
// The filter parameter must be a document containing query operators and can be used to select which documents are
// considered. It cannot be nil. An empty document (e.g. bson.D{}) should be used to select all documents.
//
// The opts parameter can be used to specify options for the operation (see the options.DistinctOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/distinct/.
func (coll *Collection) Distinct(ctx context.Context, fieldName string, filter interface{},
	opts ...*options.DistinctOptions) ([]interface{}, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := transformBsoncoreDocument(coll.registry, filter)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)

	if sess == nil && coll.client.sessionPool != nil {
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return nil, err
		}
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return nil, err
	}

	rc := coll.readConcern
	if sess.TransactionRunning() {
		rc = nil
	}

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	option := options.MergeDistinctOptions(opts...)

	op := operation.NewDistinct(fieldName, bsoncore.Document(f)).
		Session(sess).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).CommandMonitor(coll.client.monitor).
		Deployment(coll.client.deployment).ReadConcern(rc).ReadPreference(coll.readPreference).
		ServerSelector(selector).Crypt(coll.client.crypt)

	if option.Collation != nil {
		op.Collation(bsoncore.Document(option.Collation.ToDocument()))
	}
	if option.MaxTime != nil {
		op.MaxTimeMS(int64(*option.MaxTime / time.Millisecond))
	}
	retry := driver.RetryNone
	if coll.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil {
		return nil, replaceErrors(err)
	}

	arr, ok := op.Result().Values.ArrayOK()
	if !ok {
		return nil, fmt.Errorf("response field 'values' is type array, but received BSON type %s", op.Result().Values.Type)
	}

	values, err := arr.Values()
	if err != nil {
		return nil, err
	}

	retArray := make([]interface{}, len(values))

	for i, val := range values {
		raw := bson.RawValue{Type: val.Type, Value: val.Data}
		err = raw.Unmarshal(&retArray[i])
		if err != nil {
			return nil, err
		}
	}

	return retArray, replaceErrors(err)
}

// Find executes a find command and returns a Cursor over the matching documents in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select which documents are
// included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all documents.
//
// The opts parameter can be used to specify options for the operation (see the options.FindOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/find/.
func (coll *Collection) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (*Cursor, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := transformBsoncoreDocument(coll.registry, filter)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return nil, err
		}
	}

	err = coll.client.validSession(sess)
	if err != nil {
		closeImplicitSession(sess)
		return nil, err
	}

	rc := coll.readConcern
	if sess.TransactionRunning() {
		rc = nil
	}

	selector := makeReadPrefSelector(sess, coll.readSelector, coll.client.localThreshold)
	op := operation.NewFind(f).
		Session(sess).ReadConcern(rc).ReadPreference(coll.readPreference).
		CommandMonitor(coll.client.monitor).ServerSelector(selector).
		ClusterClock(coll.client.clock).Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.crypt)

	fo := options.MergeFindOptions(opts...)
	cursorOpts := driver.CursorOptions{
		CommandMonitor: coll.client.monitor,
		Crypt:          coll.client.crypt,
	}

	if fo.AllowDiskUse != nil {
		op.AllowDiskUse(*fo.AllowDiskUse)
	}
	if fo.AllowPartialResults != nil {
		op.AllowPartialResults(*fo.AllowPartialResults)
	}
	if fo.BatchSize != nil {
		cursorOpts.BatchSize = *fo.BatchSize
		op.BatchSize(*fo.BatchSize)
	}
	if fo.Collation != nil {
		op.Collation(bsoncore.Document(fo.Collation.ToDocument()))
	}
	if fo.Comment != nil {
		op.Comment(*fo.Comment)
	}
	if fo.CursorType != nil {
		switch *fo.CursorType {
		case options.Tailable:
			op.Tailable(true)
		case options.TailableAwait:
			op.Tailable(true)
			op.AwaitData(true)
		}
	}
	if fo.Hint != nil {
		hint, err := transformValue(coll.registry, fo.Hint)
		if err != nil {
			closeImplicitSession(sess)
			return nil, err
		}
		op.Hint(hint)
	}
	if fo.Limit != nil {
		limit := *fo.Limit
		if limit < 0 {
			limit = -1 * limit
			op.SingleBatch(true)
		}
		cursorOpts.Limit = int32(limit)
		op.Limit(limit)
	}
	if fo.Max != nil {
		max, err := transformBsoncoreDocument(coll.registry, fo.Max)
		if err != nil {
			closeImplicitSession(sess)
			return nil, err
		}
		op.Max(max)
	}
	if fo.MaxAwaitTime != nil {
		cursorOpts.MaxTimeMS = int64(*fo.MaxAwaitTime / time.Millisecond)
	}
	if fo.MaxTime != nil {
		op.MaxTimeMS(int64(*fo.MaxTime / time.Millisecond))
	}
	if fo.Min != nil {
		min, err := transformBsoncoreDocument(coll.registry, fo.Min)
		if err != nil {
			closeImplicitSession(sess)
			return nil, err
		}
		op.Min(min)
	}
	if fo.NoCursorTimeout != nil {
		op.NoCursorTimeout(*fo.NoCursorTimeout)
	}
	if fo.OplogReplay != nil {
		op.OplogReplay(*fo.OplogReplay)
	}
	if fo.Projection != nil {
		proj, err := transformBsoncoreDocument(coll.registry, fo.Projection)
		if err != nil {
			closeImplicitSession(sess)
			return nil, err
		}
		op.Projection(proj)
	}
	if fo.ReturnKey != nil {
		op.ReturnKey(*fo.ReturnKey)
	}
	if fo.ShowRecordID != nil {
		op.ShowRecordID(*fo.ShowRecordID)
	}
	if fo.Skip != nil {
		op.Skip(*fo.Skip)
	}
	if fo.Snapshot != nil {
		op.Snapshot(*fo.Snapshot)
	}
	if fo.Sort != nil {
		sort, err := transformBsoncoreDocument(coll.registry, fo.Sort)
		if err != nil {
			closeImplicitSession(sess)
			return nil, err
		}
		op.Sort(sort)
	}
	retry := driver.RetryNone
	if coll.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	if err = op.Execute(ctx); err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}

	bc, err := op.Result(cursorOpts)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}
	return newCursorWithSession(bc, coll.registry, sess)
}

// FindOne executes a find command and returns a SingleResult for one document in the collection.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// returned. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
// ErrNoDocuments will be returned. If the filter matches multiple documents, one will be selected from the matched set.
//
// The opts parameter can be used to specify options for this operation (see the options.FindOneOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/find/.
func (coll *Collection) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *SingleResult {

	if ctx == nil {
		ctx = context.Background()
	}

	findOpts := make([]*options.FindOptions, len(opts))
	for i, opt := range opts {
		findOpts[i] = &options.FindOptions{
			AllowPartialResults: opt.AllowPartialResults,
			BatchSize:           opt.BatchSize,
			Collation:           opt.Collation,
			Comment:             opt.Comment,
			CursorType:          opt.CursorType,
			Hint:                opt.Hint,
			Max:                 opt.Max,
			MaxAwaitTime:        opt.MaxAwaitTime,
			MaxTime:             opt.MaxTime,
			Min:                 opt.Min,
			NoCursorTimeout:     opt.NoCursorTimeout,
			OplogReplay:         opt.OplogReplay,
			Projection:          opt.Projection,
			ReturnKey:           opt.ReturnKey,
			ShowRecordID:        opt.ShowRecordID,
			Skip:                opt.Skip,
			Snapshot:            opt.Snapshot,
			Sort:                opt.Sort,
		}
	}
	// Unconditionally send a limit to make sure only one document is returned and the cursor is not kept open
	// by the server.
	findOpts = append(findOpts, options.Find().SetLimit(-1))

	cursor, err := coll.Find(ctx, filter, findOpts...)
	return &SingleResult{cur: cursor, reg: coll.registry, err: replaceErrors(err)}
}

func (coll *Collection) findAndModify(ctx context.Context, op *operation.FindAndModify) *SingleResult {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	var err error
	if sess == nil && coll.client.sessionPool != nil {
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return &SingleResult{err: err}
		}
		defer sess.EndSession()
	}

	err = coll.client.validSession(sess)
	if err != nil {
		return &SingleResult{err: err}
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	retry := driver.RetryNone
	if coll.client.retryWrites {
		retry = driver.RetryOnce
	}

	op = op.Session(sess).
		WriteConcern(wc).
		CommandMonitor(coll.client.monitor).
		ServerSelector(selector).
		ClusterClock(coll.client.clock).
		Database(coll.db.name).
		Collection(coll.name).
		Deployment(coll.client.deployment).
		Retry(retry).
		Crypt(coll.client.crypt)

	_, err = processWriteError(op.Execute(ctx))
	if err != nil {
		return &SingleResult{err: err}
	}

	return &SingleResult{rdr: bson.Raw(op.Result().Value), reg: coll.registry}
}

// FindOneAndDelete executes a findAndModify command to delete at most one document in the collection. and returns the
// document as it appeared before deletion.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// deleted. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
//
// The opts parameter can be used to specify options for the operation (see the options.FindOneAndDeleteOptions
// documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/findAndModify/.
func (coll *Collection) FindOneAndDelete(ctx context.Context, filter interface{},
	opts ...*options.FindOneAndDeleteOptions) *SingleResult {

	f, err := transformBsoncoreDocument(coll.registry, filter)
	if err != nil {
		return &SingleResult{err: err}
	}
	fod := options.MergeFindOneAndDeleteOptions(opts...)
	op := operation.NewFindAndModify(f).Remove(true)
	if fod.Collation != nil {
		op = op.Collation(bsoncore.Document(fod.Collation.ToDocument()))
	}
	if fod.MaxTime != nil {
		op = op.MaxTimeMS(int64(*fod.MaxTime / time.Millisecond))
	}
	if fod.Projection != nil {
		proj, err := transformBsoncoreDocument(coll.registry, fod.Projection)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Fields(proj)
	}
	if fod.Sort != nil {
		sort, err := transformBsoncoreDocument(coll.registry, fod.Sort)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Sort(sort)
	}
	if fod.Hint != nil {
		hint, err := transformValue(coll.registry, fod.Hint)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Hint(hint)
	}

	return coll.findAndModify(ctx, op)
}

// FindOneAndReplace executes a findAndModify command to replace at most one document in the collection
// and returns the document as it appeared before replacement.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// replaced. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
//
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators (https://docs.mongodb.com/manual/reference/operator/update/).
//
// The opts parameter can be used to specify options for the operation (see the options.FindOneAndReplaceOptions
// documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/findAndModify/.
func (coll *Collection) FindOneAndReplace(ctx context.Context, filter interface{},
	replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *SingleResult {

	f, err := transformBsoncoreDocument(coll.registry, filter)
	if err != nil {
		return &SingleResult{err: err}
	}
	r, err := transformBsoncoreDocument(coll.registry, replacement)
	if err != nil {
		return &SingleResult{err: err}
	}
	if firstElem, err := r.IndexErr(0); err == nil && strings.HasPrefix(firstElem.Key(), "$") {
		return &SingleResult{err: errors.New("replacement document cannot contain keys beginning with '$'")}
	}

	fo := options.MergeFindOneAndReplaceOptions(opts...)
	op := operation.NewFindAndModify(f).Update(bsoncore.Value{Type: bsontype.EmbeddedDocument, Data: r})
	if fo.BypassDocumentValidation != nil && *fo.BypassDocumentValidation {
		op = op.BypassDocumentValidation(*fo.BypassDocumentValidation)
	}
	if fo.Collation != nil {
		op = op.Collation(bsoncore.Document(fo.Collation.ToDocument()))
	}
	if fo.MaxTime != nil {
		op = op.MaxTimeMS(int64(*fo.MaxTime / time.Millisecond))
	}
	if fo.Projection != nil {
		proj, err := transformBsoncoreDocument(coll.registry, fo.Projection)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Fields(proj)
	}
	if fo.ReturnDocument != nil {
		op = op.NewDocument(*fo.ReturnDocument == options.After)
	}
	if fo.Sort != nil {
		sort, err := transformBsoncoreDocument(coll.registry, fo.Sort)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Sort(sort)
	}
	if fo.Upsert != nil {
		op = op.Upsert(*fo.Upsert)
	}
	if fo.Hint != nil {
		hint, err := transformValue(coll.registry, fo.Hint)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Hint(hint)
	}

	return coll.findAndModify(ctx, op)
}

// FindOneAndUpdate executes a findAndModify command to update at most one document in the collection and returns the
// document as it appeared before updating.
//
// The filter parameter must be a document containing query operators and can be used to select the document to be
// updated. It cannot be nil. If the filter does not match any documents, a SingleResult with an error set to
// ErrNoDocuments wil be returned. If the filter matches multiple documents, one will be selected from the matched set.
//
// The update parameter must be a document containing update operators
// (https://docs.mongodb.com/manual/reference/operator/update/) and can be used to specify the modifications to be made
// to the selected document. It cannot be nil or empty.
//
// The opts parameter can be used to specify options for the operation (see the options.FindOneAndUpdateOptions
// documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/findAndModify/.
func (coll *Collection) FindOneAndUpdate(ctx context.Context, filter interface{},
	update interface{}, opts ...*options.FindOneAndUpdateOptions) *SingleResult {

	if ctx == nil {
		ctx = context.Background()
	}

	f, err := transformBsoncoreDocument(coll.registry, filter)
	if err != nil {
		return &SingleResult{err: err}
	}

	fo := options.MergeFindOneAndUpdateOptions(opts...)
	op := operation.NewFindAndModify(f)

	u, err := transformUpdateValue(coll.registry, update, true)
	if err != nil {
		return &SingleResult{err: err}
	}
	op = op.Update(u)

	if fo.ArrayFilters != nil {
		filtersDoc, err := fo.ArrayFilters.ToArrayDocument()
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.ArrayFilters(bsoncore.Document(filtersDoc))
	}
	if fo.BypassDocumentValidation != nil && *fo.BypassDocumentValidation {
		op = op.BypassDocumentValidation(*fo.BypassDocumentValidation)
	}
	if fo.Collation != nil {
		op = op.Collation(bsoncore.Document(fo.Collation.ToDocument()))
	}
	if fo.MaxTime != nil {
		op = op.MaxTimeMS(int64(*fo.MaxTime / time.Millisecond))
	}
	if fo.Projection != nil {
		proj, err := transformBsoncoreDocument(coll.registry, fo.Projection)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Fields(proj)
	}
	if fo.ReturnDocument != nil {
		op = op.NewDocument(*fo.ReturnDocument == options.After)
	}
	if fo.Sort != nil {
		sort, err := transformBsoncoreDocument(coll.registry, fo.Sort)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Sort(sort)
	}
	if fo.Upsert != nil {
		op = op.Upsert(*fo.Upsert)
	}
	if fo.Hint != nil {
		hint, err := transformValue(coll.registry, fo.Hint)
		if err != nil {
			return &SingleResult{err: err}
		}
		op = op.Hint(hint)
	}

	return coll.findAndModify(ctx, op)
}

// Watch returns a change stream for all changes on the corresponding collection. See
// https://docs.mongodb.com/manual/changeStreams/ for more information about change streams.
//
// The Collection must be configured with read concern majority or no read concern for a change stream to be created
// successfully.
//
// The pipeline parameter must be an array of documents, each representing a pipeline stage. The pipeline cannot be
// nil but can be empty. The stage documents must all be non-nil. See https://docs.mongodb.com/manual/changeStreams/ for
// a list of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the
// mongo.Pipeline{} type can be used.
//
// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
// documentation).
func (coll *Collection) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) {

	csConfig := changeStreamConfig{
		readConcern:    coll.readConcern,
		readPreference: coll.readPreference,
		client:         coll.client,
		registry:       coll.registry,
		streamType:     CollectionStream,
		collectionName: coll.Name(),
		databaseName:   coll.db.Name(),
		crypt:          coll.client.crypt,
	}
	return newChangeStream(ctx, csConfig, pipeline, opts...)
}

// Indexes returns an IndexView instance that can be used to perform operations on the indexes for the collection.
func (coll *Collection) Indexes() IndexView {
	return IndexView{coll: coll}
}

// Drop drops the collection on the server. This method ignores "namespace not found" errors so it is safe to drop
// a collection that does not exist on the server.
func (coll *Collection) Drop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && coll.client.sessionPool != nil {
		var err error
		sess, err = session.NewClientSession(coll.client.sessionPool, coll.client.id, session.Implicit)
		if err != nil {
			return err
		}
		defer sess.EndSession()
	}

	err := coll.client.validSession(sess)
	if err != nil {
		return err
	}

	wc := coll.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, coll.writeSelector)

	op := operation.NewDropCollection().
		Session(sess).WriteConcern(wc).CommandMonitor(coll.client.monitor).
		ServerSelector(selector).ClusterClock(coll.client.clock).
		Database(coll.db.name).Collection(coll.name).
		Deployment(coll.client.deployment).Crypt(coll.client.crypt)
	err = op.Execute(ctx)

	// ignore namespace not found erorrs
	driverErr, ok := err.(driver.Error)
	if !ok || (ok && !driverErr.NamespaceNotFound()) {
		return replaceErrors(err)
	}
	return nil
}

// makePinnedSelector makes a selector for a pinned session with a pinned server. Will attempt to do server selection on
// the pinned server but if that fails it will go through a list of default selectors
func makePinnedSelector(sess *session.Client, defaultSelector description.ServerSelector) description.ServerSelectorFunc {
	return func(t description.Topology, svrs []description.Server) ([]description.Server, error) {
		if sess != nil && sess.PinnedServer != nil {
			return sess.PinnedServer.SelectServer(t, svrs)
		}

		return defaultSelector.SelectServer(t, svrs)
	}
}

func makeReadPrefSelector(sess *session.Client, selector description.ServerSelector, localThreshold time.Duration) description.ServerSelectorFunc {
	if sess != nil && sess.TransactionRunning() {
		selector = description.CompositeSelector([]description.ServerSelector{
			description.ReadPrefSelector(sess.CurrentRp),
			description.LatencySelector(localThreshold),
		})
	}

	return makePinnedSelector(sess, selector)
}
