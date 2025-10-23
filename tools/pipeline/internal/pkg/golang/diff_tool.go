// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"slices"
	"strings"

	"golang.org/x/mod/modfile"
)

func diffTool(a *modfile.File, b *modfile.File) []*Diff {
	if (a == nil && b == nil) || (len(a.Tool) == 0 && len(b.Tool) == 0) {
		return nil
	}

	return slices.DeleteFunc(
		append(diffToolFindDiffs(a, b), diffToolFindDiffs(b, a)...),
		func(d *Diff) bool { return d == nil },
	)
}

func diffToolFindDiffs(a, b *modfile.File) []*Diff {
	diffs := []*Diff{}
	for _, needle := range a.Tool {
		if needle == nil {
			continue
		}

		// See if there's a matching tool
		idx := slices.IndexFunc(b.Tool, func(hay *modfile.Tool) bool {
			if hay == nil {
				return false
			}

			return needle.Path == hay.Path
		})

		if idx >= 0 {
			// We have a matching tool
			continue
		}

		diff := newDiffFromModFiles(a, b, DirectiveTool)
		if needle.Syntax != nil {
			diff.Diff.A = []string{strings.Join(needle.Syntax.Token, " ")}
		}
		diffs = append(diffs, diff)
	}

	return diffs
}
