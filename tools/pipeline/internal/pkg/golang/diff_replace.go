// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"maps"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
)

func diffReplace(a *modfile.File, b *modfile.File, strictDiff bool) []*Diff {
	if (a == nil && b == nil) || (len(a.Replace) == 0 && len(b.Replace) == 0) {
		return nil
	}

	var diffs []*Diff
	if strictDiff {
		diffs = append(diffReplaceFindMissing(a, b), diffReplaceFindMissing(b, a)...)
	}
	versionDiffsA := diffReplaceFindDifferent(a, b)
	versionDiffsB := diffReplaceFindDifferent(b, a)
	maps.Copy(versionDiffsB, versionDiffsA)

	return slices.DeleteFunc(
		append(diffs, slices.Collect(maps.Values(versionDiffsB))...),
		func(d *Diff) bool { return d == nil },
	)
}

func diffReplaceFindMissing(a, b *modfile.File) []*Diff {
	diffs := []*Diff{}
	for _, needle := range a.Replace {
		if needle == nil {
			continue
		}

		// See if there's a matching old exclude
		idx := slices.IndexFunc(b.Replace, func(hay *modfile.Replace) bool {
			if hay == nil {
				return false
			}

			return needle.Old.Path == hay.Old.Path
		})

		if idx >= 0 {
			// We have a matching old exclude, we'll handle this in the version diff check
			continue
		}

		diff := newDiffFromModFiles(a, b, DirectiveReplace)
		if needle.Syntax != nil {
			diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ")}
		}
		diffs = append(diffs, diff)
	}

	return diffs
}

func diffReplaceFindDifferent(a, b *modfile.File) map[string]*Diff {
	diffs := map[string]*Diff{}
	for _, needle := range a.Replace {
		if needle == nil {
			continue
		}

		// See if there's a matching exclude
		idx := slices.IndexFunc(b.Replace, func(hay *modfile.Replace) bool {
			if hay == nil {
				return false
			}

			return needle.Old.Path == hay.Old.Path
		})

		if idx < 0 {
			// We don't have a matching exclude
			continue
		}

		hay := b.Replace[idx]
		if replaceEqual(needle, hay) {
			continue
		}

		diff := newDiffFromModFiles(a, b, DirectiveReplace)
		if needle.Syntax != nil {
			diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ") + "\n"}
		}
		if hay.Syntax != nil {
			diff.Diff.B = []string{strings.Join(hay.Syntax.Token, " ")}
		}
		diffs[needle.Old.Path] = diff
	}

	return diffs
}

func replaceEqual(a *modfile.Replace, b *modfile.Replace) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	if a.Old.Path != b.Old.Path {
		return false
	}

	if a.Old.Version != b.Old.Version {
		return false
	}

	if a.New.Path != b.New.Path {
		return false
	}

	return a.New.Version == b.New.Version
}
