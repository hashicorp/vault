// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"maps"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
)

// diffExclude compares the exclude directives in two module files and returns
// a slice of *Diff's. When strict parsing is set to true, exclude directives
// that are ommitted from one or the other modfile will also be included.
func diffExclude(a *modfile.File, b *modfile.File, strictDiff bool) []*Diff {
	if (a == nil && b == nil) || (len(a.Exclude) == 0 && len(b.Exclude) == 0) {
		return nil
	}

	var diffs []*Diff
	if strictDiff {
		diffs = append(diffExcludeFindMissing(a, b), diffExcludeFindMissing(b, a)...)
	}
	versionDiffsA := diffExcludeFindDifferent(a, b)
	versionDiffsB := diffExcludeFindDifferent(b, a)
	maps.Copy(versionDiffsB, versionDiffsA)

	return slices.DeleteFunc(
		append(diffs, slices.Collect(maps.Values(versionDiffsB))...),
		func(d *Diff) bool { return d == nil },
	)
}

func diffExcludeFindMissing(a, b *modfile.File) []*Diff {
	diffs := []*Diff{}
	for _, needle := range a.Exclude {
		if needle == nil {
			continue
		}

		// See if there's a matching require
		idx := slices.IndexFunc(b.Exclude, func(hay *modfile.Exclude) bool {
			if hay == nil {
				return false
			}

			return needle.Mod.Path == hay.Mod.Path
		})

		if idx >= 0 {
			// We have a matching require
			continue
		}

		diff := newDiffFromModFiles(a, b, DirectiveExclude)
		if needle.Syntax != nil {
			diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ")}
		}
		diffs = append(diffs, diff)
	}

	return diffs
}

func diffExcludeFindDifferent(a, b *modfile.File) map[string]*Diff {
	diffs := map[string]*Diff{}
	for _, needle := range a.Exclude {
		if needle == nil {
			continue
		}

		// See if there's a matching require
		idx := slices.IndexFunc(b.Exclude, func(hay *modfile.Exclude) bool {
			if hay == nil {
				return false
			}

			return needle.Mod.Path == hay.Mod.Path
		})

		if idx < 0 {
			// We don't have a matching require
			continue
		}

		hay := b.Exclude[idx]
		if needle.Mod.Version != hay.Mod.Version {
			diff := newDiffFromModFiles(a, b, DirectiveExclude)
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
