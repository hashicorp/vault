// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package ebom

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestInspectEBOM_SingleSheet verifies that a single-sheet workbook
// produces one entry with an empty InferredSlug (no tier tokens in description)
// and a non-empty SlugOverride placeholder.
func TestInspectEBOM_SingleSheet(t *testing.T) {
	t.Parallel()
	path := newSingleSheetXLSX(t)

	req := &InspectEBOMReq{XLSXPath: path}
	res, err := req.Run(t.Context())
	require.NoError(t, err)
	require.Len(t, res.Sheets, 1)
	require.Equal(t, 0, res.Sheets[0].Index)
	require.Empty(t, res.Sheets[0].InferredSlug)
	require.Contains(t, res.Sheets[0].SlugOverride, "--slugs 0:<slug>")
}

// TestInspectEBOM_MultiSheet verifies that multi-sheet workbooks report
// correct indices, descriptions, and inferred slugs for each eAssembly sheet.
// When a slug is successfully inferred, SlugOverride is empty.
func TestInspectEBOM_MultiSheet(t *testing.T) {
	t.Parallel()
	path := newMultiSheetXLSX(t)

	req := &InspectEBOMReq{XLSXPath: path}
	res, err := req.Run(t.Context())
	require.NoError(t, err)
	require.Len(t, res.Sheets, 3)

	require.Equal(t, 0, res.Sheets[0].Index)
	require.Equal(t, "essentials", res.Sheets[0].InferredSlug)
	require.Empty(t, res.Sheets[0].SlugOverride)

	require.Equal(t, 1, res.Sheets[1].Index)
	require.Equal(t, "standard", res.Sheets[1].InferredSlug)

	require.Equal(t, 2, res.Sheets[2].Index)
	require.Equal(t, "premium", res.Sheets[2].InferredSlug)
}

// TestInspectEBOM_InferFail verifies that when a description produces no
// slug tokens, InferredSlug is empty and SlugOverride contains the placeholder.
func TestInspectEBOM_InferFail(t *testing.T) {
	t.Parallel()
	path := newSingleSheetXLSX(t)

	req := &InspectEBOMReq{XLSXPath: path}
	res, err := req.Run(t.Context())
	require.NoError(t, err)
	require.Len(t, res.Sheets, 1)
	require.Empty(t, res.Sheets[0].InferredSlug)
	require.Contains(t, res.Sheets[0].SlugOverride, "<slug>")
}

// TestInspectEBOM_PIDOverride verifies that PIDOverride replaces the
// workbook PID in every result sheet.
func TestInspectEBOM_PIDOverride(t *testing.T) {
	t.Parallel()
	path := newSingleSheetXLSX(t)

	req := &InspectEBOMReq{XLSXPath: path, PIDOverride: "5900-OVR"}
	res, err := req.Run(t.Context())
	require.NoError(t, err)
	require.Len(t, res.Sheets, 1)
	require.Equal(t, "5900-OVR", res.Sheets[0].PID)
}

// TestInspectEBOM_SkipsNonEBOM verifies that non-eAssembly sheets are not
// present in the result set.
func TestInspectEBOM_SkipsNonEBOM(t *testing.T) {
	t.Parallel()
	path := newMixedSheetsXLSX(t)

	req := &InspectEBOMReq{XLSXPath: path}
	res, err := req.Run(t.Context())
	require.NoError(t, err)
	require.Len(t, res.Sheets, 1)
	require.Equal(t, "Vault Essen", res.Sheets[0].SheetName)
}

// TestInspectEBOMReq_Validate verifies validation rejects missing path and
// the .xls extension.
func TestInspectEBOMReq_Validate(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		req     *InspectEBOMReq
		errFrag string
	}{
		"empty path": {
			req:     &InspectEBOMReq{},
			errFrag: "xlsx path is required",
		},
		"xls extension": {
			req:     &InspectEBOMReq{XLSXPath: "foo.xls"},
			errFrag: "not .xlsx",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := tc.req.validate()
			require.Error(t, err)
			require.Contains(t, err.Error(), tc.errFrag)
		})
	}
}

// TestInspectEBOMRes_ToJSON verifies ToJSON returns valid JSON with the
// index, inferredSlug, and slugOverride fields.
func TestInspectEBOMRes_ToJSON(t *testing.T) {
	t.Parallel()
	path := newMultiSheetXLSX(t)

	req := &InspectEBOMReq{XLSXPath: path}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	b, jsonErr := res.ToJSON()
	require.NoError(t, jsonErr)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal(b, &parsed))

	sheets, ok := parsed["sheets"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, sheets)

	first := sheets[0].(map[string]any)
	_, hasIndex := first["index"]
	_, hasSlug := first["inferredSlug"]
	_, hasOverride := first["slugOverride"]
	require.True(t, hasIndex)
	require.True(t, hasSlug)
	require.True(t, hasOverride)
}

// TestInspectEBOMRes_ToTable verifies that multi-sheet ToTable output
// contains the INDEX, INFERRED SLUG, and SLUG OVERRIDE columns.
func TestInspectEBOMRes_ToTable(t *testing.T) {
	t.Parallel()
	path := newMultiSheetUnknownSlugXLSX(t)

	req := &InspectEBOMReq{XLSXPath: path}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	out := res.ToTable()
	require.Contains(t, out, "INDEX")
	require.Contains(t, out, "INFERRED SLUG")
	require.Contains(t, out, "SLUG OVERRIDE")
}

// TestInspectEBOMRes_ToTable_AllInferred verifies that when all slugs are
// inferred the SLUG OVERRIDE column is suppressed entirely (nothing to override).
func TestInspectEBOMRes_ToTable_AllInferred(t *testing.T) {
	t.Parallel()
	path := newMultiSheetXLSX(t)

	req := &InspectEBOMReq{XLSXPath: path}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	out := res.ToTable()
	require.NotContains(t, out, "pipeline ebom convert")
	require.NotContains(t, out, "SLUG OVERRIDE")
	require.Contains(t, out, "essentials")
}

// TestInspectEBOMRes_ToTable_MissingSlug verifies that when a multi-sheet
// workbook has at least one un-inferable slug the SLUG OVERRIDE column appears
// with the --slugs placeholder for that sheet.
func TestInspectEBOMRes_ToTable_MissingSlug(t *testing.T) {
	t.Parallel()
	path := newMultiSheetUnknownSlugXLSX(t)

	req := &InspectEBOMReq{XLSXPath: path}
	res, err := req.Run(t.Context())
	require.NoError(t, err)

	out := res.ToTable()
	require.NotContains(t, out, "pipeline ebom convert")
	require.Contains(t, out, "SLUG OVERRIDE")
	require.Contains(t, out, "--slugs 1:<slug>")
}
