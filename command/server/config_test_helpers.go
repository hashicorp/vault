// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/go-test/deep"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	"github.com/hashicorp/vault/internalshared/configutil"
)

var DefaultCustomHeaders = map[string]map[string]string{
	"default": {
		"Strict-Transport-Security": configutil.StrictTransportSecurity,
	},
}

func boolPointer(x bool) *bool {
	return &x
}

func testConfigRaftRetryJoin(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/raft_retry_join.hcl")
	if err != nil {
		t.Fatal(err)
	}
	retryJoinConfig := `[{"leader_api_addr":"http://127.0.0.1:8200"},{"leader_api_addr":"http://127.0.0.2:8200"},{"leader_api_addr":"http://127.0.0.3:8200"}]`
	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:8200",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},
			DisableMlock: true,
		},

		Storage: &Storage{
			Type: "raft",
			Config: map[string]string{
				"path":       "/storage/path/raft",
				"node_id":    "raft1",
				"retry_join": retryJoinConfig,
			},
		},
	}
	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testLoadConfigFile_topLevel(t *testing.T, entropy *configutil.Entropy) {
	config, err := LoadConfigFile("./test-fixtures/config2.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:443",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},

			Telemetry: &configutil.Telemetry{
				StatsdAddr:                  "bar",
				StatsiteAddr:                "foo",
				DisableHostname:             false,
				DogStatsDAddr:               "127.0.0.1:7254",
				DogStatsDTags:               []string{"tag_1:val_1", "tag_2:val_2"},
				PrometheusRetentionTime:     30 * time.Second,
				UsageGaugePeriod:            5 * time.Minute,
				MaximumGaugeCardinality:     125,
				LeaseMetricsEpsilon:         time.Hour,
				NumLeaseMetricsTimeBuckets:  168,
				LeaseMetricsNameSpaceLabels: false,
			},

			DisableMlock: true,

			PidFile: "./pidfile",

			ClusterName: "testcluster",

			Seals: []*configutil.KMS{
				{
					Type: "nopurpose",
					Name: "nopurpose",
				},
				{
					Type:    "stringpurpose",
					Purpose: []string{"foo"},
					Name:    "stringpurpose",
				},
				{
					Type:    "commastringpurpose",
					Purpose: []string{"foo", "bar"},
					Name:    "commastringpurpose",
				},
				{
					Type:    "slicepurpose",
					Purpose: []string{"zip", "zap"},
					Name:    "slicepurpose",
				},
			},
		},

		Storage: &Storage{
			Type:         "consul",
			RedirectAddr: "top_level_api_addr",
			ClusterAddr:  "top_level_cluster_addr",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		HAStorage: &Storage{
			Type:         "consul",
			RedirectAddr: "top_level_api_addr",
			ClusterAddr:  "top_level_cluster_addr",
			Config: map[string]string{
				"bar": "baz",
			},
			DisableClustering: true,
		},

		ServiceRegistration: &ServiceRegistration{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		DisableCache:    true,
		DisableCacheRaw: true,
		EnableUI:        true,
		EnableUIRaw:     true,

		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,

		DisableSealWrap:    true,
		DisableSealWrapRaw: true,

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",

		APIAddr:     "top_level_api_addr",
		ClusterAddr: "top_level_cluster_addr",
	}
	addExpectedEntConfig(expected, []string{})

	if entropy != nil {
		expected.Entropy = entropy
	}
	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testLoadConfigFile_json2(t *testing.T, entropy *configutil.Entropy) {
	config, err := LoadConfigFile("./test-fixtures/config2.hcl.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:443",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:444",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},

			Telemetry: &configutil.Telemetry{
				StatsiteAddr:                       "foo",
				StatsdAddr:                         "bar",
				DisableHostname:                    true,
				UsageGaugePeriod:                   5 * time.Minute,
				MaximumGaugeCardinality:            125,
				CirconusAPIToken:                   "0",
				CirconusAPIApp:                     "vault",
				CirconusAPIURL:                     "http://api.circonus.com/v2",
				CirconusSubmissionInterval:         "10s",
				CirconusCheckSubmissionURL:         "https://someplace.com/metrics",
				CirconusCheckID:                    "0",
				CirconusCheckForceMetricActivation: "true",
				CirconusCheckInstanceID:            "node1:vault",
				CirconusCheckSearchTag:             "service:vault",
				CirconusCheckDisplayName:           "node1:vault",
				CirconusCheckTags:                  "cat1:tag1,cat2:tag2",
				CirconusBrokerID:                   "0",
				CirconusBrokerSelectTag:            "dc:sfo",
				PrometheusRetentionTime:            30 * time.Second,
				LeaseMetricsEpsilon:                time.Hour,
				NumLeaseMetricsTimeBuckets:         168,
				LeaseMetricsNameSpaceLabels:        false,
			},
		},

		Storage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		HAStorage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"bar": "baz",
			},
			DisableClustering: true,
		},

		ServiceRegistration: &ServiceRegistration{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		CacheSize: 45678,

		EnableUI:    true,
		EnableUIRaw: true,

		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,

		DisableSealWrap:    true,
		DisableSealWrapRaw: true,
	}
	addExpectedEntConfig(expected, []string{"http"})

	if entropy != nil {
		expected.Entropy = entropy
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testParseEntropy(t *testing.T, oss bool) {
	tests := []struct {
		inConfig   string
		outErr     error
		outEntropy configutil.Entropy
	}{
		{
			inConfig: `entropy "seal" {
				mode = "augmentation"
				}`,
			outErr:     nil,
			outEntropy: configutil.Entropy{Mode: configutil.EntropyAugmentation},
		},
		{
			inConfig: `entropy "seal" {
				mode = "a_mode_that_is_not_supported"
				}`,
			outErr: fmt.Errorf("the specified entropy mode %q is not supported", "a_mode_that_is_not_supported"),
		},
		{
			inConfig: `entropy "device_that_is_not_supported" {
				mode = "augmentation"
				}`,
			outErr: fmt.Errorf("only the %q type of external entropy is supported", "seal"),
		},
		{
			inConfig: `entropy "seal" {
				mode = "augmentation"
				}
				entropy "seal" {
				mode = "augmentation"
				}`,
			outErr: fmt.Errorf("only one %q block is permitted", "entropy"),
		},
	}

	config := Config{
		SharedConfig: &configutil.SharedConfig{},
	}

	for _, test := range tests {
		obj, _ := hcl.Parse(strings.TrimSpace(test.inConfig))
		list, _ := obj.Node.(*ast.ObjectList)
		objList := list.Filter("entropy")
		err := configutil.ParseEntropy(config.SharedConfig, objList, "entropy")
		// validate the error, both should be nil or have the same Error()
		switch {
		case oss:
			if config.Entropy != nil {
				t.Fatalf("parsing Entropy should not be possible in oss but got a non-nil config.Entropy: %#v", config.Entropy)
			}
		case err != nil && test.outErr != nil:
			if err.Error() != test.outErr.Error() {
				t.Fatalf("error mismatch: expected %#v got %#v", err, test.outErr)
			}
		case err != test.outErr:
			t.Fatalf("error mismatch: expected %#v got %#v", err, test.outErr)
		case err == nil && config.Entropy != nil && *config.Entropy != test.outEntropy:
			t.Fatalf("entropy config mismatch: expected %#v got %#v", test.outEntropy, *config.Entropy)
		}
	}
}

func testLoadConfigFileIntegerAndBooleanValues(t *testing.T) {
	testLoadConfigFileIntegerAndBooleanValuesCommon(t, "./test-fixtures/config4.hcl")
}

func testLoadConfigFileIntegerAndBooleanValuesJson(t *testing.T) {
	testLoadConfigFileIntegerAndBooleanValuesCommon(t, "./test-fixtures/config4.hcl.json")
}

func testLoadConfigFileIntegerAndBooleanValuesCommon(t *testing.T, path string) {
	config, err := LoadConfigFile(path)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:8200",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},
			DisableMlock: true,
		},

		Storage: &Storage{
			Type: "raft",
			Config: map[string]string{
				"path":                   "/storage/path/raft",
				"node_id":                "raft1",
				"performance_multiplier": "1",
				"foo":                    "bar",
				"baz":                    "true",
			},
			ClusterAddr: "127.0.0.1:8201",
		},

		ClusterAddr: "127.0.0.1:8201",

		DisableCache:    true,
		DisableCacheRaw: true,
		EnableUI:        true,
		EnableUIRaw:     true,
	}

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testLoadConfigFile(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:443",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},

			Telemetry: &configutil.Telemetry{
				StatsdAddr:                  "bar",
				StatsiteAddr:                "foo",
				DisableHostname:             false,
				UsageGaugePeriod:            5 * time.Minute,
				MaximumGaugeCardinality:     100,
				DogStatsDAddr:               "127.0.0.1:7254",
				DogStatsDTags:               []string{"tag_1:val_1", "tag_2:val_2"},
				PrometheusRetentionTime:     configutil.PrometheusDefaultRetentionTime,
				MetricsPrefix:               "myprefix",
				LeaseMetricsEpsilon:         time.Hour,
				NumLeaseMetricsTimeBuckets:  168,
				LeaseMetricsNameSpaceLabels: false,
			},

			DisableMlock: true,

			Entropy: nil,

			PidFile: "./pidfile",

			ClusterName: "testcluster",
		},

		Storage: &Storage{
			Type:         "consul",
			RedirectAddr: "foo",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		HAStorage: &Storage{
			Type:         "consul",
			RedirectAddr: "snafu",
			Config: map[string]string{
				"bar": "baz",
			},
			DisableClustering: true,
		},

		ServiceRegistration: &ServiceRegistration{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		DisableCache:             true,
		DisableCacheRaw:          true,
		DisablePrintableCheckRaw: true,
		DisablePrintableCheck:    true,
		EnableUI:                 true,
		EnableUIRaw:              true,

		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,

		EnableIntrospectionEndpoint:    true,
		EnableIntrospectionEndpointRaw: true,

		DisableSealWrap:    true,
		DisableSealWrapRaw: true,

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",

		EnableResponseHeaderHostname:      true,
		EnableResponseHeaderHostnameRaw:   true,
		EnableResponseHeaderRaftNodeID:    true,
		EnableResponseHeaderRaftNodeIDRaw: true,

		LicensePath: "/path/to/license",
	}

	addExpectedEntConfig(expected, []string{})

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testUnknownFieldValidationStorageAndListener(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/storage-listener-config.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(config.UnusedKeys) != 0 {
		t.Fatalf("unused keys for valid config are %+v\n", config.UnusedKeys)
	}
}

func testUnknownFieldValidation(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []configutil.ConfigError{
		{
			Problem: "unknown or unsupported field bad_value found in configuration",
			Position: token.Pos{
				Filename: "./test-fixtures/config.hcl",
				Offset:   652,
				Line:     37,
				Column:   5,
			},
		},
	}
	errors := config.Validate("./test-fixtures/config.hcl")

	for _, er1 := range errors {
		found := false
		if strings.Contains(er1.String(), "sentinel") {
			// This happens on OSS, and is fine
			continue
		}
		for _, ex := range expected {
			// TODO: Only test the string, pos may change
			if ex.Problem == er1.Problem && reflect.DeepEqual(ex.Position, er1.Position) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("found unexpected error: %v", er1.String())
		}
	}
	for _, ex := range expected {
		found := false
		for _, er1 := range errors {
			if ex.Problem == er1.Problem && reflect.DeepEqual(ex.Position, er1.Position) {
				found = true
			}
		}
		if !found {
			t.Fatalf("could not find expected error: %v", ex.String())
		}
	}
}

// testUnknownFieldValidationJson tests that this valid json config does not result in
// errors. Prior to VAULT-8519, it reported errors even with a valid config that was
// parsed properly.
func testUnknownFieldValidationJson(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config_small.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	errors := config.Validate("./test-fixtures/config_small.json")
	if errors != nil {
		t.Fatal(errors)
	}
}

// testUnknownFieldValidationHcl tests that this valid hcl config does not result in
// errors. Prior to VAULT-8519, the json version of this config reported errors even
// with a valid config that was parsed properly.
// In short, this ensures the same for HCL as we test in testUnknownFieldValidationJson
func testUnknownFieldValidationHcl(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config_small.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	errors := config.Validate("./test-fixtures/config_small.hcl")
	if errors != nil {
		t.Fatal(errors)
	}
}

// testConfigWithAdministrativeNamespaceJson tests that a config with a valid administrative namespace path is correctly validated and loaded.
func testConfigWithAdministrativeNamespaceJson(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config_with_valid_admin_ns.json")
	require.NoError(t, err)

	configErrors := config.Validate("./test-fixtures/config_with_valid_admin_ns.json")
	require.Empty(t, configErrors)

	require.NotEmpty(t, config.AdministrativeNamespacePath)
}

// testConfigWithAdministrativeNamespaceHcl tests that a config with a valid administrative namespace path is correctly validated and loaded.
func testConfigWithAdministrativeNamespaceHcl(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config_with_valid_admin_ns.hcl")
	require.NoError(t, err)

	configErrors := config.Validate("./test-fixtures/config_with_valid_admin_ns.hcl")
	require.Empty(t, configErrors)

	require.NotEmpty(t, config.AdministrativeNamespacePath)
}

func testLoadConfigFile_json(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config.hcl.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:443",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},

			Telemetry: &configutil.Telemetry{
				StatsiteAddr:                       "baz",
				StatsdAddr:                         "",
				DisableHostname:                    false,
				UsageGaugePeriod:                   5 * time.Minute,
				MaximumGaugeCardinality:            100,
				CirconusAPIToken:                   "",
				CirconusAPIApp:                     "",
				CirconusAPIURL:                     "",
				CirconusSubmissionInterval:         "",
				CirconusCheckSubmissionURL:         "",
				CirconusCheckID:                    "",
				CirconusCheckForceMetricActivation: "",
				CirconusCheckInstanceID:            "",
				CirconusCheckSearchTag:             "",
				CirconusCheckDisplayName:           "",
				CirconusCheckTags:                  "",
				CirconusBrokerID:                   "",
				CirconusBrokerSelectTag:            "",
				PrometheusRetentionTime:            configutil.PrometheusDefaultRetentionTime,
				LeaseMetricsEpsilon:                time.Hour,
				NumLeaseMetricsTimeBuckets:         168,
				LeaseMetricsNameSpaceLabels:        false,
			},

			PidFile:     "./pidfile",
			Entropy:     nil,
			ClusterName: "testcluster",
		},

		Storage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
			DisableClustering: true,
		},

		ServiceRegistration: &ServiceRegistration{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		ClusterCipherSuites: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA",

		MaxLeaseTTL:          10 * time.Hour,
		MaxLeaseTTLRaw:       "10h",
		DefaultLeaseTTL:      10 * time.Hour,
		DefaultLeaseTTLRaw:   "10h",
		DisableCacheRaw:      interface{}(nil),
		EnableUI:             true,
		EnableUIRaw:          true,
		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,
		DisableSealWrap:      true,
		DisableSealWrapRaw:   true,
	}

	addExpectedEntConfig(expected, []string{})

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testLoadConfigDir(t *testing.T) {
	config, err := LoadConfigDir("./test-fixtures/config-dir")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			DisableMlock: true,

			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:443",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},

			Telemetry: &configutil.Telemetry{
				StatsiteAddr:                "qux",
				StatsdAddr:                  "baz",
				DisableHostname:             true,
				UsageGaugePeriod:            5 * time.Minute,
				MaximumGaugeCardinality:     100,
				PrometheusRetentionTime:     configutil.PrometheusDefaultRetentionTime,
				LeaseMetricsEpsilon:         time.Hour,
				NumLeaseMetricsTimeBuckets:  168,
				LeaseMetricsNameSpaceLabels: false,
			},
			ClusterName: "testcluster",
		},

		DisableCache:         true,
		DisableClustering:    false,
		DisableClusteringRaw: false,

		APIAddr:     "https://vault.local",
		ClusterAddr: "https://127.0.0.1:444",

		Storage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
			RedirectAddr:      "https://vault.local",
			ClusterAddr:       "https://127.0.0.1:444",
			DisableClustering: false,
		},

		EnableUI: true,

		EnableRawEndpoint: true,

		MaxLeaseTTL:     10 * time.Hour,
		DefaultLeaseTTL: 10 * time.Hour,
	}

	addExpectedEntConfig(expected, []string{"http"})

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testConfig_Sanitized(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config3.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	sanitizedConfig := config.Sanitized()

	expected := map[string]interface{}{
		"api_addr":                            "top_level_api_addr",
		"cache_size":                          0,
		"cluster_addr":                        "top_level_cluster_addr",
		"cluster_cipher_suites":               "",
		"cluster_name":                        "testcluster",
		"default_lease_ttl":                   (365 * 24 * time.Hour) / time.Second,
		"default_max_request_duration":        0 * time.Second,
		"disable_cache":                       true,
		"disable_clustering":                  false,
		"disable_indexing":                    false,
		"disable_mlock":                       true,
		"disable_performance_standby":         false,
		"experiments":                         []string(nil),
		"plugin_file_uid":                     0,
		"plugin_file_permissions":             0,
		"disable_printable_check":             false,
		"disable_sealwrap":                    true,
		"raw_storage_endpoint":                true,
		"introspection_endpoint":              false,
		"disable_sentinel_trace":              true,
		"detect_deadlocks":                    "",
		"enable_ui":                           true,
		"enable_response_header_hostname":     false,
		"enable_response_header_raft_node_id": false,
		"log_requests_level":                  "basic",
		"ha_storage": map[string]interface{}{
			"cluster_addr":       "top_level_cluster_addr",
			"disable_clustering": true,
			"redirect_addr":      "top_level_api_addr",
			"type":               "consul",
		},
		"listeners": []interface{}{
			map[string]interface{}{
				"config": map[string]interface{}{
					"address":          "127.0.0.1:443",
					"chroot_namespace": "admin/",
				},
				"type": "tcp",
			},
		},
		"log_format":       "",
		"log_level":        "",
		"max_lease_ttl":    (30 * 24 * time.Hour) / time.Second,
		"pid_file":         "./pidfile",
		"plugin_directory": "",
		"seals": []interface{}{
			map[string]interface{}{
				"disabled": false,
				"type":     "awskms",
				"name":     "awskms",
			},
		},
		"storage": map[string]interface{}{
			"cluster_addr":       "top_level_cluster_addr",
			"disable_clustering": false,
			"redirect_addr":      "top_level_api_addr",
			"type":               "consul",
		},
		"service_registration": map[string]interface{}{
			"type": "consul",
		},
		"telemetry": map[string]interface{}{
			"usage_gauge_period":                     5 * time.Minute,
			"maximum_gauge_cardinality":              100,
			"circonus_api_app":                       "",
			"circonus_api_token":                     "",
			"circonus_api_url":                       "",
			"circonus_broker_id":                     "",
			"circonus_broker_select_tag":             "",
			"circonus_check_display_name":            "",
			"circonus_check_force_metric_activation": "",
			"circonus_check_id":                      "",
			"circonus_check_instance_id":             "",
			"circonus_check_search_tag":              "",
			"circonus_submission_url":                "",
			"circonus_check_tags":                    "",
			"circonus_submission_interval":           "",
			"disable_hostname":                       false,
			"metrics_prefix":                         "pfx",
			"dogstatsd_addr":                         "",
			"dogstatsd_tags":                         []string(nil),
			"prometheus_retention_time":              24 * time.Hour,
			"stackdriver_location":                   "",
			"stackdriver_namespace":                  "",
			"stackdriver_project_id":                 "",
			"stackdriver_debug_logs":                 false,
			"statsd_address":                         "bar",
			"statsite_address":                       "",
			"lease_metrics_epsilon":                  time.Hour,
			"num_lease_metrics_buckets":              168,
			"add_lease_metrics_namespace_labels":     false,
			"add_mount_point_rollback_metrics":       false,
		},
		"administrative_namespace_path": "admin/",
		"imprecise_lease_role_tracking": false,
	}

	addExpectedEntSanitizedConfig(expected, []string{"http"})

	config.Prune()
	if diff := deep.Equal(sanitizedConfig, expected); len(diff) > 0 {
		t.Fatalf("bad, diff: %#v", diff)
	}
}

func testParseListeners(t *testing.T) {
	obj, _ := hcl.Parse(strings.TrimSpace(`
listener "tcp" {
  address = "127.0.0.1:443"
  cluster_address = "127.0.0.1:8201"
  tls_disable = false
  tls_cert_file = "./certs/server.crt"
  tls_key_file = "./certs/server.key"
  tls_client_ca_file = "./certs/rootca.crt"
  tls_min_version = "tls12"
  tls_max_version = "tls13"
  tls_require_and_verify_client_cert = true
  tls_disable_client_certs = true
  telemetry {
    unauthenticated_metrics_access = true
  }
  profiling {
    unauthenticated_pprof_access = true
  }
  agent_api {
    enable_quit = true
  }
  proxy_api {
    enable_quit = true
  }
  chroot_namespace = "admin"
}`))

	config := Config{
		SharedConfig: &configutil.SharedConfig{},
	}
	list, _ := obj.Node.(*ast.ObjectList)
	objList := list.Filter("listener")
	configutil.ParseListeners(config.SharedConfig, objList)
	listeners := config.Listeners
	if len(listeners) == 0 {
		t.Fatalf("expected at least one listener in the config")
	}
	listener := listeners[0]
	if listener.Type != "tcp" {
		t.Fatalf("expected tcp listener in the config")
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                          "tcp",
					Address:                       "127.0.0.1:443",
					ClusterAddress:                "127.0.0.1:8201",
					TLSCertFile:                   "./certs/server.crt",
					TLSKeyFile:                    "./certs/server.key",
					TLSClientCAFile:               "./certs/rootca.crt",
					TLSMinVersion:                 "tls12",
					TLSMaxVersion:                 "tls13",
					TLSRequireAndVerifyClientCert: true,
					TLSDisableClientCerts:         true,
					Telemetry: configutil.ListenerTelemetry{
						UnauthenticatedMetricsAccess: true,
					},
					Profiling: configutil.ListenerProfiling{
						UnauthenticatedPProfAccess: true,
					},
					AgentAPI: &configutil.AgentAPI{
						EnableQuit: true,
					},
					ProxyAPI: &configutil.ProxyAPI{
						EnableQuit: true,
					},
					CustomResponseHeaders: DefaultCustomHeaders,
					ChrootNamespace:       "admin/",
				},
			},
		},
	}
	config.Prune()
	if diff := deep.Equal(config, *expected); diff != nil {
		t.Fatal(diff)
	}
}

func testParseUserLockouts(t *testing.T) {
	obj, _ := hcl.Parse(strings.TrimSpace(`
	user_lockout "all" {
		lockout_duration = "40m"
		lockout_counter_reset = "45m"
		disable_lockout = "false"
	}
	  user_lockout "userpass" {
	     lockout_threshold = "100"
	     lockout_duration = "20m"
	  }
	  user_lockout "ldap" {
		disable_lockout = "true"
	 }`))

	config := Config{
		SharedConfig: &configutil.SharedConfig{},
	}
	list, _ := obj.Node.(*ast.ObjectList)
	objList := list.Filter("user_lockout")
	configutil.ParseUserLockouts(config.SharedConfig, objList)

	sort.Slice(config.SharedConfig.UserLockouts[:], func(i, j int) bool {
		return config.SharedConfig.UserLockouts[i].Type < config.SharedConfig.UserLockouts[j].Type
	})

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			UserLockouts: []*configutil.UserLockout{
				{
					Type:                "all",
					LockoutThreshold:    5,
					LockoutDuration:     2400000000000,
					LockoutCounterReset: 2700000000000,
					DisableLockout:      false,
				},
				{
					Type:                "userpass",
					LockoutThreshold:    100,
					LockoutDuration:     1200000000000,
					LockoutCounterReset: 2700000000000,
					DisableLockout:      false,
				},
				{
					Type:                "ldap",
					LockoutThreshold:    5,
					LockoutDuration:     2400000000000,
					LockoutCounterReset: 2700000000000,
					DisableLockout:      true,
				},
			},
		},
	}

	sort.Slice(expected.SharedConfig.UserLockouts[:], func(i, j int) bool {
		return expected.SharedConfig.UserLockouts[i].Type < expected.SharedConfig.UserLockouts[j].Type
	})
	config.Prune()
	require.Equal(t, config, *expected)
}

func testParseSockaddrTemplate(t *testing.T) {
	config, err := ParseConfig(`
api_addr = <<EOF
{{- GetAllInterfaces | include "flags" "loopback" | include "type" "ipv4" | attr "address" -}}
EOF

listener "tcp" {
	address = <<EOF
{{- GetAllInterfaces | include "flags" "loopback" | include "type" "ipv4" | attr "address" -}}:443
EOF
	cluster_address = <<EOF
{{- GetAllInterfaces | include "flags" "loopback" | include "type" "ipv4" | attr "address" -}}:8201
EOF
	tls_disable = true
}`, "")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		APIAddr: "127.0.0.1",
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:443",
					ClusterAddress:        "127.0.0.1:8201",
					TLSDisable:            true,
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},
		},
	}
	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testParseStorageTemplate(t *testing.T) {
	config, err := ParseConfig(`
storage "consul" {

	disable_registration = false
	path = "tmp/"

}
ha_storage "consul" {
	tls_skip_verify = true
	scheme = "http"
	max_parallel = 128
}

`, "")
	if err != nil {
		t.Fatal(err)
	}

	expected := &Config{
		Storage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"disable_registration": "false",
				"path":                 "tmp/",
			},
		},
		HAStorage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"tls_skip_verify": "true",
				"scheme":          "http",
				"max_parallel":    "128",
			},
		},
		SharedConfig: &configutil.SharedConfig{},
	}
	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testParseSeals(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config_seals.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	config.Listeners[0].RawConfig = nil

	expected := &Config{
		Storage: &Storage{
			Type:   "consul",
			Config: map[string]string{},
		},
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:443",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},
			Seals: []*configutil.KMS{
				{
					Type:    "pkcs11",
					Purpose: []string{"many", "purposes"},
					Config: map[string]string{
						"lib":                    "/usr/lib/libcklog2.so",
						"slot":                   "0.0",
						"pin":                    "XXXXXXXX",
						"key_label":              "HASHICORP",
						"mechanism":              "0x1082",
						"hmac_mechanism":         "0x0251",
						"hmac_key_label":         "vault-hsm-hmac-key",
						"default_hmac_key_label": "vault-hsm-hmac-key",
						"generate_key":           "true",
					},
					Name: "pkcs11",
				},
				{
					Type:     "pkcs11",
					Purpose:  []string{"single"},
					Disabled: true,
					Config: map[string]string{
						"lib":                    "/usr/lib/libcklog2.so",
						"slot":                   "0.0",
						"pin":                    "XXXXXXXX",
						"key_label":              "HASHICORP",
						"mechanism":              "4226",
						"hmac_mechanism":         "593",
						"hmac_key_label":         "vault-hsm-hmac-key",
						"default_hmac_key_label": "vault-hsm-hmac-key",
						"generate_key":           "true",
					},
					Name: "pkcs11-disabled",
				},
			},
		},
	}
	addExpectedDefaultEntConfig(expected)
	config.Prune()
	require.Equal(t, config, expected)
}

func testLoadConfigFileLeaseMetrics(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config5.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:                  "tcp",
					Address:               "127.0.0.1:443",
					CustomResponseHeaders: DefaultCustomHeaders,
				},
			},

			Telemetry: &configutil.Telemetry{
				StatsdAddr:                  "bar",
				StatsiteAddr:                "foo",
				DisableHostname:             false,
				UsageGaugePeriod:            5 * time.Minute,
				MaximumGaugeCardinality:     100,
				DogStatsDAddr:               "127.0.0.1:7254",
				DogStatsDTags:               []string{"tag_1:val_1", "tag_2:val_2"},
				PrometheusRetentionTime:     configutil.PrometheusDefaultRetentionTime,
				MetricsPrefix:               "myprefix",
				LeaseMetricsEpsilon:         time.Hour,
				NumLeaseMetricsTimeBuckets:  2,
				LeaseMetricsNameSpaceLabels: true,
			},

			DisableMlock: true,

			Entropy: nil,

			PidFile: "./pidfile",

			ClusterName: "testcluster",
		},

		Storage: &Storage{
			Type:         "consul",
			RedirectAddr: "foo",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		HAStorage: &Storage{
			Type:         "consul",
			RedirectAddr: "snafu",
			Config: map[string]string{
				"bar": "baz",
			},

			DisableClustering: true,
		},

		ServiceRegistration: &ServiceRegistration{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		DisableCache:             true,
		DisableCacheRaw:          true,
		DisablePrintableCheckRaw: true,
		DisablePrintableCheck:    true,
		EnableUI:                 true,
		EnableUIRaw:              true,

		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,

		DisableSealWrap:    true,
		DisableSealWrapRaw: true,

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",
	}

	addExpectedEntConfig(expected, []string{})

	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}

func testConfigRaftAutopilot(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/raft_autopilot.hcl")
	if err != nil {
		t.Fatal(err)
	}

	autopilotConfig := `[{"cleanup_dead_servers":true,"last_contact_threshold":"500ms","max_trailing_logs":250,"min_quorum":3,"server_stabilization_time":"10s"}]`
	expected := &Config{
		SharedConfig: &configutil.SharedConfig{
			Listeners: []*configutil.Listener{
				{
					Type:    "tcp",
					Address: "127.0.0.1:8200",
				},
			},
			DisableMlock: true,
		},

		Storage: &Storage{
			Type: "raft",
			Config: map[string]string{
				"path":      "/storage/path/raft",
				"node_id":   "raft1",
				"autopilot": autopilotConfig,
			},
		},
	}
	config.Prune()
	if diff := deep.Equal(config, expected); diff != nil {
		t.Fatal(diff)
	}
}
