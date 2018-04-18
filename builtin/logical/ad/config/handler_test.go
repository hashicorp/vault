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

	h := Handler(hclog.NewNullLogger())

	var configReader Reader
	configReader = h

	// we should start with no config
	config, err := configReader.Read(ctx, storage)
	if err != nil {
		_, ok := err.(*Unset)
		if !ok {
			t.Fatal(err)
		}
	}

	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      InboundPath,
		Storage:   storage,
	}

	fieldData := &framework.FieldData{
		Schema: h.Path().Fields,
		Raw: map[string]interface{}{
			"username": "tester",
			"password": "pa$$w0rd",
			"urls":     "ldap://138.91.247.105",
			"dn":       "example,com",
		},
	}

	_, err = h.updateOperation(ctx, req, fieldData)
	if err != nil {
		t.Fatal(err)
	}

	// now that we've updated the config, we should be able to readOperation it
	config, err = configReader.Read(ctx, storage)
	if err != nil {
		t.Fatal(err)
	}

	if config.ADConf.Username != "tester" {
		t.Fatal("returned config is not populated as expected")
	}

	req = &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      InboundPath,
		Storage:   storage,
	}

	_, err = h.deleteOperation(ctx, req, nil)
	if err != nil {
		t.Fatal(err)
	}

	// now that we've deleted the config, it should be unset again
	config, err = configReader.Read(ctx, storage)
	if err != nil {
		_, ok := err.(*Unset)
		if !ok {
			t.Fatal(err)
		}
	}
}
