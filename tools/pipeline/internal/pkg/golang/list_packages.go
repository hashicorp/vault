// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// ListPackagesReq is a request to list the packages in each module of a Go
// workspace.
type ListPackagesReq struct {
	// GoWork is the path to go.work. When empty, the first go.work in the
	// current directory or its parents is used.
	GoWork string
	// Tags are comma- or space-separated build tags passed to go list.
	Tags string
	// Modules are module directories relative to the workspace root. When
	// empty, every module in go.work is listed.
	Modules []string
}

// ListPackagesRes is the response to a ListPackagesReq.
type ListPackagesRes struct {
	List *ModulePackagesList
}

// Run lists the packages of each selected module with GOWORK=off, from the
// module's own directory, so the result matches what that module's tests
// build.
func (r *ListPackagesReq) Run(ctx context.Context) (*ListPackagesRes, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	ws, err := LoadWorkspaceOrFind(ctx, r.GoWork)
	if err != nil {
		return nil, err
	}

	modules, err := ws.SelectModules(r.Modules)
	if err != nil {
		return nil, err
	}

	args := []string{"list"}
	if tags := NormalizeTags(r.Tags); len(tags) > 0 {
		args = append(args, "-tags="+strings.Join(tags, ","))
	}
	args = append(args, "./...")

	res := &ListPackagesRes{List: &ModulePackagesList{Modules: []*ModulePackages{}}}
	for _, m := range modules {
		// Don't use -e: a module that doesn't load must fail rather than drop
		// its packages from the test run.
		out, err := runGo(ctx, ws.AbsDir(m), nil, args...)
		if err != nil {
			return nil, fmt.Errorf("listing packages in module %s: %w\n%s",
				m.Dir, err, strings.TrimSpace(string(out.Stderr)))
		}

		packages := strings.Fields(string(out.Stdout))
		if packages == nil {
			packages = []string{}
		}
		res.List.Modules = append(res.List.Modules, &ModulePackages{
			Path:     m.Path,
			Dir:      m.Dir,
			Packages: packages,
		})
	}

	return res, nil
}

// ToJSON marshals the module packages to JSON.
func (r *ListPackagesRes) ToJSON() ([]byte, error) {
	if r == nil || r.List == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r.List)
	if err != nil {
		return nil, fmt.Errorf("marshaling module packages to JSON: %w", err)
	}

	return b, nil
}

// ToTable renders one row per package as a text table.
func (r *ListPackagesRes) ToTable() string {
	return r.table().Render()
}

// ToMarkdown renders one row per package as a markdown table.
func (r *ListPackagesRes) ToMarkdown() string {
	t := r.table()
	t.SetTitle("Go packages")
	return t.RenderMarkdown()
}

func (r *ListPackagesRes) table() table.Writer {
	t := newTable()
	t.AppendHeader(table.Row{"Module", "Package"})
	if r == nil || r.List == nil {
		return t
	}

	for _, m := range r.List.Modules {
		for _, pkg := range m.Packages {
			t.AppendRow(table.Row{m.Dir, pkg})
		}
	}

	return t
}
