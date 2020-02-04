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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
)

var (
	defaultRunCmdOpts = []*options.RunCmdOptions{options.RunCmd().SetReadPreference(readpref.Primary())}
)

// Database is a handle to a MongoDB database. It is safe for concurrent use by multiple goroutines.
type Database struct {
	client         *Client
	name           string
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	readPreference *readpref.ReadPref
	readSelector   description.ServerSelector
	writeSelector  description.ServerSelector
	registry       *bsoncodec.Registry
}

func newDatabase(client *Client, name string, opts ...*options.DatabaseOptions) *Database {
	dbOpt := options.MergeDatabaseOptions(opts...)

	rc := client.readConcern
	if dbOpt.ReadConcern != nil {
		rc = dbOpt.ReadConcern
	}

	rp := client.readPreference
	if dbOpt.ReadPreference != nil {
		rp = dbOpt.ReadPreference
	}

	wc := client.writeConcern
	if dbOpt.WriteConcern != nil {
		wc = dbOpt.WriteConcern
	}

	reg := client.registry
	if dbOpt.Registry != nil {
		reg = dbOpt.Registry
	}

	db := &Database{
		client:         client,
		name:           name,
		readPreference: rp,
		readConcern:    rc,
		writeConcern:   wc,
		registry:       reg,
	}

	db.readSelector = description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(db.readPreference),
		description.LatencySelector(db.client.localThreshold),
	})

	db.writeSelector = description.CompositeSelector([]description.ServerSelector{
		description.WriteSelector(),
		description.LatencySelector(db.client.localThreshold),
	})

	return db
}

// Client returns the Client the Database was created from.
func (db *Database) Client() *Client {
	return db.client
}

// Name returns the name of the database.
func (db *Database) Name() string {
	return db.name
}

// Collection gets a handle for a collection with the given name configured with the given CollectionOptions.
func (db *Database) Collection(name string, opts ...*options.CollectionOptions) *Collection {
	return newCollection(db, name, opts...)
}

// Aggregate executes an aggregate command the database. This requires MongoDB version >= 3.6 and driver version >=
// 1.1.0.
//
// The pipeline parameter must be a slice of documents, each representing an aggregation stage. The pipeline
// cannot be nil but can be empty. The stage documents must all be non-nil. For a pipeline of bson.D documents, the
// mongo.Pipeline type can be used. See
// https://docs.mongodb.com/manual/reference/operator/aggregation-pipeline/#db-aggregate-stages for a list of valid
// stages in database-level aggregations.
//
// The opts parameter can be used to specify options for this operation (see the options.AggregateOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/aggregate/.
func (db *Database) Aggregate(ctx context.Context, pipeline interface{},
	opts ...*options.AggregateOptions) (*Cursor, error) {
	a := aggregateParams{
		ctx:            ctx,
		pipeline:       pipeline,
		client:         db.client,
		registry:       db.registry,
		readConcern:    db.readConcern,
		writeConcern:   db.writeConcern,
		retryRead:      db.client.retryReads,
		db:             db.name,
		readSelector:   db.readSelector,
		writeSelector:  db.writeSelector,
		readPreference: db.readPreference,
		opts:           opts,
	}
	return aggregate(a)
}

func (db *Database) processRunCommand(ctx context.Context, cmd interface{},
	opts ...*options.RunCmdOptions) (*operation.Command, *session.Client, error) {
	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		var err error
		sess, err = session.NewClientSession(db.client.sessionPool, db.client.id, session.Implicit)
		if err != nil {
			return nil, sess, err
		}
	}

	err := db.client.validSession(sess)
	if err != nil {
		return nil, sess, err
	}

	ro := options.MergeRunCmdOptions(append(defaultRunCmdOpts, opts...)...)
	if sess != nil && sess.TransactionRunning() && ro.ReadPreference != nil && ro.ReadPreference.Mode() != readpref.PrimaryMode {
		return nil, sess, errors.New("read preference in a transaction must be primary")
	}

	runCmdDoc, err := transformBsoncoreDocument(db.registry, cmd)
	if err != nil {
		return nil, sess, err
	}
	readSelect := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(ro.ReadPreference),
		description.LatencySelector(db.client.localThreshold),
	})
	if sess != nil && sess.PinnedServer != nil {
		readSelect = sess.PinnedServer
	}

	return operation.NewCommand(runCmdDoc).
		Session(sess).CommandMonitor(db.client.monitor).
		ServerSelector(readSelect).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).ReadConcern(db.readConcern).Crypt(db.client.crypt), sess, nil
}

// RunCommand executes the given command against the database.
//
// The runCommand parameter must be a document for the command to be executed. It cannot be nil.
// This must be an order-preserving type such as bson.D. Map types such as bson.M are not valid.
// If the command document contains a session ID or any transaction-specific fields, the behavior is undefined.
//
// The opts parameter can be used to specify options for this operation (see the options.RunCmdOptions documentation).
func (db *Database) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *SingleResult {
	if ctx == nil {
		ctx = context.Background()
	}

	op, sess, err := db.processRunCommand(ctx, runCommand, opts...)
	defer closeImplicitSession(sess)
	if err != nil {
		return &SingleResult{err: err}
	}

	err = op.Execute(ctx)
	return &SingleResult{
		err: replaceErrors(err),
		rdr: bson.Raw(op.Result()),
		reg: db.registry,
	}
}

// RunCommandCursor executes the given command against the database and parses the response as a cursor. If the command
// being executed does not return a cursor (e.g. insert), the command will be executed on the server and an error
// will be returned because the server response cannot be parsed as a cursor.
//
// The runCommand parameter must be a document for the command to be executed. It cannot be nil.
// This must be an order-preserving type such as bson.D. Map types such as bson.M are not valid.
// If the command document contains a session ID or any transaction-specific fields, the behavior is undefined.
//
// The opts parameter can be used to specify options for this operation (see the options.RunCmdOptions documentation).
func (db *Database) RunCommandCursor(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (*Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	op, sess, err := db.processRunCommand(ctx, runCommand, opts...)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}

	if err = op.Execute(ctx); err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}

	bc, err := op.ResultCursor(driver.CursorOptions{})
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}
	cursor, err := newCursorWithSession(bc, db.registry, sess)
	return cursor, replaceErrors(err)
}

// Drop drops the database on the server. This method ignores "namespace not found" errors so it is safe to drop
// a database that does not exist on the server.
func (db *Database) Drop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		var err error
		sess, err = session.NewClientSession(db.client.sessionPool, db.client.id, session.Implicit)
		if err != nil {
			return err
		}
		defer sess.EndSession()
	}

	err := db.client.validSession(sess)
	if err != nil {
		return err
	}

	wc := db.writeConcern
	if sess.TransactionRunning() {
		wc = nil
	}
	if !writeconcern.AckWrite(wc) {
		sess = nil
	}

	selector := makePinnedSelector(sess, db.writeSelector)

	op := operation.NewDropDatabase().
		Session(sess).WriteConcern(wc).CommandMonitor(db.client.monitor).
		ServerSelector(selector).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).Crypt(db.client.crypt)

	err = op.Execute(ctx)

	driverErr, ok := err.(driver.Error)
	if err != nil && (!ok || !driverErr.NamespaceNotFound()) {
		return replaceErrors(err)
	}
	return nil
}

// ListCollections executes a listCollections command and returns a cursor over the collections in the database.
//
// The filter parameter must be a document containing query operators and can be used to select which collections
// are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all
// collections.
//
// The opts parameter can be used to specify options for the operation (see the options.ListCollectionsOptions
// documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/listCollections/.
func (db *Database) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (*Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	filterDoc, err := transformBsoncoreDocument(db.registry, filter)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		sess, err = session.NewClientSession(db.client.sessionPool, db.client.id, session.Implicit)
		if err != nil {
			return nil, err
		}
	}

	err = db.client.validSession(sess)
	if err != nil {
		closeImplicitSession(sess)
		return nil, err
	}

	selector := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(readpref.Primary()),
		description.LatencySelector(db.client.localThreshold),
	})
	selector = makeReadPrefSelector(sess, selector, db.client.localThreshold)

	lco := options.MergeListCollectionsOptions(opts...)
	op := operation.NewListCollections(filterDoc).
		Session(sess).ReadPreference(db.readPreference).CommandMonitor(db.client.monitor).
		ServerSelector(selector).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).Crypt(db.client.crypt)
	if lco.NameOnly != nil {
		op = op.NameOnly(*lco.NameOnly)
	}
	retry := driver.RetryNone
	if db.client.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op = op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}

	bc, err := op.Result(driver.CursorOptions{Crypt: db.client.crypt})
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}
	cursor, err := newCursorWithSession(bc, db.registry, sess)
	return cursor, replaceErrors(err)
}

// ListCollectionNames executes a listCollections command and returns a slice containing the names of the collections
// in the database. This method requires driver version >= 1.1.0.
//
// The filter parameter must be a document containing query operators and can be used to select which collections
// are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all
// collections.
//
// The opts parameter can be used to specify options for the operation (see the options.ListCollectionsOptions
// documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/listCollections/.
func (db *Database) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	opts = append(opts, options.ListCollections().SetNameOnly(true))

	res, err := db.ListCollections(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)

	names := make([]string, 0)
	for res.Next(ctx) {
		next := &bsonx.Doc{}
		err = res.Decode(next)
		if err != nil {
			return nil, err
		}

		elem, err := next.LookupErr("name")
		if err != nil {
			return nil, err
		}

		if elem.Type() != bson.TypeString {
			return nil, fmt.Errorf("incorrect type for 'name'. got %v. want %v", elem.Type(), bson.TypeString)
		}

		elemName := elem.StringValue()
		names = append(names, elemName)
	}

	res.Close(ctx)
	return names, nil
}

// ReadConcern returns the read concern used to configure the Database object.
func (db *Database) ReadConcern() *readconcern.ReadConcern {
	return db.readConcern
}

// ReadPreference returns the read preference used to configure the Database object.
func (db *Database) ReadPreference() *readpref.ReadPref {
	return db.readPreference
}

// WriteConcern returns the write concern used to configure the Database object.
func (db *Database) WriteConcern() *writeconcern.WriteConcern {
	return db.writeConcern
}

// Watch returns a change stream for all changes to the corresponding database. See
// https://docs.mongodb.com/manual/changeStreams/ for more information about change streams.
//
// The Database must be configured with read concern majority or no read concern for a change stream to be created
// successfully.
//
// The pipeline parameter must be a slice of documents, each representing a pipeline stage. The pipeline cannot be
// nil but can be empty. The stage documents must all be non-nil. See https://docs.mongodb.com/manual/changeStreams/ for
// a list of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the
// mongo.Pipeline{} type can be used.
//
// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
// documentation).
func (db *Database) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) {

	csConfig := changeStreamConfig{
		readConcern:    db.readConcern,
		readPreference: db.readPreference,
		client:         db.client,
		registry:       db.registry,
		streamType:     DatabaseStream,
		databaseName:   db.Name(),
	}
	return newChangeStream(ctx, csConfig, pipeline, opts...)
}
