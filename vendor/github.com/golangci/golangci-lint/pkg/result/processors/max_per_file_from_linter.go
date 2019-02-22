package processors

import (
	"github.com/golangci/golangci-lint/pkg/result"
)

type linterToCountMap map[string]int
type fileToLinterToCountMap map[string]linterToCountMap

type MaxPerFileFromLinter struct {
	flc fileToLinterToCountMap
}

var _ Processor = &MaxPerFileFromLinter{}

func NewMaxPerFileFromLinter() *MaxPerFileFromLinter {
	return &MaxPerFileFromLinter{
		flc: fileToLinterToCountMap{},
	}
}

func (p MaxPerFileFromLinter) Name() string {
	return "max_per_file_from_linter"
}

var maxPerFileFromLinterConfig = map[string]int{
	"gofmt":     1,
	"goimports": 1,
	"typecheck": 3,
}

func (p *MaxPerFileFromLinter) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssues(issues, func(i *result.Issue) bool {
		limit := maxPerFileFromLinterConfig[i.FromLinter]
		if limit == 0 {
			return true
		}

		lm := p.flc[i.FilePath()]
		if lm == nil {
			p.flc[i.FilePath()] = linterToCountMap{}
		}
		count := p.flc[i.FilePath()][i.FromLinter]
		if count >= limit {
			return false
		}

		p.flc[i.FilePath()][i.FromLinter]++
		return true
	}), nil
}

func (p MaxPerFileFromLinter) Finish() {}
