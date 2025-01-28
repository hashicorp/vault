// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/automatedrotationutil"
	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/rotation"
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
		"access_key":                 "AKIAEXAMPLE",
		"secret_key":                 "RandomData",
		"region":                     "us-west-2",
		"iam_endpoint":               "https://iam.amazonaws.com",
		"sts_endpoint":               "https://sts.us-west-2.amazonaws.com",
		"sts_region":                 "",
		"sts_fallback_endpoints":     []string{},
		"sts_fallback_regions":       []string{},
		"max_retries":                10,
		"username_template":          defaultUserNameTemplate,
		"role_arn":                   "",
		"identity_token_audience":    "",
		"identity_token_ttl":         int64(0),
		"rotation_schedule":          "",
		"rotation_window":            0,
		"disable_automated_rotation": false,
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
	// remove rotation_period from response for comparison with original config
	delete(resp.Data, "rotation_period")
	require.Equal(t, configData, resp.Data)
	if !reflect.DeepEqual(resp.Data, configData) {
		t.Errorf("bad: expected to read config root as %#v, got %#v instead", configData, resp.Data)
	}
}

// TestBackend_PathConfigRoot_STSFallback tests valid versions of STS fallback parameters - slice and csv
func TestBackend_PathConfigRoot_STSFallback(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"access_key":                 "AKIAEXAMPLE",
		"secret_key":                 "RandomData",
		"region":                     "us-west-2",
		"iam_endpoint":               "https://iam.amazonaws.com",
		"sts_endpoint":               "https://sts.us-west-2.amazonaws.com",
		"sts_region":                 "",
		"sts_fallback_endpoints":     []string{"192.168.1.1", "127.0.0.1"},
		"sts_fallback_regions":       []string{"my-house-1", "my-house-2"},
		"max_retries":                10,
		"username_template":          defaultUserNameTemplate,
		"role_arn":                   "",
		"identity_token_audience":    "",
		"identity_token_ttl":         int64(0),
		"rotation_schedule":          "",
		"rotation_window":            0,
		"disable_automated_rotation": false,
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
	// remove rotation_period from response for comparison with original config
	delete(resp.Data, "rotation_period")
	require.Equal(t, configData, resp.Data)
	if !reflect.DeepEqual(resp.Data, configData) {
		t.Errorf("bad: expected to read config root as %#v, got %#v instead", configData, resp.Data)
	}

	// test we can handle comma separated strings, per CommaStringSlice
	configData = map[string]interface{}{
		"access_key":                 "AKIAEXAMPLE",
		"secret_key":                 "RandomData",
		"region":                     "us-west-2",
		"iam_endpoint":               "https://iam.amazonaws.com",
		"sts_endpoint":               "https://sts.us-west-2.amazonaws.com",
		"sts_region":                 "",
		"sts_fallback_endpoints":     "1.1.1.1,8.8.8.8",
		"sts_fallback_regions":       "zone-1,zone-2",
		"max_retries":                10,
		"username_template":          defaultUserNameTemplate,
		"role_arn":                   "",
		"identity_token_audience":    "",
		"identity_token_ttl":         int64(0),
		"rotation_schedule":          "",
		"rotation_window":            0,
		"disable_automated_rotation": false,
	}

	configReq = &logical.Request{
		Operation: logical.UpdateOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
		Data:      configData,
	}

	resp, err = b.HandleRequest(context.Background(), configReq)
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
	// remove rotation_period from response for comparison with original config
	delete(resp.Data, "rotation_period")
	configData["sts_fallback_endpoints"] = []string{"1.1.1.1", "8.8.8.8"}
	configData["sts_fallback_regions"] = []string{"zone-1", "zone-2"}
	require.Equal(t, configData, resp.Data)
	if !reflect.DeepEqual(resp.Data, configData) {
		t.Errorf("bad: expected to read config root as %#v, got %#v instead", configData, resp.Data)
	}
}

// TestBackend_PathConfigRoot_STSFallback_mismatchedfallback ensures configuration writing will fail if the
// region/endpoint entries are different lengths
func TestBackend_PathConfigRoot_STSFallback_mismatchedfallback(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	// test we can handle comma separated strings, per CommaStringSlice
	configData := map[string]interface{}{
		"access_key":              "AKIAEXAMPLE",
		"secret_key":              "RandomData",
		"region":                  "us-west-2",
		"iam_endpoint":            "https://iam.amazonaws.com",
		"sts_endpoint":            "https://sts.us-west-2.amazonaws.com",
		"sts_region":              "",
		"sts_fallback_endpoints":  "1.1.1.1,8.8.8.8",
		"sts_fallback_regions":    "zone-1,zone-2",
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
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}
	if resp != nil && !resp.IsError() {
		t.Fatalf("expected an error, but it successfully wrote")
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

// TestBackend_PathConfigRoot_RegisterRootRotation tests that configuration
// and registering a root credential returns an immediate error.
func TestBackend_PathConfigRoot_RegisterRootRotation(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}

	nsCtx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	b := Backend(config)
	if err := b.Setup(nsCtx, config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"access_key":        "access-key",
		"secret_key":        "secret-key",
		"rotation_schedule": "*/30 * * * * *",
		"rotation_window":   60,
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
	assert.ErrorContains(t, resp.Error(), automatedrotationutil.ErrRotationManagerUnsupported.Error())
}

type testSystemView struct {
	logical.StaticSystemView
}

func (d testSystemView) GenerateIdentityToken(_ context.Context, _ *pluginutil.IdentityTokenRequest) (*pluginutil.IdentityTokenResponse, error) {
	return nil, pluginidentityutil.ErrPluginWorkloadIdentityUnsupported
}

func (d testSystemView) RegisterRotationJob(_ context.Context, _ *rotation.RotationJobConfigureRequest) (string, error) {
	return "", automatedrotationutil.ErrRotationManagerUnsupported
}
