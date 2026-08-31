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

// createInspectTestXLSX creates a minimal multi-sheet EBOM fixture for command-layer tests.
func createInspectTestXLSX(t *testing.T, path string) {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()

	type sheetDef struct {
		tab  string
		desc string
	}
	sheets := []sheetDef{
		{"Vault Essen", "IBM Vault Self-Managed Essentials V1.20.4 Multiplatform English"},
		{"Vault Prem", "IBM Vault Self-Managed Premium V1.20.4 Multiplatform English"},
	}

	f.SetSheetName("Sheet1", sheets[0].tab)
	_, newErr := f.NewSheet(sheets[1].tab)
	require.NoError(t, newErr)

	for _, sd := range sheets {
		rows := [][]string{
			{"eAssemblies", "", "", "", ""},
			{"PID", "5900-TST"},
			{"eAssembly", sd.desc, "", "eAssembly", ""},
			{"Part No", "Platform Code", "Platform Description", "eGA Date", "Qty"},
			{"", "TST01", "All Platforms", "2024-01-01", "1"},
		}
		for i, row := range rows {
			require.NoError(t, f.SetSheetRow(sd.tab, fmt.Sprintf("A%d", i+1), &row))
		}
	}

	require.NoError(t, f.SaveAs(path))
}

// TestEbomInspectXLSXCmd_XLSError verifies that passing a .xls file returns
// the resave error message before any file I/O.
func TestEbomInspectXLSXCmd_XLSError(t *testing.T) {
	oldStdout, oldStderr := os.Stdout, os.Stderr
	_, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout, os.Stderr = w, w
	defer func() {
		w.Close()
		os.Stdout, os.Stderr = oldStdout, oldStderr
	}()

	inspectEBOMReq = &ebom.InspectEBOMReq{}
	cmd := newEBOMInspectCmd()
	cmd.SetArgs([]string{"foo.xls"})
	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "not .xlsx")
}

// TestEbomInspectXLSXCmd_FormatJSON verifies that --format json produces valid
// JSON output containing the index and inferredSlug fields.
func TestEbomInspectXLSXCmd_FormatJSON(t *testing.T) {
	tmpDir := t.TempDir()
	xlsxPath := filepath.Join(tmpDir, "test.xlsx")
	createInspectTestXLSX(t, xlsxPath)

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	inspectEBOMReq = &ebom.InspectEBOMReq{}
	rootCfg.format = "json"
	defer func() { rootCfg.format = "table" }()

	cmd := newEBOMInspectCmd()
	cmd.SetArgs([]string{xlsxPath})
	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &parsed))
	sheets, ok := parsed["sheets"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, sheets)

	first := sheets[0].(map[string]any)
	_, hasIndex := first["index"]
	_, hasSlug := first["inferredSlug"]
	require.True(t, hasIndex)
	require.True(t, hasSlug)
}

// TestEbomInspectXLSXCmd_FormatTable verifies that default (table) output
// contains the INDEX and INFERRED SLUG columns, and that when all slugs are
// inferred the SLUG OVERRIDE column is absent.
func TestEbomInspectXLSXCmd_FormatTable(t *testing.T) {
	tmpDir := t.TempDir()
	xlsxPath := filepath.Join(tmpDir, "test.xlsx")
	createInspectTestXLSX(t, xlsxPath)

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	inspectEBOMReq = &ebom.InspectEBOMReq{}
	rootCfg.format = "table"
	defer func() { rootCfg.format = "table" }()

	cmd := newEBOMInspectCmd()
	cmd.SetArgs([]string{xlsxPath})
	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	require.Contains(t, buf.String(), "INDEX")
	require.Contains(t, buf.String(), "INFERRED SLUG")
	require.NotContains(t, buf.String(), "SLUG OVERRIDE")
}
