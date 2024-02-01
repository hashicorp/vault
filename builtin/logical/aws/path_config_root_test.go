// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestBackend_PathConfigRoot(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"access_key":              "AKIAEXAMPLE",
		"secret_key":              "RandomData",
		"region":                  "us-west-2",
		"iam_endpoint":            "https://iam.amazonaws.com",
		"sts_endpoint":            "https://sts.us-west-2.amazonaws.com",
		"max_retries":             10,
		"username_template":       defaultUserNameTemplate,
		"role_arn":                "",
		"identity_token_audience": "",
		"identity_token_ttl":      int64(0),
	}

	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
		Data:      configData,
	}

	resp, err := b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config writing failed: resp:%#v\n err: %v", resp, err)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config reading failed: resp:%#v\n err: %v", resp, err)
	}

	delete(configData, "secret_key")
	require.Equal(t, configData, resp.Data)
	if !reflect.DeepEqual(resp.Data, configData) {
		t.Errorf("bad: expected to read config root as %#v, got %#v instead", configData, resp.Data)
	}
}

// TestBackend_PathConfigRoot_PluginIdentityToken tests parsing and validation of
// configuration used to set the secret engine up for web identity federation using
// plugin identity tokens.
func TestBackend_PathConfigRoot_PluginIdentityToken(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
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
		Path:      "config/root",
		Data:      configData,
	}

	resp, err := b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config writing failed: resp:%#v\n err: %v", resp, err)
	}

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
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
		t.Errorf("bad: expected to read config root as %#v, got %#v instead", configData, resp.Data)
	}

	// mutually exclusive fields must result in an error
	configData = map[string]interface{}{
		"identity_token_audience": "test-aud",
		"access_key":              "ASIAIO10230XVB",
	}

	configReq = &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
		Data:      configData,
	}

	resp, err = b.HandleRequest(context.Background(), configReq)
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
		Path:      "config/root",
		Data:      configData,
	}

	resp, err = b.HandleRequest(context.Background(), configReq)
	if !resp.IsError() {
		t.Fatalf("expected an error but got nil")
	}
	expectedError = "missing required 'role_arn' when 'identity_token_audience' is set"
	if !strings.Contains(resp.Error().Error(), expectedError) {
		t.Fatalf("expected err %s, got %s", expectedError, resp.Error())
	}
}
