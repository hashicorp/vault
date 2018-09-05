package keysutil

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
)

var compilerOpt []string

func TestEncrytedKeysStorage_BadPolicy(t *testing.T) {
	policy := NewPolicy(PolicyConfig{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              false,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	})

	_, err := NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != ErrPolicyDerivedKeys {
		t.Fatalf("Unexpected Error: %s", err)
	}

	policy = NewPolicy(PolicyConfig{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: false,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	})

	_, err = NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != ErrPolicyConvergentEncryption {
		t.Fatalf("Unexpected Error: %s", err)
	}

	policy = NewPolicy(PolicyConfig{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	})
	_, err = NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != nil {
		t.Fatalf("Unexpected Error: %s", err)
	}
}

func TestEncryptedKeysStorage_List(t *testing.T) {
	s := &logical.InmemStorage{}
	policy := NewPolicy(PolicyConfig{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	})

	ctx := context.Background()

	err := policy.Rotate(ctx, s)
	if err != nil {
		t.Fatal(err)
	}

	es, err := NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = es.Wrap(s).Put(ctx, &logical.StorageEntry{
		Key:   "test",
		Value: []byte("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = es.Wrap(s).Put(ctx, &logical.StorageEntry{
		Key:   "test/foo",
		Value: []byte("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = es.Wrap(s).Put(ctx, &logical.StorageEntry{
		Key:   "test/foo1/test",
		Value: []byte("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	keys, err := es.Wrap(s).List(ctx, "test/")
	if err != nil {
		t.Fatal(err)
	}

	// Test prefixed with "/"
	keys, err = es.Wrap(s).List(ctx, "/test/")
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 2 || keys[1] != "foo1/" || keys[0] != "foo" {
		t.Fatalf("bad keys: %#v", keys)
	}

	keys, err = es.Wrap(s).List(ctx, "/")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 || keys[0] != "test" || keys[1] != "test/" {
		t.Fatalf("bad keys: %#v", keys)
	}

	keys, err = es.Wrap(s).List(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 || keys[0] != "test" || keys[1] != "test/" {
		t.Fatalf("bad keys: %#v", keys)
	}
}

func TestEncryptedKeysStorage_CRUD(t *testing.T) {
	s := &logical.InmemStorage{}
	policy := NewPolicy(PolicyConfig{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	})

	ctx := context.Background()

	err := policy.Rotate(ctx, s)
	if err != nil {
		t.Fatal(err)
	}

	es, err := NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = es.Wrap(s).Put(ctx, &logical.StorageEntry{
		Key:   "test/foo",
		Value: []byte("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = es.Wrap(s).Put(ctx, &logical.StorageEntry{
		Key:   "test/foo1/test",
		Value: []byte("test"),
	})
	if err != nil {
		t.Fatal(err)
	}

	keys, err := es.Wrap(s).List(ctx, "test/")
	if err != nil {
		t.Fatal(err)
	}

	// Test prefixed with "/"
	keys, err = es.Wrap(s).List(ctx, "/test/")
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 2 || !strutil.StrListContains(keys, "foo1/") || !strutil.StrListContains(keys, "foo") {
		t.Fatalf("bad keys: %#v", keys)
	}

	// Test the cached value is correct
	keys, err = es.Wrap(s).List(ctx, "test/")
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 2 || !strutil.StrListContains(keys, "foo1/") || !strutil.StrListContains(keys, "foo") {
		t.Fatalf("bad keys: %#v", keys)
	}

	data, err := es.Wrap(s).Get(ctx, "test/foo")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(data.Value, []byte("test")) {
		t.Fatalf("bad data: %#v", data)
	}

	err = es.Wrap(s).Delete(ctx, "test/foo")
	if err != nil {
		t.Fatal(err)
	}

	data, err = es.Wrap(s).Get(ctx, "test/foo")
	if err != nil {
		t.Fatal(err)
	}
	if data != nil {
		t.Fatal("data should be nil")
	}

}

func BenchmarkEncrytedKeyStorage_List(b *testing.B) {
	s := &logical.InmemStorage{}
	policy := NewPolicy(PolicyConfig{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	})

	ctx := context.Background()

	err := policy.Rotate(ctx, s)
	if err != nil {
		b.Fatal(err)
	}

	es, err := NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 10000; i++ {
		err = es.Wrap(s).Put(ctx, &logical.StorageEntry{
			Key:   fmt.Sprintf("test/%d", i),
			Value: []byte("test"),
		})
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		keys, err := es.Wrap(s).List(ctx, "test/")
		if err != nil {
			b.Fatal(err)
		}
		compilerOpt = keys
	}
}

func BenchmarkEncrytedKeyStorage_Put(b *testing.B) {
	s := &logical.InmemStorage{}
	policy := NewPolicy(PolicyConfig{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
	})

	ctx := context.Background()

	err := policy.Rotate(ctx, s)
	if err != nil {
		b.Fatal(err)
	}

	es, err := NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err = es.Wrap(s).Put(ctx, &logical.StorageEntry{
			Key:   fmt.Sprintf("test/%d", i),
			Value: []byte("test"),
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
