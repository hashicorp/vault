package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/kr/text"

	"github.com/hashicorp/vault/sdk/framework"
)

const (
	doNotEdit     = `<!-- This section is auto-generated and must be updated in the Go source file. -->`
	maxLineLength = 100
	hangingIndent = 2
)

// Don't use the default flag set since it will pick up all of the -test
// flags that get pulled in when importing the logical packages.
var (
	flags      = flag.NewFlagSet("apidoc", flag.ExitOnError)
	sourceFile = flags.String("f", "openapi.json", "OpenAPI reference file")
	root       = flags.String("r", "website/source/api", "Root folder for documentation source")
	verbose    = flags.Bool("v", false, "Verbose output")
	analyze    = flags.Bool("a", false, "Analyze auto-generate coverage")
)

var typeDefaults = map[string]string{
	"string":  `""`,
	"integer": `0`,
	"array":   `[]`,
	"object":  `{}`,
}

func main() {
	flags.Parse(os.Args[1:])

	oas, err := loadOpenAPI(*sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := processTree(*root, oas); err != nil {
		fmt.Println(err)
	}
}

// loadOpenAPI parses an OpenAPI JSON file into an OASDocument.
func loadOpenAPI(filename string) (*framework.OASDocument, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var d map[string]interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}

	oas, err := framework.NewOASDocumentFromMap(d)
	if err != nil {
		return nil, err
	}

	return oas, nil
}

// processTree processes all Markdown files under root
func processTree(root string, oas *framework.OASDocument) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".md") {
			return processFile(path, oas)
		}
		return nil
	})

	return err
}

// processFile parses and updates a single Markdown file.
func processFile(filename string, oas *framework.OASDocument) error {
	var totalLines, autoLines int
	var output []string

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// exit now if there are no autogen tags
	if !*analyze && bytes.Index(data, []byte("<!--apidoc_start")) == -1 {
		return nil
	}

	lines := strings.Split(string(data), "\n")

	blockStart := -1
	blockFound := false
	for i, line := range lines {
		totalLines++
		if strings.HasPrefix(line, "<!--apidoc_start") && strings.HasSuffix(line, "-->") {
			if blockStart != -1 {
				return fmt.Errorf("%s:%d: error: nested documentation block not allowed", filename, i+1)
			}
			blockStart = i
		}
		if strings.HasPrefix(line, "<!--apidoc_end") && strings.HasSuffix(line, "-->") {
			if blockStart == -1 {
				return fmt.Errorf("%s:%d: error: end of block without start", filename, i+1)
			}

			blockFound = true
			origLines := lines[blockStart : i+1]
			newLines, err := generateBlock(origLines, oas, filename, blockStart+1)
			if err != nil {
				return err
			}

			output = append(output, newLines...)
			blockStart = -1
			totalLines += (len(newLines) - len(origLines))
			autoLines += len(newLines)
			continue
		}
		if blockStart == -1 {
			output = append(output, line)
		}
	}

	if blockStart != -1 {
		return fmt.Errorf("%s:%d: error: end of file reached without closing tag", filename, blockStart+1)
	}

	if *analyze {
		fmt.Printf("%2.f%% %s\n", math.Round(100.0*float64(autoLines)/float64(totalLines)), filename)
		return nil
	}

	if blockFound {
		outputData := []byte(strings.Join(output, "\n"))
		if !bytes.Equal(outputData, data) {
			fmt.Println("Updating " + filename)
			if err := ioutil.WriteFile(filename, outputData, 0644); err != nil {
				return err
			}
		} else {
			if *verbose {
				fmt.Printf("Processed %s (no change)\n", filename)
				fmt.Println(autoLines, totalLines)
			}
		}
	}

	return nil
}

// generateBlock outputs Markdown content for a since apidoc tag
func generateBlock(lines []string, oas *framework.OASDocument, filename string, offset int) ([]string, error) {
	output := []string{
		lines[0],
		doNotEdit,
		"",
	}

	// Parse apidoc tag configuration
	var path, method string
	re := regexp.MustCompile(`(path|method|disable):(\S+)`)

	matches := re.FindAllStringSubmatch(lines[0], -1)
	outputTypes := map[string]bool{
		"description": true,
		"table":       true,
		"parameters":  true,
		"payload":     true,
		"request":     true,
	}

	for _, m := range matches {
		switch m[1] {
		case "path":
			path = m[2]
		case "method":
			method = m[2]
		case "disable":
			o := strings.TrimSuffix(m[2], "-->")
			for _, t := range strings.Split(o, ",") {
				switch t {
				case "description", "table", "parameters", "payload", "request":
					outputTypes[t] = false
				default:
					return nil, fmt.Errorf("%s:%d: unexpected 'disable' value: %q", filename, offset, t)
				}
			}
		}
	}

	p, ok := oas.Paths[path]
	if !ok {
		return nil, fmt.Errorf("%s:%d: path '%s' not found\n", filename, offset, path)
	}

	var o *framework.OASOperation
	switch method {
	case "post":
		o = p.Post
	case "get":
		o = p.Get
	case "delete":
		o = p.Delete
	}

	if o == nil {
		return nil, fmt.Errorf("%s:%d: method '%s' for path '%s' not found\n", filename, offset, method, path)
	}

	// Generate description
	if outputTypes["description"] {
		if o.Description != "" {
			output = append(output, wrapAtLengthWithPadding(o.Description, maxLineLength, 0))
			output = append(output, "")
		}
	}

	// Generate summary table
	if outputTypes["table"] {
		output = append(output, table(
			fmt.Sprintf("`%s`", strings.ToUpper(method)),
			fmt.Sprintf("`%s`", path))...)
	}

	// Generate parameter list
	if outputTypes["parameters"] {

		// Extract path-level parameters.
		for _, p := range p.Parameters {
			// TODO: handle query parameters
			line := fmt.Sprintf("- `%s` (`%s: `<required>`) - %s", p.Name, p.Schema.Type, formatParamDescription(p.Description))
			line = wrapAtLengthWithPadding(line, maxLineLength, hangingIndent)
			output = append(output, line)
		}

		// Extract request body parameters for POST.
		if method == "post" {
			schema := o.RequestBody.Content["application/json"].Schema
			var properties []string
			for p := range schema.Properties {
				properties = append(properties, p)
			}
			sortProperties(properties, lines)

			for _, name := range properties {
				s := schema.Properties[name]

				dflt := "<required>"
				if !contains(schema.Required, name) {
					if s.Default == nil {
						dflt = typeDefaults[s.Type]
						if dflt == "" {
							return nil, fmt.Errorf("unexpected type %q in OpenAPI definition for path %q", s.Type, path)
						}
					} else {
						dflt = fmt.Sprintf("%v", s.Default)
					}
				}

				line := fmt.Sprintf("- `%s` (`%s: %s`) - %s", name, s.Type, dflt, formatParamDescription(s.Description))
				line = wrapAtLengthWithPadding(line, maxLineLength, hangingIndent)
				output = append(output, line)
			}
		}
	}

	schema := o.RequestBody.Content["application/json"].Schema
	if schema.Example != nil {
		m, err := json.MarshalIndent(schema.Example, "", "  ")
		if err != nil {
			return nil, err
		}

		if outputTypes["payload"] {
			output = append(output,
				"",
				"### Sample Payload",
				"",
				"```json",
				string(m),
				"```")
		}

		if outputTypes["request"] {
			output = append(output,
				"",
				"### Sample Request",
				"",
				"```",
				"$ curl \\",
				"    --header \"X-Vault-Token: ...\" ",
				"    --request POST \\",
				"    --data @payload.json \\",
				fmt.Sprintf("    https://127.0.0.1:8200/v1%s", path),
				"```")
		}
	}

	output = append(output, "", "<!--apidoc_end -->")

	var final []string
	for _, el := range output {
		final = append(final, strings.Split(el, "\n")...)
	}

	return final, nil
}

// table generates a Markdown table for a method & path
func table(method, path string) []string {
	var output bytes.Buffer

	w := tabwriter.NewWriter(&output, 0, 0, 1, ' ', 0)
	fmt.Fprintln(w, "|\tMethod\t|\tPath\t|")
	fmt.Fprintln(w, fmt.Sprintf("|\t:%s\t|\t:%s\t|",
		strings.Repeat("-", len(method)),
		strings.Repeat("-", len(path))))

	fmt.Fprintln(w, fmt.Sprintf("|\t%s\t|\t%s\t|", method, path))
	w.Flush()
	return strings.Split(output.String(), "\n")
}

// sortParameters sorts a list of properties using an existing set of lines
// (usually the current parameter list) as a reference to maintain the current order.
func sortProperties(properties, lines []string) {
	re := regexp.MustCompile("- +`([\\w_-]+)`")

	var existing []string
	for _, line := range lines {
		m := re.FindAllStringSubmatch(line, 1)
		if m != nil {
			existing = append(existing, m[0][1])
		}
	}

	sort.Slice(properties, func(i, j int) bool {
		for _, e := range existing {
			switch {
			case e == properties[i]:
				return true
			case e == properties[j]:
				return false
			}
		}
		return i < j
	})
}

var paragraphSplit = regexp.MustCompile(`\n\n+`)

// wrapAtLengthWithPadding wraps the given text at the maxLineLength, taking
// into account any provided left padding and retaining paragraph splits.
func wrapAtLengthWithPadding(s string, maxLineLength, hangingIndent int) string {
	paragraphs := paragraphSplit.Split(s, -1)

	for i, p := range paragraphs {
		wrapped := text.Wrap(p, maxLineLength-hangingIndent)
		lines := strings.Split(wrapped, "\n")
		indent := strings.Repeat(" ", hangingIndent)
		for i, line := range lines {
			if i > 0 {
				lines[i] = indent + line
			}
		}
		paragraphs[i] = strings.Join(lines, "\n")
	}

	return strings.Join(paragraphs, "\n\n")
}

func formatParamDescription(s string) string {
	if !strings.HasSuffix(s, ".") {
		s += "."
	}
	return s
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
