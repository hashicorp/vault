// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"text/template"
)

// workflowRunTemplate is our template for rendering workflow runs in human
// readable text.
var workflowRunTemplate = template.Must(template.New("workflow_run").Funcs(template.FuncMap{
	"boldify":              boldify,
	"format_log_lines":     formatLogLines,
	"intensify_status":     intensifyStatus,
	"intensify_annotation": intensifyAnnotationLevel,
	"redify":               redify,
	"splitlines":           splitLines,
}).ParseFS(templates, "templates/workflow-run-text.tmpl"))

func summarizeWorkflowRun(r *WorkflowRun) (string, error) {
	if r == nil {
		return "", errors.New("uninitialized workflow run")
	}

	if r.summary != "" {
		return r.summary, nil
	}

	b := &bytes.Buffer{}
	err := workflowRunTemplate.Execute(b, r)
	if err != nil {
		return "", err
	}
	r.summary = b.String()

	return r.summary, nil
}

func formatLogLines(log []byte) []string {
	if len(log) == 0 {
		return nil
	}

	lines := []string{}
	lastLine := ""
	for _, line := range strings.Split(string(log), "\n") {
		// Strip out duplicate lines and blank lines so the summaries
		// are more clear. Log lines include timestamps with microseconds
		// so we compare for duplicates by comparing line values without
		// the timestamp.
		newLineNoTimestamp := ""
		newLineParts := strings.SplitN(line, " ", 2)
		if len(newLineParts) < 2 {
			// blank lines may only have a timestamp
			continue
		}
		newLineNoTimestamp = newLineParts[1]
		if strings.TrimSpace(newLineNoTimestamp) == "" {
			// Don't write otherwise blank lines
			continue
		}
		if lastLine == newLineNoTimestamp {
			continue
		}

		lines = append(lines, line)
		lastLine = newLineNoTimestamp
	}

	return lines
}

func intensifyStatus(in string) string {
	switch in {
	case "completed", "success":
		return "\x1b[1;32;49m" + in + "\x1b[0m"
	case "cancelled":
		return "\x1b[1;33;49m" + in + "\x1b[0m"
	case "failure":
		return "\x1b[1;31;49m" + in + "\x1b[0m"
	case "skipped":
		return "\x1b[1;37;49m" + in + "\x1b[0m"
	case "in_progress":
		return "\x1b[1;37;49m" + in + "\x1b[0m"
	case "warning":
		return "\x1b[1;33;49m" + in + "\x1b[0m"
	default:
		return in
	}
}

func intensifyAnnotationLevel(in string) string {
	switch in {
	case "failure":
		return "\x1b[1;31;49m" + in + "\x1b[0m"
	case "warning":
		return "\x1b[1;33;49m" + in + "\x1b[0m"
	default:
		return in
	}
}

func boldify(in string) string {
	return "\x1b[1;39m" + in + "\x1b[0m"
}

func redify(in string) string {
	return "\x1b[1;31m" + in + "\x1b[0m"
}

func splitLines(in string) []string {
	return slices.DeleteFunc(strings.Split(in, "\n"), func(s string) bool {
		if s == "\n" || s == "" {
			return true
		}
		return false
	})
}
