package mongodbatlas

import (
	"context"
	"sync"

	"github.com/Sectorbob/mlab-ns2/gae/ns/digest"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"go.mongodb.org/atlas/mongodbatlas"
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

func (c *mongoDBAtlasConnectionProducer) secretValues() map[string]string {
	return map[string]string{
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
	client.UserAgent = useragent.String()

	c.client = client

	return c.client, nil
}
