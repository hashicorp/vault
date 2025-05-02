// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package connutil

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenSnowflake(t *testing.T) {
	// Generate a new RSA key for testing
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("Failed to marshal private key: %v", err)
	}

	pemBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}
	var pemKey bytes.Buffer
	pem.Encode(&pemKey, pemBlock)

	db, err := openSnowflake("account.snowflakecomputing.com/db", "user", pemKey.String())
	if err != nil {
		t.Fatalf("Failed to open Snowflake connection: %v", err)
	}

	require.NotNil(t, db.Stats())
}

func TestParseSnowflakeFieldsFromURL(t *testing.T) {
	tests := map[string]struct {
		connectionURL string
		wantAccount   string
		wantDB        string
		wantErr       error
	}{
		"valid URL": {
			connectionURL: "account.snowflakecomputing.com/db",
			wantAccount:   "account",
			wantDB:        "db",
			wantErr:       nil,
		},
		"invalid URL": {
			connectionURL: "invalid-url",
			wantAccount:   "",
			wantDB:        "",
			wantErr:       ErrInvalidSnowflakeURL,
		},
		"missing account name": {
			connectionURL: ".snowflakecomputing.com/db",
			wantAccount:   "",
			wantDB:        "",
			wantErr:       ErrInvalidSnowflakeURL,
		},
		"missing database name": {
			connectionURL: "account.snowflakecomputing.com/",
			wantAccount:   "",
			wantDB:        "",
			wantErr:       ErrInvalidSnowflakeURL,
		},
		"missing domain": {
			connectionURL: "account..com/db",
			wantAccount:   "",
			wantDB:        "",
			wantErr:       ErrInvalidSnowflakeURL,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			user, db, err := parseSnowflakeFieldsFromURL(tt.connectionURL)

			require.Equal(t, tt.wantAccount, user)
			require.Equal(t, tt.wantDB, db)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestGetPrivateKey(t *testing.T) {
	tests := map[string]struct {
		providedPrivateKey string
		wantErr            error
	}{
		"valid private key string": {
			providedPrivateKey: "-----BEGIN PRIVATE KEY-----\n",
			wantErr:            nil,
		},
		"valid private key filepath": {
			providedPrivateKey: "",
			wantErr:            nil,
		},
		"invalid private key": {
			providedPrivateKey: "-----BEGIN PRIVATE KEY-----\ninvalid\n",
			wantErr:            fmt.Errorf("failed to decode the private key value"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := getPrivateKey(tt.providedPrivateKey)

			require.Equal(t, tt.wantErr, err)
		})
	}
}
