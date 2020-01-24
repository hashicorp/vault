package mongodb

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/errwrap"
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
	Username      string `json:"username" structs:"username" mapstructure:"username"`
	Password      string `json:"password" structs:"password" mapstructure:"password"`

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

func (c *mongoDBConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := c.Init(ctx, conf, verifyConnection)
	return err
}

// Initialize parses connection configuration.
func (c *mongoDBConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	c.Lock()
	defer c.Unlock()

	c.RawConfig = conf

	err := mapstructure.WeakDecode(conf, c)
	if err != nil {
		return nil, err
	}

	if len(c.ConnectionURL) == 0 {
		return nil, fmt.Errorf("connection_url cannot be empty")
	}

	c.ConnectionURL = dbutil.QueryHelper(c.ConnectionURL, map[string]string{
		"username": c.Username,
		"password": c.Password,
	})

	if c.WriteConcern != "" {
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

		c.clientOptions = &options.ClientOptions{
			WriteConcern: writeConcern,
		}
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	if verifyConnection {
		if _, err := c.Connection(ctx); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}

		if err := c.client.Ping(ctx, readpref.Primary()); err != nil {
			return nil, errwrap.Wrapf("error verifying connection: {{err}}", err)
		}
	}

	return conf, nil
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

	if c.clientOptions == nil {
		c.clientOptions = options.Client()
	}
	c.clientOptions.SetSocketTimeout(1 * time.Minute)
	c.clientOptions.SetConnectTimeout(1 * time.Minute)

	var err error
	c.client, err = mongo.Connect(ctx, c.clientOptions.ApplyURI(c.ConnectionURL))
	if err != nil {
		return nil, err
	}
	return c.client, nil
}

// Close terminates the database connection.
func (c *mongoDBConnectionProducer) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Minute)
		defer cancel()
		if err := c.client.Disconnect(ctx); err != nil {
			return err
		}
	}

	c.client = nil

	return nil
}

func (c *mongoDBConnectionProducer) secretValues() map[string]interface{} {
	return map[string]interface{}{
		c.Password: "[password]",
	}
}
