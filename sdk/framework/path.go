// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/helper/license"
	"github.com/hashicorp/vault/sdk/logical"
)

// Helper which returns a generic regex string for creating endpoint patterns
// that are identified by the given name in the backends
func GenericNameRegex(name string) string {
	return fmt.Sprintf("(?P<%s>\\w(([\\w-.]+)?\\w)?)", name)
}

// GenericNameWithAtRegex returns a generic regex that allows alphanumeric
// characters along with -, . and @.
func GenericNameWithAtRegex(name string) string {
	return fmt.Sprintf("(?P<%s>\\w(([\\w-.@]+)?\\w)?)", name)
}

// Helper which returns a regex string for optionally accepting the a field
// from the API URL
func OptionalParamRegex(name string) string {
	return fmt.Sprintf("(/(?P<%s>.+))?", name)
}

// Helper which returns a regex string for capturing an entire endpoint path
// as the given name.
func MatchAllRegex(name string) string {
	return fmt.Sprintf(`(?P<%s>.*)`, name)
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
	//
	// The pattern will automatically have a ^ prepended and a $ appended before
	// use, if these are not already present, so these may be omitted for clarity.
	//
	// If a ListOperation is being defined, the pattern must end with /? to match
	// a trailing slash optionally, as ListOperations are always processed with a
	// trailing slash added to the path if not already present. The match must not
	// require the presence of a trailing slash, as HelpOperations, even for a
	// path which only implements ListOperation, are processed without a trailing
	// slash - so failure to make the trailing slash optional will break the
	// `vault path-help` command for the path.
	Pattern string

	// Fields is the mapping of data fields to a schema describing that
	// field.
	//
	// Field values are obtained from:
	//
	// - Named captures in the Pattern.
	//
	// - Parameters in the HTTP request body, for HTTP methods where a
	//   request body is expected, i.e. PUT/POST/PATCH. The request body is
	//   typically formatted as JSON, though
	//   "application/x-www-form-urlencoded" format can also be accepted.
	//
	// - Parameters in the HTTP URL query-string, for HTTP methods where
	//   there is no request body, i.e. GET/LIST/DELETE. The query-string
	//   is *not* parsed at all for PUT/POST/PATCH requests.
	//
	// Should the same field be specified both as a named capture and as
	// a parameter, the named capture takes precedence, and a warning is
	// returned.
	Fields map[string]*FieldSchema

	// Operations is the set of operations supported and the associated OperationsHandler.
	//
	// If both Create and Update operations are present, documentation and examples from
	// the Update definition will be used. Similarly if both Read and List are present,
	// Read will be used for documentation.
	Operations map[logical.Operation]OperationHandler

	// Callbacks are the set of callbacks that are called for a given
	// operation. If a callback for a specific operation is not present,
	// then logical.ErrUnsupportedOperation is automatically generated.
	//
	// The help operation is the only operation that the Path will
	// automatically handle if the Help field is set. If both the Help
	// field is set and there is a callback registered here, then the
	// callback will be called.
	//
	// Deprecated: Operations should be used instead and will take priority if present.
	Callbacks map[logical.Operation]OperationFunc

	// ExistenceCheck, if implemented, is used to query whether a given
	// resource exists or not. This is used for ACL purposes: if an Update
	// action is specified, and the existence check returns false, the action
	// is not allowed since the resource must first be created. The reverse is
	// also true. If not specified, the Update action is forced and the user
	// must have UpdateCapability on the path.
	ExistenceCheck ExistenceFunc

	// FeatureRequired, if implemented, will validate if the given feature is
	// enabled for the set of paths
	FeatureRequired license.Features

	// Deprecated denotes that this path is considered deprecated. This may
	// be reflected in help and documentation.
	Deprecated bool

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

	// DisplayAttrs provides hints for UI and documentation generators. They
	// will be included in OpenAPI output if set.
	DisplayAttrs *DisplayAttributes

	// TakesArbitraryInput is used for endpoints that take arbitrary input, instead
	// of or as well as their Fields. This is taken into account when printing
	// warnings about ignored fields. If this is set, we will not warn when data is
	// provided that is not part of the Fields declaration.
	TakesArbitraryInput bool
}

// OperationHandler defines and describes a specific operation handler.
type OperationHandler interface {
	Handler() OperationFunc
	Properties() OperationProperties
}

// OperationProperties describes an operation for documentation, help text,
// and other clients. A Summary should always be provided, whereas other
// fields can be populated as needed.
type OperationProperties struct {
	// Summary is a brief (usually one line) description of the operation.
	Summary string

	// Description is extended documentation of the operation and may contain
	// Markdown-formatted text markup.
	Description string

	// Examples provides samples of the expected request data. The most
	// relevant example should be first in the list, as it will be shown in
	// documentation that supports only a single example.
	Examples []RequestExample

	// Responses provides a list of response description for a given response
	// code. The most relevant response should be first in the list, as it will
	// be shown in documentation that only allows a single example.
	Responses map[int][]Response

	// Unpublished indicates that this operation should not appear in public
	// documentation or help text. The operation may still have documentation
	// attached that can be used internally.
	Unpublished bool

	// Deprecated indicates that this operation should be avoided.
	Deprecated bool

	// The ForwardPerformance* parameters tell the router to unconditionally forward requests
	// to this path if the processing node is a performance secondary/standby. This is generally
	// *not* needed as there is already handling in place to automatically forward requests
	// that try to write to storage. But there are a few cases where explicit forwarding is needed,
	// for example:
	//
	// * The handler makes requests to other systems (e.g. an external API, database, ...) that
	//   change external state somehow, and subsequently writes to storage. In this case the
	//   default forwarding logic could result in multiple mutative calls to the external system.
	//
	// * The operation spans multiple requests (e.g. an OIDC callback), in-memory caching used,
	//   and the same node (and therefore cache) should process both steps.
	//
	// If explicit forwarding is needed, it is usually true that forwarding from both performance
	// standbys and performance secondaries should be enabled.
	//
	// ForwardPerformanceStandby indicates that this path should not be processed
	// on a performance standby node, and should be forwarded to the active node instead.
	ForwardPerformanceStandby bool

	// ForwardPerformanceSecondary indicates that this path should not be processed
	// on a performance secondary node, and should be forwarded to the active node instead.
	ForwardPerformanceSecondary bool

	// DisplayAttrs provides hints for UI and documentation generators. They
	// will be included in OpenAPI output if set.
	DisplayAttrs *DisplayAttributes
}

type DisplayAttributes struct {
	// Name is the name of the field suitable as a label or documentation heading.
	Name string `json:"name,omitempty"`

	// Description of the field that renders as tooltip help text beside the label (name) in the UI.
	// This may be used to replace descriptions that reference comma separation but correspond
	// to UI inputs where only arrays are valid. For example params with Type: framework.TypeCommaStringSlice
	Description string `json:"description,omitempty"`

	// Value is a sample value to display for this field. This may be used
	// to indicate a default value, but it is for display only and completely separate
	// from any Default member handling.
	Value interface{} `json:"value,omitempty"`

	// Sensitive indicates that the value should be masked by default in the UI.
	Sensitive bool `json:"sensitive,omitempty"`

	// Navigation indicates that the path should be available as a navigation tab
	Navigation bool `json:"navigation,omitempty"`

	// ItemType is the type of item this path operates on
	ItemType string `json:"itemType,omitempty"`

	// Group is the suggested UI group to place this field in.
	Group string `json:"group,omitempty"`

	// Action is the verb to use for the operation.
	Action string `json:"action,omitempty"`

	// OperationPrefix is a hyphenated lower-case string used to construct
	// OpenAPI OperationID (prefix + verb + suffix). OperationPrefix is
	// typically a human-readable name of the plugin or a prefix shared by
	// multiple related endpoints.
	OperationPrefix string `json:"operationPrefix,omitempty"`

	// OperationVerb is a hyphenated lower-case string used to construct
	// OpenAPI OperationID (prefix + verb + suffix). OperationVerb is typically
	// an action to be performed (e.g. "generate", "sign", "login", etc.). If
	// not specified, the verb defaults to `logical.Operation.String()`
	// (e.g. "read", "list", "delete", "write" for Create/Update)
	OperationVerb string `json:"operationVerb,omitempty"`

	// OperationSuffix is a hyphenated lower-case string used to construct
	// OpenAPI OperationID (prefix + verb + suffix). It is typically the name
	// of the resource on which the action is performed (e.g. "role",
	// "credentials", etc.). A pipe (|) separator can be used to list different
	// suffixes for various permutations of the `Path.Pattern` regular
	// expression. If not specified, the suffix defaults to the `Path.Pattern`
	// split by dashes.
	OperationSuffix string `json:"operationSuffix,omitempty"`

	// EditType is the optional type of form field needed for a property
	// This is only necessary for a "textarea" or "file"
	EditType string `json:"editType,omitempty"`
}

// RequestExample is example of request data.
type RequestExample struct {
	Description string                 // optional description of the request
	Data        map[string]interface{} // map version of sample JSON request data

	// Optional example response to the sample request. This approach is considered
	// provisional for now, and this field may be changed or removed.
	Response *Response
}

// Response describes and optional demonstrations an operation response.
type Response struct {
	Description string                  // summary of the response and should always be provided
	MediaType   string                  // media type of the response, defaulting to "application/json" if empty
	Fields      map[string]*FieldSchema // the fields present in this response, used to generate openapi response
	Example     *logical.Response       // example response data
}

// PathOperation is a concrete implementation of OperationHandler.
type PathOperation struct {
	Callback                    OperationFunc
	Summary                     string
	Description                 string
	Examples                    []RequestExample
	Responses                   map[int][]Response
	Unpublished                 bool
	Deprecated                  bool
	ForwardPerformanceSecondary bool
	ForwardPerformanceStandby   bool
	DisplayAttrs                *DisplayAttributes
}

func (p *PathOperation) Handler() OperationFunc {
	return p.Callback
}

func (p *PathOperation) Properties() OperationProperties {
	return OperationProperties{
		Summary:                     strings.TrimSpace(p.Summary),
		Description:                 strings.TrimSpace(p.Description),
		Responses:                   p.Responses,
		Examples:                    p.Examples,
		Unpublished:                 p.Unpublished,
		Deprecated:                  p.Deprecated,
		ForwardPerformanceSecondary: p.ForwardPerformanceSecondary,
		ForwardPerformanceStandby:   p.ForwardPerformanceStandby,
		DisplayAttrs:                p.DisplayAttrs,
	}
}

func (p *Path) helpCallback(b *Backend) OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
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
		for k := range p.Fields {
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
				Deprecated:  schema.Deprecated,
			}
		}

		help, err := executeTemplate(pathHelpTemplate, &tplData)
		if err != nil {
			return nil, errwrap.Wrapf("error executing template: {{err}}", err)
		}

		// The plugin type (e.g. "kv", "cubbyhole") is only assigned at the time
		// the plugin is enabled (mounted). If specified in the request, the type
		// will be used as part of the request/response names in the OAS document
		var requestResponsePrefix string
		if v, ok := req.Data["requestResponsePrefix"]; ok {
			requestResponsePrefix = v.(string)
		}

		// Build OpenAPI response for this path
		vaultVersion := "unknown"
		if b.System() != nil {
			// b.System() should always be non-nil, except tests might create a
			// Backend without one.
			env, err := b.System().PluginEnv(context.Background())
			if err != nil {
				return nil, err
			}
			if env != nil {
				vaultVersion = env.VaultVersion
			}
		}
		redactVersion, _, _, _ := logical.CtxRedactionSettingsValue(ctx)
		if redactVersion {
			vaultVersion = ""
		}
		doc := NewOASDocument(vaultVersion)
		if err := documentPath(p, b, requestResponsePrefix, doc); err != nil {
			b.Logger().Warn("error generating OpenAPI", "error", err)
		}

		return logical.HelpResponse(help, nil, doc), nil
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
	Deprecated  bool
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
{{if .Deprecated}}
{{printf "(DEPRECATED) %s" .Description | indent 8}}
{{else}}
{{indent 8 .Description}}
{{end}}{{end}}{{end}}
## DESCRIPTION

{{.Description}}
`
