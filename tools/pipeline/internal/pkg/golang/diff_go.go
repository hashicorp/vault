// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"strings"

	"golang.org/x/mod/modfile"
)

// diffGo compares the go directives in two modules files and returns
// a slice of *Diff's.
func diffGo(a *modfile.File, b *modfile.File) *Diff {
	if (a == nil && b == nil) ||
		(a.Go == nil && b.Go == nil) ||
		((a.Go != nil && b.Go != nil) && (a.Go.Version == b.Go.Version)) {

		return nil
	}

	diff := newDiffFromModFiles(a, b, DirectiveGo)
	if a != nil && a.Go != nil && a.Go.Syntax != nil {
		diff.Diff.A = []string{strings.Join(a.Go.Syntax.Token, " ") + "\n"}
	}
	if b != nil && b.Go != nil && b.Go.Syntax != nil {
		diff.Diff.B = []string{strings.Join(b.Go.Syntax.Token, " ")}
	}

	return diff
}
