// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package config

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v81/github"
	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/changed"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

// TestConfig_Decode tests decoding of the fixtures/pipeline.hcl and verifies
// that the changed_files grouping matches our pre-pipeline.hcl build-in
// checkers.
func TestConfig_Decode(t *testing.T) {
	t.Parallel()
	cfg, err := DecodeFile("./fixtures/pipeline.hcl")
	require.NoError(t, err)

	// Verify changed file groups based on the changed_files stanza in the pipeline.hcl
	for filename, groups := range map[string]changed.FileGroups{
		".build/entrypoint.sh":                      {changed.FileGroup("pipeline")},
		".github/actions/changed-files/actions.yml": {changed.FileGroup("github"), changed.FileGroup("pipeline")},
		".github/workflows/build.yml":               {changed.FileGroup("github"), changed.FileGroup("pipeline")},
		".github/workflows/build-artifacts-ce.yml":  {changed.FileGroup("community"), changed.FileGroup("github"), changed.FileGroup("pipeline")},
		// NOTE: no "enterprise" for build-artifacts-ent.yml as it's ignored
		".github/workflows/build-artifacts-ent.yml":      {changed.FileGroup("github"), changed.FileGroup("pipeline")},
		".github/workflows/backport-ce-ent.yml":          {changed.FileGroup("community"), changed.FileGroup("github"), changed.FileGroup("pipeline")},
		".github/scripts/pr_comment.sh":                  {changed.FileGroup("github"), changed.FileGroup("pipeline")},
		".github/CODEOWNERS":                             {changed.FileGroup("github")},
		".go-version":                                    {changed.FileGroup("gotoolchain")},
		".hooks/pre-push":                                {changed.FileGroup("pipeline")},
		".release/ibm-pao/eboms/5900-BJ8.essentials.csv": {changed.FileGroup("enterprise"), changed.FileGroup("pipeline")},
		".release/docker/ubi-docker-entrypoint.sh":       {changed.FileGroup("pipeline")},
		"audit/backend_ce.go":                            {changed.FileGroup("app"), changed.FileGroup("community")},
		"audit/backend_config_ent.go":                    {changed.FileGroup("app"), changed.FileGroup("enterprise")},
		"builtin/logical/transit/something_ent.go":       {changed.FileGroup("app"), changed.FileGroup("enterprise")},
		"buf.yml":                                                                {changed.FileGroup("proto")},
		"changelog/1726.txt":                                                     {changed.FileGroup("changelog")},
		"changelog/_1726.txt":                                                    {changed.FileGroup("changelog")},
		"command/server/config.go":                                               {changed.FileGroup("app")},
		"command/operator_raft_autopilot_state.go":                               {changed.FileGroup("app"), changed.FileGroup("autopilot")},
		"command/agent_ent_test.go":                                              {changed.FileGroup("app"), changed.FileGroup("enterprise")},
		"enos/enos-samples-ce-build.hcl":                                         {changed.FileGroup("community"), changed.FileGroup("enos")},
		"enos/enos-samples-ent-build.hcl":                                        {changed.FileGroup("enos"), changed.FileGroup("enterprise")},
		"enos/enos-scenario-smoke.hcl":                                           {changed.FileGroup("enos")},
		"enos/enos-scenario-autopilot-ent.hcl":                                   {changed.FileGroup("enos"), changed.FileGroup("enterprise")},
		"enos/modules/softhsm_create_vault_keys/scripts/create-keys.sh":          {changed.FileGroup("enos")},
		"enos/modules/softhsm_create_vault_keys/scripts/get-keys.sh":             {changed.FileGroup("enos")},
		"enos/modules/softhsm_distribute_vault_keys/main.tf":                     {changed.FileGroup("enos")},
		"enos/modules/softhsm_distribute_vault_keys/scripts/distribute-token.sh": {changed.FileGroup("enos")},
		"enos/modules/softhsm_init/main.tf":                                      {changed.FileGroup("enos")},
		"enos/modules/softhsm_init/scripts/init-softhsm.sh":                      {changed.FileGroup("enos")},
		"enos/modules/softhsm_install/main.tf":                                   {changed.FileGroup("enos")},
		"enos/modules/softhsm_install/scripts/find-shared-object.sh":             {changed.FileGroup("enos")},
		"enos/modules/verify_secrets_engines/scripts/identity-verify-entity.sh":  {changed.FileGroup("enos")},
		"go.mod":                                            {changed.FileGroup("app"), changed.FileGroup("gotoolchain")},
		"go.sum":                                            {changed.FileGroup("app"), changed.FileGroup("gotoolchain")},
		"helper/identity/mfa/types.proto":                   {changed.FileGroup("proto")},
		"http/util_stubs_oss.go":                            {changed.FileGroup("app"), changed.FileGroup("community")},
		"physical/raft/raft_autopilot.go":                   {changed.FileGroup("app"), changed.FileGroup("autopilot")},
		"physical/raft/types.proto":                         {changed.FileGroup("proto")},
		"scripts/ci-helper.sh":                              {changed.FileGroup("pipeline")},
		"scripts/cross/Dockerfile-ent":                      {changed.FileGroup("enterprise"), changed.FileGroup("pipeline")},
		"scripts/cross/Dockerfile-ent-hsm":                  {changed.FileGroup("enterprise"), changed.FileGroup("pipeline")},
		"scripts/dev/hsm/README.md":                         {changed.FileGroup("docs"), changed.FileGroup("enterprise"), changed.FileGroup("pipeline")},
		"scripts/dist-ent.sh":                               {changed.FileGroup("enterprise"), changed.FileGroup("pipeline")},
		"scripts/docker/docker-entrypoint.sh":               {changed.FileGroup("pipeline")},
		"scripts/testing/test-vault-license.sh":             {changed.FileGroup("enterprise"), changed.FileGroup("pipeline")},
		"scripts/testing/upgrade/README.md":                 {changed.FileGroup("docs"), changed.FileGroup("enterprise"), changed.FileGroup("pipeline")},
		"sdk/database/dbplugin/v5/proto/database_ent.pb.go": {changed.FileGroup("app"), changed.FileGroup("enterprise")},
		"sdk/database/dbplugin/v5/proto/database_ent.proto": {changed.FileGroup("enterprise"), changed.FileGroup("proto")},
		"specs/merkle-tree/spec.md":                         {changed.FileGroup("enterprise")},
		"tools/pipeline/main.go":                            {changed.FileGroup("pipeline")},
		"ui/lib/ldap/index.js":                              {changed.FileGroup("ui")},
		"vault/acl.go":                                      {changed.FileGroup("app")},
		"vault/activity_log_util_ent.go":                    {changed.FileGroup("app"), changed.FileGroup("enterprise")},
		"vault/identity_store_ent_test.go":                  {changed.FileGroup("app"), changed.FileGroup("enterprise")},
		"vault_ent/go.mod":                                  {changed.FileGroup("app"), changed.FileGroup("enterprise"), changed.FileGroup("gotoolchain")},
		"vault_ent/go.sum":                                  {changed.FileGroup("app"), changed.FileGroup("enterprise"), changed.FileGroup("gotoolchain")},
		"vault_ent/requires_ent.go":                         {changed.FileGroup("app"), changed.FileGroup("enterprise")},
		"website/content/api-docs/index.mdx":                {changed.FileGroup("docs")},
		"CHANGELOG.md":                                      {changed.FileGroup("changelog")},
		"Dockerfile":                                        {changed.FileGroup("pipeline")},
		"Makefile":                                          {changed.FileGroup("pipeline")},
		"README.md":                                         {changed.FileGroup("docs")},
	} {
		t.Run("file name only: "+filename, func(t *testing.T) {
			t.Parallel()
			file := &changed.File{Filename: filename}
			changed.Group(context.Background(), file, cfg.ChangedFiles.FileGroups)
			require.Equal(t, groups, file.Groups)
		})

		t.Run("github file: "+filename, func(t *testing.T) {
			t.Parallel()
			file := &changed.File{GithubCommitFile: &github.CommitFile{Filename: &filename}}
			changed.Group(context.Background(), file, cfg.ChangedFiles.FileGroups)
			require.Equal(t, groups, file.Groups)
		})
	}
}

func TestDecodeFile_ErrorCases(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		path string
	}{
		{
			name: "file not found",
			path: "./fixtures/nonexistent.hcl",
		},
		{
			name: "empty path",
			path: "",
		},
		{
			name: "directory instead of file",
			path: "./fixtures",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := DecodeFile(tt.path)
			require.Error(t, err)
			require.Nil(t, cfg)
		})
	}
}

func TestDecode_ValidHCL(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name     string
		hcl      string
		validate func(t *testing.T, cfg *Config)
	}{
		{
			name: "empty config",
			hcl:  ``,
			validate: func(t *testing.T, cfg *Config) {
				require.NotNil(t, cfg)
				require.Nil(t, cfg.ChangedFiles)
			},
		},
		{
			name: "config with changed_files block",
			hcl: `
changed_files {
  group "test" {
    match {
      extension = [".go"]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.NotNil(t, cfg)
				require.NotNil(t, cfg.ChangedFiles)
				require.Len(t, cfg.ChangedFiles.Groups, 1)
				require.Equal(t, "test", cfg.ChangedFiles.Groups[0].Name)
			},
		},
		{
			name: "config with multiple groups",
			hcl: `
changed_files {
  group "test1" {
    match {
      extension = [".go"]
    }
  }
  group "test2" {
    match {
      extension = [".md"]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.NotNil(t, cfg)
				require.NotNil(t, cfg.ChangedFiles)
				require.Len(t, cfg.ChangedFiles.Groups, 2)
				require.Equal(t, "test1", cfg.ChangedFiles.Groups[0].Name)
				require.Equal(t, "test2", cfg.ChangedFiles.Groups[1].Name)
			},
		},
		{
			name: "config with joinpath function in file attribute",
			hcl: `
changed_files {
  group "test" {
    match {
      file = [joinpath("src", "main.go")]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.NotNil(t, cfg)
				require.NotNil(t, cfg.ChangedFiles)
				require.Len(t, cfg.ChangedFiles.Groups, 1)
				require.Equal(t, "test", cfg.ChangedFiles.Groups[0].Name)
				// joinpath should use OS path separator, then convert to forward slash
				expectedPath := filepath.ToSlash(filepath.Join("src", "main.go"))
				require.Equal(t, []string{expectedPath}, cfg.ChangedFiles.Groups[0].Match[0].File)
			},
		},
		{
			name: "config with multiple match blocks",
			hcl: `
changed_files {
  group "test" {
    match {
      extension = [".go"]
    }
    match {
      extension = [".md"]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.NotNil(t, cfg)
				require.NotNil(t, cfg.ChangedFiles)
				require.Len(t, cfg.ChangedFiles.Groups, 1)
				require.Len(t, cfg.ChangedFiles.Groups[0].Match, 2)
				require.Equal(t, []string{".go"}, cfg.ChangedFiles.Groups[0].Match[0].Extension)
				require.Equal(t, []string{".md"}, cfg.ChangedFiles.Groups[0].Match[1].Extension)
			},
		},
		{
			name: "config with match having multiple fields",
			hcl: `
changed_files {
  group "test" {
    match {
      extension = [".go"]
      base_dir = ["src"]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.NotNil(t, cfg)
				require.NotNil(t, cfg.ChangedFiles)
				require.Len(t, cfg.ChangedFiles.Groups, 1)
				require.Len(t, cfg.ChangedFiles.Groups[0].Match, 1)
				require.Equal(t, []string{".go"}, cfg.ChangedFiles.Groups[0].Match[0].Extension)
				require.Equal(t, []string{"src"}, cfg.ChangedFiles.Groups[0].Match[0].BaseDir)
			},
		},
		{
			name: "config with all match fields set",
			hcl: `
changed_files {
  group "test" {
    match {
      file = ["src/main.go"]
      base_dir = ["src"]
      base_name = ["main.go"]
      base_name_prefix = ["main"]
      contains = ["main"]
      extension = [".go"]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.NotNil(t, cfg)
				require.NotNil(t, cfg.ChangedFiles)
				require.Len(t, cfg.ChangedFiles.Groups, 1)
				matcher := cfg.ChangedFiles.Groups[0].Match[0]
				require.Equal(t, []string{"src/main.go"}, matcher.File)
				require.Equal(t, []string{"src"}, matcher.BaseDir)
				require.Equal(t, []string{"main.go"}, matcher.BaseName)
				require.Equal(t, []string{"main"}, matcher.BaseNamePrefix)
				require.Equal(t, []string{"main"}, matcher.Contains)
				require.Equal(t, []string{".go"}, matcher.Extension)
			},
		},
		{
			name: "config with ignore block",
			hcl: `
changed_files {
  group "test" {
    ignore {
      extension = [".test"]
    }
    match {
      extension = [".go"]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.NotNil(t, cfg)
				require.NotNil(t, cfg.ChangedFiles)
				require.Len(t, cfg.ChangedFiles.Groups, 1)
				require.Len(t, cfg.ChangedFiles.Groups[0].Ignore, 1)
				require.Equal(t, []string{".test"}, cfg.ChangedFiles.Groups[0].Ignore[0].Extension)
			},
		},
		{
			name: "config with multiple ignore blocks",
			hcl: `
changed_files {
  group "test" {
    ignore {
      extension = [".test"]
    }
    ignore {
      base_name = ["test.go"]
    }
    match {
      extension = [".go"]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.NotNil(t, cfg)
				require.NotNil(t, cfg.ChangedFiles)
				require.Len(t, cfg.ChangedFiles.Groups, 1)
				require.Len(t, cfg.ChangedFiles.Groups[0].Ignore, 2)
				require.Equal(t, []string{".test"}, cfg.ChangedFiles.Groups[0].Ignore[0].Extension)
				require.Equal(t, []string{"test.go"}, cfg.ChangedFiles.Groups[0].Ignore[1].BaseName)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := DecodeBytes([]byte(tt.hcl))
			require.NoError(t, err)
			tt.validate(t, cfg)
		})
	}
}

func TestDecode_InvalidHCL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		hcl         string
		expectError bool
		errorMsg    string
	}{
		{
			name: "invalid HCL syntax - missing closing bracket",
			hcl: `
changed_files {
  group "test" {
    match {
      extension = [".go"
    }
  }
}
`,
			expectError: true,
			errorMsg:    "Missing item separator",
		},
		{
			name: "invalid block structure",
			hcl: `
invalid_block {
  something = "value"
}
`,
			expectError: true,
			errorMsg:    "Unsupported block type",
		},
		{
			name: "invalid attribute",
			hcl: `
changed_files {
  invalid_attr = "value"
}
`,
			expectError: true,
			errorMsg:    "Unsupported argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := DecodeBytes([]byte(tt.hcl))
			if tt.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
				require.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)
			}
		})
	}
}

func TestJoinPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		args     []cty.Value
		expected string
	}{
		{
			name:     "single path",
			args:     []cty.Value{cty.StringVal("test")},
			expected: "test",
		},
		{
			name:     "two paths",
			args:     []cty.Value{cty.StringVal("src"), cty.StringVal("main.go")},
			expected: filepath.ToSlash(filepath.Join("src", "main.go")),
		},
		{
			name:     "three paths",
			args:     []cty.Value{cty.StringVal("src"), cty.StringVal("pkg"), cty.StringVal("config.go")},
			expected: filepath.ToSlash(filepath.Join("src", "pkg", "config.go")),
		},
		{
			name:     "empty path",
			args:     []cty.Value{cty.StringVal("")},
			expected: "",
		},
		{
			name:     "paths with separators",
			args:     []cty.Value{cty.StringVal("src/pkg"), cty.StringVal("config.go")},
			expected: filepath.ToSlash(filepath.Join("src/pkg", "config.go")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := joinPath.Call(tt.args)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result.AsString())
		})
	}
}

func TestJoinPath_UnknownValue(t *testing.T) {
	t.Parallel()

	// Test with unknown value
	args := []cty.Value{cty.UnknownVal(cty.String), cty.StringVal("test")}
	result, err := joinPath.Call(args)
	require.NoError(t, err)
	require.False(t, result.IsKnown())
	require.Equal(t, cty.String, result.Type())
}

func TestJoinPath_NoArgs(t *testing.T) {
	t.Parallel()

	// Test with no arguments
	args := []cty.Value{}
	result, err := joinPath.Call(args)
	require.NoError(t, err)
	require.Equal(t, "", result.AsString())
}

func TestDecodeFile_WithTempFile(t *testing.T) {
	t.Parallel()

	// Create a temporary file with valid HCL
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.hcl")
	content := `
changed_files {
  group "test" {
    match {
      extension = [".go"]
    }
  }
}
`
	err := os.WriteFile(tmpFile, []byte(content), 0o644)
	require.NoError(t, err)

	// Test DecodeFile with the temporary file
	cfg, err := DecodeFile(tmpFile)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.NotNil(t, cfg.ChangedFiles)
	require.Len(t, cfg.ChangedFiles.Groups, 1)
	require.Equal(t, "test", cfg.ChangedFiles.Groups[0].Name)
}

func TestDecodeFile_InvalidHCLFile(t *testing.T) {
	t.Parallel()

	// Create a temporary file with invalid HCL
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.hcl")
	content := `
changed_files {
  group "test" {
    match {
      extension = [".go"
    }
  }
}
`
	err := os.WriteFile(tmpFile, []byte(content), 0o644)
	require.NoError(t, err)

	// Test DecodeFile with the invalid file
	cfg, err := DecodeFile(tmpFile)
	require.Error(t, err)
	require.Nil(t, cfg)
	require.Contains(t, err.Error(), "Missing item separator")
}

func TestDecode_WithStdlibFunctions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		hcl      string
		validate func(t *testing.T, cfg *Config)
	}{
		{
			name: "concat function in extension",
			hcl: `
changed_files {
  group "test" {
    match {
      extension = concat([".go"], [".md"])
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, []string{".go", ".md"}, cfg.ChangedFiles.Groups[0].Match[0].Extension)
			},
		},
		{
			name: "upper function in file path",
			hcl: `
changed_files {
  group "test" {
    match {
      file = [upper("readme.md")]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, []string{"README.MD"}, cfg.ChangedFiles.Groups[0].Match[0].File)
			},
		},
		{
			name: "lower function in base_name",
			hcl: `
changed_files {
  group "test" {
    match {
      base_name = [lower("README.MD")]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, []string{"readme.md"}, cfg.ChangedFiles.Groups[0].Match[0].BaseName)
			},
		},
		{
			name: "format function in contains",
			hcl: `
changed_files {
  group "test" {
    match {
      contains = [format("test-%s", "name")]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, []string{"test-name"}, cfg.ChangedFiles.Groups[0].Match[0].Contains)
			},
		},
		{
			name: "join function in base_name_prefix",
			hcl: `
changed_files {
  group "test" {
    match {
      base_name_prefix = [join("-", ["test", "name"])]
    }
  }
}
`,
			validate: func(t *testing.T, cfg *Config) {
				require.Equal(t, []string{"test-name"}, cfg.ChangedFiles.Groups[0].Match[0].BaseNamePrefix)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cfg, err := DecodeBytes([]byte(tt.hcl))
			require.NoError(t, err)
			require.NotNil(t, cfg)
			require.NotNil(t, cfg.ChangedFiles)
			tt.validate(t, cfg)
		})
	}
}

func TestEvalContext(t *testing.T) {
	t.Parallel()

	// Test that evalContext returns a valid context with expected functions
	ctx := evalContext()
	require.NotNil(t, ctx)
	require.NotNil(t, ctx.Functions)

	// Test a sample of expected functions
	expectedFunctions := []string{
		"abs", "absolute", "add", "and", "upper", "lower",
		"concat", "join", "joinpath", "format", "split",
	}

	for _, funcName := range expectedFunctions {
		t.Run("function_"+funcName, func(t *testing.T) {
			_, exists := ctx.Functions[funcName]
			require.True(t, exists, "Expected function %s to exist in evalContext", funcName)
		})
	}
}

func TestDecode_NilConfig(t *testing.T) {
	t.Parallel()

	// Test that decoding empty bytes returns a valid config with nil ChangedFiles
	cfg, err := DecodeBytes([]byte(""))
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Nil(t, cfg.ChangedFiles)
}

func TestDecode_MultipleChangedFilesBlocks(t *testing.T) {
	t.Parallel()

	// HCL should only allow one changed_files block
	hcl := `
changed_files {
  group "test1" {
    match {
      extension = [".go"]
    }
  }
}
changed_files {
  group "test2" {
    match {
      extension = [".md"]
    }
  }
}
`
	cfg, err := DecodeBytes([]byte(hcl))
	// This should error because only one changed_files block is allowed
	require.Error(t, err)
	require.Nil(t, cfg)
}

// TestDecodeRes_Validate tests the Validate method of DecodeRes with various scenarios.
func TestDecodeRes_Validate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		decodeRes  *DecodeRes
		shouldFail bool
	}{
		{
			name:       "nil",
			decodeRes:  nil,
			shouldFail: true,
		},
		{
			name: "with Err set",
			decodeRes: &DecodeRes{
				Err:    errors.New("decode failed"),
				ErrStr: "decode failed",
			},
			shouldFail: true,
		},
		{
			name: "with Err set and no ErrStr",
			decodeRes: &DecodeRes{
				Err: errors.New("some error"),
			},
			shouldFail: true,
		},
		{
			name: "valid with no error",
			decodeRes: &DecodeRes{
				Path:   "/path/to/pipeline.hcl",
				Config: &Config{},
			},
			shouldFail: false,
		},
		{
			name: "valid with nil Config but no error",
			decodeRes: &DecodeRes{
				Path: "/path/to/pipeline.hcl",
			},
			shouldFail: false,
		},
		{
			name: "with ErrStr but no Err (edge case)",
			decodeRes: &DecodeRes{
				Path:   "/path/to/pipeline.hcl",
				ErrStr: "some error string",
			},
			shouldFail: false,
		},
		{
			name: "with empty Path and no error",
			decodeRes: &DecodeRes{
				Config: &Config{},
			},
			shouldFail: false,
		},
		{
			name: "with all fields set and no error",
			decodeRes: &DecodeRes{
				Path: "/path/to/pipeline.hcl",
				Config: &Config{
					ChangedFiles: nil,
				},
			},
			shouldFail: false,
		},
		{
			name: "with file not found error",
			decodeRes: &DecodeRes{
				Path:   "/nonexistent/pipeline.hcl",
				Err:    errors.New("open /nonexistent/pipeline.hcl: no such file or directory"),
				ErrStr: "open /nonexistent/pipeline.hcl: no such file or directory",
			},
			shouldFail: true,
		},
		{
			name: "with HCL parse error",
			decodeRes: &DecodeRes{
				Path:   "/path/to/invalid.hcl",
				Err:    errors.New("Missing item separator"),
				ErrStr: "Missing item separator",
			},
			shouldFail: true,
		},
		{
			name: "with wrapped error",
			decodeRes: &DecodeRes{
				Err:    errors.Join(errors.New("error 1"), errors.New("error 2")),
				ErrStr: "error 1\nerror 2",
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.decodeRes.Validate(ctx)

			if tt.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
