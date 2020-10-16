// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options // import "go.mongodb.org/mongo-driver/mongo/options"

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/tag"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// ContextDialer is an interface that can be implemented by types that can create connections. It should be used to
// provide a custom dialer when configuring a Client.
//
// DialContext should return a connection to the provided address on the given network.
type ContextDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// Credential can be used to provide authentication options when configuring a Client.
//
// AuthMechanism: the mechanism to use for authentication. Supported values include "SCRAM-SHA-256", "SCRAM-SHA-1",
// "MONGODB-CR", "PLAIN", "GSSAPI", "MONGODB-X509", and "MONGODB-AWS". This can also be set through the "authMechanism"
// URI option. (e.g. "authMechanism=PLAIN"). For more information, see
// https://docs.mongodb.com/manual/core/authentication-mechanisms/.
//
// AuthMechanismProperties can be used to specify additional configuration options for certain mechanisms. They can also
// be set through the "authMechanismProperites" URI option
// (e.g. "authMechanismProperties=SERVICE_NAME:service,CANONICALIZE_HOST_NAME:true"). Supported properties are:
//
// 1. SERVICE_NAME: The service name to use for GSSAPI authentication. The default is "mongodb".
//
// 2. CANONICALIZE_HOST_NAME: If "true", the driver will canonicalize the host name for GSSAPI authentication. The default
// is "false".
//
// 3. SERVICE_REALM: The service realm for GSSAPI authentication.
//
// 4. SERVICE_HOST: The host name to use for GSSAPI authentication. This should be specified if the host name to use for
// authentication is different than the one given for Client construction.
//
// 4. AWS_SESSION_TOKEN: The AWS token for MONGODB-AWS authentication. This is optional and used for authentication with
// temporary credentials.
//
// The SERVICE_HOST and CANONICALIZE_HOST_NAME properties must not be used at the same time on Linux and Darwin
// systems.
//
// AuthSource: the name of the database to use for authentication. This defaults to "$external" for MONGODB-X509,
// GSSAPI, and PLAIN and "admin" for all other mechanisms. This can also be set through the "authSource" URI option
// (e.g. "authSource=otherDb").
//
// Username: the username for authentication. This can also be set through the URI as a username:password pair before
// the first @ character. For example, a URI for user "user", password "pwd", and host "localhost:27017" would be
// "mongodb://user:pwd@localhost:27017". This is optional for X509 authentication and will be extracted from the
// client certificate if not specified.
//
// Password: the password for authentication. This must not be specified for X509 and is optional for GSSAPI
// authentication.
//
// PasswordSet: For GSSAPI, this must be true if a password is specified, even if the password is the empty string, and
// false if no password is specified, indicating that the password should be taken from the context of the running
// process. For other mechanisms, this field is ignored.
type Credential struct {
	AuthMechanism           string
	AuthMechanismProperties map[string]string
	AuthSource              string
	Username                string
	Password                string
	PasswordSet             bool
}

// ClientOptions contains options to configure a Client instance. Each option can be set through setter functions. See
// documentation for each setter function for an explanation of the option.
type ClientOptions struct {
	AppName                  *string
	Auth                     *Credential
	AutoEncryptionOptions    *AutoEncryptionOptions
	ConnectTimeout           *time.Duration
	Compressors              []string
	Dialer                   ContextDialer
	Direct                   *bool
	DisableOCSPEndpointCheck *bool
	HeartbeatInterval        *time.Duration
	Hosts                    []string
	LocalThreshold           *time.Duration
	MaxConnIdleTime          *time.Duration
	MaxPoolSize              *uint64
	MinPoolSize              *uint64
	PoolMonitor              *event.PoolMonitor
	Monitor                  *event.CommandMonitor
	ReadConcern              *readconcern.ReadConcern
	ReadPreference           *readpref.ReadPref
	Registry                 *bsoncodec.Registry
	ReplicaSet               *string
	RetryReads               *bool
	RetryWrites              *bool
	ServerSelectionTimeout   *time.Duration
	SocketTimeout            *time.Duration
	TLSConfig                *tls.Config
	WriteConcern             *writeconcern.WriteConcern
	ZlibLevel                *int
	ZstdLevel                *int

	err error
	uri string
	cs  *connstring.ConnString

	// AuthenticateToAnything skips server type checks when deciding if authentication is possible.
	//
	// Deprecated: This option is for internal use only and should not be set. It may be changed or removed in any
	// release.
	AuthenticateToAnything *bool

	// Deployment specifies a custom deployment to use for the new Client.
	//
	// Deprecated: This option is for internal use only and should not be set. It may be changed or removed in any
	// release.
	Deployment driver.Deployment
}

// Client creates a new ClientOptions instance.
func Client() *ClientOptions {
	return new(ClientOptions)
}

// Validate validates the client options. This method will return the first error found.
func (c *ClientOptions) Validate() error {
	c.validateAndSetError()
	return c.err
}

func (c *ClientOptions) validateAndSetError() {
	if c.err != nil {
		return
	}

	// Direct connections cannot be made if multiple hosts are specified or an SRV URI is used.
	if c.Direct != nil && *c.Direct {
		if len(c.Hosts) > 1 {
			c.err = errors.New("a direct connection cannot be made if multiple hosts are specified")
			return
		}
		if c.cs != nil && c.cs.Scheme == connstring.SchemeMongoDBSRV {
			c.err = errors.New("a direct connection cannot be made if an SRV URI is used")
			return
		}
	}
}

// GetURI returns the original URI used to configure the ClientOptions instance. If ApplyURI was not called during
// construction, this returns "".
func (c *ClientOptions) GetURI() string {
	return c.uri
}

// ApplyURI parses the given URI and sets options accordingly. The URI can contain host names, IPv4/IPv6 literals, or
// an SRV record that will be resolved when the Client is created. When using an SRV record, TLS support is
// implictly enabled. Specify the "tls=false" URI option to override this.
//
// If the connection string contains any options that have previously been set, it will overwrite them. Options that
// correspond to multiple URI parameters, such as WriteConcern, will be completely overwritten if any of the query
// parameters are specified. If an option is set on ClientOptions after this method is called, that option will override
// any option applied via the connection string.
//
// If the URI format is incorrect or there are conflicing options specified in the URI an error will be recorded and
// can be retrieved by calling Validate.
//
// For more information about the URI format, see https://docs.mongodb.com/manual/reference/connection-string/. See
// mongo.Connect documentation for examples of using URIs for different Client configurations.
func (c *ClientOptions) ApplyURI(uri string) *ClientOptions {
	if c.err != nil {
		return c
	}

	c.uri = uri
	cs, err := connstring.ParseAndValidate(uri)
	if err != nil {
		c.err = err
		return c
	}
	c.cs = &cs

	if cs.AppName != "" {
		c.AppName = &cs.AppName
	}

	// Only create a Credential if there is a request for authentication via non-empty credentials in the URI.
	if cs.HasAuthParameters() {
		c.Auth = &Credential{
			AuthMechanism:           cs.AuthMechanism,
			AuthMechanismProperties: cs.AuthMechanismProperties,
			AuthSource:              cs.AuthSource,
			Username:                cs.Username,
			Password:                cs.Password,
			PasswordSet:             cs.PasswordSet,
		}
	}

	if cs.ConnectSet {
		direct := cs.Connect == connstring.SingleConnect
		c.Direct = &direct
	}

	if cs.DirectConnectionSet {
		c.Direct = &cs.DirectConnection
	}

	if cs.ConnectTimeoutSet {
		c.ConnectTimeout = &cs.ConnectTimeout
	}

	if len(cs.Compressors) > 0 {
		c.Compressors = cs.Compressors
		if stringSliceContains(c.Compressors, "zlib") {
			defaultLevel := wiremessage.DefaultZlibLevel
			c.ZlibLevel = &defaultLevel
		}
		if stringSliceContains(c.Compressors, "zstd") {
			defaultLevel := wiremessage.DefaultZstdLevel
			c.ZstdLevel = &defaultLevel
		}
	}

	if cs.HeartbeatIntervalSet {
		c.HeartbeatInterval = &cs.HeartbeatInterval
	}

	c.Hosts = cs.Hosts

	if cs.LocalThresholdSet {
		c.LocalThreshold = &cs.LocalThreshold
	}

	if cs.MaxConnIdleTimeSet {
		c.MaxConnIdleTime = &cs.MaxConnIdleTime
	}

	if cs.MaxPoolSizeSet {
		c.MaxPoolSize = &cs.MaxPoolSize
	}

	if cs.MinPoolSizeSet {
		c.MinPoolSize = &cs.MinPoolSize
	}

	if cs.ReadConcernLevel != "" {
		c.ReadConcern = readconcern.New(readconcern.Level(cs.ReadConcernLevel))
	}

	if cs.ReadPreference != "" || len(cs.ReadPreferenceTagSets) > 0 || cs.MaxStalenessSet {
		opts := make([]readpref.Option, 0, 1)

		tagSets := tag.NewTagSetsFromMaps(cs.ReadPreferenceTagSets)
		if len(tagSets) > 0 {
			opts = append(opts, readpref.WithTagSets(tagSets...))
		}

		if cs.MaxStaleness != 0 {
			opts = append(opts, readpref.WithMaxStaleness(cs.MaxStaleness))
		}

		mode, err := readpref.ModeFromString(cs.ReadPreference)
		if err != nil {
			c.err = err
			return c
		}

		c.ReadPreference, c.err = readpref.New(mode, opts...)
		if c.err != nil {
			return c
		}
	}

	if cs.RetryWritesSet {
		c.RetryWrites = &cs.RetryWrites
	}

	if cs.RetryReadsSet {
		c.RetryReads = &cs.RetryReads
	}

	if cs.ReplicaSet != "" {
		c.ReplicaSet = &cs.ReplicaSet
	}

	if cs.ServerSelectionTimeoutSet {
		c.ServerSelectionTimeout = &cs.ServerSelectionTimeout
	}

	if cs.SocketTimeoutSet {
		c.SocketTimeout = &cs.SocketTimeout
	}

	if cs.SSL {
		tlsConfig := new(tls.Config)

		if cs.SSLCaFileSet {
			c.err = addCACertFromFile(tlsConfig, cs.SSLCaFile)
			if c.err != nil {
				return c
			}
		}

		if cs.SSLInsecure {
			tlsConfig.InsecureSkipVerify = true
		}

		var x509Subject string
		var keyPasswd string
		if cs.SSLClientCertificateKeyPasswordSet && cs.SSLClientCertificateKeyPassword != nil {
			keyPasswd = cs.SSLClientCertificateKeyPassword()
		}
		if cs.SSLClientCertificateKeyFileSet {
			x509Subject, err = addClientCertFromConcatenatedFile(tlsConfig, cs.SSLClientCertificateKeyFile, keyPasswd)
		} else if cs.SSLCertificateFileSet || cs.SSLPrivateKeyFileSet {
			x509Subject, err = addClientCertFromSeparateFiles(tlsConfig, cs.SSLCertificateFile,
				cs.SSLPrivateKeyFile, keyPasswd)
		}
		if err != nil {
			c.err = err
			return c
		}

		// If a username wasn't specified fork x509, add one from the certificate.
		if c.Auth != nil && strings.ToLower(c.Auth.AuthMechanism) == "mongodb-x509" &&
			c.Auth.Username == "" {

			// The Go x509 package gives the subject with the pairs in reverse order that we want.
			c.Auth.Username = extractX509UsernameFromSubject(x509Subject)
		}

		c.TLSConfig = tlsConfig
	}

	if cs.JSet || cs.WString != "" || cs.WNumberSet || cs.WTimeoutSet {
		opts := make([]writeconcern.Option, 0, 1)

		if len(cs.WString) > 0 {
			opts = append(opts, writeconcern.WTagSet(cs.WString))
		} else if cs.WNumberSet {
			opts = append(opts, writeconcern.W(cs.WNumber))
		}

		if cs.JSet {
			opts = append(opts, writeconcern.J(cs.J))
		}

		if cs.WTimeoutSet {
			opts = append(opts, writeconcern.WTimeout(cs.WTimeout))
		}

		c.WriteConcern = writeconcern.New(opts...)
	}

	if cs.ZlibLevelSet {
		c.ZlibLevel = &cs.ZlibLevel
	}
	if cs.ZstdLevelSet {
		c.ZstdLevel = &cs.ZstdLevel
	}

	if cs.SSLDisableOCSPEndpointCheckSet {
		c.DisableOCSPEndpointCheck = &cs.SSLDisableOCSPEndpointCheck
	}

	return c
}

// SetAppName specifies an application name that is sent to the server when creating new connections. It is used by the
// server to log connection and profiling information (e.g. slow query logs). This can also be set through the "appName"
// URI option (e.g "appName=example_application"). The default is empty, meaning no app name will be sent.
func (c *ClientOptions) SetAppName(s string) *ClientOptions {
	c.AppName = &s
	return c
}

// SetAuth specifies a Credential containing options for configuring authentication. See the options.Credential
// documentation for more information about Credential fields. The default is an empty Credential, meaning no
// authentication will be configured.
func (c *ClientOptions) SetAuth(auth Credential) *ClientOptions {
	c.Auth = &auth
	return c
}

// SetCompressors sets the compressors that can be used when communicating with a server. Valid values are:
//
// 1. "snappy" - requires server version >= 3.4
//
// 2. "zlib" - requires server version >= 3.6
//
// 3. "zstd" - requires server version >= 4.2, and driver version >= 1.2.0 with cgo support enabled or driver version >= 1.3.0
//    without cgo
//
// If this option is specified, the driver will perform a negotiation with the server to determine a common list of of
// compressors and will use the first one in that list when performing operations. See
// https://docs.mongodb.com/manual/reference/program/mongod/#cmdoption-mongod-networkmessagecompressors for more
// information about configuring compression on the server and the server-side defaults.
//
// This can also be set through the "compressors" URI option (e.g. "compressors=zstd,zlib,snappy"). The default is
// an empty slice, meaning no compression will be enabled.
func (c *ClientOptions) SetCompressors(comps []string) *ClientOptions {
	c.Compressors = comps

	return c
}

// SetConnectTimeout specifies a timeout that is used for creating connections to the server. If a custom Dialer is
// specified through SetDialer, this option must not be used. This can be set through ApplyURI with the
// "connectTimeoutMS" (e.g "connectTimeoutMS=30") option. If set to 0, no timeout will be used. The default is 30
// seconds.
func (c *ClientOptions) SetConnectTimeout(d time.Duration) *ClientOptions {
	c.ConnectTimeout = &d
	return c
}

// SetDialer specifies a custom ContextDialer to be used to create new connections to the server. The default is a
// net.Dialer with the Timeout field set to ConnectTimeout. See https://golang.org/pkg/net/#Dialer for more information
// about the net.Dialer type.
func (c *ClientOptions) SetDialer(d ContextDialer) *ClientOptions {
	c.Dialer = d
	return c
}

// SetDirect specifies whether or not a direct connect should be made. If set to true, the driver will only connect to
// the host provided in the URI and will not discover other hosts in the cluster. This can also be set through the
// "directConnection" URI option. This option cannot be set to true if multiple hosts are specified, either through
// ApplyURI or SetHosts, or an SRV URI is used.
//
// As of driver version 1.4, the "connect" URI option has been deprecated and replaced with "directConnection". The
// "connect" URI option has two values:
//
// 1. "connect=direct" for direct connections. This corresponds to "directConnection=true".
//
// 2. "connect=automatic" for automatic discovery. This corresponds to "directConnection=false"
//
// If the "connect" and "directConnection" URI options are both specified in the connection string, their values must
// not conflict. Direct connections are not valid if multiple hosts are specified or an SRV URI is used. The default
// value for this option is false.
func (c *ClientOptions) SetDirect(b bool) *ClientOptions {
	c.Direct = &b
	return c
}

// SetHeartbeatInterval specifies the amount of time to wait between periodic background server checks. This can also be
// set through the "heartbeatIntervalMS" URI option (e.g. "heartbeatIntervalMS=10000"). The default is 10 seconds.
func (c *ClientOptions) SetHeartbeatInterval(d time.Duration) *ClientOptions {
	c.HeartbeatInterval = &d
	return c
}

// SetHosts specifies a list of host names or IP addresses for servers in a cluster. Both IPv4 and IPv6 addresses are
// supported. IPv6 literals must be enclosed in '[]' following RFC-2732 syntax.
//
// Hosts can also be specified as a comma-separated list in a URI. For example, to include "localhost:27017" and
// "localhost:27018", a URI could be "mongodb://localhost:27017,localhost:27018". The default is ["localhost:27017"]
func (c *ClientOptions) SetHosts(s []string) *ClientOptions {
	c.Hosts = s
	return c
}

// SetLocalThreshold specifies the width of the 'latency window': when choosing between multiple suitable servers for an
// operation, this is the acceptable non-negative delta between shortest and longest average round-trip times. A server
// within the latency window is selected randomly. This can also be set through the "localThresholdMS" URI option (e.g.
// "localThresholdMS=15000"). The default is 15 milliseconds.
func (c *ClientOptions) SetLocalThreshold(d time.Duration) *ClientOptions {
	c.LocalThreshold = &d
	return c
}

// SetMaxConnIdleTime specifies the maximum amount of time that a connection will remain idle in a connection pool
// before it is removed from the pool and closed. This can also be set through the "maxIdleTimeMS" URI option (e.g.
// "maxIdleTimeMS=10000"). The default is 0, meaning a connection can remain unused indefinitely.
func (c *ClientOptions) SetMaxConnIdleTime(d time.Duration) *ClientOptions {
	c.MaxConnIdleTime = &d
	return c
}

// SetMaxPoolSize specifies that maximum number of connections allowed in the driver's connection pool to each server.
// Requests to a server will block if this maximum is reached. This can also be set through the "maxPoolSize" URI option
// (e.g. "maxPoolSize=100"). The default is 100. If this is 0, it will be set to math.MaxInt64.
func (c *ClientOptions) SetMaxPoolSize(u uint64) *ClientOptions {
	c.MaxPoolSize = &u
	return c
}

// SetMinPoolSize specifies the minimum number of connections allowed in the driver's connection pool to each server. If
// this is non-zero, each server's pool will be maintained in the background to ensure that the size does not fall below
// the minimum. This can also be set through the "minPoolSize" URI option (e.g. "minPoolSize=100"). The default is 0.
func (c *ClientOptions) SetMinPoolSize(u uint64) *ClientOptions {
	c.MinPoolSize = &u
	return c
}

// SetPoolMonitor specifies a PoolMonitor to receive connection pool events. See the event.PoolMonitor documentation
// for more information about the structure of the monitor and events that can be received.
func (c *ClientOptions) SetPoolMonitor(m *event.PoolMonitor) *ClientOptions {
	c.PoolMonitor = m
	return c
}

// SetMonitor specifies a CommandMonitor to receive command events. See the event.CommandMonitor documentation for more
// information about the structure of the monitor and events that can be received.
func (c *ClientOptions) SetMonitor(m *event.CommandMonitor) *ClientOptions {
	c.Monitor = m
	return c
}

// SetReadConcern specifies the read concern to use for read operations. A read concern level can also be set through
// the "readConcernLevel" URI option (e.g. "readConcernLevel=majority"). The default is nil, meaning the server will use
// its configured default.
func (c *ClientOptions) SetReadConcern(rc *readconcern.ReadConcern) *ClientOptions {
	c.ReadConcern = rc

	return c
}

// SetReadPreference specifies the read preference to use for read operations. This can also be set through the
// following URI options:
//
// 1. "readPreference" - Specifiy the read preference mode (e.g. "readPreference=primary").
//
// 2. "readPreferenceTags": Specify one or more read preference tags
// (e.g. "readPreferenceTags=region:south,datacenter:A").
//
// 3. "maxStalenessSeconds" (or "maxStaleness"): Specify a maximum replication lag for reads from secondaries in a
// replica set (e.g. "maxStalenessSeconds=10").
//
// The default is readpref.Primary(). See https://docs.mongodb.com/manual/core/read-preference/#read-preference for
// more information about read preferences.
func (c *ClientOptions) SetReadPreference(rp *readpref.ReadPref) *ClientOptions {
	c.ReadPreference = rp

	return c
}

// SetRegistry specifies the BSON registry to use for BSON marshalling/unmarshalling operations. The default is
// bson.DefaultRegistry.
func (c *ClientOptions) SetRegistry(registry *bsoncodec.Registry) *ClientOptions {
	c.Registry = registry
	return c
}

// SetReplicaSet specifies the replica set name for the cluster. If specified, the cluster will be treated as a replica
// set and the driver will automatically discover all servers in the set, starting with the nodes specified through
// ApplyURI or SetHosts. All nodes in the replica set must have the same replica set name, or they will not be
// considered as part of the set by the Client. This can also be set through the "replicaSet" URI option (e.g.
// "replicaSet=replset"). The default is empty.
func (c *ClientOptions) SetReplicaSet(s string) *ClientOptions {
	c.ReplicaSet = &s
	return c
}

// SetRetryWrites specifies whether supported write operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are InsertOne, UpdateOne, ReplaceOne, DeleteOne, FindOneAndDelete, FindOneAndReplace,
// FindOneAndDelete, InsertMany, and BulkWrite. Note that BulkWrite requests must not include UpdateManyModel or
// DeleteManyModel instances to be considered retryable. Unacknowledged writes will not be retried, even if this option
// is set to true.
//
// This option requires server version >= 3.6 and a replica set or sharded cluster and will be ignored for any other
// cluster type. This can also be set through the "retryWrites" URI option (e.g. "retryWrites=true"). The default is
// true.
func (c *ClientOptions) SetRetryWrites(b bool) *ClientOptions {
	c.RetryWrites = &b

	return c
}

// SetRetryReads specifies whether supported read operations should be retried once on certain errors, such as network
// errors.
//
// Supported operations are Find, FindOne, Aggregate without a $out stage, Distinct, CountDocuments,
// EstimatedDocumentCount, Watch (for Client, Database, and Collection), ListCollections, and ListDatabases. Note that
// operations run through RunCommand are not retried.
//
// This option requires server version >= 3.6 and driver version >= 1.1.0. The default is true.
func (c *ClientOptions) SetRetryReads(b bool) *ClientOptions {
	c.RetryReads = &b
	return c
}

// SetServerSelectionTimeout specifies how long the driver will wait to find an available, suitable server to execute an
// operation. This can also be set through the "serverSelectionTimeoutMS" URI option (e.g.
// "serverSelectionTimeoutMS=30000"). The default value is 30 seconds.
func (c *ClientOptions) SetServerSelectionTimeout(d time.Duration) *ClientOptions {
	c.ServerSelectionTimeout = &d
	return c
}

// SetSocketTimeout specifies how long the driver will wait for a socket read or write to return before returning a
// network error. This can also be set through the "socketTimeoutMS" URI option (e.g. "socketTimeoutMS=1000"). The
// default value is 0, meaning no timeout is used and socket operations can block indefinitely.
func (c *ClientOptions) SetSocketTimeout(d time.Duration) *ClientOptions {
	c.SocketTimeout = &d
	return c
}

// SetTLSConfig specifies a tls.Config instance to use use to configure TLS on all connections created to the cluster.
// This can also be set through the following URI options:
//
// 1. "tls" (or "ssl"): Specify if TLS should be used (e.g. "tls=true").
//
// 2. Either "tlsCertificateKeyFile" (or "sslClientCertificateKeyFile") or a combination of "tlsCertificateFile" and
// "tlsPrivateKeyFile". The "tlsCertificateKeyFile" option specifies a path to the client certificate and private key,
// which must be concatenated into one file. The "tlsCertificateFile" and "tlsPrivateKey" combination specifies separate
// paths to the client certificate and private key, respectively. Note that if "tlsCertificateKeyFile" is used, the
// other two options must not be specified.
//
// 3. "tlsCertificateKeyFilePassword" (or "sslClientCertificateKeyPassword"): Specify the password to decrypt the client
// private key file (e.g. "tlsCertificateKeyFilePassword=password").
//
// 4. "tlsCaFile" (or "sslCertificateAuthorityFile"): Specify the path to a single or bundle of certificate authorities
// to be considered trusted when making a TLS connection (e.g. "tlsCaFile=/path/to/caFile").
//
// 5. "tlsInsecure" (or "sslInsecure"): Specifies whether or not certificates and hostnames received from the server
// should be validated. If true (e.g. "tlsInsecure=true"), the TLS library will accept any certificate presented by the
// server and any host name in that certificate. Note that setting this to true makes TLS susceptible to
// man-in-the-middle attacks and should only be done for testing.
//
// The default is nil, meaning no TLS will be enabled.
func (c *ClientOptions) SetTLSConfig(cfg *tls.Config) *ClientOptions {
	c.TLSConfig = cfg
	return c
}

// SetWriteConcern specifies the write concern to use to for write operations. This can also be set through the following
// URI options:
//
// 1. "w": Specify the number of nodes in the cluster that must acknowledge write operations before the operation
// returns or "majority" to specify that a majority of the nodes must acknowledge writes. This can either be an integer
// (e.g. "w=10") or the string "majority" (e.g. "w=majority").
//
// 2. "wTimeoutMS": Specify how long write operations should wait for the correct number of nodes to acknowledge the
// operation (e.g. "wTimeoutMS=1000").
//
// 3. "journal": Specifies whether or not write operations should be written to an on-disk journal on the server before
// returning (e.g. "journal=true").
//
// The default is nil, meaning the server will use its configured default.
func (c *ClientOptions) SetWriteConcern(wc *writeconcern.WriteConcern) *ClientOptions {
	c.WriteConcern = wc

	return c
}

// SetZlibLevel specifies the level for the zlib compressor. This option is ignored if zlib is not specified as a
// compressor through ApplyURI or SetCompressors. Supported values are -1 through 9, inclusive. -1 tells the zlib
// library to use its default, 0 means no compression, 1 means best speed, and 9 means best compression.
// This can also be set through the "zlibCompressionLevel" URI option (e.g. "zlibCompressionLevel=-1"). Defaults to -1.
func (c *ClientOptions) SetZlibLevel(level int) *ClientOptions {
	c.ZlibLevel = &level

	return c
}

// SetZstdLevel sets the level for the zstd compressor. This option is ignored if zstd is not specified as a compressor
// through ApplyURI or SetCompressors. Supported values are 1 through 20, inclusive. 1 means best speed and 20 means
// best compression. This can also be set through the "zstdCompressionLevel" URI option. Defaults to 6.
func (c *ClientOptions) SetZstdLevel(level int) *ClientOptions {
	c.ZstdLevel = &level
	return c
}

// SetAutoEncryptionOptions specifies an AutoEncryptionOptions instance to automatically encrypt and decrypt commands
// and their results. See the options.AutoEncryptionOptions documentation for more information about the supported
// options.
func (c *ClientOptions) SetAutoEncryptionOptions(opts *AutoEncryptionOptions) *ClientOptions {
	c.AutoEncryptionOptions = opts
	return c
}

// SetDisableOCSPEndpointCheck specifies whether or not the driver should reach out to OCSP responders to verify the
// certificate status for certificates presented by the server that contain a list of OCSP responders.
//
// If set to true, the driver will verify the status of the certificate using a response stapled by the server, if there
// is one, but will not send an HTTP request to any responders if there is no staple. In this case, the driver will
// continue the connection even though the certificate status is not known.
//
// This can also be set through the tlsDisableOCSPEndpointCheck URI option. Both this URI option and tlsInsecure must
// not be set at the same time and will error if they are. The default value is false.
func (c *ClientOptions) SetDisableOCSPEndpointCheck(disableCheck bool) *ClientOptions {
	c.DisableOCSPEndpointCheck = &disableCheck
	return c
}

// MergeClientOptions combines the given *ClientOptions into a single *ClientOptions in a last one wins fashion.
// The specified options are merged with the existing options on the collection, with the specified options taking
// precedence.
func MergeClientOptions(opts ...*ClientOptions) *ClientOptions {
	c := Client()

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.Dialer != nil {
			c.Dialer = opt.Dialer
		}
		if opt.AppName != nil {
			c.AppName = opt.AppName
		}
		if opt.Auth != nil {
			c.Auth = opt.Auth
		}
		if opt.AuthenticateToAnything != nil {
			c.AuthenticateToAnything = opt.AuthenticateToAnything
		}
		if opt.Compressors != nil {
			c.Compressors = opt.Compressors
		}
		if opt.ConnectTimeout != nil {
			c.ConnectTimeout = opt.ConnectTimeout
		}
		if opt.HeartbeatInterval != nil {
			c.HeartbeatInterval = opt.HeartbeatInterval
		}
		if len(opt.Hosts) > 0 {
			c.Hosts = opt.Hosts
		}
		if opt.LocalThreshold != nil {
			c.LocalThreshold = opt.LocalThreshold
		}
		if opt.MaxConnIdleTime != nil {
			c.MaxConnIdleTime = opt.MaxConnIdleTime
		}
		if opt.MaxPoolSize != nil {
			c.MaxPoolSize = opt.MaxPoolSize
		}
		if opt.MinPoolSize != nil {
			c.MinPoolSize = opt.MinPoolSize
		}
		if opt.PoolMonitor != nil {
			c.PoolMonitor = opt.PoolMonitor
		}
		if opt.Monitor != nil {
			c.Monitor = opt.Monitor
		}
		if opt.ReadConcern != nil {
			c.ReadConcern = opt.ReadConcern
		}
		if opt.ReadPreference != nil {
			c.ReadPreference = opt.ReadPreference
		}
		if opt.Registry != nil {
			c.Registry = opt.Registry
		}
		if opt.ReplicaSet != nil {
			c.ReplicaSet = opt.ReplicaSet
		}
		if opt.RetryWrites != nil {
			c.RetryWrites = opt.RetryWrites
		}
		if opt.RetryReads != nil {
			c.RetryReads = opt.RetryReads
		}
		if opt.ServerSelectionTimeout != nil {
			c.ServerSelectionTimeout = opt.ServerSelectionTimeout
		}
		if opt.Direct != nil {
			c.Direct = opt.Direct
		}
		if opt.SocketTimeout != nil {
			c.SocketTimeout = opt.SocketTimeout
		}
		if opt.TLSConfig != nil {
			c.TLSConfig = opt.TLSConfig
		}
		if opt.WriteConcern != nil {
			c.WriteConcern = opt.WriteConcern
		}
		if opt.ZlibLevel != nil {
			c.ZlibLevel = opt.ZlibLevel
		}
		if opt.ZstdLevel != nil {
			c.ZstdLevel = opt.ZstdLevel
		}
		if opt.AutoEncryptionOptions != nil {
			c.AutoEncryptionOptions = opt.AutoEncryptionOptions
		}
		if opt.Deployment != nil {
			c.Deployment = opt.Deployment
		}
		if opt.DisableOCSPEndpointCheck != nil {
			c.DisableOCSPEndpointCheck = opt.DisableOCSPEndpointCheck
		}
		if opt.err != nil {
			c.err = opt.err
		}

	}

	return c
}

// addCACertFromFile adds a root CA certificate to the configuration given a path
// to the containing file.
func addCACertFromFile(cfg *tls.Config, file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if cfg.RootCAs == nil {
		cfg.RootCAs = x509.NewCertPool()
	}
	if !cfg.RootCAs.AppendCertsFromPEM(data) {
		return errors.New("the specified CA file does not contain any valid certificates")
	}

	return nil
}

func addClientCertFromSeparateFiles(cfg *tls.Config, keyFile, certFile, keyPassword string) (string, error) {
	keyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return "", err
	}
	certData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return "", err
	}

	data := append(keyData, '\n')
	data = append(data, certData...)
	return addClientCertFromBytes(cfg, data, keyPassword)
}

func addClientCertFromConcatenatedFile(cfg *tls.Config, certKeyFile, keyPassword string) (string, error) {
	data, err := ioutil.ReadFile(certKeyFile)
	if err != nil {
		return "", err
	}

	return addClientCertFromBytes(cfg, data, keyPassword)
}

// addClientCertFromBytes adds a client certificate to the configuration given a path to the
// containing file and returns the certificate's subject name.
func addClientCertFromBytes(cfg *tls.Config, data []byte, keyPasswd string) (string, error) {
	var currentBlock *pem.Block
	var certBlock, certDecodedBlock, keyBlock []byte

	remaining := data
	start := 0
	for {
		currentBlock, remaining = pem.Decode(remaining)
		if currentBlock == nil {
			break
		}

		if currentBlock.Type == "CERTIFICATE" {
			certBlock = data[start : len(data)-len(remaining)]
			certDecodedBlock = currentBlock.Bytes
			start += len(certBlock)
		} else if strings.HasSuffix(currentBlock.Type, "PRIVATE KEY") {
			if keyPasswd != "" && x509.IsEncryptedPEMBlock(currentBlock) {
				var encoded bytes.Buffer
				buf, err := x509.DecryptPEMBlock(currentBlock, []byte(keyPasswd))
				if err != nil {
					return "", err
				}

				pem.Encode(&encoded, &pem.Block{Type: currentBlock.Type, Bytes: buf})
				keyBlock = encoded.Bytes()
				start = len(data) - len(remaining)
			} else {
				keyBlock = data[start : len(data)-len(remaining)]
				start += len(keyBlock)
			}
		}
	}
	if len(certBlock) == 0 {
		return "", fmt.Errorf("failed to find CERTIFICATE")
	}
	if len(keyBlock) == 0 {
		return "", fmt.Errorf("failed to find PRIVATE KEY")
	}

	cert, err := tls.X509KeyPair(certBlock, keyBlock)
	if err != nil {
		return "", err
	}

	cfg.Certificates = append(cfg.Certificates, cert)

	// The documentation for the tls.X509KeyPair indicates that the Leaf certificate is not
	// retained.
	crt, err := x509.ParseCertificate(certDecodedBlock)
	if err != nil {
		return "", err
	}

	return x509CertSubject(crt), nil
}

func stringSliceContains(source []string, target string) bool {
	for _, str := range source {
		if str == target {
			return true
		}
	}
	return false
}

// create a username for x509 authentication from an x509 certificate subject.
func extractX509UsernameFromSubject(subject string) string {
	// the Go x509 package gives the subject with the pairs in the reverse order from what we want.
	pairs := strings.Split(subject, ",")
	for left, right := 0, len(pairs)-1; left < right; left, right = left+1, right-1 {
		pairs[left], pairs[right] = pairs[right], pairs[left]
	}

	return strings.Join(pairs, ",")
}
