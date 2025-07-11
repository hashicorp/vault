// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cliconfig

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, duplicate, err := loadConfig(filepath.Join("testdata", "config.hcl"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &defaultConfig{
		TokenHelper: "foo",
	}
	if !reflect.DeepEqual(expected, config) {
		t.Fatalf("bad: %#v", config)
	}

	if duplicate {
		t.Fatal("expected no duplicate")
	}
}

func TestLoadConfig_noExist(t *testing.T) {
	config, duplicate, err := loadConfig("nope/not-once/.never")
	if err != nil {
		t.Fatal(err)
	}

	if config.TokenHelper != "" {
		t.Errorf("expected %q to be %q", config.TokenHelper, "")
	}

	if duplicate {
		t.Fatal("expected no duplicate")
	}
}

func TestParseConfig_badKeys(t *testing.T) {
	_, duplicate, err := parseConfig(`
token_helper = "/token"
nope = "true"
`)
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), `invalid key "nope" on line 3`) {
		t.Errorf("bad error: %s", err.Error())
	}

	if duplicate {
		t.Fatal("expected no duplicate")
	}
}

// TestParseConfig_HclDuplicateKey tests the parsing of HCL files with duplicate keys.
// TODO (HCL_DUP_KEYS_DEPRECATION): on full removal change this test to ensure that duplicate attributes cannot be parsed
// under any circumstances.
func TestParseConfig_HclDuplicateKey(t *testing.T) {
	t.Run("fail parsing without env var", func(t *testing.T) {
		_, _, err := parseConfig(`
token_helper = "/token"
token_helper = "/token"
`)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("fail parsing with env var set to false", func(t *testing.T) {
		t.Setenv(allowHclDuplicatesEnvVar, "false")
		_, _, err := parseConfig(`
token_helper = "/token"
token_helper = "/token"
`)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("succeed parsing with env var set to true", func(t *testing.T) {
		t.Setenv(allowHclDuplicatesEnvVar, "true")
		_, duplicate, err := parseConfig(`
token_helper = "/token"
token_helper = "/token"
`)
		if err != nil {
			t.Fatal("expected no error")
		}

		if !duplicate {
			t.Fatal("expected duplicate")
		}
	})
}
