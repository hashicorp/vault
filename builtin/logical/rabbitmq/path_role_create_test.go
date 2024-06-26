// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package rabbitmq

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestBackend_RoleCreate_DefaultUsernameTemplate(t *testing.T) {
	cleanup, connectionURI := prepareRabbitMQTestContainer(t)
	defer cleanup()

	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b := Backend()
	if err = b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"connection_uri":    connectionURI,
		"username":          "guest",
		"password":          "guest",
		"username_template": "",
	}
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}

	roleData := map[string]interface{}{
		"name": "foo",
		"tags": "bar",
	}
	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/foo",
		Storage:   config.StorageView,
		Data:      roleData,
	}
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}

	credsReq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "creds/foo",
		Storage:     config.StorageView,
		DisplayName: "token",
	}
	resp, err = b.HandleRequest(context.Background(), credsReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp == nil {
		t.Fatal("missing creds response")
	}
	if resp.Data == nil {
		t.Fatalf("missing creds data")
	}

	username, exists := resp.Data["username"]
	if !exists {
		t.Fatalf("missing username in response")
	}

	require.Regexp(t, `^token-[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}$`, username)
}

func TestBackend_RoleCreate_CustomUsernameTemplate(t *testing.T) {
	cleanup, connectionURI := prepareRabbitMQTestContainer(t)
	defer cleanup()

	var resp *logical.Response
	var err error
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b := Backend()
	if err = b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"connection_uri":    connectionURI,
		"username":          "guest",
		"password":          "guest",
		"username_template": "foo-{{ .DisplayName }}",
	}
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/connection",
		Storage:   config.StorageView,
		Data:      configData,
	}
	resp, err = b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}

	roleData := map[string]interface{}{
		"name": "foo",
		"tags": "bar",
	}
	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/foo",
		Storage:   config.StorageView,
		Data:      roleData,
	}
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp != nil {
		t.Fatal("expected a nil response")
	}

	credsReq := &logical.Request{
		Operation:   logical.ReadOperation,
		Path:        "creds/foo",
		Storage:     config.StorageView,
		DisplayName: "token",
	}
	resp, err = b.HandleRequest(context.Background(), credsReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: resp: %#v\nerr:%s", resp, err)
	}
	if resp == nil {
		t.Fatal("missing creds response")
	}
	if resp.Data == nil {
		t.Fatalf("missing creds data")
	}

	username, exists := resp.Data["username"]
	if !exists {
		t.Fatalf("missing username in response")
	}

	require.Regexp(t, `^foo-token$`, username)
}
