// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package dbplugin

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
)

// stubLookRunnerUtil is a non-functional pluginutil.LookRunnerUtil used only to
// supply a non-nil Sys. The embedded nil interface is never dereferenced because
// config validation rejects the input before any plugin lookup happens.
type stubLookRunnerUtil struct {
	pluginutil.LookRunnerUtil
}

func TestPluginFactoryConfig_Validate(t *testing.T) {
	type testCase struct {
		config PluginFactoryConfig

		// expectedErrContains is empty when the config is expected to be valid.
		expectedErrContains string
	}

	tests := map[string]testCase{
		"missing sys": {
			config: PluginFactoryConfig{
				PluginName: "postgresql-database-plugin",
				Logger:     hclog.NewNullLogger(),
			},
			expectedErrContains: "system view",
		},
		"missing logger": {
			config: PluginFactoryConfig{
				PluginName: "postgresql-database-plugin",
				Sys:        stubLookRunnerUtil{},
			},
			expectedErrContains: "logger",
		},
		"fully populated config is valid": {
			config: PluginFactoryConfig{
				PluginName: "postgresql-database-plugin",
				Sys:        stubLookRunnerUtil{},
				Logger:     hclog.NewNullLogger(),
			},
			expectedErrContains: "",
		},
		// An empty plugin name is left for the plugin catalog lookup to report,
		// because some callers rely on it to resolve a default plugin.
		"empty plugin name is not rejected here": {
			config: PluginFactoryConfig{
				Sys:    stubLookRunnerUtil{},
				Logger: hclog.NewNullLogger(),
			},
			expectedErrContains: "",
		},
		"telemetry labels are optional": {
			config: PluginFactoryConfig{
				PluginName:     "postgresql-database-plugin",
				Sys:            stubLookRunnerUtil{},
				Logger:         hclog.NewNullLogger(),
				Namespace:      "root",
				MountPoint:     "database/",
				ConnectionName: "prod-primary",
			},
			expectedErrContains: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.config.validate()

			if test.expectedErrContains == "" {
				if err != nil {
					t.Fatalf("Expected no error, but got: %s", err)
				}
				return
			}

			if err == nil {
				t.Fatal("Expected an error, but got none")
			}
			if !strings.Contains(err.Error(), test.expectedErrContains) {
				t.Fatalf("Expected error mentioning %q, but got: %s", test.expectedErrContains, err)
			}
		})
	}
}

// TestPluginFactoryWithConfig_InvalidConfigReturnsError asserts that an
// incomplete config produces an error from the exported factory rather than a
// nil dereference panic. The factory takes a struct literal, so a caller can
// omit a required field and still compile.
func TestPluginFactoryWithConfig_InvalidConfigReturnsError(t *testing.T) {
	tests := map[string]PluginFactoryConfig{
		"empty config":   {},
		"missing sys":    {PluginName: "postgresql-database-plugin", Logger: hclog.NewNullLogger()},
		"missing logger": {PluginName: "postgresql-database-plugin", Sys: stubLookRunnerUtil{}},
	}

	for name, config := range tests {
		t.Run(name, func(t *testing.T) {
			// A panic here is a test failure, which is the regression this
			// guards against.
			db, err := PluginFactoryWithConfig(context.Background(), config)
			if err == nil {
				t.Fatal("Expected an error, but got none")
			}
			if db != nil {
				t.Fatalf("Expected a nil Database on error, but got: %#v", db)
			}
		})
	}
}

// TestPluginFactoryVersion_RejectsNilInput asserts the positional wrappers
// inherit the same validation, since they previously panicked on nil input too.
func TestPluginFactoryVersion_RejectsNilInput(t *testing.T) {
	db, err := PluginFactoryVersion(context.Background(), "postgresql-database-plugin", "", nil, hclog.NewNullLogger())
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}
	if db != nil {
		t.Fatalf("Expected a nil Database on error, but got: %#v", db)
	}
}
