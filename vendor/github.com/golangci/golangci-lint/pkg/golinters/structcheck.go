package golinters // nolint:dupl

import (
	"context"
	"fmt"

	structcheckAPI "github.com/golangci/check/cmd/structcheck"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Structcheck struct{}

func (Structcheck) Name() string {
	return "structcheck"
}

func (Structcheck) Desc() string {
	return "Finds an unused struct fields"
}

func (s Structcheck) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	issues := structcheckAPI.Run(lintCtx.Program, lintCtx.Settings().Structcheck.CheckExportedFields)
	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.FieldName, lintCtx.Cfg)),
			FromLinter: s.Name(),
		})
	}
	return res, nil
}
