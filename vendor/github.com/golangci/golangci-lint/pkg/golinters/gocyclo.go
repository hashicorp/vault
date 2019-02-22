package golinters

import (
	"context"
	"fmt"
	"sort"

	gocycloAPI "github.com/golangci/gocyclo/pkg/gocyclo"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gocyclo struct{}

func (Gocyclo) Name() string {
	return "gocyclo"
}

func (Gocyclo) Desc() string {
	return "Computes and checks the cyclomatic complexity of functions"
}

func (g Gocyclo) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var stats []gocycloAPI.Stat
	for _, f := range lintCtx.ASTCache.GetAllValidFiles() {
		stats = gocycloAPI.BuildStats(f.F, f.Fset, stats)
	}
	if len(stats) == 0 {
		return nil, nil
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Complexity > stats[j].Complexity
	})

	res := make([]result.Issue, 0, len(stats))
	for _, s := range stats {
		if s.Complexity <= lintCtx.Settings().Gocyclo.MinComplexity {
			break // Break as the stats is already sorted from greatest to least
		}

		res = append(res, result.Issue{
			Pos: s.Pos,
			Text: fmt.Sprintf("cyclomatic complexity %d of func %s is high (> %d)",
				s.Complexity, formatCode(s.FuncName, lintCtx.Cfg), lintCtx.Settings().Gocyclo.MinComplexity),
			FromLinter: g.Name(),
		})
	}

	return res, nil
}
