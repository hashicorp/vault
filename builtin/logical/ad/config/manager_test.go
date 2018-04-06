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

	config, err := configReader.Config(ctx, storage)
	if err != nil {
		t.FailNow()
	}
	if config.ADConf.Username != "" {
		t.FailNow()
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

	config, err = configReader.Config(ctx, storage)
	if err != nil {
		t.Error(err)
	}

	if config.ADConf.Username != "tester" {
		t.FailNow()
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

	config, err = configReader.Config(ctx, storage)
	if err != nil {
		t.Error(err)
	}
	if config.ADConf.Username != "" {
		t.FailNow()
	}
}
