package config

import (
	"context"

	"github.com/hashicorp/vault/helper/activedirectory"
	"github.com/hashicorp/vault/logical"
)

const StorageKey = "config"

func readConfig(ctx context.Context, storage logical.Storage) (*EngineConf, error) {

	config := &EngineConf{&PasswordConf{}, &activedirectory.Configuration{}}

	entry, err := storage.Get(ctx, StorageKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	if err := entry.DecodeJSON(config); err != nil {
		return nil, err
	}
	return config, nil
}

func writeConfig(ctx context.Context, storage logical.Storage, config *EngineConf) error {
	entry, err := logical.StorageEntryJSON(StorageKey, config)
	if err != nil {
		return err
	}
	return storage.Put(ctx, entry)
}

func deleteConfig(ctx context.Context, storage logical.Storage) error {
	return storage.Delete(ctx, StorageKey)
}
