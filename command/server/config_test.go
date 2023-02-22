package server

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
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
