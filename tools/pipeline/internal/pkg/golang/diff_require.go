// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"maps"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
)

func diffRequire(a *modfile.File, b *modfile.File, strictDiff bool, excludeRequire []string) []*Diff {
	if (a == nil && b == nil) || (len(a.Require) == 0 && len(b.Require) == 0) {
		return nil
	}

	var diffs []*Diff
	if strictDiff {
		diffs = append(diffRequireFindMissing(a, b, excludeRequire), diffRequireFindMissing(b, a, excludeRequire)...)
	}
	versionDiffsA := diffRequireFindDifferent(a, b, excludeRequire)
	versionDiffsB := diffRequireFindDifferent(b, a, excludeRequire)
	maps.Copy(versionDiffsB, versionDiffsA)

	return slices.DeleteFunc(
		append(diffs, slices.Collect(maps.Values(versionDiffsB))...),
		func(d *Diff) bool { return d == nil },
	)
}

func diffRequireFindMissing(a, b *modfile.File, excludeRequire []string) []*Diff {
	diffs := []*Diff{}
	for _, needle := range a.Require {
		if needle == nil {
			continue
		}

		if matchesAnyExcludePattern(excludeRequire, needle.Mod.Path) {
			continue
		}

		// See if there's a matching require
		idx := slices.IndexFunc(b.Require, func(hay *modfile.Require) bool {
			if hay == nil {
				return false
			}

			return needle.Mod.Path == hay.Mod.Path
		})

		if idx >= 0 {
			// We have a matching require
			continue
		}

		diff := newDiffFromModFiles(a, b, DirectiveRequire)
		if needle.Syntax != nil {
			diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ")}
		}
		diffs = append(diffs, diff)
	}

	return diffs
}

func diffRequireFindDifferent(a, b *modfile.File, excludeRequire []string) map[string]*Diff {
	diffs := map[string]*Diff{}
	for _, needle := range a.Require {
		if needle == nil {
			continue
		}

		if matchesAnyExcludePattern(excludeRequire, needle.Mod.Path) {
			continue
		}

		// See if there's a matching require
		idx := slices.IndexFunc(b.Require, func(hay *modfile.Require) bool {
			if hay == nil {
				return false
			}

			return needle.Mod.Path == hay.Mod.Path
		})

		if idx < 0 {
			// We don't have a matching require
			continue
		}

		hay := b.Require[idx]
		if needle.Mod.Version != hay.Mod.Version {
			diff := newDiffFromModFiles(a, b, DirectiveRequire)
			if needle.Syntax != nil {
				diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ") + "\n"}
			}
			if hay.Syntax != nil {
				diff.Diff.B = []string{strings.Join(hay.Syntax.Token, " ")}
			}
			diffs[needle.Mod.Path] = diff
		}
	}

	return diffs
}
