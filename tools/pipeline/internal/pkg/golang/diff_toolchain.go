// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"strings"

	"golang.org/x/mod/modfile"
)

func diffToolchain(a *modfile.File, b *modfile.File) *Diff {
	if (a == nil && b == nil) ||
		(a.Toolchain == nil && b.Toolchain == nil) ||
		((a.Toolchain != nil && b.Toolchain != nil) && (a.Toolchain.Name == b.Toolchain.Name)) {

		return nil
	}

	diff := newDiffFromModFiles(a, b, DirectiveToolchain)
	if a != nil && a.Toolchain != nil && a.Toolchain.Syntax != nil {
		diff.Diff.A = []string{strings.Join(a.Toolchain.Syntax.Token, " ") + "\n"}
	}
	if b != nil && b.Toolchain != nil && b.Toolchain.Syntax != nil {
		diff.Diff.B = []string{strings.Join(b.Toolchain.Syntax.Token, " ")}
	}

	return diff
}
