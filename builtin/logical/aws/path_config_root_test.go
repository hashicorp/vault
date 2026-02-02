// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/automatedrotationutil"
	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/observations"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/rotation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackend_PathConfigRoot(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or
	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	// Create operation
	configData := map[string]interface{}{
		"access_key":                 "AKIAEXAMPLE",
		"secret_key":                 "RandomData",
		"region":                     "us-west-2",
		"iam_endpoint":               "https://iam.amazonaws.com",
		"sts_endpoint":               "https://sts.us-west-2.amazonaws.com",
		"sts_region":                 "",
		"sts_fallback_endpoints":     []string{},
		"sts_fallback_regions":       []string{},
		"role_arn":                   "",
		"identity_token_audience":    "",
		"identity_token_ttl":         int64(0),
		"rotation_schedule":          "",
		"rotation_policy":            "",
		"rotation_period":            time.Duration(0).Seconds(),
		"rotation_window":            time.Duration(0).Seconds(),
		"disable_automated_rotation": false,
	}

	configReq := &logical.Request{
		Operation: logical.CreateOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
		Data:      configData,
	}

	resp, err := b.HandleRequest(context.Background(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config writing failed: resp:%#v\n err: %v", resp, err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config reading failed: resp:%#v\n err: %v", resp, err)
	}
	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigRead))

	// Ensure default values are enforced
	configData["max_retries"] = -1
	configData["username_template"] = defaultUserNameTemplate

	delete(configData, "secret_key")
	require.Equal(t, configData, resp.Data)
	if !reflect.DeepEqual(resp.Data, configData) {
		t.Errorf("bad: expected to read config root as %#v, got %#v instead", configData, resp.Data)
	}

	// Update operation
	configData = map[string]interface{}{
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
		"rotation_policy":            "",
		"rotation_period":            time.Duration(0).Seconds(),
		"rotation_window":            time.Duration(0).Seconds(),
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
	require.Equal(t, 2, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))

	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config reading failed: resp:%#v\n err: %v", resp, err)
	}
	require.Equal(t, 2, or.NumObservationsByType(ObservationTypeAWSRootConfigRead))

	delete(configData, "secret_key")
	require.Equal(t, configData, resp.Data)
	if !reflect.DeepEqual(resp.Data, configData) {
		t.Errorf("bad: expected to read config root as %#v, got %#v instead", configData, resp.Data)
	}
}

// TestBackend_PathConfigRoot_STSFallback tests valid versions of STS fallback parameters - slice and csv
func TestBackend_PathConfigRoot_STSFallback(t *testing.T) {
	config := logical.TestBackendConfig()
	or := observations.NewTestObservationRecorder()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	config.ObservationRecorder = or

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
		"rotation_policy":            "",
		"rotation_window":            time.Duration(0).Seconds(),
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

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config reading failed: resp:%#v\n err: %v", resp, err)
	}
	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigRead))

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
		"rotation_policy":            "",
		"rotation_window":            time.Duration(0).Seconds(),
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

	require.Equal(t, 2, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Storage:   config.StorageView,
		Path:      "config/root",
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config reading failed: resp:%#v\n err: %v", resp, err)
	}

	require.Equal(t, 2, or.NumObservationsByType(ObservationTypeAWSRootConfigRead))

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
// sts fallback regions and sts fallback endpoints entries are different lengths
func TestBackend_PathConfigRoot_STSFallback_mismatchedfallback(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	// sts fallback endpoints has 2 entries, regions has 1
	configData := map[string]interface{}{
		"access_key":              "AKIAEXAMPLE",
		"secret_key":              "RandomData",
		"region":                  "us-west-2",
		"iam_endpoint":            "https://iam.amazonaws.com",
		"sts_endpoint":            "https://sts.us-west-2.amazonaws.com",
		"sts_region":              "",
		"sts_fallback_endpoints":  "1.1.1.1,8.8.8.8",
		"sts_fallback_regions":    "zone-1",
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
	require.NotNil(t, resp)
	require.True(t, resp.IsError())
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
}

// TestBackend_PathConfigRoot_STSFallback_defaultEndpointRegion ensures that if no endpoints are specified, we can
// still make a config with the appropriate values.
func TestBackend_PathConfigRoot_STSFallback_defaultEndpointRegion(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"access_key":              "AKIAEXAMPLE",
		"secret_key":              "RandomData",
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

	_, err := b.HandleRequest(context.Background(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	cfgs, err := b.getRootSTSConfigs(context.Background(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get STS configs with default region/endpoints: %v", err)
	}
	if len(cfgs) != 1 {
		t.Fatalf("got %d configs, but expected 1", len(cfgs))
	} else {
		cfg := cfgs[0]
		if *(cfg.Endpoint) != matchingSTSEndpoint(*(cfg.Region)) {
			t.Fatalf("region and endpoint didn't match: %s vs. %s", *(cfg.Region), *(cfg.Endpoint))
		}
	}
	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
}

// TestBackend_PathConfigRoot_IAM_specifiedRegion ensures that if a region is set, we get a good config (with a blank
// endpoint)
func TestBackend_PathConfigRoot_IAM_specifiedRegion(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	desiredRegion := "us-west-2"

	configData := map[string]interface{}{
		"access_key":              "AKIAEXAMPLE",
		"secret_key":              "RandomData",
		"max_retries":             10,
		"username_template":       defaultUserNameTemplate,
		"region":                  desiredRegion,
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

	_, err := b.HandleRequest(context.Background(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	cfg, err := b.getRootIAMConfig(context.Background(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get IAM configs with default region/endpoints: %v", err)
	}
	if *(cfg.Endpoint) != "" {
		t.Fatalf("endpoint should have remained blank but it became %s", *(cfg.Endpoint))
	}
	if *(cfg.Region) != desiredRegion {
		t.Fatalf("region changed from config: %s became %s", desiredRegion, *(cfg.Region))
	}
}

// TestBackend_PathConfigRoot_IAM_specifiedRegionAndEndpoint ensures that if a region and endpoint are set, we get a
// good config
func TestBackend_PathConfigRoot_IAM_specifiedRegionAndEndpoint(t *testing.T) {
	config := logical.TestBackendConfig()
	or := observations.NewTestObservationRecorder()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	config.ObservationRecorder = or

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	desiredRegion := "custom-region"
	desiredEndpoint := "https://custom-endpoint.local"

	configData := map[string]interface{}{
		"access_key":              "AKIAEXAMPLE",
		"secret_key":              "RandomData",
		"max_retries":             10,
		"username_template":       defaultUserNameTemplate,
		"region":                  desiredRegion,
		"iam_endpoint":            desiredEndpoint,
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

	_, err := b.HandleRequest(context.Background(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	cfg, err := b.getRootIAMConfig(context.Background(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get IAM configs with default region/endpoints: %v", err)
	}

	if *(cfg.Endpoint) != desiredEndpoint {
		t.Fatalf("endpoint should have been %s but it became %s", desiredEndpoint, *(cfg.Endpoint))
	}
	if *(cfg.Region) != desiredRegion {
		t.Fatalf("region changed from config: %s became %s", desiredRegion, *(cfg.Region))
	}
}

// TestBackend_PathConfigRoot_IAM_defaultEndpointRegion ensures that if no endpoints are specified, we can still
// make a config with the appropriate values.
func TestBackend_PathConfigRoot_IAM_defaultEndpointRegion(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"access_key":              "AKIAEXAMPLE",
		"secret_key":              "RandomData",
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

	_, err := b.HandleRequest(context.Background(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	cfg, err := b.getRootIAMConfig(context.Background(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get IAM configs with default region/endpoints: %v", err)
	}
	// ensure endpoint is blank, because AWS wants that
	if *(cfg.Endpoint) != "" {
		t.Fatalf("expected endpoint to be blank but it was %s", *(cfg.Endpoint))
	}
}

// TestBackend_PathConfigRoot_STSIAM_SetEverything ensures that if both IAM and STS are configured, they interact
// correctly.
func TestBackend_PathConfigRoot_STSIAM_SetEverything(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or

	b := Backend(config)
	if err := b.Setup(context.Background(), config); err != nil {
		t.Fatal(err)
	}

	desiredRegion := "us-west-2"
	stsRegion := "us-east-1"
	stsEndpoint := "https://sts.us-east-1.amazonaws.com"

	configData := map[string]interface{}{
		"access_key":              "AKIAEXAMPLE",
		"secret_key":              "RandomData",
		"max_retries":             10,
		"username_template":       defaultUserNameTemplate,
		"region":                  desiredRegion,
		"sts_region":              stsRegion,
		"sts_endpoint":            stsEndpoint,
		"sts_fallback_regions":    "ap-west-1,fake-region-2",
		"sts_fallback_endpoints":  "1.1.1.1,192.168.2.3",
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

	_, err := b.HandleRequest(context.Background(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	// get IAM
	cfg, err := b.getRootIAMConfig(context.Background(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get IAM configs with default region/endpoints: %v", err)
	}

	if *(cfg.Endpoint) != "" {
		t.Fatalf("endpoint should have remained blank but it became %s", *(cfg.Endpoint))
	}
	if *(cfg.Region) != desiredRegion {
		t.Fatalf("region changed from config: %s became %s", desiredRegion, *(cfg.Region))
	}

	// get STS
	cfgs, err := b.getRootSTSConfigs(context.Background(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get IAM configs with default region/endpoints: %v", err)
	}
	if len(cfgs) != 3 {
		t.Fatalf("got %d configs, but expected 3", len(cfgs))
	}
}

// TestBackend_PathConfigRoot_PluginIdentityToken tests that configuration
// of plugin WIF returns an immediate error.
func TestBackend_PathConfigRoot_PluginIdentityToken(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or

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
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
}

// TestBackend_PathConfigRoot_RegisterRootRotation tests that configuration
// and registering a root credential returns an immediate error.
func TestBackend_PathConfigRoot_RegisterRootRotation(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or

	nsCtx := namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace)

	b := Backend(config)
	if err := b.Setup(nsCtx, config); err != nil {
		t.Fatal(err)
	}

	configData := map[string]interface{}{
		"access_key":        "access-key",
		"secret_key":        "secret-key",
		"rotation_schedule": "*/1 * * * *",
		"rotation_window":   120,
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
	require.Equal(t, 0, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
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

func (d testSystemView) DeregisterRotationJob(_ context.Context, _ *rotation.RotationJobDeregisterRequest) error {
	return nil
}
