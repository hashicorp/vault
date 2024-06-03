// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	StorageKeyConfig = "config/keys"
)

type KeyConfigEntry struct {
	DefaultKeyId KeyID `json:"default"`
}

func SetKeysConfig(ctx context.Context, s logical.Storage, config *KeyConfigEntry) error {
	json, err := logical.StorageEntryJSON(StorageKeyConfig, config)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func GetKeysConfig(ctx context.Context, s logical.Storage) (*KeyConfigEntry, error) {
	entry, err := s.Get(ctx, StorageKeyConfig)
	if err != nil {
		return nil, err
	}

	keyConfig := &KeyConfigEntry{}
	if entry != nil {
		if err := entry.DecodeJSON(keyConfig); err != nil {
			return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode key configuration: %v", err)}
		}
	}

	return keyConfig, nil
}
