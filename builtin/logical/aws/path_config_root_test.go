// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
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
		"sts_region":              "",
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

// TestBackend_PathConfigRoot_PluginIdentityToken tests that configuration
// of plugin WIF returns an immediate error.
func TestBackend_PathConfigRoot_PluginIdentityToken(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}

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
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.ErrorContains(t, resp.Error(), pluginidentityutil.ErrPluginWorkloadIdentityUnsupported.Error())
}

type testSystemView struct {
	logical.StaticSystemView
}

func (d testSystemView) GenerateIdentityToken(_ context.Context, _ *pluginutil.IdentityTokenRequest) (*pluginutil.IdentityTokenResponse, error) {
	return nil, pluginidentityutil.ErrPluginWorkloadIdentityUnsupported
}
