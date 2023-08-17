// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mongodb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/database/helper/dbutil"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// mongoDBConnectionProducer implements ConnectionProducer and provides an
// interface for databases to make connections.
type mongoDBConnectionProducer struct {
	ConnectionURL string `json:"connection_url" structs:"connection_url" mapstructure:"connection_url"`
	WriteConcern  string `json:"write_concern" structs:"write_concern" mapstructure:"write_concern"`

	Username string `json:"username" structs:"username" mapstructure:"username"`
	Password string `json:"password" structs:"password" mapstructure:"password"`

	TLSCertificateKeyData []byte `json:"tls_certificate_key" structs:"-" mapstructure:"tls_certificate_key"`
	TLSCAData             []byte `json:"tls_ca"              structs:"-" mapstructure:"tls_ca"`

	SocketTimeout          time.Duration `json:"socket_timeout"           structs:"-" mapstructure:"socket_timeout"`
	ConnectTimeout         time.Duration `json:"connect_timeout"          structs:"-" mapstructure:"connect_timeout"`
	ServerSelectionTimeout time.Duration `json:"server_selection_timeout" structs:"-" mapstructure:"server_selection_timeout"`

	Initialized   bool
	RawConfig     map[string]interface{}
	Type          string
	clientOptions *options.ClientOptions
	client        *mongo.Client
	sync.Mutex
}

// writeConcern defines the write concern options
type writeConcern struct {
	W        int    // Min # of servers to ack before success
	WMode    string // Write mode for MongoDB 2.0+ (e.g. "majority")
	WTimeout int    // Milliseconds to wait for W before timing out
	FSync    bool   // DEPRECATED: Is now handled by J. See: https://jira.mongodb.org/browse/CXX-910
	J        bool   // Sync via the journal if present
}

func (c *mongoDBConnectionProducer) loadConfig(cfg map[string]interface{}) error {
	err := mapstructure.WeakDecode(cfg, c)
	if err != nil {
		return err
	}

	if len(c.ConnectionURL) == 0 {
		return fmt.Errorf("connection_url cannot be empty")
	}

	if c.SocketTimeout < 0 {
		return fmt.Errorf("socket_timeout must be >= 0")
	}
	if c.ConnectTimeout < 0 {
		return fmt.Errorf("connect_timeout must be >= 0")
	}
	if c.ServerSelectionTimeout < 0 {
		return fmt.Errorf("server_selection_timeout must be >= 0")
	}

	opts, err := c.makeClientOpts()
	if err != nil {
		return err
	}

	c.clientOptions = opts

	return nil
}

// Connection creates or returns an existing a database connection. If the session fails
// on a ping check, the session will be closed and then re-created.
// This method does locks the mutex on its own.
func (c *mongoDBConnectionProducer) Connection(ctx context.Context) (*mongo.Client, error) {
	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if c.client != nil {
		if err := c.client.Ping(ctx, readpref.Primary()); err == nil {
			return c.client, nil
		}
		// Ignore error on purpose since we want to re-create a session
		_ = c.client.Disconnect(ctx)
	}

	client, err := c.createClient(ctx)
	if err != nil {
		return nil, err
	}
	c.client = client
	return c.client, nil
}

func (c *mongoDBConnectionProducer) createClient(ctx context.Context) (client *mongo.Client, err error) {
	if !c.Initialized {
		return nil, fmt.Errorf("failed to create client: connection producer is not initialized")
	}
	if c.clientOptions == nil {
		return nil, fmt.Errorf("missing client options")
	}
	client, err = mongo.Connect(ctx, options.MergeClientOptions(options.Client().ApplyURI(c.getConnectionURL()), c.clientOptions))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Close terminates the database connection.
func (c *mongoDBConnectionProducer) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		if err := c.client.Disconnect(ctx); err != nil {
			return err
		}
	}

	c.client = nil

	return nil
}

func (c *mongoDBConnectionProducer) secretValues() map[string]string {
	return map[string]string{
		c.Password: "[password]",
	}
}

func (c *mongoDBConnectionProducer) getConnectionURL() (connURL string) {
	connURL = dbutil.QueryHelper(c.ConnectionURL, map[string]string{
		"username": c.Username,
		"password": c.Password,
	})
	return connURL
}

func (c *mongoDBConnectionProducer) makeClientOpts() (*options.ClientOptions, error) {
	writeOpts, err := c.getWriteConcern()
	if err != nil {
		return nil, err
	}

	authOpts, err := c.getTLSAuth()
	if err != nil {
		return nil, err
	}

	timeoutOpts, err := c.timeoutOpts()
	if err != nil {
		return nil, err
	}

	opts := options.MergeClientOptions(writeOpts, authOpts, timeoutOpts)
	return opts, nil
}

func (c *mongoDBConnectionProducer) getWriteConcern() (opts *options.ClientOptions, err error) {
	if c.WriteConcern == "" {
		return nil, nil
	}

	input := c.WriteConcern

	// Try to base64 decode the input. If successful, consider the decoded
	// value as input.
	inputBytes, err := base64.StdEncoding.DecodeString(input)
	if err == nil {
		input = string(inputBytes)
	}

	concern := &writeConcern{}
	err = json.Unmarshal([]byte(input), concern)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling write_concern: %w", err)
	}

	// Translate write concern to mongo options
	var w writeconcern.Option
	switch {
	case concern.W != 0:
		w = writeconcern.W(concern.W)
	case concern.WMode != "":
		w = writeconcern.WTagSet(concern.WMode)
	default:
		w = writeconcern.WMajority()
	}

	var j writeconcern.Option
	switch {
	case concern.FSync:
		j = writeconcern.J(concern.FSync)
	case concern.J:
		j = writeconcern.J(concern.J)
	default:
		j = writeconcern.J(false)
	}

	writeConcern := writeconcern.New(
		w,
		j,
		writeconcern.WTimeout(time.Duration(concern.WTimeout)*time.Millisecond))

	opts = options.Client()
	opts.SetWriteConcern(writeConcern)
	return opts, nil
}

func (c *mongoDBConnectionProducer) getTLSAuth() (opts *options.ClientOptions, err error) {
	if len(c.TLSCAData) == 0 && len(c.TLSCertificateKeyData) == 0 {
		return nil, nil
	}

	opts = options.Client()

	tlsConfig := &tls.Config{}

	if len(c.TLSCAData) > 0 {
		tlsConfig.RootCAs = x509.NewCertPool()

		ok := tlsConfig.RootCAs.AppendCertsFromPEM(c.TLSCAData)
		if !ok {
			return nil, fmt.Errorf("failed to append CA to client options")
		}
	}

	if len(c.TLSCertificateKeyData) > 0 {
		certificate, err := tls.X509KeyPair(c.TLSCertificateKeyData, c.TLSCertificateKeyData)
		if err != nil {
			return nil, fmt.Errorf("unable to load tls_certificate_key_data: %w", err)
		}

		opts.SetAuth(options.Credential{
			AuthMechanism: "MONGODB-X509",
			Username:      c.Username,
		})

		tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
	}

	opts.SetTLSConfig(tlsConfig)
	return opts, nil
}

func (c *mongoDBConnectionProducer) timeoutOpts() (opts *options.ClientOptions, err error) {
	opts = options.Client()

	if c.SocketTimeout < 0 {
		return nil, fmt.Errorf("socket_timeout must be >= 0")
	}

	if c.SocketTimeout == 0 {
		opts.SetSocketTimeout(1 * time.Minute)
	} else {
		opts.SetSocketTimeout(c.SocketTimeout)
	}

	if c.ConnectTimeout == 0 {
		opts.SetConnectTimeout(1 * time.Minute)
	} else {
		opts.SetConnectTimeout(c.ConnectTimeout)
	}

	if c.ServerSelectionTimeout != 0 {
		opts.SetServerSelectionTimeout(c.ServerSelectionTimeout)
	}

	return opts, nil
}
