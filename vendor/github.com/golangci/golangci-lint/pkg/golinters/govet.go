package golinters

import (
	"context"
	"go/ast"
	"go/token"

	govetAPI "github.com/golangci/govet"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Govet struct{}

func (Govet) Name() string {
	return "govet"
}

func (Govet) Desc() string {
	return "Vet examines Go source code and reports suspicious constructs, " +
		"such as Printf calls whose arguments do not align with the format string"
}

func (g Govet) Run(_ context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var govetIssues []govetAPI.Issue
	var err error
	govetIssues, err = g.runImpl(lintCtx)
	if err != nil {
		return nil, err
	}

	if len(govetIssues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(govetIssues))
	for _, i := range govetIssues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       markIdentifiers(i.Message),
			FromLinter: g.Name(),
		})
	}
	return res, nil
}

func (g Govet) runImpl(lintCtx *linter.Context) ([]govetAPI.Issue, error) {
	// TODO: check .S asm files: govet can do it if pass dirs
	var govetIssues []govetAPI.Issue
	for _, pkg := range lintCtx.Program.InitialPackages() {
		if len(pkg.Files) == 0 {
			continue
		}

		issues, err := govetAPI.Analyze(pkg.Files, lintCtx.Program.Fset, pkg,
			lintCtx.Settings().Govet.CheckShadowing, getPath)
		if err != nil {
			return nil, err
		}
		govetIssues = append(govetIssues, issues...)
	}

	return govetIssues, nil
}

func getPath(f *ast.File, fset *token.FileSet) (string, error) {
	return fsutils.ShortestRelPath(fset.Position(f.Pos()).Filename, "")
}
