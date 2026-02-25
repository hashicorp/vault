// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package changed

import (
	"context"
	"path/filepath"
	"testing"

	gh "github.com/google/go-github/v81/github"
	"github.com/stretchr/testify/require"
)

// TestConfig_FileGroups tests the Config.FileGroups method with various scenarios
func TestConfig_FileGroups(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config   *Config
		file     *File
		expected FileGroups
	}{
		"nil config": {
			config:   nil,
			file:     &File{Filename: "test.go"},
			expected: nil,
		},
		"empty groups": {
			config:   &Config{Groups: []*GroupConfig{}},
			file:     &File{Filename: "test.go"},
			expected: nil,
		},
		"file with empty filename": {
			config: &Config{
				Groups: []*GroupConfig{
					{
						Name: "go",
						Match: Matchers{
							{Extension: []string{".go"}},
						},
					},
				},
			},
			file:     &File{Filename: ""},
			expected: nil,
		},
		"single group match": {
			config: &Config{
				Groups: []*GroupConfig{
					{
						Name: "go",
						Match: Matchers{
							{Extension: []string{".go"}},
						},
					},
				},
			},
			file:     &File{Filename: "main.go"},
			expected: FileGroups{"go"},
		},
		"multiple groups match": {
			config: &Config{
				Groups: []*GroupConfig{
					{
						Name: "go",
						Match: Matchers{
							{Extension: []string{".go"}},
						},
					},
					{
						Name: "source",
						Match: Matchers{
							{BaseDir: []string{"src"}},
						},
					},
				},
			},
			file:     &File{Filename: "src/main.go"},
			expected: FileGroups{"go", "source"},
		},
		"ignore takes precedence over match": {
			config: &Config{
				Groups: []*GroupConfig{
					{
						Name: "go",
						Ignore: Matchers{
							{BaseName: []string{"main.go"}},
						},
						Match: Matchers{
							{Extension: []string{".go"}},
						},
					},
				},
			},
			file:     &File{Filename: "main.go"},
			expected: FileGroups{},
		},
		"ignore does not match, match succeeds": {
			config: &Config{
				Groups: []*GroupConfig{
					{
						Name: "go",
						Ignore: Matchers{
							{BaseName: []string{"test.go"}},
						},
						Match: Matchers{
							{Extension: []string{".go"}},
						},
					},
				},
			},
			file:     &File{Filename: "main.go"},
			expected: FileGroups{"go"},
		},
		"github commit file takes precedence": {
			config: &Config{
				Groups: []*GroupConfig{
					{
						Name: "go",
						Match: Matchers{
							{Extension: []string{".go"}},
						},
					},
				},
			},
			file: &File{
				Filename:         "wrong.txt",
				GithubCommitFile: &gh.CommitFile{Filename: gh.Ptr("correct.go")},
			},
			expected: FileGroups{"go"},
		},
		"complex multi-group scenario": {
			config: &Config{
				Groups: []*GroupConfig{
					{
						Name: "backend",
						Match: Matchers{
							{BaseDir: []string{"backend"}},
						},
					},
					{
						Name: "go",
						Match: Matchers{
							{Extension: []string{".go"}},
						},
					},
					{
						Name: "tests",
						Match: Matchers{
							{BaseNamePrefix: []string{"test_"}},
						},
					},
					{
						Name: "ignored",
						Ignore: Matchers{
							{Contains: []string{"vendor"}},
						},
						Match: Matchers{
							{Extension: []string{".go"}},
						},
					},
				},
			},
			file:     &File{Filename: "backend/test_handler.go"},
			expected: FileGroups{"backend", "go", "ignored", "tests"},
		},
		"no match in any group": {
			config: &Config{
				Groups: []*GroupConfig{
					{
						Name: "go",
						Match: Matchers{
							{Extension: []string{".go"}},
						},
					},
					{
						Name: "python",
						Match: Matchers{
							{Extension: []string{".py"}},
						},
					},
				},
			},
			file:     &File{Filename: "readme.md"},
			expected: FileGroups{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			result := test.config.FileGroups(ctx, test.file)
			require.Len(t, result, len(test.expected))
			require.Equal(t, test.expected, result)
		})
	}
}

// TestMatchers_Match tests the Matchers.Match method
func TestMatchers_Match(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		matchers Matchers
		path     string
		expected bool
	}{
		"nil matchers": {
			matchers: nil,
			path:     "test.go",
			expected: false,
		},
		"empty matchers": {
			matchers: Matchers{},
			path:     "test.go",
			expected: false,
		},
		"single matcher matches": {
			matchers: Matchers{
				{Extension: []string{".go"}},
			},
			path:     "test.go",
			expected: true,
		},
		"single matcher does not match": {
			matchers: Matchers{
				{Extension: []string{".py"}},
			},
			path:     "test.go",
			expected: false,
		},
		"multiple matchers, first matches": {
			matchers: Matchers{
				{Extension: []string{".go"}},
				{Extension: []string{".py"}},
			},
			path:     "test.go",
			expected: true,
		},
		"multiple matchers, second matches": {
			matchers: Matchers{
				{Extension: []string{".py"}},
				{Extension: []string{".go"}},
			},
			path:     "test.go",
			expected: true,
		},
		"multiple matchers, none match": {
			matchers: Matchers{
				{Extension: []string{".py"}},
				{Extension: []string{".js"}},
			},
			path:     "test.go",
			expected: false,
		},
		"matcher with nil config": {
			matchers: Matchers{
				nil,
				{Extension: []string{".go"}},
			},
			path:     "test.go",
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.expected, test.matchers.Match(test.path))
		})
	}
}

// TestMatchConfig_Match tests the MatchConfig.Match method with individual fields
func TestMatchConfig_Match(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config   *MatchConfig
		path     string
		expected bool
	}{
		"nil config": {
			config:   nil,
			path:     "test.go",
			expected: false,
		},
		"empty config": {
			config:   &MatchConfig{},
			path:     "test.go",
			expected: false,
		},
		// File field tests
		"file matches exactly": {
			config: &MatchConfig{
				File: []string{"test.go"},
			},
			path:     "test.go",
			expected: true,
		},
		"file does not match": {
			config: &MatchConfig{
				File: []string{"other.go"},
			},
			path:     "test.go",
			expected: false,
		},
		"file matches one of multiple": {
			config: &MatchConfig{
				File: []string{"other.go", "test.go", "another.go"},
			},
			path:     "test.go",
			expected: true,
		},
		"file with path matches": {
			config: &MatchConfig{
				File: []string{"src/test.go"},
			},
			path:     "src/test.go",
			expected: true,
		},
		// BaseDir field tests
		"base_dir matches": {
			config: &MatchConfig{
				BaseDir: []string{"src"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: true,
		},
		"base_dir does not match": {
			config: &MatchConfig{
				BaseDir: []string{"lib"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: false,
		},
		"base_dir matches nested path": {
			config: &MatchConfig{
				BaseDir: []string{"src"},
			},
			path:     filepath.Join("src", "pkg", "test.go"),
			expected: true,
		},
		"base_dir matches one of multiple": {
			config: &MatchConfig{
				BaseDir: []string{"lib", "src", "pkg"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: true,
		},
		"base_dir does not match without separator": {
			config: &MatchConfig{
				BaseDir: []string{"src"},
			},
			path:     "srctest.go",
			expected: false,
		},
		// BaseName field tests
		"base_name matches": {
			config: &MatchConfig{
				BaseName: []string{"test.go"},
			},
			path:     "test.go",
			expected: true,
		},
		"base_name matches with path": {
			config: &MatchConfig{
				BaseName: []string{"test.go"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: true,
		},
		"base_name does not match": {
			config: &MatchConfig{
				BaseName: []string{"other.go"},
			},
			path:     "test.go",
			expected: false,
		},
		"base_name matches one of multiple": {
			config: &MatchConfig{
				BaseName: []string{"other.go", "test.go", "another.go"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: true,
		},
		// BaseNamePrefix field tests
		"base_name_prefix matches": {
			config: &MatchConfig{
				BaseNamePrefix: []string{"test"},
			},
			path:     "test.go",
			expected: true,
		},
		"base_name_prefix matches with path": {
			config: &MatchConfig{
				BaseNamePrefix: []string{"test"},
			},
			path:     filepath.Join("src", "test_handler.go"),
			expected: true,
		},
		"base_name_prefix does not match": {
			config: &MatchConfig{
				BaseNamePrefix: []string{"prod"},
			},
			path:     "test.go",
			expected: false,
		},
		"base_name_prefix matches one of multiple": {
			config: &MatchConfig{
				BaseNamePrefix: []string{"prod", "test", "dev"},
			},
			path:     filepath.Join("src", "test_handler.go"),
			expected: true,
		},
		"base_name_prefix matches substring": {
			config: &MatchConfig{
				BaseNamePrefix: []string{"hand"},
			},
			path:     "handler.go",
			expected: true,
		},
		// Contains field tests
		"contains matches": {
			config: &MatchConfig{
				Contains: []string{"test"},
			},
			path:     "src/test/handler.go",
			expected: true,
		},
		"contains does not match": {
			config: &MatchConfig{
				Contains: []string{"prod"},
			},
			path:     "src/test/handler.go",
			expected: false,
		},
		"contains matches one of multiple": {
			config: &MatchConfig{
				Contains: []string{"prod", "test", "dev"},
			},
			path:     "src/test/handler.go",
			expected: true,
		},
		"contains matches in filename": {
			config: &MatchConfig{
				Contains: []string{"handler"},
			},
			path:     "test_handler.go",
			expected: true,
		},
		// Extension field tests
		"extension matches": {
			config: &MatchConfig{
				Extension: []string{".go"},
			},
			path:     "test.go",
			expected: true,
		},
		"extension matches with path": {
			config: &MatchConfig{
				Extension: []string{".go"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: true,
		},
		"extension does not match": {
			config: &MatchConfig{
				Extension: []string{".py"},
			},
			path:     "test.go",
			expected: false,
		},
		"extension matches one of multiple": {
			config: &MatchConfig{
				Extension: []string{".py", ".go", ".js"},
			},
			path:     "test.go",
			expected: true,
		},
		"extension with no extension in path": {
			config: &MatchConfig{
				Extension: []string{".go"},
			},
			path:     "Makefile",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.expected, test.config.Match(test.path))
		})
	}
}

// TestMatchConfig_Match_MultipleFields tests MatchConfig.Match with multiple fields set
func TestMatchConfig_Match_MultipleFields(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config   *MatchConfig
		path     string
		expected bool
	}{
		"file and extension both match": {
			config: &MatchConfig{
				File:      []string{"test.go"},
				Extension: []string{".go"},
			},
			path:     "test.go",
			expected: true,
		},
		"file matches but extension does not": {
			config: &MatchConfig{
				File:      []string{"test.go"},
				Extension: []string{".py"},
			},
			path:     "test.go",
			expected: false,
		},
		"extension matches but file does not": {
			config: &MatchConfig{
				File:      []string{"other.go"},
				Extension: []string{".go"},
			},
			path:     "test.go",
			expected: false,
		},
		"base_dir and extension both match": {
			config: &MatchConfig{
				BaseDir:   []string{"src"},
				Extension: []string{".go"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: true,
		},
		"base_dir matches but extension does not": {
			config: &MatchConfig{
				BaseDir:   []string{"src"},
				Extension: []string{".py"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: false,
		},
		"base_name and extension both match": {
			config: &MatchConfig{
				BaseName:  []string{"test.go"},
				Extension: []string{".go"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: true,
		},
		"base_name matches but extension does not": {
			config: &MatchConfig{
				BaseName:  []string{"test.go"},
				Extension: []string{".py"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: false,
		},
		"base_name_prefix and extension both match": {
			config: &MatchConfig{
				BaseNamePrefix: []string{"test"},
				Extension:      []string{".go"},
			},
			path:     "test_handler.go",
			expected: true,
		},
		"base_name_prefix matches but extension does not": {
			config: &MatchConfig{
				BaseNamePrefix: []string{"test"},
				Extension:      []string{".py"},
			},
			path:     "test_handler.go",
			expected: false,
		},
		"contains and extension both match": {
			config: &MatchConfig{
				Contains:  []string{"handler"},
				Extension: []string{".go"},
			},
			path:     filepath.Join("src", "test_handler.go"),
			expected: true,
		},
		"contains matches but extension does not": {
			config: &MatchConfig{
				Contains:  []string{"handler"},
				Extension: []string{".py"},
			},
			path:     filepath.Join("src", "test_handler.go"),
			expected: false,
		},
		"base_dir, base_name, and extension all match": {
			config: &MatchConfig{
				BaseDir:   []string{"src"},
				BaseName:  []string{"test.go"},
				Extension: []string{".go"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: true,
		},
		"base_dir and base_name match but extension does not": {
			config: &MatchConfig{
				BaseDir:   []string{"src"},
				BaseName:  []string{"test.go"},
				Extension: []string{".py"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: false,
		},
		"all fields match": {
			config: &MatchConfig{
				File:           []string{filepath.Join("src", "test_handler.go")},
				BaseDir:        []string{"src"},
				BaseName:       []string{"test_handler.go"},
				BaseNamePrefix: []string{"test"},
				Contains:       []string{"handler"},
				Extension:      []string{".go"},
			},
			path:     filepath.Join("src", "test_handler.go"),
			expected: true,
		},
		"all fields except one match": {
			config: &MatchConfig{
				File:           []string{filepath.Join("src", "test_handler.go")},
				BaseDir:        []string{"src"},
				BaseName:       []string{"test_handler.go"},
				BaseNamePrefix: []string{"test"},
				Contains:       []string{"handler"},
				Extension:      []string{".py"}, // This doesn't match
			},
			path:     filepath.Join("src", "test_handler.go"),
			expected: false,
		},
		"base_dir and contains both match": {
			config: &MatchConfig{
				BaseDir:  []string{"src"},
				Contains: []string{"test"},
			},
			path:     filepath.Join("src", "test", "handler.go"),
			expected: true,
		},
		"base_dir matches but contains does not": {
			config: &MatchConfig{
				BaseDir:  []string{"src"},
				Contains: []string{"prod"},
			},
			path:     filepath.Join("src", "test", "handler.go"),
			expected: false,
		},
		"base_name_prefix and contains both match": {
			config: &MatchConfig{
				BaseNamePrefix: []string{"test"},
				Contains:       []string{"handler"},
			},
			path:     filepath.Join("src", "test_handler.go"),
			expected: true,
		},
		"file and base_dir both match": {
			config: &MatchConfig{
				File:    []string{filepath.Join("src", "pkg", "test.go")},
				BaseDir: []string{"src"},
			},
			path:     filepath.Join("src", "pkg", "test.go"),
			expected: true,
		},
		"file matches but base_dir does not": {
			config: &MatchConfig{
				File:    []string{filepath.Join("src", "pkg", "test.go")},
				BaseDir: []string{"lib"},
			},
			path:     filepath.Join("src", "pkg", "test.go"),
			expected: false,
		},
		"multiple values in each field, all match": {
			config: &MatchConfig{
				BaseDir:        []string{"lib", "src", "pkg"},
				BaseName:       []string{"main.go", "test.go", "handler.go"},
				BaseNamePrefix: []string{"main", "test", "handler"},
				Contains:       []string{"src", "test"},
				Extension:      []string{".py", ".go", ".js"},
			},
			path:     filepath.Join("src", "test.go"),
			expected: true,
		},
		"multiple values in each field, one field does not match": {
			config: &MatchConfig{
				BaseDir:        []string{"lib", "src", "pkg"},
				BaseName:       []string{"main.go", "test.go", "handler.go"},
				BaseNamePrefix: []string{"main", "test", "handler"},
				Contains:       []string{"src", "test"},
				Extension:      []string{".py", ".js"}, // .go is not in the list
			},
			path:     filepath.Join("src", "test.go"),
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.expected, test.config.Match(test.path))
		})
	}
}

// TestMatchConfig_Match_EdgeCases tests edge cases for MatchConfig.Match
func TestMatchConfig_Match_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		config   *MatchConfig
		path     string
		expected bool
	}{
		"empty string in file list": {
			config: &MatchConfig{
				File: []string{""},
			},
			path:     "",
			expected: true,
		},
		"empty string in base_dir list": {
			config: &MatchConfig{
				BaseDir: []string{""},
			},
			path:     "test.go",
			expected: false,
		},
		"empty string in base_name list": {
			config: &MatchConfig{
				BaseName: []string{""},
			},
			path:     "",
			expected: false,
		},
		"empty string in base_name_prefix list": {
			config: &MatchConfig{
				BaseNamePrefix: []string{""},
			},
			path:     "test.go",
			expected: true, // Empty string is contained in any string
		},
		"empty string in contains list": {
			config: &MatchConfig{
				Contains: []string{""},
			},
			path:     "test.go",
			expected: true, // Empty string is contained in any string
		},
		"empty string in extension list": {
			config: &MatchConfig{
				Extension: []string{""},
			},
			path:     "Makefile",
			expected: true, // Files without extension have empty extension
		},
		"path with special characters": {
			config: &MatchConfig{
				Contains: []string{"test-file"},
			},
			path:     "src/test-file_v2.go",
			expected: true,
		},
		"path with spaces": {
			config: &MatchConfig{
				Contains: []string{"my file"},
			},
			path:     "src/my file.go",
			expected: true,
		},
		"case sensitive file match": {
			config: &MatchConfig{
				File: []string{"Test.go"},
			},
			path:     "test.go",
			expected: false,
		},
		"case sensitive extension match": {
			config: &MatchConfig{
				Extension: []string{".GO"},
			},
			path:     "test.go",
			expected: false,
		},
		"extension with multiple dots": {
			config: &MatchConfig{
				Extension: []string{".gz"},
			},
			path:     "archive.tar.gz",
			expected: true,
		},
		"base_dir with trailing separator": {
			config: &MatchConfig{
				BaseDir: []string{"src" + string(filepath.Separator)},
			},
			path:     filepath.Join("src", "test.go"),
			expected: false, // hasBaseDir expects no trailing separator
		},
		"deeply nested path": {
			config: &MatchConfig{
				BaseDir: []string{"a"},
			},
			path:     filepath.Join("a", "b", "c", "d", "e", "f", "test.go"),
			expected: true,
		},
		"single character matches": {
			config: &MatchConfig{
				BaseNamePrefix: []string{"t"},
				Contains:       []string{"t"},
				Extension:      []string{".t"},
			},
			path:     "test.t",
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.expected, test.config.Match(test.path))
		})
	}
}
