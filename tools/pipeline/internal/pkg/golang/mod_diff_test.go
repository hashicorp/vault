// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DiffModFiles_Equal(t *testing.T) {
	t.Parallel()

	// We intentionally do not use the .mod suffix to avoid
	// tooling from automatically tidying these files.
	modA, err := os.ReadFile("./fixtures/go.moda")
	require.NoError(t, err)

	modB, err := os.ReadFile("./fixtures/go.modb")
	require.NoError(t, err)

	for desc, test := range map[string]struct {
		as   *ModSource
		bs   *ModSource
		opts *DiffOpts
	}{
		"moda not lax not strict": {
			&ModSource{Name: "moda-1", Data: modA},
			&ModSource{Name: "moda-2", Data: modA},
			&DiffOpts{StrictDiffRequire: false, ParseLax: false},
		},
		"moda lax not strict": {
			&ModSource{Name: "moda-1", Data: modA},
			&ModSource{Name: "moda-2", Data: modA},
			&DiffOpts{StrictDiffRequire: false, ParseLax: true},
		},
		"moda not lax strict": {
			&ModSource{Name: "moda-1", Data: modA},
			&ModSource{Name: "moda-2", Data: modA},
			&DiffOpts{StrictDiffRequire: true, ParseLax: false},
		},
		"moda lax strict": {
			&ModSource{Name: "moda-1", Data: modA},
			&ModSource{Name: "moda-2", Data: modA},
			&DiffOpts{StrictDiffRequire: true, ParseLax: true},
		},
		"modb not lax not strict": {
			&ModSource{Name: "modb-1", Data: modB},
			&ModSource{Name: "modb-2", Data: modB},
			&DiffOpts{StrictDiffRequire: false, ParseLax: false},
		},
		"modb lax not strict": {
			&ModSource{Name: "modb-1", Data: modB},
			&ModSource{Name: "modb-2", Data: modB},
			&DiffOpts{StrictDiffRequire: false, ParseLax: true},
		},
		"modb not lax strict": {
			&ModSource{Name: "modb-1", Data: modB},
			&ModSource{Name: "modb-2", Data: modB},
			&DiffOpts{StrictDiffRequire: true, ParseLax: false},
		},
		"modb lax strict": {
			&ModSource{Name: "modb-1", Data: modB},
			&ModSource{Name: "modb-2", Data: modB},
			&DiffOpts{StrictDiffRequire: true, ParseLax: true},
		},
	} {
		t.Run(desc, func(t *testing.T) {
			t.Parallel()
			diff, err := DiffModFiles(test.as, test.bs, test.opts)
			require.NoError(t, err)
			require.Nil(t, diff, "expected no diff, got:\n%v", printModDiff(diff))
		})
	}
}

func Test_DiffModFiles_Diff(t *testing.T) {
	t.Parallel()

	// We intentionally do not use the .mod suffix to avoid
	// tooling from automatically tidying these files.
	modA, err := os.ReadFile("./fixtures/go.moda")
	require.NoError(t, err)

	modB, err := os.ReadFile("./fixtures/go.modb")
	require.NoError(t, err)

	as := &ModSource{Name: "moda", Data: modA}
	bs := &ModSource{Name: "modb", Data: modB}

	for desc, test := range map[string]struct {
		opts      *DiffOpts
		condition func(t *testing.T, got ModDiff)
	}{
		"strict parse lax diff": {
			&DiffOpts{
				ParseLax:          false,
				Module:            true,
				Go:                true,
				Toolchain:         true,
				Godebug:           true,
				Require:           true,
				Exclude:           true,
				Replace:           true,
				Retract:           true,
				Tool:              true,
				Ignore:            true,
				StrictDiffRequire: false,
				StrictDiffExclude: false,
				StrictDiffReplace: false,
				StrictDiffRetract: false,
			},
			func(t *testing.T, got ModDiff) {
				hasDiffMatching(t, got, DirectiveModule, []string{
					"-module github.com/hashicorp/vault/pipeline/golang/moda",
					"+module github.com/hashicorp/vault/pipeline/golang/modb",
				})
				hasDiffMatching(t, got, DirectiveGo, []string{
					"-go 1.25",
					"+go 1.25.2",
				})
				hasDiffMatching(t, got, DirectiveToolchain, []string{
					"+toolchain go1.24",
				})
				hasDiffMatching(t, got, DirectiveGodebug, []string{
					"-default=go1.21",
				})
				hasDiffMatching(t, got, DirectiveGodebug, []string{
					"-httpcookiemaxnum=4000",
				})
				hasDiffMatching(t, got, DirectiveGodebug, []string{
					"-panicnil=1",
				})
				hasDiffMatching(t, got, DirectiveRequire, []string{
					"-require github.com/99designs/keyring v1.2.2",
					"+require github.com/99designs/keyring v0.0.0-00010101000000-000000000000",
				})
				noDiffMatching(t, got, DirectiveRequire, []string{
					"-github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c",
				})
				noDiffMatching(t, got, DirectiveExclude, []string{
					"-exclude golang.org/x/term v0.2.0",
				})
				noDiffMatching(t, got, DirectiveReplace, []string{
					"-replace github.com/99designs/keyring => github.com/Jeffail/keyring v1.2.3",
				})
				noDiffMatching(t, got, DirectiveRetract, []string{
					"-[ v1.0.0 , v1.9.9 ]",
				})
				noDiffMatching(t, got, DirectiveRetract, []string{
					"-v0.9.0",
				})
				hasDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/bisect",
				})
				hasDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/stringer",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-./third_party/javascript",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-content/html",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-static",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-ignore ./node_modules",
				})
			},
		},
		"lax parse lax diff": {
			&DiffOpts{
				ParseLax:          true,
				Module:            true,
				Go:                true,
				Toolchain:         true,
				Godebug:           true,
				Require:           true,
				Exclude:           true,
				Replace:           true,
				Retract:           true,
				Tool:              true,
				Ignore:            true,
				StrictDiffRequire: false,
				StrictDiffExclude: false,
				StrictDiffReplace: false,
				StrictDiffRetract: false,
			},
			func(t *testing.T, got ModDiff) {
				hasDiffMatching(t, got, DirectiveModule, []string{
					"-module github.com/hashicorp/vault/pipeline/golang/moda",
					"+module github.com/hashicorp/vault/pipeline/golang/modb",
				})
				hasDiffMatching(t, got, DirectiveGo, []string{
					"-go 1.25",
					"+go 1.25.2",
				})
				noDiffMatching(t, got, DirectiveToolchain, []string{
					"+toolchain go1.24",
				})
				noDiffMatching(t, got, DirectiveGodebug, []string{
					"-default=go1.21",
				})
				noDiffMatching(t, got, DirectiveGodebug, []string{
					"-httpcookiemaxnum=4000",
				})
				noDiffMatching(t, got, DirectiveGodebug, []string{
					"-panicnil=1",
				})
				hasDiffMatching(t, got, DirectiveRequire, []string{
					"-require github.com/99designs/keyring v1.2.2",
					"+require github.com/99designs/keyring v0.0.0-00010101000000-000000000000",
				})
				noDiffMatching(t, got, DirectiveRequire, []string{
					"-github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c",
				})
				noDiffMatching(t, got, DirectiveExclude, []string{
					"-exclude golang.org/x/term v0.2.0",
				})
				noDiffMatching(t, got, DirectiveReplace, []string{
					"-replace github.com/99designs/keyring => github.com/Jeffail/keyring v1.2.3",
				})
				noDiffMatching(t, got, DirectiveRetract, []string{
					"-[ v1.0.0 , v1.9.9 ]",
				})
				noDiffMatching(t, got, DirectiveRetract, []string{
					"-v0.9.0",
				})
				noDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/bisect",
				})
				noDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/stringer",
				})
				noDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/bisect",
				})
				noDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/stringer",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-./third_party/javascript",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-content/html",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-static",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-ignore ./node_modules",
				})
			},
		},
		"strict parse strict diff": {
			&DiffOpts{
				ParseLax:          false,
				Module:            true,
				Go:                true,
				Toolchain:         true,
				Godebug:           true,
				Require:           true,
				Exclude:           true,
				Replace:           true,
				Retract:           true,
				Tool:              true,
				Ignore:            true,
				StrictDiffRequire: true,
				StrictDiffExclude: true,
				StrictDiffReplace: true,
				StrictDiffRetract: true,
			},
			func(t *testing.T, got ModDiff) {
				hasDiffMatching(t, got, DirectiveModule, []string{
					"-module github.com/hashicorp/vault/pipeline/golang/moda",
					"+module github.com/hashicorp/vault/pipeline/golang/modb",
				})
				hasDiffMatching(t, got, DirectiveGo, []string{
					"-go 1.25",
					"+go 1.25.2",
				})
				hasDiffMatching(t, got, DirectiveToolchain, []string{
					"+toolchain go1.24",
				})
				hasDiffMatching(t, got, DirectiveGodebug, []string{
					"-default=go1.21",
				})
				hasDiffMatching(t, got, DirectiveGodebug, []string{
					"-httpcookiemaxnum=4000",
				})
				hasDiffMatching(t, got, DirectiveGodebug, []string{
					"-panicnil=1",
				})
				hasDiffMatching(t, got, DirectiveRequire, []string{
					"-require github.com/99designs/keyring v1.2.2",
					"+require github.com/99designs/keyring v0.0.0-00010101000000-000000000000",
				})
				hasDiffMatching(t, got, DirectiveRequire, []string{
					"-github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c",
				})
				hasDiffMatching(t, got, DirectiveExclude, []string{
					"-exclude golang.org/x/term v0.2.0",
				})
				hasDiffMatching(t, got, DirectiveReplace, []string{
					"-replace github.com/99designs/keyring => github.com/Jeffail/keyring v1.2.3",
				})
				hasDiffMatching(t, got, DirectiveRetract, []string{
					"-[ v1.0.0 , v1.9.9 ]",
				})
				hasDiffMatching(t, got, DirectiveRetract, []string{
					"-v0.9.0",
				})
				hasDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/bisect",
				})
				hasDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/stringer",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-./third_party/javascript",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-content/html",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-static",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-ignore ./node_modules",
				})
			},
		},
		"lax parse strict diff": {
			&DiffOpts{
				ParseLax:          true,
				Module:            true,
				Go:                true,
				Toolchain:         true,
				Godebug:           true,
				Require:           true,
				Exclude:           true,
				Replace:           true,
				Retract:           true,
				Tool:              true,
				Ignore:            true,
				StrictDiffRequire: true,
				StrictDiffExclude: true,
				StrictDiffReplace: true,
				StrictDiffRetract: true,
			},
			func(t *testing.T, got ModDiff) {
				hasDiffMatching(t, got, DirectiveModule, []string{
					"-module github.com/hashicorp/vault/pipeline/golang/moda",
					"+module github.com/hashicorp/vault/pipeline/golang/modb",
				})
				hasDiffMatching(t, got, DirectiveGo, []string{
					"-go 1.25",
					"+go 1.25.2",
				})
				noDiffMatching(t, got, DirectiveToolchain, []string{
					"+toolchain go1.24",
				})
				noDiffMatching(t, got, DirectiveGodebug, []string{
					"-default=go1.21",
				})
				noDiffMatching(t, got, DirectiveGodebug, []string{
					"-httpcookiemaxnum=4000",
				})
				noDiffMatching(t, got, DirectiveGodebug, []string{
					"-panicnil=1",
				})
				hasDiffMatching(t, got, DirectiveRequire, []string{
					"-require github.com/99designs/keyring v1.2.2",
					"+require github.com/99designs/keyring v0.0.0-00010101000000-000000000000",
				})
				hasDiffMatching(t, got, DirectiveRequire, []string{
					"-github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c",
				})
				noDiffMatching(t, got, DirectiveExclude, []string{
					"-exclude golang.org/x/term v0.2.0",
				})
				noDiffMatching(t, got, DirectiveReplace, []string{
					"-replace github.com/99designs/keyring => github.com/Jeffail/keyring v1.2.3",
				})
				hasDiffMatching(t, got, DirectiveRetract, []string{
					"-[ v1.0.0 , v1.9.9 ]",
				})
				hasDiffMatching(t, got, DirectiveRetract, []string{
					"-v0.9.0",
				})
				noDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/bisect",
				})
				noDiffMatching(t, got, DirectiveTool, []string{
					"-tool golang.org/x/tools/cmd/stringer",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-content/html",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-static",
				})
				hasDiffMatching(t, got, DirectiveIgnore, []string{
					"-ignore ./node_modules",
				})
			},
		},
	} {
		t.Run(desc, func(t *testing.T) {
			t.Parallel()
			diff, err := DiffModFiles(as, bs, test.opts)
			require.NoError(t, err)
			require.NotNil(t, diff, "expected a module diff")
			test.condition(t, diff)
		})
	}
}

func hasDiffMatching(t *testing.T, diff ModDiff, dir Directive, matches []string) {
	t.Helper()

	require.NotNil(t, diff)
	diffs := getDiffsForDirective(dir, diff)
	require.True(t, len(diffs) > 0, "expected %s matching %v, got diff:\n%s",
		dir,
		matches,
		printModDiff(diff),
	)
	for _, df := range diffs {
		if unifiedTextMatches(df, matches) {
			return
		}
	}
	t.Fatalf("expected %s diff matching %v, got diff:\n%s",
		dir,
		matches,
		printModDiff(diff),
	)
}

func noDiffMatching(t *testing.T, diff ModDiff, dir Directive, matches []string) {
	t.Helper()

	require.NotNil(t, diff)
	diffs := getDiffsForDirective(dir, diff)
	if len(diffs) < 1 {
		return
	}
	for _, df := range diffs {
		if unifiedTextMatches(df, matches) {
			t.Fatalf("expected no %s diff matching %v, got diff:\n%s",
				dir,
				matches,
				printModDiff(diff),
			)
		}
	}
}

func printModDiff(d ModDiff) string {
	if d == nil {
		return ""
	}

	b := strings.Builder{}
	for _, diff := range d {
		if exp := diff.Explanation(); exp != "" {
			b.WriteString(exp + "\n")
		}
		if dt := diff.UnifiedText(); dt != "" {
			b.WriteString(dt + "\n")
		}
	}

	return b.String()
}

func getDiffsForDirective(dir Directive, diff ModDiff) []*Diff {
	if len(diff) < 1 {
		return nil
	}

	diffs := []*Diff{}
	for _, d := range diff {
		if d == nil {
			continue
		}
		if d.Directive == dir {
			diffs = append(diffs, d)
		}
	}

	return diffs
}

func unifiedTextMatches(diff *Diff, matches []string) bool {
	if diff == nil && len(matches) == 0 {
		return true
	}

	if diff == nil && len(matches) > 0 {
		return false
	}

	txt := diff.UnifiedText()
	if txt == "" {
		return false
	}
	for _, m := range matches {
		if !strings.Contains(txt, m) {
			return false
		}
	}

	return true
}
