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
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/internal/httputil"
	"go.mongodb.org/mongo-driver/internal/logger"
	"go.mongodb.org/mongo-driver/internal/uuid"
	"go.mongodb.org/mongo-driver/mongo/description"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	mcopts "go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
)

const (
	defaultLocalThreshold = 15 * time.Millisecond
	defaultMaxPoolSize    = 100
)

var (
	// keyVaultCollOpts specifies options used to communicate with the key vault collection
	keyVaultCollOpts = options.Collection().SetReadConcern(readconcern.Majority()).
				SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

	endSessionsBatchSize = 10000
)

// Client is a handle representing a pool of connections to a MongoDB deployment. It is safe for concurrent use by
// multiple goroutines.
//
// The Client type opens and closes connections automatically and maintains a pool of idle connections. For
// connection pool configuration options, see documentation for the ClientOptions type in the mongo/options package.
type Client struct {
	id             uuid.UUID
	deployment     driver.Deployment
	localThreshold time.Duration
	retryWrites    bool
	retryReads     bool
	clock          *session.ClusterClock
	readPreference *readpref.ReadPref
	readConcern    *readconcern.ReadConcern
	writeConcern   *writeconcern.WriteConcern
	bsonOpts       *options.BSONOptions
	registry       *bsoncodec.Registry
	monitor        *event.CommandMonitor
	serverAPI      *driver.ServerAPIOptions
	serverMonitor  *event.ServerMonitor
	sessionPool    *session.Pool
	timeout        *time.Duration
	httpClient     *http.Client
	logger         *logger.Logger

	// client-side encryption fields
	keyVaultClientFLE  *Client
	keyVaultCollFLE    *Collection
	mongocryptdFLE     *mongocryptdClient
	cryptFLE           driver.Crypt
	metadataClientFLE  *Client
	internalClientFLE  *Client
	encryptedFieldsMap map[string]interface{}
	authenticator      driver.Authenticator
}

// Connect creates a new Client and then initializes it using the Connect method. This is equivalent to calling
// NewClient followed by Client.Connect.
//
// When creating an options.ClientOptions, the order the methods are called matters. Later Set*
// methods will overwrite the values from previous Set* method invocations. This includes the
// ApplyURI method. This allows callers to determine the order of precedence for option
// application. For instance, if ApplyURI is called before SetAuth, the Credential from
// SetAuth will overwrite the values from the connection string. If ApplyURI is called
// after SetAuth, then its values will overwrite those from SetAuth.
//
// The opts parameter is processed using options.MergeClientOptions, which will overwrite entire
// option fields of previous options, there is no partial overwriting. For example, if Username is
// set in the Auth field for the first option, and Password is set for the second but with no
// Username, after the merge the Username field will be empty.
//
// The NewClient function does not do any I/O and returns an error if the given options are invalid.
// The Client.Connect method starts background goroutines to monitor the state of the deployment and does not do
// any I/O in the main goroutine to prevent the main goroutine from blocking. Therefore, it will not error if the
// deployment is down.
//
// The Client.Ping method can be used to verify that the deployment is successfully connected and the
// Client was correctly configured.
func Connect(ctx context.Context, opts ...*options.ClientOptions) (*Client, error) {
	c, err := NewClient(opts...)
	if err != nil {
		return nil, err
	}
	err = c.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// NewClient creates a new client to connect to a deployment specified by the uri.
//
// When creating an options.ClientOptions, the order the methods are called matters. Later Set*
// methods will overwrite the values from previous Set* method invocations. This includes the
// ApplyURI method. This allows callers to determine the order of precedence for option
// application. For instance, if ApplyURI is called before SetAuth, the Credential from
// SetAuth will overwrite the values from the connection string. If ApplyURI is called
// after SetAuth, then its values will overwrite those from SetAuth.
//
// The opts parameter is processed using options.MergeClientOptions, which will overwrite entire
// option fields of previous options, there is no partial overwriting. For example, if Username is
// set in the Auth field for the first option, and Password is set for the second but with no
// Username, after the merge the Username field will be empty.
//
// Deprecated: Use [Connect] instead.
func NewClient(opts ...*options.ClientOptions) (*Client, error) {
	clientOpt := options.MergeClientOptions(opts...)

	id, err := uuid.New()
	if err != nil {
		return nil, err
	}
	client := &Client{id: id}

	// ClusterClock
	client.clock = new(session.ClusterClock)

	// LocalThreshold
	client.localThreshold = defaultLocalThreshold
	if clientOpt.LocalThreshold != nil {
		client.localThreshold = *clientOpt.LocalThreshold
	}
	// Monitor
	if clientOpt.Monitor != nil {
		client.monitor = clientOpt.Monitor
	}
	// ServerMonitor
	if clientOpt.ServerMonitor != nil {
		client.serverMonitor = clientOpt.ServerMonitor
	}
	// ReadConcern
	client.readConcern = readconcern.New()
	if clientOpt.ReadConcern != nil {
		client.readConcern = clientOpt.ReadConcern
	}
	// ReadPreference
	client.readPreference = readpref.Primary()
	if clientOpt.ReadPreference != nil {
		client.readPreference = clientOpt.ReadPreference
	}
	// BSONOptions
	if clientOpt.BSONOptions != nil {
		client.bsonOpts = clientOpt.BSONOptions
	}
	// Registry
	client.registry = bson.DefaultRegistry
	if clientOpt.Registry != nil {
		client.registry = clientOpt.Registry
	}
	// RetryWrites
	client.retryWrites = true // retry writes on by default
	if clientOpt.RetryWrites != nil {
		client.retryWrites = *clientOpt.RetryWrites
	}
	client.retryReads = true
	if clientOpt.RetryReads != nil {
		client.retryReads = *clientOpt.RetryReads
	}
	// Timeout
	client.timeout = clientOpt.Timeout
	client.httpClient = clientOpt.HTTPClient
	// WriteConcern
	if clientOpt.WriteConcern != nil {
		client.writeConcern = clientOpt.WriteConcern
	}
	// AutoEncryptionOptions
	if clientOpt.AutoEncryptionOptions != nil {
		if err := client.configureAutoEncryption(clientOpt); err != nil {
			return nil, err
		}
	} else {
		client.cryptFLE = clientOpt.Crypt
	}

	// Deployment
	if clientOpt.Deployment != nil {
		client.deployment = clientOpt.Deployment
	}

	// Set default options
	if clientOpt.MaxPoolSize == nil {
		clientOpt.SetMaxPoolSize(defaultMaxPoolSize)
	}

	if clientOpt.Auth != nil {
		client.authenticator, err = auth.CreateAuthenticator(
			clientOpt.Auth.AuthMechanism,
			topology.ConvertCreds(clientOpt.Auth),
			clientOpt.HTTPClient,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating authenticator: %w", err)
		}
	}

	cfg, err := topology.NewConfigWithAuthenticator(clientOpt, client.clock, client.authenticator)
	if err != nil {
		return nil, err
	}

	client.serverAPI = topology.ServerAPIFromServerOptions(cfg.ServerOpts)

	if client.deployment == nil {
		client.deployment, err = topology.New(cfg)
		if err != nil {
			return nil, replaceErrors(err)
		}
	}

	// Create a logger for the client.
	client.logger, err = newLogger(clientOpt.LoggerOptions)
	if err != nil {
		return nil, fmt.Errorf("invalid logger options: %w", err)
	}

	return client, nil
}

// Connect initializes the Client by starting background monitoring goroutines.
// If the Client was created using the NewClient function, this method must be called before a Client can be used.
//
// Connect starts background goroutines to monitor the state of the deployment and does not do any I/O in the main
// goroutine. The Client.Ping method can be used to verify that the connection was created successfully.
//
// Deprecated: Use [mongo.Connect] instead.
func (c *Client) Connect(ctx context.Context) error {
	if connector, ok := c.deployment.(driver.Connector); ok {
		err := connector.Connect()
		if err != nil {
			return replaceErrors(err)
		}
	}

	if c.mongocryptdFLE != nil {
		if err := c.mongocryptdFLE.connect(ctx); err != nil {
			return err
		}
	}

	if c.internalClientFLE != nil {
		if err := c.internalClientFLE.Connect(ctx); err != nil {
			return err
		}
	}

	if c.keyVaultClientFLE != nil && c.keyVaultClientFLE != c.internalClientFLE && c.keyVaultClientFLE != c {
		if err := c.keyVaultClientFLE.Connect(ctx); err != nil {
			return err
		}
	}

	if c.metadataClientFLE != nil && c.metadataClientFLE != c.internalClientFLE && c.metadataClientFLE != c {
		if err := c.metadataClientFLE.Connect(ctx); err != nil {
			return err
		}
	}

	var updateChan <-chan description.Topology
	if subscriber, ok := c.deployment.(driver.Subscriber); ok {
		sub, err := subscriber.Subscribe()
		if err != nil {
			return replaceErrors(err)
		}
		updateChan = sub.Updates
	}
	c.sessionPool = session.NewPool(updateChan)
	return nil
}

// Disconnect closes sockets to the topology referenced by this Client. It will
// shut down any monitoring goroutines, close the idle connection pool, and will
// wait until all the in use connections have been returned to the connection
// pool and closed before returning. If the context expires via cancellation,
// deadline, or timeout before the in use connections have returned, the in use
// connections will be closed, resulting in the failure of any in flight read
// or write operations. If this method returns with no errors, all connections
// associated with this Client have been closed.
func (c *Client) Disconnect(ctx context.Context) error {
	if c.logger != nil {
		defer c.logger.Close()
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if c.httpClient == httputil.DefaultHTTPClient {
		defer httputil.CloseIdleHTTPConnections(c.httpClient)
	}

	c.endSessions(ctx)
	if c.mongocryptdFLE != nil {
		if err := c.mongocryptdFLE.disconnect(ctx); err != nil {
			return err
		}
	}

	if c.internalClientFLE != nil {
		if err := c.internalClientFLE.Disconnect(ctx); err != nil {
			return err
		}
	}

	if c.keyVaultClientFLE != nil && c.keyVaultClientFLE != c.internalClientFLE && c.keyVaultClientFLE != c {
		if err := c.keyVaultClientFLE.Disconnect(ctx); err != nil {
			return err
		}
	}
	if c.metadataClientFLE != nil && c.metadataClientFLE != c.internalClientFLE && c.metadataClientFLE != c {
		if err := c.metadataClientFLE.Disconnect(ctx); err != nil {
			return err
		}
	}
	if c.cryptFLE != nil {
		c.cryptFLE.Close()
	}

	if disconnector, ok := c.deployment.(driver.Disconnector); ok {
		return replaceErrors(disconnector.Disconnect(ctx))
	}

	return nil
}

// Ping sends a ping command to verify that the client can connect to the deployment.
//
// The rp parameter is used to determine which server is selected for the operation.
// If it is nil, the client's read preference is used.
//
// If the server is down, Ping will try to select a server until the client's server selection timeout expires.
// This can be configured through the ClientOptions.SetServerSelectionTimeout option when creating a new Client.
// After the timeout expires, a server selection error is returned.
//
// Using Ping reduces application resilience because applications starting up will error if the server is temporarily
// unavailable or is failing over (e.g. during autoscaling due to a load spike).
func (c *Client) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if rp == nil {
		rp = c.readPreference
	}

	db := c.Database("admin")
	res := db.RunCommand(ctx, bson.D{
		{"ping", 1},
	}, options.RunCmd().SetReadPreference(rp))

	return replaceErrors(res.Err())
}

// StartSession starts a new session configured with the given options.
//
// StartSession does not actually communicate with the server and will not error if the client is
// disconnected.
//
// StartSession is safe to call from multiple goroutines concurrently. However, Sessions returned by StartSession are
// not safe for concurrent use by multiple goroutines.
//
// If the DefaultReadConcern, DefaultWriteConcern, or DefaultReadPreference options are not set, the client's read
// concern, write concern, or read preference will be used, respectively.
func (c *Client) StartSession(opts ...*options.SessionOptions) (Session, error) {
	if c.sessionPool == nil {
		return nil, ErrClientDisconnected
	}

	sopts := options.MergeSessionOptions(opts...)
	coreOpts := &session.ClientOptions{
		DefaultReadConcern:    c.readConcern,
		DefaultReadPreference: c.readPreference,
		DefaultWriteConcern:   c.writeConcern,
	}
	if sopts.CausalConsistency != nil {
		coreOpts.CausalConsistency = sopts.CausalConsistency
	}
	if sopts.DefaultReadConcern != nil {
		coreOpts.DefaultReadConcern = sopts.DefaultReadConcern
	}
	if sopts.DefaultWriteConcern != nil {
		coreOpts.DefaultWriteConcern = sopts.DefaultWriteConcern
	}
	if sopts.DefaultReadPreference != nil {
		coreOpts.DefaultReadPreference = sopts.DefaultReadPreference
	}
	if sopts.DefaultMaxCommitTime != nil {
		coreOpts.DefaultMaxCommitTime = sopts.DefaultMaxCommitTime
	}
	if sopts.Snapshot != nil {
		coreOpts.Snapshot = sopts.Snapshot
	}

	sess, err := session.NewClientSession(c.sessionPool, c.id, coreOpts)
	if err != nil {
		return nil, replaceErrors(err)
	}

	// Writes are not retryable on standalones, so let operation determine whether to retry
	sess.RetryWrite = false
	sess.RetryRead = c.retryReads

	return &sessionImpl{
		clientSession: sess,
		client:        c,
		deployment:    c.deployment,
	}, nil
}

func (c *Client) endSessions(ctx context.Context) {
	if c.sessionPool == nil {
		return
	}

	sessionIDs := c.sessionPool.IDSlice()
	op := operation.NewEndSessions(nil).ClusterClock(c.clock).Deployment(c.deployment).
		ServerSelector(description.ReadPrefSelector(readpref.PrimaryPreferred())).CommandMonitor(c.monitor).
		Database("admin").Crypt(c.cryptFLE).ServerAPI(c.serverAPI)

	totalNumIDs := len(sessionIDs)
	var currentBatch []bsoncore.Document
	for i := 0; i < totalNumIDs; i++ {
		currentBatch = append(currentBatch, sessionIDs[i])

		// If we are at the end of a batch or the end of the overall IDs array, execute the operation.
		if ((i+1)%endSessionsBatchSize) == 0 || i == totalNumIDs-1 {
			// Ignore all errors when ending sessions.
			_, marshalVal, err := bson.MarshalValue(currentBatch)
			if err == nil {
				_ = op.SessionIDs(marshalVal).Execute(ctx)
			}

			currentBatch = currentBatch[:0]
		}
	}
}

func (c *Client) configureAutoEncryption(clientOpts *options.ClientOptions) error {
	c.encryptedFieldsMap = clientOpts.AutoEncryptionOptions.EncryptedFieldsMap
	if err := c.configureKeyVaultClientFLE(clientOpts); err != nil {
		return err
	}
	if err := c.configureMetadataClientFLE(clientOpts); err != nil {
		return err
	}

	mc, err := c.newMongoCrypt(clientOpts.AutoEncryptionOptions)
	if err != nil {
		return err
	}

	// If the crypt_shared library was not loaded, try to spawn and connect to mongocryptd.
	if mc.CryptSharedLibVersionString() == "" {
		mongocryptdFLE, err := newMongocryptdClient(clientOpts.AutoEncryptionOptions)
		if err != nil {
			return err
		}
		c.mongocryptdFLE = mongocryptdFLE
	}

	c.configureCryptFLE(mc, clientOpts.AutoEncryptionOptions)
	return nil
}

func (c *Client) getOrCreateInternalClient(clientOpts *options.ClientOptions) (*Client, error) {
	if c.internalClientFLE != nil {
		return c.internalClientFLE, nil
	}

	internalClientOpts := options.MergeClientOptions(clientOpts)
	internalClientOpts.AutoEncryptionOptions = nil
	internalClientOpts.SetMinPoolSize(0)
	var err error
	c.internalClientFLE, err = NewClient(internalClientOpts)
	return c.internalClientFLE, err
}

func (c *Client) configureKeyVaultClientFLE(clientOpts *options.ClientOptions) error {
	// parse key vault options and create new key vault client
	var err error
	aeOpts := clientOpts.AutoEncryptionOptions
	switch {
	case aeOpts.KeyVaultClientOptions != nil:
		c.keyVaultClientFLE, err = NewClient(aeOpts.KeyVaultClientOptions)
	case clientOpts.MaxPoolSize != nil && *clientOpts.MaxPoolSize == 0:
		c.keyVaultClientFLE = c
	default:
		c.keyVaultClientFLE, err = c.getOrCreateInternalClient(clientOpts)
	}

	if err != nil {
		return err
	}

	dbName, collName := splitNamespace(aeOpts.KeyVaultNamespace)
	c.keyVaultCollFLE = c.keyVaultClientFLE.Database(dbName).Collection(collName, keyVaultCollOpts)
	return nil
}

func (c *Client) configureMetadataClientFLE(clientOpts *options.ClientOptions) error {
	// parse key vault options and create new key vault client
	aeOpts := clientOpts.AutoEncryptionOptions
	if aeOpts.BypassAutoEncryption != nil && *aeOpts.BypassAutoEncryption {
		// no need for a metadata client.
		return nil
	}
	if clientOpts.MaxPoolSize != nil && *clientOpts.MaxPoolSize == 0 {
		c.metadataClientFLE = c
		return nil
	}

	var err error
	c.metadataClientFLE, err = c.getOrCreateInternalClient(clientOpts)
	return err
}

func (c *Client) newMongoCrypt(opts *options.AutoEncryptionOptions) (*mongocrypt.MongoCrypt, error) {
	// convert schemas in SchemaMap to bsoncore documents
	cryptSchemaMap := make(map[string]bsoncore.Document)
	for k, v := range opts.SchemaMap {
		schema, err := marshal(v, c.bsonOpts, c.registry)
		if err != nil {
			return nil, err
		}
		cryptSchemaMap[k] = schema
	}

	// convert schemas in EncryptedFieldsMap to bsoncore documents
	cryptEncryptedFieldsMap := make(map[string]bsoncore.Document)
	for k, v := range opts.EncryptedFieldsMap {
		encryptedFields, err := marshal(v, c.bsonOpts, c.registry)
		if err != nil {
			return nil, err
		}
		cryptEncryptedFieldsMap[k] = encryptedFields
	}

	kmsProviders, err := marshal(opts.KmsProviders, c.bsonOpts, c.registry)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS providers document: %w", err)
	}

	// Set the crypt_shared library override path from the "cryptSharedLibPath" extra option if one
	// was set.
	cryptSharedLibPath := ""
	if val, ok := opts.ExtraOptions["cryptSharedLibPath"]; ok {
		str, ok := val.(string)
		if !ok {
			return nil, fmt.Errorf(
				`expected AutoEncryption extra option "cryptSharedLibPath" to be a string, but is a %T`, val)
		}
		cryptSharedLibPath = str
	}

	// Explicitly disable loading the crypt_shared library if requested. Note that this is ONLY
	// intended for use from tests; there is no supported public API for explicitly disabling
	// loading the crypt_shared library.
	cryptSharedLibDisabled := false
	if v, ok := opts.ExtraOptions["__cryptSharedLibDisabledForTestOnly"]; ok {
		cryptSharedLibDisabled = v.(bool)
	}

	bypassAutoEncryption := opts.BypassAutoEncryption != nil && *opts.BypassAutoEncryption
	bypassQueryAnalysis := opts.BypassQueryAnalysis != nil && *opts.BypassQueryAnalysis

	mc, err := mongocrypt.NewMongoCrypt(mcopts.MongoCrypt().
		SetKmsProviders(kmsProviders).
		SetLocalSchemaMap(cryptSchemaMap).
		SetBypassQueryAnalysis(bypassQueryAnalysis).
		SetEncryptedFieldsMap(cryptEncryptedFieldsMap).
		SetCryptSharedLibDisabled(cryptSharedLibDisabled || bypassAutoEncryption).
		SetCryptSharedLibOverridePath(cryptSharedLibPath).
		SetHTTPClient(opts.HTTPClient))
	if err != nil {
		return nil, err
	}

	var cryptSharedLibRequired bool
	if val, ok := opts.ExtraOptions["cryptSharedLibRequired"]; ok {
		b, ok := val.(bool)
		if !ok {
			return nil, fmt.Errorf(
				`expected AutoEncryption extra option "cryptSharedLibRequired" to be a bool, but is a %T`, val)
		}
		cryptSharedLibRequired = b
	}

	// If the "cryptSharedLibRequired" extra option is set to true, check the MongoCrypt version
	// string to confirm that the library was successfully loaded. If the version string is empty,
	// return an error indicating that we couldn't load the crypt_shared library.
	if cryptSharedLibRequired && mc.CryptSharedLibVersionString() == "" {
		return nil, errors.New(
			`AutoEncryption extra option "cryptSharedLibRequired" is true, but we failed to load the crypt_shared library`)
	}

	return mc, nil
}

//nolint:unused // the unused linter thinks that this function is unreachable because "c.newMongoCrypt" always panics without the "cse" build tag set.
func (c *Client) configureCryptFLE(mc *mongocrypt.MongoCrypt, opts *options.AutoEncryptionOptions) {
	bypass := opts.BypassAutoEncryption != nil && *opts.BypassAutoEncryption
	kr := keyRetriever{coll: c.keyVaultCollFLE}
	var cir collInfoRetriever
	// If bypass is true, c.metadataClientFLE is nil and the collInfoRetriever
	// will not be used. If bypass is false, to the parent client or the
	// internal client.
	if !bypass {
		cir = collInfoRetriever{client: c.metadataClientFLE}
	}

	c.cryptFLE = driver.NewCrypt(&driver.CryptOptions{
		MongoCrypt:           mc,
		CollInfoFn:           cir.cryptCollInfo,
		KeyFn:                kr.cryptKeys,
		MarkFn:               c.mongocryptdFLE.markCommand,
		TLSConfig:            opts.TLSConfig,
		BypassAutoEncryption: bypass,
	})
}

// validSession returns an error if the session doesn't belong to the client
func (c *Client) validSession(sess *session.Client) error {
	if sess != nil && sess.ClientID != c.id {
		return ErrWrongClient
	}
	return nil
}

// Database returns a handle for a database with the given name configured with the given DatabaseOptions.
func (c *Client) Database(name string, opts ...*options.DatabaseOptions) *Database {
	return newDatabase(c, name, opts...)
}

// ListDatabases executes a listDatabases command and returns the result.
//
// The filter parameter must be a document containing query operators and can be used to select which
// databases are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include
// all databases.
//
// The opts parameter can be used to specify options for this operation (see the options.ListDatabasesOptions documentation).
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listDatabases/.
func (c *Client) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (ListDatabasesResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)

	err := c.validSession(sess)
	if err != nil {
		return ListDatabasesResult{}, err
	}
	if sess == nil && c.sessionPool != nil {
		sess = session.NewImplicitClientSession(c.sessionPool, c.id)
		defer sess.EndSession()
	}

	err = c.validSession(sess)
	if err != nil {
		return ListDatabasesResult{}, err
	}

	filterDoc, err := marshal(filter, c.bsonOpts, c.registry)
	if err != nil {
		return ListDatabasesResult{}, err
	}

	selector := description.CompositeSelector([]description.ServerSelector{
		description.ReadPrefSelector(readpref.Primary()),
		description.LatencySelector(c.localThreshold),
	})
	selector = makeReadPrefSelector(sess, selector, c.localThreshold)

	ldo := options.MergeListDatabasesOptions(opts...)
	op := operation.NewListDatabases(filterDoc).
		Session(sess).ReadPreference(c.readPreference).CommandMonitor(c.monitor).
		ServerSelector(selector).ClusterClock(c.clock).Database("admin").Deployment(c.deployment).Crypt(c.cryptFLE).
		ServerAPI(c.serverAPI).Timeout(c.timeout).Authenticator(c.authenticator)

	if ldo.NameOnly != nil {
		op = op.NameOnly(*ldo.NameOnly)
	}
	if ldo.AuthorizedDatabases != nil {
		op = op.AuthorizedDatabases(*ldo.AuthorizedDatabases)
	}

	retry := driver.RetryNone
	if c.retryReads {
		retry = driver.RetryOncePerCommand
	}
	op.Retry(retry)

	err = op.Execute(ctx)
	if err != nil {
		return ListDatabasesResult{}, replaceErrors(err)
	}

	return newListDatabasesResultFromOperation(op.Result()), nil
}

// ListDatabaseNames executes a listDatabases command and returns a slice containing the names of all of the databases
// on the server.
//
// The filter parameter must be a document containing query operators and can be used to select which databases
// are included in the result. It cannot be nil. An empty document (e.g. bson.D{}) should be used to include all
// databases.
//
// The opts parameter can be used to specify options for this operation (see the options.ListDatabasesOptions
// documentation.)
//
// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/listDatabases/.
func (c *Client) ListDatabaseNames(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) ([]string, error) {
	opts = append(opts, options.ListDatabases().SetNameOnly(true))

	res, err := c.ListDatabases(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	for _, spec := range res.Databases {
		names = append(names, spec.Name)
	}

	return names, nil
}

// WithSession creates a new SessionContext from the ctx and sess parameters and uses it to call the fn callback. The
// SessionContext must be used as the Context parameter for any operations in the fn callback that should be executed
// under the session.
//
// WithSession is safe to call from multiple goroutines concurrently. However, the SessionContext passed to the
// WithSession callback function is not safe for concurrent use by multiple goroutines.
//
// If the ctx parameter already contains a Session, that Session will be replaced with the one provided.
//
// Any error returned by the fn callback will be returned without any modifications.
func WithSession(ctx context.Context, sess Session, fn func(SessionContext) error) error {
	return fn(NewSessionContext(ctx, sess))
}

// UseSession creates a new Session and uses it to create a new SessionContext, which is used to call the fn callback.
// The SessionContext parameter must be used as the Context parameter for any operations in the fn callback that should
// be executed under a session. After the callback returns, the created Session is ended, meaning that any in-progress
// transactions started by fn will be aborted even if fn returns an error.
//
// UseSession is safe to call from multiple goroutines concurrently. However, the SessionContext passed to the
// UseSession callback function is not safe for concurrent use by multiple goroutines.
//
// If the ctx parameter already contains a Session, that Session will be replaced with the newly created one.
//
// Any error returned by the fn callback will be returned without any modifications.
func (c *Client) UseSession(ctx context.Context, fn func(SessionContext) error) error {
	return c.UseSessionWithOptions(ctx, options.Session(), fn)
}

// UseSessionWithOptions operates like UseSession but uses the given SessionOptions to create the Session.
//
// UseSessionWithOptions is safe to call from multiple goroutines concurrently. However, the SessionContext passed to
// the UseSessionWithOptions callback function is not safe for concurrent use by multiple goroutines.
func (c *Client) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(SessionContext) error) error {
	defaultSess, err := c.StartSession(opts)
	if err != nil {
		return err
	}

	defer defaultSess.EndSession(ctx)
	return fn(NewSessionContext(ctx, defaultSess))
}

// Watch returns a change stream for all changes on the deployment. See
// https://www.mongodb.com/docs/manual/changeStreams/ for more information about change streams.
//
// The client must be configured with read concern majority or no read concern for a change stream to be created
// successfully.
//
// The pipeline parameter must be an array of documents, each representing a pipeline stage. The pipeline cannot be
// nil or empty. The stage documents must all be non-nil. See https://www.mongodb.com/docs/manual/changeStreams/ for a list
// of pipeline stages that can be used with change streams. For a pipeline of bson.D documents, the mongo.Pipeline{}
// type can be used.
//
// The opts parameter can be used to specify options for change stream creation (see the options.ChangeStreamOptions
// documentation).
func (c *Client) Watch(ctx context.Context, pipeline interface{},
	opts ...*options.ChangeStreamOptions) (*ChangeStream, error) {
	if c.sessionPool == nil {
		return nil, ErrClientDisconnected
	}

	csConfig := changeStreamConfig{
		readConcern:    c.readConcern,
		readPreference: c.readPreference,
		client:         c,
		bsonOpts:       c.bsonOpts,
		registry:       c.registry,
		streamType:     ClientStream,
		crypt:          c.cryptFLE,
	}

	return newChangeStream(ctx, csConfig, pipeline, opts...)
}

// NumberSessionsInProgress returns the number of sessions that have been started for this client but have not been
// closed (i.e. EndSession has not been called).
func (c *Client) NumberSessionsInProgress() int {
	// The underlying session pool uses an int64 for checkedOut to allow atomic
	// access. We convert to an int here to maintain backward compatibility with
	// older versions of the driver that did not atomically access checkedOut.
	return int(c.sessionPool.CheckedOut())
}

// Timeout returns the timeout set for this client.
func (c *Client) Timeout() *time.Duration {
	return c.timeout
}

func (c *Client) createBaseCursorOptions() driver.CursorOptions {
	return driver.CursorOptions{
		CommandMonitor: c.monitor,
		Crypt:          c.cryptFLE,
		ServerAPI:      c.serverAPI,
	}
}

// newLogger will use the LoggerOptions to create an internal logger and publish
// messages using a LogSink.
func newLogger(opts *options.LoggerOptions) (*logger.Logger, error) {
	// If there are no logger options, then create a default logger.
	if opts == nil {
		opts = options.Logger()
	}

	// If there are no component-level options and the environment does not
	// contain component variables, then do nothing.
	if (opts.ComponentLevels == nil || len(opts.ComponentLevels) == 0) &&
		!logger.EnvHasComponentVariables() {

		return nil, nil
	}

	// Otherwise, collect the component-level options and create a logger.
	componentLevels := make(map[logger.Component]logger.Level)
	for component, level := range opts.ComponentLevels {
		componentLevels[logger.Component(component)] = logger.Level(level)
	}

	return logger.New(opts.Sink, opts.MaxDocumentLength, componentLevels)
}
