// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package sarif

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestCVEReportReq_Run verifies that CVE extraction from SARIF files works correctly
// across various scenarios including empty results, single/multiple CVEs, deduplication,
// and error handling for invalid inputs.
func TestCVEReportReq_Run(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		sarifFile     string
		expectedCount int
		shouldFail    bool
	}{
		"no CVE relationships": {
			sarifFile:     "testdata/no_cve.json",
			expectedCount: 0,
			shouldFail:    false,
		},
		"single CVE finding": {
			sarifFile:     "testdata/single_cve.json",
			expectedCount: 1,
			shouldFail:    false,
		},
		"multiple CVE findings": {
			sarifFile:     "testdata/multiple_cve.json",
			expectedCount: 3,
			shouldFail:    false,
		},
		"duplicate CVE findings (deduplication)": {
			sarifFile:     "testdata/duplicate_cve.json",
			expectedCount: 2, // Should deduplicate
			shouldFail:    false,
		},
		"invalid file path": {
			sarifFile:  "nonexistent.json",
			shouldFail: true,
		},
		"invalid JSON": {
			sarifFile:  "testdata/invalid.json",
			shouldFail: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &CVEReportReq{
				SarifPath:  tt.sarifFile,
				OutputPath: "",
			}

			res, err := req.Run(context.Background())

			if tt.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)
			require.Equal(t, tt.expectedCount, len(res.CVEs))
		})
	}
}

// TestCVEReportRes_ToJSON verifies that CVE report results can be correctly serialized
// to JSON format and that the output can be parsed back into the expected structure.
func TestCVEReportRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &CVEReportRes{
		CVEs: []*CVE{
			{
				Vulnerability:         "CVE-2024-1234",
				PackageLibraryName:    "",
				PackageLibraryVersion: "",
				RuleID:                "rule-1",
			},
		},
	}

	output, err := res.ToJSON()
	require.NoError(t, err)

	var parsed []*CVE
	err = json.Unmarshal(output, &parsed)
	require.NoError(t, err)
	require.Equal(t, len(res.CVEs), len(parsed))
	require.Equal(t, res.CVEs[0].Vulnerability, parsed[0].Vulnerability)
}

// TestCVEReportRes_ToCSV verifies that CVE report results are correctly converted to CSV
// format with proper headers and formatting for empty, single, and multiple CVE scenarios.
func TestCVEReportRes_ToCSV(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cves     []*CVE
		expected string
	}{
		"empty cves": {
			cves:     []*CVE{},
			expected: "",
		},
		"single finding": {
			cves: []*CVE{
				{
					Vulnerability:         "CVE-2024-1234",
					PackageLibraryName:    "",
					PackageLibraryVersion: "",
				},
			},
			expected: "vulnerability,packageLibraryName,packageLibraryVersion\nCVE-2024-1234,,\n",
		},
		"multiple cves": {
			cves: []*CVE{
				{Vulnerability: "CVE-2024-1234"},
				{Vulnerability: "CVE-2024-5678"},
			},
			expected: "vulnerability,packageLibraryName,packageLibraryVersion\nCVE-2024-1234,,\nCVE-2024-5678,,\n",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res := &CVEReportRes{CVEs: tt.cves}
			output, err := res.ToCSV()
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(output))
		})
	}
}

// TestCVEReportRes_ToTable verifies that CVE report results are correctly rendered as
// formatted tables, including proper handling of empty results and CVE data display.
func TestCVEReportRes_ToTable(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cves     []*CVE
		contains []string
	}{
		"empty cves": {
			cves:     []*CVE{},
			contains: []string{"No CVE-linked findings"},
		},
		"single finding": {
			cves: []*CVE{
				{Vulnerability: "CVE-2024-1234"},
			},
			contains: []string{"VULNERABILITY", "CVE-2024-1234"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res := &CVEReportRes{CVEs: tt.cves}
			tbl, err := res.ToTable(nil)
			require.NoError(t, err)

			output := tbl.Render()
			for _, expected := range tt.contains {
				require.Contains(t, output, expected)
			}
		})
	}
}

// TestCSVOutputMatchesExpected verifies that CSV output generated from SARIF files
// matches pre-generated expected outputs for various test cases including no CVEs,
// single CVE, multiple CVEs, and duplicate CVE scenarios.
func TestCSVOutputMatchesExpected(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "testdata/no_cve.json",
			expected: "testdata/expected/no_cve.csv",
		},
		{
			input:    "testdata/single_cve.json",
			expected: "testdata/expected/single_cve.csv",
		},
		{
			input:    "testdata/multiple_cve.json",
			expected: "testdata/expected/multiple_cve.csv",
		},
		{
			input:    "testdata/duplicate_cve.json",
			expected: "testdata/expected/duplicate_cve.csv",
		},
	}

	for _, tc := range testCases {
		t.Run(filepath.Base(tc.input), func(t *testing.T) {
			t.Parallel()

			req := &CVEReportReq{SarifPath: tc.input}
			res, err := req.Run(context.Background())
			require.NoError(t, err)

			goCSV, err := res.ToCSV()
			require.NoError(t, err)

			expectedCSV, err := os.ReadFile(tc.expected)
			require.NoError(t, err)

			// Compare outputs
			require.Equal(t, string(expectedCSV), string(goCSV),
				"CSV output should match expected for %s", tc.input)
		})
	}
}
