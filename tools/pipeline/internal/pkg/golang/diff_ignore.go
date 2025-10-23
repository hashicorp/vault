// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
)

// diffIgnore compares the ignore directives in two module files and returns
// a slice of *Diff's. When strict parsing is set to true, ignore directives
// that are ommitted from one or the other modfile will also be included.
func diffIgnore(a *modfile.File, b *modfile.File) []*Diff {
	if (a == nil && b == nil) || (len(a.Ignore) == 0 && len(b.Ignore) == 0) {
		return nil
	}

	return slices.DeleteFunc(
		append(diffIgnoreFindDiffs(a, b), diffIgnoreFindDiffs(b, a)...),
		func(d *Diff) bool { return d == nil },
	)
}

func diffIgnoreFindDiffs(a, b *modfile.File) []*Diff {
	diffs := []*Diff{}
	for _, needle := range a.Ignore {
		if needle == nil {
			continue
		}

		// See if there's a matching ignore
		idx := slices.IndexFunc(b.Ignore, func(hay *modfile.Ignore) bool {
			if hay == nil {
				return false
			}

			return needle.Path == hay.Path
		})

		if idx >= 0 {
			// We have a matching ignore
			continue
		}

		diff := newDiffFromModFiles(a, b, DirectiveIgnore)
		if needle.Syntax != nil {
			diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ")}
		}
		diffs = append(diffs, diff)
	}

	return diffs
}
