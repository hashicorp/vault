package golinters

import (
	"bytes"
	"context"
	"fmt"
	"go/token"

	gofmtAPI "github.com/golangci/gofmt/gofmt"
	goimportsAPI "github.com/golangci/gofmt/goimports"
	"golang.org/x/tools/imports"
	diffpkg "sourcegraph.com/sourcegraph/go-diff/diff"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Gofmt struct {
	UseGoimports bool
}

func (g Gofmt) Name() string {
	if g.UseGoimports {
		return "goimports"
	}

	return "gofmt"
}

func (g Gofmt) Desc() string {
	if g.UseGoimports {
		return "Goimports does everything that gofmt does. Additionally it checks unused imports"
	}

	return "Gofmt checks whether code was gofmt-ed. By default " +
		"this tool runs with -s option to check for code simplification"
}

func getFirstDeletedAndAddedLineNumberInHunk(h *diffpkg.Hunk) (firstDeleted, firstAdded int, err error) {
	lines := bytes.Split(h.Body, []byte{'\n'})
	lineNumber := int(h.OrigStartLine - 1)
	firstAddedLineNumber := -1
	for _, line := range lines {
		lineNumber++

		if len(line) == 0 {
			continue
		}
		if line[0] == '+' && firstAddedLineNumber == -1 {
			firstAddedLineNumber = lineNumber
		}
		if line[0] == '-' {
			return lineNumber, firstAddedLineNumber, nil
		}
	}

	return 0, firstAddedLineNumber, fmt.Errorf("didn't find deletion line in hunk %s", string(h.Body))
}

func (g Gofmt) extractIssuesFromPatch(patch string, log logutils.Log) ([]result.Issue, error) {
	diffs, err := diffpkg.ParseMultiFileDiff([]byte(patch))
	if err != nil {
		return nil, fmt.Errorf("can't parse patch: %s", err)
	}

	if len(diffs) == 0 {
		return nil, fmt.Errorf("got no diffs from patch parser: %v", diffs)
	}

	issues := []result.Issue{}
	for _, d := range diffs {
		if len(d.Hunks) == 0 {
			log.Warnf("Got no hunks in diff %+v", d)
			continue
		}

		for _, hunk := range d.Hunks {
			deletedLine, addedLine, err := getFirstDeletedAndAddedLineNumberInHunk(hunk)
			if err != nil {
				if addedLine > 1 {
					deletedLine = addedLine - 1 // use previous line, TODO: use both prev and next lines
				} else {
					deletedLine = 1
				}
			}

			text := "File is not `gofmt`-ed with `-s`"
			if g.UseGoimports {
				text = "File is not `goimports`-ed"
			}
			i := result.Issue{
				FromLinter: g.Name(),
				Pos: token.Position{
					Filename: d.NewName,
					Line:     deletedLine,
				},
				Text: text,
			}
			issues = append(issues, i)
		}
	}

	return issues, nil
}

func (g Gofmt) Run(ctx context.Context, lintCtx *linter.Context) ([]result.Issue, error) {
	var issues []result.Issue

	for _, f := range getAllFileNames(lintCtx) {
		var diff []byte
		var err error
		if g.UseGoimports {
			imports.LocalPrefix = lintCtx.Settings().Goimports.LocalPrefixes
			diff, err = goimportsAPI.Run(f)
		} else {
			diff, err = gofmtAPI.Run(f, lintCtx.Settings().Gofmt.Simplify)
		}
		if err != nil { // TODO: skip
			return nil, err
		}
		if diff == nil {
			continue
		}

		is, err := g.extractIssuesFromPatch(string(diff), lintCtx.Log)
		if err != nil {
			return nil, fmt.Errorf("can't extract issues from gofmt diff output %q: %s", string(diff), err)
		}

		issues = append(issues, is...)
	}

	return issues, nil
}
