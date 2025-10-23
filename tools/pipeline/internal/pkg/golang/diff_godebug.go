// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
)

// diffGodebug compares the godebug directives in two modules files and returns
// a slice of *Diff's. When strict parsing is set to true, godebug directives
// that are ommitted from one or the other modfile will also be included.
func diffGodebug(a *modfile.File, b *modfile.File) []*Diff {
	if (a == nil && b == nil) || (len(a.Godebug) == 0 && len(b.Godebug) == 0) {
		return nil
	}

	return slices.DeleteFunc(
		append(diffGodebugFindDiffs(a, b), diffGodebugFindDiffs(b, a)...),
		func(d *Diff) bool { return d == nil },
	)
}

func diffGodebugFindDiffs(a, b *modfile.File) []*Diff {
	diffs := []*Diff{}
	for _, needle := range a.Godebug {
		idx := slices.IndexFunc(b.Godebug, func(hay *modfile.Godebug) bool {
			return godebugMatchingKey(needle, hay)
		})

		if idx < 0 {
			// We don't have matching godebug with the same key, create a single
			// sided diff
			diff := newDiffFromModFiles(a, b, DirectiveGodebug)
			if needle != nil && needle.Syntax != nil {
				diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ")}
			}
			diffs = append(diffs, diff)
			continue
		}

		// We have a matching key
		hay := b.Godebug[idx]
		if godebugEqual(needle, hay) {
			continue
		}

		// We have differing values. Create a double sided diff
		diff := newDiffFromModFiles(a, b, DirectiveGodebug)
		if needle != nil && needle.Syntax != nil {
			diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ") + "\n"}
		}
		if hay != nil && hay.Syntax != nil {
			diff.Diff.B = []string{strings.Join(hay.Syntax.Token, " ")}
		}

		diffs = append(diffs, diff)
	}

	return diffs
}

func godebugEqual(a *modfile.Godebug, b *modfile.Godebug) bool {
	if a == nil && b == nil {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	if a.Key != b.Key {
		return false
	}

	return a.Value == b.Value
}

func godebugMatchingKey(a *modfile.Godebug, b *modfile.Godebug) bool {
	if a == nil && b == nil {
		return false
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	return a.Key == b.Key
}
