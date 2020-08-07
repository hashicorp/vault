package mongodbatlas

import (
	"context"
	"errors"
	"sync"

	"github.com/Sectorbob/mlab-ns2/gae/ns/digest"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/mitchellh/mapstructure"
	"github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

type mongoDBAtlasConnectionProducer struct {
	PublicKey  string `json:"public_key" structs:"public_key" mapstructure:"public_key"`
	PrivateKey string `json:"private_key" structs:"private_key" mapstructure:"private_key"`
	ProjectID  string `json:"project_id" structs:"project_id" mapstructure:"project_id"`

	Initialized bool
	RawConfig   map[string]interface{}
	Type        string
	client      *mongodbatlas.Client
	sync.Mutex
}

func (c *mongoDBAtlasConnectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := c.Init(ctx, conf, verifyConnection)
	return err
}

// Initialize parses connection configuration.
func (c *mongoDBAtlasConnectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	c.Lock()
	defer c.Unlock()

	err := mapstructure.WeakDecode(conf, c)
	if err != nil {
		return nil, err
	}

	if len(c.PublicKey) == 0 {
		return nil, errors.New("public Key is not set")
	}

	if len(c.PrivateKey) == 0 {
		return nil, errors.New("private Key is not set")
	}

	c.RawConfig = conf

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	c.Initialized = true

	return conf, nil
}

func (c *mongoDBAtlasConnectionProducer) secretValues() map[string]interface{} {
	return map[string]interface{}{
		c.PrivateKey: "[private_key]",
	}
}

// Close terminates the database connection.
func (c *mongoDBAtlasConnectionProducer) Close() error {
	c.Lock()
	defer c.Unlock()

	c.client = nil

	return nil
}

func (c *mongoDBAtlasConnectionProducer) Connection(_ context.Context) (interface{}, error) {
	// This is intentionally not grabbing the lock since the calling functions (e.g. CreateUser)
	// are claiming it. (The locking patterns could be refactored to be more consistent/clear.)

	if !c.Initialized {
		return nil, connutil.ErrNotInitialized
	}

	if c.client != nil {
		return c.client, nil
	}

	transport := digest.NewTransport(c.PublicKey, c.PrivateKey)
	cl, err := transport.Client()
	if err != nil {
		return nil, err
	}

	client, err := mongodbatlas.New(cl)
	if err != nil {
		return nil, err
	}

	c.client = client

	return c.client, nil
}
