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

// TestConvertZapReq_Run tests the request with several different property
// configurations.
func TestConvertZapReq_Run(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		sarifFile     string
		expectedSites int
		shouldFail    bool
	}{
		"basic SARIF with webRequest": {
			sarifFile:     "testdata/zap_basic.json",
			expectedSites: 1,
			shouldFail:    false,
		},
		"SARIF without webRequest": {
			sarifFile:     "testdata/zap_no_webrequest.json",
			expectedSites: 1,
			shouldFail:    false,
		},
		"multiple sites": {
			sarifFile:     "testdata/zap_multiple_sites.json",
			expectedSites: 2,
			shouldFail:    false,
		},
		"multiple risk levels": {
			sarifFile:     "testdata/zap_multiple_levels.json",
			expectedSites: 1,
			shouldFail:    false,
		},
		"empty results": {
			sarifFile:     "testdata/no_cve.json",
			expectedSites: 1,
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

			req := &ConvertZapReq{
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
			require.NotNil(t, res.Report)
			require.Equal(t, tt.expectedSites, len(res.Report.Sites))
		})
	}
}

// TestConvertZapRes_ToJSON validates JSON marshaling of ZAP report structure
func TestConvertZapRes_ToJSON(t *testing.T) {
	t.Parallel()

	res := &ConvertZapRes{
		Report: &ZapReport{
			ProgramName: "ZAP",
			Version:     "2.15.0",
			Generated:   "Thu, 02 Jul 2026 17:00:00 GMT",
			Sites: []*ZapSite{
				{
					Name: "https://example.com",
					Host: "example.com",
					Port: "443",
					SSL:  "true",
					Alerts: []*ZapAlert{
						{
							PluginID:   "CVE-2024-1234",
							AlertRef:   "CVE-2024-1234",
							Alert:      "Test vulnerability",
							Name:       "Test vulnerability",
							RiskCode:   "3",
							Confidence: "2",
							RiskDesc:   "High (Medium)",
							Desc:       "<p>Test vulnerability</p>",
							Instances: []*ZapInstance{
								{
									URI:       "https://example.com/test",
									Method:    "GET",
									Param:     "",
									Attack:    "",
									Evidence:  "",
									OtherInfo: "",
								},
							},
							Count:     "1",
							Solution:  "",
							OtherInfo: "",
							Reference: "",
							CWEID:     "-1",
							WASCID:    "-1",
							SourceID:  "0",
						},
					},
				},
			},
		},
	}

	output, err := res.ToJSON()
	require.NoError(t, err)

	var parsed ZapReport
	err = json.Unmarshal(output, &parsed)
	require.NoError(t, err)
	require.Equal(t, res.Report.ProgramName, parsed.ProgramName)
	require.Equal(t, len(res.Report.Sites), len(parsed.Sites))
}

// TestConvertZapRes_ToTable validates table formatting of ZAP report for
// human-readable output
func TestConvertZapRes_ToTable(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		report   *ZapReport
		contains []string
	}{
		"empty report": {
			report:   &ZapReport{Sites: []*ZapSite{}},
			contains: []string{"No sites found"},
		},
		"single site": {
			report: &ZapReport{
				Sites: []*ZapSite{
					{
						Name: "https://example.com",
						Host: "example.com",
						Port: "443",
						SSL:  "true",
						Alerts: []*ZapAlert{
							{
								Instances: []*ZapInstance{{}, {}},
							},
						},
					},
				},
			},
			contains: []string{"SITE", "example.com", "443", "true", "1", "2"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res := &ConvertZapRes{Report: tt.report}
			tbl, err := res.ToTable(nil)
			require.NoError(t, err)

			output := tbl.Render()
			for _, expected := range tt.contains {
				require.Contains(t, output, expected)
			}
		})
	}
}

// TestParseSite validates URI parsing and site information extraction for
// various URL formats
func TestParseSite(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		uri      string
		expected siteInfo
	}{
		"https with default port": {
			uri: "https://example.com/path",
			expected: siteInfo{
				name: "https://example.com",
				host: "example.com",
				port: "443",
				ssl:  "true",
			},
		},
		"https with custom port": {
			uri: "https://example.com:8443/path",
			expected: siteInfo{
				name: "https://example.com:8443",
				host: "example.com",
				port: "8443",
				ssl:  "true",
			},
		},
		"http with default port": {
			uri: "http://example.com/path",
			expected: siteInfo{
				name: "http://example.com",
				host: "example.com",
				port: "80",
				ssl:  "false",
			},
		},
		"http with custom port": {
			uri: "http://example.com:8080/path",
			expected: siteInfo{
				name: "http://example.com:8080",
				host: "example.com",
				port: "8080",
				ssl:  "false",
			},
		},
		"invalid URI": {
			uri: "not a valid uri",
			expected: siteInfo{
				name: "unknown",
				host: "unknown",
				port: "80",
				ssl:  "false",
			},
		},
		"localhost": {
			uri: "http://localhost:3000/api",
			expected: siteInfo{
				name: "http://localhost:3000",
				host: "localhost",
				port: "3000",
				ssl:  "false",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := parseSite(tt.uri)
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestGetRiskInfo validates SARIF severity level to ZAP risk code mapping
func TestGetRiskInfo(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		level    string
		expected riskInfo
	}{
		"note": {
			level:    "note",
			expected: riskInfo{code: "0", label: "Informational"},
		},
		"warning": {
			level:    "warning",
			expected: riskInfo{code: "1", label: "Low"},
		},
		"error": {
			level:    "error",
			expected: riskInfo{code: "2", label: "Medium"},
		},
		"critical": {
			level:    "critical",
			expected: riskInfo{code: "3", label: "High"},
		},
		"unknown": {
			level:    "unknown",
			expected: riskInfo{code: "0", label: "Informational"},
		},
		"uppercase": {
			level:    "ERROR",
			expected: riskInfo{code: "2", label: "Medium"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := getRiskInfo(tt.level)
			require.Equal(t, tt.expected, result)
		})
	}
}

// TestZapOutputMatchesNode validates Go implementation output matches Node.js reference implementation
// by comparing converted ZAP reports against pre-generated expected outputs
func TestZapOutputMatchesNode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "testdata/zap_basic.json",
			expected: "testdata/expected/zap_basic.json",
		},
		{
			input:    "testdata/zap_no_webrequest.json",
			expected: "testdata/expected/zap_no_webrequest.json",
		},
		{
			input:    "testdata/zap_multiple_sites.json",
			expected: "testdata/expected/zap_multiple_sites.json",
		},
		{
			input:    "testdata/zap_multiple_levels.json",
			expected: "testdata/expected/zap_multiple_levels.json",
		},
	}

	for _, tc := range testCases {
		t.Run(filepath.Base(tc.input), func(t *testing.T) {
			t.Parallel()

			// Skip if expected file doesn't exist yet
			if _, err := os.Stat(tc.expected); os.IsNotExist(err) {
				t.Skipf("Expected output file not found: %s", tc.expected)
				return
			}

			req := &ConvertZapReq{SarifPath: tc.input}
			res, err := req.Run(context.Background())
			require.NoError(t, err)

			goJSON, err := res.ToJSON()
			require.NoError(t, err)

			expectedJSON, err := os.ReadFile(tc.expected)
			require.NoError(t, err)

			// Parse both JSONs to compare structure (ignore formatting differences)
			var goReport, expectedReport ZapReport
			err = json.Unmarshal(goJSON, &goReport)
			require.NoError(t, err)
			err = json.Unmarshal(expectedJSON, &expectedReport)
			require.NoError(t, err)

			// Compare key fields (ignore Generated timestamp)
			require.Equal(t, expectedReport.ProgramName, goReport.ProgramName)
			require.Equal(t, expectedReport.Version, goReport.Version)
			require.Equal(t, len(expectedReport.Sites), len(goReport.Sites))

			// Compare sites and alerts
			for i, expectedSite := range expectedReport.Sites {
				goSite := goReport.Sites[i]
				require.Equal(t, expectedSite.Name, goSite.Name)
				require.Equal(t, expectedSite.Host, goSite.Host)
				require.Equal(t, expectedSite.Port, goSite.Port)
				require.Equal(t, expectedSite.SSL, goSite.SSL)
				require.Equal(t, len(expectedSite.Alerts), len(goSite.Alerts))

				for j, expectedAlert := range expectedSite.Alerts {
					goAlert := goSite.Alerts[j]
					require.Equal(t, expectedAlert.PluginID, goAlert.PluginID)
					require.Equal(t, expectedAlert.Alert, goAlert.Alert)
					require.Equal(t, expectedAlert.RiskCode, goAlert.RiskCode)
					require.Equal(t, expectedAlert.Count, goAlert.Count)
					require.Equal(t, len(expectedAlert.Instances), len(goAlert.Instances))
				}
			}
		})
	}
}
