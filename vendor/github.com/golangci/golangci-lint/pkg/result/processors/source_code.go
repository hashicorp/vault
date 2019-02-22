package processors

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type linesCache [][]byte
type filesLineCache map[string]linesCache

type SourceCode struct {
	cache filesLineCache
	log   logutils.Log
}

var _ Processor = SourceCode{}

func NewSourceCode(log logutils.Log) *SourceCode {
	return &SourceCode{
		cache: filesLineCache{},
		log:   log,
	}
}

func (p SourceCode) Name() string {
	return "source_code"
}

func (p SourceCode) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(i *result.Issue) *result.Issue {
		lines, err := p.getFileLinesForIssue(i)
		if err != nil {
			p.log.Warnf("Failed to get lines for file %s: %s", i.FilePath(), err)
			return i
		}

		newI := *i

		lineRange := i.GetLineRange()
		var lineStr string
		for line := lineRange.From; line <= lineRange.To; line++ {
			if line == 0 { // some linters, e.g. gosec can do it: it really means first line
				line = 1
			}

			zeroIndexedLine := line - 1
			if zeroIndexedLine >= len(lines) {
				p.log.Warnf("No line %d in file %s", line, i.FilePath())
				break
			}

			lineStr = string(bytes.Trim(lines[zeroIndexedLine], "\r"))
			newI.SourceLines = append(newI.SourceLines, lineStr)
		}

		return &newI
	}), nil
}

func (p *SourceCode) getFileLinesForIssue(i *result.Issue) (linesCache, error) {
	fc := p.cache[i.FilePath()]
	if fc != nil {
		return fc, nil
	}

	// TODO: make more optimal algorithm: don't load all files into memory
	fileBytes, err := ioutil.ReadFile(i.FilePath())
	if err != nil {
		return nil, fmt.Errorf("can't read file %s for printing issued line: %s", i.FilePath(), err)
	}
	lines := bytes.Split(fileBytes, []byte("\n")) // TODO: what about \r\n?
	fc = lines
	p.cache[i.FilePath()] = fc
	return fc, nil
}

func (p SourceCode) Finish() {}
