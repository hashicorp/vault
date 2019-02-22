package golinters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/golangci/go-tools/config"
	"github.com/golangci/go-tools/stylecheck"

	"github.com/golangci/go-tools/lint"
	"github.com/golangci/go-tools/lint/lintutil"
	"github.com/golangci/go-tools/simple"
	"github.com/golangci/go-tools/staticcheck"
	"github.com/golangci/go-tools/unused"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	libpackages "github.com/golangci/golangci-lint/pkg/packages"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	MegacheckParentName      = "megacheck"
	MegacheckStaticcheckName = "staticcheck"
	MegacheckUnusedName      = "unused"
	MegacheckGosimpleName    = "gosimple"
	MegacheckStylecheckName  = "stylecheck"
)

type Staticcheck struct {
	megacheck
}

func NewStaticcheck() *Staticcheck {
	return &Staticcheck{
		megacheck: megacheck{
			staticcheckEnabled: true,
		},
	}
}

func (Staticcheck) Name() string { return MegacheckStaticcheckName }
func (Staticcheck) Desc() string {
	return "Staticcheck is a go vet on steroids, applying a ton of static analysis checks"
}

type Gosimple struct {
	megacheck
}

func NewGosimple() *Gosimple {
	return &Gosimple{
		megacheck: megacheck{
			gosimpleEnabled: true,
		},
	}
}

func (Gosimple) Name() string { return MegacheckGosimpleName }
func (Gosimple) Desc() string {
	return "Linter for Go source code that specializes in simplifying a code"
}

type Unused struct {
	megacheck
}

func NewUnused() *Unused {
	return &Unused{
		megacheck: megacheck{
			unusedEnabled: true,
		},
	}
}

func (Unused) Name() string { return MegacheckUnusedName }
func (Unused) Desc() string {
	return "Checks Go code for unused constants, variables, functions and types"
}

type Stylecheck struct {
	megacheck
}

func NewStylecheck() *Stylecheck {
	return &Stylecheck{
		megacheck: megacheck{
			stylecheckEnabled: true,
		},
	}
}

func (Stylecheck) Name() string { return MegacheckStylecheckName }
func (Stylecheck) Desc() string { return "Stylecheck is a replacement for golint" }

type megacheck struct {
	unusedEnabled      bool
	gosimpleEnabled    bool
	staticcheckEnabled bool
	stylecheckEnabled  bool
}

func (megacheck) Name() string {
	return MegacheckParentName
}

func (megacheck) Desc() string {
	return "" // shouldn't be called
}

func (m *megacheck) enableChildLinter(name string) error {
	switch name {
	case MegacheckStaticcheckName:
		m.staticcheckEnabled = true
	case MegacheckGosimpleName:
		m.gosimpleEnabled = true
	case MegacheckUnusedName:
		m.unusedEnabled = true
	case MegacheckStylecheckName:
		m.stylecheckEnabled = true
	default:
		return fmt.Errorf("invalid child linter name %s for metalinter %s", name, m.Name())
	}

	return nil
}

type MegacheckMetalinter struct{}

func (MegacheckMetalinter) Name() string {
	return MegacheckParentName
}

func (MegacheckMetalinter) BuildLinterConfig(enabledChildren []string) (*linter.Config, error) {
	var m megacheck
	for _, name := range enabledChildren {
		if err := m.enableChildLinter(name); err != nil {
			return nil, err
		}
	}

	// TODO: merge linter.Config and linter.Linter or refactor it in another way
	return &linter.Config{
		Linter:           m,
		EnabledByDefault: false,
		NeedsTypeInfo:    true,
		NeedsSSARepr:     true,
		InPresets:        []string{linter.PresetStyle, linter.PresetBugs, linter.PresetUnused},
		Speed:            1,
		AlternativeNames: nil,
		OriginalURL:      "",
		ParentLinterName: "",
	}, nil
}

func (MegacheckMetalinter) DefaultChildLinterNames() []string {
	// no stylecheck here for backwards compatibility for users who enabled megacheck: don't enable extra
	// linter for them
	return []string{MegacheckStaticcheckName, MegacheckGosimpleName, MegacheckUnusedName}
}

func (m MegacheckMetalinter) AllChildLinterNames() []string {
	return append(m.DefaultChildLinterNames(), MegacheckStylecheckName)
}

func (m MegacheckMetalinter) isValidChild(name string) bool {
	for _, child := range m.AllChildLinterNames() {
		if child == name {
			return true
		}
	}

	return false
}

func prettifyCompilationError(err packages.Error) error {
	i, _ := TypeCheck{}.parseError(err)
	if i == nil {
		return err
	}

	shortFilename, pathErr := fsutils.ShortestRelPath(i.Pos.Filename, "")
	if pathErr != nil {
		return err
	}

	errText := shortFilename
	if i.Line() != 0 {
		errText += fmt.Sprintf(":%d", i.Line())
	}
	errText += fmt.Sprintf(": %s", i.Text)
	return errors.New(errText)
}

func (m megacheck) canAnalyze(lintCtx *linter.Context) bool {
	if len(lintCtx.NotCompilingPackages) == 0 {
		return true
	}

	var errPkgs []string
	var errs []packages.Error
	for _, p := range lintCtx.NotCompilingPackages {
		if p.Name == "main" {
			// megacheck crashes on not compiling packages but main packages
			// aren't reachable by megacheck: other packages can't depend on them.
			continue
		}

		errPkgs = append(errPkgs, p.String())
		errs = append(errs, libpackages.ExtractErrors(p, lintCtx.ASTCache)...)
	}

	if len(errPkgs) == 0 { // only main packages do not compile
		return true
	}

	// TODO: print real linter names in this message
	warnText := fmt.Sprintf("Can't run megacheck because of compilation errors in packages %s", errPkgs)
	if len(errs) != 0 {
		warnText += fmt.Sprintf(": %s", prettifyCompilationError(errs[0]))
		if len(errs) > 1 {
			const runCmd = "golangci-lint run --no-config --disable-all -E typecheck"
			warnText += fmt.Sprintf(" and %d more errors: run `%s` to see all errors", len(errs)-1, runCmd)
		}
	}
	lintCtx.Log.Warnf("%s", warnText)

	// megacheck crashes if there are not compiling packages
	return false
}

func (m megacheck) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	if !m.canAnalyze(lintCtx) {
		return nil, nil
	}

	issues, err := m.runMegacheck(lintCtx.Packages, lintCtx.Settings().Unused.CheckExported)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run megacheck")
	}

	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	meta := MegacheckMetalinter{}
	for _, i := range issues {
		if !meta.isValidChild(i.Checker) {
			lintCtx.Log.Warnf("Bad megacheck checker name %q", i.Checker)
			continue
		}

		res = append(res, result.Issue{
			Pos:        i.Position,
			Text:       markIdentifiers(i.Text),
			FromLinter: i.Checker,
		})
	}
	return res, nil
}

func (m megacheck) runMegacheck(workingPkgs []*packages.Package, checkExportedUnused bool) ([]lint.Problem, error) {
	var checkers []lint.Checker

	if m.gosimpleEnabled {
		checkers = append(checkers, simple.NewChecker())
	}
	if m.staticcheckEnabled {
		checkers = append(checkers, staticcheck.NewChecker())
	}
	if m.stylecheckEnabled {
		checkers = append(checkers, stylecheck.NewChecker())
	}
	if m.unusedEnabled {
		uc := unused.NewChecker(unused.CheckAll)
		uc.ConsiderReflection = true
		uc.WholeProgram = checkExportedUnused
		checkers = append(checkers, unused.NewLintChecker(uc))
	}

	if len(checkers) == 0 {
		return nil, nil
	}

	cfg := config.Config{}
	opts := &lintutil.Options{
		// TODO: get current go version, but now it doesn't matter,
		// may be needed after next updates of megacheck
		GoVersion: 11,

		Config: cfg,
		// TODO: support Ignores option
	}

	return runMegacheckCheckers(checkers, opts, workingPkgs)
}

// parseIgnore is a copy from megacheck code just to not fork megacheck
func parseIgnore(s string) ([]lint.Ignore, error) {
	var out []lint.Ignore
	if s == "" {
		return nil, nil
	}
	for _, part := range strings.Fields(s) {
		p := strings.Split(part, ":")
		if len(p) != 2 {
			return nil, errors.New("malformed ignore string")
		}
		path := p[0]
		checks := strings.Split(p[1], ",")
		out = append(out, &lint.GlobIgnore{Pattern: path, Checks: checks})
	}
	return out, nil
}

func runMegacheckCheckers(cs []lint.Checker, opt *lintutil.Options, workingPkgs []*packages.Package) ([]lint.Problem, error) {
	stats := lint.PerfStats{
		CheckerInits: map[string]time.Duration{},
	}

	ignores, err := parseIgnore(opt.Ignores)
	if err != nil {
		return nil, err
	}

	var problems []lint.Problem
	if len(workingPkgs) == 0 {
		return problems, nil
	}

	l := &lint.Linter{
		Checkers:      cs,
		Ignores:       ignores,
		GoVersion:     opt.GoVersion,
		ReturnIgnored: opt.ReturnIgnored,
		Config:        opt.Config,

		MaxConcurrentJobs: opt.MaxConcurrentJobs,
		PrintStats:        opt.PrintStats,
	}
	problems = append(problems, l.Lint(workingPkgs, &stats)...)

	return problems, nil
}
