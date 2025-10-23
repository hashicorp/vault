// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"strings"

	"golang.org/x/mod/modfile"
)

// diffModule compares the module directives in two modules files and returns
// a slice of *Diff's.
func diffModule(a *modfile.File, b *modfile.File) *Diff {
	if (a == nil && b == nil) ||
		(a.Module == nil && b.Module == nil) ||
		(a.Module.Mod.String() == b.Module.Mod.String()) {

		return nil
	}

	diff := newDiffFromModFiles(a, b, DirectiveModule)
	if a != nil && a.Module != nil && a.Module.Syntax != nil {
		diff.Diff.A = []string{strings.Join(a.Module.Syntax.Token, " ") + "\n"}
	}
	if b != nil && b.Module != nil && b.Module.Syntax != nil {
		diff.Diff.B = []string{strings.Join(b.Module.Syntax.Token, " ")}
	}

	return diff
}
