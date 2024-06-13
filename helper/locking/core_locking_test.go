// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package locking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseDetectDeadlockConfigParameter verifies that all types of strings
// that could be obtained from the configuration file, are correctly parsed
// into a slice of string elements.
func TestParseDetectDeadlockConfigParameter(t *testing.T) {
	for _, tc := range []struct {
		name                          string
		detectDeadlockConfigParameter string
		expectedResult                []string
	}{
		{
			name: "empty-string",
		},
		{
			name:                          "single-value",
			detectDeadlockConfigParameter: "bar",
			expectedResult:                []string{"bar"},
		},
		{
			name:                          "single-value-mixed-case",
			detectDeadlockConfigParameter: "BaR",
			expectedResult:                []string{"bar"},
		},
		{
			name:                          "multiple-values",
			detectDeadlockConfigParameter: "bar,BAZ,fIZ",
			expectedResult:                []string{"bar", "baz", "fiz"},
		},
		{
			name:                          "non-canonical-string-list",
			detectDeadlockConfigParameter: "bar  ,  baz, ",
			expectedResult:                []string{"bar", "baz", ""},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseDetectDeadlockConfigParameter(tc.detectDeadlockConfigParameter)
			assert.ElementsMatch(t, tc.expectedResult, result)
		})
	}
}

// TestCreateConfigurableRWMutex verifies the correct behaviour in determining
// whether a deadlock detecting RWMutex should be returned or not based on the
// input arguments for the CreateConfigurableRWMutex function.
func TestCreateConfigurableRWMutex(t *testing.T) {
	mutexTypes := map[bool]string{
		false: "locking.SyncRWMutex",
		true:  "locking.DeadlockRWMutex",
	}

	for _, tc := range []struct {
		name               string
		detectDeadlocks    []string
		lock               string
		expectDeadlockLock bool
	}{
		{
			name: "no-lock-types-specified",
			lock: "foo",
		},
		{
			name:            "single-lock-specified-no-match",
			detectDeadlocks: []string{"bar"},
			lock:            "foo",
		},
		{
			name:               "single-lock-specified-match",
			detectDeadlocks:    []string{"foo"},
			lock:               "foo",
			expectDeadlockLock: true,
		},
		{
			name:            "multiple-locks-specified-no-match",
			detectDeadlocks: []string{"bar", "baz", "fiz"},
			lock:            "foo",
		},
		{
			name:               "multiple-locks-specified-match",
			detectDeadlocks:    []string{"bar", "foo", "baz"},
			lock:               "foo",
			expectDeadlockLock: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := CreateConfigurableRWMutex(tc.detectDeadlocks, tc.lock)

			_, ok := m.(*DeadlockRWMutex)
			if tc.expectDeadlockLock != ok {
				t.Fatalf("unexpected RWMutex type returned, expected: %s got %s", mutexTypes[tc.expectDeadlockLock], mutexTypes[ok])
			}
		})
	}
}
