package connutil

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/strutil"
)

func TestBuildConnectionString(t *testing.T) {
	var tables = []struct {
		format          Format
		conf            map[string]string
		expectedConnStr string
		expectedOK      bool
	}{
		{
			NOTIMPLEMENTED,
			map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
			"",
			false,
		},
		{
			ADO,
			map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
			"key1=value1;key2=value2;key3=value3",
			true,
		},
	}

	for _, table := range tables {
		actualConnStr, actualOK := BuildConnectionString(table.conf, table.format)
		// BuildConnectionString doesn't guarantee ordering for ADO style strings, "key1=value1;key2=value2" and
		// "key2=value2;key1=value1" should be considered equivalent
		if !strutil.EquivalentSlices(
			strings.Split(table.expectedConnStr, ";"),
			strings.Split(actualConnStr, ";")) {
			t.Errorf("Expected connection string to be:\n\"%s\" but got:\n\"%s\"", table.expectedConnStr, actualConnStr)
		}
		if table.expectedOK != actualOK {
			t.Errorf("Expected OK flag to be:\n%t but got:\n%t", table.expectedOK, actualOK)
		}
	}
}

func TestInjectDefaultsIntoConnectionString(t *testing.T) {
	var tables = []struct {
		conf         map[string]string
		defaults     map[string]string
		expectedConf map[string]string
	}{
		{
			map[string]string{
				"keySpecifiedInConfOnly":        "confValue",
				"keySpecifiedInConfAndDefaults": "confValue",
			},
			map[string]string{
				"keySpecifiedInDefaultsOnly":    "defaultsValue",
				"keySpecifiedInConfAndDefaults": "defaultsValue",
			},
			map[string]string{
				"keySpecifiedInConfOnly":        "confValue",
				"keySpecifiedInConfAndDefaults": "confValue",
				"keySpecifiedInDefaultsOnly":    "defaultsValue",
			},
		},
	}

	for _, table := range tables {
		actualConf := InjectDefaultsIntoConnectionString(table.defaults, table.conf)
		if !strutil.EqualStringMaps(table.expectedConf, actualConf) {
			t.Errorf("Expected conf to be:\n%v but got:\n%v", table.expectedConf, actualConf)
		}
	}
}

func TestUpgradeParameterKeysInConnectionString(t *testing.T) {
	var tables = []struct {
		conf         map[string]string
		backcomp     map[string]string
		expectedConf map[string]string
	}{
		{
			map[string]string{
				"key":                "value",
				"legacyKey":          "value",
				"duplicateLegacyKey": "legacyValue",
				"duplicateNewKey":    "newValue",
			},
			map[string]string{
				"legacyKey":          "newKey",
				"unusedLegacyKey":    "unusedNewKey",
				"duplicateLegacyKey": "duplicateNewKey",
			},
			map[string]string{
				"key":             "value",
				"newKey":          "value",
				"duplicateNewKey": "legacyValue",
			},
		},
	}

	for _, table := range tables {
		actualConf := UpgradeParameterKeysInConnectionString(table.backcomp, table.conf)
		if !strutil.EqualStringMaps(table.expectedConf, actualConf) {
			t.Errorf("Expected conf to be:\n%v but got:\n%v", table.expectedConf, actualConf)
		}
	}
}
