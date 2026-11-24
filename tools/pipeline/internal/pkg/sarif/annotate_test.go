// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package sarif

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnnotateReq_Run(t *testing.T) {
	tests := map[string]struct {
		inputFile     string
		keyValuePairs []string
		wantErr       bool
		wantRuns      int
	}{
		"basic annotation with flat keys": {
			inputFile:     "annotate_basic.json",
			keyValuePairs: []string{"product:vault-enterprise", "version:1.21.2"},
			wantErr:       false,
			wantRuns:      1,
		},
		"annotation with nested keys": {
			inputFile:     "annotate_basic.json",
			keyValuePairs: []string{"product.name:vault", "product.version:1.21.0"},
			wantErr:       false,
			wantRuns:      1,
		},
		"multiple runs": {
			inputFile:     "annotate_multiple_runs.json",
			keyValuePairs: []string{"product:vault", "version:1.21.2"},
			wantErr:       false,
			wantRuns:      2,
		},
		"preserve existing properties": {
			inputFile:     "annotate_existing_props.json",
			keyValuePairs: []string{"product:vault", "version:1.21.2"},
			wantErr:       false,
			wantRuns:      1,
		},
		"no key:value pairs": {
			inputFile:     "annotate_basic.json",
			keyValuePairs: []string{},
			wantErr:       true,
		},
		"invalid format - no colon": {
			inputFile:     "annotate_basic.json",
			keyValuePairs: []string{"productvalue"},
			wantErr:       true,
		},
		"invalid format - empty key": {
			inputFile:     "annotate_basic.json",
			keyValuePairs: []string{":value"},
			wantErr:       true,
		},
		"invalid format - empty value": {
			inputFile:     "annotate_basic.json",
			keyValuePairs: []string{"key:"},
			wantErr:       true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := &AnnotateReq{
				SarifPath:     filepath.Join("testdata", tt.inputFile),
				KeyValuePairs: tt.keyValuePairs,
			}

			res, err := req.Run(context.Background())

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)

			// Verify properties were added to all runs
			for _, run := range res.Report.Runs {
				require.NotNil(t, run.Properties)
				require.NotNil(t, run.Properties.Properties)

				// Check that properties from key:value pairs exist
				for _, kvp := range tt.keyValuePairs {
					// Parse the key:value pair
					key, _, _ := parseKeyValue(kvp)
					if key != "" {
						// For nested keys, just check the top-level key exists
						topKey, _, _ := strings.Cut(key, ".")
						_, exists := run.Properties.Properties[topKey]
						require.True(t, exists, "property %s should exist", topKey)
					}
				}
			}
		})
	}
}

func TestAnnotateReq_PreserveExistingProperties(t *testing.T) {
	req := &AnnotateReq{
		SarifPath:     filepath.Join("testdata", "annotate_existing_props.json"),
		KeyValuePairs: []string{"product:vault", "version:1.21.0"},
	}

	res, err := req.Run(context.Background())
	require.NoError(t, err)

	// Verify existing properties are preserved
	run := res.Report.Runs[0]
	existingVal, ok := run.Properties.Properties["existingKey"]
	require.True(t, ok, "existing property should be preserved")
	require.Equal(t, "existingValue", existingVal)

	anotherVal, ok := run.Properties.Properties["anotherKey"]
	require.True(t, ok, "another existing property should be preserved")
	require.Equal(t, "anotherValue", anotherVal)

	// Verify new properties were added
	productVal, ok := run.Properties.Properties["product"]
	require.True(t, ok, "product property should be added")
	require.Equal(t, "vault", productVal)

	versionVal, ok := run.Properties.Properties["version"]
	require.True(t, ok, "version property should be added")
	require.Equal(t, "1.21.0", versionVal)
}

func TestAnnotateReq_NestedProperties(t *testing.T) {
	req := &AnnotateReq{
		SarifPath:     filepath.Join("testdata", "annotate_basic.json"),
		KeyValuePairs: []string{"product.name:vault-enterprise", "product.version:1.21.2"},
	}

	res, err := req.Run(context.Background())
	require.NoError(t, err)

	// Verify nested structure was created
	run := res.Report.Runs[0]
	productProp, ok := run.Properties.Properties["product"]
	require.True(t, ok, "product property should exist")

	productMap, ok := productProp.(map[string]any)
	require.True(t, ok, "product should be a map")
	require.Equal(t, "vault-enterprise", productMap["name"])
	require.Equal(t, "1.21.2", productMap["version"])
}

func TestAnnotateReq_ToJSON(t *testing.T) {
	req := &AnnotateReq{
		SarifPath:     filepath.Join("testdata", "annotate_basic.json"),
		KeyValuePairs: []string{"product:vault", "version:1.21.2"},
	}

	res, err := req.Run(context.Background())
	require.NoError(t, err)

	jsonBytes, err := res.ToJSON()
	require.NoError(t, err)
	require.NotEmpty(t, jsonBytes)

	// Verify it's valid JSON
	var parsed map[string]any
	err = json.Unmarshal(jsonBytes, &parsed)
	require.NoError(t, err)
}

func TestAnnotateReq_WriteInPlace(t *testing.T) {
	// Create a temporary copy of the test file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.json")

	srcData, err := os.ReadFile(filepath.Join("testdata", "annotate_basic.json"))
	require.NoError(t, err)

	err = os.WriteFile(tmpFile, srcData, 0o644)
	require.NoError(t, err)

	req := &AnnotateReq{
		SarifPath:     tmpFile,
		KeyValuePairs: []string{"product:vault", "version:1.21.2"},
		Write:         true,
	}

	res, err := req.Run(context.Background())
	require.NoError(t, err)
	require.True(t, res.Updated)
	require.Equal(t, tmpFile, res.OutputPath)

	// Write the output
	output, err := res.ToJSON()
	require.NoError(t, err)

	err = os.WriteFile(res.OutputPath, output, 0o644)
	require.NoError(t, err)

	// Verify the file was updated
	updatedData, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	var parsed map[string]any
	err = json.Unmarshal(updatedData, &parsed)
	require.NoError(t, err)

	runs, ok := parsed["runs"].([]any)
	require.True(t, ok)
	require.Len(t, runs, 1)

	run, ok := runs[0].(map[string]any)
	require.True(t, ok)

	props, ok := run["properties"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, "vault", props["product"])
	require.Equal(t, "1.21.2", props["version"])
}

// Helper functions for tests
func parseKeyValue(kvp string) (string, string, error) {
	parts := strings.SplitN(kvp, ":", 2)
	if len(parts) != 2 {
		return "", "", nil
	}
	return parts[0], parts[1], nil
}
