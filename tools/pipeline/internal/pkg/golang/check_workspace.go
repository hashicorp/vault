// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"

	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/mod/modfile"
)

// WorkspaceViolationKind is the kind of problem found by a workspace check.
type WorkspaceViolationKind string

const (
	// WorkspaceViolationUnregisteredModule is a committed go.mod whose
	// directory isn't used in go.work.
	WorkspaceViolationUnregisteredModule WorkspaceViolationKind = "unregistered-module"
	// WorkspaceViolationMissingLocalReplace is a workspace module that
	// requires another workspace module without replacing it with that
	// module's directory.
	WorkspaceViolationMissingLocalReplace WorkspaceViolationKind = "missing-local-replace"
	// WorkspaceViolationInvalidLocalReplace is a replace of a workspace module
	// that doesn't point every version at that module's directory.
	WorkspaceViolationInvalidLocalReplace WorkspaceViolationKind = "invalid-local-replace"
	// WorkspaceViolationPublishedReplace is a published module that replaces
	// another workspace module. Go ignores replace directives in imported
	// modules, so the replace hides what importers actually build against.
	WorkspaceViolationPublishedReplace WorkspaceViolationKind = "published-replace"
)

// WorkspaceViolation is a problem found by a workspace check.
type WorkspaceViolation struct {
	Kind WorkspaceViolationKind `json:"kind"`
	// Module is the module directory relative to the workspace root.
	Module string `json:"module"`
	Detail string `json:"detail"`
}

// CheckWorkspaceReq is a request to check that go.work and the committed
// go.mod files agree, that every unpublished workspace module builds against
// the on-branch code of the workspace modules it requires, and that every
// published module builds against the released versions it requires.
type CheckWorkspaceReq struct {
	// GoWork is the path to go.work. When empty, the first go.work in the
	// current directory or its parents is used.
	GoWork string
	// Published are path.Match patterns of module directories, relative to
	// the workspace root, for modules that are published for others to
	// import. Go ignores replace directives outside the main module, so a
	// published module must not replace workspace modules: its GOWORK=off
	// builds would then test code that its importers never get. Each pattern
	// must match at least one module.
	Published []string
}

// CheckWorkspaceRes is the response to a CheckWorkspaceReq.
type CheckWorkspaceRes struct {
	Workspace  *Workspace            `json:"workspace"`
	Violations []*WorkspaceViolation `json:"violations"`
}

// Run checks the workspace.
func (r *CheckWorkspaceReq) Run(ctx context.Context) (*CheckWorkspaceRes, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	ws, err := LoadWorkspaceOrFind(ctx, r.GoWork)
	if err != nil {
		return nil, err
	}

	published, err := publishedModules(ws, r.Published)
	if err != nil {
		return nil, err
	}

	goMods, err := committedGoMods(ctx, ws)
	if err != nil {
		return nil, err
	}

	res := &CheckWorkspaceRes{
		Workspace:  ws,
		Violations: findUnregisteredModules(ws, goMods),
	}

	for _, m := range ws.Modules {
		goModPath := filepath.Join(ws.AbsDir(m), "go.mod")
		slog.Default().DebugContext(ctx, "checking local replaces", slog.String("go_mod", goModPath))
		data, err := os.ReadFile(goModPath)
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", goModPath, err)
		}

		// ParseLax drops replace directives, so parse strictly.
		modFile, err := modfile.Parse(goModPath, data, nil)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", goModPath, err)
		}

		if _, ok := published[m]; ok {
			res.Violations = append(res.Violations, findPublishedReplaceViolations(ws, m, modFile)...)
			continue
		}
		res.Violations = append(res.Violations, findReplaceViolations(ws, m, modFile)...)
	}

	return res, nil
}

// publishedModules returns the workspace modules whose directories match any
// of patterns. It's an error for a pattern to be malformed or to match no
// module, so that a typo can't silently exempt nothing.
func publishedModules(ws *Workspace, patterns []string) (map[*WorkspaceModule]struct{}, error) {
	published := map[*WorkspaceModule]struct{}{}
	for _, pattern := range patterns {
		pattern = path.Clean(filepath.ToSlash(pattern))
		matched := false
		for _, m := range ws.Modules {
			ok, err := path.Match(pattern, m.Dir)
			if err != nil {
				return nil, fmt.Errorf("invalid published module pattern %q: %w", pattern, err)
			}
			if ok {
				published[m] = struct{}{}
				matched = true
			}
		}
		if !matched {
			return nil, fmt.Errorf("published module pattern %q matches no module in %s", pattern, ws.GoWork)
		}
	}

	return published, nil
}

// committedGoMods returns the go.mod files in the git index under the
// workspace root, as slash-separated paths relative to that root.
func committedGoMods(ctx context.Context, ws *Workspace) ([]string, error) {
	git := libgit.NewClient(libgit.WithDir(ws.Root))
	res, err := git.LsFiles(ctx, &libgit.LsFilesOpts{
		NullTerminated: true,
		PathSpec:       []string{":(glob)**/go.mod"},
	})
	if err != nil {
		return nil, fmt.Errorf("listing committed go.mod files: %w: %s", err, res.String())
	}

	goMods := []string{}
	for _, p := range bytes.Split(res.Stdout, []byte{0}) {
		if len(p) > 0 {
			goMods = append(goMods, string(p))
		}
	}

	return goMods, nil
}

// findUnregisteredModules returns a violation for each go.mod in goMods whose
// directory isn't a workspace module. goMods are slash-separated paths
// relative to the workspace root. Directories the go command ignores are
// skipped.
func findUnregisteredModules(ws *Workspace, goMods []string) []*WorkspaceViolation {
	registered := map[string]struct{}{}
	for _, m := range ws.Modules {
		registered[m.Dir] = struct{}{}
	}

	violations := []*WorkspaceViolation{}
	for _, goMod := range goMods {
		dir := path.Dir(goMod)
		if isIgnoredGoDir(dir) {
			continue
		}
		if _, ok := registered[dir]; ok {
			continue
		}

		violations = append(violations, &WorkspaceViolation{
			Kind:   WorkspaceViolationUnregisteredModule,
			Module: dir,
			Detail: fmt.Sprintf("%s is committed but %s has no 'use ./%s'",
				goMod, filepath.Base(ws.GoWork), dir),
		})
	}

	return violations
}

// isIgnoredGoDir reports whether the go command ignores a directory, which is
// the case when any element is testdata or vendor, or starts with "." or "_".
func isIgnoredGoDir(dir string) bool {
	if dir == "." {
		return false
	}

	for _, elem := range strings.Split(dir, "/") {
		if elem == "testdata" || elem == "vendor" || strings.HasPrefix(elem, ".") || strings.HasPrefix(elem, "_") {
			return true
		}
	}

	return false
}

// findReplaceViolations checks that unpublished module m, whose go.mod is
// modFile, replaces every workspace module it requires with that module's
// directory. Without that replace, a GOWORK=off build of m uses the required
// published version instead of the code on the branch.
func findReplaceViolations(ws *Workspace, m *WorkspaceModule, modFile *modfile.File) []*WorkspaceViolation {
	violations := []*WorkspaceViolation{}
	for _, req := range modFile.Require {
		dep := workspaceModuleByPath(ws, req.Mod.Path)
		if dep == nil || dep == m {
			continue
		}

		want := localReplacePath(ws, m, dep)
		fix := fmt.Sprintf("use 'replace %s => %s'", dep.Path, want)
		var replaces []*modfile.Replace
		for _, rep := range modFile.Replace {
			if rep.Old.Path == dep.Path {
				replaces = append(replaces, rep)
			}
		}

		if len(replaces) == 0 {
			violations = append(violations, &WorkspaceViolation{
				Kind:   WorkspaceViolationMissingLocalReplace,
				Module: m.Dir,
				Detail: fmt.Sprintf("requires %s %s without a replace; %s",
					dep.Path, req.Mod.Version, fix),
			})
			continue
		}

		for _, rep := range replaces {
			if detail := invalidReplaceDetail(ws, m, dep, rep); detail != "" {
				violations = append(violations, &WorkspaceViolation{
					Kind:   WorkspaceViolationInvalidLocalReplace,
					Module: m.Dir,
					Detail: detail + "; " + fix,
				})
			}
		}
	}

	return violations
}

// findPublishedReplaceViolations checks that published module m, whose go.mod
// is modFile, doesn't replace any other workspace module. Its importers ignore
// the replace and build against the version m requires, so a GOWORK=off build
// of m has to as well for its tests to mean anything.
func findPublishedReplaceViolations(ws *Workspace, m *WorkspaceModule, modFile *modfile.File) []*WorkspaceViolation {
	violations := []*WorkspaceViolation{}
	for _, rep := range modFile.Replace {
		dep := workspaceModuleByPath(ws, rep.Old.Path)
		if dep == nil || dep == m {
			continue
		}

		old := rep.Old.Path
		if rep.Old.Version != "" {
			old += " " + rep.Old.Version
		}
		violations = append(violations, &WorkspaceViolation{
			Kind:   WorkspaceViolationPublishedReplace,
			Module: m.Dir,
			Detail: fmt.Sprintf("published module replaces %s, which its importers ignore; remove 'replace %s => %s' and require a released version of %s",
				dep.Path, old, formatReplaceTarget(rep), dep.Path),
		})
	}

	return violations
}

// invalidReplaceDetail explains why rep, a replace of workspace module dep in
// module m, doesn't point every version of dep at dep's directory. It returns
// an empty string when rep is valid.
func invalidReplaceDetail(ws *Workspace, m, dep *WorkspaceModule, rep *modfile.Replace) string {
	switch {
	case rep.Old.Version != "":
		return fmt.Sprintf("replace %s %s => %s only applies to version %s",
			dep.Path, rep.Old.Version, formatReplaceTarget(rep), rep.Old.Version)
	case rep.New.Version != "":
		return fmt.Sprintf("replace %s => %s isn't a directory replace",
			dep.Path, formatReplaceTarget(rep))
	}

	target := filepath.FromSlash(rep.New.Path)
	if !filepath.IsAbs(target) {
		target = filepath.Join(ws.AbsDir(m), target)
	}
	if filepath.Clean(target) == filepath.Clean(ws.AbsDir(dep)) {
		return ""
	}

	resolved := target
	if rel, err := filepath.Rel(ws.Root, target); err == nil {
		resolved = filepath.ToSlash(rel)
	}

	return fmt.Sprintf("replace %s => %s resolves to %s, not %s",
		dep.Path, rep.New.Path, resolved, dep.Dir)
}

func formatReplaceTarget(rep *modfile.Replace) string {
	if rep.New.Version == "" {
		return rep.New.Path
	}

	return rep.New.Path + " " + rep.New.Version
}

// localReplacePath returns the replace target, relative to m, for dep's
// directory, e.g. "../api" or "./vault/hcp_link/proto".
func localReplacePath(ws *Workspace, m, dep *WorkspaceModule) string {
	rel, err := filepath.Rel(ws.AbsDir(m), ws.AbsDir(dep))
	if err != nil {
		return ws.AbsDir(dep)
	}

	rel = filepath.ToSlash(rel)
	if rel != ".." && !strings.HasPrefix(rel, "../") {
		rel = "./" + rel
	}

	return rel
}

func workspaceModuleByPath(ws *Workspace, modPath string) *WorkspaceModule {
	for _, m := range ws.Modules {
		if m.Path == modPath {
			return m
		}
	}

	return nil
}

// Success reports whether the workspace has no violations.
func (r *CheckWorkspaceRes) Success() bool {
	return r != nil && len(r.Violations) == 0
}

// String summarizes the result in one sentence, without a trailing period.
func (r *CheckWorkspaceRes) String() string {
	switch {
	case r == nil:
		return "no result"
	case len(r.Violations) == 1:
		return "go.work and the go.mod files have 1 problem"
	case len(r.Violations) > 1:
		return fmt.Sprintf("go.work and the go.mod files have %d problems", len(r.Violations))
	case r.Workspace == nil:
		return "go.work and the go.mod files agree"
	case len(r.Workspace.Modules) == 1:
		return "go.work and the go.mod file of 1 module agree"
	default:
		return fmt.Sprintf("go.work and the go.mod files of %d modules agree", len(r.Workspace.Modules))
	}
}

// ToJSON marshals the result to JSON.
func (r *CheckWorkspaceRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling workspace check to JSON: %w", err)
	}

	return b, nil
}

// ToTable renders the summary, followed by a table of problems if there are
// any.
func (r *CheckWorkspaceRes) ToTable() string {
	if r == nil || len(r.Violations) == 0 {
		return r.String() + "."
	}

	return r.String() + ".\n\n" + r.table().Render()
}

// ToMarkdown renders a heading and the summary, followed by a markdown table
// of problems if there are any. The heading matches the tidy and generate
// checks, whose output shares the CI job summary.
func (r *CheckWorkspaceRes) ToMarkdown() string {
	out := "# Go workspace check\n\n" + r.String() + "."
	if r == nil || len(r.Violations) == 0 {
		return out
	}

	return out + "\n\n" + r.table().RenderMarkdown()
}

func (r *CheckWorkspaceRes) table() table.Writer {
	t := newTable()
	t.AppendHeader(table.Row{"Module", "Problem"})
	for _, v := range r.Violations {
		t.AppendRow(table.Row{v.Module, v.Detail})
	}

	return t
}
