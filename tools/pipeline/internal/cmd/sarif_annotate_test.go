// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/sarif"
	"github.com/stretchr/testify/require"
)

// TestSarifAnnotateCmd_DefaultTableOutput verifies that the command outputs a
// formatted table to stdout by default when no --format flag is specified.
func TestSarifAnnotateCmd_DefaultTableOutput(t *testing.T) {
	testFile := filepath.Join("..", "pkg", "sarif", "testdata", "annotate_basic.json")

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	annotateReq = &sarif.AnnotateReq{}

	cmd := newSarifAnnotateCmd()
	cmd.SetArgs([]string{testFile, "product:vault", "version:1.21.2"})

	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	output := buf.String()

	// Verify it's table format (not JSON)
	require.NotEmpty(t, output)
	require.Contains(t, output, "KEY")
	require.Contains(t, output, "VALUE")
	require.Contains(t, output, "product")
	require.Contains(t, output, "vault")
	require.Contains(t, output, "version")
	require.Contains(t, output, "1.21.2")

	// Should NOT be valid JSON
	var jsonCheck map[string]any
	err = json.Unmarshal([]byte(output), &jsonCheck)
	require.Error(t, err, "default output should be table, not JSON")
}

// TestSarifAnnotateCmd_FormatJSON verifies that the --format json flag causes
// the command to output valid SARIF JSON to stdout with annotations applied.
func TestSarifAnnotateCmd_FormatJSON(t *testing.T) {
	testFile := filepath.Join("..", "pkg", "sarif", "testdata", "annotate_basic.json")

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	annotateReq = &sarif.AnnotateReq{}
	rootCfg.format = "json"
	defer func() { rootCfg.format = "" }()

	cmd := newSarifAnnotateCmd()
	cmd.SetArgs([]string{testFile, "product:vault", "version:1.21.2"})

	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	output := buf.String()

	require.NotEmpty(t, output)
	var parsed map[string]any
	err = json.Unmarshal([]byte(output), &parsed)
	require.NoError(t, err, "output should be valid JSON")

	require.Contains(t, parsed, "version")
	require.Contains(t, parsed, "runs")

	runs, ok := parsed["runs"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, runs)

	run := runs[0].(map[string]any)
	props, ok := run["properties"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "vault", props["product"])
	require.Equal(t, "1.21.2", props["version"])
}

// TestSarifAnnotateCmd_FormatMarkdown verifies that the --format markdown flag
// causes the command to output a markdown-formatted table to stdout.
func TestSarifAnnotateCmd_FormatMarkdown(t *testing.T) {
	testFile := filepath.Join("..", "pkg", "sarif", "testdata", "annotate_basic.json")

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	annotateReq = &sarif.AnnotateReq{}
	rootCfg.format = "markdown"
	defer func() { rootCfg.format = "" }()

	cmd := newSarifAnnotateCmd()
	cmd.SetArgs([]string{testFile, "product:vault", "version:1.21.2"})

	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	output := buf.String()

	require.NotEmpty(t, output)
	require.Contains(t, output, "|")
	require.Contains(t, output, "KEY")
	require.Contains(t, output, "VALUE")
	require.Contains(t, output, "product")
	require.Contains(t, output, "vault")
}

// TestSarifAnnotateCmd_WriteInPlace verifies that the --write flag updates the
// source SARIF file in-place with annotations while displaying a table summary to stdout.
func TestSarifAnnotateCmd_WriteInPlace(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	srcData, err := os.ReadFile(filepath.Join("..", "pkg", "sarif", "testdata", "annotate_basic.json"))
	require.NoError(t, err)
	err = os.WriteFile(testFile, srcData, 0o644)
	require.NoError(t, err)

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	annotateReq = &sarif.AnnotateReq{}

	cmd := newSarifAnnotateCmd()
	cmd.SetArgs([]string{testFile, "product:vault", "version:1.21.2", "--write"})

	err = cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	stdoutOutput := buf.String()

	// Verify stdout shows table summary (not JSON)
	require.NotEmpty(t, stdoutOutput)
	require.Contains(t, stdoutOutput, "Updated in-place")

	updatedData, err := os.ReadFile(testFile)
	require.NoError(t, err)

	var parsed map[string]any
	err = json.Unmarshal(updatedData, &parsed)
	require.NoError(t, err, "updated file should be valid JSON")

	runs, ok := parsed["runs"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, runs)

	run := runs[0].(map[string]any)
	props, ok := run["properties"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "vault", props["product"])
	require.Equal(t, "1.21.2", props["version"])
}

// TestSarifAnnotateCmd_OutputPath verifies that the --out flag writes the annotated
// SARIF as JSON to a different file while leaving the original file unchanged.
func TestSarifAnnotateCmd_OutputPath(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join("..", "pkg", "sarif", "testdata", "annotate_basic.json")
	outFile := filepath.Join(tmpDir, "output.json")

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	annotateReq = &sarif.AnnotateReq{}

	cmd := newSarifAnnotateCmd()
	cmd.SetArgs([]string{testFile, "product:vault", "version:1.21.2", "--out", outFile})

	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	stdoutOutput := buf.String()

	// Verify stdout shows table summary with output path
	require.NotEmpty(t, stdoutOutput)
	require.Contains(t, stdoutOutput, "Output")
	require.Contains(t, stdoutOutput, outFile)

	outputData, err := os.ReadFile(outFile)
	require.NoError(t, err)

	var parsed map[string]any
	err = json.Unmarshal(outputData, &parsed)
	require.NoError(t, err, "output file should be valid JSON")

	runs, ok := parsed["runs"].([]any)
	require.True(t, ok)
	require.NotEmpty(t, runs)

	run := runs[0].(map[string]any)
	props, ok := run["properties"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "vault", props["product"])
	require.Equal(t, "1.21.2", props["version"])

	// Verify original file was NOT modified
	origData, err := os.ReadFile(testFile)
	require.NoError(t, err)

	var origParsed map[string]any
	err = json.Unmarshal(origData, &origParsed)
	require.NoError(t, err)

	origRuns := origParsed["runs"].([]any)
	origRun := origRuns[0].(map[string]any)
	origProps, hasProps := origRun["properties"].(map[string]any)

	if hasProps {
		_, hasProduct := origProps["product"]
		_, hasVersion := origProps["version"]
		require.False(t, hasProduct && hasVersion, "original file should not be modified")
	}
}

// TestSarifAnnotateCmd_WriteAndOutMutuallyExclusive verifies that the --write
// and --out flags are mutually exclusive and produce an error when used together.
func TestSarifAnnotateCmd_WriteAndOutMutuallyExclusive(t *testing.T) {
	testFile := filepath.Join("..", "pkg", "sarif", "testdata", "annotate_basic.json")
	outFile := filepath.Join(t.TempDir(), "output.json")

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	_, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	os.Stderr = w
	defer func() {
		w.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	annotateReq = &sarif.AnnotateReq{}

	cmd := newSarifAnnotateCmd()
	cmd.SetArgs([]string{testFile, "product:vault", "version:1.21.2", "--write", "--out", outFile})

	err := cmd.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "write")
	require.Contains(t, err.Error(), "out")
}

// TestSarifAnnotateCmd_FormatJSONWithWrite verifies that combining --format json
// with --write outputs JSON to both stdout and the source file.
func TestSarifAnnotateCmd_FormatJSONWithWrite(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	srcData, err := os.ReadFile(filepath.Join("..", "pkg", "sarif", "testdata", "annotate_basic.json"))
	require.NoError(t, err)
	err = os.WriteFile(testFile, srcData, 0o644)
	require.NoError(t, err)

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	annotateReq = &sarif.AnnotateReq{}
	rootCfg.format = "json"
	defer func() { rootCfg.format = "" }()

	cmd := newSarifAnnotateCmd()
	cmd.SetArgs([]string{testFile, "product:vault", "version:1.21.2", "--write"})

	err = cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	stdoutOutput := buf.String()

	// Verify stdout contains JSON (not table)
	require.NotEmpty(t, stdoutOutput)
	var stdoutParsed map[string]any
	err = json.Unmarshal([]byte(stdoutOutput), &stdoutParsed)
	require.NoError(t, err, "stdout should be valid JSON when --format json is used")

	updatedData, err := os.ReadFile(testFile)
	require.NoError(t, err)

	var fileParsed map[string]any
	err = json.Unmarshal(updatedData, &fileParsed)
	require.NoError(t, err, "updated file should be valid JSON")

	stdoutRuns := stdoutParsed["runs"].([]any)
	stdoutRun := stdoutRuns[0].(map[string]any)
	stdoutProps := stdoutRun["properties"].(map[string]any)

	fileRuns := fileParsed["runs"].([]any)
	fileRun := fileRuns[0].(map[string]any)
	fileProps := fileRun["properties"].(map[string]any)

	require.Equal(t, "vault", stdoutProps["product"])
	require.Equal(t, "vault", fileProps["product"])
	require.Equal(t, "1.21.2", stdoutProps["version"])
	require.Equal(t, "1.21.2", fileProps["version"])
}

// TestSarifAnnotateCmd_PreserveFileMode verifies that the --write flag preserves
// the original file permissions when updating a file in-place.
func TestSarifAnnotateCmd_PreserveFileMode(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	srcData, err := os.ReadFile(filepath.Join("..", "pkg", "sarif", "testdata", "annotate_basic.json"))
	require.NoError(t, err)
	err = os.WriteFile(testFile, srcData, 0o600)
	require.NoError(t, err)

	info, err := os.Stat(testFile)
	require.NoError(t, err)
	originalMode := info.Mode()
	require.Equal(t, os.FileMode(0o600), originalMode.Perm())

	annotateReq = &sarif.AnnotateReq{}

	cmd := newSarifAnnotateCmd()
	cmd.SetArgs([]string{testFile, "product:vault", "version:1.21.2", "--write"})

	oldStdout := os.Stdout
	_, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = oldStdout
	}()

	err = cmd.Execute()
	require.NoError(t, err)

	info, err = os.Stat(testFile)
	require.NoError(t, err)
	newMode := info.Mode()
	require.Equal(t, originalMode.Perm(), newMode.Perm(), "file mode should be preserved")
}

// TestSarifAnnotateCmd_NestedProperties verifies that dot notation in key:value
// pairs correctly creates nested property structures in the SARIF output.
func TestSarifAnnotateCmd_NestedProperties(t *testing.T) {
	testFile := filepath.Join("..", "pkg", "sarif", "testdata", "annotate_basic.json")

	oldStdout := os.Stdout
	r, w, pipeErr := os.Pipe()
	require.NoError(t, pipeErr)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	annotateReq = &sarif.AnnotateReq{}
	rootCfg.format = "json"
	defer func() { rootCfg.format = "" }()

	cmd := newSarifAnnotateCmd()
	cmd.SetArgs([]string{testFile, "product.name:vault", "product.version:1.21.2", "environment:prod"})

	err := cmd.Execute()
	require.NoError(t, err)

	w.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	output := buf.String()

	var parsed map[string]any
	err = json.Unmarshal([]byte(output), &parsed)
	require.NoError(t, err)

	runs := parsed["runs"].([]any)
	run := runs[0].(map[string]any)
	props := run["properties"].(map[string]any)

	product, ok := props["product"].(map[string]any)
	require.True(t, ok, "product should be a nested map")
	require.Equal(t, "vault", product["name"])
	require.Equal(t, "1.21.2", product["version"])

	require.Equal(t, "prod", props["environment"])
}
