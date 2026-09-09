// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package ebom

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

// newSingleSheetXLSX creates a minimal single-sheet EBOM workbook in t.TempDir()
// and returns its path.  The eAssembly description has no tier/variant tokens so
// inferSlug returns ("", false), matching 2.x single-sheet behaviour.
func newSingleSheetXLSX(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "single_sheet.xlsx")
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	rows := [][]string{
		{"eAssemblies", "", "", "", ""},
		{"PID", "5900-TST"},
		{"eAssembly", "IBM Vault Self-Managed V2.0.0 Multiplatform English", "", "eAssembly", ""},
		{"Part No", "Platform Code", "Platform Description", "eGA Date", "Qty"},
		{"", "TST01     ", "All Platforms", "2024-01-01", "1"},
		{"", "TST02     ", "Linux x86_64", "2024-01-01", "1"},
		{"", "", "", "", ""},
		{"", "TST03     ", "Linux ARM64", "2024-01-01", "1"},
	}
	for i, row := range rows {
		require.NoError(t, f.SetSheetRow(sheet, fmt.Sprintf("A%d", i+1), &row))
	}
	require.NoError(t, f.SaveAs(path))
	return path
}

// newMultiSheetXLSX creates a three-eAssembly-sheet workbook in t.TempDir()
// and returns its path.
func newMultiSheetXLSX(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "multi_sheet.xlsx")
	f := excelize.NewFile()
	defer f.Close()

	type sheetDef struct {
		tab  string
		desc string
	}
	sheets := []sheetDef{
		{"Vault Essen", "IBM Vault Self-Managed Essentials V1.20.4 Multiplatform English"},
		{"Vault Std", "IBM Vault Self-Managed Standard V1.20.4 Multiplatform English"},
		{"Vault Prem", "IBM Vault Self-Managed Premium V1.20.4 Multiplatform English"},
	}

	f.SetSheetName("Sheet1", sheets[0].tab)
	for _, sd := range sheets[1:] {
		_, newErr := f.NewSheet(sd.tab)
		require.NoError(t, newErr)
	}

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
	return path
}

// newBlankRowsXLSX creates a single-sheet workbook with interspersed blank rows
// in t.TempDir() and returns its path.
func newBlankRowsXLSX(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "blank_rows.xlsx")
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	rows := [][]string{
		{"eAssemblies", "", "", "", ""},
		{"PID", "5900-TST"},
		{"eAssembly", "IBM Vault Self-Managed V2.0.0 Multiplatform English", "", "eAssembly", ""},
		{"Part No", "Platform Code", "Platform Description", "eGA Date", "Qty"},
		{"", "", "", "", ""},
		{"", "TST01", "All Platforms", "2024-01-01", "1"},
		{"", "", "", "", ""},
		{"", "TST02", "Linux x86_64", "2024-01-01", "1"},
	}
	for i, row := range rows {
		require.NoError(t, f.SetSheetRow(sheet, fmt.Sprintf("A%d", i+1), &row))
	}
	require.NoError(t, f.SaveAs(path))
	return path
}

// newMixedSheetsXLSX creates a workbook with one eAssembly sheet and one
// non-eAssembly sheet in t.TempDir() and returns its path.
func newMixedSheetsXLSX(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "mixed_sheets.xlsx")
	f := excelize.NewFile()
	defer f.Close()

	f.SetSheetName("Sheet1", "Vault Essen")
	_, newErr := f.NewSheet("pBom - NA")
	require.NoError(t, newErr)

	eaRows := [][]string{
		{"eAssemblies", "", "", "", ""},
		{"PID", "5900-TST"},
		{"eAssembly", "IBM Vault Self-Managed Essentials V1.20.4 Multiplatform English", "", "eAssembly", ""},
		{"Part No", "Platform Code", "Platform Description", "eGA Date", "Qty"},
		{"", "TST01", "All Platforms", "2024-01-01", "1"},
	}
	for i, row := range eaRows {
		require.NoError(t, f.SetSheetRow("Vault Essen", fmt.Sprintf("A%d", i+1), &row))
	}

	pbomRows := [][]string{
		{"pBom", "something else"},
	}
	for i, row := range pbomRows {
		require.NoError(t, f.SetSheetRow("pBom - NA", fmt.Sprintf("A%d", i+1), &row))
	}

	require.NoError(t, f.SaveAs(path))
	return path
}

// newMultiSheetUnknownSlugXLSX creates a two-eAssembly-sheet workbook where
// the second sheet's description contains no slug tokens, so inferSlug returns
// ("", false) for it.  Used to exercise the missing-slug footer path.
func newMultiSheetUnknownSlugXLSX(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "multi_sheet_unknown.xlsx")
	f := excelize.NewFile()
	defer f.Close()

	type sheetDef struct {
		tab  string
		desc string
	}
	sheets := []sheetDef{
		{"Vault Essen", "IBM Vault Self-Managed Essentials V1.20.4 Multiplatform English"},
		{"V$version", "IBM Vault Self-Managed V2.1.0 Multiplatform English"},
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
	return path
}

// newNoEAssemblySheetsXLSX creates a workbook with no eAssembly sheets in
// t.TempDir() and returns its path.
func newNoEAssemblySheetsXLSX(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "no_ea.xlsx")
	f := excelize.NewFile()
	defer f.Close()

	f.SetSheetName("Sheet1", "pBom - NA")
	rows := [][]string{
		{"pBom", "something else"},
	}
	for i, row := range rows {
		require.NoError(t, f.SetSheetRow("pBom - NA", fmt.Sprintf("A%d", i+1), &row))
	}
	require.NoError(t, f.SaveAs(path))
	return path
}
