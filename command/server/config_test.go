// +build !enterprise

package server

import (
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	testLoadConfigFile(t)
}

func TestLoadConfigFile_topLevel(t *testing.T) {
	testLoadConfigFile_topLevel(t, nil)
}

func TestLoadConfigFile_json(t *testing.T) {
	testLoadConfigFile_json(t)
}

func TestLoadConfigFile_json2(t *testing.T) {
	testLoadConfigFile_json2(t, nil)
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

func TestParseEntropy(t *testing.T) {
	testParseEntropy(t, true)
}

func TestConfigRaftRetryJoin(t *testing.T) {
	testConfigRaftRetryJoin(t)
}
