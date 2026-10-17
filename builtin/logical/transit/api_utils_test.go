// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package transit

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/keysutil"
)

// Test_parsePaddingSchemeArg validate the various use cases we have around parsing
// the various padding_scheme arg possible values.
func Test_parsePaddingSchemeArg(t *testing.T) {
	type args struct {
		keyType keysutil.KeyType
		rawPs   any
	}
	tests := []struct {
		name    string
		args    args
		want    keysutil.PaddingScheme
		wantErr bool
	}{
		// Error cases
		{name: "nil-ps", args: args{keyType: keysutil.KeyType_RSA2048, rawPs: nil}, wantErr: true},
		{name: "nonstring-ps", args: args{keyType: keysutil.KeyType_RSA2048, rawPs: 5}, wantErr: true},
		{name: "invalid-ps", args: args{keyType: keysutil.KeyType_RSA2048, rawPs: "unknown"}, wantErr: true},
		{name: "bad-keytype-oaep", args: args{keyType: keysutil.KeyType_AES128_CMAC, rawPs: "oaep"}, wantErr: true},
		{name: "bad-keytype-pkcs1", args: args{keyType: keysutil.KeyType_ECDSA_P256, rawPs: "pkcs1v15"}, wantErr: true},
		{name: "oaep-capped", args: args{keyType: keysutil.KeyType_RSA4096, rawPs: "OAEP"}, wantErr: true},
		{name: "pkcs1-whitespace", args: args{keyType: keysutil.KeyType_RSA3072, rawPs: "   pkcs1v15    "}, wantErr: true},

		// Valid cases — native key
		{name: "oaep-2048", args: args{keyType: keysutil.KeyType_RSA2048, rawPs: "oaep"}, want: keysutil.PaddingScheme_OAEP},
		{name: "oaep-3072", args: args{keyType: keysutil.KeyType_RSA3072, rawPs: "oaep"}, want: keysutil.PaddingScheme_OAEP},
		{name: "oaep-4096", args: args{keyType: keysutil.KeyType_RSA4096, rawPs: "oaep"}, want: keysutil.PaddingScheme_OAEP},
		{name: "pkcs1", args: args{keyType: keysutil.KeyType_RSA3072, rawPs: "pkcs1v15"}, want: keysutil.PaddingScheme_PKCS1v15},

		// Valid cases — managed key
		{name: "managed-key-oaep", args: args{keyType: keysutil.KeyType_MANAGED_KEY, rawPs: "oaep"}, want: keysutil.PaddingScheme_OAEP},
		{name: "managed-key-pkcs1v15", args: args{keyType: keysutil.KeyType_MANAGED_KEY, rawPs: "pkcs1v15"}, want: keysutil.PaddingScheme_PKCS1v15},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePaddingSchemeArg(tt.args.keyType, tt.args.rawPs)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePaddingSchemeArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePaddingSchemeArg() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_parseHashAlgorithmArg validates the various use cases around parsing the hash_algorithm arg.
func Test_parseHashAlgorithmArg(t *testing.T) {
	tests := []struct {
		name    string
		keyType keysutil.KeyType
		raw     string
		want    keysutil.HashType
		wantErr bool
	}{
		// Unknown hash string — error for all key types
		{name: "unknown-string", keyType: keysutil.KeyType_MANAGED_KEY, raw: "bogus", wantErr: true},
		{name: "capped-sha", keyType: keysutil.KeyType_MANAGED_KEY, raw: "SHA2-256", wantErr: true},

		// Non-RSA key type — error regardless of hash value
		{name: "aes-sha256", keyType: keysutil.KeyType_AES256_GCM96, raw: "sha2-256", wantErr: true},

		// Native RSA — "none" is a valid value that maps to HashTypeNone (caller applies sha2-256 default)
		{name: "rsa2048-none", keyType: keysutil.KeyType_RSA2048, raw: "none", want: keysutil.HashTypeNone},
		{name: "rsa2048-sha1", keyType: keysutil.KeyType_RSA2048, raw: "sha1", want: keysutil.HashTypeSHA1},
		{name: "rsa2048-sha2-256", keyType: keysutil.KeyType_RSA2048, raw: "sha2-256", want: keysutil.HashTypeSHA2256},
		{name: "rsa2048-sha2-384", keyType: keysutil.KeyType_RSA2048, raw: "sha2-384", want: keysutil.HashTypeSHA2384},
		{name: "rsa2048-sha2-512", keyType: keysutil.KeyType_RSA2048, raw: "sha2-512", want: keysutil.HashTypeSHA2512},
		{name: "rsa3072-sha2-256", keyType: keysutil.KeyType_RSA3072, raw: "sha2-256", want: keysutil.HashTypeSHA2256},
		{name: "rsa4096-sha2-256", keyType: keysutil.KeyType_RSA4096, raw: "sha2-256", want: keysutil.HashTypeSHA2256},

		// Managed key — accepts any valid hash string; per-provider restrictions are enforced downstream
		{name: "managed-key-none", keyType: keysutil.KeyType_MANAGED_KEY, raw: "none", want: keysutil.HashTypeNone},
		{name: "managed-key-sha2-224", keyType: keysutil.KeyType_MANAGED_KEY, raw: "sha2-224", want: keysutil.HashTypeSHA2224},
		{name: "managed-key-sha2-384", keyType: keysutil.KeyType_MANAGED_KEY, raw: "sha2-384", want: keysutil.HashTypeSHA2384},
		{name: "managed-key-sha2-512", keyType: keysutil.KeyType_MANAGED_KEY, raw: "sha2-512", want: keysutil.HashTypeSHA2512},
		{name: "managed-key-sha3-256", keyType: keysutil.KeyType_MANAGED_KEY, raw: "sha3-256", want: keysutil.HashTypeSHA3256},
		{name: "managed-key-sha1", keyType: keysutil.KeyType_MANAGED_KEY, raw: "sha1", want: keysutil.HashTypeSHA1},
		{name: "managed-key-sha2-256", keyType: keysutil.KeyType_MANAGED_KEY, raw: "sha2-256", want: keysutil.HashTypeSHA2256},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHashAlgorithmArg(tt.keyType, tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHashAlgorithmArg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseHashAlgorithmArg() got = %v, want %v", got, tt.want)
			}
		})
	}
}
