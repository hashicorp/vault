// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

// ListModulesReq is a request to list the modules in a Go workspace.
type ListModulesReq struct {
	// GoWork is the path to go.work. When empty, the first go.work in the
	// current directory or its parents is used.
	GoWork string
}

// ListModulesRes is the response to a ListModulesReq.
type ListModulesRes struct {
	Workspace *Workspace
}

// Run lists the workspace's modules.
func (r *ListModulesReq) Run(ctx context.Context) (*ListModulesRes, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	ws, err := LoadWorkspaceOrFind(ctx, r.GoWork)
	if err != nil {
		return nil, err
	}

	return &ListModulesRes{Workspace: ws}, nil
}

// ToJSON marshals the workspace to JSON.
func (r *ListModulesRes) ToJSON() ([]byte, error) {
	if r == nil || r.Workspace == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r.Workspace)
	if err != nil {
		return nil, fmt.Errorf("marshaling workspace modules to JSON: %w", err)
	}

	return b, nil
}

// ToTable renders the modules as a text table.
func (r *ListModulesRes) ToTable() string {
	return r.table().Render()
}

// ToMarkdown renders the modules as a markdown table.
func (r *ListModulesRes) ToMarkdown() string {
	t := r.table()
	t.SetTitle("Go workspace modules")
	return t.RenderMarkdown()
}

func (r *ListModulesRes) table() table.Writer {
	t := newTable()
	t.AppendHeader(table.Row{"Path", "Dir"})
	if r == nil || r.Workspace == nil {
		return t
	}

	for _, m := range r.Workspace.Modules {
		t.AppendRow(table.Row{m.Path, m.Dir})
	}

	return t
}

// newTable returns a borderless table writer for command output.
func newTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.SuppressTrailingSpaces()

	return t
}
