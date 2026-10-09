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

// CheckTidyReq is a request to check that workspace modules are tidy.
type CheckTidyReq struct {
	// GoWork is the path to go.work. When empty, the first go.work in the
	// current directory or its parents is used.
	GoWork string
	// Modules are module directories relative to the workspace root. When
	// empty, every module in go.work is checked.
	Modules []string
}

// ModuleTidyResult is the tidy check result for one module.
type ModuleTidyResult struct {
	Dir  string `json:"dir"`
	Path string `json:"path"`
	Tidy bool   `json:"tidy"`
	// Diff is the change 'go mod tidy' would make.
	Diff string `json:"diff,omitempty"`
	// Error is set when tidy failed without producing a diff.
	Error string `json:"error,omitempty"`
}

// CheckTidyRes is the response to a CheckTidyReq.
type CheckTidyRes struct {
	Modules []*ModuleTidyResult `json:"modules"`
}

// Run runs 'go mod tidy -diff' with GOWORK=off in each selected module.
func (r *CheckTidyReq) Run(ctx context.Context) (*CheckTidyRes, error) {
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

	res := &CheckTidyRes{Modules: []*ModuleTidyResult{}}
	for _, m := range modules {
		result := &ModuleTidyResult{Dir: m.Dir, Path: m.Path}
		out, err := runGo(ctx, ws.AbsDir(m), nil, "mod", "tidy", "-diff")
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}

		// tidy -diff exits non-zero and prints a diff when the module isn't
		// tidy. A failure without a diff means tidy itself failed.
		switch {
		case err == nil:
			result.Tidy = true
		case len(out.Stdout) > 0:
			result.Diff = string(out.Stdout)
		case len(out.Stderr) > 0:
			result.Error = strings.TrimSpace(string(out.Stderr))
		default:
			result.Error = err.Error()
		}

		res.Modules = append(res.Modules, result)
	}

	return res, nil
}

// Success reports whether every module is tidy.
func (r *CheckTidyRes) Success() bool {
	if r == nil {
		return false
	}

	for _, m := range r.Modules {
		if !m.Tidy {
			return false
		}
	}

	return true
}

// String summarizes the result.
func (r *CheckTidyRes) String() string {
	if r == nil {
		return "no result"
	}

	var untidy, failed int
	for _, m := range r.Modules {
		switch {
		case m.Error != "":
			failed++
		case !m.Tidy:
			untidy++
		}
	}

	return fmt.Sprintf("%d of %d modules not tidy, %d failed",
		untidy, len(r.Modules), failed)
}

// ToJSON marshals the result to JSON.
func (r *CheckTidyRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling tidy check to JSON: %w", err)
	}

	return b, nil
}

// ToTable renders one row per module as a text table.
func (r *CheckTidyRes) ToTable() string {
	return r.table(true).Render()
}

// ToMarkdown renders one row per module as a markdown table, followed by the
// diff of each module that isn't tidy. A diff in a markdown table cell loses
// its line breaks and indentation, so diffs get their own code blocks.
func (r *CheckTidyRes) ToMarkdown() string {
	out := r.table(false).RenderMarkdown()
	if r == nil {
		return out
	}

	for _, m := range r.Modules {
		if m.Tidy || m.Error != "" {
			continue
		}
		out += fmt.Sprintf("\n\n`%s`:\n\n```diff\n%s\n```", m.Dir, strings.TrimRight(m.Diff, "\n"))
	}

	return out
}

// table renders one row per module. The last column holds tidy's error, plus
// the diff when withDiff is true.
func (r *CheckTidyRes) table(withDiff bool) table.Writer {
	t := newTable()
	t.SetTitle("Go module tidy check")
	if withDiff {
		t.AppendHeader(table.Row{"Module", "Result", "Diff"})
	} else {
		t.AppendHeader(table.Row{"Module", "Result", "Error"})
	}
	if r == nil {
		return t
	}

	for _, m := range r.Modules {
		switch {
		case m.Tidy:
			t.AppendRow(table.Row{m.Dir, "tidy", ""})
		case m.Error != "":
			t.AppendRow(table.Row{m.Dir, "error", m.Error})
		case withDiff:
			t.AppendRow(table.Row{m.Dir, "not tidy", strings.TrimRight(m.Diff, "\n")})
		default:
			t.AppendRow(table.Row{m.Dir, "not tidy", ""})
		}
	}

	return t
}
