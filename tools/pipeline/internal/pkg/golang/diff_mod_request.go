// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pmezard/go-difflib/difflib"
)

type DiffModReq struct {
	A    *ModSource `json:"a,omitempty"`
	B    *ModSource `json:"b,omitempty"`
	Opts *DiffOpts  `json:"opts,omitempty"`
}

func (r *DiffModReq) Run(ctx context.Context) (*DiffModRes, error) {
	res := &DiffModRes{}
	var err error

	res.ModDiff, err = DiffModFiles(r.A, r.B, r.Opts)

	return res, err
}

type DiffModRes struct {
	ModDiff
}

func (r *DiffModRes) ToJSON() ([]byte, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling latest HCP image response to JSON: %w", err)
	}

	return b, nil
}

// ToTable marshals the response to a text table.
func (r *DiffModRes) ToTable(err error) (table.Writer, error) {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false

	if r == nil || r.ModDiff == nil || len(r.ModDiff) == 0 || err != nil {
		if err != nil {
			t.AppendHeader(table.Row{"error"})
			t.AppendRow(table.Row{err.Error()})
		}
		return t, err
	}

	t.AppendHeader(table.Row{"explanation", "diff"})
	for _, diff := range r.ModDiff {
		if diff == nil {
			continue
		}
		if diff.Diff == nil {
			return nil, fmt.Errorf("missing unified diff: %v", diff)
		}

		diffText, err := difflib.GetUnifiedDiffString(*diff.Diff)
		if err != nil {
			return nil, err
		}
		t.AppendRow(table.Row{diff.Directive.Explanation(), diffText})
	}
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	return t, nil
}
