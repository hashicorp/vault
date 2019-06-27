package framework

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/wrapping"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/version"
	"github.com/mitchellh/mapstructure"
)

// OpenAPI specification (OAS): https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.2.md
const OASVersion = "3.0.2"

// NewOASDocument returns an empty OpenAPI document.
func NewOASDocument() *OASDocument {
	return &OASDocument{
		Version: OASVersion,
		Info: OASInfo{
			Title:       "HashiCorp Vault API",
			Description: "HTTP API that gives you full access to Vault. All API routes are prefixed with `/v1/`.",
			Version:     version.GetVersion().Version,
			License: OASLicense{
				Name: "Mozilla Public License 2.0",
				URL:  "https://www.mozilla.org/en-US/MPL/2.0",
			},
		},
		Paths: make(map[string]*OASPathItem),
	}
}

// NewOASDocumentFromMap builds an OASDocument from an existing map version of a document.
// If a document has been decoded from JSON or received from a plugin, it will be as a map[string]interface{}
// and needs special handling beyond the default mapstructure decoding.
func NewOASDocumentFromMap(input map[string]interface{}) (*OASDocument, error) {

	// The Responses map uses integer keys (the response code), but once translated into JSON
	// (e.g. during the plugin transport) these become strings. mapstructure will not coerce these back
	// to integers without a custom decode hook.
	decodeHook := func(src reflect.Type, tgt reflect.Type, inputRaw interface{}) (interface{}, error) {

		// Only alter data if:
		//  1. going from string to int
		//  2. string represent an int in status code range (100-599)
		if src.Kind() == reflect.String && tgt.Kind() == reflect.Int {
			if input, ok := inputRaw.(string); ok {
				if intval, err := strconv.Atoi(input); err == nil {
					if intval >= 100 && intval < 600 {
						return intval, nil
					}
				}
			}
		}
		return inputRaw, nil
	}

	doc := new(OASDocument)

	config := &mapstructure.DecoderConfig{
		DecodeHook: decodeHook,
		Result:     doc,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(input); err != nil {
		return nil, err
	}

	return doc, nil
}

type OASDocument struct {
	Version string                  `json:"openapi" mapstructure:"openapi"`
	Info    OASInfo                 `json:"info"`
	Paths   map[string]*OASPathItem `json:"paths"`
}

type OASInfo struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	License     OASLicense `json:"license"`
}

type OASLicense struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type OASPathItem struct {
	Description       string             `json:"description,omitempty"`
	Parameters        []OASParameter     `json:"parameters,omitempty"`
	Sudo              bool               `json:"x-vault-sudo,omitempty" mapstructure:"x-vault-sudo"`
	Unauthenticated   bool               `json:"x-vault-unauthenticated,omitempty" mapstructure:"x-vault-unauthenticated"`
	CreateSupported   bool               `json:"x-vault-createSupported,omitempty" mapstructure:"x-vault-createSupported"`
	DisplayNavigation bool               `json:"x-vault-displayNavigation,omitempty" mapstructure:"x-vault-displayNavigation"`
	DisplayAttrs      *DisplayAttributes `json:"x-vault-displayAttrs,omitempty" mapstructure:"x-vault-displayAttrs"`

	Get    *OASOperation `json:"get,omitempty"`
	Post   *OASOperation `json:"post,omitempty"`
	Delete *OASOperation `json:"delete,omitempty"`
}

// NewOASOperation creates an empty OpenAPI Operations object.
func NewOASOperation() *OASOperation {
	return &OASOperation{
		Responses: make(map[int]*OASResponse),
	}
}

type OASOperation struct {
	Summary     string               `json:"summary,omitempty"`
	Description string               `json:"description,omitempty"`
	OperationID string               `json:"operationId,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Parameters  []OASParameter       `json:"parameters,omitempty"`
	RequestBody *OASRequestBody      `json:"requestBody,omitempty"`
	Responses   map[int]*OASResponse `json:"responses"`
	Deprecated  bool                 `json:"deprecated,omitempty"`
}

type OASParameter struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	In          string     `json:"in"`
	Schema      *OASSchema `json:"schema,omitempty"`
	Required    bool       `json:"required,omitempty"`
	Deprecated  bool       `json:"deprecated,omitempty"`
}

type OASRequestBody struct {
	Description string     `json:"description,omitempty"`
	Content     OASContent `json:"content,omitempty"`
}

type OASContent map[string]*OASMediaTypeObject

type OASMediaTypeObject struct {
	Schema *OASSchema `json:"schema,omitempty"`
}

type OASSchema struct {
	Type        string                `json:"type,omitempty"`
	Description string                `json:"description,omitempty"`
	Properties  map[string]*OASSchema `json:"properties,omitempty"`

	// Required is a list of keys in Properties that are required to be present. This is a different
	// approach than OASParameter (unfortunately), but is how JSONSchema handles 'required'.
	Required []string `json:"required,omitempty"`

	Items      *OASSchema    `json:"items,omitempty"`
	Format     string        `json:"format,omitempty"`
	Pattern    string        `json:"pattern,omitempty"`
	Enum       []interface{} `json:"enum,omitempty"`
	Default    interface{}   `json:"default,omitempty"`
	Example    interface{}   `json:"example,omitempty"`
	Deprecated bool          `json:"deprecated,omitempty"`
	//DisplayName      string             `json:"x-vault-displayName,omitempty" mapstructure:"x-vault-displayName,omitempty"`
	DisplayValue     interface{}        `json:"x-vault-displayValue,omitempty" mapstructure:"x-vault-displayValue,omitempty"`
	DisplaySensitive bool               `json:"x-vault-displaySensitive,omitempty" mapstructure:"x-vault-displaySensitive,omitempty"`
	DisplayGroup     string             `json:"x-vault-displayGroup,omitempty" mapstructure:"x-vault-displayGroup,omitempty"`
	DisplayAttrs     *DisplayAttributes `json:"x-vault-displayAttrs,omitempty" mapstructure:"x-vault-displayAttrs,omitempty"`
}

type OASResponse struct {
	Description string     `json:"description"`
	Content     OASContent `json:"content,omitempty"`
}

var OASStdRespOK = &OASResponse{
	Description: "OK",
}

var OASStdRespNoContent = &OASResponse{
	Description: "empty body",
}

// Regex for handling optional and named parameters in paths, and string cleanup.
// Predefined here to avoid substantial recompilation.

// Capture optional path elements in ungreedy (?U) fashion
// Both "(leases/)?renew" and "(/(?P<name>.+))?" formats are detected
var optRe = regexp.MustCompile(`(?U)\([^(]*\)\?|\(/\(\?P<[^(]*\)\)\?`)

var reqdRe = regexp.MustCompile(`\(?\?P<(\w+)>[^)]*\)?`)             // Capture required parameters, e.g. "(?P<name>regex)"
var altRe = regexp.MustCompile(`\((.*)\|(.*)\)`)                     // Capture alternation elements, e.g. "(raw/?$|raw/(?P<path>.+))"
var pathFieldsRe = regexp.MustCompile(`{(\w+)}`)                     // Capture OpenAPI-style named parameters, e.g. "lookup/{urltoken}",
var cleanCharsRe = regexp.MustCompile("[()^$?]")                     // Set of regex characters that will be stripped during cleaning
var cleanSuffixRe = regexp.MustCompile(`/\?\$?$`)                    // Path suffix patterns that will be stripped during cleaning
var wsRe = regexp.MustCompile(`\s+`)                                 // Match whitespace, to be compressed during cleaning
var altFieldsGroupRe = regexp.MustCompile(`\(\?P<\w+>\w+(\|\w+)+\)`) // Match named groups that limit options, e.g. "(?<foo>a|b|c)"
var altFieldsRe = regexp.MustCompile(`\w+(\|\w+)+`)                  // Match an options set, e.g. "a|b|c"
var nonWordRe = regexp.MustCompile(`[^\w]+`)                         // Match a sequence of non-word characters

// documentPaths parses all paths in a framework.Backend into OpenAPI paths.
func documentPaths(backend *Backend, doc *OASDocument) error {
	for _, p := range backend.Paths {
		if err := documentPath(p, backend.SpecialPaths(), backend.BackendType, doc); err != nil {
			return err
		}
	}

	return nil
}

// documentPath parses a framework.Path into one or more OpenAPI paths.
func documentPath(p *Path, specialPaths *logical.Paths, backendType logical.BackendType, doc *OASDocument) error {
	var sudoPaths []string
	var unauthPaths []string

	if specialPaths != nil {
		sudoPaths = specialPaths.Root
		unauthPaths = specialPaths.Unauthenticated
	}

	// Convert optional parameters into distinct patterns to be process independently.
	paths := expandPattern(p.Pattern)

	for _, path := range paths {
		// Construct a top level PathItem which will be populated as the path is processed.
		pi := OASPathItem{
			Description: cleanString(p.HelpSynopsis),
		}

		pi.Sudo = specialPathMatch(path, sudoPaths)
		pi.Unauthenticated = specialPathMatch(path, unauthPaths)
		pi.DisplayAttrs = p.DisplayAttrs

		// If the newer style Operations map isn't defined, create one from the legacy fields.
		operations := p.Operations
		if operations == nil {
			operations = make(map[logical.Operation]OperationHandler)

			for opType, cb := range p.Callbacks {
				operations[opType] = &PathOperation{
					Callback: cb,
					Summary:  p.HelpSynopsis,
				}
			}
		}

		// Process path and header parameters, which are common to all operations.
		// Body fields will be added to individual operations.
		pathFields, bodyFields := splitFields(p.Fields, path)

		for name, field := range pathFields {
			location := "path"
			required := true

			if field.Query {
				location = "query"
				required = false
			}

			t := convertType(field.Type)
			p := OASParameter{
				Name:        name,
				Description: cleanString(field.Description),
				In:          location,
				Schema: &OASSchema{
					Type:         t.baseType,
					Pattern:      t.pattern,
					Enum:         field.AllowedValues,
					Default:      field.Default,
					DisplayAttrs: field.DisplayAttrs,
				},
				Required:   required,
				Deprecated: field.Deprecated,
			}
			pi.Parameters = append(pi.Parameters, p)
		}

		// Sort parameters for a stable output
		sort.Slice(pi.Parameters, func(i, j int) bool {
			return strings.ToLower(pi.Parameters[i].Name) < strings.ToLower(pi.Parameters[j].Name)
		})

		// Process each supported operation by building up an Operation object
		// with descriptions, properties and examples from the framework.Path data.
		for opType, opHandler := range operations {
			props := opHandler.Properties()
			if props.Unpublished {
				continue
			}

			if opType == logical.CreateOperation {
				pi.CreateSupported = true

				// If both Create and Update are defined, only process Update.
				if operations[logical.UpdateOperation] != nil {
					continue
				}
			}

			// If both List and Read are defined, only process Read.
			if opType == logical.ListOperation && operations[logical.ReadOperation] != nil {
				continue
			}

			op := NewOASOperation()

			op.Summary = props.Summary
			op.Description = props.Description
			op.Deprecated = props.Deprecated

			// Add any fields not present in the path as body parameters for POST.
			if opType == logical.CreateOperation || opType == logical.UpdateOperation {
				s := &OASSchema{
					Type:       "object",
					Properties: make(map[string]*OASSchema),
					Required:   make([]string, 0),
				}

				for name, field := range bodyFields {
					openapiField := convertType(field.Type)
					if field.Required {
						s.Required = append(s.Required, name)
					}

					p := OASSchema{
						Type:         openapiField.baseType,
						Description:  cleanString(field.Description),
						Format:       openapiField.format,
						Pattern:      openapiField.pattern,
						Enum:         field.AllowedValues,
						Default:      field.Default,
						Deprecated:   field.Deprecated,
						DisplayAttrs: field.DisplayAttrs,
					}
					if openapiField.baseType == "array" {
						p.Items = &OASSchema{
							Type: openapiField.items,
						}
					}
					s.Properties[name] = &p
				}

				// If examples were given, use the first one as the sample
				// of this schema.
				if len(props.Examples) > 0 {
					s.Example = props.Examples[0].Data
				}

				// Set the final request body. Only JSON request data is supported.
				if len(s.Properties) > 0 || s.Example != nil {
					op.RequestBody = &OASRequestBody{
						Content: OASContent{
							"application/json": &OASMediaTypeObject{
								Schema: s,
							},
						},
					}
				}
			}

			// LIST is represented as GET with a `list` query parameter
			if opType == logical.ListOperation || (opType == logical.ReadOperation && operations[logical.ListOperation] != nil) {
				op.Parameters = append(op.Parameters, OASParameter{
					Name:        "list",
					Description: "Return a list if `true`",
					In:          "query",
					Schema:      &OASSchema{Type: "string"},
				})
			}

			// Add tags based on backend type
			var tags []string
			switch backendType {
			case logical.TypeLogical:
				tags = []string{"secrets"}
			case logical.TypeCredential:
				tags = []string{"auth"}
			}

			op.Tags = append(op.Tags, tags...)

			// Set default responses.
			if len(props.Responses) == 0 {
				if opType == logical.DeleteOperation {
					op.Responses[204] = OASStdRespNoContent
				} else {
					op.Responses[200] = OASStdRespOK
				}
			}

			// Add any defined response details.
			for code, responses := range props.Responses {
				var description string
				content := make(OASContent)

				for i, resp := range responses {
					if i == 0 {
						description = resp.Description
					}
					if resp.Example != nil {
						mediaType := resp.MediaType
						if mediaType == "" {
							mediaType = "application/json"
						}

						// create a version of the response that will not emit null items
						cr := cleanResponse(resp.Example)

						// Only one example per media type is allowed, so first one wins
						if _, ok := content[mediaType]; !ok {
							content[mediaType] = &OASMediaTypeObject{
								Schema: &OASSchema{
									Example: cr,
								},
							}
						}
					}
				}

				op.Responses[code] = &OASResponse{
					Description: description,
					Content:     content,
				}
			}

			switch opType {
			case logical.CreateOperation, logical.UpdateOperation:
				pi.Post = op
			case logical.ReadOperation, logical.ListOperation:
				pi.Get = op
			case logical.DeleteOperation:
				pi.Delete = op
			}
		}

		doc.Paths["/"+path] = &pi
	}

	return nil
}

func specialPathMatch(path string, specialPaths []string) bool {
	// Test for exact or prefix match of special paths.
	for _, sp := range specialPaths {
		if sp == path ||
			(strings.HasSuffix(sp, "*") && strings.HasPrefix(path, sp[0:len(sp)-1])) {
			return true
		}
	}
	return false
}

// expandPattern expands a regex pattern by generating permutations of any optional parameters
// and changing named parameters into their {openapi} equivalents.
func expandPattern(pattern string) []string {
	var paths []string

	// GenericNameRegex adds a regex that complicates our parsing. It is much easier to
	// detect and remove it now than to compensate for in the other regexes.
	//
	// example: (?P<foo>\\w(([\\w-.]+)?\\w)?) -> (?P<foo>)
	base := GenericNameRegex("")
	start := strings.Index(base, ">")
	end := strings.LastIndex(base, ")")
	regexToRemove := ""
	if start != -1 && end != -1 && end > start {
		regexToRemove = base[start+1 : end]
	}
	pattern = strings.Replace(pattern, regexToRemove, "", -1)

	// Simplify named fields that have limited options, e.g. (?P<foo>a|b|c) -> (<P<foo>.+)
	pattern = altFieldsGroupRe.ReplaceAllStringFunc(pattern, func(s string) string {
		return altFieldsRe.ReplaceAllString(s, ".+")
	})

	// Initialize paths with the original pattern or the halves of an
	// alternation, which is also present in some patterns.
	matches := altRe.FindAllStringSubmatch(pattern, -1)
	if len(matches) > 0 {
		paths = []string{matches[0][1], matches[0][2]}
	} else {
		paths = []string{pattern}
	}

	// Expand all optional regex elements into two paths. This approach is really only useful up to 2 optional
	// groups, but we probably don't want to deal with the exponential increase beyond that anyway.
	for i := 0; i < len(paths); i++ {
		p := paths[i]

		// match is a 2-element slice that will have a start and end index
		// for the left-most match of a regex of form: (lease/)?
		match := optRe.FindStringIndex(p)

		if match != nil {
			// create a path that includes the optional element but without
			// parenthesis or the '?' character.
			paths[i] = p[:match[0]] + p[match[0]+1:match[1]-2] + p[match[1]:]

			// create a path that excludes the optional element.
			paths = append(paths, p[:match[0]]+p[match[1]:])
			i--
		}
	}

	// Replace named parameters (?P<foo>) with {foo}
	var replacedPaths []string

	for _, path := range paths {
		result := reqdRe.FindAllStringSubmatch(path, -1)
		if result != nil {
			for _, p := range result {
				par := p[1]
				path = strings.Replace(path, p[0], fmt.Sprintf("{%s}", par), 1)
			}
		}
		// Final cleanup
		path = cleanSuffixRe.ReplaceAllString(path, "")
		path = cleanCharsRe.ReplaceAllString(path, "")
		replacedPaths = append(replacedPaths, path)
	}

	return replacedPaths
}

// schemaType is a subset of the JSON Schema elements used as a target
// for conversions from Vault's standard FieldTypes.
type schemaType struct {
	baseType string
	items    string
	format   string
	pattern  string
}

// convertType translates a FieldType into an OpenAPI type.
// In the case of arrays, a subtype is returned as well.
func convertType(t FieldType) schemaType {
	ret := schemaType{}

	switch t {
	case TypeString, TypeHeader:
		ret.baseType = "string"
	case TypeNameString:
		ret.baseType = "string"
		ret.pattern = `\w([\w-.]*\w)?`
	case TypeLowerCaseString:
		ret.baseType = "string"
		ret.format = "lowercase"
	case TypeInt:
		ret.baseType = "integer"
	case TypeDurationSecond, TypeSignedDurationSecond:
		ret.baseType = "integer"
		ret.format = "seconds"
	case TypeBool:
		ret.baseType = "boolean"
	case TypeMap:
		ret.baseType = "object"
		ret.format = "map"
	case TypeKVPairs:
		ret.baseType = "object"
		ret.format = "kvpairs"
	case TypeSlice:
		ret.baseType = "array"
		ret.items = "object"
	case TypeStringSlice, TypeCommaStringSlice:
		ret.baseType = "array"
		ret.items = "string"
	case TypeCommaIntSlice:
		ret.baseType = "array"
		ret.items = "integer"
	default:
		log.L().Warn("error parsing field type", "type", t)
		ret.format = "unknown"
	}

	return ret
}

// cleanString prepares s for inclusion in the output
func cleanString(s string) string {
	// clean leading/trailing whitespace, and replace whitespace runs into a single space
	s = strings.TrimSpace(s)
	s = wsRe.ReplaceAllString(s, " ")
	return s
}

// splitFields partitions fields into path and body groups
// The input pattern is expected to have been run through expandPattern,
// with paths parameters denotes in {braces}.
func splitFields(allFields map[string]*FieldSchema, pattern string) (pathFields, bodyFields map[string]*FieldSchema) {
	pathFields = make(map[string]*FieldSchema)
	bodyFields = make(map[string]*FieldSchema)

	for _, match := range pathFieldsRe.FindAllStringSubmatch(pattern, -1) {
		name := match[1]
		pathFields[name] = allFields[name]
	}

	for name, field := range allFields {
		if _, ok := pathFields[name]; !ok {
			if field.Query {
				pathFields[name] = field
			} else {
				bodyFields[name] = field
			}
		}
	}

	return pathFields, bodyFields
}

// cleanedResponse is identical to logical.Response but with nulls
// removed from from JSON encoding
type cleanedResponse struct {
	Secret   *logical.Secret            `json:"secret,omitempty"`
	Auth     *logical.Auth              `json:"auth,omitempty"`
	Data     map[string]interface{}     `json:"data,omitempty"`
	Redirect string                     `json:"redirect,omitempty"`
	Warnings []string                   `json:"warnings,omitempty"`
	WrapInfo *wrapping.ResponseWrapInfo `json:"wrap_info,omitempty"`
	Headers  map[string][]string        `json:"headers,omitempty"`
}

func cleanResponse(resp *logical.Response) *cleanedResponse {
	return &cleanedResponse{
		Secret:   resp.Secret,
		Auth:     resp.Auth,
		Data:     resp.Data,
		Redirect: resp.Redirect,
		Warnings: resp.Warnings,
		WrapInfo: resp.WrapInfo,
		Headers:  resp.Headers,
	}
}

// CreateOperationIDs generates unique operationIds for all paths/methods.
// The transform will convert path/method into camelcase. e.g.:
//
// /sys/tools/random/{urlbytes} -> postSysToolsRandomUrlbytes
//
// In the unlikely case of a duplicate ids, a numeric suffix is added:
//   postSysToolsRandomUrlbytes_2
//
// An optional user-provided suffix ("context") may also be appended.
func (d *OASDocument) CreateOperationIDs(context string) {
	opIDCount := make(map[string]int)
	var paths []string

	// traverse paths in a stable order to ensure stable output
	for path := range d.Paths {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		pi := d.Paths[path]
		for _, method := range []string{"get", "post", "delete"} {
			var oasOperation *OASOperation
			switch method {
			case "get":
				oasOperation = pi.Get
			case "post":
				oasOperation = pi.Post
			case "delete":
				oasOperation = pi.Delete
			}

			if oasOperation == nil {
				continue
			}

			// Space-split on non-words, title case everything, recombine
			opID := nonWordRe.ReplaceAllString(strings.ToLower(path), " ")
			opID = strings.Title(opID)
			opID = method + strings.Replace(opID, " ", "", -1)

			// deduplicate operationIds. This is a safeguard, since generated IDs should
			// already be unique given our current path naming conventions.
			opIDCount[opID]++
			if opIDCount[opID] > 1 {
				opID = fmt.Sprintf("%s_%d", opID, opIDCount[opID])
			}

			if context != "" {
				opID += "_" + context
			}

			oasOperation.OperationID = opID
		}
	}
}
