// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/ebom"
	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

// createConvertTestSingleXLSX creates a single-sheet EBOM fixture for convert command tests.
func createConvertTestSingleXLSX(t *testing.T, path string) {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	rows := [][]string{
		{"eAssemblies", "", "", "", ""},
		{"PID", "5900-TST"},
		{"eAssembly", "IBM Vault Self-Managed V2.0.0 Multiplatform English", "", "eAssembly", ""},
		{"Part No", "Platform Code", "Platform Description", "eGA Date", "Qty"},
		{"", "TST01", "All Platforms", "2024-01-01", "1"},
	}
	for i, row := range rows {
		require.NoError(t, f.SetSheetRow(sheet, fmt.Sprintf("A%d", i+1), &row))
	}
	require.NoError(t, f.SaveAs(path))
}

// TestEbomConvertXLSXCmd_XLSError verifies that passing a .xls file returns
// the resave error before any file I/O.
func TestEbomConvertXLSXCmd_XLSError(t *testing.T) {
	oldStdout, oldStderr := os.Stdout, os.Stderr
	_, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout, os.Stderr = w, w
	defer func() {
		w.Close()
		os.Stdout, os.Stderr = oldStdout, oldStderr
	}()

	convertEBOMToCSVReq = &ebom.ConvertEBOMToCSVReq{}
	cmd := newEBOMConvertCmd()
	cmd.SetArgs([]string{"foo.xls"})
	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "not .xlsx")
}

// TestEbomConvertXLSXCmd_SlugsDotRejected verifies that a --slugs value
// containing a dot is rejected before any file is written.
func TestEbomConvertXLSXCmd_SlugsDotRejected(t *testing.T) {
	oldStdout, oldStderr := os.Stdout, os.Stderr
	_, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout, os.Stderr = w, w
	defer func() {
		w.Close()
		os.Stdout, os.Stderr = oldStdout, oldStderr
	}()

	convertEBOMToCSVReq = &ebom.ConvertEBOMToCSVReq{}
	cmd := newEBOMConvertCmd()
	cmd.SetArgs([]string{"dummy.xlsx", "--slugs", "0:essentials.openshift"})
	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "dots are not permitted")
}

// TestEbomConvertXLSXCmd_FormatJSON verifies that --format json produces valid
// JSON on stdout.
func TestEbomConvertXLSXCmd_FormatJSON(t *testing.T) {
	tmpDir := t.TempDir()
	xlsxPath := filepath.Join(tmpDir, "test.xlsx")
	createConvertTestSingleXLSX(t, xlsxPath)
	outDir := t.TempDir()

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	convertEBOMToCSVReq = &ebom.ConvertEBOMToCSVReq{}
	rootCfg.format = "json"
	defer func() { rootCfg.format = "table" }()

	cmd := newEBOMConvertCmd()
	cmd.SetArgs([]string{xlsxPath, "--out-dir", outDir})
	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &parsed))
	_, hasSheets := parsed["sheets"]
	require.True(t, hasSheets)
}

// TestEbomConvertXLSXCmd_FormatMarkdown verifies that --format markdown
// produces Markdown table markers on stdout.
func TestEbomConvertXLSXCmd_FormatMarkdown(t *testing.T) {
	tmpDir := t.TempDir()
	xlsxPath := filepath.Join(tmpDir, "test.xlsx")
	createConvertTestSingleXLSX(t, xlsxPath)
	outDir := t.TempDir()

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	convertEBOMToCSVReq = &ebom.ConvertEBOMToCSVReq{}
	rootCfg.format = "markdown"
	defer func() { rootCfg.format = "table" }()

	cmd := newEBOMConvertCmd()
	cmd.SetArgs([]string{xlsxPath, "--out-dir", outDir})
	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	require.Contains(t, buf.String(), "|")
}

// TestEbomConvertXLSXCmd_OutDir verifies that --out-dir places the CSV in the
// specified directory.
func TestEbomConvertXLSXCmd_OutDir(t *testing.T) {
	tmpDir := t.TempDir()
	xlsxPath := filepath.Join(tmpDir, "test.xlsx")
	createConvertTestSingleXLSX(t, xlsxPath)
	outDir := t.TempDir()

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	convertEBOMToCSVReq = &ebom.ConvertEBOMToCSVReq{}
	rootCfg.format = "table"
	defer func() { rootCfg.format = "table" }()

	cmd := newEBOMConvertCmd()
	cmd.SetArgs([]string{xlsxPath, "--out-dir", outDir})
	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	entries, readErr := os.ReadDir(outDir)
	require.NoError(t, readErr)
	require.NotEmpty(t, entries)
	require.Equal(t, "5900-TST.csv", entries[0].Name())
}
