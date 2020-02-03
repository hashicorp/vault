// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"crypto/tls"
	"errors"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/auth"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/description"
	"go.mongodb.org/mongo-driver/x/mongo/driver/operation"
	"go.mongodb.org/mongo-driver/x/mongo/driver/session"
	"go.mongodb.org/mongo-driver/x/mongo/driver/topology"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

const defaultLocalThreshold = 15 * time.Millisecond
const batchSize = 10000

// keyVaultCollOpts specifies options used to communicate with the key vault collection
var keyVaultCollOpts = options.Collection().SetReadConcern(readconcern.Majority()).
	SetWriteConcern(writeconcern.New(writeconcern.WMajority()))

// Client is a handle representing a pool of connections to a MongoDB deployment. It is safe for concurrent use by
// multiple goroutines.
//
// The Client type opens and closes connections automatically and maintains a pool of idle connections. For
// connection pool configuration options, see documentation for the ClientOptions type in the mongo/options package.
type Client struct {
	id              uuid.UUID
	topologyOptions []topology.Option
	deployment      driver.Deployment
	connString      connstring.ConnString
	localThreshold  time.Duration
	retryWrites     bool
	retryReads      bool
	clock           *session.ClusterClock
	readPreference  *readpref.ReadPref
	readConcern     *readconcern.ReadConcern
	writeConcern    *writeconcern.WriteConcern
	registry        *bsoncodec.Registry
	marshaller      BSONAppender
	monitor         *event.CommandMonitor
	sessionPool     *session.Pool

	// client-side encryption fields
	keyVaultClient *Client
	keyVaultColl   *Collection
	mongocryptd    *mcryptClient
	crypt          *driver.Crypt
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
func NewClient(opts ...*options.ClientOptions) (*Client, error) {
	clientOpt := options.MergeClientOptions(opts...)

	id, err := uuid.New()
	if err != nil {
		return nil, err
	}
	client := &Client{id: id}

	err = client.configure(clientOpt)
	if err != nil {
		return nil, err
	}

	if client.deployment == nil {
		client.deployment, err = topology.New(client.topologyOptions...)
		if err != nil {
			return nil, replaceErrors(err)
		}
	}
	return client, nil
}

// Connect initializes the Client by starting background monitoring goroutines.
// If the Client was created using the NewClient function, this method must be called before a Client can be used.
//
// Connect starts background goroutines to monitor the state of the deployment and does not do any I/O in the main
// goroutine. The Client.Ping method can be used to verify that the connection was created successfully.
func (c *Client) Connect(ctx context.Context) error {
	if connector, ok := c.deployment.(driver.Connector); ok {
		err := connector.Connect()
		if err != nil {
			return replaceErrors(err)
		}
	}

	if c.mongocryptd != nil {
		if err := c.mongocryptd.connect(ctx); err != nil {
			return err
		}
	}
	if c.keyVaultClient != nil {
		if err := c.keyVaultClient.Connect(ctx); err != nil {
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
	if ctx == nil {
		ctx = context.Background()
	}

	c.endSessions(ctx)
	if c.mongocryptd != nil {
		if err := c.mongocryptd.disconnect(ctx); err != nil {
			return err
		}
	}
	if c.keyVaultClient != nil {
		if err := c.keyVaultClient.Disconnect(ctx); err != nil {
			return err
		}
	}
	if c.crypt != nil {
		c.crypt.Close()
	}

	if disconnector, ok := c.deployment.(driver.Disconnector); ok {
		return replaceErrors(disconnector.Disconnect(ctx))
	}
	return nil
}

// Ping sends a ping command to verify that the client can connect to the deployment.
//
// The rp paramter is used to determine which server is selected for the operation.
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

	sess, err := session.NewClientSession(c.sessionPool, c.id, session.Explicit, coreOpts)
	if err != nil {
		return nil, replaceErrors(err)
	}

	sess.RetryWrite = c.retryWrites
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

	ids := c.sessionPool.IDSlice()
	idx, idArray := bsoncore.AppendArrayStart(nil)
	for i, id := range ids {
		idDoc, _ := id.MarshalBSON()
		idArray = bsoncore.AppendDocumentElement(idArray, strconv.Itoa(i), idDoc)
	}
	idArray, _ = bsoncore.AppendArrayEnd(idArray, idx)

	op := operation.NewEndSessions(idArray).ClusterClock(c.clock).Deployment(c.deployment).
		ServerSelector(description.ReadPrefSelector(readpref.PrimaryPreferred())).CommandMonitor(c.monitor).
		Database("admin").Crypt(c.crypt)

	idx, idArray = bsoncore.AppendArrayStart(nil)
	totalNumIDs := len(ids)
	for i := 0; i < totalNumIDs; i++ {
		idDoc, _ := ids[i].MarshalBSON()
		idArray = bsoncore.AppendDocumentElement(idArray, strconv.Itoa(i), idDoc)
		if ((i+1)%batchSize) == 0 || i == totalNumIDs-1 {
			idArray, _ = bsoncore.AppendArrayEnd(idArray, idx)
			_ = op.SessionIDs(idArray).Execute(ctx)
			idArray = idArray[:0]
			idx = 0
		}
	}

}

func (c *Client) configure(opts *options.ClientOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	var connOpts []topology.ConnectionOption
	var serverOpts []topology.ServerOption
	var topologyOpts []topology.Option

	// TODO(GODRIVER-814): Add tests for topology, server, and connection related options.

	// AppName
	var appName string
	if opts.AppName != nil {
		appName = *opts.AppName
	}
	// Compressors & ZlibLevel
	var comps []string
	if len(opts.Compressors) > 0 {
		comps = opts.Compressors

		connOpts = append(connOpts, topology.WithCompressors(
			func(compressors []string) []string {
				return append(compressors, comps...)
			},
		))

		for _, comp := range comps {
			switch comp {
			case "zlib":
				connOpts = append(connOpts, topology.WithZlibLevel(func(level *int) *int {
					return opts.ZlibLevel
				}))
			case "zstd":
				connOpts = append(connOpts, topology.WithZstdLevel(func(level *int) *int {
					return opts.ZstdLevel
				}))
			}
		}

		serverOpts = append(serverOpts, topology.WithCompressionOptions(
			func(opts ...string) []string { return append(opts, comps...) },
		))
	}
	// Handshaker
	var handshaker = func(driver.Handshaker) driver.Handshaker {
		return operation.NewIsMaster().AppName(appName).Compressors(comps)
	}
	// Auth & Database & Password & Username
	if opts.Auth != nil {
		cred := &auth.Cred{
			Username:    opts.Auth.Username,
			Password:    opts.Auth.Password,
			PasswordSet: opts.Auth.PasswordSet,
			Props:       opts.Auth.AuthMechanismProperties,
			Source:      opts.Auth.AuthSource,
		}
		mechanism := opts.Auth.AuthMechanism

		if len(cred.Source) == 0 {
			switch strings.ToUpper(mechanism) {
			case auth.MongoDBX509, auth.GSSAPI, auth.PLAIN:
				cred.Source = "$external"
			default:
				cred.Source = "admin"
			}
		}

		authenticator, err := auth.CreateAuthenticator(mechanism, cred)
		if err != nil {
			return err
		}

		handshakeOpts := &auth.HandshakeOptions{
			AppName:       appName,
			Authenticator: authenticator,
			Compressors:   comps,
		}
		if mechanism == "" {
			// Required for SASL mechanism negotiation during handshake
			handshakeOpts.DBUser = cred.Source + "." + cred.Username
		}
		if opts.AuthenticateToAnything != nil && *opts.AuthenticateToAnything {
			// Authenticate arbiters
			handshakeOpts.PerformAuthentication = func(serv description.Server) bool {
				return true
			}
		}

		handshaker = func(driver.Handshaker) driver.Handshaker {
			return auth.Handshaker(nil, handshakeOpts)
		}
	}
	connOpts = append(connOpts, topology.WithHandshaker(handshaker))
	// ConnectTimeout
	if opts.ConnectTimeout != nil {
		serverOpts = append(serverOpts, topology.WithHeartbeatTimeout(
			func(time.Duration) time.Duration { return *opts.ConnectTimeout },
		))
		connOpts = append(connOpts, topology.WithConnectTimeout(
			func(time.Duration) time.Duration { return *opts.ConnectTimeout },
		))
	}
	// Dialer
	if opts.Dialer != nil {
		connOpts = append(connOpts, topology.WithDialer(
			func(topology.Dialer) topology.Dialer { return opts.Dialer },
		))
	}
	// Direct
	if opts.Direct != nil && *opts.Direct {
		topologyOpts = append(topologyOpts, topology.WithMode(
			func(topology.MonitorMode) topology.MonitorMode { return topology.SingleMode },
		))
	}
	// HeartbeatInterval
	if opts.HeartbeatInterval != nil {
		serverOpts = append(serverOpts, topology.WithHeartbeatInterval(
			func(time.Duration) time.Duration { return *opts.HeartbeatInterval },
		))
	}
	// Hosts
	hosts := []string{"localhost:27017"} // default host
	if len(opts.Hosts) > 0 {
		hosts = opts.Hosts
	}
	topologyOpts = append(topologyOpts, topology.WithSeedList(
		func(...string) []string { return hosts },
	))
	// LocalThreshold
	c.localThreshold = defaultLocalThreshold
	if opts.LocalThreshold != nil {
		c.localThreshold = *opts.LocalThreshold
	}
	// MaxConIdleTime
	if opts.MaxConnIdleTime != nil {
		connOpts = append(connOpts, topology.WithIdleTimeout(
			func(time.Duration) time.Duration { return *opts.MaxConnIdleTime },
		))
	}
	// MaxPoolSize
	if opts.MaxPoolSize != nil {
		serverOpts = append(
			serverOpts,
			topology.WithMaxConnections(func(uint64) uint64 { return *opts.MaxPoolSize }),
		)
	}
	// MinPoolSize
	if opts.MinPoolSize != nil {
		serverOpts = append(
			serverOpts,
			topology.WithMinConnections(func(uint64) uint64 { return *opts.MinPoolSize }),
		)
	}
	// PoolMonitor
	if opts.PoolMonitor != nil {
		serverOpts = append(
			serverOpts,
			topology.WithConnectionPoolMonitor(func(*event.PoolMonitor) *event.PoolMonitor { return opts.PoolMonitor }),
		)
	}
	// Monitor
	if opts.Monitor != nil {
		c.monitor = opts.Monitor
		connOpts = append(connOpts, topology.WithMonitor(
			func(*event.CommandMonitor) *event.CommandMonitor { return opts.Monitor },
		))
	}
	// ReadConcern
	c.readConcern = readconcern.New()
	if opts.ReadConcern != nil {
		c.readConcern = opts.ReadConcern
	}
	// ReadPreference
	c.readPreference = readpref.Primary()
	if opts.ReadPreference != nil {
		c.readPreference = opts.ReadPreference
	}
	// Registry
	c.registry = bson.DefaultRegistry
	if opts.Registry != nil {
		c.registry = opts.Registry
	}
	// ReplicaSet
	if opts.ReplicaSet != nil {
		topologyOpts = append(topologyOpts, topology.WithReplicaSetName(
			func(string) string { return *opts.ReplicaSet },
		))
	}
	// RetryWrites
	c.retryWrites = true // retry writes on by default
	if opts.RetryWrites != nil {
		c.retryWrites = *opts.RetryWrites
	}
	c.retryReads = true
	if opts.RetryReads != nil {
		c.retryReads = *opts.RetryReads
	}
	// ServerSelectionTimeout
	if opts.ServerSelectionTimeout != nil {
		topologyOpts = append(topologyOpts, topology.WithServerSelectionTimeout(
			func(time.Duration) time.Duration { return *opts.ServerSelectionTimeout },
		))
	}
	// SocketTimeout
	if opts.SocketTimeout != nil {
		connOpts = append(
			connOpts,
			topology.WithReadTimeout(func(time.Duration) time.Duration { return *opts.SocketTimeout }),
			topology.WithWriteTimeout(func(time.Duration) time.Duration { return *opts.SocketTimeout }),
		)
	}
	// TLSConfig
	if opts.TLSConfig != nil {
		connOpts = append(connOpts, topology.WithTLSConfig(
			func(*tls.Config) *tls.Config {
				return opts.TLSConfig
			},
		))
	}
	// WriteConcern
	if opts.WriteConcern != nil {
		c.writeConcern = opts.WriteConcern
	}
	// AutoEncryptionOptions
	if opts.AutoEncryptionOptions != nil {
		if err := c.configureAutoEncryption(opts.AutoEncryptionOptions); err != nil {
			return err
		}
	}

	// ClusterClock
	c.clock = new(session.ClusterClock)

	serverOpts = append(
		serverOpts,
		topology.WithClock(func(*session.ClusterClock) *session.ClusterClock { return c.clock }),
		topology.WithConnectionOptions(func(...topology.ConnectionOption) []topology.ConnectionOption { return connOpts }),
	)
	c.topologyOptions = append(topologyOpts, topology.WithServerOptions(
		func(...topology.ServerOption) []topology.ServerOption { return serverOpts },
	))

	// Deployment
	if opts.Deployment != nil {
		if len(serverOpts) > 2 || len(topologyOpts) > 1 {
			return errors.New("cannot specify topology or server options with a deployment")
		}
		c.deployment = opts.Deployment
	}

	return nil
}

func (c *Client) configureAutoEncryption(opts *options.AutoEncryptionOptions) error {
	if err := c.configureKeyVault(opts); err != nil {
		return err
	}
	if err := c.configureMongocryptd(opts); err != nil {
		return err
	}
	return c.configureCrypt(opts)
}

func (c *Client) configureKeyVault(opts *options.AutoEncryptionOptions) error {
	// parse key vault options and create new client if necessary
	if opts.KeyVaultClientOptions != nil {
		var err error
		c.keyVaultClient, err = NewClient(opts.KeyVaultClientOptions)
		if err != nil {
			return err
		}
	}

	dbName, collName := splitNamespace(opts.KeyVaultNamespace)
	client := c.keyVaultClient
	if client == nil {
		client = c
	}
	c.keyVaultColl = client.Database(dbName).Collection(collName, keyVaultCollOpts)
	return nil
}

func (c *Client) configureMongocryptd(opts *options.AutoEncryptionOptions) error {
	var err error
	c.mongocryptd, err = newMcryptClient(opts.ExtraOptions)
	return err
}

func (c *Client) configureCrypt(opts *options.AutoEncryptionOptions) error {
	// convert schemas in SchemaMap to bsoncore documents
	cryptSchemaMap := make(map[string]bsoncore.Document)
	for k, v := range opts.SchemaMap {
		schema, err := transformBsoncoreDocument(c.registry, v)
		if err != nil {
			return err
		}
		cryptSchemaMap[k] = schema
	}

	// configure options
	var bypass bool
	if opts.BypassAutoEncryption != nil {
		bypass = *opts.BypassAutoEncryption
	}
	kr := keyRetriever{coll: c.keyVaultColl}
	cir := collInfoRetriever{client: c}
	cryptOpts := &driver.CryptOptions{
		CollInfoFn:           cir.cryptCollInfo,
		KeyFn:                kr.cryptKeys,
		MarkFn:               c.mongocryptd.markCommand,
		KmsProviders:         opts.KmsProviders,
		BypassAutoEncryption: bypass,
		SchemaMap:            cryptSchemaMap,
	}

	var err error
	c.crypt, err = driver.NewCrypt(cryptOpts)
	return err
}

// validSession returns an error if the session doesn't belong to the client
func (c *Client) validSession(sess *session.Client) error {
	if sess != nil && !uuid.Equal(sess.ClientID, c.id) {
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
// The opts paramter can be used to specify options for this operation (see the options.ListDatabasesOptions documentation).
//
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/listDatabases/.
func (c *Client) ListDatabases(ctx context.Context, filter interface{}, opts ...*options.ListDatabasesOptions) (ListDatabasesResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	sess := sessionFromContext(ctx)

	err := c.validSession(sess)
	if sess == nil && c.sessionPool != nil {
		sess, err = session.NewClientSession(c.sessionPool, c.id, session.Implicit)
		if err != nil {
			return ListDatabasesResult{}, err
		}
		defer sess.EndSession()
	}

	err = c.validSession(sess)
	if err != nil {
		return ListDatabasesResult{}, err
	}

	filterDoc, err := transformBsoncoreDocument(c.registry, filter)
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
		ServerSelector(selector).ClusterClock(c.clock).Database("admin").Deployment(c.deployment).Crypt(c.crypt)
	if ldo.NameOnly != nil {
		op = op.NameOnly(*ldo.NameOnly)
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
// For more information about the command, see https://docs.mongodb.com/manual/reference/command/listDatabases/.
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
// If the ctx parameter already contains a Session, that Session will be replaced with the one provided.
//
// Any error returned by the fn callback will be returned without any modifications.
func WithSession(ctx context.Context, sess Session, fn func(SessionContext) error) error {
	return fn(contextWithSession(ctx, sess))
}

// UseSession creates a new Session and uses it to create a new SessionContext, which is used to call the fn callback.
// The SessionContext parameter must be used as the Context parameter for any operations in the fn callback that should
// be executed under a session. After the callback returns, the created Session is ended, meaning that any in-progress
// transactions started by fn will be aborted even if fn returns an error.
//
// If the ctx parameter already contains a Session, that Session will be replaced with the newly created one.
//
// Any error returned by the fn callback will be returned without any modifications.
func (c *Client) UseSession(ctx context.Context, fn func(SessionContext) error) error {
	return c.UseSessionWithOptions(ctx, options.Session(), fn)
}

// UseSessionWithOptions operates like UseSession but uses the given SessionOptions to create the Session.
func (c *Client) UseSessionWithOptions(ctx context.Context, opts *options.SessionOptions, fn func(SessionContext) error) error {
	defaultSess, err := c.StartSession(opts)
	if err != nil {
		return err
	}

	defer defaultSess.EndSession(ctx)

	sessCtx := sessionContext{
		Context: context.WithValue(ctx, sessionKey{}, defaultSess),
		Session: defaultSess,
	}

	return fn(sessCtx)
}

// Watch returns a change stream for all changes on the deployment. See
// https://docs.mongodb.com/manual/changeStreams/ for more information about change streams.
//
// The client must be configured with read concern majority or no read concern for a change stream to be created
// successfully.
//
// The pipeline parameter must be an array of documents, each representing a pipeline stage. The pipeline cannot be
// nil or empty. The stage documents must all be non-nil. See https://docs.mongodb.com/manual/changeStreams/ for a list
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
		registry:       c.registry,
		streamType:     ClientStream,
	}

	return newChangeStream(ctx, csConfig, pipeline, opts...)
}

// NumberSessionsInProgress returns the number of sessions that have been started for this client but have not been
// closed (i.e. EndSession has not been called).
func (c *Client) NumberSessionsInProgress() int {
	return c.sessionPool.CheckedOut()
}
