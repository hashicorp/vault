package framework

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

func executeTemplate(tpl string, data interface{}) (string, error) {
	// Parse the help template
	t, err := template.New("root").Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %s", err)
	}

	// Execute the template and store the output
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error executing template: %s", err)
	}

	return strings.TrimSpace(buf.String()), nil
}
