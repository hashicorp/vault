// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package awsauth

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
)

func TestBackend_pathConfigClient(t *testing.T) {
	config := logical.TestBackendConfig()
	storage := &logical.InmemStorage{}
	config.StorageView = storage

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// make sure we start with empty roles, which gives us confidence that the read later
	// actually is the two roles we created
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/client",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	// at this point, resp == nil is valid as no client config exists
	// if resp != nil, then resp.Data must have EndPoint and IAMServerIdHeaderValue as nil
	if resp != nil {
		if resp.IsError() {
			t.Fatalf("failed to read client config entry")
		} else if resp.Data["endpoint"] != nil || resp.Data["iam_server_id_header_value"] != nil {
			t.Fatalf("returned endpoint or iam_server_id_header_value non-nil")
		}
	}

	data := map[string]interface{}{
		"sts_endpoint":               "https://my-custom-sts-endpoint.example.com",
		"sts_region":                 "us-east-2",
		"iam_server_id_header_value": "vault_server_identification_314159",
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config/client",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal("failed to create the client config entry")
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/client",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("failed to read the client config entry")
	}
	if resp.Data["iam_server_id_header_value"] != data["iam_server_id_header_value"] {
		t.Fatalf("expected iam_server_id_header_value: '%#v'; returned iam_server_id_header_value: '%#v'",
			data["iam_server_id_header_value"], resp.Data["iam_server_id_header_value"])
	}
	if resp.Data["sts_endpoint"] != data["sts_endpoint"] {
		t.Fatalf("expected sts_endpoint: '%#v'; returned sts_endpoint: '%#v'",
			data["sts_endpoint"], resp.Data["sts_endpoint"])
	}
	if resp.Data["sts_region"] != data["sts_region"] {
		t.Fatalf("expected sts_region: '%#v'; returned sts_region: '%#v'",
			data["sts_region"], resp.Data["sts_region"])
	}

	data = map[string]interface{}{
		"sts_endpoint":               "https://my-custom-sts-endpoint2.example.com",
		"sts_region":                 "us-west-1",
		"iam_server_id_header_value": "vault_server_identification_2718281",
	}
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/client",
		Data:      data,
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.IsError() {
		t.Fatal("failed to update the client config entry")
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config/client",
		Storage:   storage,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil || resp.IsError() {
		t.Fatal("failed to read the client config entry")
	}
	if resp.Data["iam_server_id_header_value"] != data["iam_server_id_header_value"] {
		t.Fatalf("expected iam_server_id_header_value: '%#v'; returned iam_server_id_header_value: '%#v'",
			data["iam_server_id_header_value"], resp.Data["iam_server_id_header_value"])
	}
	if resp.Data["sts_endpoint"] != data["sts_endpoint"] {
		t.Fatalf("expected sts_endpoint: '%#v'; returned sts_endpoint: '%#v'",
			data["sts_endpoint"], resp.Data["sts_endpoint"])
	}
	if resp.Data["sts_region"] != data["sts_region"] {
		t.Fatalf("expected sts_region: '%#v'; returned sts_region: '%#v'",
			data["sts_region"], resp.Data["sts_region"])
	}
}

// TestBackend_PathConfigRoot_PluginIdentityToken tests parsing and validation of
// configuration used to set the secret engine up for web identity federation using
// plugin identity tokens.
func TestBackend_PathConfigRoot_PluginIdentityToken(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"identity_token_ttl":      int64(10),
		"identity_token_audience": "test-aud",
		"role_arn":                "test-role-arn",
	}

	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Path:      "config/client",
		Data:      configData,
	}

	resp, err := b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config writing failed: resp:%#v\n err: %v", resp, err)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   config.StorageView,
		Path:      "config/client",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config reading failed: resp:%#v\n err: %v", resp, err)
	}

	// Grab the subset of fields from the response we care to look at for this case
	got := map[string]interface{}{
		"identity_token_ttl":      resp.Data["identity_token_ttl"],
		"identity_token_audience": resp.Data["identity_token_audience"],
		"role_arn":                resp.Data["role_arn"],
	}

	if !reflect.DeepEqual(got, configData) {
		t.Errorf("bad: expected to read config client as %#v, got %#v instead", configData, resp.Data)
	}
}

func TestBackend_PathConfigRoot_PluginIdentityTokenWantErr(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b, err := Backend(config)
	if err != nil {
		t.Fatal(err)
	}

	err = b.Setup(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	// setting both audience and access key must result in an error due to mutual exclusivity
	configData := map[string]interface{}{
		"identity_token_audience": "test-aud",
		"access_key":              "ASIAIO10230XVB",
	}

	configReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   config.StorageView,
		Path:      "config/client",
		Data:      configData,
	}

	resp, err := b.HandleRequest(context.Background(), configReq)
	if !resp.IsError() {
		t.Fatalf("expected an error but got nil")
	}
	expectedError := "only one of 'access_key' or 'identity_token_audience' can be set"
	if !strings.Contains(resp.Error().Error(), expectedError) {
		t.Fatalf("expected err %s, got %s", expectedError, resp.Error())
	}

	// missing role arn with audience must result in an error
	configData = map[string]interface{}{
		"identity_token_audience": "test-aud",
	}

	configReq = &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Path:      "config/client",
		Data:      configData,
	}

	resp, err = b.HandleRequest(context.Background(), configReq)
	if !resp.IsError() {
		t.Fatalf("expected an error but got nil")
	}
	expectedError = "role_arn must be set when identity_token_audience is set"
	if !strings.Contains(resp.Error().Error(), expectedError) {
		t.Fatalf("expected err %s, got %s", expectedError, resp.Error())
	}
}
