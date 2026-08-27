// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package ebom

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/xuri/excelize/v2"
)

// ConvertEBOMToCSVReq is the request to convert an IBM EBOM .xlsx workbook to CSV files.
type ConvertEBOMToCSVReq struct {
	XLSXPath    string
	OutDir      string
	RepoRoot    string // git repository root; used to derive the default OutDir when OutDir is empty
	PIDOverride string
	// Slugs holds raw "idx:slug" override strings in N:slug format where N is
	// the 0-based index across all sheets in the workbook (matching inspect output).
	Slugs []string
}

// ConvertEBOMToCSVRes is the response from converting an EBOM workbook to CSV.
type ConvertEBOMToCSVRes struct {
	Sheets []*EBOMCSV `json:"sheets"`
}

// EBOMCSV is the EBOM converted to a compatible CSV. Non-eAssembly sheets are
// appended to every output file.
type EBOMCSV struct {
	SheetName  string `json:"sheet_name,omitempty"`
	PID        string `json:"pid,omitempty"`
	Slug       string `json:"slug,omitempty"`
	OutputPath string `json:"output_path,omitempty"` // empty for non-eAssembly sheets
	AppendedTo string `json:"appended_to,omitempty"` // output path this sheet was appended to; non-empty for non-eAssembly sheets
	RowCount   int    `json:"row_count,omitempty"`
}

// ebomSheet is an internal wrapper for a single sheet read from the workbook.
type ebomSheet struct {
	name      string
	rows      [][]string
	sheetType string // trimmed first cell of row 0; non-empty only for eAssembly sheets (e.g. "eAssemblies")
}

// eAssembliesMarker is matched against the first cell of the first row to
// detect eAssembly sheets.
const eAssembliesMarker = "eAssemblies"

// versionRe matches version tokens such as V1.20.4, V2.1.0, $2.1.0.
var versionRe = regexp.MustCompile(`(?i)^(\$?v\d[\d.]*|\$\d[\d.]*)$`)

// slugRe validates the slug portion of a --slugs value. Dots are explicitly
// rejected — they are structurally ambiguous against the separator in
// <pid>.<slug>.csv.
var slugRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

// slugInferenceNoiseWords are stripped individually when inferring a slug.
// Includes "ibm", known product names, and "self-managed" so that inference
// works regardless of whether those words appear in the description.
var slugInferenceNoiseWords = map[string]struct{}{
	"ibm":           {},
	"vault":         {},
	"nomad":         {},
	"self-managed":  {},
	"multiplatform": {},
	"english":       {},
	"multi":         {},
	"data":          {},
	"center":        {},
	"install":       {},
	"entitlement":   {},
}

// excelErrors is the set of error strings Excel emits for broken formulas.
// These are written literally by excelize when the formula cannot be resolved
// and should be treated as empty cells in the output CSV.
var excelErrors = map[string]struct{}{
	"#REF!":   {},
	"#NAME?":  {},
	"#VALUE!": {},
	"#DIV/0!": {},
	"#NULL!":  {},
	"#N/A":    {},
	"#NUM!":   {},
}

// Run converts each eAssembly sheet in the workbook to a CSV file.
// Non-eAssembly sheets are appended to every output file in the order they
// appear in the workbook.
func (r *ConvertEBOMToCSVReq) Run(ctx context.Context) (_ *ConvertEBOMToCSVRes, err error) {
	slugOverrides, err := r.validate()
	if err != nil {
		return nil, err
	}

	outDir := r.resolveOutDir()
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating output directory %q: %w", outDir, err)
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

	var sheets []*ebomSheet
	for _, name := range f.GetSheetList() {
		raw, err := f.GetRows(name)
		if err != nil {
			return nil, fmt.Errorf("reading sheet %q: %w", name, err)
		}
		s := &ebomSheet{name: name}
		s.processRows(raw)
		sheets = append(sheets, s)
	}

	var eaCount int
	for _, s := range sheets {
		if s.sheetType != "" {
			eaCount++
		}
	}

	if eaCount == 0 {
		return nil, errors.New("no eAssembly sheets found in workbook")
	}

	isMultiSheet := eaCount > 1
	res := &ConvertEBOMToCSVRes{}

	for idx, ea := range sheets {
		if ea.sheetType == "" {
			continue
		}

		slog.Default().DebugContext(ctx, "processing eAssembly sheet",
			slog.String("sheet", ea.name),
			slog.Int("index", idx))

		pid := r.PIDOverride
		if pid == "" {
			for _, row := range ea.rows {
				if len(row) >= 2 && strings.EqualFold(strings.TrimSpace(row[0]), "pid") {
					pid = strings.TrimSpace(row[1])
					break
				}
			}
		}
		if pid == "" {
			return nil, fmt.Errorf("PID row not found in sheet %q", ea.name)
		}

		var slug string
		if isMultiSheet {
			if override, ok := slugOverrides[idx]; ok {
				slug = override
			} else {
				var desc string
				for _, row := range ea.rows {
					if len(row) >= 4 && strings.EqualFold(strings.TrimSpace(row[3]), "eAssembly") {
						desc = strings.TrimSpace(row[1])
						break
					}
				}
				if desc == "" {
					return nil, fmt.Errorf(
						"eAssembly row not found in sheet %q (index %d); use --slugs %d:<slug>",
						ea.name, idx, idx,
					)
				}
				var ok bool
				slug, ok = inferSlug(desc)
				if !ok {
					return nil, fmt.Errorf(
						"could not infer slug for sheet %q (index %d); use --slugs %d:<slug>",
						ea.name, idx, idx,
					)
				}
			}
		}

		var outputPath string
		if slug == "" {
			outputPath = filepath.Join(outDir, pid+".csv")
		} else {
			outputPath = filepath.Join(outDir, pid+"."+slug+".csv")
		}

		// Build combined sections: the eAssembly sheet rows form the first
		// section. Each non-eAssembly sheet (in workbook order) becomes its
		// own section. Each section is padded independently to its own column
		// width so that the compact pBom and reference-code sections are not
		// widened to match the eAssembly section.
		sections := [][][]string{ea.rows}
		for _, other := range sheets {
			if other.sheetType != "" || len(other.rows) == 0 {
				continue
			}
			sections = append(sections, other.rows)
		}

		totalRows := 0
		for _, sec := range sections {
			totalRows += len(sec)
		}

		slog.Default().DebugContext(ctx, "writing EBOM CSV",
			slog.String("path", outputPath),
			slog.Int("rows", totalRows))
		if err := writeCSVSections(outputPath, sections); err != nil {
			return nil, err
		}

		res.Sheets = append(res.Sheets, &EBOMCSV{
			SheetName:  ea.name,
			PID:        pid,
			Slug:       slug,
			OutputPath: outputPath,
			RowCount:   totalRows,
		})
	}

	// Record non-eAssembly sheets in the response so the table output can show
	// where they ended up. Each non-eAssembly sheet is appended to every output
	// file, so AppendedTo lists all output paths.
	for _, s := range sheets {
		if s.sheetType != "" || len(s.rows) == 0 {
			continue
		}
		for _, out := range res.Sheets {
			if out.OutputPath == "" {
				continue
			}
			res.Sheets = append(res.Sheets, &EBOMCSV{
				SheetName:  s.name,
				AppendedTo: out.OutputPath,
				RowCount:   len(s.rows),
			})
		}
	}

	return res, nil
}

// ToJSON marshals the response to JSON.
func (r *ConvertEBOMToCSVRes) ToJSON() ([]byte, error) {
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
func (r *ConvertEBOMToCSVRes) ToTable() string {
	if r == nil {
		return ""
	}
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateHeader = false
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	t.AppendHeader(table.Row{"SHEET", "PRODUCT ID", "SLUG", "OUTPUT", "ROWS"})
	for _, s := range r.Sheets {
		if s.AppendedTo != "" {
			t.AppendRow(table.Row{s.SheetName, "", "", "(appended to " + filepath.Base(s.AppendedTo) + ")", s.RowCount})
			continue
		}
		t.AppendRow(table.Row{s.SheetName, s.PID, s.Slug, s.OutputPath, s.RowCount})
	}
	return t.Render()
}

// ToMarkdown renders the response as a Markdown table.
func (r *ConvertEBOMToCSVRes) ToMarkdown() string {
	if r == nil {
		return ""
	}
	t := table.NewWriter()
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()

	t.AppendHeader(table.Row{"SHEET", "PRODUCT ID", "SLUG", "OUTPUT", "ROWS"})
	for _, s := range r.Sheets {
		if s.AppendedTo != "" {
			t.AppendRow(table.Row{s.SheetName, "", "", "(appended to " + filepath.Base(s.AppendedTo) + ")", s.RowCount})
			continue
		}
		t.AppendRow(table.Row{s.SheetName, s.PID, s.Slug, s.OutputPath, s.RowCount})
	}
	return t.RenderMarkdown()
}

// validate checks the request fields and returns the parsed slug overrides.
// Slug parsing is done first so flag errors are reported without requiring the
// file to exist.
func (r *ConvertEBOMToCSVReq) validate() (map[int]string, error) {
	if r == nil {
		return nil, errors.New("uninitialized request")
	}
	if r.XLSXPath == "" {
		return nil, errors.New("xlsx path is required")
	}
	slugOverrides, err := parseSlugs(r.Slugs)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(filepath.Ext(r.XLSXPath), ".xlsx") {
		return nil, errors.New("input file is not .xlsx; resave from Excel or LibreOffice as .xlsx before running this command")
	}
	if _, err := os.Stat(r.XLSXPath); err != nil {
		return nil, fmt.Errorf("xlsx file %q: %w", r.XLSXPath, err)
	}
	return slugOverrides, nil
}

// processRows trims cells, clears Excel error strings (e.g. "#REF!"), strips
// leading and trailing blank rows, preserves internal blank rows, and detects
// the sheet type. sheetType is set to the trimmed first cell of the first
// non-blank row when that cell matches eAssembliesMarker (case-insensitive).
func (s *ebomSheet) processRows(raw [][]string) {
	trimmed := make([][]string, 0, len(raw))
	for _, row := range raw {
		t := make([]string, len(row))
		for i, cell := range row {
			v := strings.TrimSpace(cell)
			if _, ok := excelErrors[v]; ok {
				v = ""
			}
			t[i] = v
		}
		trimmed = append(trimmed, t)
	}

	isBlank := func(row []string) bool {
		for _, cell := range row {
			if cell != "" {
				return false
			}
		}
		return true
	}

	// Strip leading blank rows.
	for len(trimmed) > 0 && isBlank(trimmed[0]) {
		trimmed = trimmed[1:]
	}
	// Strip trailing blank rows.
	for len(trimmed) > 0 && isBlank(trimmed[len(trimmed)-1]) {
		trimmed = trimmed[:len(trimmed)-1]
	}

	s.rows = trimmed
	if len(s.rows) > 0 && len(s.rows[0]) > 0 {
		if strings.EqualFold(s.rows[0][0], eAssembliesMarker) {
			s.sheetType = s.rows[0][0]
		}
	}
}

// parseSlugs converts a slice of "N:slug" strings into a map[int]string.
// Returns an error for any malformed entry.
func parseSlugs(vals []string) (map[int]string, error) {
	m := make(map[int]string, len(vals))
	for _, v := range vals {
		idx := strings.IndexByte(v, ':')
		if idx <= 0 {
			return nil, fmt.Errorf("invalid --slugs value %q: expected N:slug", v)
		}
		n, convErr := strconv.Atoi(v[:idx])
		if convErr != nil || n < 0 {
			return nil, fmt.Errorf("invalid --slugs value %q: N must be a non-negative integer", v)
		}
		slug := v[idx+1:]
		if !slugRe.MatchString(slug) {
			return nil, fmt.Errorf(
				"invalid --slugs value %q: slug %q must match ^[a-z0-9][a-z0-9-]*$ (lowercase, hyphens only — dots are not permitted)",
				v, slug,
			)
		}
		m[n] = slug
	}
	return m, nil
}

func (r *ConvertEBOMToCSVReq) resolveOutDir() string {
	if r.OutDir != "" {
		return r.OutDir
	}
	if r.RepoRoot == "" {
		slog.Default().Debug("repo root not set; writing EBOM CSVs to current directory")
		return "."
	}
	return filepath.Join(r.RepoRoot, ".release", "ibm-pao", "eboms")
}

// inferSlug derives the output filename slug from an eAssembly description by
// stripping noise words (ibm, product name, self-managed), version tokens, and
// platform/language words, then joining the remaining tokens with hyphens.
// Returns ("", false) when the result is empty — legitimate for single-sheet
// workbooks; an inference failure for multi-sheet workbooks.
func inferSlug(desc string) (string, bool) {
	var tokens []string
	for word := range strings.FieldsSeq(desc) {
		lower := strings.ToLower(word)
		if _, noise := slugInferenceNoiseWords[lower]; noise {
			continue
		}
		if versionRe.MatchString(word) {
			continue
		}
		tokens = append(tokens, lower)
	}

	if len(tokens) == 0 {
		return "", false
	}
	return strings.Join(tokens, "-"), true
}

// writeCSVSections writes sections to outPath contiguously. Each section is
// padded independently to its own maximum column count so that narrow sections
// (pBom, reference-code tables) are not artificially widened to match the wide
// eAssembly section.
func writeCSVSections(outPath string, sections [][][]string) (err error) {
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating output file %q: %w", outPath, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	w := csv.NewWriter(f)
	for _, rows := range sections {
		// Pad every row within this section to the section's own max column count.
		maxCols := 0
		for _, row := range rows {
			if len(row) > maxCols {
				maxCols = len(row)
			}
		}
		pad := make([]string, maxCols)

		for _, row := range rows {
			out := row
			if len(row) < maxCols {
				out = append(row[:len(row):len(row)], pad[len(row):]...)
			}
			if err := w.Write(out); err != nil {
				return fmt.Errorf("writing CSV row to %q: %w", outPath, err)
			}
		}
	}
	w.Flush()
	return w.Error()
}
