package config

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

var (
	ctx     = context.Background()
	storage = &logical.InmemStorage{}
)

func TestCacheReader(t *testing.T) {

	m, err := NewManager(ctx, &logical.BackendConfig{
		Logger:      hclog.NewNullLogger(),
		StorageView: storage,
	})
	if err != nil {
		t.Error(err)
	}

	var configReader Reader
	configReader = m

	// we should start with no config
	config, err := configReader.Config(ctx, storage)
	if err != nil {
		t.Error(err)
	}
	if config != nil {
		t.Error("config should initially be nil because it's unset as of yet")
	}

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      BackendPath,
		Storage:   storage,
	}

	fieldData := &framework.FieldData{
		Schema: m.Path().Fields,
		Raw: map[string]interface{}{
			"username": "tester",
			"password": "pa$$w0rd",
			"urls":     "ldap://138.91.247.105",
			"dn":       "example,com",
		},
	}

	_, err = m.update(ctx, req, fieldData)
	if err != nil {
		t.Error(err)
	}

	// now that we've updated the config, we should be able to read it
	config, err = configReader.Config(ctx, storage)
	if err != nil {
		t.Error(err)
	}

	if config.ADConf.Username != "tester" {
		t.Error("returned config is not populated as expected")
	}

	req = &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      BackendPath,
		Storage:   storage,
	}

	_, err = m.delete(ctx, req, nil)
	if err != nil {
		t.Error(err)
	}

	// now that we've deleted the config, it should be unset again
	config, err = configReader.Config(ctx, storage)
	if err != nil {
		t.Error(err)
	}
	if config != nil {
		t.Error("config should again be nil because after it's been deleted, it's again unset by the user")
	}
}
