package server

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
)

func testConfigRaftRetryJoin(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/raft_retry_join.hcl")
	if err != nil {
		t.Fatal(err)
	}
	retryJoinConfig := `[{"leader_api_addr":"http://127.0.0.1:8200"},{"leader_api_addr":"http://127.0.0.2:8200"},{"leader_api_addr":"http://127.0.0.3:8200"}]` + "\n"
	expected := &Config{
		Listeners: []*Listener{
			{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:8200",
				},
			},
		},

		Storage: &Storage{
			Type: "raft",
			Config: map[string]string{
				"path":       "/storage/path/raft",
				"node_id":    "raft1",
				"retry_join": retryJoinConfig,
			},
		},
		DisableMlock:    true,
		DisableMlockRaw: true,
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("\nexpected: %#v\n actual:%#v\n", config, expected)
	}
}

func testLoadConfigFile_topLevel(t *testing.T, entropy *Entropy) {
	config, err := LoadConfigFile("./test-fixtures/config2.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
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

		Telemetry: &Telemetry{
			StatsdAddr:                 "bar",
			StatsiteAddr:               "foo",
			DisableHostname:            false,
			DogStatsDAddr:              "127.0.0.1:7254",
			DogStatsDTags:              []string{"tag_1:val_1", "tag_2:val_2"},
			PrometheusRetentionTime:    30 * time.Second,
			PrometheusRetentionTimeRaw: "30s",
		},

		DisableCache:    true,
		DisableCacheRaw: true,
		DisableMlock:    true,
		DisableMlockRaw: true,
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
		ClusterName:        "testcluster",

		PidFile: "./pidfile",

		APIAddr:     "top_level_api_addr",
		ClusterAddr: "top_level_cluster_addr",
	}
	if entropy != nil {
		expected.Entropy = entropy
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func testLoadConfigFile_json2(t *testing.T, entropy *Entropy) {
	config, err := LoadConfigFile("./test-fixtures/config2.hcl.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
				},
			},
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:444",
				},
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

		Telemetry: &Telemetry{
			StatsiteAddr:                       "foo",
			StatsdAddr:                         "bar",
			DisableHostname:                    true,
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
			PrometheusRetentionTimeRaw:         "30s",
		},
	}
	if entropy != nil {
		expected.Entropy = entropy
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func testParseEntropy(t *testing.T, oss bool) {
	var tests = []struct {
		inConfig   string
		outErr     error
		outEntropy Entropy
	}{
		{
			inConfig: `entropy "seal" {
				mode = "augmentation"
				}`,
			outErr:     nil,
			outEntropy: Entropy{Augmentation},
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

	var config Config

	for _, test := range tests {
		obj, _ := hcl.Parse(strings.TrimSpace(test.inConfig))
		list, _ := obj.Node.(*ast.ObjectList)
		objList := list.Filter("entropy")
		err := parseEntropy(&config, objList, "entropy")
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
			fmt.Printf("\n config.Entropy: %#v", config.Entropy)
			t.Fatalf("entropy config mismatch: expected %#v got %#v", test.outEntropy, *config.Entropy)
		}
	}
}

func testLoadConfigFile(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
				},
			},
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

		Telemetry: &Telemetry{
			StatsdAddr:              "bar",
			StatsiteAddr:            "foo",
			DisableHostname:         false,
			DogStatsDAddr:           "127.0.0.1:7254",
			DogStatsDTags:           []string{"tag_1:val_1", "tag_2:val_2"},
			PrometheusRetentionTime: prometheusDefaultRetentionTime,
			MetricsPrefix:           "myprefix",
		},

		DisableCache:             true,
		DisableCacheRaw:          true,
		DisableMlock:             true,
		DisableMlockRaw:          true,
		DisablePrintableCheckRaw: true,
		DisablePrintableCheck:    true,
		EnableUI:                 true,
		EnableUIRaw:              true,

		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,

		DisableSealWrap:    true,
		DisableSealWrapRaw: true,

		Entropy: nil,

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",
		ClusterName:        "testcluster",

		PidFile: "./pidfile",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func testLoadConfigFile_json(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config.hcl.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
				},
			},
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

		Telemetry: &Telemetry{
			StatsiteAddr:                       "baz",
			StatsdAddr:                         "",
			DisableHostname:                    false,
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
			PrometheusRetentionTime:            prometheusDefaultRetentionTime,
		},

		MaxLeaseTTL:          10 * time.Hour,
		MaxLeaseTTLRaw:       "10h",
		DefaultLeaseTTL:      10 * time.Hour,
		DefaultLeaseTTLRaw:   "10h",
		ClusterName:          "testcluster",
		DisableCacheRaw:      interface{}(nil),
		DisableMlockRaw:      interface{}(nil),
		EnableUI:             true,
		EnableUIRaw:          true,
		PidFile:              "./pidfile",
		EnableRawEndpoint:    true,
		EnableRawEndpointRaw: true,
		DisableSealWrap:      true,
		DisableSealWrapRaw:   true,
		Entropy:              nil,
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func testLoadConfigDir(t *testing.T) {
	config, err := LoadConfigDir("./test-fixtures/config-dir")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		DisableCache: true,
		DisableMlock: true,

		DisableClustering:    false,
		DisableClusteringRaw: false,

		APIAddr:     "https://vault.local",
		ClusterAddr: "https://127.0.0.1:444",

		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address": "127.0.0.1:443",
				},
			},
		},

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

		Telemetry: &Telemetry{
			StatsiteAddr:            "qux",
			StatsdAddr:              "baz",
			DisableHostname:         true,
			PrometheusRetentionTime: prometheusDefaultRetentionTime,
		},

		MaxLeaseTTL:     10 * time.Hour,
		DefaultLeaseTTL: 10 * time.Hour,
		ClusterName:     "testcluster",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func testConfig_Sanitized(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config3.hcl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	sanitizedConfig := config.Sanitized()

	expected := map[string]interface{}{
		"api_addr":                     "top_level_api_addr",
		"cache_size":                   0,
		"cluster_addr":                 "top_level_cluster_addr",
		"cluster_cipher_suites":        "",
		"cluster_name":                 "testcluster",
		"default_lease_ttl":            10 * time.Hour,
		"default_max_request_duration": 0 * time.Second,
		"disable_cache":                true,
		"disable_clustering":           false,
		"disable_indexing":             false,
		"disable_mlock":                true,
		"disable_performance_standby":  false,
		"disable_printable_check":      false,
		"disable_sealwrap":             true,
		"raw_storage_endpoint":         true,
		"enable_ui":                    true,
		"ha_storage": map[string]interface{}{
			"cluster_addr":       "top_level_cluster_addr",
			"disable_clustering": true,
			"redirect_addr":      "top_level_api_addr",
			"type":               "consul"},
		"listeners": []interface{}{
			map[string]interface{}{
				"config": map[string]interface{}{
					"address": "127.0.0.1:443",
				},
				"type": "tcp",
			},
		},
		"log_format":       "",
		"log_level":        "",
		"max_lease_ttl":    10 * time.Hour,
		"pid_file":         "./pidfile",
		"plugin_directory": "",
		"seals": []interface{}{
			map[string]interface{}{
				"disabled": false,
				"type":     "awskms",
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
			"statsite_address":                       ""},
	}

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
	tls_require_and_verify_client_cert = true
	tls_disable_client_certs = true
}`))

	var config Config
	list, _ := obj.Node.(*ast.ObjectList)
	objList := list.Filter("listener")
	parseListeners(&config, objList)
	listeners := config.Listeners
	if len(listeners) == 0 {
		t.Fatalf("expected at least one listener in the config")
	}
	listener := listeners[0]
	if listener.Type != "tcp" {
		t.Fatalf("expected tcp listener in the config")
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]interface{}{
					"address":                            "127.0.0.1:443",
					"cluster_address":                    "127.0.0.1:8201",
					"tls_disable":                        false,
					"tls_cert_file":                      "./certs/server.crt",
					"tls_key_file":                       "./certs/server.key",
					"tls_client_ca_file":                 "./certs/rootca.crt",
					"tls_min_version":                    "tls12",
					"tls_require_and_verify_client_cert": true,
					"tls_disable_client_certs":           true,
				},
			},
		},
	}

	if !reflect.DeepEqual(config, *expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, *expected)
	}

}
