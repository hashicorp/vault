// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package ebom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/xuri/excelize/v2"
)

// InspectEBOMReq is the request to inspect eAssembly sheets in an EBOM .xlsx workbook.
type InspectEBOMReq struct {
	XLSXPath    string // path to the input .xlsx file (required)
	PIDOverride string // if non-empty, override the PID shown in output
}

// InspectEBOMRes is the response from inspecting an EBOM workbook.
type InspectEBOMRes struct {
	Sheets []*InspectEBOMSheet `json:"sheets"`
}

// InspectEBOMSheet is the result for a single eAssembly sheet.
type InspectEBOMSheet struct {
	Index        int    `json:"index"` // 0-based index across all sheets in the workbook (matches --slugs N)
	SheetName    string `json:"sheetName"`
	PID          string `json:"pid"`
	Description  string `json:"description"`  // eAssembly description; empty if row not found
	InferredSlug string `json:"inferredSlug"` // empty string if inference fails
	SlugOverride string `json:"slugOverride"` // "--slugs N:<slug>" when inference fails; empty when inferred
}

// Run inspects each eAssembly sheet and returns metadata without erroring on
// slug inference failures.
func (r *InspectEBOMReq) Run(ctx context.Context) (_ *InspectEBOMRes, err error) {
	if err = r.validate(); err != nil {
		return nil, err
	}

	slog.Default().DebugContext(ctx, "opening EBOM", slog.String("path", r.XLSXPath))
	var f *excelize.File
	f, err = excelize.OpenFile(r.XLSXPath)
	if err != nil {
		return nil, fmt.Errorf("opening %q: %w", r.XLSXPath, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	res := &InspectEBOMRes{}

	for idx, name := range f.GetSheetList() {
		raw, err := f.GetRows(name)
		if err != nil {
			return nil, fmt.Errorf("reading sheet %q: %w", name, err)
		}
		s := &ebomSheet{name: name}
		s.processRows(raw)
		if s.sheetType == "" {
			continue
		}

		pid := r.PIDOverride
		if pid == "" {
			for _, row := range s.rows {
				if len(row) >= 2 && strings.EqualFold(strings.TrimSpace(row[0]), "pid") {
					pid = strings.TrimSpace(row[1])
					break
				}
			}
		}

		var desc string
		for _, row := range s.rows {
			if len(row) >= 4 && strings.EqualFold(strings.TrimSpace(row[3]), "eAssembly") {
				desc = strings.TrimSpace(row[1])
				break
			}
		}

		inferredSlug, _ := inferSlug(desc)

		var slugOverride string
		if inferredSlug == "" {
			slugOverride = fmt.Sprintf("--slugs %d:<slug>", idx)
		}

		res.Sheets = append(res.Sheets, &InspectEBOMSheet{
			Index:        idx,
			SheetName:    name,
			PID:          pid,
			Description:  desc,
			InferredSlug: inferredSlug,
			SlugOverride: slugOverride,
		})
	}

	return res, nil
}

// ToJSON marshals the response to JSON.
func (r *InspectEBOMRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized response")
	}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling to JSON: %w", err)
	}
	return b, nil
}

// ToTable renders the response as a plain text table.
//
// For single-sheet workbooks the index/sheet/slug columns are omitted — there
// is nothing to disambiguate and --slugs is never needed.  For multi-sheet
// workbooks all columns are shown; the SLUG OVERRIDE column is suppressed when
// every slug was successfully inferred.
func (r *InspectEBOMRes) ToTable() string {
	if r == nil {
		return ""
	}
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateHeader = false
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	if len(r.Sheets) == 1 {
		t.AppendHeader(table.Row{"PRODUCT ID", "DESCRIPTION"})
		s := r.Sheets[0]
		t.AppendRow(table.Row{s.PID, s.Description})
		return t.Render()
	}

	t.AppendHeader(table.Row{"INDEX", "SHEET", "PRODUCT ID", "DESCRIPTION", "INFERRED SLUG", "SLUG OVERRIDE"})
	for _, s := range r.Sheets {
		slug := s.InferredSlug
		if slug == "" {
			slug = "(no match)"
		}
		t.AppendRow(table.Row{s.Index, s.SheetName, s.PID, s.Description, slug, s.SlugOverride})
	}
	return t.Render()
}

// ToMarkdown renders the response as a Markdown table.
//
// Single-sheet workbooks omit the index/sheet/slug columns for the same
// reasons described in ToTable.
func (r *InspectEBOMRes) ToMarkdown() string {
	if r == nil {
		return ""
	}
	t := table.NewWriter()
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	if len(r.Sheets) == 1 {
		t.AppendHeader(table.Row{"PRODUCT ID", "DESCRIPTION"})
		s := r.Sheets[0]
		t.AppendRow(table.Row{s.PID, s.Description})
		return t.RenderMarkdown()
	}

	t.AppendHeader(table.Row{"INDEX", "SHEET", "PRODUCT ID", "DESCRIPTION", "INFERRED SLUG", "SLUG OVERRIDE"})
	for _, s := range r.Sheets {
		slug := s.InferredSlug
		if slug == "" {
			slug = "(no match)"
		}
		t.AppendRow(table.Row{s.Index, s.SheetName, s.PID, s.Description, slug, s.SlugOverride})
	}
	return t.RenderMarkdown()
}

func (r *InspectEBOMReq) validate() error {
	if r == nil {
		return errors.New("uninitialized request")
	}
	if r.XLSXPath == "" {
		return errors.New("xlsx path is required")
	}
	if !strings.EqualFold(filepath.Ext(r.XLSXPath), ".xlsx") {
		return errors.New("input file is not .xlsx; resave from Excel or LibreOffice as .xlsx before running this command")
	}
	if _, err := os.Stat(r.XLSXPath); err != nil {
		return fmt.Errorf("xlsx file %q: %w", r.XLSXPath, err)
	}
	return nil
}
