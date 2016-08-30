package framework

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/vault/logical"
)

// Helper which returns a generic regex string for creating endpoint patterns
// that are identified by the given name in the backends
func GenericNameRegex(name string) string {
	return fmt.Sprintf("(?P<%s>\\w(([\\w-.]+)?\\w)?)", name)
}

// Helper which returns a regex string for optionally accepting the a field
// from the API URL
func OptionalParamRegex(name string) string {
	return fmt.Sprintf("(/(?P<%s>.+))?", name)
}

// PathAppend is a helper for appending lists of paths into a single
// list.
func PathAppend(paths ...[]*Path) []*Path {
	result := make([]*Path, 0, 10)
	for _, ps := range paths {
		result = append(result, ps...)
	}

	return result
}

// Path is a single path that the backend responds to.
type Path struct {
	// Pattern is the pattern of the URL that matches this path.
	//
	// This should be a valid regular expression. Named captures will be
	// exposed as fields that should map to a schema in Fields. If a named
	// capture is not a field in the Fields map, then it will be ignored.
	Pattern string

	// Fields is the mapping of data fields to a schema describing that
	// field. Named captures in the Pattern also map to fields. If a named
	// capture name matches a PUT body name, the named capture takes
	// priority.
	//
	// Note that only named capture fields are available in every operation,
	// whereas all fields are available in the Write operation.
	Fields map[string]*FieldSchema

	// Callbacks are the set of callbacks that are called for a given
	// operation. If a callback for a specific operation is not present,
	// then logical.ErrUnsupportedOperation is automatically generated.
	//
	// The help operation is the only operation that the Path will
	// automatically handle if the Help field is set. If both the Help
	// field is set and there is a callback registered here, then the
	// callback will be called.
	Callbacks map[logical.Operation]OperationFunc

	// ExistenceCheck, if implemented, is used to query whether a given
	// resource exists or not. This is used for ACL purposes: if an Update
	// action is specified, and the existence check returns false, the action
	// is not allowed since the resource must first be created. The reverse is
	// also true. If not specified, the Update action is forced and the user
	// must have UpdateCapability on the path.
	ExistenceCheck func(*logical.Request, *FieldData) (bool, error)

	// Help is text describing how to use this path. This will be used
	// to auto-generate the help operation. The Path will automatically
	// generate a parameter listing and URL structure based on the
	// regular expression, so the help text should just contain a description
	// of what happens.
	//
	// HelpSynopsis is a one-sentence description of the path. This will
	// be automatically line-wrapped at 80 characters.
	//
	// HelpDescription is a long-form description of the path. This will
	// be automatically line-wrapped at 80 characters.
	HelpSynopsis    string
	HelpDescription string
}

func (p *Path) helpCallback(
	req *logical.Request, data *FieldData) (*logical.Response, error) {
	var tplData pathTemplateData
	tplData.Request = req.Path
	tplData.RoutePattern = p.Pattern
	tplData.Synopsis = strings.TrimSpace(p.HelpSynopsis)
	if tplData.Synopsis == "" {
		tplData.Synopsis = "<no synopsis>"
	}
	tplData.Description = strings.TrimSpace(p.HelpDescription)
	if tplData.Description == "" {
		tplData.Description = "<no description>"
	}

	// Alphabetize the fields
	fieldKeys := make([]string, 0, len(p.Fields))
	for k, _ := range p.Fields {
		fieldKeys = append(fieldKeys, k)
	}
	sort.Strings(fieldKeys)

	// Build the field help
	tplData.Fields = make([]pathTemplateFieldData, len(fieldKeys))
	for i, k := range fieldKeys {
		schema := p.Fields[k]
		description := strings.TrimSpace(schema.Description)
		if description == "" {
			description = "<no description>"
		}

		tplData.Fields[i] = pathTemplateFieldData{
			Key:         k,
			Type:        schema.Type.String(),
			Description: description,
		}
	}

	help, err := executeTemplate(pathHelpTemplate, &tplData)
	if err != nil {
		return nil, fmt.Errorf("error executing template: %s", err)
	}

	return logical.HelpResponse(help, nil), nil
}

type pathTemplateData struct {
	Request      string
	RoutePattern string
	Synopsis     string
	Description  string
	Fields       []pathTemplateFieldData
}

type pathTemplateFieldData struct {
	Key         string
	Type        string
	Description string
	URL         bool
}

const pathHelpTemplate = `
Request:        {{.Request}}
Matching Route: {{.RoutePattern}}

{{.Synopsis}}

{{ if .Fields -}}
## PARAMETERS
{{range .Fields}}
{{indent 4 .Key}} ({{.Type}})
{{indent 8 .Description}}
{{end}}{{end}}
## DESCRIPTION

{{.Description}}
`
