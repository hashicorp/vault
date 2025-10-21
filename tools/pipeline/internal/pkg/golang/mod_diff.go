// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/mod/modfile"
)

// DiffOpts are options for the module diff.
type DiffOpts struct {
	ParseLax bool

	Module    bool
	Go        bool
	Toolchain bool
	Godebug   bool
	Require   bool
	Exclude   bool
	Replace   bool
	Retract   bool
	Tool      bool
	Ignore    bool

	StrictDiffRequire bool
	StrictDiffExclude bool
	StrictDiffReplace bool
	StrictDiffRetract bool
}

func DefaultDiffOpts() *DiffOpts {
	return &DiffOpts{
		ParseLax:          false,
		Module:            true,
		Go:                true,
		Toolchain:         true,
		Godebug:           true,
		Require:           true,
		Exclude:           true,
		Replace:           true,
		Retract:           true,
		Tool:              true,
		Ignore:            true,
		StrictDiffRequire: true,
		StrictDiffExclude: true,
		StrictDiffReplace: true,
		StrictDiffRetract: true,
	}
}

// ModDiff is the result of comparing two Go modules.
type ModDiff []*Diff

// ModSource is a go.mod file
type ModSource struct {
	Name string
	Data []byte
}

type Diff struct {
	Directive Directive
	Diff      *difflib.UnifiedDiff
}

type Directive string

const (
	DirectiveModule    Directive = "module"
	DirectiveGo        Directive = "go"
	DirectiveToolchain Directive = "toolchain"
	DirectiveGodebug   Directive = "godebug"
	DirectiveRequire   Directive = "require"
	DirectiveExclude   Directive = "exclude"
	DirectiveReplace   Directive = "replace"
	DirectiveRetract   Directive = "retract"
	DirectiveTool      Directive = "tool"
	DirectiveIgnore    Directive = "ignore"
)

func (d *Diff) Explanation() string {
	if d == nil || d.Directive == "" {
		return ""
	}

	return d.Directive.Explanation()
}

func (d *Diff) UnifiedText() string {
	if d == nil || d.Diff == nil {
		return ""
	}

	txt, _ := difflib.GetUnifiedDiffString(*d.Diff)

	return txt
}

func (d Directive) Explanation() string {
	return fmt.Sprintf("The '%s' directives do not match", d)
}

// DiffModFiles diffs two go.mod "files" and returns a ModDiff.
func DiffModFiles(as *ModSource, bs *ModSource, opts *DiffOpts) (ModDiff, error) {
	if as == nil {
		return nil, errors.New("missing a mod source")
	}

	if bs == nil {
		return nil, errors.New("missing b mod source")
	}

	var af *modfile.File
	var bf *modfile.File
	var err error
	if opts.ParseLax {
		af, err = modfile.ParseLax(as.Name, as.Data, nil)
		if err != nil {
			return nil, fmt.Errorf("parsing %s contents: %w", as.Name, err)
		}

		bf, err = modfile.ParseLax(bs.Name, bs.Data, nil)
		if err != nil {
			return nil, fmt.Errorf("parsing %s contents: %w", bs.Name, err)
		}
	} else {
		af, err = modfile.Parse(as.Name, as.Data, nil)
		if err != nil {
			return nil, fmt.Errorf("parsing %s contents: %w", as.Name, err)
		}

		bf, err = modfile.Parse(bs.Name, bs.Data, nil)
		if err != nil {
			return nil, fmt.Errorf("parsing %s contents: %w", bs.Name, err)
		}
	}

	return diffModFiles(af, bf, opts)
}

func diffModFiles(a *modfile.File, b *modfile.File, opts *DiffOpts) (ModDiff, error) {
	if a == nil {
		return nil, errors.New("missing a mod file")
	}

	if b == nil {
		return nil, errors.New("missing b mod file")
	}

	diff := ModDiff{}
	if opts.Module {
		diff = append(diff, diffModule(a, b))
	}
	if opts.Go {
		diff = append(diff, diffGo(a, b))
	}
	if opts.Toolchain {
		diff = append(diff, diffToolchain(a, b))
	}
	if opts.Godebug {
		diff = append(diff, diffGodebug(a, b)...)
	}
	if opts.Require || opts.StrictDiffRequire {
		diff = append(diff, diffRequire(a, b, opts.StrictDiffRequire)...)
	}
	if opts.Exclude || opts.StrictDiffExclude {
		diff = append(diff, diffExclude(a, b, opts.StrictDiffExclude)...)
	}
	if opts.Replace || opts.StrictDiffReplace {
		diff = append(diff, diffReplace(a, b, opts.StrictDiffReplace)...)
	}
	if opts.Retract || opts.StrictDiffRetract {
		diff = append(diff, diffRetract(a, b, opts.StrictDiffRetract)...)
	}
	if opts.Tool {
		diff = append(diff, diffTool(a, b)...)
	}
	if opts.Ignore {
		diff = append(diff, diffIgnore(a, b)...)
	}

	if len(diff) == 0 {
		return nil, nil
	}

	diff = slices.DeleteFunc(diff, func(d *Diff) bool { return d == nil })
	slices.SortStableFunc(diff, func(a *Diff, b *Diff) int {
		// a < b = 1
		// a == b = 0
		// a > b = -1
		if a == nil && b == nil {
			return 0
		}

		if a == nil && b != nil {
			return 1
		}

		if a != nil && b == nil {
			return -1
		}

		if n := strings.Compare(string(a.Directive), string(b.Directive)); n != 0 {
			return n
		}

		if a.Diff == nil && b.Diff == nil {
			return 0
		}

		if a.Diff == nil && b.Diff != nil {
			return 1
		}

		if a.Diff != nil && b.Diff == nil {
			return -1
		}

		atxt, aerr := difflib.GetUnifiedDiffString(*a.Diff)
		btxt, berr := difflib.GetUnifiedDiffString(*b.Diff)

		if aerr == nil && berr != nil {
			return 1
		}

		if aerr != nil && berr == nil {
			return -1
		}

		return strings.Compare(atxt, btxt)
	})

	return diff, nil
}

func newDiffFromModFiles(a, b *modfile.File, typ Directive) *Diff {
	res := &Diff{
		Directive: typ,
		Diff:      &difflib.UnifiedDiff{Context: 1},
	}
	if a != nil && a.Syntax != nil {
		res.Diff.FromFile = a.Syntax.Name
	}
	if b != nil && b.Syntax != nil {
		res.Diff.ToFile = b.Syntax.Name
	}

	return res
}
