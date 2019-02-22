package golinters // nolint:dupl

import (
	"context"
	"fmt"

	varcheckAPI "github.com/golangci/check/cmd/varcheck"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Varcheck struct{}

func (Varcheck) Name() string {
	return "varcheck"
}

func (Varcheck) Desc() string {
	return "Finds unused global variables and constants"
}

func (v Varcheck) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	issues := varcheckAPI.Run(lintCtx.Program, lintCtx.Settings().Varcheck.CheckExportedFields)
	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.VarName, lintCtx.Cfg)),
			FromLinter: v.Name(),
		})
	}
	return res, nil
}
