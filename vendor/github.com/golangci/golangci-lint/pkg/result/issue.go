package result

import "go/token"

type Range struct {
	From, To int
}

type Issue struct {
	FromLinter string
	Text       string

	Pos       token.Position
	LineRange *Range `json:",omitempty"`
	HunkPos   int    `json:",omitempty"`

	SourceLines []string
}

func (i *Issue) FilePath() string {
	return i.Pos.Filename
}

func (i *Issue) Line() int {
	return i.Pos.Line
}

func (i *Issue) Column() int {
	return i.Pos.Column
}

func (i *Issue) GetLineRange() Range {
	if i.LineRange == nil {
		return Range{
			From: i.Line(),
			To:   i.Line(),
		}
	}

	if i.LineRange.From == 0 {
		return Range{
			From: i.Line(),
			To:   i.Line(),
		}
	}

	return *i.LineRange
}
