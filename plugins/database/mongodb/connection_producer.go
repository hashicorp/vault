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

	"github.com/hashicorp/errwrap"
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

// Connection creates or returns an existing a database connection. If the session fails
// on a ping check, the session will be closed and then re-created.
// This method does not lock the mutex and it is intended that this is the callers
// responsibility.
func (c *mongoDBConnectionProducer) Connection(ctx context.Context) (interface{}, error) {
	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	if c.client != nil {
		if err := c.client.Ping(ctx, readpref.Primary()); err == nil {
			return c.client, nil
		}
		// Ignore error on purpose since we want to re-create a session
		_ = c.client.Disconnect(ctx)
	}

	connURL := c.getConnectionURL()
	client, err := createClient(ctx, connURL, c.clientOptions)
	if err != nil {
		return nil, err
	}
	c.client = client
	return c.client, nil
}

func createClient(ctx context.Context, connURL string, clientOptions *options.ClientOptions) (client *mongo.Client, err error) {
	if clientOptions == nil {
		clientOptions = options.Client()
	}
	clientOptions.SetSocketTimeout(1 * time.Minute)
	clientOptions.SetConnectTimeout(1 * time.Minute)

	client, err = mongo.Connect(ctx, options.MergeClientOptions(options.Client().ApplyURI(connURL), clientOptions))
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
		return nil, errwrap.Wrapf("error unmarshalling write_concern: {{err}}", err)
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
