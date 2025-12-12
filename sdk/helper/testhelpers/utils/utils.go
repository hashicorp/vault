// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package utils

import (
	"os"
	"strings"
	"testing"
)

// SkipUnlessEnvVarsSet skips the test unless all of the given environment
// variables are set
func SkipUnlessEnvVarsSet(t testing.TB, envVars []string) {
	t.Helper()

	for _, i := range envVars {
		if os.Getenv(i) == "" {
			t.Skipf("%s must be set for this test to run", strings.Join(envVars, " "))
		}
	}
}
