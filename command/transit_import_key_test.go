// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// Validate the `vault transit import` command works.
func TestTransitImport(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	if err := client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	}); err != nil {
		t.Fatalf("transit mount error: %#v", err)
	}

	// Force the generation of the Transit wrapping key now with a longer context
	// to help the 32bit nightly tests. This creates a 4096-bit RSA key which can take
	// a while on an overloaded system
	genWrappingKeyCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if _, err := client.Logical().ReadWithContext(genWrappingKeyCtx, "transit/wrapping_key"); err != nil {
		t.Fatalf("transit failed generating wrapping key: %#v", err)
	}

	rsa1, rsa2, aes128, aes256 := generateKeys(t)

	type testCase struct {
		variant    string
		path       string
		key        []byte
		args       []string
		shouldFail bool
	}
	tests := []testCase{
		{
			"import",
			"transit/keys/rsa1",
			rsa1,
			[]string{"type=rsa-2048"},
			false, /* first import */
		},
		{
			"import",
			"transit/keys/rsa1",
			rsa2,
			[]string{"type=rsa-2048"},
			true, /* already exists */
		},
		{
			"import-version",
			"transit/keys/rsa1",
			rsa2,
			[]string{"type=rsa-2048"},
			false, /* new version */
		},
		{
			"import",
			"transit/keys/rsa2",
			rsa2,
			[]string{"type=rsa-4096"},
			true, /* wrong type */
		},
		{
			"import",
			"transit/keys/rsa2",
			rsa2,
			[]string{"type=rsa-2048"},
			false, /* new name */
		},
		{
			"import",
			"transit/keys/aes1",
			aes128,
			[]string{"type=aes128-gcm96"},
			false, /* first import */
		},
		{
			"import",
			"transit/keys/aes1",
			aes256,
			[]string{"type=aes256-gcm96"},
			true, /* already exists */
		},
		{
			"import-version",
			"transit/keys/aes1",
			aes256,
			[]string{"type=aes256-gcm96"},
			true, /* new version, different type */
		},
		{
			"import-version",
			"transit/keys/aes1",
			aes128,
			[]string{"type=aes128-gcm96"},
			false, /* new version */
		},
		{
			"import",
			"transit/keys/aes2",
			aes256,
			[]string{"type=aes128-gcm96"},
			true, /* wrong type */
		},
		{
			"import",
			"transit/keys/aes2",
			aes256,
			[]string{"type=aes256-gcm96"},
			false, /* new name */
		},
	}

	for index, tc := range tests {
		t.Logf("Running test case %d: %v", index, tc)
		execTransitImport(t, client, tc.variant, tc.path, tc.key, tc.args, tc.shouldFail)
	}
}

func execTransitImport(t *testing.T, client *api.Client, method string, path string, key []byte, data []string, expectFailure bool) {
	t.Helper()

	keyBase64 := base64.StdEncoding.EncodeToString(key)

	var args []string
	args = append(args, "transit")
	args = append(args, method)
	args = append(args, path)
	args = append(args, keyBase64)
	args = append(args, data...)

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	runOpts := &RunOptions{
		Stdout: stdout,
		Stderr: stderr,
		Client: client,
	}

	code := RunCustom(args, runOpts)
	combined := stdout.String() + stderr.String()

	if code != 0 {
		if !expectFailure {
			t.Fatalf("Got unexpected failure from test (ret %d): %v", code, combined)
		}
	} else {
		if expectFailure {
			t.Fatalf("Expected failure, got success from test (ret %d): %v", code, combined)
		}
	}
}

func generateKeys(t *testing.T) (rsa1 []byte, rsa2 []byte, aes128 []byte, aes256 []byte) {
	t.Helper()

	priv1, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NotNil(t, priv1, "failed generating RSA 1 key")
	require.NoError(t, err, "failed generating RSA 1 key")

	rsa1, err = x509.MarshalPKCS8PrivateKey(priv1)
	require.NotNil(t, rsa1, "failed marshaling RSA 1 key")
	require.NoError(t, err, "failed marshaling RSA 1 key")

	priv2, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NotNil(t, priv2, "failed generating RSA 2 key")
	require.NoError(t, err, "failed generating RSA 2 key")

	rsa2, err = x509.MarshalPKCS8PrivateKey(priv2)
	require.NotNil(t, rsa2, "failed marshaling RSA 2 key")
	require.NoError(t, err, "failed marshaling RSA 2 key")

	aes128 = make([]byte, 128/8)
	_, err = rand.Read(aes128)
	require.NoError(t, err, "failed generating AES 128 key")

	aes256 = make([]byte, 256/8)
	_, err = rand.Read(aes256)
	require.NoError(t, err, "failed generating AES 256 key")

	return
}
