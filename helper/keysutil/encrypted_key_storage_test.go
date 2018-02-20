package keysutil

import (
	"context"
	"fmt"
	"reflect"
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
			"ynyL",
		},
		{
			"5d5746d044b9a9429249966c9e3fee178ca679b91487b11d4b73c9865202104c",
			"13QTmOdiTPQ6FFk4VPoLP6tnK7rrH4ofcupCjcsVetvyHwSHwFTQkn1oOnl2y0MwOm1fbRKaDSzSvvxrOSnAuQEP",
		},
		{
			"5ba33e16d742f3c785f6e7e8bb6f5fe82346ffa1c47aa8e95da4ddd5a55bb334",
			"13Qrv6yytI2utCmFN5v9jl1lV8UsmDnHjESc1idr9Tvf0GNJ26oAANUaUduRizC8Rgq7t3foeNFxAK6rNnvDhSha",
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
			"3pr",
		},
		{
			"11",
			"3H7",
		},
		{
			"abc",
			"wFbx",
		},
		{
			"1234598760",
			"2IhN69KruT1tMA",
		},
		{
			"abcdefghijklmnopqrstuvwxyz",
			"2UTr2Q0FDv7tHDPGi81CyhnbA3awFFq0R14C",
		},
	}

	for _, c := range tCases {
		e := Base58Encode([]byte(c.in))
		d := string(Base58Decode(e))

		if d != c.in {
			t.Fatalf("decoded value didn't match input %#v %#v", c.in, d)
		}

		if e != c.out {
			t.Fatalf("encoded value didn't match expected %#v, %#v", e, c.out)
		}
	}

	d := Base58Decode("!0000/")
	if len(d) != 0 {
		t.Fatalf("Decode of invalid string should be empty, got %#v", d)
	}
}

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
	}

	_, err = NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: s,
		Policy:  policy,
		Prefix:  "prefix",
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
	}

	_, err = NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: s,
		Policy:  policy,
		Prefix:  "prefix",
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
	}

	_, err = NewEncryptedKeyStorage(EncryptedKeyStorageConfig{
		Storage: nil,
		Policy:  policy,
		Prefix:  "prefix",
	})
	if err != ErrNilStorage {
		t.Fatalf("Unexpected Error: %s", err)
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
