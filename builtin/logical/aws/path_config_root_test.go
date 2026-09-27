// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package aws

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/vault/internalshared/namespace"
	"github.com/hashicorp/vault/sdk/helper/automatedrotationutil"
	"github.com/hashicorp/vault/sdk/helper/pluginidentityutil"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/observations"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/rotation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBackend_NoRootIMDSSTSConnectionReuse verifies that, without root configuration,
// IMDS credentials are cached and separate AssumeRole calls reuse the STS connection.
func TestBackend_NoRootIMDSSTSConnectionReuse(t *testing.T) {
	// SDK configuration reads process-wide environment variables, so this test cannot run in parallel.
	t.Setenv("AWS_CONFIG_FILE", t.TempDir()+"/config")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", t.TempDir()+"/credentials")
	t.Setenv("AWS_REGION", "us-east-1")
	for _, key := range []string{
		"AWS_PROFILE", "AWS_DEFAULT_PROFILE",
		"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_SESSION_TOKEN",
		"AWS_ROLE_ARN", "AWS_WEB_IDENTITY_TOKEN_FILE",
		"AWS_CONTAINER_CREDENTIALS_RELATIVE_URI", "AWS_CONTAINER_CREDENTIALS_FULL_URI",
	} {
		t.Setenv(key, "")
	}
	t.Setenv("AWS_EC2_METADATA_DISABLED", "false")
	t.Setenv("AWS_EC2_METADATA_V1_DISABLED", "true")

	var credentialFetches atomic.Int32
	imds := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/latest/api/token" {
			assert.Equal(t, http.MethodPut, r.Method)
			w.Header().Set("X-aws-ec2-metadata-token-ttl-seconds", "21600")
			fmt.Fprint(w, "test-imds-token")
			return
		}
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "test-imds-token", r.Header.Get("X-aws-ec2-metadata-token"))
		switch r.URL.Path {
		case "/latest/meta-data/iam/security-credentials/":
			fmt.Fprint(w, "test-instance-role")
		case "/latest/meta-data/iam/security-credentials/test-instance-role":
			credentialFetches.Add(1)
			fmt.Fprintf(w, `{"Code":"Success","AccessKeyId":"imds-access","SecretAccessKey":"imds-secret","Token":"imds-session","Expiration":%q}`,
				time.Now().Add(time.Hour).UTC().Format(time.RFC3339))
		default:
			http.Error(w, "unexpected metadata path", http.StatusNotFound)
		}
	}))
	t.Cleanup(imds.Close)
	t.Setenv("AWS_EC2_METADATA_SERVICE_ENDPOINT", imds.URL)

	var connections, calls atomic.Int32
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.Header.Get("Authorization"), "Credential=imds-access/")
		assert.Equal(t, "imds-session", r.Header.Get("X-Amz-Security-Token"))
		if assert.NoError(t, r.ParseForm()) {
			assert.Equal(t, "AssumeRole", r.Form.Get("Action"))
		}
		calls.Add(1)
		w.Header().Set("Content-Type", "text/xml")
		fmt.Fprint(w, `<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/"><AssumeRoleResult><Credentials><AccessKeyId>test-access</AccessKeyId><SecretAccessKey>test-secret</SecretAccessKey><SessionToken>test-token</SessionToken><Expiration>2099-01-01T00:00:00Z</Expiration></Credentials></AssumeRoleResult></AssumeRoleResponse>`)
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			connections.Add(1)
		}
	}
	server.Start()
	t.Cleanup(server.Close)

	conf := logical.TestBackendConfig()
	b := Backend(conf)
	b.clientMutex.RLock()
	configs, err := b.getRootSTSConfigs(t.Context(), &logical.InmemStorage{}, conf.Logger)
	b.clientMutex.RUnlock()
	require.NoError(t, err, "no-root STS configuration must load")
	require.Len(t, configs, 1, "no-root path must return one configuration")
	client, ok := configs[0].HTTPClient.(*http.Client)
	require.True(t, ok, "SDK configuration must retain the HTTP client")
	t.Cleanup(client.CloseIdleConnections)

	// Redirect only STS; leave the IMDS credential chain and pooled transport intact.
	configs[0].BaseEndpoint = aws.String(server.URL)
	stsClient := sts.NewFromConfig(*configs[0])
	for range 2 {
		response, err := stsClient.AssumeRole(t.Context(), &sts.AssumeRoleInput{
			RoleArn:         aws.String("arn:aws:iam::123456789012:role/test"),
			RoleSessionName: aws.String("connection-reuse-test"),
		})
		require.NoError(t, err, "each STS call must succeed")
		require.NotNil(t, response.Credentials, "each STS call must return credentials")
	}
	require.EqualValues(t, 2, calls.Load(), "issued credentials must not be cached")
	require.EqualValues(t, 1, connections.Load(), "both STS calls must reuse one TCP connection")
	require.EqualValues(t, 1, credentialFetches.Load(), "source IMDS credentials must be cached across STS calls")
}

func TestBackend_PathConfigRoot(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	config.System = &testSystemView{}
	or := observations.NewTestObservationRecorder()
	config.ObservationRecorder = or
	b := Backend(config)
	if err := b.Setup(t.Context(), config); err != nil {
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

	resp, err := b.HandleRequest(t.Context(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config writing failed: resp:%#v\n err: %v", resp, err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	resp, err = b.HandleRequest(t.Context(), &logical.Request{
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

	resp, err = b.HandleRequest(t.Context(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config writing failed: resp:%#v\n err: %v", resp, err)
	}
	require.Equal(t, 2, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))

	resp, err = b.HandleRequest(t.Context(), &logical.Request{
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
	if err := b.Setup(t.Context(), config); err != nil {
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

	resp, err := b.HandleRequest(t.Context(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config writing failed: resp:%#v\n err: %v", resp, err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	resp, err = b.HandleRequest(t.Context(), &logical.Request{
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

	resp, err = b.HandleRequest(t.Context(), configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: config writing failed: resp:%#v\n err: %v", resp, err)
	}

	require.Equal(t, 2, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	resp, err = b.HandleRequest(t.Context(), &logical.Request{
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
	if err := b.Setup(t.Context(), config); err != nil {
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

	resp, err := b.HandleRequest(t.Context(), configReq)
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
	if err := b.Setup(t.Context(), config); err != nil {
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

	_, err := b.HandleRequest(t.Context(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	cfgs, err := b.getRootSTSConfigs(t.Context(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get STS configs with default region/endpoints: %v", err)
	}
	if len(cfgs) != 1 {
		t.Fatalf("got %d configs, but expected 1", len(cfgs))
	} else {
		cfg := cfgs[0]
		if aws.ToString(cfg.BaseEndpoint) != matchingSTSEndpoint(cfg.Region) {
			t.Fatalf("region and endpoint didn't match: %s vs. %s", cfg.Region, aws.ToString(cfg.BaseEndpoint))
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
	if err := b.Setup(t.Context(), config); err != nil {
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

	_, err := b.HandleRequest(t.Context(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	cfg, err := b.getRootIAMConfig(t.Context(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get IAM configs with default region/endpoints: %v", err)
	}
	if aws.ToString(cfg.BaseEndpoint) != "" {
		t.Fatalf("endpoint should have remained blank but it became %s", aws.ToString(cfg.BaseEndpoint))
	}
	if cfg.Region != desiredRegion {
		t.Fatalf("region changed from config: %s became %s", desiredRegion, cfg.Region)
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
	if err := b.Setup(t.Context(), config); err != nil {
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

	_, err := b.HandleRequest(t.Context(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	cfg, err := b.getRootIAMConfig(t.Context(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get IAM configs with default region/endpoints: %v", err)
	}

	if aws.ToString(cfg.BaseEndpoint) != desiredEndpoint {
		t.Fatalf("endpoint should have been %s but it became %s", desiredEndpoint, aws.ToString(cfg.BaseEndpoint))
	}
	if cfg.Region != desiredRegion {
		t.Fatalf("region changed from config: %s became %s", desiredRegion, cfg.Region)
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
	if err := b.Setup(t.Context(), config); err != nil {
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

	_, err := b.HandleRequest(t.Context(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	cfg, err := b.getRootIAMConfig(t.Context(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get IAM configs with default region/endpoints: %v", err)
	}
	// ensure endpoint is blank, because AWS wants that
	if aws.ToString(cfg.BaseEndpoint) != "" {
		t.Fatalf("expected endpoint to be blank but it was %s", aws.ToString(cfg.BaseEndpoint))
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
	if err := b.Setup(t.Context(), config); err != nil {
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

	_, err := b.HandleRequest(t.Context(), configReq)
	if err != nil {
		t.Fatalf("bad: config writing failed: err: %v", err)
	}

	require.Equal(t, 1, or.NumObservationsByType(ObservationTypeAWSRootConfigWrite))
	// get IAM
	cfg, err := b.getRootIAMConfig(t.Context(), config.StorageView, b.Logger())
	if err != nil {
		t.Fatalf("couldn't get IAM configs with default region/endpoints: %v", err)
	}

	if aws.ToString(cfg.BaseEndpoint) != "" {
		t.Fatalf("endpoint should have remained blank but it became %s", aws.ToString(cfg.BaseEndpoint))
	}
	if cfg.Region != desiredRegion {
		t.Fatalf("region changed from config: %s became %s", desiredRegion, cfg.Region)
	}

	// get STS
	cfgs, err := b.getRootSTSConfigs(t.Context(), config.StorageView, b.Logger())
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
	if err := b.Setup(t.Context(), config); err != nil {
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

	resp, err := b.HandleRequest(t.Context(), configReq)
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

	nsCtx := namespace.ContextWithNamespace(t.Context(), namespace.RootNamespace)

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

	resp, err := b.HandleRequest(t.Context(), configReq)
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
