package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	f, err := os.Getwd()
	if err != nil {
		t.Fatal("failed to get working directory", err)
	}
	analysistest.Run(t, filepath.Join(f, "testdata"), Analyzer, ".")
}
