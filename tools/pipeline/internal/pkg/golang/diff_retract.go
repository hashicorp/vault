// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"maps"
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
)

func diffRetract(a *modfile.File, b *modfile.File, strictDiff bool) []*Diff {
	if (a == nil && b == nil) || (len(a.Retract) == 0 && len(b.Retract) == 0) {
		return nil
	}

	var diffs []*Diff
	if strictDiff {
		diffs = append(diffRetractFindMissing(a, b), diffRetractFindMissing(b, a)...)
	}
	versionDiffsA := diffRetractFindDifferent(a, b)
	versionDiffsB := diffRetractFindDifferent(b, a)
	maps.Copy(versionDiffsB, versionDiffsA)

	return slices.DeleteFunc(
		append(diffs, slices.Collect(maps.Values(versionDiffsB))...),
		func(d *Diff) bool { return d == nil },
	)
}

func diffRetractFindMissing(a, b *modfile.File) []*Diff {
	diffs := []*Diff{}
	for _, needle := range a.Retract {
		if needle == nil {
			continue
		}

		// See if there's a matching retract
		idx := slices.IndexFunc(b.Retract, func(hay *modfile.Retract) bool {
			if hay == nil {
				return false
			}

			return needle.VersionInterval.Low == hay.VersionInterval.Low
		})

		if idx >= 0 {
			// We have a matching retract
			continue
		}

		diff := newDiffFromModFiles(a, b, DirectiveRetract)
		if needle.Syntax != nil {
			diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ")}
		}
		diffs = append(diffs, diff)
	}

	return diffs
}

func diffRetractFindDifferent(a, b *modfile.File) map[string]*Diff {
	diffs := map[string]*Diff{}
	for _, needle := range a.Retract {
		if needle == nil {
			continue
		}

		// See if there's a matching require
		idx := slices.IndexFunc(b.Retract, func(hay *modfile.Retract) bool {
			if hay == nil {
				return false
			}

			return needle.VersionInterval.Low == hay.VersionInterval.Low
		})

		if idx < 0 {
			// We don't have a matching require
			continue
		}

		hay := b.Retract[idx]
		if retractEqual(needle, hay) {
			continue
		}

		diff := newDiffFromModFiles(a, b, DirectiveRetract)
		if needle.Syntax != nil {
			diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ") + "\n"}
		}
		if hay.Syntax != nil {
			diff.Diff.B = []string{strings.Join(hay.Syntax.Token, " ")}
		}
		diffs[needle.VersionInterval.Low] = diff
	}

	return diffs
}

func retractEqual(a *modfile.Retract, b *modfile.Retract) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	if a.VersionInterval.Low != b.VersionInterval.Low {
		return false
	}

	if a.VersionInterval.High != b.VersionInterval.High {
		return false
	}

	return a.Rationale == b.Rationale
}
