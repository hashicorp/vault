// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package ebom

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xuri/excelize/v2"
)

// TestEbomSheet_ProcessRows verifies that processRows trims cells, strips only
// leading and trailing blank rows, and preserves internal blank rows.
func TestEbomSheet_ProcessRows(t *testing.T) {
	t.Parallel()
	raw := [][]string{
		{"", "", ""}, // leading blank — stripped
		{"eAssemblies", "", ""},
		{"PID", "5900-TST"},
		{"", "", ""}, // internal blank — preserved
		{"TST01     ", "All Platforms", ""},
		{"TST02"},
		{"", "", ""}, // trailing blank — stripped
	}
	s := &ebomSheet{name: "Sheet1"}
	s.processRows(raw)

	require.Len(t, s.rows, 5)
	require.Equal(t, "eAssemblies", s.rows[0][0])
	require.Equal(t, []string{"", "", ""}, s.rows[2]) // internal blank preserved
	require.Equal(t, "TST01", s.rows[3][0])
}

// TestEbomSheet_ProcessRows_ExcelErrors verifies that Excel error strings such
// as #REF! are cleared to empty strings rather than written literally to the
// output CSV.
func TestEbomSheet_ProcessRows_ExcelErrors(t *testing.T) {
	t.Parallel()
	raw := [][]string{
		{"eAssemblies", "", ""},
		{"PID", "5900-TST"},
		{"", "Media Pack P/N", "#REF!", "#REF!", "#REF!"},
		{"", "#NAME?", "#VALUE!", "#DIV/0!", "#NULL!"},
		{"", "normal", "#N/A", "#NUM!", ""},
	}
	s := &ebomSheet{name: "Sheet1"}
	s.processRows(raw)

	require.Len(t, s.rows, 5)
	// All error strings in row 2 must be cleared.
	require.Equal(t, []string{"", "Media Pack P/N", "", "", ""}, s.rows[2])
	// Row 3: #NAME? and remaining errors cleared.
	require.Equal(t, []string{"", "", "", "", ""}, s.rows[3])
	// Row 4: #N/A and #NUM! cleared; normal value and empty preserved.
	require.Equal(t, []string{"", "normal", "", "", ""}, s.rows[4])
}

// TestConvertEBOMToCSV_SingleSheet verifies that a single-sheet workbook
// produces <pid>.csv with no slug in the filename.
func TestConvertEBOMToCSV_SingleSheet(t *testing.T) {
	t.Parallel()
	path := newSingleSheetXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath: path,
		OutDir:   outDir,
	}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	var written []*EBOMCSV
	for _, s := range res.Sheets {
		if s.AppendedTo == "" {
			written = append(written, s)
		}
	}
	require.Len(t, written, 1)
	require.Equal(t, "5900-TST.csv", filepath.Base(written[0].OutputPath))
	require.Empty(t, written[0].Slug)
	require.FileExists(t, written[0].OutputPath)
}

// TestConvertEBOMToCSV_MultiSheet verifies that a multi-sheet workbook
// produces one <pid>.<slug>.csv per eAssembly sheet with slugs inferred.
func TestConvertEBOMToCSV_MultiSheet(t *testing.T) {
	t.Parallel()
	path := newMultiSheetXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath: path,
		OutDir:   outDir,
	}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	var written []*EBOMCSV
	for _, s := range res.Sheets {
		if s.AppendedTo == "" {
			written = append(written, s)
		}
	}
	require.Len(t, written, 3)
	require.Equal(t, "essentials", written[0].Slug)
	require.Equal(t, "standard", written[1].Slug)
	require.Equal(t, "premium", written[2].Slug)

	for _, s := range written {
		require.FileExists(t, s.OutputPath)
		require.Equal(t, "5900-TST."+s.Slug+".csv", filepath.Base(s.OutputPath))
	}
}

// TestConvertEBOMToCSV_MultiSheetSlugOverride verifies that a --slugs override
// for index 0 is used while other sheets fall back to inference.
func TestConvertEBOMToCSV_MultiSheetSlugOverride(t *testing.T) {
	t.Parallel()
	path := newMultiSheetXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath: path,
		OutDir:   outDir,
		Slugs:    []string{"0:essen"},
	}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	var written []*EBOMCSV
	for _, s := range res.Sheets {
		if s.AppendedTo == "" {
			written = append(written, s)
		}
	}
	require.Len(t, written, 3)
	require.Equal(t, "essen", written[0].Slug)
	require.Equal(t, "standard", written[1].Slug)
	require.Equal(t, "premium", written[2].Slug)
}

// TestConvertEBOMToCSV_MultiSheetAllOverrides verifies that when all sheets
// have explicit overrides, inference is never called.
func TestConvertEBOMToCSV_MultiSheetAllOverrides(t *testing.T) {
	t.Parallel()
	path := newMultiSheetXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath: path,
		OutDir:   outDir,
		Slugs:    []string{"0:a", "1:b", "2:c"},
	}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	var written []*EBOMCSV
	for _, s := range res.Sheets {
		if s.AppendedTo == "" {
			written = append(written, s)
		}
	}
	require.Len(t, written, 3)
	require.Equal(t, "a", written[0].Slug)
	require.Equal(t, "b", written[1].Slug)
	require.Equal(t, "c", written[2].Slug)
}

// TestConvertEBOMToCSV_PIDOverride verifies that --pid replaces the workbook
// PID in the output filename.
func TestConvertEBOMToCSV_PIDOverride(t *testing.T) {
	t.Parallel()
	path := newSingleSheetXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath:    path,
		OutDir:      outDir,
		PIDOverride: "5900-OVR",
	}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	var written []*EBOMCSV
	for _, s := range res.Sheets {
		if s.AppendedTo == "" {
			written = append(written, s)
		}
	}
	require.Len(t, written, 1)
	require.Equal(t, "5900-OVR", written[0].PID)
	require.Equal(t, "5900-OVR.csv", filepath.Base(written[0].OutputPath))
}

// TestConvertEBOMToCSV_PIDFromWorkbook verifies that the PID is read from the
// "PID" row in the sheet when no override is given.
func TestConvertEBOMToCSV_PIDFromWorkbook(t *testing.T) {
	t.Parallel()
	path := newSingleSheetXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath: path,
		OutDir:   outDir,
	}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	for _, s := range res.Sheets {
		if s.AppendedTo == "" {
			require.Equal(t, "5900-TST", s.PID)
		}
	}
}

// TestConvertEBOMToCSV_NonEBOMSheetsAppended verifies that non-eAssembly sheets
// are appended to the output CSV rather than skipped, and that the response
// records them with AppendedTo pointing at the output file.
func TestConvertEBOMToCSV_NonEBOMSheetsAppended(t *testing.T) {
	t.Parallel()
	path := newMixedSheetsXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath: path,
		OutDir:   outDir,
	}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	var primary, appended int
	for _, s := range res.Sheets {
		if s.AppendedTo != "" {
			appended++
			require.NotEmpty(t, s.AppendedTo)
		} else {
			primary++
		}
	}
	require.Equal(t, 1, primary)
	require.Equal(t, 1, appended)
}

// TestConvertEBOMToCSV_SectionColumnWidths verifies that the appended
// non-eAssembly section is written with its own column width and is not padded
// to the wider eAssembly section's column count.
func TestConvertEBOMToCSV_SectionColumnWidths(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "width_mismatch.xlsx")
	f := excelize.NewFile()
	defer f.Close()

	f.SetSheetName("Sheet1", "EA Sheet")
	_, newErr := f.NewSheet("pBom")
	require.NoError(t, newErr)

	eaRows := [][]string{
		{"eAssemblies", "", "", "", "", "", "", "", "", ""},
		{"PID", "5900-TST", "", "", "", "", "", "", "", ""},
		{"", "IBM Vault Self-Managed V2.0.0 Multiplatform English", "", "eAssembly", "", "", "", "", "", ""},
		{"col1", "col2", "col3", "col4", "col5", "col6", "col7", "col8", "col9", "col10"},
		{"", "TST01", "", "", "", "", "", "", "", ""},
	}
	for i, row := range eaRows {
		require.NoError(t, f.SetSheetRow("EA Sheet", fmt.Sprintf("A%d", i+1), &row))
	}

	pbomRows := [][]string{
		{"Level", "Description", "NA English"},
		{"1", "Hard copy LI", "xxxx"},
		{"1", "Quick Start", "CF3DQZZ"},
	}
	for i, row := range pbomRows {
		require.NoError(t, f.SetSheetRow("pBom", fmt.Sprintf("A%d", i+1), &row))
	}

	require.NoError(t, f.SaveAs(path))

	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{XLSXPath: path, OutDir: outDir}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	var csvPath string
	for _, s := range res.Sheets {
		if s.AppendedTo == "" {
			csvPath = s.OutputPath
		}
	}
	require.NotEmpty(t, csvPath)

	data, readErr := os.ReadFile(csvPath)
	require.NoError(t, readErr)

	r := csv.NewReader(strings.NewReader(string(data)))
	r.FieldsPerRecord = -1
	records, parseErr := r.ReadAll()
	require.NoError(t, parseErr)

	// The fixture has 5 eAssembly rows (10 cols) followed immediately by 3
	// pBom rows (3 cols). Verify section boundaries by column count.
	require.Len(t, records, 8)
	for _, rec := range records[:5] {
		require.Len(t, rec, 10, "eAssembly rows should have 10 columns")
	}
	for _, rec := range records[5:] {
		require.Len(t, rec, 3, "pBom rows should have 3 columns, not padded to eAssembly width")
	}
}

// TestConvertEBOMToCSV_BlankRowHandling verifies that internal blank rows are
// preserved in the output CSV while leading and trailing blank rows are stripped.
func TestConvertEBOMToCSV_BlankRowHandling(t *testing.T) {
	t.Parallel()
	path := newBlankRowsXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath: path,
		OutDir:   outDir,
	}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	var csvPath string
	for _, s := range res.Sheets {
		if s.AppendedTo == "" {
			csvPath = s.OutputPath
		}
	}
	require.NotEmpty(t, csvPath)

	data, readErr := os.ReadFile(csvPath)
	require.NoError(t, readErr)

	r := csv.NewReader(strings.NewReader(string(data)))
	r.FieldsPerRecord = -1
	records, parseErr := r.ReadAll()
	require.NoError(t, parseErr)

	// The fixture has 2 internal blank rows between data rows — they must be
	// present in the output.
	blankCount := 0
	for _, rec := range records {
		allEmpty := true
		for _, cell := range rec {
			if cell != "" {
				allEmpty = false
				break
			}
		}
		if allEmpty {
			blankCount++
		}
	}
	require.Equal(t, 2, blankCount, "expected 2 internal blank rows in CSV output")

	// Leading and trailing rows must not be blank.
	require.NotEmpty(t, records)
	firstNonEmpty := func(row []string) bool {
		for _, c := range row {
			if c != "" {
				return true
			}
		}
		return false
	}
	require.True(t, firstNonEmpty(records[0]), "first row must not be blank")
	require.True(t, firstNonEmpty(records[len(records)-1]), "last row must not be blank")
}

// TestConvertEBOMToCSV_CellTrimming verifies that platform codes with trailing
// spaces are trimmed in the output CSV.
func TestConvertEBOMToCSV_CellTrimming(t *testing.T) {
	t.Parallel()
	path := newSingleSheetXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath: path,
		OutDir:   outDir,
	}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	var csvPath string
	for _, s := range res.Sheets {
		if s.AppendedTo == "" {
			csvPath = s.OutputPath
		}
	}
	require.NotEmpty(t, csvPath)

	data, readErr := os.ReadFile(csvPath)
	require.NoError(t, readErr)

	// "TST01     " should appear trimmed
	require.Contains(t, string(data), "TST01")
	require.NotContains(t, string(data), "TST01     ")
}

// TestConvertEBOMToCSV_NoEAssemblySheets verifies that a workbook with no
// eAssembly sheets returns an error.
func TestConvertEBOMToCSV_NoEAssemblySheets(t *testing.T) {
	t.Parallel()
	path := newNoEAssemblySheetsXLSX(t)
	outDir := t.TempDir()

	req := &ConvertEBOMToCSVReq{
		XLSXPath: path,
		OutDir:   outDir,
	}
	_, err := req.Run(t.Context())
	require.Error(t, err)
	require.Contains(t, err.Error(), "no eAssembly sheets")
}

// TestConvertEBOMToCSVReq_Validate verifies validation rejects missing path,
// .xls extension, and malformed --slugs values.
func TestConvertEBOMToCSVReq_Validate(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		req     *ConvertEBOMToCSVReq
		errFrag string
	}{
		"empty path": {
			req:     &ConvertEBOMToCSVReq{},
			errFrag: "xlsx path is required",
		},
		"xls extension": {
			req:     &ConvertEBOMToCSVReq{XLSXPath: "foo.xls"},
			errFrag: "not .xlsx",
		},
		"nonexistent file": {
			req:     &ConvertEBOMToCSVReq{XLSXPath: "nonexistent.xlsx"},
			errFrag: "nonexistent.xlsx",
		},
		"dot in slug": {
			req:     &ConvertEBOMToCSVReq{XLSXPath: "foo.xlsx", Slugs: []string{"0:essentials.openshift"}},
			errFrag: "dots are not permitted",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := tc.req.validate()
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.errFrag)
		})
	}
}

// TestParseSlugs verifies that parseSlugs accepts valid N:slug values and
// rejects malformed input including dots in the slug portion.
func TestParseSlugs(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		vals    []string
		want    map[int]string
		errFrag string
	}{
		"empty": {
			vals: nil,
			want: map[int]string{},
		},
		"single valid": {
			vals: []string{"0:essentials"},
			want: map[int]string{0: "essentials"},
		},
		"multiple valid": {
			vals: []string{"0:essentials", "1:standard", "2:premium"},
			want: map[int]string{0: "essentials", 1: "standard", 2: "premium"},
		},
		"hyphens allowed": {
			vals: []string{"3:essentials-for-red-hat-openshift"},
			want: map[int]string{3: "essentials-for-red-hat-openshift"},
		},
		"dot in slug rejected": {
			vals:    []string{"0:essentials.openshift"},
			errFrag: "dots are not permitted",
		},
		"no colon": {
			vals:    []string{"0essentials"},
			errFrag: "expected N:slug",
		},
		"colon at position zero": {
			vals:    []string{":essentials"},
			errFrag: "expected N:slug",
		},
		"N not an integer": {
			vals:    []string{"abc:essentials"},
			errFrag: "non-negative integer",
		},
		"negative N": {
			vals:    []string{"-1:essentials"},
			errFrag: "non-negative integer",
		},
		"leading hyphen in slug": {
			vals:    []string{"0:-essentials"},
			errFrag: "^[a-z0-9][a-z0-9-]*$",
		},
		"uppercase in slug": {
			vals:    []string{"0:Essentials"},
			errFrag: "^[a-z0-9][a-z0-9-]*$",
		},
		"empty slug": {
			vals:    []string{"0:"},
			errFrag: "^[a-z0-9][a-z0-9-]*$",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := parseSlugs(tc.vals)
			if tc.errFrag != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errFrag)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

// TestConvertEBOMToCSVRes_ToJSON verifies ToJSON returns valid JSON with
// expected fields.
func TestConvertEBOMToCSVRes_ToJSON(t *testing.T) {
	t.Parallel()
	res := &ConvertEBOMToCSVRes{
		Sheets: []*EBOMCSV{
			{SheetName: "Test", PID: "5900-TST", Slug: "essentials", OutputPath: "/tmp/out.csv", RowCount: 5},
		},
	}
	b, err := res.ToJSON()
	require.NoError(t, err)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))
	sheets, ok := parsed["sheets"]
	require.True(t, ok)
	require.NotEmpty(t, sheets)
}

// TestConvertEBOMToCSVRes_ToTable verifies ToTable contains expected column headers.
func TestConvertEBOMToCSVRes_ToTable(t *testing.T) {
	t.Parallel()
	res := &ConvertEBOMToCSVRes{
		Sheets: []*EBOMCSV{
			{SheetName: "Test", PID: "5900-TST", Slug: "essentials", OutputPath: "/tmp/out.csv", RowCount: 5},
		},
	}
	out := res.ToTable()
	require.Contains(t, out, "SHEET")
	require.Contains(t, out, "PRODUCT ID")
	require.Contains(t, out, "SLUG")
	require.Contains(t, out, "OUTPUT")
	require.Contains(t, out, "ROWS")
}

// TestConvertEBOMToCSVRes_ToMarkdown verifies ToMarkdown returns a Markdown table.
func TestConvertEBOMToCSVRes_ToMarkdown(t *testing.T) {
	t.Parallel()
	res := &ConvertEBOMToCSVRes{
		Sheets: []*EBOMCSV{
			{SheetName: "Test", PID: "5900-TST", Slug: "essentials", OutputPath: "/tmp/out.csv", RowCount: 5},
		},
	}
	out := res.ToMarkdown()
	require.Contains(t, out, "|")
	require.Contains(t, out, "SHEET")
}
