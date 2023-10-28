// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const FixturePath = "../test-fixtures"

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig(filepath.Join(FixturePath, "config.hcl"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := &DefaultConfig{
		TokenHelper: "foo",
	}
	if !reflect.DeepEqual(expected, config) {
		t.Fatalf("bad: %#v", config)
	}
}

func TestLoadConfig_noExist(t *testing.T) {
	config, err := LoadConfig("nope/not-once/.never")
	if err != nil {
		t.Fatal(err)
	}

	if config.TokenHelper != "" {
		t.Errorf("expected %q to be %q", config.TokenHelper, "")
	}
}

func TestParseConfig_badKeys(t *testing.T) {
	_, err := ParseConfig(`
token_helper = "/token"
nope = "true"
`)
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), `invalid key "nope" on line 3`) {
		t.Errorf("bad error: %s", err.Error())
	}
}

func TestParseClientContext(t *testing.T) {
	conf := strings.TrimSpace(`
	current_context {
		name = "ns1"
		cluster_token = "foo"
		cluster_addr = "http://127.0.0.1:8200"
		namespace_path = "root"
	}
    client_context ctx1 {
	     name = "ns2"
		cluster_token = "bar"
		cluster_addr = "http://127.0.0.2:8200"
		namespace_path = "root"
	  } 
    client_context ctx2 {
	     name = "ns3"
		cluster_token = "baz"
		cluster_addr = "http://127.0.0.3:8200"
		namespace_path = "asdf"
	  }
	  `)
	expected := ClientContextConfig{
		CurrentContext: ContextInfo{
			Name:          "ns1",
			VaultAddr:     "http://127.0.0.1:8200",
			ClusterToken:  "foo",
			NamespacePath: "root",
		},
		ClientContexts: []ContextInfo{
			{
				Name:          "ns2",
				VaultAddr:     "http://127.0.0.2:8200",
				ClusterToken:  "bar",
				NamespacePath: "root",
			},
			{
				Name:          "ns3",
				VaultAddr:     "http://127.0.0.3:8200",
				ClusterToken:  "baz",
				NamespacePath: "asdf",
			},
		},
	}
	p, err := ParseClientContextConfig(conf)
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(expected.ClientContexts[:], func(i, j int) bool {
		return expected.ClientContexts[i].Name < expected.ClientContexts[j].Name
	})

	require.Equal(t, p, expected)
	path := "/tmp/config_context.hcl"
	defer os.RemoveAll(path)
	err = WriteClientContextConfig(path, expected)
	require.NoError(t, err)
}

func TestParseInvalidNumCurrentContext(t *testing.T) {
	conf := strings.TrimSpace(`
	current_context {
		name = "ns1"
		cluster_token = "foo"
		cluster_addr = "http://127.0.0.1:8200"
		namespace_path = "root"
	}
    current_context {
	     name = "ns2"
		cluster_token = "bar"
		cluster_addr = "http://127.0.0.2:8200"
		namespace_path = "root"
	  }
    client_context {
	     name = "ns3"
		cluster_token = "baz"
		cluster_addr = "http://127.0.0.3:8200"
		namespace_path = "asdf"
	  }
	  `)
	_, err := ParseClientContextConfig(conf)
	require.Error(t, err)
}
