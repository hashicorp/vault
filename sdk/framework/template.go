package framework

import (
	"bufio"
	"bytes"
	"strings"
	"text/template"

	"github.com/hashicorp/errwrap"
)

func executeTemplate(tpl string, data interface{}) (string, error) {
	// Define the functions
	funcs := map[string]interface{}{
		"indent": funcIndent,
	}

	// Parse the help template
	t, err := template.New("root").Funcs(funcs).Parse(tpl)
	if err != nil {
		return "", errwrap.Wrapf("error parsing template: {{err}}", err)
	}

	// Execute the template and store the output
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", errwrap.Wrapf("error executing template: {{err}}", err)
	}

	return strings.TrimSpace(buf.String()), nil
}

func funcIndent(count int, text string) string {
	var buf bytes.Buffer
	prefix := strings.Repeat(" ", count)
	scan := bufio.NewScanner(strings.NewReader(text))
	for scan.Scan() {
		buf.WriteString(prefix + scan.Text() + "\n")
	}

	return strings.TrimRight(buf.String(), "\n")
}
