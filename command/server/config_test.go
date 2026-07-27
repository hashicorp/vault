// Copyright IBM Corp. 2016, 2025
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

// TestParseStorageURLConformance tests that all config attrs whose values can be
// URLs, IP addresses, or host:port addresses, when configured with an IPv6
// address, the normalized to be conformant with RFC-5942 §4
// See: https://rfc-editor.org/rfc/rfc5952.html
func TestParseStorageURLConformance(t *testing.T) {
	testParseStorageURLConformance(t)
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

// Test_ReportingScanDirectory makes sure that the reporting scan directory is correctly parsed
func Test_ReportingScanDirectory(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/reporting_directory.hcl")
	require.NoError(t, err)
	require.NotNil(t, config)
	require.NotEmpty(t, config.ReportingScanDirectory)
	require.Equal(t, "/foo/bar/", config.ReportingScanDirectory)
}

// Test_ObservationSystemConfig makes sure that the observation system config
// is properly loaded.
func Test_ObservationSystemConfig(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/observations.hcl")
	require.NoError(t, err)
	require.NotNil(t, config)
	require.NotNil(t, config.Observations)
	require.Equal(t, "/var/ledger.log", config.Observations.LedgerPath)
	require.Empty(t, config.Observations.TypePrefixAllowlist)
	require.Empty(t, config.Observations.TypePrefixDenylist)
}

// Test_ObservationSystemConfigAllowDenyList makes sure that the observation system config
// is properly loaded with an allowlist and denylist.
func Test_ObservationSystemConfigAllowDenyList(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/observations_allow_deny.hcl")
	require.NoError(t, err)
	require.NotNil(t, config)
	require.NotNil(t, config.Observations)
	require.Equal(t, "/var/ledger.log", config.Observations.LedgerPath)
	require.Equal(t, []string{"deny1", "deny2"}, config.Observations.TypePrefixDenylist)
	require.Equal(t, []string{"allow1", "allow2", "allow3"}, config.Observations.TypePrefixAllowlist)
	require.Equal(t, "0777", config.Observations.FileMode)
}

// Test_ObservationSystemConfigMerge checks merge for observation system config
func Test_ObservationSystemConfigMerge(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/observations.hcl")
	require.NoError(t, err)
	require.NotNil(t, config)

	config2, err := LoadConfigFile("./test-fixtures/observations_allow_deny.hcl")
	require.NoError(t, err)
	require.NotNil(t, config2)

	merged := config.Merge(config2)
	require.NotNil(t, merged)
	require.NotNil(t, merged.Observations)
	require.Equal(t, "/var/ledger.log", merged.Observations.LedgerPath)
	require.Equal(t, []string{"deny1", "deny2"}, merged.Observations.TypePrefixDenylist)
	require.Equal(t, []string{"allow1", "allow2", "allow3"}, merged.Observations.TypePrefixAllowlist)
	require.Equal(t, "0777", merged.Observations.FileMode)
}

// Test_ObservationSystemConfigMergeFromNoObservations checks merge for observation system config from a config
// without an observation system defined
func Test_ObservationSystemConfigMergeFromNoObservations(t *testing.T) {
	config, err := LoadConfigFile("./test-fixtures/config.hcl")
	require.NoError(t, err)
	require.NotNil(t, config)

	config2, err := LoadConfigFile("./test-fixtures/observations_allow_deny.hcl")
	require.NoError(t, err)
	require.NotNil(t, config2)

	merged := config.Merge(config2)
	require.NotNil(t, merged)
	require.NotNil(t, merged.Observations)
	require.Equal(t, "/var/ledger.log", merged.Observations.LedgerPath)
	require.Equal(t, []string{"deny1", "deny2"}, merged.Observations.TypePrefixDenylist)
	require.Equal(t, []string{"allow1", "allow2", "allow3"}, merged.Observations.TypePrefixAllowlist)
	require.Equal(t, "0777", merged.Observations.FileMode)
	require.Equal(t, true, merged.EnableUI)
}

// Test_UISettingsParsing checks that the ui_settings stanza is parsed correctly,
// normalizing ui_telemetry into UITelemetry/UITelemetrySet. An absent stanza leaves
// UISettings nil, and an empty stanza leaves telemetry unset so that a later config
// file merging in an explicit value is not treated as an override of a real choice.
func Test_UISettingsParsing(t *testing.T) {
	tests := []struct {
		name        string
		hcl         string
		wantNil     bool
		wantEnabled bool
		wantSet     bool
	}{
		{
			name:    "no stanza",
			hcl:     `disable_cache = true`,
			wantNil: true,
		},
		{name: "explicit true", hcl: "ui_settings { ui_telemetry = true }", wantEnabled: true, wantSet: true},
		{name: "explicit false", hcl: "ui_settings { ui_telemetry = false }", wantEnabled: false, wantSet: true},
		{name: "empty stanza leaves telemetry unset", hcl: "ui_settings { }", wantEnabled: false, wantSet: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config, err := ParseConfig(tc.hcl, "")
			require.NoError(t, err)
			if tc.wantNil {
				require.Nil(t, config.UISettings)
				return
			}
			require.NotNil(t, config.UISettings)
			require.Equal(t, tc.wantEnabled, config.UISettings.UITelemetry)
			require.Equal(t, tc.wantSet, config.UISettings.UITelemetrySet)
		})
	}
}

// Test_UISettingsValidate checks that unknown keys inside the ui_settings stanza are
// surfaced as warnings by Config.Validate, so a misspelled key is reported rather than
// silently ignored.
func Test_UISettingsValidate(t *testing.T) {
	config, err := ParseConfig("ui_settings { ui_telemetry = true }", "")
	require.NoError(t, err)
	require.Empty(t, config.Validate(""), "known keys should not produce warnings")

	config, err = ParseConfig("ui_settings { ui_telemeeeetry = true }", "")
	require.NoError(t, err)
	results := config.Validate("")
	require.NotEmpty(t, results, "a misspelled/unknown key should produce a warning")
	require.Contains(t, results[0].String(), "ui_telemeeeetry")
}

// Test_UISettingsSanitized checks that the ui_settings stanza is surfaced by
// Config.Sanitized (and therefore sys/config/state/sanitized) only when it is actually
// configured, so operators can confirm whether the stanza was loaded.
func Test_UISettingsSanitized(t *testing.T) {
	configured, err := ParseConfig(`ui_settings { ui_telemetry = true }`, "")
	require.NoError(t, err)
	sanitized := configured.Sanitized()
	require.Contains(t, sanitized, "ui_settings")
	require.Equal(t, map[string]interface{}{"ui_telemetry": true}, sanitized["ui_settings"])

	absent, err := ParseConfig(`disable_cache = true`, "")
	require.NoError(t, err)
	require.NotContains(t, absent.Sanitized(), "ui_settings")
}

// Test_UISettingsMerge checks that merging two configs respects UITelemetrySet, so an
// explicitly set value overrides an earlier one while an unset value never clobbers a
// previously configured one.
func Test_UISettingsMerge(t *testing.T) {
	mustParse := func(hcl string) *Config {
		config, err := ParseConfig(hcl, "")
		require.NoError(t, err)
		return config
	}

	enabledTrue := `ui_settings { ui_telemetry = true }`
	enabledFalse := `ui_settings { ui_telemetry = false }`
	empty := `ui_settings {}`
	noStanza := `disable_cache = true`

	t.Run("second explicit value overrides first", func(t *testing.T) {
		merged := mustParse(enabledTrue).Merge(mustParse(enabledFalse))
		require.NotNil(t, merged.UISettings)
		require.False(t, merged.UISettings.UITelemetry)
		require.True(t, merged.UISettings.UITelemetrySet)
	})

	t.Run("unset second does not clobber first", func(t *testing.T) {
		merged := mustParse(enabledTrue).Merge(mustParse(empty))
		require.NotNil(t, merged.UISettings)
		require.True(t, merged.UISettings.UITelemetry)
		require.True(t, merged.UISettings.UITelemetrySet)
	})

	t.Run("nil first takes second", func(t *testing.T) {
		merged := mustParse(noStanza).Merge(mustParse(enabledTrue))
		require.NotNil(t, merged.UISettings)
		require.True(t, merged.UISettings.UITelemetry)
		require.True(t, merged.UISettings.UITelemetrySet)
	})

	t.Run("nil second keeps first", func(t *testing.T) {
		merged := mustParse(enabledTrue).Merge(mustParse(noStanza))
		require.NotNil(t, merged.UISettings)
		require.True(t, merged.UISettings.UITelemetry)
	})
}

// TestDuplicateKeyValidationHcl checks that loading an HCL config file with duplicate keys returns an error.
func TestDuplicateKeyValidationHcl(t *testing.T) {
	testDuplicateKeyValidationHcl(t)
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
			_, err := CheckConfig(tt.config)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
