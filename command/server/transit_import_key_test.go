// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
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
			[]string{},
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
			aes128,
			[]string{},
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
	runOpts := &base.RunOptions{
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

	priv1, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
	require.NotNil(t, priv1, "failed generating RSA 1 key")
	require.NoError(t, err, "failed generating RSA 1 key")

	rsa1, err = x509.MarshalPKCS8PrivateKey(priv1)
	require.NotNil(t, rsa1, "failed marshaling RSA 1 key")
	require.NoError(t, err, "failed marshaling RSA 1 key")

	priv2, err := cryptoutil.GenerateRSAKey(rand.Reader, 2048)
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

// TestTransitImport_Stdin verifies that transit import reads K=V data given as
// "-" from the stdin in RunOptions rather than the process's stdin: a 128-bit
// key only imports with the key type named in that input, since the default
// type is aes256-gcm96.
func TestTransitImport_Stdin(t *testing.T) {
	t.Parallel()

	client, closer := testVaultServer(t)
	defer closer()

	require.NoError(t, client.Sys().Mount("transit", &api.MountInput{Type: "transit"}))

	// Generate the 4096-bit RSA wrapping key with a longer timeout, as
	// TestTransitImport does, so it cannot time out the import's own request.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	_, err := client.Logical().ReadWithContext(ctx, "transit/wrapping_key")
	require.NoError(t, err)

	key := make([]byte, 16)
	_, err = rand.Read(key)
	require.NoError(t, err)

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	code := RunCustom([]string{
		"transit", "import", "transit/keys/stdin",
		base64.StdEncoding.EncodeToString(key),
		"-",
	}, &base.RunOptions{
		Stdin:  strings.NewReader(`{"type": "aes128-gcm96"}`),
		Stdout: stdout,
		Stderr: stderr,
		Client: client,
	})
	require.Equal(t, 0, code, "stdout: %s\nstderr: %s", stdout, stderr)

	resp, err := client.Logical().Read("transit/keys/stdin")
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "aes128-gcm96", resp.Data["type"])
}
