package keysutil

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/hashicorp/vault/logical"
)

var compilerOpt []string

func TestBase58(t *testing.T) {
	tCases := []struct {
		in  string
		out string
	}{
		{
			"",
			"0",
		},
		{
			"foo",
			"sapp",
		},
		{
			"5d5746d044b9a9429249966c9e3fee178ca679b91487b11d4b73c9865202104c",
			"cozMP2pOYdDiNGeFQ2afKAOGIzO0HVpJ8OPFXuVPNbHasFyenK9CzIIPuOG7EFWOCy4YWvKGZa671N4kRSoaxZ",
		},
		{
			"5ba33e16d742f3c785f6e7e8bb6f5fe82346ffa1c47aa8e95da4ddd5a55bb334",
			"cotpEJPnhuTRofLi4lDe5iKw2fkSGc6TpUYeuWoBp8eLYJBWLRUVDZI414OjOCWXKZ0AI8gqNMoxd4eLOklwYk",
		},
		{
			" ",
			"w",
		},
		{
			"-",
			"J",
		},
		{
			"0",
			"M",
		},
		{
			"1",
			"N",
		},
		{
			"-1",
			"30B",
		},
		{
			"11",
			"3h7",
		},
		{
			"abc",
			"qMin",
		},
		{
			"1234598760",
			"1a0AFzKIPnihTq",
		},
		{
			"abcdefghijklmnopqrstuvwxyz",
			"hUBXsgd3F2swSlEgbVi2p0Ncr6kzVeJTLaW",
		},
	}

	for _, c := range tCases {
		e := Base62Encode([]byte(c.in))
		d := string(Base62Decode(e))

		if d != c.in {
			t.Fatalf("decoded value didn't match input %#v %#v", c.in, d)
		}

		if e != c.out {
			t.Fatalf("encoded value didn't match expected %#v, %#v", e, c.out)
		}
	}

	d := Base62Decode("!0000/")
	if len(d) != 0 {
		t.Fatalf("Decode of invalid string should be empty, got %#v", d)
	}
}

func TestEncrytedKeysStorage_BadPolicy(t *testing.T) {
	policy := &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              false,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
		versionPrefixCache:   &sync.Map{},
	}

	_, err := NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != ErrPolicyDerivedKeys {
		t.Fatalf("Unexpected Error: %s", err)
	}

	policy = &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: false,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
		versionPrefixCache:   &sync.Map{},
	}

	_, err = NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != ErrPolicyConvergentEncryption {
		t.Fatalf("Unexpected Error: %s", err)
	}

	policy = &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    1,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
		versionPrefixCache:   &sync.Map{},
	}

	_, err = NewEncryptedKeyStorageWrapper(EncryptedKeyStorageConfig{
		Policy: policy,
		Prefix: "prefix",
	})
	if err != ErrPolicyConvergentVersion {
		t.Fatalf("Unexpected Error: %s", err)
	}

	policy = &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
		versionPrefixCache:   &sync.Map{},
	}
}

func TestEncryptedKeysStorage_List(t *testing.T) {
	s := &logical.InmemStorage{}
	policy := &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
		versionPrefixCache:   &sync.Map{},
	}

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

	if len(keys) != 2 || keys[0] != "foo1/" || keys[1] != "foo" {
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
	policy := &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
		versionPrefixCache:   &sync.Map{},
	}

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

	if len(keys) != 2 || keys[0] != "foo1/" || keys[1] != "foo" {
		t.Fatalf("bad keys: %#v", keys)
	}

	// Test the cached value is correct
	keys, err = es.Wrap(s).List(ctx, "test/")
	if err != nil {
		t.Fatal(err)
	}

	if len(keys) != 2 || keys[0] != "foo1/" || keys[1] != "foo" {
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
	policy := &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
		versionPrefixCache:   &sync.Map{},
	}

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
	policy := &Policy{
		Name:                 "metadata",
		Type:                 KeyType_AES256_GCM96,
		Derived:              true,
		KDF:                  Kdf_hkdf_sha256,
		ConvergentEncryption: true,
		ConvergentVersion:    2,
		VersionTemplate:      EncryptedKeyPolicyVersionTpl,
		versionPrefixCache:   &sync.Map{},
	}

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
