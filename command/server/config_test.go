package server

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/helper/logformat"
	log "github.com/mgutz/logxi/v1"
)

func TestLoadConfigFile(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	config, err := LoadConfigFile("./test-fixtures/config.hcl", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "atlas",
				Config: map[string]string{
					"token":          "foobar",
					"infrastructure": "foo/bar",
					"endpoint":       "https://foo.bar:1111",
					"node_id":        "foo_node",
				},
			},
			&Listener{
				Type: "tcp",
				Config: map[string]string{
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

		Telemetry: &Telemetry{
			StatsdAddr:      "bar",
			StatsiteAddr:    "foo",
			DisableHostname: false,
		},

		DisableCache:    true,
		DisableCacheRaw: true,
		DisableMlock:    true,
		DisableMlockRaw: true,
		EnableUI:        true,
		EnableUIRaw:     true,

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",
		ClusterName:        "testcluster",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func TestLoadConfigFile_json(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	config, err := LoadConfigFile("./test-fixtures/config.hcl.json", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]string{
					"address": "127.0.0.1:443",
				},
			},
			&Listener{
				Type: "atlas",
				Config: map[string]string{
					"token":          "foobar",
					"infrastructure": "foo/bar",
					"endpoint":       "https://foo.bar:1111",
					"node_id":        "foo_node",
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
		},

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",
		ClusterName:        "testcluster",
		DisableCacheRaw:    interface{}(nil),
		DisableMlockRaw:    interface{}(nil),
		EnableUI:           true,
		EnableUIRaw:        true,
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func TestLoadConfigFile_json2(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	config, err := LoadConfigFile("./test-fixtures/config2.hcl.json", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]string{
					"address": "127.0.0.1:443",
				},
			},
			&Listener{
				Type: "tcp",
				Config: map[string]string{
					"address": "127.0.0.1:444",
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

		HAStorage: &Storage{
			Type: "consul",
			Config: map[string]string{
				"bar": "baz",
			},
		},

		CacheSize: 45678,

		EnableUI: true,

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
		},
	}
	if !reflect.DeepEqual(config, expected) {
	}
}

func TestLoadConfigDir(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	config, err := LoadConfigDir("./test-fixtures/config-dir", logger)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &Config{
		DisableCache: true,
		DisableMlock: true,

		Listeners: []*Listener{
			&Listener{
				Type: "tcp",
				Config: map[string]string{
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

		EnableUI: true,

		Telemetry: &Telemetry{
			StatsiteAddr:    "qux",
			StatsdAddr:      "baz",
			DisableHostname: true,
		},

		MaxLeaseTTL:     10 * time.Hour,
		DefaultLeaseTTL: 10 * time.Hour,
		ClusterName:     "testcluster",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("expected \n\n%#v\n\n to be \n\n%#v\n\n", config, expected)
	}
}

func TestParseConfig_badTopLevel(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	_, err := ParseConfig(strings.TrimSpace(`
backend {}
bad  = "one"
nope = "yes"
`), logger)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "invalid key 'bad' on line 2") {
		t.Errorf("bad error: %q", err)
	}

	if !strings.Contains(err.Error(), "invalid key 'nope' on line 3") {
		t.Errorf("bad error: %q", err)
	}
}

func TestParseConfig_badListener(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	_, err := ParseConfig(strings.TrimSpace(`
listener "tcp" {
	address = "1.2.3.3"
	bad  = "one"
	nope = "yes"
}
`), logger)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "listeners.tcp: invalid key 'bad' on line 3") {
		t.Errorf("bad error: %q", err)
	}

	if !strings.Contains(err.Error(), "listeners.tcp: invalid key 'nope' on line 4") {
		t.Errorf("bad error: %q", err)
	}
}

func TestParseConfig_badTelemetry(t *testing.T) {
	logger := logformat.NewVaultLogger(log.LevelTrace)

	_, err := ParseConfig(strings.TrimSpace(`
telemetry {
	statsd_address = "1.2.3.3"
	bad  = "one"
	nope = "yes"
}
`), logger)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "telemetry: invalid key 'bad' on line 3") {
		t.Errorf("bad error: %q", err)
	}

	if !strings.Contains(err.Error(), "telemetry: invalid key 'nope' on line 4") {
		t.Errorf("bad error: %q", err)
	}
}
