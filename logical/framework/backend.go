package framework

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"sync"
	"text/template"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/go-wordwrap"
)

// Backend is an implementation of logical.Backend that allows
// the implementer to code a backend using a much more programmer-friendly
// framework that handles a lot of the routing and validation for you.
//
// This is recommended over implementing logical.Backend directly.
type Backend struct {
	// Paths are the various routes that the backend responds to.
	// This cannot be modified after construction (i.e. dynamically changing
	// paths, including adding or removing, is not allowed once the
	// backend is in use).
	Paths []*Path

	// PathsRoot is the list of path patterns that denote the
	// paths above that require root-level privileges. These can't be
	// regular expressions, it is either exact match or prefix match.
	// For prefix match, append '*' as a suffix.
	PathsRoot []string

	// Rollback is called when a WAL entry (see wal.go) has to be rolled
	// back. It is called with the data from the entry. Boolean true should
	// be returned on success. Errors should just be logged.
	Rollback       func(kind string, data interface{}) bool
	RollbackMinAge time.Duration

	once    sync.Once
	pathsRe []*regexp.Regexp
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
	// whereas all fields are avaiable in the Write operation.
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

// OperationFunc is the callback called for an operation on a path.
type OperationFunc func(*logical.Request, *FieldData) (*logical.Response, error)

// logical.Backend impl.
func (b *Backend) HandleRequest(req *logical.Request) (*logical.Response, error) {
	// Rollbacks are treated outside of the normal request cycle since
	// the path doesn't matter for them.
	if req.Operation == logical.RollbackOperation {
		return b.handleRollback(req)
	}

	// Find the matching route
	path, captures := b.route(req.Path)
	if path == nil {
		return nil, logical.ErrUnsupportedPath
	}

	// Build up the data for the route, with the URL taking priority
	// for the fields over the PUT data.
	raw := make(map[string]interface{}, len(path.Fields))
	for k, v := range req.Data {
		raw[k] = v
	}
	for k, v := range captures {
		raw[k] = v
	}

	// Look up the callback for this operation
	var callback OperationFunc
	var ok bool
	if path.Callbacks != nil {
		callback, ok = path.Callbacks[req.Operation]
	}
	if !ok {
		if req.Operation == logical.HelpOperation && path.HelpSynopsis != "" {
			callback = path.helpCallback
			ok = true
		}
	}
	if !ok {
		return nil, logical.ErrUnsupportedOperation
	}

	// Call the callback with the request and the data
	return callback(req, &FieldData{
		Raw:    raw,
		Schema: path.Fields,
	})
}

// logical.Backend impl.
func (b *Backend) RootPaths() []string {
	return b.PathsRoot
}

// Route looks up the path that would be used for a given path string.
func (b *Backend) Route(path string) *Path {
	result, _ := b.route(path)
	return result
}

func (b *Backend) init() {
	b.pathsRe = make([]*regexp.Regexp, len(b.Paths))
	for i, p := range b.Paths {
		b.pathsRe[i] = regexp.MustCompile(p.Pattern)
	}
}

func (b *Backend) route(path string) (*Path, map[string]string) {
	b.once.Do(b.init)

	for i, re := range b.pathsRe {
		matches := re.FindStringSubmatch(path)
		if matches == nil {
			continue
		}

		// We have a match, determine the mapping of the captures and
		// store that for returning.
		var captures map[string]string
		path := b.Paths[i]
		if captureNames := re.SubexpNames(); len(captureNames) > 1 {
			captures = make(map[string]string, len(captureNames))
			for i, name := range captureNames {
				if name != "" {
					captures[name] = matches[i]
				}
			}
		}

		return path, captures
	}

	return nil, nil
}

func (b *Backend) handleRollback(
	req *logical.Request) (*logical.Response, error) {
	if b.Rollback == nil {
		return nil, logical.ErrUnsupportedOperation
	}

	var merr error
	keys, err := ListWAL(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if len(keys) == 0 {
		return nil, nil
	}

	// Calculate the minimum time that the WAL entries could be
	// created in order to be rolled back.
	age := b.RollbackMinAge
	if age == 0 {
		age = 10 * time.Minute
	}
	minAge := time.Now().UTC().Add(-1 * age)

	for _, k := range keys {
		entry, err := GetWAL(req.Storage, k)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}
		if entry == nil {
			continue
		}

		// If the entry isn't old enough, then don't roll it back
		if !time.Unix(entry.CreatedAt, 0).Before(minAge) {
			continue
		}

		// Attempt a rollback
		if b.Rollback(entry.Kind, entry.Data) {
			if err := DeleteWAL(req.Storage, k); err != nil {
				merr = multierror.Append(merr, err)
			}
		}
	}

	if merr == nil {
		return nil, nil
	}

	return logical.ErrorResponse(merr.Error()), nil
}

func (p *Path) helpCallback(req *logical.Request, data *FieldData) (*logical.Response, error) {
	var tplData pathTemplateData
	tplData.Request = req.Path
	tplData.RoutePattern = p.Pattern
	tplData.Synopsis = wordwrap.WrapString(p.HelpSynopsis, 80)
	tplData.Description = wordwrap.WrapString(p.HelpDescription, 80)

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
		description := wordwrap.WrapString(schema.Description, 60)
		if description == "" {
			description = "<no description>"
		}

		tplData.Fields[i] = pathTemplateFieldData{
			Key:         k,
			Type:        schema.Type.String(),
			Description: description,
		}
	}

	// Parse the help template
	tpl, err := template.New("root").Parse(pathHelpTemplate)
	if err != nil {
		return nil, fmt.Errorf("error parsing template: %s", err)
	}

	// Execute the template and store the output
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, &tplData); err != nil {
		return nil, fmt.Errorf("error executing template: %s", err)
	}

	return logical.HelpResponse(buf.String(), nil), nil
}

// FieldSchema is a basic schema to describe the format of a path field.
type FieldSchema struct {
	Type        FieldType
	Default     interface{}
	Description string
}

// DefaultOrZero returns the default value if it is set, or otherwise
// the zero value of the type.
func (s *FieldSchema) DefaultOrZero() interface{} {
	if s.Default != nil {
		return s.Default
	}

	return s.Type.Zero()
}

func (t FieldType) Zero() interface{} {
	switch t {
	case TypeString:
		return ""
	case TypeInt:
		return 0
	case TypeBool:
		return false
	default:
		panic("unknown type: " + t.String())
	}
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

## Parameters

{{range .Fields}}
### {{.Key}} (type: {{.Type}})

{{.Description}}

{{end}}
## Description

{{.Description}}
`
