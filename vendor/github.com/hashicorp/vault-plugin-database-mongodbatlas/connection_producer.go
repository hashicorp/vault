// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodbatlas

import (
	"context"
	"errors"
	"sync"

	dbplugin "github.com/hashicorp/vault/sdk/database/dbplugin/v5"
	"github.com/hashicorp/vault/sdk/database/helper/connutil"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/mongodb-forks/digest"
	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	// TODO: The Vault version used in the user agent string is hard coded until
	//  it's possible for database plugins to use the system view to obtain correct
	//  Vault version information via the plugin environment.
	userAgentVaultVersion = "1.13.0"
	userAgentPluginName   = "database-mongodbatlas"
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

	// TODO: Obtain the plugin environment from the system view.
	env := &logical.PluginEnvironment{
		VaultVersion: userAgentVaultVersion,
	}
	client.UserAgent = useragent.PluginString(env, userAgentPluginName)

	c.client = client

	return c.client, nil
}

func (m *mongoDBAtlasConnectionProducer) Initialize(ctx context.Context, req dbplugin.InitializeRequest) error {
	m.Lock()
	defer m.Unlock()

	m.RawConfig = req.Config

	err := mapstructure.WeakDecode(req.Config, m)
	if err != nil {
		return err
	}

	if len(m.PublicKey) == 0 {
		return errors.New("public Key is not set")
	}

	if len(m.PrivateKey) == 0 {
		return errors.New("private Key is not set")
	}

	// Set initialized to true at this point since all fields are set,
	// and the connection can be established at a later time.
	m.Initialized = true

	return nil
}
