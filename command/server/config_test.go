package server

import (
	"reflect"
	"testing"
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

		Telemetry: &Telemetry{
			StatsdAddr: "bar",
			StatsiteAddr: "foo",
			DisableHostname: false,
		},

		DisableMlock: true,
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
			StatsiteAddr: "baz",
			StatsdAddr: "",
			DisableHostname: false,
		},
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

		Telemetry: &Telemetry{
			StatsiteAddr: "foo",
			StatsdAddr: "bar",
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
			StatsiteAddr: "qux",
			StatsdAddr: "baz",
			DisableHostname: true,
		},
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("bad: %#v", config)
	}
}
