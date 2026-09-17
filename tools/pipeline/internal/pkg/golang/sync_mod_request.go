// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"slices"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/mod/modfile"
)

// SyncAction records what happened to a module entry during a sync.
type SyncAction string

const (
	// SyncActionUpdated means the version changed.
	SyncActionUpdated SyncAction = "updated"
	// SyncActionAdded means the entry was absent in dest and we added it.
	SyncActionAdded SyncAction = "added"
)

// SyncModReq is a request to sync a destination go.mod with a source go.mod,
// applying A's directives to B in-place.
type SyncModReq struct {
	A        *ModSource // source (read-only)
	B        *ModSource // dest (read and written; B.Name is the on-disk file path)
	DiffOpts *DiffOpts  // options governing which directives to compare and sync
	DryRun   bool       // compute changes but skip writing to disk
}

// SyncModRes is the result of a single go.mod sync operation.
type SyncModRes struct {
	Path    string        `json:"path"`
	Changes []*SyncChange `json:"changes,omitempty"`
	Err     error         `json:"-"`
	Error   string        `json:"error,omitempty"`
}

// SyncChange is one module version change applied to a go.mod during a sync.
type SyncChange struct {
	Directive  Directive  `json:"directive"`
	Module     string     `json:"module"`
	OldVersion string     `json:"old_version,omitempty"`
	NewVersion string     `json:"new_version"`
	Action     SyncAction `json:"action"`
}

// DefaultSyncDiffOpts returns DiffOpts with the defaults used by all three sync commands.
func DefaultSyncDiffOpts() *DiffOpts {
	return &DiffOpts{
		ParseLax:          false,
		Require:           true,
		Replace:           true,
		Go:                false,
		Toolchain:         false,
		StrictDiffRequire: true,
		StrictDiffReplace: true,
	}
}

// Run applies A's directives to B in-place. It computes the diff, mutates B's
// modfile AST, and writes the result to disk unless DryRun is set.
func (r *SyncModReq) Run(ctx context.Context) (*SyncModRes, error) {
	if r.A == nil {
		return nil, errors.New("sync mod: source (A) is nil")
	}
	if r.B == nil {
		return nil, errors.New("sync mod: dest (B) is nil")
	}

	res := &SyncModRes{
		Path: r.B.Name,
	}

	if r.DiffOpts == nil {
		r.DiffOpts = DefaultSyncDiffOpts()
	}

	// Gate strict flags behind their directive enable flags so the invariant that
	// disabling a directive skips it entirely is preserved across diff and apply.
	diffOpts := *r.DiffOpts
	diffOpts.StrictDiffRequire = r.DiffOpts.Require && r.DiffOpts.StrictDiffRequire
	diffOpts.StrictDiffReplace = r.DiffOpts.Replace && r.DiffOpts.StrictDiffReplace

	// Compute what differs between A and B before mutating anything.
	modDiff, af, bf, err := DiffModFiles(r.A, r.B, &diffOpts)
	if err != nil {
		res.Err = fmt.Errorf("diffing %s and %s: %w", r.A.Name, r.B.Name, err)
		res.Error = res.Err.Error()
		return res, res.Err
	}

	// Nothing to sync.
	if len(modDiff) == 0 {
		return res, nil
	}

	// Walk the diff and apply mutations to bf. Each apply function iterates all of
	// af's entries for the directive, so we only need to invoke it once per
	// directive — use a seen set to deduplicate.
	seen := make(map[Directive]struct{})
	for _, d := range modDiff {
		if d == nil {
			continue
		}
		if _, ok := seen[d.Directive]; ok {
			continue
		}
		seen[d.Directive] = struct{}{}

		switch d.Directive {
		case DirectiveRequire:
			changes, err := applyRequireDiff(ctx, af, bf, &diffOpts)
			if err != nil {
				res.Err = err
				res.Error = err.Error()
				return res, err
			}
			res.Changes = append(res.Changes, changes...)

		case DirectiveReplace:
			changes, err := applyReplaceDiff(ctx, af, bf, &diffOpts)
			if err != nil {
				res.Err = err
				res.Error = err.Error()
				return res, err
			}
			res.Changes = append(res.Changes, changes...)

		case DirectiveGo:
			c, err := applyGoSync(ctx, af, bf)
			if err != nil {
				res.Err = err
				res.Error = err.Error()
				return res, err
			}
			if c != nil {
				res.Changes = append(res.Changes, c)
			}

		case DirectiveToolchain:
			c, err := applyToolchainSync(ctx, af, bf)
			if err != nil {
				res.Err = err
				res.Error = err.Error()
				return res, err
			}
			if c != nil {
				res.Changes = append(res.Changes, c)
			}
		}
	}

	// Write the updated file back to disk.
	if len(res.Changes) > 0 {
		if r.DryRun {
			slog.Default().DebugContext(
				ctx,
				"dry run: skipping write to disk",
				slog.String("path", r.B.Name),
				slog.Int("changes", len(res.Changes)),
			)
			return res, nil
		}

		formatted, fmtErr := bf.Format()
		if fmtErr != nil {
			res.Err = fmt.Errorf("formatting dest %s: %w", r.B.Name, fmtErr)
			res.Error = res.Err.Error()
			return res, res.Err
		}

		slog.Default().DebugContext(
			ctx,
			"writing updated go.mod",
			slog.String("path", r.B.Name),
			slog.Int("changes", len(res.Changes)),
		)
		if writeErr := os.WriteFile(r.B.Name, formatted, 0o644); writeErr != nil {
			res.Err = fmt.Errorf("writing dest %s: %w", r.B.Name, writeErr)
			res.Error = res.Err.Error()
			return res, res.Err
		}
	}

	return res, nil
}

// applyRequireDiff walks af's require entries and applies them to bf.
// Returns the list of changes made and any accumulated errors.
func applyRequireDiff(ctx context.Context, af, bf *modfile.File, opts *DiffOpts) ([]*SyncChange, error) {
	var changes []*SyncChange
	var errs []error

	for _, aReq := range af.Require {
		if aReq == nil {
			continue
		}

		if matchesAnyExcludePattern(opts.ExcludeRequire, aReq.Mod.Path) {
			continue
		}

		// Look for this module in bf.
		idx := slices.IndexFunc(bf.Require, func(bReq *modfile.Require) bool {
			return bReq != nil && bReq.Mod.Path == aReq.Mod.Path
		})

		if idx < 0 {
			// Not in bf — add it only when strict mode is enabled.
			if !opts.StrictDiffRequire {
				continue
			}
			slog.Default().DebugContext(ctx, "adding require",
				slog.String("module", aReq.Mod.Path),
				slog.String("version", aReq.Mod.Version),
			)
			bf.AddNewRequire(aReq.Mod.Path, aReq.Mod.Version, aReq.Indirect)
			changes = append(changes, &SyncChange{
				Directive:  DirectiveRequire,
				Module:     aReq.Mod.Path,
				NewVersion: aReq.Mod.Version,
				Action:     SyncActionAdded,
			})
			continue
		}

		bReq := bf.Require[idx]
		if bReq.Mod.Version == aReq.Mod.Version {
			continue
		}

		oldVersion := bReq.Mod.Version
		slog.Default().DebugContext(ctx, "updating require",
			slog.String("module", aReq.Mod.Path),
			slog.String("old", oldVersion),
			slog.String("new", aReq.Mod.Version),
		)
		if err := bf.AddRequire(aReq.Mod.Path, aReq.Mod.Version); err != nil {
			errs = append(errs, fmt.Errorf("sync mod: add require %s: %w", aReq.Mod.Path, err))
			continue
		}
		changes = append(changes, &SyncChange{
			Directive:  DirectiveRequire,
			Module:     aReq.Mod.Path,
			OldVersion: oldVersion,
			NewVersion: aReq.Mod.Version,
			Action:     SyncActionUpdated,
		})
	}

	return changes, errors.Join(errs...)
}

// applyReplaceDiff walks af's replace entries and applies them to bf.
// Returns the list of changes made and accumulated errors, failing last.
func applyReplaceDiff(ctx context.Context, af, bf *modfile.File, opts *DiffOpts) ([]*SyncChange, error) {
	var changes []*SyncChange
	var errs []error

	for _, aRep := range af.Replace {
		if aRep == nil {
			continue
		}

		if matchesAnyExcludePattern(opts.ExcludeReplace, aRep.Old.Path, aRep.New.Path) {
			continue
		}

		// Match on old path AND old version — both must agree.
		idx := slices.IndexFunc(bf.Replace, func(bRep *modfile.Replace) bool {
			return bRep != nil &&
				bRep.Old.Path == aRep.Old.Path &&
				bRep.Old.Version == aRep.Old.Version
		})

		if idx < 0 {
			// Not in bf — add it only when strict mode is enabled.
			if !opts.StrictDiffReplace {
				continue
			}
			if err := bf.AddReplace(
				aRep.Old.Path, aRep.Old.Version,
				aRep.New.Path, aRep.New.Version,
			); err != nil {
				errs = append(errs, fmt.Errorf("sync mod: add replace %s: %w", aRep.Old.Path, err))
				continue
			}
			changes = append(changes, &SyncChange{
				Directive:  DirectiveReplace,
				Module:     aRep.Old.Path,
				NewVersion: aRep.New.Version,
				Action:     SyncActionAdded,
			})
			continue
		}

		bRep := bf.Replace[idx]
		if replaceEqual(aRep, bRep) {
			continue
		}

		oldVersion := bRep.New.Version

		// Add the new entry first so we're never left without a replace in place.
		if err := bf.AddReplace(
			aRep.Old.Path, aRep.Old.Version,
			aRep.New.Path, aRep.New.Version,
		); err != nil {
			errs = append(errs, fmt.Errorf("sync mod: add replace %s: %w", aRep.Old.Path, err))
			continue
		}
		if err := bf.DropReplace(bRep.Old.Path, bRep.Old.Version); err != nil {
			// AddReplace already succeeded; log and move on. modfile.Format()
			// will deduplicate any leftover entry.
			slog.Default().DebugContext(ctx, "drop replace after add",
				slog.String("path", bRep.Old.Path),
				slog.String("error", err.Error()),
			)
		}
		changes = append(changes, &SyncChange{
			Directive:  DirectiveReplace,
			Module:     aRep.Old.Path,
			OldVersion: oldVersion,
			NewVersion: aRep.New.Version,
			Action:     SyncActionUpdated,
		})
	}

	return changes, errors.Join(errs...)
}

// applyGoSync syncs the go directive from a into b.
func applyGoSync(ctx context.Context, a, b *modfile.File) (*SyncChange, error) {
	if a.Go == nil {
		return nil, nil
	}

	oldVersion := ""
	if b.Go != nil {
		if b.Go.Version == a.Go.Version {
			return nil, nil
		}
		oldVersion = b.Go.Version
	}

	slog.Default().DebugContext(ctx, "syncing go directive", slog.String("version", a.Go.Version))

	if err := b.AddGoStmt(a.Go.Version); err != nil {
		return nil, fmt.Errorf("sync mod: add go stmt %s: %w", a.Go.Version, err)
	}

	action := SyncActionUpdated
	if oldVersion == "" {
		action = SyncActionAdded
	}

	return &SyncChange{
		Directive:  DirectiveGo,
		Module:     "go",
		OldVersion: oldVersion,
		NewVersion: a.Go.Version,
		Action:     action,
	}, nil
}

// applyToolchainSync syncs the toolchain directive from a into b.
func applyToolchainSync(ctx context.Context, a, b *modfile.File) (*SyncChange, error) {
	if a.Toolchain == nil {
		return nil, nil
	}

	oldName := ""
	if b.Toolchain != nil {
		if b.Toolchain.Name == a.Toolchain.Name {
			return nil, nil
		}
		oldName = b.Toolchain.Name
	}

	slog.Default().DebugContext(ctx, "syncing toolchain directive", slog.String("name", a.Toolchain.Name))

	if err := b.AddToolchainStmt(a.Toolchain.Name); err != nil {
		return nil, fmt.Errorf("sync mod: add toolchain stmt %s: %w", a.Toolchain.Name, err)
	}

	action := SyncActionUpdated
	if oldName == "" {
		action = SyncActionAdded
	}

	return &SyncChange{
		Directive:  DirectiveToolchain,
		Module:     "toolchain",
		OldVersion: oldName,
		NewVersion: a.Toolchain.Name,
		Action:     action,
	}, nil
}

// matchesExcludePattern reports whether pattern matches p using path.Match.
// We use path.Match (not filepath.Match) because Go module paths always use /
// as a separator. If path.Match returns an error we fall back to literal equality.
func matchesExcludePattern(pattern, p string) bool {
	matched, err := path.Match(pattern, p)
	if err != nil {
		return pattern == p
	}
	return matched
}

// matchesAnyExcludePattern reports whether any of the patterns matches any of
// the provided paths.
func matchesAnyExcludePattern(patterns []string, paths ...string) bool {
	for _, pat := range patterns {
		for _, p := range paths {
			if matchesExcludePattern(pat, p) {
				return true
			}
		}
	}
	return false
}

// newSyncResultTable returns a borderless table.Writer used by ToTable and ToMarkdown.
func newSyncResultTable() table.Writer {
	t := table.NewWriter()
	t.Style().Options.DrawBorder = false
	t.Style().Options.SeparateColumns = false
	t.Style().Options.SeparateFooter = false
	t.Style().Options.SeparateHeader = false
	t.Style().Options.SeparateRows = false
	t.SuppressEmptyColumns()
	t.SuppressTrailingSpaces()
	return t
}

// ToTable renders the result as a go-pretty table. When changes are present the
// columns are directive | action | from | to | module. When there are no changes
// (or on error) the columns collapse to changes | error.
func (res *SyncModRes) TableWriter() table.Writer {
	t := newSyncResultTable()
	if res == nil {
		return t
	}

	if len(res.Changes) > 0 {
		t.AppendHeader(table.Row{"directive", "action", "from", "to", "module"})
		for _, c := range res.Changes {
			if c == nil {
				continue
			}
			t.AppendRow(table.Row{c.Directive, c.Action, c.OldVersion, c.NewVersion, c.Module})
		}
	} else {
		t.AppendHeader(table.Row{"changes", "error"})
		errStr := ""
		if res.Error != "" {
			errStr = res.Error
		}
		t.AppendRow(table.Row{len(res.Changes), errStr})
	}

	return t
}

func (res *SyncModRes) ToTable() string {
	if res == nil {
		return ""
	}
	return res.TableWriter().Render()
}

// ToMarkdown renders the result as a Markdown table using the same column layout
// as ToTable.
func (res *SyncModRes) ToMarkdown() string {
	if res == nil {
		return ""
	}
	return res.TableWriter().RenderMarkdown()
}

// ToJSON marshals the result to JSON.
func (res *SyncModRes) ToJSON() ([]byte, error) {
	if res == nil {
		return nil, errors.New("sync mod result is nil")
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("marshaling sync mod result to JSON: %w", err)
	}
	return b, nil
}
