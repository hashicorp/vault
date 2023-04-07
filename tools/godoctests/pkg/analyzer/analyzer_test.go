// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer runs the analyzer on the test functions in testdata/funcs.go. The report from the analyzer is compared against
// the comments in funcs.go beginning with "want." If there is no comment beginning with "want", then the analyzer is expected
// not to report anything.
func TestAnalyzer(t *testing.T) {
	f, err := os.Getwd()
	if err != nil {
		t.Fatal("failed to get working directory", err)
	}
	analysistest.Run(t, filepath.Join(f, "testdata"), Analyzer, ".")
}
