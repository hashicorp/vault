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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/internal/csfle"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
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
	bsonOpts       *options.BSONOptions
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

	bsonOpts := client.bsonOpts
	if dbOpt.BSONOptions != nil {
		bsonOpts = dbOpt.BSONOptions
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
		bsonOpts:       bsonOpts,
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
// https://www.mongodb.com/docs/manual/reference/operator/aggregation-pipeline/#db-aggregate-stages for a list of valid
// stages in database-level aggregations.
//
// The opts parameter can be used to specify options for this operation (see the options.AggregateOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/aggregate/.
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
	cursorCommand bool, opts ...*options.RunCmdOptions) (*operation.Command, *session.Client, error) {
	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(db.client.sessionPool, db.client.id)
	}

	err := db.client.validSession(sess)
	if err != nil {
		return nil, sess, err
	}

	ro := options.MergeRunCmdOptions(append(defaultRunCmdOpts, opts...)...)
	if sess != nil && sess.TransactionRunning() && ro.ReadPreference != nil && ro.ReadPreference.Mode() != readpref.PrimaryMode {
		return nil, sess, errors.New("read preference in a transaction must be primary")
	}

	if isUnorderedMap(cmd) {
		return nil, sess, ErrMapForOrderedArgument{"cmd"}
	}

	runCmdDoc, err := marshal(cmd, db.bsonOpts, db.registry)
	if err != nil {
		return nil, sess, err
	}
	readSelect := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(ro.ReadPreference),
		description.LatencySelector(db.client.localThreshold),
	})
	if sess != nil && sess.PinnedServer != nil {
		readSelect = makePinnedSelector(sess, readSelect)
	}

	var op *operation.Command
	switch cursorCommand {
	case true:
		cursorOpts := db.client.createBaseCursorOptions()

		cursorOpts.MarshalValueEncoderFn = newEncoderFn(db.bsonOpts, db.registry)

		op = operation.NewCursorCommand(runCmdDoc, cursorOpts)
	default:
		op = operation.NewCommand(runCmdDoc)
	}

	return op.Session(sess).CommandMonitor(db.client.monitor).
		ServerSelector(readSelect).ClusterClock(db.client.clock).
		Database(db.name).Deployment(db.client.deployment).
		Crypt(db.client.cryptFLE).ReadPreference(ro.ReadPreference).ServerAPI(db.client.serverAPI).
		Timeout(db.client.timeout).Logger(db.client.logger).Authenticator(db.client.authenticator), sess, nil
}

// RunCommand executes the given command against the database.
//
// This function does not obey the Database's readPreference. To specify a read
// preference, the RunCmdOptions.ReadPreference option must be used.
//
// This function does not obey the Database's readConcern or writeConcern. A
// user must supply these values manually in the user-provided runCommand
// parameter.
//
// The runCommand parameter must be a document for the command to be executed. It cannot be nil.
// This must be an order-preserving type such as bson.D. Map types such as bson.M are not valid.
//
// The opts parameter can be used to specify options for this operation (see the options.RunCmdOptions documentation).
//
// The behavior of RunCommand is undefined if the command document contains any of the following:
// - A session ID or any transaction-specific fields
// - API versioning options when an API version is already declared on the Client
// - maxTimeMS when Timeout is set on the Client
func (db *Database) RunCommand(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) *SingleResult {
	if ctx == nil {
		ctx = context.Background()
	}

	op, sess, err := db.processRunCommand(ctx, runCommand, false, opts...)
	defer closeImplicitSession(sess)
	if err != nil {
		return &SingleResult{err: err}
	}

	err = op.Execute(ctx)
	// RunCommand can be used to run a write, thus execute may return a write error
	_, convErr := processWriteError(err)
	return &SingleResult{
		ctx:      ctx,
		err:      convErr,
		rdr:      bson.Raw(op.Result()),
		bsonOpts: db.bsonOpts,
		reg:      db.registry,
	}
}

// RunCommandCursor executes the given command against the database and parses the response as a cursor. If the command
// being executed does not return a cursor (e.g. insert), the command will be executed on the server and an error will
// be returned because the server response cannot be parsed as a cursor. This function does not obey the Database's read
// preference. To specify a read preference, the RunCmdOptions.ReadPreference option must be used.
//
// The runCommand parameter must be a document for the command to be executed. It cannot be nil.
// This must be an order-preserving type such as bson.D. Map types such as bson.M are not valid.
//
// The opts parameter can be used to specify options for this operation (see the options.RunCmdOptions documentation).
//
// The behavior of RunCommandCursor is undefined if the command document contains any of the following:
// - A session ID or any transaction-specific fields
// - API versioning options when an API version is already declared on the Client
// - maxTimeMS when Timeout is set on the Client
func (db *Database) RunCommandCursor(ctx context.Context, runCommand interface{}, opts ...*options.RunCmdOptions) (*Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	op, sess, err := db.processRunCommand(ctx, runCommand, true, opts...)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}

	if err = op.Execute(ctx); err != nil {
		closeImplicitSession(sess)
		if errors.Is(err, driver.ErrNoCursor) {
			return nil, errors.New(
				"database response does not contain a cursor; try using RunCommand instead")
		}
		return nil, replaceErrors(err)
	}

	bc, err := op.ResultCursor()
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}
	cursor, err := newCursorWithSession(bc, db.bsonOpts, db.registry, sess)
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
		sess = session.NewImplicitClientSession(db.client.sessionPool, db.client.id)
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
		Database(db.name).Deployment(db.client.deployment).Crypt(db.client.cryptFLE).
		ServerAPI(db.client.serverAPI).Authenticator(db.client.authenticator)

	err = op.Execute(ctx)

	driverErr, ok := err.(driver.Error)
	if err != nil && (!ok || !driverErr.NamespaceNotFound()) {
		return replaceErrors(err)
	}
	return nil
}

// ListCollectionSpecifications executes a listCollections command and returns a slice of CollectionSpecification
// instances representing the collections in the database.
//
// The filter parameter must be a document containing query operators and can be used to select which collections
// are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all
// collections.
//
// The opts parameter can be used to specify options for the operation (see the options.ListCollectionsOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listCollections/.
//
// BUG(benjirewis): ListCollectionSpecifications prevents listing more than 100 collections per database when running
// against MongoDB version 2.6.
func (db *Database) ListCollectionSpecifications(ctx context.Context, filter interface{},
	opts ...*options.ListCollectionsOptions) ([]*CollectionSpecification, error) {

	cursor, err := db.ListCollections(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	var specs []*CollectionSpecification
	err = cursor.All(ctx, &specs)
	if err != nil {
		return nil, err
	}

	for _, spec := range specs {
		// Pre-4.4 servers report a namespace in their responses, so we only set Namespace manually if it was not in
		// the response.
		if spec.IDIndex != nil && spec.IDIndex.Namespace == "" {
			spec.IDIndex.Namespace = db.name + "." + spec.Name
		}
	}
	return specs, nil
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
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listCollections/.
//
// BUG(benjirewis): ListCollections prevents listing more than 100 collections per database when running against
// MongoDB version 2.6.
func (db *Database) ListCollections(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) (*Cursor, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	filterDoc, err := marshal(filter, db.bsonOpts, db.registry)
	if err != nil {
		return nil, err
	}

	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(db.client.sessionPool, db.client.id)
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
		Database(db.name).Deployment(db.client.deployment).Crypt(db.client.cryptFLE).
		ServerAPI(db.client.serverAPI).Timeout(db.client.timeout).Authenticator(db.client.authenticator)

	cursorOpts := db.client.createBaseCursorOptions()

	cursorOpts.MarshalValueEncoderFn = newEncoderFn(db.bsonOpts, db.registry)

	if lco.NameOnly != nil {
		op = op.NameOnly(*lco.NameOnly)
	}
	if lco.BatchSize != nil {
		cursorOpts.BatchSize = *lco.BatchSize
		op = op.BatchSize(*lco.BatchSize)
	}
	if lco.AuthorizedCollections != nil {
		op = op.AuthorizedCollections(*lco.AuthorizedCollections)
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

	bc, err := op.Result(cursorOpts)
	if err != nil {
		closeImplicitSession(sess)
		return nil, replaceErrors(err)
	}
	cursor, err := newCursorWithSession(bc, db.bsonOpts, db.registry, sess)
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
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listCollections/.
//
// BUG(benjirewis): ListCollectionNames prevents listing more than 100 collections per database when running against
// MongoDB version 2.6.
func (db *Database) ListCollectionNames(ctx context.Context, filter interface{}, opts ...*options.ListCollectionsOptions) ([]string, error) {
	opts = append(opts, options.ListCollections().SetNameOnly(true))

	res, err := db.ListCollections(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)

	names := make([]string, 0)
	for res.Next(ctx) {
		elem, err := res.Current.LookupErr("name")
		if err != nil {
			return nil, err
		}

		if elem.Type != bson.TypeString {
			return nil, fmt.Errorf("incorrect type for 'name'. got %v. want %v", elem.Type, bson.TypeString)
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
// https://www.mongodb.com/docs/manual/changeStreams/ for more information about change streams.
//
// The Database must be configured with read concern majority or no read concern for a change stream to be created
// successfully.
//
// The pipeline parameter must be a slice of documents, each representing a pipeline stage. The pipeline cannot be
// nil but can be empty. The stage documents must all be non-nil. See https://www.mongodb.com/docs/manual/changeStreams/ for
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
		crypt:          db.client.cryptFLE,
	}
	return newChangeStream(ctx, csConfig, pipeline, opts...)
}

// CreateCollection executes a create command to explicitly create a new collection with the specified name on the
// server. If the collection being created already exists, this method will return a mongo.CommandError. This method
// requires driver version 1.4.0 or higher.
//
// The opts parameter can be used to specify options for the operation (see the options.CreateCollectionOptions
// documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/create/.
func (db *Database) CreateCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	cco := options.MergeCreateCollectionOptions(opts...)
	// Follow Client-Side Encryption specification to check for encryptedFields.
	// Check for encryptedFields from create options.
	ef := cco.EncryptedFields
	// Check for encryptedFields from the client EncryptedFieldsMap.
	if ef == nil {
		ef = db.getEncryptedFieldsFromMap(name)
	}
	if ef != nil {
		return db.createCollectionWithEncryptedFields(ctx, name, ef, opts...)
	}

	return db.createCollection(ctx, name, opts...)
}

// getEncryptedFieldsFromServer tries to get an "encryptedFields" document associated with collectionName by running the "listCollections" command.
// Returns nil and no error if the listCollections command succeeds, but "encryptedFields" is not present.
func (db *Database) getEncryptedFieldsFromServer(ctx context.Context, collectionName string) (interface{}, error) {
	// Check if collection has an EncryptedFields configured server-side.
	collSpecs, err := db.ListCollectionSpecifications(ctx, bson.D{{"name", collectionName}})
	if err != nil {
		return nil, err
	}
	if len(collSpecs) == 0 {
		return nil, nil
	}
	if len(collSpecs) > 1 {
		return nil, fmt.Errorf("expected 1 or 0 results from listCollections, got %v", len(collSpecs))
	}
	collSpec := collSpecs[0]
	rawValue, err := collSpec.Options.LookupErr("encryptedFields")
	if errors.Is(err, bsoncore.ErrElementNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	encryptedFields, ok := rawValue.DocumentOK()
	if !ok {
		return nil, fmt.Errorf("expected encryptedFields of %v to be document, got %v", collectionName, rawValue.Type)
	}

	return encryptedFields, nil
}

// getEncryptedFieldsFromMap tries to get an "encryptedFields" document associated with collectionName by checking the client EncryptedFieldsMap.
// Returns nil and no error if an EncryptedFieldsMap is not configured, or does not contain an entry for collectionName.
func (db *Database) getEncryptedFieldsFromMap(collectionName string) interface{} {
	// Check the EncryptedFieldsMap
	efMap := db.client.encryptedFieldsMap
	if efMap == nil {
		return nil
	}

	namespace := db.name + "." + collectionName

	ef, ok := efMap[namespace]
	if ok {
		return ef
	}
	return nil
}

// createCollectionWithEncryptedFields creates a collection with an EncryptedFields.
func (db *Database) createCollectionWithEncryptedFields(ctx context.Context, name string, ef interface{}, opts ...*options.CreateCollectionOptions) error {
	efBSON, err := marshal(ef, db.bsonOpts, db.registry)
	if err != nil {
		return fmt.Errorf("error transforming document: %w", err)
	}

	// Check the wire version to ensure server is 7.0.0 or newer.
	// After the wire version check, and before creating the collections, it is possible the server state changes.
	// That is OK. This wire version check is a best effort to inform users earlier if using a QEv2 driver with a QEv1 server.
	{
		const QEv2WireVersion = 21
		server, err := db.client.deployment.SelectServer(ctx, description.WriteSelector())
		if err != nil {
			return fmt.Errorf("error selecting server to check maxWireVersion: %w", err)
		}
		conn, err := server.Connection(ctx)
		if err != nil {
			return fmt.Errorf("error getting connection to check maxWireVersion: %w", err)
		}
		defer conn.Close()
		wireVersionRange := conn.Description().WireVersion
		if wireVersionRange.Max < QEv2WireVersion {
			return fmt.Errorf("Driver support of Queryable Encryption is incompatible with server. Upgrade server to use Queryable Encryption. Got maxWireVersion %v but need maxWireVersion >= %v", wireVersionRange.Max, QEv2WireVersion)
		}
	}

	// Create the two encryption-related, associated collections: `escCollection` and `ecocCollection`.

	stateCollectionOpts := options.CreateCollection().
		SetClusteredIndex(bson.D{{"key", bson.D{{"_id", 1}}}, {"unique", true}})
	// Create ESCCollection.
	escCollection, err := csfle.GetEncryptedStateCollectionName(efBSON, name, csfle.EncryptedStateCollection)
	if err != nil {
		return err
	}

	if err := db.createCollection(ctx, escCollection, stateCollectionOpts); err != nil {
		return err
	}

	// Create ECOCCollection.
	ecocCollection, err := csfle.GetEncryptedStateCollectionName(efBSON, name, csfle.EncryptedCompactionCollection)
	if err != nil {
		return err
	}

	if err := db.createCollection(ctx, ecocCollection, stateCollectionOpts); err != nil {
		return err
	}

	// Create a data collection with the 'encryptedFields' option.
	op, err := db.createCollectionOperation(name, opts...)
	if err != nil {
		return err
	}

	op.EncryptedFields(efBSON)
	if err := db.executeCreateOperation(ctx, op); err != nil {
		return err
	}

	// Create an index on the __safeContent__ field in the collection @collectionName.
	if _, err := db.Collection(name).Indexes().CreateOne(ctx, IndexModel{Keys: bson.D{{"__safeContent__", 1}}}); err != nil {
		return fmt.Errorf("error creating safeContent index: %w", err)
	}

	return nil
}

// createCollection creates a collection without EncryptedFields.
func (db *Database) createCollection(ctx context.Context, name string, opts ...*options.CreateCollectionOptions) error {
	op, err := db.createCollectionOperation(name, opts...)
	if err != nil {
		return err
	}
	return db.executeCreateOperation(ctx, op)
}

func (db *Database) createCollectionOperation(name string, opts ...*options.CreateCollectionOptions) (*operation.Create, error) {
	cco := options.MergeCreateCollectionOptions(opts...)
	op := operation.NewCreate(name).ServerAPI(db.client.serverAPI).Authenticator(db.client.authenticator)

	if cco.Capped != nil {
		op.Capped(*cco.Capped)
	}
	if cco.Collation != nil {
		op.Collation(bsoncore.Document(cco.Collation.ToDocument()))
	}
	if cco.ChangeStreamPreAndPostImages != nil {
		csppi, err := marshal(cco.ChangeStreamPreAndPostImages, db.bsonOpts, db.registry)
		if err != nil {
			return nil, err
		}
		op.ChangeStreamPreAndPostImages(csppi)
	}
	if cco.DefaultIndexOptions != nil {
		idx, doc := bsoncore.AppendDocumentStart(nil)
		if cco.DefaultIndexOptions.StorageEngine != nil {
			storageEngine, err := marshal(cco.DefaultIndexOptions.StorageEngine, db.bsonOpts, db.registry)
			if err != nil {
				return nil, err
			}

			doc = bsoncore.AppendDocumentElement(doc, "storageEngine", storageEngine)
		}
		doc, err := bsoncore.AppendDocumentEnd(doc, idx)
		if err != nil {
			return nil, err
		}

		op.IndexOptionDefaults(doc)
	}
	if cco.MaxDocuments != nil {
		op.Max(*cco.MaxDocuments)
	}
	if cco.SizeInBytes != nil {
		op.Size(*cco.SizeInBytes)
	}
	if cco.StorageEngine != nil {
		storageEngine, err := marshal(cco.StorageEngine, db.bsonOpts, db.registry)
		if err != nil {
			return nil, err
		}
		op.StorageEngine(storageEngine)
	}
	if cco.ValidationAction != nil {
		op.ValidationAction(*cco.ValidationAction)
	}
	if cco.ValidationLevel != nil {
		op.ValidationLevel(*cco.ValidationLevel)
	}
	if cco.Validator != nil {
		validator, err := marshal(cco.Validator, db.bsonOpts, db.registry)
		if err != nil {
			return nil, err
		}
		op.Validator(validator)
	}
	if cco.ExpireAfterSeconds != nil {
		op.ExpireAfterSeconds(*cco.ExpireAfterSeconds)
	}
	if cco.TimeSeriesOptions != nil {
		idx, doc := bsoncore.AppendDocumentStart(nil)
		doc = bsoncore.AppendStringElement(doc, "timeField", cco.TimeSeriesOptions.TimeField)

		if cco.TimeSeriesOptions.MetaField != nil {
			doc = bsoncore.AppendStringElement(doc, "metaField", *cco.TimeSeriesOptions.MetaField)
		}
		if cco.TimeSeriesOptions.Granularity != nil {
			doc = bsoncore.AppendStringElement(doc, "granularity", *cco.TimeSeriesOptions.Granularity)
		}

		if cco.TimeSeriesOptions.BucketMaxSpan != nil {
			bmss := int64(*cco.TimeSeriesOptions.BucketMaxSpan / time.Second)

			doc = bsoncore.AppendInt64Element(doc, "bucketMaxSpanSeconds", bmss)
		}

		if cco.TimeSeriesOptions.BucketRounding != nil {
			brs := int64(*cco.TimeSeriesOptions.BucketRounding / time.Second)

			doc = bsoncore.AppendInt64Element(doc, "bucketRoundingSeconds", brs)
		}

		doc, err := bsoncore.AppendDocumentEnd(doc, idx)
		if err != nil {
			return nil, err
		}

		op.TimeSeries(doc)
	}
	if cco.ClusteredIndex != nil {
		clusteredIndex, err := marshal(cco.ClusteredIndex, db.bsonOpts, db.registry)
		if err != nil {
			return nil, err
		}
		op.ClusteredIndex(clusteredIndex)
	}

	return op, nil
}

// CreateView executes a create command to explicitly create a view on the server. See
// https://www.mongodb.com/docs/manual/core/views/ for more information about views. This method requires driver version >=
// 1.4.0 and MongoDB version >= 3.4.
//
// The viewName parameter specifies the name of the view to create.
//
// # The viewOn parameter specifies the name of the collection or view on which this view will be created
//
// The pipeline parameter specifies an aggregation pipeline that will be exececuted against the source collection or
// view to create this view.
//
// The opts parameter can be used to specify options for the operation (see the options.CreateViewOptions
// documentation).
func (db *Database) CreateView(ctx context.Context, viewName, viewOn string, pipeline interface{},
	opts ...*options.CreateViewOptions) error {

	pipelineArray, _, err := marshalAggregatePipeline(pipeline, db.bsonOpts, db.registry)
	if err != nil {
		return err
	}

	op := operation.NewCreate(viewName).
		ViewOn(viewOn).
		Pipeline(pipelineArray).
		ServerAPI(db.client.serverAPI).
		Authenticator(db.client.authenticator)
	cvo := options.MergeCreateViewOptions(opts...)
	if cvo.Collation != nil {
		op.Collation(bsoncore.Document(cvo.Collation.ToDocument()))
	}

	return db.executeCreateOperation(ctx, op)
}

func (db *Database) executeCreateOperation(ctx context.Context, op *operation.Create) error {
	sess := sessionFromContext(ctx)
	if sess == nil && db.client.sessionPool != nil {
		sess = session.NewImplicitClientSession(db.client.sessionPool, db.client.id)
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
	op = op.Session(sess).
		WriteConcern(wc).
		CommandMonitor(db.client.monitor).
		ServerSelector(selector).
		ClusterClock(db.client.clock).
		Database(db.name).
		Deployment(db.client.deployment).
		Crypt(db.client.cryptFLE)

	return replaceErrors(op.Execute(ctx))
}
