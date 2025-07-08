// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package github

import (
	"bytes"
	"embed"
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
