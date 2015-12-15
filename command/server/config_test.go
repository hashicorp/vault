package server

import (
	"reflect"
	"testing"
	"time"
)

func TestLoadConfigFile(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config.hcl")
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
		},

		Backend: &Backend{
			Type:          "consul",
			AdvertiseAddr: "foo",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		HABackend: &Backend{
			Type:          "consul",
			AdvertiseAddr: "snafu",
			Config: map[string]string{
				"bar": "baz",
			},
		},

		Telemetry: &Telemetry{
			StatsdAddr:      "bar",
			StatsiteAddr:    "foo",
			DisableHostname: false,
		},

		DisableCache: true,
		DisableMlock: true,

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("bad: %#v", config)
	}
}

func TestLoadConfigFile_json(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config.hcl.json")
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
		},

		Backend: &Backend{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		Telemetry: &Telemetry{
			StatsiteAddr:    "baz",
			StatsdAddr:      "",
			DisableHostname: false,
		},

		MaxLeaseTTL:        10 * time.Hour,
		MaxLeaseTTLRaw:     "10h",
		DefaultLeaseTTL:    10 * time.Hour,
		DefaultLeaseTTLRaw: "10h",
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("bad: %#v", config)
	}
}

func TestLoadConfigFile_json2(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config2.hcl.json")
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
		},

		Backend: &Backend{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		HABackend: &Backend{
			Type: "consul",
			Config: map[string]string{
				"bar": "baz",
			},
		},

		Telemetry: &Telemetry{
			StatsiteAddr:    "foo",
			StatsdAddr:      "bar",
			DisableHostname: true,
		},
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("bad: %#v", config)
	}
}

func TestLoadConfigDir(t *testing.T) {
	config, err := LoadConfigDir("./test-fixtures/config-dir")
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

		Backend: &Backend{
			Type: "consul",
			Config: map[string]string{
				"foo": "bar",
			},
		},

		Telemetry: &Telemetry{
			StatsiteAddr:    "qux",
			StatsdAddr:      "baz",
			DisableHostname: true,
		},

		MaxLeaseTTL:     10 * time.Hour,
		DefaultLeaseTTL: 10 * time.Hour,
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("bad: %#v", config)
	}
}
