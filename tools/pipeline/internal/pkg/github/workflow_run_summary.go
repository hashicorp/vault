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
}).Parse(workflowRunTextTemplate))

// workflowRunTextTemplate is the actual template text of our human readable
// workflow run output.
const workflowRunTextTemplate = `
{{ .Run.Name }} (ID: {{ .Run.ID }})
  Title: {{ boldify .Run.DisplayTitle }}
  URL: {{ .Run.HTMLURL }}
  HEAD Branch: {{ .Run.HeadBranch }}
  HEAD SHA: {{ .Run.HeadSHA }}
  Author: {{ .Run.HeadCommit.Author.Name }}
  Actor: {{ .Run.Actor.Login }}
  Attempt: {{ .Run.RunAttempt }}
  {{- if .CheckRuns }}
  {{ boldify "Annotations" }}
    {{- range $cr := .CheckRuns }}
    {{- range $a := $cr.Annotations }}
    {{ intensify_annotation $a.AnnotationLevel }}{{- if $a.Title }} {{ boldify $a.Title }}{{- end -}}
    {{- if $a.Message }}
    {{- range $am := splitlines $a.Message }}
      {{ boldify $am }}
    {{- end -}}
    {{- end -}}
    {{- if $a.RawDetails }}
    {{- range $ad := splitlines $a.RawDetails }}
      {{ boldify $ad }}
    {{- end -}}
    {{- end -}}
    {{- end -}}
    {{- end -}}
  {{- end }}
  Status: {{ intensify_status .Run.Status }}
  Conclusion: {{ intensify_status .Run.Conclusion }}
  {{- if .Jobs -}}
  {{- range $j := .Jobs }}
    Job: {{ boldify $j.Job.Name }}
      URL: {{ $j.Job.HTMLURL }}
      Status: {{ $j.Job.Status }}
      Conclusion: {{ $j.Job.Conclusion }}
      CreatedAt: {{ $j.Job.CreatedAt }}
      StartedAt: {{ $j.Job.StartedAt }}
      CompletedAt: {{ $j.Job.CompletedAt }}
      {{- if .UnsuccessfulSteps }}
      {{ boldify "Unsuccessful Steps:" }}
      {{- range $s := .UnsuccessfulSteps }}
        Step: {{ boldify $s.Name }}
          Status: {{ intensify_status $s.Status }}
          Conclusion: {{ intensify_status $s.Conclusion }}
      {{- end -}}
      {{- end -}}
      {{- if .LogEntries }}
      {{ boldify "Unsuccessful Step Log Summaries:" }}
      {{- range $l := .LogEntries }}
        Step: {{ boldify $l.StepName }}
        {{- range $sl := format_log_lines $l.SetupLog }}
          {{ $sl }}
        {{- end -}}
        {{- range $bl := format_log_lines $l.BodyLog }}
          {{ $bl }}
        {{- end -}}
        {{- range $el := format_log_lines $l.ErrorLog }}
          {{ redify $el }}
        {{- end -}}
      {{- end -}}
      {{- end -}}
  {{- end -}}
  {{- end -}}`

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
