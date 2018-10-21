package framework

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/openapi"
	"github.com/hashicorp/vault/helper/wrapping"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/mapstructure"
)

// Regex for handling optional and named parameters in paths, and string cleanup.
// Predefined here to avoid substantial recompilation.
var reqdRe = regexp.MustCompile(`\(?\?P<(\w+)>[^)]*\)?`) // capture required named parameters
var optRe = regexp.MustCompile(`(?U)\(.*\)\?`)           // capture optional named parameters in ungreedy (?U) fashion
var altRe = regexp.MustCompile(`\((.*)\|(.*)\)`)         // capture alternation elements
var pathFieldsRe = regexp.MustCompile(`{(\w+)}`)         // capture OpenAPI-format named parameters {example}
var cleanCharsRe = regexp.MustCompile("[()^$?]")         // regex characters that will be stripped during cleaning
var cleanSuffixRe = regexp.MustCompile(`/\?\$?$`)        // path suffix patterns that will be stripped during cleaning
var wsRe = regexp.MustCompile(`\s+`)                     // match whitespace, to be compressed during cleaning

// documentPaths parses all paths in a framework.Backend into OpenAPI paths.
func documentPaths(backend *Backend, doc *openapi.Document) error {
	var sudoPaths []string

	sp := backend.SpecialPaths()
	if sp != nil {
		sudoPaths = sp.Root
	}

	for _, p := range backend.Paths {
		if err := documentPath(p, sudoPaths, doc); err != nil {
			return err
		}
	}

	return nil
}

// documentPath parses a framework.Path into one or more OpenAPI paths.
func documentPath(p *Path, sudoPaths []string, doc *openapi.Document) error {

	// Convert optional parameters into distinct patterns to be process independently.
	paths := expandPattern(p.Pattern)

	for _, path := range paths {
		// Construct a top level PathItem which will be populated as the path is processed.
		pi := openapi.PathItem{
			Description: cleanString(p.HelpSynopsis),
		}

		// Test for exact or prefix match of root paths.
		for _, root := range sudoPaths {
			if root == path ||
				(strings.HasSuffix(root, "*") && strings.HasPrefix(path, root[0:len(root)-1])) {
				pi.Sudo = true
				break
			}
		}

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

		pathFields, bodyFields := splitFields(p.Fields, path)

		// Process each supported operation by building up an Operation object
		// with descriptions, properties and examples from the framework.Path data.
		for opType, opHandler := range operations {
			props := opHandler.Properties()
			if props.Unpublished {
				continue
			}

			if opType == logical.CreateOperation {
				pi.Create = true

				// If both Create and Update are defined, only process Update.
				if operations[logical.UpdateOperation] != nil {
					continue
				}
			}

			// If both List and Read are defined, only process Read.
			if opType == logical.ListOperation && operations[logical.ReadOperation] != nil {
				continue
			}

			op := openapi.NewOperation()

			op.Summary = props.Summary
			op.Description = props.Description
			op.Deprecated = props.Deprecated

			for name, field := range pathFields {
				location := "path"
				required := true

				// Header parameters are part of the Parameters group but with
				// a dedicated "header" location and are not required.
				if field.Type == TypeHeader {
					location = "header"
					required = false
				}

				t := convertType(field.Type)
				p := openapi.Parameter{
					Name:        name,
					Description: cleanString(field.Description),
					In:          location,
					Schema:      &openapi.Schema{Type: t.baseType},
					Required:    required,
					Deprecated:  field.Deprecated,
				}
				op.Parameters = append(op.Parameters, p)
			}

			// LIST is represented as GET with a `list` query parameter
			if opType == logical.ListOperation || (opType == logical.ReadOperation && operations[logical.ListOperation] != nil) {
				op.Parameters = append(op.Parameters, openapi.Parameter{
					Name:        "list",
					Description: "Return a list if `true`",
					In:          "query",
					Schema:      &openapi.Schema{Type: "string"},
					Required:    true,
				})
			}

			// Add any fields not present in the path as body parameters for POST.
			if opType == logical.CreateOperation || opType == logical.UpdateOperation {
				s := &openapi.Schema{
					Type:       "object",
					Properties: make(map[string]*openapi.Schema),
				}

				for name, field := range bodyFields {
					openapiField := convertType(field.Type)
					p := openapi.Schema{
						Type:        openapiField.baseType,
						Description: cleanString(field.Description),
						Format:      openapiField.format,
						Deprecated:  field.Deprecated,
					}
					if openapiField.baseType == "array" {
						p.Items = &openapi.Schema{
							Type: openapiField.items,
						}
					}
					s.Properties[name] = &p
				}

				// If examples were given, use the first one as the sample
				// of this schema.
				if len(props.Examples) > 0 {
					s.Example = props.Examples[0].Value
				}

				// Set the final request body. Only JSON request data is supported.
				if len(s.Properties) > 0 || s.Example != nil {
					op.RequestBody = &openapi.RequestBody{
						Content: &openapi.Content{
							"application/json": &openapi.MediaTypeObject{
								Schema: s,
							},
						},
					}
				}
			}

			// Set default responses.
			if len(props.Responses) == 0 {
				if opType == logical.DeleteOperation {
					op.Responses[204] = openapi.StdRespNoContent
				} else {
					op.Responses[200] = openapi.StdRespOK
				}
			}

			// Add any defined response details.
			for code, responses := range props.Responses {
				var description string
				content := make(openapi.Content)

				for i, resp := range responses {
					if i == 0 {
						description = resp.Description
					}
					if resp.Example != nil {
						mediaType := resp.MediaType
						if mediaType == "" {
							mediaType = "application/json"
						}

						// Only one example per media type is allowed, so first one wins
						if _, ok := content[mediaType]; !ok {
							content[mediaType] = &openapi.MediaTypeObject{
								Schema: &openapi.Schema{
									Example: cleanResponse(resp.Example),
								},
							}
						}
					}
				}

				// a nil pointer when empty is needed for omitempty to work
				var c *openapi.Content
				if len(content) > 0 {
					c = &content
				}
				op.Responses[code] = &openapi.Response{
					Description: description,
					Content:     c,
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

// expandPattern expands a regex pattern by generating permutations of any optional parameters
// and changing named parameters into their {openapi} equivalents.
func expandPattern(pattern string) []string {
	var paths []string

	// This construct is added by GenericNameRegex and is much easier to remove now
	// than to compensate for in the other regexes.
	pattern = strings.Replace(pattern, `\w(([\w-.]+)?\w)?`, "", -1)

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
		match := optRe.FindStringIndex(p)
		if match != nil {
			paths[i] = p[0:match[0]] + p[match[0]+1:match[1]-2] + p[match[1]:]
			paths = append(paths, p[0:match[0]]+p[match[1]:])
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
}

// convertType translates a FieldType into an OpenAPI type.
// In the case of arrays, a subtype is returned as well.
func convertType(t FieldType) schemaType {
	ret := schemaType{}

	switch t {
	case TypeString, TypeNameString, TypeLowerCaseString, TypeHeader:
		ret.baseType = "string"
	case TypeInt:
		ret.baseType = "number"
	case TypeDurationSecond:
		ret.baseType = "number"
		ret.format = "seconds"
	case TypeBool:
		ret.baseType = "boolean"
	case TypeMap:
		ret.baseType = "object"
		ret.format = "map"
	case TypeKVPairs:
		ret.baseType = "object"
		ret.format = "kvpairs"
	case TypeSlice, TypeStringSlice, TypeCommaStringSlice:
		ret.baseType = "array"
		ret.items = "string"
	case TypeCommaIntSlice:
		ret.baseType = "array"
		ret.items = "number"
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
			// Header fields are in "parameters" with other path fields
			if field.Type == TypeHeader {
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
}

func cleanResponse(resp *logical.Response) *cleanedResponse {
	var r cleanedResponse

	if mapstructure.Decode(resp, &r) != nil {
		return nil
	}

	return &r
}
