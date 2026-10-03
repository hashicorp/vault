// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

// renderEmbeddedTemplateToTmpFile renders an embedded template to a temporary
// file on disk and returns the path. The caller is responsible for handling
// the file thereafter.
func renderEmbeddedTemplateToTmpFile(name string, data any) (*os.File, error) {
	s, err := os.CreateTemp("", name)
	if err != nil {
		return nil, err
	}

	err = renderEmbeddedTemplateTo(name, data, s)
	if err != nil {
		return nil, err
	}

	// Rename the file as it forces writes to be flushed
	dst := s.Name() + "d"
	err = os.Rename(s.Name(), dst)
	if err != nil {
		return nil, err
	}

	return os.Open(dst)
}

// renderEmbeddedTemplateTo renders an embedded template to an io.Writer. The
// caller is responsible for closing the writer.
func renderEmbeddedTemplateTo(name string, data any, writer io.Writer) error {
	body, err := renderEmbeddedTemplate(name, data)
	if err != nil {
		return err
	}

	_, err = io.WriteString(writer, body)

	return err
}

// renderEmbeddedTemplate renders an embedded template to a string
func renderEmbeddedTemplate(name string, data any) (string, error) {
	tmpl, err := templates.ReadFile("templates/" + name)
	if err != nil {
		return "", err
	}

	t, err := template.New(name).Parse(string(tmpl))
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// limitCharacters truncates a string to maxGitHubMessageChars while preserving
// markdown formatting integrity. It adds a truncation notice if truncation occurs.
const maxGitHubMessageChars = 65536

func limitCharacters(content string) string {
	if len(content) <= maxGitHubMessageChars {
		return content
	}

	// Define the truncation notice to calculate its length
	truncationNotice := "\n\n---\n\n:scissors: **Message truncated due to GitHub's character limit**\n\n" +
		fmt.Sprintf("This message was automatically truncated because it exceeded GitHub's %d character limit for "+
			"comments and pull request descriptions.", maxGitHubMessageChars)

	// Calculate how much space the notice needs
	noticeLen := len(truncationNotice)
	maxContentLen := maxGitHubMessageChars - noticeLen

	// Find a good truncation point to avoid breaking markdown
	truncateAt := maxContentLen
	for truncateAt > 0 && content[truncateAt] != '\n' {
		truncateAt--
	}

	// If we can't find a newline within reasonable bounds, just truncate
	if truncateAt == 0 || (maxContentLen-truncateAt) > 1000 {
		truncateAt = maxContentLen
	}

	// Ensure we don't go out of bounds
	if truncateAt > len(content) {
		truncateAt = len(content)
	}

	return content[:truncateAt] + truncationNotice
}
