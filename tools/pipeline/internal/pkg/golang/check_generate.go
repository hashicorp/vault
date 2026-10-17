// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package golang

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"

	libgit "github.com/hashicorp/vault/tools/pipeline/internal/pkg/git/client"
	"github.com/jedib0t/go-pretty/v6/table"
)

// GeneratedChangeKind describes how go generate changed a path.
type GeneratedChangeKind string

const (
	GeneratedChangeModified GeneratedChangeKind = "modified"
	GeneratedChangeAdded    GeneratedChangeKind = "added"
	GeneratedChangeDeleted  GeneratedChangeKind = "deleted"
	// GeneratedChangeAddedIgnored is a new file that git ignores. Builds that
	// skip generation would never see it, so it's a failure too.
	GeneratedChangeAddedIgnored GeneratedChangeKind = "added-ignored"
)

// maxGenerateDiffBytes caps the diff in the response so step summaries and
// logs stay readable.
const maxGenerateDiffBytes = 64 << 10

// CheckGenerateReq is a request to check that committed generated code
// matches what go generate produces.
type CheckGenerateReq struct {
	// GoWork is the path to go.work. When empty, the first go.work in the
	// current directory or its parents is used.
	GoWork string
	// Tags are comma- or space-separated build tags for go generate and the
	// generators it runs.
	Tags string
	// Modules are module directories relative to the workspace root. When
	// empty, every module in go.work is generated.
	Modules []string
}

// GenerateModuleResult is the go generate result for one module.
type GenerateModuleResult struct {
	Dir   string `json:"dir"`
	Path  string `json:"path"`
	Error string `json:"error,omitempty"`
}

// GeneratedChange is a path that go generate changed.
type GeneratedChange struct {
	// Path is slash-separated and relative to the git repository root.
	Path   string              `json:"path"`
	Change GeneratedChangeKind `json:"change"`
}

// CheckGenerateRes is the response to a CheckGenerateReq.
type CheckGenerateRes struct {
	Modules []*GenerateModuleResult `json:"modules"`
	Changes []*GeneratedChange      `json:"changes"`
	// Diff is the git diff of changed tracked files, truncated to 64 KiB.
	Diff string `json:"diff"`
}

// statusEntry is a path's git status before or after generation.
type statusEntry struct {
	// code is the two-letter porcelain v1 status, e.g. " M", "??" or "!!".
	code string
	// hash is the SHA-256 of a non-ignored regular file's contents.
	hash string
}

// Run runs go generate in each selected module and reports every file it
// added, changed or deleted. Changes are left in the working tree.
func (r *CheckGenerateReq) Run(ctx context.Context) (*CheckGenerateRes, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	ws, err := LoadWorkspaceOrFind(ctx, r.GoWork)
	if err != nil {
		return nil, err
	}

	modules, err := ws.SelectModules(r.Modules)
	if err != nil {
		return nil, err
	}

	topLevel, err := gitTopLevel(ctx, ws.Root)
	if err != nil {
		return nil, err
	}
	git := libgit.NewClient(libgit.WithDir(topLevel))

	before, err := statusSnapshot(ctx, git, topLevel)
	if err != nil {
		return nil, err
	}

	env := generateEnv(os.Getenv("GOFLAGS"), NormalizeTags(r.Tags))
	res := &CheckGenerateRes{Modules: []*GenerateModuleResult{}}
	for _, m := range modules {
		result := &GenerateModuleResult{Dir: m.Dir, Path: m.Path}
		out, err := runGo(ctx, ws.AbsDir(m), env, "generate", "./...")
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		if err != nil {
			result.Error = strings.TrimSpace(string(out.Stderr))
			if result.Error == "" {
				result.Error = err.Error()
			}
		}
		res.Modules = append(res.Modules, result)
	}

	after, err := statusSnapshot(ctx, git, topLevel)
	if err != nil {
		return nil, err
	}

	res.Changes = compareSnapshots(before, after)
	res.Diff, err = trackedDiff(ctx, git, res.Changes, before, after)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *CheckGenerateReq) validate() error {
	if r == nil {
		return errors.New("uninitialized")
	}

	// tools/tools.go uses the tools build tag to keep its generate directive,
	// which installs every tool, out of normal runs.
	if slices.Contains(NormalizeTags(r.Tags), "tools") {
		return errors.New("the tools build tag would run tools/tools.go's install directive; remove it from --tags")
	}

	return nil
}

// generateEnv returns the environment overrides for go generate. Tags go in
// GOFLAGS rather than -tags so generators that load packages with the go
// command, such as enumer, see the same tags. Target platform variables are
// cleared so the output doesn't depend on the job's cross-compile settings.
func generateEnv(goFlags string, tags []string) []string {
	env := []string{"GOOS=", "GOARCH=", "CC=", "CC_FOR_TARGET="}
	if len(tags) > 0 {
		goFlags = strings.TrimSpace(goFlags + " -tags=" + strings.Join(tags, ","))
		env = append(env, "GOFLAGS="+goFlags)
	}

	return env
}

func gitTopLevel(ctx context.Context, dir string) (string, error) {
	res, err := libgit.NewClient(libgit.WithDir(dir)).RevParse(ctx, &libgit.RevParseOpts{ShowTopLevel: true})
	if err != nil {
		return "", fmt.Errorf("finding the git repository root for %s: %w: %s", dir, err, res.String())
	}

	return strings.TrimSpace(string(res.Stdout)), nil
}

// statusSnapshot records the git status of every changed, untracked and
// ignored path, plus a content hash of each non-ignored regular file.
func statusSnapshot(ctx context.Context, git *libgit.Client, topLevel string) (map[string]*statusEntry, error) {
	res, err := git.Status(ctx, &libgit.StatusOpts{
		Porcelain:      true,
		NullTerminated: true,
		NoRenames:      true,
		UntrackedFiles: libgit.UntrackedFilesAll,
		Ignored:        libgit.IgnoredModeMatching,
	})
	if err != nil {
		return nil, fmt.Errorf("getting git status: %w: %s", err, res.String())
	}

	codes, err := parsePorcelainZ(res.Stdout)
	if err != nil {
		return nil, err
	}

	snapshot := make(map[string]*statusEntry, len(codes))
	for p, code := range codes {
		entry := &statusEntry{code: code}
		if code != "!!" {
			entry.hash, err = hashRegularFile(filepath.Join(topLevel, filepath.FromSlash(p)))
			if err != nil {
				return nil, err
			}
		}
		snapshot[p] = entry
	}
	slog.Default().DebugContext(ctx, "recorded git status snapshot", slog.Int("entries", len(snapshot)))

	return snapshot, nil
}

// parsePorcelainZ parses 'git status --porcelain=v1 -z --no-renames' output
// into a map of path to two-letter status code.
func parsePorcelainZ(out []byte) (map[string]string, error) {
	codes := map[string]string{}
	for _, entry := range bytes.Split(out, []byte{0}) {
		if len(entry) == 0 {
			continue
		}
		if len(entry) < 4 || entry[2] != ' ' {
			return nil, fmt.Errorf("unexpected git status entry: %q", entry)
		}
		codes[string(entry[3:])] = string(entry[:2])
	}

	return codes, nil
}

// hashRegularFile returns the hex SHA-256 of a regular file's contents, or an
// empty string when the path is missing or isn't a regular file.
func hashRegularFile(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", fmt.Errorf("checking %s: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return "", nil
	}

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hashing %s: %w", path, err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// compareSnapshots returns every path whose status or contents differ between
// the snapshots, sorted by path.
func compareSnapshots(before, after map[string]*statusEntry) []*GeneratedChange {
	paths := slices.Collect(maps.Keys(before))
	for p := range after {
		if _, ok := before[p]; !ok {
			paths = append(paths, p)
		}
	}
	slices.Sort(paths)

	changes := []*GeneratedChange{}
	for _, p := range paths {
		b, a := before[p], after[p]
		if b != nil && a != nil && *b == *a {
			continue
		}
		changes = append(changes, &GeneratedChange{Path: p, Change: classifyChange(b, a)})
	}

	return changes
}

func classifyChange(before, after *statusEntry) GeneratedChangeKind {
	switch {
	case after == nil:
		// An untracked file was removed, or a modified tracked file was
		// restored to its committed contents.
		if before.code == "??" || before.code == "!!" {
			return GeneratedChangeDeleted
		}
		return GeneratedChangeModified
	case after.code == "!!":
		return GeneratedChangeAddedIgnored
	case after.code == "??" && before == nil:
		return GeneratedChangeAdded
	case strings.Contains(after.code, "D"):
		return GeneratedChangeDeleted
	default:
		return GeneratedChangeModified
	}
}

// trackedDiff returns the git diff of the changed paths that git tracks.
func trackedDiff(
	ctx context.Context,
	git *libgit.Client,
	changes []*GeneratedChange,
	before, after map[string]*statusEntry,
) (string, error) {
	untracked := func(e *statusEntry) bool {
		return e != nil && (e.code == "??" || e.code == "!!")
	}

	pathSpec := []string{}
	for _, c := range changes {
		if untracked(before[c.Path]) || untracked(after[c.Path]) {
			continue
		}
		pathSpec = append(pathSpec, ":(literal)"+c.Path)
	}
	if len(pathSpec) == 0 {
		return "", nil
	}

	res, err := git.Diff(ctx, &libgit.DiffOpts{
		// Ignore the user's diff config so the output is a plain patch.
		NoColor:   true,
		NoExtDiff: true,
		SrcPrefix: "a/",
		DstPrefix: "b/",
		PathSpec:  pathSpec,
	})
	if err != nil {
		return "", fmt.Errorf("diffing generated files: %w: %s", err, res.String())
	}

	return truncateDiff(string(res.Stdout), maxGenerateDiffBytes), nil
}

// truncateDiff caps diff at limit bytes, backing up to a UTF-8 rune boundary
// so the output stays valid text.
func truncateDiff(diff string, limit int) string {
	if len(diff) <= limit {
		return diff
	}

	cut := limit
	for cut > 0 && !utf8.RuneStart(diff[cut]) {
		cut--
	}

	return diff[:cut] + fmt.Sprintf("\n... diff truncated at %d KiB ...\n", limit>>10)
}

// Success reports whether every module generated without errors or changes.
func (r *CheckGenerateRes) Success() bool {
	return r != nil && len(r.Changes) == 0 && len(r.failedModules()) == 0
}

func (r *CheckGenerateRes) failedModules() []*GenerateModuleResult {
	failed := []*GenerateModuleResult{}
	for _, m := range r.Modules {
		if m.Error != "" {
			failed = append(failed, m)
		}
	}

	return failed
}

// String summarizes the result.
func (r *CheckGenerateRes) String() string {
	if r == nil {
		return "no result"
	}

	return fmt.Sprintf("changed paths: %d, modules that failed to generate: %d of %d",
		len(r.Changes), len(r.failedModules()), len(r.Modules))
}

// ToJSON marshals the result to JSON.
func (r *CheckGenerateRes) ToJSON() ([]byte, error) {
	if r == nil {
		return nil, errors.New("uninitialized")
	}

	b, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("marshaling generate check to JSON: %w", err)
	}

	return b, nil
}

// ToTable renders module errors and changes as a text table, followed by the
// diff.
func (r *CheckGenerateRes) ToTable() string {
	out := r.table().Render()
	if r != nil && r.Diff != "" {
		out += "\n\n" + strings.TrimRight(r.Diff, "\n")
	}

	return out
}

// ToMarkdown renders module errors and changes as a markdown table, followed
// by the diff in a code block.
func (r *CheckGenerateRes) ToMarkdown() string {
	out := r.table().RenderMarkdown()
	if r != nil && r.Diff != "" {
		// Generated Go can contain backtick runs; a longer fence keeps them
		// inside the block.
		out += "\n\n````diff\n" + strings.TrimRight(r.Diff, "\n") + "\n````"
	}

	return out
}

func (r *CheckGenerateRes) table() table.Writer {
	t := newTable()
	t.SetTitle("Go generate check")
	t.AppendHeader(table.Row{"Result", "Path", "Detail"})
	if r == nil {
		return t
	}

	for _, m := range r.failedModules() {
		t.AppendRow(table.Row{"error", m.Dir, m.Error})
	}
	for _, c := range r.Changes {
		t.AppendRow(table.Row{c.Change, c.Path, ""})
	}

	if r.Success() {
		t.SetCaption(fmt.Sprintf("Generated code is up to date in %d modules.", len(r.Modules)))
	}

	return t
}
