// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package generate

import (
	"context"
	"os"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/hashicorp/vault/tools/pipeline/internal/pkg/releases"
	"github.com/stretchr/testify/require"
)

// Test_GenerateTemplateReq_Run tests running the Test_GenerateTemplateReq's
// Run with mocked version data.
func Test_GenerateTemplateReq_Run(t *testing.T) {
	t.Parallel()

	mockVersions := []string{
		"1.16.10+ent",
		"1.17.6+ent",
		"1.17.5+ent",
		"1.16.9+ent",
		"1.17.4+ent",
		"1.16.8+ent",
	}

	for name, test := range map[string]struct {
		req        *GenerateTemplateReq
		template   string
		expected   string
		shouldFail bool
	}{
		"template with version context": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "1.18.0",
				VersionLister: releases.NewMockClient(mockVersions),
			},
			template: `Current Version: {{ .Version }}`,
			expected: `Current Version: 1.18.0`,
		},
		"template without version context": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "",
				VersionLister: releases.NewMockClient(mockVersions),
			},
			template: `{{- if .Version }}Version: {{ .Version }}{{ else }}No version{{ end }}`,
			expected: `No version`,
		},
		"template using version in function": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "1.18.0",
				VersionLister: releases.NewMockClient(mockVersions),
			},
			template: `{{- $versions := VersionsNMinus "vault" "enterprise" .Version 2 "minor" }}{{ range $versions }}{{ . }} {{ end }}`,
			expected: `1.16.8 1.16.9 1.16.10 1.17.4 1.17.5 1.17.6`,
		},
		"template with GeneratedAt": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "",
				VersionLister: releases.NewMockClient(mockVersions),
			},
			template: `Generated: {{ if .GeneratedAt.IsZero }}zero{{ else }}{{ .GeneratedAt.Format "2006" }}{{ end }}`,
			expected: `Generated: zero`,
		},
		"template with version utilities": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "1.18.0",
				VersionLister: releases.NewMockClient(mockVersions),
			},
			template: `{{- $v := ParseVersion .Version }}{{ $v.Major }}.{{ $v.Minor }}`,
			expected: `1.18`,
		},
		"template with version comparison": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "1.18.0",
				VersionLister: releases.NewMockClient(mockVersions),
			},
			template: `{{- $cmp := CompareVersions .Version "1.17.0" }}{{ if eq $cmp 1 }}newer{{ else }}older{{ end }}`,
			expected: `newer`,
		},
		"template with bounded versions": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "",
				VersionLister: releases.NewMockClient(mockVersions),
			},
			template: `{{- $versions := VersionsBounded "vault" "enterprise" "1.17.6" "1.17.4" "minor" }}{{ range $versions }}{{ . }} {{ end }}`,
			expected: `1.17.4 1.17.5 1.17.6`,
		},
		"template with filter versions": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "",
				VersionLister: releases.NewMockClient(mockVersions),
			},
			template: `{{- $all := VersionsNMinus "vault" "enterprise" "1.18.0" 5 "minor" }}{{- $filtered := FilterVersions $all "~1.17.0" }}{{ range $filtered }}{{ . }} {{ end }}`,
			expected: `1.17.4 1.17.5 1.17.6`,
		},
		"complex template with version": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "1.18.0",
				VersionLister: releases.NewMockClient(mockVersions),
			},
			template: `# Version: {{ .Version }}
{{- $versions := VersionsNMinus "vault" "enterprise" .Version 3 "minor" }}
Versions:
{{- range $versions }}
  - {{ . }}
{{- end }}`,
			expected: `# Version: 1.18.0
Versions:
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.4
  - 1.17.5
  - 1.17.6`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Create a temporary template by writing to stdin simulation
			// For testing, we'll directly set the template content
			tc := &TemplateContext{
				VersionLister: test.req.VersionLister,
				Version:       test.req.Version,
			}

			// Parse and render template
			rendered, err := renderTemplateString(t.Context(), tc, test.template)
			if test.shouldFail {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expected, strings.TrimSpace(string(rendered)))
		})
	}
}

// Test_GenerateTemplateReq_Validate tests running the Test_GenerateTemplateReq's
// Validate with mocked version data.
func Test_GenerateTemplateReq_Validate(t *testing.T) {
	t.Parallel()

	for name, test := range map[string]struct {
		req        *GenerateTemplateReq
		shouldFail bool
	}{
		"valid request with version": {
			req: &GenerateTemplateReq{
				TemplatePath:  "test.tmpl",
				OutputPath:    "output.txt",
				Version:       "1.18.0",
				VersionLister: releases.NewMockClient([]string{}),
			},
		},
		"valid request without version": {
			req: &GenerateTemplateReq{
				TemplatePath:  "test.tmpl",
				OutputPath:    "output.txt",
				Version:       "",
				VersionLister: releases.NewMockClient([]string{}),
			},
		},
		"valid request with stdin": {
			req: &GenerateTemplateReq{
				TemplatePath:  "-",
				OutputPath:    "",
				Version:       "1.18.0",
				VersionLister: releases.NewMockClient([]string{}),
			},
		},
		"valid request with stdout": {
			req: &GenerateTemplateReq{
				TemplatePath:  "test.tmpl",
				OutputPath:    "",
				Version:       "1.18.0",
				VersionLister: releases.NewMockClient([]string{}),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if test.shouldFail {
				require.Nil(t, test.req)
			} else {
				require.NotNil(t, test.req)
				require.NotNil(t, test.req.VersionLister)
			}
		})
	}
}

// Test_TemplateFunctions tests the template functions.
func Test_TemplateFunctions(t *testing.T) {
	t.Parallel()

	mockVersions := []string{
		"1.14.3+ent",
		"1.15.6+ent",
		"1.16.10+ent",
		"1.17.6+ent",
		"1.17.5+ent",
		"1.16.9+ent",
		"1.17.4+ent",
	}

	ctx := context.Background()
	tc := &TemplateContext{
		VersionLister: releases.NewMockClient(mockVersions),
		Version:       "1.18.0",
	}

	t.Run("ParseVersion", func(t *testing.T) {
		t.Parallel()
		v, err := parseVersionFunc("1.18.0")
		require.NoError(t, err)
		require.Equal(t, int64(1), v.Major())
		require.Equal(t, int64(18), v.Minor())
		require.Equal(t, int64(0), v.Patch())
	})

	t.Run("ParseVersion invalid", func(t *testing.T) {
		t.Parallel()
		_, err := parseVersionFunc("invalid")
		require.Error(t, err)
	})

	t.Run("CompareVersions", func(t *testing.T) {
		t.Parallel()
		cmp, err := compareVersionsFunc("1.18.0", "1.17.0")
		require.NoError(t, err)
		require.Equal(t, 1, cmp)

		cmp, err = compareVersionsFunc("1.17.0", "1.18.0")
		require.NoError(t, err)
		require.Equal(t, -1, cmp)

		cmp, err = compareVersionsFunc("1.18.0", "1.18.0")
		require.NoError(t, err)
		require.Equal(t, 0, cmp)
	})

	t.Run("FilterVersions", func(t *testing.T) {
		t.Parallel()
		versions := []string{"1.16.10", "1.17.6", "1.17.5", "1.17.4"}
		filtered, err := filterVersionsFunc(versions, "~1.17.0")
		require.NoError(t, err)
		require.Equal(t, []string{"1.17.6", "1.17.5", "1.17.4"}, filtered)
	})

	t.Run("VersionsNMinus", func(t *testing.T) {
		t.Parallel()
		fn := versionsNMinusFunc(ctx, tc)
		versions, err := fn("vault", "enterprise", "1.18.0", 3, "minor")
		require.NoError(t, err)
		require.Equal(t, []string{"1.15.6", "1.16.9", "1.16.10", "1.17.4", "1.17.5", "1.17.6"}, versions)
	})

	t.Run("VersionsBounded", func(t *testing.T) {
		t.Parallel()
		fn := versionsBoundedFunc(ctx, tc)
		versions, err := fn("vault", "enterprise", "1.17.6", "1.17.4", "minor")
		require.NoError(t, err)
		require.Equal(t, []string{"1.17.4", "1.17.5", "1.17.6"}, versions)
	})

	t.Run("VersionsNMinusTransition", func(t *testing.T) {
		t.Parallel()
		fn := versionsNMinusTransitionFunc(ctx, tc)
		versions, err := fn("vault", "enterprise", "1.18.0", 3, "minor", "", "")
		require.NoError(t, err)
		require.Equal(t, []string{"1.15.6", "1.16.9", "1.16.10", "1.17.4", "1.17.5", "1.17.6"}, versions)
	})

	t.Run("VersionsNMinusTransition with transition", func(t *testing.T) {
		t.Parallel()

		// Mock versions that simulate a cadence transition
		transitionVersions := []string{
			"1.15.0+ent",
			"1.15.1+ent",
			"1.15.2+ent",
			"1.16.0+ent",
			"1.17.0+ent",
			"1.19.3+ent",
			"1.19.4+ent",
			"1.20.1+ent",
			"1.21.1+ent",
			"1.21.2+ent",
			"2.0.0-beta1+ent",
			"2.0.0+ent",
			"2.1.0+ent",
		}

		tc := &TemplateContext{
			VersionLister: releases.NewMockClient(transitionVersions),
			Version:       "2.1.1+ent",
		}
		fn := versionsNMinusTransitionFunc(ctx, tc)

		// Test major cadence with transition from minor
		versions, err := fn("vault", "enterprise", "2.1.1", 3, "major", "1.21.2", "minor")
		require.NoError(t, err)
		require.Equal(t, []string{"1.19.3", "1.19.4", "1.20.1", "1.21.1", "1.21.2", "2.0.0-beta1", "2.0.0", "2.1.0"}, versions)
	})

	t.Run("VersionsBoundedTransition", func(t *testing.T) {
		t.Parallel()
		// Mock versions that simulate a cadence transition
		transitionVersions := []string{
			"1.15.0+ent",
			"1.15.1+ent",
			"1.15.2+ent",
			"1.16.0+ent",
			"1.17.0+ent",
			"1.19.3+ent",
			"1.19.4+ent",
			"1.20.1+ent",
			"1.21.1+ent",
			"1.21.2+ent",
			"2.0.0-beta1+ent",
			"2.0.0+ent",
			"2.1.0+ent",
		}

		tc := &TemplateContext{
			VersionLister: releases.NewMockClient(transitionVersions),
			Version:       "2.1.1+ent",
		}
		fn := versionsBoundedTransitionFunc(ctx, tc)

		// Test major cadence with transition from minor
		versions, err := fn("vault", "enterprise", "2.1.1", "1.20.1", "major", "1.21.2", "minor")
		require.NoError(t, err)

		// Test with minor cadence
		require.NoError(t, err)
		require.Equal(t, []string{"1.20.1", "1.21.1", "1.21.2", "2.0.0-beta1", "2.0.0", "2.1.0"}, versions)
	})
}

// Test_FixtureTemplates tests rendering our templates in ./fixtures
func Test_FixtureTemplates(t *testing.T) {
	t.Parallel()

	mockVersions := []string{
		"1.16.10+ent",
		"1.17.6+ent",
		"1.17.5+ent",
		"1.16.9+ent",
		"1.17.4+ent",
		"1.16.8+ent",
		"1.17.0+ent",
		"1.17.1+ent",
		"1.17.2+ent",
		"1.17.3+ent",
		"1.18.0+ent",
	}

	// Fixed time for consistent test output
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	for name, test := range map[string]struct {
		fixturePath string
		version     string
		expected    string
	}{
		"test-template.tmpl": {
			fixturePath: "fixtures/test-template.tmpl",
			version:     "1.18.0",
			expected: `# Generated at 2024-01-15 10:30:00
# Current Version: 1.18.0

# Get last 3 versions before current version
Versions (N-Minus 3):
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0

# Example: Get versions between bounds
Versions (Bounded 1.17.0 to 1.18.0):
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0
`,
		},
		"test-functions.tmpl": {
			fixturePath: "fixtures/test-functions.tmpl",
			version:     "1.18.0",
			expected: `# Template Functions Test
# Generated at 2024-01-15 10:30:00

## Current Version Context
- Current Version: 1.18.0

## Version Utilities

### ParseVersion
- ParseVersion 1.16.0: 1.16.0
- ParseVersion 1.17.0: 1.17.0

### CompareVersions
- CompareVersions 1.16.0 vs 1.17.0: -1
  (Returns: -1 if v1 < v2, 0 if equal, 1 if v1 > v2)

## Version Listing Functions

### VersionsNMinus
Get last 3 versions before current version (1.18.0):
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0

### VersionsBounded
Get versions between 1.17.0 and 1.18.0:
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0

### FilterVersions
All versions (last 5):
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0

Filtered to ~1.17.0:
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6

### Skip Versions
Get versions skipping 1.17.5:
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.6
  - 1.18.0

### Version Comparison in Template Logic
1.18.0 is newer than 1.17.0
`,
		},
		"test-cadence.tmpl": {
			fixturePath: "fixtures/test-cadence.tmpl",
			version:     "1.18.0",
			expected: `# Cadence Template Test
# Generated at 2024-01-15 10:30:00

## Current Version Context
- Current Version: 1.18.0

## Cadence Functions

### VersionsNMinus - Minor Cadence
Get last 3 versions with minor cadence:
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0

### VersionsNMinus - Major Cadence
Get last 1 version with major cadence:
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0

### VersionsNMinusTransition
Get versions with cadence transition (major cadence, transitioned from minor at 1.15.2):
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0

### VersionsBoundedTransition - With Transition
Get versions with cadence transition (major cadence, transitioned from minor at 1.15.2):
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0

## Comparison: Regular vs Cadence Functions

### Regular VersionsNMinus (no cadence awareness)
Regular (11 versions):
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0

### Cadence-aware VersionsNMinusTransition (minor cadence)
Cadence-aware (11 versions):
  - 1.16.8
  - 1.16.9
  - 1.16.10
  - 1.17.0
  - 1.17.1
  - 1.17.2
  - 1.17.3
  - 1.17.4
  - 1.17.5
  - 1.17.6
  - 1.18.0
`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			tc := &TemplateContext{
				GeneratedAt:   fixedTime,
				VersionLister: releases.NewMockClient(mockVersions),
				Version:       test.version,
			}

			// Read and render the template directly
			body, err := os.ReadFile(test.fixturePath)
			require.NoError(t, err)

			tmpl, err := template.New("test").Funcs(templateFuncsFor(ctx, tc)).Parse(string(body))
			require.NoError(t, err)

			var buf strings.Builder
			err = tmpl.Execute(&buf, tc)
			require.NoError(t, err)
			require.Equal(t, test.expected, buf.String())
		})
	}
}

// renderTemplateString is a helper function to render a template string for testing
func renderTemplateString(ctx context.Context, tc *TemplateContext, templateStr string) ([]byte, error) {
	// We need to simulate the template rendering without file I/O
	tmpl, err := template.New("test").Funcs(templateFuncsFor(ctx, tc)).Parse(templateStr)
	if err != nil {
		return nil, err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, tc); err != nil {
		return nil, err
	}

	return []byte(buf.String()), nil
}
