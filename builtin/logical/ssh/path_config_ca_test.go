// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestSSH_ConfigCAStorageUpgrade(t *testing.T) {
	var err error

	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// Store at an older path
	err = config.StorageView.Put(context.Background(), &logical.StorageEntry{
		Key:   caPrivateKeyStoragePathDeprecated,
		Value: []byte(testCAPrivateKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reading it should return the key as well as upgrade the storage path
	privateKeyEntry, err := readStoredKey(context.Background(), config.StorageView, caPrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	if privateKeyEntry == nil || privateKeyEntry.Key == "" {
		t.Fatalf("failed to read the stored private key")
	}

	entry, err := config.StorageView.Get(context.Background(), caPrivateKeyStoragePathDeprecated)
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatalf("bad: expected a nil entry after upgrade")
	}

	entry, err = config.StorageView.Get(context.Background(), caPrivateKeyStoragePath)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatalf("bad: expected a non-nil entry after upgrade")
	}

	// Store at an older path
	err = config.StorageView.Put(context.Background(), &logical.StorageEntry{
		Key:   caPublicKeyStoragePathDeprecated,
		Value: []byte(testCAPublicKey),
	})
	if err != nil {
		t.Fatal(err)
	}

	// Reading it should return the key as well as upgrade the storage path
	publicKeyEntry, err := readStoredKey(context.Background(), config.StorageView, caPublicKey)
	if err != nil {
		t.Fatal(err)
	}
	if publicKeyEntry == nil || publicKeyEntry.Key == "" {
		t.Fatalf("failed to read the stored public key")
	}

	entry, err = config.StorageView.Get(context.Background(), caPublicKeyStoragePathDeprecated)
	if err != nil {
		t.Fatal(err)
	}
	if entry != nil {
		t.Fatalf("bad: expected a nil entry after upgrade")
	}

	entry, err = config.StorageView.Get(context.Background(), caPublicKeyStoragePath)
	if err != nil {
		t.Fatal(err)
	}
	if entry == nil {
		t.Fatalf("bad: expected a non-nil entry after upgrade")
	}
}

func TestSSH_ConfigCAUpdateDelete(t *testing.T) {
	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	caReq := &logical.Request{
		Path:      "config/ca",
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
	}

	// Auto-generate the keys
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Fail to overwrite it
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error, got %#v", *resp)
	}

	caReq.Operation = logical.DeleteOperation
	// Delete the configured keys
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	caReq.Operation = logical.UpdateOperation
	caReq.Data = map[string]interface{}{
		"public_key":  testCAPublicKey,
		"private_key": testCAPrivateKey,
	}

	// Successfully create a new one
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Fail to overwrite it
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.IsError() {
		t.Fatalf("expected an error, got %#v", *resp)
	}

	caReq.Operation = logical.DeleteOperation
	// Delete the configured keys
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	caReq.Operation = logical.UpdateOperation
	caReq.Data = nil

	// Successfully create a new one
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}

	// Delete the configured keys
	caReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v, resp:%v", err, resp)
	}
}

func createDeleteHelper(t *testing.T, b logical.Backend, config *logical.BackendConfig, index int, keyType string, keyBits int) {
	// Check that we can create a new key of the specified type
	caReq := &logical.Request{
		Path:      "config/ca",
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
	}
	caReq.Data = map[string]interface{}{
		"generate_signing_key": true,
		"key_type":             keyType,
		"key_bits":             keyBits,
	}
	resp, err := b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad case %v: err: %v, resp: %v", index, err, resp)
	}
	if !strings.Contains(resp.Data["public_key"].(string), caReq.Data["key_type"].(string)) {
		t.Fatalf("bad case %v: expected public key of type %v but was %v", index, caReq.Data["key_type"], resp.Data["public_key"])
	}

	issueOptions := map[string]interface{}{
		"public_key": testCAPublicKeyEd25519,
	}
	issueReq := &logical.Request{
		Path:      "sign/ca-issuance",
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Data:      issueOptions,
	}
	resp, err = b.HandleRequest(context.Background(), issueReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad case %v: err: %v, resp: %v", index, err, resp)
	}

	// Delete the configured keys
	caReq.Operation = logical.DeleteOperation
	resp, err = b.HandleRequest(context.Background(), caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad case %v: err: %v, resp: %v", index, err, resp)
	}
}

func TestSSH_ConfigCAKeyTypes(t *testing.T) {
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatalf("Cannot create backend: %s", err)
	}

	cases := []struct {
		keyType string
		keyBits int
	}{
		{"ssh-rsa", 2048},
		{"ssh-rsa", 4096},
		{"ssh-rsa", 0},
		{"rsa", 2048},
		{"rsa", 4096},
		{"ecdsa-sha2-nistp256", 0},
		{"ecdsa-sha2-nistp384", 0},
		{"ecdsa-sha2-nistp521", 0},
		{"ec", 256},
		{"ec", 384},
		{"ec", 521},
		{"ec", 0},
		{"ssh-ed25519", 0},
		{"ed25519", 0},
	}

	// Create a role for ssh signing.
	roleOptions := map[string]interface{}{
		"allow_user_certificates": true,
		"allowed_users":           "*",
		"key_type":                "ca",
		"ttl":                     "30s",
		"not_before_duration":     "2h",
		"allow_empty_principals":  true,
	}
	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/ca-issuance",
		Data:      roleOptions,
		Storage:   config.StorageView,
	}
	_, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil {
		t.Fatalf("Cannot create role to issue against: %s", err)
	}

	for index, scenario := range cases {
		createDeleteHelper(t, b, config, index, scenario.keyType, scenario.keyBits)
	}
}

func TestReadManagedKey(t *testing.T) {
	t.Parallel()

	storage := &logical.InmemStorage{}
	entry, err := readManagedKey(context.Background(), storage)
	if err != nil {
		t.Fatalf("error reading managed key: %s", err)
	}

	if entry != nil {
		t.Fatal("expected nil, but got a non-nil return")
	}

	err = writeKey(context.Background(), storage, caManagedKeyStoragePath, "test-managed-key")
	if err != nil {
		t.Fatalf("error writing test key: %s", err)
	}

	entry, err = readManagedKey(context.Background(), storage)
	if err != nil {
		t.Fatalf("error reading managed key: %s", err)
	}

	if entry == nil {
		t.Fatal("unexpected nil entry")
	}

	if entry.PublicKey != "test-managed-key" {
		t.Fatalf("key value mismatch: expected %s, got %s", "test-managed-key", entry.PublicKey)
	}
}

func TestReadStoredKey(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		privateKeyStoragePath string
		publicKeyStoragePath  string
		publicKey             string
		privateKey            string
	}{
		"stored-keys-configured": {
			privateKeyStoragePath: caPrivateKeyStoragePath,
			publicKeyStoragePath:  caPublicKeyStoragePath,
			publicKey:             testCAPublicKey,
			privateKey:            testCAPrivateKey,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			storage := &logical.InmemStorage{}

			if err := writeKey(ctx, storage, tt.privateKeyStoragePath, tt.privateKey); err != nil {
				t.Fatalf("error writing private key: %s", err)
			}

			if err := writeKey(ctx, storage, tt.publicKeyStoragePath, tt.publicKey); err != nil {
				t.Fatalf("error writing public key: %s", err)
			}

			publicKeyEntry, err := readStoredKey(context.Background(), storage, caPublicKey)
			if err != nil {
				t.Fatalf("error reading public key: %s", err)
			}

			if publicKeyEntry.Key != tt.publicKey {
				t.Fatalf("returned key does not match: expected %s, got %s", tt.publicKey, publicKeyEntry.Key)
			}

			privateKeyEntry, err := readStoredKey(context.Background(), storage, caPrivateKey)
			if err != nil {
				t.Fatalf("error reading private key: %s", err)
			}

			if privateKeyEntry.Key != tt.privateKey {
				t.Fatalf("returned key does not match: expected %s, got %s", tt.privateKey, privateKeyEntry.Key)
			}
		})
	}
}

func TestGetCAPublicKey(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		publicKeyStoragePath string
		publicKey            string
	}{
		"stored-keys-configured": {
			publicKeyStoragePath: caPublicKeyStoragePath,
			publicKey:            testCAPublicKey,
		},
		"managed-key-configured": {
			publicKeyStoragePath: caManagedKeyStoragePath,
			publicKey:            testCAPublicKey,
		},
		"no-keys-configured": {},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			storage := &logical.InmemStorage{}
			err := writeKey(ctx, storage, tt.publicKeyStoragePath, tt.publicKey)
			if err != nil {
				t.Fatalf("error writing key: %s", err)
			}

			key, err := getCAPublicKey(ctx, storage)
			if err != nil {
				t.Fatalf("error retrieving public key: %s", err)
			}

			if key != tt.publicKey {
				t.Fatalf("key values do not match: expected %s, got %s", tt.publicKey, key)
			}
		})
	}
}

func TestCreateStoredKey(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		publicKey  string
		privateKey string
		expectErr  bool
	}{
		"both-keys-provided": {
			publicKey:  testCAPublicKey,
			privateKey: testCAPrivateKey,
		},
		"only-public-key": {
			publicKey: testCAPublicKey,
			expectErr: true,
		},
		"only-private-key": {
			privateKey: testCAPrivateKey,
			expectErr:  true,
		},
		"empty keys": {
			expectErr: true,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			storage := &logical.InmemStorage{}
			err := createStoredKey(context.Background(), storage, tt.publicKey, tt.privateKey)
			if err != nil && !tt.expectErr {
				t.Fatalf("unexpected error: %s", err)
			} else if err == nil && tt.expectErr {
				t.Fatal("expected error, got nil")
			}

			if !tt.expectErr {
				err = readKey(context.Background(), storage, caPublicKeyStoragePath)
				if err != nil {
					t.Fatalf("error reading public key: %s", err)
				}

				err = readKey(context.Background(), storage, caPrivateKeyStoragePath)
				if err != nil {
					t.Fatalf("error reading private key: %s", err)
				}
			}
		})
	}
}

func TestCAKeysConfigured(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		privateKeyStoragePath string
		publicKeyStoragePath  string
		publicKey             string
		privateKey            string
		expectedValue         bool
	}{
		"stored-keys-configured": {
			privateKeyStoragePath: caPrivateKeyStoragePath,
			publicKeyStoragePath:  caPublicKeyStoragePath,
			publicKey:             testCAPublicKey,
			privateKey:            testCAPrivateKey,
			expectedValue:         true,
		},
		"deprecated-path-keys-configured": {
			privateKeyStoragePath: caPrivateKeyStoragePathDeprecated,
			publicKeyStoragePath:  caPublicKeyStoragePathDeprecated,
			publicKey:             testCAPublicKey,
			privateKey:            testCAPrivateKey,
			expectedValue:         true,
		},
		"managed-key-configured": {
			publicKeyStoragePath: caManagedKeyStoragePath,
			publicKey:            testCAPublicKey,
			expectedValue:        true,
		},
		"stored-keys-empty": {
			privateKeyStoragePath: caPrivateKeyStoragePath,
			publicKeyStoragePath:  caPublicKeyStoragePath,
			expectedValue:         false,
		},
		"no-storage-entry": {
			expectedValue: false,
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			storage := &logical.InmemStorage{}

			if err := writeKey(ctx, storage, tt.privateKeyStoragePath, tt.privateKey); err != nil {
				t.Fatalf("error writing private key: %s", err)
			}

			if err := writeKey(ctx, storage, tt.publicKeyStoragePath, tt.publicKey); err != nil {
				t.Fatalf("error writing public key: %s", err)
			}

			keysConfigured, err := caKeysConfigured(context.Background(), storage)
			if err != nil {
				t.Fatalf("error checking for configured keys: %s", err)
			}

			if tt.expectedValue != keysConfigured {
				t.Fatalf("unexpected return value: expected %v, got %v", tt.expectedValue, keysConfigured)
			}
		})
	}
}

func writeKey(ctx context.Context, s logical.Storage, path, key string) error {
	if path == "" {
		return nil
	}

	var entry *logical.StorageEntry
	var err error
	switch path {
	case caPublicKeyStoragePath, caPrivateKeyStoragePath:
		entry, err = logical.StorageEntryJSON(path, &keyStorageEntry{Key: key})
		if err != nil {
			return err
		}
	case caPublicKeyStoragePathDeprecated, caPrivateKeyStoragePathDeprecated:
		entry, err = logical.StorageEntryJSON(path, []byte(key))
		if err != nil {
			return err
		}
	case caManagedKeyStoragePath:
		entry, err = logical.StorageEntryJSON(path, &managedKeyStorageEntry{
			KeyId:     "test-key-id",
			KeyName:   "test-key-name",
			PublicKey: key,
		})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected storage path %s", path)
	}

	return s.Put(ctx, entry)
}

func readKey(ctx context.Context, s logical.Storage, path string) error {
	switch path {
	case caPublicKeyStoragePath, caPrivateKeyStoragePath:
		var entry keyStorageEntry

		storageEntry, err := s.Get(ctx, path)
		if err != nil {
			return fmt.Errorf("error reading public key from storage: %s", err)
		}

		err = storageEntry.DecodeJSON(&entry)
		if err != nil {
			return fmt.Errorf("error decoding storage entry: %s", err)
		}

		if entry.Key == "" {
			return errors.New("stored key was empty")
		}
	case caManagedKeyStoragePath:
		var entry managedKeyStorageEntry

		storageEntry, err := s.Get(ctx, path)
		if err != nil {
			return fmt.Errorf("error reading managed key from storage: %s", err)
		}

		err = storageEntry.DecodeJSON(&entry)
		if err != nil {
			return fmt.Errorf("error decoding storage entry: %s", err)
		}

		if entry.KeyId == "" || entry.KeyName == "" || entry.PublicKey == "" {
			return errors.New("managed key storage fields were empty")
		}
	default:
		return fmt.Errorf("unexpected storage path %s", path)
	}

	return nil
}
