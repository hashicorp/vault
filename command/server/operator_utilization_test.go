// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build enterprise

package command

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestGetOutputFileName tests the getOutputFileName function for different values of output flag
func TestGetOutputFileName(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)
	inputTime := time.Now().UTC()
	formattedTime := inputTime.Format(fileFriendlyTimeFormat)
	cases := []struct {
		name             string
		flagOutput       string
		expectedFileName string
		expectedError    string
	}{
		{
			"empty output flag",
			"",
			fmt.Sprintf("vault-utilization-%s.json", formattedTime),
			"",
		},
		{
			"output file name ending with .json ",
			"test.json",
			"test.json",
			"",
		},
		{
			"output file path ending with .json ",
			fmt.Sprintf("%s/test.json", tempDir),
			fmt.Sprintf("%s/test.json", tempDir),
			"",
		},
		{
			"output file path ending with .jsn ",
			fmt.Sprintf("%s/test.jsn", tempDir),
			"",
			"invalid file extension .jsn, must be .json",
		},
		{
			"output file path ending with .json/ ",
			fmt.Sprintf("%s/test.json/", tempDir),
			fmt.Sprintf("%s/test.json", tempDir),
			"",
		},
		{
			"output flag only containing dir name",
			fmt.Sprintf("%s", tempDir),
			fmt.Sprintf("%s/vault-utilization-%s.json", tempDir, formattedTime),
			"",
		},
		{
			"output flag only containing dir name with trailing /",
			fmt.Sprintf("%s/", tempDir),
			fmt.Sprintf("%s/vault-utilization-%s.json", tempDir, formattedTime),
			"",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			filename, err := getOutputFileName(inputTime, tc.flagOutput)
			if tc.expectedError == "" {
				require.NoError(t, err)
				require.Equal(t, tc.expectedFileName, filename)

				// check if the temp dir has been created if needed
				dirName := filepath.Dir(filename)
				if len(dirName) != 0 {
					require.DirExists(t, dirName)
				}
			} else {
				require.Equal(t, tc.expectedError, err.Error())
			}
		})
	}
}
