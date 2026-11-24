// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// GroupPackagesReq is a request to group Go packages by the workspace module
// that contains them.
type GroupPackagesReq struct {
	// GoWork is the path to go.work. When empty, the first go.work in the
	// current directory or its parents is used.
	GoWork string
	// Dir is the directory that relative paths in Packages are resolved
	// against. When empty, the current directory is used.
	Dir string
	// Packages are package import paths, package directories or .go files. A
	// .go file stands for the package in its directory. Entries that exist on
	// disk are paths; all others are import paths.
	Packages []string
}

// GroupPackagesRes is the response to a GroupPackagesReq.
type GroupPackagesRes struct {
	List *ModulePackagesList
}

// Run groups the packages by module without running go. Modules are in
// go.work order and only include modules with at least one package. Packages
// are import paths and keep their input order with duplicates removed. Run
// fails if any entry isn't a package directory in a workspace module.
func (r *GroupPackagesReq) Run(ctx context.Context) (*GroupPackagesRes, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	ws, err := LoadWorkspaceOrFind(ctx, r.GoWork)
	if err != nil {
		return nil, err
	}

	baseDir := r.Dir
	if baseDir == "" {
		baseDir, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("getting working directory: %w", err)
		}
	}

	byDir := map[string][]string{}
	seen := map[string]struct{}{}
	problems := []string{}
	for _, entry := range r.Packages {
		m, pkg, err := ws.resolvePackage(baseDir, entry)
		if err != nil {
			problems = append(problems, err.Error())
			continue
		}

		if _, ok := seen[pkg]; ok {
			continue
		}
		seen[pkg] = struct{}{}
		byDir[m.Dir] = append(byDir[m.Dir], pkg)
	}

	if len(problems) > 0 {
		return nil, fmt.Errorf("resolving packages in %s: %s", ws.GoWork, strings.Join(problems, "; "))
	}

	res := &GroupPackagesRes{List: &ModulePackagesList{Modules: []*ModulePackages{}}}
	for _, m := range ws.Modules {
		if packages, ok := byDir[m.Dir]; ok {
			res.List.Modules = append(res.List.Modules, &ModulePackages{
				Path:     m.Path,
				Dir:      m.Dir,
				Packages: packages,
			})
		}
	}

	return res, nil
}

// ToJSON marshals the grouped packages to JSON.
func (r *GroupPackagesRes) ToJSON() ([]byte, error) {
	if r == nil || r.List == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r.List)
	if err != nil {
		return nil, fmt.Errorf("marshaling grouped packages to JSON: %w", err)
	}

	return b, nil
}

// ToTable renders the packages of each module as a text table.
func (r *GroupPackagesRes) ToTable() string {
	return r.table().Render()
}

// ToMarkdown renders the packages of each module as a markdown table.
func (r *GroupPackagesRes) ToMarkdown() string {
	t := r.table()
	t.SetTitle("Go packages by module")
	return t.RenderMarkdown()
}

func (r *GroupPackagesRes) table() table.Writer {
	t := newTable()
	t.AppendHeader(table.Row{"Module dir", "Packages"})
	if r == nil || r.List == nil {
		return t
	}

	for _, m := range r.List.Modules {
		t.AppendRow(table.Row{m.Dir, strings.Join(m.Packages, "\n")})
	}

	return t
}

// resolvePackage returns the module and import path of entry, which is a
// package import path, a package directory or a .go file. Relative paths are
// resolved against baseDir. An entry that exists on disk is a path, so "vault"
// run from the repository root is the vault directory, not an import path.
func (w *Workspace) resolvePackage(baseDir, entry string) (*WorkspaceModule, string, error) {
	path := filepath.FromSlash(entry)
	if !filepath.IsAbs(path) {
		path = filepath.Join(baseDir, path)
	}
	path = filepath.Clean(path)

	var dir string
	info, err := os.Stat(path)
	switch {
	case err == nil && info.IsDir():
		dir = path
	case err == nil && filepath.Ext(path) == ".go":
		dir = filepath.Dir(path)
	case err == nil:
		return nil, "", fmt.Errorf("%s is not a directory or a .go file", entry)
	case !errors.Is(err, fs.ErrNotExist):
		return nil, "", fmt.Errorf("checking %s: %w", entry, err)
	}

	var m *WorkspaceModule
	if dir != "" {
		m = w.moduleForDir(dir)
		if m == nil {
			return nil, "", fmt.Errorf("%s is outside every module in %s", entry, filepath.Base(w.GoWork))
		}
	} else {
		m = w.ModuleForPackage(entry)
		if m == nil {
			return nil, "", fmt.Errorf("%s doesn't exist and isn't a package in any module in %s", entry, filepath.Base(w.GoWork))
		}
		dir = filepath.Join(w.AbsDir(m), filepath.FromSlash(strings.TrimPrefix(entry, m.Path)))
	}

	pkg, err := w.packageImportPath(m, dir)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", entry, err)
	}

	return m, pkg, nil
}

// moduleForDir returns the workspace module whose directory contains dir, an
// absolute path, or nil if no module does. Nested modules win over their
// parents.
func (w *Workspace) moduleForDir(dir string) *WorkspaceModule {
	var found *WorkspaceModule
	for _, m := range w.Modules {
		if !dirContains(w.AbsDir(m), dir) {
			continue
		}
		if found == nil || len(w.AbsDir(m)) > len(w.AbsDir(found)) {
			found = m
		}
	}

	return found
}

// packageImportPath returns the import path of the package in dir, an
// absolute directory inside module m. It fails if dir isn't a package that
// the go command would build as part of m.
func (w *Workspace) packageImportPath(m *WorkspaceModule, dir string) (string, error) {
	dir = filepath.Clean(dir)
	modDir := w.AbsDir(m)
	if !dirContains(modDir, dir) {
		return "", fmt.Errorf("%s is outside module %s", dir, m.Dir)
	}

	rel, err := filepath.Rel(modDir, dir)
	if err != nil {
		return "", fmt.Errorf("making %s relative to module %s: %w", dir, m.Dir, err)
	}
	rel = filepath.ToSlash(rel)
	if isIgnoredGoDir(rel) {
		return "", fmt.Errorf("the go command ignores %s in module %s", rel, m.Dir)
	}

	// A go.mod between dir and the module root makes dir part of another
	// module, which go.work doesn't use. dirContains guarantees that walking up
	// from dir reaches modDir.
	for d := dir; d != modDir; d = filepath.Dir(d) {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return "", fmt.Errorf("%s is in a module that isn't in %s", rel, filepath.Base(w.GoWork))
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("module %s has no directory %s", m.Dir, rel)
	}
	if !containsGoFile(entries) {
		return "", fmt.Errorf("module %s has no Go files in %s", m.Dir, rel)
	}

	if rel == "." {
		return m.Path, nil
	}

	return m.Path + "/" + rel, nil
}

// dirContains reports whether dir is parent or one of its subdirectories.
// Both must be clean absolute paths.
func dirContains(parent, dir string) bool {
	rel, err := filepath.Rel(parent, dir)
	if err != nil {
		return false
	}

	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func containsGoFile(entries []fs.DirEntry) bool {
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			return true
		}
	}

	return false
}
