// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

func Test_ParseLogFormat(t *testing.T) {
	type testData struct {
		format      string
		expected    LogFormat
		expectedErr error
	}

	tests := []testData{
		{format: "", expected: UnspecifiedFormat, expectedErr: nil},
		{format: " ", expected: UnspecifiedFormat, expectedErr: nil},
		{format: "standard", expected: StandardFormat, expectedErr: nil},
		{format: "STANDARD", expected: StandardFormat, expectedErr: nil},
		{format: "json", expected: JSONFormat, expectedErr: nil},
		{format: " json ", expected: JSONFormat, expectedErr: nil},
		{format: "bogus", expected: UnspecifiedFormat, expectedErr: errors.New("unknown log format: bogus")},
	}

	for _, test := range tests {
		result, err := ParseLogFormat(test.format)
		if test.expected != result {
			t.Errorf("expected %s, got %s", test.expected, result)
		}
		if !reflect.DeepEqual(test.expectedErr, err) {
			t.Errorf("expected error %v, got %v", test.expectedErr, err)
		}
	}
}

func Test_ParseEnv_VAULT_LOG_FORMAT(t *testing.T) {
	oldVLF := os.Getenv("VAULT_LOG_FORMAT")
	defer os.Setenv("VAULT_LOG_FORMAT", oldVLF)

	testParseEnvLogFormat(t, "VAULT_LOG_FORMAT")
}

func testParseEnvLogFormat(t *testing.T, name string) {
	env := []string{
		"json", "vauLT_Json", "VAULT-JSON", "vaulTJSon",
		"standard", "STANDARD",
		"bogus",
	}

	formats := []LogFormat{
		JSONFormat, JSONFormat, JSONFormat, JSONFormat,
		StandardFormat, StandardFormat,
		UnspecifiedFormat,
	}

	for i, e := range env {
		os.Setenv(name, e)
		if lf := ParseEnvLogFormat(); formats[i] != lf {
			t.Errorf("expected %s, got %s", formats[i], lf)
		}
	}
}
