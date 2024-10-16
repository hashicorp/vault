// Copyright (c) HashiCorp, Inc.
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

		// Valid cases
		{name: "oaep-2048", args: args{keyType: keysutil.KeyType_RSA2048, rawPs: "oaep"}, want: keysutil.PaddingScheme_OAEP},
		{name: "oaep-3072", args: args{keyType: keysutil.KeyType_RSA3072, rawPs: "oaep"}, want: keysutil.PaddingScheme_OAEP},
		{name: "oaep-4096", args: args{keyType: keysutil.KeyType_RSA4096, rawPs: "oaep"}, want: keysutil.PaddingScheme_OAEP},
		{name: "pkcs1", args: args{keyType: keysutil.KeyType_RSA3072, rawPs: "pkcs1v15"}, want: keysutil.PaddingScheme_PKCS1v15},
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
