// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/mod/modfile"
)

// Workspace describes a go.work file and the modules in its use directives.
type Workspace struct {
	// GoWork is the absolute path to go.work.
	GoWork string `json:"go_work"`
	// Root is the absolute directory containing go.work.
	Root string `json:"root"`
	// Modules are in go.work order.
	Modules []*WorkspaceModule `json:"modules"`
}

// WorkspaceModule is a module from a go.work use directive.
type WorkspaceModule struct {
	// Path is the module path from the module's go.mod.
	Path string `json:"path"`
	// Dir is slash-separated and relative to the workspace root; "." for the
	// root module.
	Dir string `json:"dir"`
	// GoMod is the slash-separated go.mod path relative to the workspace root,
	// e.g. "go.mod" or "sdk/go.mod".
	GoMod string `json:"go_mod"`
}

// ModulePackages is a module and a set of its package import paths.
type ModulePackages struct {
	Path     string   `json:"path"`
	Dir      string   `json:"dir"`
	Packages []string `json:"packages"`
}

// ModulePackagesList is a list of modules and their packages.
type ModulePackagesList struct {
	Modules []*ModulePackages `json:"modules"`
}

// FindGoWork walks up from startDir and returns the path of the first go.work
// file it finds. It ignores the GOWORK environment variable, which CI sets to
// "off" for module-mode builds.
func FindGoWork(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("resolving %s: %w", startDir, err)
	}

	for {
		path := filepath.Join(dir, "go.work")
		info, err := os.Stat(path)
		switch {
		case err == nil && !info.IsDir():
			return path, nil
		case err != nil && !errors.Is(err, fs.ErrNotExist):
			return "", fmt.Errorf("checking for %s: %w", path, err)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.work file found in %s or any parent directory", startDir)
		}
		dir = parent
	}
}

// LoadWorkspace parses goWork and the go.mod of every module in its use
// directives.
func LoadWorkspace(ctx context.Context, goWork string) (*Workspace, error) {
	goWork, err := filepath.Abs(goWork)
	if err != nil {
		return nil, fmt.Errorf("resolving go.work path: %w", err)
	}

	slog.Default().DebugContext(ctx, "loading go workspace", slog.String("go_work", goWork))
	data, err := os.ReadFile(goWork)
	if err != nil {
		return nil, fmt.Errorf("reading go.work: %w", err)
	}

	workFile, err := modfile.ParseWork(goWork, data, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing go.work: %w", err)
	}

	ws := &Workspace{
		GoWork:  goWork,
		Root:    filepath.Dir(goWork),
		Modules: []*WorkspaceModule{},
	}

	dirsByPath := map[string]string{}
	for _, use := range workFile.Use {
		dir, err := ws.relDir(use.Path)
		if err != nil {
			return nil, err
		}

		goModPath := filepath.Join(ws.Root, filepath.FromSlash(dir), "go.mod")
		slog.Default().DebugContext(ctx, "loading workspace module", slog.String("go_mod", goModPath))
		modData, err := os.ReadFile(goModPath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("go.work uses %s, which has no go.mod", use.Path)
			}
			return nil, fmt.Errorf("reading %s: %w", goModPath, err)
		}

		modFile, err := modfile.ParseLax(goModPath, modData, nil)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", goModPath, err)
		}
		if modFile.Module == nil || modFile.Module.Mod.Path == "" {
			return nil, fmt.Errorf("%s has no module directive", goModPath)
		}

		modPath := modFile.Module.Mod.Path
		if prev, ok := dirsByPath[modPath]; ok {
			return nil, fmt.Errorf("module %s is used twice in go.work: %s and %s", modPath, prev, dir)
		}
		dirsByPath[modPath] = dir

		goMod := "go.mod"
		if dir != "." {
			goMod = dir + "/go.mod"
		}
		ws.Modules = append(ws.Modules, &WorkspaceModule{Path: modPath, Dir: dir, GoMod: goMod})
	}

	return ws, nil
}

// LoadWorkspaceOrFind loads goWork, or the go.work found by FindGoWork from
// the current working directory when goWork is empty.
func LoadWorkspaceOrFind(ctx context.Context, goWork string) (*Workspace, error) {
	if goWork == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("getting working directory: %w", err)
		}
		goWork, err = FindGoWork(cwd)
		if err != nil {
			return nil, err
		}
	}

	return LoadWorkspace(ctx, goWork)
}

// ModuleByDir returns the workspace module whose directory is dir. dir is
// relative to the workspace root, e.g. "sdk", "./sdk" or ".".
func (w *Workspace) ModuleByDir(dir string) (*WorkspaceModule, error) {
	if w == nil {
		return nil, errors.New("uninitialized workspace")
	}

	want, err := w.relDir(dir)
	if err != nil {
		return nil, err
	}

	for _, m := range w.Modules {
		if m.Dir == want {
			return m, nil
		}
	}

	return nil, fmt.Errorf("%s is not a module in %s", dir, w.GoWork)
}

// SelectModules returns the modules whose directories are in dirs, in go.work
// order. It returns every module when dirs is empty.
func (w *Workspace) SelectModules(dirs []string) ([]*WorkspaceModule, error) {
	if w == nil {
		return nil, errors.New("uninitialized workspace")
	}

	if len(dirs) == 0 {
		return w.Modules, nil
	}

	selected := map[string]struct{}{}
	for _, dir := range dirs {
		m, err := w.ModuleByDir(dir)
		if err != nil {
			return nil, err
		}
		selected[m.Dir] = struct{}{}
	}

	modules := []*WorkspaceModule{}
	for _, m := range w.Modules {
		if _, ok := selected[m.Dir]; ok {
			modules = append(modules, m)
		}
	}

	return modules, nil
}

// ModuleForPackage returns the workspace module that contains the package
// importPath, or nil if no module does. Nested modules win over their parents,
// e.g. "github.com/hashicorp/vault/api/auth/approle" belongs to the approle
// module, not the api module.
func (w *Workspace) ModuleForPackage(importPath string) *WorkspaceModule {
	if w == nil {
		return nil
	}

	var found *WorkspaceModule
	for _, m := range w.Modules {
		if importPath != m.Path && !strings.HasPrefix(importPath, m.Path+"/") {
			continue
		}
		if found == nil || len(m.Path) > len(found.Path) {
			found = m
		}
	}

	return found
}

// AbsDir returns the absolute directory of a workspace module.
func (w *Workspace) AbsDir(m *WorkspaceModule) string {
	return filepath.Join(w.Root, filepath.FromSlash(m.Dir))
}

// relDir normalizes a module directory from go.work or a flag to a clean,
// slash-separated path relative to the workspace root.
func (w *Workspace) relDir(dir string) (string, error) {
	native := filepath.FromSlash(dir)
	if filepath.IsAbs(native) {
		rel, err := filepath.Rel(w.Root, native)
		if err != nil {
			return "", fmt.Errorf("making %s relative to %s: %w", dir, w.Root, err)
		}
		native = rel
	}

	return filepath.ToSlash(filepath.Clean(native)), nil
}

// NormalizeTags splits a comma- or space-separated list of build tags, drops
// empty and duplicate entries, and keeps the original order. For example,
// ",deadlock" becomes ["deadlock"].
func NormalizeTags(tags string) []string {
	fields := strings.FieldsFunc(tags, func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	})

	normalized := []string{}
	for _, tag := range fields {
		if !slices.Contains(normalized, tag) {
			normalized = append(normalized, tag)
		}
	}

	return normalized
}
