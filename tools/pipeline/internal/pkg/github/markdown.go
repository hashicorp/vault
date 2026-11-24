// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"fmt"
	"html"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

var markdownTableCellReplacer = strings.NewReplacer(
	"\\", "&#92;", "`", "&#96;", "*", "&#42;", "_", "&#95;", "[", "&#91;", "]", "&#93;",
)

// escapeMarkdownTableCell keeps cell text from being treated as Markdown.
// RenderMarkdown already handles pipes and line breaks for us.
func escapeMarkdownTableCell(s string) string {
	return markdownTableCellReplacer.Replace(html.EscapeString(s))
}

// escapeMarkdownTableColumns escapes the columns we'll render as Markdown.
// Leave out columns where we build Markdown links ourselves.
func escapeMarkdownTableColumns(t table.Writer, columns ...int) table.Writer {
	transform := func(value any) string {
		return escapeMarkdownTableCell(fmt.Sprint(value))
	}
	configs := make([]table.ColumnConfig, 0, len(columns))
	for _, column := range columns {
		configs = append(configs, table.ColumnConfig{
			Number: column, Transformer: transform, TransformerFooter: transform,
		})
	}
	t.SetColumnConfigs(configs)
	return t
}
