// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

type mockCmdFetcher struct {
	data *releases.ListBinaryVersionsRes
	err  error
}

func (m *mockCmdFetcher) Fetch(ctx context.Context, versions string) (*releases.ListBinaryVersionsRes, error) {
	return m.data, m.err
}

func makeTestBinaryData(version string, ceActive bool) *releases.ListBinaryVersionsRes {
	allOS := releases.RequiredOSList
	linuxOnly := []string{"linux"}

	var variants []releases.VariantInfo
	if ceActive {
		variants = append(variants, releases.VariantInfo{Variant: version, OS: allOS})
	}
	variants = append(variants,
		releases.VariantInfo{Variant: version + "+ent", OS: allOS},
		releases.VariantInfo{Variant: version + "+ent.fips1403", OS: linuxOnly},
		releases.VariantInfo{Variant: version + "+ent.hsm", OS: linuxOnly},
		releases.VariantInfo{Variant: version + "+ent.hsm.fips1403", OS: linuxOnly},
	)

	return &releases.ListBinaryVersionsRes{
		AllVersions: []string{version},
		ValidVersions: map[string]struct {
			Status   string                 `json:"status"`
			Variants []releases.VariantInfo `json:"variants"`
		}{
			version: {
				Status:   "valid",
				Variants: variants,
			},
		},
		InvalidVersions: []string{},
	}
}

func TestReleasesVerifyBinariesCmd_FlagParsing(t *testing.T) {
	for name, tc := range map[string]struct {
		args        []string
		expectError bool
	}{
		"versions flag": {
			args: []string{"--versions", "2.2.1 1.21.13"},
		},
		"github output and fail on missing": {
			args: []string{"--versions", "2.2.1", "--github-output", "--fail-on-missing"},
		},
		"missing required versions flag returns error": {
			args:        []string{},
			expectError: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			cmd := newReleasesVerifyBinariesCmd()
			cmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
			cmd.SetArgs(tc.args)

			err := cmd.Execute()
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestReleasesVerifyBinariesCmd_ExecutionWithMockFetcher(t *testing.T) {
	validData := makeTestBinaryData("2.2.1", true)

	tmpOutput := filepath.Join(t.TempDir(), "github_output")
	t.Setenv("GITHUB_OUTPUT", tmpOutput)

	verifyBinariesCustomFetcher = &mockCmdFetcher{data: validData}
	t.Cleanup(func() {
		verifyBinariesCustomFetcher = nil
	})

	cmd := newReleasesVerifyBinariesCmd()
	cmd.SetArgs([]string{"--versions", "2.2.1", "--github-output"})

	err := cmd.Execute()
	require.NoError(t, err)

	outputContent, err := os.ReadFile(tmpOutput)
	require.NoError(t, err)
	require.Contains(t, string(outputContent), "release_binary_verification")
	require.Contains(t, string(outputContent), "*Available Versions:*")
}

func TestReleasesVerifyBinariesCmd_FormatJSON(t *testing.T) {
	validData := makeTestBinaryData("2.2.1", true)

	rootCfg.format = "json"
	t.Cleanup(func() {
		rootCfg.format = "table"
	})

	verifyBinariesCustomFetcher = &mockCmdFetcher{data: validData}
	t.Cleanup(func() {
		verifyBinariesCustomFetcher = nil
	})

	// Capture stdout with guaranteed restoration
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w
	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	cmd := newReleasesVerifyBinariesCmd()
	cmd.SetArgs([]string{"--versions", "2.2.1"})

	execErr := cmd.Execute()

	_ = w.Close()
	os.Stdout = oldStdout

	var outBuf bytes.Buffer
	_, _ = io.Copy(&outBuf, r)

	require.NoError(t, execErr)

	var jsonRes releases.VerifyBinariesRes
	err = json.Unmarshal(outBuf.Bytes(), &jsonRes)
	require.NoError(t, err, "stdout should be valid JSON: %s", outBuf.String())
	require.True(t, jsonRes.Success)
	require.False(t, jsonRes.HasMissing)
}

func TestReleasesVerifyBinariesCmd_EnterpriseOnlyLabelViaCLI(t *testing.T) {
	data := makeTestBinaryData("1.21.13", false)
	entData := &releases.ListBinaryVersionsRes{
		AllVersions: []string{"1.21.13-ent"},
		ValidVersions: map[string]struct {
			Status   string                 `json:"status"`
			Variants []releases.VariantInfo `json:"variants"`
		}{"1.21.13-ent": data.ValidVersions["1.21.13"]},
		InvalidVersions: []string{},
	}

	verifyBinariesCustomFetcher = &mockCmdFetcher{data: entData}
	t.Cleanup(func() { verifyBinariesCustomFetcher = nil })

	cmd := newReleasesVerifyBinariesCmd()
	cmd.SetArgs([]string{"--versions", "1.21.13-ent", "--fail-on-missing"})

	require.NoError(t, cmd.Execute())
}

func TestReleasesVerifyBinariesCmd_FailOnMissingFlag(t *testing.T) {
	invalidData := &releases.ListBinaryVersionsRes{
		AllVersions: []string{"9.9.9"},
		ValidVersions: map[string]struct {
			Status   string                 `json:"status"`
			Variants []releases.VariantInfo `json:"variants"`
		}{},
		InvalidVersions: []string{"9.9.9"},
	}

	verifyBinariesCustomFetcher = &mockCmdFetcher{data: invalidData}
	t.Cleanup(func() {
		verifyBinariesCustomFetcher = nil
	})

	cmd := newReleasesVerifyBinariesCmd()
	cmd.SetArgs([]string{"--versions", "9.9.9", "--fail-on-missing"})

	err := cmd.Execute()
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "binary verification failed"))
}
