// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package issuing

import (
	"context"
	"crypto"
	"encoding/pem"
	"fmt"

	"github.com/hashicorp/vault/builtin/logical/pki/managed_key"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	KeyPrefix      = "config/key/"
	KeyRefNotFound = KeyID("not-found")
)

type KeyID string

func (p KeyID) String() string {
	return string(p)
}

type KeyEntry struct {
	ID             KeyID                   `json:"id"`
	Name           string                  `json:"name"`
	PrivateKeyType certutil.PrivateKeyType `json:"private_key_type"`
	PrivateKey     string                  `json:"private_key"`
}

func (e KeyEntry) IsManagedPrivateKey() bool {
	return e.PrivateKeyType == certutil.ManagedPrivateKey
}

func ListKeys(ctx context.Context, s logical.Storage) ([]KeyID, error) {
	strList, err := s.List(ctx, KeyPrefix)
	if err != nil {
		return nil, err
	}

	keyIds := make([]KeyID, 0, len(strList))
	for _, entry := range strList {
		keyIds = append(keyIds, KeyID(entry))
	}

	return keyIds, nil
}

func FetchKeyById(ctx context.Context, s logical.Storage, keyId KeyID) (*KeyEntry, error) {
	if len(keyId) == 0 {
		return nil, errutil.InternalError{Err: "unable to fetch pki key: empty key identifier"}
	}

	entry, err := s.Get(ctx, KeyPrefix+keyId.String())
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch pki key: %v", err)}
	}
	if entry == nil {
		return nil, errutil.UserError{Err: fmt.Sprintf("pki key id %s does not exist", keyId.String())}
	}

	var key KeyEntry
	if err := entry.DecodeJSON(&key); err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("unable to decode pki key with id %s: %v", keyId.String(), err)}
	}

	return &key, nil
}

func WriteKey(ctx context.Context, s logical.Storage, key KeyEntry) error {
	keyId := key.ID

	json, err := logical.StorageEntryJSON(KeyPrefix+keyId.String(), key)
	if err != nil {
		return err
	}

	return s.Put(ctx, json)
}

func DeleteKey(ctx context.Context, s logical.Storage, id KeyID) (bool, error) {
	config, err := GetKeysConfig(ctx, s)
	if err != nil {
		return false, err
	}

	wasDefault := false
	if config.DefaultKeyId == id {
		wasDefault = true
		config.DefaultKeyId = KeyID("")
		if err := SetKeysConfig(ctx, s, config); err != nil {
			return wasDefault, err
		}
	}

	return wasDefault, s.Delete(ctx, KeyPrefix+id.String())
}

func ResolveKeyReference(ctx context.Context, s logical.Storage, reference string) (KeyID, error) {
	if reference == DefaultRef {
		// Handle fetching the default key.
		config, err := GetKeysConfig(ctx, s)
		if err != nil {
			return KeyID("config-error"), err
		}
		if len(config.DefaultKeyId) == 0 {
			return KeyRefNotFound, fmt.Errorf("no default key currently configured")
		}

		return config.DefaultKeyId, nil
	}

	// Lookup by a direct get first to see if our reference is an ID, this is quick and cached.
	if len(reference) == uuidLength {
		entry, err := s.Get(ctx, KeyPrefix+reference)
		if err != nil {
			return KeyID("key-read"), err
		}
		if entry != nil {
			return KeyID(reference), nil
		}
	}

	// ... than to pull all keys from storage.
	keys, err := ListKeys(ctx, s)
	if err != nil {
		return KeyID("list-error"), err
	}
	for _, keyId := range keys {
		key, err := FetchKeyById(ctx, s, keyId)
		if err != nil {
			return KeyID("key-read"), err
		}

		if key.Name == reference {
			return key.ID, nil
		}
	}

	// Otherwise, we must not have found the key.
	return KeyRefNotFound, errutil.UserError{Err: fmt.Sprintf("unable to find PKI key for reference: %v", reference)}
}

func GetManagedKeyUUID(key *KeyEntry) (managed_key.UUIDKey, error) {
	if !key.IsManagedPrivateKey() {
		return "", errutil.InternalError{Err: "getManagedKeyUUID called on a key id %s (%s) "}
	}
	return managed_key.ExtractManagedKeyId([]byte(key.PrivateKey))
}

func GetSignerFromKeyEntry(ctx context.Context, mkv managed_key.PkiManagedKeyView, keyEntry *KeyEntry) (crypto.Signer, certutil.PrivateKeyType, error) {
	if keyEntry.PrivateKeyType == certutil.UnknownPrivateKey {
		return nil, certutil.UnknownPrivateKey, fmt.Errorf("unsupported unknown private key type for key: %s (%s)", keyEntry.ID, keyEntry.Name)
	}

	if keyEntry.IsManagedPrivateKey() {
		managedKeyId, err := GetManagedKeyUUID(keyEntry)
		if err != nil {
			return nil, certutil.UnknownPrivateKey, fmt.Errorf("unable to get managed key uuid: %w", err)
		}
		bundle, actualKeyType, err := managed_key.CreateKmsKeyBundle(ctx, mkv, managedKeyId)
		if err != nil {
			return nil, certutil.UnknownPrivateKey, fmt.Errorf("failed to create kms key bundle from managed key uuid %s: %w", managedKeyId, err)
		}

		// The bundle's PrivateKeyType value is set to a ManagedKeyType so use the actual key type value
		return bundle.PrivateKey, actualKeyType, nil
	}

	pemBlock, _ := pem.Decode([]byte(keyEntry.PrivateKey))
	if pemBlock == nil {
		return nil, certutil.UnknownPrivateKey, fmt.Errorf("no data found in PEM block")
	}

	signer, _, err := certutil.ParseDERKey(pemBlock.Bytes)
	if err != nil {
		return nil, certutil.UnknownPrivateKey, fmt.Errorf("failed to parse PEM block: %w", err)
	}
	return signer, keyEntry.PrivateKeyType, nil
}
