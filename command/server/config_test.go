// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/vault/internalshared/configutil"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFile(t *testing.T) {
	testLoadConfigFile(t)
}

func TestLoadConfigFile_json(t *testing.T) {
	testLoadConfigFile_json(t)
}

func TestLoadConfigFileIntegerAndBooleanValues(t *testing.T) {
	testLoadConfigFileIntegerAndBooleanValues(t)
}

func TestLoadConfigFileIntegerAndBooleanValuesJson(t *testing.T) {
	testLoadConfigFileIntegerAndBooleanValuesJson(t)
}

func TestLoadConfigFileWithLeaseMetricTelemetry(t *testing.T) {
	testLoadConfigFileLeaseMetrics(t)
}

func TestLoadConfigDir(t *testing.T) {
	testLoadConfigDir(t)
}

func TestConfig_Sanitized(t *testing.T) {
	testConfig_Sanitized(t)
}

func TestParseListeners(t *testing.T) {
	testParseListeners(t)
}

func TestParseUserLockouts(t *testing.T) {
	testParseUserLockouts(t)
}

func TestParseSockaddrTemplate(t *testing.T) {
	testParseSockaddrTemplate(t)
}

func TestConfigRaftRetryJoin(t *testing.T) {
	testConfigRaftRetryJoin(t)
}

func TestParseSeals(t *testing.T) {
	testParseSeals(t)
}

func TestParseStorage(t *testing.T) {
	testParseStorageTemplate(t)
}

// TestConfigWithAdministrativeNamespace tests that .hcl and .json configurations are correctly parsed when the administrative_namespace_path is present.
func TestConfigWithAdministrativeNamespace(t *testing.T) {
	testConfigWithAdministrativeNamespaceHcl(t)
	testConfigWithAdministrativeNamespaceJson(t)
}

func TestUnknownFieldValidation(t *testing.T) {
	testUnknownFieldValidation(t)
}

func TestUnknownFieldValidationJson(t *testing.T) {
	testUnknownFieldValidationJson(t)
}

func TestUnknownFieldValidationHcl(t *testing.T) {
	testUnknownFieldValidationHcl(t)
}

func TestUnknownFieldValidationListenerAndStorage(t *testing.T) {
	testUnknownFieldValidationStorageAndListener(t)
}

func TestExperimentsConfigParsing(t *testing.T) {
	const envKey = "VAULT_EXPERIMENTS"
	originalValue := validExperiments
	validExperiments = []string{"foo", "bar", "baz"}
	t.Cleanup(func() {
		validExperiments = originalValue
	})

	for name, tc := range map[string]struct {
		fromConfig    []string
		fromEnv       []string
		fromCLI       []string
		expected      []string
		expectedError string
	}{
		// Multiple sources.
		"duplication":  {[]string{"foo"}, []string{"foo"}, []string{"foo"}, []string{"foo"}, ""},
		"disjoint set": {[]string{"foo"}, []string{"bar"}, []string{"baz"}, []string{"foo", "bar", "baz"}, ""},

		// Single source.
		"config only": {[]string{"foo"}, nil, nil, []string{"foo"}, ""},
		"env only":    {nil, []string{"foo"}, nil, []string{"foo"}, ""},
		"CLI only":    {nil, nil, []string{"foo"}, []string{"foo"}, ""},

		// Validation errors.
		"config invalid": {[]string{"invalid"}, nil, nil, nil, "from config"},
		"env invalid":    {nil, []string{"invalid"}, nil, nil, "from environment variable"},
		"CLI invalid":    {nil, nil, []string{"invalid"}, nil, "from command line flag"},
	} {
		t.Run(name, func(t *testing.T) {
			var configString string
			t.Setenv(envKey, strings.Join(tc.fromEnv, ","))
			if len(tc.fromConfig) != 0 {
				configString = fmt.Sprintf("experiments = [\"%s\"]", strings.Join(tc.fromConfig, "\", \""))
			}
			config, err := ParseConfig(configString, "")
			if err == nil {
				err = ExperimentsFromEnvAndCLI(config, envKey, tc.fromCLI)
			}

			switch tc.expectedError {
			case "":
				if err != nil {
					t.Fatal(err)
				}

			default:
				if err == nil || !strings.Contains(err.Error(), tc.expectedError) {
					t.Fatalf("Expected error to contain %q, but got: %s", tc.expectedError, err)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	originalValue := validExperiments
	for name, tc := range map[string]struct {
		validSet    []string
		input       []string
		expectError bool
	}{
		// Valid cases
		"minimal valid": {[]string{"foo"}, []string{"foo"}, false},
		"valid subset":  {[]string{"foo", "bar"}, []string{"bar"}, false},
		"repeated":      {[]string{"foo"}, []string{"foo", "foo"}, false},

		// Error cases
		"partially valid":      {[]string{"foo", "bar"}, []string{"foo", "baz"}, true},
		"empty":                {[]string{"foo"}, []string{""}, true},
		"no valid experiments": {[]string{}, []string{"foo"}, true},
	} {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(func() {
				validExperiments = originalValue
			})

			validExperiments = tc.validSet
			err := validateExperiments(tc.input)
			if tc.expectError && err == nil {
				t.Fatal("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Fatal("Did not expect error but got", err)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	for name, tc := range map[string]struct {
		left     []string
		right    []string
		expected []string
	}{
		"disjoint":    {[]string{"foo"}, []string{"bar"}, []string{"foo", "bar"}},
		"empty left":  {[]string{}, []string{"foo"}, []string{"foo"}},
		"empty right": {[]string{"foo"}, []string{}, []string{"foo"}},
		"overlapping": {[]string{"foo", "bar"}, []string{"foo", "baz"}, []string{"foo", "bar", "baz"}},
	} {
		t.Run(name, func(t *testing.T) {
			result := mergeExperiments(tc.left, tc.right)
			if !reflect.DeepEqual(tc.expected, result) {
				t.Fatalf("Expected %v but got %v", tc.expected, result)
			}
		})
	}
}

// Test_parseDevTLSConfig verifies that both Windows and Unix directories are correctly escaped when creating a dev TLS
// configuration in HCL
func Test_parseDevTLSConfig(t *testing.T) {
	tests := []struct {
		name          string
		certDirectory string
	}{
		{
			name:          "windows path",
			certDirectory: `C:\Users\ADMINI~1\AppData\Local\Temp\2\vault-tls4169358130`,
		},
		{
			name:          "unix path",
			certDirectory: "/tmp/vault-tls4169358130",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseDevTLSConfig("file", tt.certDirectory)
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf("%s/%s", tt.certDirectory, VaultDevCertFilename), cfg.Listeners[0].TLSCertFile)
			require.Equal(t, fmt.Sprintf("%s/%s", tt.certDirectory, VaultDevKeyFilename), cfg.Listeners[0].TLSKeyFile)
		})
	}
}

func TestCheckConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name:        "no-seals-configured",
			config:      &Config{SharedConfig: &configutil.SharedConfig{Seals: []*configutil.KMS{}}},
			expectError: false,
		},
		{
			name: "seal-with-empty-name",
			config: &Config{SharedConfig: &configutil.SharedConfig{
				Seals: []*configutil.KMS{
					{
						Type:     "awskms",
						Disabled: false,
					},
				},
			}},
			expectError: true,
		},
		{
			name: "seals-with-unique-names",
			config: &Config{SharedConfig: &configutil.SharedConfig{
				Seals: []*configutil.KMS{
					{
						Type:     "awskms",
						Disabled: false,
						Name:     "enabled-awskms",
					},
					{
						Type:     "awskms",
						Disabled: true,
						Name:     "disabled-awskms",
					},
				},
			}},
			expectError: false,
		},
		{
			name: "seals-with-same-names",
			config: &Config{SharedConfig: &configutil.SharedConfig{
				Seals: []*configutil.KMS{
					{
						Type:     "awskms",
						Disabled: false,
						Name:     "awskms",
					},
					{
						Type:     "awskms",
						Disabled: true,
						Name:     "awskms",
					},
				},
			}},
			expectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CheckConfig(tt.config, nil)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
