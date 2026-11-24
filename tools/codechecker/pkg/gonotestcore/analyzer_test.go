// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package gonotestcore

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer runs the analyzer on the functions in testdata/funcs.go. The reports from the analyzer are compared
// against the comments in funcs.go beginning with "want". If there is no comment beginning with "want" on a line, then
// the analyzer is expected not to report anything there.
func TestAnalyzer(t *testing.T) {
	f, err := os.Getwd()
	if err != nil {
		t.Fatal("failed to get working directory", err)
	}
	analysistest.Run(t, filepath.Join(f, "testdata"), Analyzer, ".")
}
