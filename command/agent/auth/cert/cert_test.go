// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cert

import (
	"context"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/agent/auth"
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

func TestCertAuthMethod_AuthClient_withCertsReload(t *testing.T) {
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
			"name":        "with-certs-reloaded",
			"client_cert": clientCert.Name(),
			"client_key":  clientKey.Name(),
			"reload":      true,
		},
	}

	method, err := NewCertAuthMethod(config)
	if err != nil {
		t.Fatal(err)
	}

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
}
