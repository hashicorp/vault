// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package releases

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const validVersionsHCL = `schema = 1
active_versions {
  version "1.19.x" {
    ce_active = true
  }

  version "1.18.x" {
    ce_active = false
  }

  version "1.17.x" {
    ce_active = false
  }

  version "1.16.x" {
    ce_active = false
    lts       = true
  }
}
`

// TestDecode tests the Decode function with various HCL inputs
func TestDecode(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		input      string
		expected   *VersionsConfig
		shouldFail bool
	}{
		"valid config with schema": {
			input: validVersionsHCL,
			expected: &VersionsConfig{
				Schema: 1,
				ActiveVersion: &ActiveVersion{
					Versions: map[string]*Version{
						"1.19.x": {CEActive: true, LTS: false},
						"1.18.x": {CEActive: false, LTS: false},
						"1.17.x": {CEActive: false, LTS: false},
						"1.16.x": {CEActive: false, LTS: true},
					},
				},
			},
			shouldFail: false,
		},
		"config without schema field": {
			input: `active_versions {
  version "1.19.x" {
    ce_active = true
    lts       = true
  }
}
`,
			shouldFail: true,
		},
		"empty active_versions": {
			input: `active_versions {
} `,
			shouldFail: true,
		},
		"invalid HCL syntax": {
			input:      `active_versions { version "1.19.x" { ce_active = }`,
			shouldFail: true,
		},
		"invalid schema type": {
			input: `schema = "not a number"
active_versions {
  version "1.19.x" {
    ce_active = true
  }
}`,
			shouldFail: true,
		},
		"invalid ce_active type": {
			input: `active_versions {
  version "1.19.x" {
    ce_active = "not a bool"
  }
}`,
			shouldFail: true,
		},
		"invalid lts type": {
			input: `active_versions {
  version "1.19.x" {
    lts = "not a bool"
  }
}`,
			shouldFail: true,
		},
		"unknown attribute in version": {
			input: `active_versions {
  version "1.19.x" {
    ce_active = true
    lts = false
    unknown_field = true
  }
}`,
			shouldFail: true,
		},
		"empty input": {
			input:      "",
			shouldFail: true,
		},
		"only whitespace": {
			input:      "   \n\t  \n  ",
			shouldFail: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result, err := DecodeBytes([]byte(test.input))
			if test.shouldFail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, test.expected.Schema, result.Schema)
			require.NotNil(t, result.ActiveVersion)
			require.EqualValues(t, test.expected.ActiveVersion, result.ActiveVersion)
		})
	}
}

// TestDecodeFile tests the DecodeFile function
func TestDecodeFile(t *testing.T) {
	t.Parallel()
	for name, test := range map[string]struct {
		content    string
		shouldFail bool
	}{
		"valid file": {
			content:    validVersionsHCL,
			shouldFail: false,
		},
		"invalid HCL": {
			content:    `active_versions { invalid }`,
			shouldFail: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create a temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "versions.hcl")
			err := os.WriteFile(tmpFile, []byte(test.content), 0o644)
			require.NoError(t, err)

			result, err := DecodeFile(tmpFile)

			if test.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotNil(t, result.ActiveVersion)
		})
	}
}

// TestDecodeFile_FileErrors tests DecodeFile with file I/O errors
func TestDecodeFile_FileErrors(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		setupFunc  func(t *testing.T) string
		shouldFail bool
	}{
		"non-existent file": {
			setupFunc: func(t *testing.T) string {
				return "/non/existent/path/versions.hcl"
			},
			shouldFail: true,
		},
		"directory instead of file": {
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				return tmpDir
			},
			shouldFail: true,
		},
		"empty file path": {
			setupFunc: func(t *testing.T) string {
				return ""
			},
			shouldFail: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			path := test.setupFunc(t)
			result, err := DecodeFile(path)
			if test.shouldFail {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
		})
	}
}

// TestLoad tests the Load function
func TestLoad(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		content    string
		shouldFail bool
	}{
		"valid config": {
			content:    validVersionsHCL,
			shouldFail: false,
		},
		"invalid HCL": {
			content:    `invalid hcl content`,
			shouldFail: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "versions.hcl")
			err := os.WriteFile(tmpFile, []byte(test.content), 0o644)
			require.NoError(t, err)
			result, err := Load(context.Background(), tmpFile)
			if test.shouldFail {
				require.Error(t, err)
				require.Nil(t, result)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotNil(t, result.ActiveVersion)
		})
	}
}

// TestDecodeRes_Validate tests the Validate method of DecodeRes
func TestDecodeRes_Validate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	for name, test := range map[string]struct {
		decodeRes  *DecodeRes
		shouldFail bool
	}{
		"nil DecodeRes": {
			decodeRes:  nil,
			shouldFail: true,
		},
		"DecodeRes with embedded error": {
			decodeRes: &DecodeRes{
				Path:   "/some/path/versions.hcl",
				Err:    errors.New("failed to decode configuration"),
				ErrStr: "failed to decode configuration",
			},
			shouldFail: true,
		},
		"DecodeRes with file not found error": {
			decodeRes: &DecodeRes{
				Path:   "/non/existent/versions.hcl",
				Err:    errors.New("open /non/existent/versions.hcl: no such file or directory"),
				ErrStr: "open /non/existent/versions.hcl: no such file or directory",
			},
			shouldFail: true,
		},
		"DecodeRes with parse error": {
			decodeRes: &DecodeRes{
				Path:   "/some/path/versions.hcl",
				Err:    errors.New("versions.hcl:1,1-1: Invalid expression; Expected the start of an expression, but found an invalid expression token."),
				ErrStr: "versions.hcl:1,1-1: Invalid expression; Expected the start of an expression, but found an invalid expression token.",
			},
			shouldFail: true,
		},
		"DecodeRes with validation error": {
			decodeRes: &DecodeRes{
				Path:   "/some/path/versions.hcl",
				Err:    errors.New("no active_versions stanza has been defined"),
				ErrStr: "no active_versions stanza has been defined",
			},
			shouldFail: true,
		},
		"DecodeRes with multiple joined errors": {
			decodeRes: &DecodeRes{
				Path:   "/some/path/versions.hcl",
				Err:    errors.Join(errors.New("error 1"), errors.New("error 2")),
				ErrStr: "error 1\nerror 2",
			},
			shouldFail: true,
		},
		"valid DecodeRes with config": {
			decodeRes: &DecodeRes{
				Path: "/some/path/versions.hcl",
				Config: &VersionsConfig{
					Schema: 1,
					ActiveVersion: &ActiveVersion{
						Versions: map[string]*Version{
							"1.19.x": {CEActive: true, LTS: false},
						},
					},
				},
				Err:    nil,
				ErrStr: "",
			},
			shouldFail: false,
		},
		"valid DecodeRes with minimal config": {
			decodeRes: &DecodeRes{
				Path: "/some/path/versions.hcl",
				Config: &VersionsConfig{
					ActiveVersion: &ActiveVersion{
						Versions: map[string]*Version{},
					},
				},
				Err:    nil,
				ErrStr: "",
			},
			shouldFail: false,
		},
		"valid DecodeRes with empty path": {
			decodeRes: &DecodeRes{
				Path: "",
				Config: &VersionsConfig{
					Schema: 1,
					ActiveVersion: &ActiveVersion{
						Versions: map[string]*Version{
							"1.19.x": {CEActive: true, LTS: false},
						},
					},
				},
				Err:    nil,
				ErrStr: "",
			},
			shouldFail: false,
		},
		"valid DecodeRes with nil config but no error": {
			decodeRes: &DecodeRes{
				Path:   "/some/path/versions.hcl",
				Config: nil,
				Err:    nil,
				ErrStr: "",
			},
			shouldFail: false,
		},
		"DecodeRes with error but empty error string": {
			decodeRes: &DecodeRes{
				Path:   "/some/path/versions.hcl",
				Err:    errors.New("some error"),
				ErrStr: "",
			},
			shouldFail: true,
		},
		"DecodeRes with error string but nil error": {
			decodeRes: &DecodeRes{
				Path:   "/some/path/versions.hcl",
				Err:    nil,
				ErrStr: "some error string",
			},
			shouldFail: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := test.decodeRes.Validate(ctx)

			if test.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
