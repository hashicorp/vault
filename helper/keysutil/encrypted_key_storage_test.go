package keysutil

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/logical"
)

var compilerOpt []string

func TestEncrytedKeysStorage_BadPolicy(t *testing.T) {
	s := &logical.InmemStorage{}
	policy := &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              false,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	}

	_, err := NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: s,
		Policy:  policy,
		Prefix:  "prefix",
	})
	if err != ErrPolicyDerivedKeys {
		t.Fatal("Unexpected Error: %s", err)
	}

	policy = &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: false,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	}

	_, err = NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: s,
		Policy:  policy,
		Prefix:  "prefix",
	})
	if err != ErrPolicyConvergentEncryption {
		t.Fatal("Unexpected Error: %s", err)
	}

	policy = &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    1,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	}

	_, err = NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: s,
		Policy:  policy,
		Prefix:  "prefix",
	})
	if err != ErrPolicyConvergentVersion {
		t.Fatal("Unexpected Error: %s", err)
	}

	policy = &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	}

	_, err = NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: nil,
		Policy:  policy,
		Prefix:  "prefix",
	})
	if err != ErrNilStorage {
		t.Fatal("Unexpected Error: %s", err)
	}
}

func TestEncrytedKeysStorage_CRUD(t *testing.T) {
	s := &logical.InmemStorage{}
	policy := &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	}

	ctx := context.Background()

	err := policy.Rotate(ctx, s)
	if err != nil {
		t.Fatal(err)
	}

	es, err := NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: s,
		Policy:  policy,
		Prefix:  "prefix",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = es.Put(ctx, &logical.StorageEntry{
		Key:   "test/foo",
		Value: []byte("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = es.Put(ctx, &logical.StorageEntry{
		Key:   "test/foo1/test",
		Value: []byte("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	keys, err := es.List(ctx, "test/")
	if err != nil {
		t.Fatal(err)
	}

	// Test prefixed with "/"
	keys, err = es.List(ctx, "/test/")
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 2 || keys[0] != "foo1/" || keys[1] != "foo" {
		t.Fatalf("bad keys: %#v", keys)
	}

	// Test the cached value is correct
	keys, err = es.List(ctx, "test/")
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 2 || keys[0] != "foo1/" || keys[1] != "foo" {
		t.Fatalf("bad keys: %#v", keys)
	}

	data, err := es.Get(ctx, "test/foo")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(data.Value, []byte("test")) {
		t.Fatalf("bad data: %#v", data)
	}

	err = es.Delete(ctx, "test/foo")
	if err != nil {
		t.Fatal(err)
	}

	data, err = es.Get(ctx, "test/foo")
	if err != nil {
		t.Fatal(err)
	}
	if data != nil {
		t.Fatal("data should be nil")
	}

}

func BenchmarkEncrytedKeyStorage_List(b *testing.B) {
	s := &logical.InmemStorage{}
	policy := &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	}

	ctx := context.Background()

	err := policy.Rotate(ctx, s)
	if err != nil {
		b.Fatal(err)
	}

	es, err := NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: s,
		Policy:  policy,
		Prefix:  "prefix",
	})
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 10000; i++ {
		err = es.Put(ctx, &logical.StorageEntry{
			Key:   fmt.Sprintf("test/%d", i),
			Value: []byte("test"),
		})
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		keys, err := es.List(ctx, "test/")
		if err != nil {
			b.Fatal(err)
		}
		compilerOpt = keys
	}
}

func BenchmarkEncrytedKeyStorage_Put(b *testing.B) {
	s := &logical.InmemStorage{}
	policy := &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	}

	ctx := context.Background()

	err := policy.Rotate(ctx, s)
	if err != nil {
		b.Fatal(err)
	}

	es, err := NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: s,
		Policy:  policy,
		Prefix:  "prefix",
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = es.Put(ctx, &logical.StorageEntry{
			Key:   fmt.Sprintf("test/%d", i),
			Value: []byte("test"),
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
