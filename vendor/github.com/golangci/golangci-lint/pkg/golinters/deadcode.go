package golinters

import (
	"context"
	"fmt"

	deadcodeAPI "github.com/golangci/go-misc/deadcode"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Deadcode struct{}

func (Deadcode) Name() string {
	return "deadcode"
}

func (Deadcode) Desc() string {
	return "Finds unused code"
}

func (d Deadcode) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	issues, err := deadcodeAPI.Run(lintCtx.Program)
	if err != nil {
		return nil, err
	}

	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]result.Issue, 0, len(issues))
	for _, i := range issues {
		res = append(res, result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.UnusedIdentName, lintCtx.Cfg)),
			FromLinter: d.Name(),
		})
	}
	return res, nil
}
