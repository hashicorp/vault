// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package cert

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agentproxyshared/auth"
)

func TestCertAuthMethod_Authenticate(t *testing.T) {
	config := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "cert-test",
		Config: map[string]interface{}{
			"name": "foo",
		},
	}

	method, err := NewCertAuthMethod(config)
	if err != nil {
		t.Fatal(err)
	}
	defer method.Shutdown()

	client, err := api.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	loginPath, _, authMap, err := method.Authenticate(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}

	expectedLoginPath := path.Join(config.MountPath, "/login")
	if loginPath != expectedLoginPath {
		t.Fatalf("mismatch on login path: got: %s, expected: %s", loginPath, expectedLoginPath)
	}

	expectedAuthMap := map[string]interface{}{
		"name": config.Config["name"],
	}
	if !reflect.DeepEqual(authMap, expectedAuthMap) {
		t.Fatalf("mismatch on login path:\ngot:\n\t%v\nexpected:\n\t%v", authMap, expectedAuthMap)
	}
}

func TestCertAuthMethod_AuthClient_withoutCerts(t *testing.T) {
	config := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "cert-test",
		Config: map[string]interface{}{
			"name": "without-certs",
		},
	}

	method, err := NewCertAuthMethod(config)
	if err != nil {
		t.Fatal(err)
	}
	defer method.Shutdown()

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	clientToUse, err := method.(auth.AuthMethodWithClient).AuthClient(client)
	if err != nil {
		t.Fatal(err)
	}

	if client != clientToUse {
		t.Fatal("error: expected AuthClient to return back original client")
	}
}

func TestCertAuthMethod_AuthClient_withCerts(t *testing.T) {
	clientCert, err := os.Open("./test-fixtures/keys/cert.pem")
	if err != nil {
		t.Fatal(err)
	}
	defer clientCert.Close()

	clientKey, err := os.Open("./test-fixtures/keys/key.pem")
	if err != nil {
		t.Fatal(err)
	}
	defer clientKey.Close()

	config := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "cert-test",
		Config: map[string]interface{}{
			"name":        "with-certs",
			"client_cert": clientCert.Name(),
			"client_key":  clientKey.Name(),
		},
	}

	method, err := NewCertAuthMethod(config)
	if err != nil {
		t.Fatal(err)
	}
	defer method.Shutdown()

	client, err := api.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	clientToUse, err := method.(auth.AuthMethodWithClient).AuthClient(client)
	if err != nil {
		t.Fatal(err)
	}

	if client == clientToUse {
		t.Fatal("expected client from AuthClient to be different from original client")
	}

	// Call AuthClient again to get back the cached client
	cachedClient, err := method.(auth.AuthMethodWithClient).AuthClient(client)
	if err != nil {
		t.Fatal(err)
	}

	if cachedClient != clientToUse {
		t.Fatal("expected client from AuthClient to return back a cached client")
	}
}

func copyFile(from, to string) error {
	data, err := os.ReadFile(from)
	if err != nil {
		return err
	}

	return os.WriteFile(to, data, 0o600)
}

// TestCertAuthMethod_AuthClient_withCertsReload makes the file change and ensures the cert auth method deliver the event.
func TestCertAuthMethod_AuthClient_withCertsReload(t *testing.T) {
	// Initial the cert/key pair to temp path
	certPath := filepath.Join(os.TempDir(), "app.crt")
	keyPath := filepath.Join(os.TempDir(), "app.key")
	if err := copyFile("./test-fixtures/keys/cert.pem", certPath); err != nil {
		t.Fatal("copy cert file failed", err)
	}
	defer os.Remove(certPath)
	if err := copyFile("./test-fixtures/keys/key.pem", keyPath); err != nil {
		t.Fatal("copy key file failed", err)
	}
	defer os.Remove(keyPath)

	config := &auth.AuthConfig{
		Logger:    hclog.NewNullLogger(),
		MountPath: "cert-test",
		Config: map[string]interface{}{
			"name":          "with-certs-reloaded",
			"client_cert":   certPath,
			"client_key":    keyPath,
			"reload":        true,
			"reload_period": 1,
		},
	}

	method, err := NewCertAuthMethod(config)
	if err != nil {
		t.Fatal(err)
	}
	defer method.Shutdown()

	client, err := api.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}

	clientToUse, err := method.(auth.AuthMethodWithClient).AuthClient(client)
	if err != nil {
		t.Fatal(err)
	}

	if client == clientToUse {
		t.Fatal("expected client from AuthClient to be different from original client")
	}

	// Call AuthClient again to get back a new client with reloaded certificates
	reloadedClient, err := method.(auth.AuthMethodWithClient).AuthClient(client)
	if err != nil {
		t.Fatal(err)
	}

	if reloadedClient == clientToUse {
		t.Fatal("expected client from AuthClient to return back a new client")
	}

	method.CredSuccess()
	// Only make a change to the cert file, it doesn't match the key file so the client won't pick and load them.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err = copyFile("./test-fixtures/keys/cert1.pem", certPath); err != nil {
		t.Fatal("update cert file failed", err)
	}

	select {
	case <-ctx.Done():
	case <-method.NewCreds():
		cancel()
		t.Fatal("malformed cert should not be observed as a change")
	}

	// Make a change to the key file and now they are good to be picked.
	if err = copyFile("./test-fixtures/keys/key1.pem", keyPath); err != nil {
		t.Fatal("update key file failed", err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	select {
	case <-ctx.Done():
		t.Fatal("failed to watch the cert change: timeout")
	case <-method.NewCreds():
		cancel()
	}
}

// TestCertAuthMethod_hashCert_withEmptyPaths tests hashCert() if it works well with optional options.
func TestCertAuthMethod_hashCert_withEmptyPaths(t *testing.T) {
	c := &certMethod{
		logger: hclog.NewNullLogger(),
	}

	// It skips empty file paths
	sum, err := c.hashCert("", "", "")
	if sum == "" || err != nil {
		t.Fatal("hashCert() should skip empty file paths and succeed.")
	}
	emptySum := sum

	// Only present ca cert
	sum, err = c.hashCert("", "", "./test-fixtures/root/rootcacert.pem")
	if sum == "" || err != nil {
		t.Fatal("hashCert() should succeed when only present ca cert.")
	}

	// Only present client cert/key
	sum, err = c.hashCert("./test-fixtures/keys/cert.pem", "./test-fixtures/keys/key.pem", "")
	if sum == "" || err != nil {
		fmt.Println(sum, err)
		t.Fatal("hashCert() should succeed when only present client cert/key.")
	}

	// The client cert/key should be presented together or will be skipped
	sum, err = c.hashCert("./test-fixtures/keys/cert.pem", "", "")
	if sum == "" || err != nil {
		t.Fatal("hashCert() should succeed when only present client cert.")
	} else if sum != emptySum {
		t.Fatal("hashCert() should skip the client cert/key when only present client cert.")
	}
}

// TestCertAuthMethod_hashCert_withInvalidClientCert adds test cases for invalid input for hashCert().
func TestCertAuthMethod_hashCert_withInvalidClientCert(t *testing.T) {
	c := &certMethod{
		logger: hclog.NewNullLogger(),
	}

	// With mismatched cert/key pair
	sum, err := c.hashCert("./test-fixtures/keys/cert1.pem", "./test-fixtures/keys/key.pem", "")
	if sum != "" || err == nil {
		t.Fatal("hashCert() should fail with invalid client cert.")
	}

	// With non-existed paths
	sum, err = c.hashCert("./test-fixtures/keys/cert2.pem", "./test-fixtures/keys/key.pem", "")
	if sum != "" || err == nil {
		t.Fatal("hashCert() should fail with non-existed client cert path.")
	}
}

// TestCertAuthMethod_hashCert_withChange tests hashCert() if it detects changes from both client cert/key and ca cert.
func TestCertAuthMethod_hashCert_withChange(t *testing.T) {
	c := &certMethod{
		logger: hclog.NewNullLogger(),
	}

	// A good first case.
	sum, err := c.hashCert("./test-fixtures/keys/cert.pem", "./test-fixtures/keys/key.pem", "./test-fixtures/root/rootcacert.pem")
	if sum == "" || err != nil {
		t.Fatal("hashCert() shouldn't fail with a valid pair of cert/key.")
	}

	// Only change the ca cert from the first case.
	sum1, err := c.hashCert("./test-fixtures/keys/cert.pem", "./test-fixtures/keys/key.pem", "./test-fixtures/keys/cert.pem")
	if sum1 == "" || err != nil {
		t.Fatal("hashCert() shouldn't fail with valid pair of cert/key.")
	} else if sum == sum1 {
		t.Fatal("The hash should be different with a different ca cert.")
	}

	// Only change the cert/key pair from the first case.
	sum2, err := c.hashCert("./test-fixtures/keys/cert1.pem", "./test-fixtures/keys/key1.pem", "./test-fixtures/root/rootcacert.pem")
	if sum2 == "" || err != nil {
		t.Fatal("hashCert() shouldn't fail with a valid cert/key pair")
	} else if sum == sum2 || sum1 == sum2 {
		t.Fatal("The hash should be different with a different pair of cert/key.")
	}
}
