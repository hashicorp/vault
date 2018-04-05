package config

import (
	"context"

	"github.com/hashicorp/vault/logical"
)

const storageKey = "config"

func readConfig(ctx context.Context, storage logical.Storage) (*EngineConf, error) {

	engineConf := newUnsetEngineConf()

	entry, err := storage.Get(ctx, storageKey)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return engineConf, nil
	}

	if err := entry.DecodeJSON(engineConf); err != nil {
		return nil, err
	}
	return engineConf, nil
}

func writeConfig(ctx context.Context, storage logical.Storage, config *EngineConf) error {
	entry, err := logical.StorageEntryJSON(storageKey, config)
	if err != nil {
		return err
	}
	return storage.Put(ctx, entry)
}

func deleteConfig(ctx context.Context, storage logical.Storage) error {
	return storage.Delete(ctx, storageKey)
}
