package server

import (
	"fmt"
	"testing"

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
